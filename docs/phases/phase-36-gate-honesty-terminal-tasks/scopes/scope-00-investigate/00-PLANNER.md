# P36-S00-00 — Scope planner (investigate)

## Metadata
- id: P36-S00-00
- todo_ids: [P36-S00-00]
- role: planner
- skills: [diagnosing-bugs, planning-and-task-breakdown, test-driven-development]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock investigation scope so a fresh subagent can produce `INVESTIGATION.md`: feet-seller repro, **Trace vs agent vs harness verdict**, two planning systems (plan-change vs progressive planner), MCP gap, install/enforce posture, transition path without gate, mega-goal pattern. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law **19**
- [Phase 36 README](../../README.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live: `internal/loop/gate.go`, `internal/loop/policy.go`, `internal/deliberation/select.go`, `cmd/trace/plan.go`, `internal/mcp/server.go`, `internal/install/`, `web/src/screens/TaskDetail.tsx`, `web/src/components/GateStrip.tsx`

## Session start

Follow agent-loop-protocol Session start. Fixture path is locked.

## Locked defaults

| Item | Value |
|------|-------|
| Output of S00-01 | `scopes/scope-00-investigate/INVESTIGATION.md` |
| Product / CLI / web edits | **No** on S00 |
| Fixture | `/home/ali/Desktop/feet seller telegram app` — read-only DB |
| Feedback loop | Prefer CLI `trace loop gate` JSON asserts; GUI TaskDetail screenshot secondary |
| Sequence | S00 → S01 serial |

## Must answer (handoff to 01)

1. **Verdict table:** Trace product vs agent misuse vs harness — evidence per row.
2. Prove plan-change count (11) vs progressive planner empty (0) — cite export/DB.
3. Trace transition path: how 123 tasks reached DONE without `PlanExists` ever true (enforce off? `as_operator` + reviews? gate not called on MCP path?).
4. MCP surface audit: what agents *can* vs *must* do for PLAN phase — `trace_add kind=plan-change` satisfies `PlanCritiqued` on apply but **not** `PlanExists` (`policy.go:60–62` vs `:45–49`); `plan.go` documents "No MCP plan tools".
5. Install audit: feet-seller missing AGENTS.md / `.trace/config.json` / hooks — impact on enforce + bootstrap discoverability.
6. Goal structure: 123 tasks / 1 goal — product guidance gap?
7. GUI secondary: terminal DONE + plan_missing copy (defer to S01 unless blocking).
8. Red-capable tests for S02: greenfield bootstrap vs legacy dogfood.

## Planner gate

- [ ] `01-investigate.md` runnable
- [ ] `02-review.md` checklist vs DESIGN-LOCKS + INTAKE
- [ ] `SCOPE-TODOS.md` accurate
- [ ] Do **not** write `INVESTIGATION.md` in this planner row

## Next

`P36-S00-01`
