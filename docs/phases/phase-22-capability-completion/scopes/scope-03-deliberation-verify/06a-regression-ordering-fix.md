# P22-S03-06a — Implement: C13 regression ordering fix

## Metadata
- id: P22-S03-06a
- todo_ids: [P22-S03-06a]
- role: implementer
- skills: [incremental-implementation, test-driven-development, debugging-and-error-recovery]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Fix **C13** regression detection so `DetectTestRegression` reliably compares the **last two stored** `kind=test` rows per `test_name` in **insertion/chronological order**. Spawned from **P22-S03-06** review: keeper `TestRegressionDetectedVsPriorPassingTest` fails ~60% of runs (`-count=1` loop).

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [06-verification-cycle-review.md](06-verification-cycle-review.md) — spawn rationale
- [05-verification-cycle.md](05-verification-cycle.md) — original deliverable
- `internal/domain/coordinate.go` — `DetectTestRegression`
- `internal/store/outcomes.go` — `ListOutcomeResultsByTaskKind`, `nowRFC3339`

## Root cause (locked)

`created_at` is **second-precision** RFC3339. Two rapid `RecordTestOutcome` calls share a timestamp; tie-break uses **UUID string order**, not insertion order. `DetectTestRegression` may treat `fail→pass` as prior/current → `Detected=false`.

## Locked defaults

| Item | Value |
|------|-------|
| Scope | Ordering fix only — no behavior change to C11 status block or C36 coordinator order |
| Preferred fix | Preserve true chronological order for same-`test_name` rows: e.g. `ORDER BY created_at ASC, rowid ASC` in store list query **and** drop redundant re-sort in `DetectTestRegression`, **or** sub-ms timestamps on insert — pick minimal diff |
| Schema / compat | **23** — no migration unless unavoidable |
| DONE policy | Unchanged (D-22-20) |

## Requirements

1. `TestRegressionDetectedVsPriorPassingTest` PASS `-count=30` (non-flaky).
2. `DetectAnyTestRegression` inherits fix (calls `DetectTestRegression`).
3. No `TransitionTask` / auto-DONE in coordinator paths.

## Named tests

| Test | Proves |
|------|--------|
| `TestRegressionDetectedVsPriorPassingTest` | pass then fail → regression true (stable) |
| `TestRegressionDetectedSameSecondTimestamps` | **new** — force equal `created_at`, assert pass→fail detected |
| Keepers | `TestCoordinateVerification*`, `TestVerificationCycle*` unchanged PASS |

```bash
go test ./internal/domain/... -count=30 -run TestRegressionDetectedVsPriorPassingTest
go test ./internal/loop/... ./internal/domain/... -count=1 -run 'TestVerificationCycle|TestCoordinateVerification|TestRegressionDetected|TestTestPassAloneCannotSatisfyVerificationGate'
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestLoopStatus|TestVerifyRun'
```

## Exit criteria

- [ ] C13 ordering fix landed
- [ ] Keeper non-flaky (`-count=30`)
- [ ] Board status + notes only
- [ ] Checklist **not** boxed (S03-06b closes)

## Minimal todos

- [ ] Fix store ordering or timestamp precision
- [ ] Simplify `DetectTestRegression` pairing
- [ ] Add same-second regression test
- [ ] Board notes
