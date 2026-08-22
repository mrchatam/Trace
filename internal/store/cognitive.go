package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Uncertainty lifecycle (not provenance ACTIVE/STALE).
const (
	UncertaintyStatusOpen       = "OPEN"
	UncertaintyStatusResolved   = "RESOLVED"
	UncertaintyStatusSuperseded = "SUPERSEDED"
)

// Uncertainty severity: INFO | BLOCKING only (PLAN_AFFECTING is Discovery-only).
const (
	UncertaintySeverityINFO     = "INFO"
	UncertaintySeverityBlocking = "BLOCKING"
)

// Uncertainty kind merge (Risk lives here as kind=risk).
const (
	UncertaintyKindNone    = ""
	UncertaintyKindRisk    = "risk"
	UncertaintyKindGap     = "gap"
	UncertaintyKindUnknown = "unknown"
)

// Hypothesis lifecycle.
const (
	HypothesisStatusOpen       = "OPEN"
	HypothesisStatusConfirmed  = "CONFIRMED"
	HypothesisStatusRejected   = "REJECTED"
	HypothesisStatusSuperseded = "SUPERSEDED"
)

// Decision reconsideration trigger vocabulary.
const (
	ReconsiderTriggerContradictedEffect    = "contradicted_effect"
	ReconsiderTriggerNewEvidence           = "new_evidence"
	ReconsiderTriggerInvalidatedAssumption = "invalidated_assumption"
)

// Decision reconsideration status: OPEN = watch; FIRED = trigger satisfied now.
const (
	ReconsiderStatusOpen  = "OPEN"
	ReconsiderStatusFired = "FIRED"
)

// Uncertainty is a question/gap (title = question text).
type Uncertainty struct {
	ID             string
	Title          string
	Body           string
	Severity       string
	Status         string
	Kind           string
	Confidence     float64
	SourceType     string
	Resolution     string
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// Hypothesis is a testable statement (title = statement).
type Hypothesis struct {
	ID             string
	Title          string
	Body           string
	Status         string
	Confidence     float64
	SourceType     string
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// DecisionReconsideration is an append-only child row on a decision.
type DecisionReconsideration struct {
	ID           string
	DecisionID   string
	Trigger      string
	Status       string
	Reason       string
	RelatedType  string
	RelatedID    string
	ReconsiderAt string
	CreatedAt    string
}

// UpsertUncertainty inserts or replaces an uncertainty by id. Empty ID allocates a UUID.
func (s *Store) UpsertUncertainty(u Uncertainty) (Uncertainty, error) {
	now := nowRFC3339()
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if u.Title == "" {
		return Uncertainty{}, fmt.Errorf("store: upsert uncertainty: title required")
	}
	if u.Severity == "" {
		u.Severity = UncertaintySeverityINFO
	}
	if u.Status == "" {
		u.Status = UncertaintyStatusOpen
	}
	if u.CreatedAt == "" {
		u.CreatedAt = now
	}
	u.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO uncertainties(
			id, title, body, severity, status, kind, confidence, source_type,
			resolution, created_at, updated_at, last_verified_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body = excluded.body,
			severity = excluded.severity,
			status = excluded.status,
			kind = excluded.kind,
			confidence = excluded.confidence,
			source_type = excluded.source_type,
			resolution = excluded.resolution,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, u.ID, u.Title, u.Body, u.Severity, u.Status, u.Kind, u.Confidence, u.SourceType,
		u.Resolution, u.CreatedAt, u.UpdatedAt, nullStr(u.LastVerifiedAt))
	if err != nil {
		return Uncertainty{}, fmt.Errorf("store: upsert uncertainty: %w", err)
	}
	out, err := s.GetUncertainty(u.ID)
	if err != nil {
		return Uncertainty{}, err
	}
	if err := s.SyncEntityFTS("uncertainty", out.ID); err != nil {
		return Uncertainty{}, err
	}
	return out, nil
}

// GetUncertainty loads an uncertainty by id.
func (s *Store) GetUncertainty(id string) (Uncertainty, error) {
	if id == "" {
		return Uncertainty{}, fmt.Errorf("store: get uncertainty: id required")
	}
	var u Uncertainty
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, title, body, severity, status, kind, confidence, source_type,
			resolution, created_at, updated_at, last_verified_at
		FROM uncertainties WHERE id = ?
	`, id).Scan(
		&u.ID, &u.Title, &u.Body, &u.Severity, &u.Status, &u.Kind, &u.Confidence, &u.SourceType,
		&u.Resolution, &u.CreatedAt, &u.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Uncertainty{}, fmt.Errorf("store: uncertainty %q: %w", id, err)
	}
	if err != nil {
		return Uncertainty{}, fmt.Errorf("store: get uncertainty: %w", err)
	}
	u.LastVerifiedAt = nullStrPtr(lastVerified)
	return u, nil
}

// UpsertHypothesis inserts or replaces a hypothesis by id. Empty ID allocates a UUID.
func (s *Store) UpsertHypothesis(h Hypothesis) (Hypothesis, error) {
	now := nowRFC3339()
	if h.ID == "" {
		h.ID = uuid.NewString()
	}
	if h.Title == "" {
		return Hypothesis{}, fmt.Errorf("store: upsert hypothesis: title required")
	}
	if h.Status == "" {
		h.Status = HypothesisStatusOpen
	}
	if h.CreatedAt == "" {
		h.CreatedAt = now
	}
	h.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO hypotheses(
			id, title, body, status, confidence, source_type,
			created_at, updated_at, last_verified_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body = excluded.body,
			status = excluded.status,
			confidence = excluded.confidence,
			source_type = excluded.source_type,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, h.ID, h.Title, h.Body, h.Status, h.Confidence, h.SourceType,
		h.CreatedAt, h.UpdatedAt, nullStr(h.LastVerifiedAt))
	if err != nil {
		return Hypothesis{}, fmt.Errorf("store: upsert hypothesis: %w", err)
	}
	out, err := s.GetHypothesis(h.ID)
	if err != nil {
		return Hypothesis{}, err
	}
	if err := s.SyncEntityFTS("hypothesis", out.ID); err != nil {
		return Hypothesis{}, err
	}
	return out, nil
}

// GetHypothesis loads a hypothesis by id.
func (s *Store) GetHypothesis(id string) (Hypothesis, error) {
	if id == "" {
		return Hypothesis{}, fmt.Errorf("store: get hypothesis: id required")
	}
	var h Hypothesis
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, title, body, status, confidence, source_type,
			created_at, updated_at, last_verified_at
		FROM hypotheses WHERE id = ?
	`, id).Scan(
		&h.ID, &h.Title, &h.Body, &h.Status, &h.Confidence, &h.SourceType,
		&h.CreatedAt, &h.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Hypothesis{}, fmt.Errorf("store: hypothesis %q: %w", id, err)
	}
	if err != nil {
		return Hypothesis{}, fmt.Errorf("store: get hypothesis: %w", err)
	}
	h.LastVerifiedAt = nullStrPtr(lastVerified)
	return h, nil
}

// InsertDecisionReconsideration appends a reconsideration row. Empty ID allocates a UUID.
// There is no update/delete API (Law 11).
func (s *Store) InsertDecisionReconsideration(r DecisionReconsideration) (DecisionReconsideration, error) {
	now := nowRFC3339()
	if r.DecisionID == "" {
		return DecisionReconsideration{}, fmt.Errorf("store: insert decision reconsideration: decision_id required")
	}
	if r.Trigger == "" {
		return DecisionReconsideration{}, fmt.Errorf("store: insert decision reconsideration: trigger required")
	}
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.Status == "" {
		r.Status = ReconsiderStatusFired
	}
	if r.ReconsiderAt == "" {
		r.ReconsiderAt = now
	}
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}

	_, err := s.db.Exec(`
		INSERT INTO decision_reconsiderations(
			id, decision_id, trigger, status, reason, related_type, related_id,
			reconsider_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.ID, r.DecisionID, r.Trigger, r.Status, r.Reason, r.RelatedType, r.RelatedID,
		r.ReconsiderAt, r.CreatedAt)
	if err != nil {
		return DecisionReconsideration{}, fmt.Errorf("store: insert decision reconsideration: %w", err)
	}
	return s.GetDecisionReconsideration(r.ID)
}

// GetDecisionReconsideration loads a reconsideration by id.
func (s *Store) GetDecisionReconsideration(id string) (DecisionReconsideration, error) {
	if id == "" {
		return DecisionReconsideration{}, fmt.Errorf("store: get decision reconsideration: id required")
	}
	var r DecisionReconsideration
	err := s.db.QueryRow(`
		SELECT id, decision_id, trigger, status, reason, related_type, related_id,
			reconsider_at, created_at
		FROM decision_reconsiderations WHERE id = ?
	`, id).Scan(
		&r.ID, &r.DecisionID, &r.Trigger, &r.Status, &r.Reason, &r.RelatedType, &r.RelatedID,
		&r.ReconsiderAt, &r.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return DecisionReconsideration{}, fmt.Errorf("store: decision reconsideration %q: %w", id, err)
	}
	if err != nil {
		return DecisionReconsideration{}, fmt.Errorf("store: get decision reconsideration: %w", err)
	}
	return r, nil
}

// ListAllUncertainties returns every uncertainty row.
func (s *Store) ListAllUncertainties() ([]Uncertainty, error) {
	rows, err := s.db.Query(`
		SELECT id, title, body, severity, status, kind, confidence, source_type,
			resolution, created_at, updated_at, last_verified_at
		FROM uncertainties
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all uncertainties: %w", err)
	}
	defer rows.Close()
	var out []Uncertainty
	for rows.Next() {
		var u Uncertainty
		var lastVerified sql.NullString
		if err := rows.Scan(
			&u.ID, &u.Title, &u.Body, &u.Severity, &u.Status, &u.Kind, &u.Confidence, &u.SourceType,
			&u.Resolution, &u.CreatedAt, &u.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan uncertainty: %w", err)
		}
		u.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, u)
	}
	if out == nil {
		out = []Uncertainty{}
	}
	return out, rows.Err()
}

// ListAllHypotheses returns every hypothesis row.
func (s *Store) ListAllHypotheses() ([]Hypothesis, error) {
	rows, err := s.db.Query(`
		SELECT id, title, body, status, confidence, source_type,
			created_at, updated_at, last_verified_at
		FROM hypotheses
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all hypotheses: %w", err)
	}
	defer rows.Close()
	var out []Hypothesis
	for rows.Next() {
		var h Hypothesis
		var lastVerified sql.NullString
		if err := rows.Scan(
			&h.ID, &h.Title, &h.Body, &h.Status, &h.Confidence, &h.SourceType,
			&h.CreatedAt, &h.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan hypothesis: %w", err)
		}
		h.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, h)
	}
	if out == nil {
		out = []Hypothesis{}
	}
	return out, rows.Err()
}

// ListAllDecisionReconsiderations returns every decision reconsideration row.
func (s *Store) ListAllDecisionReconsiderations() ([]DecisionReconsideration, error) {
	rows, err := s.db.Query(`
		SELECT id, decision_id, trigger, status, reason, related_type, related_id,
			reconsider_at, created_at
		FROM decision_reconsiderations
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all decision reconsiderations: %w", err)
	}
	defer rows.Close()
	var out []DecisionReconsideration
	for rows.Next() {
		var r DecisionReconsideration
		if err := rows.Scan(
			&r.ID, &r.DecisionID, &r.Trigger, &r.Status, &r.Reason, &r.RelatedType, &r.RelatedID,
			&r.ReconsiderAt, &r.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan decision reconsideration: %w", err)
		}
		out = append(out, r)
	}
	if out == nil {
		out = []DecisionReconsideration{}
	}
	return out, rows.Err()
}

// UpsertDecisionReconsideration inserts or replaces a reconsideration by id.
func (s *Store) UpsertDecisionReconsideration(r DecisionReconsideration) (DecisionReconsideration, error) {
	now := nowRFC3339()
	if r.DecisionID == "" {
		return DecisionReconsideration{}, fmt.Errorf("store: upsert decision reconsideration: decision_id required")
	}
	if r.Trigger == "" {
		return DecisionReconsideration{}, fmt.Errorf("store: upsert decision reconsideration: trigger required")
	}
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.Status == "" {
		r.Status = ReconsiderStatusFired
	}
	if r.ReconsiderAt == "" {
		r.ReconsiderAt = now
	}
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}

	_, err := s.db.Exec(`
		INSERT INTO decision_reconsiderations(
			id, decision_id, trigger, status, reason, related_type, related_id,
			reconsider_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			decision_id = excluded.decision_id,
			trigger = excluded.trigger,
			status = excluded.status,
			reason = excluded.reason,
			related_type = excluded.related_type,
			related_id = excluded.related_id,
			reconsider_at = excluded.reconsider_at
	`, r.ID, r.DecisionID, r.Trigger, r.Status, r.Reason, r.RelatedType, r.RelatedID,
		r.ReconsiderAt, r.CreatedAt)
	if err != nil {
		return DecisionReconsideration{}, fmt.Errorf("store: upsert decision reconsideration: %w", err)
	}
	return s.GetDecisionReconsideration(r.ID)
}

// ListDecisionReconsiderationsByDecisionID returns append-only reconsideration rows for a decision.
func (s *Store) ListDecisionReconsiderationsByDecisionID(decisionID string) ([]DecisionReconsideration, error) {
	if decisionID == "" {
		return nil, fmt.Errorf("store: list decision reconsiderations: decision_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, decision_id, trigger, status, reason, related_type, related_id,
			reconsider_at, created_at
		FROM decision_reconsiderations
		WHERE decision_id = ?
		ORDER BY created_at ASC, id ASC
	`, decisionID)
	if err != nil {
		return nil, fmt.Errorf("store: list decision reconsiderations: %w", err)
	}
	defer rows.Close()
	var out []DecisionReconsideration
	for rows.Next() {
		var r DecisionReconsideration
		if err := rows.Scan(
			&r.ID, &r.DecisionID, &r.Trigger, &r.Status, &r.Reason, &r.RelatedType, &r.RelatedID,
			&r.ReconsiderAt, &r.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan decision reconsideration: %w", err)
		}
		out = append(out, r)
	}
	if out == nil {
		out = []DecisionReconsideration{}
	}
	return out, rows.Err()
}

// OpenUncertaintyRow is a bounded open uncertainty for loop packets.
type OpenUncertaintyRow struct {
	ID       string
	Title    string
	Severity string
	Status   string
	Kind     string
}

// ListOpenUncertaintiesByTaskID returns OPEN uncertainties scoped to a task (via
// uncertainty_blocks_task) or its goal (INFO via uncertainty_affects_goal). BLOCKING
// rows sort before INFO; capped at limit.
func (s *Store) ListOpenUncertaintiesByTaskID(taskID string, limit int) ([]OpenUncertaintyRow, error) {
	if taskID == "" {
		return nil, fmt.Errorf("store: list open uncertainties: task_id required")
	}
	if limit <= 0 {
		limit = 16
	}
	task, err := s.GetTask(taskID)
	if err != nil {
		return nil, err
	}
	goalID := ""
	if task.GoalID != nil {
		goalID = *task.GoalID
	}
	rows, err := s.db.Query(`
		SELECT DISTINCT u.id, u.title, u.severity, u.status, u.kind,
			CASE WHEN u.severity = 'BLOCKING' THEN 0 ELSE 1 END AS sort_rank
		FROM uncertainties u
		INNER JOIN entity_links l ON l.from_type = 'uncertainty' AND l.from_id = u.id
		WHERE u.status = 'OPEN'
		  AND (
		    (l.rel = 'uncertainty_blocks_task' AND l.to_type = 'task' AND l.to_id = ?)
		    OR (? != '' AND l.rel = 'uncertainty_affects_goal' AND l.to_type = 'goal' AND l.to_id = ? AND u.severity = 'INFO')
		  )
		ORDER BY sort_rank ASC, u.created_at ASC, u.id ASC
		LIMIT ?
	`, taskID, goalID, goalID, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list open uncertainties: %w", err)
	}
	defer rows.Close()
	var out []OpenUncertaintyRow
	for rows.Next() {
		var row OpenUncertaintyRow
		var sortRank int
		if err := rows.Scan(&row.ID, &row.Title, &row.Severity, &row.Status, &row.Kind, &sortRank); err != nil {
			return nil, fmt.Errorf("store: scan open uncertainty: %w", err)
		}
		out = append(out, row)
	}
	if out == nil {
		out = []OpenUncertaintyRow{}
	}
	return out, rows.Err()
}

// CountOpenBlockingUncertaintiesByTaskID counts OPEN BLOCKING uncertainties linked
// uncertainty_blocks_task to the given task. Empty taskID fails closed.
func (s *Store) CountOpenBlockingUncertaintiesByTaskID(taskID string) (int, error) {
	if taskID == "" {
		return 0, fmt.Errorf("store: count open blocking uncertainties: task_id required")
	}
	var n int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM uncertainties u
		INNER JOIN entity_links l
		  ON l.from_type = 'uncertainty' AND l.from_id = u.id
		WHERE l.rel = 'uncertainty_blocks_task'
		  AND l.to_type = 'task' AND l.to_id = ?
		  AND u.severity = 'BLOCKING'
		  AND u.status = 'OPEN'
	`, taskID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("store: count open blocking uncertainties: %w", err)
	}
	return n, nil
}
