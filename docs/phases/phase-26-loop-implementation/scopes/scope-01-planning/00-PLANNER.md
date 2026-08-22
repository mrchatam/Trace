# P26-S01-00 — Scope planner (task breakdown)

## Metadata
- id: P26-S01-00
- todo_ids: [P26-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Convert `AUDIT.md` into concrete, testable tasks for S02 (P25-A), S03 (P25-B), and S04 (installer). Gate **P26-S01-01**. No product code.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [AUDIT.md](../scope-00-loop-audit/AUDIT.md) (must exist before this row)
- [Phase 26 README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Input | `scopes/scope-00-loop-audit/AUDIT.md` |
| Output from S01-01 | `scopes/scope-01-planning/PLAN.md` |
| Product Go | **No** |
| Thresholds | List options; implementer + tests decide |
| Scope order | S02 → S03 → S04 (serial default; S04 may parallel S03 only after S01) |

## Planner gate

- [ ] `AUDIT.md` has INT-01/02/05/06/09 + installer sections
- [ ] `01-task-breakdown.md` runnable
- [ ] `SCOPE-TODOS.md` current

## Exit criteria

- [ ] Planning implementer prompt locked against AUDIT.md
- [ ] Own board row Notes updated; do not start S02

## Todo updates

Status + notes on **P26-S01-00** only.
