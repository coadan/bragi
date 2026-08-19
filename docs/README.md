# Bragi documentation

Read in this order:

1. [Proposal](proposal.md) — problem, direction, and scope.
2. [Evidence](evidence.md) — current evidence, alternatives, and readiness.
3. [Specification](spec.md) — normative Bragi 1.0 semantics.
4. [Conformance](conformance.md) — reference implementation and verification.
5. [WebSocket binding](websocket-binding.md) — replayable client delivery.
6. [Benchmark](benchmark.md) — the test that can invalidate adoption.
7. [Decisions](decisions/README.md) — durable architectural choices:
   [model language versus transport](decisions/0001-separate-model-language-from-transport.md)
   [runtime-accepted commits](decisions/0002-runtime-accepted-commits-before-effects.md),
   the [stable v1 source language](decisions/0003-stabilize-the-v1-source-language.md),
   and [bounded deterministic recovery](decisions/0004-recover-deterministic-source-variations.md).

The [Midgard profile](../profiles/midgard-v1.md) explains the coding-harness
mapping; its [machine-readable form](../profiles/midgard-v1.json) is used by
the reference implementation. The [EBNF grammar](../grammar/bragi.ebnf) and
[examples](../examples/) accompany the specification.
