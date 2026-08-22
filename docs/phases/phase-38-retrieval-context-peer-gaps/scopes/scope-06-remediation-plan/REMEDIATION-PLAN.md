# REMEDIATION-PLAN — Phase 38 retrieval & context peer gaps

**Author:** P38-S06-01 (2026-08-22)  
**Status:** Planning artifact only — **no product code, no implement in P38**  
**Authority:** [SATURATION-NOTES.md](../scope-05-saturation-gate/SATURATION-NOTES.md) (APPROVE saturated) · [GAP-REGISTRY.md](../scope-04-gap-registry/GAP-REGISTRY.md) · [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)  
**Evidence root (S06):** [`experiments/runs/2026-08-22-p38-s06-666/evidence/`](../../../../../experiments/runs/2026-08-22-p38-s06-666/evidence/)  
**Upstream evidence:** [GAP-REGISTRY §7](../scope-04-gap-registry/GAP-REGISTRY.md) — link S01–S05 `$EV/` paths only; no duplicate JSON

---

## §1 Executive summary

Phase 38 investigation is **saturated** (S05 6/6 PASS, `ready_for_REMEDIATION_PLAN: true`). Eleven evidence-backed gaps (G-001…G-011) consolidate into **nine ranked remediation themes (G1–G9)** for human-promoted **Phase 39+** implementation.

**Trace moat (M-001) preserved:** Progressive task packet, task loop + gate + review, enforcement harness, causal `trace_why`, plan tree, local-first progressive caps (Laws 6–7), and explicit defer honesty remain **non-negotiable strengths** peers lack. Remediation merges peer read-orient patterns **into** the task moat — never replaces it.

**Phase 39 recommendation (human promotion):** Entry co-wave **G1 + G3 + G4** — query+task context merge, MCP/harness orient playbook, and doc-only Trace+Codegraph dual-stack recipe. Compose-first read-surface UX (G2 primary path) ships in Phase 39 docs/harness; unified `trace_explore` deferred to Phase 40+ after G1 and law review.

**Confidence:** **high** — synthesis of APPROVED GAP-REGISTRY (row 661) and SATURATION-NOTES (row 664); no new gap discovery required.

---

## §2 Ranked themes G1–G9

### Summary table

| Rank | Theme | GAP ids | impact | law_fit | effort | score | Phase sketch |
|------|-------|---------|--------|---------|--------|-------|--------------|
| 1 | **G1** Query+task orient merge | G-001, G-002 | 5 | 5 | 3 | 8.33 | Phase 39 |
| 2 | **G3** MCP discovery & harness orient | G-006, G-010 | 5 | 5 | 2 | 12.50 | Phase 39 |
| 3 | **G4** Dual-stack documentation (H11) | G-011 | 4 | 5 | 1 | 20.00 | Phase 39 docs |
| 4 | **G5** Graph-first onboarding UX | G-008 | 4 | 4 | 3 | 5.33 | Phase 39–40 |
| 5 | **G2** Read-surface strategy (H7) | G-007 | 5 | 4 | 4 | 5.00 | Phase 39 compose / 40+ explore |
| 6 | **G6** Non-semantic concept retrieval | G-004b | 4 | 3 | 4 | 3.00 | Phase 40+ |
| 7 | **G7** Index freshness & language coverage | G-005 | 3 | 4 | 4 | 3.00 | Phase 40+ |
| 8 | **G8** Progressive layers L2–L3 | G-003 | 4 | 4 | 4 | 4.00 | Phase 41+ |
| 9 | **G9** Intent pipeline | G-009 | 3 | 3 | 5 | 1.80 | Phase 41+ or doc-revise |

**Rank order lock:** G1 → G3 → G4 → G5 → G2 → G6 → G7 → G8 → G9.

**G1/G3 tie-breaker note:** G3 raw score (12.50) exceeds G1 (8.33) on lower effort, but G1 covers more agent-blocking gaps (G-001+G-002 cluster). **G1 ranks first** per planner lock and tie-breaker #1 (gap coverage).

**Phase 39 co-wave:** G1 + G3 + G4 (human promotion).

---

### G1 — Query+task orient merge

| Field | Detail |
|-------|--------|
| **Problem** | Agents cannot merge agent query + task packet in one orient step; compiler FTS uses `task.Title` only; `trace search` hits do not auto-merge into context packets. |
| **GAP ids** | G-001, G-002 |
| **Peer pattern** | UA `buildChatContext(query)` → SearchEngine + 1-hop expand (`context-builder.ts` L25–79); MP `wake_up()` L0+L1 identity packet (`layers.py` L404–431); CG query on `codegraph_explore` (contrast — task-agnostic). |
| **Phase sketch** | **Phase 39 — Context orient merge:** optional `query` on `trace_context` / compiler path; merge search hits into packet; preserve task UUID + gates + reason_codes. |
| **Risks** | Over-merge → dump risk (Law 6); query-only drift abandons moat (M-001). |
| **Not P38** | No compiler/MCP/Go code in P38. |

**Evidence:** [h1-trace-partial.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h1-trace-partial.md) · [h2-compiler-fts.txt](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h2-compiler-fts.txt) · GAP-REGISTRY §2.2 G-001/G-002

---

### G3 — MCP discovery & harness orient

| Field | Detail |
|-------|--------|
| **Problem** | 16 MCP tools without ranked “start here”; moat under-promoted vs CG 1-tool + MP 44-tool extremes; Cursor session saw 9/16 tools (harness hygiene). |
| **GAP ids** | G-006, G-010 (+ harness 9/16 fold) |
| **Peer pattern** | CG `DEFAULT_MCP_TOOLS = ['explore']` + `SERVER_INSTRUCTIONS` (“One tool: codegraph_explore”); MP categorized READ sets (contrast — do not copy tool count). |
| **Phase sketch** | **Phase 39 — Harness orient:** MCP init playbook, install/bootstrap moat-first messaging (`trace_loop`/`trace_review` lead), Cursor 9/16 registration hygiene doc. |
| **Risks** | Copying MP 44-tool surface; hiding task/write tools behind read-only facade. |
| **Not P38** | No MCP `server.go` changes in P38. |

**Evidence:** [h6-mcp-surface.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h6-mcp-surface.md) · [h6-single-tool-ux.md](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h6-single-tool-ux.md) · [h10-install-moat.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h10-install-moat.md)

---

### G4 — Dual-stack documentation (H11)

| Field | Detail |
|-------|--------|
| **Problem** | No user-facing doc for Trace + Codegraph complementary workflow; CONTRIBUTING/README/AGENTS are Trace-only. |
| **GAP ids** | G-011 |
| **Peer pattern** | PEER-CG §5 complement note; PEER-FIXTURES fixture paths; [`h11-stack-docs.md`](../../../../../experiments/runs/2026-08-22-p38-s04-660/evidence/h11-stack-docs.md). |
| **Phase sketch** | **Phase 39 — Docs only:** CONTRIBUTING/AGENTS section — when to `codegraph init` vs `trace index`, Law 19 adapter boundaries, storage separation (`.trace/` vs `.codegraph/`). |
| **Risks** | Product dual-index default → complexity; adapter logic fork in core. |
| **Not P38** | **Doc-only** — investigation conclusion: not product integration. No bundled MCP or default dual-index in P38 or Phase 39 product code. |

**H11 verdict (LOCK):** Doc-only dual-stack recipe **Accept**; product integration (default dual-index, bundled MCP) **Reject**.

---

### G5 — Graph-first onboarding UX

| Field | Detail |
|-------|--------|
| **Problem** | GUI `/` → Graph route without committed orient artifact, install hook, or confidence-labeled orient UX. |
| **GAP ids** | G-008 |
| **Peer pattern** | GF `worked/rsl-siege-manager/graph.html` + `/graphify` hook; UA `onboard-builder.ts`; MP `onboarding.py` + `wake_up()` identity story. |
| **Phase sketch** | **Phase 39–40 — GUI orient adapter:** graph route content, optional static orient artifact, install hook narrative (Law 19 — adapter calls library, no business logic in `web/`). |
| **Risks** | Business logic in `web/`; full Graphify port; graph-only product drift. |
| **Not P38** | No `web/` product code in P38. |

**Evidence:** [h8-gui-partial.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h8-gui-partial.md) · [h8-gf-onboarding-ux.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h8-gf-onboarding-ux.md)

---

### G2 — Read-surface strategy (H7 owner)

| Field | Detail |
|-------|--------|
| **Problem** | No unified read-orient; desk-check **not equivalent** to CG `codegraph_explore` on 7/7 dimensions (single call, query, verbatim source, call path, blast, task-merge, tool-count). |
| **GAP ids** | G-007 (+ G-006 partial via discovery cost) |
| **Peer pattern** | CG `codegraph_explore` single capped call (`tools.ts` L1163–1181, L3193+); Trace 16-tool split compose. |
| **Phase sketch** | **Phase 39 (compose-first):** SERVER_INSTRUCTIONS-style orient recipe, ranked read tools, optional orchestration doc — **before** **Phase 40+ (unified `trace_explore`):** task-aware capped explore after G1 + law spike. |
| **Risks** | Unified tool without task scope → H1 anti-pattern; mega-tool hides write surface; query-only trap. |
| **Not P38** | No `trace_explore` implement in P38; live spike optional Phase 39 pre-gate only. |

**H7 owner decision (LOCK — SATURATION-NOTES §4 defer owner = S06):**

| Option | Rank | Rationale |
|--------|------|-----------|
| **Compose-first UX** (orient recipe + G1 merge + ranked read tools) | **Primary — Phase 39** | Preserves 16-tool write surface + task moat (M-001); lower effort; addresses G-006 discovery without mega-tool |
| **Unified `trace_explore`** (task-aware, capped, P24 transfer) | **Secondary — Phase 40+** | Requires G1 query+task merge first; law review + optional live spike; CG parity without query-only trap |

**Reject for remediation:** Claiming Trace multi-tool compose already equivalent to CG explore ([h7-compose-desk-check.md](../../../../../experiments/runs/2026-08-22-p38-s05-663/evidence/h7-compose-desk-check.md) closed that path).

---

### G6 — Non-semantic concept retrieval

| Field | Detail |
|-------|--------|
| **Problem** | Title-token FTS misses concept/summary graph channels; DR-NOSSEM forbids vector leg. |
| **GAP ids** | G-004b only (**G-004a vector → §4 defer, not theme**) |
| **Peer pattern** | GF EXTRACTED/INFERRED edges + GRAPH_REPORT; MP BM25 text leg in `_hybrid_rank` (not vector). |
| **Phase sketch** | **Phase 40+ — Graph-label retrieval:** summary/label channel under DR-NOSSEM; law review gate before build. |
| **Risks** | DR-NOSSEM slip into embeddings (G-004a); semantic creep; packet bloat. |
| **Not P38** | G-004a vector **reject** — not remediation theme. No retrieval implement in P38. |

**Evidence:** [h4-gf-extracted-inferred.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h4-gf-extracted-inferred.md) · GAP-REGISTRY §2.3

---

### G7 — Index freshness & language coverage

| Field | Detail |
|-------|--------|
| **Problem** | 5 language IDs; manual `trace index`; `hook_installed: false`; no watcher. |
| **GAP ids** | G-005 |
| **Peer pattern** | CG watcher debounce 300ms + 29 extractors (study idea — not daemon stack). |
| **Phase sketch** | **Phase 40+ — Index ergonomics:** lang expansion policy; optional watch/hook path (local-first, no always-on daemon). |
| **Risks** | CG detached daemon anti-pattern; lang sprawl without policy; full-rebuild-on-change architecture. |
| **Not P38** | No analyzer/index code in P38. |

**Evidence:** [h5-index-langs.txt](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h5-index-langs.txt) · [h5-index-watch-contrast.md](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h5-index-watch-contrast.md)

---

### G8 — Progressive layers L2–L3

| Field | Detail |
|-------|--------|
| **Problem** | Layers 2–3 designed (`doc.go` L7) not in live packet; `--depth 2` expands graph within L0–1 only. |
| **GAP ids** | G-003 |
| **Peer pattern** | MP 4-layer stack — L2 on-demand / L3 deep search designed (`layers.py` L3–17). |
| **Phase sketch** | **Phase 41+ — Layer expansion:** ship L2–L3 in compiler or revise spec with documented alternative. |
| **Risks** | Packet bloat; layer honesty regression; silent omission vs explicit defer strength. |
| **Not P38** | No compiler layer ship in P38. |

**Evidence:** [h3-layers-designed-vs-shipped.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h3-layers-designed-vs-shipped.md)

---

### G9 — Intent pipeline

| Field | Detail |
|-------|--------|
| **Problem** | RETRIEVAL_AND_CONTEXT §3 intent pipeline — zero code in `internal/retrieval/`. |
| **GAP ids** | G-009 |
| **Peer pattern** | MP `fact_checker.py` L55–78 (contrast — offline check pipeline, not direct port). |
| **Phase sketch** | **Phase 41+ or doc-revise:** implement intent extraction **or** mark doc aspirational + supersede. |
| **Risks** | Large scope; overlaps G1/G6; law review for new retrieval channel. |
| **Not P38** | No retrieval implement in P38. |

**Evidence:** [h9-intent-pipeline.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h9-intent-pipeline.md) · [h9-mp-fact-check-contrast.md](../../../../../experiments/runs/2026-08-22-p38-s03-657/evidence/h9-mp-fact-check-contrast.md)

---

## §3 Proposed future phase sketches

**No board rows spawned in P38.** Titles + scope bullets only — human promotes Phase 39+.

### Phase 39 — Context orient & harness (entry co-wave: G1 + G3 + G4)

- **G1 Context orient merge** — optional query param; merge FTS/search into task packet; preserve gates
- **G3 Harness orient** — MCP playbook, moat-first bootstrap, 9/16 Cursor hygiene
- **G4 Dual-stack docs** — CONTRIBUTING/AGENTS Trace+CG recipe (doc-only)
- **G5 start** — graph route orient narrative; static artifact sketch
- **G2 compose-first** — orient recipe doc; ranked read-tool sequence (no unified explore yet)

### Phase 40+ — Read surface & retrieval depth

- **G2 unified `trace_explore`** — task-aware, capped, after G1 + law spike (secondary H7 path)
- **G6 Graph-label retrieval** — G-004b non-semantic channel
- **G7 Index ergonomics** — lang policy; optional local watch/hook
- **G5 complete** — GUI orient adapter; install hook narrative

### Phase 41+ — Layers & intent

- **G8 Layer expansion** — L2–L3 ship or spec revise
- **G9 Intent pipeline** — implement or doc-revise aspirational §3

---

## §4 Reject / defer registry

Minimum 12 rejects. Sources: PEER-CG §4, SATURATION-NOTES §5, GAP-REGISTRY §4–§5, plan-specific.

| # | Reject / defer | Rationale |
|---|----------------|-----------|
| 1 | **CG detached MCP daemon as P0** | `daemon.ts` L1–25; P24 §4 — no always-on daemon on core path |
| 2 | **MCP-only core loop** | CG install centers MCP; Trace CLI+library remain authoritative |
| 3 | **Full-graph dump defaults** | Violates Law 6; progressive caps required (Law 7) |
| 4 | **Graph-only product direction** | Abandons task loop, gates, evidence moat (M-001) |
| 5 | **Query-only replaces task packet** | Gap is *merge* query+task, not discard task scope (G-001) |
| 6 | **CG benchmark % as remediation proof** | Marketing/harness claims; mechanism study only in P38 |
| 7 | **Copy MP 44-tool MCP surface** | Opposite UX failure; no orient ranking (G-006) |
| 8 | **Embedding/vector semantic channel (G-004a)** | DR-NOSSEM — product law defer, not backlog gap |
| 9 | **Implement any remediation in P38** | DESIGN-LOCKS investigation-only phase |
| 10 | **Product default dual-index integration** | H11 investigation — blurs moat; doc-only accepted |
| 11 | **Claim compose ≈ CG explore** | h7-compose-desk-check: 7/7 dimensions **not equivalent** |
| 12 | **Always-on network daemon** | Local-first law; CG daemon is study-only anti-pattern |
| 13 | **Implement `trace_explore` during P38** | Out of scope — investigation-only |
| 14 | **Re-audit S01–S03 live CLI wave** | Duplicate — saturated registry |
| 15 | **Semantic embedding spike** | Law conflict — G-004a explicit defer |

### Defer (not theme)

| Item | Verdict | Owner |
|------|---------|-------|
| **G-004a** vector/embedding | **defer** | Phase 41+ law review only if policy changes |
| **M-001** task moat | **non-gap** | §1 executive — preserved, not remediated |

---

## §5 Open questions for human promotion

1. **Phase 39 scope cut** — Can G5 (GUI start) slip to Phase 40 if G1+G3+G4 co-wave is capacity-bound?
2. **G2 spike gate** — Is a Phase 39 pre-implement live MCP spike required before Phase 40 `trace_explore` board spawn?
3. **G9 implement vs doc-revise** — Does RETRIEVAL_AND_CONTEXT §3 remain roadmap or become superseded aspirational doc?
4. **G7 lang policy** — Which languages beyond go/js/ts/tsx/py justify analyzer investment vs defer?
5. **Harness 9/16** — Is Cursor stale-server fix docs-only (G3) or requires MCP registration code change in Phase 39?

---

## §6 Successor recommendation for S07 DR-HANDOFF

- **Human-promote Phase 39** with entry themes **G1 + G3 + G4** (context merge, harness orient, dual-stack docs).
- **Close Phase 38 investigation** at S07 VERIFY — no implement rows in P38 board history.
- **Preserve M-001** in Phase 39 charter — all remediation themes merge peer patterns into task moat, never replace it.
- **Forward artifacts:** This plan + GAP-REGISTRY + SATURATION-NOTES → Phase 39 INTAKE/planner (human-authored).
- **Next board row:** P38-S06-02 (review).

---

## Evidence index (S06)

| File | Task |
|------|------|
| [t0-preflight.md](../../../../../experiments/runs/2026-08-22-p38-s06-666/evidence/t0-preflight.md) | T0 |
| [t1-gap-theme-coverage.md](../../../../../experiments/runs/2026-08-22-p38-s06-666/evidence/t1-gap-theme-coverage.md) | T1 |
| [t2-ranking-scores.md](../../../../../experiments/runs/2026-08-22-p38-s06-666/evidence/t2-ranking-scores.md) | T2 |
| [t3-reject-registry.md](../../../../../experiments/runs/2026-08-22-p38-s06-666/evidence/t3-reject-registry.md) | T3 |
| [t4-phase-sketches.md](../../../../../experiments/runs/2026-08-22-p38-s06-666/evidence/t4-phase-sketches.md) | T4 |
| [t5-h7-h11-decisions.md](../../../../../experiments/runs/2026-08-22-p38-s06-666/evidence/t5-h7-h11-decisions.md) | T5 |

**Upstream (read-only):** [GAP-REGISTRY §7](../scope-04-gap-registry/GAP-REGISTRY.md) · [SATURATION-NOTES §7](../scope-05-saturation-gate/SATURATION-NOTES.md)

---

## T7 self-check

- [x] Every G-001…G-011 in theme or §4
- [x] No "implement in P38" language
- [x] No Go/TS/web diff
- [x] Rubric table present (§2 summary)
- [x] H7 + H11 locked decisions documented (G2, G4, §5)
- [x] Links to GAP-REGISTRY evidence paths (not duplicate JSON)
- [x] ≥12 rejects in §4 (15 listed)
- [x] M-001 moat in §1 executive
