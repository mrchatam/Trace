package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrchatam/Trace/internal/store"
)

func TestListPlansClampsLimitToMax(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		if _, err := st.UpsertGoal(store.Goal{Title: fmt.Sprintf("g%d", i), Body: "b"}); err != nil {
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
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/plans?limit=999999", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d %s", rr.Code, rr.Body.String())
	}
	var page map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	items, _ := page["items"].([]any)
	// Only 5 goals exist; clamp must not error and must return them all.
	if len(items) != 5 {
		t.Fatalf("items=%d want 5 (clamp must accept huge limit): %v", len(items), page)
	}
}
