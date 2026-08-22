package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	maxKnowledgeBodyBytes     = 8192
	maxKnowledgeEvidenceIDs   = 32
	knowledgeTopicDecision    = "decision"
	knowledgeTopicReflection  = "reflection"
	knowledgeTopicPattern     = "pattern"
	knowledgeTopicImprovement = "improvement"

	knowledgeSourceDecisionReconsideration = "decision_reconsideration"
	knowledgeSourceReflection              = "reflection"
	knowledgeSourceChangePattern           = "change_pattern"
	knowledgeSourceImprovement             = "improvement"
)

// KnowledgeBody is the structured body_json payload for engineering_knowledge rows.
type KnowledgeBody struct {
	DecisionID       string `json:"decision_id,omitempty"`
	Pattern          string `json:"pattern,omitempty"`
	Summary          string `json:"summary,omitempty"`
	SourceEntityType string `json:"source_entity_type,omitempty"`
	SourceEntityID   string `json:"source_entity_id,omitempty"`
	ChangeKind       string `json:"change_kind,omitempty"`
	OutcomeKind      string `json:"outcome_kind,omitempty"`
}

// EngineeringKnowledgeInput upserts one engineering_knowledge row.
type EngineeringKnowledgeInput struct {
	ID          string
	Title       string
	Body        KnowledgeBody
	Topic       string
	EvidenceIDs []string
	Confidence  float64
	Status      string
	SourceType  string
}

// ListEngineeringKnowledgeOpts filters list queries.
type ListEngineeringKnowledgeOpts struct {
	Topic  string
	Status string
	Limit  int
}

// SynthesizeKnowledgeResult counts rows created vs updated during synthesis.
type SynthesizeKnowledgeResult struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
}

func knowledgeStableID(sourceType, sourceEntityID string) string {
	return "ek:" + sourceType + ":" + sourceEntityID
}

func marshalKnowledgeBody(body KnowledgeBody) (string, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return "", &ErrValidation{Msg: "body must be a JSON object"}
	}
	if len(raw) > maxKnowledgeBodyBytes {
		return "", &ErrValidation{Msg: "body_json exceeds 8192 bytes"}
	}
	return string(raw), nil
}

func parseKnowledgeBody(raw string) (KnowledgeBody, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return KnowledgeBody{}, nil
	}
	var body KnowledgeBody
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		return KnowledgeBody{}, &ErrValidation{Msg: "body_json must be a JSON object"}
	}
	return body, nil
}

func validateKnowledgeEvidenceIDs(s *Service, ids []string) ([]string, error) {
	if len(ids) > maxKnowledgeEvidenceIDs {
		return nil, &ErrValidation{Msg: "evidence_ids exceeds 32 items"}
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			return nil, &ErrValidation{Msg: "evidence id is required"}
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		if _, err := s.store.GetEvidence(id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}

// UpsertEngineeringKnowledge persists one engineering_knowledge row.
func (s *Service) UpsertEngineeringKnowledge(ctx context.Context, in EngineeringKnowledgeInput) (store.EngineeringKnowledge, error) {
	_ = ctx
	bodyJSON, err := marshalKnowledgeBody(in.Body)
	if err != nil {
		return store.EngineeringKnowledge{}, err
	}
	evIDs, err := validateKnowledgeEvidenceIDs(s, in.EvidenceIDs)
	if err != nil {
		return store.EngineeringKnowledge{}, err
	}
	evJSON, err := marshalEvidenceIDs(evIDs)
	if err != nil {
		return store.EngineeringKnowledge{}, err
	}
	if err := validateConfidence01(in.Confidence); err != nil {
		return store.EngineeringKnowledge{}, err
	}
	topic := strings.TrimSpace(in.Topic)
	if topic == "" {
		return store.EngineeringKnowledge{}, &ErrValidation{Msg: "topic is required"}
	}
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return store.EngineeringKnowledge{}, &ErrValidation{Msg: "title is required"}
	}
	status := strings.TrimSpace(in.Status)
	if status == "" {
		status = store.KnowledgeStatusActive
	}
	if status != store.KnowledgeStatusActive && status != store.KnowledgeStatusSuperseded {
		return store.EngineeringKnowledge{}, &ErrValidation{Msg: "status must be active or superseded"}
	}

	id := strings.TrimSpace(in.ID)
	if id == "" {
		id = uuid.NewString()
	}

	return s.store.UpsertEngineeringKnowledge(store.EngineeringKnowledge{
		ID:              id,
		Title:           title,
		BodyJSON:        bodyJSON,
		Topic:           topic,
		EvidenceIDsJSON: evJSON,
		Confidence:      in.Confidence,
		Status:          status,
		SourceType:      defaultSource(in.SourceType),
	})
}

// GetEngineeringKnowledge loads one engineering_knowledge row by id.
func (s *Service) GetEngineeringKnowledge(ctx context.Context, id string) (store.EngineeringKnowledge, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	if id == "" {
		return store.EngineeringKnowledge{}, &ErrValidation{Msg: "knowledge id is required"}
	}
	return s.store.GetEngineeringKnowledge(id)
}

// ListEngineeringKnowledge returns knowledge rows with optional filters.
func (s *Service) ListEngineeringKnowledge(ctx context.Context, opts ListEngineeringKnowledgeOpts) ([]store.EngineeringKnowledge, error) {
	_ = ctx
	return s.store.ListEngineeringKnowledge(opts.Topic, opts.Status, opts.Limit)
}

func (s *Service) upsertSynthesizedKnowledge(in EngineeringKnowledgeInput) (created bool, err error) {
	sourceEntityID := strings.TrimSpace(in.Body.SourceEntityID)
	if sourceEntityID == "" {
		return false, &ErrValidation{Msg: "source_entity_id is required for synthesized knowledge"}
	}
	sourceType := strings.TrimSpace(in.SourceType)
	if sourceType == "" {
		return false, &ErrValidation{Msg: "source_type is required for synthesized knowledge"}
	}
	stableID := knowledgeStableID(sourceType, sourceEntityID)
	in.ID = stableID
	in.SourceType = sourceType

	_, getErr := s.store.GetEngineeringKnowledge(stableID)
	existed := getErr == nil

	row, err := s.UpsertEngineeringKnowledge(context.Background(), in)
	if err != nil {
		return false, err
	}
	_ = row
	return !existed, nil
}

// SynthesizeKnowledge refreshes change patterns then upserts knowledge from live sources.
func (s *Service) SynthesizeKnowledge(ctx context.Context) (SynthesizeKnowledgeResult, error) {
	if _, err := s.RefreshChangePatterns(ctx); err != nil {
		return SynthesizeKnowledgeResult{}, err
	}

	var result SynthesizeKnowledgeResult

	recons, err := s.store.ListAllDecisionReconsiderations()
	if err != nil {
		return SynthesizeKnowledgeResult{}, err
	}
	for _, r := range recons {
		body := KnowledgeBody{
			DecisionID:       r.DecisionID,
			Summary:          r.Reason,
			SourceEntityType: knowledgeSourceDecisionReconsideration,
			SourceEntityID:   r.ID,
		}
		title := fmt.Sprintf("Decision reconsideration: %s", strings.TrimSpace(r.Trigger))
		if title == "Decision reconsideration: " {
			title = "Decision reconsideration"
		}
		created, err := s.upsertSynthesizedKnowledge(EngineeringKnowledgeInput{
			Title:      title,
			Body:       body,
			Topic:      knowledgeTopicDecision,
			SourceType: knowledgeSourceDecisionReconsideration,
			Confidence: 0.7,
		})
		if err != nil {
			return SynthesizeKnowledgeResult{}, err
		}
		if created {
			result.Created++
		} else {
			result.Updated++
		}
	}

	reflections, err := s.store.ListAllReflections()
	if err != nil {
		return SynthesizeKnowledgeResult{}, err
	}
	for _, r := range reflections {
		summary := strings.TrimSpace(r.Summary)
		if summary == "" {
			continue
		}
		body := KnowledgeBody{
			Summary:          summary,
			SourceEntityType: EntityReflection,
			SourceEntityID:   r.ID,
		}
		created, err := s.upsertSynthesizedKnowledge(EngineeringKnowledgeInput{
			Title:      truncateTitle(summary, 120),
			Body:       body,
			Topic:      knowledgeTopicReflection,
			SourceType: knowledgeSourceReflection,
			Confidence: r.Confidence,
		})
		if err != nil {
			return SynthesizeKnowledgeResult{}, err
		}
		if created {
			result.Created++
		} else {
			result.Updated++
		}
	}

	patterns, err := s.store.ListChangePatterns(64)
	if err != nil {
		return SynthesizeKnowledgeResult{}, err
	}
	for _, p := range patterns {
		if p.CountPositive < 2 && p.CountNegative < 2 {
			continue
		}
		entityID := p.ChangeKind + "\x00" + p.OutcomeKind
		patternLabel := fmt.Sprintf("%s → %s", p.ChangeKind, p.OutcomeKind)
		body := KnowledgeBody{
			Pattern:          patternLabel,
			ChangeKind:       p.ChangeKind,
			OutcomeKind:      p.OutcomeKind,
			SourceEntityType: "change_pattern",
			SourceEntityID:   entityID,
		}
		created, err := s.upsertSynthesizedKnowledge(EngineeringKnowledgeInput{
			Title:      "Change pattern: " + patternLabel,
			Body:       body,
			Topic:      knowledgeTopicPattern,
			SourceType: knowledgeSourceChangePattern,
			Confidence: patternConfidence(p),
		})
		if err != nil {
			return SynthesizeKnowledgeResult{}, err
		}
		if created {
			result.Created++
		} else {
			result.Updated++
		}
	}

	improvements, err := s.store.ListAllImprovements()
	if err != nil {
		return SynthesizeKnowledgeResult{}, err
	}
	for _, im := range improvements {
		summary := strings.TrimSpace(im.Summary)
		if summary == "" {
			continue
		}
		body := KnowledgeBody{
			Summary:          summary,
			SourceEntityType: EntityImprovement,
			SourceEntityID:   im.ID,
		}
		evIDs, err := parseEvidenceIDsJSON(im.EvidenceIDsJSON)
		if err != nil {
			return SynthesizeKnowledgeResult{}, err
		}
		created, err := s.upsertSynthesizedKnowledge(EngineeringKnowledgeInput{
			Title:       truncateTitle(summary, 120),
			Body:        body,
			Topic:       knowledgeTopicImprovement,
			EvidenceIDs: evIDs,
			SourceType:  knowledgeSourceImprovement,
			Confidence:  im.Confidence,
		})
		if err != nil {
			return SynthesizeKnowledgeResult{}, err
		}
		if created {
			result.Created++
		} else {
			result.Updated++
		}
	}

	return result, nil
}

func patternConfidence(p store.ChangePattern) float64 {
	total := p.CountPositive + p.CountNegative
	if total <= 0 {
		return 0.5
	}
	if total >= 4 {
		return 0.9
	}
	return 0.7
}

func truncateTitle(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

// SuccessfulApproachRow merges worked outcomes and active knowledge rows.
type SuccessfulApproachRow struct {
	ID        string `json:"id"`
	Source    string `json:"source"`
	Kind      string `json:"kind"`
	Title     string `json:"title,omitempty"`
	Summary   string `json:"summary,omitempty"`
	TaskID    string `json:"task_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

// SuccessfulApproachesOpts bounds merged successful-approach queries.
type SuccessfulApproachesOpts struct {
	TaskID string
	Limit  int
}

type successfulApproachCandidate struct {
	id, source, kind, title, summary, taskID, createdAt string
}

func knowledgeRowIsImproveDirection(body KnowledgeBody) bool {
	if body.OutcomeKind == "" {
		return false
	}
	return isPositiveOutcome(body.OutcomeKind)
}

// ListSuccessfulApproaches merges worked outcomes and active improvement/pattern knowledge.
func (s *Service) ListSuccessfulApproaches(ctx context.Context, opts SuccessfulApproachesOpts) ([]SuccessfulApproachRow, error) {
	limit := clampEvidenceLimit(opts.Limit)
	worked, err := s.ListWorkedApproaches(ctx, EvidenceQueryOpts{
		TaskID: opts.TaskID,
		Limit:  0,
	})
	if err != nil {
		return nil, err
	}
	var candidates []successfulApproachCandidate
	for _, w := range worked {
		candidates = append(candidates, successfulApproachCandidate{
			id: w.ID, source: "worked", kind: w.Kind,
			summary: w.Summary, taskID: w.TaskID, createdAt: w.CreatedAt,
		})
	}
	for _, topic := range []string{knowledgeTopicImprovement, knowledgeTopicPattern} {
		rows, err := s.store.ListEngineeringKnowledge(topic, store.KnowledgeStatusActive, 64)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			body, err := parseKnowledgeBody(row.BodyJSON)
			if err != nil {
				return nil, err
			}
			if topic == knowledgeTopicPattern && !knowledgeRowIsImproveDirection(body) {
				continue
			}
			summary := strings.TrimSpace(body.Summary)
			if summary == "" {
				summary = strings.TrimSpace(body.Pattern)
			}
			candidates = append(candidates, successfulApproachCandidate{
				id: row.ID, source: "knowledge", kind: topic,
				title: row.Title, summary: summary, createdAt: row.UpdatedAt,
			})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if a.createdAt != b.createdAt {
			return a.createdAt > b.createdAt
		}
		return a.id > b.id
	})
	seen := map[string]struct{}{}
	out := make([]SuccessfulApproachRow, 0, limit)
	for _, c := range candidates {
		if _, dup := seen[c.id]; dup {
			continue
		}
		seen[c.id] = struct{}{}
		out = append(out, SuccessfulApproachRow{
			ID: c.id, Source: c.source, Kind: c.kind,
			Title: c.title, Summary: c.summary, TaskID: c.taskID, CreatedAt: c.createdAt,
		})
		if len(out) >= limit {
			break
		}
	}
	if out == nil {
		out = []SuccessfulApproachRow{}
	}
	return out, nil
}

func parseEvidenceIDsJSON(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	var ids []string
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return nil, &ErrValidation{Msg: "evidence_ids_json must be a JSON array"}
	}
	return ids, nil
}
