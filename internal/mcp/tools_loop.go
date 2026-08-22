package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/retrieval"
)

// LoopInput mirrors `trace loop next|apply|status|gate`.
type LoopInput struct {
	Project  string          `json:"project,omitempty" jsonschema:"optional project root override"`
	Action   string          `json:"action" jsonschema:"next|apply|status|gate"`
	TaskID   string          `json:"task_id,omitempty" jsonschema:"next|status|gate: seed task UUID"`
	GoalID   string          `json:"goal_id,omitempty" jsonschema:"status: optional seed goal UUID (derived from task when omitted)"`
	For      string          `json:"for,omitempty" jsonschema:"gate: orient|edit|execute|done|export (default edit)"`
	Envelope json.RawMessage `json:"envelope,omitempty" jsonschema:"apply: trace.loop.apply.v1 JSON object or string"`
}

const mcpGateSchemaVersion = "trace.loop.gate.v1"

type mcpGateEnvelope struct {
	SchemaVersion    string           `json:"schema_version"`
	TaskID           string           `json:"task_id"`
	For              string           `json:"for"`
	Allowed          bool             `json:"allowed"`
	RecommendedPhase string           `json:"recommended_phase,omitempty"`
	ReasonCode       string           `json:"reason_code,omitempty"`
	Violations       []loop.Violation `json:"violations"`
}

func (s *Server) toolLoop(ctx context.Context, _ *sdkmcp.CallToolRequest, in LoopInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	if err := assertMCPToolAllowed(ctx, st, "trace_loop"); err != nil {
		_ = st.Close()
		return nil, nil, err
	}
	if err := st.Close(); err != nil {
		return nil, nil, err
	}

	switch strings.ToLower(strings.TrimSpace(in.Action)) {
	case "next":
		return s.loopNext(ctx, in)
	case "apply":
		return s.loopApply(ctx, in)
	case "status":
		return s.loopStatus(ctx, in)
	case "gate":
		return s.loopGate(ctx, in)
	default:
		return nil, nil, fmt.Errorf("trace_loop: action must be next|apply|status|gate")
	}
}

func (s *Server) loopNext(ctx context.Context, in LoopInput) (*sdkmcp.CallToolResult, any, error) {
	taskID := strings.TrimSpace(in.TaskID)
	if taskID == "" {
		return nil, nil, fmt.Errorf("trace_loop next: task_id is required")
	}

	abs, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()

	eng := retrieval.New(st)
	if repo, rerr := tryOpenGit(abs, st); rerr == nil {
		defer repo.Close()
		eng = eng.WithVCS(repo)
	}
	comp := compiler.New(st).WithRetrieval(eng)

	packet, err := loop.BuildNextPacket(ctx, loop.BuildNextInput{
		TaskID:    taskID,
		Store:     st,
		Planner:   planner.New(st),
		Retrieval: eng,
		Compiler:  comp,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_loop next: %w", err)
	}
	b, err := json.Marshal(packet)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) loopApply(ctx context.Context, in LoopInput) (*sdkmcp.CallToolResult, any, error) {
	raw, err := parseLoopEnvelope(in.Envelope)
	if err != nil {
		return nil, nil, err
	}
	env, err := loop.ParseApplyEnvelope(raw)
	if err != nil {
		return nil, nil, err
	}

	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()

	res, err := loop.Apply(ctx, st, planner.New(st), env)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_loop apply: %w", err)
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) loopStatus(ctx context.Context, in LoopInput) (*sdkmcp.CallToolResult, any, error) {
	taskID := strings.TrimSpace(in.TaskID)
	if taskID == "" {
		return nil, nil, fmt.Errorf("trace_loop status: task_id is required")
	}

	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()

	res, err := loop.Status(ctx, st, planner.New(st), loop.ApplySeed{
		TaskID: taskID,
		GoalID: strings.TrimSpace(in.GoalID),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_loop status: %w", err)
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) loopGate(ctx context.Context, in LoopInput) (*sdkmcp.CallToolResult, any, error) {
	taskID := strings.TrimSpace(in.TaskID)
	if taskID == "" {
		return nil, nil, fmt.Errorf("trace_loop gate: task_id is required")
	}
	gateFor := strings.TrimSpace(in.For)
	if gateFor == "" {
		gateFor = "edit"
	}
	gf, err := parseMCPLoopGateFor(gateFor)
	if err != nil {
		return nil, nil, err
	}

	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()

	allowed, violations, err := loop.EvaluateGate(ctx, domain.New(st), planner.New(st), st, taskID, gf)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_loop gate: %w", err)
	}
	env := mcpGateEnvelope{
		SchemaVersion: mcpGateSchemaVersion,
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
	b, err := json.Marshal(env)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func parseMCPLoopGateFor(s string) (loop.GateFor, error) {
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
		return "", fmt.Errorf("trace_loop gate: invalid for %q (want orient|edit|execute|done|export)", s)
	}
}

func parseLoopEnvelope(raw json.RawMessage) ([]byte, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("trace_loop apply: envelope is required")
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, fmt.Errorf("trace_loop apply: envelope is required")
		}
		return []byte(s), nil
	}
	return raw, nil
}
