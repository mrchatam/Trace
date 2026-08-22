package loop

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mrchatam/Trace/internal/agents"
	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	NextSchemaVersion = "trace.loop.next.v1"

	FreshnessFresh   = "fresh"
	FreshnessDirty   = "dirty"
	FreshnessStale   = "stale"
	FreshnessUnknown = "unknown"
)

type TaskSummary struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	WorkState string  `json:"work_state"`
	GoalID    *string `json:"goal_id"`
}

type PlanSnapshot struct {
	GoalID           string                    `json:"goal_id"`
	CurrentScopeID   *string                   `json:"current_scope_id"`
	Phases           []planner.PhaseView       `json:"phases"`
	CurrentDeepPlan  *planner.DeepPlanDocument `json:"current_deep_plan"`
	LookaheadScopeID string                    `json:"lookahead_scope_id"`
	LookaheadSummary string                    `json:"lookahead_summary"`
	Tasks            []TaskSummary             `json:"tasks"`
}

type WhySnapshot struct {
	retrieval.WhyResult
	Impact []domain.DecisionImpact `json:"impact,omitempty"`
}

type NextPacket struct {
	SchemaVersion           string                         `json:"schema_version"`
	GeneratedAt             time.Time                      `json:"generated_at"`
	Seed                    SeedSection                    `json:"seed"`
	Tasks                   TasksSection                   `json:"tasks"`
	Plan                    PlanSection                    `json:"plan"`
	Why                     WhySection                     `json:"why"`
	Context                 ContextSection                 `json:"context"`
	Related                 RelatedSection                 `json:"related"`
	LoopHints               LoopHintsSection               `json:"loop_hints"`
	PromotionCandidates     []PromotionCandidate           `json:"promotion_candidates"`
	Deliberation            DeliberationSection            `json:"deliberation"`
	OpenUncertainties       OpenUncertaintiesSection       `json:"open_uncertainties"`
	VerificationDebt        VerificationDebtSection        `json:"verification_debt"`
	RecentChanges           RecentChangesSection           `json:"recent_changes"`
	HistoricalRelationships HistoricalRelationshipsSection `json:"historical_relationships"`
	PlanningEvidence        PlanningEvidenceSection        `json:"planning_evidence"`
	Tendencies              TendenciesSection              `json:"tendencies"`
	SuccessfulApproaches    SuccessfulApproachesSection    `json:"successful_approaches"`
	SimilarChanges          SimilarChangesSection          `json:"similar_changes"`
	RiskHints               RiskHintsSection               `json:"risk_hints"`
	WorkConflicts           WorkConflictsSection           `json:"work_conflicts"`
	HarnessRecommendations  HarnessRecommendationsSection  `json:"harness_recommendations"`
}

// PromotionCandidate mirrors domain.PromotionCandidate for NextPacket JSON.
type PromotionCandidate = domain.PromotionCandidate

// HarnessRecommendationsSection surfaces ranked harness agent suggestions (UNTRUSTED_DATA).
type HarnessRecommendationsSection struct {
	Freshness string                  `json:"freshness"`
	Items     []agents.Recommendation `json:"items"`
}

// PlanningEvidenceSection surfaces task-scoped evaluations, reflections, and mixed planning signals.
type PlanningEvidenceSection struct {
	Freshness        string                          `json:"freshness"`
	Evaluations      []compiler.EvaluationItem       `json:"evaluations"`
	Reflections      []compiler.ReflectionItem       `json:"reflections"`
	PlanningEvidence []compiler.PlanningEvidenceItem `json:"planning_evidence"`
}

// TendenciesSection surfaces project-wide help/hurt tendencies from change patterns.
type TendenciesSection struct {
	Freshness string                  `json:"freshness"`
	Items     []compiler.TendencyItem `json:"items"`
}

// SuccessfulApproachesSection surfaces merged worked outcomes and knowledge rows.
type SuccessfulApproachesSection struct {
	Freshness string                            `json:"freshness"`
	Items     []compiler.SuccessfulApproachItem `json:"items"`
}

// SimilarChangesSection surfaces prior changes similar to the task's latest change kind.
type SimilarChangesSection struct {
	Freshness string                    `json:"freshness"`
	Items     []domain.SimilarChangeRow `json:"items"`
}

// WorkConflictsSection surfaces advisory overlaps with other active tasks.
type WorkConflictsSection struct {
	Freshness string                `json:"freshness"`
	Items     []domain.WorkConflict `json:"items"`
}

type SeedSection struct {
	Freshness string `json:"freshness"`
	TaskID    string `json:"task_id"`
	Title     string `json:"title"`
	WorkState string `json:"work_state"`
	GoalID    string `json:"goal_id"`
}

type TasksSection struct {
	Freshness string        `json:"freshness"`
	GoalID    string        `json:"goal_id"`
	Items     []TaskSummary `json:"items"`
}

type PlanSection struct {
	Freshness string       `json:"freshness"`
	Snapshot  PlanSnapshot `json:"snapshot"`
}

type WhySection struct {
	Freshness string      `json:"freshness"`
	Snapshot  WhySnapshot `json:"snapshot"`
}

type ContextSection struct {
	Freshness string          `json:"freshness"`
	Snapshot  compiler.Packet `json:"snapshot"`
}

type RelatedSection struct {
	Freshness         string                      `json:"freshness"`
	Available         bool                        `json:"available"`
	Seeds             []retrieval.ImpactSeed      `json:"seeds"`
	Snapshot          *retrieval.ImpactWalkResult `json:"snapshot,omitempty"`
	UnavailableReason string                      `json:"unavailable_reason,omitempty"`
}

type LoopHintsSection struct {
	Freshness         string `json:"freshness"`
	Available         bool   `json:"available"`
	UnavailableReason string `json:"unavailable_reason,omitempty"`
}

type BuildNextInput struct {
	TaskID    string
	Store     *store.Store
	Planner   *planner.Service
	Retrieval *retrieval.Engine
	Compiler  *compiler.Compiler
}

func BuildNextPacket(ctx context.Context, in BuildNextInput) (NextPacket, error) {
	if in.Store == nil {
		return NextPacket{}, fmt.Errorf("loop next: store is required")
	}
	if in.Planner == nil {
		return NextPacket{}, fmt.Errorf("loop next: planner is required")
	}
	if in.Retrieval == nil {
		return NextPacket{}, fmt.Errorf("loop next: retrieval is required")
	}
	if in.Compiler == nil {
		return NextPacket{}, fmt.Errorf("loop next: compiler is required")
	}
	if in.TaskID == "" {
		return NextPacket{}, fmt.Errorf("loop next: task id is required")
	}

	task, err := in.Store.GetTask(in.TaskID)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: load seed task: %w", err)
	}
	if task.GoalID == nil || *task.GoalID == "" {
		return NextPacket{}, fmt.Errorf("loop next: task %q has no goal_id", task.ID)
	}
	goalID := *task.GoalID

	tasks, err := in.Store.ListTasksByGoalID(goalID)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: list tasks for goal %q: %w", goalID, err)
	}
	taskRows := make([]TaskSummary, 0, len(tasks))
	for _, row := range tasks {
		taskRows = append(taskRows, TaskSummary{
			ID:        row.ID,
			Title:     row.Title,
			WorkState: row.WorkState,
			GoalID:    row.GoalID,
		})
	}

	planView, err := in.Planner.GetPlan(ctx, goalID)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: plan for goal %q: %w", goalID, err)
	}
	if planView.CurrentScopeID == nil || *planView.CurrentScopeID == "" || planView.CurrentDeepPlan == nil {
		return NextPacket{}, fmt.Errorf("loop next: missing goal plan context for goal %q", goalID)
	}
	lookaheadSummary := planView.LookaheadSummary
	planSnapshot := PlanSnapshot{
		GoalID:           planView.GoalID,
		CurrentScopeID:   planView.CurrentScopeID,
		Phases:           planView.Phases,
		CurrentDeepPlan:  planView.CurrentDeepPlan,
		LookaheadScopeID: planView.LookaheadScopeID,
		LookaheadSummary: lookaheadSummary,
		Tasks:            taskRows,
	}

	dom := domain.New(in.Store)
	seed := ApplySeed{TaskID: task.ID, GoalID: goalID}
	promotionCandidates, err := buildPromotionCandidates(in.Store)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: promotion candidates: %w", err)
	}
	p19Sat := p19SaturatedFromLastStep(in.Store, seed)
	inputs, err := BuildPolicyInputs(ctx, dom, in.Planner, task.ID, goalID, nil, p19Sat)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: policy inputs: %w", err)
	}

	dState := loadDeliberationState(ctx, dom, task.ID, goalID)
	recommended, _, _ := deliberation.SelectNext(dState, inputs)
	profile := PhaseContextProfile(recommended)
	if profile.TrimLookahead {
		planSnapshot.LookaheadSummary = trimLookaheadSummary(planSnapshot.LookaheadSummary, investigateLookaheadMax)
	}

	ctxOpts := compiler.ContextOptions{
		IncludeWhy:      profile.ContextIncludeWhy,
		IncludeMarkdown: profile.ContextIncludeMD,
		MaxItems:        profile.ContextMaxItems,
	}
	contextPacket, err := in.Compiler.TaskContext(ctx, task.ID, ctxOpts)
	if err != nil {
		if inputs.BlockingUncertaintyCount > 0 || inputs.OpenRegression {
			contextPacket, err = minimalTaskContextPacket(in.Store, task, ctxOpts.MaxItems)
		}
	}
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: context for task %q: %w", task.ID, err)
	}

	whyResult, err := in.Retrieval.Why(ctx, "task", task.ID)
	if err != nil && (inputs.BlockingUncertaintyCount > 0 || inputs.OpenRegression) {
		whyResult = retrieval.WhyResult{SeedType: "task", SeedID: task.ID}
		err = nil
	}
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: why for task %q: %w", task.ID, err)
	}
	impact, err := dom.ImpactSummariesForWhySeed(ctx, "task", task.ID)
	if err != nil && (inputs.BlockingUncertaintyCount > 0 || inputs.OpenRegression) {
		impact = nil
		err = nil
	}
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: why impact for task %q: %w", task.ID, err)
	}

	seedFreshness := freshnessFromStatus(task.Status)
	tasksFreshness := tasksSectionFreshness(tasks)
	planFreshness := planSectionFreshness(planView, in.Store)
	contextFreshness := contextSectionFreshness(contextPacket, seedFreshness)
	related := buildRelatedSection(ctx, in.Retrieval, contextPacket, contextFreshness, profile.RelatedDepth)

	delibSec := buildDeliberationSection(dState, inputs, profile, seedFreshness)
	if profile.IncludeOpenRegress {
		regs, err := buildOpenRegressionsSummary(ctx, dom, task.ID)
		if err != nil {
			return NextPacket{}, fmt.Errorf("loop next: open regressions: %w", err)
		}
		delibSec.OpenRegressions = regs
	}

	openUnc, err := buildOpenUncertaintiesSection(in.Store, task.ID, profile.OpenUncertaintyCap, seedFreshness)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: open uncertainties: %w", err)
	}

	verifyDebt, err := buildVerificationDebtSection(ctx, dom, task.ID, seedFreshness)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: verification debt: %w", err)
	}

	recentCap := profile.RecentChangesCap
	if recentCap <= 0 {
		recentCap = maxRecentChangesCap
	}
	recent, err := buildRecentChangesSection(in.Store, task.ID, recentCap, seedFreshness)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: recent changes: %w", err)
	}

	historical, err := buildHistoricalRelationshipsSection(in.Store, seedFreshness)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: historical relationships: %w", err)
	}

	riskHints, err := buildRiskHintsSection(ctx, dom, in.Store, task.ID, seedFreshness)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: risk hints: %w", err)
	}

	evidence := compiler.BuildEvidenceSections(in.Store, task.ID)
	planningEvidence := PlanningEvidenceSection{
		Freshness:        contextFreshness,
		Evaluations:      evidence.Evaluations,
		Reflections:      evidence.Reflections,
		PlanningEvidence: evidence.PlanningEvidence,
	}
	tendencies := TendenciesSection{
		Freshness: contextFreshness,
		Items:     evidence.Tendencies,
	}
	successfulApproaches := SuccessfulApproachesSection{
		Freshness: contextFreshness,
		Items:     evidence.SuccessfulApproaches,
	}
	similarChanges, err := buildSimilarChangesSection(ctx, dom, in.Store, task.ID, contextFreshness)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: similar changes: %w", err)
	}

	workConflicts, err := buildWorkConflictsSection(ctx, dom, task.ID, seedFreshness)
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: work conflicts: %w", err)
	}

	goalTitle := ""
	if g, err := in.Store.GetGoal(goalID); err == nil {
		goalTitle = g.Title
	}
	harnessRecs, err := buildHarnessRecommendationsSection(ctx, dom, in.Store, recommended, task.Title, goalKeywordsFromTitle(goalTitle))
	if err != nil {
		return NextPacket{}, fmt.Errorf("loop next: harness recommendations: %w", err)
	}

	return NextPacket{
		SchemaVersion: NextSchemaVersion,
		GeneratedAt:   time.Now().UTC(),
		Seed: SeedSection{
			Freshness: seedFreshness,
			TaskID:    task.ID,
			Title:     task.Title,
			WorkState: task.WorkState,
			GoalID:    goalID,
		},
		Tasks: TasksSection{
			Freshness: tasksFreshness,
			GoalID:    goalID,
			Items:     taskRows,
		},
		Plan: PlanSection{
			Freshness: planFreshness,
			Snapshot:  planSnapshot,
		},
		Why: WhySection{
			Freshness: seedFreshness,
			Snapshot: WhySnapshot{
				WhyResult: whyResult,
				Impact:    impact,
			},
		},
		Context: ContextSection{
			Freshness: contextFreshness,
			Snapshot:  contextPacket,
		},
		Related:                 related,
		PromotionCandidates:     promotionCandidates,
		Deliberation:            delibSec,
		OpenUncertainties:       openUnc,
		VerificationDebt:        verifyDebt,
		RecentChanges:           recent,
		HistoricalRelationships: historical,
		PlanningEvidence:        planningEvidence,
		Tendencies:              tendencies,
		SuccessfulApproaches:    successfulApproaches,
		SimilarChanges:          similarChanges,
		RiskHints:               riskHints,
		WorkConflicts:           workConflicts,
		HarnessRecommendations:  harnessRecs,
		LoopHints: LoopHintsSection{
			Freshness:         FreshnessUnknown,
			Available:         false,
			UnavailableReason: "iteration metadata unavailable in current store",
		},
	}, nil
}

func buildPromotionCandidates(st *store.Store) ([]PromotionCandidate, error) {
	return domain.New(st).ListPromotionCandidates()
}

func buildRelatedSection(ctx context.Context, eng *retrieval.Engine, pkt compiler.Packet, upstream string, depth int) RelatedSection {
	if depth <= 0 {
		depth = retrieval.DefaultImpactDepth()
	}
	seeds := discoverImpactSeeds(pkt)
	if len(seeds) == 0 {
		return RelatedSection{
			Freshness:         FreshnessUnknown,
			Available:         false,
			Seeds:             []retrieval.ImpactSeed{},
			UnavailableReason: "no file or symbol seeds discovered from context",
		}
	}
	res, err := eng.ImpactWalk(ctx, seeds, depth)
	if err != nil {
		return RelatedSection{
			Freshness:         FreshnessUnknown,
			Available:         false,
			Seeds:             seeds,
			UnavailableReason: err.Error(),
		}
	}
	return RelatedSection{
		Freshness: upstream,
		Available: true,
		Seeds:     seeds,
		Snapshot:  res,
	}
}

func discoverImpactSeeds(pkt compiler.Packet) []retrieval.ImpactSeed {
	out := make([]retrieval.ImpactSeed, 0)
	seen := map[string]struct{}{}
	for _, item := range pkt.Items {
		if item.EntityType != "file" && item.EntityType != "symbol" {
			continue
		}
		key := item.EntityType + "\x00" + item.EntityID
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, retrieval.ImpactSeed{EntityType: item.EntityType, EntityID: item.EntityID})
	}
	return out
}

func freshnessFromStatus(status string) string {
	switch status {
	case store.StatusStale, store.StatusSuperseded:
		return FreshnessStale
	default:
		return FreshnessFresh
	}
}

func tasksSectionFreshness(tasks []store.Task) string {
	for _, task := range tasks {
		if freshnessFromStatus(task.Status) == FreshnessStale {
			return FreshnessStale
		}
	}
	return FreshnessFresh
}

func planSectionFreshness(view planner.PlanView, st *store.Store) string {
	if view.CurrentScopeID == nil || *view.CurrentScopeID == "" {
		return FreshnessUnknown
	}
	scope, err := st.GetPlanScope(*view.CurrentScopeID)
	if err != nil {
		return FreshnessUnknown
	}
	switch scope.Status {
	case store.StatusStale, store.StatusSuperseded:
		return FreshnessStale
	default:
		return FreshnessFresh
	}
}

func contextSectionFreshness(pkt compiler.Packet, fallback string) string {
	if pkt.IndexHonesty != nil && pkt.IndexHonesty.StaleTotal > 0 {
		return FreshnessDirty
	}
	if pkt.GraphSyncHonesty != nil && pkt.GraphSyncHonesty.StaleCommit {
		return FreshnessDirty
	}
	return fallback
}

func minimalTaskContextPacket(st *store.Store, task store.Task, maxItems int) (compiler.Packet, error) {
	if maxItems <= 0 {
		maxItems = compiler.DefaultMaxItems
	}
	d0 := 0
	items := []compiler.Item{
		{
			EntityType: "task", EntityID: task.ID, Title: task.Title,
			ReasonCode: retrieval.ReasonDirectTaskScope, Distance: &d0,
			Trust: compiler.TrustUntrustedData, Layer: 0,
		},
	}
	if task.GoalID != nil && *task.GoalID != "" {
		if g, err := st.GetGoal(*task.GoalID); err == nil {
			items = append(items, compiler.Item{
				EntityType: "goal", EntityID: g.ID, Title: g.Title,
				ReasonCode: retrieval.ReasonGoalHasTask, Distance: &d0,
				Trust: compiler.TrustUntrustedData, Layer: 0,
			})
		}
	}
	if len(items) > maxItems {
		items = items[:maxItems]
	}
	evidence := compiler.BuildEvidenceSections(st, task.ID)
	return compiler.Packet{
		SchemaVersion: compiler.SchemaVersion,
		TaskID:        task.ID,
		Budget: compiler.Budget{
			MaxItems: maxItems, ItemsTotal: len(items), ItemsKept: len(items),
		},
		Items:                items,
		Evaluations:          evidence.Evaluations,
		Reflections:          evidence.Reflections,
		PlanningEvidence:     evidence.PlanningEvidence,
		Tendencies:           evidence.Tendencies,
		SuccessfulApproaches: evidence.SuccessfulApproaches,
	}, nil
}

func buildSimilarChangesSection(ctx context.Context, dom *domain.Service, st *store.Store, taskID, freshness string) (SimilarChangesSection, error) {
	empty := SimilarChangesSection{Freshness: freshness, Items: []domain.SimilarChangeRow{}}
	changes, err := st.ListChangesByTaskID(taskID)
	if err != nil {
		return SimilarChangesSection{}, err
	}
	var latest *store.Change
	for i := range changes {
		c := changes[i]
		if c.Status == store.ChangeStatusSuperseded {
			continue
		}
		if latest == nil || c.CreatedAt > latest.CreatedAt {
			latest = &c
		}
	}
	if latest == nil {
		return empty, nil
	}
	paths, err := st.ListChangePaths(latest.ID)
	if err != nil {
		return SimilarChangesSection{}, err
	}
	if len(paths) == 0 {
		return empty, nil
	}
	pathStrs := make([]string, 0, len(paths))
	for _, p := range paths {
		pathStrs = append(pathStrs, p.Path)
	}
	kind := domain.InferChangeKind(pathStrs)
	result, err := dom.QuerySimilarChanges(ctx, domain.SimilarChangesOpts{
		ChangeKind: kind,
		Limit:      8,
	})
	if err != nil {
		return SimilarChangesSection{}, err
	}
	items := result.Changes
	if len(items) > 8 {
		items = items[:8]
	}
	if items == nil {
		items = []domain.SimilarChangeRow{}
	}
	return SimilarChangesSection{Freshness: freshness, Items: items}, nil
}

const subagentHookSlug = "hook:harness:subagent"

func buildHarnessRecommendationsSection(
	ctx context.Context,
	dom *domain.Service,
	st *store.Store,
	phase deliberation.Phase,
	taskTitle string,
	goalKeywords []string,
) (HarnessRecommendationsSection, error) {
	catalog, err := st.ListHarnessAgents()
	if err != nil {
		return HarnessRecommendationsSection{}, err
	}
	if len(catalog) == 0 {
		return HarnessRecommendationsSection{
			Freshness: FreshnessUnknown,
			Items:     []agents.Recommendation{},
		}, nil
	}

	harnessCaps, err := loadHarnessCapsForRouting(ctx, dom, st)
	if err != nil {
		return HarnessRecommendationsSection{}, err
	}

	recs, err := agents.RecommendAgents(ctx, st, agents.RecommendInput{
		Phase:        string(phase),
		TaskTitle:    taskTitle,
		GoalKeywords: goalKeywords,
		HarnessCaps:  harnessCaps,
	})
	if err != nil {
		return HarnessRecommendationsSection{}, err
	}
	if recs == nil {
		recs = []agents.Recommendation{}
	}
	return HarnessRecommendationsSection{
		Freshness: FreshnessFresh,
		Items:     recs,
	}, nil
}

func loadHarnessCapsForRouting(ctx context.Context, dom *domain.Service, st *store.Store) (map[string]string, error) {
	slugSet := map[string]struct{}{subagentHookSlug: {}}
	agents, err := st.ListHarnessAgents()
	if err != nil {
		return nil, err
	}
	for _, agent := range agents {
		reqs, err := st.ListHarnessAgentRequirements(agent.ID)
		if err != nil {
			return nil, err
		}
		for _, req := range reqs {
			slugSet[req.RequiredCapabilitySlug] = struct{}{}
		}
	}

	caps := make(map[string]string, len(slugSet))
	for slug := range slugSet {
		cap, err := dom.GetCapabilityBySlug(ctx, slug)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				caps[slug] = store.CapabilityStatusUnknown
				continue
			}
			return nil, err
		}
		caps[slug] = cap.Status
	}
	if status, ok := caps[subagentHookSlug]; ok {
		caps["harness:subagent"] = status
	}
	return caps, nil
}

func goalKeywordsFromTitle(title string) []string {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil
	}
	return []string{strings.ToLower(title)}
}

func buildWorkConflictsSection(ctx context.Context, dom *domain.Service, seedTaskID, freshness string) (WorkConflictsSection, error) {
	const workConflictsCap = 8
	conflicts, err := dom.DetectWorkConflicts(ctx, domain.DetectWorkConflictsOpts{
		TaskID: seedTaskID,
		Limit:  workConflictsCap,
	})
	if err != nil {
		return WorkConflictsSection{}, err
	}
	if conflicts == nil {
		conflicts = []domain.WorkConflict{}
	}
	return WorkConflictsSection{Freshness: freshness, Items: conflicts}, nil
}
