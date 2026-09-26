package domain

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/store"
)

func TestCreateScopeAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	id := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	injected := errors.New("injected fail after create scope upsert")
	svc.afterCreateMutateHook = func() error { return injected }

	_, err = svc.CreateScope(ctx, ScopeInput{ID: id, Slug: "auth", Title: "Auth", Kind: store.ScopeKindFeature})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	if _, err := st.GetScope(id); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("scope must roll back, err=%v", err)
	}
	evs, err := st.ListEventsByEntity(EntityScope, id)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type == EventEntityCreated {
			t.Fatalf("unexpected entity.created: %+v", e)
		}
	}
}

func TestLinkScopeMemberAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	scope, err := svc.CreateScope(ctx, ScopeInput{Slug: "auth", Title: "Auth", Kind: store.ScopeKindFeature})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, TaskInput{Title: "fix login"})
	if err != nil {
		t.Fatal(err)
	}

	injected := errors.New("injected fail after link upsert")
	svc.afterLinkMutateHook = func() error { return injected }

	err = svc.LinkScopeMember(ctx, task.ID, scope.ID, LinkMeta{})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	links, err := st.ListLinksFrom(EntityTask, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range links {
		if l.Rel == RelScopeMember {
			t.Fatalf("unexpected link after rollback: %+v", l)
		}
	}
	evs, err := st.ListEventsByEntity(EntityTask, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type == EventEntityLinked {
			t.Fatalf("unexpected entity.linked: %+v", e)
		}
	}
}

func TestLinkAPIContractAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	task1, err := svc.CreateTask(ctx, TaskInput{Title: "task 1"})
	if err != nil {
		t.Fatal(err)
	}
	task2, err := svc.CreateTask(ctx, TaskInput{Title: "task 2"})
	if err != nil {
		t.Fatal(err)
	}

	injected := errors.New("injected fail after link upsert")
	svc.afterLinkMutateHook = func() error { return injected }

	err = svc.LinkAPIContract(ctx, task1.ID, task2.ID, LinkMeta{})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	links, err := st.ListLinksFrom(EntityTask, task1.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range links {
		if l.Rel == RelAPIContract {
			t.Fatalf("unexpected link after rollback: %+v", l)
		}
	}
	evs, err := st.ListEventsByEntity(EntityTask, task1.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type == EventEntityLinked {
			t.Fatalf("unexpected entity.linked: %+v", e)
		}
	}
}

func TestLinkImplementsAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	task1, err := svc.CreateTask(ctx, TaskInput{Title: "task 1"})
	if err != nil {
		t.Fatal(err)
	}
	task2, err := svc.CreateTask(ctx, TaskInput{Title: "task 2"})
	if err != nil {
		t.Fatal(err)
	}

	injected := errors.New("injected fail after link upsert")
	svc.afterLinkMutateHook = func() error { return injected }

	err = svc.LinkImplements(ctx, task1.ID, task2.ID, LinkMeta{})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	links, err := st.ListLinksFrom(EntityTask, task1.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range links {
		if l.Rel == RelImplements {
			t.Fatalf("unexpected link after rollback: %+v", l)
		}
	}
	evs, err := st.ListEventsByEntity(EntityTask, task1.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type == EventEntityLinked {
			t.Fatalf("unexpected entity.linked: %+v", e)
		}
	}
}

func TestLinkBlocksAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	task1, err := svc.CreateTask(ctx, TaskInput{Title: "task 1"})
	if err != nil {
		t.Fatal(err)
	}
	task2, err := svc.CreateTask(ctx, TaskInput{Title: "task 2"})
	if err != nil {
		t.Fatal(err)
	}

	injected := errors.New("injected fail after link upsert")
	svc.afterLinkMutateHook = func() error { return injected }

	err = svc.LinkBlocks(ctx, task1.ID, task2.ID, LinkMeta{})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	links, err := st.ListLinksFrom(EntityTask, task1.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range links {
		if l.Rel == RelBlocks {
			t.Fatalf("unexpected link after rollback: %+v", l)
		}
	}
	evs, err := st.ListEventsByEntity(EntityTask, task1.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type == EventEntityLinked {
			t.Fatalf("unexpected entity.linked: %+v", e)
		}
	}
}

// TestInferScopesEmitsLinkedEvent verifies that InferScopes emits entity.linked events
// for inserted scope_member links.
func TestInferScopesEmitsLinkedEvent(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	scope, err := svc.CreateScope(ctx, ScopeInput{Slug: "auth", Title: "Auth", Kind: store.ScopeKindFeature})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, TaskInput{Title: "fix login session"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.InferScopes(ctx, InferOptions{DryRun: false})
	if err != nil {
		t.Fatal(err)
	}

	evs, err := st.ListEventsByEntity(EntityTask, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range evs {
		if e.Type == EventEntityLinked {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected entity.linked event for scope_member, got none")
	}

	links, err := st.ListLinksByRel(RelScopeMember)
	if err != nil {
		t.Fatal(err)
	}
	foundLink := false
	for _, l := range links {
		if l.FromID == task.ID && l.ToID == scope.ID && l.Rel == RelScopeMember {
			foundLink = true
			break
		}
	}
	if !foundLink {
		t.Fatalf("expected scope_member link from task to scope, got none")
	}
}

// TestInferScopesAtomicOnMidpathFailure verifies that when afterLinkMutateHook fails,
// all inserted scope_member links are rolled back.
func TestInferScopesAtomicOnMidpathFailure(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := New(st)
	ctx := context.Background()

	scope, err := svc.CreateScope(ctx, ScopeInput{Slug: "auth", Title: "Auth", Kind: store.ScopeKindFeature})
	if err != nil {
		t.Fatal(err)
	}
	task1, err := svc.CreateTask(ctx, TaskInput{Title: "fix login session"})
	if err != nil {
		t.Fatal(err)
	}
	task2, err := svc.CreateTask(ctx, TaskInput{Title: "fix auth token"})
	if err != nil {
		t.Fatal(err)
	}

	injected := errors.New("injected fail after first link upsert")
	svc.afterLinkMutateHook = func() error { return injected }

	rep, err := svc.InferScopes(ctx, InferOptions{DryRun: false})
	if !errors.Is(err, injected) {
		t.Fatalf("want injected, got %v", err)
	}
	if rep.Inserted != 0 {
		t.Fatalf("report.Inserted must be 0 after rollback, got %d", rep.Inserted)
	}

	links, err := st.ListLinksByRel(RelScopeMember)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range links {
		if l.Rel == RelScopeMember && (l.FromID == task1.ID || l.FromID == task2.ID) && l.ToID == scope.ID {
			t.Fatalf("unexpected scope_member link after rollback: %+v", l)
		}
	}

	for _, id := range []string{task1.ID, task2.ID} {
		evs, err := st.ListEventsByEntity(EntityTask, id)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range evs {
			if e.Type == EventEntityLinked {
				t.Fatalf("unexpected entity.linked for %s after rollback: %+v", id, e)
			}
		}
	}
}
