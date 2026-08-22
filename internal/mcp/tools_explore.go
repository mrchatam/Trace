package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/retrieval"
)

// ExploreInput is the argument schema for trace_explore.
type ExploreInput struct {
	Project  string  `json:"project,omitempty" jsonschema:"optional project root override (else server -C / cwd)"`
	TaskID   string  `json:"task_id" jsonschema:"task UUID (required)"`
	Query    string  `json:"query,omitempty" jsonschema:"optional agent query merged into task context (G1 path)"`
	Limit    float64 `json:"limit,omitempty" jsonschema:"optional FTS max hits (default 32, cap 64)"`
	MaxNodes float64 `json:"max_nodes,omitempty" jsonschema:"optional neighborhood max_nodes (default 100)"`
	Depth    float64 `json:"depth,omitempty" jsonschema:"optional neighborhood depth (default 2)"`
}

func (s *Server) toolExplore(ctx context.Context, _ *sdkmcp.CallToolRequest, in ExploreInput) (*sdkmcp.CallToolResult, any, error) {
	if strings.TrimSpace(in.TaskID) == "" {
		return nil, nil, fmt.Errorf("trace_explore: task_id is required")
	}

	abs, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_explore"); err != nil {
		return nil, nil, err
	}

	eng := retrieval.New(st)
	if repo, rerr := tryOpenGit(abs, st); rerr == nil {
		defer repo.Close()
		eng = eng.WithVCS(repo)
	}
	comp := compiler.New(st).WithRetrieval(eng)

	out, err := compiler.Explore(ctx, comp, eng, compiler.ExploreOpts{
		TaskID:   in.TaskID,
		Query:    in.Query,
		Limit:    int(in.Limit),
		MaxNodes: int(in.MaxNodes),
		Depth:    int(in.Depth),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_explore: %w", err)
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("trace_explore: encode: %w", err)
	}
	return textResult(string(b)), nil, nil
}
