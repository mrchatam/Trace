# G6 LAW-REVIEW — Non-semantic concept retrieval (desk-check)

**Row:** P42-S00-00 · **Date:** 2026-08-22 · **GAP:** G-004b · **Verdict:** **PASS**

Desk-check only — no product code in this row. Implement row **P42-S00-01** unblocked.

## Authority

| Source | Lock |
|--------|------|
| [INTAKE.md](../../INTAKE.md) | G-004a vector **out** (permanent DR-NOSSEM defer) |
| [REMEDIATION-PLAN §2 G6](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) | Accept G-004b graph-label channel; law gate before build |
| [GAP-REGISTRY G-004b](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md) | Label/summary/concept via graph — **gap**, non-semantic OK |
| [00-PLANNER.md](00-PLANNER.md) | Locked defaults for S00-01 |

## Proposed mechanism (S00-01 scope)

```text
task (+ optional query)
  → G9 ExtractIntent (keywords + entity hints)
  → SearchGraphLabels: bounded FTS5 over concept entity types only
  → reason_code graph_label_match (distinct from fts_match)
  → merge into compile/explore candidates (fail-open DF-87)
  → existing G8 layer admission + budget caps
```

**Primary deliverable:** FTS over `{discovery, assumption, decision, goal, claim}` bodies/titles indexed in `internal/store/fts.go`.

**Explicit defer (not required for G6 closure):** optional 1-hop task-graph label attach via existing Expand neighbors — peer study only ([h4-gf-extracted-inferred.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h4-gf-extracted-inferred.md)); Trace already carries `edge_provenance` on structural import hops (`expand.go`); do **not** port GF confidence vocabulary or full GRAPH_REPORT semantics in S00.

## Desk-check matrix

| Check | Criterion | Evidence | Result |
|-------|-----------|----------|--------|
| **DR-NOSSEM** | No embeddings, vector index, or `semantic_match` reason_code | `internal/retrieval/doc.go:8–9` forbids semantic_match; no vector deps in retrieval package; G-004a locked **reject** in INTAKE + P42-00 | **PASS** |
| **Channel** | Lexical FTS on graph-adjacent entity text only — not semantic similarity | Proposed `SearchGraphLabels` reuses `store.SearchFTS` (same SQLite FTS5 leg as `search.go`); filters to concept entity types; new locked reason_code `graph_label_match` documented in `doc.go` | **PASS** |
| **Caps** | limit ≤ 64; same candidate pool discipline as Search | `SearchOptions` hard cap 64 (`types.go:60`, `search.go:19–20`); compiler uses `Limit: 16` today (`compiler.go:161`); concept channel inherits SearchOptions — no unbounded entity scan | **PASS** |
| **Moat (M-001)** | Requires task_id on compile/explore merge path; enriches packet — never query-only replacement | `compiler.Context` requires task (`compiler.go:155–180`); explore requires `opts.TaskID` (`explore.go:115–117`); MCP direct search remains Intent-free legacy path — G6 does not add query-only concept API | **PASS** |
| **Law 6–7** | No dump API; progressive budget competition | No new export/dump surface; concept hits compete in existing candidate merge → layer admission; `TestNoDumpAPI` regression preserved | **PASS** |
| **Law 19** | Library first; adapters thin | Logic confined to `internal/retrieval/concept.go`; `compiler.go` / `explore.go` append-only merge (mirror title FTS fail-open) | **PASS** |
| **G9 boundary** | Complementary — G9 feeds terms; G6 adds distinct channel | G9 `ExtractIntent` precedes channels (`intent.go:49–53`); G9 hits stay `fts_match`; G6 adds entity-type-filtered leg with own reason_code — does not replace ExtractIntent or G1 query merge | **PASS** |
| **Peer pattern** | GF EXTRACTED/INFERRED study only — not vector leg | MP BM25 text leg acceptable analog; GF provenance tiers inform optional future hop labeling only — **no** embedding leg, **no** semantic_match | **PASS** |

## Live repo gap (re-verified 2026-08-22)

| Item | Shipped | S00-01 action |
|------|---------|---------------|
| `graph_label_match` reason_code | Absent (`types.go:8–26`) | Add `ReasonGraphLabelMatch` + `doc.go` honesty |
| Concept entity FTS channel | Generic `Search` → all types → `fts_match` | New `SearchGraphLabels` with type filter |
| Compiler concept merge | Title + G9 intent FTS only | Append concept hits after intent FTS, before file-seed expand |
| Explore concept merge | Search hits only | Append concept hits to `SearchHits` (dedupe by entity) |
| §2 graph-label bullet | Aspirational “Semantic” section; no graph-label shipped line | REVISE `RETRIEVAL_AND_CONTEXT.md` §2 |

## Accept / reject (locked for implement)

| Decision | Item |
|----------|------|
| **Accept** | `internal/retrieval/concept.go` — `SearchGraphLabels(ctx, intent, opts)` |
| **Accept** | `ReasonGraphLabelMatch = "graph_label_match"` in `types.go` + `doc.go` |
| **Accept** | Compile/explore merge (fail-open DF-87) |
| **Accept** | Tests G6-C1–C7 |
| **Accept** | §2 doc revise (graph-label shipped, semantic still deferred) |
| **Reject** | G-004a vector / embeddings / `semantic_match` |
| **Reject** | Query-only concept search without task_id on compile path |
| **Reject** | Unbounded entity-type scan / full-graph dump |
| **Reject** | Replacing G9 intent or G1 query merge |
| **Reject** | New MCP tool / HTTP concept endpoint in S00 |
| **Reject** | LLM concept extraction |

## Risk review

| Risk | Mitigation |
|------|------------|
| DR-NOSSEM slip into embeddings | Desk-check + G6-C5 grep gate; no new deps in `concept.go` |
| Packet bloat | Same limit cap as Search; budget/layer admission unchanged |
| Duplicate hits (fts_match + graph_label_match) | Dedupe by entity key; prefer graph_label_match when both match (implement detail in `01-implement.md`) |
| Semantic creep via “concept” naming | Channel is **graph-label lexical match** — document as non-semantic in `doc.go` + §2 |

## Overall verdict

**PASS** — G6 graph-label channel under DR-NOSSEM is lawful, bounded, moat-preserving, and G9-complementary. **P42-S00-01** may proceed.

## Next

`P42-S00-01` — cite this file PASS in implement Notes.
