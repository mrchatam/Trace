package compiler

import (
	"context"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

const evidenceSectionCap = 8

const maxScoresJSONInPacket = 512

// EvaluationItem is a compact evaluation outcome for context packets.
type EvaluationItem struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	Summary    string `json:"summary,omitempty"`
	ScoresJSON string `json:"scores_json,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// ReflectionItem is a compact reflection for context packets (no assumption blobs).
type ReflectionItem struct {
	ID        string `json:"id"`
	Summary   string `json:"summary"`
	CreatedAt string `json:"created_at"`
}

// PlanningEvidenceItem is one mixed planning signal (regression, failed test, improvement).
type PlanningEvidenceItem struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	CreatedAt  string `json:"created_at"`
}

// TendencyItem is a compact tendency row for context packets.
type TendencyItem struct {
	ChangeKind    string `json:"change_kind"`
	OutcomeKind   string `json:"outcome_kind"`
	Direction     string `json:"direction"`
	CountPositive int    `json:"count_positive"`
	CountNegative int    `json:"count_negative"`
	LastSeen      string `json:"last_seen"`
}

// SuccessfulApproachItem is a worked outcome or knowledge row for context packets.
type SuccessfulApproachItem struct {
	ID        string `json:"id"`
	Source    string `json:"source"`
	Kind      string `json:"kind"`
	Title     string `json:"title,omitempty"`
	Summary   string `json:"summary,omitempty"`
	TaskID    string `json:"task_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

// EvidenceBundle holds task-scoped evaluations, reflections, and mixed planning evidence.
type EvidenceBundle struct {
	Evaluations          []EvaluationItem
	Reflections          []ReflectionItem
	PlanningEvidence     []PlanningEvidenceItem
	Tendencies           []TendencyItem
	SuccessfulApproaches []SuccessfulApproachItem
}

func attachEvidenceSections(pkt *Packet, st *store.Store, taskID string) {
	bundle := buildEvidenceSections(st, taskID)
	pkt.Evaluations = bundle.Evaluations
	pkt.Reflections = bundle.Reflections
	pkt.PlanningEvidence = bundle.PlanningEvidence
	pkt.Tendencies = bundle.Tendencies
	pkt.SuccessfulApproaches = bundle.SuccessfulApproaches
}

// BuildEvidenceSections loads capped task-scoped evidence for context and loop next.
func BuildEvidenceSections(st *store.Store, taskID string) EvidenceBundle {
	return buildEvidenceSections(st, taskID)
}

func buildEvidenceSections(st *store.Store, taskID string) EvidenceBundle {
	return EvidenceBundle{
		Evaluations:          buildEvaluations(st, taskID),
		Reflections:          buildReflections(st, taskID),
		PlanningEvidence:     buildPlanningEvidence(st, taskID),
		Tendencies:           buildTendencies(st),
		SuccessfulApproaches: buildSuccessfulApproaches(st),
	}
}

func buildEvaluations(st *store.Store, taskID string) []EvaluationItem {
	rows, err := st.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindEvaluation)
	if err != nil || len(rows) == 0 {
		return []EvaluationItem{}
	}
	sortOutcomeResultsNewestFirst(rows)
	if len(rows) > evidenceSectionCap {
		rows = rows[:evidenceSectionCap]
	}
	out := make([]EvaluationItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, EvaluationItem{
			ID:         r.ID,
			TaskID:     r.TaskID,
			Summary:    r.Summary,
			ScoresJSON: truncateScoresJSON(r.ScoresJSON),
			CreatedAt:  r.CreatedAt,
		})
	}
	return out
}

func buildReflections(st *store.Store, taskID string) []ReflectionItem {
	rows, err := st.ListReflectionsByTaskID(taskID)
	if err != nil || len(rows) == 0 {
		return []ReflectionItem{}
	}
	sortReflectionsNewestFirst(rows)
	if len(rows) > evidenceSectionCap {
		rows = rows[:evidenceSectionCap]
	}
	out := make([]ReflectionItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, ReflectionItem{
			ID:        r.ID,
			Summary:   r.Summary,
			CreatedAt: r.CreatedAt,
		})
	}
	return out
}

func buildPlanningEvidence(st *store.Store, taskID string) []PlanningEvidenceItem {
	var merged []PlanningEvidenceItem

	if regs, err := st.ListOpenRegressions(taskID); err == nil {
		for _, r := range regs {
			title := strings.TrimSpace(r.Dimension)
			if title == "" {
				title = "regression"
			}
			merged = append(merged, PlanningEvidenceItem{
				EntityType: "regression",
				EntityID:   r.ID,
				Title:      title,
				Summary:    r.Summary,
				CreatedAt:  r.CreatedAt,
			})
		}
	}

	if fails, err := st.ListFailedTestOutcomes(evidenceSectionCap, taskID); err == nil {
		for _, o := range fails {
			title := strings.TrimSpace(o.TestName)
			if title == "" {
				title = "test"
			}
			summary := strings.TrimSpace(o.Summary)
			if summary == "" {
				summary = o.TestStatus
			}
			merged = append(merged, PlanningEvidenceItem{
				EntityType: "outcome_result",
				EntityID:   o.ID,
				Title:      title,
				Summary:    summary,
				CreatedAt:  o.CreatedAt,
			})
		}
	}

	if imps, err := st.ListImprovementsByTaskID(taskID); err == nil {
		for _, imp := range imps {
			title := strings.TrimSpace(imp.Dimension)
			if title == "" {
				title = "improvement"
			}
			merged = append(merged, PlanningEvidenceItem{
				EntityType: "improvement",
				EntityID:   imp.ID,
				Title:      title,
				Summary:    imp.Summary,
				CreatedAt:  imp.CreatedAt,
			})
		}
	}

	if len(merged) == 0 {
		return []PlanningEvidenceItem{}
	}
	sortPlanningEvidenceNewestFirst(merged)
	if len(merged) > evidenceSectionCap {
		merged = merged[:evidenceSectionCap]
	}
	return merged
}

func buildTendencies(st *store.Store) []TendencyItem {
	svc := domain.New(st)
	rows, err := svc.ListTendencies(context.Background(), evidenceSectionCap)
	if err != nil || len(rows) == 0 {
		return []TendencyItem{}
	}
	out := make([]TendencyItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, TendencyItem{
			ChangeKind: r.ChangeKind, OutcomeKind: r.OutcomeKind, Direction: r.Direction,
			CountPositive: r.CountPositive, CountNegative: r.CountNegative, LastSeen: r.LastSeen,
		})
	}
	return out
}

func buildSuccessfulApproaches(st *store.Store) []SuccessfulApproachItem {
	svc := domain.New(st)
	rows, err := svc.ListSuccessfulApproaches(context.Background(), domain.SuccessfulApproachesOpts{
		Limit: evidenceSectionCap,
	})
	if err != nil || len(rows) == 0 {
		return []SuccessfulApproachItem{}
	}
	out := make([]SuccessfulApproachItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, SuccessfulApproachItem{
			ID: r.ID, Source: r.Source, Kind: r.Kind,
			Title: r.Title, Summary: r.Summary, TaskID: r.TaskID, CreatedAt: r.CreatedAt,
		})
	}
	return out
}

func sortOutcomeResultsNewestFirst(rows []store.OutcomeResult) {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].CreatedAt != rows[j].CreatedAt {
			return rows[i].CreatedAt > rows[j].CreatedAt
		}
		return rows[i].ID > rows[j].ID
	})
}

func sortReflectionsNewestFirst(rows []store.Reflection) {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].CreatedAt != rows[j].CreatedAt {
			return rows[i].CreatedAt > rows[j].CreatedAt
		}
		return rows[i].ID > rows[j].ID
	})
}

func sortPlanningEvidenceNewestFirst(rows []PlanningEvidenceItem) {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].CreatedAt != rows[j].CreatedAt {
			return rows[i].CreatedAt > rows[j].CreatedAt
		}
		return rows[i].EntityID > rows[j].EntityID
	})
}

func truncateScoresJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || len(raw) <= maxScoresJSONInPacket {
		return raw
	}
	return raw[:maxScoresJSONInPacket] + "…"
}
