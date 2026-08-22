# Scope 02 — board map

**S02 implement** — ship PLAN + tests + review. Serial: **P35-S02-00 → P35-S02-01 → P35-S02-02**. Requires `PLAN.md` from S01.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 616 | P35-S02-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock B + P1→P2→P3a + node:test cmd; thicken 01/02 |
| 617 | P35-S02-01 | [01-implement.md](01-implement.md) | Implementer | `pickActiveTask` + fetch-for-pick + Overview/Loop + unit tests |
| 618 | P35-S02-02 | [02-review.md](02-review.md) | Reviewer | Diff vs PLAN + Law 19; re-run node:test |

## Locked for implement (S02-00)

- Placement **B**: `web/src/lib/pickActiveTask.ts`
- Semantics **P1→P2→P3a**; `?task_id=` wins
- Fetch-for-pick via `listTasksForPick` (no HTTP pagination in-slice)
- Cmd: `cd web && node --experimental-strip-types --test src/lib/pickActiveTask.test.ts`

## Out of this scope

- Authoring PLAN (S01); live VERIFY dogfood write-up (S03) — though implementer may run local checks.
- Go library pick (A); OpenAPI/`handlers_tasks.go` limit/cursor compliance.
