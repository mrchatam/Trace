# P28-S01-00 — Scope planner (integration test matrix)

## Metadata
- id: P28-S01-00
- todo_ids: [P28-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, test-driven-development]
- mcps: [user-codegraph]
- verification: automated
- hooks: []

## Objective

Plan automated regression for P25-A/B/C/D/E implementations — not dogfood (S02). Lock test homes, minimum cases, and `01-implement.md` / `02-review.md` from `RESIDUAL-AUDIT.md` R7 seeds.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) (after S00)
- [Phase 28 README](../../README.md)
- Existing tests: `internal/loop/*_test.go`, `cmd/trace/enforce_test.go`, `internal/install/enforcement_test.go`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Area | Test home | Minimum cases |
|------|-----------|---------------|
| Promotion (P25-A) | `internal/loop/` or `evals/` | BLOCKING discovery → spawned task |
| Reset (P25-B) | `internal/loop/` | saturate → reset → gate allow |
| Saturation (P25-B) | `internal/loop/` | 1 empty no STOP; 2 empty STOP |
| Honesty (P25-E) | `cmd/trace/` | thin export enforce blocked; rich passes |
| Install (P25-C) | `internal/install/` | gap pass + orchestrator phrases |
| Protocol (P25-D) | `experiments/.../score.sh` smoke or `evals/` | P25-3a/3b labels (optional script test) |

Prefer extending existing test files over new packages. `evals/p28-regression/` acceptable for CLI golden driver.

## Planner gate

- [ ] `RESIDUAL-AUDIT.md` S01 seeds reviewed
- [ ] `01-implement.md` + `02-review.md` runnable
- [ ] `SCOPE-TODOS.md` lists S01 board rows

## Exit criteria

- [ ] Implementer prompt locked for fresh subagent
- [ ] Board row P28-S01-00 Notes cite what was verified/thickened
- [ ] Next runnable **P28-S01-01**

## Todo updates

Status + notes on **P28-S01-00** only.

## Next

`P28-S01-01`
