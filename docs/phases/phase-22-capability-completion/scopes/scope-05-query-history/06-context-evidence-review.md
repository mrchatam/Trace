# P22-S05-06 — Review: context evidence + S05 close

## Metadata
- id: P22-S05-06
- todo_ids: [P22-S05-06]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C35, C42-surface** and **S05 scope close** — all owned caps (C17, C29–C34, C37) hold with S05-02/S05-04 evidence.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

### S05-05 (C35, C42-surface)

1. `compiler.Packet` includes **`evaluations`**, **`reflections`**, **`planning_evidence`** — task-scoped, cap **8** each.
2. `trace.loop.next.v1` adds **`planning_evidence`** section; schema version string unchanged.
3. MCP `trace_context` output includes new fields (inherits compiler — no fork).
4. No score/runner blobs — truncated summaries only.
5. `TestContextIncludesEvaluationsAndReflections`, `TestLoopNextPlanningEvidenceSection` PASS.

### S05 regression (all owned caps)

6. **C29/C30/C37:** S05-02 closed — search + changes CLI/MCP.
7. **C17/C31–C34:** S05-04 closed — evidence queries + regressions MCP.
8. MCP catalog final **13** tools; compat **24**; schema **24** sql files.
9. S01–S04 keepers spot-check: `TestImpactWalkIncludesAffectedTests`, `TestChangesCompare`, `TestRegressionDetectedVsPriorPassingTest` (30/30 if run), `TestRecordImprovement`.
10. No **025+** migration; no daemon/HTTP added.

## Spawn policy

If any S05 capability unmet: spawn **`P22-S05-06a` + `P22-S05-06b`**. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/compiler/... ./internal/loop/... -count=1 -run 'TestContextIncludesEvaluations|TestLoopNextPlanningEvidence|TestTaskContextAndBudgets'
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/mcp/... -count=1 -run 'TestCLISearch|TestCLIChanges|TestTestsVerifying|TestOutcomes|TestRegressions|TestContext|TestLoopNext|TestToolNamesRegistered'
go test ./internal/domain/... -count=30 -run TestRegressionDetectedVsPriorPassingTest
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 24
```

## Exit criteria

- [ ] All S05 capabilities (C17, C29–C35, C37, C42-surface) closed or spawned
- [ ] Confidence **high** | **medium** (must spawn if medium+unmet)
- [ ] Board Notes: **S05 complete** → next `P22-S06-00`
- [ ] Checklist: box all S05-owned `[ ]` lines with evidence summary
