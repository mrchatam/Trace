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

// SearchInput mirrors `trace search <query> [--limit N] [--matcher and|ops]`.
type SearchInput struct {
	Project string  `json:"project,omitempty" jsonschema:"optional project root override"`
	Query   string  `json:"query" jsonschema:"FTS query string"`
	Limit   float64 `json:"limit,omitempty" jsonschema:"optional max hits (default 32, cap 64)"`
	Matcher string  `json:"matcher,omitempty" jsonschema:"and|ops; default and (AND-joined terms); ops honors AND|OR|NOT|NEAR, \"phrase\", and trailing *"`
}

type searchResponse struct {
	OK               bool                       `json:"ok"`
	Hits             []retrieval.Hit            `json:"hits"`
	Count            int                        `json:"count"`
	Limit            int                        `json:"limit"`
	Matcher          string                     `json:"matcher"`
	QueryNote        string                     `json:"query_note,omitempty"`
	GraphSyncHonesty *compiler.GraphSyncHonesty `json:"graph_sync_honesty,omitempty"`
}

// searchOperatorsNotice explains the AND-mode rewrite so an operator query is
// never silently narrowed to zero hits (#149).
const searchOperatorsNotice = "query contains FTS operators (AND/OR/NOT/NEAR/*); under matcher=and they are AND-joined as search terms — pass matcher=ops to use them as FTS5 operators"

// searchEffectiveLimit clamps a requested limit to the store's 64 cap so the
// response reports what actually ran.
func searchEffectiveLimit(requested float64) int {
	limit := int(requested)
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}
	return limit
}

func (s *Server) toolSearch(ctx context.Context, _ *sdkmcp.CallToolRequest, in SearchInput) (*sdkmcp.CallToolResult, any, error) {
	matcher, err := store.ParseMatchMode(in.Matcher)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_search: %w", err)
	}

	// One open for the whole call: the assert and the query must observe the
	// same store lifetime.
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_search"); err != nil {
		return nil, nil, err
	}

	limit := searchEffectiveLimit(in.Limit)
	resp := searchResponse{
		OK:               true,
		Hits:             []retrieval.Hit{},
		Limit:            limit,
		Matcher:          matcher.String(),
		GraphSyncHonesty: compiler.GraphSyncHonestyFor(st),
	}

	query := strings.TrimSpace(in.Query)
	if query != "" {
		hits, compiled, err := retrieval.New(st).SearchWithQuery(ctx, query, retrieval.SearchOptions{
			Limit: limit, Matcher: matcher,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("trace_search: %w", err)
		}
		if hits == nil {
			hits = []retrieval.Hit{}
		}
		resp.Hits = hits
		resp.Count = len(hits)
		resp.QueryNote = searchQueryNote(matcher, compiled)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

// searchQueryNote reports any rewrite the compiler applied: an operator query
// that the AND mode turned into terms, or an ops-mode fallback. Empty when the
// query ran exactly as written.
func searchQueryNote(matcher store.MatchMode, compiled store.FTSQuery) string {
	if compiled.Fallback != "" {
		return compiled.Fallback
	}
	if compiled.SawOperator && matcher != store.MatchOps {
		return searchOperatorsNotice
	}
	return ""
}
