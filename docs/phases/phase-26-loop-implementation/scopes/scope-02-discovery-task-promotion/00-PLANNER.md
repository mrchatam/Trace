# P26-S02-00 — Scope planner (P25-A discovery→task)

## Metadata
- id: P26-S02-00
- todo_ids: [P26-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Verify `PLAN.md` S02 section and lock implement/review prompts for INT-01 + INT-06. No product code on this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [PLAN.md](../scope-01-planning/PLAN.md)
- [AUDIT.md](../scope-00-loop-audit/AUDIT.md)
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) §4 (human gate)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Deliverables | SpawnedTasks on apply; CLI visibility; MCP description nudge; tests for BLOCKING→task |
| Human gate | No silent background spawn — `loop apply` or explicit `trace add task` only |
| Out of scope | S03 saturation/reset; S04 ParentOrchestratorRule wiring |
| Tests | `go test ./internal/...` must stay green |

## Problem (context)

G1 recorded discoveries but **0 new tasks** (FM-10). `loop apply` lacks `spawned_tasks[]`; MCP does not nudge promotion after discovery.

## Deliverables (INT-01 + INT-06)

| D | Description |
|---|-------------|
| D1 | `SpawnedTasks []string` (or equiv.) on apply result for promoted task IDs |
| D2 | CLI `trace loop apply` surfaces spawned tasks |
| D3 | MCP `trace_add` description: promotion guidance before generic add |
| D4 | Install/harness nudge after BLOCKING discovery |
| D5 | Unit/integration test: BLOCKING discovery → linked task |

## Planner gate

- [ ] `PLAN.md` S02 section exists
- [ ] `01-implement.md` + `02-review.md` runnable
- [ ] `SCOPE-TODOS.md` current

## Exit criteria

- [ ] Implementer has locked paths/acceptance from PLAN.md
- [ ] Own row Notes only; do not start P26-S02-01 until this row is `done`
