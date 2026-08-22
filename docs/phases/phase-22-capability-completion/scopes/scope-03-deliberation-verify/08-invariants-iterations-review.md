# P22-S03-08 — Review: invariants + S03 close

## Metadata
- id: P22-S03-08
- todo_ids: [P22-S03-08]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C14, C15** and that all S03-owned capabilities still hold: **C09, C11, C12, C13, C14, C15, C36, C38-CLI**.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## S03 capability matrix

| Cap | Row | Evidence |
|-----|-----|----------|
| C09 | S03-01/02 | Live cycle flags in `BuildPolicyInputs` |
| C11 | S03-05/06 | Status cycle gate + deliberation cannot skip TEST/VERIFY/EVALUATE |
| C12 | S03-03/04 | `trace test run` + stored outcomes |
| C13 | S03-05/06 | Test regression vs prior pass |
| C14 | S03-07 | Invariant check on change paths |
| C15 | S03-07 | Iteration outcome compare |
| C36 | S03-05/06 | `CoordinateVerification` order |
| C38-CLI | S03-03/05/07 | test/verify/outcomes CLI (MCP half deferred S08) |

## Re-run keeper floor

```bash
go test ./internal/loop/... ./internal/deliberation/... ./internal/domain/... ./internal/testrun/... -count=1 -run 'TestBuildPolicyInputs|TestSelectNextNeverExecuteOnBlockingUncertainty|TestTestRun|TestVerificationCycle|TestCoordinateVerification|TestInvariant|TestCompareIterationOutcomes|TestTestPassAloneCannotSatisfyVerificationGate|TestPromotionGateRequiresStoredTestNotAgentClaim'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestLoopNextExecute|TestTestRun|TestVerifyRun|TestVerifyInvariants|TestOutcomesCompare'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered
go test ./internal/retrieval/... -count=1 -run TestImpactWalkIncludesAffectedTests
```

Checklist: box **C09, C11–C15** when all true (C38-CLI test/verify only — note MCP gap for S08).

## Spawn policy

If any owned capability unmet: spawn **`P22-S03-08a` + `P22-S03-08b`**. Do not close with residuals.

## Exit criteria

- [ ] All S03 capabilities closed or spawned
- [ ] Confidence **high**
- [ ] Board Notes: **S03 complete** → next `P22-S04-00`
