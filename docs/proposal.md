# Proposal: Bragi model-native structured streaming

Bragi 1.0 is a small model emission language for incrementally constructing
typed entities, repairing draft state, and explicitly proposing commits. A
server validates those records into an append-only event log and remains
the authority for effects, workflow transitions, and task completion.

## Why

- Autoregressive models produce useful information sequentially, but common
  structured interfaces treat incomplete output as unusable until a complete
  JSON object or tool call arrives.
- Midgard's tagged stream already proves the value of mixed reports, payloads,
  commands, budgets, and repair, but its prose and control syntax share one
  line grammar and repairs can require replacement generations.
- Existing provider and orchestration streams expose partial values or runtime
  events; Bragi's narrower hypothesis is that the model itself can author
  replayable state changes and repair them before committing effects.

## Core proposal

Use four primary model operators:

```text
+ <target> <value>   add an entity, field, or collection member
~ <target> <value>   replace an existing draft value
- <target> [value]   retract a field or collection member
! <target>           seal a literal or propose an entity commit
```

Entities have stable IDs such as `@t1` and profile-defined types such as
`tool`, `artifact`, `finding`, `question`, or `completion`. Fields form a
shallow typed graph; collections contain stable references rather than array
indexes.

Every completed source record is independently accepted or rejected. Accepted
records update a speculative materialized view immediately. The source log
never changes: `~` and `-` add corrective events rather than rewriting prior
events.

Presentation is gated separately. For example, the Midgard profile withholds a
message until its speaker, audience, and channel are validated, then preserves
the entity ID while later content repairs settle in place. This borrows a
useful safety boundary from Void2 without importing its role grammar or
full-response repair protocol.

`! @t1` is a proposal, not permission. The materializer validates the entity's
profile and current revision. A `commit.accepted` event only makes it eligible
for host capability, policy, budget, and idempotency checks; an effect cannot
dispatch before both boundaries pass.

## Three distinct authorities

### Model

The model may construct and repair entities, request tools or human input,
amend a plan, report findings, and propose completion.

### Harness

The harness owns canonical sequence numbers, validation, policy, side effects,
phase transitions, budgets, review gates, persistence, and the completion
contract. Provider EOF and model completion proposals are observations, not
terminal conditions.

### Human

The human owns material scope decisions and can explicitly accept known
findings. Ordinary ambiguity should first be reduced through inspection and
bounded probes; only genuinely blocking choices become human questions.

## Coding-harness lifecycle

Bragi carries typed proposals and evidence through a harness-owned lifecycle:

```text
idea -> discovery -> scope -> plan -> commit loop -> feature review -> QA
     -> completion check -> done | needs human | accepted with findings
```

Within the commit loop, each change has an intent, bounded diff, checks,
independent review, repair cycle, and acceptance decision. The harness keeps
working while the completion contract reports remaining plan units, blocking
findings, failed checks, unresolved effects, missing deliverables, or pending
human decisions.

## Comprehension as a deliverable

Each accepted coding commit should produce a compact comprehension entity that
answers:

- what changed and why;
- which ownership or architectural boundary moved;
- which invariant or operational assumption is new;
- how the change was verified;
- what a maintainer must understand before accepting it.

These entities link to diffs, checks, findings, and ADRs rather than copying
them. The final feature summary is a projection over accepted commit
comprehension records, so rapid agent output does not leave one giant
after-the-fact explanation.

## Layering

```text
model token deltas
    -> Bragi source parser
    -> profile/schema/policy validation
    -> append-only canonical event log
    -> materialized task view
    -> effects and harness state machine
    -> WebSocket projections and replay
```

This separation keeps the model language optimized for generation while the
server and clients use explicit typed JSON. It also lets Midgard replace the
model syntax, WebSocket implementation, or UI independently.

## Non-goals

Bragi 1.0 does not attempt to:

- replace MCP, JSON-RPC, WebSocket, SSE, or provider APIs;
- encode private chain of thought;
- let a model authorize its own side effects or completion;
- be a general-purpose storage serialization format;
- stream arbitrary binary data;
- standardize Midgard's entire workflow in the core grammar;
- claim novelty for streaming, patching, constrained decoding, or event logs
  individually.

## Novelty claim to test

The defensible hypothesis is the combination: a model-native, append-oriented
language whose accepted records materialize heterogeneous typed state, whose
model-authored repairs update that state in-band, and whose explicit commits
separate speculative intent from runtime-authorized effects.

The claim becomes useful only if experiments show better recovery,
time-to-useful-state, or small-model reliability than existing interfaces.

## V1 boundary

V1 stabilizes the four operators, literal framing, profile format,
revision-pinned materialization, commit semantics, canonical core events, and
negotiation rules. It includes a reference decoder, materializer, replay path,
profile, fixtures, and conformance command.

It does not establish an adoption claim. The next step is an optional Midgard
adapter followed by identical-trace benchmarks against current tags, native
tool calls, and JSON/JSONL. Bragi should remain experimental in a host until
that gate passes.

## Watch

The largest uncertainty is whether operation and path overhead outweighs the
saved retries, especially for large text artifacts. Literal encoding must be
measured by tokenizer, model, task shape, and constrained-decoding mode.
