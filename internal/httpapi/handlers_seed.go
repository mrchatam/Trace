package httpapi

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
)

type seedExportRequest struct {
	OutputPath string `json:"output_path"`
	Strict     bool   `json:"strict"`
	TaskID     string `json:"task_id"`
}

type seedImportRequest struct {
	InputPath string `json:"input_path"`
}

type seedJobStatus struct {
	Status      string `json:"status"`
	Summary     string `json:"summary,omitempty"`
	Path        string `json:"path,omitempty"`
	EntityCount *int   `json:"entity_count,omitempty"`
	LinkCount   *int   `json:"link_count,omitempty"`
	Error       string `json:"error,omitempty"`
}

func (s *Server) handleSeedStatus(w http.ResponseWriter, r *http.Request) {
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	ready := true
	writeJSON(w, http.StatusOK, map[string]any{
		"ready": ready,
		"notes": "Seed HTTP returns status/summary only; full graph.json is a local file concern",
	})
}

func (s *Server) handleSeedExport(w http.ResponseWriter, r *http.Request) {
	var in seedExportRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&in) // empty body OK
		_ = r.Body.Close()
	}
	if in.Strict || strings.TrimSpace(in.TaskID) != "" {
		// Honesty/gate export options stay CLI-parity until a library helper is shared; do not silently ignore.
		writeEnvelope(w, http.StatusNotImplemented, "NOT_IMPLEMENTED",
			"seed export strict/task_id not supported over HTTP yet; use `trace seed export --strict`", nil)
		return
	}
	outPath := strings.TrimSpace(in.OutputPath)
	if outPath == "" {
		outPath = "trace/graph.json"
	}
	confined, err := confineUnderRoot(s.root, outPath)
	if err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	doc, err := domain.BuildSeedDocument(r.Context(), st, domain.ExportOpts{ProjectRoot: s.root})
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	raw = append(raw, '\n')
	if err := os.MkdirAll(filepath.Dir(confined), 0o755); err != nil {
		mapDomainErr(w, err)
		return
	}
	if err := os.WriteFile(confined, raw, 0o644); err != nil {
		mapDomainErr(w, err)
		return
	}
	entityCount := len(doc.Goals) + len(doc.Tasks) + len(doc.Decisions) + len(doc.Assumptions) +
		len(doc.Discoveries) + len(doc.PlanChanges) + len(doc.Claims) + len(doc.Evidence)
	linkCount := len(doc.Links)
	rel, _ := filepath.Rel(s.root, confined)
	writeJSON(w, http.StatusOK, seedJobStatus{
		Status:      "ok",
		Summary:     "seed export written (status/summary only; body is not the graph document)",
		Path:        rel,
		EntityCount: &entityCount,
		LinkCount:   &linkCount,
	})
}

func (s *Server) handleSeedImport(w http.ResponseWriter, r *http.Request) {
	var in seedImportRequest
	if err := decodeJSON(r, &in); err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON body", nil)
		return
	}
	if strings.TrimSpace(in.InputPath) == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "input_path is required", nil)
		return
	}
	confined, err := confineUnderRoot(s.root, in.InputPath)
	if err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	raw, err := os.ReadFile(confined)
	if err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	var doc domain.SeedDocument
	if err := json.Unmarshal(raw, &doc); err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid seed JSON", nil)
		return
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	summary, err := domain.New(st).ImportSeedDocument(r.Context(), doc)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	linkCount := summary.Links
	rel, _ := filepath.Rel(s.root, confined)
	writeJSON(w, http.StatusOK, seedJobStatus{
		Status:    "ok",
		Summary:   "seed import completed",
		Path:      rel,
		LinkCount: &linkCount,
	})
}
