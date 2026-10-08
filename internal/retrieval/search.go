package retrieval

import (
	"context"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// Search wraps store FTS5 lexical search and attaches reason_code fts_match.
// opts.Limit defaults to 32 and is hard-capped at 64 (pre-compile candidate cap).
// When opts.Intent is set, the FTS query is built from ExtractIntent (G9 pre-channel);
// otherwise the raw q string is used (backward compatible MCP path).
func (e *Engine) Search(ctx context.Context, q string, opts SearchOptions) ([]Hit, error) {
	hits, _, err := e.SearchWithQuery(ctx, q, opts)
	return hits, err
}

// SearchWithQuery is Search plus the compiled FTS query, so callers can report
// an operator rewrite or an ops-mode fallback. The compiled form is only
// meaningful for a single-term (non-Intent) search; the intent fan-out path
// returns a zero FTSQuery because no single MATCH expression was issued.
func (e *Engine) SearchWithQuery(ctx context.Context, q string, opts SearchOptions) ([]Hit, store.FTSQuery, error) {
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
		return nil, store.FTSQuery{}, nil
	}

	if len(terms) == 1 {
		fts, compiled, err := e.store.SearchFTSOpts(store.SearchFTSOpts{
			Query: terms[0], Limit: limit, Mode: opts.Matcher,
		})
		if err != nil {
			return nil, compiled, err
		}
		out := e.hitsFromFTS(fts, limit)
		sortHits(out)
		return out, compiled, nil
	}

	// Intent fan-out: one MATCH per term, unioned by the caller. Boolean
	// operator mode does not apply across separately issued queries.
	hits, err := e.searchFTSMultiOr(terms, limit)
	return hits, store.FTSQuery{}, err
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
		seen = appendUniqueHits(seen, &out, fts, limit)
	}
	sortHits(out)
	return out, nil
}

// hitsFromFTS converts FTS rows to engine hits with a fresh dedupe set.
func (e *Engine) hitsFromFTS(fts []store.FTSHit, limit int) []Hit {
	out := make([]Hit, 0, len(fts))
	appendUniqueHits(map[string]struct{}{}, &out, fts, limit)
	return out
}

// appendUniqueHits appends FTS rows to out, skipping entity ids already seen
// and stopping at limit. Local seen is threaded by the caller so a multi-term
// fan-out dedupes across terms.
func appendUniqueHits(seen map[string]struct{}, out *[]Hit, fts []store.FTSHit, limit int) map[string]struct{} {
	for _, fh := range fts {
		if len(*out) >= limit {
			break
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
		*out = append(*out, Hit{
			EntityType: fh.EntityType,
			EntityID:   fh.EntityID,
			Title:      title,
			Excerpt:    excerpt(fh.Body),
			Path:       fh.Path,
			ReasonCode: ReasonFTSMatch,
			Score:      fh.Score,
			Distance:   0,
		})
	}
	return seen
}

// sortHits orders by descending score with entity id tiebreaks so paging is
// deterministic.
func sortHits(out []Hit) {
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		if out[i].EntityType != out[j].EntityType {
			return out[i].EntityType < out[j].EntityType
		}
		return out[i].EntityID < out[j].EntityID
	})
}
