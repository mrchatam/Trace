# P36-S02-00 — Scope planner (implement)

## Metadata
- id: P36-S02-00
- todo_ids: [P36-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, incremental-implementation]
- verification: automated
- hooks: []

## Objective

Lock implement scope from `PLAN.md`. Thicken `01-implement.md` + `02-review.md` with touch-list, test files, and Law 19 boundaries. **No product code in this row.**

## References

- [PLAN.md](../scope-01-plan/PLAN.md) (when exists)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live touch areas (typical — PLAN.md is SoT):
  - `internal/loop/policy.go`, `internal/loop/gate.go`
  - `internal/mcp/` (new plan tools if chosen)
  - `cmd/trace/plan.go`
  - `internal/install/`, bundled AGENTS templates
  - `internal/httpapi/` (gate endpoint if needed)
  - `web/src/components/GateStrip.tsx`, `web/src/screens/TaskDetail.tsx`

## Session start

Follow agent-loop-protocol Session start. Waits for S01-01 PLAN.md.

## Must lock (for 01-implement)

1. Exact files from PLAN.md touch-list
2. Test strategy — unit tests in `internal/loop/`; MCP parity tests in `internal/mcp/`; CLI tests in `cmd/trace/`
3. Feet-seller recovery steps (if in PLAN) — read-only vs explicit backfill command
4. Greenfield fixture — temp dir with init + goal + MCP bootstrap path
5. Non-goals — no unrelated refactors; no global gate weaken

## Planner gate

- [ ] `01-implement.md` runnable alone
- [ ] `02-review.md` checklist matches PLAN.md acceptance tests
- [ ] Do **not** implement product code in this row

## Next

`P36-S02-01`
