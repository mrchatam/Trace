package domain_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestInferChangeKind(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{"empty", nil, "seg:unknown"},
		{"lex smallest cmd", []string{"internal/domain/foo.go", "cmd/trace/main.go"}, "seg:cmd"},
		{"single segment", []string{"main.go"}, "seg:main.go"},
		{"internal only", []string{"internal/store/patterns.go"}, "seg:internal"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := domain.InferChangeKind(tc.paths); got != tc.want {
				t.Fatalf("InferChangeKind(%v) = %q want %q", tc.paths, got, tc.want)
			}
		})
	}
}

func TestClassifyChangeOutcomePriority(t *testing.T) {
	got := domain.ClassifyChangeOutcome(domain.ChangeOutcomeSignals{
		HasRegression:      true,
		HasImprovement:     true,
		HasEffectSupported: true,
		HasTestPass:        true,
	})
	if got != domain.OutcomeKindRegression {
		t.Fatalf("regression wins: got %q", got)
	}
	got = domain.ClassifyChangeOutcome(domain.ChangeOutcomeSignals{
		HasEffectContradicted: true,
		HasImprovement:        true,
	})
	if got != domain.OutcomeKindEffectContradicted {
		t.Fatalf("contradicted wins: got %q", got)
	}
	got = domain.ClassifyChangeOutcome(domain.ChangeOutcomeSignals{
		HasImprovement:     true,
		HasEffectSupported: true,
	})
	if got != domain.OutcomeKindImprovement {
		t.Fatalf("improvement wins: got %q", got)
	}
	got = domain.ClassifyChangeOutcome(domain.ChangeOutcomeSignals{
		HasEffectSupported: true,
		HasTestPass:        true,
	})
	if got != domain.OutcomeKindEffectSupported {
		t.Fatalf("effect_supported wins: got %q", got)
	}
	got = domain.ClassifyChangeOutcome(domain.ChangeOutcomeSignals{HasTestFail: true, HasTestPass: true})
	if got != domain.OutcomeKindTestFail {
		t.Fatalf("test_fail wins: got %q", got)
	}
	got = domain.ClassifyChangeOutcome(domain.ChangeOutcomeSignals{HasTestPass: true})
	if got != domain.OutcomeKindTestPass {
		t.Fatalf("test_pass: got %q", got)
	}
	got = domain.ClassifyChangeOutcome(domain.ChangeOutcomeSignals{})
	if got != domain.OutcomeKindNeutral {
		t.Fatalf("neutral: got %q", got)
	}
}

func TestPatternCountsFromChangesAndOutcomes(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	chgImprove, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "internal improvement",
		Paths:  []domain.ChangePathInput{{Path: "internal/a.go", Status: "M"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
		ChangeID: chgImprove.ID,
		TaskID:   task.ID,
		Summary:  "faster path",
	}); err != nil {
		t.Fatal(err)
	}

	chgRegression, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "internal regression",
		Paths:  []domain.ChangePathInput{{Path: "internal/b.go", Status: "M"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := mustEvalWithRegression(t, svc, task.ID)
	reg, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID,
		TaskID:    task.ID,
		Summary:   "score dropped",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.AssociateRegressionWithChange(ctx, reg.ID, chgRegression.ID); err != nil {
		t.Fatal(err)
	}

	chgSupported, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "cmd supported",
		Paths:  []domain.ChangePathInput{{Path: "cmd/trace/main.go", Status: "M"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordExpectedEffect(ctx, chgSupported.ID, domain.ExpectedEffectInput{
		Dimension: "latency", Expected: "lower",
	}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.RecordActualEffect(ctx, chgSupported.ID, domain.RecordActualEffectInput{
		Dimension: "latency", Actual: "lower", Comparison: store.EffectComparisonSupported,
	}); err != nil {
		t.Fatal(err)
	}

	n, err := svc.RefreshChangePatterns(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("RefreshChangePatterns rows: got %d want 3", n)
	}

	patterns, err := svc.ListChangePatterns(ctx, 64)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]struct {
		pos int
		neg int
	}{
		"seg:internal|" + domain.OutcomeKindImprovement: {pos: 1, neg: 0},
		"seg:internal|" + domain.OutcomeKindRegression:  {pos: 0, neg: 1},
		"seg:cmd|" + domain.OutcomeKindEffectSupported:  {pos: 1, neg: 0},
	}
	for _, p := range patterns {
		key := p.ChangeKind + "|" + p.OutcomeKind
		exp, ok := want[key]
		if !ok {
			t.Fatalf("unexpected pattern row: %+v", p)
		}
		if p.CountPositive != exp.pos || p.CountNegative != exp.neg {
			t.Fatalf("counts for %s: %+v want pos=%d neg=%d", key, p, exp.pos, exp.neg)
		}
		delete(want, key)
	}
	if len(want) != 0 {
		t.Fatalf("missing pattern rows: %+v got %+v", want, patterns)
	}
}

func TestQuerySimilarChanges(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	first, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "older internal",
		Paths:  []domain.ChangePathInput{{Path: "internal/alpha.go", Status: "M"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "newer internal",
		Paths:  []domain.ChangePathInput{{Path: "internal/beta.go", Status: "A"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordExpectedEffect(ctx, second.ID, domain.ExpectedEffectInput{
		Dimension: "tests", Expected: "pass",
	}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.RecordActualEffect(ctx, second.ID, domain.RecordActualEffectInput{
		Dimension: "tests", Actual: "pass", Comparison: store.EffectComparisonSupported,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "cmd only",
		Paths:  []domain.ChangePathInput{{Path: "cmd/main.go", Status: "M"}},
	}); err != nil {
		t.Fatal(err)
	}

	byKind, err := svc.QuerySimilarChanges(ctx, domain.SimilarChangesOpts{ChangeKind: "seg:internal"})
	if err != nil {
		t.Fatal(err)
	}
	if len(byKind.Changes) != 2 {
		t.Fatalf("kind filter: got %d changes", len(byKind.Changes))
	}
	if byKind.Changes[0].ID != second.ID || byKind.Changes[1].ID != first.ID {
		t.Fatalf("newest-first: %+v", byKind.Changes)
	}
	for _, row := range byKind.Changes {
		if row.ChangeKind != "seg:internal" {
			t.Fatalf("change_kind: %+v", row)
		}
	}
	if byKind.Changes[0].OutcomeKind != domain.OutcomeKindEffectSupported {
		t.Fatalf("effects summary outcome: %+v", byKind.Changes[0])
	}
	if len(byKind.Changes[0].Effects) != 1 || byKind.Changes[0].Effects[0].Dimension != "tests" {
		t.Fatalf("effects: %+v", byKind.Changes[0].Effects)
	}

	byPath, err := svc.QuerySimilarChanges(ctx, domain.SimilarChangesOpts{PathPrefix: "internal/"})
	if err != nil {
		t.Fatal(err)
	}
	if len(byPath.Changes) != 2 {
		t.Fatalf("path prefix: got %d", len(byPath.Changes))
	}
	for _, row := range byPath.Changes {
		found := false
		for _, p := range row.Paths {
			if strings.HasPrefix(p, "internal/") {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("path prefix mismatch: %+v", row)
		}
	}
}

func TestQuerySimilarChangesSameSecondTimestamps(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	const sameSecond = "2026-08-18T12:00:00Z"
	// UUID lex order would put zzzz before aaaa; rowid preserves insertion order.
	older, err := st.UpsertChange(store.Change{
		ID:         "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz",
		TaskID:     task.ID,
		Reason:     "older internal",
		Status:     store.ChangeStatusOpen,
		SourceType: "USER_ASSERTED",
		CreatedAt:  sameSecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.InsertChangePath(store.ChangePath{
		ChangeID: older.ID, Path: "internal/alpha.go", Status: "M",
	}); err != nil {
		t.Fatal(err)
	}
	newer, err := st.UpsertChange(store.Change{
		ID:         "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		TaskID:     task.ID,
		Reason:     "newer internal",
		Status:     store.ChangeStatusOpen,
		SourceType: "USER_ASSERTED",
		CreatedAt:  sameSecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.InsertChangePath(store.ChangePath{
		ChangeID: newer.ID, Path: "internal/beta.go", Status: "A",
	}); err != nil {
		t.Fatal(err)
	}

	byKind, err := svc.QuerySimilarChanges(ctx, domain.SimilarChangesOpts{ChangeKind: "seg:internal"})
	if err != nil {
		t.Fatal(err)
	}
	if len(byKind.Changes) != 2 {
		t.Fatalf("kind filter: got %d changes", len(byKind.Changes))
	}
	if byKind.Changes[0].ID != newer.ID || byKind.Changes[1].ID != older.ID {
		t.Fatalf("newest-first by rowid: %+v", byKind.Changes)
	}

	byPath, err := svc.QuerySimilarChanges(ctx, domain.SimilarChangesOpts{PathPrefix: "internal/"})
	if err != nil {
		t.Fatal(err)
	}
	if len(byPath.Changes) != 2 {
		t.Fatalf("path prefix: got %d", len(byPath.Changes))
	}
	if byPath.Changes[0].ID != newer.ID || byPath.Changes[1].ID != older.ID {
		t.Fatalf("path prefix newest-first by rowid: %+v", byPath.Changes)
	}
}

func TestQuerySimilarChangesFailClosed(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	_, err := svc.QuerySimilarChanges(ctx, domain.SimilarChangesOpts{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	var v *domain.ErrValidation
	if !errors.As(err, &v) {
		t.Fatalf("want ErrValidation, got %v", err)
	}

	_, err = svc.QuerySimilarChanges(ctx, domain.SimilarChangesOpts{
		PathPrefix: "internal/",
		ChangeKind: "seg:internal",
	})
	if err == nil {
		t.Fatal("expected mutual exclusion error")
	}
}
