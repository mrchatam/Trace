# P22-S05-04 — Review: evidence queries

## Metadata
- id: P22-S05-04
- todo_ids: [P22-S05-04]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C17, C31, C32, C33, C34** — library queries + CLI; MCP for regressions only.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

1. **C31:** `trace tests verifying` uses **`validates` graph** (`ListValidatesForSymbol`/`ForFile`), not grep or regex on filenames.
2. **C32:** `trace outcomes failed` reads **stored** `outcome_results` (test fail/error); no subprocess re-run.
3. **C33:** `trace outcomes worked` includes **improvements** (S04) and/or passing test rows — honest “what worked” surface.
4. **C17/C34:** `trace regressions list` + MCP `trace_regressions` query stored regressions; change filter uses **`regression_associated_change`** (S04), not invented columns.
5. G19: domain helpers own SQL; CLI/MCP encode only.
6. Limits 32/64 enforced on list endpoints.
7. MCP catalog **13**; S05-01 tools unchanged.
8. Schema **24**; compat PASS.
9. S01 keeper: `TestValidatesEdgeExtractedFromImport` or `TestTestsVerifyingQuery` uses real validates edge.
10. S04 keeper: `TestListRegressionsByChangeID` PASS.

## Spawn policy

If any unmet: spawn **`P22-S05-04a` + `P22-S05-04b`**. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/domain/... ./internal/analyzers/... -count=1 -run 'TestTestsVerifying|TestOutcomesFailed|TestRegressionsList|TestListRegressionsByChange|TestValidatesEdgeExtractedFromImport'
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/mcp/... -count=1 -run 'TestTestsVerifying|TestOutcomes|TestRegressions|TestMCPRegressions|TestToolNamesRegistered'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 24
```

## Exit criteria

- [ ] C17, C31–C34 closed or spawned
- [ ] Confidence **high** | **medium** (must spawn if medium+unmet)
- [ ] Board Notes: checklist boxed when closed
