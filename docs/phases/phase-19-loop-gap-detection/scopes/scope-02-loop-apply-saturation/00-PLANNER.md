# P19-S02-00 — Loop apply/status planner

## Metadata
- id: P19-S02-00
- todo_ids: [P19-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Grep, Glob, Write]
- verification: automated

## Objective

Lock the MVP behavior for `trace loop apply` and `trace loop status`. This scope owns the tasks-saturated stop rule.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [../../00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- Live seams: `cmd/trace/plan.go`, `cmd/trace/impact.go`, `cmd/trace/tasks.go`, existing add/link/transition/review surfaces as needed

## Locked defaults

- `trace loop apply` records structured outputs from an agent:
  - discoveries
  - plan changes
  - spawned tasks
  - optional stop reason
- `trace loop status` reports whether the loop is saturated
- MVP saturated = zero new tasks and zero new plan changes since last loop step (or max iteration reached)
- Keep transport harness-agnostic; file/stdin ingestion is acceptable only as input to apply, not as the primary core interaction model
- Reuse existing domain/planner writes where possible; do not invent a parallel persistence model for discoveries, plan changes, or tasks
- Apply must be structured and fail closed: partial writes without an explicit success envelope are a defect

## Planner questions to lock

S02-00 must leave S02-01 with explicit answers for:

- accepted input envelope shape for `loop apply`
- whether input is stdin, file path, argument, or a combination, while keeping stdout-first overall loop semantics
- minimum writeback set: discoveries, plan changes, spawned tasks, optional stop reason
- replay/idempotency policy for repeated apply on the same packet
- exact status output fields for `loop status`
- evidence source for saturation: what counts as "no new tasks" and "no new plan changes"

## Locked answers (resolved against live repo seams)

### `trace loop apply` input envelope

- Schema: `trace.loop.apply.v1`
- Accepted input modes: **exactly one** of:
  - `--in <path>` (JSON file)
  - stdin JSON stream (when `--in` omitted)
- Top-level required fields:
  - `schema_version` (must equal `trace.loop.apply.v1`)
  - `apply_id` (UUID; replay key)
  - `seed` object with `task_id` and `goal_id`
  - `writes` object (may contain empty arrays but object must exist)
- `writes` supports:
  - `discoveries[]`
  - `plan_changes[]`
  - `spawned_tasks[]`
  - optional `stop` object (`reason`, optional `max_iterations_reached`)

### Replay/idempotency policy (MVP lock)

- `apply_id` is required and must be persisted as the loop-step identity.
- Every write item must carry explicit stable `id`; if missing, fail closed.
- Entity writes reuse existing upsert-by-id behavior (`Create*` + `store.Upsert*`) for deterministic replay.
- Link writes must use dedupe-safe inserts (`InsertLinkOrIgnore`) so repeated apply packets do not multiply identical edges.
- Replaying the same envelope is expected to be a no-op at semantic level (same entity ids + deduped links + same loop-step summary).
- This is **replay-safe idempotency by caller-supplied IDs**, not fuzzy dedupe by title/body.

### Apply writeback contract

- Apply executes in three phases:
  1. Parse + validate full envelope (including required IDs and enum checks).
  2. Apply writes through existing primitives (discoveries/tasks/plan-change + planner apply-discovery where requested).
  3. Persist one loop-step summary event only after successful write phase.
- Failure mode is fail-closed:
  - malformed/invalid payload -> non-zero, no writes.
  - write/apply error -> non-zero, no success envelope and no step summary event.
- Success envelope must report counts:
  - `new_discoveries`
  - `new_plan_changes`
  - `new_spawned_tasks`
  - `saturated`
  - `replay` (`true` only when prior step with same `apply_id` already exists)

### `trace loop status` output contract

- Schema: `trace.loop.status.v1`
- Required output fields:
  - `schema_version`
  - `seed.task_id`
  - `seed.goal_id`
  - `last_apply_id` (or empty when none)
  - `new_plan_changes_since_last_step`
  - `new_tasks_since_last_step`
  - `max_iterations_reached`
  - `saturated`
  - `reason` (`tasks_and_plan_unchanged`, `max_iterations_reached`, or `insufficient_history`)

### Saturation evidence source (MVP lock)

- Source of truth is persisted loop-step events (append-only `events` row with loop payload), not process memory.
- Saturated when either:
  - `new_tasks_since_last_step == 0` AND `new_plan_changes_since_last_step == 0`, or
  - `max_iterations_reached == true`.
- If no prior loop step exists for the seed/goal, `saturated=false` with reason `insufficient_history`.

## Exit criteria

- [x] S02 implement/review stubs are FINAL
- [x] saturation semantics locked
- [x] apply envelope and replay stance locked
- [x] no product Go
