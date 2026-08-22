package domain

import (
	"context"
	"os"
	"path/filepath"

	"github.com/mrchatam/Trace/internal/gitcli"
	"github.com/mrchatam/Trace/internal/store"
)

// SeedDocument is seed JSON v1 for export/import round-trip.
type SeedDocument struct {
	Version                  int                           `json:"version"`
	Goals                    []SeedEntity                  `json:"goals"`
	Tasks                    []SeedTask                    `json:"tasks"`
	Decisions                []SeedEntity                  `json:"decisions"`
	Assumptions              []SeedEntity                  `json:"assumptions"`
	Discoveries              []SeedEntity                  `json:"discoveries"`
	PlanChanges              []SeedEntity                  `json:"plan_changes"`
	Claims                   []SeedEntity                  `json:"claims"`
	Evidence                 []SeedEntity                  `json:"evidence"`
	Links                    []SeedLink                    `json:"links"`
	Findings                 []SeedFinding                 `json:"findings"`
	Alternatives             []SeedAlternative             `json:"alternatives"`
	PlanPhases               []SeedPlanPhase               `json:"plan_phases"`
	PlanScopes               []SeedPlanScope               `json:"plan_scopes"`
	ScopeDeepPlans           []SeedScopeDeepPlan           `json:"scope_deep_plans"`
	GoalPlanState            []SeedGoalPlanState           `json:"goal_plan_state"`
	DeliberationStates       []SeedDeliberationState       `json:"deliberation_states,omitempty"`
	Uncertainties            []SeedUncertainty             `json:"uncertainties,omitempty"`
	Hypotheses               []SeedHypothesis              `json:"hypotheses,omitempty"`
	DecisionReconsiderations []SeedDecisionReconsideration `json:"decision_reconsiderations,omitempty"`
	Changes                  []SeedChange                  `json:"changes,omitempty"`
	Effects                  []SeedEffect                  `json:"effects,omitempty"`
	OutcomeResults           []SeedOutcomeResult           `json:"outcome_results,omitempty"`
	Baselines                []SeedBaseline                `json:"baselines,omitempty"`
	Regressions              []SeedRegression              `json:"regressions,omitempty"`
	Improvements             []SeedImprovement             `json:"improvements,omitempty"`
	Reflections              []SeedReflection              `json:"reflections,omitempty"`
	ChangePatterns           []SeedChangePattern           `json:"change_patterns,omitempty"`
	EngineeringKnowledge     []SeedEngineeringKnowledge    `json:"engineering_knowledge,omitempty"`
	HarnessAgents            []SeedHarnessAgent            `json:"harness_agents,omitempty"`
	EvalRulesPath            string                        `json:"eval_rules_path,omitempty"`
	ExportedAtCommit         string                        `json:"exported_at_commit,omitempty"`
	Transitions              []SeedTransition              `json:"transitions,omitempty"`
}

type SeedEntity struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// SeedTask is portable identity (id, title, body, goal_id). Default export
// omits work_state so clone import lands PENDING (DF-88; keep TestSeedExportOmitsDeniedSurfaces).
type SeedTask struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	GoalID string `json:"goal_id"`
}

type SeedLink struct {
	Rel    string `json:"rel"`
	From   string `json:"from"`
	To     string `json:"to"`
	FromID string `json:"from_id,omitempty"`
	ToID   string `json:"to_id,omitempty"`
}

type SeedFinding struct {
	ID          string `json:"id"`
	DecisionID  string `json:"decision_id"`
	ImpactClass string `json:"impact_class"`
	Kind        string `json:"kind"`
	Uncertainty string `json:"uncertainty,omitempty"`
	Body        string `json:"body"`
	RelatedType string `json:"related_type,omitempty"`
	RelatedID   string `json:"related_id,omitempty"`
}

type SeedAlternative struct {
	ID          string `json:"id"`
	DecisionID  string `json:"decision_id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Recommended bool   `json:"recommended"`
}

type SeedPlanPhase struct {
	ID     string `json:"id"`
	GoalID string `json:"goal_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Ord    int    `json:"ord"`
	Status string `json:"status"`
}

type SeedPlanScope struct {
	ID              string `json:"id"`
	PhaseID         string `json:"phase_id"`
	Title           string `json:"title"`
	Body            string `json:"body"`
	Ord             int    `json:"ord"`
	Status          string `json:"status"`
	AutoReplanCount int    `json:"auto_replan_count"`
}

type SeedScopeDeepPlan struct {
	ID          string `json:"id"`
	ScopeID     string `json:"scope_id"`
	ContentJSON string `json:"content_json"`
	Status      string `json:"status"`
}

type SeedGoalPlanState struct {
	GoalID         string  `json:"goal_id"`
	CurrentScopeID *string `json:"current_scope_id"`
}

type SeedTransition struct {
	TaskID           string `json:"task_id"`
	To               string `json:"to"`
	Actor            string `json:"actor"`
	Reason           string `json:"reason"`
	AllowDone        bool   `json:"allow_done"`
	AsOperator       bool   `json:"as_operator"`
	AllowMissingCaps bool   `json:"allow_missing_caps"`
}

type SeedDeliberationState struct {
	TaskID                  string `json:"task_id"`
	GoalID                  string `json:"goal_id"`
	CurrentPhase            string `json:"current_phase"`
	HopCount                int    `json:"hop_count"`
	LastPhase               string `json:"last_phase"`
	PlanCritiqued           bool   `json:"plan_critiqued"`
	Stopped                 bool   `json:"stopped"`
	StopReason              string `json:"stop_reason"`
	ConsecutiveEmptyApplies int    `json:"consecutive_empty_applies,omitempty"`
	UpdatedAt               string `json:"updated_at"`
}

type SeedUncertainty struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Body           string  `json:"body"`
	Severity       string  `json:"severity"`
	Status         string  `json:"status"`
	Kind           string  `json:"kind"`
	Confidence     float64 `json:"confidence"`
	SourceType     string  `json:"source_type"`
	Resolution     string  `json:"resolution"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	LastVerifiedAt *string `json:"last_verified_at,omitempty"`
}

type SeedHypothesis struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Body           string  `json:"body"`
	Status         string  `json:"status"`
	Confidence     float64 `json:"confidence"`
	SourceType     string  `json:"source_type"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	LastVerifiedAt *string `json:"last_verified_at,omitempty"`
}

type SeedDecisionReconsideration struct {
	ID           string `json:"id"`
	DecisionID   string `json:"decision_id"`
	Trigger      string `json:"trigger"`
	Status       string `json:"status"`
	Reason       string `json:"reason"`
	RelatedType  string `json:"related_type"`
	RelatedID    string `json:"related_id"`
	ReconsiderAt string `json:"reconsider_at"`
	CreatedAt    string `json:"created_at"`
}

type SeedChangePath struct {
	Path     string `json:"path"`
	Status   string `json:"status"`
	SymbolID string `json:"symbol_id"`
}

type SeedChange struct {
	ID             string           `json:"id"`
	TaskID         string           `json:"task_id"`
	GitCommit      string           `json:"git_commit"`
	ParentChangeID string           `json:"parent_change_id"`
	Actor          string           `json:"actor"`
	Reason         string           `json:"reason"`
	Status         string           `json:"status"`
	SourceType     string           `json:"source_type"`
	Confidence     float64          `json:"confidence"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
	LastVerifiedAt *string          `json:"last_verified_at,omitempty"`
	Paths          []SeedChangePath `json:"paths,omitempty"`
}

type SeedEffect struct {
	ID         string  `json:"id"`
	ChangeID   string  `json:"change_id"`
	Dimension  string  `json:"dimension"`
	Expected   string  `json:"expected"`
	Actual     string  `json:"actual"`
	Comparison string  `json:"comparison"`
	Confidence float64 `json:"confidence"`
	SourceType string  `json:"source_type"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type SeedBaseline struct {
	ID           string `json:"id"`
	GitCommit    string `json:"git_commit"`
	ScoresJSON   string `json:"scores_json"`
	Label        string `json:"label"`
	SourceType   string `json:"source_type"`
	Status       string `json:"status,omitempty"`
	SupersedesID string `json:"supersedes_id,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type SeedOutcomeResult struct {
	ID                 string  `json:"id"`
	TaskID             string  `json:"task_id"`
	Kind               string  `json:"kind"`
	TestName           string  `json:"test_name"`
	TestStatus         string  `json:"test_status"`
	GoalID             string  `json:"goal_id"`
	VerificationStatus string  `json:"verification_status"`
	BaselineID         string  `json:"baseline_id"`
	ScoresJSON         string  `json:"scores_json"`
	ComparisonJSON     string  `json:"comparison_json"`
	Summary            string  `json:"summary"`
	Actor              string  `json:"actor"`
	SourceType         string  `json:"source_type"`
	Confidence         float64 `json:"confidence"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type SeedRegression struct {
	ID          string  `json:"id"`
	TaskID      string  `json:"task_id"`
	SourceKind  string  `json:"source_kind"`
	SourceID    string  `json:"source_id"`
	Dimension   string  `json:"dimension"`
	Attribution string  `json:"attribution"`
	Status      string  `json:"status"`
	Summary     string  `json:"summary"`
	Actor       string  `json:"actor"`
	SourceType  string  `json:"source_type"`
	Confidence  float64 `json:"confidence"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type SeedImprovement struct {
	ID              string  `json:"id"`
	ChangeID        string  `json:"change_id"`
	TaskID          string  `json:"task_id"`
	Dimension       string  `json:"dimension"`
	Summary         string  `json:"summary"`
	EvidenceIDsJSON string  `json:"evidence_ids_json"`
	SourceType      string  `json:"source_type"`
	Confidence      float64 `json:"confidence"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type SeedReflection struct {
	ID                         string  `json:"id"`
	TaskID                     string  `json:"task_id"`
	Summary                    string  `json:"summary"`
	InvalidatedAssumptionsJSON string  `json:"invalidated_assumptions_json"`
	NewDependenciesJSON        string  `json:"new_dependencies_json"`
	UsefulTestsJSON            string  `json:"useful_tests_json"`
	BroadenTestsNote           string  `json:"broaden_tests_note"`
	Actor                      string  `json:"actor"`
	SourceType                 string  `json:"source_type"`
	Confidence                 float64 `json:"confidence"`
	CreatedAt                  string  `json:"created_at"`
	UpdatedAt                  string  `json:"updated_at"`
}

type SeedChangePattern struct {
	ChangeKind    string `json:"change_kind"`
	OutcomeKind   string `json:"outcome_kind"`
	CountPositive int    `json:"count_positive"`
	CountNegative int    `json:"count_negative"`
	LastSeen      string `json:"last_seen"`
}

type SeedEngineeringKnowledge struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	BodyJSON        string  `json:"body_json"`
	Topic           string  `json:"topic"`
	EvidenceIDsJSON string  `json:"evidence_ids_json"`
	Confidence      float64 `json:"confidence"`
	Status          string  `json:"status"`
	SourceType      string  `json:"source_type"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type SeedHarnessAgent struct {
	ID                 string                        `json:"id"`
	Slug               string                        `json:"slug"`
	Title              string                        `json:"title"`
	Description        string                        `json:"description"`
	SubagentType       string                        `json:"subagent_type"`
	DeliberationPhases string                        `json:"deliberation_phases"`
	TaskKeywords       string                        `json:"task_keywords"`
	RecommendSubagent  bool                          `json:"recommend_subagent"`
	RegistrySource     string                        `json:"registry_source"`
	RegistryVersion    string                        `json:"registry_version"`
	ExternalURL        string                        `json:"external_url,omitempty"`
	CreatedAt          string                        `json:"created_at"`
	UpdatedAt          string                        `json:"updated_at"`
	Requirements       []SeedHarnessAgentRequirement `json:"requirements,omitempty"`
}

type SeedHarnessAgentRequirement struct {
	RequiredCapabilitySlug string `json:"required_capability_slug"`
}

// ExportOpts controls seed export behavior.
type ExportOpts struct {
	ProjectRoot string
}

var seedExportLinkRels = []string{
	RelDecisionAffectsTask,
	RelDiscoveryCausesPlanChange,
	RelClaimHasEvidence,
	RelDiscoveryMentionsTask,
}

// BuildSeedDocument assembles seed JSON v1 from the store (causal entities, links, plan tree, impact).
// Default export omits reviews, transitions, and task work_state (DF-88). Clone PENDING is expected.
func BuildSeedDocument(ctx context.Context, st *store.Store, opts ExportOpts) (SeedDocument, error) {
	if st == nil {
		return SeedDocument{}, &ErrValidation{Msg: "store is nil"}
	}

	doc := SeedDocument{
		Version:        1,
		Goals:          []SeedEntity{},
		Tasks:          []SeedTask{},
		Decisions:      []SeedEntity{},
		Assumptions:    []SeedEntity{},
		Discoveries:    []SeedEntity{},
		PlanChanges:    []SeedEntity{},
		Claims:         []SeedEntity{},
		Evidence:       []SeedEntity{},
		Links:          []SeedLink{},
		Findings:       []SeedFinding{},
		Alternatives:   []SeedAlternative{},
		PlanPhases:     []SeedPlanPhase{},
		PlanScopes:     []SeedPlanScope{},
		ScopeDeepPlans: []SeedScopeDeepPlan{},
		GoalPlanState:  []SeedGoalPlanState{},
	}

	goals, err := st.ListGoals()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, g := range goals {
		doc.Goals = append(doc.Goals, SeedEntity{ID: g.ID, Title: g.Title, Body: g.Body})
	}

	tasks, err := st.ListTasks()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, t := range tasks {
		st := SeedTask{ID: t.ID, Title: t.Title, Body: t.Body}
		if t.GoalID != nil {
			st.GoalID = *t.GoalID
		}
		doc.Tasks = append(doc.Tasks, st)
		if st.GoalID != "" {
			doc.Links = append(doc.Links, SeedLink{
				Rel:  RelGoalHasTaskEvent,
				From: st.GoalID,
				To:   st.ID,
			})
		}
	}

	decisions, err := st.ListDecisions()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, d := range decisions {
		doc.Decisions = append(doc.Decisions, SeedEntity{ID: d.ID, Title: d.Title, Body: d.Body})
	}

	assumptions, err := st.ListAssumptions()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, a := range assumptions {
		doc.Assumptions = append(doc.Assumptions, SeedEntity{ID: a.ID, Title: a.Title, Body: a.Body})
	}

	discoveries, err := st.ListDiscoveries()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, d := range discoveries {
		doc.Discoveries = append(doc.Discoveries, SeedEntity{ID: d.ID, Title: d.Title, Body: d.Body})
	}

	planChanges, err := st.ListPlanChanges()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, p := range planChanges {
		doc.PlanChanges = append(doc.PlanChanges, SeedEntity{ID: p.ID, Title: p.Title, Body: p.Body})
	}

	claims, err := st.ListClaims()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, c := range claims {
		doc.Claims = append(doc.Claims, SeedEntity{ID: c.ID, Title: c.Title, Body: c.Body})
	}

	evidence, err := st.ListEvidence()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, e := range evidence {
		doc.Evidence = append(doc.Evidence, SeedEntity{ID: e.ID, Title: e.Title, Body: e.Body})
	}

	for _, rel := range seedExportLinkRels {
		links, err := st.ListLinksByRel(rel)
		if err != nil {
			return SeedDocument{}, err
		}
		for _, l := range links {
			doc.Links = append(doc.Links, SeedLink{Rel: l.Rel, From: l.FromID, To: l.ToID})
		}
	}

	findings, err := st.ListAllDecisionImpactFindings()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, f := range findings {
		doc.Findings = append(doc.Findings, SeedFinding{
			ID: f.ID, DecisionID: f.DecisionID, ImpactClass: f.ImpactClass,
			Kind: f.Kind, Uncertainty: f.Uncertainty, Body: f.Body,
			RelatedType: f.RelatedType, RelatedID: f.RelatedID,
		})
	}

	alts, err := st.ListAllDecisionAlternatives()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, a := range alts {
		doc.Alternatives = append(doc.Alternatives, SeedAlternative{
			ID: a.ID, DecisionID: a.DecisionID, Title: a.Title, Body: a.Body,
			Recommended: a.IsRecommended,
		})
	}

	phases, err := st.ListAllPlanPhases()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, p := range phases {
		doc.PlanPhases = append(doc.PlanPhases, SeedPlanPhase{
			ID: p.ID, GoalID: p.GoalID, Title: p.Title, Body: p.Body, Ord: p.Ord, Status: p.Status,
		})
	}

	scopes, err := st.ListAllPlanScopes()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, sc := range scopes {
		doc.PlanScopes = append(doc.PlanScopes, SeedPlanScope{
			ID: sc.ID, PhaseID: sc.PhaseID, Title: sc.Title, Body: sc.Body,
			Ord: sc.Ord, Status: sc.Status, AutoReplanCount: sc.AutoReplanCount,
		})
	}

	deepPlans, err := st.ListAllScopeDeepPlans()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, d := range deepPlans {
		doc.ScopeDeepPlans = append(doc.ScopeDeepPlans, SeedScopeDeepPlan{
			ID: d.ID, ScopeID: d.ScopeID, ContentJSON: d.ContentJSON, Status: d.Status,
		})
	}

	planStates, err := st.ListAllGoalPlanStates()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, st := range planStates {
		doc.GoalPlanState = append(doc.GoalPlanState, SeedGoalPlanState{
			GoalID: st.GoalID, CurrentScopeID: st.CurrentScopeID,
		})
	}

	baselines, err := st.ListAllBaselines()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, b := range baselines {
		sb := SeedBaseline{
			ID: b.ID, GitCommit: b.GitCommit, ScoresJSON: b.ScoresJSON, Label: b.Label,
			SourceType: b.SourceType, CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt,
		}
		if b.Status != "" && b.Status != store.BaselineStatusActive {
			sb.Status = b.Status
		}
		if b.SupersedesID != "" {
			sb.SupersedesID = b.SupersedesID
		}
		doc.Baselines = append(doc.Baselines, sb)
	}

	outcomes, err := st.ListAllOutcomeResults()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, o := range outcomes {
		doc.OutcomeResults = append(doc.OutcomeResults, SeedOutcomeResult{
			ID: o.ID, TaskID: o.TaskID, Kind: o.Kind, TestName: o.TestName, TestStatus: o.TestStatus,
			GoalID: o.GoalID, VerificationStatus: o.VerificationStatus, BaselineID: o.BaselineID,
			ScoresJSON: o.ScoresJSON, ComparisonJSON: o.ComparisonJSON, Summary: o.Summary,
			Actor: o.Actor, SourceType: o.SourceType, Confidence: o.Confidence,
			CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt,
		})
	}

	changes, err := st.ListAllChanges()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, c := range changes {
		sc := SeedChange{
			ID: c.ID, TaskID: c.TaskID, GitCommit: c.GitCommit, ParentChangeID: c.ParentChangeID,
			Actor: c.Actor, Reason: c.Reason, Status: c.Status, SourceType: c.SourceType,
			Confidence: c.Confidence, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
			LastVerifiedAt: c.LastVerifiedAt,
		}
		paths, err := st.ListChangePaths(c.ID)
		if err != nil {
			return SeedDocument{}, err
		}
		for _, p := range paths {
			sc.Paths = append(sc.Paths, SeedChangePath{Path: p.Path, Status: p.Status, SymbolID: p.SymbolID})
		}
		doc.Changes = append(doc.Changes, sc)
	}

	effects, err := st.ListAllEffects()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, e := range effects {
		doc.Effects = append(doc.Effects, SeedEffect{
			ID: e.ID, ChangeID: e.ChangeID, Dimension: e.Dimension, Expected: e.Expected,
			Actual: e.Actual, Comparison: e.Comparison, Confidence: e.Confidence,
			SourceType: e.SourceType, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
		})
	}

	uncertainties, err := st.ListAllUncertainties()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, u := range uncertainties {
		doc.Uncertainties = append(doc.Uncertainties, SeedUncertainty{
			ID: u.ID, Title: u.Title, Body: u.Body, Severity: u.Severity, Status: u.Status,
			Kind: u.Kind, Confidence: u.Confidence, SourceType: u.SourceType, Resolution: u.Resolution,
			CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt, LastVerifiedAt: u.LastVerifiedAt,
		})
	}

	hypotheses, err := st.ListAllHypotheses()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, h := range hypotheses {
		doc.Hypotheses = append(doc.Hypotheses, SeedHypothesis{
			ID: h.ID, Title: h.Title, Body: h.Body, Status: h.Status, Confidence: h.Confidence,
			SourceType: h.SourceType, CreatedAt: h.CreatedAt, UpdatedAt: h.UpdatedAt,
			LastVerifiedAt: h.LastVerifiedAt,
		})
	}

	recons, err := st.ListAllDecisionReconsiderations()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, r := range recons {
		doc.DecisionReconsiderations = append(doc.DecisionReconsiderations, SeedDecisionReconsideration{
			ID: r.ID, DecisionID: r.DecisionID, Trigger: r.Trigger, Status: r.Status,
			Reason: r.Reason, RelatedType: r.RelatedType, RelatedID: r.RelatedID,
			ReconsiderAt: r.ReconsiderAt, CreatedAt: r.CreatedAt,
		})
	}

	regressions, err := st.ListAllRegressions()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, r := range regressions {
		doc.Regressions = append(doc.Regressions, SeedRegression{
			ID: r.ID, TaskID: r.TaskID, SourceKind: r.SourceKind, SourceID: r.SourceID,
			Dimension: r.Dimension, Attribution: r.Attribution, Status: r.Status, Summary: r.Summary,
			Actor: r.Actor, SourceType: r.SourceType, Confidence: r.Confidence,
			CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
		})
	}

	improvements, err := st.ListAllImprovements()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, im := range improvements {
		doc.Improvements = append(doc.Improvements, SeedImprovement{
			ID: im.ID, ChangeID: im.ChangeID, TaskID: im.TaskID, Dimension: im.Dimension,
			Summary: im.Summary, EvidenceIDsJSON: im.EvidenceIDsJSON, SourceType: im.SourceType,
			Confidence: im.Confidence, CreatedAt: im.CreatedAt, UpdatedAt: im.UpdatedAt,
		})
	}

	reflections, err := st.ListAllReflections()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, r := range reflections {
		doc.Reflections = append(doc.Reflections, SeedReflection{
			ID: r.ID, TaskID: r.TaskID, Summary: r.Summary,
			InvalidatedAssumptionsJSON: r.InvalidatedAssumptionsJSON,
			NewDependenciesJSON:        r.NewDependenciesJSON, UsefulTestsJSON: r.UsefulTestsJSON,
			BroadenTestsNote: r.BroadenTestsNote, Actor: r.Actor, SourceType: r.SourceType,
			Confidence: r.Confidence, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
		})
	}

	patterns, err := st.ListAllChangePatterns()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, p := range patterns {
		doc.ChangePatterns = append(doc.ChangePatterns, SeedChangePattern{
			ChangeKind: p.ChangeKind, OutcomeKind: p.OutcomeKind,
			CountPositive: p.CountPositive, CountNegative: p.CountNegative, LastSeen: p.LastSeen,
		})
	}

	knowledge, err := st.ListAllEngineeringKnowledge()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, k := range knowledge {
		doc.EngineeringKnowledge = append(doc.EngineeringKnowledge, SeedEngineeringKnowledge{
			ID: k.ID, Title: k.Title, BodyJSON: k.BodyJSON, Topic: k.Topic,
			EvidenceIDsJSON: k.EvidenceIDsJSON, Confidence: k.Confidence, Status: k.Status,
			SourceType: k.SourceType, CreatedAt: k.CreatedAt, UpdatedAt: k.UpdatedAt,
		})
	}

	harnessAgents, err := st.ListAllHarnessAgents()
	if err != nil {
		return SeedDocument{}, err
	}
	reqsByAgent := map[string][]SeedHarnessAgentRequirement{}
	allReqs, err := st.ListAllHarnessAgentRequirements()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, r := range allReqs {
		reqsByAgent[r.AgentID] = append(reqsByAgent[r.AgentID], SeedHarnessAgentRequirement{
			RequiredCapabilitySlug: r.RequiredCapabilitySlug,
		})
	}
	for _, a := range harnessAgents {
		doc.HarnessAgents = append(doc.HarnessAgents, SeedHarnessAgent{
			ID: a.ID, Slug: a.Slug, Title: a.Title, Description: a.Description,
			SubagentType: a.SubagentType, DeliberationPhases: a.DeliberationPhases,
			TaskKeywords: a.TaskKeywords, RecommendSubagent: a.RecommendSubagent,
			RegistrySource: a.RegistrySource, RegistryVersion: a.RegistryVersion,
			ExternalURL: a.ExternalURL, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
			Requirements: reqsByAgent[a.ID],
		})
	}

	delibStates, err := st.ListAllDeliberationStates()
	if err != nil {
		return SeedDocument{}, err
	}
	for _, ds := range delibStates {
		doc.DeliberationStates = append(doc.DeliberationStates, SeedDeliberationState{
			TaskID: ds.TaskID, GoalID: ds.GoalID, CurrentPhase: ds.CurrentPhase,
			HopCount: ds.HopCount, LastPhase: ds.LastPhase, PlanCritiqued: ds.PlanCritiqued,
			Stopped: ds.Stopped, StopReason: ds.StopReason,
			ConsecutiveEmptyApplies: ds.ConsecutiveEmptyApplies, UpdatedAt: ds.UpdatedAt,
		})
	}

	if opts.ProjectRoot != "" {
		if head, err := gitcli.HeadAtRoot(ctx, opts.ProjectRoot); err == nil && head != "" {
			doc.ExportedAtCommit = head
		}
		rulesPath := filepath.Join(opts.ProjectRoot, "trace", "eval-rules.json")
		if _, err := os.Stat(rulesPath); err == nil {
			doc.EvalRulesPath = "trace/eval-rules.json"
		}
	}

	return doc, nil
}
