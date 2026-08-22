package domain_test

import (
	"context"
	"errors"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestCreateExperimentLinksOutcome(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	exp, err := svc.CreateExperiment(ctx, domain.ExperimentInput{
		TaskID:            task.ID,
		Label:             "baseline vs candidate A",
		HypothesisSummary: "candidate A improves latency",
	})
	if err != nil {
		t.Fatalf("CreateExperiment: %v", err)
	}
	if exp.Status != store.ExperimentStatusPlanned || exp.OutcomeResultID != "" {
		t.Fatalf("defaults: %+v", exp)
	}

	outcome, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:     task.ID,
		TestName:   "TestCandidateA",
		TestStatus: store.TestStatusPass,
	})
	if err != nil {
		t.Fatalf("RecordTestOutcome: %v", err)
	}

	if err := svc.LinkExperimentOutcome(ctx, exp.ID, outcome.ID); err != nil {
		t.Fatalf("LinkExperimentOutcome: %v", err)
	}
	got, err := svc.GetExperiment(ctx, exp.ID)
	if err != nil {
		t.Fatalf("GetExperiment: %v", err)
	}
	if got.OutcomeResultID != outcome.ID {
		t.Fatalf("outcome link: got %q want %q", got.OutcomeResultID, outcome.ID)
	}
}

func TestExperimentStatusLifecycle(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	exp, err := svc.CreateExperiment(ctx, domain.ExperimentInput{
		TaskID: task.ID,
		Label:  "lifecycle",
	})
	if err != nil {
		t.Fatalf("CreateExperiment: %v", err)
	}

	for _, status := range []string{
		store.ExperimentStatusRunning,
		store.ExperimentStatusCompleted,
	} {
		if err := svc.SetExperimentStatus(ctx, exp.ID, status); err != nil {
			t.Fatalf("SetExperimentStatus(%q): %v", status, err)
		}
		got, err := svc.GetExperiment(ctx, exp.ID)
		if err != nil || got.Status != status {
			t.Fatalf("after %q: %+v err=%v", status, got, err)
		}
	}

	if err := svc.SetExperimentStatus(ctx, exp.ID, store.ExperimentStatusCompleted); err != nil {
		t.Fatalf("idempotent completed: %v", err)
	}

	if err := svc.SetExperimentStatus(ctx, exp.ID, store.ExperimentStatusRunning); err == nil {
		t.Fatal("backward transition must fail closed")
	} else {
		var ve *domain.ErrValidation
		if !errors.As(err, &ve) {
			t.Fatalf("backward transition: want ErrValidation, got %v", err)
		}
	}

	if err := svc.SetExperimentStatus(ctx, exp.ID, "bogus"); err == nil {
		t.Fatal("bogus status must fail closed")
	} else {
		var ve *domain.ErrValidation
		if !errors.As(err, &ve) {
			t.Fatalf("bogus status: want ErrValidation, got %v", err)
		}
	}
}

func TestNoExperimentRunnerInvoked(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	experimentsPath := filepath.Join(filepath.Dir(file), "experiments.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, experimentsPath, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse experiments.go: %v", err)
	}
	for _, imp := range f.Imports {
		if imp.Path.Value == `"os/exec"` {
			t.Fatal("experiments.go must not import os/exec")
		}
	}

	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	exp, err := svc.CreateExperiment(ctx, domain.ExperimentInput{TaskID: task.ID, Label: "no-runner"})
	if err != nil {
		t.Fatalf("CreateExperiment: %v", err)
	}
	outcome, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "TestNoRunner", TestStatus: store.TestStatusPass,
	})
	if err != nil {
		t.Fatalf("RecordTestOutcome: %v", err)
	}
	if err := svc.LinkExperimentOutcome(ctx, exp.ID, outcome.ID); err != nil {
		t.Fatalf("LinkExperimentOutcome: %v", err)
	}
}
