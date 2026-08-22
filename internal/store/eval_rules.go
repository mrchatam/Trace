package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

const EvalRuleSetDefaultID = "default"

// EvalRuleSet is a cached snapshot of a committed eval-rules file.
type EvalRuleSet struct {
	ID         string
	SourcePath string
	BodyJSON   string
	UpdatedAt  string
}

func requireEvalRuleBodyJSON(raw string) (string, error) {
	if raw == "" {
		return "{}", nil
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return "", fmt.Errorf("store: eval_rule_sets: body_json must be valid JSON")
	}
	if _, ok := v.(map[string]any); !ok {
		return "", fmt.Errorf("store: eval_rule_sets: body_json must be a JSON object")
	}
	return raw, nil
}

// UpsertEvalRuleSet inserts or replaces an eval_rule_sets row by id.
func (s *Store) UpsertEvalRuleSet(row EvalRuleSet) (EvalRuleSet, error) {
	if strings.TrimSpace(row.ID) == "" {
		row.ID = EvalRuleSetDefaultID
	}
	bodyJSON, err := requireEvalRuleBodyJSON(row.BodyJSON)
	if err != nil {
		return EvalRuleSet{}, err
	}
	now := nowRFC3339()
	if row.UpdatedAt == "" {
		row.UpdatedAt = now
	}
	row.BodyJSON = bodyJSON

	_, err = s.db.Exec(`
		INSERT INTO eval_rule_sets(id, source_path, body_json, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			source_path = excluded.source_path,
			body_json = excluded.body_json,
			updated_at = excluded.updated_at
	`, row.ID, row.SourcePath, row.BodyJSON, row.UpdatedAt)
	if err != nil {
		return EvalRuleSet{}, fmt.Errorf("store: upsert eval_rule_sets: %w", err)
	}
	return row, nil
}

// GetEvalRuleSet loads one eval_rule_sets row by id.
func (s *Store) GetEvalRuleSet(id string) (EvalRuleSet, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		id = EvalRuleSetDefaultID
	}
	var row EvalRuleSet
	err := s.db.QueryRow(`
		SELECT id, source_path, body_json, updated_at
		FROM eval_rule_sets WHERE id = ?
	`, id).Scan(&row.ID, &row.SourcePath, &row.BodyJSON, &row.UpdatedAt)
	if err == sql.ErrNoRows {
		return EvalRuleSet{}, fmt.Errorf("store: eval_rule_sets %q: %w", id, err)
	}
	if err != nil {
		return EvalRuleSet{}, fmt.Errorf("store: get eval_rule_sets: %w", err)
	}
	return row, nil
}
