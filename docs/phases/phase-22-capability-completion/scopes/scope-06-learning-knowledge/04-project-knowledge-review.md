# P22-S06-04 — Review: project knowledge

## Metadata
- id: P22-S06-04
- todo_ids: [P22-S06-04]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C10, C21, C26, C27** and seed portability (W-26).

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [03-project-knowledge.md](03-project-knowledge.md) — implementer deliverables

## Review checklist

1. **C10:** Reconsideration-sourced knowledge links **`decision_id`**; `TestKnowledgeLinksDecision` PASS.
2. **C21:** Knowledge synthesized from historical changes/patterns/reflections/improvements — not hand-authored only.
3. **C26:** Rows persist with `created_at`/`updated_at`; synthesize updates in place; accumulates over multiple runs.
4. **C27:** Second synthesize adds/updates rows — project-specific corpus grows (`TestSynthesizeKnowledgeFromPatterns` + idempotency spot-check).
5. **No LLM (W-27):** grep knowledge paths — zero model/API calls.
6. **Provenance:** `evidence_ids_json` validated; `source_type` set per source (decision_reconsideration, reflection, pattern, improvement).
7. **Seed (D-22-19):** `change_patterns[]` + `engineering_knowledge[]` export/import; `TestSeedExportIncludesKnowledge` PASS; S04 `improvements[]` keeper PASS.
8. **Schema:** still **25** sql; no **026+**; compat **25** PASS.
9. **S06-01 hold:** `TestPatternCountsFromChangesAndOutcomes`, `TestQuerySimilarChanges` PASS.
10. **G19:** CLI thin; no MCP SQL fork.
11. **No blobs:** `body_json` structured metadata only; no source file content.

## Spawn policy

If unmet: spawn **`P22-S06-04a` + `P22-S06-04b`**. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/domain/... -count=1 -run 'TestUpsertEngineeringKnowledge|TestSynthesizeKnowledgeFromPatterns|TestKnowledgeLinksDecision|TestSeedExportIncludesKnowledge|TestSeedExportIncludesImprovements|TestPatternCountsFromChangesAndOutcomes'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestKnowledge|TestPatterns'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 25
rg -i 'openai|anthropic|llm|chat\.completions' internal/domain/knowledge.go internal/store/knowledge.go cmd/trace/knowledge.go || true
```

## Exit criteria

- [ ] C10, C21, C26, C27 closed or spawned
- [ ] Confidence **high** | **medium** (must spawn if medium+unmet)
- [ ] Board Notes: findings + confidence + checklist boxed when closed
