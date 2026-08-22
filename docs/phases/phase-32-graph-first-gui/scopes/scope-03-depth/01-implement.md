# P32-S03-01 — Implement depth

## Metadata
- id: P32-S03-01
- todo_ids: [P32-S03-01]
- role: implementer
- skills: [frontend-ui-engineering, incremental-implementation, vercel-react-best-practices]
- mcps: []
- verification: automated
- hooks: []

## Objective

Evolve existing `web/` into **hybrid C graph-home** + structured **inspector depth** per [`UX-IA.md`](../scope-01-ux-ia/UX-IA.md): summary → why → context → impact → reviews → links (+ optional loop strip on tasks). **Select opens inspector; expand/re-center is a distinct affordance** (today’s click→re-center only is the gap). `/v1` adapters only (Law 19). Keep `DEFAULT_MAX=50` / `UI_CAP=100` + truncation honesty (Laws 6–7). **Do not invent `/v1/path`.** Wire Impact via shipped S02 `getImpact(taskId, opt?)`. **Defer visual craft to S04** (functional layout OK; no typography/motion polish pass).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [00-PLANNER.md](00-PLANNER.md) — Must-answer locks below are **final**
- S01 [`UX-IA.md`](../scope-01-ux-ia/UX-IA.md) — shell + inspector depth map (PASS)
- S02 shipped: `getImpact` in `web/src/api/ops.ts` → `GET /v1/impact?task_id=`; `NO-GAPS.md` (no new library `/v1`)
- Live: `web/src/App.tsx`, `web/src/layout/Nav.tsx`, `web/src/screens/Graph.tsx`, `web/src/screens/TaskDetail.tsx` (why/context drawer patterns), `web/src/api/ops.ts`, `web/e2e/s05-gates.spec.ts` (smoke nav `/graph`)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: proceed without waiting for plan confirmation.

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Law 19 | `/v1` via `ops.ts` only — format library packets; no browser SQLite; no business-logic fork |
| Laws 6–7 | Every graph load: `center` + `max_nodes`; DEFAULT_MAX=50 / UI_CAP=100; truncation banner; **no** “load entire graph” CTA |
| Stack | Existing Vite/React/`@xyflow/react` — **no** Three.js / 3D |
| Home | Graph canvas is **index** (`/`); Overview demoted to `/overview` |
| `/graph` | Keep as **alias** → same Graph home (redirect or same element) so bookmarks + existing e2e still work |
| Select vs expand | **Distinct** — see locks below |
| Path API | **Out** — do not invent `/v1/path` |
| Impact | Import existing `getImpact(taskId, opt?)` — do **not** re-wrap |
| Kind/search filters | **Keep** on Graph (already live) |
| Visual | Functional depth shell only — S04 owns craft |

### Must-answer locks (encode in implementation + Notes)

| # | Question | Locked answer |
|---|----------|---------------|
| 1 | Graph-home file/route changes? | See **A — Graph-home shell** below. |
| 2 | Inspector data loading map? | See **B — Inspector load map**. |
| 3 | E2E/smoke minimum? | See **C — Build + smoke**. |

### A — Graph-home shell (concrete files)

| File | Change |
|------|--------|
| `web/src/App.tsx` | `index` → `<Graph />`; add `path="overview"` → `<Overview />`; keep `path="graph"` as `<Navigate to="/" replace />` **or** same `<Graph />` (prefer redirect so one home URL) |
| `web/src/layout/Nav.tsx` | Reweight `LINKS`: **Explore** (or Graph) `{ to: '/', end: true }` **first**; Overview → `/overview` demoted (after Explore / before or among ops); keep Tasks, Loop, Reviews, Discoveries, Seed, Settings as secondary |
| `web/src/screens/Graph.tsx` | Evolve into home shell: search → center → budgeted XYFlow + **inspector region**; stop click-only-recenter |
| `web/src/components/Inspector.tsx` (new, preferred) | Structured sections per UX-IA order; deep-links to `/tasks/:id`, `/reviews/:id`, `/loop?task_id=` |
| CSS | Minimal layout (split pane / sheet) in existing stylesheet — **no** S04 visual polish |

**Chrome:** Shell brand already links `/` — after this row that lands on graph home (good). Do not invent a second SPA.

### B — Inspector load map (Law 19)

On **selection** (not on expand), load sections for the selected node. Prefer parallel fetches; show per-section loading/error; structured UI over raw `<pre>` as primary (collapsible raw dump OK as secondary).

| Order | Section | When | `/v1` / client | Notes |
|------:|---------|------|----------------|-------|
| 1 | Summary | always | `getEntity(id)`; if `kind === 'task'` also `getTask(id)` | Header: kind, title, id, work_state when task |
| 2 | Why | always | `getWhy(entityType, id)` — `entityType` from node/entity `kind` | Structured Why packet |
| 3 | Context | **task only** | `getContext(taskId, depth)` depth **1\|2** control | Hide/collapse for non-task |
| 4 | Impact | **task only** | **`getImpact(taskId)`** (S02) | Omit/honest empty for non-task — do not invent decision_id client |
| 5 | Reviews | **task only** | `listReviews({ task_id })` | Deep-link `getReview` / `/reviews/:id` |
| 6 | Links | always | Edges from **current** `getGraph` neighborhood involving selected id + deep-links | No `/v1/path` |
| — | Loop strip | **task only** | `getLoopStatus` / `getLoopGate` | Compact strip + “Open Loop” → `/loop?task_id=` |

**Out:** `listChanges` / `listRegressions`; inventing path API; full-graph dump.

### Select ≠ expand (behavior lock)

| Affordance | Behavior |
|------------|----------|
| **Select** (single-click node, or equivalent) | Set `selectedId`; highlight; open/focus inspector; load depth map. **Does not** call `getGraph` / change `center` |
| **Expand / re-center** | Explicit only: double-click **or** inspector/control “Use as center” / “Expand”. Calls `loadGraph(id)` with **same** budget (cap ≤ UI_CAP). If `truncated`, show honesty banner |
| Today’s gap | `Graph.tsx` `onNodeClick` → `onExpand` when `id !== center` — **replace** with select; move expand to distinct control |

### C — Build + smoke

1. `cd web && npm run build` (or `npx tsc -b` + vite build per package scripts) **PASS**.
2. E2E minimum (Playwright under `web/e2e/`):
   - Update smoke that assumes Overview-at-`/` / Graph-only-at-`/graph` (`s05-gates.spec.ts` smoke nav): home or `/graph` alias still shows graph controls (`#graph-budget` or successor).
   - Add **depth smoke** (extend smoke or `e2e/s03-depth.spec.ts`): pick task or search → load neighborhood → **single-click** node → inspector visible with at least Summary (+ Why if fetchable) → assert **center unchanged** on select → trigger **expand** → center updates (or graph reload for new center).
3. If Playwright binary missing in env: document in Notes; still ship build PASS + manual checklist evidence. Prefer automated when `test:e2e` works.

## Preflight (confirm in Notes; then change)

1. `App.tsx`: index = Overview; Graph = `/graph` only.
2. `Nav.tsx`: Overview first among 8 equal links.
3. `Graph.tsx`: click non-center → re-center; **no** inspector; budgets 50/100.
4. `ops.ts`: `getImpact` **present** (S02); why/context/graph/search/reviews/entity/task/loop present.
5. No `/v1/path` in OpenAPI — do not add.

## Role work

1. **Routes + Nav** — graph home per table A; Overview secondary.
2. **Graph shell** — select vs expand; keep search/kind/`max_nodes`; truncation banner; no full-dump CTA.
3. **Inspector** — sections B in order; task-scoped omit rules; Impact via S02 `getImpact`.
4. **Deep-links** — Open task / review / Loop (and Seed only if selection-relevant — optional).
5. **Verify** — build + e2e/smoke per C.
6. **Board** — status + Notes only on **P32-S03-01** (cite files + commands).

## Exit criteria

- [ ] `/` is graph-home (hybrid C); Overview at `/overview` (or equivalent demotion); `/graph` still reaches graph
- [ ] Inspector sections per UX-IA order; task-scoped omit for Context/Impact/Reviews/Loop
- [ ] Select ≠ expand/re-center (evidence in Notes or e2e)
- [ ] Budgets DEFAULT_MAX=50 / UI_CAP=100; truncation honesty; no full-dump CTA
- [ ] Impact uses existing `getImpact`; **no** `/v1/path`; Law 19 adapters only
- [ ] `npm run build` (web) PASS; e2e/smoke per C (or documented env blocker + manual checklist)
- [ ] Visual polish **not** required (S04); Board Notes cite evidence

## Minimal todos

- [ ] App routes + Nav reweight (graph home; Overview demoted; `/graph` alias)
- [ ] Graph: select → inspector; distinct expand/re-center; keep budgets/filters
- [ ] Inspector component: summary → why → context → impact → reviews → links (+ loop strip)
- [ ] Wire ops including S02 `getImpact`; deep-links to secondary routes
- [ ] Build + e2e/smoke update
- [ ] Update board row **P32-S03-01** only

## Todo updates

Status + notes on **P32-S03-01** only. Do not start **P32-S03-02**.

## Next

`P32-S03-02`
