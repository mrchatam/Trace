# P38-S02-01 — Codegraph peer investigation

## Metadata
- id: P38-S02-01
- todo_ids: [P38-S02-01]
- role: implementer
- skills: [research, graphify, code-explorer]
- mcps: [user-codegraph]
- verification: mixed
- hooks: none

## Objective

Read-only investigation of **Codegraph** as a peer. Author **`PEER-CG.md`** with mechanism cites (peer `file:line`), hypothesis verdicts for **H1 (partial), H5, H6, H7**, and explicit anti-patterns Trace must not copy. **Do not implement fixes or product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §2 verify/reject for H1, H5, H6, H7
- [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md) — Trace baseline (S01 APPROVED)
- [INTAKE.md](../../INTAKE.md), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md), [PEER-FIXTURES.md](../../PEER-FIXTURES.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — planner locks from P38-S02-00
- P24 [EXTERNAL-RESEARCH.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-03-external-research/EXTERNAL-RESEARCH.md) — CG row + §4 anti-patterns
- Peer root: [`similar projects/codegraph/`](../../../../../similar%20projects/codegraph/)
- CG anchors (spot-check before cite drift):
  - [`src/mcp/tools.ts`](../../../../../similar%20projects/codegraph/src/mcp/tools.ts) — `codegraph_explore` schema ~L1163–1181, `handleExplore` ~L3193+, `buildBlastRadiusSection` ~L3017+
  - [`src/mcp/tools.ts`](../../../../../similar%20projects/codegraph/src/mcp/tools.ts) — `DEFAULT_MCP_TOOLS` ~L1275–1286 (single-tool default surface)
  - [`src/sync/watcher.ts`](../../../../../similar%20projects/codegraph/src/sync/watcher.ts) — debounced auto-sync ~L1–80
  - [`src/mcp/daemon.ts`](../../../../../similar%20projects/codegraph/src/mcp/daemon.ts) — detached MCP daemon ~L1–40
  - [`README.md`](../../../../../similar%20projects/codegraph/README.md) — init/watch L108–127, benchmark L196–218

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (P38-S02-00 — do not re-debate)

| Item | Value |
|------|-------|
| Output path | `scopes/scope-02-codegraph-peer/PEER-CG.md` |
| Peer root | `similar projects/codegraph/` (read-only) |
| Product edits | **Forbidden** (Go/TS/web in Trace repo) |
| Method | **Peer file:line + optional CG MCP** — not README-only for major mechanism claims |
| Evidence root | `experiments/runs/YYYY-MM-DD-p38-s02-654/evidence/` |
| Dogfood fixture | Trace repo has **no** `.codegraph/` (confirmed S02 planner spot-check) — optional MCP live on **any** indexed peer subpath or sample; absence is evidence, not blocker |
| Hypotheses (S02-owned) | **H1 (partial), H5, H6, H7** |
| Verdict vocabulary | `supported` \| `weakened` \| `rejected` \| `inconclusive` (+ defer note if law/DR fit) |
| P24 stance | **Extend** EXTERNAL-RESEARCH CG row — cite whether `trace_explore` / MCP consolidation transfer is still open |
| Spawn | Unbounded CG slice → note in §6 + reviewer may spawn S02-01a before S03 |
| Non-goals | No copying CG into Trace; no implement `trace_explore`; no ranked REMEDIATION-PLAN |

## Must answer (planner handoff — embed in PEER-CG.md)

1. **Mechanism of `codegraph_explore`:** inputs (`query`, `maxFiles`, `projectPath`), outputs (verbatim source sections, call paths, blast radius), caps (`MAX_OUTPUT_LENGTH`, adaptive budget).
2. **Index/watch vs Trace `trace index`:** CG `init` + watcher debounce vs Trace manual index + `hook_installed: false` (TRACE-AUDIT H5).
3. **Single-tool UX vs Trace 16-tool:** default MCP surface (`DEFAULT_MCP_TOOLS`), discovery implications vs Trace FM-08 (TRACE-AUDIT H6).
4. **P24 transfer items:** Is “consolidate Trace MCP read surface toward one explore tool” still **deferred**? Evidence from P24 EXTERNAL-RESEARCH CG row, INTERVENTION-MATRIX INT-06/FM-08, and live Trace MCP inventory.
5. **What Trace must NOT adopt:** Law 6/7 (full-graph dump defaults), always-on daemon as P0 requirement, task/backlog absence as product direction.

---

## Investigation todos (run in order; do not skip)

### T0 — Preflight + evidence folder

```bash
EV=experiments/runs/$(date +%Y-%F)-p38-s02-654/evidence
mkdir -p "$EV"
# Peer presence
test -f "similar projects/codegraph/src/mcp/tools.ts" && echo ok | tee "$EV/t0-peer-root.txt"
# Trace has no CG index on dogfood repo
test ! -d .codegraph && echo "no .codegraph on Trace" | tee "$EV/t0-no-codegraph-on-trace.txt"
# Lang inventory (CG README vs Trace 5 ids)
rg -n "language" "similar projects/codegraph/README.md" | head -20 | tee "$EV/t0-cg-readme-lang-refs.txt"
```

- Record date + row id in every evidence file header.
- Skim TRACE-AUDIT §2 H5/H6/H1 rows — S02 compares **peer side** only; do not re-litigate Trace live captures unless contrasting in tables.

### T1 — H7: `codegraph_explore` mechanism (inputs, outputs, blast radius)

**Peer code read (required file:line cites):**

| Mechanism | Start here |
|-----------|------------|
| Tool schema (inputs) | [`tools.ts`](../../../../../similar%20projects/codegraph/src/mcp/tools.ts) L1163–1181 — `query` required; `maxFiles` default 12; `projectPath` |
| Handler entry | [`tools.ts`](../../../../../similar%20projects/codegraph/src/mcp/tools.ts) L3193–3224 — `handleExplore`: graph traversal → file sections; adaptive `getExploreOutputBudget` |
| Output cap | [`tools.ts`](../../../../../similar%20projects/codegraph/src/mcp/tools.ts) L86–87 — `MAX_OUTPUT_LENGTH = 15000` |
| Blast radius section | [`tools.ts`](../../../../../similar%20projects/codegraph/src/mcp/tools.ts) L3017–3073 — `buildBlastRadiusSection`: callers + test files, locations only |
| Dispatch switch | [`tools.ts`](../../../../../similar%20projects/codegraph/src/mcp/tools.ts) ~L2098 — `case 'codegraph_explore'` |

**Document in evidence:**

- Request shape vs MCP schema (`user-codegraph` tool descriptor — compare to peer source).
- Response sections: symbol grouping, call path, blast radius, truncation/epilogue behavior.
- Progressive/capped behavior (not full repo dump) — tie to Law 6/7 **peer compliance**.

**Trace contrast (read-only, cite TRACE-AUDIT):**

- No `trace_explore`; composition=`trace_search` + `trace_context` + `trace_why` (16 tools).
- [`internal/mcp/server.go`](../../../../../internal/mcp/server.go) `RegisteredToolNames` — list for §2 table.

**Verdict target:** `h7-explore-mechanism.md` + `h7-explore-gap.md` — does CG single-call subsume Trace read orient? P24 transfer still open?

### T2 — H6: Single-tool MCP UX vs Trace 16-tool discovery

**Peer code read:**

- [`tools.ts`](../../../../../similar%20projects/codegraph/src/mcp/tools.ts) L1275–1286 — `DEFAULT_MCP_TOOLS = new Set(['explore'])`; other tools exist but **not listed** by default.
- [`tools.ts`](../../../../../similar%20projects/codegraph/src/mcp/tools.ts) L1266–1272 — `getStaticTools()` / `CODEGRAPH_MCP_TOOLS` env override.
- [`src/mcp/server-instructions.ts`](../../../../../similar%20projects/codegraph/src/mcp/server-instructions.ts) — agent steering text (if present: “PRIMARY TOOL”, call-first language).

**Compare to Trace (TRACE-AUDIT, do not re-run unless spot-check):**

- 16 tools; no ranked default; FM-08 paralysis **plausible** (S01 verdict).
- Observation only: CG **removes choice** by hiding narrow tools; Trace **exposes choice** — agent discovery tradeoff, not winner declaration.

**Verdict target:** `h6-single-tool-ux.md` — table: dimension \| CG \| Trace \| H6 implication.

### T3 — H5: Index, watch, languages vs Trace index

**Peer code read:**

- [`src/sync/watcher.ts`](../../../../../similar%20projects/codegraph/src/sync/watcher.ts) L1–80 — debounced auto-sync; platform strategies; degrade paths.
- [`README.md`](../../../../../similar%20projects/codegraph/README.md) L108–127 — `codegraph init` one-shot full index; auto-sync default.
- CG language list: README supported-languages section (~L145+) — count families vs Trace 5 ids (`language_adapter.go` per TRACE-AUDIT).

**Trace baseline (cite TRACE-AUDIT §2 H5, do not duplicate full live run):**

- Manual `trace index`; `hook_installed: false`; incremental file index, no watcher.

**Optional peer CLI (read-only, if binary available in env):**

```bash
# Only if codegraph CLI on PATH — else skip with note in §3
# codegraph status --help 2>&1 | tee "$EV/h5-cg-status-help.txt"
```

**Verdict target:** `h5-index-watch-contrast.md` — mechanism table + lang count evidence (README vs extractor dirs under `src/extraction/` if needed).

### T4 — H1 (partial): Query-driven orient vs Trace task packet

**Peer mechanism:**

- `codegraph_explore` accepts **natural-language or symbol bag** in `query` (tools.ts L1168–1170) — no task UUID, no backlog.
- README agent workflow L206+ — orient-before-edit pattern (benchmark narrative).

**Trace baseline (TRACE-AUDIT H1 partial):**

- `trace_context` requires `task_id`; no `query` param; task packet ≠ query neighborhood in one call.

**Side-by-side (required in PEER-CG §2):**

| Step | CG | Trace |
|------|----|-------|
| Agent question input | `query` string | separate `trace_search` (not in packet) |
| Task/plan scope | **none** (read-only graph) | task packet L0 |
| Single-call orient | `codegraph_explore` | partial — needs compose |

**Verdict target:** `h1-cg-partial.md` — structural gap **supported/weakened** for CG side of H1; full matrix deferred S04.

### T5 — P24 transfer items: still deferred?

**Doc read (required cites):**

- P24 EXTERNAL-RESEARCH §2 **CG row** — transfer: “Consolidate Trace MCP read surface toward one high-signal explore tool…”
- P24 INTERVENTION-MATRIX — INT-06 (MCP reorder / FM-08); no shipped `trace_explore` in Phase 37 closure.
- Phase 38 INTAKE H7 — hypothesis register owner S02.
- Grep Trace repo: `trace_explore` — expect **zero** product symbol.

```bash
rg -n 'trace_explore' . --glob '*.go' --glob '*.md' 2>&1 | tee "$EV/h7-p24-transfer-grep.txt"
```

**Answer explicitly:** Is P24 CG transfer **still deferred** (investigation-only in P38)? Evidence: no implement row, no tool in MCP registry, FM-08 addressed as reorder not consolidation.

**Verdict target:** `h7-p24-deferred-evidence.md`.

### T6 — README / benchmark claims (verifiable vs marketing)

**Read:**

- [`README.md`](../../../../../similar%20projects/codegraph/README.md) L196–218 — 43 tool calls without index vs 1–4 explores; cost/token claims.
- [`docs/benchmarks/`](../../../../../similar%20projects/codegraph/docs/benchmarks/) — if present, cite methodology file path.

**Classify each major claim:**

| Claim | Verifiable how | PEER-CG treatment |
|-------|----------------|-------------------|
| 43 vs 1–4 tool calls | README cites harness; optional read benchmark doc | Cite methodology path or mark **observation / not re-run in P38** |
| Auto-sync latency | README L267+; watcher.ts mechanism | Mechanism cite suffices; skip live timing unless easy |
| 20+ languages | README icons + `src/extraction/` | Count or cite README list |

**Verdict target:** `h6-benchmark-claims.md` — do not treat marketing numbers as Trace gap proof without mechanism backing.

### T7 — Anti-patterns: what Trace must NOT adopt

**Peer read (for contrast, not copy list):**

- [`src/mcp/daemon.ts`](../../../../../similar%20projects/codegraph/src/mcp/daemon.ts) L1–25 — detached `codegraph serve --mcp` daemon, shared watcher/SQLite across clients.
- CG explore caps vs hypothetical full dump — tie to **G Law 6/7** and DESIGN-LOCKS.
- CG has **no** task loop / gate / evidence — moat row seed (Trace strength).

**Trace law alignment (cite project rules / G_PROJECT_LAWS):**

- No always-on daemon as P0-X requirement (P24 §4 anti-patterns).
- No full-graph dump as default API/MCP behavior.
- MCP optional on harness — must not become Trace core path requirement.

**Verdict target:** `anti-patterns-not-for-trace.md` — bullet list with peer cite + law cite.

### T8 — Optional live: Codegraph MCP `codegraph_explore`

**Only if** an indexed project exists (peer clone with `.codegraph/`, or user-built index elsewhere):

```
# Example — adjust projectPath to indexed tree
# MCP codegraph_explore query="handleExplore ToolHandler" projectPath="/path/with/.codegraph"
```

- Save redacted response snippet → `$EV/h8-live-explore-sample.txt`
- Note latency / orient feel (observation paragraph in §1).

**If no index (expected on Trace repo):** Record in §3 log as `skipped — no .codegraph/`; T1–T7 file:line evidence satisfies exit.

### T9 — Synthesize PEER-CG.md

Author deliverable using T0–T8 evidence.

---

## Deliverable shape (PEER-CG.md)

### §1 Executive summary
Investigation-only; link INVESTIGATION-INDEX + evidence folder; 3–5 bullets on CG mechanism vs Trace contrast.

### §2 Comparison table vs Trace (mechanism, not feature checklist)

Columns suggested: **Dimension | Codegraph mechanism (file:line) | Trace mechanism (file:line or TRACE-AUDIT) | Notes**

Minimum rows: explore/orient, index freshness, MCP surface, task/moat, caps/progressive context.

### §3 Gap hypotheses **supported / weakened / rejected**

| H | Verdict | Peer evidence | Trace contrast | Confidence |
|---|---------|---------------|----------------|------------|
| H1 partial | … | … | TRACE-AUDIT | … |
| H5 | … | … | … | … |
| H6 | … | … | … | … |
| H7 | … | … | … | … |

Use INVESTIGATION-INDEX §2 **verified if / rejected if** when picking verdicts.

### §4 Anti-patterns Trace must not copy
Daemon-as-P0, full dump defaults, abandoning task/gate/evidence moat for graph-only orient.

### §5 Trace strengths peers lack (moat row seed)
Tasks, gates, evidence, plan tree — CG is read-only graph.

### §6 Evidence appendix + live command log
Peer paths, `$EV/` files, optional MCP capture; P24 deferral answer for Q4.

---

## Exit criteria

- [ ] PEER-CG.md §§1–6 complete
- [ ] Every S02-owned H* row has verdict + **peer file:line** evidence (README alone insufficient for T1/T2/T3/T7)
- [ ] Planner must-answer Q1–Q5 addressed in §§2–4
- [ ] P24 transfer/deferral stated with grep or matrix cite
- [ ] No Go/TS/web product diff in Trace repo
- [ ] Board row P38-S02-01 → `done` with Notes (confidence + next P38-S02-02)

## Next

`P38-S02-02`
