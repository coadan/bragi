package bragi

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestToolCanRepairBeforeCommitButNotAfter(t *testing.T) {
	profile := loadMidgardProfile(t)
	materializer := mustMaterializer(t, profile)
	applySource(t, materializer, "+ @t1 tool\n+ @t1.name \"search\"\n~ @t1.name \"repo.inspect\"\n+ @t1.reason \"Inspect first\"\n! @t1\n")
	entity, ok := materializer.Entity("@t1")
	if !ok || entity.Status != StatusCommitted || entity.Fields["name"].Scalar.String != "repo.inspect" {
		t.Fatalf("committed tool = %#v", entity)
	}
	rejected := materializer.Apply(Record{Operation: OpReplace, Target: "@t1.name", Value: &Value{Kind: ValueString, String: "shell.run"}, Line: 6})
	if len(rejected) != 1 || rejected[0].Diagnostic == nil || rejected[0].Diagnostic.Code != "immutable_entity" {
		t.Fatalf("post-commit mutation = %#v", rejected)
	}
}

func TestRevisionPinsCommittedReferences(t *testing.T) {
	profile := loadMidgardProfile(t)
	materializer := mustMaterializer(t, profile)
	applySource(t, materializer, "+ @s1 plan_step\n+ @s1.intent \"First\"\n+ @s1.acceptance \"Passes\"\n+ @s1.status \"planned\"\n! @s1\n+ @c1 commit_unit\n+ @c1.intent \"Implement\"\n+ @c1.plan_steps @s1\n! @c1\n~ @s1.status \"done\"\n! @s1\n")
	unit, _ := materializer.Entity("@c1")
	refs := unit.Committed[0].References["plan_steps"]
	if len(refs) != 1 || refs[0].Revision != 1 {
		t.Fatalf("pinned refs = %#v", refs)
	}
	step, _ := materializer.Entity("@s1")
	if step.Revision != 2 || len(step.Committed) != 2 {
		t.Fatalf("revised step = %#v", step)
	}
}

func TestRejectedRecordDoesNotMutateView(t *testing.T) {
	profile := loadMidgardProfile(t)
	materializer := mustMaterializer(t, profile)
	applySource(t, materializer, "+ @s1 plan_step\n+ @s1.intent \"First\"\n")
	before, _ := materializer.Entity("@s1")
	events := materializer.Apply(Record{Operation: OpAdd, Target: "@s1.unknown", Value: &Value{Kind: ValueString, String: "bad"}, Line: 3})
	after, _ := materializer.Entity("@s1")
	if events[0].Kind != "op.rejected" || !reflect.DeepEqual(before, after) {
		t.Fatalf("rejected event/view = %#v / %#v", events, after)
	}
}

func TestMessagePublicationGate(t *testing.T) {
	materializer := mustMaterializer(t, loadMidgardProfile(t))
	applySource(t, materializer, "+ @m1 message\n+ @m1.content \"hello\"\n")
	if materializer.Publishable("@m1") {
		t.Fatal("message became publishable before routing fields")
	}
	applySource(t, materializer, "+ @m1.speaker \"assistant\"\n+ @m1.audience \"user\"\n+ @m1.channel \"commentary\"\n")
	if !materializer.Publishable("@m1") {
		t.Fatal("message did not become publishable after routing fields")
	}
}

func TestRecoveredCaseStillEnforcesCanonicalCollisionsAndSchema(t *testing.T) {
	materializer := mustMaterializer(t, loadMidgardProfile(t))
	decoder := NewDecoder(DecoderOptions{})
	records, issues := decoder.Write([]byte("+ @T1 Tool\n+ @t1 tool\n+ @T1.Unknown \"value\"\n"))
	if len(issues) != 0 {
		t.Fatalf("source diagnostics: %#v", issues)
	}
	if events := materializer.Apply(records[0]); events[0].Kind != "op.accepted" {
		t.Fatalf("first create = %#v", events)
	}
	if events := materializer.Apply(records[1]); events[0].Diagnostic == nil || events[0].Diagnostic.Code != "entity_exists" {
		t.Fatalf("case collision = %#v", events)
	}
	if events := materializer.Apply(records[2]); events[0].Diagnostic == nil || events[0].Diagnostic.Code != "unknown_field" {
		t.Fatalf("unknown field = %#v", events)
	}
}

func TestExamplesConformAndReplayDeterministically(t *testing.T) {
	profile := loadMidgardProfile(t)
	paths, err := filepath.Glob("examples/*.bragi")
	if err != nil || len(paths) == 0 {
		t.Fatalf("examples: %v, %v", paths, err)
	}
	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			first := mustMaterializer(t, profile)
			applySource(t, first, string(source))
			replayed, err := Replay(profile, first.Events())
			if err != nil {
				t.Fatal(err)
			}
			if issues := first.ValidateComplete(); len(issues) > 0 {
				t.Fatalf("incomplete: %#v", issues)
			}
			if !reflect.DeepEqual(first.Events(), replayed.Events()) {
				t.Fatal("canonical replay is not deterministic")
			}
		})
	}
}

func TestReplayRejectsSequenceGap(t *testing.T) {
	profile := loadMidgardProfile(t)
	_, err := Replay(profile, []Event{{Sequence: 2, Kind: "source.rejected", Diagnostic: &Diagnostic{Code: "bad"}}})
	if err == nil {
		t.Fatal("replay accepted a sequence gap")
	}
}

func loadMidgardProfile(t *testing.T) Profile {
	t.Helper()
	file, err := os.Open("profiles/midgard-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	profile, err := LoadProfile(file)
	if err != nil {
		t.Fatal(err)
	}
	return profile
}

func mustMaterializer(t *testing.T, profile Profile) *Materializer {
	t.Helper()
	materializer, err := NewMaterializer(profile)
	if err != nil {
		t.Fatal(err)
	}
	return materializer
}

func applySource(t *testing.T, materializer *Materializer, source string) {
	t.Helper()
	decoder := NewDecoder(DecoderOptions{MaxLineBytes: materializer.profile.Limits.MaxLineBytes, AllowCRLF: true})
	records, issues := decoder.Write([]byte(source))
	issues = append(issues, decoder.Finish()...)
	if len(issues) > 0 {
		t.Fatalf("source diagnostics: %#v", issues)
	}
	for _, record := range records {
		for _, event := range materializer.Apply(record) {
			if event.Kind == "op.rejected" || event.Kind == "commit.rejected" {
				t.Fatalf("rejected source event: %#v", event)
			}
		}
	}
}
