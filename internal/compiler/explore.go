package compiler

import (
	"context"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/retrieval"
)

const (
	DefaultExploreSearchLimit = 32
	MaxExploreSearchLimit     = 64
	DefaultExploreWhyTopN     = 3
	DefaultExploreMaxNodes    = 100
	DefaultExploreDepth       = 2
)

// ExploreOpts controls unified task-aware explore compose.
type ExploreOpts struct {
	TaskID   string
	Query    string
	Limit    int // FTS search limit; default 32, cap 64
	MaxNodes int // neighborhood max_nodes; default 100
	Depth    int // neighborhood depth; default 2
}

func (o ExploreOpts) withDefaults() ExploreOpts {
	if o.Limit <= 0 {
		o.Limit = DefaultExploreSearchLimit
	}
	if o.Limit > MaxExploreSearchLimit {
		o.Limit = MaxExploreSearchLimit
	}
	if o.MaxNodes <= 0 {
		o.MaxNodes = DefaultExploreMaxNodes
	}
	if o.MaxNodes > retrieval.MaxNeighborhoodNodes {
		o.MaxNodes = retrieval.MaxNeighborhoodNodes
	}
	if o.Depth <= 0 {
		o.Depth = DefaultExploreDepth
	}
	if o.Depth > 8 {
		o.Depth = 8
	}
	return o
}

// TaskSummary is the required Layer-0 task identity section.
type TaskSummary struct {
	TaskID    string `json:"task_id"`
	Title     string `json:"title,omitempty"`
	WorkState string `json:"work_state,omitempty"`
}

// ExploreResult is the unified capped explore response (task moat required).
type ExploreResult struct {
	TaskSummary  TaskSummary             `json:"task_summary"`
	PacketBudget Budget                  `json:"packet_budget"`
	SearchHits   []retrieval.Hit         `json:"search_hits"`
	WhySlices    []retrieval.WhyResult   `json:"why_slices"`
	Neighborhood *retrieval.BoundedGraph `json:"neighborhood,omitempty"`
	Truncated    bool                    `json:"truncated"`
}

// ExploreEngine is the retrieval surface used by Explore (tests may stub Search).
type ExploreEngine interface {
	Search(ctx context.Context, q string, opts retrieval.SearchOptions) ([]retrieval.Hit, error)
	SearchGraphLabels(ctx context.Context, intent retrieval.Intent, opts retrieval.SearchOptions) ([]retrieval.Hit, error)
	Why(ctx context.Context, entityType, entityID string) (retrieval.WhyResult, error)
	Neighborhood(ctx context.Context, opts retrieval.NeighborhoodOpts) (*retrieval.BoundedGraph, error)
}

// Explore composes task context (G1 query merge) + capped search + why + neighborhood.
func Explore(ctx context.Context, comp *Compiler, eng ExploreEngine, opts ExploreOpts) (ExploreResult, error) {
	opts = opts.withDefaults()
	if strings.TrimSpace(opts.TaskID) == "" {
		return ExploreResult{}, fmt.Errorf("compiler: Explore: task_id is required")
	}
	if comp == nil {
		return ExploreResult{}, fmt.Errorf("compiler: Explore: compiler is nil")
	}
	if eng == nil {
		return ExploreResult{}, fmt.Errorf("compiler: Explore: engine is nil")
	}

	pkt, err := comp.TaskContext(ctx, opts.TaskID, ContextOptions{Query: opts.Query})
	if err != nil {
		return ExploreResult{}, err
	}

	summary := TaskSummary{TaskID: opts.TaskID}
	for _, it := range pkt.Items {
		if it.EntityType == "task" && it.EntityID == opts.TaskID {
			summary.Title = it.Title
		}
		if it.EntityType == "task_state" && it.EntityID == opts.TaskID {
			summary.WorkState = it.Excerpt
		}
	}

	out := ExploreResult{
		TaskSummary:  summary,
		PacketBudget: pkt.Budget,
		SearchHits:   []retrieval.Hit{},
		WhySlices:    []retrieval.WhyResult{},
		Truncated:    pkt.Budget.Truncated || pkt.Budget.CandidatesCapped,
	}

	searchQ := strings.TrimSpace(opts.Query)
	if searchQ == "" {
		searchQ = strings.TrimSpace(summary.Title)
	}
	if searchQ != "" {
		intentIn := comp.taskIntentInput(opts.TaskID, opts.Query)
		searchOpts := retrieval.SearchOptions{Limit: opts.Limit, Intent: &intentIn}
		hits, serr := eng.Search(ctx, searchQ, searchOpts)
		if serr != nil {
			hits = nil
		}
		if hits == nil {
			hits = []retrieval.Hit{}
		}
		out.SearchHits = hits
		labelHits, lerr := eng.SearchGraphLabels(ctx, retrieval.ExtractIntent(intentIn), searchOpts)
		if lerr != nil {
			labelHits = nil
		}
		out.SearchHits = retrieval.MergeConceptHits(out.SearchHits, labelHits)
		if len(out.SearchHits) >= opts.Limit {
			out.Truncated = true
		}
	}

	topN := DefaultExploreWhyTopN
	if topN > len(out.SearchHits) {
		topN = len(out.SearchHits)
	}
	for i := 0; i < topN; i++ {
		h := out.SearchHits[i]
		why, werr := eng.Why(ctx, h.EntityType, h.EntityID)
		if werr != nil {
			continue
		}
		out.WhySlices = append(out.WhySlices, why)
	}

	nb, nerr := eng.Neighborhood(ctx, retrieval.NeighborhoodOpts{
		Center:   opts.TaskID,
		MaxNodes: opts.MaxNodes,
		Depth:    opts.Depth,
	})
	if nerr == nil && nb != nil {
		if nb.Nodes == nil {
			nb.Nodes = []retrieval.GraphNode{}
		}
		if nb.Edges == nil {
			nb.Edges = []retrieval.GraphEdge{}
		}
		out.Neighborhood = nb
		if nb.Truncated {
			out.Truncated = true
		}
	}

	return out, nil
}
