package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
