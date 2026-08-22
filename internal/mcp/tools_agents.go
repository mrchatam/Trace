package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/agents"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

// AgentsInput mirrors trace agents list|describe|recommend.
type AgentsInput struct {
	Project  string `json:"project,omitempty" jsonschema:"optional project root override"`
	Action   string `json:"action" jsonschema:"list|describe|recommend"`
	Slug     string `json:"slug,omitempty" jsonschema:"describe: agent slug"`
	TaskID   string `json:"task_id,omitempty" jsonschema:"recommend: seed task UUID"`
	Phase    string `json:"phase,omitempty" jsonschema:"recommend: deliberation phase e.g. CRITIQUE"`
	Keywords string `json:"keywords,omitempty" jsonschema:"recommend: optional keyword injection"`
}

func (s *Server) toolAgents(ctx context.Context, _ *sdkmcp.CallToolRequest, in AgentsInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	if err := assertMCPToolAllowed(ctx, st, "trace_agents"); err != nil {
		_ = st.Close()
		return nil, nil, err
	}
	defer st.Close()

	switch strings.ToLower(strings.TrimSpace(in.Action)) {
	case "list":
		return s.agentsList(ctx, st)
	case "describe":
		return s.agentsDescribe(ctx, st, in.Slug)
	case "recommend":
		return s.agentsRecommend(ctx, st, in)
	default:
		return nil, nil, fmt.Errorf("trace_agents: action must be list|describe|recommend")
	}
}

func (s *Server) agentsList(ctx context.Context, st *store.Store) (*sdkmcp.CallToolResult, any, error) {
	items, err := agents.ListAgentSummaries(ctx, st)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_agents list: %w", err)
	}
	b, err := json.Marshal(items)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) agentsDescribe(ctx context.Context, st *store.Store, slug string) (*sdkmcp.CallToolResult, any, error) {
	profile, err := agents.DescribeAgent(ctx, st, slug)
	if err != nil {
		return nil, nil, err
	}
	b, err := json.Marshal(profile)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) agentsRecommend(ctx context.Context, st *store.Store, in AgentsInput) (*sdkmcp.CallToolResult, any, error) {
	dom := domain.New(st)
	recs, err := loop.RecommendHarness(ctx, dom, st, planner.New(st), loop.RecommendHarnessInput{
		TaskID:   in.TaskID,
		Phase:    in.Phase,
		Keywords: in.Keywords,
	})
	if err != nil {
		return nil, nil, err
	}
	if recs == nil {
		recs = []agents.Recommendation{}
	}
	b, err := json.Marshal(recs)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}
