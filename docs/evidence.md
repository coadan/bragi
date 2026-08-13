# Evidence and readiness

## Current state

- **Fact:** Midgard v1 accepts top-level report, payload, edit, reference,
  command, result, and error tags, with byte/frame budgets and bounded repairs.
  Evidence: `internal/stream/` and `internal/model/packet.go` in the local
  Midgard repository at design time.
- **Fact:** Midgard has recovery logic for inline controls, malformed known
  tags, missing results, open payloads, disallowed role frames, and result-only
  repair. Evidence: `internal/stream/parser_test.go` and
  `internal/stream/repair.go`.
- **Fact:** Provider APIs stream typed lifecycle and content-block deltas, and
  tool arguments may remain partial JSON until a block closes. Evidence:
  [OpenAI streaming events](https://platform.openai.com/docs/api-reference/responses-streaming/response/refusal/delta),
  [Anthropic streaming messages](https://platform.claude.com/docs/en/build-with-claude/streaming),
  and [DeepSeek Chat Completions](https://api-docs.deepseek.com/api/create-chat-completion/).
- **Fact:** Sequential patch operations and incremental constrained parsing are
  established techniques. Evidence: [RFC 6902 JSON Patch](https://www.rfc-editor.org/rfc/rfc6902.html),
  [PICARD](https://aclanthology.org/2021.emnlp-main.779/), and
  [Efficient Guided Generation](https://arxiv.org/abs/2307.09702).
- **Fact:** Agent runtimes already expose typed state updates and persisted
  checkpoints. Evidence: [LangGraph streaming](https://docs.langchain.com/oss/python/langgraph/streaming)
  and [persistence](https://docs.langchain.com/oss/javascript/langgraph/persistence).
- **Fact:** Void2's local protocol separates provider transport from canonical
  reducer inputs, withholds partial speaker/audience/channel headers from
  presentation, preserves stable visible-frame identity across correction,
  and limits normalization to incidental syntax. Evidence inspected at design
  time: `docs/architecture-model-protocol-reference.md`,
  `worker/line-protocol.ts`, `worker/protocol-parser.ts`,
  `worker/protocol-visible-frames.ts`, and `worker/repairable-stream.ts` in the
  local Void2 repository.
- **Fact:** Bragi's reference tests now cover arbitrary chunk boundaries,
  truncation, atomic rejection, repair, publication gates, immutable effects,
  pinned revisions, example conformance, and canonical replay. Evidence:
  `decoder_test.go`, `materializer_test.go`, and `cmd/bragi-check`.
- **Unknown:** Whether a model-authored patch language beats native strict tool
  calling after counting tokens, semantic errors, validation latency, and
  recovery cost. Next probe: the comparative benchmark.

## Diagnosis

The supported problem is not simply that JSON can be malformed. The deeper
problem is that model generation state, harness state, effect authorization,
and client presentation are commonly forced through one completion-shaped
interface. Partial information then has weak semantics, corrections require
regeneration or parser inference, and a model stop can be mistaken for task
completion.

Bragi separates those concerns: the model authors tentative state changes; the
runtime validates and commits canonical events; policy authorizes effects; the
workflow evaluates completion; clients render projections.

## Scope

**Problem:** Coding harnesses cannot safely consume, expose, repair, and resume
heterogeneous model output as it is generated because current model-facing
formats conflate partial syntax, mutable intent, committed effects, and final
task state.

**In scope:** a model emission grammar, typed profiles, incremental validation,
materialization, repair, commit semantics, canonical events, replay, and a
WebSocket binding.

**Out of scope:** general RPC, provider transport, workflow policy, binary
artifact transfer, model training recipes, and a claim that every output should
use Bragi.

**Invariants:** history is append-only; stable IDs replace indexes; rejected
records do not alter state; side effects require accepted commits; completion
is harness-owned; replay yields the same view.

## Direction comparison

| Criterion | Midgard tags/status quo | One streamed JSON object | Bragi patch stream |
|---|---|---|---|
| Useful partial state | Text and payload bytes stream, but controls share the report grammar | Usually requires partial parsing and later closure | Each accepted line updates a typed view |
| Repair | Bounded continuation or replacement repair | Parser repair or regenerate/patch the object | Model-authored `~`/`-` records preserve history |
| Side-effect boundary | Command/result ordering conventions | Tool-call completion | Explicit runtime-accepted entity commit |
| Truncation | Draft artifacts plus repair packet | Often incomplete structure | Accepted records survive; incomplete record is ignored |
| Model entropy | Small tags, but flexible fields and mixed modes | Braces, nesting, quoting, schema order | Four primary operators and profile-constrained paths |
| Compatibility | Already implemented in Midgard | Broad provider support | Requires parser/profile adapters |
| Main uncertainty | Repair complexity as features grow | Partial JSON and retry cost | Token overhead and semantic patch errors |

## V1 release gate

**Verdict:** protocol boundary passed; host adoption remains conditional.  
**Commitment level:** consequential but reversible while Bragi remains an
optional Midgard adapter.

**Why deepen now:** Midgard supplies a concrete baseline and failure corpus;
Void2 provides a second local comparison; the source/runtime/client boundary
and commit invariant now have executable checks; and v1 remains optional to
Midgard.

**Could invalidate this:** native strict tool calling may equal or beat Bragi
on validity, latency, and recovery; models may misuse patch targets; literal
records may cost too many tokens.

**Stop or pivot:** do not adopt Bragi as Midgard's default if it fails to
improve either recovery cost or time-to-committed-action without a material
semantic-validity regression across two model families.

**Smallest reversible test completed:** the reference parser, materializer,
replay fixtures, profile validator, and pre-commit correction tests.

**Next gate:** a side-by-side benchmark with forced truncation and identical
semantic tasks. V1 protocol stability does not authorize Midgard adoption.

**Independent verification:** deterministic replay and schema checks, token
counts from each provider tokenizer, and benchmark grading from harness-owned
outcomes rather than model claims.

## Decision ledger

| ID/date | Decision or question | Evidence | Status | Revisit trigger |
|---|---|---|---|---|
| D-001 / 2026-08-08 | Separate model syntax, canonical events, and client transport | Different owners and failure modes; Midgard and Void2 both require provider adapters | accepted for 1.x | Adapters duplicate semantics or obscure provenance |
| D-002 / 2026-08-08 | Require runtime acceptance of explicit commits before effects | Tested correction and immutable post-commit behavior | accepted for 1.x | Commit latency dominates and a narrower rule proves equivalent safety |
| D-003 / 2026-08-08 | Stabilize `+`, `~`, `-`, `!` with stable entity IDs | Chunk-independent parser and deterministic fixtures | accepted for 1.x | Incompatible change is justified for Bragi 2 |
| D-004 / 2026-08-08 | Prefix literal continuation lines with `|` | Exact parser behavior and payload fixtures | accepted for 1.x; adoption experimental | Benchmark shows material token or reliability loss; replace in Bragi 2 |
| D-005 / 2026-08-08 | Recover only a closed set of deterministic source variations before validation | Void2 casing/CRLF normalization and Bragi conformance tests | accepted for 1.0 | Recovery changes payload meaning, hides ambiguity, or harms cross-runtime replay |
