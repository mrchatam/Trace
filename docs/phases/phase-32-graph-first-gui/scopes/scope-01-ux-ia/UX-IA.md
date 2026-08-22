# Phase 32 S01 — UX-IA

**Date:** 2026-08-21  
**Row:** P32-S01-01  
**Status:** Spec for S03+ implementers — docs only (no `web/` edits this row)

**Preflight (live, 2026-08-21 — confirmed, not changed this row):**

| Fact | Evidence |
|------|----------|
| Index = Overview; Graph = `/graph` | `web/src/App.tsx` routes |
| Nav = 8 CRUD links (Overview first) | `web/src/layout/Nav.tsx` `LINKS` |
| Graph budgets `DEFAULT_MAX=50`, `UI_CAP=100`; click re-centers; **no inspector** | `web/src/screens/Graph.tsx` |
| OpenAPI why/context/impact/graph/search/reviews/entity | RESEARCH API reuse map; `web/src/api/schema.d.ts` |
| `ops.ts` has getWhy/getContext/getGraph/…; **no `getImpact`** | `web/src/api/ops.ts` |

---

## Summary

Phase 32 evolves the existing `web/` SPA into a **hybrid C** explorer: the **budgeted XYFlow graph canvas is home**, not Overview. Selection opens a structured **inspector** (summary → why → context → impact → reviews → links, plus an optional loop strip for tasks) that formats `/v1` packets only (Law 19). Neighborhood loads always require `center` + `max_nodes` with live defaults **50 / UI_CAP 100**; truncation is honest and there is no full-graph dump CTA (Laws 6–7). Ops surfaces (tasks, loop, reviews, seed, etc.) stay reachable as **deep panels on selection** and/or **secondary routes** — not a second SPA. S02 inherits a **client-only** `getImpact` glue gap plus **mandatory P32-PORT**; no new library `/v1` ops and no `/v1/path` are required by this IA. S03 implements inspector depth; S04 owns visual craft after depth ships.

## Product framing

| Principle | Implication |
|-----------|-------------|
| Hybrid C (DESIGN-LOCKS) | Graph home; ops as deep panels / secondary routes |
| Law 19 | `/v1` adapters only — no browser SQLite, no business-logic fork in `web/` |
| Laws 6–7 | Always `center` + `max_nodes`; keep DEFAULT_MAX=50 / UI_CAP=100; truncation banner; no full-graph dump default / primary CTA |
| Peer bar | Between Graphify (relation explore) and Understand-Anything (search + graph + rich inspector); keep Trace plan/task/decision/review/loop semantics |
| Out | 3D / Three.js, second SPA, hosted SaaS, MCP `/rpc` in browser, P32-PORT design (S02), inventing `/v1/path`, cloning UA tours/persona or Graphify AST-only home |

## Shell layout

### Target chrome

```text
[Chrome: project title · health/version · search entry · Explore-first nav]
├── Graph canvas (HOME) ──────────────┬── Inspector (on selection)
│   search → center → budgeted XYFlow │   summary → why → context →
│   select ≠ expand                   │   impact → reviews → links
│   truncation honesty                │   (+ loop strip if task)
└── Secondary routes / deep panels ←──┘
    Overview · Tasks · Loop · Reviews · Discoveries · Seed · Settings
```

- **Primary plane:** full-height (or dominant) **graph canvas** at `/` (or equivalent primary route).
- **Inspector region:** side panel (desktop) / sheet (narrow) opened by **node selection**; structured sections, not raw JSON as the primary presentation.
- **Secondary ops access:** reweighted nav (Explore/Graph + Search first); deep-link buttons from inspector into existing detail routes; selection-relevant task/loop/review/seed actions as panels when useful.

### Route / panel map

| Surface | Today | Target (home / panel / secondary route) | Notes |
|---------|-------|----------------------------------------|-------|
| Graph | `/graph` | **Home** (`/` or primary) | XYFlow canvas; budgeted `getGraph` |
| Overview | `/` | Secondary route and/or selection panel | Demote — active task / gate strip may surface in chrome or inspector, not as landing |
| Tasks | `/tasks` | Secondary route; list entry from search/nav | Not nav-CRUD-primary |
| TaskDetail | `/tasks/:taskId` | Secondary route + deep-link from inspector | Valid for refresh / share; prefer inspector for in-canvas depth |
| Loop | `/loop` | Secondary route; **deep panel** when task selected | Optional inspector Loop strip → “Open Loop” |
| Reviews | `/reviews` | Secondary route; inspector section filters by `task_id` | List remains for browse |
| ReviewDetail | `/reviews/:reviewId` | Secondary route; deep-link from inspector Reviews | Keep |
| Discoveries | `/discoveries` | Secondary route | Unchanged role |
| Seed | `/seed` | Secondary route; seed-adjacent actions only when relevant to selection | Honesty about HTTP export/import unchanged |
| Settings | `/settings` | Secondary route | Keep (token / theme / bind display) |
| Nav | 8 equal CRUD links | **Reweight** Explore/Graph + Search first; Overview demoted | Do not invent a second SPA |

## Inspector depth map

Ordered depth map for S03. UI formats library packets only (Law 19). Structured sections beat raw `<pre>` dumps as the primary presentation.

| Section (order) | `/v1` op(s) | Client today | S03 notes |
|-----------------|-------------|--------------|-----------|
| 1. Entity summary | `getEntity` / `getTask` (task row hydrate) | ops yes | Header: kind, title, id, status when task |
| 2. Why | `getWhy` | ops yes (TaskDetail drawer) | Structured Why packet; not raw JSON primary |
| 3. Context | `getContext` | ops yes (TaskDetail drawer) | Task-scoped; depth **1\|2** control |
| 4. Impact | `getImpact` | **ops missing** | Flag **S02 client wrapper**; OpenAPI already present (`task_id` query optional — prefer when selection is a task; else omit section or show empty/honest state) |
| 5. Reviews | `listReviews` / `getReview` | ops yes (Reviews screens) | Filter by `task_id` when task-selected; deep-link detail |
| 6. Links / relations | edges from current `getGraph` neighborhood + entity deep-links | Graph canvas only | From loaded neighborhood; link out to entity/task routes |
| Loop (optional strip) | `getLoopStatus` / `getLoopGate` | ops yes | **Task selection only** — deep panel, not a graph dump |

**Task-scoped sections:** Context, Reviews (filtered), Impact (prefer `task_id`), and Loop apply when the selection resolves to a **task** (or a task-linked entity hydrate). For non-task centers, still show Entity summary + Why + Links; hide or collapse Context/Loop/Reviews-filter rather than inventing non-library packets.

**Not in this map (default omit):** `listChanges` / `listRegressions` (no proven library-backed IA need); `/v1/path` (does not exist — do not invent).

## Navigation / selection / search flows

1. **Search → center → budgeted graph**  
   User searches (`search`) or picks a task → set `center` → `getGraph(center, max_nodes)` with budget ≤ UI_CAP. Empty search results: prompt to refine query or pick a task.

2. **Select → inspector**  
   Single-click (or equivalent) on a node → **selection highlight** + open/focus inspector; load depth-map sections for that entity/task. Does **not** by itself re-center the neighborhood.

3. **Expand / re-center (distinct from select)**  
   Explicit affordance only: double-click, “Expand” / “Use as center” control, or equivalent. Progressive expand = **new center + same budget**; if response `truncated`, show honesty banner. Today’s Graph behavior (every click re-centers, no inspector) is the gap this flow closes.

4. **Deep-link to secondary routes from inspector**  
   “Open task”, “Open review”, “Open Loop”, “Open Seed”, etc. navigate to existing routes for refresh-safe detail / write flows while keeping canvas as home.

5. **Empty / no-center / truncated states**  
   - No center: prompt to search or pick a task — do not call `getGraph` without `center` + `max_nodes`.  
   - Empty neighborhood: “No edges in budget — try another center.”  
   - Truncated: banner when `truncated: true` — “Budgeted neighborhood — not the full project graph.” **No** “load entire graph” CTA.

## Budget rules (Laws 6–7)

| Rule | Value |
|------|-------|
| Required params | `center` + `max_nodes` on every graph load |
| Default / cap | `DEFAULT_MAX=50` / `UI_CAP=100` (keep live values; do not raise unbounded dump) |
| Truncation | Banner when `truncated`; no “load entire graph” primary CTA |
| Progressive expand | Re-center with **same** budget (raise only within UI_CAP via explicit budget control) |
| Graph tech | 2D `@xyflow/react` only — no Three.js / 3D this phase |

## API gap list for S02

| Gap | Kind | Action |
|-----|------|--------|
| `getImpact` `ops.ts` wrapper | **Client glue** | S02 add helper — OpenAPI `/v1/impact` already exists; not a new core API |
| Library-backed HTTP ops for inspector map | **None** | Why/context/graph/search/reviews/entity/task/loop already in OpenAPI + mostly in `ops.ts` |
| `/v1/path`, `listChanges`, `listRegressions` | **Not required** | No IA proof of need — omit; do not invent path API |
| **P32-PORT** | Serve UX (not an API gap) | **S02 still ships #1 min** (friendly `EADDRINUSE` + `--addr` examples) even if API story is `NO-GAPS.md` |

**Explicit:** There are **no library-backed HTTP gaps** for this inspector map. S02 may author `NO-GAPS.md` for **API** **and** must still ship **P32-PORT**. The only product glue called out here is the missing `getImpact` client wrapper.

## Non-goals (this phase / this artifact)

- Writing or rewriting SPA code
- S04 visual craft (typography/motion/density) beyond naming that it follows depth
- Designing port auto-bind / P32-PORT UI copy (S02/S05)
- Inventing `/v1/path` or full-graph dump endpoints
- Cloning UA tours/persona or Graphify AST-only home
