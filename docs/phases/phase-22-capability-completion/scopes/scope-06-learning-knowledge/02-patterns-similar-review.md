# P22-S06-02 — Review: patterns + similar

## Metadata
- id: P22-S06-02
- todo_ids: [P22-S06-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C19, C20** without ML (D-22-11, W-27).

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-patterns-similar.md](01-patterns-similar.md) — implementer deliverables

## Review checklist

1. **C19 (patterns):** `RefreshChangePatterns` aggregates `(change_kind, outcome_kind)` with positive/negative counts; **`TestPatternCountsFromChangesAndOutcomes`** PASS; deterministic (same DB → same counts).
2. **C20 (similar):** `QuerySimilarChanges` returns prior changes + effects before new work; CLI `trace changes similar` + domain API; **`TestQuerySimilarChanges`** PASS.
3. **No ML:** grep `internal/domain/patterns.go`, `internal/store/patterns.go`, CLI — zero ML libs, zero LLM calls, zero statistical models beyond counting.
4. **change_kind lock:** `InferChangeKind` uses first segment of lexicographically smallest path → `seg:<segment>`; tested.
5. **outcome_kind lock:** priority order matches planner; one bucket per change.
6. **Schema:** exactly **25** sql files; **`025_engineering_knowledge.sql`** only new mig; compat/embed **25** PASS.
7. **Both tables in 025:** `change_patterns` populated; `engineering_knowledge` may be empty (S06-03 owns rows) — table must exist.
8. **G19:** MCP handlers unchanged; CLI calls domain only.
9. **Limits:** similar query default 32, cap 64.
10. **No blobs:** similar JSON paths-only; no file content, patches, or diffs.
11. **S05 keepers:** `TestCLIChangesList`, `TestChangesCompare`, `TestSyncEntityFTSChange` still PASS.
12. **MCP catalog:** **13** tools; `TestToolNamesRegistered` unchanged.

## Spawn policy

If unmet: spawn **`P22-S06-02a` + `P22-S06-02b`**. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/domain/... ./internal/store/... -count=1 -run 'TestPatternCountsFromChangesAndOutcomes|TestQuerySimilarChanges|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestPatterns|TestChangesSimilar|TestChangesCompare|TestCLIChanges'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered
ls internal/store/schema/*.sql | wc -l  # expect 25
rg -i 'openai|anthropic|tensorflow|torch|sklearn|embedding model' internal/domain/patterns.go internal/store/patterns.go cmd/trace/patterns.go || true
```

## Exit criteria

- [ ] C19, C20 closed or spawned
- [ ] Confidence **high** | **medium** (must spawn if medium+unmet)
- [ ] Board Notes: findings + confidence + checklist boxed when closed
