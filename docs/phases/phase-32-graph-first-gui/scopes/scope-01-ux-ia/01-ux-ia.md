# P32-S01-01 — UX IA

## Metadata
- id: P32-S01-01
- todo_ids: [P32-S01-01]
- role: implementer
- skills: [frontend-design, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Author `UX-IA.md` in this folder: graph-home information architecture, inspector depth map (mapped to `/v1`), panel vs route map, Law 6–7 budgets, and an explicit API gap list for S02. **Docs only — no SPA rewrite, no product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [Phase 32 README](../../README.md)
- [00-PLANNER.md](00-PLANNER.md)
- S00 [`RESEARCH.md`](../scope-00-research/RESEARCH.md) (**PASS**, high confidence) — Handoff to S01, Gaps → IA, API reuse map
- Live baseline: `web/src/App.tsx` (index=Overview; Graph=`/graph`), `web/src/layout/Nav.tsx` (8-link CRUD nav), `web/src/screens/Graph.tsx` (`DEFAULT_MAX=50`, `UI_CAP=100`, click=re-center, **no inspector**)
- Phase 29 IA shape (ops-era, **supersede nav-home only**): [`../../phase-29-http-api-browser-gui/scopes/scope-02-ux-ia/UX-IA.md`](../../../phase-29-http-api-browser-gui/scopes/scope-02-ux-ia/UX-IA.md) — borrow section discipline; do **not** keep Overview-as-home

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: proceed without waiting for plan confirmation.

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Artifact | `scopes/scope-01-ux-ia/UX-IA.md` only |
| Product / SPA code | **None** this row |
| Shell strategy | **Hybrid C** — evolve `web/` into graph-home; not a second SPA |
| Graph tech | 2D `@xyflow/react` only; **no** Three.js / 3D |
| Budgets (Laws 6–7) | Always `center` + `max_nodes`; keep live `DEFAULT_MAX=50`, `UI_CAP=100`; truncation honesty banner; **no** full-graph dump default / primary CTA |
| Law 19 | UI formats `/v1` packets only — no browser SQLite, no business-logic fork |
| Port / serve | **Out of S01** — P32-PORT ships in **S02** (#1 min per RESEARCH) |
| Depth vs craft | IA targets **S03 depth**; leave S04 visual craft to later (note only) |

### Must-answer locks (encode these in `UX-IA.md`)

| # | Question | Locked answer |
|---|----------|---------------|
| 1 | Default route / home? | **Graph canvas is home.** Target: index `/` (or equivalent primary route) renders the budgeted XYFlow explorer. Overview is **demoted** (secondary route and/or selection panel) — not the default landing. |
| 2 | Inspector sections? | Ordered depth map (structured UI, not raw JSON dumps as the primary presentation): **(1) Entity summary** (`getEntity` / task row hydrate) → **(2) Why** (`getWhy`) → **(3) Context** (`getContext`, depth 1\|2) → **(4) Impact** (`getImpact`) → **(5) Reviews** (`listReviews` / `getReview`, filter by `task_id` when task-selected) → **(6) Links / relations** (edges from current neighborhood + deep-link to entity routes). Optional **Loop strip** when selection is a task (`getLoopStatus` / `getLoopGate`) as a deep panel, not graph dump. |
| 3 | Screens → panels vs routes? | **Canvas home:** Graph. **Deep panels on selection (preferred entry):** task ops / loop / linked reviews / seed-adjacent actions when relevant to selection. **Remain secondary routes** (reachable, not nav-CRUD-primary): Overview, Tasks list, TaskDetail, Loop, Discoveries, Reviews list, ReviewDetail, Seed, Settings. Nav reweights toward Explore/Graph + Search; do not invent a second SPA. Detail routes stay valid for deep links / refresh. |
| 4 | Selection / search → center → expand? | **Search or task pick** → set `center` → `getGraph(center, max_nodes)`. **Select node** → open/focus inspector (selection highlight). **Expand / re-center** is a distinct affordance (e.g. double-click, explicit Expand control, or “use as center”) — do **not** make every click only re-center with no inspector (today’s Graph behavior is the gap). Progressive expand = new center + same budget; show truncation when `truncated`. |
| 5 | Missing `/v1` for S02? | **Library/OpenAPI:** no new core API required for the inspector map above. **Client glue:** `ops.ts` lacks `getImpact` — list as **S02 client wrapper** (not a new library op). **Do not** invent `/v1/path` or require `listChanges` / `listRegressions` unless IA explicitly proves a library-backed need (default: omit → expect S02 `NO-GAPS.md` for API + still **always** P32-PORT). |

## Preflight (confirm in Notes / UX-IA, do not change code)

1. `App.tsx`: index → Overview; `path="graph"` → Graph (baseline to supersede in IA, not in this row’s code).
2. `Graph.tsx`: `DEFAULT_MAX=50`, `UI_CAP=100`; click re-centers; no inspector.
3. RESEARCH API reuse map still accurate for why/context/impact/graph/search/reviews/entity.

## Role work

Write `UX-IA.md` using the **template sections below** (headings required). Stay within locked defaults. Cite live paths / RESEARCH where claims need evidence.

### UX-IA.md template (required headings)

```markdown
# Phase 32 S01 — UX-IA

**Date:** YYYY-MM-DD
**Row:** P32-S01-01
**Status:** Spec for S03+ implementers — docs only (no `web/` edits this row)

## Summary
(3–6 sentences: graph-home hybrid C; inspector job; budgets; what S02 inherits.)

## Product framing
| Principle | Implication |
|-----------|-------------|
| Hybrid C (DESIGN-LOCKS) | Graph home; ops as deep panels / secondary routes |
| Law 19 | `/v1` adapters only |
| Laws 6–7 | center + max_nodes; UI_CAP; no full-graph dump default |
| Peer bar | Between Graphify and Understand-Anything; keep Trace plan/task/decision/review/loop |
| Out | 3D, second SPA, hosted SaaS, MCP /rpc in browser, P32-PORT design (S02) |

## Shell layout
### Target chrome
(canvas primary + inspector region + secondary ops access)

### Route / panel map
| Surface | Today | Target (home / panel / secondary route) | Notes |
|---------|-------|----------------------------------------|-------|
| Graph | `/graph` | **Home** (`/` or primary) | XYFlow canvas |
| Overview | `/` | Secondary / panel | Demote |
| Tasks / TaskDetail | … | … | … |
| Loop | … | … | … |
| Reviews / ReviewDetail | … | … | … |
| Discoveries | … | … | … |
| Seed | … | … | … |
| Settings | … | … | Keep |
| Nav | 8 CRUD links | Reweight | Explore-first |

## Inspector depth map
| Section (order) | `/v1` op(s) | Client today | S03 notes |
|-----------------|-------------|--------------|-----------|
| Entity summary | getEntity / getTask | … | Header |
| Why | getWhy | ops yes | Structured, not raw `<pre>` primary |
| Context | getContext | ops yes | depth 1\|2 |
| Impact | getImpact | **ops missing** | Flag S02 wrapper |
| Reviews | listReviews / getReview | ops yes | task_id filter |
| Links / relations | getGraph edges (+ deep links) | Graph only | From neighborhood |
| Loop (optional strip) | getLoopStatus / getLoopGate | ops yes | Task selection only |

## Navigation / selection / search flows
1. Search → center → budgeted graph
2. Select → inspector
3. Expand / re-center (distinct from select)
4. Deep-link to secondary routes from inspector
5. Empty / no-center / truncated states (honesty copy)

## Budget rules (Laws 6–7)
| Rule | Value |
|------|-------|
| Required params | center + max_nodes |
| Default / cap | DEFAULT_MAX=50 / UI_CAP=100 (keep unless UX-IA justifies change **without** raising unbounded dump) |
| Truncation | Banner when truncated; no “load entire graph” CTA |
| Progressive expand | Re-center with same budget |

## API gap list for S02
| Gap | Kind | Action |
|-----|------|--------|
| getImpact ops.ts wrapper | Client glue | S02 add helper (OpenAPI already exists) |
| (others) | … | Only if library-backed and proven here |
| P32-PORT | Serve UX | **Not an API gap** — S02 still ships #1 min even if this table is otherwise empty / NO-GAPS |

If no library-backed HTTP gaps: state explicitly that S02 may author `NO-GAPS.md` for API **and** must still ship P32-PORT.

## Non-goals (this phase / this artifact)
- Writing or rewriting SPA code
- S04 visual craft (typography/motion/density) beyond naming that it follows depth
- Designing port auto-bind / P32-PORT UI copy (S02/S05)
- Inventing `/v1/path` or full-graph dump endpoints
- Cloning UA tours/persona or Graphify AST-only home
```

## Exit criteria

- [ ] `UX-IA.md` exists with **all** template headings above
- [ ] Must-answer locks (1–5) encoded (graph home; inspector map; panel/route; select vs expand; API/P32-PORT gaps)
- [ ] Laws 6–7 budgets explicit; Law 19 stated
- [ ] API gap list explicit (may be client-only / empty library gaps) + P32-PORT reminder
- [ ] **No** product code changes
- [ ] Board Notes cite artifact path

## Minimal todos

- [ ] Confirm preflight live facts in Notes
- [ ] Draft `UX-IA.md` from template + locks
- [ ] Map inspector sections → `/v1` + client gaps
- [ ] List API gaps / NO-GAPS expectation + P32-PORT still required
- [ ] Update board row **P32-S01-01** status + Notes only

## Todo updates

Status + notes on **P32-S01-01** only. Do not start **P32-S01-02**. Do not write product code.

## Next

`P32-S01-02`
