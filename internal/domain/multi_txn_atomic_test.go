package domain

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/store"
)

// Regression for #89: link/event failure after task upsert must roll back promote.
func TestPromoteBlockingDiscoveryAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	goal, err := svc.CreateGoal(ctx, GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	disc, err := svc.CreateDiscovery(ctx, DiscoveryInput{Title: "block", Severity: SeverityBlocking})
	if err != nil {
		t.Fatal(err)
	}

	injected := errors.New("injected fail after promote task")
	svc.afterPromoteTaskHook = func() error { return injected }

	_, _, err = svc.PromoteBlockingDiscovery(ctx, disc.ID, goal.ID)
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	if _, err := st.GetTask(disc.ID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("orphan task must roll back, GetTask err=%v", err)
	}
	links, err := st.ListLinksFrom(EntityDiscovery, disc.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range links {
		if l.Rel == RelDiscoveryMentionsTask {
			t.Fatalf("unexpected link after rollback: %+v", l)
		}
	}
}

func TestPromoteBlockingDiscoverySucceedsAtomically(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()
	goal, err := svc.CreateGoal(ctx, GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	disc, err := svc.CreateDiscovery(ctx, DiscoveryInput{Title: "block", Severity: SeverityBlocking})
	if err != nil {
		t.Fatal(err)
	}
	taskID, inserted, err := svc.PromoteBlockingDiscovery(ctx, disc.ID, goal.ID)
	if err != nil || !inserted || taskID == "" {
		t.Fatalf("promote: id=%s inserted=%v err=%v", taskID, inserted, err)
	}
	if _, err := st.GetTask(taskID); err != nil {
		t.Fatal(err)
	}
}

// Regression for #90: deliberation upsert + event must share one tx.
func TestApplyDeliberationTransitionAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()
	goal, err := svc.CreateGoal(ctx, GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, TaskInput{Title: "t", GoalID: &goal.ID})
	if err != nil {
		t.Fatal(err)
	}

	injected := errors.New("injected fail after deliberation upsert")
	svc.afterDeliberationUpsertHook = func() error { return injected }

	_, _, err = svc.ApplyDeliberationTransition(ctx, task.ID, goal.ID, deliberation.PolicyInputs{
		BlockingUncertaintyCount: 1, PlanExists: true, PlanCritiqued: true,
	})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	if _, err := st.GetDeliberationState(task.ID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("deliberation state must roll back, err=%v", err)
	}
	evs, err := st.ListEventsByEntity(EntityTask, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type == EventDeliberationTransition {
			t.Fatalf("unexpected deliberation.transition event: %+v", e)
		}
	}
}

// Regression for #91: delete+reinsert requirements must be atomic.
func TestUpsertHarnessAgentAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	agentID := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	if _, err := svc.UpsertHarnessAgent(ctx, HarnessAgentInput{
		ID: agentID, Slug: "agent:atomic", Title: "A",
		SubagentType: "generalPurpose", DeliberationPhases: `["ORIENT"]`,
		Requirements: []string{"skill:one", "skill:two"},
	}); err != nil {
		t.Fatal(err)
	}
	reqs, err := st.ListHarnessAgentRequirements(agentID)
	if err != nil || len(reqs) != 2 {
		t.Fatalf("seed reqs: %+v err=%v", reqs, err)
	}

	injected := errors.New("injected fail after delete requirements")
	svc.afterHarnessDeleteHook = func() error { return injected }

	_, err = svc.UpsertHarnessAgent(ctx, HarnessAgentInput{
		ID: agentID, Slug: "agent:atomic", Title: "A",
		SubagentType: "generalPurpose", DeliberationPhases: `["ORIENT"]`,
		Requirements: []string{"skill:three"},
	})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	reqs, err = st.ListHarnessAgentRequirements(agentID)
	if err != nil {
		t.Fatal(err)
	}
	if len(reqs) != 2 {
		t.Fatalf("requirements must roll back to 2, got %+v", reqs)
	}
}

// Regression for #92: SetReviewResult / CreateReview / MarkStale entity+event atomic.
func TestSetReviewResultAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()
	rev, err := svc.CreateReview(ctx, ReviewInput{Title: "r"})
	if err != nil {
		t.Fatal(err)
	}

	injected := errors.New("injected fail after review upsert")
	svc.afterReviewMutateHook = func() error { return injected }

	err = svc.SetReviewResult(ctx, rev.ID, store.ReviewResultPass, ReviewResultOptions{
		Actor: "a", Reason: "ok",
	})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	got, err := st.GetReview(rev.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Result != store.ReviewResultOpen {
		t.Fatalf("result=%q want OPEN (rollback)", got.Result)
	}
}

func TestCreateReviewAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	id := "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	injected := errors.New("injected fail after create review upsert")
	svc.afterReviewMutateHook = func() error { return injected }

	_, err = svc.CreateReview(ctx, ReviewInput{ID: id, Title: "new"})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	if _, err := st.GetReview(id); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("review must roll back, err=%v", err)
	}
}

func TestMarkStaleAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}

	injected := errors.New("injected fail after mark stale upsert")
	svc.afterMarkStaleUpsertHook = func() error { return injected }

	err = svc.MarkStale(ctx, EntityGoal, g.ID, "stale reason")
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	got, err := st.GetGoal(g.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status == store.StatusStale {
		t.Fatalf("status must not stay STALE after rollback, got %q", got.Status)
	}
	evs, err := st.ListEventsByEntity(EntityGoal, g.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type == "entity.stale" {
			t.Fatalf("unexpected entity.stale: %+v", e)
		}
	}
}
