# P21-S03-00 — Planner: full SelectNext phase cycle

## Metadata
- id: P21-S03-00
- todo_ids: [P21-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock extended `SelectNext` priority table for EXECUTE, TEST, EVALUATE, REFLECT, REPLAN (+ optional EXPLORE). Optional thin PolicyInputs enrichment (D-18). **No product Go this row.**

## References
- [DECISION-LOG.md](../../DECISION-LOG.md) D-03, D-04, D-18
- [WORK-MAP.md](../../WORK-MAP.md) W-04, W-15
- P20 S01 table: [00-PLANNER.md](../../../phase-20-cognitive-deliberation/scopes/scope-01-deliberation-controller/00-PLANNER.md)
- Live: `internal/deliberation/{select.go,types.go,select_test.go}`, `internal/loop/next.go`, `internal/domain/outcomes.go`

## Live inventory (confirmed 2026-08-18)

| Surface | Today (live read) | S03 action |
|---------|-------------------|------------|
| `select.go` L5–28 | **8-row MVP**; comment "EXECUTE never selected"; rows 1–5 + `!plan_critiqued`→CRITIQUE + verify + default ORIENT | **Extend** to 14 rows (insert EXPLORE row 6; rows 8–13 cycle) |
| `types.go` PolicyInputs | **6 fields** (`blocking_uncertainty_count`, `plan_exists`, `plan_critiqued`, `verification_incomplete`, `open_regression`, `p19_saturated`) | **Add** 8 fields (below) |
| `types.go` ReasonCode | **8 codes**; no `execute_pending` etc. | **Add** 6 reason codes |
| `select_test.go` | `TestSelectNext` (10 cases), `TestSelectNextNeverExecuteOnBlockingUncertainty`, ApplyTransition keepers | **Add** 9 named tests |
| `deliberation_packet.go` | `PhaseContextProfile` already maps EXECUTE/TEST/EVALUATE/REFLECT/REPLAN; EXPLORE uses default | **Optional** EXPLORE profile; no schema bump |
| `policy.go` BuildPolicyInputs | Populates 6 legacy fields only | **Stub** new fields (S06 wires live queries) |
| `loop_test.go` L660 | `TestLoopNextExecuteEmphasizesContextAndRelated` tests profile directly, not SelectNext EXECUTE | Keeper unchanged |
| Hop budget | `HopBudget = 12` | Unchanged |
| Schema / compat | max mig **019** (19 files); compat ceiling **19** | No SQL migration |

## FINAL locked defaults (S03-01 must not re-debate)

| Item | Value |
|------|-------|
| EXECUTE gate | **Never** when `blocking_uncertainty_count > 0` (unchanged) |
| ML scoring | **Forbidden** — deterministic first-match table only |
| EXPLORE | **Optional row 12** when `open_decision_alternatives > 0` AND `plan_critiqued=false` — else skip (D-04 keep merged) |
| Float enrichment | **Optional** `plan_confidence` / `requirement_coverage` in `PolicyInputs` JSON for observability only — **must not** change priority order in S03-01 unless explicit threshold locked below |

### New PolicyInputs fields (caller-populated; S06 wires queries in follow-on if needed)

| Field | Type | Semantics |
|-------|------|-----------|
| `ExecutePending` | bool | Plan critiqued, no blocking uncertainty, no open regression, implementation not recorded (no COMPARED change for task) |
| `TestPending` | bool | Execute satisfied; no `kind=test` outcome for task since last change |
| `EvaluationPending` | bool | Test pass recorded; no `kind=evaluation` outcome |
| `ReflectPending` | bool | Evaluation recorded; no reflection row for task |
| `ReplanNeeded` | bool | Open plan-affecting discovery/plan_change OR open reconsideration OR user replan flag |
| `OpenDecisionAlternatives` | int | Count of non-recommended alternatives on open decisions linked to task goal |
| `PlanConfidence` | float64 | Optional 0..1 derived: `1.0` when plan_critiqued && !verification_debt; else `0.5` stub OK in tests |
| `RequirementCoverage` | float64 | Optional 0..1 stub (`0.0` default) — observability only |

### SelectNext priority table (FINAL — extends P20 S01)

Evaluated top-to-bottom; first match wins:

| Pri | Condition | Phase | reason_code |
|-----|-----------|-------|-------------|
| 1 | `hop_count >= 12` OR `stopped` | STOP | hop_budget_exceeded |
| 2 | `p19_saturated` | STOP | p19_saturated |
| 3 | `blocking_uncertainty_count > 0` | INVESTIGATE | blocking_uncertainty |
| 4 | `open_regression` | INVESTIGATE | open_regression |
| 5 | `!plan_exists` | PLAN | plan_missing |
| 6 | `open_decision_alternatives > 0` AND `!plan_critiqued` | EXPLORE | explore_alternatives |
| 7 | `plan_exists && !plan_critiqued` | CRITIQUE | plan_uncritiqued |
| 8 | `execute_pending` AND blocking=0 AND !open_regression | EXECUTE | execute_pending |
| 9 | `test_pending` | TEST | test_pending |
| 10 | `verification_incomplete` | VERIFY | verification_incomplete |
| 11 | `evaluation_pending` | EVALUATE | evaluation_pending |
| 12 | `reflect_pending` | REFLECT | reflect_pending |
| 13 | `replan_needed` | REPLAN | replan_needed |
| 14 | default | ORIENT | continue_orient |

**Note:** Row 6 EXPLORE only when alternatives exist; otherwise row 7 CRITIQUE applies.

## Named tests (S03-01)

| Test | Proves |
|------|--------|
| `TestSelectNextExecuteWhenPending` | execute_pending → EXECUTE |
| `TestSelectNextTestWhenExecuteDone` | test_pending → TEST |
| `TestSelectNextEvaluateWhenTestPass` | evaluation_pending → EVALUATE |
| `TestSelectNextReflectWhenEvaluationDone` | reflect_pending → REFLECT |
| `TestSelectNextReplanWhenFlagged` | replan_needed → REPLAN |
| `TestSelectNextExploreWhenAlternatives` | alternatives + !critiqued → EXPLORE |
| `TestSelectNextNeverExecuteOnBlockingUncertainty` | **Keep** P20 keeper |
| `TestSelectNextFullCycleOrdering` | Table rows 8→12 fire in sequence when flags toggled |
| `TestTransitionPayloadIncludesOptionalScores` | policy_inputs JSON includes floats when set |

## Touch files

- `internal/deliberation/select.go`
- `internal/deliberation/types.go` — PolicyInputs + reason codes
- `internal/deliberation/select_test.go`
- `internal/loop/next.go` — phase emphasis comments/options only (minimal)
- `internal/domain/deliberation.go` — pass-through only if payload shape changes

## Planner work

1. [x] Live inventory `select.go` / `types.go` / `select_test.go` / `policy.go` / `deliberation_packet.go`.
2. [x] Lock W-04 (full SelectNext cycle) + W-15 (PolicyInputs enrichment) defaults — table below unchanged.
3. [x] Thicken `01-full-phase-cycle.md` + `02-scope-review.md` with live before-state, 9 named tests, keeper floor.
4. [x] Update `SCOPE-TODOS.md`.

## Exit criteria

- [x] Extended table + tests locked
- [x] 01/02 thickened enough to implement alone
- [x] No product Go

## Next

**P21-S03-01**
