# Scope 00 — board map

**S00 triage** — accept/defer/reject per R1–R11. Serial: **P37-S00-00 → P37-S00-01 → P37-S00-02**. Artifact: `RESIDUALS.md` (written in **S00-01**, reviewed in **S00-02**). Do **not** start S01 until S00-02 PASS. Do **not** write product code. Planner (**S00-00**) does **not** author `RESIDUALS.md`.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 635 | P37-S00-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 + this file |
| 636 | P37-S00-01 | [01-triage.md](01-triage.md) | Implementer | Author `RESIDUALS.md` + live cites |
| 637 | P37-S00-02 | [02-review.md](02-review.md) | Reviewer | Checklist vs DESIGN-LOCKS + INTAKE + live cites |

## Planner-locked live facts (for 01) — verified 2026-08-22 (P37-S00-00)

Re-verify during S00-01; override only with evidence.

| ID | Live fact | Cite |
|----|-----------|------|
| R1 | No PlanExists bridge; no bootstrap advisory in loop status | P36 PLAN §2.4 defer; `internal/loop/policy.go` sets `PlanExists` from store only (`:45–48`, `:91`); no plan-change-count advisory |
| R2 | HTTP plan read-only | `internal/httpapi/server.go:282` — `GET /v1/plans`; `handlers_p1.go:115` `handleListPlans`; no POST plan routes; `api/openapi.yaml:760` GET only |
| R3 | No MCP gate action | `internal/mcp/tools_loop.go:39–47` — `next\|apply\|status`; CLI `cmd/trace/loop.go:122–176` `cmdLoopGate` exists; HTTP `GET /v1/loop/gate` in `handlers_loop.go:104` |
| R4 | Bootstrap shipped; refinement note gap | `cmd/trace/plan.go:55–72` help + `:334–360` bootstrap stderr; lacks PLAN §2.2 “human refinement via create-coarse/deep expected” |
| R5 | Warning helper partial | `internal/planner/advisory.go:13–41` `GoalStructureWarning`; tests `advisory_test.go`; wired `plan.go:297–298`, `tools_plan.go:144–156`; **gap:** `StatusResult` no `advisories` (`internal/loop/apply.go:196–209`) |
| R6 | Warn helper untested | `internal/config/enforce.go:43–57`; called `cmd/trace/init.go:50`; no `enforce_test.go` |
| R7 | Enforce doc-only in P36 | P36 PLAN §2.6 accept document-only; `LoadEnforceMode` missing config → `EnforceOff` (`enforce.go:25–29`); stderr nudge shipped |
| R8 | TaskDetail > Overview for plan UX | `web/src/screens/TaskDetail.tsx:205–209` bootstrap advisory; `Overview.tsx:72–78` GateStrip only |
| R9 | Feet-seller post-P36 bootstrap | VERIFY Block 6 — minimal progressive plan; `plan_uncritiqued` post-bootstrap expected |
| R10 | GUI verify deferred | VERIFY Block 4 — pre-bootstrap API + GateStrip sufficient; live browser deferred |
| R11 | CLI critique gap | VERIFY Block 1 partial — CLI chain → `plan_uncritiqued`; Block 0 MCP greenfield covers full path |

## Partial P36 S02 shipment (S00-01 must confirm)

| Component | Shipped P36 S02 | Residual |
|-----------|-----------------|----------|
| Bootstrap CLI/MCP | `internal/planner/bootstrap.go`, `cmd/trace/plan.go`, `tools_plan.go` | R4 help note |
| Goal structure warning | `advisory.go`, plan show + MCP show | R5 loop status `advisories[]` |
| Enforce nudge | `WarnIfTraceDirWithoutConfig`, init call | R6 test; R7 default flip deferred |
| Terminal honesty GUI | TaskDetail bootstrap copy, GateStrip adapter | R8 Overview surface |
| PlanExists bridge | — | R1 full gap |
| HTTP/MCP gate parity | GET gate HTTP only | R2 POST plans, R3 MCP gate |

## Preliminary triage (S00-00 locks — 01 implements in RESIDUALS.md)

| ID | Decision | Effort | Risk | S02? |
|----|----------|--------|------|------|
| R1 | **accept** advisory-only | M | low | Y — Wave A |
| R2 | **accept** HTTP POST plan routes | M | low–med | Y — Wave B |
| R3 | **accept** MCP `trace_loop gate` | S | low | Y — Wave B |
| R4 | **accept** bootstrap help | S | low | Y — Wave B |
| R5 | **accept** status `advisories[]` | S | low | Y — Wave A |
| R6 | **accept** unit test | S | low | Y — Wave C |
| R7 | **defer** default enforce flip | — | med | N — §6 re-defer |
| R8 | **accept** Overview minimal | M | low | Y — Wave D |
| R9 | **defer** → S03 dogfood | S | low | N |
| R10 | **accept** → S03 browser | S | low | N |
| R11 | **accept** doc path | S | low | Y — Wave C |

**Reject always:** silent PlanExists bridge; weakening active `plan_missing`; business-logic fork in `web/`; feet-seller history rewrite.

## R1 advisory signal (locked for 01)

| Field | Value |
|-------|-------|
| Trigger | `!PlanExists` ∧ goal-linked `plan_changes ≥ 1` |
| N | **1** |
| Surface | `trace.loop.status.v1` → `advisories[]` code `bootstrap_recommended` |
| Must not | Set `PlanExists`, mutate planner store |

## R2 / R3 Law 19 adapter targets (locked for 01)

| Residual | Adapter file | Library call |
|----------|--------------|--------------|
| R2 POST bootstrap (min) | `internal/httpapi/handlers_p1.go` (new handler) | `planner.Service.BootstrapFromPlanChanges` — mirror `cmd/trace/plan.go:351` |
| R2 POST create-coarse (if S01 expands) | same pattern | `planner.Service.CreateCoarsePlan` — mirror `cmd/trace/plan.go:88` |
| R3 MCP gate | `internal/mcp/tools_loop.go` switch + helper | `loop.EvaluateGate` — mirror `cmd/trace/loop.go:155–162` |
| R3 HTTP (exists) | `internal/httpapi/handlers_loop.go:104` | precedent only — do not duplicate |

## R11 path lock (locked for 01)

**Accept:** agent-workflow / verify documentation — critique after bootstrap = `trace loop apply` with `plan_changes` envelope (see P36 Block 0 MCP test pattern).

**Reject:** new MCP `trace_plan action=critique-seed` or similar — duplicates existing loop apply path.

## S02 wave order (locked for RESIDUALS §5)

```text
A: R5 advisories[] + R1 bootstrap_recommended
B: R3 MCP gate | R2 HTTP POST plans | R4 help  (parallel)
C: R6 test | R11 docs
D: R8 Overview GUI
S03: R9 dogfood | R10 browser
Out: R7 defer
```

## Source docs (must read)

- Phase 36 [`DR-HANDOFF.md`](../../../phase-36-gate-honesty-terminal-tasks/DR-HANDOFF.md)
- Phase 36 [`VERIFY-NOTES.md`](../../../phase-36-gate-honesty-terminal-tasks/scopes/scope-03-verify/VERIFY-NOTES.md) § Residuals
- Phase 36 [`PLAN.md`](../../../phase-36-gate-honesty-terminal-tasks/scopes/scope-01-plan/PLAN.md) §2.4, §2.6, §2.7, defer touch-list
- [`INTAKE.md`](../../INTAKE.md) R1–R11
- [`DESIGN-LOCKS.md`](../../DESIGN-LOCKS.md)

## Triage rejects (01 must document if proposed)

1. Silent PlanExists bridge (fake progressive plan from plan-changes alone).
2. Weakening active-work `plan_missing` enforcement.
3. Business-logic fork in `web/` (Law 19).
4. Rewriting feet-seller task history.
5. New MCP critique-seed tool (R11 — prefer doc path).

## Re-defer registry seeds (RESIDUALS §6)

| ID | Owner | Trigger |
|----|-------|---------|
| R7 | Human / product | Explicit decision to flip missing-config default to `warn` |
| R8 full plan screen | Phase 38 planner | S03 shows Overview minimal insufficient |
| R9 deep refinement | Human dogfood | Documented path; verify quality on feet-seller |

## Out of scope

- Product implementation (S02)
- Authoring `RESIDUALS.md` in S00-00 (planner row)
- Re-litigating Phase 36 core fixes (MCP plan, bootstrap, terminal advisory)
