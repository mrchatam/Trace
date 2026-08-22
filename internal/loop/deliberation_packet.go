package loop

import (
	"context"
	"database/sql"
	"errors"
	"sort"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	maxOpenUncertaintiesCap       = 16
	maxVerificationDebtCap        = 8
	maxRecentChangesCap           = 8
	maxChangePathsCap             = 16
	maxChangeEffectsCap           = 8
	maxOpenRegressionsCap         = 8
	maxHistoricalRelationshipsCap = 8
	investigateLookaheadMax       = 256
)

// PhaseProfile holds S06 phase-specific compiler and section caps.
type PhaseProfile struct {
	ContextMaxItems     int
	ContextIncludeWhy   bool
	ContextIncludeMD    bool
	RelatedDepth        int
	OpenUncertaintyCap  int
	RecentChangesCap    int
	TrimLookahead       bool
	IncludeOpenRegress  bool
	EmphasizeVerifyDebt bool
}

// PhaseContextProfile returns locked S06 emphasis for a recommended deliberation phase.
func PhaseContextProfile(phase deliberation.Phase) PhaseProfile {
	switch phase {
	case deliberation.PhaseInvestigate:
		return PhaseProfile{
			ContextMaxItems:    24,
			ContextIncludeWhy:  false, // why section built separately in next.go
			RelatedDepth:       1,
			OpenUncertaintyCap: 16,
			TrimLookahead:      true,
		}
	case deliberation.PhaseExplore:
		return PhaseProfile{
			ContextMaxItems:    24,
			RelatedDepth:       2,
			OpenUncertaintyCap: 12,
		}
	case deliberation.PhaseExecute:
		return PhaseProfile{
			ContextMaxItems:    32,
			ContextIncludeMD:   true,
			RelatedDepth:       2,
			OpenUncertaintyCap: 4,
			RecentChangesCap:   8,
		}
	case deliberation.PhaseVerify:
		return PhaseProfile{
			ContextMaxItems:     20,
			ContextIncludeWhy:   true,
			EmphasizeVerifyDebt: true,
			OpenUncertaintyCap:  8,
		}
	case deliberation.PhaseTest, deliberation.PhaseEvaluate:
		return PhaseProfile{
			RecentChangesCap:   8,
			OpenUncertaintyCap: 8,
			RelatedDepth:       2,
		}
	case deliberation.PhaseReflect, deliberation.PhaseReplan:
		return PhaseProfile{
			IncludeOpenRegress: true,
			OpenUncertaintyCap: 8,
			RelatedDepth:       2,
		}
	default:
		return PhaseProfile{
			OpenUncertaintyCap: 8,
			RelatedDepth:       2,
		}
	}
}

type OpenUncertaintyItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Severity string `json:"severity"`
	Status   string `json:"status"`
	Kind     string `json:"kind"`
}

type OpenUncertaintiesSection struct {
	Freshness string                `json:"freshness"`
	Items     []OpenUncertaintyItem `json:"items"`
}

type VerificationDebtSection struct {
	Freshness string           `json:"freshness"`
	Present   bool             `json:"present"`
	Items     []store.DebtItem `json:"items"`
}

type RecentChangeEffect struct {
	Dimension  string `json:"dimension"`
	Comparison string `json:"comparison"`
}

type RecentChangeItem struct {
	ID        string               `json:"id"`
	GitCommit string               `json:"git_commit,omitempty"`
	Status    string               `json:"status"`
	Paths     []string             `json:"paths"`
	Effects   []RecentChangeEffect `json:"effects"`
}

type RecentChangesSection struct {
	Freshness string             `json:"freshness"`
	Items     []RecentChangeItem `json:"items"`
}

type HistoricalRelationshipsSection struct {
	Freshness string                       `json:"freshness"`
	Items     []HistoricalRelationshipItem `json:"items"`
}

type HistoricalRelationshipItem struct {
	Rel        string  `json:"rel"`
	FromType   string  `json:"from_type"`
	FromID     string  `json:"from_id"`
	ToType     string  `json:"to_type"`
	ToID       string  `json:"to_id"`
	Confidence float64 `json:"confidence"`
}

type OpenRegressionItem struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Attribution string `json:"attribution"`
	Status      string `json:"status"`
}

type DeliberationSection struct {
	Freshness       string                    `json:"freshness"`
	Phase           deliberation.Phase        `json:"phase"`
	WhySelected     deliberation.ReasonCode   `json:"why_selected"`
	PolicyInputs    deliberation.PolicyInputs `json:"policy_inputs"`
	HopCount        int                       `json:"hop_count"`
	Stopped         bool                      `json:"stopped"`
	StopReason      string                    `json:"stop_reason,omitempty"`
	LastPhase       deliberation.Phase        `json:"last_phase,omitempty"`
	OpenRegressions []OpenRegressionItem      `json:"open_regressions,omitempty"`
}

type StatusDeliberation struct {
	Phase            deliberation.Phase        `json:"phase"`
	RecommendedPhase deliberation.Phase        `json:"recommended_phase"`
	WhySelected      deliberation.ReasonCode   `json:"why_selected"`
	StopReason       string                    `json:"stop_reason,omitempty"`
	HopCount         int                       `json:"hop_count"`
	Stopped          bool                      `json:"stopped"`
	Blocked          bool                      `json:"blocked"`
	NeedsPhase       deliberation.Phase        `json:"needs_phase,omitempty"`
	PolicyInputs     deliberation.PolicyInputs `json:"policy_inputs"`
}

func loadDeliberationState(ctx context.Context, dom *domain.Service, taskID, goalID string) deliberation.State {
	st, err := dom.GetDeliberationState(ctx, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return deliberation.InitialState(taskID, goalID)
		}
		return deliberation.InitialState(taskID, goalID)
	}
	return st
}

func buildDeliberationSection(
	dState deliberation.State,
	inputs deliberation.PolicyInputs,
	profile PhaseProfile,
	freshness string,
) DeliberationSection {
	phase, reason, stopped := deliberation.SelectNext(dState, inputs)
	sec := DeliberationSection{
		Freshness:    freshness,
		Phase:        phase,
		WhySelected:  reason,
		PolicyInputs: inputs,
		HopCount:     dState.HopCount,
		Stopped:      stopped || dState.Stopped,
		LastPhase:    dState.LastPhase,
	}
	if dState.Stopped {
		sec.StopReason = dState.StopReason
	} else if stopped {
		sec.StopReason = string(reason)
	}
	return sec
}

func buildOpenUncertaintiesSection(st *store.Store, taskID string, cap int, freshness string) (OpenUncertaintiesSection, error) {
	if cap <= 0 {
		cap = maxOpenUncertaintiesCap
	}
	rows, err := st.ListOpenUncertaintiesByTaskID(taskID, cap)
	if err != nil {
		return OpenUncertaintiesSection{}, err
	}
	items := make([]OpenUncertaintyItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, OpenUncertaintyItem{
			ID: row.ID, Title: row.Title, Severity: row.Severity,
			Status: row.Status, Kind: row.Kind,
		})
	}
	return OpenUncertaintiesSection{Freshness: freshness, Items: items}, nil
}

func buildVerificationDebtSection(ctx context.Context, dom *domain.Service, taskID string, freshness string) (VerificationDebtSection, error) {
	present, err := dom.HasVerificationDebt(ctx, taskID)
	if err != nil {
		return VerificationDebtSection{}, err
	}
	items, err := dom.ListVerificationDebtSummary(ctx, taskID)
	if err != nil {
		return VerificationDebtSection{}, err
	}
	if len(items) > maxVerificationDebtCap {
		items = items[:maxVerificationDebtCap]
	}
	return VerificationDebtSection{
		Freshness: freshness,
		Present:   present,
		Items:     items,
	}, nil
}

func buildRecentChangesSection(st *store.Store, taskID string, cap int, freshness string) (RecentChangesSection, error) {
	if cap <= 0 {
		return RecentChangesSection{Freshness: freshness, Items: []RecentChangeItem{}}, nil
	}
	if cap > maxRecentChangesCap {
		cap = maxRecentChangesCap
	}
	changes, err := st.ListChangesByTaskID(taskID)
	if err != nil {
		return RecentChangesSection{}, err
	}
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].CreatedAt == changes[j].CreatedAt {
			return changes[i].ID > changes[j].ID
		}
		return changes[i].CreatedAt > changes[j].CreatedAt
	})
	if len(changes) > cap {
		changes = changes[:cap]
	}
	items := make([]RecentChangeItem, 0, len(changes))
	for _, ch := range changes {
		paths, err := st.ListChangePaths(ch.ID)
		if err != nil {
			return RecentChangesSection{}, err
		}
		pathStrs := make([]string, 0, len(paths))
		for i, p := range paths {
			if i >= maxChangePathsCap {
				break
			}
			pathStrs = append(pathStrs, p.Path)
		}
		effects, err := st.ListEffectsByChangeID(ch.ID)
		if err != nil {
			return RecentChangesSection{}, err
		}
		effItems := make([]RecentChangeEffect, 0, len(effects))
		for i, e := range effects {
			if i >= maxChangeEffectsCap {
				break
			}
			if e.Comparison == "" {
				continue
			}
			effItems = append(effItems, RecentChangeEffect{
				Dimension:  e.Dimension,
				Comparison: e.Comparison,
			})
		}
		items = append(items, RecentChangeItem{
			ID: ch.ID, GitCommit: ch.GitCommit, Status: ch.Status,
			Paths: pathStrs, Effects: effItems,
		})
	}
	return RecentChangesSection{Freshness: freshness, Items: items}, nil
}

func buildHistoricalRelationshipsSection(st *store.Store, freshness string) (HistoricalRelationshipsSection, error) {
	empty := HistoricalRelationshipsSection{Freshness: freshness, Items: []HistoricalRelationshipItem{}}
	if st == nil {
		return empty, nil
	}
	observed, err := st.ListLinksByRel(domain.RelObservedRelationship)
	if err != nil {
		return HistoricalRelationshipsSection{}, err
	}
	caused, err := st.ListLinksByRel(domain.RelCausedBy)
	if err != nil {
		return HistoricalRelationshipsSection{}, err
	}
	merged := make([]store.EntityLink, 0, len(observed)+len(caused))
	merged = append(merged, observed...)
	for _, l := range caused {
		fromLinks, err := st.ListLinksFrom(l.FromType, l.FromID)
		if err != nil {
			return HistoricalRelationshipsSection{}, err
		}
		hasEvidence := false
		for _, fl := range fromLinks {
			if fl.Rel == domain.RelRelationshipSupportedBy {
				hasEvidence = true
				break
			}
		}
		if hasEvidence {
			merged = append(merged, l)
		}
	}
	sort.Slice(merged, func(i, j int) bool {
		if merged[i].CreatedAt == merged[j].CreatedAt {
			return merged[i].ID > merged[j].ID
		}
		return merged[i].CreatedAt > merged[j].CreatedAt
	})
	if len(merged) > maxHistoricalRelationshipsCap {
		merged = merged[:maxHistoricalRelationshipsCap]
	}
	items := make([]HistoricalRelationshipItem, 0, len(merged))
	for _, l := range merged {
		items = append(items, HistoricalRelationshipItem{
			Rel:        l.Rel,
			FromType:   l.FromType,
			FromID:     l.FromID,
			ToType:     l.ToType,
			ToID:       l.ToID,
			Confidence: l.Confidence,
		})
	}
	return HistoricalRelationshipsSection{Freshness: freshness, Items: items}, nil
}

func buildOpenRegressionsSummary(ctx context.Context, dom *domain.Service, taskID string) ([]OpenRegressionItem, error) {
	regs, err := dom.ListOpenRegressions(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if len(regs) > maxOpenRegressionsCap {
		regs = regs[:maxOpenRegressionsCap]
	}
	out := make([]OpenRegressionItem, 0, len(regs))
	for _, r := range regs {
		out = append(out, OpenRegressionItem{
			ID: r.ID, Summary: r.Summary, Attribution: r.Attribution, Status: r.Status,
		})
	}
	return out, nil
}

func trimLookaheadSummary(summary string, max int) string {
	if max <= 0 || len(summary) <= max {
		return summary
	}
	return summary[:max]
}

func statusNeedsPhase(recommended deliberation.Phase, blocked, stopped, p19Sat bool) deliberation.Phase {
	if blocked && !stopped && !p19Sat {
		return recommended
	}
	return ""
}

func statusBlocked(inputs deliberation.PolicyInputs) bool {
	return inputs.BlockingUncertaintyCount > 0 || inputs.OpenRegression || inputs.VerificationIncomplete
}
