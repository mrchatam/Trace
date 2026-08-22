# P19-S02-02 — Review `trace loop apply/status`

## Metadata
- id: P19-S02-02
- todo_ids: [P19-S02-02]
- role: reviewer
- skills: [documentation-and-adrs, writing-for-agents]
- mcps: []
- verification: automated

## Objective

Review loop writeback and saturation semantics for correctness, replay safety, and usefulness to future looped planning.

## Review focus

The reviewer should verify:

- writeback reused existing Trace entities/links/planner transitions instead of inventing shadow state
- malformed input fails closed
- replay/idempotency behavior matches the planner lock and is evidenced by tests
- `loop status` saturation is computed from explicit repo state, not hand-waved process memory
- forward-only history rules were preserved
- no daemon / hosted / stdin-primary loop contract slipped in
- apply envelope contract is enforced (`trace.loop.apply.v1`, required `apply_id`, required stable item IDs)
- loop-step evidence is persisted append-only and queried by status path
- saturation reasons are explicit (`tasks_and_plan_unchanged`, `max_iterations_reached`, `insufficient_history`)

## Required evidence checks

- CLI proof for file input and stdin input forms of `loop apply`
- replay proof: apply same packet twice and assert no semantic duplication
- status proof with:
  - no prior step
  - zero-delta saturation
  - max-iteration saturation
- tests assert that errors never emit success-style partial envelopes

## Exit criteria

- blocker/high issues fixed or spawned forward
- saturation rule is evidenced, not hand-waved
