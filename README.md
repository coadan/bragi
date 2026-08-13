# Bragi

Bragi 1.0 is a protocol for models that need to stream structured state,
repair it in-band, and explicitly commit actions while generation is still in
progress.

The core idea is deliberately small:

```text
+ @t1 tool
+ @t1.name "search"
+ @t1.arguments.query "Qwen 3 benchmarks"
~ @t1.arguments.query "Qwen3.8 local inference benchmarks"
! @t1
```

`+` adds, `~` replaces, `-` retracts, and `!` proposes a commit. A runtime
materializes the current view as records arrive. It may accept or reject a
commit; only an accepted commit can dispatch a side effect.

Bragi is not another WebSocket or RPC format. The model emits a low-entropy
line language. A server validates that source into an append-only canonical
event log, and conventional transports project the log to clients.

Bragi 1.0 also has a closed deterministic recovery envelope for common model
surface mistakes: structural ASCII casing, CRLF, and a final LF omitted on
authoritative normal completion. Recovery is recorded, never changes payload
content, and grants no effect authority.

## Start here

- [Focused proposal](docs/proposal.md)
- [Bragi 1.0 specification](docs/spec.md)
- [Conformance and reference implementation](docs/conformance.md)
- [Evidence and readiness](docs/evidence.md)
- [WebSocket binding](docs/websocket-binding.md)
- [Midgard profile](profiles/midgard-v1.md)
- [Benchmark plan](docs/benchmark.md)
- [Decision records](docs/decisions/README.md)

## Status

The core grammar and semantics are stable for the Bragi 1.x line. The Go
reference implementation validates the shipped examples, materializes typed
state, and replays canonical events deterministically.

Bragi's performance advantage is still a hypothesis. It should not replace an
existing model protocol until the [benchmark gate](docs/benchmark.md) passes
across at least two model families. Literal-text efficiency remains a specific
benchmark target.

Run the conformance checks with:

```sh
go test ./...
go run ./cmd/bragi-check -profile profiles/midgard-v1.json examples/*.bragi
```

No license has been selected yet.
