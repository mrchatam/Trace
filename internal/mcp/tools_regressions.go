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
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_regressions"); err != nil {
		return nil, nil, err
	}

	if strings.ToLower(strings.TrimSpace(in.Action)) != "list" {
		return nil, nil, fmt.Errorf("trace_regressions: action must be list")
	}
	return s.regressionsList(ctx, st, in)
}

func (s *Server) regressionsList(ctx context.Context, st *store.Store, in RegressionsInput) (*sdkmcp.CallToolResult, any, error) {
	svc := domain.New(st)
	taskID := strings.TrimSpace(in.TaskID)
	if taskID != "" {
		// Fail closed on a task that does not exist (#154 class): an empty
		// regression list is only authoritative when the task is real.
		if _, err := svc.GetTask(ctx, taskID); err != nil {
			return nil, nil, wrapToolErr("trace_regressions list", err)
		}
	}
	rows, err := svc.ListRegressions(ctx, domain.EvidenceQueryOpts{
		TaskID: in.TaskID, ChangeID: in.ChangeID, Limit: int(in.Limit),
	})
	if err != nil {
		return nil, nil, wrapToolErr("trace_regressions list", err)
	}
	effLimit := clampEvidenceLimitForMCP(int(in.Limit))
	// truncated is the page-fill signal: the page returned exactly the
	// effective limit, so there may be more rows. The store has no cursor
	// plumbing, so a definitive answer is not available here.
	b, err := json.Marshal(newListEnvelope(rows, len(rows), len(rows) == effLimit).withLimit(effLimit))
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

// clampEvidenceLimitForMCP mirrors domain.clampEvidenceLimit so the response
// can report the effective limit rather than silently clamping.
func clampEvidenceLimitForMCP(limit int) int {
	switch {
	case limit <= 0:
		return 32
	case limit > 64:
		return 64
	default:
		return limit
	}
}
