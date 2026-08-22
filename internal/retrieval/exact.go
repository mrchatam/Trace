package retrieval

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// Exact resolves entities by UUID+type, file path, and/or symbol name.
// Misses return an empty slice (and nil error) unless a hard store failure occurs.
func (e *Engine) Exact(ctx context.Context, q ExactQuery) ([]Hit, error) {
	_ = ctx
	var out []Hit
	seen := map[string]struct{}{}

	add := func(h Hit) {
		k := hitKey(h.EntityType, h.EntityID)
		if _, ok := seen[k]; ok {
			return
		}
		seen[k] = struct{}{}
		out = append(out, h)
	}

	if q.EntityID != "" {
		typ := NormalizeEntityType(q.EntityType)
		if typ == "" {
			return nil, fmt.Errorf("retrieval: Exact: EntityType required when EntityID set")
		}
		h, err := e.lookupEntity(typ, q.EntityID, ReasonExactID, 0, 1.0)
		if err != nil {
			if isNotFound(err) {
				return out, nil
			}
			return nil, err
		}
		add(h)
	}

	if q.Path != "" {
		path := store.NormalizePath(q.Path)
		f, err := e.store.GetFileByPath(path)
		if err != nil {
			if !isNotFound(err) {
				return nil, err
			}
		} else {
			add(Hit{
				EntityType: "file",
				EntityID:   f.ID,
				Title:      f.Path,
				Path:       f.Path,
				ReasonCode: ReasonExactPath,
				Score:      1.0,
				Distance:   0,
			})
			if q.SymbolName != "" {
				syms, err := e.store.ListSymbolsByPath(path)
				if err != nil {
					return nil, err
				}
				for _, sym := range syms {
					if sym.Name == q.SymbolName {
						add(Hit{
							EntityType: "symbol",
							EntityID:   sym.ID,
							Title:      sym.Name,
							Excerpt:    sym.Kind,
							Path:       path,
							ReasonCode: ReasonExactSymbol,
							Score:      1.0,
							Distance:   0,
						})
					}
				}
			}
		}
	} else if q.SymbolName != "" {
		// Unambiguous single-name lookup via FTS then exact name filter.
		fts, err := e.store.SearchFTS(q.SymbolName, 16)
		if err != nil {
			return nil, err
		}
		var matches []Hit
		for _, fh := range fts {
			if fh.EntityType != "symbol" || fh.SymbolName != q.SymbolName {
				continue
			}
			matches = append(matches, Hit{
				EntityType: "symbol",
				EntityID:   fh.EntityID,
				Title:      fh.SymbolName,
				Excerpt:    fh.SymbolKind,
				Path:       fh.Path,
				ReasonCode: ReasonExactSymbol,
				Score:      1.0,
				Distance:   0,
			})
		}
		if len(matches) == 1 {
			add(matches[0])
		}
	}

	return out, nil
}

func (e *Engine) lookupEntity(entityType, id, reason string, distance int, score float64) (Hit, error) {
	entityType = NormalizeEntityType(entityType)
	switch entityType {
	case "goal":
		g, err := e.store.GetGoal(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{EntityType: "goal", EntityID: g.ID, Title: g.Title, Excerpt: excerpt(g.Body), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "task":
		t, err := e.store.GetTask(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{EntityType: "task", EntityID: t.ID, Title: t.Title, Excerpt: excerpt(t.Body), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "decision":
		d, err := e.store.GetDecision(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{EntityType: "decision", EntityID: d.ID, Title: d.Title, Excerpt: excerpt(d.Body), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "assumption":
		a, err := e.store.GetAssumption(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{EntityType: "assumption", EntityID: a.ID, Title: a.Title, Excerpt: excerpt(a.Body), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "discovery":
		d, err := e.store.GetDiscovery(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{EntityType: "discovery", EntityID: d.ID, Title: d.Title, Excerpt: excerpt(d.Body), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "plan_change":
		p, err := e.store.GetPlanChange(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{EntityType: "plan_change", EntityID: p.ID, Title: p.Title, Excerpt: excerpt(p.Body), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "claim":
		c, err := e.store.GetClaim(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{EntityType: "claim", EntityID: c.ID, Title: c.Title, Excerpt: excerpt(c.Body), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "evidence":
		ev, err := e.store.GetEvidence(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{EntityType: "evidence", EntityID: ev.ID, Title: ev.Title, Excerpt: excerpt(ev.Body), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "review":
		r, err := e.store.GetReview(id)
		if err != nil {
			return Hit{}, err
		}
		ex := r.Result
		if ex == "" {
			ex = excerpt(r.Body)
		}
		return Hit{EntityType: "review", EntityID: r.ID, Title: r.Title, Excerpt: ex, ReasonCode: reason, Score: score, Distance: distance}, nil
	case "capability":
		c, err := e.store.GetCapability(id)
		if err != nil {
			return Hit{}, err
		}
		title := c.Title
		if title == "" {
			title = c.Slug
		}
		ex := c.Status
		if ex == "" {
			ex = excerpt(c.Body)
		}
		return Hit{EntityType: "capability", EntityID: c.ID, Title: title, Excerpt: ex, ReasonCode: reason, Score: score, Distance: distance}, nil
	case "file":
		f, err := e.store.GetFileByID(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{EntityType: "file", EntityID: f.ID, Title: f.Path, Path: f.Path, ReasonCode: reason, Score: score, Distance: distance}, nil
	case "symbol":
		sym, path, err := e.store.GetSymbolByID(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{
			EntityType: "symbol",
			EntityID:   sym.ID,
			Title:      sym.Name,
			Excerpt:    sym.Kind,
			Path:       path,
			ReasonCode: reason,
			Score:      score,
			Distance:   distance,
		}, nil
	case "uncertainty":
		u, err := e.store.GetUncertainty(id)
		if err != nil {
			return Hit{}, err
		}
		ex := excerpt(u.Body)
		if prefix := uncertaintyExcerptPrefix(u.Severity, u.Kind); prefix != "" {
			if ex != "" {
				ex = prefix + " " + ex
			} else {
				ex = prefix
			}
		}
		return Hit{EntityType: "uncertainty", EntityID: u.ID, Title: u.Title, Excerpt: ex, ReasonCode: reason, Score: score, Distance: distance}, nil
	case "hypothesis":
		h, err := e.store.GetHypothesis(id)
		if err != nil {
			return Hit{}, err
		}
		ex := excerpt(h.Body)
		if h.Status != "" {
			if ex != "" {
				ex = h.Status + " " + ex
			} else {
				ex = h.Status
			}
		}
		return Hit{EntityType: "hypothesis", EntityID: h.ID, Title: h.Title, Excerpt: ex, ReasonCode: reason, Score: score, Distance: distance}, nil
	case "change":
		c, err := e.store.GetChange(id)
		if err != nil {
			return Hit{}, err
		}
		title := c.GitCommit
		if title == "" && len(c.ID) >= 8 {
			title = c.ID[:8]
		}
		ex := strings.TrimSpace(c.Reason)
		if c.Status != "" {
			if ex != "" {
				ex = ex + " " + c.Status
			} else {
				ex = c.Status
			}
		}
		return Hit{EntityType: "change", EntityID: c.ID, Title: title, Excerpt: excerpt(ex), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "effect":
		ef, err := e.store.GetEffect(id)
		if err != nil {
			return Hit{}, err
		}
		ex := ef.Comparison
		if ef.Expected != "" {
			if ex != "" {
				ex = ex + " expected=" + excerpt(ef.Expected)
			} else {
				ex = "expected=" + excerpt(ef.Expected)
			}
		}
		if ef.Actual != "" {
			if ex != "" {
				ex = ex + " actual=" + excerpt(ef.Actual)
			} else {
				ex = "actual=" + excerpt(ef.Actual)
			}
		}
		return Hit{EntityType: "effect", EntityID: ef.ID, Title: ef.Dimension, Excerpt: excerpt(ex), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "regression":
		r, err := e.store.GetRegression(id)
		if err != nil {
			return Hit{}, err
		}
		ex := strings.TrimSpace(r.Summary)
		if r.Attribution != "" {
			if ex != "" {
				ex = ex + " " + r.Attribution
			} else {
				ex = r.Attribution
			}
		}
		return Hit{EntityType: "regression", EntityID: r.ID, Title: r.Dimension, Excerpt: excerpt(ex), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "reflection":
		r, err := e.store.GetReflection(id)
		if err != nil {
			return Hit{}, err
		}
		return Hit{EntityType: "reflection", EntityID: r.ID, Title: r.Summary, Excerpt: excerpt(r.InvalidatedAssumptionsJSON), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "baseline":
		b, err := e.store.GetBaseline(id)
		if err != nil {
			return Hit{}, err
		}
		title := b.Label
		if title == "" {
			title = b.GitCommit
		}
		return Hit{EntityType: "baseline", EntityID: b.ID, Title: title, Excerpt: excerpt(b.ScoresJSON), ReasonCode: reason, Score: score, Distance: distance}, nil
	case "outcome_result":
		o, err := e.store.GetOutcomeResult(id)
		if err != nil {
			return Hit{}, err
		}
		title := o.TestName
		if title == "" {
			title = o.Kind
		}
		return Hit{EntityType: "outcome_result", EntityID: o.ID, Title: title, Excerpt: excerpt(o.Summary), ReasonCode: reason, Score: score, Distance: distance}, nil
	default:
		return Hit{}, fmt.Errorf("retrieval: unknown entity type %q", entityType)
	}
}

func uncertaintyExcerptPrefix(severity, kind string) string {
	var parts []string
	if severity != "" {
		parts = append(parts, severity)
	}
	if kind != "" {
		parts = append(parts, kind)
	}
	return strings.Join(parts, " ")
}

func excerpt(body string) string {
	body = strings.TrimSpace(body)
	if len(body) <= 240 {
		return body
	}
	return body[:240] + "…"
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "sql: no rows") || strings.Contains(msg, ": "+sql.ErrNoRows.Error())
}
