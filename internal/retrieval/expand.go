package retrieval

import (
	"context"
	"fmt"
	"sort"
)

const (
	maxExpandDepth     = 2
	defaultExpandDepth = 1
)

// Expand walks causal links + goal_id and structural file↔symbol/import neighbors
// up to depth (1..2). Seed hits are distance 0; each hop increments distance.
// Depth > 2 or < 1 is rejected.
func (e *Engine) Expand(ctx context.Context, seeds []Hit, depth int) ([]Hit, error) {
	_ = ctx
	if depth < 1 || depth > maxExpandDepth {
		return nil, fmt.Errorf("retrieval: Expand: depth must be 1..%d, got %d", maxExpandDepth, depth)
	}

	seen := map[string]Hit{}
	type frontierItem struct {
		h Hit
	}
	var frontier []frontierItem

	for _, s := range seeds {
		s.Distance = 0
		if s.ReasonCode == "" {
			s.ReasonCode = ReasonExactID
		}
		k := hitKey(s.EntityType, s.EntityID)
		seen[k] = s
		frontier = append(frontier, frontierItem{h: s})
	}

	for d := 1; d <= depth; d++ {
		var next []frontierItem
		for _, fi := range frontier {
			neighbors, err := e.neighbors(fi.h)
			if err != nil {
				return nil, err
			}
			for _, n := range neighbors {
				n.Distance = d
				k := hitKey(n.EntityType, n.EntityID)
				if existing, ok := seen[k]; ok {
					if n.Distance < existing.Distance {
						seen[k] = n
					}
					continue
				}
				seen[k] = n
				next = append(next, frontierItem{h: n})
			}
		}
		frontier = next
	}

	// DF-19: goal-scoped discovery↔plan_change attach for task seeds (not global dump).
	for _, s := range seeds {
		if s.EntityType != "task" {
			continue
		}
		dpc, err := e.discoveryPlanChangeHitsForTask(s.EntityID, seen)
		if err != nil {
			return nil, err
		}
		for _, h := range dpc {
			k := hitKey(h.EntityType, h.EntityID)
			if existing, ok := seen[k]; ok {
				if h.Distance < existing.Distance {
					seen[k] = h
				}
				continue
			}
			seen[k] = h
		}
	}

	out := make([]Hit, 0, len(seen))
	for _, h := range seen {
		out = append(out, h)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Distance != out[j].Distance {
			return out[i].Distance < out[j].Distance
		}
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		if out[i].EntityType != out[j].EntityType {
			return out[i].EntityType < out[j].EntityType
		}
		return out[i].EntityID < out[j].EntityID
	})
	return out, nil
}

func (e *Engine) neighbors(h Hit) ([]Hit, error) {
	var out []Hit

	// Causal: links from/to
	fromLinks, err := e.store.ListLinksFrom(h.EntityType, h.EntityID)
	if err != nil {
		return nil, err
	}
	for _, l := range fromLinks {
		nh, err := e.hitFromLinkNeighbor(l.ToType, l.ToID, l.Rel, 0.6)
		if err != nil {
			if isNotFound(err) {
				continue
			}
			return nil, err
		}
		out = append(out, nh)
	}
	toLinks, err := e.store.ListLinksTo(h.EntityType, h.EntityID)
	if err != nil {
		return nil, err
	}
	for _, l := range toLinks {
		nh, err := e.hitFromLinkNeighbor(l.FromType, l.FromID, l.Rel, 0.6)
		if err != nil {
			if isNotFound(err) {
				continue
			}
			return nil, err
		}
		out = append(out, nh)
	}

	// Goal↔Task via goal_id
	switch h.EntityType {
	case "task":
		t, err := e.store.GetTask(h.EntityID)
		if err != nil {
			return nil, err
		}
		if t.GoalID != nil && *t.GoalID != "" {
			gh, err := e.lookupEntity("goal", *t.GoalID, ReasonGoalHasTask, 0, 0.9)
			if err != nil {
				if !isNotFound(err) {
					return nil, err
				}
			} else {
				out = append(out, gh)
			}
		}
	case "goal":
		tasks, err := e.store.ListTasksByGoalID(h.EntityID)
		if err != nil {
			return nil, err
		}
		for _, t := range tasks {
			// DF-35: title/id OK; never attach sibling task body via goal expand.
			out = append(out, Hit{
				EntityType: "task",
				EntityID:   t.ID,
				Title:      t.Title,
				Excerpt:    "",
				ReasonCode: ReasonGoalHasTask,
				Score:      0.9,
			})
		}
	case "file":
		path := h.Path
		if path == "" {
			f, err := e.store.GetFileByID(h.EntityID)
			if err != nil {
				return nil, err
			}
			path = f.Path
		}
		syms, err := e.store.ListSymbolsByPath(path)
		if err != nil {
			return nil, err
		}
		for _, sym := range syms {
			out = append(out, Hit{
				EntityType: "symbol",
				EntityID:   sym.ID,
				Title:      sym.Name,
				Excerpt:    sym.Kind,
				Path:       path,
				ReasonCode: ReasonGraphNeighbor,
				Score:      0.5,
			})
		}
		imps, err := e.store.ListImportsByPath(path)
		if err != nil {
			return nil, err
		}
		for _, imp := range imps {
			nf, err := e.resolveImportedFile(path, imp.ImportedPath)
			if err != nil {
				if isNotFound(err) {
					continue
				}
				return nil, err
			}
			out = append(out, Hit{
				EntityType:     "file",
				EntityID:       nf.ID,
				Title:          nf.Path,
				Path:           nf.Path,
				ReasonCode:     ReasonGraphNeighbor,
				Score:          0.45,
				EdgeProvenance: imp.Provenance,
			})
		}
	case "symbol":
		// Prefer path on hit; else resolve via FTS/store not needed if Expand seeds carry path.
		if h.Path != "" {
			f, err := e.store.GetFileByPath(h.Path)
			if err != nil {
				if !isNotFound(err) {
					return nil, err
				}
			} else {
				out = append(out, Hit{
					EntityType: "file",
					EntityID:   f.ID,
					Title:      f.Path,
					Path:       f.Path,
					ReasonCode: ReasonGraphNeighbor,
					Score:      0.5,
				})
			}
		}
	}

	return out, nil
}

func (e *Engine) hitFromLinkNeighbor(entityType, entityID, rel string, score float64) (Hit, error) {
	reason := ReasonGraphNeighbor
	switch rel {
	case ReasonDecisionAffectsTask:
		reason = ReasonDecisionAffectsTask
	case ReasonDiscoveryCausesPlanChg:
		reason = ReasonDiscoveryCausesPlanChg
	case ReasonClaimHasEvidence:
		reason = ReasonClaimHasEvidence
	case ReasonReviewJudgesTask:
		reason = ReasonReviewJudgesTask
	case ReasonReviewJudgesScope:
		reason = ReasonReviewJudgesScope
	}
	return e.lookupEntity(entityType, entityID, reason, 0, score)
}
