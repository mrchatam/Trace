package httpapi

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/mrchatam/Trace/internal/store"
)

// Regression for #93: concurrent SetToken + authenticated traffic must not race
// (go test -race) and must preserve #72 (mint does not flip requireToken).
func TestSetTokenConcurrentWithAuth(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv, err := New(Options{Root: dir, Addr: "127.0.0.1:0", Token: "initial-token"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(srv.CloseStore)
	h := srv.Handler()

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			srv.SetToken("tok-" + string(rune('a'+i%26)))
		}(i)
		go func() {
			defer wg.Done()
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
			req.Header.Set("Authorization", "Bearer "+srv.Token())
			h.ServeHTTP(rr, req)
			// With requireToken true (token opted in on loopback), health needs bearer.
			if rr.Code != http.StatusOK && rr.Code != http.StatusUnauthorized {
				t.Errorf("unexpected status %d", rr.Code)
			}
		}()
	}
	wg.Wait()

	// #72: mint/SetToken must not change requireToken lockout semantics beyond New().
	// Server was constructed with Token set → requireToken already true; SetToken
	// still must not panic or tear the string under -race.
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("without bearer want 401 got %d (requireToken should stay as New configured)", rr.Code)
	}
}

func TestSetTokenDoesNotFlipRequireTokenOnLoopbackTrust(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv, err := New(Options{Root: dir, Addr: "127.0.0.1:0"}) // no token → loopback trust
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(srv.CloseStore)
	h := srv.Handler()
	srv.SetToken("minted")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("loopback trust after SetToken: %d %s (#72)", rr.Code, rr.Body.String())
	}
}
