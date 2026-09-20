package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/store"
)

// Regression for #74: failure after PASS invalidation must not leave UNCERTAIN
// reviews while the task remains DONE (invalidate + UpsertTask + AppendEvent
// share one WithTx).
func TestTransitionTaskReopenAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, TaskInput{Title: "reopen-atomic"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, TransitionOptions{
		Actor: "a", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
	rev, err := svc.CreateReview(ctx, ReviewInput{Title: "ship-it"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, rev.ID, task.ID, LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, rev.ID, store.ReviewResultPass, ReviewResultOptions{
		Actor: "a", Reason: "looks good",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, TransitionOptions{
		Actor: "a", Reason: "done", AllowOperatorDone: true,
	}); err != nil {
		t.Fatal(err)
	}

	injected := errors.New("injected fail after invalidate")
	svc.afterPassInvalidateHook = func() error { return injected }

	err = svc.TransitionTask(ctx, task.ID, store.WorkStatePending, TransitionOptions{
		Actor: "a", Reason: "reopen",
	})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected error, got %v", err)
	}

	gotTask, err := st.GetTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotTask.WorkState != store.WorkStateDone {
		t.Fatalf("task work_state=%q want DONE (rollback)", gotTask.WorkState)
	}
	gotRev, err := st.GetReview(rev.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotRev.Result != store.ReviewResultPass {
		t.Fatalf("review result=%q want PASS (rollback; must not leave UNCERTAIN with DONE task)", gotRev.Result)
	}
}

func TestTransitionTaskReopenSucceedsAtomically(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, TaskInput{Title: "reopen-ok"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, TransitionOptions{
		Actor: "a", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
	rev, err := svc.CreateReview(ctx, ReviewInput{Title: "ok"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, rev.ID, task.ID, LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, rev.ID, store.ReviewResultPass, ReviewResultOptions{
		Actor: "a", Reason: "pass",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, TransitionOptions{
		Actor: "a", Reason: "done", AllowOperatorDone: true,
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStatePending, TransitionOptions{
		Actor: "a", Reason: "reopen",
	}); err != nil {
		t.Fatalf("reopen: %v", err)
	}
	gotTask, err := st.GetTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotTask.WorkState != store.WorkStatePending {
		t.Fatalf("work_state=%q want PENDING", gotTask.WorkState)
	}
	gotRev, err := st.GetReview(rev.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotRev.Result != store.ReviewResultUncertain {
		t.Fatalf("review result=%q want UNCERTAIN", gotRev.Result)
	}
}
