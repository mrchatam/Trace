# P22-S04-04 — Review: regression attribution

## Metadata
- id: P22-S04-04
- todo_ids: [P22-S04-04]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C16** without auto-**`caused`** (Law 5 / D-22-08). Change association must be **queryable**, not inferred only from effect source at read time.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

| Check | Pass condition |
|-------|----------------|
| Law 5 | Grep: no code path sets `attribution=caused` without `SetRegressionAttributionCaused` + evidence |
| Create default | New regressions still **`correlated`** only |
| Link persisted | `regression_associated_change` in `entity_links` (not nullable column hack) |
| Auto-link | Contradicted-effect regression has link; attribution stays correlated |
| Caused test | `TestRegressionLinkedToChangeCausedWithEvidence` PASS |
| Query | `ListRegressionsByChangeID` returns linked rows; empty when no link |
| Schema | Still **24** sql files; **no 025+** |
| Keepers | `TestCorrelationAndContradictionNeverAutoSetCaused`, `TestSetAttributionCausedFailClosedWithoutEvidence` PASS |
| C08 hold | S04-01 predict/compare tests still PASS |

## Spawn policy

If C16 unmet: spawn **`P22-S04-04a` + `P22-S04-04b`**. Do not close with residuals.

## Named tests (re-run)

```bash
go test ./internal/domain/... -count=1 -run 'TestRegression|TestSetAttributionCaused|TestCorrelationAndContradiction|TestListRegressionsByChange'
go test ./internal/domain/... -count=1 -run 'TestRecordPredictedImpact|TestImpactCompare'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [ ] C16 closed or spawned
- [ ] Confidence **high** \| **medium** \| **low**
- [ ] Board Notes
