package retrieval

import (
	"context"
)

// HistoricalFileHits returns L3 temporal hits (historical_vcs) for indexed file paths.
// Requires optional VCS; returns nil when VCS is absent.
func (e *Engine) HistoricalFileHits(ctx context.Context, paths []string) ([]Hit, error) {
	if e.vcs == nil {
		return nil, nil
	}
	var out []Hit
	seen := map[string]struct{}{}
	for _, path := range paths {
		if path == "" {
			continue
		}
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		meta, err := e.vcs.LastChanged(ctx, path)
		if err != nil || meta.OID == "" {
			continue
		}
		out = append(out, Hit{
			EntityType: "commit",
			EntityID:   meta.OID,
			Title:      meta.Subject,
			ReasonCode: ReasonHistoricalVCS,
			Distance:   3,
			Score:      0.3,
		})
	}
	return out, nil
}

// BlastHitsToLayer converts ImpactWalk blast nodes to retrieval hits for layer admission.
// minHop filters blast entries (1 for L2 dependents, 2 for L3 cross-module).
func BlastHitsToLayer(res *ImpactWalkResult, minHop int) []Hit {
	if res == nil {
		return nil
	}
	out := make([]Hit, 0, len(res.Blast))
	for _, b := range res.Blast {
		if b.Hop < minHop {
			continue
		}
		dist := b.Hop + 1
		out = append(out, Hit{
			EntityType:     b.EntityType,
			EntityID:       b.EntityID,
			Title:          b.Title,
			Path:           b.Path,
			ReasonCode:     ReasonGraphNeighbor,
			Distance:       dist,
			EdgeProvenance: b.EdgeProvenance,
			Score:          0.4,
		})
	}
	return out
}
