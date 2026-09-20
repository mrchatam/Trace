package httpapi

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mrchatam/Trace/internal/agents"
	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/gitcli"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

type reviewSummary struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Body   string  `json:"body,omitempty"`
	Result *string `json:"result,omitempty"`
	TaskID *string `json:"task_id,omitempty"`
}

func (s *Server) handleListReviews(w http.ResponseWriter, r *http.Request) {
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	svc := domain.New(st)
	taskID := strings.TrimSpace(r.URL.Query().Get("task_id"))
	cursor := strings.TrimSpace(r.URL.Query().Get("cursor"))
	limit := store.HTTPDefaultReviewListLimit
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		n, perr := strconv.Atoi(raw)
		if perr != nil || n < 1 {
			writeEnvelope(w, http.StatusBadRequest, "BAD_REQUEST", "limit must be a positive integer", nil)
			return
		}
		limit = n
	}
	res, err := svc.ListReviewsFiltered(r.Context(), store.ReviewListFilter{
		TaskID: taskID,
		Cursor: cursor,
		Limit:  limit,
	})
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	items := make([]reviewSummary, 0, len(res.Reviews))
	for _, rv := range res.Reviews {
		// List omits body; get keeps full body.
		item := reviewSummary{ID: rv.ID, Title: rv.Title}
		if rv.Result != "" {
			result := rv.Result
			item.Result = &result
		}
		items = append(items, item)
	}
	out := map[string]any{"items": items, "truncated": res.Truncated}
	if res.NextCursor != "" {
		out["next_cursor"] = res.NextCursor
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleCreateReview(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title  string `json:"title"`
		Body   string `json:"body"`
		TaskID string `json:"task_id"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON body", nil)
		return
	}
	if strings.TrimSpace(in.Title) == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "title is required", nil)
		return
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	svc := domain.New(st)
	rv, err := svc.CreateReview(r.Context(), domain.ReviewInput{Title: in.Title, Body: in.Body})
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	var taskPtr *string
	if strings.TrimSpace(in.TaskID) != "" {
		if err := svc.LinkReviewTask(r.Context(), rv.ID, in.TaskID, domain.LinkMeta{}); err != nil {
			mapDomainErr(w, err)
			return
		}
		t := in.TaskID
		taskPtr = &t
	}
	writeJSON(w, http.StatusCreated, reviewSummary{ID: rv.ID, Title: rv.Title, Body: rv.Body, TaskID: taskPtr})
}

func (s *Server) handleGetReview(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("review_id")
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	rv, err := domain.New(st).GetReview(r.Context(), id)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	out := reviewSummary{ID: rv.ID, Title: rv.Title, Body: rv.Body}
	if rv.Result != "" {
		res := rv.Result
		out.Result = &res
	}
	writeJSON(w, http.StatusOK, out)
}

const (
	defaultPlanListGoalLimit = 20
	maxPlanListGoalLimit     = 100 // plans are expensive (GetPlan per goal)
)

func (s *Server) handleListPlans(w http.ResponseWriter, r *http.Request) {
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}

	goalID := strings.TrimSpace(r.URL.Query().Get("goal_id"))
	ps := planner.New(st)

	if goalID != "" {
		view, err := ps.GetPlan(r.Context(), goalID)
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": []any{view}, "count": 1, "truncated": false})
		return
	}

	limit := defaultPlanListGoalLimit
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		n, perr := strconv.Atoi(raw)
		if perr != nil || n < 1 {
			writeEnvelope(w, http.StatusBadRequest, "BAD_REQUEST", "limit must be a positive integer", nil)
			return
		}
		limit = n
		if limit > maxPlanListGoalLimit {
			limit = maxPlanListGoalLimit
		}
	}

	// Fetch limit+1 so truncation is decided in SQL, not load-all-then-slice.
	goals, err := st.ListGoalsLimited(limit + 1)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	truncated := false
	if len(goals) > limit {
		truncated = true
		goals = goals[:limit]
	}
	items := make([]any, 0, len(goals))
	for _, g := range goals {
		view, err := ps.GetPlan(r.Context(), g.ID)
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		items = append(items, view)
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "count": len(items), "truncated": truncated})
}

func (s *Server) handlePlanBootstrap(w http.ResponseWriter, r *http.Request) {
	var in struct {
		GoalID string `json:"goal_id"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON body", nil)
		return
	}
	goalID := strings.TrimSpace(in.GoalID)
	if goalID == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "goal_id is required", nil)
		return
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	res, err := planner.New(st).BootstrapFromPlanChanges(r.Context(), goalID, "http")
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleListCapability(w http.ResponseWriter, r *http.Request) {
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	svc := domain.New(st)
	action := strings.TrimSpace(r.URL.Query().Get("action"))
	if action == "" {
		action = "list"
	}
	taskID := strings.TrimSpace(r.URL.Query().Get("task_id"))
	switch action {
	case "list":
		capLimit := store.DefaultTaskListLimit
		if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
			n, perr := strconv.Atoi(raw)
			if perr != nil || n < 1 {
				writeEnvelope(w, http.StatusBadRequest, "BAD_REQUEST", "limit must be a positive integer", nil)
				return
			}
			capLimit = n
			if capLimit > store.MaxTaskListLimit {
				capLimit = store.MaxTaskListLimit
			}
		}
		if strings.EqualFold(strings.TrimSpace(r.URL.Query().Get("all")), "true") || r.URL.Query().Get("all") == "1" {
			capLimit = 0 // unbounded (domain: Limit 0)
		}
		list, err := svc.ListCapabilities(r.Context(), domain.ListCapabilitiesFilter{Limit: capLimit})
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		if list == nil {
			list = []store.Capability{}
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "capabilities": list, "count": len(list)})
	case "missing":
		if taskID == "" {
			writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "task_id is required for action=missing", nil)
			return
		}
		list, err := svc.MissingCapabilities(r.Context(), taskID)
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		if list == nil {
			list = []store.Capability{}
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "capabilities": list, "count": len(list)})
	default:
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "action must be list|missing", nil)
	}
}

func (s *Server) handleGetImpact(w http.ResponseWriter, r *http.Request) {
	taskID := strings.TrimSpace(r.URL.Query().Get("task_id"))
	decisionID := strings.TrimSpace(r.URL.Query().Get("decision_id"))
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	svc := domain.New(st)
	if decisionID != "" {
		rep, err := svc.ImpactReport(r.Context(), decisionID)
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"ok": true, "decision_id": rep.Decision.ID,
			"affected_task_ids": rep.AffectedTaskIDs, "findings": rep.Findings,
			"alternatives": rep.Alternatives, "overall_class": rep.OverallClass,
			"overall_uncertainty": rep.OverallUncertainty, "has_unknown": rep.HasUnknown,
			"incomplete": rep.Incomplete,
		})
		return
	}
	if taskID == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "task_id or decision_id is required", nil)
		return
	}
	summaries, err := svc.ImpactSummariesForTask(r.Context(), taskID)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "task_id": taskID, "items": summaries})
}

func (s *Server) handleListChanges(w http.ResponseWriter, r *http.Request) {
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	taskID := strings.TrimSpace(r.URL.Query().Get("task_id"))
	limit, _ := queryInt(r, "limit", 32)
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}
	list, err := domain.New(st).ListChanges(r.Context(), limit, taskID)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	if list == nil {
		list = []store.Change{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": list, "count": len(list)})
}

func (s *Server) handleListRegressions(w http.ResponseWriter, r *http.Request) {
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	taskID := strings.TrimSpace(r.URL.Query().Get("task_id"))
	changeID := strings.TrimSpace(r.URL.Query().Get("change_id"))
	limit, _ := queryInt(r, "limit", 32)
	rows, err := domain.New(st).ListRegressions(r.Context(), domain.EvidenceQueryOpts{
		TaskID: taskID, ChangeID: changeID, Limit: limit,
	})
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": rows, "count": len(rows)})
}

func (s *Server) handleListAgents(w http.ResponseWriter, r *http.Request) {
	action := strings.TrimSpace(r.URL.Query().Get("action"))
	if action == "" {
		action = "list"
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	switch action {
	case "list":
		items, err := agents.ListAgentSummaries(r.Context(), st)
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		if items == nil {
			items = []agents.AgentListItem{}
		}
		out := map[string]any{"ok": true, "items": items, "count": len(items)}
		if len(items) == 0 {
			out["hint"] = agents.EmptyCatalogHint
		}
		writeJSON(w, http.StatusOK, out)
	case "recommend":
		taskID := strings.TrimSpace(r.URL.Query().Get("task_id"))
		phase := strings.TrimSpace(r.URL.Query().Get("phase"))
		if taskID == "" && phase == "" {
			writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "task_id or phase is required for action=recommend", nil)
			return
		}
		recs, err := loop.RecommendHarness(r.Context(), domain.New(st), st, planner.New(st), loop.RecommendHarnessInput{
			TaskID: taskID,
			Phase:  phase,
		})
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		if recs == nil {
			recs = []agents.Recommendation{}
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": recs, "count": len(recs)})
	default:
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "action must be list|recommend", nil)
	}
}

func (s *Server) handleIndexStatus(w http.ResponseWriter, r *http.Request) {
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	state, err := st.GetGraphSyncState()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	out := map[string]any{
		"head":                "",
		"last_indexed_commit": state.LastIndexedCommit,
		"stale":               false,
		"hook_installed":      state.HookInstalled,
		"supported_languages": analyzers.SupportedLanguages(),
	}
	if repo, rerr := gitcli.OpenWithStore(s.root, st); rerr == nil {
		defer repo.Close()
		head, herr := repo.Head(context.Background())
		if herr == nil {
			out["head"] = head
			out["stale"] = head != state.LastIndexedCommit
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleAuthToken(w http.ResponseWriter, r *http.Request) {
	tok, err := GenerateToken()
	if err != nil {
		writeEnvelope(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to generate token", nil)
		return
	}
	s.SetToken(tok)
	writeJSON(w, http.StatusOK, map[string]any{
		"token":      tok,
		"expires_at": nil,
		"issued_at":  time.Now().UTC().Format(time.RFC3339),
	})
}
