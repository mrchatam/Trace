package planner

import (
	"errors"
)

// ErrNotFound is returned when a required plan row is missing or unset.
var ErrNotFound = errors.New("planner: not found")

// ErrReplanBudgetExceeded is returned when auto_replan_count >= MaxAutoReplans.
var ErrReplanBudgetExceeded = errors.New("planner: auto-replan budget exceeded")

// DefaultMaxAutoReplans is the G16 / DR-CHURN default (N=5).
const DefaultMaxAutoReplans = 5

// Severity constants re-exported from store/domain vocabulary.
const (
	SeverityINFO          = "INFO"
	SeverityPlanAffecting = "PLAN_AFFECTING"
	SeverityBlocking      = "BLOCKING"
)

// ErrValidation is returned for invalid caller input.
type ErrValidation struct {
	Msg string
}

func (e *ErrValidation) Error() string {
	return "planner: " + e.Msg
}

// ErrNotCurrent is returned when DeepPlan targets a non-current scope.
type ErrNotCurrent struct {
	ScopeID   string
	CurrentID string
}

func (e *ErrNotCurrent) Error() string {
	return "planner: deep plan requires current scope (got " + e.ScopeID + ", current " + e.CurrentID + ")"
}

// CoarsePlanInput creates a coarse phase/scope tree under an existing goal.
type CoarsePlanInput struct {
	GoalID string
	Phases []PhaseInput
	Actor  string // optional; event payload
}

// PhaseInput is one coarse phase with nested scopes.
type PhaseInput struct {
	Title  string
	Body   string
	Scopes []ScopeInput
}

// ScopeInput is one coarse scope (title + optional body summary).
type ScopeInput struct {
	Title string
	Body  string
}

// CoarsePlan is the created hierarchy.
type CoarsePlan struct {
	GoalID string      `json:"goal_id"`
	Phases []PhaseView `json:"phases"`
}

// PhaseView is a phase with nested scopes.
type PhaseView struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	Body   string      `json:"body"`
	Ord    int         `json:"ord"`
	Status string      `json:"status"`
	Scopes []ScopeView `json:"scopes"`
}

// ScopeView is a coarse scope row.
type ScopeView struct {
	ID              string `json:"id"`
	PhaseID         string `json:"phase_id"`
	Title           string `json:"title"`
	Body            string `json:"body"`
	Ord             int    `json:"ord"`
	Status          string `json:"status"`
	AutoReplanCount int    `json:"auto_replan_count"`
}

// ScopeRef is a lightweight ordered scope identity.
type ScopeRef struct {
	ID       string
	PhaseID  string
	GoalID   string
	Title    string
	Body     string
	Ord      int // scope ord within phase
	PhaseOrd int
}

// WorkItem is a caller-supplied deep-plan work item.
type WorkItem struct {
	Title string `json:"title"`
	Notes string `json:"notes"`
}

// DeepPlanDocument is the locked JSON shape stored in scope_deep_plans.content_json.
type DeepPlanDocument struct {
	ScopeID          string     `json:"scope_id"`
	ExitCriteria     []string   `json:"exit_criteria"`
	Constraints      []string   `json:"constraints"`
	WorkItems        []WorkItem `json:"work_items"`
	LookaheadScopeID string     `json:"lookahead_scope_id"`
	LookaheadSummary string     `json:"lookahead_summary"`
}

// DeepPlanInput writes/replaces the ACTIVE deep plan for the current scope.
type DeepPlanInput struct {
	ScopeID          string
	ExitCriteria     []string
	Constraints      []string
	WorkItems        []WorkItem
	LookaheadSummary string // optional; when non-empty, may set lookahead scope body
	Actor            string
}

// SupersedeInput replaces the ACTIVE deep plan for a scope (S02 hook; not current-gated).
type SupersedeInput struct {
	ScopeID          string
	ExitCriteria     []string
	Constraints      []string
	WorkItems        []WorkItem
	LookaheadScopeID string
	LookaheadSummary string
	Actor            string
}

// DeepPlanResult is the written revision plus resolved lookahead identity.
type DeepPlanResult struct {
	RevisionID       string
	Document         DeepPlanDocument
	SupersededCount  int64
	LookaheadScopeID string
}

// PlanView is the GetPlan snapshot.
type PlanView struct {
	GoalID           string            `json:"goal_id"`
	CurrentScopeID   *string           `json:"current_scope_id"`
	Phases           []PhaseView       `json:"phases"`
	CurrentDeepPlan  *DeepPlanDocument `json:"current_deep_plan"`
	LookaheadScopeID string            `json:"lookahead_scope_id"`
	LookaheadSummary string            `json:"lookahead_summary"`
}

// ApplyDiscoveryReplanInput drives discovery→PlanChange replan with churn controls.
type ApplyDiscoveryReplanInput struct {
	DiscoveryID      string
	ScopeID          string // deep plan target (need not be current)
	PlanChangeID     string // optional: if empty, create from title/body
	PlanChangeTitle  string
	PlanChangeBody   string
	ExitCriteria     []string
	Constraints      []string
	WorkItems        []WorkItem
	LookaheadScopeID string
	LookaheadSummary string
	MaxAutoReplans   int // 0 → DefaultMaxAutoReplans
	Actor            string
}

// ApplyDiscoveryReplanResult reports whether auto-replan ran and the resulting state.
type ApplyDiscoveryReplanResult struct {
	DiscoveryID       string
	ScopeID           string
	PlanChangeID      string
	AutoReplanApplied bool
	Reason            string // e.g. severity_info
	RevisionID        string
	AutoReplanCount   int
	SupersededCount   int64
}
