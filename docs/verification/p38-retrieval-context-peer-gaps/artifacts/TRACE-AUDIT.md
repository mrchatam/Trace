# TRACE-AUDIT — Phase 38 S01 live audit

**Author:** P38-S01-01 (2026-08-22)  
**Status:** Investigation only — no product changes  
**Authority:** [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) §2  
**Evidence root:** [`experiments/runs/2026-08-22-p38-s01-651/evidence/`](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/)

Dogfood fixture: Trace repo root `/home/ali/Desktop/Trace`. Active task UUID for live captures: `6eb156e9-a3ca-4975-b6da-f2742b68812b`.

---

## §1 Executive summary

Live CLI + MCP investigation of Trace as shipped confirms most S01-owned hypotheses as **gaps** between design docs and runtime behavior, plus harness/discovery friction. Trace’s strengths (task loop, progressive packet, local-first, enforcement hooks) are real but often **not** the first message agents see.

**Headline findings:**

1. **FTS/query (H2):** `ExpandContext` FTS uses `task.Title` only — not agent questions. `trace search` is a separate path; hits do not auto-merge into context packets.
2. **Layers (H3):** Only layers 0–1 ship in JSON; layers 2–3 are documented and explicitly deferred. `--depth 2` expands graph within L0–1, not L2–3 content.
3. **Intent pipeline (H9):** `RETRIEVAL_AND_CONTEXT.md` §3 describes intent extraction; **zero** implementation in `internal/retrieval/`.
4. **Index (H5):** Five language IDs (Go, JS, TS, TSX, Python); manual `trace index` only; no watcher (`hook_installed: false`).
5. **MCP (H6):** Sixteen registered tools, no ranked “start here”; Cursor session exposed 9/16 tools (stale/partial server risk).
6. **Unified orient (H1 partial):** No `query` on `trace_context`; one call does not merge code neighborhood with task packet.
7. **Install moat (H10):** Capability exists in enforcement/install paths; **detect** and AGENTS.md do not lead with task/gate/evidence vs peer graph-first READMEs.
8. **GUI (H8 partial):** `/` → Graph, `/overview` → task loop — defer full peer UX compare to S03.

---

## §2 Findings table

| H | Verdict | Evidence | Notes |
|---|---------|----------|-------|
| **H1** | **confirmed gap** (partial) | [h1-trace-partial.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h1-trace-partial.md); MCP schema; h2-search-compiler.json vs h1-trace-context-depth2.json | Trace-only slice; UA/CG compare → S02/S03/S04 |
| **H2** | **confirmed gap** | [h2-compiler-fts.txt](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h2-compiler-fts.txt); `internal/compiler/compiler.go` L146–151; UA context-builder.ts L32–34 | Title-only FTS; no `--query` / MCP `query` |
| **H3** | **confirmed gap** | [h3-layers-designed-vs-shipped.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h3-layers-designed-vs-shipped.md); `internal/compiler/doc.go` L7; h3-layers-packet-depth2.json | Max item.layer=1; no L2/L3 sections |
| **H5** | **confirmed gap** | [h5-index-langs.txt](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h5-index-langs.txt); `language_adapter.go` L17–25; h5-index-status.json | INTAKE “3 langs” stale → 5 IDs; manual index; no watcher |
| **H6** | **confirmed gap** | [h6-mcp-surface.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h6-mcp-surface.md); [h6-mcp-live-probes.json](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h6-mcp-live-probes.json); `server.go` L227–235 | 16 tools; discovery friction; 9/16 visible in Cursor MCP |
| **H8** | **inconclusive** (partial) | [h8-gui-partial.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h8-gui-partial.md); `web/src/App.tsx` L21–22 | Graph default route; defer Graphify/UA hook compare → S03 |
| **H9** | **confirmed gap** | [h9-intent-pipeline.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h9-intent-pipeline.md); RETRIEVAL_AND_CONTEXT.md §3; h9-intent-grep.txt | Doc pipeline vs zero retrieval code |
| **H10** | **confirmed gap** | [h10-install-moat.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h10-install-moat.md); h10-install-detect.json; AGENTS.md L13–20 | Moat buried in install subcommands |

**Confidence:** H2, H3, H5, H9, H1-partial — **high** (code + live). H6 — **high** (test + MCP). H10 — **medium** (doc/read). H8 — **medium-low**, defer S03.

---

## §3 Live command log

| # | Command | Exit | Evidence file |
|---|---------|------|---------------|
| T0 | `go test ./internal/mcp/ -run TestRegisteredToolNames_IncludesTracePlan -count=1` | 0 | t0-mcp-tool-names.txt |
| T0 | `trace version` | 0 | t0-trace-version.txt |
| T0 | `trace tasks` | 0 | t0-trace-tasks.json |
| T1 | `trace context 6eb156e9-… --format json` | 0 | h2-context-packet.json |
| T1 | `trace search "retrieval FTS compiler" --limit 8` | 0 | h2-search-same-topic.json |
| T1 | `trace search "compiler" --limit 8` | 0 | h2-search-compiler.json |
| T1 | `trace search "e2e-promoted" --limit 8` | 0 | h2-search-task-title.json |
| T2 | `trace context 6eb156e9-… --depth 2 --format json` | 0 | h3-layers-packet-depth2.json, h1-trace-context-depth2.json |
| T4 | `trace index status` | 0 | h5-index-status.json |
| T4 | `trace index internal/compiler/compiler.go` | 0 | h5-index-single-file.txt |
| T5 | `go test ./internal/mcp/ -run TestToolNamesRegistered -count=1` | 0 | h6-mcp-tool-list.txt |
| T6 | `trace install detect` | 0 | h10-install-detect.json |
| T5 | `rg -n intent internal/retrieval internal/compiler --glob '*.go'` | 0 | h9-intent-grep.txt |
| MCP | `trace_version`, `trace_capability`, `trace_context`, `trace_why` | 0* | h6-mcp-live-probes.json |

\* `trace_context` without `project` failed (wrong cwd); succeeded with `project=/home/ali/Desktop/Trace`.

---

## §4 Designed vs shipped (layers + pipeline)

### Layer 0–1 vs 2–3

| Layer | Designed (RETRIEVAL_AND_CONTEXT.md §4) | Shipped (live packet) |
|-------|----------------------------------------|------------------------|
| 0 | Task, objective, exit criteria, state | **Yes** — task, task_state, goal items with `layer: 0` |
| 1 | Files, symbols, deps, decisions, assumptions | **Partial** — graph neighbors (discovery, sibling tasks); sparse symbol/file in e2e task fixture |
| 2 | Dependents, discoveries, future tasks, neighbors | **No** — absent from JSON; doc.go L7 defers |
| 3 | Historical/cross-module depth | **No** — absent |

`--depth 2` runs `ExpandContext` (more graph hops + title FTS) but packet `layer` field stays ≤ 1.

### FTS / query input source

| Stage | Input text | Source |
|-------|------------|--------|
| Graph expand seed | Task entity | `compiler.go` task load |
| FTS inside compiler | **`task.Title`** | `compiler.go` L148 |
| CLI `trace search` | User query string | `cmd/trace/search.go` — **not wired into compiler** |
| MCP `trace_context` | `task_id` only | No `query` in schema |

**Answer:** Compiler retrieval for context packets is **title-token FTS + graph expand**, not agent natural-language query. UA peer uses `SearchEngine(graph, query)` (context-builder.ts L32–34).

### Intent pipeline (H9)

Designed: task/query → **intent extraction** → hybrid retrieval → compiler (RETRIEVAL_AND_CONTEXT.md §3).  
Shipped: direct graph expand + title FTS; “intent” appears only as trust banner text in `packet.go` L200–201.

---

## §5 Non-gaps (Trace strengths)

Observations for S04 moat row — not peer gaps:

| Strength | Evidence |
|----------|----------|
| **Progressive task context packet** | Bounded JSON with budget, reason_codes, trust labels (h2-context-packet.json) |
| **Task loop + gate + review** | MCP trace_loop, trace_review, trace_transition; Overview GateStrip (web/src/screens/Overview.tsx) |
| **Enforcement harness** | internal/install/enforcement.go — TRACE_TASK_ID, loop gate, plan bootstrap |
| **Local-first / no dump API** | Packet caps (items_kept, token_limit); Laws 6–7 honored in compiler |
| **Why / causal trace** | trace_why live probe returned discovery→task chain (h6-mcp-live-probes.json) |
| **Explicit layer honesty** | doc.go L7 documents L2–3 deferral rather than silent omission |
| **Version probe for stale MCP** | trace_version {ok,name,version} |

---

## §6 Open questions → spawn list

| Item | Trigger | Suggested owner |
|------|---------|-----------------|
| Symbol/file inclusion in context for indexed tasks | Re-run with task linked to symbols or richer seed | Optional S01-01a if reviewer wants stronger H1/H3 symbol evidence |
| Cursor MCP 9/16 tool exposure | Restart trace-mcp after rebuild; document in S04 harness row | S04 / install docs (plan only) |
| GUI onboarding vs Graphify graph.html | Screenshot + worked example compare | S03 (H8) |
| CG watcher + lang matrix | Peer live read | S02 (H5 partial) |
| Unified explore tool (H7) | CG codegraph_explore vs trace_search+why | S02 |
| Semantic channel (H4) | Graphify contrast | S03 |

No spawn inserted by implementer; reviewer may add S01-01a per §6 first row.

### Potential improvements (unranked — not build tasks)

- Optional `query` param on `trace context` / MCP merging search hits into packet
- Ranked MCP “orient recipe” in tool descriptions or trace_capability
- Layers 2–3 explicit request API when P0-X scope opens
- Install detect output summarizing moat (tasks → context → loop gate)
- Index watcher or post-commit hook parity messaging vs CG

---

## Appendix — Hypothesis verdict summary

| H | Verdict | Confidence |
|---|---------|------------|
| H1 (partial) | confirmed gap | high |
| H2 | confirmed gap | high |
| H3 | confirmed gap | high |
| H5 | confirmed gap | high |
| H6 | confirmed gap | high |
| H8 (partial) | inconclusive | medium-low |
| H9 | confirmed gap | high |
| H10 | confirmed gap | medium |
