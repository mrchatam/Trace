package mcp

import (
	"context"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/retrieval"
)

// ContextInput is the argument schema for trace_context.
type ContextInput struct {
	Project    string  `json:"project,omitempty" jsonschema:"optional project root override (else server -C / cwd)"`
	TaskID     string  `json:"task_id" jsonschema:"task UUID"`
	Depth      float64 `json:"depth,omitempty" jsonschema:"expand depth 1=TaskContext 2=ExpandContext; default 1; max 2"`
	MaxLayer   float64 `json:"max_layer,omitempty" jsonschema:"progressive layer ceiling 1=L0–L1 (default) 2|3=opt-in deeper layers"`
	Format     string  `json:"format,omitempty" jsonschema:"json|markdown|both; default json"`
	IncludeWhy bool    `json:"include_why,omitempty" jsonschema:"include why_trace; default false"`
	Query      string  `json:"query,omitempty" jsonschema:"optional agent query merged into context (task moat preserved)"`
}

func (s *Server) toolContext(ctx context.Context, _ *sdkmcp.CallToolRequest, in ContextInput) (*sdkmcp.CallToolResult, any, error) {
	if strings.TrimSpace(in.TaskID) == "" {
		return nil, nil, fmt.Errorf("trace_context: task_id is required")
	}
	depth := int(in.Depth)
	if in.Depth == 0 {
		depth = 1
	}
	if depth < 1 || depth > 2 {
		return nil, nil, fmt.Errorf("trace_context: depth must be 1 or 2")
	}
	maxLayer := int(in.MaxLayer)
	if in.MaxLayer == 0 {
		maxLayer = 1
	}
	if maxLayer < 1 || maxLayer > 3 {
		return nil, nil, fmt.Errorf("trace_context: max_layer must be 1, 2, or 3")
	}
	format := in.Format
	if format == "" {
		format = "json"
	}
	switch format {
	case "json", "markdown", "both":
	default:
		return nil, nil, fmt.Errorf("trace_context: format must be json|markdown|both")
	}

	abs, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_context"); err != nil {
		return nil, nil, err
	}

	eng := retrieval.New(st)
	if repo, rerr := tryOpenGit(abs, st); rerr == nil {
		defer repo.Close()
		eng = eng.WithVCS(repo)
	}
	comp := compiler.New(st).WithRetrieval(eng)
	opts := compiler.ContextOptions{
		MaxLayer:        maxLayer,
		IncludeWhy:      in.IncludeWhy,
		IncludeMarkdown: format == "markdown" || format == "both",
		Query:           in.Query,
	}

	var pkt compiler.Packet
	if depth == 1 {
		pkt, err = comp.TaskContext(ctx, in.TaskID, opts)
	} else {
		pkt, err = comp.ExpandContext(ctx, in.TaskID, depth, opts)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("trace_context: %w", err)
	}

	var out string
	switch format {
	case "json":
		b, err := pkt.JSON()
		if err != nil {
			return nil, nil, fmt.Errorf("trace_context: %w", err)
		}
		out = string(b)
	case "markdown":
		out = pkt.Markdown()
	case "both":
		b, err := pkt.JSON()
		if err != nil {
			return nil, nil, fmt.Errorf("trace_context: %w", err)
		}
		out = string(b) + "\n---\n" + pkt.Markdown()
	}
	return textResult(out), nil, nil
}
