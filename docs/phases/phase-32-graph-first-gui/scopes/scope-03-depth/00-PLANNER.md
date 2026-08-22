# P32-S03-00 — Scope planner (depth)

## Metadata
- id: P32-S03-00
- todo_ids: [P32-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, frontend-design, incremental-implementation]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock depth implement/review so S03-01 evolves `web/` into graph-home + inspector depth per `UX-IA.md`: summary → why → context → impact → reviews → links (+ optional loop strip); **select ≠ expand/re-center**; Laws 6–7 budgets. **Do not invent `/v1/path`.** Kind/search filters already on Graph stay; richer path only if library-backed later. Visual polish deferred to S04. **No deep visual craft here.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- `../scope-01-ux-ia/UX-IA.md` (PASS after P32-S01-02)
- Live: `web/src/`
- Depends: S02 **shipped** `getImpact(taskId, opt?)` in `web/src/api/ops.ts` → `GET /v1/impact?task_id=` (P32-S02-01/02); wire inspector Impact section to it — do not re-invent

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Shell | Evolve existing `web/` — graph home (index); Overview demoted |
| API | `/v1` only (Law 19); reuse UX-IA op map |
| Graph | Budgeted `@xyflow/react`; DEFAULT_MAX=50 / UI_CAP=100; no 3D |
| Selection | Select → inspector; expand/re-center = distinct affordance |
| Path | **Out** unless library gains `/v1/path` — omit inventing |
| Visual polish | Defer to S04 |

## Must answer (handoff to 01) — LOCKED

1. **Graph-home file/route:** `App.tsx` index → `<Graph />`; Overview at `/overview`; `/graph` alias/redirect to home; `Nav.tsx` Explore/Graph first; evolve `Graph.tsx` + new `components/Inspector.tsx` (preferred).
2. **Inspector load map:** summary=`getEntity`(+`getTask` if task) → why=`getWhy` → context=`getContext` (task) → impact=`getImpact` (task, S02) → reviews=`listReviews({task_id})` (task) → links=neighborhood edges → loop strip=`getLoopStatus`/`getLoopGate` (task). Non-task: omit/collapse Context/Impact/Reviews/Loop.
3. **E2E/smoke:** update nav smoke for graph-home/`/graph` alias; depth smoke: pick/search → center → select → inspector; select does not change center; expand does. `npm run build` required.

## Planner gate

- [x] `01-implement.md` runnable with exit criteria
- [x] `02-review.md` depth checklist (not visual)

## Exit criteria

- [x] Next **P32-S03-01** runnable (thickened)

## Todo updates

Status + notes on **P32-S03-00** only.

## Next

`P32-S03-01`
