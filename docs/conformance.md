# Conformance

Bragi 1.0 ships a small Go reference implementation at the repository root.
It owns four independently testable behaviors:

- `Decoder` incrementally turns arbitrary byte chunks into complete source
  records without applying partial lines.
- `Profile` loads and validates the exact negotiated machine-readable schema.
- `Materializer` atomically accepts or rejects records and produces canonical
  events and revision-pinned state.
- `Replay` reconstructs the same materialized state from canonical events
  without executing effects.

`cmd/bragi-check` is the conformance entry point for profiles and source
fixtures:

```sh
go run ./cmd/bragi-check \
  -profile profiles/midgard-v1.json \
  examples/*.bragi
```

Use `-json` to inspect canonical events as JSON Lines. A conforming stream exits
zero only when every record and commit is accepted, every literal is sealed,
and no draft entity remains at EOF.

## Required properties

The automated suite verifies:

- identical decoding across whole-stream and one-byte chunks;
- auditable canonical recovery of structural ASCII case and CRLF without
  changing string or literal payloads;
- terminal-LF recovery only after authoritative normal provider completion;
- survival of the accepted prefix after truncation;
- byte-exact literal continuation and canonical-equivalent sealing;
- pre-commit repair and post-commit immutability for effects;
- atomic rejection without view mutation;
- revision-pinned references;
- publication gating before routing fields are complete;
- deterministic canonical replay; and
- conformance of every shipped example to the Midgard profile.

Run all checks with `go test ./...`. These checks establish semantic
conformance, not model quality or token efficiency; those claims belong to the
[benchmark](benchmark.md).
