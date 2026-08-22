package domain_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func indexValidatesFixture(t *testing.T, st *store.Store) (fooSymbolID string) {
	t.Helper()
	ctx := context.Background()
	root := st.ProjectRoot()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/mod\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fooGo := analyzersTestdata(t, "foo.go")
	fooTest := analyzersTestdata(t, "foo_test.go")
	if err := analyzers.IndexFile(ctx, st, "pkg/foo.go", fooGo, analyzers.IndexOptions{}); err != nil {
		t.Fatalf("IndexFile foo.go: %v", err)
	}
	if err := analyzers.IndexFile(ctx, st, "pkg/foo_test.go", fooTest, analyzers.IndexOptions{}); err != nil {
		t.Fatalf("IndexFile foo_test.go: %v", err)
	}
	prod, err := st.ListSymbolsByPath("pkg/foo.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range prod {
		if s.Name == "Foo" {
			return s.ID
		}
	}
	t.Fatalf("Foo symbol missing from pkg/foo.go: %v", prod)
	return ""
}

func TestTestsVerifyingQuery(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	fooID := indexValidatesFixture(t, st)

	bySymbol, err := svc.ListTestsValidatingSymbol(ctx, fooID)
	if err != nil {
		t.Fatalf("ListTestsValidatingSymbol: %v", err)
	}
	if len(bySymbol) != 1 {
		t.Fatalf("by symbol: %+v", bySymbol)
	}
	if bySymbol[0].TestFilePath != "pkg/foo_test.go" || bySymbol[0].TestSymbolName != "TestFoo" {
		t.Fatalf("by symbol row: %+v", bySymbol[0])
	}
	if bySymbol[0].EdgeProvenance != store.ImportProvenanceInferred {
		t.Fatalf("provenance=%q", bySymbol[0].EdgeProvenance)
	}

	byFile, err := svc.ListTestsValidatingFile(ctx, "pkg/foo.go")
	if err != nil {
		t.Fatalf("ListTestsValidatingFile: %v", err)
	}
	if len(byFile) != 1 || byFile[0].TestSymbolName != "TestFoo" {
		t.Fatalf("by file: %+v", byFile)
	}

	if _, err := svc.ListTestsValidatingSymbol(ctx, "00000000-0000-0000-0000-000000000000"); err == nil {
		t.Fatal("unknown symbol must fail closed")
	}
	if _, err := svc.ListTestsValidatingFile(ctx, "missing.go"); err == nil {
		t.Fatal("unknown file must fail closed")
	}
}

func TestOutcomesFailedAndWorked(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	failRow, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "TestAlpha", TestStatus: store.TestStatusFail, Summary: "boom",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID: task.ID, TestName: "TestBeta", TestStatus: store.TestStatusPass,
	}); err != nil {
		t.Fatal(err)
	}
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "abc1234",
		Paths: []domain.ChangePathInput{{Path: "a.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	imp, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
		ChangeID: chg.ID, TaskID: task.ID, Summary: "cache hit",
	})
	if err != nil {
		t.Fatal(err)
	}

	failed, err := svc.ListFailedOutcomes(ctx, domain.EvidenceQueryOpts{TaskID: task.ID})
	if err != nil {
		t.Fatalf("ListFailedOutcomes: %v", err)
	}
	if len(failed) != 1 || failed[0].ID != failRow.ID || failed[0].TestStatus != store.TestStatusFail {
		t.Fatalf("failed: %+v", failed)
	}

	worked, err := svc.ListWorkedApproaches(ctx, domain.EvidenceQueryOpts{TaskID: task.ID})
	if err != nil {
		t.Fatalf("ListWorkedApproaches: %v", err)
	}
	if len(worked) < 2 {
		t.Fatalf("worked: %+v", worked)
	}
	kinds := map[string]bool{}
	for _, w := range worked {
		kinds[w.Kind] = true
		if w.ID == imp.ID && w.Kind != "improvement" {
			t.Fatalf("improvement kind: %+v", w)
		}
		if w.Kind == "test_pass" && w.TestName != "TestBeta" {
			t.Fatalf("pass row: %+v", w)
		}
	}
	if !kinds["improvement"] || !kinds["test_pass"] {
		t.Fatalf("kinds=%v worked=%+v", kinds, worked)
	}
}

func TestRegressionsListQueryable(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	out := mustEvalWithRegression(t, svc, task.ID)
	reg, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID, GitCommit: "def5678",
		Paths: []domain.ChangePathInput{{Path: "b.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.AssociateRegressionWithChange(ctx, reg.ID, chg.ID); err != nil {
		t.Fatal(err)
	}

	byChange, err := svc.ListRegressions(ctx, domain.EvidenceQueryOpts{ChangeID: chg.ID})
	if err != nil {
		t.Fatalf("ListRegressions by change: %v", err)
	}
	if len(byChange) != 1 || byChange[0].ID != reg.ID {
		t.Fatalf("by change: %+v", byChange)
	}
	if len(byChange[0].ChangeIDs) != 1 || byChange[0].ChangeIDs[0] != chg.ID {
		t.Fatalf("change ids: %+v", byChange[0].ChangeIDs)
	}

	byTask, err := svc.ListRegressions(ctx, domain.EvidenceQueryOpts{TaskID: task.ID})
	if err != nil {
		t.Fatalf("ListRegressions by task: %v", err)
	}
	if len(byTask) < 1 {
		t.Fatalf("by task empty")
	}
	found := false
	for _, r := range byTask {
		if r.ID == reg.ID {
			found = true
			if r.Attribution == "" || r.Dimension == "" {
				t.Fatalf("missing fields: %+v", r)
			}
		}
	}
	if !found {
		t.Fatalf("regression not in task list: %+v", byTask)
	}
}
