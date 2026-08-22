# P22-S06-02b — Review: similar-changes ordering fix

## Metadata
- id: P22-S06-02b
- todo_ids: [P22-S06-02b]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C19, C20** closed after **P22-S06-02a**. Re-run full **P22-S06-02** checklist. Box checklist when both true.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

1. `TestQuerySimilarChanges` PASS `-count=30` (non-flaky).
2. Same-second / tie-break test exists and PASS.
3. **C19 hold:** `TestPatternCountsFromChangesAndOutcomes` PASS; deterministic refresh.
4. **No ML:** grep patterns paths — zero ML libs/LLM calls.
5. Schema **25**; compat **25**; MCP **13** unchanged.
6. Similar JSON paths-only; limits 32/64; S05 keepers PASS.

## Re-run commands

```bash
go test ./internal/domain/... -count=30 -run TestQuerySimilarChanges$
go test ./internal/domain/... ./internal/store/... -count=1 -run 'TestPatternCountsFromChangesAndOutcomes|TestQuerySimilarChanges|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestPatterns|TestChangesSimilar|TestChangesCompare|TestCLIChanges'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered
go test ./internal/store/... -count=1 -run TestSyncEntityFTSChange
ls internal/store/schema/*.sql | wc -l
```

## Spawn policy

If C20 still flaky or C19 regressed: spawn **`P22-S06-02c` + `P22-S06-02d`**.

## Exit criteria

- [ ] C19, C20 closed (checklist boxed)
- [ ] Confidence **high**
- [ ] Board Notes
