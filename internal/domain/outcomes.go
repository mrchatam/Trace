package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	maxOutcomeSummaryBytes = 4096
	maxScoresJSONBytes     = 16384
	maxBaselineLabelBytes  = 256
)

// Experiments (TRACE_THOUGHTPROCESS §16) and risk-adaptive verification (§18)
// remain Future — this package records outcomes; it does not embed a test runner.

// BaselineInput creates a thin baseline (git OID + scores JSON only).
type BaselineInput struct {
	GitCommit  string
	ScoresJSON string
	Label      string
	SourceType string
}

// TestOutcomeInput records kind=test. A stored row is required for CheckTestGate.
type TestOutcomeInput struct {
	TaskID      string
	TestName    string
	TestStatus  string
	Summary     string
	EvidenceIDs []string
	Actor       string
	SourceType  string
}

// VerificationOutcomeInput records kind=verification. Requires goal + evidence.
type VerificationOutcomeInput struct {
	TaskID             string
	GoalID             string
	VerificationStatus string
	EvidenceIDs        []string
	Summary            string
	Actor              string
	SourceType         string
}

// EvaluationOutcomeInput records kind=evaluation. comparison_json is library-computed.
type EvaluationOutcomeInput struct {
	TaskID     string
	BaselineID string
	ScoresJSON string
	Actor      string
	SourceType string
}

// ScoreDimensionComparison is one dimension in comparison_json.
type ScoreDimensionComparison struct {
	Baseline   any      `json:"baseline"`
	Current    any      `json:"current"`
	Delta      *float64 `json:"delta"`
	Regression bool     `json:"regression"`
	Note       string   `json:"note,omitempty"`
}

// ScoreComparison is the structured evaluation result (not a boolean PASS).
type ScoreComparison struct {
	BaselineID        string                              `json:"baseline_id"`
	GitCommit         string                              `json:"git_commit"`
	Dimensions        map[string]ScoreDimensionComparison `json:"dimensions"`
	OverallRegression bool                                `json:"overall_regression"`
}

func normalizeTestStatus(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	switch s {
	case store.TestStatusPass, store.TestStatusFail, store.TestStatusSkip, store.TestStatusError:
		return s, nil
	default:
		return "", &ErrValidation{Msg: "test_status must be pass, fail, skip, or error"}
	}
}

func normalizeVerificationStatus(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	switch s {
	case store.VerificationStatusVerified, store.VerificationStatusFailed, store.VerificationStatusPartial:
		return s, nil
	default:
		return "", &ErrValidation{Msg: "verification_status must be verified, failed, or partial"}
	}
}

func parseScoresObject(raw, label string) (map[string]json.RawMessage, error) {
	s := strings.TrimSpace(raw)
	if err := failClosedMaxBytes(label, s, maxScoresJSONBytes); err != nil {
		return nil, err
	}
	if s == "" {
		return nil, &ErrValidation{Msg: label + " is required"}
	}
	var root any
	if err := json.Unmarshal([]byte(s), &root); err != nil {
		return nil, &ErrValidation{Msg: label + " must be valid JSON"}
	}
	obj, ok := root.(map[string]any)
	if !ok {
		return nil, &ErrValidation{Msg: label + " must be a JSON object"}
	}
	if len(obj) == 0 {
		return nil, &ErrValidation{Msg: label + " must be a non-empty JSON object"}
	}
	rawMap := map[string]json.RawMessage{}
	if err := json.Unmarshal([]byte(s), &rawMap); err != nil {
		return nil, &ErrValidation{Msg: label + " must be a JSON object"}
	}
	return rawMap, nil
}

func compactJSONObject(raw string) (string, error) {
	var obj map[string]any
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return "", &ErrValidation{Msg: "invalid JSON object"}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return "", &ErrValidation{Msg: "invalid JSON object"}
	}
	return string(out), nil
}

func scoreAsNumber(raw json.RawMessage) (float64, bool) {
	var n float64
	if err := json.Unmarshal(raw, &n); err != nil {
		return 0, false
	}
	return n, true
}

func scoreAsValue(raw json.RawMessage) any {
	if n, ok := scoreAsNumber(raw); ok {
		return n
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var v any
	_ = json.Unmarshal(raw, &v)
	return v
}

func lowerIsBetterDimension(key string) bool {
	k := strings.ToLower(key)
	return strings.Contains(k, "_ms") || strings.Contains(k, "latency")
}

// CompareScoresToBaseline returns structured per-dimension deltas. Not a boolean PASS.
func CompareScoresToBaseline(currentJSON string, baseline store.Baseline) (ScoreComparison, error) {
	current, err := parseScoresObject(currentJSON, "scores_json")
	if err != nil {
		return ScoreComparison{}, err
	}
	base, err := parseScoresObject(baseline.ScoresJSON, "baseline scores_json")
	if err != nil {
		return ScoreComparison{}, err
	}
	dims := make(map[string]ScoreDimensionComparison, len(current)+len(base))
	overall := false
	seen := map[string]struct{}{}
	for key, curRaw := range current {
		seen[key] = struct{}{}
		dim := ScoreDimensionComparison{Current: scoreAsValue(curRaw)}
		baseRaw, ok := base[key]
		if !ok {
			dim.Baseline = nil
			dim.Delta = nil
			dim.Regression = false
			dim.Note = "missing baseline dimension"
			dims[key] = dim
			continue
		}
		dim.Baseline = scoreAsValue(baseRaw)
		curN, curOK := scoreAsNumber(curRaw)
		baseN, baseOK := scoreAsNumber(baseRaw)
		if !curOK || !baseOK {
			dim.Delta = nil
			dim.Regression = false
			dim.Note = "non-numeric"
			dims[key] = dim
			continue
		}
		delta := curN - baseN
		dim.Delta = &delta
		if lowerIsBetterDimension(key) {
			dim.Regression = curN > baseN
		} else {
			dim.Regression = curN < baseN
		}
		if dim.Regression {
			overall = true
		}
		dims[key] = dim
	}
	for key, baseRaw := range base {
		if _, ok := seen[key]; ok {
			continue
		}
		dims[key] = ScoreDimensionComparison{
			Baseline:   scoreAsValue(baseRaw),
			Current:    nil,
			Delta:      nil,
			Regression: false,
			Note:       "missing current dimension",
		}
	}
	return ScoreComparison{
		BaselineID:        baseline.ID,
		GitCommit:         baseline.GitCommit,
		Dimensions:        dims,
		OverallRegression: overall,
	}, nil
}

func rejectBooleanPassComparison(cmp ScoreComparison) error {
	raw, err := json.Marshal(cmp)
	if err != nil {
		return err
	}
	var probe map[string]any
	if err := json.Unmarshal(raw, &probe); err != nil {
		return err
	}
	if pass, ok := probe["pass"].(bool); ok {
		return &ErrValidation{Msg: "comparison_json must not be a boolean pass; got pass=" + boolString(pass)}
	}
	return nil
}

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func (s *Service) requireTask(taskID string) (store.Task, error) {
	if taskID == "" {
		return store.Task{}, &ErrValidation{Msg: "task_id is required"}
	}
	return s.store.GetTask(taskID)
}

func (s *Service) linkOutcomeEvidence(outcomeID string, evidenceIDs []string) error {
	seen := map[string]struct{}{}
	for _, eid := range evidenceIDs {
		eid = strings.TrimSpace(eid)
		if eid == "" {
			return &ErrValidation{Msg: "evidence id is required"}
		}
		if _, dup := seen[eid]; dup {
			continue
		}
		seen[eid] = struct{}{}
		if _, err := s.store.GetEvidence(eid); err != nil {
			return err
		}
		if err := s.insertTypedLink(EntityOutcomeResult, outcomeID, RelOutcomeSupportedBy, EntityEvidence, eid); err != nil {
			return err
		}
	}
	return nil
}

func defaultSource(src string) string {
	if strings.TrimSpace(src) == "" {
		return DefaultSourceType
	}
	return strings.TrimSpace(src)
}

// CreateBaseline persists git_commit + scores_json. No DeleteBaseline (Law 11).
func (s *Service) CreateBaseline(ctx context.Context, in BaselineInput) (store.Baseline, error) {
	_ = ctx
	sha, err := normalizeGitCommit(in.GitCommit)
	if err != nil {
		return store.Baseline{}, err
	}
	if sha == "" {
		return store.Baseline{}, &ErrValidation{Msg: "git_commit is required"}
	}
	if _, err := parseScoresObject(in.ScoresJSON, "scores_json"); err != nil {
		return store.Baseline{}, err
	}
	compact, err := compactJSONObject(in.ScoresJSON)
	if err != nil {
		return store.Baseline{}, err
	}
	if err := failClosedMaxBytes("label", in.Label, maxBaselineLabelBytes); err != nil {
		return store.Baseline{}, err
	}
	b, err := s.store.UpsertBaseline(store.Baseline{
		ID:           uuid.NewString(),
		GitCommit:    sha,
		ScoresJSON:   compact,
		Label:        in.Label,
		SourceType:   defaultSource(in.SourceType),
		Status:       store.BaselineStatusActive,
		SupersedesID: "",
	})
	if err != nil {
		return store.Baseline{}, err
	}
	if err := s.appendNamed(EventBaselineCreated, EntityBaseline, b.ID, map[string]string{
		"git_commit": b.GitCommit,
	}); err != nil {
		return store.Baseline{}, err
	}
	return b, nil
}

// GetBaseline loads a baseline by id.
func (s *Service) GetBaseline(ctx context.Context, id string) (store.Baseline, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	if id == "" {
		return store.Baseline{}, &ErrValidation{Msg: "baseline id is required"}
	}
	return s.store.GetBaseline(id)
}

// SupersedeBaseline marks a baseline superseded (Law 11 — no delete).
func (s *Service) SupersedeBaseline(ctx context.Context, baselineID string) error {
	_ = ctx
	baselineID = strings.TrimSpace(baselineID)
	if baselineID == "" {
		return &ErrValidation{Msg: "baseline id is required"}
	}
	if _, err := s.store.GetBaseline(baselineID); err != nil {
		return err
	}
	return s.store.SetBaselinePromotion(baselineID, store.BaselineStatusSuperseded, "")
}

// PromoteBaseline marks baselineID active and supersedes prior active baseline for same git_commit+label.
func (s *Service) PromoteBaseline(ctx context.Context, baselineID string) (store.Baseline, error) {
	_ = ctx
	baselineID = strings.TrimSpace(baselineID)
	if baselineID == "" {
		return store.Baseline{}, &ErrValidation{Msg: "baseline id is required"}
	}
	target, err := s.store.GetBaseline(baselineID)
	if err != nil {
		return store.Baseline{}, err
	}
	prior, priorErr := s.store.GetActiveBaselineByCommitLabelExcluding(target.GitCommit, target.Label, target.ID)
	if priorErr != nil {
		if !errors.Is(priorErr, sql.ErrNoRows) {
			return store.Baseline{}, priorErr
		}
		if target.Status == store.BaselineStatusActive {
			return target, nil
		}
		if err := s.store.SetBaselinePromotion(target.ID, store.BaselineStatusActive, ""); err != nil {
			return store.Baseline{}, err
		}
		return s.store.GetBaseline(target.ID)
	}
	if target.Status == store.BaselineStatusActive && target.SupersedesID == prior.ID {
		return target, nil
	}
	if err := s.SupersedeBaseline(ctx, prior.ID); err != nil {
		return store.Baseline{}, err
	}
	if err := s.store.SetBaselinePromotion(target.ID, store.BaselineStatusActive, prior.ID); err != nil {
		return store.Baseline{}, err
	}
	promoted, err := s.store.GetBaseline(target.ID)
	if err != nil {
		return store.Baseline{}, err
	}
	if err := s.appendNamed(EventBaselinePromoted, EntityBaseline, promoted.ID, map[string]string{
		"git_commit":    promoted.GitCommit,
		"label":         promoted.Label,
		"supersedes_id": prior.ID,
	}); err != nil {
		return store.Baseline{}, err
	}
	return promoted, nil
}

func latestComputedEvaluation(rows []store.OutcomeResult) (store.OutcomeResult, bool) {
	var latest store.OutcomeResult
	var found bool
	for _, o := range rows {
		if !comparisonComputed(o.ComparisonJSON) {
			continue
		}
		if !found || o.CreatedAt > latest.CreatedAt || (o.CreatedAt == latest.CreatedAt && o.ID > latest.ID) {
			latest = o
			found = true
		}
	}
	return latest, found
}

func latestRecordedOrComparedChange(changes []store.Change) (store.Change, bool) {
	var latest store.Change
	var found bool
	for _, c := range changes {
		if c.Status != store.ChangeStatusRecorded && c.Status != store.ChangeStatusCompared {
			continue
		}
		if !found || c.CreatedAt > latest.CreatedAt || (c.CreatedAt == latest.CreatedAt && c.ID > latest.ID) {
			latest = c
			found = true
		}
	}
	return latest, found
}

func comparisonOverallRegression(raw string) (bool, error) {
	var cmp ScoreComparison
	if err := json.Unmarshal([]byte(raw), &cmp); err != nil {
		return false, err
	}
	return cmp.OverallRegression, nil
}

// CheckPromotionGate is advisory — does not mutate work_state or auto-DONE.
func (s *Service) CheckPromotionGate(ctx context.Context, taskID, baselineID string) (allowed bool, reason string, err error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	baselineID = strings.TrimSpace(baselineID)
	if taskID == "" {
		return false, "task_id is required", &ErrValidation{Msg: "task_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return false, "task not found", err
	}
	if baselineID != "" {
		if _, err := s.store.GetBaseline(baselineID); err != nil {
			return false, "baseline_not_found", err
		}
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindEvaluation)
	if err != nil {
		return false, "lookup failed", err
	}
	latest, ok := latestComputedEvaluation(rows)
	if !ok {
		return false, "no_stored_evaluation", nil
	}
	regressed, err := comparisonOverallRegression(latest.ComparisonJSON)
	if err != nil {
		return false, "lookup failed", err
	}
	if regressed {
		return false, "eval_regression", nil
	}
	return true, "", nil
}

// RecordTestOutcome persists kind=test. Agent claim without a row cannot pass CheckTestGate.
func (s *Service) RecordTestOutcome(ctx context.Context, in TestOutcomeInput) (store.OutcomeResult, error) {
	_ = ctx
	task, err := s.requireTask(strings.TrimSpace(in.TaskID))
	if err != nil {
		return store.OutcomeResult{}, err
	}
	name := strings.TrimSpace(in.TestName)
	if name == "" {
		return store.OutcomeResult{}, &ErrValidation{Msg: "test_name is required"}
	}
	status, err := normalizeTestStatus(in.TestStatus)
	if err != nil {
		return store.OutcomeResult{}, err
	}
	if err := failClosedMaxBytes("summary", in.Summary, maxOutcomeSummaryBytes); err != nil {
		return store.OutcomeResult{}, err
	}
	for _, eid := range in.EvidenceIDs {
		if strings.TrimSpace(eid) == "" {
			return store.OutcomeResult{}, &ErrValidation{Msg: "evidence id is required"}
		}
		if _, err := s.store.GetEvidence(strings.TrimSpace(eid)); err != nil {
			return store.OutcomeResult{}, err
		}
	}
	o, err := s.store.UpsertOutcomeResult(store.OutcomeResult{
		ID:             uuid.NewString(),
		TaskID:         task.ID,
		Kind:           store.OutcomeKindTest,
		TestName:       name,
		TestStatus:     status,
		ScoresJSON:     "{}",
		ComparisonJSON: "{}",
		Summary:        in.Summary,
		Actor:          strings.TrimSpace(in.Actor),
		SourceType:     defaultSource(in.SourceType),
	})
	if err != nil {
		return store.OutcomeResult{}, err
	}
	if err := s.linkOutcomeEvidence(o.ID, in.EvidenceIDs); err != nil {
		return store.OutcomeResult{}, err
	}
	if err := s.appendNamed(EventOutcomeRecorded, EntityOutcomeResult, o.ID, map[string]string{
		"kind":        o.Kind,
		"test_name":   o.TestName,
		"test_status": o.TestStatus,
		"task_id":     o.TaskID,
	}); err != nil {
		return store.OutcomeResult{}, err
	}
	return o, nil
}

// RecordVerificationOutcome persists kind=verification with goal_id + evidence links.
func (s *Service) RecordVerificationOutcome(ctx context.Context, in VerificationOutcomeInput) (store.OutcomeResult, error) {
	_ = ctx
	task, err := s.requireTask(strings.TrimSpace(in.TaskID))
	if err != nil {
		return store.OutcomeResult{}, err
	}
	goalID := strings.TrimSpace(in.GoalID)
	if goalID == "" {
		return store.OutcomeResult{}, &ErrValidation{Msg: "goal_id is required"}
	}
	if _, err := s.store.GetGoal(goalID); err != nil {
		return store.OutcomeResult{}, err
	}
	if task.GoalID != nil && strings.TrimSpace(*task.GoalID) != "" && *task.GoalID != goalID {
		return store.OutcomeResult{}, &ErrValidation{Msg: "goal_id must match task.goal_id"}
	}
	status, err := normalizeVerificationStatus(in.VerificationStatus)
	if err != nil {
		return store.OutcomeResult{}, err
	}
	if len(in.EvidenceIDs) < 1 {
		return store.OutcomeResult{}, &ErrValidation{Msg: "at least one evidence id is required"}
	}
	if err := failClosedMaxBytes("summary", in.Summary, maxOutcomeSummaryBytes); err != nil {
		return store.OutcomeResult{}, err
	}
	for _, eid := range in.EvidenceIDs {
		eid = strings.TrimSpace(eid)
		if eid == "" {
			return store.OutcomeResult{}, &ErrValidation{Msg: "evidence id is required"}
		}
		if _, err := s.store.GetEvidence(eid); err != nil {
			return store.OutcomeResult{}, err
		}
	}
	o, err := s.store.UpsertOutcomeResult(store.OutcomeResult{
		ID:                 uuid.NewString(),
		TaskID:             task.ID,
		Kind:               store.OutcomeKindVerification,
		GoalID:             goalID,
		VerificationStatus: status,
		ScoresJSON:         "{}",
		ComparisonJSON:     "{}",
		Summary:            in.Summary,
		Actor:              strings.TrimSpace(in.Actor),
		SourceType:         defaultSource(in.SourceType),
	})
	if err != nil {
		return store.OutcomeResult{}, err
	}
	if err := s.linkOutcomeEvidence(o.ID, in.EvidenceIDs); err != nil {
		return store.OutcomeResult{}, err
	}
	if err := s.appendNamed(EventOutcomeRecorded, EntityOutcomeResult, o.ID, map[string]string{
		"kind":                o.Kind,
		"goal_id":             o.GoalID,
		"verification_status": o.VerificationStatus,
		"task_id":             o.TaskID,
	}); err != nil {
		return store.OutcomeResult{}, err
	}
	return o, nil
}

// RecordEvaluationOutcome persists kind=evaluation and library-computed comparison_json.
func (s *Service) RecordEvaluationOutcome(ctx context.Context, in EvaluationOutcomeInput) (store.OutcomeResult, error) {
	_ = ctx
	task, err := s.requireTask(strings.TrimSpace(in.TaskID))
	if err != nil {
		return store.OutcomeResult{}, err
	}
	baselineID := strings.TrimSpace(in.BaselineID)
	if baselineID == "" {
		return store.OutcomeResult{}, &ErrValidation{Msg: "baseline_id is required"}
	}
	baseline, err := s.store.GetBaseline(baselineID)
	if err != nil {
		return store.OutcomeResult{}, err
	}
	if _, err := parseScoresObject(in.ScoresJSON, "scores_json"); err != nil {
		return store.OutcomeResult{}, err
	}
	compact, err := compactJSONObject(in.ScoresJSON)
	if err != nil {
		return store.OutcomeResult{}, err
	}
	cmp, err := CompareScoresToBaseline(compact, baseline)
	if err != nil {
		return store.OutcomeResult{}, err
	}
	if err := rejectBooleanPassComparison(cmp); err != nil {
		return store.OutcomeResult{}, err
	}
	cmpJSON, err := json.Marshal(cmp)
	if err != nil {
		return store.OutcomeResult{}, err
	}
	if err := failClosedMaxBytes("comparison_json", string(cmpJSON), maxScoresJSONBytes); err != nil {
		return store.OutcomeResult{}, err
	}
	o, err := s.store.UpsertOutcomeResult(store.OutcomeResult{
		ID:             uuid.NewString(),
		TaskID:         task.ID,
		Kind:           store.OutcomeKindEvaluation,
		BaselineID:     baseline.ID,
		ScoresJSON:     compact,
		ComparisonJSON: string(cmpJSON),
		Actor:          strings.TrimSpace(in.Actor),
		SourceType:     defaultSource(in.SourceType),
	})
	if err != nil {
		return store.OutcomeResult{}, err
	}
	if err := s.appendNamed(EventOutcomeRecorded, EntityOutcomeResult, o.ID, map[string]string{
		"kind":        o.Kind,
		"baseline_id": o.BaselineID,
		"task_id":     o.TaskID,
	}); err != nil {
		return store.OutcomeResult{}, err
	}
	if err := s.appendNamed(EventEvaluationCompared, EntityOutcomeResult, o.ID, map[string]string{
		"baseline_id": o.BaselineID,
		"git_commit":  baseline.GitCommit,
	}); err != nil {
		return store.OutcomeResult{}, err
	}
	return o, nil
}

// GetOutcomeResult loads an outcome by id.
func (s *Service) GetOutcomeResult(ctx context.Context, id string) (store.OutcomeResult, error) {
	_ = ctx
	id = strings.TrimSpace(id)
	if id == "" {
		return store.OutcomeResult{}, &ErrValidation{Msg: "outcome id is required"}
	}
	return s.store.GetOutcomeResult(id)
}

// ListEvalOutcomeResults returns test, verification, and evaluation outcomes for a task,
// newest first (created_at DESC, id DESC tiebreak). Default limit 32, cap 64.
func (s *Service) ListEvalOutcomeResults(ctx context.Context, taskID string, limit int) ([]store.OutcomeResult, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, &ErrValidation{Msg: "task_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return nil, err
	}
	limit = clampEvidenceLimit(limit)
	kinds := []string{store.OutcomeKindTest, store.OutcomeKindVerification, store.OutcomeKindEvaluation}
	var rows []store.OutcomeResult
	for _, kind := range kinds {
		part, err := s.store.ListOutcomeResultsByTaskKind(taskID, kind)
		if err != nil {
			return nil, err
		}
		rows = append(rows, part...)
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].CreatedAt != rows[j].CreatedAt {
			return rows[i].CreatedAt > rows[j].CreatedAt
		}
		return rows[i].ID > rows[j].ID
	})
	if len(rows) > limit {
		rows = rows[:limit]
	}
	return rows, nil
}

func comparisonComputed(raw string) bool {
	s := strings.TrimSpace(raw)
	if s == "" || s == "{}" {
		return false
	}
	var probe map[string]any
	if err := json.Unmarshal([]byte(s), &probe); err != nil {
		return false
	}
	if _, ok := probe["pass"].(bool); ok {
		return false
	}
	dims, _ := probe["dimensions"].(map[string]any)
	return len(dims) > 0
}

// CheckTestGate is true only when a stored kind=test row exists with matching name and pass.
func (s *Service) CheckTestGate(ctx context.Context, taskID, testName string) (bool, string, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	testName = strings.TrimSpace(testName)
	if taskID == "" {
		return false, "task_id is required", &ErrValidation{Msg: "task_id is required"}
	}
	if testName == "" {
		return false, "test_name is required", &ErrValidation{Msg: "test_name is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return false, "task not found", err
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindTest)
	if err != nil {
		return false, "lookup failed", err
	}
	for _, o := range rows {
		if o.TestName == testName && o.TestStatus == store.TestStatusPass {
			return true, "", nil
		}
	}
	return false, "no passing stored test result", nil
}

// CheckVerificationGate is true when a verified outcome exists with goal_id and evidence.
func (s *Service) CheckVerificationGate(ctx context.Context, taskID string) (bool, string, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return false, "task_id is required", &ErrValidation{Msg: "task_id is required"}
	}
	task, err := s.store.GetTask(taskID)
	if err != nil {
		return false, "task not found", err
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindVerification)
	if err != nil {
		return false, "lookup failed", err
	}
	for _, o := range rows {
		if o.VerificationStatus != store.VerificationStatusVerified || strings.TrimSpace(o.GoalID) == "" {
			continue
		}
		if task.GoalID != nil && strings.TrimSpace(*task.GoalID) != "" && o.GoalID != *task.GoalID {
			continue
		}
		n, err := s.store.CountOutcomeSupportedByEvidence(o.ID)
		if err != nil {
			return false, "lookup failed", err
		}
		if n >= 1 {
			return true, "", nil
		}
	}
	return false, "no verified outcome with goal_id and evidence", nil
}

// CheckEvaluationGate is true when an evaluation for baselineID has computed comparison_json.
func (s *Service) CheckEvaluationGate(ctx context.Context, taskID, baselineID string) (bool, string, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	baselineID = strings.TrimSpace(baselineID)
	if taskID == "" {
		return false, "task_id is required", &ErrValidation{Msg: "task_id is required"}
	}
	if baselineID == "" {
		return false, "baseline_id is required", &ErrValidation{Msg: "baseline_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return false, "task not found", err
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindEvaluation)
	if err != nil {
		return false, "lookup failed", err
	}
	for _, o := range rows {
		if o.BaselineID == baselineID && comparisonComputed(o.ComparisonJSON) {
			return true, "", nil
		}
	}
	return false, "no evaluation comparison for baseline", nil
}

// HasImplementationSignal reports a RECORDED or COMPARED change for the task.
func (s *Service) HasImplementationSignal(ctx context.Context, taskID string) (bool, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return false, &ErrValidation{Msg: "task_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return false, err
	}
	return s.store.HasImplementationSignal(taskID)
}

// HasVerificationDebt reports implementation-without-satisfactory-verification.
func (s *Service) HasVerificationDebt(ctx context.Context, taskID string) (bool, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return false, &ErrValidation{Msg: "task_id is required"}
	}
	return s.store.HasVerificationDebt(taskID)
}

// HasTestOutcomeSinceLatestChange is true when a kind=test outcome exists at or
// after the latest RECORDED/COMPARED change. No such change → false.
func (s *Service) HasTestOutcomeSinceLatestChange(ctx context.Context, taskID string) (bool, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return false, &ErrValidation{Msg: "task_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return false, err
	}
	changes, err := s.store.ListChangesByTaskID(taskID)
	if err != nil {
		return false, err
	}
	latest, ok := latestRecordedOrComparedChange(changes)
	if !ok {
		return false, nil
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindTest)
	if err != nil {
		return false, err
	}
	for _, o := range rows {
		if o.CreatedAt >= latest.CreatedAt {
			return true, nil
		}
	}
	return false, nil
}

// HasComputedEvaluation is true when a kind=evaluation row has computed comparison_json.
func (s *Service) HasComputedEvaluation(ctx context.Context, taskID string) (bool, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return false, &ErrValidation{Msg: "task_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return false, err
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindEvaluation)
	if err != nil {
		return false, err
	}
	_, ok := latestComputedEvaluation(rows)
	return ok, nil
}

// HasVerificationOutcome reports any kind=verification row for the task.
func (s *Service) HasVerificationOutcome(ctx context.Context, taskID string) (bool, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return false, &ErrValidation{Msg: "task_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return false, err
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindVerification)
	if err != nil {
		return false, err
	}
	return len(rows) > 0, nil
}

// HasReflectionSinceEvaluation is true when a reflection exists at or after the
// latest computed evaluation. No computed evaluation → false.
func (s *Service) HasReflectionSinceEvaluation(ctx context.Context, taskID string) (bool, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return false, &ErrValidation{Msg: "task_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return false, err
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, store.OutcomeKindEvaluation)
	if err != nil {
		return false, err
	}
	eval, ok := latestComputedEvaluation(rows)
	if !ok {
		return false, nil
	}
	refs, err := s.store.ListReflectionsByTaskID(taskID)
	if err != nil {
		return false, err
	}
	for _, r := range refs {
		if r.CreatedAt >= eval.CreatedAt {
			return true, nil
		}
	}
	return false, nil
}

// ListVerificationDebtSummary returns bounded debt items for S06 packets.
func (s *Service) ListVerificationDebtSummary(ctx context.Context, taskID string) ([]store.DebtItem, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, &ErrValidation{Msg: "task_id is required"}
	}
	return s.store.ListVerificationDebtSummary(taskID)
}

// CountTasksWithVerificationDebt counts tasks currently in verification debt.
func (s *Service) CountTasksWithVerificationDebt(ctx context.Context) (int, error) {
	_ = ctx
	return s.store.CountTasksWithVerificationDebt()
}

// IterationSnapshot is one stored outcome in an iteration compare.
type IterationSnapshot struct {
	ID         string         `json:"id"`
	Kind       string         `json:"kind"`
	TestName   string         `json:"test_name,omitempty"`
	TestStatus string         `json:"test_status,omitempty"`
	Scores     map[string]any `json:"scores,omitempty"`
	CreatedAt  string         `json:"created_at"`
}

// TestStatusDelta is a from→to test_status change between iterations.
type TestStatusDelta struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// IterationDelta is the compare payload (test_status and/or score dimensions).
type IterationDelta struct {
	TestStatus *TestStatusDelta                    `json:"test_status,omitempty"`
	Dimensions map[string]ScoreDimensionComparison `json:"dimensions,omitempty"`
}

// IterationCompareResult is {previous, current, delta} JSON for CompareIterationOutcomes.
type IterationCompareResult struct {
	Previous IterationSnapshot `json:"previous"`
	Current  IterationSnapshot `json:"current"`
	Delta    IterationDelta    `json:"delta"`
}

func normalizeCompareKind(raw string) (string, error) {
	k := strings.TrimSpace(raw)
	switch k {
	case store.OutcomeKindTest, store.OutcomeKindEvaluation:
		return k, nil
	default:
		return "", &ErrValidation{Msg: "kind must be test or evaluation"}
	}
}

func snapshotOutcome(o store.OutcomeResult) IterationSnapshot {
	snap := IterationSnapshot{
		ID:         o.ID,
		Kind:       o.Kind,
		TestName:   o.TestName,
		TestStatus: o.TestStatus,
		CreatedAt:  o.CreatedAt,
	}
	if strings.TrimSpace(o.ScoresJSON) != "" && o.ScoresJSON != "{}" {
		var scores map[string]any
		if err := json.Unmarshal([]byte(o.ScoresJSON), &scores); err == nil && len(scores) > 0 {
			snap.Scores = scores
		}
	}
	return snap
}

// CompareIterationOutcomes compares the last two stored outcomes of kind
// (test or evaluation) by created_at. Pure on stored rows — no test re-exec.
func (s *Service) CompareIterationOutcomes(ctx context.Context, taskID, kind string) (IterationCompareResult, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return IterationCompareResult{}, &ErrValidation{Msg: "task_id is required"}
	}
	kind, err := normalizeCompareKind(kind)
	if err != nil {
		return IterationCompareResult{}, err
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return IterationCompareResult{}, err
	}
	rows, err := s.store.ListOutcomeResultsByTaskKind(taskID, kind)
	if err != nil {
		return IterationCompareResult{}, err
	}
	if len(rows) < 2 {
		return IterationCompareResult{}, &ErrValidation{Msg: "need at least two outcomes of this kind"}
	}
	prev := rows[len(rows)-2]
	curr := rows[len(rows)-1]
	out := IterationCompareResult{
		Previous: snapshotOutcome(prev),
		Current:  snapshotOutcome(curr),
	}
	switch kind {
	case store.OutcomeKindTest:
		out.Delta.TestStatus = &TestStatusDelta{From: prev.TestStatus, To: curr.TestStatus}
	case store.OutcomeKindEvaluation:
		cmp, err := CompareScoresToBaseline(curr.ScoresJSON, store.Baseline{ID: prev.ID, ScoresJSON: prev.ScoresJSON})
		if err != nil {
			return IterationCompareResult{}, err
		}
		out.Delta.Dimensions = cmp.Dimensions
	}
	return out, nil
}
