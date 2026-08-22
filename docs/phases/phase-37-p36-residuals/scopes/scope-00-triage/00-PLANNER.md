# P37-S00-00 — Scope planner (triage)

## Metadata
- id: P37-S00-00
- todo_ids: [P37-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, code-review-and-quality]
- mcps: [user-trace]
- verification: automated

## Objective

Lock S00 triage scope for R1–R11. Thicken `01-triage.md` + `02-review.md` + `SCOPE-TODOS.md`. **No product code. No RESIDUALS.md in this row.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law **19**
- [Phase 37 README](../../README.md)
- [INTAKE.md](../../INTAKE.md), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Phase 36: [DR-HANDOFF.md](../../../phase-36-gate-honesty-terminal-tasks/DR-HANDOFF.md), [VERIFY-NOTES.md](../../../phase-36-gate-honesty-terminal-tasks/scopes/scope-03-verify/VERIFY-NOTES.md), [PLAN.md](../../../phase-36-gate-honesty-terminal-tasks/scopes/scope-01-plan/PLAN.md) §2.4, §2.6
- Live: `internal/planner/advisory.go`, `internal/loop/apply.go`, `internal/mcp/tools_loop.go`, `internal/httpapi/server.go`, `internal/config/enforce.go`, `cmd/trace/plan.go`, `web/src/screens/TaskDetail.tsx`, `web/src/screens/Overview.tsx`

## Session start

Follow agent-loop-protocol Session start. Re-read live code before locking triage questions — P36 may have partially shipped low findings.

## Locked defaults

| Item | Value |
|------|-------|
| Output of S00-01 | `scopes/scope-00-triage/RESIDUALS.md` |
| Product / CLI / web edits | **No** on S00 |
| Fixture | `/home/ali/Desktop/feet seller telegram app` — read-only |
| R1 stance | Advisory-only bridge OK; **reject** silent `PlanExists=true` |
| R7 stance | Default enforce stays `off` unless S01 locks explicit product change |
| Sequence | S00 → S01 serial |

## Must answer (handoff to 01-triage)

1. Per R1–R11: **accept / defer / reject** with effort (S/M/L) and risk.
2. Re-read live code — which residuals already partially shipped in P36 S02?
3. Dependency order for S02 (e.g. R5 `advisories[]` before R8 GUI surfaces).
4. Explicit re-deferrals need owner + trigger (not silent orphan).
5. R1: if accept advisory bridge, document exact signal (plan-change count threshold?) without setting `PlanExists`.
6. R2/R3: Law 19 — thin adapter over library; cite handler/tool targets.
7. R11: MCP critique seed vs agent-workflow doc — triage which path.

## Planner gate

- [ ] `01-triage.md` runnable alone
- [ ] `02-review.md` checklist vs DESIGN-LOCKS + INTAKE
- [ ] `SCOPE-TODOS.md` board IDs accurate
- [ ] Do **not** write `RESIDUALS.md` in this planner row

## Exit criteria

- [ ] S00-01 and S00-02 prompts thickened with locked defaults
- [ ] Live code anchors in SCOPE-TODOS verified against repo
- [ ] Board row `P37-S00-00` → `done` with Notes

## Minimal todos

- [ ] Re-read INTAKE R1–R11 + P36 DR-HANDOFF residuals paragraph
- [ ] Spot-check live code for partial P36 shipment (see SCOPE-TODOS)
- [ ] Thicken `01-triage.md` + `02-review.md`
- [ ] Update board status + notes only

## Next

`P37-S00-01`
