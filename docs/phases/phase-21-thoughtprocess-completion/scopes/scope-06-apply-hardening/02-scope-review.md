# P21-S06-02 — Review: apply hardening

## Metadata
- id: P21-S06-02
- todo_ids: [P21-S06-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective
Independent review: transactional apply, goal_id guards, no partial writes, test location fixed (D-15). Schema/compat unchanged at **20**.

## Session start
**Fresh subagent** (not S06-01 session). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: status + notes only; spawn Na/Nb if gap found.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-apply-hardening.md](01-apply-hardening.md) — implementer deliverable
- [DECISION-LOG.md](../../DECISION-LOG.md) D-08, D-13, D-15
- [WORK-MAP.md](../../WORK-MAP.md) W-09, W-10, W-11

## Review checklist

| Check | Evidence |
|-------|----------|
| `store.WithTx` | Exists; Begin/Commit/Rollback; callback error → rollback |
| Tx scope | Discoveries through `loop.step.applied` inside one tx |
| Pre-tx boundary | Validation, seed goal check, replay **outside** tx |
| Rollback proof | `TestLoopApplyTransactionalRollbackOnFailure` green — 0 rows after mid-failure |
| Loop goal guard | `apply.go` seed mismatch before tx; error text inspectable |
| Domain goal guard | `ApplyDeliberationTransition` rejects goal ≠ task.GoalID |
| No partial validation writes | `TestLoopApplyNoPartialWritesOnValidationFailure` green |
| Replay unchanged | `TestLoopApplyReplaySkipsDuplicateTransition` — single transition |
| internal/loop test count | **≥8** named apply tests in `internal/loop/*_test.go` |
| D-15 hygiene | Moved tests not duplicated verbatim in cmd/trace (or documented thin CLI smoke) |
| Schema | **20** files; max **020**; **no 021+** |
| Compat ceiling **20** | `evals/compat` + `TestMigrationStatusReportsEmbedMax` green |
| Apply schema | `trace.loop.apply.v1` unchanged |
| P19 keepers | `TestLoopApplyMalformedInputFailsClosed`, `TestLoopApplyReplayAndStatusFlow` green |
| P20 keepers | `TestLoopApplyUncertaintyWriteAffectsNextSelectNext`, `TestLoopApplyRegressionWriteAffectsPolicyInputs` green |
| S05 keeper | `TestWhyTaskIncludesDeliberationTransition` green |

## D-08 / D-13 / D-15 closure

- **D-08 promote:** Transactional apply — grep `WithTx` in `Apply`; rollback test proves no partial discoveries.
- **D-13 promote:** Domain + loop goal_id validation — `TestApplyDeliberationTransitionRequiresMatchingGoalID` + `TestLoopApplyGoalIDMismatchFailsClosed`.
- **D-15 promote:** `./internal/loop/...` non-vacuous — count ≥8 `TestLoopApply*` / `TestValidateApply*` in package.

## Keeper command floor

```bash
go test ./internal/loop/... -count=1 -run 'TestLoopApplyTransactionalRollbackOnFailure|TestLoopApplyGoalIDMismatchFailsClosed|TestLoopApplyDeliberationTransitionEvent|TestLoopApplyNoPartialWritesOnValidationFailure|TestLoopApplyReplaySkipsDuplicateTransition|TestLoopApplyUnknownWriteKeyFailsClosed|TestValidateApplyEnvelopeSpawnedTaskGoalMismatch|TestLoopApplySuccessPersistsLoopStepEvent'
go test ./internal/domain/... -count=1 -run 'TestApplyDeliberationTransitionRequiresMatchingGoalID|TestApplyDeliberationTransitionRequiresIDs|TestApplyDeliberationTransitionPersistsEvent'
go test ./cmd/trace -count=1 -run 'TestLoopApplyUncertaintyWriteAffectsNextSelectNext|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopNextPacketShape|TestWhyTaskIncludesDeliberationTransition|TestLoopApplyRegressionWriteAffectsPolicyInputs'
go test ./internal/store/... -count=1 -run TestMigrationStatusReportsEmbedMax
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

Document internal/loop test count in board Notes (e.g. `go test ./internal/loop/... -list '.*' | wc -l` or manual count).

## Review focus

- **Tx leak:** Failed apply must not leave deliberation state advanced without matching writes (or vice versa).
- **Replay vs tx:** Replay path must not open tx or double-append events.
- **Goal binding:** Wrong goal rejected at loop envelope **and** direct domain transition call.
- **Blast radius:** No changes to `select.go`, MCP catalog, or seed export shape.
- **Test duplication:** cmd/trace should not maintain full copies of 8 moved tests — D-15 intent is package-local proof.

## Spawn policy

- **Na (implement):** No tx wrapper, rollback test fails, goal guard missing, <8 internal/loop tests, partial writes on failure, compat ≠ 20
- **Nb (review):** re-review after Na
- Do **not** spawn for optional dedicated store `WithTx` unit test if loop rollback test covers behavior

## Exit criteria

- [ ] No blocker/high without spawn or inline fix
- [ ] Confidence **high** with test output pasted in board Notes
- [ ] D-08 + D-13 + D-15 closure evidenced
- [ ] internal/loop test count ≥8 documented in Notes
- [ ] Spawn Na/Nb only if apply hardening gap found

## Next

**P21-S07-00** (unless Na spawned)
