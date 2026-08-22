# P38-S03-01 — UA + Graphify + Mempalace peer investigation

## Metadata
- id: P38-S03-01
- todo_ids: [P38-S03-01]
- role: implementer
- skills: [research, graphify, code-explorer]
- mcps: [user-codegraph]
- verification: mixed
- hooks: none

## Objective

Read-only investigation of **Understand Anything (UA)**, **Graphify**, and **Mempalace** as peers. Author **`PEER-UA-GF.md`** with **§1 UA**, **§2 Graphify**, **§3 Mempalace**, mechanism cites (peer `file:line`), hypothesis verdicts for **H1 (partial), H4, H6 (MP slice), H8, H9 (MP contrast partial)**, and moat seed. **Do not implement fixes or product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §2 verify/reject for H1, H4, H6, H8, H9 (MP contrast)
- [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md) — Trace baseline (S01 APPROVED)
- [PEER-CG.md](../scope-02-codegraph-peer/PEER-CG.md) — CG leg of H1/H6 (do not re-litigate)
- [INTAKE.md](../../INTAKE.md), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md), [PEER-FIXTURES.md](../../PEER-FIXTURES.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — planner locks from P38-S03-00

### Peer roots (read-only)

| Peer | Root | Primary anchors (spot-check before cite drift) |
|------|------|-----------------------------------------------|
| **UA** | [`similar projects/Understand-Anything/`](../../../../../similar%20projects/Understand-Anything/) | [`context-builder.ts`](../../../../../similar%20projects/Understand-Anything/understand-anything-plugin/src/context-builder.ts) L9–79; [`search.ts`](../../../../../similar%20projects/Understand-Anything/understand-anything-plugin/packages/core/src/search.ts) L27–59; [`onboard-builder.ts`](../../../../../similar%20projects/Understand-Anything/understand-anything-plugin/src/onboard-builder.ts) L7+ |
| **Graphify** | [`similar projects/graphify/`](../../../../../similar%20projects/graphify/) | [`validate.py`](../../../../../similar%20projects/graphify/graphify/validate.py) L5–7; [`symbol_resolution.py`](../../../../../similar%20projects/graphify/graphify/symbol_resolution.py) L289–370; [`exporters/html.py`](../../../../../similar%20projects/graphify/graphify/exporters/html.py) L30+; worked [`rsl-siege-manager/graph.html`](../../../../../similar%20projects/graphify/worked/rsl-siege-manager/graph.html) L67–68 |
| **Mempalace** | [`similar projects/mempalace/`](../../../../../similar%20projects/mempalace/) | [`searcher.py`](../../../../../similar%20projects/mempalace/mempalace/searcher.py) L5–9, L276–329, L1652+; [`layers.py`](../../../../../similar%20projects/mempalace/mempalace/layers.py) L3–17, L404–431; [`mcp_server.py`](../../../../../similar%20projects/mempalace/mempalace/mcp_server.py) L2369+, L4765+; [`service.py`](../../../../../similar%20projects/mempalace/mempalace/service.py) L60–102; [`knowledge_graph.py`](../../../../../similar%20projects/mempalace/mempalace/knowledge_graph.py) L1–36; [`fact_checker.py`](../../../../../similar%20projects/mempalace/mempalace/fact_checker.py) L1–78; [`onboarding.py`](../../../../../similar%20projects/mempalace/mempalace/onboarding.py) L3–17 |

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (P38-S03-00 — do not re-debate)

| Item | Value |
|------|-------|
| Output path | `scopes/scope-03-ua-graphify-peer/PEER-UA-GF.md` |
| Artifact sections | §1 UA · §2 Graphify · **§3 Mempalace** · §4 Hypothesis verdicts · §5 Moat · §6 Evidence · §7 Spawn |
| Product edits | **Forbidden** (Go/TS/web in Trace repo) |
| Method | **Peer file:line** — not README-only for major mechanism claims; website docs OK as orientation only |
| Evidence root | `experiments/runs/YYYY-MM-DD-p38-s03-657/evidence/` |
| Hypotheses (S03-owned) | **H1 (partial), H4, H8** + **H6, H9 (Mempalace contrast slices only)** |
| Verdict vocabulary | `supported` \| `weakened` \| `rejected` \| `inconclusive` (+ defer note if law/DR fit) |
| DR-NOSSEM | Trace [`internal/retrieval/doc.go`](../../../../../internal/retrieval/doc.go) L8–9 — semantic channel forbidden; H4 compares **peer capability**, not auto-build |
| H12 spawn | **Not required** — mempalace mapped to existing H* (human-added peer 2026-08-22); spawn H12+ only if investigation finds uncovered slice |
| Non-goals | No copying peers into Trace; no REMEDIATION-PLAN; no full monorepo scan |

## Must answer (planner handoff — embed in PEER-UA-GF.md)

1. **UA `context-builder.ts`:** `ChatContext` packet shape (query + 1-hop expand + layers) vs Trace `Compiler` task packet — side-by-side table.
2. **UA SearchEngine:** Query-driven fuzzy retrieval (`fuse.js`, weighted keys) — file:line vs Trace title FTS ([`compiler.go`](../../../../../internal/compiler/compiler.go) L146–151).
3. **Graphify EXTRACTED vs INFERRED:** Mechanism in source + worked example counts — relevance to H4 concept discovery vs Trace FTS-only.
4. **Graphify `graph.html` / worked examples:** Orient/onboarding UX (committed artifact, GRAPH_REPORT, interactive viz) vs Trace GUI ([`web/src/App.tsx`](../../../../../web/src/App.tsx) L21–22).
5. **Mempalace hybrid search:** BM25 + vector rerank (`_hybrid_rank`, `search_memories`) vs Trace retrieval — H4 contrast with DR-NOSSEM note.
6. **Mempalace MCP + memory stack:** Tool surface (`TOOLS`, `READ_TOOLS`/`WRITE_TOOLS`) and 4-layer wake-up (`layers.wake_up`) vs Trace 16-tool + compiler L0–1 — H1/H6/H8 partial.
7. **Mempalace KG + fact_checker:** Shipped contradiction pipeline vs Trace aspirational intent ([`RETRIEVAL_AND_CONTEXT.md`](../../../../RETRIEVAL_AND_CONTEXT.md)) — H9 contrast only (S01 owns Trace live).
8. **Moat seed:** What UA/GF/MP lack (task loop, gates, evidence, plan tree).

---

## Investigation todos (run in order; do not skip)

### T0 — Preflight + evidence folder

```bash
EV=experiments/runs/$(date +%Y-%F)-p38-s03-657/evidence
mkdir -p "$EV"
# Peer presence
test -f "similar projects/Understand-Anything/understand-anything-plugin/src/context-builder.ts" && echo ok | tee "$EV/t0-ua-root.txt"
test -f "similar projects/graphify/graphify/validate.py" && echo ok | tee "$EV/t0-gf-root.txt"
test -f "similar projects/mempalace/mempalace/searcher.py" && echo ok | tee "$EV/t0-mp-root.txt"
# Trace baseline pointers (do not re-run full S01)
test -f internal/retrieval/doc.go && rg -n "DR-NOSSEM|semantic" internal/retrieval/doc.go | tee "$EV/t0-trace-nossem.txt"
```

- Record date + row id in every evidence file header.
- Skim TRACE-AUDIT §2 H1/H8 rows — S03 compares **peer side**; cite TRACE-AUDIT for Trace contrast.

### T1 — H1 (partial): UA query+task context packet

**Peer code read (required file:line):**

| Mechanism | Start here |
|-----------|------------|
| `ChatContext` shape | [`context-builder.ts`](../../../../../similar%20projects/Understand-Anything/understand-anything-plugin/src/context-builder.ts) L9–18 — fields include `query`, nodes, edges, layers |
| Build pipeline | Same file L25–79 — SearchEngine → 1-hop edge expand → layer filter → return packet |
| Prompt formatting | Same file L85–147 — `formatContextForPrompt` markdown sections |

**Trace contrast (TRACE-AUDIT H1 partial):**

- `trace_context` / MCP `trace_context` requires `task_id`; no unified `query` in one call.
- [`internal/compiler/compiler.go`](../../../../../internal/compiler/compiler.go) — task-scoped packet; FTS on title L146–151.

**Side-by-side table (required in PEER-UA-GF §1 + §4):**

| Step | UA | Trace |
|------|----|-------|
| Agent question | `query` string in `buildChatContext` | separate `trace_search` (not merged into packet) |
| Neighborhood | 1-hop edge expand from search hits | graph expand at depth; symbol/file gaps per TRACE-AUDIT |
| Task/plan scope | **none** (graph chat) | L0 task packet |
| Single-call orient | `buildChatContext` + format | compose tools |

**Verdict target:** `h1-ua-partial.md`

### T2 — UA SearchEngine vs Trace compiler FTS

**Peer code read:**

- [`search.ts`](../../../../../similar%20projects/Understand-Anything/understand-anything-plugin/packages/core/src/search.ts) L14–25 — Fuse keys/weights (`name`, `tags`, `summary`, `languageNotes`).
- Same file L36–58 — extended OR token query; type filter; limit.
- Optional: [`embedding-search.ts`](../../../../../similar%20projects/Understand-Anything/understand-anything-plugin/packages/core/src/embedding-search.ts) — semantic path when embeddings exist (note vs Trace DR-NOSSEM).

**Trace contrast:**

- TRACE-AUDIT H2 — title-only FTS in compiler; document as **related** not S03-owned verdict.

**Verdict target:** `h1-ua-search-mechanism.md` — feeds H1 partial evidence; do not duplicate H2 verdict row.

### T3 — H4: Graphify EXTRACTED vs INFERRED edges

**Peer code read (required — not README-only):**

| Mechanism | Start here |
|-----------|------------|
| Schema / valid confidences | [`validate.py`](../../../../../similar%20projects/graphify/graphify/validate.py) L5–7 — `EXTRACTED`, `INFERRED`, `AMBIGUOUS` |
| EXTRACTED emission | [`symbol_resolution.py`](../../../../../similar%20projects/graphify/graphify/symbol_resolution.py) L289–290, L498–499 |
| INFERRED emission | Same file L318–370 — name-based / unresolved call edges |
| Worked example counts | [`worked/rsl-siege-manager/GRAPH_REPORT.md`](../../../../../similar%20projects/graphify/worked/rsl-siege-manager/GRAPH_REPORT.md) — EXTRACTED/INFERRED ratio; [`review.md`](../../../../../similar%20projects/graphify/worked/rsl-siege-manager/review.md) L68–91 cross-language INFERRED critique |

**Trace contrast:**

- [`internal/retrieval/doc.go`](../../../../../internal/retrieval/doc.go) L8–9 — no semantic; FTS + expand only.
- Document **query pair** (concept-style question): Graphify label/summary vs Trace `trace search` — evidence in `h4-semantic-contrast.md`.

**Verdict target:** `h4-gf-extracted-inferred.md`, `h4-semantic-contrast.md`

### T4 — H8: Graphify graph.html + worked examples (onboarding/orient)

**Peer read:**

- [`worked/example/README.md`](../../../../../similar%20projects/graphify/worked/example/README.md) — reproducible corpus, `/graphify` hook, post-run questions.
- [`worked/rsl-siege-manager/README.md`](../../../../../similar%20projects/graphify/worked/rsl-siege-manager/README.md) — committed `graph.html`, `GRAPH_REPORT.md`, reproduction steps.
- [`exporters/html.py`](../../../../../similar%20projects/graphify/graphify/exporters/html.py) L30+ — sidebar search, legend, community viz (mechanism cite).
- Worked [`graph.html`](../../../../../similar%20projects/graphify/worked/rsl-siege-manager/graph.html) L67–68 — edge `confidence: EXTRACTED` in RAW_EDGES (spot-check).

**UA onboarding (H8 secondary peer):**

- [`onboard-builder.ts`](../../../../../similar%20projects/Understand-Anything/understand-anything-plugin/src/onboard-builder.ts) L7+ — graph-derived onboarding markdown.

**Trace contrast (TRACE-AUDIT H8 partial):**

- [`web/src/App.tsx`](../../../../../web/src/App.tsx) L21–22 — `/` → Graph; `/overview` → task loop.
- Observation: Trace has graph route but lacks committed peer-style worked artifact + install hook narrative.

**Verdict target:** `h8-gf-onboarding-ux.md`, optional `h8-ua-onboard.md`

### T5 — H8: Trace GUI observation (optional live)

**If GUI available:** screenshot or route list for Graph vs Overview — compare to GF `graph.html` affordances (search sidebar, community legend, committed artifact).

**If skipped:** Cite TRACE-AUDIT `h8-gui-partial.md` + App.tsx only.

**Verdict target:** fold into `h8-gf-onboarding-ux.md` — do not re-litigate S01.

### T6 — H4: Mempalace hybrid search (vector + BM25)

**Peer code read (required file:line):**

| Mechanism | Start here |
|-----------|------------|
| Hybrid design docstring | [`searcher.py`](../../../../../similar%20projects/mempalace/mempalace/searcher.py) L5–9 |
| `_hybrid_rank` | Same file L276–329 — vector_weight 0.6 / bm25_weight 0.4 |
| `search_memories` API | Same file L1652–1718 — query param, candidate_strategy, vector_disabled fallback |
| MCP entry | [`mcp_server.py`](../../../../../similar%20projects/mempalace/mempalace/mcp_server.py) L2369–2408 — `tool_search` → `search_memories` + `sanitize_query` |

**Trace contrast:**

- DR-NOSSEM in `doc.go`; TRACE-AUDIT H4 defer to S03 — state whether gap is **product law** vs **missing channel**.

**Verdict target:** `h4-mp-hybrid-search.md`

### T7 — H6 (MP slice): Mempalace MCP tool surface vs Trace 16-tool

**Peer code read:**

- [`service.py`](../../../../../similar%20projects/mempalace/mempalace/service.py) L60–102 — `READ_TOOLS`, `WRITE_TOOLS`, `MAINTENANCE_TOOLS` frozensets (count tools).
- [`mcp_server.py`](../../../../../similar%20projects/mempalace/mempalace/mcp_server.py) L4765+ — `TOOLS` dict registration + descriptions.
- Compare discovery: ranked "start here" vs flat list — observation only.

**Trace contrast (TRACE-AUDIT H6 — cite, do not re-run):**

- 16 tools; FM-08 paralysis plausible; CG single-tool in PEER-CG §3.

**Verdict target:** `h6-mp-mcp-surface.md` — **MP slice**; full H6 matrix deferred S04.

### T8 — H1/H8 (partial): Mempalace memory stack + context packet

**Peer code read:**

- [`layers.py`](../../../../../similar%20projects/mempalace/mempalace/layers.py) L3–17 — L0–L3 stack definition; wake-up token budget.
- Same file L404–431 — `wake_up()` (L0 identity + L1 essential story); `recall()` L2; `search()` L3.
- [`knowledge_graph.py`](../../../../../similar%20projects/mempalace/mempalace/knowledge_graph.py) L1–36 — temporal entity graph (SQLite).
- [`onboarding.py`](../../../../../similar%20projects/mempalace/mempalace/onboarding.py) L3–17 — first-run entity/wing seeding.

**Compare to Trace:**

- Compiler L0–1 only (TRACE-AUDIT H3); no memory/wake-up packet; task loop on Overview.

**Verdict target:** `h1-mp-context-packet.md`, `h8-mp-onboarding.md`

### T9 — H9 (MP contrast partial): fact_checker vs Trace intent pipeline

**Peer code read:**

- [`fact_checker.py`](../../../../../similar%20projects/mempalace/mempalace/fact_checker.py) L1–78 — `check_text`, relationship_mismatch, stale_fact classes.
- KG query backing: [`knowledge_graph.py`](../../../../../similar%20projects/mempalace/mempalace/knowledge_graph.py) — temporal triples.

**Trace contrast (TRACE-AUDIT H9 — cite only):**

- `rg intent internal/retrieval` zero; RETRIEVAL_AND_CONTEXT aspirational.

**Verdict target:** `h9-mp-fact-check-contrast.md` — **contrast slice**; S01 H9 verdict stands for Trace live.

### T10 — Moat seed: what UA / GF / MP lack

**Document (required §5):**

| Peer | Lacks (Trace moat) | Evidence |
|------|-------------------|----------|
| UA | Task loop, gate enforcement, evidence chain, plan tree | No task UUID in context-builder; read-only graph chat |
| Graphify | Same + progressive planning | Graph-only orient; no loop MCP |
| Mempalace | Task/plan/gate/evidence (has memory + hooks) | No Trace-style loop gate; memory-first not task-first |

**Verdict target:** `moat-peers-lack.md`

### T11 — Synthesize PEER-UA-GF.md

Author deliverable using T0–T10 evidence.

---

## Deliverable shape (PEER-UA-GF.md)

### §1 UA mechanism
Query-driven `ChatContext`; SearchEngine; onboard-builder; contrast Trace compiler/task packet.

### §2 Graphify mechanism
EXTRACTED/INFERRED pipeline; worked examples; `graph.html` UX; concept-query contrast.

### §3 Mempalace mechanism
Hybrid search; 4-layer memory stack; MCP surface; KG + fact_checker; onboarding.

### §4 Hypothesis verdicts

| H | Verdict | Primary peer evidence | Trace contrast | Confidence |
|---|---------|----------------------|----------------|------------|
| H1 partial | … | UA + MP rows | TRACE-AUDIT + PEER-CG | … |
| H4 | … | GF + MP | doc.go DR-NOSSEM | … |
| H6 (MP slice) | … | MP TOOLS | TRACE-AUDIT H6 | … |
| H8 | … | GF + UA + MP | TRACE-AUDIT H8 | … |
| H9 (MP contrast) | … | fact_checker | TRACE-AUDIT H9 | … |

Use INVESTIGATION-INDEX §2 **verified if / rejected if**.

### §5 Trace strengths peers lack (moat row seed)
Tasks, gates, evidence, plan tree — peers are graph/memory orient-first.

### §6 Evidence appendix + command log
Peer paths with line numbers; `$EV/` files per todo.

### §7 Spawn list
Empty or trigger + owner; H12+ only if uncovered slice found.

---

## Exit criteria

- [ ] PEER-UA-GF.md §§1–7 complete
- [ ] Every S03-owned H* row has verdict + **peer file:line** (README alone insufficient for T1–T4, T6–T9)
- [ ] Mempalace §3 present with hybrid search + layers + MCP cites
- [ ] Planner must-answer Q1–Q8 addressed
- [ ] No Go/TS/web product diff in Trace repo
- [ ] Board row P38-S03-01 → `done` with Notes (confidence + next P38-S03-02)

## Next

`P38-S03-02`
