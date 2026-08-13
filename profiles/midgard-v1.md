# Midgard profile 1.0

This profile maps Bragi entities to a server-owned coding-harness lifecycle. It
is not part of the core grammar. The normative machine-readable schema is
[`midgard-v1.json`](midgard-v1.json).

## Lifecycle and authority

```text
idea -> discovery -> scope -> plan -> commit loop -> feature review -> QA
     -> completion check -> done | needs-human | accepted-with-findings
```

The server owns phase changes and completion evaluation. Models propose typed
entities; they do not assign canonical task state. One coding agent remains the
normal path. Planner and reviewer roles can be layered on as explicit harness
experiments rather than encoded in the core protocol.

## Entity groups

- `message` carries `speaker`, `audience`, `channel`, and `content`. The first
  three are publication guards: clients receive no provisional content until
  the server has validated its presentation route.
- `artifact` owns reports, patches, plans, ADRs, comprehension notes, and QA
  evidence. Large text uses literal mode.
- `tool` requests a host capability through shallow `arguments.*` fields.
  Accepted tools are immutable and effectful.
- `plan_step`, `finding`, `check`, `commit_unit`, and `comprehension` capture
  work and evidence with stable references.
- `option` and `question` support blocking human decisions. A recommended
  option is explicit data, not presentation convention.
- `completion` requests evaluation of the server's completion contract; it
  never ends a task directly.

The JSON profile contains the exact fields, required sets, mutation classes,
effects, recommended generation order, and resource limits. Capability names
and argument-specific validation remain server-owned because they vary by
deployment.

## Completion contract

The server evaluates at least whether all in-scope plan steps are resolved,
blocking findings are closed, required checks and QA passed, effects are
settled, human decisions are answered, workspace state is coherent, and final
deliverables and comprehension summaries exist.

Evaluation outcomes are `done`, `incomplete`, `needs-human`, or
`accepted-with-findings`. Only the last outcome requires a human acceptance
record naming the residual findings.

## Commit loop

Each planned commit unit is implemented, checked, independently reviewed,
repaired as needed, and accepted before dependent work begins. A material plan
amendment returns to scope or plan approval instead of silently expanding the
feature.
