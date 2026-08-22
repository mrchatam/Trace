# P19-S02-01 — Implement `trace loop apply` + `trace loop status`

## Metadata
- id: P19-S02-01
- todo_ids: [P19-S02-01]
- role: implementer
- skills: [writing-for-agents]
- mcps: []
- verification: automated

## Objective

Implement structured loop writeback and saturation reporting.

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. This row may add product Go and tests, but must preserve forward-only planning history and keep board edits to status + notes.

## Locked defaults

| Item | Value |
|------|-------|
| Loop surface | `trace loop apply` writes structured results; `trace loop status` reports whether another planning iteration is warranted |
| Persistence | reuse existing add/link/plan/review/task state primitives where possible |
| Transport | harness-agnostic; apply accepts stdin JSON or `--in <path>` JSON, while loop remains stdout-first overall |
| Failure mode | malformed input must fail closed with non-zero exit and no ambiguous partial-success story |
| Saturation rule | saturated when no new tasks and no new plan changes were produced since the last loop step, or when max iterations is reached if that counter is part of the chosen design |
| Apply schema | required top-level `schema_version=trace.loop.apply.v1`, `apply_id` UUID, `seed.task_id`, `seed.goal_id`, and `writes` object |
| Replay lock | caller supplies stable IDs; writes replay safely via upsert-by-id + deduped links; no fuzzy dedupe by title/body |

## Envelope contract (locked)

`trace loop apply` must require this envelope shape:

- `schema_version`: string (`trace.loop.apply.v1`)
- `apply_id`: UUID (required, replay key)
- `seed`: object with required `task_id`, `goal_id` (plus optional scope metadata)
- `writes`: object with arrays:
  - `discoveries[]` (each requires stable `id`, `title`; optional body/severity/links)
  - `plan_changes[]` (each requires stable `id`; may reference a discovery and optional replan payload)
  - `spawned_tasks[]` (each requires stable `id`, `title`; goal-scoped)
  - optional `stop` object (`reason`, optional `max_iterations_reached`)

Reject envelopes missing required IDs or required top-level fields.

## Requirements

- `trace loop apply` ingests structured agent output and records it into Trace entities/links/transitions
- `trace loop status` reports whether more planning iterations are needed
- fail closed on malformed input
- preserve forward-only project history rules
- avoid creating duplicate entities on replay where reasonable

At minimum, the structured apply path should have a clear story for:

- discoveries
- plan changes or replan acknowledgments
- spawned tasks
- optional stop / saturation hints coming back from the agent

If replay safety is approximate rather than perfect, the implementation and tests must state the exact guarantee.

### Replay semantics to implement

- A repeated packet with identical `apply_id` must not create additional semantic writes.
- Use existing upsert-by-id semantics for entities and `InsertLinkOrIgnore` for duplicate-safe links.
- Persist one append-only loop-step summary event only after successful apply.
- `loop status` must compute from persisted loop-step evidence, never from in-memory counters.

## Verification expectations

- tests for malformed input, replay behavior, and saturation state transitions
- at least one plain CLI proof that `apply` can feed `status`
- no hidden dependency on a daemon, hosted coordinator, or special harness runtime
- explicit tests for:
  - missing `apply_id` / schema mismatch / missing item IDs => fail closed
  - replay same packet => stable counts, no duplicate links/tasks/plan_changes
  - status with no prior step => `saturated=false` + `reason=insufficient_history`
  - status saturated via zero-delta step and via max-iteration signal

## Exit criteria

- structured apply path works from plain CLI
- status correctly reports saturated / not saturated
- tests cover malformed input and replay/idempotency behavior chosen by the planner
