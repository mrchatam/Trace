package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const placeholderHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/>
<title>Trace</title>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<style>
body{font-family:system-ui,sans-serif;max-width:40rem;margin:3rem auto;padding:0 1rem;line-height:1.5;color:#1a1a1a}
code{background:#f4f4f4;padding:0.1em 0.35em;border-radius:3px}
</style>
</head>
<body>
<h1>Trace</h1>
<p>Local HTTP API is running, but the embedded GUI assets are missing or incomplete
(packaging mistake — run <code>scripts/embed-gui.sh</code> before <code>go build</code>
in a Trace checkout). Consumer projects only need <code>.trace/</code>; they do not
require a project <code>web/</code>.</p>
<p>JSON API prefix: <code>/v1</code></p>
</body>
</html>
`

func (s *Server) serveStaticOrPlaceholder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		writeEnvelope(w, http.StatusMethodNotAllowed, "BAD_REQUEST", "method not allowed", nil)
		return
	}
	path := r.URL.Path
	if path == "" {
		path = "/"
	}

	// Prefer disk StaticDir when index.html is present.
	if diskIndexExists(s.staticDir) {
		s.setStaticCSP(w)
		if path == "/" {
			http.ServeFile(w, r, filepath.Join(s.staticDir, "index.html"))
			return
		}
		rel := strings.TrimPrefix(path, "/")
		candidate := filepath.Join(s.staticDir, filepath.Clean(rel))
		// Rel-based confinement (avoids strings.HasPrefix sibling prefix tricks).
		relToStatic, rerr := filepath.Rel(s.staticDir, candidate)
		if rerr != nil || relToStatic == ".." || strings.HasPrefix(relToStatic, ".."+string(filepath.Separator)) {
			writeEnvelope(w, http.StatusNotFound, "NOT_FOUND", "not found", nil)
			return
		}
		if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
			http.ServeFile(w, r, candidate)
			return
		}
		// SPA fallback for non-/v1 GET
		http.ServeFile(w, r, filepath.Join(s.staticDir, "index.html"))
		return
	}

	// Embed fallback when disk dist is missing.
	if embeddedIndexExists() {
		s.serveEmbedded(w, r, path)
		return
	}

	s.setStaticCSP(w)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		_, _ = w.Write([]byte(placeholderHTML))
	}
}
