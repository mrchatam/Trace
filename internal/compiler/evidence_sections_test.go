package compiler_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestContextIncludesEvaluationsAndReflections(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	c := compiler.New(st)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "evidence goal"})
	if err != nil {
		t.Fatal(err)
	}
	taskA, err := svc.CreateTask(ctx, domain.TaskInput{Title: "task A", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	taskB, err := svc.CreateTask(ctx, domain.TaskInput{Title: "task B", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}

	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness": 0.9}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	evalA, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskA.ID, BaselineID: b.ID, ScoresJSON: `{"correctness": 0.85}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: taskB.ID, BaselineID: b.ID, ScoresJSON: `{"correctness": 0.1}`,
	}); err != nil {
		t.Fatal(err)
	}

	refA, err := svc.CreateReflection(ctx, domain.ReflectionInput{
		TaskID: taskA.ID, Summary: "reflection A", UsefulTests: []string{"TestCache"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateReflection(ctx, domain.ReflectionInput{
		TaskID: taskB.ID, Summary: "foreign reflection", UsefulTests: []string{"TestOther"},
	}); err != nil {
		t.Fatal(err)
	}

	// Planning evidence on task A: open regression + failed test + improvement.
	out := mustEvalWithRegressionForCompiler(t, svc, taskA.ID)
	reg, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: taskA.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: taskA.ID, TestName: "TestFail", TestStatus: store.TestStatusFail, Summary: "boom",
	}); err != nil {
		t.Fatal(err)
	}
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskA.ID, GitCommit: "def5678",
		Paths: []domain.ChangePathInput{{Path: "a.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
		ChangeID: chg.ID, TaskID: taskA.ID, Summary: "retry backoff",
	}); err != nil {
		t.Fatal(err)
	}

	pkt, err := c.TaskContext(ctx, taskA.ID, compiler.ContextOptions{IncludeMarkdown: true})
	if err != nil {
		t.Fatalf("TaskContext: %v", err)
	}

	if len(pkt.Evaluations) == 0 {
		t.Fatalf("evaluations empty")
	}
	foundEval := false
	for _, ev := range pkt.Evaluations {
		if ev.ID == evalA.ID {
			foundEval = true
		}
		if ev.TaskID != taskA.ID {
			t.Fatalf("evaluation task_id leak: %+v", ev)
		}
	}
	if !foundEval {
		t.Fatalf("evaluations missing evalA: %+v", pkt.Evaluations)
	}
	if !strings.Contains(pkt.Evaluations[0].ScoresJSON, "correctness") && !strings.Contains(pkt.Evaluations[len(pkt.Evaluations)-1].ScoresJSON, "correctness") {
		t.Fatalf("scores_json missing: %+v", pkt.Evaluations)
	}

	if len(pkt.Reflections) == 0 || pkt.Reflections[0].ID != refA.ID {
		t.Fatalf("reflections: %+v want id=%s", pkt.Reflections, refA.ID)
	}

	if len(pkt.PlanningEvidence) < 3 {
		t.Fatalf("planning_evidence: %+v want regression+fail+improvement", pkt.PlanningEvidence)
	}
	kinds := map[string]bool{}
	for _, pe := range pkt.PlanningEvidence {
		kinds[pe.EntityType] = true
		if pe.EntityID == reg.ID && pe.EntityType != "regression" {
			t.Fatalf("regression item: %+v", pe)
		}
	}
	if !kinds["regression"] || !kinds["outcome_result"] || !kinds["improvement"] {
		t.Fatalf("planning kinds=%v items=%+v", kinds, pkt.PlanningEvidence)
	}

	raw, err := pkt.JSON()
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"evaluations", "reflections", "planning_evidence"} {
		arr, ok := decoded[key].([]any)
		if !ok || len(arr) == 0 {
			t.Fatalf("JSON %s: %#v", key, decoded[key])
		}
	}

	foreign, err := c.TaskContext(ctx, taskB.ID, compiler.ContextOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, ev := range foreign.Evaluations {
		if ev.ID == evalA.ID {
			t.Fatalf("cross-task evaluation leak: %+v", ev)
		}
	}
	for _, rf := range foreign.Reflections {
		if rf.ID == refA.ID {
			t.Fatalf("cross-task reflection leak: %+v", rf)
		}
	}

	md := pkt.Markdown()
	for _, heading := range []string{"## Evaluations", "## Reflections", "## Planning evidence"} {
		if !strings.Contains(md, heading) {
			t.Fatalf("markdown missing %q:\n%s", heading, md)
		}
	}
}

func mustEvalWithRegressionForCompiler(t *testing.T, svc *domain.Service, taskID string) store.OutcomeResult {
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
		ScoresJSON: `{"correctness": 0.50}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestTendHelpHurtInContext(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	c := compiler.New(st)
	ctx := context.Background()

	task := mustChangeTaskForCompiler(t, svc)
	for i := 0; i < 2; i++ {
		chg, err := svc.CreateChange(ctx, domain.ChangeInput{
			TaskID: task.ID,
			Reason: "internal improvement",
			Paths:  []domain.ChangePathInput{{Path: "internal/a.go", Status: "M"}},
		})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
			ChangeID: chg.ID, TaskID: task.ID, Summary: "faster path",
		}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.RefreshChangePatterns(ctx); err != nil {
		t.Fatal(err)
	}

	pkt, err := c.TaskContext(ctx, task.ID, compiler.ContextOptions{IncludeMarkdown: true})
	if err != nil {
		t.Fatalf("TaskContext: %v", err)
	}
	if len(pkt.Tendencies) == 0 {
		t.Fatalf("tendencies empty")
	}
	foundImprove := false
	for _, td := range pkt.Tendencies {
		if td.Direction == domain.TendencyDirectionImprove && td.CountPositive >= 2 {
			foundImprove = true
		}
	}
	if !foundImprove {
		t.Fatalf("tendencies missing improve threshold: %+v", pkt.Tendencies)
	}

	raw, err := pkt.JSON()
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	arr, ok := decoded["tendencies"].([]any)
	if !ok || len(arr) == 0 {
		t.Fatalf("JSON tendencies: %#v", decoded["tendencies"])
	}
	if md := pkt.Markdown(); !strings.Contains(md, "## Tendencies") {
		t.Fatalf("markdown missing tendencies:\n%s", md)
	}
}

func TestSuccessfulApproachesSurfaced(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	c := compiler.New(st)
	ctx := context.Background()

	task := mustChangeTaskForCompiler(t, svc)
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "improve cache",
		Paths:  []domain.ChangePathInput{{Path: "internal/cache.go", Status: "M"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	imp, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
		ChangeID: chg.ID, TaskID: task.ID, Summary: "LRU cache helped",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "TestCache", TestStatus: store.TestStatusPass, Summary: "pass",
	}); err != nil {
		t.Fatal(err)
	}
	know, err := svc.UpsertEngineeringKnowledge(ctx, domain.EngineeringKnowledgeInput{
		Title: "Bounded cache",
		Body: domain.KnowledgeBody{
			Summary:          "Use bounded LRU",
			SourceEntityType: "manual",
			SourceEntityID:   "note-cache",
		},
		Topic:      "improvement",
		SourceType: "USER_ASSERTED",
	})
	if err != nil {
		t.Fatal(err)
	}

	pkt, err := c.TaskContext(ctx, task.ID, compiler.ContextOptions{IncludeMarkdown: true})
	if err != nil {
		t.Fatalf("TaskContext: %v", err)
	}
	if len(pkt.SuccessfulApproaches) < 2 {
		t.Fatalf("successful_approaches: %+v", pkt.SuccessfulApproaches)
	}
	ids := map[string]bool{}
	for _, sa := range pkt.SuccessfulApproaches {
		ids[sa.ID] = true
	}
	if !ids[imp.ID] || !ids[know.ID] {
		t.Fatalf("missing worked/knowledge ids: got=%v want imp=%s know=%s", ids, imp.ID, know.ID)
	}

	raw, err := pkt.JSON()
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	arr, ok := decoded["successful_approaches"].([]any)
	if !ok || len(arr) == 0 {
		t.Fatalf("JSON successful_approaches: %#v", decoded["successful_approaches"])
	}
	if md := pkt.Markdown(); !strings.Contains(md, "## Successful approaches") {
		t.Fatalf("markdown missing successful approaches:\n%s", md)
	}
}

func mustChangeTaskForCompiler(t *testing.T, svc *domain.Service) store.Task {
	t.Helper()
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "patterns goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "patterns task", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	return task
}
