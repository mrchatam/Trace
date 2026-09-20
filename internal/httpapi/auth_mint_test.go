package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrchatam/Trace/internal/store"
)

// Regression for #72: POST /v1/auth/token on loopback must not flip requireToken
// and brick /v1/health / GUI (curl repro in the issue).
func TestAuthTokenMintDoesNotLockOutHealth(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv, err := New(Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	h := srv.Handler()

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("health before mint: %d %s", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/auth/token", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("mint: %d %s", rr.Code, rr.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("mint json: %v", err)
	}
	tok, _ := body["token"].(string)
	if tok == "" {
		t.Fatalf("mint missing token: %v", body)
	}
	if srv.Token() != tok {
		t.Fatalf("Token()=%q want %q", srv.Token(), tok)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("health after mint must stay 200, got %d %s (mint must not enable requireToken)", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/tasks", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("tasks after mint must stay 200, got %d %s", rr.Code, rr.Body.String())
	}
}
