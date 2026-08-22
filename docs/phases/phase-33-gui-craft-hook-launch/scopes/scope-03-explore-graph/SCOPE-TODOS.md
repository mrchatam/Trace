# Scope 03 — board map

**S03 Explore hook** — interactive project graph on `/` (D+B+C seed compose). Serial: **P33-S03-00 → P33-S03-01 → P33-S03-02**. Depends on S01 design + S02 CLI land (`trace gui` → `/`). Skills mandatory on 01+02.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 580 | P33-S03-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock data path, budgets, keyboard, files; thicken 01/02 |
| 581 | P33-S03-01 | [01-implement.md](01-implement.md) | Implementer | Overview-on-open in `Graph.tsx` (+ optional `overviewCompose` helper); e2e/smoke |
| 582 | P33-S03-02 | [02-review.md](02-review.md) | Reviewer | Theme B + Laws 6–7/19 + S02 `/` land; write `REVIEW.md` |

## Locked leans (from S00/S01/S02 + this planner)

| Topic | Lock |
|-------|------|
| Model | **(D)+(B)+(C)** — seeds → parallel budgeted `getGraph` → progressive expand |
| Route | Explore = **`/`** ≠ `/overview`; S02 opens `http://{addr}/` |
| Budgets | Seeds target **6** ≤**8**; per-seed **40**/depth **2**; **UI_CAP=100**; expand ≤**50** |
| API | `reuse` ops.ts only; no Leiden; no seed-export-as-graph-body; no fake project node |
| Keyboard | List+inspector required; canvas best-effort or residual documented |
| Craft | Hooks `data-kind`/`data-state`; **S04** owns shell colorize |

## Out of this scope

- Full shell colorize (S04), docs primary flip (S05), PATH/`trace gui` changes (S02 done).
