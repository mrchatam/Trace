package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// TestCapabilityListSnakeCase covers #105: GET /v1/capability must emit
// snake_case keys (id/kind/slug/…), not Go PascalCase field names.
func TestCapabilityListSnakeCase(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	if _, err := svc.UpsertCapability(context.Background(), domain.CapabilityInput{
		Kind: domain.CapabilityKindMCP, Slug: "mcp:demo", Title: "Demo",
		Status: domain.CapabilityStatusAvailable,
	}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	srv, err := New(Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(srv.CloseStore)

	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/capability", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rr.Code, rr.Body.String())
	}
	body := rr.Body.String()
	for _, bad := range []string{`"ID"`, `"Kind"`, `"Slug"`, `"CreatedAt"`, `"UpdatedAt"`} {
		if strings.Contains(body, bad) {
			t.Fatalf("PascalCase key %s in body: %s", bad, body)
		}
	}
	var env struct {
		OK           bool               `json:"ok"`
		Capabilities []store.Capability `json:"capabilities"`
		Count        int                `json:"count"`
	}
	if err := json.Unmarshal([]byte(body), &env); err != nil {
		t.Fatal(err)
	}
	if !env.OK || env.Count < 1 || len(env.Capabilities) < 1 {
		t.Fatalf("envelope: %+v body=%s", env, body)
	}
	c := env.Capabilities[0]
	if c.ID == "" || c.Kind == "" || c.Slug == "" || c.CreatedAt == "" {
		t.Fatalf("snake_case decode failed: %+v body=%s", c, body)
	}
	var raw map[string]any
	_ = json.Unmarshal([]byte(body), &raw)
	arr, _ := raw["capabilities"].([]any)
	row, _ := arr[0].(map[string]any)
	for _, key := range []string{"id", "kind", "slug", "title", "status", "created_at", "updated_at"} {
		if _, ok := row[key]; !ok {
			t.Fatalf("missing snake_case key %q in %#v", key, row)
		}
	}
}
