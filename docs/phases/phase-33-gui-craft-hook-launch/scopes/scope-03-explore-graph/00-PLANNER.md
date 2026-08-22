# P33-S03-00 — Scope planner (Explore hook graph)

## Metadata
- id: P33-S03-00
- todo_ids: [P33-S03-00]
- role: planner
- skills: [frontend-design, impeccable, ui-ux-pro-max, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Finalize S03 so Explore (`/`) ships the **interactive project overview graph hook** per S01 UX-IA + S00 budget rules. Keep `@xyflow/react` 2D. Skills on implement+review. Full shell colorize remains **S04** (S03 may apply token placeholders / kind colors).

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- `../scope-00-research/RESEARCH.md`
- `../scope-01-design-ux/` design artifacts
- Live: `web/src/screens/Graph.tsx`, `web/src/App.tsx`, `web/src/api/ops.ts`, `api/openapi.yaml`

## Session start

Follow agent-loop-protocol Session start.

## Locked defaults

| Item | Value |
|------|-------|
| Surface | Evolve existing Explore/`Graph.tsx` — not a second SPA |
| IA SoT | [`../scope-01-design-ux/UX-IA.md`](../scope-01-design-ux/UX-IA.md) — model **(D)+(B)+(C)**; Explore=`/` ≠ `/overview` |
| CLI land (S02) | `trace gui` opens **`http://{addr}/`** after listen — **keep Explore at `/`**; do not move the graph home to `/overview` (would break open URL without a CLI change; out of S03) |
| Seed pipeline | `getProject` → `listTasks` (IN_PROGRESS then active) → `search` fill → dedupe; **target 6, ≤8** |
| Budgets | Per-seed `getGraph(max_nodes=40, depth=2)`; merge honor **`UI_CAP=100`**; expand re-center `max_nodes≤50`; no load-all |
| API | **reuse** `ops.ts` (`getProject` / `listTasks` / `search` / `getGraph`); no seed-export-as-graph-body; no Leiden |
| Keyboard | Chrome tab order + **canvas selection keyboard path** (see Must answer #5); visible focus on selected node |
| S04 hooks | Attach `data-kind` / `data-state` (or equiv.); do **not** invent palette — tokens in DESIGN.md |
| Skills | frontend-design + impeccable + ui-ux-pro-max on 01 and 02 |
| Craft depth | Hook + usable coloring OK; S04 owns shell-wide colorize/bolder |
| Stack | xyflow 2D only |

## Must answer (handoff to 01) — resolved 2026-08-21

1. **Data path:** Mount `/` → loading → (D) `getProject` chrome only (no fake node; live `ProjectResponse` lacks graph id) + `listTasks` priority (`IN_PROGRESS` → active → DONE/SKIPPED only if needed) + non-empty `search("goal"|"capability"|"decision"|"discovery")` fill → dedupe ≤8 (target 6) → (B) parallel `getGraph(seed, 40, 2)` → merge/trim **UI_CAP=100** (seeds first) → first paint → (C) user expand ≤50 replace-prefer.
2. **Interaction:** Pan/zoom/click; expand via double-click / Expand / Use as center; smoke = overview nodes without mandatory “Pick center” EmptyState (empty store → UX-IA no-seeds copy).
3. **States:** Loading / happy / no-seeds / partial (subgraph + Retry) / hard-fail per UX-IA table — copy locked in `01-implement.md`.
4. **Files:** Primary `web/src/screens/Graph.tsx`; optional `web/src/lib/overviewCompose.ts` (or equiv.); hook styles only; no `App.tsx` route move; no `Overview.tsx` redesign.
5. **Keyboard:** Required chrome → node-list → inspector (visible focus, no trap); canvas best-effort; residual risk OK if documented in implementer Notes.

## Planner gate

- [x] `01-implement.md` runnable with exit criteria + e2e smoke
- [x] `02-review.md` Laws 6–7 + skills checklist
- [x] `SCOPE-TODOS.md` accurate

## Exit criteria

- [x] Implementer locked; next **P33-S03-01**

## Todo updates

Status + notes on **P33-S03-00** only.

## Next

`P33-S03-01`
