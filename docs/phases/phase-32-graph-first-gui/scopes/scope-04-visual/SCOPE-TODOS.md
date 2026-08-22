# Scope 04 — board map

Visual craft on the **S03 depth shell** (layout, typography, PacketView density, motion, canvas chrome). Serial: **S04-00 → S04-01 → S04-02**. Do not strip inspector depth; do not change IA/API contracts; no Three.js.

| Board ID | Row | Prompt | Role |
|----------|-----|--------|------|
| 561 | P32-S04-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner |
| 562 | P32-S04-01 | [01-implement.md](01-implement.md) | Implementer |
| 563 | P32-S04-02 | [02-review.md](02-review.md) | Reviewer |

## Craft locks (from S04-00)

| Axis | Lock |
|------|------|
| Brand | Evolve IBM Plex + forest/sage tokens — do not replace |
| Layout | Canvas-first `.graph-shell`; taller canvas; inspector ~18–26rem |
| Density | `PacketView` structured primary, tighter rows; raw JSON in `details` |
| Motion | 2–3 intentional + `prefers-reduced-motion` |
| Depth | Preserve S03 order / select≠expand / budgets / `getImpact` / no `/v1/path` |
| Out | No Three.js; no route reshuffle; no ops-nav shell |

## Regression gate

`npm run build` + `e2e/s03-depth.spec.ts` + `e2e/s05-gates.spec.ts` — must stay green after craft.

## S03 residuals owned here

1. Dense PacketView / layout craft (primary)
2. E2e list-vs-canvas select nit (optional; not blocker)
