# P22-S04-02 — Review: predicted vs actual

## Metadata
- id: P22-S04-02
- todo_ids: [P22-S04-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C08** with **stored** prediction vs actual blast-key diff — not a one-off CLI print or effects-table reuse.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md), [01-predicted-vs-actual.md](01-predicted-vs-actual.md)
- Live post-01: `impact_predictions`, domain compare, CLI `predict`/`compare`

## Review checklist

| Check | Pass condition |
|-------|----------------|
| Schema | Exactly **24** sql files; **`024_impact_compare.sql`** only new mig; **no 025+** |
| Compat | `TestCompatibilitySecurityChecklist` ceiling **24** |
| Storage | Row keyed by `change_id`; `predicted_json` + persisted `compare_json` after compare |
| No blobs | Grep predicted_json / compare_json tests — no file contents, patches, or source |
| Walk reuse | Compare calls `retrieval.ImpactWalk` (or shared helper), not ad-hoc SQL graph |
| Affected tests | `affected_test_keys` in snapshot; S01 keeper `TestImpactWalkIncludesAffectedTests` PASS |
| Fail-closed | Compare without predict fails; unindexed path predict fails (named test or grep + code path) |
| Effects boundary | C08 **not** implemented solely via `effects.expected/actual` |
| MCP catalog | Still **10** tools (`TestToolNamesRegistered`) |
| S01/S03 hold | Spot-check `TestRegressionDetectedVsPriorPassingTest` `-count=30` if time permits |

## Spawn policy

If C08 unmet: spawn **`P22-S04-02a`** (implement fix) + **`P22-S04-02b`** (re-review) immediately below. Do not close with residuals.

## Named tests (re-run)

```bash
go test ./internal/domain/... ./internal/store/... ./internal/retrieval/... -count=1 -run 'TestRecordPredictedImpact|TestImpactCompare|TestImpactWalkIncludesAffectedTests'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestImpact'
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered
ls internal/store/schema/*.sql | wc -l
```

## Exit criteria

- [ ] C08 closed (checklist boxed) or **`P22-S04-02a/b`** spawned
- [ ] Confidence **high** \| **medium** \| **low** with inline fix or spawn
- [ ] Board Notes: verdict + test output
