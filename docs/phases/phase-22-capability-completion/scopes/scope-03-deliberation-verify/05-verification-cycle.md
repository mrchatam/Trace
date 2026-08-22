# P22-S03-05 — Implement: verification cycle + scoring coordination

## Metadata
- id: P22-S03-05
- todo_ids: [P22-S03-05]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

**Require** a test/verification cycle (**C11**), detect regression vs prior passing tests (**C13**), and **coordinate** testing, verification, and scoring (**C36**). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- D-22-20: verification cycle is deliberation/status — **do not** auto-`TransitionTask` DONE
- S03-03: `internal/testrun` (coordinator may call it)
- Domain gates: `CheckVerificationGate`, `RecordVerificationOutcome`, `RecordEvaluationOutcome`, `CompareScoresToBaseline`
- Loop: `internal/loop/apply.go` `Status()`, `StatusResult`, `statusBlocked`

## Locked defaults

| Item | Value |
|------|-------|
| DONE policy | **Unchanged** — operator + Review PASS only |
| Status schema | Add to `trace.loop.status.v1`: `verification_cycle: { execute_pending, test_pending, verification_incomplete, evaluation_pending, reflect_pending, incomplete_reason }` — mirror `PolicyInputs` + human reason string |
| Cycle incomplete | `incomplete_reason` when any of test/verification/evaluation pending flags true OR `verification_debt.present` |
| Coordinator | `domain.CoordinateVerification(ctx, taskID, CoordinateOptions)` — order: (1) testrun if test pending, (2) if tests failed → stop unless `ForceEval`, (3) record/find verification with goal+evidence if debt, (4) evaluation vs **active baseline** for task goal commit if eval pending |
| Regression (C13) | `DetectTestRegression(taskID, testName)` — compare latest two `kind=test` rows for same `test_name`: prior `pass` + current `fail` → regression signal (return struct; may set open regression or `ReplanNeeded` advisory in status only — **no** silent DONE) |
| CLI | `trace verify run --task <id> [--force-eval]` → coordinator JSON on stdout |
| Schema / compat | **23** — no migration |

## Requirements

1. Status reports cycle incompleteness (JSON block, not docs-only).
2. Coordinator enforces order; default **fail-closed** skip evaluation when latest test run failed; `--force-eval` overrides.
3. C13 regression helper + wire into status or coordinator result payload (`regression_detected: bool`, `test_name`).
4. `SelectNext` + `statusBlocked` behavior unchanged except new visibility fields.
5. Named tests + keepers.

## Touch files

- `internal/loop/apply.go` — extend `StatusResult`
- `internal/domain/coordinate.go` or `outcomes.go` — coordinator + regression helper (**new** file OK)
- `internal/testrun` — import from coordinator
- `cmd/trace/verify.go`, `verify_test.go` (**new**)
- `cmd/trace/root.go`, `help.go`, `capability.go`
- `internal/loop/*_test.go`

## Named tests

| Test | Proves |
|------|--------|
| `TestVerificationCycleBlocksSkipInStatus` | status `verification_cycle.incomplete_reason` when no test row |
| `TestCoordinateVerificationOrder` | mock/stub: test → verify → eval call order |
| `TestRegressionDetectedVsPriorPassingTest` | pass then fail → regression true |
| `TestTestPassAloneCannotSatisfyVerificationGate` | keeper |
| `TestHasVerificationDebt` / `TestVerificationDebtWhenImplementationWithoutVerification` | keepers |

```bash
go test ./internal/loop/... ./internal/domain/... -count=1 -run 'TestVerificationCycle|TestCoordinateVerification|TestRegressionDetectedVsPriorPassingTest|TestTestPassAloneCannotSatisfyVerificationGate|TestVerificationDebt'
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestLoopStatus|TestVerifyRun'
```

## Exit criteria

- [ ] C11, C13, C36 true
- [ ] Named tests PASS; compat **23**
- [ ] Checklist boxes **not** checked until S03-06 review
- [ ] Board Notes

## Minimal todos

- [ ] Status `verification_cycle` block
- [ ] Coordinator + regression helper
- [ ] CLI `verify run`
- [ ] Tests + board notes
