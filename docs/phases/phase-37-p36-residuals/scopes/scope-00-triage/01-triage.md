# P37-S00-01 — Triage residuals

## Metadata
- id: P37-S00-01
- todo_ids: [P37-S00-01]
- role: implementer
- skills: [planning-and-task-breakdown, code-review-and-quality, source-driven-development]
- mcps: [user-trace, user-codegraph]
- verification: automated

## Objective

Author **only** `scopes/scope-00-triage/RESIDUALS.md`: triage table for R1–R11 with accept/defer/reject, effort (S/M/L), risk, live code cites, S02 wave order, and explicit re-defers. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law **19**
- [INTAKE.md](../../INTAKE.md), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- Phase 36: [DR-HANDOFF.md](../../../phase-36-gate-honesty-terminal-tasks/DR-HANDOFF.md), [VERIFY-NOTES.md](../../../phase-36-gate-honesty-terminal-tasks/scopes/scope-03-verify/VERIFY-NOTES.md) § Residuals, [PLAN.md](../../../phase-36-gate-honesty-terminal-tasks/scopes/scope-01-plan/PLAN.md) §2.4, §2.6, §2.7, defer touch-list
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — planner-locked live facts + **preliminary triage** (S00-00)
- Live anchors: `internal/planner/advisory.go`, `internal/loop/apply.go`, `internal/loop/policy.go`, `internal/mcp/tools_loop.go`, `internal/mcp/tools_plan.go`, `internal/httpapi/server.go`, `internal/httpapi/handlers_p1.go`, `internal/httpapi/handlers_loop.go`, `internal/config/enforce.go`, `cmd/trace/plan.go`, `cmd/trace/loop.go`, `web/src/screens/TaskDetail.tsx`, `web/src/screens/Overview.tsx`, `api/openapi.yaml`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). **Grep/read live repo** before deciding accept/defer — do not rely on INTAKE defaults or SCOPE-TODOS alone. Confirm partial P36 S02 shipment (see Locked defaults § Partial shipment).

## Locked defaults

| Item | Value |
|------|-------|
| Artifact | `RESIDUALS.md` in this scope folder |
| Fixture | `/home/ali/Desktop/feet seller telegram app` — read-only |
| R1 | **Accept** advisory-only; **reject** silent `PlanExists=true`. Signal: `advisories[]` entry when goal-linked `plan_changes ≥ 1` and `!PlanExists` (see § R1 signal below). Does **not** mutate store or `policy.go` plan flags. |
| R2 | **Accept** thin HTTP POST plan routes (Law 19). Targets: new handlers in `internal/httpapi/handlers_p1.go` (or sibling) calling `planner.Service` — mirror `cmd/trace/plan.go` subcommands. OpenAPI: `api/openapi.yaml`. |
| R3 | **Accept**. `trace_loop action=gate` mirrors `cmd/trace/loop.go` `cmdLoopGate` → `loop.EvaluateGate`; HTTP precedent: `GET /v1/loop/gate` in `handlers_loop.go:104`. |
| R4 | **Accept** (S). Help-text only: `cmd/trace/plan.go` bootstrap usage + top-level help — add human-refinement note per P36 PLAN §2.2 risk paragraph. |
| R5 | **Accept** (S). Wire existing `planner.GoalStructureWarning` into `loop.Status` → `StatusResult.advisories[]`. Helper + tests already in P36 S02; gap is status JSON only (`apply.go:196–209` has no `advisories`). |
| R6 | **Accept** (S). Unit test `WarnIfTraceDirWithoutConfig` in `internal/config/enforce_test.go`. Helper shipped P36 (`enforce.go:43–57`, called from `cmd/trace/init.go:50`). |
| R7 | **Defer** (re-defer). Default enforce stays **`off`** when config missing (`LoadEnforceMode` → `EnforceOff`, `enforce.go:25–29`). P36 §2.6 shipped stderr nudge only. Owner + trigger required in RESIDUALS §6. |
| R8 | **Accept** minimal (M). At least **Overview** plan-gap surface beyond TaskDetail bootstrap paragraph. Law 19: consume API (`GET /v1/loop/gate` or status advisories); no planner logic in `web/`. Full plan screen → defer with owner. |
| R9 | **Defer** to S03 verify (dogfood). Document refinement path (`create-coarse` / `deep` after bootstrap); no mandatory feet-seller history rewrite. Optional S02 doc pointer only. |
| R10 | **Accept** in S03 (S). Live browser spot-check terminal + plan surfaces; pin evidence under `docs/verification/` or `experiments/runs/`. |
| R11 | **Accept** doc path (S); **reject** new MCP critique-seed tool. Block 0 MCP greenfield test already seeds critique via `loop apply`; gap is locked CLI verify script + agent-workflow discoverability. |

### Partial P36 S02 shipment (01 must verify; do not re-implement)

| ID | Shipped in P36 S02 | Still open |
|----|-------------------|------------|
| R4 | `trace plan bootstrap` + MCP `trace_plan bootstrap` | Help omits explicit human-refinement note (VERIFY-NOTES low #1) |
| R5 | `GoalStructureWarning` + unit test; wired to `plan show` stderr + MCP `trace_plan show` `goal_structure_warning` | `trace loop status` has no `advisories[]` (`StatusResult` in `apply.go:196–209`) |
| R6 | `WarnIfTraceDirWithoutConfig` + `init` call | No dedicated unit test |
| R7 | Stderr nudge helper (§2.6 document-only accept) | Default enforce still `off`; no auto-flip to `warn` |
| R8 | `TaskDetail.tsx` terminal bootstrap advisory paragraph | `Overview.tsx` GateStrip only — no plan-gap copy |
| R1, R2, R3 | — | Full gap |
| R9–R11 | MCP greenfield path (R11 Block 0) | Feet refinement quality (R9); live GUI (R10); CLI verify critique gap (R11) |

### R1 advisory signal (locked — no PlanExists)

When **accept**ing R1, RESIDUALS §3 must specify:

| Field | Lock |
|-------|------|
| Trigger | `!PlanExists(goal)` **and** `len(goalLinkedPlanChanges) ≥ 1` (reuse counting logic in `advisory.go` `goalLinkedPlanChangeIDs` or equivalent store read) |
| Threshold N | **1** (feet-seller had 11 plan-changes with 0 planner rows — any linked plan-change without progressive plan warrants nudge) |
| Surface | `trace.loop.status.v1` → `advisories[]` entry |
| Advisory code | `bootstrap_recommended` (or equivalent stable snake_case) |
| Message shape | Recommend `trace plan bootstrap --goal <id>` or MCP `trace_plan action=bootstrap` |
| Must not | Set `PlanExists`, write planner rows, or change `policy.go` deliberation inputs |

R5 (`goal_structure_warning`) uses **task count > 15** — orthogonal to R1; both may appear in `advisories[]`.

### S02 dependency order (locked recommendation for RESIDUALS §5)

```text
Wave A — status schema (blocks R1, R8 data):
  R5 advisories[] field + GoalStructureWarning wire
  R1 bootstrap_recommended advisory (same channel)

Wave B — adapters (parallel after A):
  R3 MCP trace_loop gate
  R2 HTTP POST plan routes (bootstrap minimum; create-coarse/deep if S01 expands)
  R4 bootstrap help text

Wave C — tests + docs:
  R6 WarnIfTraceDirWithoutConfig test
  R11 agent-workflow / verify-script doc for critique via loop apply

Wave D — GUI (after A):
  R8 Overview plan-gap / advisories surface

Out of S02:
  R7 defer (product default flip)
  R9, R10 → S03 verify
```

### Re-defer registry (locked seeds for RESIDUALS §6)

| ID | Owner | Trigger |
|----|-------|---------|
| R7 enforce default `warn` | Human / product | Explicit S01+ product decision to change `LoadEnforceMode` missing-config behavior; until then default **off** |
| R8 full plan screen / plan tree GUI | Phase 38 planner (or human) | S03 VERIFY shows Overview minimal surface insufficient for operator needs |
| R9 feet-seller deep refinement | Human dogfood | Post-bootstrap `create-coarse`/`deep` path documented; live quality acceptable for verify |

## RESIDUALS.md template

1. **Summary** — counts accept / defer / reject; note partial P36 shipment
2. **Table** — ID, item, P36 source, decision, effort (S/M/L), risk (low/med/high), rationale, S02 scope (Y/N), test sketch, live cite
3. **§3 R1 PlanExists advisory bridge** — locked signal above; explicit reject of silent satisfy
4. **§4 R7 enforce warn default** — defer prose; default off; nudge already shipped
5. **§5 S02 wave order** — Waves A–D with dependencies
6. **§6 Re-defer registry** — owner + trigger for every defer/reject that stays out of S02

## Role work

1. Read P36 VERIFY-NOTES § Residuals + DR-HANDOFF residuals paragraph + PLAN §2.4/§2.6/§2.7.
2. For each R1–R11: grep/read cited files; confirm partial shipment table; override locked defaults **only** with live evidence (document override in Notes).
3. Apply DESIGN-LOCKS accept sketch; cite Law 19 for R2/R8.
4. For R11: document **doc path** (extend agent-loop / verify docs: critique = `loop apply` with plan_change) — not new MCP tool.
5. Write `RESIDUALS.md`; self-check every INTAKE row has decision + effort + risk.

## Exit criteria

- [ ] `RESIDUALS.md` exists with all R1–R11 decided (accept/defer/reject + S/M/L + risk)
- [ ] Partial P36 S02 shipment documented in §1 or §2
- [ ] Every accept row has S02 scope + test sketch
- [ ] Every defer/reject has rationale; re-defers have owner + trigger in §6
- [ ] R1 cannot silently set PlanExists; R1 signal documented in §3
- [ ] R2/R3 Law 19 adapter targets named (HTTP handlers / MCP tool switch)
- [ ] S02 wave order in §5 matches dependency locks
- [ ] Board row `P37-S00-01` → `done` with Notes (decision counts)

## Minimal todos

- [ ] Verify partial shipment table against live repo
- [ ] Draft RESIDUALS.md §1–§6 using locked defaults
- [ ] Live-code cite each non-trivial decision (file:line or symbol)
- [ ] Record S02 wave order + re-defer registry
- [ ] Update board status + notes only

## Next

`P37-S00-02`
