package domain

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

// Regression for #75: a bad link after valid entities must roll back the whole
// import (no torn graph) and leave summary.OK false.
func TestImportSeedDocumentAtomicBadLink(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	goalID := uuid.NewString()
	taskID := uuid.NewString()
	missingTask := uuid.NewString()

	doc := SeedDocument{
		Version: 1,
		Goals: []SeedEntity{
			{ID: goalID, Title: "G", Body: "goal"},
		},
		Tasks: []SeedTask{
			{ID: taskID, Title: "T", Body: "task", GoalID: goalID},
		},
		Links: []SeedLink{
			{Rel: "goal-task", From: goalID, To: missingTask},
		},
	}

	summary, err := svc.ImportSeedDocument(ctx, doc)
	if err == nil {
		t.Fatal("expected import error for missing task link")
	}
	if summary.OK {
		t.Fatalf("summary.OK want false on failure, got %+v", summary)
	}
	if len(summary.Created) != 0 {
		t.Fatalf("Created should be cleared on failure, got %+v", summary.Created)
	}

	if _, err := st.GetGoal(goalID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("goal should be rolled back, GetGoal err=%v", err)
	}
	if _, err := st.GetTask(taskID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("task should be rolled back, GetTask err=%v", err)
	}
}

func TestImportSeedDocumentOKStillTrueOnSuccess(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	goalID := uuid.NewString()
	taskID := uuid.NewString()
	doc := SeedDocument{
		Version: 1,
		Goals:   []SeedEntity{{ID: goalID, Title: "G", Body: "goal"}},
		Tasks:   []SeedTask{{ID: taskID, Title: "T", Body: "task", GoalID: goalID}},
		Links:   []SeedLink{{Rel: "goal-task", From: goalID, To: taskID}},
	}
	summary, err := svc.ImportSeedDocument(ctx, doc)
	if err != nil {
		t.Fatalf("import: %v", err)
	}
	if !summary.OK {
		t.Fatalf("summary.OK want true, got %+v", summary)
	}
	if _, err := st.GetGoal(goalID); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetTask(taskID); err != nil {
		t.Fatal(err)
	}
}
