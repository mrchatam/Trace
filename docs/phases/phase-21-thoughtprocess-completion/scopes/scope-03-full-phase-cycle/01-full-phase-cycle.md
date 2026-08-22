# P21-S03-01 — Implement: full SelectNext phase cycle

## Metadata
- id: P21-S03-01
- todo_ids: [P21-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective
Extend `SelectNext` from the P20 **8-row MVP** to the S03-00 **14-row FINAL** table: add EXPLORE (optional) plus EXECUTE, TEST, EVALUATE, REFLECT, REPLAN rows; extend `PolicyInputs` with cycle-pending booleans/counts and optional observability floats (D-18). Keep policy **pure** (no I/O). Wire new fields in unit tests only — live `BuildPolicyInputs` population deferred to S06 unless trivial pass-through is zero-cost.

## Session start
Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: **status + notes only**.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- [DECISION-LOG.md](../../DECISION-LOG.md) D-03, D-04, D-18
- [WORK-MAP.md](../../WORK-MAP.md) W-04, W-15
- P20 baseline: [scope-01/00-PLANNER.md](../../../phase-20-cognitive-deliberation/scopes/scope-01-deliberation-controller/00-PLANNER.md)
- Live: `internal/deliberation/{select.go,types.go,select_test.go}`, `internal/loop/{policy.go,deliberation_packet.go,next.go}`, `internal/domain/deliberation.go`, `cmd/trace/loop_test.go`

## Locked defaults (from S03-00 — do not re-debate)

| Item | Value |
|------|-------|
| Priority model | Deterministic **first-match** table only — **no ML scoring** (D-21) |
| Hop budget | **12** (`HopBudget` unchanged) |
| EXECUTE gate | **Never** when `blocking_uncertainty_count > 0` (rows 3 wins; row 8 redundant guard OK) |
| EXPLORE (D-04) | Row **6** only when `open_decision_alternatives > 0` AND `!plan_critiqued`; else fall through to CRITIQUE |
| Float fields (D-18) | `plan_confidence` / `requirement_coverage` — **observability only**; must **not** affect priority order |
| Schema / compat | Max mig **019**; compat ceiling **19**; no SQL migration |
| `BuildPolicyInputs` | May leave new fields at zero/false; S06 wires live queries later |
| Loop packet schema | **No bump** — `trace.loop.next.v1` unchanged; `PhaseContextProfile` emphasis only |

### SelectNext priority table (FINAL — 14 rows)

Evaluated top-to-bottom; first match wins:

| Pri | Condition | Phase | reason_code |
|-----|-----------|-------|-------------|
| 1 | `hop_count >= 12` OR `stopped` | STOP | `hop_budget_exceeded` |
| 2 | `p19_saturated` | STOP | `p19_saturated` |
| 3 | `blocking_uncertainty_count > 0` | INVESTIGATE | `blocking_uncertainty` |
| 4 | `open_regression` | INVESTIGATE | `open_regression` |
| 5 | `!plan_exists` | PLAN | `plan_missing` |
| 6 | `open_decision_alternatives > 0` AND `!plan_critiqued` | EXPLORE | `explore_alternatives` |
| 7 | `plan_exists && !plan_critiqued` | CRITIQUE | `plan_uncritiqued` |
| 8 | `execute_pending` AND blocking=0 AND `!open_regression` | EXECUTE | `execute_pending` |
| 9 | `test_pending` | TEST | `test_pending` |
| 10 | `verification_incomplete` | VERIFY | `verification_incomplete` |
| 11 | `evaluation_pending` | EVALUATE | `evaluation_pending` |
| 12 | `reflect_pending` | REFLECT | `reflect_pending` |
| 13 | `replan_needed` | REPLAN | `replan_needed` |
| 14 | default | ORIENT | `continue_orient` |

**Note:** Row 6 EXPLORE requires alternatives; when `!plan_critiqued` without alternatives, row 7 CRITIQUE applies.

### PolicyInputs extensions

Add to `internal/deliberation/types.go` (`PolicyInputs` struct):

| Field | Go type | JSON key | Semantics (caller-populated) |
|-------|---------|----------|------------------------------|
| `ExecutePending` | bool | `execute_pending` | Plan critiqued, no blocking uncertainty, no open regression, implementation not recorded (no COMPARED change for task) |
| `TestPending` | bool | `test_pending` | Execute satisfied; no `kind=test` outcome for task since last change |
| `EvaluationPending` | bool | `evaluation_pending` | Test pass recorded; no `kind=evaluation` outcome |
| `ReflectPending` | bool | `reflect_pending` | Evaluation recorded; no reflection row for task |
| `ReplanNeeded` | bool | `replan_needed` | Open plan-affecting discovery/plan_change OR open reconsideration OR user replan flag |
| `OpenDecisionAlternatives` | int | `open_decision_alternatives` | Count of non-recommended alternatives on open decisions linked to task goal |
| `PlanConfidence` | float64 | `plan_confidence` | Optional 0..1 derived stub: `1.0` when plan_critiqued && !verification_debt; else `0.5` OK in tests |
| `RequirementCoverage` | float64 | `requirement_coverage` | Optional 0..1 stub (`0.0` default) — observability only |

Add matching `ReasonCode` constants:

```go
ReasonExploreAlternatives  = "explore_alternatives"
ReasonExecutePending       = "execute_pending"
ReasonTestPending          = "test_pending"
ReasonEvaluationPending    = "evaluation_pending"
ReasonReflectPending       = "reflect_pending"
ReasonReplanNeeded         = "replan_needed"
```

## Live inventory (before — confirmed S03-00)

| Surface | Location | Today |
|---------|----------|-------|
| SelectNext | `select.go` L5–28 | **8 rows**; comment "EXECUTE never selected"; no EXPLORE/EXECUTE/TEST/EVALUATE/REFLECT/REPLAN |
| PolicyInputs | `types.go` L47–54 | **6 fields** only |
| Reason codes | `types.go` L35–44 | **8 codes**; no cycle-pending codes |
| Tests | `select_test.go` | `TestSelectNext` table (10 cases), `TestSelectNextNeverExecuteOnBlockingUncertainty`, ApplyTransition keepers |
| Phase profiles | `deliberation_packet.go` L38–81 | EXECUTE/TEST/EVALUATE/REFLECT/REPLAN profiles **exist** but SelectNext never returns those phases today |
| BuildPolicyInputs | `policy.go` L64–71 | Populates 6 legacy fields only |
| Loop EXECUTE test | `loop_test.go` L660 | Tests `PhaseContextProfile(PhaseExecute)` directly — not SelectNext integration |

## Requirements

1. **`types.go`** — extend `PolicyInputs` + `ReasonCode` constants per locked table above. Use `omitempty` on optional float JSON tags if idiomatic; floats must serialize when non-zero in transition payload tests.

2. **`select.go`** — replace MVP table with 14-row FINAL:
   - Insert row 6 EXPLORE before CRITIQUE.
   - Change CRITIQUE to `plan_exists && !plan_critiqued` (row 5 already catches missing plan).
   - Insert rows 8–13 between CRITIQUE and VERIFY/default.
   - Update package comment: EXECUTE **is** selected when `execute_pending` and upstream gates clear.
   - Keep function pure — no store/planner calls.

3. **`select_test.go`** — add **9 named tests** (keep all existing tests green):

   | Test | Proves |
   |------|--------|
   | `TestSelectNextExecuteWhenPending` | `execute_pending=true`, plan critiqued, clear blocking/regression → EXECUTE / `execute_pending` |
   | `TestSelectNextTestWhenExecuteDone` | `test_pending=true`, upstream clear → TEST / `test_pending` |
   | `TestSelectNextEvaluateWhenTestPass` | `evaluation_pending=true` → EVALUATE / `evaluation_pending` |
   | `TestSelectNextReflectWhenEvaluationDone` | `reflect_pending=true` → REFLECT / `reflect_pending` |
   | `TestSelectNextReplanWhenFlagged` | `replan_needed=true` → REPLAN / `replan_needed` |
   | `TestSelectNextExploreWhenAlternatives` | `open_decision_alternatives=2`, `plan_critiqued=false`, plan exists → EXPLORE / `explore_alternatives` |
   | `TestSelectNextNeverExecuteOnBlockingUncertainty` | **Keep** P20 keeper — must not regress |
   | `TestSelectNextFullCycleOrdering` | Sequential table: set only the highest-priority true flag among rows 8→12; prove EXECUTE→TEST→EVALUATE→REFLECT order as flags advance |
   | `TestTransitionPayloadIncludesOptionalScores` | `ApplyTransition` JSON: `policy_inputs.plan_confidence` and `requirement_coverage` present when set |

   Update `TestTransitionPayloadJSONRequiredFields` policy_inputs key list to include new boolean/int fields (not necessarily floats when zero).

4. **`deliberation_packet.go`** (minimal) — confirm `PhaseContextProfile` covers EXPLORE if reasonable (default profile OK); no schema changes. Optional: add `PhaseExplore` case with `OpenUncertaintyCap: 12`, `RelatedDepth: 2` if cheap.

5. **`policy.go`** — **no live wiring required** this scope. New fields remain zero/false until S06. Do not break existing 6-field population.

6. **`internal/domain/deliberation.go`** — pass-through only if `TransitionPayload` / event JSON needs no changes (should inherit extended `PolicyInputs` automatically).

7. **No migration, seed, or compat changes.**

## Touch files

- `internal/deliberation/types.go` — PolicyInputs + ReasonCode
- `internal/deliberation/select.go` — 14-row table
- `internal/deliberation/select_test.go` — 9 new + updated JSON field test
- `internal/loop/deliberation_packet.go` — optional EXPLORE profile
- `internal/loop/policy.go` — verify-only (no wiring unless zero-cost)
- `internal/domain/deliberation.go` — verify-only

## Keeper floor

```bash
go test ./internal/deliberation/... -count=1 -run 'TestSelectNext|TestSelectNextNeverExecuteOnBlockingUncertainty|TestSelectNextExecuteWhenPending|TestSelectNextTestWhenExecuteDone|TestSelectNextEvaluateWhenTestPass|TestSelectNextReflectWhenEvaluationDone|TestSelectNextReplanWhenFlagged|TestSelectNextExploreWhenAlternatives|TestSelectNextFullCycleOrdering|TestTransitionPayloadIncludesOptionalScores|TestApplyTransition|TestTransitionPayloadJSONRequiredFields'
go test ./cmd/trace -count=1 -run 'TestLoopNextDeliberationSectionPresent|TestLoopNextExecuteEmphasizesContextAndRelated|TestLoopApplyDeliberationTransitionEvent'
go test ./internal/domain/... -count=1 -run 'TestApplyDeliberationTransition|TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition|TestHasOpenRegressionFeedsApplyDeliberationTransition'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Minimal todos

- [ ] Extend `PolicyInputs` + reason codes in `types.go`
- [ ] Implement 14-row `SelectNext` in `select.go`
- [ ] Add 9 named tests + update JSON required-fields test
- [ ] Optional EXPLORE phase profile in `deliberation_packet.go`
- [ ] Run keeper floor; paste PASS evidence in board Notes

## Exit criteria

- [ ] 9 new + all legacy SelectNext / ApplyTransition tests PASS
- [ ] Priority table matches S03-00 (14 rows) in code
- [ ] EXECUTE never when `blocking_uncertainty_count > 0`
- [ ] Float fields do not affect phase selection
- [ ] No mig 020+; compat ceiling 19 unchanged
- [ ] Board row status + Notes only

## Next

**P21-S03-02**
