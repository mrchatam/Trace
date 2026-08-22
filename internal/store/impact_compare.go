package store

import (
	"database/sql"
	"fmt"
)

// ImpactPrediction is a stored impact walk snapshot for one change.
type ImpactPrediction struct {
	ChangeID      string
	PredictedJSON string
	CompareJSON   string
	Depth         int
	CreatedAt     string
	ComparedAt    string
}

// UpsertImpactPrediction inserts or replaces the prediction row for change_id.
func (s *Store) UpsertImpactPrediction(p ImpactPrediction) (ImpactPrediction, error) {
	if p.ChangeID == "" {
		return ImpactPrediction{}, fmt.Errorf("store: upsert impact prediction: change_id required")
	}
	if p.PredictedJSON == "" {
		return ImpactPrediction{}, fmt.Errorf("store: upsert impact prediction: predicted_json required")
	}
	if p.Depth < 1 {
		return ImpactPrediction{}, fmt.Errorf("store: upsert impact prediction: depth required")
	}
	now := nowRFC3339Nano()
	p.CreatedAt = now
	if p.CompareJSON == "" {
		p.CompareJSON = ""
	}
	if p.ComparedAt == "" {
		p.ComparedAt = ""
	}

	_, err := s.db.Exec(`
		INSERT INTO impact_predictions(change_id, predicted_json, compare_json, depth, created_at, compared_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(change_id) DO UPDATE SET
			predicted_json = excluded.predicted_json,
			compare_json = '',
			depth = excluded.depth,
			created_at = excluded.created_at,
			compared_at = ''
	`, p.ChangeID, p.PredictedJSON, p.CompareJSON, p.Depth, p.CreatedAt, p.ComparedAt)
	if err != nil {
		return ImpactPrediction{}, fmt.Errorf("store: upsert impact prediction: %w", err)
	}
	return s.GetImpactPrediction(p.ChangeID)
}

// GetImpactPrediction loads the prediction row for change_id. Missing returns sql.ErrNoRows.
func (s *Store) GetImpactPrediction(changeID string) (ImpactPrediction, error) {
	if changeID == "" {
		return ImpactPrediction{}, fmt.Errorf("store: get impact prediction: change_id required")
	}
	var p ImpactPrediction
	err := s.db.QueryRow(`
		SELECT change_id, predicted_json, compare_json, depth, created_at, compared_at
		FROM impact_predictions WHERE change_id = ?
	`, changeID).Scan(&p.ChangeID, &p.PredictedJSON, &p.CompareJSON, &p.Depth, &p.CreatedAt, &p.ComparedAt)
	if err == sql.ErrNoRows {
		return ImpactPrediction{}, err
	}
	if err != nil {
		return ImpactPrediction{}, fmt.Errorf("store: get impact prediction: %w", err)
	}
	return p, nil
}

// UpdateImpactPredictionCompare persists compare_json and compared_at on an existing row.
func (s *Store) UpdateImpactPredictionCompare(changeID, compareJSON, comparedAt string) (ImpactPrediction, error) {
	if changeID == "" {
		return ImpactPrediction{}, fmt.Errorf("store: update impact compare: change_id required")
	}
	if compareJSON == "" {
		return ImpactPrediction{}, fmt.Errorf("store: update impact compare: compare_json required")
	}
	if comparedAt == "" {
		comparedAt = nowRFC3339()
	}
	res, err := s.db.Exec(`
		UPDATE impact_predictions
		SET compare_json = ?, compared_at = ?
		WHERE change_id = ?
	`, compareJSON, comparedAt, changeID)
	if err != nil {
		return ImpactPrediction{}, fmt.Errorf("store: update impact compare: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return ImpactPrediction{}, fmt.Errorf("store: update impact compare rows: %w", err)
	}
	if n == 0 {
		return ImpactPrediction{}, sql.ErrNoRows
	}
	return s.GetImpactPrediction(changeID)
}
