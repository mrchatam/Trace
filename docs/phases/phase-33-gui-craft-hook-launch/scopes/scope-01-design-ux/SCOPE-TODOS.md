# Scope 01 — board map

**S01 design + UX** — tokens + Explore IA only (markdown artifacts). Serial: **P33-S01-00 → P33-S01-01 → P33-S01-02**. Depends on S00-02 PASS. Skills mandatory on 01/02. No product CSS/TSX shipping.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 574 | P33-S01-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock artifacts + thicken 01/02 (**done** after planner run) |
| 575 | P33-S01-01 | [01-design-ux.md](01-design-ux.md) | Implementer | Author [`DESIGN.md`](DESIGN.md) + [`UX-IA.md`](UX-IA.md) |
| 576 | P33-S01-02 | [02-review.md](02-review.md) | Reviewer | Craft + a11y + Laws 6–7 + ownership → [`REVIEW.md`](REVIEW.md) **PASS** (high); thickened S02–S04 |

## Planner locks (P33-S01-00 — binding)

| Lock | Value |
|------|-------|
| Artifacts | **Two files**: `DESIGN.md` (tokens/craft) + `UX-IA.md` (Explore IA) |
| Mode / authority | Operate; **refine** forest-moss / IBM Plex `tokens.css` — not purple/cream/broadsheet, not HUD/glow, not wholesale slate+#22C55E |
| IA model | S00 **(D)+(B)+(C)** — seed compose → budgeted `getGraph` → progressive expand |
| Budgets | Seeds target **6** (4–8, ≤8); per-seed **40** (≤50); merged honor **`UI_CAP=100`** (argue ≤120 only); depth **2**; user-driven expand |
| Seed priority | `getProject` → `listTasks` (IN_PROGRESS then active) → `search` fill → dedupe ≤8 |
| Routes | Explore = `/` Graph ≠ Nav `/overview` |
| API | Prefer `reuse`; never full dump / seed-export-as-graph-body |
| Ownership | **S03** implements IA; **S04** lands tokens + shell colorize; S01 markdown only |

## S00 research leans (binding)

- Overview model: **(D)+(B)+(C)** from [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md).
- Clusters (A): visual grouping inspiration only — no Leiden API in Phase 33.
- Color: tokens here; S04 owns full shell.

## Out of this scope

- `trace gui` / PATH (S02), Explore UI ship (S03), shell colorize (S04), docs primary flip (S05).
