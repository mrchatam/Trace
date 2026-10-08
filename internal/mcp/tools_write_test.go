package mcp_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	tracemcp "github.com/mrchatam/Trace/internal/mcp"
	"github.com/mrchatam/Trace/internal/store"
)

// writeFixture creates an initialized project and returns its root and context.
func writeFixture(t *testing.T) (string, context.Context) {
	t.Helper()
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	return dir, context.Background()
}

func addServer(t *testing.T, dir string) *tracemcp.Server {
	t.Helper()
	return tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
}

// TestAddRejectsExistingIDWithoutReplace is the #150 proof: the second write
// used to silently destroy the first row while reporting ok:true.
func TestAddRejectsExistingIDWithoutReplace(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	id := "11111111-2222-3333-4444-555555555555"

	first, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{
		Kind: "task", Title: "original title", ID: id,
	})
	if err != nil {
		t.Fatalf("first add: %v", err)
	}
	if got := addCreated(t, tracemcp.ResultText(t, first)); !got {
		t.Fatal("first add must report created=true")
	}

	_, _, err = tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{
		Kind: "task", Title: "overwritten title", ID: id,
	})
	if err == nil {
		t.Fatal("want a refusal when the id already exists")
	}
	if !strings.Contains(err.Error(), "already exists") || !strings.Contains(err.Error(), "replace=true") {
		t.Fatalf("error must name the remedy, got %v", err)
	}

	// The original row must be untouched.
	st := mustOpenStore(t, dir)
	task, err := st.GetTask(id)
	if err != nil {
		t.Fatal(err)
	}
	if task.Title != "original title" {
		t.Fatalf("row was mutated: %q", task.Title)
	}
}

// TestAddReplaceTrueReportsCreatedFalse proves the opt-in overwrite works and
// is honestly labelled.
func TestAddReplaceTrueReportsCreatedFalse(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	id := "22222222-3333-4444-5555-666666666666"

	if _, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{
		Kind: "task", Title: "first", ID: id,
	}); err != nil {
		t.Fatal(err)
	}
	res, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{
		Kind: "task", Title: "second", ID: id, Replace: true,
	})
	if err != nil {
		t.Fatalf("replace=true must be honoured: %v", err)
	}
	if got := addCreated(t, tracemcp.ResultText(t, res)); got {
		t.Fatal("replace=true must report created=false")
	}

	st := mustOpenStore(t, dir)
	task, err := st.GetTask(id)
	if err != nil {
		t.Fatal(err)
	}
	if task.Title != "second" {
		t.Fatalf("replace did not overwrite: %q", task.Title)
	}
}

// TestAddRejectsNonUUIDID is the #151 proof for the id half.
func TestAddRejectsNonUUIDID(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	_, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{
		Kind: "task", Title: "dangling", ID: "not-a-uuid",
	})
	if err == nil || !strings.Contains(err.Error(), "id must be a UUID") {
		t.Fatalf("want a UUID validation error, got %v", err)
	}

	st := mustOpenStore(t, dir)
	rows, err := st.ListTasksFiltered(store.TaskListFilter{Limit: 50})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Tasks) != 0 {
		t.Fatalf("a rejected add must not write: %+v", rows.Tasks)
	}
}

// TestAddRejectsUnresolvableGoalID proves tasks cannot gain a dangling goal_id
// (tasks.goal_id has no FK).
func TestAddRejectsUnresolvableGoalID(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	_, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{
		Kind: "task", Title: "orphan", GoalID: "00000000-0000-0000-0000-0000000000ff",
	})
	if err == nil || !strings.Contains(err.Error(), "goal_id") {
		t.Fatalf("want a goal_id resolution error, got %v", err)
	}

	st := mustOpenStore(t, dir)
	rows, _ := st.ListTasksFiltered(store.TaskListFilter{Limit: 50})
	if len(rows.Tasks) != 0 {
		t.Fatalf("a rejected add must not write: %+v", rows.Tasks)
	}
}

// TestAddAcceptsResolvableGoalID is the keeper for the happy path.
func TestAddAcceptsResolvableGoalID(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	goalRes, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{Kind: "goal", Title: "parent"})
	if err != nil {
		t.Fatal(err)
	}
	var goal struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(tracemcp.ResultText(t, goalRes)), &goal); err != nil {
		t.Fatal(err)
	}
	if _, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{
		Kind: "task", Title: "child", GoalID: goal.ID,
	}); err != nil {
		t.Fatalf("a resolvable goal_id must be accepted: %v", err)
	}
}

// TestTasksRejectsInvalidWorkState is the #154 proof: an invalid filter used to
// return an authoritative-looking empty list.
func TestTasksRejectsInvalidWorkState(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)

	_, _, err := tracemcp.CallTasks(srv, ctx, tracemcp.TasksInput{WorkState: "ARCHIVED"})
	if err == nil || !strings.Contains(err.Error(), "work_state must be PENDING") {
		t.Fatalf("want a work_state validation error, got %v", err)
	}
	_, _, err = tracemcp.CallTasks(srv, ctx, tracemcp.TasksInput{WorkStates: []string{"PENDING", "NOPE"}})
	if err == nil || !strings.Contains(err.Error(), "work_states[1]") {
		t.Fatalf("want an indexed work_states error, got %v", err)
	}
}

// TestTasksAcceptsValidWorkState is the keeper.
func TestTasksAcceptsValidWorkState(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	if _, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{Kind: "task", Title: "t"}); err != nil {
		t.Fatal(err)
	}
	for _, ws := range []string{"PENDING", "in_progress", "DONE"} {
		if _, _, err := tracemcp.CallTasks(srv, ctx, tracemcp.TasksInput{WorkState: ws}); err != nil {
			t.Fatalf("work_state %q must be accepted: %v", ws, err)
		}
	}
}

// TestLinkRejectsSelfReference is the #162 proof, including the blocks case
// recorded in the issue.
func TestLinkRejectsSelfReference(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	taskRes, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{Kind: "task", Title: "self"})
	if err != nil {
		t.Fatal(err)
	}
	var task struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(tracemcp.ResultText(t, taskRes)), &task); err != nil {
		t.Fatal(err)
	}

	for _, rel := range []string{"blocks", "api-contract"} {
		_, _, err := tracemcp.CallLink(srv, ctx, tracemcp.LinkInput{
			Rel: rel, From: task.ID, To: task.ID,
		})
		if err == nil || !strings.Contains(err.Error(), "cannot link an entity to itself") {
			t.Fatalf("rel %s: want a self-link refusal, got %v", rel, err)
		}
		if !strings.Contains(err.Error(), rel) {
			t.Fatalf("rel %s: error must name the rel, got %v", rel, err)
		}
	}

	// The row must not exist: no entity_links row may reference the task twice.
	st := mustOpenStore(t, dir)
	links, err := st.ListLinksFrom("task", task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 0 {
		t.Fatalf("a rejected self-link must not persist: %+v", links)
	}
}

// TestLinkBlocksDomainRejectsSelfLink covers the CLI/evals path that does not
// go through the MCP adapter.
func TestLinkBlocksDomainRejectsSelfLink(t *testing.T) {
	st := mustOpenStore(t, t.TempDir())
	svc := domain.New(st)
	ctx := context.Background()
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "solo"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkBlocks(ctx, task.ID, task.ID, domain.LinkMeta{}); err == nil {
		t.Fatal("want a domain-level self-link refusal")
	}
	if err := svc.LinkAPIContract(ctx, task.ID, task.ID, domain.LinkMeta{}); err == nil {
		t.Fatal("want a domain-level self-link refusal for api-contract")
	}
}

func addCreated(t *testing.T, raw string) bool {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatalf("add json: %v\n%s", err, raw)
	}
	v, ok := m["created"].(bool)
	if !ok {
		t.Fatalf("response must carry a boolean created field: %s", raw)
	}
	return v
}

func seedGoal(t *testing.T, st *store.Store) string {
	t.Helper()
	g, err := domain.New(st).CreateGoal(context.Background(), domain.GoalInput{Title: "envelope goal"})
	if err != nil {
		t.Fatal(err)
	}
	return g.ID
}

func seedTask(t *testing.T, st *store.Store, goalID string) string {
	t.Helper()
	task, err := domain.New(st).CreateTask(context.Background(), domain.TaskInput{
		Title: "envelope task", GoalID: &goalID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return task.ID
}
