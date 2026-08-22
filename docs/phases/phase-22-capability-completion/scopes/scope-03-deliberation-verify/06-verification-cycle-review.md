# P22-S03-06 — Review: verification cycle

## Metadata
- id: P22-S03-06
- todo_ids: [P22-S03-06]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C11, C13, C36** are real gates/coordination, not documentation. DONE / Review PASS policy unchanged (D-22-20).

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

1. `trace loop status` JSON includes `verification_cycle` with honest flags (grep struct + CLI test).
2. `CoordinateVerification` exists and orders test → verification → evaluation (read call graph or unit test).
3. C13: regression detection compares **two stored test outcomes**, not agent claims.
4. Grep `TransitionTask` / `work_state` — coordinator must **not** auto-DONE on cycle complete.
5. `TestTestPassAloneCannotSatisfyVerificationGate` still PASS.
6. C09/C12 hold from prior rows (spot-check `BuildPolicyInputs` + testrun).
7. Schema **23**; compat PASS.

## Spawn policy

If any of C11/C13/C36 unmet: spawn **`P22-S03-06a` + `P22-S03-06b`**. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/loop/... ./internal/domain/... -count=1 -run 'TestVerificationCycle|TestCoordinateVerification|TestRegressionDetectedVsPriorPassingTest|TestTestPassAloneCannotSatisfyVerificationGate|TestBuildPolicyInputs'
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestLoopStatus|TestVerifyRun'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [x] C11, C13, C36 closed or spawned
- [ ] Confidence **high** (medium — C13 blocker spawned 06a/06b)
- [x] Board Notes
