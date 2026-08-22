# P32-S04-01 — Implement visuals

## Metadata
- id: P32-S04-01
- todo_ids: [P32-S04-01]
- role: implementer
- skills: [frontend-design, improve-animations, ui-styling]
- mcps: []
- verification: automated
- hooks: []

## Objective

Craft pass on the **S03 depth shell** only: typography, layout density, motion, canvas chrome. Evolve existing tokens/CSS — do **not** replace brand, strip inspector depth, reshuffle routes/IA, or invent API clients. **No Three.js.** Keep a11y (focus, named landmarks, contrast, `prefers-reduced-motion`).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [00-PLANNER.md](00-PLANNER.md) — craft locks below are **final**
- S01 [`UX-IA.md`](../scope-01-ux-ia/UX-IA.md) — shell/IA unchanged
- S03 review residuals → this scope: dense `PacketView` / layout craft; e2e list-vs-canvas select nit
- Live anchors (S03 PASS): `web/src/App.tsx`, `layout/Nav.tsx`, `screens/Graph.tsx`, `components/Inspector.tsx` (`PacketView`), `styles/tokens.css`, `styles/app.css` (`.graph-shell` / `.graph-inspector` / `.inspector-dl*`), `e2e/s03-depth.spec.ts`, `e2e/s05-gates.spec.ts`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: proceed without waiting for plan confirmation.

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Scope | Visual craft **only** on existing explorer depth shell |
| Brand | **Evolve** live tokens — IBM Plex Sans/Mono; forest/sage accent (`#2f5d3a` / dark twin); do **not** invent a new brand or swap to Inter/purple/glow/cream-serif clusters |
| Depth | Must remain: section order summary→why→context→impact→reviews→links (+loop tasks); task-scoped omit; select≠expand; budgets DEFAULT_MAX=50 / UI_CAP=100; Impact=`getImpact` only; no `/v1/path` |
| Routes / IA | **No** graph-home demotion; **no** Nav re-equalize to CRUD-first; **no** ops-nav shell swap |
| Law 19 | Format `/v1` packets in UI only — no browser SoT / business-logic fork |
| Graph | 2D `@xyflow/react` — **no** Three.js / 3D default |
| Regression gate | `cd web && npm run build` + `npm run test:e2e -- e2e/s03-depth.spec.ts e2e/s05-gates.spec.ts` still **PASS** |

### Must-answer craft locks (encode in implementation + Notes)

| # | Axis | Locked target |
|---|------|---------------|
| 1 | **Layout / shell** | See **A** |
| 2 | **Typography + PacketView density** | See **B** |
| 3 | **Canvas chrome + motion + a11y** | See **C** |

### A — Layout / graph-shell split

| Surface | Craft target |
|---------|--------------|
| `.graph-shell` | Canvas-first split: canvas column dominates; inspector ~`minmax(18rem, 26rem)` (or equivalent token). Avoid equal dual “dashboard cards.” |
| `.graph-canvas` | Taller / more viewport-aware than today’s `min(32rem, 70vh)` — explorer feel; still scrolls safely on short viewports. |
| `.graph-inspector` | Readable max-height aligned to canvas; sticky `__chrome`; section rhythm with clear separation **without** card-stack chrome. |
| Narrow (`max-width` ~56rem) | Stack canvas then inspector; no horizontal crush of dl columns. |
| Out | Do not move Explore off `/`; do not add a second SPA layout. |

### B — Typography + PacketView density (S03 residual)

| Surface | Craft target |
|---------|--------------|
| Type scale | Explicit roles in CSS/tokens: inspector section title (small caps / tracked muted), body/value readable, mono keys/ids quieter and smaller than values. Keep IBM Plex pair. |
| `PacketView` | Primary = structured `dl` / field rows with **tighter density** (less sparse padding; nested objects summarized cleanly). Raw JSON stays in `<details>` — secondary, not the default reading path. |
| Walls of JSON | Prefer truncated structured rows over large always-open `<pre>` as primary. Cap lengths may stay; presentation must not feel like a debug dump. |
| Section loading/error | Compact, same visual language as banners — do not balloon vertical rhythm. |

### C — Canvas chrome, motion, a11y

| Surface | Craft target |
|---------|--------------|
| Nodes | Center vs selected states distinct and calm (outline/border using existing accent) — **no** multi-layer glow / neon. |
| Empty inspector | Intentional quiet empty state (short copy + muted type), not a blank card void. |
| Motion (ship **2–3** intentional) | e.g. (1) inspector content appear / selection settle, (2) node selected/center transition, (3) section open or sticky chrome — short, purposeful. Honor `prefers-reduced-motion: reduce` (instant or opacity-only). |
| Focus / landmarks | Visible `:focus-visible`; inspector remains a named region (`aria-*` / landmark already or improve); contrast on any new surfaces ≥ existing token floor. |
| E2E nit (optional) | Prefer adding a **canvas node click** select path in `s03-depth` if cheap (same `onSelect`); list pick remains valid — **not** a craft blocker if list-only stays green. |

### Out of scope (defer / forbid)

- Three.js / 3D; WebGL decorative backgrounds
- Route/IA/API contract changes; new `/v1` clients beyond formatting
- Full redesign of ops screens (Tasks/Loop/Reviews/…); light token inheritance OK if shared CSS improves
- Replacing brand fonts/palette wholesale
- “Load entire graph” CTA or budget relaxation

## Preflight (confirm in Notes; then change)

1. Graph home still `/` + `/graph` alias; Nav Explore-first.
2. `Inspector.tsx` depth order + `PacketView` + `getImpact` intact.
3. `tokens.css` / `app.css` are the craft surfaces; no new design-system package.

## Exit criteria

- [ ] Craft locks A/B/C landed without depth/IA/API regression
- [ ] A11y basics held (focus, landmarks, contrast, reduced-motion)
- [ ] `npm run build` PASS
- [ ] `npm run test:e2e -- e2e/s03-depth.spec.ts e2e/s05-gates.spec.ts` PASS
- [ ] Board Notes cite evidence + which 2–3 motions shipped

## Minimal todos

- [ ] Tokens / typography scale + graph-shell / canvas / inspector layout
- [ ] PacketView density + section rhythm
- [ ] Canvas chrome (nodes, empty state)
- [ ] 2–3 motions + `prefers-reduced-motion`
- [ ] Verify build + s03-depth + s05-gates
- [ ] Update board row **P32-S04-01** only

## Todo updates

Status + notes on **P32-S04-01** only.

## Next

`P32-S04-02`
