# Bragi WebSocket binding 1.0

This binding projects canonical server events to clients. It does not change
the model source language and is optional for Bragi conformance.

## Handshake and subscription

After transport authentication, the server announces its binding and task
protocol tuple:

```json
{
  "type": "hello",
  "v": "bragi.ws/1.0",
  "protocol": "bragi/1.0",
  "profile": "midgard/1.0",
  "profile_fingerprint": "sha256:..."
}
```

A client subscribes to one task after the highest canonical sequence it has
durably applied:

```json
{"type":"subscribe","task":"task_123","after":40}
```

The server sends events in canonical sequence order. Delivery may be at least
once; clients deduplicate by `(task, seq)`. If the requested sequence is no
longer replayable, the server sends a snapshot followed by its tail.

## Envelope

```json
{
  "v": "bragi.ws/1.0",
  "task": "task_123",
  "seq": 41,
  "kind": "op.accepted",
  "origin": "model",
  "time": "2026-08-08T12:00:00Z",
  "data": {
    "operation": "replace",
    "target": "@t1.arguments.query",
    "value": {"kind":"string","string":"Qwen3.8 local inference benchmarks"}
  }
}
```

JSON is suitable here because the server owns schema validity, sequence,
escaping, and versioning. Raw provider bytes and source lines may be retained
as access-controlled audit artifacts but are not forwarded as authority.

## Event families

The binding carries Bragi core events plus host-owned families such as:

- `effect.queued`, `effect.started`, `effect.progress`, `effect.finished`;
- `phase.changed` and `plan.amended`;
- `review.finding` and `review.resolved`;
- `completion.evaluated`;
- `human.input.required` and `human.input.received`;
- `budget.updated`; and
- `snapshot`.

An additive event kind requires binding negotiation. A client MUST disconnect
or request a supported projection when it receives an unnegotiated kind; it
must not silently invent its meaning.

## Projection and cadence

Clients are projections, not orchestration owners. They may render provisional
draft fields only after the entity's publication guards pass, and must
distinguish draft, committed, effect, rejection, and harness-verified terminal
states.

Entity IDs and revisions are stable presentation identity. Corrections update
the same component instead of reparsing a growing Markdown document. Servers
may coalesce high-frequency content changes for presentation cadence, provided
canonical persistence and replay remain lossless. Coalescing is never model
syntax or canonical state.

Clients acknowledge the highest contiguous sequence they applied. Slow clients
may move to snapshot-plus-tail delivery. Client backpressure does not pause
canonical persistence or grant control of model/tool execution unless the
harness exposes an explicit authorized pause command.

## Security

- Authenticate and authorize subscription independently of model roles.
- Treat every client command as an untrusted request with an idempotency key.
- Never treat a client message as committed model or runtime state.
- Gate source, command output, secrets, and private reasoning separately.
- Bound event and snapshot sizes.
