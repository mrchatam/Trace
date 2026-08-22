package domain

import (
	"context"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

const (
	changeKindUnknown = "seg:unknown"
	changeKindPrefix  = "seg:"

	OutcomeKindRegression         = "regression"
	OutcomeKindEffectContradicted = "effect_contradicted"
	OutcomeKindImprovement        = "improvement"
	OutcomeKindEffectSupported    = "effect_supported"
	OutcomeKindTestFail           = "test_fail"
	OutcomeKindTestPass           = "test_pass"
	OutcomeKindNeutral            = "neutral"
)

// ChangeOutcomeSignals drives ClassifyChangeOutcome priority (one bucket per change).
type ChangeOutcomeSignals struct {
	HasRegression         bool
	HasEffectContradicted bool
	HasImprovement        bool
	HasEffectSupported    bool
	HasTestFail           bool
	HasTestPass           bool
}

// InferChangeKind returns seg:<first-path-segment> from the lexicographically smallest path.
// No paths → seg:unknown.
func InferChangeKind(paths []string) string {
	if len(paths) == 0 {
		return changeKindUnknown
	}
	sorted := append([]string(nil), paths...)
	sort.Strings(sorted)
	seg := firstPathSegment(sorted[0])
	if seg == "" {
		return changeKindUnknown
	}
	return changeKindPrefix + seg
}

func firstPathSegment(path string) string {
	path = store.NormalizePath(strings.TrimSpace(path))
	if path == "" {
		return ""
	}
	if i := strings.Index(path, "/"); i >= 0 {
		return path[:i]
	}
	return path
}

// ClassifyChangeOutcome returns one outcome bucket per change.
// Priority: regression > effect_contradicted > improvement > effect_supported > test_fail > test_pass > neutral.
func ClassifyChangeOutcome(s ChangeOutcomeSignals) string {
	switch {
	case s.HasRegression:
		return OutcomeKindRegression
	case s.HasEffectContradicted:
		return OutcomeKindEffectContradicted
	case s.HasImprovement:
		return OutcomeKindImprovement
	case s.HasEffectSupported:
		return OutcomeKindEffectSupported
	case s.HasTestFail:
		return OutcomeKindTestFail
	case s.HasTestPass:
		return OutcomeKindTestPass
	default:
		return OutcomeKindNeutral
	}
}

func isPositiveOutcome(outcome string) bool {
	switch outcome {
	case OutcomeKindImprovement, OutcomeKindEffectSupported, OutcomeKindTestPass:
		return true
	default:
		return false
	}
}

func isNegativeOutcome(outcome string) bool {
	switch outcome {
	case OutcomeKindRegression, OutcomeKindEffectContradicted, OutcomeKindTestFail:
		return true
	default:
		return false
	}
}

// SimilarChangesOpts filters similar historical changes by path prefix or change kind (mutually exclusive).
type SimilarChangesOpts struct {
	PathPrefix string
	ChangeKind string
	Limit      int
}

// SimilarEffectSummary is a compact effect row for similar-change queries.
type SimilarEffectSummary struct {
	Dimension  string `json:"dimension"`
	Comparison string `json:"comparison"`
}

// SimilarChangeRow is one prior change in QuerySimilarChanges results.
type SimilarChangeRow struct {
	ID          string                 `json:"id"`
	TaskID      string                 `json:"task_id"`
	ChangeKind  string                 `json:"change_kind"`
	Reason      string                 `json:"reason"`
	CreatedAt   string                 `json:"created_at"`
	Paths       []string               `json:"paths"`
	Effects     []SimilarEffectSummary `json:"effects"`
	OutcomeKind string                 `json:"outcome_kind"`
}

// SimilarChangesResult is the compact JSON payload for similar-change queries.
type SimilarChangesResult struct {
	Changes  []SimilarChangeRow    `json:"changes"`
	Patterns []store.ChangePattern `json:"patterns"`
}

// RefreshChangePatterns performs a deterministic full rebuild of change_patterns.
func (s *Service) RefreshChangePatterns(ctx context.Context) (int, error) {
	_ = ctx
	changes, err := s.store.ListActiveChanges()
	if err != nil {
		return 0, err
	}

	type bucketKey struct {
		changeKind  string
		outcomeKind string
	}
	agg := map[bucketKey]*store.ChangePattern{}

	for _, c := range changes {
		paths, err := s.store.ListChangePaths(c.ID)
		if err != nil {
			return 0, err
		}
		pathStrs := make([]string, 0, len(paths))
		for _, p := range paths {
			pathStrs = append(pathStrs, p.Path)
		}
		kind := InferChangeKind(pathStrs)
		signals, err := s.gatherChangeOutcomeSignals(c)
		if err != nil {
			return 0, err
		}
		outcome := ClassifyChangeOutcome(signals)
		key := bucketKey{changeKind: kind, outcomeKind: outcome}
		row, ok := agg[key]
		if !ok {
			row = &store.ChangePattern{ChangeKind: kind, OutcomeKind: outcome}
			agg[key] = row
		}
		if isPositiveOutcome(outcome) {
			row.CountPositive++
		}
		if isNegativeOutcome(outcome) {
			row.CountNegative++
		}
		if c.CreatedAt > row.LastSeen {
			row.LastSeen = c.CreatedAt
		}
	}

	rows := make([]store.ChangePattern, 0, len(agg))
	for _, row := range agg {
		rows = append(rows, *row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].ChangeKind != rows[j].ChangeKind {
			return rows[i].ChangeKind < rows[j].ChangeKind
		}
		return rows[i].OutcomeKind < rows[j].OutcomeKind
	})
	if err := s.store.ReplaceChangePatterns(rows); err != nil {
		return 0, err
	}
	return len(rows), nil
}

func (s *Service) gatherChangeOutcomeSignals(c store.Change) (ChangeOutcomeSignals, error) {
	var sig ChangeOutcomeSignals

	regs, err := s.store.ListRegressionsByChangeID(c.ID)
	if err != nil {
		return sig, err
	}
	sig.HasRegression = len(regs) > 0

	effects, err := s.store.ListEffectsByChangeID(c.ID)
	if err != nil {
		return sig, err
	}
	for _, e := range effects {
		switch e.Comparison {
		case store.EffectComparisonContradicted:
			sig.HasEffectContradicted = true
		case store.EffectComparisonSupported:
			sig.HasEffectSupported = true
		}
	}

	imps, err := s.store.ListImprovementsByChangeID(c.ID)
	if err != nil {
		return sig, err
	}
	sig.HasImprovement = len(imps) > 0

	tests, err := s.testOutcomesForChange(c)
	if err != nil {
		return sig, err
	}
	for _, t := range tests {
		switch t.TestStatus {
		case store.TestStatusFail, store.TestStatusError:
			sig.HasTestFail = true
		case store.TestStatusPass:
			sig.HasTestPass = true
		}
	}
	return sig, nil
}

func (s *Service) testOutcomesForChange(c store.Change) ([]store.OutcomeResult, error) {
	all, err := s.store.ListOutcomeResultsByTaskKind(c.TaskID, store.OutcomeKindTest)
	if err != nil {
		return nil, err
	}
	taskChanges, err := s.store.ListChangesByTaskID(c.TaskID)
	if err != nil {
		return nil, err
	}
	var nextCreatedAt string
	for _, tc := range taskChanges {
		if tc.Status == store.ChangeStatusSuperseded {
			continue
		}
		if tc.CreatedAt > c.CreatedAt && (nextCreatedAt == "" || tc.CreatedAt < nextCreatedAt) {
			nextCreatedAt = tc.CreatedAt
		}
	}
	var out []store.OutcomeResult
	for _, row := range all {
		if row.CreatedAt < c.CreatedAt {
			continue
		}
		if nextCreatedAt != "" && row.CreatedAt >= nextCreatedAt {
			continue
		}
		out = append(out, row)
	}
	return out, nil
}

// ListChangePatterns reads stored pattern rows (default limit 32, cap 64).
func (s *Service) ListChangePatterns(ctx context.Context, limit int) ([]store.ChangePattern, error) {
	_ = ctx
	return s.store.ListChangePatterns(limit)
}

// QuerySimilarChanges returns prior changes and related patterns for a path prefix or change kind.
func (s *Service) QuerySimilarChanges(ctx context.Context, opts SimilarChangesOpts) (SimilarChangesResult, error) {
	_ = ctx
	pathPrefix := strings.TrimSpace(opts.PathPrefix)
	changeKind := strings.TrimSpace(opts.ChangeKind)
	if pathPrefix != "" && changeKind != "" {
		return SimilarChangesResult{}, &ErrValidation{Msg: "PathPrefix and ChangeKind are mutually exclusive"}
	}
	if pathPrefix == "" && changeKind == "" {
		return SimilarChangesResult{}, &ErrValidation{Msg: "PathPrefix or ChangeKind is required"}
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}

	var (
		changes []store.Change
		err     error
	)
	if pathPrefix != "" {
		changes, err = s.store.ListChangesByPathPrefix(pathPrefix, limit)
	} else {
		changes, err = s.listChangesByKind(changeKind, limit)
	}
	if err != nil {
		return SimilarChangesResult{}, err
	}

	out := SimilarChangesResult{
		Changes:  []SimilarChangeRow{},
		Patterns: []store.ChangePattern{},
	}
	kindSet := map[string]struct{}{}
	for _, c := range changes {
		paths, err := s.store.ListChangePaths(c.ID)
		if err != nil {
			return SimilarChangesResult{}, err
		}
		pathStrs := make([]string, 0, len(paths))
		for _, p := range paths {
			pathStrs = append(pathStrs, p.Path)
		}
		kind := InferChangeKind(pathStrs)
		kindSet[kind] = struct{}{}

		effects, err := s.store.ListEffectsByChangeID(c.ID)
		if err != nil {
			return SimilarChangesResult{}, err
		}
		effectRows := make([]SimilarEffectSummary, 0, len(effects))
		for _, e := range effects {
			if e.Comparison == store.EffectComparisonNone {
				continue
			}
			effectRows = append(effectRows, SimilarEffectSummary{
				Dimension:  e.Dimension,
				Comparison: e.Comparison,
			})
		}
		signals, err := s.gatherChangeOutcomeSignals(c)
		if err != nil {
			return SimilarChangesResult{}, err
		}
		out.Changes = append(out.Changes, SimilarChangeRow{
			ID:          c.ID,
			TaskID:      c.TaskID,
			ChangeKind:  kind,
			Reason:      c.Reason,
			CreatedAt:   c.CreatedAt,
			Paths:       pathStrs,
			Effects:     effectRows,
			OutcomeKind: ClassifyChangeOutcome(signals),
		})
	}

	if changeKind != "" {
		patterns, err := s.store.ListChangePatternsByKind(changeKind, limit)
		if err != nil {
			return SimilarChangesResult{}, err
		}
		out.Patterns = patterns
	} else {
		seenPattern := map[string]struct{}{}
		for kind := range kindSet {
			patterns, err := s.store.ListChangePatternsByKind(kind, limit)
			if err != nil {
				return SimilarChangesResult{}, err
			}
			for _, p := range patterns {
				key := p.ChangeKind + "\x00" + p.OutcomeKind
				if _, dup := seenPattern[key]; dup {
					continue
				}
				seenPattern[key] = struct{}{}
				out.Patterns = append(out.Patterns, p)
			}
		}
		sort.Slice(out.Patterns, func(i, j int) bool {
			if out.Patterns[i].ChangeKind != out.Patterns[j].ChangeKind {
				return out.Patterns[i].ChangeKind < out.Patterns[j].ChangeKind
			}
			return out.Patterns[i].OutcomeKind < out.Patterns[j].OutcomeKind
		})
	}
	return out, nil
}

const (
	TendencyDirectionImprove = "improve"
	TendencyDirectionDamage  = "damage"
)

// TendencyRow describes what tends to help or hurt from aggregated change patterns.
type TendencyRow struct {
	ChangeKind    string `json:"change_kind"`
	OutcomeKind   string `json:"outcome_kind"`
	Direction     string `json:"direction"`
	CountPositive int    `json:"count_positive"`
	CountNegative int    `json:"count_negative"`
	LastSeen      string `json:"last_seen"`
}

// ListTendencies returns change-pattern rows where count_positive or count_negative ≥ 2.
func (s *Service) ListTendencies(ctx context.Context, limit int) ([]TendencyRow, error) {
	_ = ctx
	if limit <= 0 {
		limit = evidenceQueryDefaultLimit
	}
	if limit > evidenceQueryMaxLimit {
		limit = evidenceQueryMaxLimit
	}
	patterns, err := s.store.ListChangePatterns(64)
	if err != nil {
		return nil, err
	}
	var rows []TendencyRow
	for _, p := range patterns {
		if p.CountPositive >= 2 {
			rows = append(rows, TendencyRow{
				ChangeKind: p.ChangeKind, OutcomeKind: p.OutcomeKind,
				Direction:     TendencyDirectionImprove,
				CountPositive: p.CountPositive, CountNegative: p.CountNegative,
				LastSeen: p.LastSeen,
			})
		}
		if p.CountNegative >= 2 {
			rows = append(rows, TendencyRow{
				ChangeKind: p.ChangeKind, OutcomeKind: p.OutcomeKind,
				Direction:     TendencyDirectionDamage,
				CountPositive: p.CountPositive, CountNegative: p.CountNegative,
				LastSeen: p.LastSeen,
			})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		if a.LastSeen != b.LastSeen {
			return a.LastSeen > b.LastSeen
		}
		if a.ChangeKind != b.ChangeKind {
			return a.ChangeKind < b.ChangeKind
		}
		if a.OutcomeKind != b.OutcomeKind {
			return a.OutcomeKind < b.OutcomeKind
		}
		return a.Direction < b.Direction
	})
	if len(rows) > limit {
		rows = rows[:limit]
	}
	if rows == nil {
		rows = []TendencyRow{}
	}
	return rows, nil
}

func (s *Service) listChangesByKind(changeKind string, limit int) ([]store.Change, error) {
	all, err := s.store.ListActiveChanges()
	if err != nil {
		return nil, err
	}
	var matched []store.Change
	for _, c := range all {
		paths, err := s.store.ListChangePaths(c.ID)
		if err != nil {
			return nil, err
		}
		pathStrs := make([]string, 0, len(paths))
		for _, p := range paths {
			pathStrs = append(pathStrs, p.Path)
		}
		if InferChangeKind(pathStrs) != changeKind {
			continue
		}
		matched = append(matched, c)
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].CreatedAt != matched[j].CreatedAt {
			return matched[i].CreatedAt > matched[j].CreatedAt
		}
		return matched[i].RowID > matched[j].RowID
	})
	if len(matched) > limit {
		matched = matched[:limit]
	}
	if matched == nil {
		matched = []store.Change{}
	}
	return matched, nil
}
