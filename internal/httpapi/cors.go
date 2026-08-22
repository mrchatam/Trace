package httpapi

import (
	"net/http"
	"net/url"
	"strings"
)

// applyCORS enforces deny-by-default CORS: never sets Access-Control-Allow-Origin: *.
// When corsOrigin is non-empty, exact Origin match is reflected (+ Vary: Origin).
// Same-origin SPA needs no CORS headers; unmatched foreign origins get none.
func applyCORS(corsOrigin string, next http.Handler) http.Handler {
	corsOrigin = strings.TrimSpace(corsOrigin)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if corsOrigin != "" && origin != "" && origin == corsOrigin {
			w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
			w.Header().Set("Vary", "Origin")
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
				w.Header().Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		} else if r.Method == http.MethodOptions {
			// Explicitly do not set ACAO wildcard. Preflight without allow-* so browsers deny.
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// corsConnectSrcHost returns a CSP connect-src host (scheme://host[:port]) from an exact Origin URL,
// or empty if unset/invalid. Never returns "*".
func corsConnectSrcHost(corsOrigin string) string {
	corsOrigin = strings.TrimSpace(corsOrigin)
	if corsOrigin == "" || corsOrigin == "*" {
		return ""
	}
	u, err := url.Parse(corsOrigin)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	return u.Scheme + "://" + u.Host
}
