package httpapi

import (
	"net/http"
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
	defer st.Close()

	goalID := strings.TrimSpace(r.URL.Query().Get("goal_id"))
	workState := strings.TrimSpace(r.URL.Query().Get("work_state"))

	var tasks []store.Task
	if goalID != "" {
		tasks, err = st.ListTasksByGoalID(goalID)
	} else {
		tasks, err = st.ListTasks()
	}
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	items := make([]taskListRow, 0, len(tasks))
	for _, t := range tasks {
		if workState != "" && t.WorkState != workState {
			continue
		}
		items = append(items, taskListRow{
			ID: t.ID, Title: t.Title, WorkState: t.WorkState, GoalID: t.GoalID,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("task_id")
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	t, err := st.GetTask(id)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, taskListRow{
		ID: t.ID, Title: t.Title, WorkState: t.WorkState, GoalID: t.GoalID, Body: t.Body,
	})
}
