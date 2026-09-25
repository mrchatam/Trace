package retrieval

import (
	"context"
	"fmt"
)

const MaxNeighborhoodNodes = 5000

// BoundedGraph is a hop-limited neighborhood for HTTP GET /v1/graph.
type BoundedGraph struct {
	Mode          string      `json:"mode,omitempty"`
	Center        string      `json:"center"`
	MaxNodes      int         `json:"max_nodes"`
	TotalEntities int         `json:"total_entities,omitempty"`
	Nodes         []GraphNode `json:"nodes"`
	Edges         []GraphEdge `json:"edges"`
	Truncated     bool        `json:"truncated"`
}

// GraphNode is one node in a bounded graph response.
type GraphNode struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Title   string `json:"title"`
	GoalID  string `json:"goal_id,omitempty"`
	ScopeID string `json:"scope_id,omitempty"` // optional; from scope_member when present
}

// GraphEdge is one edge in a bounded graph response.
type GraphEdge struct {
	Rel        string `json:"rel"`
	From       string `json:"from"`
	To         string `json:"to"`
	Provenance string `json:"provenance,omitempty"` // explicit | inferred (from link source_type)
}

// NeighborhoodOpts controls Neighborhood.
type NeighborhoodOpts struct {
	Center   string
	MaxNodes int
	Depth    int // default 2; capped at 8 for HTTP
}

// ErrBudgetExceeded is returned when the neighborhood would exceed MaxNodes
// before producing a truncated response — callers may prefer 400 BUDGET_EXCEEDED.
// Neighborhood itself returns Truncated=true when it stops early; this error is
// for hard rejects (e.g. max_nodes > MaxNeighborhoodNodes).
type ErrBudgetExceeded struct {
	Message string
}

func (e *ErrBudgetExceeded) Error() string { return e.Message }

// Neighborhood returns a hop-limited causal neighborhood around center.
// Center must be a known entity UUID (type inferred by probing store).
// max_nodes is required (1..5000). depth defaults to 2 (max 8).
func (e *Engine) Neighborhood(ctx context.Context, opts NeighborhoodOpts) (*BoundedGraph, error) {
	_ = ctx
	center := opts.Center
	if center == "" {
		return nil, fmt.Errorf("retrieval: Neighborhood: center is required")
	}
	if opts.MaxNodes < 1 {
		return nil, fmt.Errorf("retrieval: Neighborhood: max_nodes is required and must be >= 1")
	}
	if opts.MaxNodes > MaxNeighborhoodNodes {
		return nil, &ErrBudgetExceeded{
			Message: fmt.Sprintf("retrieval: Neighborhood: max_nodes %d exceeds hard cap %d", opts.MaxNodes, MaxNeighborhoodNodes),
		}
	}
	depth := opts.Depth
	if depth <= 0 {
		depth = 2
	}
	if depth > 8 {
		depth = 8
	}

	seedType, seedHit, err := e.resolveCenter(center)
	if err != nil {
		return nil, err
	}
	_ = seedType

	type frontierItem struct {
		h Hit
	}
	seen := map[string]Hit{}
	edges := make([]GraphEdge, 0)
	edgeSeen := map[string]struct{}{}

	addEdge := func(ed GraphEdge) {
		k := ed.Rel + "\x00" + ed.From + "\x00" + ed.To
		if _, ok := edgeSeen[k]; ok {
			return
		}
		edgeSeen[k] = struct{}{}
		edges = append(edges, ed)
	}

	seedHit.Distance = 0
	seen[hitKey(seedHit.EntityType, seedHit.EntityID)] = seedHit
	frontier := []frontierItem{{h: seedHit}}
	truncated := false

	for d := 1; d <= depth && !truncated; d++ {
		var next []frontierItem
		for _, fi := range frontier {
			neighbors, err := e.graphWalkNeighbors(fi.h)
			if err != nil {
				return nil, err
			}
			for _, nb := range neighbors {
				addEdge(nb.edge)
				nh := nb.neighbor
				nh.Distance = d
				k := hitKey(nh.EntityType, nh.EntityID)
				if _, ok := seen[k]; ok {
					continue
				}
				if len(seen) >= opts.MaxNodes {
					truncated = true
					break
				}
				seen[k] = nh
				next = append(next, frontierItem{h: nh})
			}
			if truncated {
				break
			}
		}
		frontier = next
	}

	memberOf, _ := e.memberScopeIndex()

	nodes := make([]GraphNode, 0, len(seen))
	for _, h := range seen {
		gn := GraphNode{
			ID:    h.EntityID,
			Kind:  h.EntityType,
			Title: h.Title,
		}
		if h.EntityType == "task" {
			t, err := e.store.GetTask(h.EntityID)
			if err == nil && t.GoalID != nil && *t.GoalID != "" {
				gn.GoalID = *t.GoalID
			}
		}
		if sid, ok := memberOf[h.EntityID]; ok {
			gn.ScopeID = sid
		}
		nodes = append(nodes, gn)
	}

	return &BoundedGraph{
		Center:    center,
		MaxNodes:  opts.MaxNodes,
		Nodes:     nodes,
		Edges:     edges,
		Truncated: truncated,
	}, nil
}

func (e *Engine) resolveCenter(id string) (entityType string, hit Hit, err error) {
	types := []string{
		"task", "goal", "decision", "assumption", "discovery", "plan_change",
		"claim", "evidence", "review", "capability", "change", "regression",
	}
	for _, typ := range types {
		h, lerr := e.lookupEntity(typ, id, ReasonExactID, 0, 1.0)
		if lerr != nil {
			if isNotFound(lerr) {
				continue
			}
			return "", Hit{}, lerr
		}
		return typ, h, nil
	}
	return "", Hit{}, fmt.Errorf("retrieval: Neighborhood: center %q not found", id)
}
