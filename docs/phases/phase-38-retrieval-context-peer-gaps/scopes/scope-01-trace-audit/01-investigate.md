# P38-S01-01 — Trace live audit

## Metadata
- id: P38-S01-01
- todo_ids: [P38-S01-01]
- role: implementer
- skills: [research, graphify, code-explorer, debugging-and-error-recovery]
- mcps: [user-trace]
- verification: mixed
- hooks: none

## Objective

Investigate Trace **as shipped** (CLI + MCP + index + install + optional GUI). Author **`TRACE-AUDIT.md`** with per-hypothesis verdicts and evidence. Identify gaps **and non-gaps**. **Do not implement fixes.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §2 hypothesis register (authority for verify/reject criteria)
- [INTAKE.md](../../INTAKE.md), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md), [PEER-FIXTURES.md](../../PEER-FIXTURES.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — planner locks from P38-S01-00
- Trace anchors: [`internal/compiler/compiler.go`](../../../../../internal/compiler/compiler.go), [`internal/compiler/doc.go`](../../../../../internal/compiler/doc.go), [`internal/retrieval/doc.go`](../../../../../internal/retrieval/doc.go), [`internal/mcp/server.go`](../../../../../internal/mcp/server.go), [`internal/analyzers/language_adapter.go`](../../../../../internal/analyzers/language_adapter.go)
- Design doc: [`docs/RETRIEVAL_AND_CONTEXT.md`](../../../../RETRIEVAL_AND_CONTEXT.md)
- CLI: [`cmd/trace/context.go`](../../../../../cmd/trace/context.go), [`cmd/trace/search.go`](../../../../../cmd/trace/search.go), [`cmd/trace/index.go`](../../../../../cmd/trace/index.go), [`cmd/trace/install.go`](../../../../../cmd/trace/install.go)
- GUI (H8 partial): [`web/src/App.tsx`](../../../../../web/src/App.tsx), [`web/src/screens/Graph.tsx`](../../../../../web/src/screens/Graph.tsx), [`web/src/screens/Overview.tsx`](../../../../../web/src/screens/Overview.tsx)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (P38-S01-00 — do not re-debate)

| Item | Value |
|------|-------|
| Output path | `scopes/scope-01-trace-audit/TRACE-AUDIT.md` |
| Product edits | **Forbidden** (Go/TS/web) |
| Method | **Live CLI/MCP + file:line** — not docs-only |
| Evidence root | `experiments/runs/YYYY-MM-DD-p38-s01-651/evidence/` |
| Dogfood fixture | Trace repo root (`/home/ali/Desktop/Trace` or `-C` override); read-only on `.trace/` |
| Hypotheses (S01-owned) | **H2, H3, H5, H6, H9, H10** + **H1 partial, H8 partial** (per INVESTIGATION-INDEX §3) |
| Verdict vocabulary | `confirmed gap` \| `not a gap` \| `inconclusive` (+ defer note if DR-NOSSEM / law fit) |
| Spawn | Unbounded slice → note in §5 + reviewer may spawn S01-01a before S02 |
| Recommendations | “Potential improvement (unranked)” only — no committed build tasks |

## Must answer (planner handoff — embed in TRACE-AUDIT.md)

1. **Per hypothesis:** verdict + evidence approach (command output or `file:line`).
2. **Layer 0–1 vs 2–3:** what is **designed** (docs) vs **shipped** (live packet JSON).
3. **FTS/query:** what text the compiler **actually** uses for retrieval inside `TaskContext` / `ExpandContext`.
4. **MCP tool surface:** full list + grouping + agent discovery friction (observation only).
5. **Index:** language adapters shipped, manual vs auto trigger, freshness/staleness story.
6. **Install/harness:** how Trace moat (tasks/gates/evidence) is or isn't surfaced vs orient-first peers.

---

## Investigation todos (run in order; do not skip)

### T0 — Preflight + evidence folder

```bash
EV=experiments/runs/$(date +%Y-%m-%d)-p38-s01-651/evidence
mkdir -p "$EV"
go test ./internal/mcp/ -run TestRegisteredToolNames_IncludesTracePlan -count=1 2>&1 | tee "$EV/t0-mcp-tool-names.txt"
trace version 2>&1 | tee "$EV/t0-trace-version.txt"
trace tasks 2>&1 | tee "$EV/t0-trace-tasks.json"
```

- Pick an **active task UUID** from `trace tasks` (or seed one read-only) for T2/T3/T8 live captures.
- Record date + row id in every evidence file header.

### T1 — H2: Compiler FTS inputs (task title vs agent query)

**Code read:**

- [`internal/compiler/compiler.go`](../../../../../internal/compiler/compiler.go) ~L146–151 — `c.retr.Search(ctx, task.Title, …)` (FTS uses **task title**, not MCP/CLI query param).
- [`cmd/trace/context.go`](../../../../../cmd/trace/context.go) — flags: `depth`, `format`, `include-why`; **no `--query`**.
- [`cmd/trace/search.go`](../../../../../cmd/trace/search.go) — separate `trace search <query>` path (not wired into compiler by default).

**Live:**

```bash
TASK=<uuid-from-t0>
trace context "$TASK" --format json 2>&1 | tee "$EV/h2-context-packet.json"
trace search "retrieval FTS compiler" --limit 8 2>&1 | tee "$EV/h2-search-same-topic.json"
```

**Compare:** Do search hits for an agent-style question appear in the context packet without manual merge? UA peer cite (read-only, no run): INVESTIGATION-INDEX H2 → UA `context-builder.ts` SearchEngine(query).

**Verdict target:** `h2-compiler-fts.txt` — state whether FTS query source is title-only; note if `trace_context` MCP schema lacks query param (see `user-trace` tool schema).

### T2 — H3: Layers 0–1 shipped vs 2–3 designed

**Code/doc read:**

- [`internal/compiler/doc.go`](../../../../../internal/compiler/doc.go) L7 — “Layers 2–3 are not auto-loaded in P0-X.”
- [`docs/RETRIEVAL_AND_CONTEXT.md`](../../../../RETRIEVAL_AND_CONTEXT.md) §4 Layer 2–3 definitions.
- [`internal/compiler/packet.go`](../../../../../internal/compiler/packet.go) — `Layer` field on items; grep `Layer:` usage in compiler.

**Live:**

```bash
trace context "$TASK" --depth 2 --format json 2>&1 | tee "$EV/h3-layers-packet-depth2.json"
```

**Inspect JSON:** Max `layer` value present? Any layer-2/3-specific reason codes or sections? Depth-2 == ExpandContext only, or true L2 content?

**Verdict target:** `h3-layers-designed-vs-shipped.md` — table: Layer \| designed (doc) \| shipped (packet) \| evidence.

### T3 — H9: Intent extraction pipeline (documented vs implemented)

**Doc read:**

- [`docs/RETRIEVAL_AND_CONTEXT.md`](../../../../RETRIEVAL_AND_CONTEXT.md) §3 pipeline diagram (`task/query → intent extraction → …`).

**Code grep:**

```bash
rg -n 'intent' internal/retrieval internal/compiler --glob '*.go' 2>&1 | tee "$EV/h9-intent-grep.txt"
```

Note: `packet.go` “project intent” banner is **trust labeling**, not an extraction pipeline.

**Live:** Confirm no intent-extraction stage affects `trace context` output shape (compare T1/T2 packets).

**Verdict target:** `h9-intent-pipeline.md` — doc claims vs code paths; cite line numbers or “zero matches” in retrieval/.

### T4 — H5: Index languages, manual trigger, freshness

**Code read:**

- [`internal/analyzers/language_adapter.go`](../../../../../internal/analyzers/language_adapter.go) L17–25 — `builtinAdapters`: JS, TS, TSX, Python, Go (**5 language ids**, 4 extractor families).
- [`cmd/trace/index.go`](../../../../../cmd/trace/index.go) — manual `trace index` / path args; git refresh on run.
- [`cmd/trace/index_status.go`](../../../../../cmd/trace/index_status.go) — `stale`, `last_indexed_commit`, `hook_installed`.

**Live:**

```bash
trace index status 2>&1 | tee "$EV/h5-index-status.json"
# Optional: trace index <single-file> on a known .go path — capture indexed/skipped counts
trace index internal/compiler/compiler.go 2>&1 | tee "$EV/h5-index-single-file.txt"
```

**Reconcile:** INTAKE says “3 langs”; shipped table has Go + JS/TS/TSX + Python. Document **actual count** and **no file watcher** (contrast CG in S02).

**Verdict target:** `h5-index-langs.txt` — lang list, manual-only, staleness semantics.

### T5 — H6: MCP tool surface + discovery friction

**Code read:**

- [`internal/mcp/server.go`](../../../../../internal/mcp/server.go) L227–236 — `RegisteredToolNames()` (16 tools, locked order).
- Skim handler registrations for **read vs write** grouping.

**Live MCP (preferred) or CLI mirror:**

| Tool | Live probe |
|------|------------|
| `trace_version` | `{ok,name,version}` |
| `trace_capability` | `action=list` |
| `trace_context` | `task_id`, `depth=1`, `format=json` |
| `trace_search` | query string |
| `trace_why` | entity from packet |

```bash
# Sanity: names match test
go test ./internal/mcp/ -run TestToolNamesRegistered -count=1 2>&1 | tee "$EV/h6-mcp-tool-list.txt"
```

**Discovery notes (observation):** Is there a ranked “start here” (context vs search vs why vs plan)? Does `trace_capability` / `trace_plan bootstrap` reduce FM-08 paralysis? Compare count to CG single-tool (peer cite deferred to S02).

**Verdict target:** `h6-mcp-surface.md` — table: tool \| read/write \| orient? \| friction note.

### T6 — H10: Install / harness moat surfacing

**Read:**

- [`cmd/trace/install.go`](../../../../../cmd/trace/install.go) — subcommands: `detect`, `cursor`, `claude`, `agents`, hooks.
- [`internal/install/enforcement.go`](../../../../../internal/install/enforcement.go) — `ParentOrchestratorRule`, loop gate strings.
- [`internal/install/agents.go`](../../../../../internal/install/agents.go) — AGENTS.md / rules bundle content.
- P24 EXTERNAL-RESEARCH OH/SWE/CG harness rows ( skim — orient-first messaging ).

**Live:**

```bash
trace install detect 2>&1 | tee "$EV/h10-install-detect.json"
# Read-only: cat cmd/trace/AGENTS.md or installed .cursor/rules if present — cite paths
```

**Compare:** Do install outputs lead with **task loop / gate / evidence** before code-graph grep? vs peer README first lines (cite only).

**Verdict target:** `h10-install-moat.md`.

### T7 — H1 (partial): Query + task packet — Trace-only slice

**Scope:** Trace side only; full H1 matrix is S04.

**Live:**

```bash
trace context "$TASK" --depth 2 --format json 2>&1 | tee "$EV/h1-trace-context-depth2.json"
```

**Questions:**

- Does one `trace_context` call merge task scope + query-driven symbol neighborhood?
- Is agent query input absent from MCP schema (no `query` field on `trace_context`)?
- Can agent manually compose `trace_search` + `trace_context` — and is that documented?

**Verdict target:** `h1-trace-partial.md` — structural gap vs harness gap; defer UA/CG compare to S02/S03.

### T8 — H8 (partial): GUI orient vs graph-first onboarding

**Scope:** Skim only — full H8 in S03.

**Read:**

- [`web/src/App.tsx`](../../../../../web/src/App.tsx) — `/` → `Graph`, `/overview` → `Overview`.
- [`web/src/screens/Graph.tsx`](../../../../../web/src/screens/Graph.tsx) — overview seed compose, caps (Laws 6–7).
- Phase 32–33 docs if needed for hook/onboarding intent.

**Optional live:** `trace serve` + browser snapshot of `/` and `/overview` (screenshot → `$EV/h8-gui-*.png`). **Not required** for S01 exit if file:line + route map suffice.

**Verdict target:** `h8-gui-partial.md` — what Trace GUI offers vs “committed graph.html hook” peers; mark **partial / defer to S03**.

### T9 — Synthesize TRACE-AUDIT.md

Author deliverable using evidence from T0–T8.

---

## Deliverable shape (TRACE-AUDIT.md)

### §1 Executive summary
Investigation-only; link INVESTIGATION-INDEX + evidence folder.

### §2 Findings table

| H | Verdict | Evidence (file:line or command) | Notes |
|---|---------|-----------------------------------|-------|
| H1 | partial | … | Trace-only slice |
| H2 | … | … | |
| H3 | … | … | |
| H5 | … | … | |
| H6 | … | … | |
| H8 | partial | … | GUI skim |
| H9 | … | … | |
| H10 | … | … | |

Use INVESTIGATION-INDEX §2 **verified if / rejected if** columns when picking verdicts.

### §3 Live command log
Redacted commands + exit status; pointer to `$EV/` files.

### §4 Designed vs shipped (layers + pipeline)
Dedicated subsection answering planner Q2–Q3 (L0–1 vs L2–3; FTS input source).

### §5 Non-gaps (Trace strengths)
Moat observations: task loop, gate, evidence, progressive packet, local-first.

### §6 Open questions → spawn list
Items for S05 / optional S01-01a (scoped slice + trigger).

---

## Exit criteria

- [ ] TRACE-AUDIT.md §§1–6 complete
- [ ] Every S01-owned H* row has verdict + evidence
- [ ] At least **one live Trace CLI or MCP invocation** per major gap claim (stored in `$EV/`)
- [ ] Layer 0–1 vs 2–3 and FTS/query inputs explicitly answered (§4)
- [ ] MCP tool inventory + discovery notes present (§2 H6)
- [ ] No implement recommendations as committed work
- [ ] Board row P38-S01-01 → `done` with Notes (confidence)

## Next

`P38-S01-02`
