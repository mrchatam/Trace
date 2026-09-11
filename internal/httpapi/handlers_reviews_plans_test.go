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

func TestListReviewsHonorsLimitAndCursor(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		if _, err := st.UpsertReview(store.Review{Title: fmt.Sprintf("r%d", i), Body: "full-body"}); err != nil {
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
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/reviews?limit=2", nil))
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
	row0, _ := items1[0].(map[string]any)
	if _, has := row0["body"]; has && row0["body"] != "" && row0["body"] != nil {
		t.Fatalf("list must omit body: %#v", row0)
	}
	cur, _ := page1["next_cursor"].(string)
	if cur == "" {
		t.Fatal("want next_cursor")
	}

	rr = httptest.NewRecorder()
	u := "/v1/reviews?limit=2&cursor=" + url.QueryEscape(cur)
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

func TestListPlansRequiresGoalOrBoundsGoals(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	var firstID string
	for i := 0; i < 3; i++ {
		g, err := st.UpsertGoal(store.Goal{Title: fmt.Sprintf("g%d", i)})
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			firstID = g.ID
		}
	}
	st.Close()

	srv, err := New(Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatal(err)
	}
	h := srv.Handler()

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/plans?limit=2", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	var page map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	items, _ := page["items"].([]any)
	if len(items) != 2 {
		t.Fatalf("want 2 plans got %v", page)
	}
	if page["truncated"] != true {
		t.Fatalf("want truncated: %v", page)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/plans?goal_id="+firstID, nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("goal_id status %d %s", rr.Code, rr.Body.String())
	}
	var one map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &one); err != nil {
		t.Fatal(err)
	}
	oneItems, _ := one["items"].([]any)
	if len(oneItems) != 1 || one["truncated"] != false {
		t.Fatalf("goal_id page: %v", one)
	}
}
