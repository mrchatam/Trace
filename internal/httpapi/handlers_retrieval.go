package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/gitcli"
	"github.com/mrchatam/Trace/internal/retrieval"
)

func (s *Server) handleContext(w http.ResponseWriter, r *http.Request) {
	taskID := strings.TrimSpace(r.URL.Query().Get("task_id"))
	if taskID == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "task_id is required", nil)
		return
	}
	depth, err := queryInt(r, "depth", 1)
	if err != nil || depth < 1 || depth > 2 {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "depth must be 1 or 2", nil)
		return
	}
	format := strings.TrimSpace(r.URL.Query().Get("format"))
	if format == "" {
		format = "json"
	}
	includeWhy := r.URL.Query().Get("include_why") == "true" || r.URL.Query().Get("include_why") == "1"

	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	eng := retrieval.New(st)
	if repo, rerr := gitcli.OpenWithStore(s.root, st); rerr == nil {
		defer repo.Close()
		eng = eng.WithVCS(repo)
	}
	comp := compiler.New(st).WithRetrieval(eng)
	opts := compiler.ContextOptions{
		IncludeWhy:      includeWhy,
		IncludeMarkdown: format == "markdown" || format == "both",
	}
	var pkt compiler.Packet
	if depth == 1 {
		pkt, err = comp.TaskContext(r.Context(), taskID, opts)
	} else {
		pkt, err = comp.ExpandContext(r.Context(), taskID, depth, opts)
	}
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	switch format {
	case "json":
		b, err := pkt.JSON()
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		writeRawJSON(w, http.StatusOK, b)
	case "markdown":
		writeJSON(w, http.StatusOK, map[string]any{"markdown": pkt.Markdown()})
	case "both":
		b, err := pkt.JSON()
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		var obj map[string]any
		if err := json.Unmarshal(b, &obj); err != nil {
			mapDomainErr(w, err)
			return
		}
		obj["markdown"] = pkt.Markdown()
		writeJSON(w, http.StatusOK, obj)
	default:
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "format must be json|markdown|both", nil)
	}
}

func (s *Server) handleWhy(w http.ResponseWriter, r *http.Request) {
	entityType := strings.TrimSpace(r.URL.Query().Get("entity_type"))
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if entityType == "" || id == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "entity_type and id are required", nil)
		return
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	eng := retrieval.New(st)
	if repo, rerr := gitcli.OpenWithStore(s.root, st); rerr == nil {
		defer repo.Close()
		eng = eng.WithVCS(repo)
	}
	res, err := eng.Why(r.Context(), entityType, id)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	impact, err := domain.New(st).ImpactSummariesForWhySeed(r.Context(), retrieval.NormalizeEntityType(entityType), id)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	out := struct {
		retrieval.WhyResult
		Impact []domain.DecisionImpact `json:"impact,omitempty"`
	}{WhyResult: res, Impact: impact}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "q is required", nil)
		return
	}
	limit, err := queryInt(r, "limit", 20)
	if err != nil || limit < 1 {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid limit", nil)
		return
	}
	if limit > 100 {
		limit = 100
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	hits, err := retrieval.New(st).Search(r.Context(), q, retrieval.SearchOptions{Limit: limit})
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	if hits == nil {
		hits = []retrieval.Hit{}
	}
	items := make([]map[string]any, 0, len(hits))
	for _, h := range hits {
		items = append(items, map[string]any{
			"id": h.EntityID, "kind": h.EntityType, "title": h.Title, "snippet": h.Excerpt,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *Server) handleGraph(w http.ResponseWriter, r *http.Request) {
	center := strings.TrimSpace(r.URL.Query().Get("center"))
	maxRaw := strings.TrimSpace(r.URL.Query().Get("max_nodes"))
	if center == "" || maxRaw == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "center and max_nodes are required", nil)
		return
	}
	maxNodes, err := queryInt(r, "max_nodes", 0)
	if err != nil || maxNodes < 1 {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "max_nodes must be a positive integer", nil)
		return
	}
	depth, err := queryInt(r, "depth", 2)
	if err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid depth", nil)
		return
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	g, err := retrieval.New(st).Neighborhood(r.Context(), retrieval.NeighborhoodOpts{
		Center: center, MaxNodes: maxNodes, Depth: depth,
	})
	if err != nil {
		var bud *retrieval.ErrBudgetExceeded
		if errors.As(err, &bud) {
			writeEnvelope(w, http.StatusBadRequest, "BUDGET_EXCEEDED", bud.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			writeEnvelope(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "max_nodes") {
			writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			return
		}
		mapDomainErr(w, err)
		return
	}
	if g.Nodes == nil {
		g.Nodes = []retrieval.GraphNode{}
	}
	if g.Edges == nil {
		g.Edges = []retrieval.GraphEdge{}
	}
	writeJSON(w, http.StatusOK, g)
}
