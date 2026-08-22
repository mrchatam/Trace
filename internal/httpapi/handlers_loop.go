package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/gitcli"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/retrieval"
)

const gateSchemaVersion = "trace.loop.gate.v1"

type gateEnvelope struct {
	SchemaVersion    string           `json:"schema_version"`
	TaskID           string           `json:"task_id"`
	For              string           `json:"for"`
	Allowed          bool             `json:"allowed"`
	RecommendedPhase string           `json:"recommended_phase,omitempty"`
	ReasonCode       string           `json:"reason_code,omitempty"`
	Violations       []loop.Violation `json:"violations"`
}

func (s *Server) handleLoopStatus(w http.ResponseWriter, r *http.Request) {
	taskID := strings.TrimSpace(r.URL.Query().Get("task_id"))
	if taskID == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "task_id is required", nil)
		return
	}
	goalID := strings.TrimSpace(r.URL.Query().Get("goal_id"))
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	res, err := loop.Status(r.Context(), st, planner.New(st), loop.ApplySeed{TaskID: taskID, GoalID: goalID})
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleLoopNext(w http.ResponseWriter, r *http.Request) {
	taskID := strings.TrimSpace(r.URL.Query().Get("task_id"))
	if taskID == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "task_id is required", nil)
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
	comp := compiler.New(st).WithRetrieval(eng)
	packet, err := loop.BuildNextPacket(r.Context(), loop.BuildNextInput{
		TaskID: taskID, Store: st, Planner: planner.New(st), Retrieval: eng, Compiler: comp,
	})
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, packet)
}

func (s *Server) handleLoopApply(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	_ = r.Body.Close()
	if err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "failed to read body", nil)
		return
	}
	env, err := loop.ParseApplyEnvelope(raw)
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
	res, err := loop.Apply(r.Context(), st, planner.New(st), env)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleLoopGate(w http.ResponseWriter, r *http.Request) {
	taskID := strings.TrimSpace(r.URL.Query().Get("task_id"))
	if taskID == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "task_id is required", nil)
		return
	}
	gateFor := strings.TrimSpace(r.URL.Query().Get("for"))
	if gateFor == "" {
		gateFor = "edit"
	}
	gf, err := parseLoopGateFor(gateFor)
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
	allowed, violations, err := loop.EvaluateGate(r.Context(), domain.New(st), planner.New(st), st, taskID, gf)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	env := gateEnvelope{
		SchemaVersion: gateSchemaVersion,
		TaskID:        taskID,
		For:           gateFor,
		Allowed:       allowed,
		Violations:    []loop.Violation{},
	}
	if len(violations) > 0 {
		env.Violations = violations
	}
	if !allowed && len(violations) == 1 {
		env.RecommendedPhase = violations[0].RecommendedPhase
		env.ReasonCode = violations[0].ReasonCode
	}
	writeJSON(w, http.StatusOK, env)
}

func parseLoopGateFor(s string) (loop.GateFor, error) {
	switch s {
	case string(loop.GateForOrient):
		return loop.GateForOrient, nil
	case string(loop.GateForEdit):
		return loop.GateForEdit, nil
	case string(loop.GateForExecute):
		return loop.GateForExecute, nil
	case string(loop.GateForDone):
		return loop.GateForDone, nil
	case string(loop.GateForExport):
		return loop.GateForExport, nil
	default:
		return "", errString("invalid for (want orient|edit|execute|done|export)")
	}
}

func (s *Server) handleLoopReset(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TaskID string `json:"task_id"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<16)).Decode(&body); err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON body", nil)
		return
	}
	if strings.TrimSpace(body.TaskID) == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "task_id is required", nil)
		return
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	next, err := domain.New(st).ResetDeliberationState(r.Context(), body.TaskID)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, next)
}
