package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Impact class vocabulary (canonical; domain rejects unknown/empty).
const (
	ImpactClassSAFE        = "SAFE"
	ImpactClassCaution     = "CAUTION"
	ImpactClassHigh        = "HIGH"
	ImpactClassDestructive = "DESTRUCTIVE"
	ImpactClassReversal    = "REVERSAL"
)

// Uncertainty vocabulary (canonical; empty on create → UNKNOWN at domain).
const (
	UncertaintyKNOWN    = "KNOWN"
	UncertaintyLIKELY   = "LIKELY"
	UncertaintyPOSSIBLE = "POSSIBLE"
	UncertaintyUNKNOWN  = "UNKNOWN"
)

// Finding kind vocabulary (canonical; domain rejects unknown/empty).
const (
	FindingKindAffectedWork          = "AFFECTED_WORK"
	FindingKindInvalidatedAssumption = "INVALIDATED_ASSUMPTION"
	FindingKindWorkAtRisk            = "WORK_AT_RISK"
	FindingKindNewWork               = "NEW_WORK"
	FindingKindUnresolved            = "UNRESOLVED"
)

// DecisionImpactFinding is a planted/manual impact finding on a decision.
type DecisionImpactFinding struct {
	ID          string `json:"id"`
	DecisionID  string `json:"decision_id"`
	ImpactClass string `json:"impact_class"`
	Uncertainty string `json:"uncertainty"`
	Kind        string `json:"kind"`
	Body        string `json:"body"`
	RelatedType string `json:"related_type"`
	RelatedID   string `json:"related_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// DecisionAlternative is a thin alternative row on a decision.
type DecisionAlternative struct {
	ID            string `json:"id"`
	DecisionID    string `json:"decision_id"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	IsRecommended bool   `json:"is_recommended"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// InsertDecisionImpactFinding inserts a finding. Empty ID allocates a UUID.
func (s *Store) InsertDecisionImpactFinding(f DecisionImpactFinding) (DecisionImpactFinding, error) {
	now := nowRFC3339()
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	if f.DecisionID == "" {
		return DecisionImpactFinding{}, fmt.Errorf("store: insert decision impact finding: decision_id required")
	}
	if f.ImpactClass == "" {
		return DecisionImpactFinding{}, fmt.Errorf("store: insert decision impact finding: impact_class required")
	}
	if f.Kind == "" {
		return DecisionImpactFinding{}, fmt.Errorf("store: insert decision impact finding: kind required")
	}
	if f.Uncertainty == "" {
		f.Uncertainty = UncertaintyUNKNOWN
	}
	if f.CreatedAt == "" {
		f.CreatedAt = now
	}
	f.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO decision_impact_findings(
			id, decision_id, impact_class, uncertainty, kind, body,
			related_type, related_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, f.ID, f.DecisionID, f.ImpactClass, f.Uncertainty, f.Kind, f.Body,
		f.RelatedType, f.RelatedID, f.CreatedAt, f.UpdatedAt)
	if err != nil {
		return DecisionImpactFinding{}, fmt.Errorf("store: insert decision impact finding: %w", err)
	}
	return s.GetDecisionImpactFinding(f.ID)
}

// UpsertDecisionImpactFinding inserts or updates a finding by id.
func (s *Store) UpsertDecisionImpactFinding(f DecisionImpactFinding) (DecisionImpactFinding, error) {
	now := nowRFC3339()
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	if f.DecisionID == "" {
		return DecisionImpactFinding{}, fmt.Errorf("store: upsert decision impact finding: decision_id required")
	}
	if f.ImpactClass == "" {
		return DecisionImpactFinding{}, fmt.Errorf("store: upsert decision impact finding: impact_class required")
	}
	if f.Kind == "" {
		return DecisionImpactFinding{}, fmt.Errorf("store: upsert decision impact finding: kind required")
	}
	if f.Uncertainty == "" {
		f.Uncertainty = UncertaintyUNKNOWN
	}
	if f.CreatedAt == "" {
		f.CreatedAt = now
	}
	f.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO decision_impact_findings(
			id, decision_id, impact_class, uncertainty, kind, body,
			related_type, related_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			decision_id = excluded.decision_id,
			impact_class = excluded.impact_class,
			uncertainty = excluded.uncertainty,
			kind = excluded.kind,
			body = excluded.body,
			related_type = excluded.related_type,
			related_id = excluded.related_id,
			updated_at = excluded.updated_at
	`, f.ID, f.DecisionID, f.ImpactClass, f.Uncertainty, f.Kind, f.Body,
		f.RelatedType, f.RelatedID, f.CreatedAt, f.UpdatedAt)
	if err != nil {
		return DecisionImpactFinding{}, fmt.Errorf("store: upsert decision impact finding: %w", err)
	}
	return s.GetDecisionImpactFinding(f.ID)
}

// GetDecisionImpactFinding loads a finding by id.
func (s *Store) GetDecisionImpactFinding(id string) (DecisionImpactFinding, error) {
	var f DecisionImpactFinding
	err := s.db.QueryRow(`
		SELECT id, decision_id, impact_class, uncertainty, kind, body,
			related_type, related_id, created_at, updated_at
		FROM decision_impact_findings WHERE id = ?
	`, id).Scan(
		&f.ID, &f.DecisionID, &f.ImpactClass, &f.Uncertainty, &f.Kind, &f.Body,
		&f.RelatedType, &f.RelatedID, &f.CreatedAt, &f.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return DecisionImpactFinding{}, fmt.Errorf("store: decision impact finding %q: %w", id, err)
	}
	if err != nil {
		return DecisionImpactFinding{}, fmt.Errorf("store: get decision impact finding: %w", err)
	}
	return f, nil
}

// ListDecisionImpactFindingsByDecisionID returns findings for a decision, ordered by created_at.
func (s *Store) ListDecisionImpactFindingsByDecisionID(decisionID string) ([]DecisionImpactFinding, error) {
	if decisionID == "" {
		return nil, fmt.Errorf("store: list decision impact findings: decision_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, decision_id, impact_class, uncertainty, kind, body,
			related_type, related_id, created_at, updated_at
		FROM decision_impact_findings
		WHERE decision_id = ?
		ORDER BY created_at ASC, id ASC
	`, decisionID)
	if err != nil {
		return nil, fmt.Errorf("store: list decision impact findings: %w", err)
	}
	defer rows.Close()
	var out []DecisionImpactFinding
	for rows.Next() {
		var f DecisionImpactFinding
		if err := rows.Scan(
			&f.ID, &f.DecisionID, &f.ImpactClass, &f.Uncertainty, &f.Kind, &f.Body,
			&f.RelatedType, &f.RelatedID, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan decision impact finding: %w", err)
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// InsertDecisionAlternative inserts an alternative. Empty ID allocates a UUID.
func (s *Store) InsertDecisionAlternative(a DecisionAlternative) (DecisionAlternative, error) {
	now := nowRFC3339()
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	if a.DecisionID == "" {
		return DecisionAlternative{}, fmt.Errorf("store: insert decision alternative: decision_id required")
	}
	if a.Title == "" {
		return DecisionAlternative{}, fmt.Errorf("store: insert decision alternative: title required")
	}
	if a.CreatedAt == "" {
		a.CreatedAt = now
	}
	a.UpdatedAt = now
	rec := 0
	if a.IsRecommended {
		rec = 1
	}

	_, err := s.db.Exec(`
		INSERT INTO decision_alternatives(
			id, decision_id, title, body, is_recommended, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, a.ID, a.DecisionID, a.Title, a.Body, rec, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return DecisionAlternative{}, fmt.Errorf("store: insert decision alternative: %w", err)
	}
	return s.GetDecisionAlternative(a.ID)
}

// UpsertDecisionAlternative inserts or updates an alternative by id.
func (s *Store) UpsertDecisionAlternative(a DecisionAlternative) (DecisionAlternative, error) {
	now := nowRFC3339()
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	if a.DecisionID == "" {
		return DecisionAlternative{}, fmt.Errorf("store: upsert decision alternative: decision_id required")
	}
	if a.Title == "" {
		return DecisionAlternative{}, fmt.Errorf("store: upsert decision alternative: title required")
	}
	if a.CreatedAt == "" {
		a.CreatedAt = now
	}
	a.UpdatedAt = now
	rec := 0
	if a.IsRecommended {
		rec = 1
	}

	_, err := s.db.Exec(`
		INSERT INTO decision_alternatives(
			id, decision_id, title, body, is_recommended, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			decision_id = excluded.decision_id,
			title = excluded.title,
			body = excluded.body,
			is_recommended = excluded.is_recommended,
			updated_at = excluded.updated_at
	`, a.ID, a.DecisionID, a.Title, a.Body, rec, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return DecisionAlternative{}, fmt.Errorf("store: upsert decision alternative: %w", err)
	}
	return s.GetDecisionAlternative(a.ID)
}

// GetDecisionAlternative loads an alternative by id.
func (s *Store) GetDecisionAlternative(id string) (DecisionAlternative, error) {
	var a DecisionAlternative
	var rec int
	err := s.db.QueryRow(`
		SELECT id, decision_id, title, body, is_recommended, created_at, updated_at
		FROM decision_alternatives WHERE id = ?
	`, id).Scan(
		&a.ID, &a.DecisionID, &a.Title, &a.Body, &rec, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return DecisionAlternative{}, fmt.Errorf("store: decision alternative %q: %w", id, err)
	}
	if err != nil {
		return DecisionAlternative{}, fmt.Errorf("store: get decision alternative: %w", err)
	}
	a.IsRecommended = rec != 0
	return a, nil
}

// ClearRecommendedAlternatives clears is_recommended for all alternatives on a decision.
func (s *Store) ClearRecommendedAlternatives(decisionID string) error {
	if decisionID == "" {
		return fmt.Errorf("store: clear recommended alternatives: decision_id required")
	}
	_, err := s.db.Exec(`
		UPDATE decision_alternatives SET is_recommended = 0, updated_at = ?
		WHERE decision_id = ? AND is_recommended != 0
	`, nowRFC3339(), decisionID)
	if err != nil {
		return fmt.Errorf("store: clear recommended alternatives: %w", err)
	}
	return nil
}

// UpdateDecisionAlternativeRecommended sets is_recommended for one alternative.
func (s *Store) UpdateDecisionAlternativeRecommended(id string, recommended bool) error {
	if id == "" {
		return fmt.Errorf("store: update decision alternative recommended: id required")
	}
	rec := 0
	if recommended {
		rec = 1
	}
	res, err := s.db.Exec(`
		UPDATE decision_alternatives SET is_recommended = ?, updated_at = ? WHERE id = ?
	`, rec, nowRFC3339(), id)
	if err != nil {
		return fmt.Errorf("store: update decision alternative recommended: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: update decision alternative recommended: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("store: decision alternative %q: %w", id, sql.ErrNoRows)
	}
	return nil
}

// ListAllDecisionImpactFindings returns every impact finding row.
func (s *Store) ListAllDecisionImpactFindings() ([]DecisionImpactFinding, error) {
	rows, err := s.db.Query(`
		SELECT id, decision_id, impact_class, uncertainty, kind, body,
			related_type, related_id, created_at, updated_at
		FROM decision_impact_findings
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all decision impact findings: %w", err)
	}
	defer rows.Close()
	var out []DecisionImpactFinding
	for rows.Next() {
		var f DecisionImpactFinding
		if err := rows.Scan(
			&f.ID, &f.DecisionID, &f.ImpactClass, &f.Uncertainty, &f.Kind, &f.Body,
			&f.RelatedType, &f.RelatedID, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan decision impact finding: %w", err)
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// ListAllDecisionAlternatives returns every decision alternative row.
func (s *Store) ListAllDecisionAlternatives() ([]DecisionAlternative, error) {
	rows, err := s.db.Query(`
		SELECT id, decision_id, title, body, is_recommended, created_at, updated_at
		FROM decision_alternatives
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list all decision alternatives: %w", err)
	}
	defer rows.Close()
	var out []DecisionAlternative
	for rows.Next() {
		var a DecisionAlternative
		var rec int
		if err := rows.Scan(
			&a.ID, &a.DecisionID, &a.Title, &a.Body, &rec, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan decision alternative: %w", err)
		}
		a.IsRecommended = rec != 0
		out = append(out, a)
	}
	return out, rows.Err()
}

// ListDecisionAlternativesByDecisionID returns alternatives for a decision, ordered by created_at.
func (s *Store) ListDecisionAlternativesByDecisionID(decisionID string) ([]DecisionAlternative, error) {
	if decisionID == "" {
		return nil, fmt.Errorf("store: list decision alternatives: decision_id required")
	}
	rows, err := s.db.Query(`
		SELECT id, decision_id, title, body, is_recommended, created_at, updated_at
		FROM decision_alternatives
		WHERE decision_id = ?
		ORDER BY created_at ASC, id ASC
	`, decisionID)
	if err != nil {
		return nil, fmt.Errorf("store: list decision alternatives: %w", err)
	}
	defer rows.Close()
	var out []DecisionAlternative
	for rows.Next() {
		var a DecisionAlternative
		var rec int
		if err := rows.Scan(
			&a.ID, &a.DecisionID, &a.Title, &a.Body, &rec, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("store: scan decision alternative: %w", err)
		}
		a.IsRecommended = rec != 0
		out = append(out, a)
	}
	return out, rows.Err()
}
