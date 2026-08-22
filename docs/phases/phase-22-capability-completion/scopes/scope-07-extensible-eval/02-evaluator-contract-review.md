# P22-S07-02 — Review: evaluator contract

## Metadata
- id: P22-S07-02
- todo_ids: [P22-S07-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C40, C43** — multiple mechanisms via additive contract; core outcomes/changes model unchanged.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md), [01-evaluator-contract.md](01-evaluator-contract.md)
- S03 keepers: `TestCoordinateVerification`, `TestBuildPolicyInputs`, `TestInvariant`

## Review checklist

### C40 — multiple verification mechanisms

- [ ] `DefaultRegistry` lists ≥ **4** built-in ids: `stored_test`, `stored_verification`, `stored_evaluation`, `architectural_invariant`
- [ ] `TestEvalRegistryMultipleMechanisms` PASS
- [ ] `RunAll` executes filtered subset; stable sort by id
- [ ] Each built-in delegates to existing domain API (grep: no duplicate invariant SQL in eval package)

### C43 — additive without core redesign

- [ ] `TestAddMechanismWithoutSchemaChange` PASS — fake mechanism registers + runs
- [ ] Schema exactly **26** sql files; **`026_eval_rules.sql`** only new mig (no 027+)
- [ ] **Zero** ALTER on `outcome_results`; kind CHECK still `test|verification|evaluation` only
- [ ] `CoordinateVerification` / `RecordEvaluationOutcome` **unchanged** this row (grep diff scope)
- [ ] Compat/embed ceiling **26** PASS

### Landmines

- [ ] **Import cycle**: `internal/eval` imports `domain`; `domain` does **not** import `eval`
- [ ] **No daemon/subprocess** in eval package (Law 2 boundary)
- [ ] MCP catalog **13** (`TestToolNamesRegistered`)
- [ ] S03 keeper `TestRegressionDetectedVsPriorPassingTest` **30/30** stable

## Spawn policy

If unmet: spawn **`P22-S07-02a`** (fix) + **`P22-S07-02b`** (re-review + box). Do not close with residuals.

## Re-run commands

```bash
go test ./internal/eval/... ./internal/store/... -count=1 -run 'TestEvalRegistry|TestAddMechanism|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/domain/... -count=30 -run TestRegressionDetectedVsPriorPassingTest
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered
ls internal/store/schema/*.sql | wc -l
```

## Exit criteria

- [ ] C40, C43 closed or spawned
- [ ] Confidence **high**
- [ ] Checklist C40 + C43 **boxed** when closed
