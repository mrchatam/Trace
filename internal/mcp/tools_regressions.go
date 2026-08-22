package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/domain"
)

// RegressionsInput mirrors `trace regressions list`.
type RegressionsInput struct {
	Project  string  `json:"project,omitempty" jsonschema:"optional project root override"`
	Action   string  `json:"action" jsonschema:"list"`
	TaskID   string  `json:"task_id,omitempty" jsonschema:"optional task UUID filter"`
	ChangeID string  `json:"change_id,omitempty" jsonschema:"optional change UUID filter"`
	Limit    float64 `json:"limit,omitempty" jsonschema:"optional max rows (default 32, cap 64)"`
}

func (s *Server) toolRegressions(ctx context.Context, _ *sdkmcp.CallToolRequest, in RegressionsInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	if err := assertMCPToolAllowed(ctx, st, "trace_regressions"); err != nil {
		_ = st.Close()
		return nil, nil, err
	}
	if err := st.Close(); err != nil {
		return nil, nil, err
	}

	if strings.ToLower(strings.TrimSpace(in.Action)) != "list" {
		return nil, nil, fmt.Errorf("trace_regressions: action must be list")
	}
	return s.regressionsList(ctx, in)
}

func (s *Server) regressionsList(ctx context.Context, in RegressionsInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	rows, err := svc.ListRegressions(ctx, domain.EvidenceQueryOpts{
		TaskID: in.TaskID, ChangeID: in.ChangeID, Limit: int(in.Limit),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_regressions list: %w", err)
	}
	b, err := json.Marshal(map[string]any{"ok": true, "regressions": rows, "count": len(rows)})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}
