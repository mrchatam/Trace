package retrieval

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// ProjectGraphOpts controls ProjectGraph (mode=project on GET /v1/graph).
type ProjectGraphOpts struct {
	MaxNodes int
	// Scope filters to one thin scope (slug or id) plus bounded N-hop neighbors.
	// Empty = full project (subject to max_nodes).
	Scope string
	// Depth for scope-filtered expansion (default 1; allow 2; hard-capped by MaxNodes).
	Depth int
}

var projectGraphKindOrder = map[string]int{
	"goal": 0, "task": 1, "decision": 2, "assumption": 3, "discovery": 4,
	"plan_change": 5, "claim": 6, "evidence": 7, "review": 8, "capability": 9,
	"change": 10, "regression": 11, "scope": 12,
}

// ProjectGraph returns a bounded view of project entities and edges between them.
// max_nodes is required (1..5000). Truncated=true when total entities exceed the budget.
// When Scope is set, returns scope members + N-hop neighbors only (still bounded by max_nodes).
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

	if strings.TrimSpace(opts.Scope) != "" {
		return e.projectGraphScoped(opts)
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

func (e *Engine) projectGraphScoped(opts ProjectGraphOpts) (*BoundedGraph, error) {
	sc, err := e.resolveScope(opts.Scope)
	if err != nil {
		return nil, err
	}
	depth := opts.Depth
	if depth <= 0 {
		depth = 1
	}
	if depth > 2 {
		depth = 2
	}

	memberLinks, err := e.store.ListLinksTo("scope", sc.ID)
	if err != nil {
		return nil, err
	}

	// Seed: scope node + members.
	type frontierItem struct {
		h Hit
	}
	seen := map[string]Hit{}
	seedHit := Hit{EntityType: "scope", EntityID: sc.ID, Title: sc.Title, Distance: 0}
	if seedHit.Title == "" {
		seedHit.Title = sc.Slug
	}
	seen[hitKey("scope", sc.ID)] = seedHit
	frontier := []frontierItem{{h: seedHit}}

	for _, l := range memberLinks {
		if l.Rel != "scope_member" {
			continue
		}
		nh, lerr := e.lookupEntity(l.FromType, l.FromID, "scope_member", 0, 1.0)
		if lerr != nil {
			if isNotFound(lerr) {
				continue
			}
			return nil, lerr
		}
		nh.Distance = 0
		k := hitKey(nh.EntityType, nh.EntityID)
		if _, ok := seen[k]; ok {
			continue
		}
		if len(seen) >= opts.MaxNodes {
			break
		}
		seen[k] = nh
		frontier = append(frontier, frontierItem{h: nh})
	}

	truncated := len(seen) >= opts.MaxNodes
	edgeSeen := map[string]struct{}{}
	var edges []GraphEdge
	addEdge := func(ed GraphEdge) {
		k := ed.Rel + "\x00" + ed.From + "\x00" + ed.To
		if _, ok := edgeSeen[k]; ok {
			return
		}
		edgeSeen[k] = struct{}{}
		edges = append(edges, ed)
	}

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

	// Also collect edges among included nodes that may not have been walked yet.
	scopeIDs := make(map[string]string) // entityID → scope_id from membership
	for _, l := range memberLinks {
		if l.Rel == "scope_member" {
			scopeIDs[l.FromID] = sc.ID
		}
	}

	nodes := make([]GraphNode, 0, len(seen))
	included := make(map[string]struct{}, len(seen))
	for _, h := range seen {
		included[h.EntityID] = struct{}{}
		gn := GraphNode{ID: h.EntityID, Kind: h.EntityType, Title: h.Title}
		if h.EntityType == "task" {
			t, err := e.store.GetTask(h.EntityID)
			if err == nil && t.GoalID != nil && *t.GoalID != "" {
				gn.GoalID = *t.GoalID
			}
		}
		if sid, ok := scopeIDs[h.EntityID]; ok {
			gn.ScopeID = sid
		}
		nodes = append(nodes, gn)
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

	// Re-collect edges among included set for completeness (both endpoints in set).
	more, err := e.collectEdgesForNodes(nodes, included)
	if err != nil {
		return nil, err
	}
	for _, ed := range more {
		addEdge(ed)
	}

	center := sc.ID
	return &BoundedGraph{
		Mode:          "project",
		Center:        center,
		MaxNodes:      opts.MaxNodes,
		TotalEntities: len(nodes),
		Nodes:         nodes,
		Edges:         edges,
		Truncated:     truncated,
	}, nil
}

func (e *Engine) resolveScope(slugOrID string) (store.Scope, error) {
	slugOrID = strings.TrimSpace(slugOrID)
	if slugOrID == "" {
		return store.Scope{}, fmt.Errorf("retrieval: ProjectGraph: scope is required")
	}
	if sc, err := e.store.GetScope(slugOrID); err == nil {
		return sc, nil
	} else if !isNotFound(err) {
		return store.Scope{}, err
	}
	sc, err := e.store.GetScopeBySlug(slugOrID)
	if err != nil {
		if isNotFound(err) {
			return store.Scope{}, fmt.Errorf("retrieval: ProjectGraph: scope %q not found", slugOrID)
		}
		return store.Scope{}, err
	}
	return sc, nil
}

func (e *Engine) collectProjectNodes(maxNodes int) ([]GraphNode, int, error) {
	total := 0
	for _, table := range []string{
		"goals", "tasks", "decisions", "assumptions", "discoveries", "plan_changes",
		"claims", "evidence", "reviews", "capabilities", "changes", "regressions", "scopes",
	} {
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

	if rem := remaining(); rem > 0 {
		goals, err := e.store.ListGoalsLimited(rem)
		if err != nil {
			return nil, 0, err
		}
		for _, g := range goals {
			if !appendNode(GraphNode{ID: g.ID, Kind: "goal", Title: g.Title}) {
				return nodes, total, nil
			}
		}
	} else {
		return nodes, total, nil
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
	} else {
		return nodes, total, nil
	}

	if rem := remaining(); rem > 0 {
		decisions, err := e.store.ListDecisionsLimited(rem)
		if err != nil {
			return nil, 0, err
		}
		for _, d := range decisions {
			if !appendNode(GraphNode{ID: d.ID, Kind: "decision", Title: d.Title}) {
				return nodes, total, nil
			}
		}
	} else {
		return nodes, total, nil
	}

	if rem := remaining(); rem > 0 {
		assumptions, err := e.store.ListAssumptionsLimited(rem)
		if err != nil {
			return nil, 0, err
		}
		for _, a := range assumptions {
			if !appendNode(GraphNode{ID: a.ID, Kind: "assumption", Title: a.Title}) {
				return nodes, total, nil
			}
		}
	} else {
		return nodes, total, nil
	}

	if rem := remaining(); rem > 0 {
		discoveries, err := e.store.ListDiscoveriesLimited(rem)
		if err != nil {
			return nil, 0, err
		}
		for _, d := range discoveries {
			if !appendNode(GraphNode{ID: d.ID, Kind: "discovery", Title: d.Title}) {
				return nodes, total, nil
			}
		}
	} else {
		return nodes, total, nil
	}

	if rem := remaining(); rem > 0 {
		planChanges, err := e.store.ListPlanChangesLimited(rem)
		if err != nil {
			return nil, 0, err
		}
		for _, p := range planChanges {
			if !appendNode(GraphNode{ID: p.ID, Kind: "plan_change", Title: p.Title}) {
				return nodes, total, nil
			}
		}
	} else {
		return nodes, total, nil
	}

	if rem := remaining(); rem > 0 {
		claims, err := e.store.ListClaimsLimited(rem)
		if err != nil {
			return nil, 0, err
		}
		for _, c := range claims {
			if !appendNode(GraphNode{ID: c.ID, Kind: "claim", Title: c.Title}) {
				return nodes, total, nil
			}
		}
	} else {
		return nodes, total, nil
	}

	if rem := remaining(); rem > 0 {
		evidence, err := e.store.ListEvidenceLimited(rem)
		if err != nil {
			return nil, 0, err
		}
		for _, ev := range evidence {
			if !appendNode(GraphNode{ID: ev.ID, Kind: "evidence", Title: ev.Title}) {
				return nodes, total, nil
			}
		}
	} else {
		return nodes, total, nil
	}

	if rem := remaining(); rem > 0 {
		reviews, err := e.store.ListReviewsLimited(rem)
		if err != nil {
			return nil, 0, err
		}
		for _, r := range reviews {
			if !appendNode(GraphNode{ID: r.ID, Kind: "review", Title: r.Title}) {
				return nodes, total, nil
			}
		}
	} else {
		return nodes, total, nil
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
	} else {
		return nodes, total, nil
	}

	if rem := remaining(); rem > 0 {
		changes, err := e.store.ListChangesLimited(rem)
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
	} else {
		return nodes, total, nil
	}

	if rem := remaining(); rem > 0 {
		regressions, err := e.store.ListRegressionsLimited(rem)
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
	} else {
		return nodes, total, nil
	}

	if rem := remaining(); rem > 0 {
		allScopes, err := e.store.ListScopes()
		if err != nil {
			return nil, 0, err
		}
		scopes := allScopes
		if len(scopes) > rem {
			scopes = scopes[:rem]
		}
		for _, sc := range scopes {
			title := sc.Title
			if title == "" {
				title = sc.Slug
			}
			if !appendNode(GraphNode{ID: sc.ID, Kind: "scope", Title: title}) {
				return nodes, total, nil
			}
		}
	} else {
		return nodes, total, nil
	}

	// Populate scope_id on members when a single membership is cheap to map.
	memberOf, err := e.memberScopeIndex()
	if err != nil {
		return nil, 0, err
	}
	for i := range nodes {
		if sid, ok := memberOf[nodes[i].ID]; ok {
			nodes[i].ScopeID = sid
		}
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

	return nodes, total, nil
}

// memberScopeIndex maps entity id → scope id for scope_member links.
// If an entity belongs to multiple scopes, the first (stable by link order) wins.
func (e *Engine) memberScopeIndex() (map[string]string, error) {
	links, err := e.store.ListLinksByRel("scope_member")
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(links))
	for _, l := range links {
		if l.ToType != "scope" {
			continue
		}
		if _, ok := out[l.FromID]; ok {
			continue
		}
		out[l.FromID] = l.ToID
	}
	return out, nil
}

func (e *Engine) collectEdgesForNodes(nodes []GraphNode, included map[string]struct{}) ([]GraphEdge, error) {
	edgeSeen := map[string]struct{}{}
	var edges []GraphEdge
	for _, n := range nodes {
		fromType := domainEntityType(n.Kind)
		links, err := e.store.ListLinksFrom(fromType, n.ID)
		if err != nil {
			return nil, err
		}
		for _, l := range links {
			if _, ok := included[l.ToID]; !ok {
				continue
			}
			k := l.Rel + "\x00" + l.FromID + "\x00" + l.ToID
			if _, ok := edgeSeen[k]; ok {
				continue
			}
			edgeSeen[k] = struct{}{}
			edges = append(edges, GraphEdge{
				From:       l.FromID,
				To:         l.ToID,
				Rel:        l.Rel,
				Provenance: ProvenanceFromSourceType(l.SourceType),
			})
		}
		// Add goal_has_task edges from task.GoalID (not in entity_links table)
		if n.Kind == "task" && n.GoalID != "" {
			if _, ok := included[n.GoalID]; ok {
				k := "goal_has_task\x00" + n.GoalID + "\x00" + n.ID
				if _, ok := edgeSeen[k]; !ok {
					edgeSeen[k] = struct{}{}
					edges = append(edges, GraphEdge{From: n.GoalID, To: n.ID, Rel: "goal_has_task", Provenance: "explicit"})
				}
			}
		}
	}
	return edges, nil
}

// domainEntityType maps graph node kind to store entity type string.
func domainEntityType(kind string) string {
	switch kind {
	case "goal":
		return "goal"
	case "task":
		return "task"
	case "decision":
		return "decision"
	case "assumption":
		return "assumption"
	case "discovery":
		return "discovery"
	case "plan_change":
		return "plan_change"
	case "claim":
		return "claim"
	case "evidence":
		return "evidence"
	case "review":
		return "review"
	case "capability":
		return "capability"
	case "change":
		return "change"
	case "regression":
		return "regression"
	case "scope":
		return "scope"
	default:
		return "entity"
	}
}
