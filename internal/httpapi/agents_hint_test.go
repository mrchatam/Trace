package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/agents"
	"github.com/mrchatam/Trace/internal/httpapi"
	"github.com/mrchatam/Trace/internal/store"
)

func TestHTTPAgentsListEmptyHint(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	srv, err := httpapi.New(httpapi.Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(srv.CloseStore)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/agents", nil))
	if rr.Code != 200 {
		t.Fatalf("%d %s", rr.Code, rr.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["ok"] != true || body["count"] != float64(0) {
		t.Fatalf("%v", body)
	}
	hint, _ := body["hint"].(string)
	if hint != agents.EmptyCatalogHint || !strings.Contains(hint, "install agents") {
		t.Fatalf("hint=%q", hint)
	}
}
