# PEER-CG — Codegraph peer investigation

**Author:** P38-S02-01 (2026-08-22)  
**Status:** Investigation only — no product changes  
**Authority:** [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) §2 (H1 partial, H5, H6, H7)  
**Trace baseline:** [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md) (S01 APPROVED)  
**Evidence root:** [`experiments/runs/2026-08-22-p38-s02-654/evidence/`](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/)  
**Peer root:** [`similar projects/codegraph/`](../../../../../similar%20projects/codegraph/) (read-only)

---

## §1 Executive summary

Investigation-only peer read of Codegraph (CG) against Trace retrieval/MCP surfaces. No product code changed; no live CG index on Trace dogfood or peer clone (T8 skipped — file:line + MCP schema suffice).

**Headline contrasts:**

1. **`codegraph_explore` is a unified read-orient tool** — one capped call returns query-scoped subgraph, verbatim line-numbered source grouped by file, call paths (including dynamic-dispatch hops), and blast-radius locations (`tools.ts` L1163–1181, L2775+, L3017+, L3193+). Trace has no equivalent; agents compose `trace_search` + `trace_context` + `trace_why` (+ optional `trace_impact`) across 16 MCP tools.
2. **CG defaults to a single listed MCP tool** — `DEFAULT_MCP_TOOLS = ['explore']` (`tools.ts` L1275–1286) plus aggressive server instructions (`server-instructions.ts` L34–50). Trace exposes all 16 tools without ranked “start here” (TRACE-AUDIT H6).
3. **CG index freshness is watcher-driven** — debounced auto-sync (`watcher.ts` L1–80, README L125–127) vs Trace manual `trace index` and `hook_installed: false` (TRACE-AUDIT H5).
4. **CG is query-driven, task-agnostic** — `query` string orients; no task UUID or backlog. Trace task packet is a strength but lacks query merge in one call (H1 partial).
5. **P24 “consolidate toward trace_explore” remains deferred** — zero Go symbol; Phase 37 shipped INT-06 reorder only, not read-surface consolidation (P24 EXTERNAL-RESEARCH CG row; evidence `h7-p24-deferred-evidence.md`).

**Confidence:** H5, H6, H7 — **high** (mechanism cites + S01 baseline). H1 partial — **high** on structural gap; full cross-peer matrix deferred S04.

---

## §2 Comparison table vs Trace

Mechanism-focused (not feature checklist). Trace column cites TRACE-AUDIT unless noted.

| Dimension | Codegraph mechanism (file:line) | Trace mechanism (file:line or TRACE-AUDIT) | Notes |
|-----------|--------------------------------|--------------------------------------------|-------|
| **Explore / orient** | `codegraph_explore`: `query` (required), `maxFiles`, `projectPath` → `handleExplore` → `findRelevantContext` + file sections + flow + blast (`tools.ts` L1163–1181, L2098, L3193+, L2775+, L3017+) | No `trace_explore`; `trace_context` (`task_id` only), `trace_search`, `trace_why`, `trace_impact` (`server.go` L68, L164, L230–233); TRACE-AUDIT §4 | CG subsumes read-orient in one call; Trace multi-tool compose |
| **Query vs task input** | Free-form `query` — NL or symbol bag; no task scope (`tools.ts` L1168–1170) | `trace_context` requires `task_id`; compiler FTS uses `task.Title` (TRACE-AUDIT H1/H2) | Structural H1 gap on Trace side; CG read-only |
| **Index freshness** | `codegraph init` full index (README L112–117); watcher debounce 300ms quick / scoped sync (`watcher.ts` L68–76) | Manual `trace index`; `hook_installed: false` (TRACE-AUDIT H5) | Material auto-sync delta |
| **Language coverage** | 29 extractor modules under `src/extraction/languages/`; README lists 37 language icons (L145–182) | 5 language IDs: go, js, ts, tsx, py (`language_adapter.go` L19–25) | H5 lang gap supported |
| **MCP surface** | Default list: **1 tool** — explore only (`tools.ts` L1275–1286); others via `CODEGRAPH_MCP_TOOLS` | **16 tools** — `RegisteredToolNames()` (`server.go` L228–235) | CG removes choice; Trace FM-08 friction |
| **Server steering** | `SERVER_INSTRUCTIONS` — “One tool: codegraph_explore” (`server-instructions.ts` L34–50) | No ranked orient recipe in MCP init (TRACE-AUDIT H6) | CG playbook every session |
| **Caps / progressive context** | Adaptive `getExploreOutputBudget` 13k–24k chars; global `MAX_OUTPUT_LENGTH` 15k for other tools (`tools.ts` L86–87, L214–276); blast locations only (`L3017–3073`) | Packet caps, layers 0–1 shipped (`compiler`, TRACE-AUDIT §4); Laws 6–7 | Both bounded; CG explore is dense single payload |
| **Task / moat** | **None** — read-only graph; no gates, evidence, or backlog | Task loop, gate, review, plan tree, enforcement (TRACE-AUDIT §5) | Trace strength — do not abandon for graph-only |
| **Daemon / shared process** | Optional detached daemon — shared watcher/SQLite (`daemon.ts` L1–25) | No daemon; CLI/MCP per invocation | **Anti-pattern for Trace P0-X** |

### Side-by-side orient flow (H1 partial)

| Step | Codegraph | Trace |
|------|-----------|-------|
| Agent question input | `query` on `codegraph_explore` | Separate `trace_search` (not in packet) |
| Task/plan scope | **None** (read-only graph) | Task packet L0 |
| Single-call orient | `codegraph_explore` | Partial — compose read tools + task context |

---

## §3 Gap hypotheses — verdicts

Aligned with [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) §2 verify/reject criteria.

| H | Verdict | Peer evidence | Trace contrast | Confidence |
|---|---------|---------------|----------------|------------|
| **H1 partial** | **supported** | Query-only orient via `codegraph_explore` (`tools.ts` L1168–1170; `h1-cg-partial.md`) | TRACE-AUDIT H1 — no `query` on `trace_context`; compose required | **high** |
| **H5** | **supported** | Watcher auto-sync (`watcher.ts` L1–80); 29 lang extractors vs README 37 icons (`h5-index-watch-contrast.md`) | TRACE-AUDIT H5 — 5 langs, manual index, no watcher | **high** |
| **H6** | **supported** | Single default MCP tool + instructions (`tools.ts` L1275–1286; `server-instructions.ts` L34–50; `h6-single-tool-ux.md`) | TRACE-AUDIT H6 — 16 tools, discovery friction | **high** |
| **H7** | **supported** | Unified explore mechanism (`h7-explore-mechanism.md`, `h7-explore-gap.md`); P24 transfer still open (`h7-p24-deferred-evidence.md`) | No `trace_explore` in Go; INT-06 reorder only | **high** |

**H7 nuance:** Trace *could* theoretically compose read tools with evidence — S04 must test equivalence. P38 S02 finds no shipped single-tool parity and P24 consolidation **still deferred**.

**Benchmark claims:** Classified observation-only for gap proof (`h6-benchmark-claims.md`) — not re-run in P38.

---

## §4 Anti-patterns Trace must not copy

See evidence [`anti-patterns-not-for-trace.md`](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/anti-patterns-not-for-trace.md).

| Anti-pattern | Why not |
|--------------|---------|
| **Detached MCP daemon as P0 requirement** | CG `daemon.ts` L1–25; P24 §4 — no daemon on P0-X core path |
| **MCP-only core loop** | CG install centers MCP; Trace keeps library+CLI authoritative |
| **Full-graph dump defaults** | Violates G Law 6 — CG itself caps explore; Trace must keep progressive packets (G Law 7) |
| **Graph-only product direction** | CG has no task loop, gate, or evidence — abandoning Trace moat |
| **Replacing task packet with query-only orient** | Gap is *merge* (query + task), not discard task scope (H1 partial) |
| **Using CG benchmark % as remediation justification** | Marketing/harness claims; mechanism study only in P38 |

---

## §5 Trace strengths peers lack (moat row seed)

Codegraph is **read-only code intelligence** — no executable backlog, state machine, or verification loop. Trace ships (TRACE-AUDIT §5):

| Strength | Trace evidence |
|----------|----------------|
| **Progressive task context packet** | Bounded JSON, budget, reason_codes, trust labels |
| **Task loop + gate + review** | `trace_loop`, `trace_review`, `trace_transition`; Overview GateStrip |
| **Enforcement harness** | `TRACE_TASK_ID`, loop gate, plan bootstrap |
| **Why / causal trace** | `trace_why` — discovery→task chains |
| **Local-first, no dump API** | Packet caps; Laws 6–7 in compiler |
| **Plan tree + agent orchestration** | `trace_plan`, `trace_tasks`, 16-tool write surface for loop discipline |

CG complements code orientation; Trace owns **directed work with evidence**. Dual-stack documentation (H11) is S04 scope.

---

## §6 Evidence appendix + live command log

### Evidence files

| File | Todo | Description |
|------|------|-------------|
| `t0-peer-root.txt` | T0 | Peer `tools.ts` present |
| `t0-no-codegraph-on-trace.txt` | T0 | No `.codegraph/` on Trace |
| `t0-cg-readme-lang-refs.txt` | T0 | README language section refs |
| `h7-explore-mechanism.md` | T1 | Explore inputs/outputs/caps |
| `h7-explore-gap.md` | T1 | Trace composition gap |
| `h6-single-tool-ux.md` | T2 | MCP surface comparison |
| `h5-index-watch-contrast.md` | T3 | Index/watch/lang table |
| `h1-cg-partial.md` | T4 | Query orient vs task packet |
| `h7-p24-transfer-grep.txt` | T5 | `trace_explore` grep (docs only) |
| `h7-p24-deferred-evidence.md` | T5 | P24 deferral answer |
| `h6-benchmark-claims.md` | T6 | README claim classification |
| `anti-patterns-not-for-trace.md` | T7 | Law-aligned reject list |
| `h8-skipped-no-codegraph.txt` | T8 | Live MCP skipped |

### Live command log

| # | Command / action | Exit | Evidence |
|---|------------------|------|----------|
| T0 | Preflight: peer root, no `.codegraph/`, README lang rg | 0 | `t0-*` |
| T5 | `rg trace_explore . --glob '*.go' --glob '*.md'` | 0 | `h7-p24-transfer-grep.txt` |
| T8 | MCP `codegraph_explore` live | skipped | `h8-skipped-no-codegraph.txt` |
| — | MCP schema inspect (`GetMcpTools` user-codegraph) | ok | Matches `tools.ts` L1164; requires `projectPath` (no root index) |

### P24 Q4 answer (explicit)

**Is P24 CG transfer still deferred?** **Yes.** EXTERNAL-RESEARCH CG row recommends consolidation toward one explore tool; no `trace_explore` in product code; Phase 37 closed with INT-06 (MCP reorder / FM-08 nudge) not consolidation; Phase 38 DR-HANDOFF forbids implement. Successor phase owns REMEDIATION-PLAN sketch only after S05 saturation.

### Spawn notes

No S02-01a spawn required — mechanism cites fit one row. Reviewer may spawn if live CG index becomes available for T8 parity capture.

---

## Must-answer checklist (planner Q1–Q5)

| Q | Answer location |
|---|-----------------|
| Q1 Explore mechanism | §2 explore row; `h7-explore-mechanism.md` |
| Q2 Index/watch vs Trace | §2 index row; `h5-index-watch-contrast.md` |
| Q3 Single-tool vs 16-tool | §2 MCP row; `h6-single-tool-ux.md` |
| Q4 P24 deferral | §3 H7, §6 P24 answer; `h7-p24-deferred-evidence.md` |
| Q5 Must not adopt | §4; `anti-patterns-not-for-trace.md` |

---

**Next board row:** P38-S02-02 (review)
