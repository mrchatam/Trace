# P37-S02-00 — Scope planner (implement)

## Metadata
- id: P37-S02-00
- todo_ids: [P37-S02-00]
- role: planner
- skills: [incremental-implementation, planning-and-task-breakdown]
- verification: automated

## Objective

Thicken implement + review prompts from `PLAN.md` accept set and S02 wave order. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [PLAN.md](../scope-01-plan/PLAN.md) (must exist)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)

## Session start

Follow agent-loop-protocol Session start. Block if PLAN.md missing.

## Must lock for S02-01

1. Implement wave order from PLAN.md (dependencies first).
2. File touch-list per accept row.
3. Test commands (`go test ./...` subsets).
4. Phase 36 regression subset to re-run.
5. MCP/OpenAPI update expectations for R2/R3/R5.

## Exit criteria

- [ ] `01-implement.md` + `02-review.md` runnable
- [ ] Board row → `done` with Notes

## Next

`P37-S02-01`
