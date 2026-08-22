# Phase 32 S00 — RESEARCH

**Date:** 2026-08-21  
**Row:** P32-S00-01  
**Peers:** local clones under `similar projects/graphify/`, `similar projects/Understand-Anything/`  
**Baseline:** locked live facts from P32-S00-00 (Overview index; Graph `/graph`; budgets 50/100; OpenAPI why/context/impact/graph/search/reviews; `ops.ts` missing `getImpact`; serve `:7432` fail-on-conflict).

## Summary

Trace’s Phase 29 GUI is a **nav-first ops shell**: home is Overview (active task + loop gate), and Graph is a secondary `/graph` route with search → center → budgeted `getGraph` (`DEFAULT_MAX=50`, `UI_CAP=100`) and **no side inspector**. Peers set the explorer bar: Graphify emphasizes explained edges, path/query, and a clickable local `graph.html`; Understand-Anything ships search + graph + rich node inspector (summaries, relationships, tours) with a local viewer that auto-increments port and tokenizes the URL. The top **depth** gap is graph-home + selection inspector that surfaces why/context/impact/reviews without leaving the canvas — APIs largely exist; UI wiring and IA do not. For **P32-PORT**, confirm fail-on-conflict at `127.0.0.1:7432` and recommend S02 ship **#1** (friendly `EADDRINUSE` + `--addr` examples) as minimum; treat UA auto-port as optional peer pattern for #2, not Phase 32 default.

## Peer matrix

| Dimension | Graphify | Understand-Anything | Trace today | Borrow / reject |
|-----------|----------|---------------------|-------------|-----------------|
| Home surface | Local `graph.html` is the interactive map (force-directed, clickable communities); CLI also query/path/explain | Interactive dashboard is the product: graph + search + domain views; viewer opens tokenized URL | Overview nav shell; Graph is a route | **Borrow** graph-as-primary explore; **reject** replacing Trace plan/task/loop with AST-only “code map” home |
| Search → graph | Search/filter in `graph.html`; `graphify query` scoped subgraph; CLI explain | Fuzzy + semantic search → focus nodes on canvas | Graph has search + task pick → `getGraph(center, max_nodes)` | **Borrow** search-to-center; **reject** requiring semantic/LLM search as MVP (Trace `/v1/search` is enough) |
| Node inspector | Click nodes in HTML; CLI `explain` lists connections with EXTRACTED/INFERRED tags | Select node → plain-English summary, relationships, guided tours (rich side panel) | Graph: none; TaskDetail: why/context drawers (raw JSON `<pre>`) | **Borrow** UA-style side inspector on selection; **reject** cloning UA tours/persona UI wholesale |
| Path / filter | First-class `graphify path A B`; community/filter in HTML; edge confidence tags | Layer/domain filters; diff-impact view | Kind filter on Graph search only; node click re-centers; edge labels = rel | **Borrow** clearer relation/path affordances when library-backed; **reject** full community Leiden / god-node UX as Phase 32 must-have |
| Budget / dump | Ships full `graph.json` + HTML for offline browse (product is the dump artifact) | Loads generated `knowledge-graph.json` into viewer (full analyzed graph) | `getGraph` center + `max_nodes`; UI_CAP 100; truncation banner | **Keep** Laws 6–7 budgeted neighborhood; **reject** full-graph dump as GUI default |
| Local serve / ports | Static `graph.html` (open file / any static host) — no multi-project serve story | Viewer: `--port` (default 5173), **auto-increment if taken**, `127.0.0.1` + one-time token URL | Fixed `127.0.0.1:7432`; bare `net.Listen`; fail on conflict; raw `serve: %v` | **Borrow** clear local-loopback + token story already partly present; **borrow idea** of friendly conflict UX; **defer/#2** auto-increment (optional); **reject** silent always-on daemon |
| Semantics Trace must keep | — | — | plan/task/decision/review/loop | do not clone code-AST-only UX |

### Peer notes (cite)

- **Graphify:** README positions `/graphify` → `graph.html` + `GRAPH_REPORT.md` + `graph.json`; “Every edge is explained” with `EXTRACTED` / `INFERRED`; CLI `explain`, `path`, `query` for relation explore (`similar projects/graphify/README.md`). Hero is interactive force-directed HTML, not an ops CRUD shell.
- **Understand-Anything:** README features interactive knowledge graph, node select → summaries/relationships/tours, fuzzy & semantic search, domain view (`similar projects/Understand-Anything/README.md`). Viewer package: `--port` default 5173, **auto-increments if taken**, prints tokenized `http://127.0.0.1:<port>/?token=…`, loopback-only (`similar projects/Understand-Anything/understand-anything-plugin/packages/viewer/README.md`).

## web/ inventory vs explorer job

| Screen / area | Path | Role today | Fit for graph-home explorer |
|---------------|------|------------|-----------------------------|
| Overview | `/` (index) | Active-task dashboard: tasks list, goal search, loop status/gate (`Overview.tsx`) | Ops landing — demote vs canvas; keep as panel/secondary per DESIGN-LOCKS hybrid C |
| Graph | `/graph` | Search/task pick → budgeted XYFlow neighborhood; click expands center; **no inspector** (`Graph.tsx`) | Core canvas to promote to home; add selection inspector + depth panels |
| TaskDetail | `/tasks/:id` | Task CRUD/ops + why/context drawers (JSON dump) | Depth content belongs beside graph selection, not only behind Tasks nav |
| Tasks | `/tasks` | Task list | Secondary / panel entry to tasks |
| Loop | `/loop` | Loop status/next/apply UX | Deep panel on selected task (not primary nav job) |
| Discoveries | `/discoveries` | Discovery listing | Keep accessible; not graph-home |
| Reviews | `/reviews` | Review list | Inspector should surface linked reviews; list remains secondary |
| ReviewDetail | `/reviews/:id` | Single review | Link from inspector; keep route |
| Seed | `/seed` | Export/import status | Ops/settings-adjacent; not explorer home |
| Settings | `/settings` | Token / project chrome | Keep |
| ProjectGate | (Shell gate) | Store-ready gate before ops screens | Keep; unchanged by research |
| Nav | `layout/Nav.tsx` | Eight primary links; Graph mid-list | Expect S01 to reweight toward graph-home + panels-not-nav-CRUD |

**Preflight confirmed (2026-08-21):** `App.tsx` index → Overview, `path="graph"` → Graph; OpenAPI `/v1/why|context|impact|graph|search|reviews`; `bind.go` `DefaultAddr = 127.0.0.1:7432`; `server.go` `ListenAndServe` → `net.Listen` only; peers readable under `similar projects/`.

## Gaps (ordered)

### IA — S01 (before depth implement)

1. **Promote graph to home / demote CRUD nav** — hybrid C: canvas home; tasks / loop / reviews / seed as deep panels on selection (not nav-first). S01 authors `UX-IA.md`; shell wiring follows later scopes.

### Depth — S03 candidates (priority order)

1. **Graph-side node inspector** — on select (not only re-center): entity summary + why + context without navigating to TaskDetail.
2. **Wire impact + reviews (+ evidence/links) into inspector** — OpenAPI `getImpact`, `listReviews`/`getReview` exist; impact has **no** `ops.ts` helper; reviews unused on Graph.
3. **Path / relation clarity** — richer edge/path affordances when library supports; avoid inventing path API without library backing (flag for S02 only if S01 proves need).
4. **Inspector presentation** — replace raw JSON drawers with structured sections (still adapter-only; no business fork).

### Visual craft — S04 candidates (only after depth)

1. Canvas layout/density (less naive polar layout; readable labels at budget).
2. Typography / hierarchy for inspector vs canvas (explorer feel, not ops form stack).
3. Motion for selection / expand / truncation feedback (intentional, sparse).
4. Kind/rel visual encoding (legend, edge weight) without 3D.

### Explicit non-goals this phase

- 3D / Three.js, second SPA, hosted SaaS, MCP `/rpc` in browser, auto-delete root `trace.db`, full Graphiti clone, full-graph dump as product default, always-on public bind, cloning UA guided tours / persona UI as MVP.

## API reuse map (`/v1` → inspector)

| Inspector need | OpenAPI op | Client (`ops.ts`) | Used in UI today | Notes |
|----------------|------------|-------------------|------------------|-------|
| Why | getWhy | yes | TaskDetail drawer | Reuse on Graph selection; Law 19 adapter |
| Context | getContext | yes | TaskDetail drawer | Task-scoped; depth 1\|2 already supported |
| Impact | getImpact | **missing wrapper** | no | Library/OpenAPI present (`api/openapi.yaml` `/v1/impact`); **client gap only** — S02 likely add wrapper if S01 keeps impact in inspector map; not a new core API |
| Neighborhood | getGraph | yes | Graph | Keep center + max_nodes; UI_CAP 100 |
| Search | search | yes | Graph (+ Overview goals) | Enough for search→center |
| Reviews | listReviews / getReview | yes | Reviews screens | Filter by `task_id` for inspector |
| Entity summary | getEntity | yes | (available) | Useful inspector header |
| Task row | getTask / listTasks | yes | Tasks / Graph pick | Selection hydrate |
| Loop status/gate | getLoopStatus / getLoopGate | yes | Overview / Loop | Panel depth, not graph dump |
| Changes / regressions | listChanges / listRegressions | **no wrappers** | no | OpenAPI exists; **do not invent** unless S01 proves library-backed need — else S02 `NO-GAPS.md` for API |
| Path A→B | — | — | no | Peer pattern only; no Trace `/v1/path` — reject inventing without library |

**Library-backed hole for S02:** none confirmed beyond **client wrapper** for `getImpact` (and optional thin wrappers if S01 requires changes/regressions). Expect API story may be `NO-GAPS.md` + client glue; **P32-PORT is separate and still required.**

## P32-PORT

### Confirm light review

| Item | Live fact | Source |
|------|-----------|--------|
| Default | `127.0.0.1:7432` | `internal/httpapi/bind.go` `DefaultAddr`; `cmd/trace/serve.go` `--addr` default |
| Bind | `net.Listen("tcp", s.addr)` — no free-port search, no `:0` default | `internal/httpapi/server.go` `ListenAndServe` |
| On conflict | Listen error → CLI `serve: %v` and exit | `cmd/trace/serve.go` |
| Multi-project | One serve = one root; second project on same port conflicts unless distinct `--addr` | OPEN-PORT-MULTI.md + live CLI |
| Workaround | Manual `trace serve -C <proj> --addr 127.0.0.1:<other>` | OPEN-PORT-MULTI.md |

Matches [`OPEN-PORT-MULTI.md`](../../OPEN-PORT-MULTI.md) light review.

### Peer patterns

- **UA viewer:** `--port`, auto-increment when taken, tokenized loopback URL (`viewer/README.md`). **Borrow idea** for optional #2 / docs storytelling; **reject** copying auto-increment as sole Phase 32 story without also fixing error clarity (#1).
- **Graphify:** file-based `graph.html` — not a multi-`--addr` server; little to borrow for serve conflicts.

### Recommendation for S02

| Option | Ship? | Rationale |
|--------|-------|-----------|
| #1 Friendly EADDRINUSE + `--addr` examples | **Prefer yes (min)** | Matches OPEN-PORT-MULTI owner guidance; fixes silent/opaque conflict for multi-project without changing default bind policy |
| #2 Auto free-port / `:0` | optional / defer | UA peer pattern; useful but changes discoverability/URL printing — only after or with #1; document if deferred |
| #3 Docs-only multi-`--addr` | S05 docs; not sole S02 story | Needed for VERIFY polish; insufficient alone for P32-PORT |

**S02 ownership:** always ships P32-PORT even if API work is `NO-GAPS.md`. Discouraged: `NO-PORT-CHANGE.md` without written reason.

## Handoff to S01

- **Primary job hybrid C:** graph is home; tasks / loop / reviews / seed are deep panels on selection — not nav-first CRUD (`DESIGN-LOCKS.md`).
- **Depth before craft:** S03 inspector (why/context/impact/reviews) before S04 visuals.
- **Laws 6–7:** keep `center` + `max_nodes` / UI_CAP; no full-graph dump default.
- **Law 19:** reuse `/v1` ops; UI formats packets only; note `getImpact` client gap for S02 glue.
- **Keep Trace semantics:** plan/task/decision/review/loop — do not redesign as AST code-map clone.
- **P32-PORT:** S01 need not design ports; S02 ships #1 minimum per this RESEARCH.
- No UX-IA writing in this row — S01 owns `UX-IA.md`.
