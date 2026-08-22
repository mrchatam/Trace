# P22-S03-06b — Review: C13 regression ordering fix

## Metadata
- id: P22-S03-06b
- todo_ids: [P22-S03-06b]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C13** closed after **P22-S03-06a**. Re-verify **C11** and **C36** still hold. Box checklist lines when all three true.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

1. `TestRegressionDetectedVsPriorPassingTest` PASS `-count=30` (non-flaky).
2. Same-second / tie-break test exists and PASS.
3. `DetectTestRegression` compares two **stored** rows — not agent claims.
4. `CoordinateVerification` order unchanged; no `TransitionTask` auto-DONE.
5. `verification_cycle` status block still honest (`TestVerificationCycleBlocksSkipInStatus`).
6. Schema **23**; compat PASS.

## Re-run commands

```bash
go test ./internal/domain/... -count=30 -run TestRegressionDetectedVsPriorPassingTest
go test ./internal/loop/... ./internal/domain/... -count=1 -run 'TestVerificationCycle|TestCoordinateVerification|TestRegressionDetected|TestTestPassAloneCannotSatisfyVerificationGate|TestBuildPolicyInputs'
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestLoopStatus|TestVerifyRun'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Spawn policy

If C13 still flaky or C11/C36 regressed: spawn **`P22-S03-06c` + `P22-S03-06d`**.

## Exit criteria

- [ ] C11, C13, C36 closed (checklist boxed)
- [ ] Confidence **high**
- [ ] Board Notes
