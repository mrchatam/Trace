package loop_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/store"
)

func seedRiskTask(t *testing.T, dsvc *domain.Service) (context.Context, string) {
	t.Helper()
	ctx := context.Background()
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "risk goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "risk task", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	return ctx, task.ID
}

func addChangeWithPaths(t *testing.T, dsvc *domain.Service, taskID string, paths ...string) {
	t.Helper()
	inputs := make([]domain.ChangePathInput, 0, len(paths))
	for _, p := range paths {
		inputs = append(inputs, domain.ChangePathInput{Path: p})
	}
	if _, err := dsvc.CreateChange(context.Background(), domain.ChangeInput{
		TaskID:    taskID,
		GitCommit: "abc1234",
		Paths:     inputs,
	}); err != nil {
		t.Fatalf("CreateChange: %v", err)
	}
}

func TestRiskHintsManyPaths(t *testing.T) {
	st, _, dsvc := openLoopTestStore(t)
	ctx, taskID := seedRiskTask(t, dsvc)

	paths := make([]string, 9)
	for i := range paths {
		paths[i] = fmt.Sprintf("internal/pkg/file%d.go", i)
	}
	addChangeWithPaths(t, dsvc, taskID, paths...)

	sec, err := loopBuildRiskHints(ctx, dsvc, st, taskID)
	if err != nil {
		t.Fatalf("buildRiskHintsSection: %v", err)
	}
	if !riskHintContains(sec, "many_paths") {
		t.Fatalf("expected many_paths hint: %+v", sec.Items)
	}
}

func TestRiskHintsBlockingUncertainty(t *testing.T) {
	st, _, dsvc := openLoopTestStore(t)
	ctx, taskID := seedRiskTask(t, dsvc)

	if _, err := dsvc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "blocking question",
		Severity: store.UncertaintySeverityBlocking,
		TaskID:   taskID,
	}); err != nil {
		t.Fatalf("CreateUncertainty: %v", err)
	}

	sec, err := loopBuildRiskHints(ctx, dsvc, st, taskID)
	if err != nil {
		t.Fatalf("buildRiskHintsSection: %v", err)
	}
	if !riskHintContains(sec, "blocking_uncertainty") {
		t.Fatalf("expected blocking_uncertainty hint: %+v", sec.Items)
	}
}

func TestLoopNextRiskHintsBounded(t *testing.T) {
	st, _, dsvc := openLoopTestStore(t)
	ctx, taskID := seedRiskTask(t, dsvc)

	if _, err := dsvc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "blocker",
		Severity: store.UncertaintySeverityBlocking,
		TaskID:   taskID,
	}); err != nil {
		t.Fatalf("CreateUncertainty: %v", err)
	}

	churnPath := "internal/loop/apply.go"
	for i := 0; i < 3; i++ {
		addChangeWithPaths(t, dsvc, taskID, churnPath, fmt.Sprintf("internal/other/file%d.go", i))
	}

	manyPaths := make([]string, 9)
	for i := range manyPaths {
		manyPaths[i] = fmt.Sprintf("internal/wide/file%d.go", i)
	}
	addChangeWithPaths(t, dsvc, taskID, manyPaths...)

	sec, err := loopBuildRiskHints(ctx, dsvc, st, taskID)
	if err != nil {
		t.Fatalf("buildRiskHintsSection: %v", err)
	}
	if len(sec.Items) > 4 {
		t.Fatalf("risk hints cap exceeded: %d items", len(sec.Items))
	}
	if len(sec.Items) == 0 {
		t.Fatal("expected at least one risk hint")
	}
	if sec.Items[0].Code != "blocking_uncertainty" {
		t.Fatalf("priority order: first hint=%q want blocking_uncertainty; items=%+v", sec.Items[0].Code, sec.Items)
	}
}

func loopBuildRiskHints(ctx context.Context, dsvc *domain.Service, st *store.Store, taskID string) (loop.RiskHintsSection, error) {
	return loop.BuildRiskHintsSectionForTest(ctx, dsvc, st, taskID, loop.FreshnessFresh)
}

func riskHintContains(sec loop.RiskHintsSection, code string) bool {
	for _, item := range sec.Items {
		if item.Code == code {
			return true
		}
	}
	return false
}
