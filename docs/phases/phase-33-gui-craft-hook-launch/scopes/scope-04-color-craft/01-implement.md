# P33-S04-01 — Colorize / craft shell

## Metadata
- id: P33-S04-01
- todo_ids: [P33-S04-01]
- role: implementer
- skills: [impeccable, ui-ux-pro-max, frontend-design]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Land S01 DESIGN token table (kind + state + any missing semantic) into CSS and colorize the **full shell** so Explore reads as Trace forest-moss with **kind chroma strip** signature — not gray/bland. Skills required before large style edits. Visual/CSS only (Law 19).

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) Theme A
- [00-PLANNER.md](00-PLANNER.md) locked defaults + must-answer
- [`../scope-01-design-ux/DESIGN.md`](../scope-01-design-ux/DESIGN.md) — **token SoT** (do not invent competing palette)
- Live: `web/src/styles/tokens.css`, `web/src/styles/app.css`
- Hooks ready (S03 PASS): `Graph.tsx` list + canvas `data-kind` / `data-state`, `.graph-node__kind` / `__state`, `--seed` / `--center` / `--selected`
- Skills: impeccable **colorize** / **bolder** / **polish** + craft-floor; ui-ux-pro-max `--domain color|ux`; frontend-design signature = chroma strip

## Session start

Follow agent-loop-protocol Session start. **Must** load listed skills (impeccable colorize + craft-floor bans; ui-ux-pro-max contrast/color-not-only; frontend-design signature restraint) before large CSS edits. Do **not** start P33-S04-02 work.

## Locked defaults

| Item | Value |
|------|-------|
| Token SoT | DESIGN.md table → `tokens.css` light+dark: keep existing surfaces/accent/focus/semantic; **add** full `--kind-*` + `--state-*` hexes from DESIGN |
| Signature | Kind **chroma strip**: `border-left: 3px solid var(--kind-*)` on `.graph-node__inner` / list rows; soft fill `color-mix(in srgb, var(--kind-*) 14%, var(--bg-elevated))`; mono uppercase `.graph-node__kind` colored with `--kind-*`; **no** glow bloom / box-shadow halos |
| State channel | `.graph-node__state[data-state]` (and list state text) → `--state-*` text/tick; labeled always (`IN_PROGRESS`, etc.) |
| Center/seed | Keep moss `--accent` / `--accent-subtle` emphasis; **kind label still shown**; do not erase kind strip under center |
| Selected | `--focus` / accent outline; keep kind border visible (dual-ring OK; no glow) |
| Scope priority | **1)** Explore `/` — canvas nodes + node list + inspector chrome **2)** App chrome (nav active, theme, banners) **3)** Shared `.pill` / empty / error **4)** Other screens inherit via tokens — no IA rewrite |
| Files (expected) | `web/src/styles/tokens.css`, `web/src/styles/app.css`; optional tiny class hooks in `ErrorBanner` / list markup **only if** CSS attrs insufficient — prefer CSS-only |
| Forbidden edits | `overviewCompose.ts`, `loadOverview` path, `App.tsx` routes, budgets, S02 CLI; do not treat `/overview` as Explore craft target |
| Contrast | Text ≥4.5:1; strokes/focus ≥3:1; prove light **and** dark in Notes (checker URL or computed pairs for `--text`/`--bg`, `--text`/`--bg-elevated`, sample `--kind-task` label on soft fill, `--focus` on `--bg`) |
| Motion | Keep `--duration-*`; extend `prefers-reduced-motion` if new transitions added |
| Avoid | Purple/cream/broadsheet; HUD/glow; emoji; wholesale slate+#22C55E; rainbow tipping (capability/assumption stay muted + labeled) |
| Optional polish | Stronger inline Tasks CTA on no-seeds EmptyState — **visual/copy only** (S03 owns message substance) |

## Token application map (implement this)

| Target | Selectors / hooks | Tokens |
|--------|-------------------|--------|
| Kind strip + fill | `.react-flow__node.graph-node .graph-node__inner[data-kind='…']`, list `button[data-kind]` | `--kind-*` border-left 3px + `color-mix` 14% fill; default stroke may stay `--border` on other sides |
| Kind label | `.graph-node__kind` under `[data-kind]` | `color: var(--kind-*)`; keep mono uppercase |
| State chip | `.graph-node__state[data-state='…']` | `color: var(--state-*)` (optional right tick via border-inline-end) |
| Seed / center | existing `--seed` / `--center` rules | Moss accent overlay **on top of** kind strip — do not remove `data-kind` styling entirely |
| Selected | `--selected` + `:focus-visible` | `--focus` / accent outline ≥3:1 |
| Pills | `.pill` where kind/state known; `.pill--ok/bad/warn` | Align with `--ok`/`--danger`/`--warn`; kind-colored pills only if markup has kind |
| Empty / error | `.empty`, `.banner--error` | `--danger` / `--danger-bg`, `--warn` pairs; text ≥4.5:1 |
| Nav / chrome | `.nav a[aria-current]`, header | `--accent` / `--accent-subtle`; forest neutrals only |
| Fallback | unknown kind | `--kind-unknown` (= muted) |

Use attribute selectors for each DESIGN kind (`goal`, `task`, `decision`, `assumption`, `discovery`, `plan-change`, `claim`, `evidence`, `capability`, `review`) and states (`IN_PROGRESS`, `PENDING`, `DONE`, `SKIPPED`→done, `FAIL`, `BLOCKED`). Map SKIPPED → `--state-done`.

## Skills gate (evidence in Notes)

1. **impeccable:** colorize (roles) + bolder hierarchy without glow + polish; craft-floor bans honored.
2. **ui-ux-pro-max:** reconfirm `--domain ux` Color Contrast + Color Only; reject `--domain color` slate+#22C55E wholesale (S01 already rejected).
3. **frontend-design:** one signature (chroma strip); chrome stays quiet moss-neutral.

## Role work

1. Add `--kind-*` / `--state-*` to `tokens.css` (light + dark) exactly from DESIGN candidates.
2. Wire `app.css` per application map — Explore first, then chrome/pills/empty.
3. Preserve kind + state **text** (color-not-only); do not delete labels for color-only encoding.
4. Verify `prefers-reduced-motion` still disables node/inspector transitions.
5. Evidence: light+dark screenshots (Explore with ≥2 kinds visible) + contrast table in board Notes.
6. Build sanity: `npm run build` in `web/` (and keep S03 unit/e2e green if touched only CSS — no compose changes expected).

## Exit criteria

- [ ] `--kind-*` and `--state-*` present in `tokens.css` light+dark
- [ ] Explore nodes/list show chroma strip + kind label color; state chips colored + labeled
- [ ] Shell not gray-monochrome; forest-moss identity intact
- [ ] Contrast floors evidenced light **and** dark
- [ ] Focus ring visible; reduced-motion respected
- [ ] No compose/route/budget edits
- [ ] Skills evidence in Notes + visual evidence paths
- [ ] Board Notes complete → next **P33-S04-02**

## Minimal todos

- [ ] Land kind/state tokens in `tokens.css`
- [ ] Colorize Explore canvas + list via `data-kind`/`data-state` (signature strip)
- [ ] Chrome / pills / empty-error pass
- [ ] Contrast check light+dark + screenshots
- [ ] Skills evidence + board Notes

## Todo updates

Status + notes on **P33-S04-01** only.

## Next

`P33-S04-02`
