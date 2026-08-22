package retrieval

import (
	"context"
	"fmt"
	"sort"

	"github.com/mrchatam/Trace/internal/store"
)

// ProjectGraphOpts controls ProjectGraph (mode=project on GET /v1/graph).
type ProjectGraphOpts struct {
	MaxNodes int
}

var projectGraphKindOrder = map[string]int{
	"goal": 0, "task": 1, "decision": 2, "assumption": 3, "discovery": 4,
	"plan_change": 5, "claim": 6, "evidence": 7, "review": 8, "capability": 9,
	"change": 10, "regression": 11,
}

// ProjectGraph returns a bounded view of all project entities and edges between them.
// max_nodes is required (1..5000). Truncated=true when total entities exceed the budget.
func (e *Engine) ProjectGraph(ctx context.Context, opts ProjectGraphOpts) (*BoundedGraph, error) {
	_ = ctx
	if opts.MaxNodes < 1 {
		return nil, fmt.Errorf("retrieval: ProjectGraph: max_nodes is required and must be >= 1")
	}
	if opts.MaxNodes > MaxNeighborhoodNodes {
		return nil, &ErrBudgetExceeded{
			Message: fmt.Sprintf("retrieval: ProjectGraph: max_nodes %d exceeds hard cap %d", opts.MaxNodes, MaxNeighborhoodNodes),
		}
	}

	all, err := e.collectProjectNodes()
	if err != nil {
		return nil, err
	}
	total := len(all)
	truncated := total > opts.MaxNodes
	if truncated {
		all = all[:opts.MaxNodes]
	}

	included := make(map[string]struct{}, len(all))
	for _, n := range all {
		included[n.ID] = struct{}{}
	}

	edges, err := e.collectEdgesForNodes(all, included)
	if err != nil {
		return nil, err
	}

	center := ""
	for _, n := range all {
		if n.Kind == "goal" {
			center = n.ID
			break
		}
	}
	if center == "" && len(all) > 0 {
		center = all[0].ID
	}

	return &BoundedGraph{
		Mode:          "project",
		Center:        center,
		MaxNodes:      opts.MaxNodes,
		TotalEntities: total,
		Nodes:         all,
		Edges:         edges,
		Truncated:     truncated,
	}, nil
}

func (e *Engine) collectProjectNodes() ([]GraphNode, error) {
	var nodes []GraphNode

	goals, err := e.store.ListGoals()
	if err != nil {
		return nil, err
	}
	for _, g := range goals {
		nodes = append(nodes, GraphNode{ID: g.ID, Kind: "goal", Title: g.Title})
	}

	tasks, err := e.store.ListTasks()
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		gn := GraphNode{ID: t.ID, Kind: "task", Title: t.Title}
		if t.GoalID != nil && *t.GoalID != "" {
			gn.GoalID = *t.GoalID
		}
		nodes = append(nodes, gn)
	}

	decisions, err := e.store.ListDecisions()
	if err != nil {
		return nil, err
	}
	for _, d := range decisions {
		nodes = append(nodes, GraphNode{ID: d.ID, Kind: "decision", Title: d.Title})
	}

	assumptions, err := e.store.ListAssumptions()
	if err != nil {
		return nil, err
	}
	for _, a := range assumptions {
		nodes = append(nodes, GraphNode{ID: a.ID, Kind: "assumption", Title: a.Title})
	}

	discoveries, err := e.store.ListDiscoveries()
	if err != nil {
		return nil, err
	}
	for _, d := range discoveries {
		nodes = append(nodes, GraphNode{ID: d.ID, Kind: "discovery", Title: d.Title})
	}

	planChanges, err := e.store.ListPlanChanges()
	if err != nil {
		return nil, err
	}
	for _, p := range planChanges {
		nodes = append(nodes, GraphNode{ID: p.ID, Kind: "plan_change", Title: p.Title})
	}

	claims, err := e.store.ListClaims()
	if err != nil {
		return nil, err
	}
	for _, c := range claims {
		nodes = append(nodes, GraphNode{ID: c.ID, Kind: "claim", Title: c.Title})
	}

	evidence, err := e.store.ListEvidence()
	if err != nil {
		return nil, err
	}
	for _, ev := range evidence {
		nodes = append(nodes, GraphNode{ID: ev.ID, Kind: "evidence", Title: ev.Title})
	}

	reviews, err := e.store.ListReviews()
	if err != nil {
		return nil, err
	}
	for _, r := range reviews {
		nodes = append(nodes, GraphNode{ID: r.ID, Kind: "review", Title: r.Title})
	}

	caps, err := e.store.ListCapabilities(store.CapabilityListFilter{})
	if err != nil {
		return nil, err
	}
	for _, c := range caps {
		nodes = append(nodes, GraphNode{ID: c.ID, Kind: "capability", Title: c.Title})
	}

	changes, err := e.store.ListAllChanges()
	if err != nil {
		return nil, err
	}
	for _, c := range changes {
		title := c.Reason
		if title == "" {
			title = c.ID
		}
		nodes = append(nodes, GraphNode{ID: c.ID, Kind: "change", Title: title})
	}

	regressions, err := e.store.ListAllRegressions()
	if err != nil {
		return nil, err
	}
	for _, r := range regressions {
		title := r.Summary
		if title == "" {
			title = r.ID
		}
		nodes = append(nodes, GraphNode{ID: r.ID, Kind: "regression", Title: title})
	}

	sort.SliceStable(nodes, func(i, j int) bool {
		oi, oki := projectGraphKindOrder[nodes[i].Kind]
		oj, okj := projectGraphKindOrder[nodes[j].Kind]
		if oki && okj && oi != oj {
			return oi < oj
		}
		if oki != okj {
			return oki
		}
		return nodes[i].ID < nodes[j].ID
	})

	return nodes, nil
}

func (e *Engine) collectEdgesForNodes(nodes []GraphNode, included map[string]struct{}) ([]GraphEdge, error) {
	edgeSeen := map[string]struct{}{}
	var edges []GraphEdge

	addEdge := func(rel, from, to string) {
		if _, ok := included[from]; !ok {
			return
		}
		if _, ok := included[to]; !ok {
			return
		}
		k := rel + "\x00" + from + "\x00" + to
		if _, ok := edgeSeen[k]; ok {
			return
		}
		edgeSeen[k] = struct{}{}
		edges = append(edges, GraphEdge{Rel: rel, From: from, To: to})
	}

	for _, n := range nodes {
		h := Hit{EntityType: n.Kind, EntityID: n.ID, Title: n.Title}
		neighbors, err := e.graphWalkNeighbors(h)
		if err != nil {
			return nil, err
		}
		for _, nb := range neighbors {
			addEdge(nb.edge.Rel, nb.edge.From, nb.edge.To)
		}
	}

	return edges, nil
}
