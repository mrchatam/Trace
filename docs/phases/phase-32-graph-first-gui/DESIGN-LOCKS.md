# DESIGN-LOCKS — Phase 32 graph-first GUI

**Human-confirmed 2026-08-21** (grill recommendations accepted wholesale).

| Lock | Value |
|------|-------|
| Primary job | **Hybrid (C):** graph is the home; tasks / loop / reviews / seed are **deep panels** on selection — not a nav-first CRUD shell |
| Depth vs craft | **Depth first (S03), visual craft second (S04)** inside this phase |
| Sequencing | **Phase 31** (extra stray-db tests) completes before this phase runs |
| Peer bar | Between **Graphify** (clear relation explore) and **Understand-Anything** (search + graph + rich side inspector). Trace keeps plan/task/decision semantics. |
| Shell strategy | **Evolve `web/`** into graph-home (not a second SPA product). Ops screens become panels / routes secondary to the canvas. |
| Graph tech | Keep **2D `@xyflow/react`** (budgeted). **No Three.js / 3D** as default in this phase. |
| Laws | **Law 19** — UI/HTTP adapters only. **Laws 6–7** — no unbounded full-graph dump; neighborhood / `max_nodes` / progressive expand. |
| API | Prefer existing `/v1` (why, context, impact, reviews, search, graph). OpenAPI/library gaps in **S02** only if S01 proves a missing library-backed need — no business-logic fork in `web/`. |
| P32-PORT | **S02 always owns** port/multi-project serve work even when API is `NO-GAPS.md` — see [`OPEN-PORT-MULTI.md`](OPEN-PORT-MULTI.md). |
| Out of scope | Hosted SaaS, always-on daemon, public bind defaults, MCP `/rpc` in browser, auto-delete of root `trace.db`, full Graphiti clone |

## Open work (phase outcome)

| ID | Topic | Notes |
|----|-------|-------|
| **P32-PORT** | Port conflict + multi-project serve | **Done for Phase 32:** #1 + S05 docs verified; **#2 auto-port deferred** (residual, not a successor). See [`OPEN-PORT-MULTI.md`](OPEN-PORT-MULTI.md). DR-HANDOFF **CLOSED** at `P32-S06-02`. |

## Success sketch

Opening `trace serve` feels like an **explorer**: search → center graph → select node → inspector shows why/context/impact/reviews/links with enough detail to understand *why work exists*, not only task CRUD. Multi-project local use must not silently die on “address already in use” without a clear story (auto-port and/or docs + friendly error).
