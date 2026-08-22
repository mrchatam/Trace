package httpapi

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:generate ../../scripts/embed-gui.sh
//go:embed all:embeddist
var embeddedUIRoot embed.FS

// embeddedUI is the embeddist subtree (real Explore SPA after scripts/embed-gui.sh).
var embeddedUI fs.FS

func init() {
	sub, err := fs.Sub(embeddedUIRoot, "embeddist")
	if err != nil {
		embeddedUI = embeddedUIRoot
		return
	}
	embeddedUI = sub
}

func diskIndexExists(staticDir string) bool {
	fi, err := os.Stat(filepath.Join(staticDir, "index.html"))
	return err == nil && !fi.IsDir()
}

func embeddedIndexExists() bool {
	fi, err := fs.Stat(embeddedUI, "index.html")
	return err == nil && !fi.IsDir()
}

// setStaticCSP applies Content-Security-Policy for SPA/static HTML responses (not /v1 JSON).
func (s *Server) setStaticCSP(w http.ResponseWriter) {
	parts := []string{
		"default-src 'self'",
		"base-uri 'self'",
		"frame-ancestors 'none'",
		"object-src 'none'",
		"script-src 'self'",
		// Placeholder + rare inline style hooks; Vite build uses external CSS.
		"style-src 'self' 'unsafe-inline'",
	}
	connect := "connect-src 'self'"
	if host := corsConnectSrcHost(s.corsOrigin); host != "" {
		connect += " " + host
	}
	parts = append(parts, connect)
	w.Header().Set("Content-Security-Policy", strings.Join(parts, "; "))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Referrer-Policy", "no-referrer")
}

func (s *Server) serveEmbedded(w http.ResponseWriter, r *http.Request, reqPath string) {
	s.setStaticCSP(w)
	name := "index.html"
	if reqPath != "/" && reqPath != "" {
		candidate := path.Clean(strings.TrimPrefix(reqPath, "/"))
		if candidate != "." && candidate != ".." && !strings.HasPrefix(candidate, "../") {
			if fi, err := fs.Stat(embeddedUI, candidate); err == nil && !fi.IsDir() {
				name = candidate
			}
		}
	}
	f, err := embeddedUI.Open(name)
	if err != nil {
		writeEnvelope(w, http.StatusNotFound, "NOT_FOUND", "not found", nil)
		return
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		writeEnvelope(w, http.StatusNotFound, "NOT_FOUND", "not found", nil)
		return
	}
	if rs, ok := f.(io.ReadSeeker); ok {
		http.ServeContent(w, r, name, stat.ModTime(), rs)
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		writeEnvelope(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error", nil)
		return
	}
	http.ServeContent(w, r, name, stat.ModTime(), strings.NewReader(string(data)))
}
