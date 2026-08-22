# P26-S03-00 — Scope planner (P25-B loop recalibration)

## Metadata
- id: P26-S03-00
- todo_ids: [P26-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Verify `PLAN.md` S03 and lock implement/review prompts for INT-02, INT-05, INT-09. No product code on this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [PLAN.md](../scope-01-planning/PLAN.md)
- [AUDIT.md](../scope-00-loop-audit/AUDIT.md)
- [Phase 26 README](../../README.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Saturation | Consecutive empty-apply threshold ≥ 2 (named constant); first empty must not sticky-STOP greenfield |
| Reset | Clear stopped, hop_count, consecutive empty; state → EXECUTE via CLI |
| STOP UX | One canonical reason string for gate JSON + export |
| Schema | Version bump if columns added |
| Out of scope | S02 promotion; S04 installer ParentOrchestratorRule |
| Incoming S02 | S02 may add `ApplyResult.spawned_task_ids` and `spawned_tasks[].discovery_id`. Do not revert. Discoveries-only apply still saturates until S03-T01. |

## Problem (context)

1. **INT-02:** First empty apply saturates → sticky STOP / hop_budget on verify.
2. **INT-05:** No API to reset after gap pass.
3. **INT-09:** Gate vs export STOP strings disagree (`hop_budget_exceeded` vs `p19_saturated`).

## Deliverables

| D | Description |
|---|-------------|
| D1 | Saturation threshold ≥ 2 consecutive empty applies |
| D2 | `trace loop reset --task <id>` (or equivalent) restores EXECUTE |
| D3 | Unified stop reason constant across gate + export |
| D4 | Tests for first/second empty + reset |

## Planner gate

- [ ] `PLAN.md` S03 exists
- [ ] `01-implement.md` + `02-review.md` runnable
- [ ] `SCOPE-TODOS.md` current

## Exit criteria

- [ ] Implementer locked; own Notes updated; do not start implement until this row `done`
