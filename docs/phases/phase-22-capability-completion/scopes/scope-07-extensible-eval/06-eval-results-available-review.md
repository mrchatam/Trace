# P22-S07-06 — Review: eval results + S07 close

## Metadata
- id: P22-S07-06
- todo_ids: [P22-S07-06]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C42-library** and S07-owned **C40, C41, C43** closed. C42 **surface** was S05 — hold that row.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md)
- Full S07 implement rows 01, 03, 05
- S05-06 close notes for C42-surface

## Review checklist

### C42 — evaluation results for future agents

- [ ] **Surface (S05 hold):** `TestContextIncludesEvaluations` PASS; packet cap 8; task-scoped filter
- [ ] **Library (this scope):** `TestListEvaluationResultsForFutureAgents` PASS
- [ ] `TestEvalResultsIncludeMechanismID` PASS — mechanism_id on all rows
- [ ] Evaluation rows include **`comparison_json`** (library not truncated)
- [ ] **Box checklist C42 line** when both halves verified

### C40 / C41 / C43 holds

- [ ] S07-02 closed — registry + four built-ins + fake mech extensibility
- [ ] S07-04 closed — rules file + invariant override + seed pointer
- [ ] Outcome kind enum unchanged; schema **26**; compat **26**
- [ ] No daemon/subprocess in eval package

### S07 complete spot-check

- [ ] S01 keeper: `TestImpactWalkIncludesAffectedTests` PASS
- [ ] S03 keeper: `TestRegressionDetectedVsPriorPassingTest` **30/30**
- [ ] S06 keeper: `TestQuerySimilarChanges` PASS
- [ ] MCP catalog **13**
- [ ] grep: no **027+** migrations

## Spawn policy

If unmet: spawn **`P22-S07-06a`** + **`P22-S07-06b`**. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/eval/... -count=1 -run 'TestEvalRegistry|TestAddMechanism|TestProjectEvalRules|TestListEvaluationResults|TestEvalResultsIncludeMechanismID'
go test ./internal/compiler/... ./internal/loop/... -count=1 -run 'TestContextIncludesEvaluations|TestLoopNextPlanningEvidence'
go test ./internal/domain/... -count=30 -run TestRegressionDetectedVsPriorPassingTest
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestEvalRules|TestEvalResults|TestContext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered
ls internal/store/schema/*.sql | wc -l
```

## Exit criteria

- [ ] C40–C43 closed or spawned (C42 surface + library both `[x]`)
- [ ] **S07 complete** — next scope planner **P22-S08-00**
- [ ] Confidence **high**
