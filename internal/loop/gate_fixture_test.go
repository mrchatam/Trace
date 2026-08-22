package loop_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

func importFeetExportFixture(t *testing.T) (*store.Store, *planner.Service, *domain.Service) {
	t.Helper()
	st, psvc, dsvc := openLoopTestStore(t)
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	path := filepath.Join(filepath.Dir(file), "testdata", "feet-export-min.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var doc domain.SeedDocument
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}
	if _, err := dsvc.ImportSeedDocument(context.Background(), doc); err != nil {
		t.Fatalf("import fixture: %v", err)
	}
	return st, psvc, dsvc
}

func seedTerminalDoneTask(t *testing.T, st *store.Store, taskID string) {
	t.Helper()
	task, err := st.GetTask(taskID)
	if err != nil {
		t.Fatal(err)
	}
	task.WorkState = store.WorkStateDone
	if _, err := st.UpsertTask(task); err != nil {
		t.Fatal(err)
	}
}
