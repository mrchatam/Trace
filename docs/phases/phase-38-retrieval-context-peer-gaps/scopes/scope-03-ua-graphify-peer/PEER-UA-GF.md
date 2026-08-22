# PEER-UA-GF — UA + Graphify + Mempalace peer investigation

**Author:** P38-S03-01 (2026-08-22)  
**Status:** Investigation only — no product changes  
**Authority:** [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) §2 (H1 partial, H4, H6 MP slice, H8, H9 MP contrast)  
**Trace baseline:** [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md) (S01 APPROVED)  
**CG leg (do not re-litigate):** [PEER-CG.md](../scope-02-codegraph-peer/PEER-CG.md)  
**Evidence root:** [`experiments/runs/2026-08-22-p38-s03-657/evidence/`](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/)  
**Peer roots:** UA · Graphify · **Mempalace** (read-only under `similar projects/`)

---

## §1 UA mechanism

Understand Anything (UA) orients agents via a **query-driven `ChatContext` packet** assembled from a knowledge graph in one function call.

### ChatContext shape (`context-builder.ts` L9–18)

```typescript
export interface ChatContext {
  projectName: string;
  projectDescription: string;
  languages: string[];
  frameworks: string[];
  relevantNodes: GraphNode[];
  relevantEdges: GraphEdge[];
  relevantLayers: Layer[];
  query: string;
}
```

### Build pipeline (`context-builder.ts` L25–79)

1. **SearchEngine** fuzzy search on agent `query` (L32–34)
2. **1-hop edge expand** from matched nodes (L39–48)
3. **Layer filter** — layers containing any relevant node (L65–68)
4. Return unified packet including original `query` (L70–79)

### SearchEngine (`search.ts` L14–58)

| Mechanism | Lines | Detail |
|-----------|-------|--------|
| Fuse keys | L14–25 | `name` 0.4, `tags` 0.3, `summary` 0.2, `languageNotes` 0.1 |
| Extended OR | L44–47 | Token OR via `\|` join — tolerates typos |
| Optional semantic | `embedding-search.ts` L1+ | Cosine similarity when embeddings exist |

### Onboarding (`onboard-builder.ts` L7+)

Generates standalone markdown: architecture layers, key concepts, guided tour, file map, complexity hotspots — graph-derived docs artifact.

### Side-by-side: UA vs Trace task packet (Q1)

| Step | UA | Trace |
|------|----|-------|
| Agent question | `query` in `buildChatContext` L27 | Separate `trace search` — not merged (TRACE-AUDIT §4) |
| Neighborhood | 1-hop edge expand L39–48 | Graph expand at `--depth`; symbol/file gaps (TRACE-AUDIT H1) |
| Task/plan scope | **None** — graph chat only | L0 task packet via compiler/`trace_context` |
| Single-call orient | `buildChatContext` + `formatContextForPrompt` L85–147 | Compose `trace_search` + `trace_context` + optional `trace_why` |
| FTS input | Agent query via Fuse L36–58 | **`task.Title` only** — `compiler.go` L146–151 |

### Q2 — SearchEngine vs Trace compiler FTS

UA SearchEngine searches **multi-field node metadata** (name, tags, summary) driven by agent query. Trace compiler FTS inside context packets uses **`task.Title` tokens only** (`compiler.go` L148). CLI `trace search` accepts user queries but results do not auto-merge into packets (TRACE-AUDIT H2 — related, S01-owned).

---

## §2 Graphify mechanism

Graphify builds code knowledge graphs with **provenance-labeled edges** and ships **committed orient artifacts**.

### EXTRACTED vs INFERRED (`validate.py` L5–7)

```python
VALID_CONFIDENCES = {"EXTRACTED", "INFERRED", "AMBIGUOUS"}
```

| Confidence | Emission cite | Meaning |
|------------|---------------|---------|
| EXTRACTED | `symbol_resolution.py` L289–290, L498–499 | Import-guided or AST-proven edges; score 1.0 |
| INFERRED | `symbol_resolution.py` L318–370 | Name-based unresolved calls; score 0.8; docstring L318 notes "not import proof" |
| AMBIGUOUS | Schema only | Validator enum; 0% in rsl-siege-manager run |

### Worked example counts (`GRAPH_REPORT.md` L7–8)

- 1886 nodes · 3876 edges · 141 communities
- **90% EXTRACTED · 10% INFERRED** (393 INFERRED edges, avg confidence 0.62)

Cross-language INFERRED critique: `review.md` L68–91 — Python `AuthError` --uses--> TS `Member` is name similarity, not runtime use.

### Q3 — Concept discovery vs Trace FTS

Graphify graph nodes carry **labels, summaries, and typed relations** navigable by concept-style questions ("what connects auth to API?"). Trace retrieval is FTS5 + structural expand with **DR-NOSSEM** forbidding semantic channel (`doc.go` L8–9). Live probe: `trace search "compiler retrieval"` → 0 hits (`h4-trace-search-sample.txt`).

### Q4 — graph.html / worked examples vs Trace GUI

| Affordance | Graphify | Trace |
|------------|----------|-------|
| Committed artifact | `worked/rsl-siege-manager/graph.html` L68 — RAW_EDGES with `confidence: EXTRACTED` | No checked-in graph.html for dogfood |
| Report | `GRAPH_REPORT.md` — community hubs, EXTRACTED/INFERRED ratio | No GRAPH_REPORT equivalent |
| Install hook | `worked/example/README.md` L31–35 — `/graphify ./raw` | No slash-command graph hook |
| Interactive viz | `exporters/html.py` L30–68 — sidebar search, legend, community filter | `web/src/App.tsx` L21–22 — `/` → Graph route exists |
| Reproduction | `worked/rsl-siege-manager/README.md` L10–18 — clone + checkout + run steps | No worked example corpus in repo |

TRACE-AUDIT H8 was **inconclusive** (S01). S03 resolves: Trace has graph **route** but lacks peer-style **committed artifact + hook narrative + confidence-labeled orient UX**.

---

## §3 Mempalace mechanism

Mempalace (MP) provides **hybrid memory retrieval**, a **4-layer wake-up stack**, a **large MCP surface**, and a **shipped fact-checking pipeline**.

### Hybrid search (Q5)

| Component | File:line | Detail |
|-----------|-----------|--------|
| Design | `searcher.py` L5–9 | BM25 + vector semantic similarity |
| `_hybrid_rank` | `searcher.py` L276–329 | vector_weight 0.6 / bm25_weight 0.4; Okapi-BM25 min-max normalized + vector similarity |
| `search_memories` | `searcher.py` L1652–1718 | NL query; `candidate_strategy` vector/union; `vector_disabled` → BM25-only fallback |
| MCP entry | `mcp_server.py` L2369–2408 | `tool_search` → `sanitize_query` → `search_memories` |

**Trace contrast + DR-NOSSEM:** Vector/embedding channel is **product law forbidden** (`doc.go` L8–9). Gap remains for **rich-text BM25 + concept discovery** vs Trace title-token FTS.

### MCP + memory stack (Q6)

| Surface | Count | Cite |
|---------|-------|------|
| READ_TOOLS | 20 | `service.py` L60–82 |
| WRITE_TOOLS | 12 | `service.py` L85–99 |
| MAINTENANCE_TOOLS | 3 | `service.py` L102 |
| TOOLS dict | 44 handlers | `mcp_server.py` L4765+ |

**4-layer stack (`layers.py` L3–17, L404–431):**

| Layer | Tokens | API |
|-------|--------|-----|
| L0 Identity | ~100 | Always loaded — "Who am I?" |
| L1 Essential Story | ~500–800 | `wake_up()` L404–431 — L0 + L1 single packet |
| L2 On-Demand | ~200–500 | `recall()` L425–427 |
| L3 Deep Search | Unlimited | `search()` L429–431 — ChromaDB semantic |

Compare Trace: compiler L0–1 only (TRACE-AUDIT H3); 16 MCP tools (TRACE-AUDIT H6); no memory wake-up packet.

### KG + fact_checker (Q7)

| Component | File:line | Shipped behavior |
|-----------|-----------|------------------|
| Temporal KG | `knowledge_graph.py` L1–36 | SQLite entity triples with `valid_from`/`valid_to` |
| `check_text` | `fact_checker.py` L55–78 | Offline: `similar_name`, `relationship_mismatch`, `stale_fact` |
| Onboarding | `onboarding.py` L3–17 | First-run entity/wing seeding before indexing |

Trace: `RETRIEVAL_AND_CONTEXT.md` §3 intent pipeline — **zero** `intent` in `internal/retrieval/` (TRACE-AUDIT H9). MP contrast only; S01 Trace verdict stands.

---

## §4 Hypothesis verdicts

Aligned with INVESTIGATION-INDEX §2 verify/reject criteria. Cross-ref PEER-CG for CG leg of H1/H6.

| H | Verdict | Primary peer evidence | Trace contrast | Confidence |
|---|---------|----------------------|----------------|------------|
| **H1 partial** | **supported** | UA `buildChatContext` L25–79 (query + 1-hop + layers); MP `wake_up()` L404–431 (L0+L1 packet) | TRACE-AUDIT H1 — no query on `trace_context`; compose required; PEER-CG CG leg | **high** |
| **H4** | **supported** (defer vector leg) | GF EXTRACTED/INFERRED L289–370 + GRAPH_REPORT 90/10%; MP `_hybrid_rank` L276–329 | `doc.go` L8–9 DR-NOSSEM; title FTS only — embedding channel **law defer**; label/summary graph gap **real** | **high** |
| **H6 (MP slice)** | **supported** | MP 35 categorized / 44 TOOLS (`service.py` L60–102; `mcp_server.py` L4765+) | TRACE-AUDIT H6 — 16 tools; CG 1-tool (PEER-CG) — MP opposite extreme, no orient ranking | **high** |
| **H8** | **supported** | GF committed `graph.html` + `/graphify` hook; UA `onboard-builder.ts` L7+; MP `onboarding.py` + `wake_up()` | TRACE-AUDIT H8 inconclusive → **resolved supported**; App.tsx L21–22 graph route only | **high** |
| **H9 (MP contrast)** | **supported** | MP `fact_checker.py` L55–78 shipped pipeline | TRACE-AUDIT H9 — doc-only intent; zero retrieval code | **high** |

### H8 resolution note

S01 marked H8 **inconclusive** (medium-low). S03 peer read resolves: all three peers ship stronger **graph/memory-first onboarding artifacts** than Trace's task-loop-first GUI. Trace default route is Graph (`App.tsx` L21) but lacks committed worked examples, install hooks, and confidence-labeled interactive orient UX.

### H4 DR-NOSSEM nuance

| Channel | Status |
|---------|--------|
| Embedding/vector semantic_match | **Product law** — defer, not gap |
| Graphify label/summary + INFERRED concept edges | **Gap** — concept discovery beyond title FTS |
| MP hybrid BM25 over memory text | **Gap** — richer lexical channel than compiler title tokens |

---

## §5 Trace strengths peers lack (moat row seed)

Peers are **graph/memory orient-first**; Trace owns **directed work with evidence** (TRACE-AUDIT §5, PEER-CG §5).

| Peer | Lacks | Evidence |
|------|-------|----------|
| **UA** | Task loop, gates, evidence, plan tree | No task UUID in context-builder; read-only graph chat |
| **Graphify** | Same + progressive planning | Graph-only `/graphify` hook; no loop MCP |
| **Mempalace** | Task/plan/gate/evidence | Memory-first MCP; `wake_up()` is identity/story not backlog |

| Trace strength | Evidence |
|----------------|----------|
| Progressive task context packet | Budget, reason_codes, trust labels (TRACE-AUDIT §5) |
| Task loop + gate + review | `trace_loop`, `trace_review`, `trace_transition` |
| Enforcement harness | `TRACE_TASK_ID`, plan bootstrap |
| Why / causal trace | `trace_why` |
| Plan tree + orchestration | `trace_plan`, `trace_tasks`, 16-tool write surface |

**Q8 answer:** UA/GF/MP lack task loop, gate enforcement, evidence chain, and plan tree — Trace moat is real; gap is orient/retrieval merge, not abandoning task scope.

---

## §6 Evidence appendix + command log

### Evidence files

| File | Todo | Description |
|------|------|-------------|
| `t0-ua-root.txt` | T0 | UA peer present |
| `t0-gf-root.txt` | T0 | Graphify peer present |
| `t0-mp-root.txt` | T0 | Mempalace peer present |
| `t0-trace-nossem.txt` | T0 | DR-NOSSEM cite |
| `h1-ua-partial.md` | T1 | UA ChatContext vs Trace |
| `h1-ua-search-mechanism.md` | T2 | Fuse vs compiler FTS |
| `h4-gf-extracted-inferred.md` | T3 | EXTRACTED/INFERRED mechanism |
| `h4-semantic-contrast.md` | T3 | Concept query pair |
| `h4-trace-search-sample.txt` | T3 | Live Trace search probe |
| `h8-gf-onboarding-ux.md` | T4/T5 | Graphify orient UX |
| `h8-ua-onboard.md` | T4 | UA onboard-builder |
| `h4-mp-hybrid-search.md` | T6 | MP hybrid search |
| `h6-mp-mcp-surface.md` | T7 | MP MCP tool count |
| `h1-mp-context-packet.md` | T8 | MP wake_up packet |
| `h8-mp-onboarding.md` | T8 | MP onboarding |
| `h9-mp-fact-check-contrast.md` | T9 | MP fact_checker vs Trace |
| `moat-peers-lack.md` | T10 | Moat seed |

### Live command log

| # | Command / action | Exit | Evidence |
|---|------------------|------|----------|
| T0 | Preflight peer roots + DR-NOSSEM rg | 0 | `t0-*` |
| T3 | `trace search "compiler retrieval" --limit 5` | 0 | `h4-trace-search-sample.txt` |
| T5 | Trace GUI live screenshot | skipped | TRACE-AUDIT `h8-gui-partial.md` + App.tsx L21–22 |

### Must-answer checklist (Q1–Q8)

| Q | Answer location |
|---|-----------------|
| Q1 ChatContext vs Compiler packet | §1 side-by-side table |
| Q2 SearchEngine vs Trace FTS | §1 Q2; `h1-ua-search-mechanism.md` |
| Q3 GF EXTRACTED/INFERRED vs Trace | §2; `h4-gf-extracted-inferred.md` |
| Q4 graph.html vs Trace GUI | §2 Q4; `h8-gf-onboarding-ux.md` |
| Q5 MP hybrid vs Trace + DR-NOSSEM | §3; `h4-mp-hybrid-search.md` |
| Q6 MP MCP + layers vs Trace 16-tool | §3; `h6-mp-mcp-surface.md`, `h1-mp-context-packet.md` |
| Q7 MP KG + fact_checker vs Trace intent | §3; `h9-mp-fact-check-contrast.md` |
| Q8 Moat seed | §5; `moat-peers-lack.md` |

---

## §7 Spawn list

| Item | Trigger | Owner |
|------|---------|-------|
| _(empty)_ | H12+ not required — mempalace mapped to H1/H4/H6/H8/H9 per INVESTIGATION-INDEX | — |

No uncovered slice found requiring H12+ spawn. Reviewer may request live GUI screenshot in S03-02 if medium residual on H8 affordance detail.

---

**Overall confidence:** **high** — all S03-owned hypotheses have peer file:line cites; H8 upgraded from S01 inconclusive to supported; DR-NOSSEM split documented for H4.

**Next:** P38-S03-02 (independent review)
