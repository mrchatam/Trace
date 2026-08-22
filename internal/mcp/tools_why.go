package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
)

// WhyInput is the argument schema for trace_why.
type WhyInput struct {
	Project    string `json:"project,omitempty" jsonschema:"optional project root override (else server -C / cwd)"`
	EntityType string `json:"entity_type" jsonschema:"entity type (same vocabulary as trace why <type> <id>)"`
	ID         string `json:"id" jsonschema:"entity UUID"`
}

func (s *Server) toolWhy(ctx context.Context, _ *sdkmcp.CallToolRequest, in WhyInput) (*sdkmcp.CallToolResult, any, error) {
	if strings.TrimSpace(in.EntityType) == "" || strings.TrimSpace(in.ID) == "" {
		return nil, nil, fmt.Errorf("trace_why: entity_type and id are required")
	}
	abs, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_why"); err != nil {
		return nil, nil, err
	}

	eng := retrieval.New(st)
	if repo, rerr := tryOpenGit(abs, st); rerr == nil {
		defer repo.Close()
		eng = eng.WithVCS(repo)
	}

	res, err := eng.Why(ctx, in.EntityType, in.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_why: %w", err)
	}
	impact, err := domain.New(st).ImpactSummariesForWhySeed(ctx, retrieval.NormalizeEntityType(in.EntityType), in.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_why: %w", err)
	}
	out := struct {
		retrieval.WhyResult
		Impact []domain.DecisionImpact `json:"impact,omitempty"`
	}{WhyResult: res, Impact: impact}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("trace_why: encode: %w", err)
	}
	return textResult(string(b)), nil, nil
}
