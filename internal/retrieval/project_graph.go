package retrieval

import (
	"context"
	"fmt"

	"github.com/mrchatam/Trace/internal/store"
)

// ProjectGraphOpts controls ProjectGraph (mode=project on GET /v1/graph).
type ProjectGraphOpts struct {
	MaxNodes int
}

var projectGraphCountTables = []string{
	"goals", "tasks", "decisions", "assumptions", "discoveries", "plan_changes",
	"claims", "evidence", "reviews", "capabilities", "changes", "regressions",
}

// ProjectGraph returns a bounded view of project entities and edges between them.
// max_nodes is required (1..5000). Truncated=true when total entities exceed the budget.
// Nodes are collected in kind order and stop once MaxNodes is filled; tasks use SQL LIMIT.
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

	all, total, err := e.collectProjectNodes(opts.MaxNodes)
	if err != nil {
		return nil, err
	}
	truncated := total > opts.MaxNodes

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

func (e *Engine) collectProjectNodes(maxNodes int) ([]GraphNode, int, error) {
	total := 0
	for _, table := range projectGraphCountTables {
		n, err := e.store.CountInTable(table)
		if err != nil {
			return nil, 0, err
		}
		total += n
	}

	var nodes []GraphNode
	appendNode := func(n GraphNode) bool {
		if len(nodes) >= maxNodes {
			return false
		}
		nodes = append(nodes, n)
		return len(nodes) < maxNodes
	}
	remaining := func() int {
		r := maxNodes - len(nodes)
		if r < 0 {
			return 0
		}
		return r
	}

	goals, err := e.store.ListGoals()
	if err != nil {
		return nil, 0, err
	}
	for _, g := range goals {
		if !appendNode(GraphNode{ID: g.ID, Kind: "goal", Title: g.Title}) {
			return nodes, total, nil
		}
	}

	if rem := remaining(); rem > 0 {
		res, err := e.store.ListTasksFiltered(store.TaskListFilter{Limit: rem})
		if err != nil {
			return nil, 0, err
		}
		for _, t := range res.Tasks {
			gn := GraphNode{ID: t.ID, Kind: "task", Title: t.Title}
			if t.GoalID != nil && *t.GoalID != "" {
				gn.GoalID = *t.GoalID
			}
			if !appendNode(gn) {
				return nodes, total, nil
			}
		}
	}

	decisions, err := e.store.ListDecisions()
	if err != nil {
		return nil, 0, err
	}
	for _, d := range decisions {
		if !appendNode(GraphNode{ID: d.ID, Kind: "decision", Title: d.Title}) {
			return nodes, total, nil
		}
	}

	assumptions, err := e.store.ListAssumptions()
	if err != nil {
		return nil, 0, err
	}
	for _, a := range assumptions {
		if !appendNode(GraphNode{ID: a.ID, Kind: "assumption", Title: a.Title}) {
			return nodes, total, nil
		}
	}

	discoveries, err := e.store.ListDiscoveries()
	if err != nil {
		return nil, 0, err
	}
	for _, d := range discoveries {
		if !appendNode(GraphNode{ID: d.ID, Kind: "discovery", Title: d.Title}) {
			return nodes, total, nil
		}
	}

	planChanges, err := e.store.ListPlanChanges()
	if err != nil {
		return nil, 0, err
	}
	for _, p := range planChanges {
		if !appendNode(GraphNode{ID: p.ID, Kind: "plan_change", Title: p.Title}) {
			return nodes, total, nil
		}
	}

	claims, err := e.store.ListClaims()
	if err != nil {
		return nil, 0, err
	}
	for _, c := range claims {
		if !appendNode(GraphNode{ID: c.ID, Kind: "claim", Title: c.Title}) {
			return nodes, total, nil
		}
	}

	evidence, err := e.store.ListEvidence()
	if err != nil {
		return nil, 0, err
	}
	for _, ev := range evidence {
		if !appendNode(GraphNode{ID: ev.ID, Kind: "evidence", Title: ev.Title}) {
			return nodes, total, nil
		}
	}

	reviews, err := e.store.ListReviews()
	if err != nil {
		return nil, 0, err
	}
	for _, r := range reviews {
		if !appendNode(GraphNode{ID: r.ID, Kind: "review", Title: r.Title}) {
			return nodes, total, nil
		}
	}

	if rem := remaining(); rem > 0 {
		caps, err := e.store.ListCapabilities(store.CapabilityListFilter{Limit: rem})
		if err != nil {
			return nil, 0, err
		}
		for _, c := range caps {
			if !appendNode(GraphNode{ID: c.ID, Kind: "capability", Title: c.Title}) {
				return nodes, total, nil
			}
		}
	}

	changes, err := e.store.ListAllChanges()
	if err != nil {
		return nil, 0, err
	}
	for _, c := range changes {
		title := c.Reason
		if title == "" {
			title = c.ID
		}
		if !appendNode(GraphNode{ID: c.ID, Kind: "change", Title: title}) {
			return nodes, total, nil
		}
	}

	regressions, err := e.store.ListAllRegressions()
	if err != nil {
		return nil, 0, err
	}
	for _, r := range regressions {
		title := r.Summary
		if title == "" {
			title = r.ID
		}
		if !appendNode(GraphNode{ID: r.ID, Kind: "regression", Title: title}) {
			return nodes, total, nil
		}
	}

	return nodes, total, nil
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
