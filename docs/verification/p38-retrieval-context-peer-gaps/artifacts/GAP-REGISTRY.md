# GAP-REGISTRY — Phase 38 cross-matrix synthesis

**Author:** P38-S04-01 (2026-08-22)  
**Status:** Investigation only — no product changes, no REMEDIATION-PLAN  
**Authority:** [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) §2 · [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)  
**Inputs (APPROVED):** [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md) · [PEER-CG.md](../scope-02-codegraph-peer/PEER-CG.md) · [PEER-UA-GF.md](../scope-03-ua-graphify-peer/PEER-UA-GF.md)  
**Evidence root (S04):** [`experiments/runs/2026-08-22-p38-s04-660/evidence/`](../../../../../experiments/runs/2026-08-22-p38-s04-660/evidence/)

---

## §1 Purpose

Single source of truth for **evidence-backed gaps**, **non-gaps**, and **deferrals** from Phase 38 S01–S03 investigations. Cross-matrix columns: **Trace | CG | UA | GF | MP | moat row**.

- Feeds **S05 saturation gate** — reviewer signs confident exit or spawns more S01–S04 rows.
- **Not** REMEDIATION-PLAN (S06) — no ranked G1–Gn themes or build tasks here.
- **Severity** = investigation confidence only (`high` | `medium` | `low`) — lowest confidence across contributing scopes per row.

---

## §2 Gap register (cross-matrix)

### §2.1 G-ID ↔ H* mapping (locked)

| Gap ID | H* | Theme (short) |
|--------|-----|---------------|
| **G-001** | H1 | Unified query+task orient packet |
| **G-002** | H2 | Compiler FTS title-only |
| **G-003** | H3 | Layers 2–3 designed not shipped |
| **G-004** | H4 | Semantic/concept retrieval (+ DR-NOSSEM defer sub-row) |
| **G-005** | H5 | Index langs + manual vs watcher |
| **G-006** | H6 | MCP discovery surface (16 / 1 / 44) |
| **G-007** | H7 | `trace_explore` unified read missing |
| **G-008** | H8 | Graph-first onboarding hook |
| **G-009** | H9 | Intent pipeline doc-only |
| **G-010** | H10 | Moat under-promoted in install/harness |
| **G-011** | H11 | Trace+CG dual-stack undocumented |
| **M-001** | moat | Trace strengths peers lack (**non-gap**) |

### §2.2 Main matrix

| Gap ID | H* | Trace | CG | UA | GF | MP | Verdict | Severity | Law fit | Evidence (dual-side) |
|--------|-----|-------|----|----|-----|-----|---------|----------|---------|----------------------|
| **G-001** | H1 | No `query` on `trace_context`; agents compose `trace_search` + `trace_context` — hits not auto-merged (`compiler.go`; MCP schema) | `codegraph_explore`: required `query`, single-call subgraph + source + call path + blast (`tools.ts` L1168–1170, L3193+) | `buildChatContext(query)` → SearchEngine + 1-hop expand (`context-builder.ts` L25–79) | N/A — graph orient separate from task backlog | `wake_up()` L0+L1 identity packet (`layers.py` L404–431) | **gap** | **high** | G7 — merge query+task, not replace packet | Trace: [h1-trace-partial.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h1-trace-partial.md) · CG: [h1-cg-partial.md](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h1-cg-partial.md) · UA/MP: [h1-ua-partial.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h1-ua-partial.md) |
| **G-002** | H2 | Compiler FTS inside packets uses **`task.Title` only** (`compiler.go` L146–151); CLI `trace search` separate | N/A — CG orients via explore `query`, not task-title FTS | SearchEngine on agent query: name/tags/summary Fuse keys (`search.ts` L14–58) | N/A | N/A | **gap** | **high** | OK | Trace: [h2-compiler-fts.txt](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h2-compiler-fts.txt) · UA: [h1-ua-search-mechanism.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h1-ua-search-mechanism.md) |
| **G-003** | H3 | Layers 0–1 shipped; L2–3 absent from live JSON; `doc.go` L7 explicit defer | N/A | N/A | N/A | 4-layer stack; L2 on-demand / L3 deep search designed (`layers.py` L3–17) | **gap** | **high** | Honest defer — still designed-vs-shipped | Trace: [h3-layers-designed-vs-shipped.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h3-layers-designed-vs-shipped.md) · MP: PEER-UA-GF §3 |
| **G-004** | H4 | DR-NOSSEM: FTS+expand only; no semantic channel (`retrieval/doc.go` L8–9) | N/A | Optional embedding when present (`embedding-search.ts`) | EXTRACTED/INFERRED edges + GRAPH_REPORT 90/10% (`symbol_resolution.py` L289–370) | `_hybrid_rank` BM25+vector L276–329 | **gap** (+ defer sub-row) | **high** | Vector **defer**; graph-label channel **gap** | Trace: [h4-trace-search-sample.txt](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h4-trace-search-sample.txt) · GF: [h4-gf-extracted-inferred.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h4-gf-extracted-inferred.md) · MP: [h4-mp-hybrid-search.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h4-mp-hybrid-search.md) |
| **G-005** | H5 | 5 lang IDs; manual `trace index`; `hook_installed: false` | Watcher debounce 300ms; 29 extractors vs README 37 langs | N/A | N/A | N/A | **gap** | **high** | Manual index OK by law; material delta | Trace: [h5-index-langs.txt](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h5-index-langs.txt) · CG: [h5-index-watch-contrast.md](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h5-index-watch-contrast.md) |
| **G-006** | H6 | 16 MCP tools; no ranked “start here”; 9/16 visible in Cursor session | **1** default tool + `SERVER_INSTRUCTIONS` (“One tool: codegraph_explore”) | N/A | N/A | **35** categorized / **44** TOOLS handlers — opposite extreme, no orient ranking | **gap** | **high** | Do not copy MP tool count | Trace: [h6-mcp-surface.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h6-mcp-surface.md) · CG: [h6-single-tool-ux.md](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h6-single-tool-ux.md) · MP: [h6-mp-mcp-surface.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h6-mp-mcp-surface.md) |
| **G-007** | H7 | No `trace_explore`; compose `trace_search`+`trace_why`+`trace_context`(+`trace_impact`) across 16 tools | `codegraph_explore` unified read+blast in one capped call; P24 consolidation **still deferred** | N/A | N/A | N/A | **gap** | **high** | Compose-equivalence **untested** — see §6 | Trace: PEER-CG `h7-explore-gap.md` contrast · CG: [h7-explore-mechanism.md](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h7-explore-mechanism.md) |
| **G-008** | H8 | `/` → Graph route (`App.tsx` L21–22); no committed graph.html, install hook, or confidence-labeled orient UX | N/A (read-only MCP) | `onboard-builder.ts` L7+ standalone markdown artifact | `worked/rsl-siege-manager/graph.html` + `/graphify` hook | `onboarding.py` + `wake_up()` identity story | **gap** | **high** | G19 GUI is adapter; orient UX gap | Trace: [h8-gui-partial.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h8-gui-partial.md) · GF/UA/MP: [h8-gf-onboarding-ux.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h8-gf-onboarding-ux.md), [h8-ua-onboard.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h8-ua-onboard.md), [h8-mp-onboarding.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h8-mp-onboarding.md) |
| **G-009** | H9 | RETRIEVAL_AND_CONTEXT §3 intent pipeline — **zero** code in `internal/retrieval/` | N/A | N/A | N/A | `fact_checker.py` L55–78 shipped offline check pipeline | **gap** | **high** | Doc-vs-shipped | Trace: [h9-intent-pipeline.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h9-intent-pipeline.md) · MP: [h9-mp-fact-check-contrast.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h9-mp-fact-check-contrast.md) |
| **G-010** | H10 | `trace install detect` / AGENTS.md lead with board protocol — task/gate/evidence moat buried | CG README graph-first agent workflow | UA graph-chat orient | Graphify install hook narrative | MP memory-first MCP | **gap** | **medium** | Harness/docs — moat exists (M-001) | Trace: [h10-install-moat.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h10-install-moat.md) · Peers: PEER-CG §1, PEER-UA-GF §1/§2 intros |
| **G-011** | H11 | CONTRIBUTING/README/AGENTS — Trace-only workflow; no dual-index recipe | PEER-FIXTURES lists CG as investigation fixture only — not user dogfood doc | N/A | N/A | N/A | **gap** | **high** | G19 boundaries undocumented for complementary stack | Trace: [h11-stack-docs.md](../../../../../experiments/runs/2026-08-22-p38-s04-660/evidence/h11-stack-docs.md) · CG: PEER-CG §5 complement note |

### §2.3 G-004 DR-NOSSEM sub-table

| Sub-ID | Channel | Trace | Peer mechanism | Verdict | Severity | Law fit |
|--------|---------|-------|----------------|---------|----------|---------|
| **G-004a** | Embedding / vector semantic | Forbidden (`doc.go` L8–9) | MP vector leg in `_hybrid_rank`; UA optional embeddings | **defer** | **high** | DR-NOSSEM — correctly out of scope |
| **G-004b** | Label / summary / concept via graph | Title-token FTS + structural expand only | GF EXTRACTED/INFERRED edges; MP BM25 over memory text | **gap** | **high** | Non-semantic graph channels OK |

---

## §3 Moat row (M-001)

**Verdict:** **non-gap** — Trace strengths peers lack.

| ID | Trace | CG | UA | GF | MP |
|----|-------|----|----|----|-----|
| **M-001** | **Ships:** progressive task packet (L0–1), task loop + gate + review, enforcement harness, `trace_why`, plan tree, local-first caps, explicit defer honesty | **Lacks:** task UUID, gates, evidence, backlog, plan tree — read-only graph | **Lacks:** task loop, enforcement, evidence chain | **Lacks:** same; graph-only orient | **Lacks:** task/plan/gate; `wake_up()` is identity not backlog |

**Expanded strengths (deduped from TRACE-AUDIT §5 + PEER-CG §5 + PEER-UA-GF §5):**

1. Progressive bounded context packet with budget, reason_codes, trust labels  
2. Task loop + gate + review MCP surface (`trace_loop`, `trace_review`, `trace_transition`)  
3. Enforcement harness (`TRACE_TASK_ID`, loop gate, plan bootstrap)  
4. Causal `trace_why` — discovery→task chains  
5. Local-first progressive retrieval (Laws 6–7) — no full-graph dump defaults  
6. Plan tree + orchestration (`trace_plan`, `trace_tasks`)  
7. Explicit layer/defer documentation (`doc.go` L7) vs silent omission  

**Distinct from G-010:** M-001 confirms moat **exists**; G-010 is moat **under-promotion** in first-run messaging.

**Confidence:** **high**

---

## §4 Non-gaps (peers weaker)

Explicit list where peers are **weaker** — not Trace backlog items:

| Peer weakness | Why not a Trace gap | Evidence |
|---------------|---------------------|----------|
| CG / UA / GF / MP: no task loop, gates, evidence, plan tree | Trace moat (M-001) | PEER-CG §5; PEER-UA-GF §5 |
| CG: graph-only product (no directed work) | Abandons Trace core value | PEER-CG §4 |
| CG: always-on MCP daemon | Violates local-first / P24 anti-pattern | `daemon.ts` L1–25; PEER-CG §4 |
| CG: MCP-only core loop | Trace CLI+library remain authoritative | PEER-CG §4 |
| MP: 44-tool surface as remediation target | Opposite UX failure mode; no orient ranking | G-006; PEER-UA-GF §3 |

---

## §5 Deferred (law / policy / explicit design)

| Item | H/G | Verdict | Rationale |
|------|-----|---------|-----------|
| Embedding / vector semantic channel | H4 / G-004a | **defer** | DR-NOSSEM (`retrieval/doc.go` L8–9) — product law, not backlog gap |
| Replace task packet with query-only orient | H1 | **reject** | Gap is *merge* query+task, not discard task scope (PEER-CG §4) |
| CG detached daemon as P0 path | — | **reject** | Law/local-first (PEER-CG §4) |
| CG benchmark % claims as remediation proof | — | **reject** | Harness/marketing only (`h6-benchmark-claims.md`) |
| Layers 2–3 explicit defer in `doc.go` L7 | H3 / G-003 | **gap** (not defer) | Honest deferral is strength; designed-vs-shipped still verified gap per INVESTIGATION-INDEX |
| P24 `trace_explore` implement defer | H7 | **gap** until rejected | Absence is gap (G-007); compose-equivalence could reject — no evidence yet |

**PEER-CG §4 anti-patterns** (full list): [`anti-patterns-not-for-trace.md`](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/anti-patterns-not-for-trace.md)

---

## §6 Spawn list → S05

| Trigger | Condition | Owner | Status |
|---------|-----------|-------|--------|
| **H7 compose-equivalence** | No live proof that Trace multi-tool compose ≈ `codegraph_explore` | S01 or S02 live test | **Open** — documented for S05/S06; not blocking registry completeness |
| H11 inconclusive doc read | — | — | **Closed** — G-011 gap verified (high) |
| Dual-side missing on gap row | — | — | **Closed** — all G-001…G-011 have Trace + peer cites |
| H12+ uncovered peer slice | — | — | **Closed** — MP mapped to H1/H4/H6/H8/H9 |
| Cursor MCP 9/16 tool exposure | Harness stale-server risk | S06 plan / install docs | **Fold** — harness hygiene |
| Optional S01-01a symbol evidence | Richer symbol/file in packet | S01-01a | **Defer** — reviewer optional |

**S05 note:** Registry is **complete** for saturation forward-fit (DESIGN-LOCKS § Saturation exit). One **open trigger** (H7 compose-equivalence) may spawn before REMEDIATION-PLAN if reviewer requires equivalence proof.

---

## §7 Evidence index

| Scope | Board row | Evidence root | Primary artifact |
|-------|-----------|---------------|------------------|
| S01 Trace live | P38-S01-01 (651) | [`…-p38-s01-651/evidence/`](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/) | TRACE-AUDIT.md |
| S02 Codegraph | P38-S02-01 (654) | [`…-p38-s02-654/evidence/`](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/) | PEER-CG.md |
| S03 UA+GF+MP | P38-S03-01 (657) | [`…-p38-s03-657/evidence/`](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/) | PEER-UA-GF.md |
| S04 cross-matrix | P38-S04-01 (660) | [`…-p38-s04-660/evidence/`](../../../../../experiments/runs/2026-08-22-p38-s04-660/evidence/) | This file |

**S04 synthesis files:** `t1-trace-column-seed.md` · `t2-cg-column-seed.md` · `t3-ua-gf-mp-column-seed.md` · `t4-gap-id-registry.md` · `t5-matrix-h1-h4.md` · `t6-matrix-h5-h7.md` · `t7-matrix-h8-h10.md` · `h11-stack-docs.md` · `t9-moat-row-m001.md` · `t10-non-gaps-deferrals.md` · `t11-spawn-triggers.md`

No duplicate JSON captures — link S01–S03 paths only.

---

## Summary counts

| Category | Count | Confidence breakdown |
|----------|-------|----------------------|
| **Gaps (G-001…G-011)** | **11** | high: **10** · medium: **1** (G-010) |
| **Defer sub-row (G-004a)** | **1** | high |
| **Moat (M-001)** | **1** non-gap | high |
| **Peer-weaker non-gaps (§4)** | **5** items | high |
| **Law/policy deferrals (§5)** | **1** primary (vector) + rejects | high |
| **Open S05 spawn trigger** | **1** (H7 compose-equivalence) | — |

**Overall registry confidence:** **high** — all rows dual-side cited; H8 upgraded S01→S03; H11 doc-read complete.

**Next board row:** P38-S04-02 (review)
