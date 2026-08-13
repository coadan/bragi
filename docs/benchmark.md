# Benchmark plan

Bragi should earn adoption through measured behavior, not syntax preference.

## Baselines

Run identical scenarios through:

1. provider-native strict tool calling or structured output;
2. one streamed JSON/JSONL structure;
3. Midgard's current tagged stream;
4. Bragi 1.0, prompted only;
5. Bragi with constrained decoding or a distilled protocol model, when
   available.

## Scenario set

- one valid tool call;
- two independent tool calls;
- a tool argument corrected before commit;
- generation truncated after a valid field and mid-record;
- an invalid path or type followed by repair;
- a tool rejection followed by a corrected new call;
- a report plus a large patch artifact;
- a plan amendment with stable step references;
- review findings, repair, and re-review;
- a false completion proposal rejected by the harness contract.

## Metrics

- time and tokens to first accepted field;
- time and tokens to first accepted commit;
- syntactic record acceptance rate;
- schema-valid and semantically correct action rate;
- malformed/truncated recovery rate without full regeneration;
- tokens and model calls spent on repair;
- duplicate or premature effect rate;
- deterministic replay rate;
- literal payload bytes per model token;
- constrained-decoder branching and latency overhead;
- task completion under harness-owned acceptance criteria;
- protocol success after small-model fine-tuning or distillation.

Report results per model and tokenizer. Do not collapse syntax validity and
semantic correctness into one score.

## Fault injection

For every format, inject equivalent failures at controlled semantic positions:

- terminate after entity creation, field name, and half a string;
- replace one operator or path segment;
- duplicate a record;
- reorder two dependent records;
- attempt to mutate an accepted effectful entity;
- disconnect and replay from a prior canonical sequence.

The grader is deterministic and harness-owned. Model claims about success are
not evidence.

## Adoption gate

Bragi 1.0 should not become Midgard's default unless, across at least two model
families, it improves recovery cost or time-to-committed-action without a
material regression in semantic validity, effect safety, or total task cost.

Large literal artifacts are a separate gate. If `|` continuation performs
poorly, retain the entity/commit protocol and replace only the literal channel.
