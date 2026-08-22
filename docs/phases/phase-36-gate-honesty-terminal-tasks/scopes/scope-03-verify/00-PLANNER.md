# P36-S03-00 — Scope planner (verify)

## Metadata
- id: P36-S03-00
- todo_ids: [P36-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, test-driven-development]
- verification: automated
- hooks: []

## Objective

Lock VERIFY blocks for feet-seller live dogfood **and** greenfield MCP agent workflow. Thicken `01-verify.md` + `02-dr-handoff.md`.

## References

- [PLAN.md](../scope-01-plan/PLAN.md) (when exists)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Fixture: `/home/ali/Desktop/feet seller telegram app`

## Session start

Follow agent-loop-protocol Session start. Waits for S02-02 PASS.

## Verify blocks (minimum — thicken into 01-verify)

0. Unit/integration tests from S02 pass (`go test ./...` scoped per PLAN)
1. **Greenfield agent path** — temp project: goal → MCP/bootstrap coarse plan → `trace loop gate --for edit` passes
2. CLI gate: Step 1 DONE — honest terminal outcome per PLAN (not misleading `done blocked`)
3. CLI gate: Loop 112 DONE — same
4. Live TaskDetail: no misleading red strip on DONE tasks (if GUI in PLAN)
5. Active/non-terminal synthetic task still gets PLAN gate if no plan
6. Feet-seller recovery — per PLAN (backfill run or documented legacy state)
7. Residuals + successor table for DR-HANDOFF

## Planner gate

- [ ] `01-verify.md` runnable with evidence paths
- [ ] `02-dr-handoff.md` successor table aligned to VERIFY outcome
- [ ] Do **not** run VERIFY in this row

## Next

`P36-S03-01`
