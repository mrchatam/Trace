# Scope 04 — board map

**S04 colorize/craft** — shell-wide tokens + craft (Theme A). Serial: **P33-S04-00 → P33-S04-01 → P33-S04-02**. After S03 (**P33-S03-02 PASS**). Skills mandatory on implement+review.

**Live gap (planner 2026-08-21):** `tokens.css` has forest-moss surfaces/accent/focus/semantic but **no** `--kind-*` / `--state-*`. Graph nodes use muted kind labels + generic `--border` — chroma strip not landed. S03 hooks (`data-kind`/`data-state`, `.graph-node__kind`/`__state`, seed/center/selected) are ready — **CSS colorize only**; do not rework seed compose or `/` routing.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 583 | P33-S04-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock craft map + thicken 01/02 (**this row**) |
| 584 | P33-S04-01 | [01-implement.md](01-implement.md) | Implementer | Land DESIGN tokens; chroma strip; contrast evidence |
| 585 | P33-S04-02 | [02-review.md](02-review.md) | Reviewer | Theme A + a11y; artifact `REVIEW.md` |

## Planner locks (handoff)

| Question | Answer |
|----------|--------|
| Token map | DESIGN → `tokens.css` kind/state; `app.css` `[data-kind]` strip+fill, `[data-state]` chips, pills/empty/nav |
| Priority | Explore canvas+list+inspector → chrome → shared empty/error → other screens inherit |
| Verify | Light+dark screenshots + contrast pairs in Notes (≥4.5:1 text, ≥3:1 stroke/focus) |
| Reduced motion | Preserve existing `@media (prefers-reduced-motion)` on nodes/inspector |
| Color-not-only | Keep kind + state **text**; color is redundant channel |

## Review status

**P33-S04-02 PASS** (2026-08-21) — artifact [`REVIEW.md`](REVIEW.md). No spawn. Residuals for S05: EmptyState CTA polish; optional canvas screenshot (evidence list-focused).

## Out of this scope

- Overview compose / budgets / canvas arrow-roving (S03 residual)
- `trace gui` / PATH (S02)
- Docs primary flip (S05), VERIFY (S06)
