package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/retrieval"
)

// SearchInput mirrors `trace search <query> [--limit N]`.
type SearchInput struct {
	Project string  `json:"project,omitempty" jsonschema:"optional project root override"`
	Query   string  `json:"query" jsonschema:"FTS query string"`
	Limit   float64 `json:"limit,omitempty" jsonschema:"optional max hits (default 32, cap 64)"`
}

type searchResponse struct {
	OK    bool            `json:"ok"`
	Hits  []retrieval.Hit `json:"hits"`
	Count int             `json:"count"`
}

func (s *Server) toolSearch(ctx context.Context, _ *sdkmcp.CallToolRequest, in SearchInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	if err := assertMCPToolAllowed(ctx, st, "trace_search"); err != nil {
		_ = st.Close()
		return nil, nil, err
	}
	if err := st.Close(); err != nil {
		return nil, nil, err
	}

	query := strings.TrimSpace(in.Query)
	resp := searchResponse{OK: true, Hits: []retrieval.Hit{}, Count: 0}
	if query != "" {
		_, st2, err := s.openStore(in.Project)
		if err != nil {
			return nil, nil, err
		}
		defer st2.Close()
		limit := int(in.Limit)
		hits, err := retrieval.New(st2).Search(ctx, query, retrieval.SearchOptions{Limit: limit})
		if err != nil {
			return nil, nil, fmt.Errorf("trace_search: %w", err)
		}
		if hits == nil {
			hits = []retrieval.Hit{}
		}
		resp.Hits = hits
		resp.Count = len(hits)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}
