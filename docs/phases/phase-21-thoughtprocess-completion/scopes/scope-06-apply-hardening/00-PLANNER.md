# P21-S06-00 — Planner: loop apply hardening

## Metadata
- id: P21-S06-00
- todo_ids: [P21-S06-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock transactional `loop apply`, seed/task `goal_id` validation on deliberation transition, and verify-floor test location hygiene (D-15). **No product Go this row.**

## References
- [DECISION-LOG.md](../../DECISION-LOG.md) D-08, D-13, D-15
- [WORK-MAP.md](../../WORK-MAP.md) W-09, W-10, W-11
- P19/P20 apply baseline: [scope-06/00-PLANNER.md](../../../phase-20-cognitive-deliberation/scopes/scope-06-protocol-context/00-PLANNER.md)
- Live: `internal/loop/{apply.go,apply_writes.go,policy.go}`, `internal/domain/deliberation.go`, `internal/store/open.go`, `cmd/trace/loop_test.go`

## Live inventory (confirmed 2026-08-18)

| Surface | Location | Today (live read) | S06 action |
|---------|----------|-------------------|------------|
| `loop.Apply` | `apply.go` L367–517 | Sequential domain/store writes; **no transaction**; replay check before writes (L378–386) | **Wrap** all mutating steps in single store transaction |
| Seed goal guard | `apply.go` L375–377 | `task.GoalID` must equal `env.Seed.GoalID` before writes | **Keep** — runs **outside** tx (read-only) |
| Status goal guard | `apply.go` L534–542 | Same mismatch check on status path | **Keep** unchanged |
| `ApplyDeliberationTransition` | `deliberation.go` L17–65 | Requires non-empty `task_id`/`goal_id`; loads task but **does not** verify goal matches task | **Add** fail-closed: `goal_id` must equal task's stored `goal_id` |
| Cognitive writes | `apply_writes.go` L45–92 | Sequential per-artifact handlers after P19 writes | **Inside** same tx as P19 writes + transition + loop.step |
| Store transactions | `fts.go`, `vcs_index.go`, `file_graph.go` | Ad-hoc `db.Begin()` / `Commit()` / `Rollback()` in isolated methods | **Add** `store.WithTx` (or equivalent) for loop apply — no full-store refactor |
| Loop unit tests | `internal/loop/` | **0** `*_test.go` files | **Add** `apply_test.go` with **≥8** behavioral tests (D-15) |
| CLI loop tests | `cmd/trace/loop_test.go` | **27** `TestLoop*` funcs; **14** apply/status/next S06-relevant | **Move or duplicate** core apply tests; keep CLI integration keepers |
| `BuildPolicyInputs` | `policy.go` L15–72 | Live queries wired (blocking, debt, regression, plan) | **No change** — already live post-S03/S05 |
| Schema / compat | `internal/store/schema/` | **20** files, max **020** (S04 promotion); compat ceiling **20** | **No new migration**; compat **20** unchanged |
| Replay semantics | `apply.go` L378–386, L647–664 | Duplicate `apply_id` returns prior counts; skips all writes | **Unchanged** — replay detection stays **outside** tx |

### Non-transactional gap (D-08 evidence)

Today `Apply` runs: discoveries → plan_changes → spawned_tasks → cognitive writes → `ApplyDeliberationTransition` → `loop.step` event. Each store call auto-commits. Example failure after partial success: valid **discovery** inserted, then **plan_change** with `replan` missing `discovery_id` fails at `apply.go` L438–439 — discovery row **remains** (no rollback). S06 closes this.

### goal_id gap (D-13 evidence)

`Apply` validates seed goal at envelope entry. `ApplyDeliberationTransition` accepts any non-empty `goalID` even when it differs from `task.GoalID` — deliberation state could be keyed under wrong goal. S06 adds domain guard; loop path already pre-checks seed.

### Test location gap (D-15 evidence)

P20 verify ran `go test ./internal/loop/...` — vacuous pass (no tests). All behavioral apply proofs live in `cmd/trace/loop_test.go`. S06 moves minimum floor to package under test.

## W-09 / W-10 / W-11 locks

| Work ID | TRACE § | Lock |
|---------|---------|------|
| **W-09** | §29O fail-closed apply | All mutating apply steps in **one** SQLite transaction; any error → full rollback; no partial artifact rows |
| **W-10** | §6 goal binding | `seed.goal_id` must equal task's stored `goal_id` at loop **and** domain transition boundaries |
| **W-11** | §29Q verify floor | **≥8** named apply tests in `internal/loop/*_test.go`; cmd/trace retains thin CLI integration keepers only |

## FINAL locked defaults (S06-01 must not re-debate)

| Item | Value |
|------|-------|
| Transaction scope | Inside one tx: P19 writes (discoveries/plan_changes/spawned_tasks) + cognitive writes + `ApplyDeliberationTransition` (state + event) + `loop.step.applied` event |
| Pre-tx (read-only) | `ParseApplyEnvelope` / `ValidateApplyEnvelope`, task load, seed goal mismatch check, replay short-circuit (`findLoopStep`) |
| Store API | Add `store.WithTx(fn func(*Store) error) error` — `Begin`, callback, `Commit` on nil / `Rollback` on error; bind callback store to tx (precedent: `fts.go` L138–233) |
| Domain service | Construct `domain.New(txStore)` inside callback — domain uses same tx-bound store |
| Planner service | Pass through to replan inside tx when `plan_changes[].replan` present |
| goal_id — loop | Keep existing `apply.go` L375–377 check; error text contains `seed goal mismatch` |
| goal_id — domain | In `ApplyDeliberationTransition`: after `GetTask`, require `task.GoalID != nil && *task.GoalID == goalID`; else `ErrValidation` with `goal_id does not match task` |
| Replay | Unchanged — duplicate `apply_id` returns cached counts; **no** tx opened |
| Fail-closed | Validation errors before tx; domain/store errors inside tx roll back all prior writes in that apply |
| Schema / compat | Max mig **020**; compat ceiling **20** — **no mig 021** |
| MCP / CLI schema | `trace.loop.apply.v1` string unchanged |
| Test migration | Move **8** core tests to `internal/loop/apply_test.go`; delete or thin-wrap duplicates in `cmd/trace/loop_test.go` per D-15 |

### Transaction boundary sketch (FINAL)

```text
Apply(ctx, st, plan, env):
  validate envelope + load task + seed goal match     // no tx
  if replay(apply_id): return cached result           // no tx
  return st.WithTx(func(txSt *Store) error {
    dom := domain.New(txSt)
    // discoveries, plan_changes, spawned_tasks
    applyCognitiveWrites(ctx, dom, txSt, ...)
    BuildPolicyInputs(ctx, dom, plan, ...)
    dom.ApplyDeliberationTransition(...)
    txSt.AppendEvent(loop.step.applied)
    return nil
  })
```

### Rollback test scenario (FINAL)

`TestLoopApplyTransactionalRollbackOnFailure` envelope (passes `ValidateApplyEnvelope`):

1. One valid `discoveries[]` row (will insert if tx absent).
2. One `plan_changes[]` with `replan` set but **empty** `discovery_id` — fails at apply time after discovery loop.

**Assert:** after error, `ListDiscoveries()` length **0** and no `deliberation.transition` / `loop.step.applied` events.

## Named tests (S06-01)

| # | Test | Location | Proves |
|---|------|----------|--------|
| 1 | `TestLoopApplyTransactionalRollbackOnFailure` | `internal/loop/apply_test.go` | Mid-apply error → zero artifact rows |
| 2 | `TestLoopApplyGoalIDMismatchFailsClosed` | `internal/loop/apply_test.go` | Wrong `seed.goal_id` rejected; no writes |
| 3 | `TestLoopApplyDeliberationTransitionEvent` | `internal/loop/apply_test.go` | Success path persists `deliberation.transition` with required payload keys |
| 4 | `TestLoopApplyNoPartialWritesOnValidationFailure` | `internal/loop/apply_test.go` | Pre-tx validation fail → clean DB (move from cmd/trace) |
| 5 | `TestLoopApplyReplaySkipsDuplicateTransition` | `internal/loop/apply_test.go` | Replay idempotent — single transition event (move from cmd/trace) |
| 6 | `TestLoopApplyUnknownWriteKeyFailsClosed` | `internal/loop/apply_test.go` | Unknown `writes` key rejected at parse (move from cmd/trace) |
| 7 | `TestValidateApplyEnvelopeSpawnedTaskGoalMismatch` | `internal/loop/apply_test.go` | Spawned task `goal_id` ≠ seed fails validation |
| 8 | `TestLoopApplySuccessPersistsLoopStepEvent` | `internal/loop/apply_test.go` | `loop.step.applied` event on happy path |
| — | `TestApplyDeliberationTransitionRequiresMatchingGoalID` | `internal/domain/deliberation_test.go` | Domain goal guard independent of loop |
| — | `TestLoopApplyUncertaintyWriteAffectsNextSelectNext` | `cmd/trace/loop_test.go` | CLI integration keeper |
| — | `TestLoopApplyMalformedInputFailsClosed` | `cmd/trace/loop_test.go` | P19 CLI keeper |
| — | `TestLoopApplyReplayAndStatusFlow` | `cmd/trace/loop_test.go` | P19 CLI end-to-end keeper |

**Keep green (no regressions):** P20 apply keepers, S05 `TestWhyTaskIncludesDeliberationTransition`, compat `TestCompatibilitySecurityChecklist`.

## Touch files

- `internal/store/open.go` (or new `tx.go`) — `WithTx`
- `internal/loop/apply.go` — transaction wrapper
- `internal/loop/apply_test.go` — **new**; ≥8 tests + shared helpers (`openLoopTestStore`, seed task/goal/plan fixtures)
- `internal/domain/deliberation.go` — goal_id match guard
- `internal/domain/deliberation_test.go` — `TestApplyDeliberationTransitionRequiresMatchingGoalID`
- `cmd/trace/loop_test.go` — remove moved tests or replace with `t.Skip` + re-export comment; keep CLI keepers

## Planner work

1. [x] Live inventory `apply.go` / `deliberation.go` / store tx patterns / schema max **020** / compat **20**.
2. [x] Lock **W-09** (transactional apply — single commit, rollback on any error).
3. [x] Lock **W-10** (goal_id validation at loop + domain transition).
4. [x] Lock **W-11** (≥8 tests in `internal/loop`; D-15 cmd/trace trim policy).
5. [x] Thicken `01-apply-hardening.md` + `02-scope-review.md` with before-state, 8+8 tests, keeper floor.
6. [x] Update `SCOPE-TODOS.md`.

## Exit criteria

- [x] Transaction + validation + test floor locked
- [x] 8+ named internal/loop tests locked
- [x] W-09/W-10/W-11 explicit
- [x] 01/02 thickened enough to implement alone
- [x] No product Go
- [x] No mig 021

## Next

**P21-S06-01**
