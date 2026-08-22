# P23-S01-01 — Implement premature implementation gate library

## Metadata
- id: P23-S01-01
- todo_ids: [P23-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [user-codegraph]
- verification: automated

## Objective
Implement `domain.PrematureImplementation` + shared gate evaluator per **S01-00 FINAL locks**. Library-only — **no CLI, no SQL migration**.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S01-00 planner: [00-PLANNER.md](./00-PLANNER.md)
- Live reuse: `internal/loop/policy.go`, `internal/deliberation/select.go`, `internal/loop/deliberation_packet.go` (`loadDeliberationState`), `internal/domain/task_state.go` (DONE alignment reference)

## Session start
Follow agent-loop-protocol. Board edits: **status + notes only**.

## Locked defaults (from S01-00 — do not re-debate)

| Item | Value |
|------|-------|
| Evaluator package | **`internal/loop/gate.go`** (+ `gate_test.go`) — single entrypoint |
| Domain error file | **`internal/domain/gate_errors.go`** (or colocate in existing domain file if convention fits) |
| Reuse | `BuildPolicyInputs`, `deliberation.SelectNext`, `loadDeliberationState` pattern from `deliberation_packet.go` |
| Do **not** fork | SelectNext priority table lives only in `internal/deliberation/select.go` |
| p19 saturation | Mirror `p19SaturatedFromLastStep(st, seed)` from `policy.go` / `next.go` |
| SQL | **None** — pure policy over existing tables |
| CLI | **None** — S02 owns `trace loop gate` |

### Types (FINAL)

```go
// internal/loop/gate.go

type GateFor string

const (
    GateForOrient  GateFor = "orient"
    GateForEdit    GateFor = "edit"
    GateForExecute GateFor = "execute"
    GateForDone    GateFor = "done"
    GateForExport  GateFor = "export"
)

type Violation struct {
    Code             string `json:"code"`
    For              string `json:"for"`
    Message          string `json:"message"`
    RecommendedPhase string `json:"recommended_phase,omitempty"`
    ReasonCode       string `json:"reason_code"`
}

// EvaluateGate is the single harness/product gate entrypoint.
// st is required for task lookup and p19 saturation (mirrors loop next path).
func EvaluateGate(
    ctx context.Context,
    dom *domain.Service,
    plan *planner.Service,
    st *store.Store,
    taskID string,
    gateFor GateFor,
) (allowed bool, violations []Violation, err error)
```

```go
// internal/domain/gate_errors.go

type PrematureImplementation struct {
    TaskID           string
    For              string
    RecommendedPhase string
    ReasonCode       string
    Message          string
}

func (e *PrematureImplementation) Error() string { /* human stderr message */ }
func (e *PrematureImplementation) Code() string  { return "premature_implementation" }
```

- Violation `code` for premature-impl blocks: **`premature_implementation`** (matches `PrematureImplementation.Code()`).
- Orient / infrastructure failures use **`gate_orient_failed`** (or return `err` without violations — pick one path and test it; prefer violation + `allowed=false` for unknown task).
- `recommended_phase` and `reason_code` come from `SelectNext` when policy-driven; empty when N/A (e.g. unknown task).

### GateFor → policy table (FINAL — implement exactly)

Shared prelude (all GateFor except early orient fail):

1. Load task via `st.GetTask(taskID)`; missing → **orient fail** (see orient row).
2. Require non-empty `goal_id`.
3. `BuildPolicyInputs(ctx, dom, plan, taskID, goalID, nil, p19Sat)`.
4. `dState := loadDeliberationState(ctx, dom, taskID, goalID)`.
5. `phase, reason, stopped := deliberation.SelectNext(dState, inputs)`.

| GateFor | Allow when | Block when | Violation `code` | `reason_code` source |
|---------|------------|------------|------------------|----------------------|
| **`orient`** | Task exists, has `goal_id`, goal plan loadable (`CurrentScopeID` + `CurrentDeepPlan` non-nil — mirror `BuildNextPacket` L208–209) | Task not found, empty `goal_id`, or missing plan context | `gate_orient_failed` | `task_not_found`, `missing_goal_id`, `missing_plan_context` |
| **`edit`** | `phase == PhaseExecute && !stopped` | Any other `phase` OR `stopped` (incl. hop budget / P19 saturated) | `premature_implementation` | `SelectNext` `reason` string |
| **`execute`** | `phase == PhaseExecute && reason == ReasonExecutePending && inputs.ExecutePending && !stopped` | Stricter than edit: phase ≠ EXECUTE, reason ≠ execute_pending, `!inputs.ExecutePending`, or stopped | `premature_implementation` | `SelectNext` `reason` or `execute_not_pending` when phase=EXECUTE but flag false |
| **`done`** | All clear: `!inputs.VerificationIncomplete`, `!inputs.OpenRegression`, `inputs.BlockingUncertaintyCount == 0`, and `phase` ∈ `{PhaseOrient, PhaseStop}` | Verification debt, open regression, blocking uncertainty, OR `phase` ∈ `{INVESTIGATE, EXPLORE, PLAN, CRITIQUE, EXECUTE, TEST, VERIFY, EVALUATE, REFLECT, REPLAN}` | `premature_implementation` (or distinct `verification_debt` / `open_regression` codes if multiple violations returned — prefer **one violation per block** with accurate `reason_code`) | Matching signal: `verification_incomplete`, `open_regression`, `blocking_uncertainty`, or SelectNext reason |
| **`export`** | Same as **`done`** in S01 | Same as **`done`** | Same as **`done`** | Same as **`done`** (S03 adds export-honesty extensions) |

#### Edit vs SelectNext (explicit)

- **Edit gate = EXECUTE-only guard.** Do not invent a separate edit policy table.
- When `SelectNext` returns `INVESTIGATE` (blocking uncertainty / open regression), `PLAN`, `CRITIQUE`, `EXPLORE`, `TEST`, etc., edit **blocks** with `recommended_phase` = that phase and `reason_code` = SelectNext reason.
- When `SelectNext` returns `STOP`, edit **blocks** (`recommended_phase` = `STOP`, reason = hop budget or P19 saturated).
- **Allow edit** iff SelectNext would send the agent to `EXECUTE` (implementation phase).

#### Done vs TransitionTask (explicit)

- **`GateForDone` in S01** covers **deliberation / verification readiness** only (table above).
- **`domain.TransitionTask` DONE** also enforces review PASS, `--as-operator`, missing capabilities — **not duplicated in S01**. S03 `--enforce` may call `EvaluateGate(..., GateForDone)` **before** transition; transition hatches remain unless S03 wires combined enforce.
- Align verification debt signal: `inputs.VerificationIncomplete` from `BuildPolicyInputs` → `dom.HasVerificationDebt`.

## Files to create/modify

| File | Action |
|------|--------|
| `internal/loop/gate.go` | **Create** — `GateFor`, `Violation`, `EvaluateGate`, gateFor dispatch |
| `internal/loop/gate_test.go` | **Create** — table-driven tests (named below) |
| `internal/domain/gate_errors.go` | **Create** — `PrematureImplementation` with `Code()` |
| `internal/domain/gate_errors_test.go` | **Create** — `TestPrematureImplementation_Code` |

**Do not modify:** `cmd/trace/*`, `internal/deliberation/select.go`, loop CLI commands.

## Named unit tests (minimum — S01-01 must implement all)

| Test name | Setup gist | Expect |
|-----------|------------|--------|
| `TestEvaluateGate_Orient_UnknownTask` | random UUID | `allowed=false`, orient violation or error |
| `TestEvaluateGate_Orient_MissingGoal` | task without goal_id | `allowed=false`, `missing_goal_id` |
| `TestEvaluateGate_Orient_MissingPlan` | goal without deep plan | `allowed=false`, `missing_plan_context` |
| `TestEvaluateGate_Orient_OK` | seeded goal/task/plan | `allowed=true` |
| `TestEvaluateGate_Edit_BlockingUncertainty` | blocking uncertainty linked | `allowed=false`, code `premature_implementation`, reason `blocking_uncertainty`, recommended `INVESTIGATE` |
| `TestEvaluateGate_Edit_OpenRegression` | open regression | block, reason `open_regression` |
| `TestEvaluateGate_Edit_PlanMissing` | no plan | block, recommended `PLAN`, reason `plan_missing` |
| `TestEvaluateGate_Edit_PlanUncritiqued` | plan, not critiqued | block, recommended `CRITIQUE`, reason `plan_uncritiqued` |
| `TestEvaluateGate_Edit_ExecuteReady` | plan critiqued, gates clear | `allowed=true` |
| `TestEvaluateGate_Edit_HopBudgetStopped` | hop_count ≥ HopBudget | block, recommended `STOP` |
| `TestEvaluateGate_Execute_NotExecutePending` | impl signal present (past execute) | block even if phase might differ |
| `TestEvaluateGate_Execute_ExecutePendingClear` | execute_pending true | `allowed=true` |
| `TestEvaluateGate_Done_VerificationDebt` | change without verification | block, reason `verification_incomplete` |
| `TestEvaluateGate_Done_OpenRegression` | open regression | block |
| `TestEvaluateGate_Done_DeliberationIncomplete` | e.g. test_pending | block with SelectNext reason |
| `TestEvaluateGate_Done_Clean` | full cycle clear | `allowed=true` |
| `TestEvaluateGate_Export_SameAsDone` | debt vs clean | mirrors done allow/block |
| `TestPrematureImplementation_Code` | direct struct | `Code() == "premature_implementation"` |

Use existing loop test helpers (`openLoopTestStore`, `seedGoalTaskPlan`, `markPlanCritiqued` from `policy_test.go`) where possible.

## Implementation notes

- Export `EvaluateGate` only; keep gateFor-specific helpers unexported (`evaluateOrient`, `evaluateEdit`, …).
- Human `message` on violations: short, actionable (e.g. `"edit blocked: recommended phase INVESTIGATE (blocking_uncertainty)"`).
- Multiple violations: return **at most one** violation per `EvaluateGate` call for S01 (simplest JSON for S02 CLI).
- `errors.Is` support optional for `PrematureImplementation` if useful for S03.

## Exit criteria

- [ ] All named unit tests green
- [ ] `go test ./internal/loop/... ./internal/domain/...` green
- [ ] P19/P20 loop keeper tests unchanged: `go test ./cmd/trace -run 'TestLoopNext|TestLoopApply|TestLoopStatus'`
- [ ] No duplicate SelectNext table
- [ ] Violation JSON fields match ENFORCEMENT.md example
- [ ] No CLI / no SQL migration / no daemon

## Minimal todos

- [ ] Add `PrematureImplementation` domain error with stable `Code()`
- [ ] Add `internal/loop/gate.go` with `EvaluateGate` + GateFor dispatch
- [ ] Implement gateFor table exactly as locked above
- [ ] Add all named tests in `gate_test.go`
- [ ] Run keeper tests; fix regressions only in gate files
- [ ] Board row: status + notes only
