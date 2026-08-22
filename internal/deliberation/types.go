// Package deliberation is the Phase 20 state-driven deliberation controller.
//
// SelectNext is a pure policy function (no I/O). Persistence lives in
// internal/store; domain applies transitions and appends events. This package
// does not store raw chain-of-thought and has no CLI.
package deliberation

// HopBudget is the maximum number of applied phase hops per seed (N=12).
const HopBudget = 12

// SaturationEmptyThreshold is how many consecutive pure-empty applies trigger P19 STOP.
// Discoveries-only applies do not increment the counter (see NextConsecutiveEmptyApplies).
const SaturationEmptyThreshold = 2

// EventTransition is the append-only audit event type for phase hops.
const EventTransition = "deliberation.transition"

// Phase is the deliberation controller vocabulary.
type Phase string

const (
	PhaseOrient      Phase = "ORIENT"
	PhaseInvestigate Phase = "INVESTIGATE"
	PhaseExplore     Phase = "EXPLORE"
	PhasePlan        Phase = "PLAN"
	PhaseCritique    Phase = "CRITIQUE"
	PhaseExecute     Phase = "EXECUTE"
	PhaseTest        Phase = "TEST"
	PhaseVerify      Phase = "VERIFY"
	PhaseEvaluate    Phase = "EVALUATE"
	PhaseReflect     Phase = "REFLECT"
	PhaseReplan      Phase = "REPLAN"
	PhaseStop        Phase = "STOP"
)

// ReasonCode is the first-match SelectNext reason.
type ReasonCode string

const (
	ReasonHopBudgetExceeded      ReasonCode = "hop_budget_exceeded"
	ReasonP19Saturated           ReasonCode = "p19_saturated"
	ReasonBlockingUncertainty    ReasonCode = "blocking_uncertainty"
	ReasonOpenRegression         ReasonCode = "open_regression"
	ReasonPlanMissing            ReasonCode = "plan_missing"
	ReasonPlanUncritiqued        ReasonCode = "plan_uncritiqued"
	ReasonExploreAlternatives    ReasonCode = "explore_alternatives"
	ReasonExecutePending         ReasonCode = "execute_pending"
	ReasonTestPending            ReasonCode = "test_pending"
	ReasonVerificationIncomplete ReasonCode = "verification_incomplete"
	ReasonEvaluationPending      ReasonCode = "evaluation_pending"
	ReasonReflectPending         ReasonCode = "reflect_pending"
	ReasonReplanNeeded           ReasonCode = "replan_needed"
	ReasonContinueOrient         ReasonCode = "continue_orient"
)

// PolicyInputs are caller-populated signals (S06 wires queries). Zero/false is a valid stub.
type PolicyInputs struct {
	BlockingUncertaintyCount int     `json:"blocking_uncertainty_count"`
	PlanExists               bool    `json:"plan_exists"`
	PlanCritiqued            bool    `json:"plan_critiqued"`
	VerificationIncomplete   bool    `json:"verification_incomplete"`
	OpenRegression           bool    `json:"open_regression"`
	P19Saturated             bool    `json:"p19_saturated"`
	ExecutePending           bool    `json:"execute_pending"`
	TestPending              bool    `json:"test_pending"`
	EvaluationPending        bool    `json:"evaluation_pending"`
	ReflectPending           bool    `json:"reflect_pending"`
	ReplanNeeded             bool    `json:"replan_needed"`
	OpenDecisionAlternatives int     `json:"open_decision_alternatives"`
	PlanConfidence           float64 `json:"plan_confidence,omitempty"`
	RequirementCoverage      float64 `json:"requirement_coverage,omitempty"`
}

// State is in-memory deliberation_state (one row per task_id).
type State struct {
	TaskID                  string
	GoalID                  string
	CurrentPhase            Phase
	HopCount                int
	LastPhase               Phase
	PlanCritiqued           bool
	Stopped                 bool
	StopReason              string
	ConsecutiveEmptyApplies int
	UpdatedAt               string
}

// NextConsecutiveEmptyApplies updates the persisted empty-apply counter.
// Pure empty (no plan, no spawn, no discovery writes) increments; plan/spawn clears;
// discoveries-only leaves the counter unchanged.
func NextConsecutiveEmptyApplies(prev, newPlanChanges, newSpawnedTasks, discoveryWrites int) int {
	if newPlanChanges > 0 || newSpawnedTasks > 0 {
		return 0
	}
	if discoveryWrites > 0 {
		return prev
	}
	return prev + 1
}

// SaturatedFromCounter reports P19 saturation from max-iterations or consecutive empties.
func SaturatedFromCounter(maxIterationsReached bool, consecutiveEmptyApplies int) bool {
	return maxIterationsReached || consecutiveEmptyApplies >= SaturationEmptyThreshold
}

// TransitionPayload is the required JSON body of a deliberation.transition event.
type TransitionPayload struct {
	TaskID       string       `json:"task_id"`
	GoalID       string       `json:"goal_id"`
	FromPhase    Phase        `json:"from_phase"`
	ToPhase      Phase        `json:"to_phase"`
	ReasonCode   ReasonCode   `json:"reason_code"`
	HopCount     int          `json:"hop_count"`
	PolicyInputs PolicyInputs `json:"policy_inputs"`
}

// InitialState returns a new row defaulting to ORIENT with hop_count 0.
func InitialState(taskID, goalID string) State {
	return State{
		TaskID:       taskID,
		GoalID:       goalID,
		CurrentPhase: PhaseOrient,
	}
}
