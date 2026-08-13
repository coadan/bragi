# ADR 0003: Stabilize the v1 source language

Status: accepted  
Date: 2026-08-08  
Related: [specification](../spec.md), [grammar](../../grammar/bragi.ebnf),
[benchmark](../benchmark.md)

## Context

Bragi needs a compatibility boundary before independent prompts, profiles,
adapters, and constrained decoders can target it. The draft grammar is small,
incrementally parseable across arbitrary provider chunks, and supports local
repair without rewriting history. Its principal uncertainty is comparative
model efficiency, especially for literal text, rather than parser semantics.

Changing syntax after adopters train or integrate against it would create a
larger coordination cost than keeping Bragi optional while measurements run.

## Decision

We will stabilize `+`, `~`, `-`, and `!`, stable entity IDs, JSON scalar and
reference values, and `|` literal continuation as the Bragi 1.x source
language. The grammar is closed unless an extension is explicitly negotiated.
Incompatible syntax or semantic changes require Bragi 2 and a superseding ADR.

## Consequences

Implementations can build against one deterministic grammar, and failures can
be compared across models without a moving protocol target. The profile and
host policy remain independently evolvable through negotiation.

V1 stability does not imply Midgard adoption or performance superiority. If
literal continuation performs poorly, Bragi 1.x remains valid but unsuitable
for that workload; replacement of the syntax waits for a major version.

Simplicity: four operators own four distinct transitions, one explicit literal
mode owns large text, and profiles own domain vocabulary.

## Considered options

- **Keep the whole protocol draft until benchmarks finish:** preserves syntax
  freedom but makes benchmark fixtures, integrations, and training targets
  unstable.
- **Stabilize only tool calls:** reduces surface area but couples the core to
  one use case and loses the heterogeneous-state hypothesis being tested.

## Verification and disconfirming evidence

The conformance suite verifies chunk independence, truncation, literal framing,
atomic rejection, repair, revisions, references, and replay. The benchmark may
disconfirm adoption, but not those semantics.

Revisit only for a Bragi 2 proposal supported by concrete incompatibility,
security, or model-reliability evidence.
