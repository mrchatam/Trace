# P35-S00-02 — Investigation review

## Metadata
- id: P35-S00-02
- todo_ids: [P35-S00-02]
- role: reviewer
- skills: [diagnosing-bugs, code-review-and-quality]
- verification: automated

## Objective

Independent review of `INVESTIGATION.md`. Confirm DESIGN-LOCKS / INTAKE alignment, cite quality, and that S01 can plan without re-opening root cause. Spawn forward only if investigation is incomplete.

## References

- [INVESTIGATION.md](INVESTIGATION.md) (must exist from S00-01)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [01-investigate.md](01-investigate.md) — must-answer set + template
- Live paths cited in INVESTIGATION (re-open files; do not trust cites blindly)

## Checklist vs DESIGN-LOCKS

- [ ] Verdict matches dogfood symptom (Step 1 vs ~123 / Loop 112), not a different bug
- [ ] Must-fix theme covered: default ≠ first `listTasks` row when that is stale DONE and later tasks exist
- [ ] Must-test angle present: all-DONE + later tasks; `limit`/>100 honesty addressed with evidence
- [ ] Selection home guidance usable by S01 (library found **or** explicit “none → shared GUI helper”)
- [ ] Agent/`TRACE_TASK_ID` story noted; out-of-scope items not proposed as primary fix (`plan_missing` weaken, SaaS, delete dogfood)

## Checklist vs INTAKE

- [ ] Each INTAKE likely-cause bullet (1–5) explicitly **confirmed** or **refuted** with evidence
- [ ] Overview `pickActiveTask` + Loop `items[0]` cited with file paths (and lines if claimed)
- [ ] `limit`/pagination honesty verified end-to-end (HTTP length + handler), not assumed from client alone
- [ ] “No durable current work” confirmed or corrected with search notes
- [ ] Secondary `plan_missing` not treated as must-fix

## Quality / handoff

- [ ] Library “current work” search explicit (found or none) — Law 19 handoff clear for S01
- [ ] Red-capable loop sketch is agent-runnable and asserts the **user** symptom (default = Step 1)
- [ ] Rejects include “TRACE_TASK_ID-only” and “delete dogfood data” (and other SCOPE-TODOS rejects)
- [ ] No product code slipped into S00 (`git status` / diff sanity)
- [ ] Handoff to S01 is enough to draft policy options without re-investigating

## On FAIL

Insert spawn rows forward (`P35-S00-03a` / `03b`) per protocol; do not rewrite S00-01 history. Thicken spawn prompts with the specific missing evidence (e.g. unrun HTTP `limit` check).

## Exit criteria

- [ ] PASS or FAIL with Notes evidence
- [ ] On PASS, next runnable is **P35-S01-00**

## Next

`P35-S01-00`
