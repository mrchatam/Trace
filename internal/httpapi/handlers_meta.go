package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":            true,
		"name":          "trace",
		"api_version":   APIVersion,
		"trace_version": TraceVersion,
	})
}

func (s *Server) handleProject(w http.ResponseWriter, r *http.Request) {
	storePath := filepath.Join(s.root, ".trace")
	storeReady := false
	if fi, err := os.Stat(filepath.Join(storePath, "trace.db")); err == nil && fi.Mode().IsRegular() {
		storeReady = true
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"root":        s.root,
		"store_path":  storePath,
		"store_ready": storeReady,
		"flags": map[string]any{
			"gui_placeholder": true,
		},
	})
}
