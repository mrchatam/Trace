package domain

import (
	"context"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// Impact class vocabulary (canonical; matches store.ImpactClass*).
const (
	ImpactClassSAFE        = store.ImpactClassSAFE
	ImpactClassCaution     = store.ImpactClassCaution
	ImpactClassHigh        = store.ImpactClassHigh
	ImpactClassDestructive = store.ImpactClassDestructive
	ImpactClassReversal    = store.ImpactClassReversal
)

// Uncertainty vocabulary.
const (
	UncertaintyKNOWN    = store.UncertaintyKNOWN
	UncertaintyLIKELY   = store.UncertaintyLIKELY
	UncertaintyPOSSIBLE = store.UncertaintyPOSSIBLE
	UncertaintyUNKNOWN  = store.UncertaintyUNKNOWN
)

// Finding kind vocabulary.
const (
	FindingKindAffectedWork          = store.FindingKindAffectedWork
	FindingKindInvalidatedAssumption = store.FindingKindInvalidatedAssumption
	FindingKindWorkAtRisk            = store.FindingKindWorkAtRisk
	FindingKindNewWork               = store.FindingKindNewWork
	FindingKindUnresolved            = store.FindingKindUnresolved
)

// ImpactFindingInput creates a planted/manual impact finding on a decision.
type ImpactFindingInput struct {
	ID          string
	ImpactClass string
	Uncertainty string // empty → UNKNOWN
	Kind        string
	Body        string
	RelatedType string
	RelatedID   string
}

// AlternativeInput creates a thin alternative on a decision.
type AlternativeInput struct {
	ID          string
	Title       string
	Body        string
	Recommended bool
}

// ImpactReportResult aggregates decision + decision_affects_task links + findings + alternatives.
// Fail-closed: never mutates plan/tasks/decisions; surfaces HasUnknown/Incomplete.
type ImpactReportResult struct {
	Decision           store.Decision
	AffectedTaskIDs    []string
	AffectsTaskLinks   []store.EntityLink
	Findings           []store.DecisionImpactFinding
	Alternatives       []store.DecisionAlternative
	OverallClass       string
	OverallUncertainty string
	HasUnknown         bool
	Incomplete         bool
}

// NormalizeImpactClass returns a valid impact class. Empty and unknown fail closed.
func NormalizeImpactClass(class string) (string, error) {
	c := strings.ToUpper(strings.TrimSpace(class))
	if c == "" {
		return "", &ErrValidation{Msg: "impact_class is required (SAFE|CAUTION|HIGH|DESTRUCTIVE|REVERSAL)"}
	}
	switch c {
	case ImpactClassSAFE, ImpactClassCaution, ImpactClassHigh, ImpactClassDestructive, ImpactClassReversal:
		return c, nil
	default:
		return "", &ErrValidation{Msg: "impact_class must be SAFE, CAUTION, HIGH, DESTRUCTIVE, or REVERSAL"}
	}
}

// NormalizeUncertainty returns a valid uncertainty. Empty → UNKNOWN; unknown values fail closed.
func NormalizeUncertainty(u string) (string, error) {
	s := strings.ToUpper(strings.TrimSpace(u))
	if s == "" {
		return UncertaintyUNKNOWN, nil
	}
	switch s {
	case UncertaintyKNOWN, UncertaintyLIKELY, UncertaintyPOSSIBLE, UncertaintyUNKNOWN:
		return s, nil
	default:
		return "", &ErrValidation{Msg: "uncertainty must be KNOWN, LIKELY, POSSIBLE, or UNKNOWN"}
	}
}

// NormalizeFindingKind returns a valid finding kind. Empty and unknown fail closed.
func NormalizeFindingKind(kind string) (string, error) {
	k := strings.ToUpper(strings.TrimSpace(kind))
	if k == "" {
		return "", &ErrValidation{Msg: "finding kind is required (AFFECTED_WORK|INVALIDATED_ASSUMPTION|WORK_AT_RISK|NEW_WORK|UNRESOLVED)"}
	}
	switch k {
	case FindingKindAffectedWork, FindingKindInvalidatedAssumption, FindingKindWorkAtRisk,
		FindingKindNewWork, FindingKindUnresolved:
		return k, nil
	default:
		return "", &ErrValidation{Msg: "finding kind must be AFFECTED_WORK, INVALIDATED_ASSUMPTION, WORK_AT_RISK, NEW_WORK, or UNRESOLVED"}
	}
}

var impactClassRank = map[string]int{
	ImpactClassSAFE:        1,
	ImpactClassCaution:     2,
	ImpactClassHigh:        3,
	ImpactClassDestructive: 4,
	ImpactClassReversal:    5,
}

var uncertaintyRank = map[string]int{
	UncertaintyKNOWN:    1,
	UncertaintyLIKELY:   2,
	UncertaintyPOSSIBLE: 3,
	UncertaintyUNKNOWN:  4,
}

// AddImpactFinding records a finding on an existing decision (no auto links / plan mutation).
func (s *Service) AddImpactFinding(ctx context.Context, decisionID string, in ImpactFindingInput) (store.DecisionImpactFinding, error) {
	_ = ctx
	if decisionID == "" {
		return store.DecisionImpactFinding{}, &ErrValidation{Msg: "decisionID is required"}
	}
	class, err := NormalizeImpactClass(in.ImpactClass)
	if err != nil {
		return store.DecisionImpactFinding{}, err
	}
	kind, err := NormalizeFindingKind(in.Kind)
	if err != nil {
		return store.DecisionImpactFinding{}, err
	}
	unc, err := NormalizeUncertainty(in.Uncertainty)
	if err != nil {
		return store.DecisionImpactFinding{}, err
	}
	if _, err := s.store.GetDecision(decisionID); err != nil {
		return store.DecisionImpactFinding{}, err
	}
	return s.store.InsertDecisionImpactFinding(store.DecisionImpactFinding{
		ID:          in.ID,
		DecisionID:  decisionID,
		ImpactClass: class,
		Uncertainty: unc,
		Kind:        kind,
		Body:        in.Body,
		RelatedType: strings.TrimSpace(in.RelatedType),
		RelatedID:   strings.TrimSpace(in.RelatedID),
	})
}

// ListImpactFindings returns findings for a decision.
func (s *Service) ListImpactFindings(ctx context.Context, decisionID string) ([]store.DecisionImpactFinding, error) {
	_ = ctx
	if decisionID == "" {
		return nil, &ErrValidation{Msg: "decisionID is required"}
	}
	return s.store.ListDecisionImpactFindingsByDecisionID(decisionID)
}

// AddDecisionAlternative records an alternative. If Recommended, clears siblings then sets this one.
func (s *Service) AddDecisionAlternative(ctx context.Context, decisionID string, in AlternativeInput) (store.DecisionAlternative, error) {
	_ = ctx
	if decisionID == "" {
		return store.DecisionAlternative{}, &ErrValidation{Msg: "decisionID is required"}
	}
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return store.DecisionAlternative{}, &ErrValidation{Msg: "alternative title is required"}
	}
	if _, err := s.store.GetDecision(decisionID); err != nil {
		return store.DecisionAlternative{}, err
	}
	if in.Recommended {
		if err := s.store.ClearRecommendedAlternatives(decisionID); err != nil {
			return store.DecisionAlternative{}, err
		}
	}
	return s.store.InsertDecisionAlternative(store.DecisionAlternative{
		ID:            in.ID,
		DecisionID:    decisionID,
		Title:         title,
		Body:          in.Body,
		IsRecommended: in.Recommended,
	})
}

// SetRecommendedAlternative clears siblings then marks the given alternative recommended.
func (s *Service) SetRecommendedAlternative(ctx context.Context, decisionID, alternativeID string) error {
	_ = ctx
	if decisionID == "" || alternativeID == "" {
		return &ErrValidation{Msg: "decisionID and alternativeID are required"}
	}
	if _, err := s.store.GetDecision(decisionID); err != nil {
		return err
	}
	alt, err := s.store.GetDecisionAlternative(alternativeID)
	if err != nil {
		return err
	}
	if alt.DecisionID != decisionID {
		return &ErrValidation{Msg: "alternative does not belong to decision"}
	}
	if err := s.store.ClearRecommendedAlternatives(decisionID); err != nil {
		return err
	}
	return s.store.UpdateDecisionAlternativeRecommended(alternativeID, true)
}

// ListDecisionAlternatives returns alternatives for a decision.
func (s *Service) ListDecisionAlternatives(ctx context.Context, decisionID string) ([]store.DecisionAlternative, error) {
	_ = ctx
	if decisionID == "" {
		return nil, &ErrValidation{Msg: "decisionID is required"}
	}
	return s.store.ListDecisionAlternativesByDecisionID(decisionID)
}

// ImpactReport aggregates decision + decision_affects_task links + findings + alternatives.
// Fail-closed; never mutates plan / tasks / decisions.status.
func (s *Service) ImpactReport(ctx context.Context, decisionID string) (ImpactReportResult, error) {
	_ = ctx
	if decisionID == "" {
		return ImpactReportResult{}, &ErrValidation{Msg: "decisionID is required"}
	}
	dec, err := s.store.GetDecision(decisionID)
	if err != nil {
		return ImpactReportResult{}, err
	}
	allLinks, err := s.store.ListLinksFrom(EntityDecision, decisionID)
	if err != nil {
		return ImpactReportResult{}, err
	}
	var affects []store.EntityLink
	var taskIDs []string
	for _, l := range allLinks {
		if l.Rel != RelDecisionAffectsTask {
			continue
		}
		affects = append(affects, l)
		taskIDs = append(taskIDs, l.ToID)
	}
	findings, err := s.store.ListDecisionImpactFindingsByDecisionID(decisionID)
	if err != nil {
		return ImpactReportResult{}, err
	}
	if findings == nil {
		findings = []store.DecisionImpactFinding{}
	}
	alts, err := s.store.ListDecisionAlternativesByDecisionID(decisionID)
	if err != nil {
		return ImpactReportResult{}, err
	}
	if alts == nil {
		alts = []store.DecisionAlternative{}
	}

	hasUnknownFinding := false
	overallClass := ""
	overallUnc := UncertaintyUNKNOWN
	maxClass := 0
	maxUnc := 0
	for _, f := range findings {
		if f.Uncertainty == UncertaintyUNKNOWN {
			hasUnknownFinding = true
		}
		if r := impactClassRank[f.ImpactClass]; r > maxClass {
			maxClass = r
			overallClass = f.ImpactClass
		}
		if r := uncertaintyRank[f.Uncertainty]; r > maxUnc {
			maxUnc = r
			overallUnc = f.Uncertainty
		}
	}
	if len(findings) == 0 {
		overallClass = ""
		overallUnc = UncertaintyUNKNOWN
	}

	hasUnknown := hasUnknownFinding || (len(findings) == 0 && len(affects) >= 1)
	hasRecommended := false
	for _, a := range alts {
		if a.IsRecommended {
			hasRecommended = true
			break
		}
	}
	incomplete := hasUnknown ||
		(len(alts) > 0 && !hasRecommended) ||
		(len(findings) == 0 && len(affects) == 0)

	if taskIDs == nil {
		taskIDs = []string{}
	}
	if affects == nil {
		affects = []store.EntityLink{}
	}

	return ImpactReportResult{
		Decision:           dec,
		AffectedTaskIDs:    taskIDs,
		AffectsTaskLinks:   affects,
		Findings:           findings,
		Alternatives:       alts,
		OverallClass:       overallClass,
		OverallUncertainty: overallUnc,
		HasUnknown:         hasUnknown,
		Incomplete:         incomplete,
	}, nil
}

// DecisionImpact is the packet/why JSON shape for one decision with findings.
type DecisionImpact struct {
	DecisionID         string                        `json:"decision_id"`
	OverallClass       string                        `json:"overall_class"`
	OverallUncertainty string                        `json:"overall_uncertainty,omitempty"`
	HasUnknown         bool                          `json:"has_unknown"`
	Incomplete         bool                          `json:"incomplete"`
	Findings           []store.DecisionImpactFinding `json:"findings"`
}

func decisionImpactFromReport(rep ImpactReportResult) DecisionImpact {
	findings := rep.Findings
	if findings == nil {
		findings = []store.DecisionImpactFinding{}
	}
	return DecisionImpact{
		DecisionID:         rep.Decision.ID,
		OverallClass:       rep.OverallClass,
		OverallUncertainty: rep.OverallUncertainty,
		HasUnknown:         rep.HasUnknown,
		Incomplete:         rep.Incomplete,
		Findings:           findings,
	}
}

// ImpactSummariesForTask returns ImpactReport summaries for decisions linked
// decision_affects_task to the task, omitting decisions with zero findings.
func (s *Service) ImpactSummariesForTask(ctx context.Context, taskID string) ([]DecisionImpact, error) {
	if taskID == "" {
		return nil, &ErrValidation{Msg: "taskID is required"}
	}
	links, err := s.store.ListLinksTo(EntityTask, taskID)
	if err != nil {
		return nil, err
	}
	var out []DecisionImpact
	for _, l := range links {
		if l.Rel != RelDecisionAffectsTask {
			continue
		}
		rep, err := s.ImpactReport(ctx, l.FromID)
		if err != nil {
			return nil, err
		}
		if len(rep.Findings) == 0 {
			continue
		}
		out = append(out, decisionImpactFromReport(rep))
	}
	return out, nil
}

// ImpactSummariesForWhySeed returns packet-shaped impact for why JSON inherit.
// Task seeds use neighborhood reports; decision seeds include that report iff findings exist.
func (s *Service) ImpactSummariesForWhySeed(ctx context.Context, entityType, entityID string) ([]DecisionImpact, error) {
	switch entityType {
	case EntityTask:
		return s.ImpactSummariesForTask(ctx, entityID)
	case EntityDecision:
		rep, err := s.ImpactReport(ctx, entityID)
		if err != nil {
			return nil, err
		}
		if len(rep.Findings) == 0 {
			return nil, nil
		}
		return []DecisionImpact{decisionImpactFromReport(rep)}, nil
	case EntityRegression:
		return s.impactSummaryForRegressionWhySeed(entityID)
	case EntityUncertainty:
		return nil, nil
	default:
		return nil, nil
	}
}

func (s *Service) impactSummaryForRegressionWhySeed(regressionID string) ([]DecisionImpact, error) {
	reg, err := s.store.GetRegression(regressionID)
	if err != nil {
		return nil, err
	}
	links, err := s.store.ListLinksFrom(EntityRegression, regressionID)
	if err != nil {
		return nil, err
	}
	var sourceType, sourceID string
	for _, l := range links {
		switch l.Rel {
		case RelRegressionFromEvaluation:
			sourceType, sourceID = EntityOutcomeResult, l.ToID
		case RelRegressionFromEffect:
			sourceType, sourceID = EntityEffect, l.ToID
		}
		if sourceType != "" {
			break
		}
	}
	if sourceType == "" {
		return nil, nil
	}
	summary := reg.Summary
	if summary == "" {
		summary = reg.Dimension
	}
	return []DecisionImpact{{
		DecisionID:   sourceID,
		OverallClass: sourceType,
		Findings: []store.DecisionImpactFinding{{
			Kind:        "REGRESSION_SOURCE",
			Body:        summary,
			RelatedType: sourceType,
			RelatedID:   sourceID,
		}},
	}}, nil
}
