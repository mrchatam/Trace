package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// Regression for #86: concurrent store-backed handlers must share one Open and
// must not return opaque 500 INTERNAL_ERROR from exclusive flock contention.
func TestConcurrentHandlersShareStore(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	goal, err := svc.CreateGoal(context.Background(), domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(context.Background(), domain.TaskInput{Title: "t", GoalID: &goal.ID})
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	srv, err := New(Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(srv.CloseStore)
	h := srv.Handler()

	const n = 16
	var wg sync.WaitGroup
	errs := make(chan string, n*2)
	for i := 0; i < n; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/v1/loop/status?task_id="+task.ID+"&goal_id="+goal.ID, nil)
			h.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				errs <- "status: " + rr.Body.String()
			}
		}()
		go func() {
			defer wg.Done()
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/v1/loop/gate?task_id="+task.ID+"&goal_id="+goal.ID+"&for=execute", nil)
			h.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				errs <- "gate: " + rr.Body.String()
			}
		}()
	}
	wg.Wait()
	close(errs)
	for msg := range errs {
		t.Fatalf("concurrent handler failed: %s", msg)
	}

	// Shared store still usable after fan-out.
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/tasks", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("tasks after fan-out: %d %s", rr.Code, rr.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
}

func TestMapDomainErrLockedIs503(t *testing.T) {
	rr := httptest.NewRecorder()
	if !mapDomainErr(rr, store.ErrLocked) {
		t.Fatal("expected mapped")
	}
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d want 503", rr.Code)
	}
	var env errorBody
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if env.Error.Code != "LOCKED" {
		t.Fatalf("code=%q want LOCKED", env.Error.Code)
	}
}

// TestLoopStatusValidationErrors covers #104: missing goal_id / seed mismatch
// must surface as 400 VALIDATION_ERROR with the real message, not opaque 500.
func TestLoopStatusValidationErrors(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	orphan, err := svc.CreateTask(ctx, domain.TaskInput{Title: "orphan"})
	if err != nil {
		t.Fatal(err)
	}
	linked, err := svc.CreateTask(ctx, domain.TaskInput{Title: "linked", GoalID: &goal.ID})
	if err != nil {
		t.Fatal(err)
	}
	other, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "other"})
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	srv, err := New(Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(srv.CloseStore)
	h := srv.Handler()

	cases := []struct {
		name string
		url  string
		want string
	}{
		{"no_goal_id", "/v1/loop/status?task_id="+orphan.ID+"&goal_id="+goal.ID, "has no goal_id"},
		{"seed_mismatch", "/v1/loop/status?task_id="+linked.ID+"&goal_id="+other.ID, "seed goal mismatch"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, tc.url, nil))
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
			}
			var env map[string]any
			if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
				t.Fatal(err)
			}
			errObj, _ := env["error"].(map[string]any)
			if errObj["code"] != "VALIDATION_ERROR" {
				t.Fatalf("code: %#v body=%s", errObj, rr.Body.String())
			}
			msg, _ := errObj["message"].(string)
			if !strings.Contains(msg, tc.want) {
				t.Fatalf("message %q want substring %q", msg, tc.want)
			}
			if strings.Contains(msg, "internal error") {
				t.Fatalf("must not hide message: %q", msg)
			}
		})
	}
}


// TestLoopNextWhyValidationErrors covers #114: missing plan / unknown why type
// must surface as 400 VALIDATION_ERROR with the real message, not opaque 500.
func TestLoopNextWhyValidationErrors(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t", GoalID: &goal.ID})
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	srv, err := New(Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(srv.CloseStore)
	h := srv.Handler()

	cases := []struct {
		name string
		url  string
		want string
	}{
		{"loop_next_missing_plan", "/v1/loop/next?task_id=" + task.ID, "missing goal plan context"},
		{"why_unknown_type", "/v1/why?entity_type=nope&id=" + task.ID, "unknown entity type"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, tc.url, nil))
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
			}
			var env map[string]any
			if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
				t.Fatal(err)
			}
			errObj, _ := env["error"].(map[string]any)
			if errObj["code"] != "VALIDATION_ERROR" {
				t.Fatalf("code: %#v body=%s", errObj, rr.Body.String())
			}
			msg, _ := errObj["message"].(string)
			if !strings.Contains(msg, tc.want) {
				t.Fatalf("message %q want substring %q", msg, tc.want)
			}
			if strings.Contains(msg, "internal error") {
				t.Fatalf("must not hide message: %q", msg)
			}
		})
	}
}
