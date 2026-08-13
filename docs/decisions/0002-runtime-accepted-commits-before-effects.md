# ADR 0002: Runtime-accepted commits before effects

Status: accepted  
Date: 2026-08-08  
Related: [specification](../spec.md), [benchmark](../benchmark.md)

## Context

Autoregressive output may expose an incomplete tool call, then correct an
argument. Dispatching on partial state can execute an unintended effect.
Waiting for provider EOF loses the latency and recovery benefits of incremental
state. Model intent is not sufficient authority for shell commands, edits,
human interruptions, or task completion.

## Decision

We will require an explicit model commit proposal and a separate runtime
acceptance before any Bragi entity can dispatch an external effect. Accepted
effectful entities are immutable; correction creates a new entity and effect
identity.

## Consequences

Draft tool state can render and self-repair safely. The exact accepted revision
is auditable and can supply an idempotency key. Truncation before commit cannot
dispatch an effect, and replay does not re-execute effects.

Every action incurs a commit record and validation step. Models must learn when
to commit, and overly strict validation could add latency or stall useful work.

Simplicity: mutable intent and immutable effects are separate states; policy
remains in the runtime rather than being interleaved with model syntax.

## Considered options

- **Dispatch when a tool object first becomes schema-valid:** lower latency,
  but a later model repair races with an already-started effect.
- **Dispatch only at provider EOF:** simple lifecycle, but prevents early action
  and makes unrelated output gate the call.

## Verification and disconfirming evidence

Reference tests prove that partial and rejected records do not commit, repairs
can occur before commit, accepted effects are immutable, and replay does not
execute effects. Host adapters must additionally prove idempotent dispatch and
duplicate delivery handling before production use.

Revisit if commit latency materially dominates short safe calls and a narrower
profile can prove equivalent safety with automatic commit.
