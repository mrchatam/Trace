package retrieval

import (
	"context"
	"sort"
)

// conceptEntityTypes — bounded set for G6 graph-label channel (no file/symbol/task).
var conceptEntityTypes = map[string]struct{}{
	"discovery":  {},
	"assumption": {},
	"decision":   {},
	"goal":       {},
	"claim":      {},
}

func isConceptEntityType(entityType string) bool {
	_, ok := conceptEntityTypes[entityType]
	return ok
}

// SearchGraphLabels runs bounded FTS over concept entity types only (G6 graph-label channel).
// Uses G9 intent terms; hits carry reason_code graph_label_match (distinct from fts_match).
func (e *Engine) SearchGraphLabels(ctx context.Context, intent Intent, opts SearchOptions) ([]Hit, error) {
	_ = ctx
	limit := opts.Limit
	if limit <= 0 {
		limit = 16
	}
	if limit > 64 {
		limit = 64
	}
	terms := intentSearchTerms(intent)
	if len(terms) == 0 {
		return nil, nil
	}
	return e.searchGraphLabelsMultiOr(terms, limit)
}

func (e *Engine) searchGraphLabelsMultiOr(terms []string, limit int) ([]Hit, error) {
	seen := map[string]struct{}{}
	var out []Hit
	for _, term := range terms {
		if len(out) >= limit {
			break
		}
		fts, err := e.store.SearchFTS(term, limit)
		if err != nil {
			return nil, err
		}
		for _, fh := range fts {
			if !isConceptEntityType(fh.EntityType) {
				continue
			}
			key := hitKey(fh.EntityType, fh.EntityID)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			title := fh.Title
			if title == "" {
				title = fh.SymbolName
			}
			if title == "" {
				title = fh.Path
			}
			out = append(out, Hit{
				EntityType: fh.EntityType,
				EntityID:   fh.EntityID,
				Title:      title,
				Excerpt:    excerpt(fh.Body),
				Path:       fh.Path,
				ReasonCode: ReasonGraphLabelMatch,
				Score:      fh.Score,
				Distance:   0,
			})
			if len(out) >= limit {
				break
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
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

// MergeConceptHits appends concept hits and dedupes by entity key.
// When the same entity appears in base and concept hits, graph_label_match wins.
func MergeConceptHits(base, concept []Hit) []Hit {
	if len(concept) == 0 {
		return base
	}
	index := make(map[string]int, len(base))
	out := append([]Hit(nil), base...)
	for i, h := range base {
		index[hitKey(h.EntityType, h.EntityID)] = i
	}
	for _, ch := range concept {
		key := hitKey(ch.EntityType, ch.EntityID)
		if i, ok := index[key]; ok {
			if ch.ReasonCode == ReasonGraphLabelMatch {
				out[i] = ch
			}
			continue
		}
		index[key] = len(out)
		out = append(out, ch)
	}
	return out
}
