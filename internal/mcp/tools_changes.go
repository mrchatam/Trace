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

// ChangesInput mirrors `trace changes list|show|compare`.
type ChangesInput struct {
	Project  string  `json:"project,omitempty" jsonschema:"optional project root override"`
	Action   string  `json:"action" jsonschema:"list|show|compare"`
	TaskID   string  `json:"task_id,omitempty" jsonschema:"list: optional task UUID filter"`
	Limit    float64 `json:"limit,omitempty" jsonschema:"list: optional max rows (default 32, cap 64)"`
	ChangeID string  `json:"change_id,omitempty" jsonschema:"show: change UUID"`
	From     string  `json:"from,omitempty" jsonschema:"compare: start git OID"`
	To       string  `json:"to,omitempty" jsonschema:"compare: end git OID"`
}

type changeListRow struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	GitCommit  string `json:"git_commit"`
	Status     string `json:"status"`
	Reason     string `json:"reason"`
	SourceType string `json:"source_type"`
	CreatedAt  string `json:"created_at"`
}

type changeShowPathRow struct {
	Path     string `json:"path"`
	Status   string `json:"status"`
	SymbolID string `json:"symbol_id"`
}

type changeShowResponse struct {
	OK     bool                `json:"ok"`
	Change store.Change        `json:"change"`
	Paths  []changeShowPathRow `json:"paths"`
}

func (s *Server) toolChanges(ctx context.Context, _ *sdkmcp.CallToolRequest, in ChangesInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	if err := assertMCPToolAllowed(ctx, st, "trace_changes"); err != nil {
		_ = st.Close()
		return nil, nil, err
	}
	if err := st.Close(); err != nil {
		return nil, nil, err
	}

	switch strings.ToLower(strings.TrimSpace(in.Action)) {
	case "list":
		return s.changesList(ctx, in)
	case "show":
		return s.changesShow(ctx, in)
	case "compare":
		return s.changesCompare(ctx, in)
	default:
		return nil, nil, fmt.Errorf("trace_changes: action must be list|show|compare")
	}
}

func (s *Server) changesList(ctx context.Context, in ChangesInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	changes, err := svc.ListChanges(ctx, int(in.Limit), in.TaskID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_changes list: %w", err)
	}
	out := make([]changeListRow, 0, len(changes))
	for _, c := range changes {
		out = append(out, changeListRow{
			ID:         c.ID,
			TaskID:     c.TaskID,
			GitCommit:  c.GitCommit,
			Status:     c.Status,
			Reason:     c.Reason,
			SourceType: c.SourceType,
			CreatedAt:  c.CreatedAt,
		})
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) changesShow(ctx context.Context, in ChangesInput) (*sdkmcp.CallToolResult, any, error) {
	changeID := strings.TrimSpace(in.ChangeID)
	if changeID == "" {
		return nil, nil, fmt.Errorf("trace_changes show: change_id is required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	c, err := svc.GetChange(ctx, changeID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_changes show: %w", err)
	}
	paths, err := svc.ListChangePaths(ctx, changeID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_changes show: %w", err)
	}
	pathRows := make([]changeShowPathRow, 0, len(paths))
	for _, p := range paths {
		pathRows = append(pathRows, changeShowPathRow{
			Path:     p.Path,
			Status:   p.Status,
			SymbolID: p.SymbolID,
		})
	}
	resp := changeShowResponse{OK: true, Change: c, Paths: pathRows}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) changesCompare(ctx context.Context, in ChangesInput) (*sdkmcp.CallToolResult, any, error) {
	from := strings.TrimSpace(in.From)
	to := strings.TrimSpace(in.To)
	if from == "" || to == "" {
		return nil, nil, fmt.Errorf("trace_changes compare: from and to are required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	result, err := domain.New(st).CompareStates(ctx, from, to)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_changes compare: %w", err)
	}
	b, err := json.Marshal(result)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}
