package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mrchatam/Trace/internal/store"
)

func TestListTasksHonorsLimitAndCursor(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		if _, err := st.UpsertTask(store.Task{Title: fmt.Sprintf("t%d", i), WorkState: store.WorkStatePending}); err != nil {
			t.Fatal(err)
		}
	}
	st.Close()

	srv, err := New(Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatal(err)
	}
	h := srv.Handler()

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/tasks?limit=2", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	var page1 map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &page1); err != nil {
		t.Fatal(err)
	}
	items1, _ := page1["items"].([]any)
	if len(items1) != 2 {
		t.Fatalf("want 2 items got %v", page1)
	}
	if page1["truncated"] != true {
		t.Fatalf("want truncated true: %v", page1)
	}
	cur, _ := page1["next_cursor"].(string)
	if cur == "" {
		t.Fatal("want next_cursor")
	}

	rr = httptest.NewRecorder()
	u := "/v1/tasks?limit=2&cursor=" + url.QueryEscape(cur)
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, u, nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("page2 status %d %s", rr.Code, rr.Body.String())
	}
	var page2 map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &page2); err != nil {
		t.Fatal(err)
	}
	items2, _ := page2["items"].([]any)
	if len(items2) != 2 {
		t.Fatalf("page2 want 2 got %v", page2)
	}
	id1, _ := items1[0].(map[string]any)["id"].(string)
	id2, _ := items2[0].(map[string]any)["id"].(string)
	if id1 == id2 {
		t.Fatal("cursor did not advance")
	}
}
