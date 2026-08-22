# Phase 35 — Active task selection (test + fix)

**Phase planner.** Row `P35-00`.

## Metadata
- id: P35-00
- todo_ids: [P35-00]
- role: planner
- skills: [diagnosing-bugs, test-driven-development, planning-and-task-breakdown]
- verification: automated

## Mission

**Test** then **fix** wrong “current task” binding: GUI/gate shows task 1 while the project is at task ~123 (feet-seller dogfood).

Read [`INTAKE.md`](INTAKE.md) + [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md).

## Gate

Phase 34 is **done**. Proceed. Do not rewrite P34 history.

## Scope sequence (matches board)

```
S00 investigate (repro + root cause: pickActiveTask, Loop default, list limit honesty)
 → S01 plan fix (selection policy + API/GUI)
 → S02 implement + tests + review
 → S03 VERIFY (feet-seller live) + DR-HANDOFF
```

| Scope | Board rows | Artifact |
|-------|------------|----------|
| S00 | P35-S00-00 → S00-01 → S00-02 | `INVESTIGATION.md` |
| S01 | P35-S01-00 → S01-01 | `PLAN.md` |
| S02 | P35-S02-00 → S02-01 → S02-02 | Code + tests |
| S03 | P35-S03-00 → S03-01 → S03-02 | `VERIFY-NOTES.md` + CLOSED handoff |

## Hard constraints

- Law 19 — no business logic fork in `web/` beyond selection UX calling library/API
- Do not delete feet-seller data
- Prefer library-backed “current work” if one exists; else define explicit GUI policy + document agent `TRACE_TASK_ID` expectation
- Do not weaken `plan_missing` for true PLAN-phase work (secondary honesty only)

## Planner gate (P35-00)

- [x] Phase folder README + light locks + live baseline
- [x] Each scope has runnable `00-PLANNER` / implement / review (or VERIFY) stubs
- [x] `SCOPE-TODOS.md` per scope
- [x] Board scope-sequence prose matches folders S00–S03
- [x] `DR-HANDOFF.md` OPEN; close owner `P35-S03-02`
- [x] No product / GUI / serve implementation in this row

## P35-00 outcome (2026-08-21)

Gate **PASS**. Thickened README (locks, live baseline, serial scopes). Light DESIGN-LOCKS (`limit` honesty note; library-first). Protocol-thickened S00–S03 prompts + `SCOPE-TODOS.md`. Aligned DR-HANDOFF checklist to board (S00–S03). Spot-checked feet-seller: 123 tasks all DONE; index0=Step1; index122=Loop112. HTTP `handleListTasks` currently ignores client `limit` — S00 must verify end-to-end. No product code. Next: **P35-S00-00**.

## Next

`P35-S00-00`
