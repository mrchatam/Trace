# P33-S04-00 — Scope planner (colorize / craft)

## Metadata
- id: P33-S04-00
- todo_ids: [P33-S04-00]
- role: planner
- skills: [impeccable, ui-ux-pro-max, frontend-design, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Finalize S04 to apply S01 tokens across the **full shell** (nav, panels, graph nodes, empty/error) via impeccable **colorize** + **bolder/polish** + craft-floor. Skills required on implement+review. No bland gray default.

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- S01 `DESIGN.md` / tokens
- Live: `web/src/styles/app.css`, `web/src/layout/`, `web/src/screens/`, `web/src/components/`

## Session start

Follow agent-loop-protocol Session start.

## Locked defaults

| Item | Value |
|------|-------|
| Token SoT | [`../scope-01-design-ux/DESIGN.md`](../scope-01-design-ux/DESIGN.md) — land full table (surfaces, `--kind-*`, `--state-*`, focus, semantic) into `tokens.css` / `app.css` |
| Signature | Kind **chroma strip** (3px left border + mono uppercase kind label); soft fill `color-mix` ~14%; no glow bloom |
| Skills | impeccable colorize/bolder/polish + ui-ux-pro-max + frontend-design |
| Scope | Shell-wide craft on existing Explore depth from S03 (**PASS** P33-S03-02) |
| Explore hooks (do not rework) | Nodes already expose `data-kind` / `data-state`, `.graph-node__kind` / `__state`, `.graph-node--seed` / `--center` / `--selected`. Colorize via CSS attrs/classes — **no** seed-compose or routing changes |
| Home contract | Explore stays at **`/`** (S02 `trace gui` land); `/overview` remains ops Overview |
| A11y | Text ≥4.5:1 light+dark; strokes/focus ≥3:1; verify with contrast checker evidence in Notes |
| Avoid | Purple/cream/broadsheet; HUD/glow; emoji; wholesale slate+#22C55E; rainbow tipping on kind spectrum |
| Out | Redesigning IA / overview compose (S01/S03); canvas arrow-roving (accepted residual); launch CLI (S02) |

## Must answer (handoff to 01) — LOCKED

1. **Token map:** Land DESIGN `--kind-*`/`--state-*` into `tokens.css` (surfaces already present). Wire `app.css`: `[data-kind]` → 3px left strip + 14% `color-mix` fill + kind-label color; `[data-state]` → state chip color; `.empty`/`.banner--error` → semantic; `.nav`/pills → accent/ok/warn/danger. Prefer Graph hooks already shipped.
2. **Priority:** Explore `/` (canvas + list + inspector) → app chrome → shared empty/error → other screens inherit tokens.
3. **Verify:** Light+dark screenshots (Explore ≥2 kinds) + contrast pairs in Notes (≥4.5:1 text, ≥3:1 stroke/focus).
4. **Motion:** Yes — keep/extend existing `prefers-reduced-motion` block in `app.css`.
5. **Color-not-only:** Kind + state **text** stay; color is redundant. Never remove `.graph-node__kind` / `__state` labels.

## Planner gate

- [x] `01-implement.md` + `02-review.md` thickened with skills + exit criteria
- [x] `SCOPE-TODOS.md` accurate

## Exit criteria

- [x] Implementer locked; next **P33-S04-01**

## Todo updates

Status + notes on **P33-S04-00** only.

## Next

`P33-S04-01`
