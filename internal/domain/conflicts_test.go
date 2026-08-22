package domain_test

import (
	"context"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func mustInProgress(t *testing.T, svc *domain.Service, taskID string) {
	t.Helper()
	if err := svc.TransitionTask(context.Background(), taskID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "test", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestDetectOverlappingOpenTasks(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "overlap goal"})
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
	mustInProgress(t, svc, taskA.ID)
	mustInProgress(t, svc, taskB.ID)

	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskA.ID,
		Paths:  []domain.ChangePathInput{{Path: "internal/pkg/foo"}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskB.ID,
		Paths:  []domain.ChangePathInput{{Path: "internal/pkg/foo/bar.go"}},
	}); err != nil {
		t.Fatal(err)
	}
	_ = st

	conflicts, err := svc.DetectWorkConflicts(ctx, domain.DetectWorkConflictsOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(conflicts) != 1 {
		t.Fatalf("conflicts len=%d want 1: %+v", len(conflicts), conflicts)
	}
	c := conflicts[0]
	if c.TaskA > c.TaskB {
		t.Fatalf("pair not sorted: %+v", c)
	}
	if !strings.Contains(c.Reason, domain.WorkConflictReasonPathOverlap) {
		t.Fatalf("reason=%q want path_overlap", c.Reason)
	}
	if len(c.Paths) == 0 {
		t.Fatalf("paths empty: %+v", c)
	}
}

func TestRedundantSimilarTitleSameGoal(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "shared goal"})
	if err != nil {
		t.Fatal(err)
	}
	taskA, err := svc.CreateTask(ctx, domain.TaskInput{
		Title:  "implement conflict detection advisory",
		GoalID: &g.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	taskB, err := svc.CreateTask(ctx, domain.TaskInput{
		Title:  "implement conflict detection",
		GoalID: &g.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	mustInProgress(t, svc, taskA.ID)
	mustInProgress(t, svc, taskB.ID)

	conflicts, err := svc.DetectWorkConflicts(ctx, domain.DetectWorkConflictsOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(conflicts) != 1 {
		t.Fatalf("conflicts len=%d want 1: %+v", len(conflicts), conflicts)
	}
	if !strings.Contains(conflicts[0].Reason, domain.WorkConflictReasonTitleRedundancy) {
		t.Fatalf("reason=%q want title_redundancy", conflicts[0].Reason)
	}
	if conflicts[0].TaskA == taskA.ID && conflicts[0].TaskB == taskB.ID ||
		conflicts[0].TaskA == taskB.ID && conflicts[0].TaskB == taskA.ID {
		// ok
	} else {
		t.Fatalf("unexpected pair: %+v", conflicts[0])
	}
}

func TestNoConflictWhenTasksDisjoint(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	g1, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "goal one"})
	if err != nil {
		t.Fatal(err)
	}
	g2, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "goal two"})
	if err != nil {
		t.Fatal(err)
	}
	taskA, err := svc.CreateTask(ctx, domain.TaskInput{Title: "alpha wiring", GoalID: &g1.ID})
	if err != nil {
		t.Fatal(err)
	}
	taskB, err := svc.CreateTask(ctx, domain.TaskInput{Title: "beta plumbing", GoalID: &g2.ID})
	if err != nil {
		t.Fatal(err)
	}
	mustInProgress(t, svc, taskA.ID)
	mustInProgress(t, svc, taskB.ID)

	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskA.ID,
		Paths:  []domain.ChangePathInput{{Path: "cmd/trace/a.go"}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskB.ID,
		Paths:  []domain.ChangePathInput{{Path: "internal/domain/b.go"}},
	}); err != nil {
		t.Fatal(err)
	}

	conflicts, err := svc.DetectWorkConflicts(ctx, domain.DetectWorkConflictsOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(conflicts) != 0 {
		t.Fatalf("expected no conflicts: %+v", conflicts)
	}
}
