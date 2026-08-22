package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// TasksInput mirrors `trace tasks [--goal]`.
type TasksInput struct {
	Project string `json:"project,omitempty" jsonschema:"optional project root override"`
	GoalID  string `json:"goal_id,omitempty" jsonschema:"optional goal UUID filter (mirrors --goal)"`
}

// CapabilityInput mirrors `trace capability declare|list|require|unrequire|missing`.
type CapabilityInput struct {
	Project    string `json:"project,omitempty" jsonschema:"optional project root override"`
	Action     string `json:"action" jsonschema:"declare|list|require|unrequire|missing"`
	Kind       string `json:"kind,omitempty" jsonschema:"SKILL|RULE|MCP|TOOL|HOOK (declare/list)"`
	Slug       string `json:"slug,omitempty" jsonschema:"capability slug (declare)"`
	Title      string `json:"title,omitempty" jsonschema:"optional title (declare)"`
	Status     string `json:"status,omitempty" jsonschema:"AVAILABLE|UNAVAILABLE|UNKNOWN (declare/list)"`
	Body       string `json:"body,omitempty" jsonschema:"optional body (declare)"`
	ID         string `json:"id,omitempty" jsonschema:"optional UUID (declare)"`
	TaskID     string `json:"task,omitempty" jsonschema:"task UUID (require|unrequire|missing); alias task_id"`
	TaskIDAlt  string `json:"task_id,omitempty" jsonschema:"alias for task"`
	Capability string `json:"capability,omitempty" jsonschema:"capability id or slug (require|unrequire)"`
}

func (in CapabilityInput) resolvedTaskID() string {
	if strings.TrimSpace(in.TaskID) != "" {
		return in.TaskID
	}
	return strings.TrimSpace(in.TaskIDAlt)
}

// VersionInput is empty; trace_version needs no parameters.
type VersionInput struct{}

// taskListRow matches CLI `trace tasks` JSON (DF-21).
type taskListRow struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	WorkState string  `json:"work_state"`
	GoalID    *string `json:"goal_id"`
}

// capabilityListRow matches CLI capability list/missing (DF-32).
type capabilityListRow struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Slug   string `json:"slug"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func capabilityListRows(list []store.Capability) []capabilityListRow {
	out := make([]capabilityListRow, 0, len(list))
	for _, c := range list {
		out = append(out, capabilityListRow{
			ID: c.ID, Kind: c.Kind, Slug: c.Slug, Title: c.Title, Status: c.Status,
		})
	}
	return out
}

func (s *Server) toolTasks(ctx context.Context, _ *sdkmcp.CallToolRequest, in TasksInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_tasks"); err != nil {
		return nil, nil, err
	}

	var tasks []store.Task
	if strings.TrimSpace(in.GoalID) != "" {
		tasks, err = st.ListTasksByGoalID(in.GoalID)
	} else {
		tasks, err = st.ListTasks()
	}
	if err != nil {
		return nil, nil, fmt.Errorf("trace_tasks: %w", err)
	}
	out := make([]taskListRow, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, taskListRow{
			ID: t.ID, Title: t.Title, WorkState: t.WorkState, GoalID: t.GoalID,
		})
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) toolCapability(ctx context.Context, _ *sdkmcp.CallToolRequest, in CapabilityInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	if err := assertMCPToolAllowed(ctx, st, "trace_capability"); err != nil {
		_ = st.Close()
		return nil, nil, err
	}
	if err := st.Close(); err != nil {
		return nil, nil, err
	}
	switch strings.ToLower(strings.TrimSpace(in.Action)) {
	case "declare":
		return s.capabilityDeclare(ctx, in)
	case "list":
		return s.capabilityList(ctx, in)
	case "require":
		return s.capabilityRequire(ctx, in)
	case "unrequire":
		return s.capabilityUnrequire(ctx, in)
	case "missing":
		return s.capabilityMissing(ctx, in)
	default:
		return nil, nil, fmt.Errorf("trace_capability: action must be declare|list|require|unrequire|missing")
	}
}

func (s *Server) capabilityDeclare(ctx context.Context, in CapabilityInput) (*sdkmcp.CallToolResult, any, error) {
	if strings.TrimSpace(in.Kind) == "" || strings.TrimSpace(in.Slug) == "" {
		return nil, nil, fmt.Errorf("trace_capability declare: kind and slug are required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	c, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		ID: in.ID, Kind: in.Kind, Slug: in.Slug, Title: in.Title, Status: in.Status, Body: in.Body,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_capability: %w", err)
	}
	b, err := json.Marshal(map[string]any{
		"ok": true, "id": c.ID, "kind": c.Kind, "slug": c.Slug,
		"title": c.Title, "status": c.Status,
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) capabilityList(ctx context.Context, in CapabilityInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	list, err := svc.ListCapabilities(ctx, domain.ListCapabilitiesFilter{
		Kind: in.Kind, Status: in.Status,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_capability: %w", err)
	}
	if list == nil {
		list = []store.Capability{}
	}
	rows := capabilityListRows(list)
	b, err := json.Marshal(map[string]any{
		"ok": true, "capabilities": rows, "count": len(rows),
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) capabilityRequire(ctx context.Context, in CapabilityInput) (*sdkmcp.CallToolResult, any, error) {
	taskID := in.resolvedTaskID()
	if taskID == "" || strings.TrimSpace(in.Capability) == "" {
		return nil, nil, fmt.Errorf("trace_capability require: task/task_id and capability are required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	cap, err := svc.ResolveCapabilityIDOrSlug(ctx, in.Capability)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_capability: %w", err)
	}
	r, err := svc.RequireCapability(ctx, taskID, cap.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_capability: %w", err)
	}
	b, err := json.Marshal(map[string]any{
		"ok": true, "id": r.ID, "task": r.TaskID, "capability": r.CapabilityID, "slug": cap.Slug,
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) capabilityUnrequire(ctx context.Context, in CapabilityInput) (*sdkmcp.CallToolResult, any, error) {
	taskID := in.resolvedTaskID()
	if taskID == "" || strings.TrimSpace(in.Capability) == "" {
		return nil, nil, fmt.Errorf("trace_capability unrequire: task/task_id and capability are required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	cap, err := svc.ResolveCapabilityIDOrSlug(ctx, in.Capability)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_capability: %w", err)
	}
	if err := svc.UnrequireCapability(ctx, taskID, cap.ID); err != nil {
		return nil, nil, fmt.Errorf("trace_capability: %w", err)
	}
	b, err := json.Marshal(map[string]any{
		"ok": true, "task": taskID, "capability": cap.ID, "slug": cap.Slug,
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) capabilityMissing(ctx context.Context, in CapabilityInput) (*sdkmcp.CallToolResult, any, error) {
	taskID := in.resolvedTaskID()
	if taskID == "" {
		return nil, nil, fmt.Errorf("trace_capability missing: task/task_id is required; list tasks via trace_tasks")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	list, err := svc.MissingCapabilities(ctx, taskID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_capability: %w", err)
	}
	if list == nil {
		list = []store.Capability{}
	}
	rows := capabilityListRows(list)
	b, err := json.Marshal(map[string]any{
		"ok": true, "task": taskID, "missing": rows, "count": len(rows),
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) toolVersion(ctx context.Context, _ *sdkmcp.CallToolRequest, _ VersionInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore("")
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_version"); err != nil {
		return nil, nil, err
	}
	b, err := json.Marshal(map[string]any{
		"ok": true, "name": serverName, "version": serverVersion,
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}
