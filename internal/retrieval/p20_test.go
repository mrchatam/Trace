package retrieval_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func TestExactLookupUncertainty(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "task for uncertainty"})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	u, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "zephyruncertainty token",
		Body:     "need API shape",
		Severity: store.UncertaintySeverityBlocking,
		TaskID:   task.ID,
	})
	if err != nil {
		t.Fatalf("CreateUncertainty: %v", err)
	}

	hits, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "uncertainty", EntityID: u.ID})
	if err != nil || len(hits) != 1 {
		t.Fatalf("Exact uncertainty: err=%v hits=%+v", err, hits)
	}
	if hits[0].EntityType != "uncertainty" || hits[0].EntityID != u.ID {
		t.Fatalf("Exact hit: %+v", hits[0])
	}
	if hits[0].Title != "zephyruncertainty token" {
		t.Fatalf("Title: %q", hits[0].Title)
	}
}

func TestExactLookupHypothesis(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	h, err := svc.CreateHypothesis(ctx, domain.HypothesisInput{
		Title: "zephyrhypothesis token",
		Body:  "cache miss explains latency",
	})
	if err != nil {
		t.Fatalf("CreateHypothesis: %v", err)
	}

	hits, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "hypothesis", EntityID: h.ID})
	if err != nil || len(hits) != 1 {
		t.Fatalf("Exact hypothesis: err=%v hits=%+v", err, hits)
	}
	if hits[0].EntityType != "hypothesis" || hits[0].EntityID != h.ID {
		t.Fatalf("Exact hit: %+v", hits[0])
	}
	if hits[0].Title != "zephyrhypothesis token" {
		t.Fatalf("Title: %q", hits[0].Title)
	}
}

func TestWhyUncertaintyIncludesGraphSteps(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "blocked task"})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	u, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "blocking gap",
		Severity: store.UncertaintySeverityBlocking,
		TaskID:   task.ID,
	})
	if err != nil {
		t.Fatalf("CreateUncertainty: %v", err)
	}

	why, err := eng.Why(ctx, "uncertainty", u.ID)
	if err != nil {
		t.Fatalf("Why: %v", err)
	}
	if why.SeedType != "uncertainty" || why.SeedID != u.ID {
		t.Fatalf("Why seed: %+v", why)
	}
	if len(why.Steps) == 0 {
		t.Fatal("Why steps empty")
	}
	if why.Steps[0].EntityType != "uncertainty" || why.Steps[0].EntityID != u.ID {
		t.Fatalf("Why seed step: %+v", why.Steps[0])
	}
}

func TestNormalizeEntityTypeP20Aliases(t *testing.T) {
	if got := retrieval.NormalizeEntityType("outcome"); got != "outcome_result" {
		t.Fatalf("outcome alias: got %q want outcome_result", got)
	}
	if got := retrieval.NormalizeEntityType("plan-change"); got != "plan_change" {
		t.Fatalf("plan-change alias: got %q want plan_change", got)
	}
	for _, canonical := range []string{
		"uncertainty", "hypothesis", "change", "effect", "regression",
		"reflection", "baseline", "outcome_result",
	} {
		if got := retrieval.NormalizeEntityType(canonical); got != canonical {
			t.Fatalf("%q pass-through: got %q", canonical, got)
		}
	}
}
