package httpapi

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
)

// GenerateToken returns a high-entropy opaque bearer token (32 bytes hex).
func GenerateToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

// LoadTokenFile reads a token from path (trimmed whitespace). Empty file is an error.
func LoadTokenFile(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	tok := strings.TrimSpace(string(raw))
	if tok == "" {
		return "", errEmptyTokenFile
	}
	return tok, nil
}

var errEmptyTokenFile = errString("httpapi: token file is empty")

type errString string

func (e errString) Error() string { return string(e) }

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if h == "" {
		return ""
	}
	const prefix = "Bearer "
	if len(h) < len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(h[len(prefix):])
}

func tokenEqual(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// authMiddleware enforces bearer when requireToken is true.
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.requireToken {
			next.ServeHTTP(w, r)
			return
		}
		got := bearerToken(r)
		if !tokenEqual(got, s.token) {
			writeEnvelope(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing or invalid bearer token", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
