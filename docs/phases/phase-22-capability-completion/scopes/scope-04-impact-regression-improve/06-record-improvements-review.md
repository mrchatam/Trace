# P22-S04-06 — Review: improvements + S04 close

## Metadata
- id: P22-S04-06
- todo_ids: [P22-S04-06]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C18** and S04-owned **C08, C16** all closed. Scope exit gate before S05.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## S04 close checklist

| Capability | Pass condition |
|------------|----------------|
| **C08** | Stored `impact_predictions`; compare persisted; CLI predict/compare; S04-02 boxed |
| **C16** | `regression_associated_change` links; list by change; Law 5 hold; S04-04 boxed |
| **C18** | `improvements` CRUD + CLI + seed export; S04-06 boxes C18 |
| Schema | Exactly **24** sql files; **no 025+** |
| Compat | Ceiling **24** |
| MCP | Catalog **10** unchanged |
| Cross-scope | S01 `TestImpactWalkIncludesAffectedTests`; S03 `TestRegressionDetectedVsPriorPassingTest`; regression keepers PASS |

## Spawn policy

If **any** of C08/C16/C18 unmet: spawn **`P22-S04-06a` + `P22-S04-06b`**. Do not mark S04 complete with residuals.

## Full S04 re-run

```bash
go test ./internal/domain/... ./internal/store/... ./internal/retrieval/... -count=1 -run 'TestRecordPredictedImpact|TestImpactCompare|TestRegression|TestSetAttributionCaused|TestRecordImprovement|TestSeedExportIncludesImprovements'
go test ./internal/domain/... -count=30 -run TestRegressionDetectedVsPriorPassingTest
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestImpact|TestOutcomes'
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered
ls internal/store/schema/*.sql | wc -l
```

## Exit criteria

- [ ] C08, C16, C18 closed (checklist all three `[x]`) or spawned
- [ ] Confidence **high** \| **medium** \| **low**
- [ ] Board Notes: **S04 complete** → next **`P22-S05-00`**
