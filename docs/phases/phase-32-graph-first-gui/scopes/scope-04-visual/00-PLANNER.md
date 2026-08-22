# P32-S04-00 — Scope planner (visual craft)

## Metadata
- id: P32-S04-00
- todo_ids: [P32-S04-00]
- role: planner
- skills: [planning-and-task-breakdown, frontend-design, impeccable]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock visual craft pass on the **depth shell** from S03: typography, layout density, motion, canvas chrome. Must not strip inspector depth. No Three.js.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- Live: `web/src/` after S03
- S03 review residuals: dense PacketView/layout craft; e2e list-vs-canvas nit

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Base | S03 depth shell — evolve, do not replace with ops-nav shell |
| Live anchors (S03 PASS) | `App.tsx` graph `@/`; `Nav` Explore-first; `Graph.tsx` + `Inspector.tsx`; CSS `.graph-shell` / `.graph-inspector` in `web/src/styles/app.css`; e2e `e2e/s03-depth.spec.ts` |
| Depth must keep | Inspector order summary→why→context→impact→reviews→links (+loop tasks); select≠expand; budgets 50/UI_CAP 100; `getImpact` only; no `/v1/path` |
| Brand | Evolve IBM Plex + forest/sage tokens — do not invent new brand |
| 3D | Forbidden as default |
| A11y | Focus order, named landmarks, contrast, `prefers-reduced-motion` |

## Must answer (handoff to 01) — LOCKED

1. **Layout / shell:** Canvas-first `.graph-shell`; taller viewport-aware canvas; inspector ~`minmax(18rem, 26rem)`; sticky chrome; stack on narrow.
2. **Typography + PacketView:** Explicit inspector type roles; tighten structured `dl` density; raw JSON secondary in `<details>` (S03 residual).
3. **Chrome + motion:** Calm center/selected nodes; intentional empty state; **2–3** motions with reduced-motion; optional e2e canvas-click nit.

## Planner gate

- [x] `01-implement.md` / `02-review.md` runnable
- [x] Explicit: do not regress depth (re-run s03-depth + s05 smoke after craft)
- [x] Craft targets: inspector density/PacketView, graph-shell split, canvas chrome — not route reshuffle

## Exit criteria

- [x] Next **P32-S04-01** runnable (thickened)

## Todo updates

Status + notes on **P32-S04-00** only.

## Next

`P32-S04-01`
