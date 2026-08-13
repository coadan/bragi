# ADR 0004: Recover deterministic source variations

Status: accepted  
Date: 2026-08-08  
Related: [specification](../spec.md), [grammar](../../grammar/bragi.ebnf),
[conformance](../conformance.md)

## Context

Models generate nondeterministically and can vary incidental syntax even when
their intended operation is unambiguous. Rejecting every case-only identifier
variation, CRLF ending, or terminal record missing only its LF creates repair
turns without protecting meaning. Void2 showed that casing and line-ending
normalization can recover syntax while downstream validation retains authority.

Broad parser forgiveness would create the opposite problem: a runtime could
silently choose an operation, identifier, value, or tool argument the model did
not actually emit.

## Decision

We will make bounded deterministic recovery part of Bragi 1.0 ingestion. The
decoder will canonicalize ASCII case only in structural names and references,
normalize CRLF, and accept a missing final LF only after authoritative normal
provider completion. Canonical records will retain normalization provenance.

We will not normalize payload values, whitespace, operators, JSON syntax, or
unknown schema names. Recovered records pass through the same profile,
materializer, collision, commit, and host-policy checks as canonical source.

## Consequences

Common surface mistakes no longer require another model turn, and all clients
materialize one canonical namespace. String and literal bytes remain exact;
ambiguous or semantic errors still fail closed and can be repaired in-band.

Accepting previously invalid case variants expands the v1 input envelope but
does not change the meaning of any already-valid Bragi 1.0 source. Strict mode
remains useful as a diagnostic, not as normal conforming ingestion.

## Verification and disconfirming evidence

Conformance tests cover one-byte chunking, casing across IDs/types/paths/refs,
literal seals, payload preservation, CRLF provenance, strict diagnostics, and
the distinction between normal-completion and abrupt EOF.

Revisit if two implementations canonicalize the same bytes differently, if
normalization masks a semantic collision, or if evidence supports another
closed recovery that can preserve the same guarantees.
