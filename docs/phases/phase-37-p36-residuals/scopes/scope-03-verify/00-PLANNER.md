# P37-S03-00 — Scope planner (verify)

## Metadata
- id: P37-S03-00
- todo_ids: [P37-S03-00]
- role: planner
- skills: [test-driven-development, planning-and-task-breakdown]
- verification: automated

## Objective

Lock VERIFY blocks from PLAN.md § acceptance + Phase 36 regression subset. Thicken `01-verify.md` + `02-dr-handoff.md`. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [PLAN.md](../scope-01-plan/PLAN.md) § VERIFY blocks
- Phase 36 [`VERIFY-NOTES.md`](../../../phase-36-gate-honesty-terminal-tasks/scopes/scope-03-verify/VERIFY-NOTES.md)

## Session start

Follow agent-loop-protocol Session start.

## Minimum VERIFY blocks (lock in 01-verify)

| Block | Content |
|-------|---------|
| 0 | Phase 36 acceptance subset still green |
| 1 | Per accepted residual — evidence file or test name |
| 2 | Feet-seller spot-check (if R8/R9/R10 in accept set) |
| 3 | Greenfield MCP path (if R3/R11 in accept set) |
| 4 | Re-defer registry updated in VERIFY-NOTES |
| 5 | Successor table for DR-HANDOFF |

## Locked defaults

| Item | Value |
|------|-------|
| Fixture | `/home/ali/Desktop/feet seller telegram app` — read-only unless PLAN scoped mutation |
| Evidence dir | `experiments/runs/` + optional `docs/verification/phase-37-p36-residuals/` |
| R10 | Browser spot-check only if accepted — else cite Block 4 defer rationale |

## Exit criteria

- [ ] `01-verify.md` + `02-dr-handoff.md` runnable
- [ ] Board row → `done` with Notes

## Next

`P37-S03-01`
