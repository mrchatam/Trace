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

func analyzersTestdata(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "analyzers", "testdata", name))
	if err != nil {
		t.Fatalf("analyzers testdata %s: %v", name, err)
	}
	return b
}

func indexArchitectureFixture(t *testing.T, st *store.Store, internalSrc []byte) {
	t.Helper()
	ctx := context.Background()
	root := st.ProjectRoot()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/mod\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmdSrc := analyzersTestdata(t, "architecture_cmd_main.go")
	if err := analyzers.IndexFile(ctx, st, "cmd/app/main.go", cmdSrc, analyzers.IndexOptions{}); err != nil {
		t.Fatalf("IndexFile cmd: %v", err)
	}
	if err := analyzers.IndexFile(ctx, st, "internal/lib/lib.go", internalSrc, analyzers.IndexOptions{}); err != nil {
		t.Fatalf("IndexFile internal: %v", err)
	}
}

func TestInvariantFailOnForbiddenLayerImport(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	indexArchitectureFixture(t, st, analyzersTestdata(t, "architecture_internal_imports_cmd.go"))

	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abc1234",
		Paths:     []domain.ChangePathInput{{Path: "internal/lib/lib.go"}},
	}); err != nil {
		t.Fatal(err)
	}

	got, err := svc.CheckArchitecturalInvariants(ctx, task.ID)
	if err != nil {
		t.Fatalf("CheckArchitecturalInvariants: %v", err)
	}
	if got.Passed {
		t.Fatalf("internal→cmd import must fail: %+v", got)
	}
	if len(got.Violations) == 0 {
		t.Fatal("expected at least one violation")
	}
	v := got.Violations[0]
	if v.FromLayer != "internal" || v.ToLayer != "cmd" {
		t.Fatalf("layers: %+v", v)
	}
	if v.ImporterPath != "internal/lib/lib.go" {
		t.Fatalf("importer: %+v", v)
	}
}

func TestInvariantPassWhenNoCrossLayer(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)

	indexArchitectureFixture(t, st, analyzersTestdata(t, "architecture_internal_lib.go"))

	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "def5678",
		Paths: []domain.ChangePathInput{
			{Path: "cmd/app/main.go"},
			{Path: "internal/lib/lib.go"},
		},
	}); err != nil {
		t.Fatal(err)
	}

	got, err := svc.CheckArchitecturalInvariants(ctx, task.ID)
	if err != nil {
		t.Fatalf("CheckArchitecturalInvariants: %v", err)
	}
	if !got.Passed || len(got.Violations) != 0 {
		t.Fatalf("clean change must pass: %+v", got)
	}
}
