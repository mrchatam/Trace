# P37-S01-00 — Scope planner (plan)

## Metadata
- id: P37-S01-00
- todo_ids: [P37-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, spec-driven-development]
- verification: automated

## Objective

Thicken `01-plan.md` from `RESIDUALS.md` accept set. Lock touch-list order (library → MCP → HTTP → install → GUI). **No PLAN.md in this row.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [RESIDUALS.md](../scope-00-triage/RESIDUALS.md) (must exist — S00-02 PASS)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- Phase 36 [`PLAN.md`](../../../phase-36-gate-honesty-terminal-tasks/scopes/scope-01-plan/PLAN.md) touch-list pattern

## Session start

Follow agent-loop-protocol Session start. Block if RESIDUALS.md missing or S00-02 not done.

## Locked defaults

| Item | Value |
|------|-------|
| Output of S01-01 | `scopes/scope-01-plan/PLAN.md` |
| Product code | **No** on S01-00 |
| R1 | If accepted: advisory-only — never `PlanExists=true` |
| R7 | Lock explicit enforce default in PLAN § (off unless product decision) |
| Regression | PLAN must include Phase 36 acceptance subset for S03 |

## Must lock for 01-plan

1. Accepted residuals table (from RESIDUALS accept rows only).
2. Re-deferred residuals with owner/trigger.
3. Touch-list with file targets per accept.
4. Acceptance tests per accept + VERIFY block mapping.
5. Non-goals (hosted SaaS, silent bridge, history rewrite).

## Exit criteria

- [ ] `01-plan.md` runnable alone
- [ ] Board row → `done` with Notes

## Next

`P37-S01-01`
