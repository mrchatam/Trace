# P42-S00-02 — Review (G6 non-semantic concept)

## Metadata
- id: P42-S00-02
- todo_ids: [P42-S00-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter, security-and-hardening]
- mcps: [user-trace, user-codegraph]
- verification: mixed

## Objective

Fresh independent review of S00-01 G6 implementation vs REMEDIATION-PLAN G6, GAP-REGISTRY G-004b, M-001 moat, Laws 6–7/19, and S00-00 LAW-REVIEW-NOTES.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md) — G6-C1–C7 acceptance map + touch-list
- [00-PLANNER.md](00-PLANNER.md) — locks + live repo gap
- [LAW-REVIEW-NOTES.md](LAW-REVIEW-NOTES.md) — S00-00 desk-check
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [REMEDIATION-PLAN G6](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-004b](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Pre-ship baseline: [h4-gf-extracted-inferred.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h4-gf-extracted-inferred.md)

## Session start

Follow agent-loop-protocol Session start. **Fresh subagent** — do not share implementer session.

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Medium+ confidence; zero open blocker/high |
| Spawn trigger | Blocker/high → spawn 02a/02b below this row |
| LAW-REVIEW | S00-00 PASS must match shipped code (no DR-NOSSEM slip) |
| Default behavior | Concept channel enriches packet when matches exist; no behavior change when empty |

## Preflight evidence commands

Run fresh subagent — do not trust implementer Notes alone.

```bash
# LAW-REVIEW still PASS
grep -q '^\*\*PASS\*\*' docs/phases/phase-42-concept-index/scopes/scope-00-non-semantic-concept/LAW-REVIEW-NOTES.md

# Reason code shipped
grep -n 'ReasonGraphLabelMatch\|graph_label_match' internal/retrieval/types.go internal/retrieval/doc.go

# No semantic slip
rg -n 'semantic_match|embedding|vector' internal/retrieval/concept.go || true

# Compile merge wired
rg -n 'SearchGraphLabels' internal/compiler/compiler.go internal/compiler/explore.go

# Tests green
go test ./internal/retrieval/... ./internal/compiler/... -count=1 -run 'GraphLabel|G6'
```

## Review checklist

### A — G6 gap closure

- [ ] `graph_label_match` reason_code on concept channel hits
- [ ] `SearchGraphLabels` (or equivalent) in `internal/retrieval/`
- [ ] Concept entity types bounded (discovery/assumption/decision/goal/claim)
- [ ] Compiler merge after FTS, before file-seed expand (`compiler.go:~180`)
- [ ] Explore merge appends to `SearchHits` with dedupe (`explore.go:~124`)
- [ ] Fail-open on concept path errors (DF-87) — compile/explore do not abort
- [ ] Dedupe: same entity not duplicated as both `fts_match` and `graph_label_match`
- [ ] G6-C1–C7 evidenced green
- [ ] `RETRIEVAL_AND_CONTEXT.md` §2 graph-label honesty updated (semantic section still deferred)

### B — M-001 moat

- [ ] Compile/explore require task_id — concept merge not query-only
- [ ] Task Layer-0 core always present
- [ ] Concept hits merge into packet — not standalone dump API

### C — Laws 6–7

- [ ] Limit cap 64 honored (G6-C4)
- [ ] No full-graph scan (G6-C7 + `TestNoDumpAPI` regression)
- [ ] Same budget competition as other candidates

### D — Law 19

- [ ] Logic in `internal/retrieval/` + compiler merge
- [ ] No business logic in MCP/CLI beyond existing paths

### E — DR-NOSSEM / G-004a

- [ ] No vector/embedding imports or indexes
- [ ] No `semantic_match` reason_code
- [ ] LAW-REVIEW-NOTES assertions match code

### F — G9 boundary

- [ ] G6 uses G9 intent terms — does not replace ExtractIntent
- [ ] Distinct reason_code from `fts_match`

### G — Rejects

- [ ] No LLM concept extraction
- [ ] No new always-on daemon
- [ ] No default cap inflation

## Exit criteria

- APPROVE with medium+ confidence; zero open blocker/high without spawn
- Board row → `done` with verdict + confidence in Notes

## Next

`P42-S01-00`
