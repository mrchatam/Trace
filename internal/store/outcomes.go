package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Outcome result kinds (single table discriminator).
const (
	OutcomeKindTest         = "test"
	OutcomeKindVerification = "verification"
	OutcomeKindEvaluation   = "evaluation"
)

// Test status vocabulary. Empty only for non-test kinds.
const (
	TestStatusNone  = ""
	TestStatusPass  = "pass"
	TestStatusFail  = "fail"
	TestStatusSkip  = "skip"
	TestStatusError = "error"
)

// Verification status vocabulary. Empty only for non-verification kinds.
const (
	VerificationStatusNone     = ""
	VerificationStatusVerified = "verified"
	VerificationStatusFailed   = "failed"
	VerificationStatusPartial  = "partial"
)

const emptyJSONObject = "{}"

// Baseline status vocabulary.
const (
	BaselineStatusActive     = "active"
	BaselineStatusSuperseded = "superseded"
)

// Baseline is a thin git-OID + scores JSON pointer. No runner output blobs.
type Baseline struct {
	ID           string
	GitCommit    string
	ScoresJSON   string
	Label        string
	SourceType   string
	Status       string
	SupersedesID string
	CreatedAt    string
	UpdatedAt    string
}

// OutcomeResult is one recorded test, verification, or evaluation.
type OutcomeResult struct {
	ID                 string
	TaskID             string
	Kind               string
	TestName           string
	TestStatus         string
	GoalID             string
	VerificationStatus string
	BaselineID         string
	ScoresJSON         string
	ComparisonJSON     string
	Summary            string
	Actor              string
	SourceType         string
	Confidence         float64
	CreatedAt          string
	UpdatedAt          string
}

// DebtItem is a bounded packet row for verification debt (S06).
type DebtItem struct {
	TaskID  string
	GoalID  string
	Missing string
}

func emptyJSON(s string) string {
	if s == "" {
		return emptyJSONObject
	}
	return s
}

// UpsertBaseline inserts or replaces a baseline by id. Empty ID allocates a UUID.
func (s *Store) UpsertBaseline(b Baseline) (Baseline, error) {
	now := nowRFC3339()
	if b.ID == "" {
		b.ID = uuid.NewString()
	}
	if b.GitCommit == "" {
		return Baseline{}, fmt.Errorf("store: upsert baseline: git_commit required")
	}
	if b.ScoresJSON == "" {
		b.ScoresJSON = emptyJSONObject
	}
	if b.Status == "" {
		b.Status = BaselineStatusActive
	}
	if b.CreatedAt == "" {
		b.CreatedAt = now
	}
	b.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO baselines(id, git_commit, scores_json, label, source_type, status, supersedes_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			git_commit = excluded.git_commit,
			scores_json = excluded.scores_json,
			label = excluded.label,
			source_type = excluded.source_type,
			status = excluded.status,
			supersedes_id = excluded.supersedes_id,
			updated_at = excluded.updated_at
	`, b.ID, b.GitCommit, b.ScoresJSON, b.Label, b.SourceType, b.Status, b.SupersedesID, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return Baseline{}, fmt.Errorf("store: upsert baseline: %w", err)
	}
	return s.GetBaseline(b.ID)
}

// ListAllBaselines returns every baseline row.
func (s *Store) ListAllBaselines() ([]Baseline, error) {
	rows, err := s.db.Query(`
		SELECT id, git_commit, scores_json, label, source_type, status, supersedes_id, created_at, updated_at
		FROM baselines
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all baselines: %w", err)
	}
	defer rows.Close()
	var out []Baseline
	for rows.Next() {
		var b Baseline
		if err := rows.Scan(&b.ID, &b.GitCommit, &b.ScoresJSON, &b.Label, &b.SourceType, &b.Status, &b.SupersedesID, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store: scan baseline: %w", err)
		}
		out = append(out, b)
	}
	if out == nil {
		out = []Baseline{}
	}
	return out, rows.Err()
}

// ListAllOutcomeResults returns every outcome_result row.
func (s *Store) ListAllOutcomeResults() ([]OutcomeResult, error) {
	rows, err := s.db.Query(`
		SELECT id, task_id, kind, test_name, test_status, goal_id, verification_status,
			baseline_id, scores_json, comparison_json, summary, actor, source_type,
			confidence, created_at, updated_at
		FROM outcome_results
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all outcome results: %w", err)
	}
	defer rows.Close()
	return scanOutcomeResults(rows)
}

// GetBaseline loads a baseline by id.
func (s *Store) GetBaseline(id string) (Baseline, error) {
	if id == "" {
		return Baseline{}, fmt.Errorf("store: get baseline: id required")
	}
	var b Baseline
	err := s.db.QueryRow(`
		SELECT id, git_commit, scores_json, label, source_type, status, supersedes_id, created_at, updated_at
		FROM baselines WHERE id = ?
	`, id).Scan(&b.ID, &b.GitCommit, &b.ScoresJSON, &b.Label, &b.SourceType, &b.Status, &b.SupersedesID, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return Baseline{}, fmt.Errorf("store: baseline %q: %w", id, err)
	}
	if err != nil {
		return Baseline{}, fmt.Errorf("store: get baseline: %w", err)
	}
	return b, nil
}

// GetActiveBaselineByCommitLabelExcluding returns an active baseline for git_commit+label excluding excludeID.
func (s *Store) GetActiveBaselineByCommitLabelExcluding(gitCommit, label, excludeID string) (Baseline, error) {
	gitCommit = strings.TrimSpace(gitCommit)
	label = strings.TrimSpace(label)
	excludeID = strings.TrimSpace(excludeID)
	if gitCommit == "" {
		return Baseline{}, fmt.Errorf("store: get active baseline: git_commit required")
	}
	var b Baseline
	err := s.db.QueryRow(`
		SELECT id, git_commit, scores_json, label, source_type, status, supersedes_id, created_at, updated_at
		FROM baselines
		WHERE git_commit = ? AND label = ? AND status = ? AND id != ?
		LIMIT 1
	`, gitCommit, label, BaselineStatusActive, excludeID).Scan(
		&b.ID, &b.GitCommit, &b.ScoresJSON, &b.Label, &b.SourceType, &b.Status, &b.SupersedesID, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Baseline{}, fmt.Errorf("store: active baseline %q/%q excluding %q: %w", gitCommit, label, excludeID, sql.ErrNoRows)
	}
	if err != nil {
		return Baseline{}, fmt.Errorf("store: get active baseline: %w", err)
	}
	return b, nil
}

// GetActiveBaselineByCommitLabel returns the active baseline for git_commit+label.
func (s *Store) GetActiveBaselineByCommitLabel(gitCommit, label string) (Baseline, error) {
	gitCommit = strings.TrimSpace(gitCommit)
	label = strings.TrimSpace(label)
	if gitCommit == "" {
		return Baseline{}, fmt.Errorf("store: get active baseline: git_commit required")
	}
	var b Baseline
	err := s.db.QueryRow(`
		SELECT id, git_commit, scores_json, label, source_type, status, supersedes_id, created_at, updated_at
		FROM baselines
		WHERE git_commit = ? AND label = ? AND status = ?
		LIMIT 1
	`, gitCommit, label, BaselineStatusActive).Scan(
		&b.ID, &b.GitCommit, &b.ScoresJSON, &b.Label, &b.SourceType, &b.Status, &b.SupersedesID, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Baseline{}, fmt.Errorf("store: active baseline %q/%q: %w", gitCommit, label, sql.ErrNoRows)
	}
	if err != nil {
		return Baseline{}, fmt.Errorf("store: get active baseline: %w", err)
	}
	return b, nil
}

// SetBaselinePromotion updates status and supersedes_id for a baseline row.
func (s *Store) SetBaselinePromotion(id, status, supersedesID string) error {
	if id == "" {
		return fmt.Errorf("store: set baseline promotion: id required")
	}
	switch status {
	case BaselineStatusActive, BaselineStatusSuperseded:
	default:
		return fmt.Errorf("store: set baseline promotion: unknown status %q", status)
	}
	now := nowRFC3339()
	res, err := s.db.Exec(`
		UPDATE baselines SET status = ?, supersedes_id = ?, updated_at = ? WHERE id = ?
	`, status, supersedesID, now, id)
	if err != nil {
		return fmt.Errorf("store: set baseline promotion: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: set baseline promotion: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("store: baseline %q: %w", id, sql.ErrNoRows)
	}
	return nil
}

// UpsertOutcomeResult inserts or replaces an outcome by id. Empty ID allocates a UUID.
func (s *Store) UpsertOutcomeResult(o OutcomeResult) (OutcomeResult, error) {
	now := nowRFC3339()
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	if o.TaskID == "" {
		return OutcomeResult{}, fmt.Errorf("store: upsert outcome: task_id required")
	}
	if err := validateOutcomeKindColumns(o); err != nil {
		return OutcomeResult{}, err
	}
	o.ScoresJSON = emptyJSON(o.ScoresJSON)
	o.ComparisonJSON = emptyJSON(o.ComparisonJSON)
	if o.CreatedAt == "" {
		o.CreatedAt = now
	}
	o.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO outcome_results(
			id, task_id, kind, test_name, test_status, goal_id, verification_status,
			baseline_id, scores_json, comparison_json, summary, actor, source_type,
			confidence, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			task_id = excluded.task_id,
			kind = excluded.kind,
			test_name = excluded.test_name,
			test_status = excluded.test_status,
			goal_id = excluded.goal_id,
			verification_status = excluded.verification_status,
			baseline_id = excluded.baseline_id,
			scores_json = excluded.scores_json,
			comparison_json = excluded.comparison_json,
			summary = excluded.summary,
			actor = excluded.actor,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			created_at = excluded.created_at,
			updated_at = excluded.updated_at
	`, o.ID, o.TaskID, o.Kind, o.TestName, o.TestStatus, o.GoalID, o.VerificationStatus,
		o.BaselineID, o.ScoresJSON, o.ComparisonJSON, o.Summary, o.Actor, o.SourceType,
		o.Confidence, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		return OutcomeResult{}, fmt.Errorf("store: upsert outcome: %w", err)
	}
	out, err := s.GetOutcomeResult(o.ID)
	if err != nil {
		return OutcomeResult{}, err
	}
	if err := s.SyncEntityFTS("outcome_result", out.ID); err != nil {
		return OutcomeResult{}, err
	}
	return out, nil
}

func validateOutcomeKindColumns(o OutcomeResult) error {
	switch o.Kind {
	case OutcomeKindTest:
		if o.TestName == "" || o.TestStatus == TestStatusNone {
			return fmt.Errorf("store: upsert outcome: test_name and test_status required for kind=test")
		}
		switch o.TestStatus {
		case TestStatusPass, TestStatusFail, TestStatusSkip, TestStatusError:
		default:
			return fmt.Errorf("store: upsert outcome: unknown test_status %q", o.TestStatus)
		}
		if o.GoalID != "" || o.BaselineID != "" || o.VerificationStatus != "" {
			return fmt.Errorf("store: upsert outcome: kind=test forbids goal_id, baseline_id, verification_status")
		}
	case OutcomeKindVerification:
		if o.GoalID == "" || o.VerificationStatus == VerificationStatusNone {
			return fmt.Errorf("store: upsert outcome: goal_id and verification_status required for kind=verification")
		}
		switch o.VerificationStatus {
		case VerificationStatusVerified, VerificationStatusFailed, VerificationStatusPartial:
		default:
			return fmt.Errorf("store: upsert outcome: unknown verification_status %q", o.VerificationStatus)
		}
		if o.TestName != "" || o.TestStatus != "" || o.BaselineID != "" {
			return fmt.Errorf("store: upsert outcome: kind=verification forbids test_name, test_status, baseline_id")
		}
	case OutcomeKindEvaluation:
		if o.BaselineID == "" {
			return fmt.Errorf("store: upsert outcome: baseline_id required for kind=evaluation")
		}
		if o.TestName != "" || o.TestStatus != "" || o.GoalID != "" || o.VerificationStatus != "" {
			return fmt.Errorf("store: upsert outcome: kind=evaluation forbids test_name, test_status, goal_id, verification_status")
		}
	default:
		return fmt.Errorf("store: upsert outcome: unknown kind %q", o.Kind)
	}
	return nil
}

// GetOutcomeResult loads an outcome by id.
func (s *Store) GetOutcomeResult(id string) (OutcomeResult, error) {
	if id == "" {
		return OutcomeResult{}, fmt.Errorf("store: get outcome: id required")
	}
	var o OutcomeResult
	err := s.db.QueryRow(`
		SELECT id, task_id, kind, test_name, test_status, goal_id, verification_status,
			baseline_id, scores_json, comparison_json, summary, actor, source_type,
			confidence, created_at, updated_at
		FROM outcome_results WHERE id = ?
	`, id).Scan(
		&o.ID, &o.TaskID, &o.Kind, &o.TestName, &o.TestStatus, &o.GoalID, &o.VerificationStatus,
		&o.BaselineID, &o.ScoresJSON, &o.ComparisonJSON, &o.Summary, &o.Actor, &o.SourceType,
		&o.Confidence, &o.CreatedAt, &o.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return OutcomeResult{}, fmt.Errorf("store: outcome %q: %w", id, err)
	}
	if err != nil {
		return OutcomeResult{}, fmt.Errorf("store: get outcome: %w", err)
	}
	return o, nil
}

// ListOutcomeResultsByTaskID returns outcomes for a task, oldest first.
func (s *Store) ListOutcomeResultsByTaskID(taskID string) ([]OutcomeResult, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list outcomes: task_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, task_id, kind, test_name, test_status, goal_id, verification_status,
			baseline_id, scores_json, comparison_json, summary, actor, source_type,
			confidence, created_at, updated_at
		FROM outcome_results
		WHERE task_id = ?
		ORDER BY created_at ASC, id ASC
	`, taskID)
	if err != nil {
		return nil, fmt.Errorf("store: list outcomes: %w", err)
	}
	defer rows.Close()
	return scanOutcomeResults(rows)
}

// ListOutcomeResultsByTaskKind returns outcomes of one kind for a task, oldest first.
func (s *Store) ListOutcomeResultsByTaskKind(taskID, kind string) ([]OutcomeResult, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list outcomes by kind: task_id required")
	}
	if kind == "" {
		return nil, fmt.Errorf("store: list outcomes by kind: kind required")
	}
	rows, err := s.db.Query(`
		SELECT id, task_id, kind, test_name, test_status, goal_id, verification_status,
			baseline_id, scores_json, comparison_json, summary, actor, source_type,
			confidence, created_at, updated_at
		FROM outcome_results
		WHERE task_id = ? AND kind = ?
		ORDER BY created_at ASC, rowid ASC
	`, taskID, kind)
	if err != nil {
		return nil, fmt.Errorf("store: list outcomes by kind: %w", err)
	}
	defer rows.Close()
	return scanOutcomeResults(rows)
}

func scanOutcomeResults(rows *sql.Rows) ([]OutcomeResult, error) {
	var out []OutcomeResult
	for rows.Next() {
		var o OutcomeResult
		if err := rows.Scan(
			&o.ID, &o.TaskID, &o.Kind, &o.TestName, &o.TestStatus, &o.GoalID, &o.VerificationStatus,
			&o.BaselineID, &o.ScoresJSON, &o.ComparisonJSON, &o.Summary, &o.Actor, &o.SourceType,
			&o.Confidence, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan outcome: %w", err)
		}
		out = append(out, o)
	}
	if out == nil {
		out = []OutcomeResult{}
	}
	return out, rows.Err()
}

// CountOutcomeSupportedByEvidence counts outcome_supported_by → evidence links.
func (s *Store) CountOutcomeSupportedByEvidence(outcomeID string) (int, error) {
	if outcomeID == "" {
		return 0, fmt.Errorf("store: count outcome evidence: id required")
	}
	var n int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM entity_links
		WHERE from_type = 'outcome_result' AND from_id = ?
		  AND rel = 'outcome_supported_by' AND to_type = 'evidence'
	`, outcomeID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("store: count outcome evidence: %w", err)
	}
	return n, nil
}

// HasImplementationSignal reports ≥1 change for the task in RECORDED or COMPARED.
func (s *Store) HasImplementationSignal(taskID string) (bool, error) {
	if taskID == "" {
		return false, fmt.Errorf("store: implementation signal: task_id required")
	}
	var n int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM changes
		WHERE task_id = ? AND status IN ('RECORDED', 'COMPARED')
	`, taskID).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("store: implementation signal: %w", err)
	}
	return n > 0, nil
}

func taskGoalID(goal *string) string {
	if goal == nil {
		return ""
	}
	return *goal
}

func (s *Store) countSatisfiedVerifications(taskID, goalID string) (int, error) {
	var n int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM outcome_results o
		WHERE o.task_id = ?
		  AND o.kind = 'verification'
		  AND o.verification_status = 'verified'
		  AND o.goal_id = ?
		  AND EXISTS (
			SELECT 1 FROM entity_links l
			WHERE l.from_type = 'outcome_result' AND l.from_id = o.id
			  AND l.rel = 'outcome_supported_by' AND l.to_type = 'evidence'
		  )
	`, taskID, goalID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("store: count satisfied verifications: %w", err)
	}
	return n, nil
}

// HasVerificationDebt is true when implementation is recorded, the task has a
// goal, and no verified outcome with evidence exists for that goal.
func (s *Store) HasVerificationDebt(taskID string) (bool, error) {
	if taskID == "" {
		return false, fmt.Errorf("store: verification debt: task_id required")
	}
	task, err := s.GetTask(taskID)
	if err != nil {
		return false, err
	}
	goalID := taskGoalID(task.GoalID)
	if goalID == "" {
		return false, nil
	}
	signal, err := s.HasImplementationSignal(taskID)
	if err != nil {
		return false, err
	}
	if !signal {
		return false, nil
	}
	n, err := s.countSatisfiedVerifications(taskID, goalID)
	if err != nil {
		return false, err
	}
	return n == 0, nil
}

// ListVerificationDebtSummary returns bounded missing-verification labels for a task.
func (s *Store) ListVerificationDebtSummary(taskID string) ([]DebtItem, error) {
	debt, err := s.HasVerificationDebt(taskID)
	if err != nil {
		return nil, err
	}
	if !debt {
		return []DebtItem{}, nil
	}
	task, err := s.GetTask(taskID)
	if err != nil {
		return nil, err
	}
	return []DebtItem{{
		TaskID:  taskID,
		GoalID:  taskGoalID(task.GoalID),
		Missing: "verification missing for goal",
	}}, nil
}

// ListFailedTestOutcomes returns kind=test rows with fail or error status, newest first.
func (s *Store) ListFailedTestOutcomes(limit int, taskID string) ([]OutcomeResult, error) {
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}
	taskID = strings.TrimSpace(taskID)
	var (
		rows *sql.Rows
		err  error
	)
	if taskID != "" {
		rows, err = s.db.Query(`
			SELECT id, task_id, kind, test_name, test_status, goal_id, verification_status,
				baseline_id, scores_json, comparison_json, summary, actor, source_type,
				confidence, created_at, updated_at
			FROM outcome_results
			WHERE task_id = ? AND kind = ? AND test_status IN (?, ?)
			ORDER BY created_at DESC, id DESC
			LIMIT ?
		`, taskID, OutcomeKindTest, TestStatusFail, TestStatusError, limit)
	} else {
		rows, err = s.db.Query(`
			SELECT id, task_id, kind, test_name, test_status, goal_id, verification_status,
				baseline_id, scores_json, comparison_json, summary, actor, source_type,
				confidence, created_at, updated_at
			FROM outcome_results
			WHERE kind = ? AND test_status IN (?, ?)
			ORDER BY created_at DESC, id DESC
			LIMIT ?
		`, OutcomeKindTest, TestStatusFail, TestStatusError, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("store: list failed test outcomes: %w", err)
	}
	defer rows.Close()
	return scanOutcomeResults(rows)
}

// ListPassingTestOutcomes returns kind=test pass rows, newest first. limit 0 means no SQL LIMIT.
func (s *Store) ListPassingTestOutcomes(limit int, taskID string) ([]OutcomeResult, error) {
	taskID = strings.TrimSpace(taskID)
	query := `
		SELECT id, task_id, kind, test_name, test_status, goal_id, verification_status,
			baseline_id, scores_json, comparison_json, summary, actor, source_type,
			confidence, created_at, updated_at
		FROM outcome_results
		WHERE kind = ? AND test_status = ?`
	args := []any{OutcomeKindTest, TestStatusPass}
	if taskID != "" {
		query += ` AND task_id = ?`
		args = append(args, taskID)
	}
	query += ` ORDER BY created_at DESC, id DESC`
	if limit > 0 {
		if limit > 64 {
			limit = 64
		}
		query += ` LIMIT ?`
		args = append(args, limit)
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("store: list passing test outcomes: %w", err)
	}
	defer rows.Close()
	return scanOutcomeResults(rows)
}

// CountTasksWithVerificationDebt counts tasks that currently have verification debt.
func (s *Store) CountTasksWithVerificationDebt() (int, error) {
	var n int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM tasks t
		WHERE t.goal_id IS NOT NULL AND t.goal_id != ''
		  AND EXISTS (
			SELECT 1 FROM changes c
			WHERE c.task_id = t.id AND c.status IN ('RECORDED', 'COMPARED')
		  )
		  AND NOT EXISTS (
			SELECT 1 FROM outcome_results o
			WHERE o.task_id = t.id
			  AND o.kind = 'verification'
			  AND o.verification_status = 'verified'
			  AND o.goal_id = t.goal_id
			  AND EXISTS (
				SELECT 1 FROM entity_links l
				WHERE l.from_type = 'outcome_result' AND l.from_id = o.id
				  AND l.rel = 'outcome_supported_by' AND l.to_type = 'evidence'
			  )
		  )
	`).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("store: count verification debt: %w", err)
	}
	return n, nil
}
