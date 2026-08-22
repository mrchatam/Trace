package loop_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func TestLoopNextPlanningEvidenceSection(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	_, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	b, err := dsvc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.9}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskID, BaselineID: b.ID, ScoresJSON: `{"correctness": 0.7}`,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.CreateReflection(ctx, domain.ReflectionInput{
		TaskID: taskID, Summary: "loop reflection", UsefulTests: []string{"TestLoop"},
	}); err != nil {
		t.Fatal(err)
	}

	out := mustEvalOutcome(t, dsvc, taskID)
	if _, err := dsvc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: taskID,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: taskID, TestName: "TestBroken", TestStatus: store.TestStatusFail, Summary: "fail",
	}); err != nil {
		t.Fatal(err)
	}
	chg, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID, GitCommit: "fedcba9",
		Paths: []domain.ChangePathInput{{Path: "loop.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.RecordImprovement(ctx, domain.ImprovementInput{
		ChangeID: chg.ID, TaskID: taskID, Summary: "narrow scope",
	}); err != nil {
		t.Fatal(err)
	}

	pkt, err := loop.BuildNextPacket(ctx, loop.BuildNextInput{
		TaskID:    taskID,
		Store:     st,
		Planner:   psvc,
		Retrieval: retrieval.New(st),
		Compiler:  compiler.New(st).WithRetrieval(retrieval.New(st)),
	})
	if err != nil {
		t.Fatalf("BuildNextPacket: %v", err)
	}

	sec := pkt.PlanningEvidence
	if sec.Freshness == "" {
		t.Fatalf("planning_evidence freshness empty: %+v", sec)
	}
	if len(sec.Evaluations) == 0 {
		t.Fatalf("planning_evidence evaluations empty: %+v", sec)
	}
	if len(sec.Reflections) == 0 {
		t.Fatalf("planning_evidence reflections empty: %+v", sec)
	}
	if len(sec.PlanningEvidence) < 3 {
		t.Fatalf("planning_evidence items: %+v", sec.PlanningEvidence)
	}

	kinds := map[string]bool{}
	for _, pe := range sec.PlanningEvidence {
		kinds[pe.EntityType] = true
	}
	if !kinds["regression"] || !kinds["outcome_result"] || !kinds["improvement"] {
		t.Fatalf("planning kinds=%v items=%+v", kinds, sec.PlanningEvidence)
	}

	// Context snapshot mirrors compiler evidence fields.
	if len(pkt.Context.Snapshot.Evaluations) == 0 || len(pkt.Context.Snapshot.Reflections) == 0 {
		t.Fatalf("context snapshot evidence: eval=%d refl=%d",
			len(pkt.Context.Snapshot.Evaluations), len(pkt.Context.Snapshot.Reflections))
	}

	raw, err := json.Marshal(pkt)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["schema_version"] != loop.NextSchemaVersion {
		t.Fatalf("schema: %v", decoded["schema_version"])
	}
	peSec, ok := decoded["planning_evidence"].(map[string]any)
	if !ok {
		t.Fatalf("planning_evidence section: %#v", decoded["planning_evidence"])
	}
	for _, key := range []string{"evaluations", "reflections", "planning_evidence"} {
		arr, ok := peSec[key].([]any)
		if !ok || len(arr) == 0 {
			t.Fatalf("planning_evidence.%s: %#v", key, peSec[key])
		}
	}
}

func TestLoopNextIncludesEvidenceForDecisions(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	_, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	for i := 0; i < 2; i++ {
		chg, err := dsvc.CreateChange(ctx, domain.ChangeInput{
			TaskID: taskID,
			Reason: "internal improvement",
			Paths:  []domain.ChangePathInput{{Path: "internal/x.go", Status: "M"}},
		})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := dsvc.RecordImprovement(ctx, domain.ImprovementInput{
			ChangeID: chg.ID, TaskID: taskID, Summary: "retry backoff",
		}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := dsvc.RefreshChangePatterns(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.UpsertEngineeringKnowledge(ctx, domain.EngineeringKnowledgeInput{
		Title: "Cache pattern",
		Body: domain.KnowledgeBody{
			Summary:          "Use bounded LRU",
			SourceEntityType: "manual",
			SourceEntityID:   "note-1",
		},
		Topic:      "improvement",
		SourceType: "USER_ASSERTED",
	}); err != nil {
		t.Fatal(err)
	}

	b, err := dsvc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.9}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskID, BaselineID: b.ID, ScoresJSON: `{"correctness": 0.7}`,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.CreateReflection(ctx, domain.ReflectionInput{
		TaskID: taskID, Summary: "loop reflection", UsefulTests: []string{"TestLoop"},
	}); err != nil {
		t.Fatal(err)
	}
	out := mustEvalOutcome(t, dsvc, taskID)
	if _, err := dsvc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: taskID,
	}); err != nil {
		t.Fatal(err)
	}

	pkt, err := loop.BuildNextPacket(ctx, loop.BuildNextInput{
		TaskID:    taskID,
		Store:     st,
		Planner:   psvc,
		Retrieval: retrieval.New(st),
		Compiler:  compiler.New(st).WithRetrieval(retrieval.New(st)),
	})
	if err != nil {
		t.Fatalf("BuildNextPacket: %v", err)
	}

	if len(pkt.Tendencies.Items) == 0 {
		t.Fatalf("tendencies empty: %+v", pkt.Tendencies)
	}
	if len(pkt.SuccessfulApproaches.Items) == 0 {
		t.Fatalf("successful_approaches empty: %+v", pkt.SuccessfulApproaches)
	}
	if len(pkt.PlanningEvidence.Evaluations) == 0 || len(pkt.PlanningEvidence.PlanningEvidence) == 0 {
		t.Fatalf("planning_evidence incomplete: %+v", pkt.PlanningEvidence)
	}
	if len(pkt.Context.Snapshot.Tendencies) == 0 || len(pkt.Context.Snapshot.SuccessfulApproaches) == 0 {
		t.Fatalf("context snapshot learning fields: tend=%d sa=%d",
			len(pkt.Context.Snapshot.Tendencies), len(pkt.Context.Snapshot.SuccessfulApproaches))
	}

	raw, err := json.Marshal(pkt)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"tendencies", "successful_approaches", "planning_evidence", "similar_changes"} {
		sec, ok := decoded[key].(map[string]any)
		if !ok {
			t.Fatalf("section %s: %#v", key, decoded[key])
		}
		if sec["freshness"] == "" {
			t.Fatalf("%s freshness empty", key)
		}
	}
}

func mustEvalOutcome(t *testing.T, svc *domain.Service, taskID string) store.OutcomeResult {
	t.Helper()
	ctx := context.Background()
	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.98}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     taskID,
		BaselineID: b.ID,
		ScoresJSON: `{"correctness": 0.40}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestLoopNextIncludesWorkConflicts(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	gid := goalID
	other, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "parallel work", GoalID: &gid})
	if err != nil {
		t.Fatal(err)
	}
	if err := dsvc.TransitionTask(ctx, taskID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "test", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
	if err := dsvc.TransitionTask(ctx, other.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "test", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID,
		Paths:  []domain.ChangePathInput{{Path: "internal/loop/work"}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID: other.ID,
		Paths:  []domain.ChangePathInput{{Path: "internal/loop/work/helper.go"}},
	}); err != nil {
		t.Fatal(err)
	}

	pkt, err := loop.BuildNextPacket(ctx, loop.BuildNextInput{
		TaskID:    taskID,
		Store:     st,
		Planner:   psvc,
		Retrieval: retrieval.New(st),
		Compiler:  compiler.New(st).WithRetrieval(retrieval.New(st)),
	})
	if err != nil {
		t.Fatalf("BuildNextPacket: %v", err)
	}
	if len(pkt.WorkConflicts.Items) != 1 {
		t.Fatalf("work_conflicts len=%d want 1: %+v", len(pkt.WorkConflicts.Items), pkt.WorkConflicts)
	}
	if pkt.WorkConflicts.Freshness == "" {
		t.Fatalf("work_conflicts freshness empty")
	}

	raw, err := json.Marshal(pkt)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	wc, ok := decoded["work_conflicts"].(map[string]any)
	if !ok {
		t.Fatalf("work_conflicts section missing: keys=%v", decoded)
	}
	items, ok := wc["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("work_conflicts.items=%#v", wc["items"])
	}
}

func TestLoopNextPromotionCandidates(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	_, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	goalID, _, _ := seedGoalTaskPlan(t, psvc, dsvc)

	blockingUnlinked, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title:    "Promote me",
		Severity: domain.SeverityBlocking,
	})
	if err != nil {
		t.Fatal(err)
	}
	blockingLinked, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title:    "Already promoted",
		Severity: domain.SeverityBlocking,
	})
	if err != nil {
		t.Fatal(err)
	}
	linkedTask, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "child", GoalID: &goalID})
	if err != nil {
		t.Fatal(err)
	}
	if err := dsvc.LinkDiscoveryMentionsTask(ctx, blockingLinked.ID, linkedTask.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title:    "Noise info",
		Severity: domain.SeverityINFO,
	}); err != nil {
		t.Fatal(err)
	}

	pkt := buildLoopNextPacket(t, st, psvc, taskID)
	if len(pkt.PromotionCandidates) != 1 {
		t.Fatalf("promotion_candidates len=%d want 1: %+v", len(pkt.PromotionCandidates), pkt.PromotionCandidates)
	}
	if pkt.PromotionCandidates[0].DiscoveryID != blockingUnlinked.ID {
		t.Fatalf("wrong candidate: %+v", pkt.PromotionCandidates[0])
	}

	raw, err := json.Marshal(pkt)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	sec, ok := decoded["promotion_candidates"].([]any)
	if !ok {
		t.Fatalf("promotion_candidates missing from json: %#v", decoded["promotion_candidates"])
	}
	if len(sec) != 1 {
		t.Fatalf("promotion_candidates json len=%d", len(sec))
	}
}

func seedDefaultHarnessCatalog(t *testing.T, svc *domain.Service) {
	t.Helper()
	ctx := context.Background()
	catalog := []domain.HarnessAgentInput{
		{Slug: "agent:code-reviewer", Title: "Code Reviewer", SubagentType: "code-reviewer", RecommendSubagent: true, DeliberationPhases: `["CRITIQUE"]`},
		{Slug: "agent:nested-reviewer", Title: "Nested Reviewer", SubagentType: "nested-reviewer", RecommendSubagent: true, DeliberationPhases: `["CRITIQUE"]`},
		{Slug: "agent:performance-reviewer", Title: "Performance Reviewer", SubagentType: "performance-reviewer", DeliberationPhases: `["VERIFY"]`, TaskKeywords: `["perf","performance","latency"]`},
		{Slug: "agent:security-reviewer", Title: "Security Reviewer", SubagentType: "security-reviewer", DeliberationPhases: `["CRITIQUE","VERIFY"]`},
		{Slug: "agent:explore", Title: "Explore", SubagentType: "explore", DeliberationPhases: `["INVESTIGATE","ORIENT"]`},
		{Slug: "agent:generalPurpose", Title: "General Purpose", SubagentType: "generalPurpose"},
	}
	for _, in := range catalog {
		if _, err := svc.UpsertHarnessAgent(ctx, in); err != nil {
			t.Fatalf("UpsertHarnessAgent %s: %v", in.Slug, err)
		}
	}
}

func seedGoalTaskPlanWithTitle(t *testing.T, psvc *planner.Service, dsvc *domain.Service, title string) (goalID, taskID, scopeID string) {
	t.Helper()
	ctx := context.Background()
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "loop goal"})
	if err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: title, GoalID: &g.ID})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	cp, err := psvc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{
			Title:  "Phase 1",
			Scopes: []planner.ScopeInput{{Title: "Scope 1"}},
		}},
	})
	if err != nil {
		t.Fatalf("CreateCoarsePlan: %v", err)
	}
	scopeID = cp.Phases[0].Scopes[0].ID
	if err := psvc.SetCurrentScope(ctx, g.ID, scopeID); err != nil {
		t.Fatalf("SetCurrentScope: %v", err)
	}
	if _, err := psvc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID:      scopeID,
		ExitCriteria: []string{"packet ready"},
		WorkItems:    []planner.WorkItem{{Title: "emit loop packet"}},
	}); err != nil {
		t.Fatalf("DeepPlan: %v", err)
	}
	return g.ID, task.ID, scopeID
}

func seedVerifyPhase(t *testing.T, st *store.Store, dsvc *domain.Service, goalID, taskID string) {
	t.Helper()
	ctx := context.Background()
	markPlanCritiqued(t, st, taskID, goalID)
	if _, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    taskID,
		GitCommit: "abc1234",
		Paths:     []domain.ChangePathInput{{Path: "bench.go"}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := dsvc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: taskID, TestName: "TestBench", TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
}

func buildLoopNextPacket(t *testing.T, st *store.Store, psvc *planner.Service, taskID string) loop.NextPacket {
	t.Helper()
	ctx := context.Background()
	pkt, err := loop.BuildNextPacket(ctx, loop.BuildNextInput{
		TaskID:    taskID,
		Store:     st,
		Planner:   psvc,
		Retrieval: retrieval.New(st),
		Compiler:  compiler.New(st).WithRetrieval(retrieval.New(st)),
	})
	if err != nil {
		t.Fatalf("BuildNextPacket: %v", err)
	}
	return pkt
}

func TestLoopNextIncludesHarnessRecommendations(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	seedDefaultHarnessCatalog(t, dsvc)
	_, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	pkt := buildLoopNextPacket(t, st, psvc, taskID)
	sec := pkt.HarnessRecommendations
	if sec.Freshness != loop.FreshnessFresh {
		t.Fatalf("freshness=%q want fresh", sec.Freshness)
	}
	if len(sec.Items) == 0 {
		t.Fatalf("harness_recommendations empty: %+v", sec)
	}
	if pkt.Deliberation.Phase != "CRITIQUE" {
		t.Fatalf("phase=%q want CRITIQUE for uncritiqued plan", pkt.Deliberation.Phase)
	}

	raw, err := json.Marshal(pkt)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	hr, ok := decoded["harness_recommendations"].(map[string]any)
	if !ok {
		t.Fatalf("harness_recommendations section missing: keys=%v", decoded)
	}
	items, ok := hr["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatalf("harness_recommendations.items=%#v", hr["items"])
	}
}

func TestRecommendSubagentWhenAvailable(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	seedDefaultHarnessCatalog(t, dsvc)
	ctx := context.Background()
	_, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	if _, err := dsvc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: domain.CapabilityKindHook, Slug: "hook:harness:subagent",
		Title: "Harness subagent", Status: domain.CapabilityStatusAvailable,
	}); err != nil {
		t.Fatal(err)
	}

	pkt := buildLoopNextPacket(t, st, psvc, taskID)
	if len(pkt.HarnessRecommendations.Items) == 0 {
		t.Fatalf("no recommendations: %+v", pkt.HarnessRecommendations)
	}
	first := pkt.HarnessRecommendations.Items[0]
	if first.AgentSlug != "agent:code-reviewer" {
		t.Fatalf("want code-reviewer first, got %+v", first)
	}
	if !first.UseSubagent {
		t.Fatalf("use_subagent want true when hook AVAILABLE: %+v", first)
	}
	if first.PromptStub != "Fresh subagent for independent review — not the implementer session." {
		t.Fatalf("prompt_stub=%q", first.PromptStub)
	}
}

func TestRecommendSubagentHonestWhenUnavailable(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	seedDefaultHarnessCatalog(t, dsvc)
	_, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	pkt := buildLoopNextPacket(t, st, psvc, taskID)
	if len(pkt.HarnessRecommendations.Items) == 0 {
		t.Fatal("expected recommendations")
	}
	first := pkt.HarnessRecommendations.Items[0]
	if first.UseSubagent {
		t.Fatalf("use_subagent must be false when hook UNKNOWN: %+v", first)
	}
	if first.PromptStub != "" {
		t.Fatalf("prompt_stub must be empty when subagent unavailable: %+v", first)
	}
}

func TestRecommendPerformanceReviewerForPerfTask(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	seedDefaultHarnessCatalog(t, dsvc)
	goalID, taskID, _ := seedGoalTaskPlanWithTitle(t, psvc, dsvc, "Fix latency regression in benchmark suite")
	seedVerifyPhase(t, st, dsvc, goalID, taskID)

	pkt := buildLoopNextPacket(t, st, psvc, taskID)
	if pkt.Deliberation.Phase != "VERIFY" {
		t.Fatalf("phase=%q want VERIFY", pkt.Deliberation.Phase)
	}
	if len(pkt.HarnessRecommendations.Items) == 0 {
		t.Fatal("no recommendations")
	}
	if pkt.HarnessRecommendations.Items[0].AgentSlug != "agent:performance-reviewer" {
		t.Fatalf("want performance-reviewer first, got %+v", pkt.HarnessRecommendations.Items)
	}
}

func TestNoRecommendationWhenCatalogEmpty(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	_, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	pkt := buildLoopNextPacket(t, st, psvc, taskID)
	sec := pkt.HarnessRecommendations
	if sec.Freshness != loop.FreshnessUnknown {
		t.Fatalf("freshness=%q want unknown", sec.Freshness)
	}
	if len(sec.Items) != 0 {
		t.Fatalf("items=%+v want []", sec.Items)
	}
}
