# P20-S01-01 — Implement deliberation controller

## Metadata
- id: P20-S01-01
- todo_ids: [P20-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- verification: automated

## Objective
Implement `internal/deliberation` + persistence for deliberation state and transition events per S01-00 FINAL locks.

## Session start
Follow agent-loop-protocol. Unattended: execute after S01-00 is `done`. Board edits: status + notes only.

## Locked defaults (from S01-00 FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Package | `internal/deliberation` + store migration **`015_deliberation_state.sql`** |
| Hop budget | **12** |
| SelectNext | Priority table in S01-00 — pure function `SelectNext(PolicyInputs, State) → (Phase, ReasonCode, Stop)` |
| Events | `deliberation.transition` via `AppendEvent`; payload schema locked in S01-00 |
| EXECUTE | Must never be returned when `blocking_uncertainty_count > 0` |
| P19 loop | No changes to `internal/loop` or `cmd/trace/loop.go` in S01 |

## Requirements

- Library-first (Law 19). Unit tests in `internal/deliberation/*_test.go`.
- Store layer: add migration + CRUD for `deliberation_state` if S01-00 locked it.
- `SelectNext(state, inputs) → (phase, reasonCode, stop bool)` — pure function testable without CLI.
- Every transition persists `deliberation.transition` with full policy inputs JSON.
- Named table tests (minimum — map to S01-00 reason_codes):
  - `blocking_uncertainty_count=1` → `INVESTIGATE` / `blocking_uncertainty` (never EXECUTE)
  - `plan_exists=false`, uncertainties clear → `PLAN` / `plan_missing`
  - `plan_exists=true`, `plan_critiqued=false` → `CRITIQUE` / `plan_uncritiqued`
  - `open_regression=true` → `INVESTIGATE` / `open_regression`
  - `hop_count=12` → `STOP` / `hop_budget_exceeded`
  - `p19_saturated=true` → `STOP` / `p19_saturated`
- Do **not** break P19: `go test ./cmd/trace -run Loop` must pass (S01 should not touch loop CLI yet).

## Likely touch points

- `internal/deliberation/` (new)
- `internal/store/` (migration + deliberation_state if locked)
- `internal/domain/` (thin service wrapper if pattern matches existing create flows)

## Exit criteria

- [ ] SelectNext table tests from S01-00 green
- [ ] `deliberation.transition` events inspectable in store
- [ ] No daemon / no CoT blobs / no hosted MCP
- [ ] P19 Loop tests still green
