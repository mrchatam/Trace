package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	maxRegressionSummaryBytes   = 4096
	maxRegressionDimensionBytes = 64
	maxReflectionArrayItems     = 32
	maxDependencyRefBytes       = 512
	maxUsefulTestNameBytes      = 256
	maxBroadenTestsNoteBytes    = 256
)

const (
	depKindPath   = "path"
	depKindSymbol = "symbol"
	depKindFile   = "file"
)

// EvaluationRegressionInput derives a thin regression from a stored evaluation.
// Attribution is always persisted as correlated — no override on create.
type EvaluationRegressionInput struct {
	OutcomeID  string
	TaskID     string
	Actor      string
	Summary    string
	SourceType string
	Confidence float64
}

// EffectRegressionInput derives a thin regression from a contradicted effect.
type EffectRegressionInput struct {
	EffectID   string
	TaskID     string
	Actor      string
	Summary    string
	SourceType string
	Confidence float64
}

// ReflectionInput creates a structured reflection (JSON arrays, no essay body).
type ReflectionInput struct {
	TaskID                   string
	Summary                  string
	InvalidatedAssumptionIDs []string
	NewDependencies          []DependencyRef
	UsefulTests              []string
	BroadenTestsNote         string
	Actor                    string
	SourceType               string
	Confidence               float64
}

// DependencyRef is one structured dependency on a reflection.
type DependencyRef struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}

// RelInput is an observed or causal relationship edge.
type RelInput struct {
	FromType   string
	FromID     string
	ToType     string
	ToID       string
	Confidence float64
	SourceType string
}

func (s *Service) insertLink(fromType, fromID, rel, toType, toID, sourceType string, confidence float64) error {
	if sourceType == "" {
		sourceType = DefaultSourceType
	}
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   fromType,
		FromID:     fromID,
		Rel:        rel,
		ToType:     toType,
		ToID:       toID,
		SourceType: sourceType,
		Confidence: confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(fromType, fromID, rel, toType, toID, LinkMeta{SourceType: sourceType, Confidence: confidence})
}

func (s *Service) requireKnownEntity(typ, id string) error {
	id = strings.TrimSpace(id)
	if typ == "" || id == "" {
		return &ErrValidation{Msg: "relationship endpoints are required"}
	}
	var err error
	switch typ {
	case EntityGoal:
		_, err = s.store.GetGoal(id)
	case EntityTask:
		_, err = s.store.GetTask(id)
	case EntityDecision:
		_, err = s.store.GetDecision(id)
	case EntityAssumption:
		_, err = s.store.GetAssumption(id)
	case EntityDiscovery:
		_, err = s.store.GetDiscovery(id)
	case EntityPlanChange:
		_, err = s.store.GetPlanChange(id)
	case EntityClaim:
		_, err = s.store.GetClaim(id)
	case EntityEvidence:
		_, err = s.store.GetEvidence(id)
	case EntityReview:
		_, err = s.store.GetReview(id)
	case EntityPlanScope:
		_, err = s.store.GetPlanScope(id)
	case EntityUncertainty:
		_, err = s.store.GetUncertainty(id)
	case EntityHypothesis:
		_, err = s.store.GetHypothesis(id)
	case EntityChange:
		_, err = s.store.GetChange(id)
	case EntityEffect:
		_, err = s.store.GetEffect(id)
	case EntityOutcomeResult:
		_, err = s.store.GetOutcomeResult(id)
	case EntityBaseline:
		_, err = s.store.GetBaseline(id)
	case EntityRegression:
		_, err = s.store.GetRegression(id)
	case EntityImprovement:
		_, err = s.store.GetImprovement(id)
	case EntityReflection:
		_, err = s.store.GetReflection(id)
	default:
		return nil
	}
	return err
}

func validateConfidence01(c float64) error {
	if math.IsNaN(c) || math.IsInf(c, 0) || c < 0 || c > 1 {
		return &ErrValidation{Msg: "confidence must be in [0, 1]"}
	}
	return nil
}

func regressionDimensionFromComparison(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "{}" {
		return "", &ErrValidation{Msg: "evaluation has no regression flags"}
	}
	var cmp ScoreComparison
	if err := json.Unmarshal([]byte(raw), &cmp); err != nil {
		return "", &ErrValidation{Msg: "comparison_json must be valid JSON"}
	}
	if cmp.OverallRegression {
		return "overall", nil
	}
	var keys []string
	for k, d := range cmp.Dimensions {
		if d.Regression {
			keys = append(keys, k)
		}
	}
	if len(keys) == 0 {
		return "", &ErrValidation{Msg: "evaluation has no regression flags"}
	}
	sort.Strings(keys)
	return keys[0], nil
}

func existingRegression(err error, row store.Regression) (store.Regression, bool, error) {
	if err == nil {
		return row, true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return store.Regression{}, false, nil
	}
	return store.Regression{}, false, err
}

// RecordRegressionFromEvaluation inspects stored comparison_json and persists correlated.
func (s *Service) RecordRegressionFromEvaluation(ctx context.Context, in EvaluationRegressionInput) (store.Regression, error) {
	_ = ctx
	task, err := s.requireTask(strings.TrimSpace(in.TaskID))
	if err != nil {
		return store.Regression{}, err
	}
	outcomeID := strings.TrimSpace(in.OutcomeID)
	if outcomeID == "" {
		return store.Regression{}, &ErrValidation{Msg: "outcome id is required"}
	}
	out, err := s.store.GetOutcomeResult(outcomeID)
	if err != nil {
		return store.Regression{}, err
	}
	if out.Kind != store.OutcomeKindEvaluation {
		return store.Regression{}, &ErrValidation{Msg: "outcome must be kind=evaluation"}
	}
	if out.TaskID != task.ID {
		return store.Regression{}, &ErrValidation{Msg: "task_id must match outcome.task_id"}
	}
	dim, err := regressionDimensionFromComparison(out.ComparisonJSON)
	if err != nil {
		return store.Regression{}, err
	}
	if err := failClosedMaxBytes("dimension", dim, maxRegressionDimensionBytes); err != nil {
		return store.Regression{}, err
	}
	if err := failClosedMaxBytes("summary", in.Summary, maxRegressionSummaryBytes); err != nil {
		return store.Regression{}, err
	}
	if err := validateConfidence01(in.Confidence); err != nil {
		return store.Regression{}, err
	}

	got, err := s.store.GetRegressionBySource(store.RegressionSourceEvaluation, out.ID, dim)
	existing, found, err := existingRegression(err, got)
	if err != nil {
		return store.Regression{}, err
	}
	if found {
		return existing, nil
	}

	row, err := s.store.UpsertRegression(store.Regression{
		ID:          uuid.NewString(),
		TaskID:      task.ID,
		SourceKind:  store.RegressionSourceEvaluation,
		SourceID:    out.ID,
		Dimension:   dim,
		Attribution: store.RegressionAttributionCorrelated,
		Status:      store.RegressionStatusOpen,
		Summary:     in.Summary,
		Actor:       strings.TrimSpace(in.Actor),
		SourceType:  defaultSource(in.SourceType),
		Confidence:  in.Confidence,
	})
	if err != nil {
		return store.Regression{}, err
	}
	if err := s.insertLink(EntityRegression, row.ID, RelRegressionFromEvaluation, EntityOutcomeResult, out.ID, row.SourceType, row.Confidence); err != nil {
		return store.Regression{}, err
	}
	if err := s.appendCreated(EntityRegression, row.ID, row.Summary); err != nil {
		return store.Regression{}, err
	}
	if err := s.appendNamed(EventRegressionRecorded, EntityRegression, row.ID, map[string]string{
		"source_kind": row.SourceKind,
		"source_id":   row.SourceID,
		"dimension":   row.Dimension,
		"attribution": row.Attribution,
		"task_id":     row.TaskID,
	}); err != nil {
		return store.Regression{}, err
	}
	return row, nil
}

// RecordRegressionFromContradictedEffect persists correlated from a contradicted effect.
func (s *Service) RecordRegressionFromContradictedEffect(ctx context.Context, in EffectRegressionInput) (store.Regression, error) {
	_ = ctx
	task, err := s.requireTask(strings.TrimSpace(in.TaskID))
	if err != nil {
		return store.Regression{}, err
	}
	effectID := strings.TrimSpace(in.EffectID)
	if effectID == "" {
		return store.Regression{}, &ErrValidation{Msg: "effect id is required"}
	}
	eff, err := s.store.GetEffect(effectID)
	if err != nil {
		return store.Regression{}, err
	}
	if eff.Comparison != store.EffectComparisonContradicted {
		return store.Regression{}, &ErrValidation{Msg: "effect comparison must be contradicted"}
	}
	chg, err := s.store.GetChange(eff.ChangeID)
	if err != nil {
		return store.Regression{}, err
	}
	if chg.TaskID != task.ID {
		return store.Regression{}, &ErrValidation{Msg: "task_id must match parent change.task_id"}
	}
	if err := failClosedMaxBytes("summary", in.Summary, maxRegressionSummaryBytes); err != nil {
		return store.Regression{}, err
	}
	if err := validateConfidence01(in.Confidence); err != nil {
		return store.Regression{}, err
	}

	got, err := s.store.GetRegressionBySource(store.RegressionSourceContradictedEffect, eff.ID, eff.Dimension)
	existing, found, err := existingRegression(err, got)
	if err != nil {
		return store.Regression{}, err
	}
	if found {
		if err := s.associateRegressionWithChange(existing.ID, eff.ChangeID, existing.SourceType, existing.Confidence); err != nil {
			return store.Regression{}, err
		}
		return existing, nil
	}

	row, err := s.store.UpsertRegression(store.Regression{
		ID:          uuid.NewString(),
		TaskID:      task.ID,
		SourceKind:  store.RegressionSourceContradictedEffect,
		SourceID:    eff.ID,
		Dimension:   eff.Dimension,
		Attribution: store.RegressionAttributionCorrelated,
		Status:      store.RegressionStatusOpen,
		Summary:     in.Summary,
		Actor:       strings.TrimSpace(in.Actor),
		SourceType:  defaultSource(in.SourceType),
		Confidence:  in.Confidence,
	})
	if err != nil {
		return store.Regression{}, err
	}
	if err := s.insertLink(EntityRegression, row.ID, RelRegressionFromEffect, EntityEffect, eff.ID, row.SourceType, row.Confidence); err != nil {
		return store.Regression{}, err
	}
	if err := s.associateRegressionWithChange(row.ID, eff.ChangeID, row.SourceType, row.Confidence); err != nil {
		return store.Regression{}, err
	}
	if err := s.appendCreated(EntityRegression, row.ID, row.Summary); err != nil {
		return store.Regression{}, err
	}
	if err := s.appendNamed(EventRegressionRecorded, EntityRegression, row.ID, map[string]string{
		"source_kind": row.SourceKind,
		"source_id":   row.SourceID,
		"dimension":   row.Dimension,
		"attribution": row.Attribution,
		"task_id":     row.TaskID,
	}); err != nil {
		return store.Regression{}, err
	}
	return row, nil
}

func (s *Service) associateRegressionWithChange(regressionID, changeID, sourceType string, confidence float64) error {
	regressionID = strings.TrimSpace(regressionID)
	changeID = strings.TrimSpace(changeID)
	if regressionID == "" || changeID == "" {
		return &ErrValidation{Msg: "regression id and change id are required"}
	}
	if _, err := s.store.GetRegression(regressionID); err != nil {
		return err
	}
	if _, err := s.store.GetChange(changeID); err != nil {
		return err
	}
	if sourceType == "" {
		sourceType = DefaultSourceType
	}
	inserted, _, err := s.store.InsertLinkOrIgnore(store.EntityLink{
		FromType:   EntityRegression,
		FromID:     regressionID,
		Rel:        RelRegressionAssociatedChange,
		ToType:     EntityChange,
		ToID:       changeID,
		SourceType: sourceType,
		Confidence: confidence,
	})
	if err != nil {
		return err
	}
	if inserted {
		if err := s.appendLinked(EntityRegression, regressionID, RelRegressionAssociatedChange, EntityChange, changeID, LinkMeta{SourceType: sourceType, Confidence: confidence}); err != nil {
			return err
		}
	}
	return nil
}

// AssociateRegressionWithChange links a regression to a change (correlated default; never sets caused).
func (s *Service) AssociateRegressionWithChange(ctx context.Context, regressionID, changeID string) error {
	_ = ctx
	return s.associateRegressionWithChange(regressionID, changeID, DefaultSourceType, 0)
}

// ListRegressionsByChangeID returns regressions linked to the change via regression_associated_change.
func (s *Service) ListRegressionsByChangeID(ctx context.Context, changeID string) ([]store.Regression, error) {
	_ = ctx
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return nil, &ErrValidation{Msg: "change_id is required"}
	}
	if _, err := s.store.GetChange(changeID); err != nil {
		return nil, err
	}
	return s.store.ListRegressionsByChangeID(changeID)
}

// RegressionsForChange is an alias for ListRegressionsByChangeID.
func (s *Service) RegressionsForChange(ctx context.Context, changeID string) ([]store.Regression, error) {
	return s.ListRegressionsByChangeID(ctx, changeID)
}

// LinkHypothesisToRegression upgrades correlated → hypothesized. Never sets caused.
func (s *Service) LinkHypothesisToRegression(ctx context.Context, hypothesisID, regressionID string) (store.Regression, error) {
	_ = ctx
	hypothesisID = strings.TrimSpace(hypothesisID)
	regressionID = strings.TrimSpace(regressionID)
	if hypothesisID == "" || regressionID == "" {
		return store.Regression{}, &ErrValidation{Msg: "hypothesis id and regression id are required"}
	}
	h, err := s.store.GetHypothesis(hypothesisID)
	if err != nil {
		return store.Regression{}, err
	}
	if h.Status != store.HypothesisStatusOpen && h.Status != store.HypothesisStatusConfirmed {
		return store.Regression{}, &ErrValidation{Msg: "hypothesis status must be OPEN or CONFIRMED"}
	}
	row, err := s.store.GetRegression(regressionID)
	if err != nil {
		return store.Regression{}, err
	}
	if row.Attribution == store.RegressionAttributionCaused {
		return store.Regression{}, &ErrValidation{Msg: "caused attribution is terminal"}
	}

	old := row.Attribution
	if row.Attribution == store.RegressionAttributionCorrelated {
		row.Attribution = store.RegressionAttributionHypothesized
		row, err = s.store.UpsertRegression(row)
		if err != nil {
			return store.Regression{}, err
		}
		if err := s.appendNamed(EventRegressionAttributionChanged, EntityRegression, row.ID, map[string]string{
			"old": old,
			"new": row.Attribution,
		}); err != nil {
			return store.Regression{}, err
		}
	}

	inserted, _, err := s.store.InsertLinkOrIgnore(store.EntityLink{
		FromType:   EntityHypothesis,
		FromID:     h.ID,
		Rel:        RelHypothesisExplainsRegression,
		ToType:     EntityRegression,
		ToID:       row.ID,
		SourceType: DefaultSourceType,
	})
	if err != nil {
		return store.Regression{}, err
	}
	if inserted {
		if err := s.appendLinked(EntityHypothesis, h.ID, RelHypothesisExplainsRegression, EntityRegression, row.ID, LinkMeta{}.withDefaults()); err != nil {
			return store.Regression{}, err
		}
	}
	return row, nil
}

func (s *Service) linkedConfirmedHypothesis(regressionID string) (bool, error) {
	links, err := s.store.ListLinksTo(EntityRegression, regressionID)
	if err != nil {
		return false, err
	}
	for _, l := range links {
		if l.Rel != RelHypothesisExplainsRegression || l.FromType != EntityHypothesis {
			continue
		}
		h, err := s.store.GetHypothesis(l.FromID)
		if err != nil {
			return false, err
		}
		if h.Status == store.HypothesisStatusConfirmed {
			return true, nil
		}
	}
	return false, nil
}

// SetRegressionAttributionCaused applies the evidence policy. Never implicit.
func (s *Service) SetRegressionAttributionCaused(ctx context.Context, regressionID string, evidenceIDs []string) (store.Regression, error) {
	_ = ctx
	regressionID = strings.TrimSpace(regressionID)
	if regressionID == "" {
		return store.Regression{}, &ErrValidation{Msg: "regression id is required"}
	}
	if len(evidenceIDs) == 0 {
		return store.Regression{}, &ErrValidation{Msg: "evidence ids are required for attribution=caused"}
	}
	row, err := s.store.GetRegression(regressionID)
	if err != nil {
		return store.Regression{}, err
	}
	if row.Attribution == store.RegressionAttributionCaused {
		return store.Regression{}, &ErrValidation{Msg: "caused attribution is terminal"}
	}
	if row.Attribution != store.RegressionAttributionHypothesized {
		return store.Regression{}, &ErrValidation{Msg: "attribution=caused requires hypothesized first"}
	}
	confirmed, err := s.linkedConfirmedHypothesis(row.ID)
	if err != nil {
		return store.Regression{}, err
	}
	if !confirmed {
		return store.Regression{}, &ErrValidation{Msg: "attribution=caused requires a CONFIRMED hypothesis link"}
	}

	seen := map[string]struct{}{}
	for _, eid := range evidenceIDs {
		eid = strings.TrimSpace(eid)
		if eid == "" {
			return store.Regression{}, &ErrValidation{Msg: "evidence id is required"}
		}
		if _, dup := seen[eid]; dup {
			continue
		}
		seen[eid] = struct{}{}
		if _, err := s.store.GetEvidence(eid); err != nil {
			return store.Regression{}, err
		}
	}

	old := row.Attribution
	row.Attribution = store.RegressionAttributionCaused
	out, err := s.store.UpsertRegression(row)
	if err != nil {
		return store.Regression{}, err
	}
	for eid := range seen {
		if err := s.insertLink(EntityRegression, out.ID, RelRegressionSupportedBy, EntityEvidence, eid, out.SourceType, out.Confidence); err != nil {
			return store.Regression{}, err
		}
	}
	if err := s.appendNamed(EventRegressionAttributionChanged, EntityRegression, out.ID, map[string]string{
		"old": old,
		"new": out.Attribution,
	}); err != nil {
		return store.Regression{}, err
	}
	return out, nil
}

func (s *Service) transitionRegressionStatus(ctx context.Context, id, toStatus, reason string) (store.Regression, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	reason = strings.TrimSpace(reason)
	if id == "" {
		return store.Regression{}, &ErrValidation{Msg: "regression id is required"}
	}
	if reason == "" {
		return store.Regression{}, &ErrValidation{Msg: "reason is required"}
	}
	row, err := s.store.GetRegression(id)
	if err != nil {
		return store.Regression{}, err
	}
	if row.Status != store.RegressionStatusOpen {
		return store.Regression{}, &ErrValidation{Msg: "regression status transition only from OPEN"}
	}
	row.Status = toStatus
	out, err := s.store.UpsertRegression(row)
	if err != nil {
		return store.Regression{}, err
	}
	if err := s.appendNamed(EventRegressionResolved, EntityRegression, out.ID, map[string]string{
		"reason": reason,
		"status": out.Status,
	}); err != nil {
		return store.Regression{}, err
	}
	return out, nil
}

// ResolveRegression sets status RESOLVED. Reason is required. No delete.
func (s *Service) ResolveRegression(ctx context.Context, id, reason string) (store.Regression, error) {
	return s.transitionRegressionStatus(ctx, id, store.RegressionStatusResolved, reason)
}

// SupersedeRegression sets status SUPERSEDED. Reason is required. No delete.
func (s *Service) SupersedeRegression(ctx context.Context, id, reason string) (store.Regression, error) {
	return s.transitionRegressionStatus(ctx, id, store.RegressionStatusSuperseded, reason)
}

// HasOpenRegression reports any OPEN regression for the task (S06 / SelectNext input).
func (s *Service) HasOpenRegression(ctx context.Context, taskID string) (bool, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return false, &ErrValidation{Msg: "task_id is required"}
	}
	return s.store.HasOpenRegression(taskID)
}

// CountOpenRegressionsByTaskID counts OPEN regressions for the task.
func (s *Service) CountOpenRegressionsByTaskID(ctx context.Context, taskID string) (int, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return 0, &ErrValidation{Msg: "task_id is required"}
	}
	return s.store.CountOpenRegressionsByTaskID(taskID)
}

// ListOpenRegressions returns a bounded OPEN list for packets (max 32).
func (s *Service) ListOpenRegressions(ctx context.Context, taskID string) ([]store.Regression, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, &ErrValidation{Msg: "task_id is required"}
	}
	return s.store.ListOpenRegressions(taskID)
}

// GetRegression loads a regression by id.
func (s *Service) GetRegression(ctx context.Context, id string) (store.Regression, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	if id == "" {
		return store.Regression{}, &ErrValidation{Msg: "regression id is required"}
	}
	return s.store.GetRegression(id)
}

func normalizeDependencyRef(in DependencyRef) (DependencyRef, error) {
	kind := strings.TrimSpace(in.Kind)
	ref := strings.TrimSpace(in.Ref)
	if ref == "" {
		return DependencyRef{}, &ErrValidation{Msg: "dependency ref is required"}
	}
	switch kind {
	case depKindPath, depKindFile:
		ref = store.NormalizePath(ref)
		if ref == "" {
			return DependencyRef{}, &ErrValidation{Msg: "dependency ref is required"}
		}
	case depKindSymbol:
	default:
		return DependencyRef{}, &ErrValidation{Msg: "dependency kind must be path, symbol, or file"}
	}
	if err := failClosedMaxBytes("dependency ref", ref, maxDependencyRefBytes); err != nil {
		return DependencyRef{}, err
	}
	return DependencyRef{Kind: kind, Ref: ref}, nil
}

func marshalJSONArray(v any) (string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return "", &ErrValidation{Msg: "structured field must be a JSON array"}
	}
	var probe any
	if err := json.Unmarshal(raw, &probe); err != nil {
		return "", &ErrValidation{Msg: "structured field must be valid JSON"}
	}
	if _, ok := probe.([]any); !ok {
		return "", &ErrValidation{Msg: "structured field must be a JSON array"}
	}
	return string(raw), nil
}

// CreateReflection persists structured arrays. Essay-only (summary/note) fails closed.
func (s *Service) CreateReflection(ctx context.Context, in ReflectionInput) (store.Reflection, error) {
	_ = ctx
	task, err := s.requireTask(strings.TrimSpace(in.TaskID))
	if err != nil {
		return store.Reflection{}, err
	}
	if err := failClosedMaxBytes("summary", in.Summary, maxRegressionSummaryBytes); err != nil {
		return store.Reflection{}, err
	}
	if err := failClosedMaxBytes("broaden_tests_note", in.BroadenTestsNote, maxBroadenTestsNoteBytes); err != nil {
		return store.Reflection{}, err
	}
	if err := validateConfidence01(in.Confidence); err != nil {
		return store.Reflection{}, err
	}

	assumpIDs := make([]string, 0, len(in.InvalidatedAssumptionIDs))
	seenA := map[string]struct{}{}
	for _, id := range in.InvalidatedAssumptionIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			return store.Reflection{}, &ErrValidation{Msg: "assumption id is required"}
		}
		if _, dup := seenA[id]; dup {
			continue
		}
		seenA[id] = struct{}{}
		if _, err := s.store.GetAssumption(id); err != nil {
			return store.Reflection{}, err
		}
		assumpIDs = append(assumpIDs, id)
	}
	if len(assumpIDs) > maxReflectionArrayItems {
		return store.Reflection{}, &ErrValidation{Msg: "invalidated_assumptions_json exceeds 32 items"}
	}

	deps := make([]DependencyRef, 0, len(in.NewDependencies))
	for _, d := range in.NewDependencies {
		norm, err := normalizeDependencyRef(d)
		if err != nil {
			return store.Reflection{}, err
		}
		deps = append(deps, norm)
	}
	if len(deps) > maxReflectionArrayItems {
		return store.Reflection{}, &ErrValidation{Msg: "new_dependencies_json exceeds 32 items"}
	}

	tests := make([]string, 0, len(in.UsefulTests))
	for _, name := range in.UsefulTests {
		name = strings.TrimSpace(name)
		if name == "" {
			return store.Reflection{}, &ErrValidation{Msg: "test name is required"}
		}
		if err := failClosedMaxBytes("test name", name, maxUsefulTestNameBytes); err != nil {
			return store.Reflection{}, err
		}
		tests = append(tests, name)
	}
	if len(tests) > maxReflectionArrayItems {
		return store.Reflection{}, &ErrValidation{Msg: "useful_tests_json exceeds 32 items"}
	}

	if len(assumpIDs) == 0 && len(deps) == 0 && len(tests) == 0 {
		return store.Reflection{}, &ErrValidation{Msg: "reflection requires at least one structured array"}
	}

	invJSON, err := marshalJSONArray(assumpIDs)
	if err != nil {
		return store.Reflection{}, err
	}
	depJSON, err := marshalJSONArray(deps)
	if err != nil {
		return store.Reflection{}, err
	}
	testJSON, err := marshalJSONArray(tests)
	if err != nil {
		return store.Reflection{}, err
	}

	row, err := s.store.UpsertReflection(store.Reflection{
		ID:                         uuid.NewString(),
		TaskID:                     task.ID,
		Summary:                    in.Summary,
		InvalidatedAssumptionsJSON: invJSON,
		NewDependenciesJSON:        depJSON,
		UsefulTestsJSON:            testJSON,
		BroadenTestsNote:           in.BroadenTestsNote,
		Actor:                      strings.TrimSpace(in.Actor),
		SourceType:                 defaultSource(in.SourceType),
		Confidence:                 in.Confidence,
	})
	if err != nil {
		return store.Reflection{}, err
	}
	for _, aid := range assumpIDs {
		if err := s.insertLink(EntityReflection, row.ID, RelReflectionInvalidatesAssumption, EntityAssumption, aid, row.SourceType, row.Confidence); err != nil {
			return store.Reflection{}, err
		}
	}
	if err := s.appendCreated(EntityReflection, row.ID, row.Summary); err != nil {
		return store.Reflection{}, err
	}
	if err := s.appendNamed(EventReflectionRecorded, EntityReflection, row.ID, map[string]string{
		"task_id": row.TaskID,
	}); err != nil {
		return store.Reflection{}, err
	}
	return row, nil
}

// GetReflection loads a reflection by id.
func (s *Service) GetReflection(ctx context.Context, id string) (store.Reflection, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	if id == "" {
		return store.Reflection{}, &ErrValidation{Msg: "reflection id is required"}
	}
	return s.store.GetReflection(id)
}

// ListReflectionsByTaskID returns reflections for a task.
func (s *Service) ListReflectionsByTaskID(ctx context.Context, taskID string) ([]store.Reflection, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, &ErrValidation{Msg: "task_id is required"}
	}
	return s.store.ListReflectionsByTaskID(taskID)
}

func (s *Service) validateRelInput(in RelInput) error {
	if strings.TrimSpace(in.FromType) == "" || strings.TrimSpace(in.FromID) == "" ||
		strings.TrimSpace(in.ToType) == "" || strings.TrimSpace(in.ToID) == "" {
		return &ErrValidation{Msg: "from_type, from_id, to_type, and to_id are required"}
	}
	if err := validateConfidence01(in.Confidence); err != nil {
		return err
	}
	if err := s.requireKnownEntity(in.FromType, in.FromID); err != nil {
		return err
	}
	return s.requireKnownEntity(in.ToType, in.ToID)
}

// RecordObservedRelationship inserts observed_relationship with confidence. No evidence required.
func (s *Service) RecordObservedRelationship(ctx context.Context, in RelInput) (store.EntityLink, error) {
	_ = ctx
	if err := s.validateRelInput(in); err != nil {
		return store.EntityLink{}, err
	}
	src := defaultSource(in.SourceType)
	link, err := s.store.InsertLink(store.EntityLink{
		FromType:   strings.TrimSpace(in.FromType),
		FromID:     strings.TrimSpace(in.FromID),
		Rel:        RelObservedRelationship,
		ToType:     strings.TrimSpace(in.ToType),
		ToID:       strings.TrimSpace(in.ToID),
		SourceType: src,
		Confidence: in.Confidence,
	})
	if err != nil {
		return store.EntityLink{}, err
	}
	if err := s.appendLinked(link.FromType, link.FromID, RelObservedRelationship, link.ToType, link.ToID, LinkMeta{SourceType: src, Confidence: in.Confidence}); err != nil {
		return store.EntityLink{}, err
	}
	if err := s.appendNamed(EventRelationshipObserved, link.FromType, link.FromID, map[string]string{
		"rel":     RelObservedRelationship,
		"to_type": link.ToType,
		"to_id":   link.ToID,
	}); err != nil {
		return store.EntityLink{}, err
	}
	return link, nil
}

// RecordCausalRelationship inserts caused_by with required evidence. Does not set attribution.
func (s *Service) RecordCausalRelationship(ctx context.Context, in RelInput, evidenceIDs []string) (store.EntityLink, error) {
	_ = ctx
	if err := s.validateRelInput(in); err != nil {
		return store.EntityLink{}, err
	}
	if len(evidenceIDs) == 0 {
		return store.EntityLink{}, &ErrValidation{Msg: "evidence ids are required for caused_by"}
	}
	src := defaultSource(in.SourceType)
	fromType := strings.TrimSpace(in.FromType)
	fromID := strings.TrimSpace(in.FromID)

	seen := map[string]struct{}{}
	for _, eid := range evidenceIDs {
		eid = strings.TrimSpace(eid)
		if eid == "" {
			return store.EntityLink{}, &ErrValidation{Msg: "evidence id is required"}
		}
		if _, dup := seen[eid]; dup {
			continue
		}
		seen[eid] = struct{}{}
		if _, err := s.store.GetEvidence(eid); err != nil {
			return store.EntityLink{}, err
		}
	}

	link, err := s.store.InsertLink(store.EntityLink{
		FromType:   fromType,
		FromID:     fromID,
		Rel:        RelCausedBy,
		ToType:     strings.TrimSpace(in.ToType),
		ToID:       strings.TrimSpace(in.ToID),
		SourceType: src,
		Confidence: in.Confidence,
	})
	if err != nil {
		return store.EntityLink{}, err
	}
	if err := s.appendLinked(link.FromType, link.FromID, RelCausedBy, link.ToType, link.ToID, LinkMeta{SourceType: src, Confidence: in.Confidence}); err != nil {
		return store.EntityLink{}, err
	}
	for eid := range seen {
		if err := s.insertLink(fromType, fromID, RelRelationshipSupportedBy, EntityEvidence, eid, src, in.Confidence); err != nil {
			return store.EntityLink{}, err
		}
	}
	if err := s.appendNamed(EventRelationshipCaused, link.FromType, link.FromID, map[string]string{
		"rel":     RelCausedBy,
		"to_type": link.ToType,
		"to_id":   link.ToID,
	}); err != nil {
		return store.EntityLink{}, err
	}
	return link, nil
}
