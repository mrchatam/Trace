# P35-S01-00 — Scope planner (plan)

## Metadata
- id: P35-S01-00
- todo_ids: [P35-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, api-and-interface-design]
- verification: automated

## Objective

Lock planning scope so S01-01 can author `PLAN.md`: selection policy options ranked, Law 19 placement (library vs GUI helper), API/GUI touch list, and acceptance tests (feet-seller all-DONE + >100 tasks). **No product code.** Requires S00-02 PASS + `INVESTIGATION.md`.

## References

- [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Protocol + Law 19

## Locked defaults (refined 2026-08-21 against INVESTIGATION / S00-02 PASS)

Do **not** re-open root cause. Full locks live in thickened [01-plan.md](01-plan.md).

| Item | Value |
|------|-------|
| Artifact | `scopes/scope-01-plan/PLAN.md` only (authored in **P35-S01-01**) |
| Root cause | Overview ~L17–20 + Loop ~L51–54 → Step1 `33247e2d-…`; Loop112 `99d8fb92-…`; HTTP ignores limit; no library current-work |
| Policy goal | Default ≠ oldest DONE `tasks[0]` when later work exists |
| Placement set | **A** Go library + thin adapters · **B** one shared `web/src/lib` pick helper — lean **B** unless A fits one S02 row |
| Semantics lean | IN_PROGRESS → non-terminal → **last** DONE (`items[n-1]` on oldest-first); explicit `task_id` wins |
| Must cover | Overview + Loop (Graph/`prioritizeTaskSeeds` optional mention only) |
| Tests | all-DONE ≠ Step1; limit/>100 honesty; red assert seed exit 2 → green |
| Out | Weakening `plan_missing`; TRACE_TASK_ID-only; delete dogfood; Overview/Loop fork |

## Planner gate

- [x] `01-plan.md` runnable with PLAN template + locked defaults + exit criteria
- [x] `SCOPE-TODOS.md` accurate
- [x] Do **not** write `PLAN.md` in this planner row
- [x] Note effects on upcoming S02 planner seed

## Exit criteria

- [x] Plan implementer prompt locked
- [x] Next **P35-S01-01**

## Next

`P35-S01-01`
