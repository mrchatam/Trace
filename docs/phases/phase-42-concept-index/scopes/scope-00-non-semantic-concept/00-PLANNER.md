# P42-S00-00 — Scope planner (G6 non-semantic concept)

## Metadata
- id: P42-S00-00
- todo_ids: [P42-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, context-engineering, api-and-interface-design]
- mcps: [user-trace, user-codegraph]
- verification: automated

## Objective

Lock S00 **G6** against live repo: non-semantic concept retrieval (G-004b). Run **DR-NOSSEM law review desk-check** (no product code). Thicken `01-implement.md` + `02-review.md` with file targets, acceptance map, and rejects. **No product code in this row.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md) — P42-00 Q1 resolution (**desk-check at S00-00**)
- [REMEDIATION-PLAN §2 G6](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-004b](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [RETRIEVAL_AND_CONTEXT.md §2–§3](../../../../RETRIEVAL_AND_CONTEXT.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors (verified 2026-08-22 P42-00):
  - `internal/retrieval/doc.go:8–9` — DR-NOSSEM; no `semantic_match`
  - `internal/retrieval/types.go:8–26` — locked reason_codes; no graph-label code yet
  - `internal/retrieval/search.go:9–44` — FTS + G9 intent; all hits `fts_match`
  - `internal/retrieval/intent.go:49–53` — G9 precedes channels; complementary to G6
  - `internal/retrieval/expand.go:14–80` — structural/causal expand; no label/summary leg
  - `internal/store/fts.go:151–160` — FTS indexes goal/task/decision/assumption/discovery/… bodies
  - `internal/compiler/compiler.go:155–180` — compile path: Expand → intent FTS on **title** + query merge
  - Evidence: [h4-gf-extracted-inferred.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h4-gf-extracted-inferred.md)

## Session start

Follow agent-loop-protocol Session start. Unattended: INTAKE + P42-00 locks are authority.

## Locked defaults (FINAL — P42-00)

| Item | Value |
|------|-------|
| GAP ids | G-004b only (**G-004a vector forbidden**) |
| Verdict | **Accept** per REMEDIATION-PLAN G6 |
| P42-00 Q1 | **Law review desk-check at S00-00**; implement at S00-01 |
| Channel name | **Graph-label concept retrieval** — non-semantic summary/label leg |
| Mechanism | Intent keywords → bounded FTS over graph-adjacent **text entities** (discovery, assumption, decision, goal, claim) + optional 1-hop task-graph attach; distinct `graph_label_match` reason_code |
| M-001 | Merges into compile/explore task packet — never query-only replacement |
| Law 6–7 | Same caps as Search (limit default 16, hard max 64); no dump API |
| G9 boundary | G9 builds FTS terms; G6 adds **entity-type-filtered graph-label channel** with own reason_code — not duplicate intent |
| Library first | Logic in `internal/retrieval/`; compiler/explore merge only |
| Out | Vector/embeddings (G-004a); LLM concept extraction; standalone concept MCP tool; full-graph label scan |

## Live repo gap (re-verified P42-00)

| Check | Shipped | Gap |
|-------|---------|-----|
| Concept/graph-label reason_code | Absent | Need `graph_label_match` (document in `doc.go`) |
| Compiler FTS seed | Title + G9 intent terms → generic FTS | Misses graph-label channel for linked summaries |
| Entity-type concept filter | FTS searches all types equally | Need bounded concept entity set |
| EXTRACTED/INFERRED edge concept hop | Provenance on import hops only | Optional 1-hop label attach from task graph — not full GF port |
| Law review artifact | None | S00-00 produces `LAW-REVIEW-NOTES.md` desk-check |

## Accept / reject (G6)

| Decision | Item |
|----------|------|
| **Accept** | `internal/retrieval/concept.go` — `SearchGraphLabels(ctx, intent, opts)` bounded channel |
| **Accept** | New reason_code `graph_label_match` in `types.go` + `doc.go` |
| **Accept** | Wire into `compiler.go` + `explore.go` candidate merge (fail-open like DF-87) |
| **Accept** | Tests G6-C1–C7 (see thickened `01-implement.md`) |
| **Accept** | `LAW-REVIEW-NOTES.md` PASS desk-check at S00-00 (no vector, caps, moat) |
| **Accept** | REVISE `RETRIEVAL_AND_CONTEXT.md` §2 graph-label bullet (shipped honesty) |
| **Reject** | G-004a vector / `semantic_match` reason_code |
| **Reject** | Query-only concept search without task_id on compile path |
| **Reject** | Unbounded entity-type scan / full-graph dump |
| **Reject** | Replacing G9 intent or G1 query merge |

## Law review desk-check (S00-00 deliverable — no product code)

Author `scopes/scope-00-non-semantic-concept/LAW-REVIEW-NOTES.md`:

| Check | Pass criterion |
|-------|----------------|
| DR-NOSSEM | No embeddings, no vector index, no `semantic_match` |
| Channel | Lexical FTS on graph-adjacent entity text only |
| Caps | limit ≤ 64; same candidate pool as Search |
| Moat | Requires task_id on compile/explore merge path |
| Peer pattern | GF EXTRACTED/INFERRED study only — not vector leg (MP BM25 text ok as analog) |

**S00-01 blocked** until LAW-REVIEW-NOTES shows **PASS**.

## Must lock for S00-01 (delivered in thickened 01-implement)

1. Touch-list: concept.go → types/doc reason_code → compiler/explore merge → tests → §2 doc.
2. Seven acceptance tests G6-C1–C7.
3. LAW-REVIEW-NOTES PASS cited in S00-01 Notes.

## Exit criteria

- [ ] `LAW-REVIEW-NOTES.md` authored with PASS desk-check
- [ ] `01-implement.md` + `02-review.md` runnable with file targets + G6 accept map
- [ ] SCOPE-TODOS G6-0 marked done
- [ ] Board row → `done` with evidence in Notes

## Next

`P42-S00-01`
