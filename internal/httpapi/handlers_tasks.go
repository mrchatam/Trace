package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

type taskListRow struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	WorkState string  `json:"work_state"`
	GoalID    *string `json:"goal_id"`
	Body      string  `json:"body,omitempty"`
}

func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}

	goalID := strings.TrimSpace(r.URL.Query().Get("goal_id"))
	workState := strings.TrimSpace(r.URL.Query().Get("work_state"))
	cursor := strings.TrimSpace(r.URL.Query().Get("cursor"))
	limit := store.HTTPDefaultTaskListLimit
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		n, perr := strconv.Atoi(raw)
		if perr != nil || n < 1 {
			writeEnvelope(w, http.StatusBadRequest, "BAD_REQUEST", "limit must be a positive integer", nil)
			return
		}
		limit = n
	}

	filt := store.TaskListFilter{
		GoalID: goalID,
		Cursor: cursor,
		Limit:  limit,
	}
	if workState != "" {
		filt.WorkStates = []string{workState}
	}
	res, err := st.ListTasksFiltered(filt)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	items := make([]taskListRow, 0, len(res.Tasks))
	for _, t := range res.Tasks {
		items = append(items, taskListRow{
			ID: t.ID, Title: t.Title, WorkState: t.WorkState, GoalID: t.GoalID,
		})
	}
	out := map[string]any{"items": items, "truncated": res.Truncated}
	if res.NextCursor != "" {
		out["next_cursor"] = res.NextCursor
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("task_id")
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	t, err := st.GetTask(id)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, taskListRow{
		ID: t.ID, Title: t.Title, WorkState: t.WorkState, GoalID: t.GoalID, Body: t.Body,
	})
}
