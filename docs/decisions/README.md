# Decision records

Decisions are sequential and never renumbered. A superseded record remains in
place and links to its replacement.

- [ADR 0001: Separate model language from transport](0001-separate-model-language-from-transport.md) — accepted
- [ADR 0002: Runtime-accepted commits before effects](0002-runtime-accepted-commits-before-effects.md) — accepted
- [ADR 0003: Stabilize the v1 source language](0003-stabilize-the-v1-source-language.md) — accepted
- [ADR 0004: Recover deterministic source variations](0004-recover-deterministic-source-variations.md) — accepted

Protocol compatibility and host adoption are separate decisions. The source
language is stable for 1.x; [benchmark](../benchmark.md) results still decide
whether a host should adopt it.
