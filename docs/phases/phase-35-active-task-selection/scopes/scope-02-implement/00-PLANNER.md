# P35-S02-00 — Scope planner (implement)

## Metadata
- id: P35-S02-00
- todo_ids: [P35-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, test-driven-development]
- verification: automated

## Objective

Lock implement defaults from `PLAN.md` (files, test commands, Law 19 placement). Thicken `01-implement` / `02-review`. **No product code in this row.**

## References

- [PLAN.md](../scope-01-plan/PLAN.md)
- [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)

## Locked defaults (final — from PLAN.md + live repo 2026-08-21)

| Item | Value |
|------|-------|
| Spec | [`PLAN.md`](../scope-01-plan/PLAN.md) normative |
| Placement | **B** — `web/src/lib/pickActiveTask.ts` (+ `pickActiveTask.test.ts`) |
| Semantics | **P1→P2→P3a** (IN_PROGRESS → non-terminal → last list item); empty → `null` |
| Must wire | `Overview.tsx` + `Loop.tsx`; remove local/`items[0]` auto-pick fallbacks |
| Overrides | `?task_id=` wins; `TRACE_TASK_ID` agent bind (not sole GUI fix) |
| Fetch-for-pick | Prefer `listTasksForPick` in `ops.ts` (page on `next_cursor`); **no** HTTP/OpenAPI pagination in S02 |
| Deferred | Placement **A** (Go library) as Law-19 future upgrade |
| Test runner | **node:test** + strip-types (same as `overviewCompose.test.ts`); **not** vitest/`npm test` |
| Test cmd | `cd web && node --experimental-strip-types --test src/lib/pickActiveTask.test.ts` |
| Fixtures | Synthetic only; narrative Step1 / Loop112 UUIDs; never mutate feet-seller |
| Acceptance | all-DONE ≠ Step1; >100 honesty; red assert seed → green |

## Planner gate

- [x] `01-implement.md` runnable (minimal todos, exit criteria, test cmds)
- [x] `02-review.md` checklist vs PLAN
- [x] Do not implement in this row

## Next

`P35-S02-01`
