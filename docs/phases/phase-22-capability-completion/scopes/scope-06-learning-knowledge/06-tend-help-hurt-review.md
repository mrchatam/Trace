# P22-S06-06 — Review: tend help/hurt + S06 close

## Metadata
- id: P22-S06-06
- todo_ids: [P22-S06-06]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C22, C23, C24** and **all S06 owned capabilities** (C10, C19–C24, C26, C27). Close scope S06.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks + full capability list
- [05-tend-help-hurt.md](05-tend-help-hurt.md) — implementer deliverables

## Review checklist — S06-05 deliverables

1. **C22:** `tendencies` in context + loop; threshold **2**; `TestTendHelpHurtInContext` PASS; CLI `trace knowledge tendencies`.
2. **C23:** `successful_approaches` in packet/loop merges worked outcomes + knowledge; `TestSuccessfulApproachesSurfaced` PASS; S05 `trace outcomes worked` keeper PASS.
3. **C24:** Loop next combines tendencies + successful_approaches + planning_evidence; `TestLoopNextIncludesEvidenceForDecisions` PASS.
4. **Caps:** tendencies, successful_approaches, similar_changes (if present) each cap **8**.
5. **MCP:** `trace_context` inherits — no new tools; catalog **13**.

## Review checklist — S06 scope close (all owned caps)

6. **C19/C20 hold:** S06-02 boxed; pattern + similar keepers PASS.
7. **C10/C21/C26/C27 hold:** S06-04 boxed; knowledge + seed keepers PASS.
8. **W-27 no ML:** grep S06 product paths — zero ML/LLM.
9. **Schema compat:** **25** sql files; **no 026+**; compat PASS.
10. **S01–S05 spot-check:** `TestImpactWalkIncludesAffectedTests`, `TestRegressionDetectedVsPriorPassingTest` (30/30), `TestCLIChanges`, `TestSyncEntityFTSChange` PASS.
11. **Checklist:** all nine S06-owned bullets `[x]` or spawn.

## Spawn policy

If any S06 capability unmet: spawn **`P22-S06-06a` + `P22-S06-06b`**. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/compiler/... ./internal/loop/... ./internal/domain/... -count=1 -run 'TestTendHelpHurtInContext|TestSuccessfulApproachesSurfaced|TestLoopNextIncludesEvidenceForDecisions|TestPatternCountsFromChangesAndOutcomes|TestQuerySimilarChanges|TestSynthesizeKnowledgeFromPatterns|TestKnowledgeLinksDecision|TestSeedExportIncludesKnowledge'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestContext|TestLoopNext|TestKnowledge|TestPatterns|TestChangesSimilar|TestOutcomesWorked|TestCLIChanges'
go test ./internal/domain/... -count=30 -run TestRegressionDetectedVsPriorPassingTest
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered
ls internal/store/schema/*.sql | wc -l  # expect 25
```

## Exit criteria

- [ ] C10, C19–C24, C26, C27 closed or spawned
- [ ] Confidence **high** | **medium** (must spawn if medium+unmet)
- [ ] Board Notes: **S06 complete** → next `P22-S07-00`; checklist all S06 caps boxed
