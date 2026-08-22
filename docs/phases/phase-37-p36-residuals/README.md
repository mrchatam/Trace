# Phase 37 — Phase 36 residuals closure

Human-promoted **2026-08-22**. **Complete — closed at `P37-S03-02`.** Close documented deferrals and low findings from Phase 36.

## Design SoT

| Doc | Role |
|-----|------|
| [`INTAKE.md`](INTAKE.md) | Residual inventory R1–R11 |
| [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md) | Triage rules + accept sketch |
| [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) | Row `P37-00` |
| [`DR-HANDOFF.md`](DR-HANDOFF.md) | **CLOSED** at `P37-S03-02` |

Board: [`docs/TODO/phase-37.md`](../../TODO/phase-37.md).

## Predecessor evidence

- Phase 36 DR-HANDOFF: [`../phase-36-gate-honesty-terminal-tasks/DR-HANDOFF.md`](../phase-36-gate-honesty-terminal-tasks/DR-HANDOFF.md)
- VERIFY: [`../phase-36-gate-honesty-terminal-tasks/scopes/scope-03-verify/VERIFY-NOTES.md`](../phase-36-gate-honesty-terminal-tasks/scopes/scope-03-verify/VERIFY-NOTES.md)
- PLAN deferrals: [`../phase-36-gate-honesty-terminal-tasks/scopes/scope-01-plan/PLAN.md`](../phase-36-gate-honesty-terminal-tasks/scopes/scope-01-plan/PLAN.md) §2.4, §2.6, touch-list defer

## Scope sequence

```
S00 triage residuals → RESIDUALS.md (accept/defer/reject per R1–R11)
 → S01 plan → PLAN.md
 → S02 implement + review
 → S03 VERIFY + DR-HANDOFF
```

| Scope | Rows | Artifact |
|-------|------|----------|
| S00 | P37-S00-00 → 02 | `RESIDUALS.md` |
| S01 | P37-S01-00 → 01 | `PLAN.md` |
| S02 | P37-S02-00 → 02 | Code + tests |
| S03 | P37-S03-00 → 02 | `VERIFY-NOTES.md` + CLOSED |

## Residual code anchors (P37-00 spot-check — S00 must triage live)

| ID | Signal | Location | Post-P36 state |
|----|--------|----------|----------------|
| R1 | PlanExists bridge | P36 PLAN §2.4 defer | No heuristic in `internal/loop/policy.go` — plan-changes alone do not satisfy `PlanExists` |
| R2 | HTTP plan routes | `internal/httpapi/server.go` | `GET /v1/plans` only — no POST bootstrap/create parity |
| R3 | MCP loop gate | `internal/mcp/tools_loop.go:39–47` | `trace_loop` = `next\|apply\|status` — no `gate` action |
| R4 | Bootstrap help | `cmd/trace/plan.go`, `internal/mcp/tools_plan.go` | Bootstrap ships; human-refinement note (create-coarse/deep) may be thin |
| R5 | `advisories[]` | `internal/loop/apply.go` `StatusResult` | `GoalStructureWarning` in `internal/planner/advisory.go` — **not** wired to loop status JSON |
| R6 | Config warn test | `internal/config/enforce.go:43–44` | `WarnIfTraceDirWithoutConfig` — **no** dedicated unit test |
| R7 | Enforce default warn | P36 PLAN §2.6 | Document-only nudge accepted; default enforce stays **off** without config |
| R8 | Plan UX surfaces | `web/src/screens/TaskDetail.tsx`, `Overview.tsx` | TaskDetail bootstrap hint; Overview GateStrip only — no dedicated plan screen |
| R9 | Feet-seller plan quality | Dogfood fixture | Post-bootstrap minimal plan; human refinement via create-coarse/deep expected |
| R10 | Live GUI verify | P36 VERIFY Block 4 | Deferred — terminal honesty proved via pre-bootstrap API + GateStrip |
| R11 | Critique path | P36 VERIFY Block 1 partial | CLI plan chain → `plan_uncritiqued`; Block 0 MCP covers full greenfield path |

## Dogfood fixture

`/home/ali/Desktop/feet seller telegram app` — read-only unless R9 recovery explicitly scoped in S01.

## Hard constraints

- Law 19 — library first; adapters thin
- Do not weaken Phase 36 active-work PLAN gate
- No silent PlanExists bridge
