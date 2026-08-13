package bragi

import (
	"reflect"
	"strings"
	"testing"
)

func TestDecoderIsChunkIndependent(t *testing.T) {
	source := "+ @t1 tool\n+ @t1.name \"search\"\n~ @t1.name \"repo.inspect\"\n! @t1\n"
	wantRecords, wantIssues := decodeChunks(t, source, len(source))
	if len(wantIssues) != 0 {
		t.Fatalf("whole stream diagnostics: %#v", wantIssues)
	}
	gotRecords, gotIssues := decodeChunks(t, source, 1)
	if !reflect.DeepEqual(gotRecords, wantRecords) || !reflect.DeepEqual(gotIssues, wantIssues) {
		t.Fatalf("one-byte decoding differs\nrecords: %#v\nissues: %#v", gotRecords, gotIssues)
	}
}

func TestDecoderPreservesAcceptedPrefixOnTruncation(t *testing.T) {
	records, issues := decodeChunks(t, "+ @t1 tool\n+ @t1.name \"sea", 3)
	if len(records) != 1 || records[0].Operation != OpCreate {
		t.Fatalf("accepted prefix = %#v", records)
	}
	if len(issues) != 1 || issues[0].Code != "incomplete_record" {
		t.Fatalf("diagnostics = %#v", issues)
	}
}

func TestDecoderLiteralRequiresExactSeal(t *testing.T) {
	records, issues := decodeChunks(t, "+ @a1 artifact\n+ @a1.content |\n|first\n||second\n! @a1.content\n! @a1\n", 5)
	if len(issues) != 0 {
		t.Fatalf("diagnostics = %#v", issues)
	}
	operations := []Operation{OpCreate, OpLiteralOpen, OpLiteralAppend, OpLiteralAppend, OpLiteralSeal, OpCommit}
	for index, operation := range operations {
		if records[index].Operation != operation {
			t.Fatalf("record %d operation = %s, want %s", index, records[index].Operation, operation)
		}
	}
}

func TestDecoderRecoversStructuralCaseWithoutChangingPayload(t *testing.T) {
	source := "+ @T1 Tool\r\n+ @T1.Name \"Shell.Run\"\n+ @T1.Arguments.Parent @P1\n! @t1\n"
	records, issues := decodeChunks(t, source, 1)
	if len(issues) != 0 {
		t.Fatalf("diagnostics = %#v", issues)
	}
	if records[0].Target != "@t1" || records[0].EntityType != "tool" {
		t.Fatalf("create was not canonicalized: %#v", records[0])
	}
	if records[1].Target != "@t1.name" || records[1].Value.String != "Shell.Run" {
		t.Fatalf("string payload changed or target was not canonicalized: %#v", records[1])
	}
	if records[2].Target != "@t1.arguments.parent" || records[2].Value.String != "@p1" {
		t.Fatalf("reference was not canonicalized: %#v", records[2])
	}
	if len(records[0].Normalizations) != 3 {
		t.Fatalf("case/CRLF recovery was not recorded: %#v", records[0].Normalizations)
	}
}

func TestDecoderStrictModeRejectsNonCanonicalCase(t *testing.T) {
	decoder := NewDecoder(DecoderOptions{StrictSource: true})
	_, issues := decoder.Write([]byte("+ @T1 Tool\n"))
	if len(issues) != 1 || issues[0].Code != "invalid_path" {
		t.Fatalf("strict diagnostics = %#v", issues)
	}
}

func TestDecoderRecoversCaseInLiteralSeal(t *testing.T) {
	records, issues := decodeChunks(t, "+ @A1 Artifact\n+ @A1.Content |\n|Mixed Case stays\n! @a1.CONTENT\n", 2)
	if len(issues) != 0 || records[len(records)-1].Operation != OpLiteralSeal {
		t.Fatalf("literal seal recovery = %#v / %#v", records, issues)
	}
	if records[2].Value.String != "Mixed Case stays" {
		t.Fatalf("literal payload changed: %#v", records[2])
	}
}

func TestFinishCompletedRecoversOnlyTerminalLF(t *testing.T) {
	decoder := NewDecoder(DecoderOptions{})
	records, issues := decoder.Write([]byte("+ @T1 Tool"))
	if len(records) != 0 || len(issues) != 0 {
		t.Fatalf("partial write = %#v / %#v", records, issues)
	}
	records, issues = decoder.FinishCompleted()
	if len(issues) != 0 || len(records) != 1 || records[0].Target != "@t1" {
		t.Fatalf("completed finish = %#v / %#v", records, issues)
	}
	last := records[0].Normalizations[len(records[0].Normalizations)-1]
	if last.Kind != "terminal-lf" {
		t.Fatalf("terminal recovery not recorded: %#v", records[0].Normalizations)
	}

	abrupt := NewDecoder(DecoderOptions{})
	abrupt.Write([]byte("+ @t1 tool"))
	issues = abrupt.Finish()
	if len(issues) != 1 || issues[0].Code != "incomplete_record" {
		t.Fatalf("abrupt finish was recovered: %#v", issues)
	}
}

func decodeChunks(t *testing.T, source string, chunkSize int) ([]Record, []Diagnostic) {
	t.Helper()
	decoder := NewDecoder(DecoderOptions{AllowCRLF: true})
	var records []Record
	var issues []Diagnostic
	reader := strings.NewReader(source)
	buffer := make([]byte, chunkSize)
	for {
		count, err := reader.Read(buffer)
		if count > 0 {
			newRecords, newIssues := decoder.Write(buffer[:count])
			records = append(records, newRecords...)
			issues = append(issues, newIssues...)
		}
		if err != nil {
			break
		}
	}
	issues = append(issues, decoder.Finish()...)
	return records, issues
}
