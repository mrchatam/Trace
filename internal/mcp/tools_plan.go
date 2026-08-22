package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

// PlanInput mirrors trace plan create-coarse|set-current|deep|show|bootstrap.
type PlanInput struct {
	Project   string   `json:"project,omitempty" jsonschema:"optional project root override"`
	Action    string   `json:"action" jsonschema:"create-coarse|set-current|deep|show|bootstrap"`
	GoalID    string   `json:"goal_id,omitempty" jsonschema:"goal UUID"`
	Phase     string   `json:"phase,omitempty" jsonschema:"create-coarse: phase title"`
	Scope     string   `json:"scope,omitempty" jsonschema:"create-coarse: scope title"`
	ScopeID   string   `json:"scope_id,omitempty" jsonschema:"set-current|deep: scope UUID"`
	Exit      []string `json:"exit,omitempty" jsonschema:"deep: exit criteria strings"`
	Work      []string `json:"work,omitempty" jsonschema:"deep: work item titles"`
	Lookahead string   `json:"lookahead,omitempty" jsonschema:"deep: optional lookahead summary"`
}

func (s *Server) toolPlan(ctx context.Context, _ *sdkmcp.CallToolRequest, in PlanInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	if err := assertMCPToolAllowed(ctx, st, "trace_plan"); err != nil {
		_ = st.Close()
		return nil, nil, err
	}
	defer st.Close()

	psvc := planner.New(st)
	switch strings.ToLower(strings.TrimSpace(in.Action)) {
	case "create-coarse":
		return s.planCreateCoarse(ctx, psvc, in)
	case "set-current":
		return s.planSetCurrent(ctx, psvc, in)
	case "deep":
		return s.planDeep(ctx, psvc, in)
	case "show":
		return s.planShow(ctx, psvc, st, in)
	case "bootstrap":
		return s.planBootstrap(ctx, psvc, in)
	default:
		return nil, nil, fmt.Errorf("trace_plan: action must be create-coarse|set-current|deep|show|bootstrap")
	}
}

func (s *Server) planCreateCoarse(ctx context.Context, psvc *planner.Service, in PlanInput) (*sdkmcp.CallToolResult, any, error) {
	goalID := strings.TrimSpace(in.GoalID)
	phase := strings.TrimSpace(in.Phase)
	scope := strings.TrimSpace(in.Scope)
	if goalID == "" || phase == "" || scope == "" {
		return nil, nil, fmt.Errorf("trace_plan create-coarse: goal_id, phase, and scope are required")
	}
	cp, err := psvc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: goalID,
		Actor:  "mcp",
		Phases: []planner.PhaseInput{{
			Title:  phase,
			Scopes: []planner.ScopeInput{{Title: scope}},
		}},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_plan create-coarse: %w", err)
	}
	b, err := json.Marshal(cp)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) planSetCurrent(ctx context.Context, psvc *planner.Service, in PlanInput) (*sdkmcp.CallToolResult, any, error) {
	goalID := strings.TrimSpace(in.GoalID)
	scopeID := strings.TrimSpace(in.ScopeID)
	if goalID == "" || scopeID == "" {
		return nil, nil, fmt.Errorf("trace_plan set-current: goal_id and scope_id are required")
	}
	if err := psvc.SetCurrentScope(ctx, goalID, scopeID); err != nil {
		return nil, nil, fmt.Errorf("trace_plan set-current: %w", err)
	}
	b, _ := json.Marshal(map[string]any{"ok": true, "goal_id": goalID, "scope_id": scopeID})
	return textResult(string(b)), nil, nil
}

func (s *Server) planDeep(ctx context.Context, psvc *planner.Service, in PlanInput) (*sdkmcp.CallToolResult, any, error) {
	scopeID := strings.TrimSpace(in.ScopeID)
	if scopeID == "" || len(in.Exit) == 0 {
		return nil, nil, fmt.Errorf("trace_plan deep: scope_id and exit are required")
	}
	items := make([]planner.WorkItem, 0, len(in.Work))
	for _, w := range in.Work {
		if t := strings.TrimSpace(w); t != "" {
			items = append(items, planner.WorkItem{Title: t})
		}
	}
	res, err := psvc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID:          scopeID,
		ExitCriteria:     in.Exit,
		WorkItems:        items,
		LookaheadSummary: in.Lookahead,
		Actor:            "mcp",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_plan deep: %w", err)
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) planShow(ctx context.Context, psvc *planner.Service, st *store.Store, in PlanInput) (*sdkmcp.CallToolResult, any, error) {
	goalID := strings.TrimSpace(in.GoalID)
	if goalID == "" {
		return nil, nil, fmt.Errorf("trace_plan show: goal_id is required")
	}
	view, err := psvc.GetPlan(ctx, goalID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_plan show: %w", err)
	}
	if view.Phases == nil {
		view.Phases = []planner.PhaseView{}
	}
	tasks, err := st.ListTasksByGoalID(goalID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_plan show: %w", err)
	}
	taskRows := make([]map[string]any, 0, len(tasks))
	for _, t := range tasks {
		taskRows = append(taskRows, map[string]any{
			"id": t.ID, "title": t.Title, "work_state": t.WorkState, "goal_id": t.GoalID,
		})
	}
	warning, _ := psvc.GoalStructureWarning(ctx, goalID)
	out := map[string]any{
		"goal_id":            view.GoalID,
		"current_scope_id":   view.CurrentScopeID,
		"phases":             view.Phases,
		"current_deep_plan":  view.CurrentDeepPlan,
		"lookahead_scope_id": view.LookaheadScopeID,
		"lookahead_summary":  view.LookaheadSummary,
		"tasks":              taskRows,
	}
	if warning != "" {
		out["goal_structure_warning"] = warning
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) planBootstrap(ctx context.Context, psvc *planner.Service, in PlanInput) (*sdkmcp.CallToolResult, any, error) {
	goalID := strings.TrimSpace(in.GoalID)
	if goalID == "" {
		return nil, nil, fmt.Errorf("trace_plan bootstrap: goal_id is required")
	}
	res, err := psvc.BootstrapFromPlanChanges(ctx, goalID, "mcp")
	if err != nil {
		return nil, nil, fmt.Errorf("trace_plan bootstrap: %w", err)
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}
