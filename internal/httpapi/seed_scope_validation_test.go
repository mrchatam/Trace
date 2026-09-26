package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestSeedImportBadKindReturns400(t *testing.T) {
	dir, srv := openFixture(t)
	h := srv.Handler()

	seed := domain.SeedDocument{
		Scopes: []domain.SeedScope{
			{ID: "11111111-1111-1111-1111-111111111111", Slug: "bad-kind", Title: "Bad", Kind: "Feature"},
		},
	}
	inputPath := filepath.Join(dir, "seed_bad_kind.json")
	raw, _ := json.MarshalIndent(seed, "", "  ")
	if err := os.WriteFile(inputPath, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	body, _ := json.Marshal(map[string]string{"input_path": inputPath})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/seed/import", bytes.NewReader(body)))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d %s", rr.Code, rr.Body.String())
	}
	var env map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if env["error"].(map[string]any)["code"] != "VALIDATION_ERROR" {
		t.Fatalf("code: %v", env)
	}
}

func TestSeedImportSlugCollisionReturns400(t *testing.T) {
	dir, srv := openFixture(t)
	h := srv.Handler()

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	existing, err := svc.CreateScope(context.Background(), domain.ScopeInput{
		Slug:  "collision",
		Title: "Existing",
		Kind:  "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	seed := domain.SeedDocument{
		Scopes: []domain.SeedScope{
			{ID: "22222222-2222-2222-2222-222222222222", Slug: "collision", Title: "Different", Kind: "layer"},
		},
	}
	inputPath := filepath.Join(dir, "seed_collision.json")
	raw, _ := json.MarshalIndent(seed, "", "  ")
	if err := os.WriteFile(inputPath, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	body, _ := json.Marshal(map[string]string{"input_path": inputPath})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/seed/import", bytes.NewReader(body)))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d %s", rr.Code, rr.Body.String())
	}
	var env map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if env["error"].(map[string]any)["code"] != "VALIDATION_ERROR" {
		t.Fatalf("code: %v", env)
	}
	msg, _ := env["error"].(map[string]any)["message"].(string)
	if msg == "" || !contains(msg, "already exists with id") || !contains(msg, existing.ID) || !contains(msg, "22222222-2222-2222-2222-222222222222") {
		t.Fatalf("bad message: %q", msg)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
