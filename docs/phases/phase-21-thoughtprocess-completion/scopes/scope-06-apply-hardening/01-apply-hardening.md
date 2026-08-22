# P21-S06-01 — Implement: apply hardening

## Metadata
- id: P21-S06-01
- todo_ids: [P21-S06-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective
Transactional `loop apply`, goal_id validation on deliberation transition, and internal/loop test floor (D-15). **No schema migration.**

## Session start
Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: **status + notes only**.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- [DECISION-LOG.md](../../DECISION-LOG.md) D-08, D-13, D-15
- [WORK-MAP.md](../../WORK-MAP.md) W-09, W-10, W-11
- Live: `internal/loop/{apply.go,apply_writes.go}`, `internal/domain/deliberation.go`, `internal/store/open.go`, `cmd/trace/loop_test.go`

## Locked defaults (from S06-00 — do not re-debate)

| Item | Value |
|------|-------|
| Transaction | All apply writes + deliberation transition + `loop.step.applied` in **one** DB transaction |
| Pre-tx | Envelope validation, task load, seed goal mismatch, replay short-circuit |
| Store API | `store.WithTx(fn func(*Store) error) error` |
| goal_id — loop | Keep `apply.go` seed mismatch check |
| goal_id — domain | `ApplyDeliberationTransition` rejects goal ≠ task.GoalID |
| Test floor | **≥8** tests in `internal/loop/apply_test.go` |
| cmd/trace | Keep CLI integration keepers; trim moved duplicates |
| Schema / compat | Max mig **020**; ceiling **20** — no **021+** |
| Replay | Unchanged semantics |

## Live inventory (before — confirmed S06-00)

| Surface | Location | Today |
|---------|----------|-------|
| Apply writes | `apply.go` L388–516 | Sequential; partial rows on mid-failure |
| Seed goal check | `apply.go` L375–377 | Present |
| Transition goal check | `deliberation.go` L27–29 | Task existence only |
| Store tx helper | — | **Missing** — ad-hoc tx in fts/vcs only |
| Loop tests | `internal/loop/` | **0** test files |
| CLI apply tests | `cmd/trace/loop_test.go` | 14+ apply-related funcs at L191–1079 |
| Schema | 20 files, max **020** | S04 promotion landed |
| Compat | `evals/compat` | Ceiling **20** |

## Requirements

### 1. Store — `WithTx` (W-09)

Add to `internal/store/` (prefer new `tx.go` or extend `open.go`):

```go
func (s *Store) WithTx(fn func(*Store) error) error
```

- `Begin()` on `s.db`; defer `Rollback()` (no-op after commit).
- Callback receives `&Store{db: txAsQueryable, projectRoot, projectID}` — use pattern that lets existing store methods run on tx (see `fts.go` L138–233 for precedent).
- `Commit()` only when callback returns nil.
- Propagate callback error after rollback.

**Do not** refactor every store method to accept `*sql.Tx` — minimal surface for loop apply only.

### 2. Loop — wrap `Apply` in transaction (W-09)

Refactor `Apply` in `internal/loop/apply.go`:

1. Keep **outside** tx: nil store check, task load, seed goal mismatch (L375–377), replay return (L378–386).
2. Move **inside** `st.WithTx`:
   - All discovery/plan_change/spawned_task loops
   - `applyCognitiveWrites`
   - `BuildPolicyInputs` + `ApplyDeliberationTransition`
   - `AppendEvent` for `EventLoopStep`
3. Return `ApplyResult` from tx callback via closure or named result struct.
4. On tx error, return zero `ApplyResult` and wrapped error — **no** partial state.

### 3. Domain — goal_id guard (W-10)

In `ApplyDeliberationTransition` (`deliberation.go`), after `GetTask(taskID)`:

```go
task, err := s.store.GetTask(taskID)
// ...
if task.GoalID == nil || strings.TrimSpace(*task.GoalID) == "" {
    return ..., &ErrValidation{Msg: "task has no goal_id"}
}
if strings.TrimSpace(*task.GoalID) != goalID {
    return ..., &ErrValidation{Msg: "goal_id does not match task"}
}
```

Add `TestApplyDeliberationTransitionRequiresMatchingGoalID` in `deliberation_test.go` — create task under goal A, call transition with goal B → error, no event.

### 4. Tests — internal/loop floor (W-11 / D-15)

Create `internal/loop/apply_test.go` with package `loop_test` (or `loop` + helpers). Minimum **8** tests:

| Test | Implementation notes |
|------|---------------------|
| `TestLoopApplyTransactionalRollbackOnFailure` | Discovery + plan_change replan w/o discovery_id; assert 0 discoveries, 0 transitions |
| `TestLoopApplyGoalIDMismatchFailsClosed` | Task under goal G1; envelope seed.goal_id = G2; error contains `seed goal mismatch`; 0 writes |
| `TestLoopApplyDeliberationTransitionEvent` | Port logic from `cmd/trace/loop_test.go` L875–921 using `loop.Apply` directly |
| `TestLoopApplyNoPartialWritesOnValidationFailure` | Port L1047–1079 — invalid uncertainty title in envelope |
| `TestLoopApplyReplaySkipsDuplicateTransition` | Port L923–983 |
| `TestLoopApplyUnknownWriteKeyFailsClosed` | Port L758–771 |
| `TestValidateApplyEnvelopeSpawnedTaskGoalMismatch` | Spawned task goal_id ≠ seed — validation error |
| `TestLoopApplySuccessPersistsLoopStepEvent` | Happy-path apply; assert `loop.step.applied` event exists |

Shared helpers (in test file): temp store via `store.Open`, create goal/task/plan via domain, build `ApplyEnvelope` structs.

**cmd/trace trim:** Remove duplicated test bodies moved to internal/loop; optionally leave one-liner smoke that runs CLI `loop apply` for regression. **Keep** unchanged:

- `TestLoopApplyMalformedInputFailsClosed`
- `TestLoopApplyUncertaintyWriteAffectsNextSelectNext`
- `TestLoopApplyRegressionWriteAffectsPolicyInputs`
- `TestLoopApplyReplayAndStatusFlow`
- `TestLoopNextPacketShape`

## Implementation order

1. Add `TestApplyDeliberationTransitionRequiresMatchingGoalID` (domain guard — can land before tx).
2. Add `store.WithTx` + unit smoke in store if needed.
3. Wrap `Apply` in tx; run rollback test (should fail before tx, pass after).
4. Port/move 8 tests to `internal/loop/apply_test.go`.
5. Trim `cmd/trace/loop_test.go` duplicates.
6. Run keeper floor below.

## Keeper command floor

```bash
go test ./internal/loop/... -count=1 -run 'TestLoopApplyTransactionalRollbackOnFailure|TestLoopApplyGoalIDMismatchFailsClosed|TestLoopApplyDeliberationTransitionEvent|TestLoopApplyNoPartialWritesOnValidationFailure|TestLoopApplyReplaySkipsDuplicateTransition|TestLoopApplyUnknownWriteKeyFailsClosed|TestValidateApplyEnvelopeSpawnedTaskGoalMismatch|TestLoopApplySuccessPersistsLoopStepEvent'
go test ./internal/domain/... -count=1 -run 'TestApplyDeliberationTransitionRequiresMatchingGoalID|TestApplyDeliberationTransitionRequiresIDs|TestApplyDeliberationTransitionPersistsEvent'
go test ./cmd/trace -count=1 -run 'TestLoopApplyUncertaintyWriteAffectsNextSelectNext|TestLoopApplyMalformedInputFailsClosed|TestLoopNextPacketShape|TestLoopApplyReplayAndStatusFlow|TestWhyTaskIncludesDeliberationTransition'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [ ] `store.WithTx` exists; `Apply` uses single transaction for all mutating steps
- [ ] `ApplyDeliberationTransition` rejects mismatched goal_id
- [ ] **≥8** tests in `./internal/loop/...` PASS
- [ ] P19/P20 CLI keepers PASS
- [ ] Compat ceiling **20** unchanged
- [ ] No mig **021+**

## Next

**P21-S06-02**
