package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

// ExploreInput is the argument schema for trace_explore.
type ExploreInput struct {
	Project  string  `json:"project,omitempty" jsonschema:"optional project root override (else server -C / cwd)"`
	TaskID   string  `json:"task_id" jsonschema:"task UUID (required)"`
	Query    string  `json:"query,omitempty" jsonschema:"optional agent query merged into task context (G1 path)"`
	Limit    float64 `json:"limit,omitempty" jsonschema:"optional FTS max hits (default 32, cap 64)"`
	MaxNodes float64 `json:"max_nodes,omitempty" jsonschema:"optional neighborhood max_nodes (default 100)"`
	Depth    float64 `json:"depth,omitempty" jsonschema:"optional neighborhood depth (default 2)"`
	Matcher  string  `json:"matcher,omitempty" jsonschema:"and|ops; default and (AND-joined terms); ops honors AND|OR|NOT|NEAR, \"phrase\", and trailing *"`
} // exploreBounds holds the effective explore parameters. Values above the cap
// are clamped to the cap (matching the compiler/retrieval defaults), negative
// values are rejected; the clamps mirror compiler.ExploreOpts normalization so
// the tool never applies a cap the engine would not (#158).
type exploreBounds struct {
	limit    int
	maxNodes int
	depth    int
}

func resolveExploreBounds(limit, maxNodes, depth float64) (exploreBounds, error) {
	var b exploreBounds

	if limit < 0 {
		return b, fmt.Errorf("trace_explore: limit must be 1..%d", compiler.MaxExploreSearchLimit)
	}
	if maxNodes < 0 {
		return b, fmt.Errorf("trace_explore: max_nodes must be 1..%d", retrieval.MaxNeighborhoodNodes)
	}
	if depth < 0 {
		return b, fmt.Errorf("trace_explore: depth must be 1..8")
	}
	if limit == 0 {
		b.limit = compiler.DefaultExploreSearchLimit
	} else {
		b.limit = int(limit)
		if b.limit > compiler.MaxExploreSearchLimit {
			b.limit = compiler.MaxExploreSearchLimit
		}
	}
	if maxNodes == 0 {
		b.maxNodes = compiler.DefaultExploreMaxNodes
	} else {
		b.maxNodes = int(maxNodes)
		if b.maxNodes > retrieval.MaxNeighborhoodNodes {
			b.maxNodes = retrieval.MaxNeighborhoodNodes
		}
	}
	if depth == 0 {
		b.depth = compiler.DefaultExploreDepth
	} else {
		b.depth = int(depth)
		if b.depth > 8 {
			b.depth = 8
		}
	}
	return b, nil
}

func (s *Server) toolExplore(ctx context.Context, _ *sdkmcp.CallToolRequest, in ExploreInput) (*sdkmcp.CallToolResult, any, error) {
	if strings.TrimSpace(in.TaskID) == "" {
		return nil, nil, fmt.Errorf("trace_explore: task_id is required")
	}
	bounds, err := resolveExploreBounds(in.Limit, in.MaxNodes, in.Depth)
	if err != nil {
		return nil, nil, err
	}
	matcher, err := store.ParseMatchMode(in.Matcher)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_explore: %w", err)
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
		Limit:    bounds.limit,
		MaxNodes: bounds.maxNodes,
		Depth:    bounds.depth,
		Matcher:  matcher,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_explore: %w", err)
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}
