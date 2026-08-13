# ADR 0001: Separate model language from transport

Status: accepted  
Date: 2026-08-08  
Related: [proposal](../proposal.md), [WebSocket binding](../websocket-binding.md)

## Context

The model, harness, and browser have different constraints. Models benefit from
a small predictable grammar. The harness needs typed validation, policy,
durability, and replay. Clients need versioned network messages, snapshots, and
delivery semantics. Using one representation for all three couples token
generation choices to server and UI evolution.

## Decision

We will keep the Bragi model source language, canonical runtime event schema,
and server-to-client transport as three explicit interfaces with adapters
between them.

## Consequences

The model grammar can be benchmarked or replaced without changing replay and
client contracts. The WebSocket envelope can use strict JSON without forcing
the model to generate it. The runtime becomes the clear validation and
authority boundary.

The design adds adapters and requires tests proving equivalent meaning across
each boundary. Debugging must preserve source-to-canonical provenance without
leaking untrusted raw output into client authority.

Simplicity: generation, validation, persistence, and delivery remain
independently understandable; the adapter contracts are intentional coupling.

## Considered options

- **Status quo:** forward provider or tagged stream content toward clients.
  This is direct but makes UI and persistence depend on model syntax.
- **One JSON event format everywhere:** uniform tooling, but it makes models
  author transport metadata and binds grammar entropy to server concerns.

## Verification and disconfirming evidence

The Bragi 1.0 reference parser now produces source-free canonical records and
reconstructs materialized state from those events. A future Midgard adapter
must additionally prove equivalent meaning at its boundary before adoption.

Revisit if implementation shows the adapters duplicate most semantics or make
provenance materially harder to reason about than a shared representation.
