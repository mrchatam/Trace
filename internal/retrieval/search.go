package retrieval

import (
	"context"
	"sort"
	"strings"
)

// Search wraps store FTS5 lexical search and attaches reason_code fts_match.
// opts.Limit defaults to 32 and is hard-capped at 64 (pre-compile candidate cap).
// When opts.Intent is set, the FTS query is built from ExtractIntent (G9 pre-channel);
// otherwise the raw q string is used (backward compatible MCP path).
func (e *Engine) Search(ctx context.Context, q string, opts SearchOptions) ([]Hit, error) {
	_ = ctx
	limit := opts.Limit
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}
	ftsQ := strings.TrimSpace(q)
	var terms []string
	if opts.Intent != nil {
		intent := ExtractIntent(*opts.Intent)
		if ftsQ == "" {
			terms = intentSearchTerms(intent)
		} else {
			terms = []string{ftsQ}
			qLower := strings.ToLower(ftsQ)
			for _, term := range intentSearchTerms(intent) {
				if term == "" || strings.Contains(qLower, strings.ToLower(term)) {
					continue
				}
				terms = append(terms, term)
			}
		}
	} else if ftsQ != "" {
		terms = []string{ftsQ}
	}
	if len(terms) == 0 {
		return nil, nil
	}
	return e.searchFTSMultiOr(terms, limit)
}

func intentSearchTerms(intent Intent) []string {
	seen := map[string]struct{}{}
	var terms []string
	add := func(t string) {
		t = strings.TrimSpace(t)
		if t == "" {
			return
		}
		key := strings.ToLower(t)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		terms = append(terms, t)
	}
	for _, kw := range intent.Keywords {
		add(kw)
	}
	for _, h := range intent.EntityHints {
		add(h.Value)
	}
	return terms
}

func (e *Engine) searchFTSMultiOr(terms []string, limit int) ([]Hit, error) {
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
				ReasonCode: ReasonFTSMatch,
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
