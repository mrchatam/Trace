# P20-S01-00 — Deliberation controller planner

## Metadata
- id: P20-S01-00
- todo_ids: [P20-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock `internal/deliberation`: phase vocabulary, `deliberation_state`, deterministic `SelectNext`, entry/exit, hop budget, `deliberation.transition` events. **No product Go this row.** Wrap P19 loop; do not replace it.

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [../../00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- Live: `internal/loop/{next,apply}.go`, `internal/store/entities.go` (`AppendEvent`), `internal/store/migrate.go`, `internal/domain/create.go`, `internal/planner`

## Doc map
§1, 3B, 4, 5, 6, 21, 25, 28, 29D, 29N

## Live inventory to reuse

| Surface | Location | S01 use |
|---------|----------|---------|
| Loop packet/apply | `internal/loop/{next,apply}.go` | S01 does **not** change packet shape — S06 integrates phase into next/apply |
| Events | `store.AppendEvent`, `events` table | `deliberation.transition` payload on seed task |
| Planner replan budget | `internal/planner` N=5 | Separate from hop budget — do not conflate |
| P19 keeper tests | `cmd/trace/loop_test.go` | Must stay green through S01 (library-only) |
| Next migration id | `internal/store/schema/` max **014** | S01-01 adds **`015_deliberation_state.sql`** |

## FINAL locked defaults (S01-01 must not re-debate)

| Item | Value |
|------|-------|
| Phases | `ORIENT` `INVESTIGATE` `EXPLORE` `PLAN` `CRITIQUE` `EXECUTE` `TEST` `VERIFY` `EVALUATE` `REFLECT` `REPLAN` + terminal `STOP` |
| Library | `internal/deliberation` — pure policy + types; persistence in `internal/store` |
| State key | One row per `task_id` (unique); `goal_id` denormalized for queries |
| Hop budget N | **12** transitions per seed before fail-closed `STOP` |
| SelectNext | Pure function; priority table below |
| EXECUTE gate | If `blocking_uncertainty_count > 0`, SelectNext must **never** return `EXECUTE` |
| Open regression default | `open_regression=true` → `INVESTIGATE` (not REPLAN) unless replan flag set in inputs |
| Event type | `deliberation.transition` |
| Event entity | `entity_type=task`, `entity_id=seed.task_id` |
| CLI in S01 | **Library-only** — no `cmd/trace` changes; S06 owns loop integration |

### `deliberation_state` table (migration 015)

```sql
CREATE TABLE IF NOT EXISTS deliberation_state (
    task_id TEXT PRIMARY KEY,
    goal_id TEXT NOT NULL,
    current_phase TEXT NOT NULL DEFAULT 'ORIENT',
    hop_count INTEGER NOT NULL DEFAULT 0,
    last_phase TEXT NOT NULL DEFAULT '',
    plan_critiqued INTEGER NOT NULL DEFAULT 0,
    stopped INTEGER NOT NULL DEFAULT 0,
    stop_reason TEXT NOT NULL DEFAULT '',
    updated_at TEXT NOT NULL
);
```

- `plan_critiqued`: 0/1 boolean stored as INTEGER (SQLite convention).
- `stopped=1`: SelectNext returns `STOP` with `stop_reason=hop_budget_exceeded` or policy terminal.

### `deliberation.transition` event payload

JSON fields (required):

- `task_id`, `goal_id`
- `from_phase`, `to_phase`
- `reason_code` (string enum — see table below)
- `hop_count` (after transition)
- `policy_inputs` object:
  - `blocking_uncertainty_count`
  - `plan_exists`
  - `plan_critiqued`
  - `verification_incomplete`
  - `open_regression`
  - `p19_saturated`

### SelectNext priority table (FINAL)

Evaluated top-to-bottom; first match wins:

| Priority | Condition | Phase | reason_code |
|----------|-----------|-------|-------------|
| 1 | `hop_count >= 12` OR `stopped` | `STOP` | `hop_budget_exceeded` |
| 2 | `p19_saturated` | `STOP` | `p19_saturated` |
| 3 | `blocking_uncertainty_count > 0` | `INVESTIGATE` | `blocking_uncertainty` |
| 4 | `open_regression` | `INVESTIGATE` | `open_regression` |
| 5 | `!plan_exists` | `PLAN` | `plan_missing` |
| 6 | `plan_exists && !plan_critiqued` | `CRITIQUE` | `plan_uncritiqued` |
| 7 | `verification_incomplete` | `VERIFY` | `verification_incomplete` |
| 8 | default | `ORIENT` | `continue_orient` |

**EXECUTE** is not selected by this MVP table — S06 may add EXECUTE when prior phases satisfied and no blocking inputs (future row if needed). S01-01 tests must prove EXECUTE never returned when blocking uncertainty > 0 even if caller requests EXECUTE path simulation.

### PolicyInputs sourcing (S01-01)

S01 library accepts `PolicyInputs` struct populated by caller (S06 will wire queries). Fields may be stubbed to zero/false in unit tests; integration tests come in S06.

## Planner work

1. [x] Confirm hop budget N = 12.
2. [x] Lock `deliberation_state` migration 015 shape.
3. [x] Lock event payload schema.
4. [x] Thicken `01-deliberation-controller.md` and `02-scope-review.md`.
5. [x] Update `SCOPE-TODOS.md`.

## Exit criteria

- [x] 01/02/SCOPE-TODOS thickened enough to implement alone
- [x] Hop budget N locked (12)
- [x] SelectNext table locked with named test rows
- [x] P19 Loop keeper tests listed for S01-02 review
- [x] No product Go

## Next

Orchestrator: **P20-S01-01** after this row is `done`.
