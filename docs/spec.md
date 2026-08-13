# Bragi 1.0 specification

Status: stable core protocol  
Date: 2026-08-08

The key words **MUST**, **MUST NOT**, **SHOULD**, **SHOULD NOT**, and **MAY**
are to be interpreted as described by RFC 2119 and RFC 8174.

## 1. Scope

Bragi defines a model-authored source language for incremental typed state and
a canonical event interface produced by a validating runtime. It does not
define provider token transport, storage encoding, orchestration policy, or
client network transport.

The 1.0 stability claim covers syntax and semantics. It does not claim that the
format is faster, cheaper, or more reliable than native structured output;
those remain benchmark questions.

## 2. Terms

- A **source record** is one complete model-authored operation ending in LF.
- An **entity** is a typed object with a session-local stable ID such as `@t1`.
- A **profile** declares the available entity and field schema for a domain.
- A **draft** is mutable entity state not accepted for its current revision.
- A **commit proposal** is `! @id`, requesting validation of that revision.
- A **canonical event** is a server-owned sequenced fact produced while
  validating source records.
- An **effect** changes state outside the materialized view, such as running a
  command, applying an edit, asking a human, or changing task lifecycle state.
- A **projection** is a client-facing view derived from canonical state.

## 3. Negotiation and layer boundary

Before generation, the host MUST establish this session tuple out of band:

- protocol ID, exactly `bragi/1.0`;
- profile name and version;
- SHA-256 fingerprint of the exact profile bytes;
- effective resource limits; and
- any separately negotiated extensions.

A continuation MUST retain that tuple and the entity namespace or start a new
session. Bragi has no in-band negotiation record.

Implementations MUST distinguish provider deltas, Bragi source records,
canonical events, materialized state, effects, and network projections. Raw
source MUST NOT be treated as an accepted event, and client messages MUST NOT
be fed back as Bragi source without explicit conversion.

## 4. Encoding and identifiers

Source is UTF-8. The canonical grammar uses LF record endings and lowercase
structural names. Because model output is nondeterministic, a conforming Bragi
1.0 decoder MUST perform the bounded deterministic recoveries in section 7
before validation. NUL is forbidden.

Whitespace outside literal continuations is exactly one ASCII space between
grammar terms. Values are RFC 8259 strings, numbers, booleans, `null`, or
entity references. Composite JSON arrays and objects are not core values.

Canonical names begin with lowercase ASCII and continue with lowercase ASCII,
digits, `_`, or `-`. Entity IDs prefix a name with `@`. A field path is an
entity ID plus one or more dot-separated names. Canonical identifiers are
case-sensitive after recovery.

The profile MUST provide positive limits for bytes per line, bytes per literal,
entities, and decoded records. The host SHOULD impose additional generation,
event, time, token, and effect budgets.

## 5. Source operations

### 5.1 Create

```text
+ @t1 tool
```

`+ @id type` creates revision 1 in draft state. The ID MUST be unused and the
type MUST exist in the profile.

### 5.2 Add

```text
+ @t1.name "repo.inspect"
+ @plan1.steps @step1
```

For a scalar, the field MUST be absent. For a collection, `+` appends a unique
entity reference in acceptance order. Collections contain only stable entity
references; positional array values are not supported.

### 5.3 Replace

```text
~ @t1.arguments.query "Qwen3.8 local inference benchmarks"
```

`~` replaces an existing scalar with a profile-valid value. It appends history
and changes only the materialized view.

### 5.4 Remove

```text
- @t1.arguments.limit
- @plan1.steps @step2
```

Removing a scalar names no value. Removing a collection member names the
existing reference. Entity removal is not supported in Bragi 1.0.

### 5.5 Literal text

```text
+ @p1.content |
|diff --git a/a.go b/a.go
|--- a/a.go
|+++ b/a.go
|
! @p1.content
```

`+ path |` opens a new literal string. `~ path |` clears an existing field and
opens its replacement. Only one literal may be open in a source session. While
open, each line MUST be either `|...` or a `! path` seal that canonicalizes to
the open path under section 7.

The runtime appends every byte after the first `|`, then LF. Thus `|` appends a
blank line and `||text` appends `|text`. Sealing does not commit the entity.
The field and profile limits still apply.

### 5.6 Commit proposal

```text
! @t1
```

The materializer MUST verify that the entity exists in draft state, required
fields are present, no literal is open, field values satisfy the profile, and
every reference targets a committed revision. It appends `commit.accepted` or
`commit.rejected`.

Acceptance makes an entity revision canonical and eligible for host policy; it
does not execute an effect. The host MUST still validate capability, policy,
budget, idempotency, and current task state before dispatch.

## 6. Profiles

A profile is a UTF-8 JSON object with this shape:

```json
{
  "protocol": "bragi/1.0",
  "name": "example",
  "version": "1.0",
  "limits": {
    "max_line_bytes": 16384,
    "max_literal_bytes": 262144,
    "max_entities": 1000,
    "max_records": 10000
  },
  "types": {
    "tool": {
      "mutation": "immutable",
      "effect": "host-action",
      "required": ["name"],
      "publication_guards": [],
      "field_order": ["name", "arguments.*"],
      "fields": {
        "name": {"kinds": ["string"], "cardinality": "scalar"},
        "arguments.*": {
          "kinds": ["string", "number", "bool", "null", "ref"],
          "cardinality": "scalar",
          "literal": true
        }
      }
    }
  }
}
```

Unknown profile properties MUST be rejected. `mutation` is `immutable` or
`revisioned`. `effect` is `none`, `host-action`, `human-question`, or
`completion-proposal`. Cardinality is `scalar` or `collection`; collections
permit only `ref`. A field pattern ending in `.*` matches descendants beneath
that path. Required fields and publication guards are exact field names.

`field_order` recommends a low-entropy generation order but does not change
valid dependency-respecting semantics. Tool capability names, argument
constraints, lifecycle policy, and completion contracts remain host-owned.

## 7. Atomic validation and recovery

The decoder MAY buffer only the incomplete current line and open literal
state. At LF, the materializer MUST atomically accept the complete record or
reject it without changing the view. Unknown operators, types, fields, paths,
and invalid values are rejected.

A Bragi 1.0 decoder MUST recover these source variations before profile or
materializer validation:

- replace a terminal CRLF record ending with canonical LF;
- fold uppercase ASCII to lowercase in entity IDs, entity type names, field
  path segments, and entity-reference values; and
- when the provider authoritatively reports normal completion, synthesize one
  missing terminal LF and decode the buffered bytes as one final record.

The normal-completion recovery MUST NOT be used for cancellation, transport
failure, timeout, or an unknown finish state. Every recovered record MUST carry
canonical normalization metadata identifying the recovery and its source and
canonical forms. Recovery happens before ordinary validation, so canonical ID
collisions, duplicate entities, unknown lowercased types or fields, invalid
JSON, and policy failures still reject normally.

This recovery set is closed. Implementations MUST NOT trim or collapse spaces,
change operators, fuzzy-match names, repair JSON, infer an unknown type or
field, select a semantic target, or change string, number, boolean, `null`, or
literal content. For example, `@T1.Name` may become `@t1.name`, while
`"Shell.Run"`, `TRUE`, command text, paths inside strings, and literal bytes
remain untouched. A runtime MAY expose a strict diagnostic mode, but normal
Bragi 1.0 ingestion uses the recovery rules above.

Provider interruption preserves accepted canonical events. An incomplete
ordinary line is diagnostic evidence only. An open literal and uncommitted
entities remain drafts. A host MAY resume with a bounded snapshot and the last
accepted sequence.

## 8. Revisions, references, and effects

The canonical log is append-only. Repairs update materialized state, never
prior events.

For an `immutable` type, mutations after commit are rejected. For a
`revisioned` type, the first valid mutation after commit opens the next draft
revision as a copy of the committed fields. Earlier revisions remain intact.

At commit, every entity reference is resolved to the target's current
committed revision and that resolution is recorded. Later target changes do
not mutate accepted inputs indirectly.

Every effectful type MUST be immutable after acceptance. A corrected effect
uses a new entity ID and may reference the superseded entity. Effects use an
idempotency identity derived from host session/task identity, entity ID, and
accepted revision—not raw model text. Replay MUST NOT re-execute effects.

## 9. Publication gates and client identity

A profile may name fields as `publication_guards`. A projection MUST withhold
an entity until all such fields contain accepted, closed values and any
host-owned routing predicates also pass. This prevents an incomplete speaker,
audience, or channel prefix from becoming visible content.

Once publishable, provisional draft updates MAY stream to clients. Clients
SHOULD key them by stable entity ID and revision so `~` repairs settle in place
instead of rebuilding an entire document or moving the viewport. Publication
does not imply commit, effect authorization, or task completion.

## 10. Canonical core events

Core events have a positive contiguous `seq` and a `kind`. Implementations may
encode them as JSON, but model source never supplies sequence or authority.

The Bragi 1.0 core kinds are:

- `source.rejected`: decoder diagnostic;
- `op.accepted`: canonical source record and entity ID;
- `op.rejected`: source record, entity ID, and diagnostic;
- `commit.proposed`: commit record and entity ID;
- `commit.accepted`: entity ID, revision, effect class, and pinned references;
- `commit.rejected`: entity ID and diagnostic.

Canonical records exclude raw model text. A recovered record includes its
bounded normalization metadata; replay consumes the canonical fields and
retains that metadata as provenance. Diagnostics contain a stable code, brief
message, and optional source line. Host events such as effect progress, phase
changes, and completion evaluation are separate event families and MUST NOT be
confused with these core validation facts.

## 11. Replay and completion

Replaying accepted operation and commit events in sequence MUST reconstruct the
same revisions, fields, pinned references, and effect identities. Rejected
events are evidence and do not affect state. Sequence gaps, unknown core kinds,
or a commit that cannot be reproduced MUST fail replay rather than be guessed.

Model EOF, provider finish reason, a committed report, or a committed
`completion` entity MUST NOT directly end a task. The host evaluates its own
completion contract and emits a server-authored lifecycle result such as
`done`, `incomplete`, `needs-human`, or `accepted-with-findings`.

## 12. Versioning and extensions

Protocol IDs use `bragi/<major>.<minor>`. A different major may be
incompatible. A minor release may only add behavior that 1.0 peers explicitly
negotiate; absent that negotiation, the 1.0 grammar is closed and unknown
operators are rejected.

Profiles evolve independently through their name, version, and exact-byte
fingerprint. Adding a type or field therefore requires a negotiated profile,
not parser inference. An implementation MUST refuse an unsupported protocol or
profile tuple before generation.

The [EBNF](../grammar/bragi.ebnf) is the syntax companion. The normative
behavior is verified by the [reference conformance suite](conformance.md).
