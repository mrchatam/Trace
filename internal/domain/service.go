package domain

import (
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// Entity type strings (lowercase singular).
const (
	EntityGoal          = "goal"
	EntityTask          = "task"
	EntityDecision      = "decision"
	EntityAssumption    = "assumption"
	EntityDiscovery     = "discovery"
	EntityPlanChange    = "plan_change"
	EntityClaim         = "claim"
	EntityEvidence      = "evidence"
	EntityReview        = "review"
	EntityPlanScope     = "plan_scope"
	EntityUncertainty   = "uncertainty"
	EntityHypothesis    = "hypothesis"
	EntityChange        = "change"
	EntityEffect        = "effect"
	EntityOutcomeResult = "outcome_result"
	EntityBaseline      = "baseline"
	EntityRegression    = "regression"
	EntityReflection    = "reflection"
	EntityImprovement   = "improvement"
)

// Link relation values.
const (
	RelDecisionAffectsTask             = "decision_affects_task"
	RelDiscoveryCausesPlanChange       = "discovery_causes_plan_change"
	RelDiscoveryMentionsTask           = "discovery_mentions_task"
	RelClaimHasEvidence                = "claim_has_evidence"
	RelReviewJudgesTask                = "review_judges_task"
	RelReviewJudgesScope               = "review_judges_scope"
	RelGoalHasTaskEvent                = "goal_has_task" // event payload only; not entity_links
	RelUncertaintyBlocksTask           = "uncertainty_blocks_task"
	RelUncertaintyAffectsGoal          = "uncertainty_affects_goal"
	RelHypothesisSupportedBy           = "hypothesis_supported_by"
	RelHypothesisAddressesUncertainty  = "hypothesis_addresses_uncertainty"
	RelAssumptionSupportsDecision      = "assumption_supports_decision"
	RelAssumptionAffectsTask           = "assumption_affects_task"
	RelDiscoveryInvalidatesAssumption  = "discovery_invalidates_assumption"
	RelChangeImplementsDecision        = "change_implements_decision"
	RelEffectSupportedBy               = "effect_supported_by"
	RelHypothesisExplainsEffect        = "hypothesis_explains_effect"
	RelDiscoveryFromContradictedEffect = "discovery_from_contradicted_effect"
	RelOutcomeSupportedBy              = "outcome_supported_by"
	RelRegressionFromEvaluation        = "regression_from_evaluation"
	RelRegressionFromEffect            = "regression_from_effect"
	RelRegressionAssociatedChange      = "regression_associated_change"
	RelHypothesisExplainsRegression    = "hypothesis_explains_regression"
	RelRegressionSupportedBy           = "regression_supported_by"
	RelImprovementSupportedBy          = "improvement_supported_by"
	RelReflectionInvalidatesAssumption = "reflection_invalidates_assumption"
	RelObservedRelationship            = "observed_relationship"
	RelCausedBy                        = "caused_by"
	RelRelationshipSupportedBy         = "relationship_supported_by"
)

// Residual severity vocabulary (canonical; matches store.ResidualSeverity*).
const (
	ResidualSeverityINFO     = store.ResidualSeverityINFO
	ResidualSeverityWARN     = store.ResidualSeverityWARN
	ResidualSeverityBlocking = store.ResidualSeverityBlocking
)

// Residual status vocabulary.
const (
	ResidualStatusOpen     = store.ResidualStatusOpen
	ResidualStatusAcked    = store.ResidualStatusAcked
	ResidualStatusResolved = store.ResidualStatusResolved
)

// Recommended residual codes (string column; other non-empty codes OK).
const (
	ResidualCodeMissingEvidence = "MISSING_EVIDENCE"
	ResidualCodeOpenGap         = "OPEN_GAP"
	ResidualCodeContractGap     = "CONTRACT_GAP"
	ResidualCodePolicyException = "POLICY_EXCEPTION"
)

// Event type names (locked).
const (
	EventEntityCreated                = "entity.created"
	EventEntityLinked                 = "entity.linked"
	EventTaskTransition               = "task.transition"
	EventReviewResult                 = "review.result"
	EventDeliberationTransition       = "deliberation.transition"
	EventUncertaintyResolved          = "uncertainty.resolved"
	EventUncertaintySuperseded        = "uncertainty.superseded"
	EventHypothesisConfirmed          = "hypothesis.confirmed"
	EventHypothesisRejected           = "hypothesis.rejected"
	EventHypothesisSuperseded         = "hypothesis.superseded"
	EventAssumptionInvalidated        = "assumption.invalidated"
	EventDecisionReconsider           = "decision.reconsider"
	EventChangeRecorded               = "change.recorded"
	EventEffectCompared               = "effect.compared"
	EventEffectContradicted           = "effect.contradicted"
	EventOutcomeRecorded              = "outcome.recorded"
	EventBaselineCreated              = "baseline.created"
	EventBaselinePromoted             = "baseline.promoted"
	EventEvaluationCompared           = "evaluation.compared"
	EventRegressionRecorded           = "regression.recorded"
	EventRegressionAttributionChanged = "regression.attribution_changed"
	EventRegressionResolved           = "regression.resolved"
	EventReflectionRecorded           = "reflection.recorded"
	EventRelationshipObserved         = "relationship.observed"
	EventRelationshipCaused           = "relationship.caused"
)

// Discovery severity vocabulary (canonical; matches store.Severity*).
const (
	SeverityINFO          = store.SeverityINFO
	SeverityPlanAffecting = store.SeverityPlanAffecting
	SeverityBlocking      = store.SeverityBlocking
)

// DefaultSourceType when caller omits source_type.
const DefaultSourceType = "USER_ASSERTED"

// NormalizeSeverity returns a valid severity. Empty defaults to INFO; unknown values fail closed.
func NormalizeSeverity(severity string) (string, error) {
	s := strings.TrimSpace(severity)
	if s == "" {
		return SeverityINFO, nil
	}
	switch s {
	case SeverityINFO, SeverityPlanAffecting, SeverityBlocking:
		return s, nil
	default:
		return "", &ErrValidation{Msg: "severity must be INFO, PLAN_AFFECTING, or BLOCKING"}
	}
}

// Service is the work/causal API bound to a store.
type Service struct {
	store        *store.Store
	impactWalker ImpactWalker
}

// New constructs a domain Service. st must be non-nil and already opened.
func New(st *store.Store) *Service {
	if st == nil {
		panic("domain: New: store is nil")
	}
	return &Service{store: st}
}

// LinkMeta carries optional provenance for link operations.
type LinkMeta struct {
	SourceType string
	Confidence float64
}

func (m LinkMeta) withDefaults() LinkMeta {
	if m.SourceType == "" {
		m.SourceType = DefaultSourceType
	}
	return m
}

// TransitionOptions controls TransitionTask.
type TransitionOptions struct {
	Actor                    string
	Reason                   string
	EvidenceIDs              []string
	AllowDoneWithoutReview   bool // hatch: bypass PASS + operator (Gate G)
	AllowOperatorDone        bool // DF-17: required with linked PASS (Actor string ≠ auth)
	AllowMissingCapabilities bool // DF-24: override fail-closed missing-cap gate
}

// ErrInvalidTransition is returned for illegal work_state edges or DONE policy.
type ErrInvalidTransition struct {
	From, To string
	Reason   string
}

func (e *ErrInvalidTransition) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("domain: invalid transition %s → %s: %s", e.From, e.To, e.Reason)
	}
	return fmt.Sprintf("domain: invalid transition %s → %s", e.From, e.To)
}

// ErrValidation is returned for bad inputs (empty title, missing actor, etc.).
type ErrValidation struct {
	Msg string
}

func (e *ErrValidation) Error() string {
	return "domain: " + e.Msg
}
