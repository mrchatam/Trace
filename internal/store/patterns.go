package store

import (
	"database/sql"
	"fmt"
	"strings"
)

// ChangePattern aggregates outcome counts for a (change_kind, outcome_kind) bucket.
type ChangePattern struct {
	ChangeKind    string
	OutcomeKind   string
	CountPositive int
	CountNegative int
	LastSeen      string
}

// UpsertChangePattern inserts or replaces one change_patterns row.
func (s *Store) UpsertChangePattern(row ChangePattern) error {
	if strings.TrimSpace(row.ChangeKind) == "" || strings.TrimSpace(row.OutcomeKind) == "" {
		return fmt.Errorf("store: upsert change pattern: change_kind and outcome_kind required")
	}
	_, err := s.db.Exec(`
		INSERT INTO change_patterns(change_kind, outcome_kind, count_positive, count_negative, last_seen)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(change_kind, outcome_kind) DO UPDATE SET
			count_positive = excluded.count_positive,
			count_negative = excluded.count_negative,
			last_seen = excluded.last_seen
	`, row.ChangeKind, row.OutcomeKind, row.CountPositive, row.CountNegative, row.LastSeen)
	if err != nil {
		return fmt.Errorf("store: upsert change pattern: %w", err)
	}
	return nil
}

// ReplaceChangePatterns replaces all change_patterns rows (deterministic full rebuild).
func (s *Store) ReplaceChangePatterns(rows []ChangePattern) error {
	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: replace change patterns: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM change_patterns`); err != nil {
		return fmt.Errorf("store: replace change patterns: delete: %w", err)
	}
	for _, r := range rows {
		if _, err := tx.Exec(`
			INSERT INTO change_patterns(change_kind, outcome_kind, count_positive, count_negative, last_seen)
			VALUES (?, ?, ?, ?, ?)
		`, r.ChangeKind, r.OutcomeKind, r.CountPositive, r.CountNegative, r.LastSeen); err != nil {
			return fmt.Errorf("store: replace change patterns: insert: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: replace change patterns: commit: %w", err)
	}
	return nil
}

// ListChangePatterns returns stored pattern rows ordered by change_kind, outcome_kind.
func (s *Store) ListChangePatterns(limit int) ([]ChangePattern, error) {
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}
	rows, err := s.db.Query(`
		SELECT change_kind, outcome_kind, count_positive, count_negative, last_seen
		FROM change_patterns
		ORDER BY change_kind ASC, outcome_kind ASC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list change patterns: %w", err)
	}
	defer rows.Close()
	var out []ChangePattern
	for rows.Next() {
		var p ChangePattern
		if err := rows.Scan(&p.ChangeKind, &p.OutcomeKind, &p.CountPositive, &p.CountNegative, &p.LastSeen); err != nil {
			return nil, fmt.Errorf("store: scan change pattern: %w", err)
		}
		out = append(out, p)
	}
	if out == nil {
		out = []ChangePattern{}
	}
	return out, rows.Err()
}

// ListChangePatternsByKind returns pattern rows for one change_kind.
func (s *Store) ListChangePatternsByKind(changeKind string, limit int) ([]ChangePattern, error) {
	changeKind = strings.TrimSpace(changeKind)
	if changeKind == "" {
		return nil, fmt.Errorf("store: list change patterns by kind: change_kind required")
	}
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}
	rows, err := s.db.Query(`
		SELECT change_kind, outcome_kind, count_positive, count_negative, last_seen
		FROM change_patterns
		WHERE change_kind = ?
		ORDER BY outcome_kind ASC
		LIMIT ?
	`, changeKind, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list change patterns by kind: %w", err)
	}
	defer rows.Close()
	var out []ChangePattern
	for rows.Next() {
		var p ChangePattern
		if err := rows.Scan(&p.ChangeKind, &p.OutcomeKind, &p.CountPositive, &p.CountNegative, &p.LastSeen); err != nil {
			return nil, fmt.Errorf("store: scan change pattern: %w", err)
		}
		out = append(out, p)
	}
	if out == nil {
		out = []ChangePattern{}
	}
	return out, rows.Err()
}

// ListChangesByPathPrefix returns non-SUPERSEDED changes with a path matching prefix, newest first.
func (s *Store) ListChangesByPathPrefix(prefix string, limit int) ([]Change, error) {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return nil, fmt.Errorf("store: list changes by path prefix: prefix required")
	}
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}
	like := prefix
	if !strings.HasSuffix(like, "%") {
		like += "%"
	}
	rows, err := s.db.Query(`
		SELECT DISTINCT c.id, c.task_id, c.git_commit, c.parent_change_id, c.actor, c.reason, c.status,
			c.source_type, c.confidence, c.created_at, c.updated_at, c.last_verified_at, c.rowid
		FROM changes c
		JOIN change_paths cp ON cp.change_id = c.id
		WHERE c.status != ? AND cp.path LIKE ?
		ORDER BY c.created_at DESC, c.rowid DESC
		LIMIT ?
	`, ChangeStatusSuperseded, like, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list changes by path prefix: %w", err)
	}
	defer rows.Close()
	return scanChanges(rows)
}

// ListActiveChanges returns all changes except SUPERSEDED, oldest first (for pattern refresh).
func (s *Store) ListActiveChanges() ([]Change, error) {
	rows, err := s.db.Query(`
		SELECT id, task_id, git_commit, parent_change_id, actor, reason, status,
			source_type, confidence, created_at, updated_at, last_verified_at, rowid
		FROM changes
		WHERE status != ?
		ORDER BY created_at ASC, rowid ASC
	`, ChangeStatusSuperseded)
	if err != nil {
		return nil, fmt.Errorf("store: list active changes: %w", err)
	}
	defer rows.Close()
	return scanChanges(rows)
}

func scanChanges(rows *sql.Rows) ([]Change, error) {
	var out []Change
	for rows.Next() {
		var c Change
		var lastVerified sql.NullString
		if err := rows.Scan(
			&c.ID, &c.TaskID, &c.GitCommit, &c.ParentChangeID, &c.Actor, &c.Reason, &c.Status,
			&c.SourceType, &c.Confidence, &c.CreatedAt, &c.UpdatedAt, &lastVerified, &c.RowID,
		); err != nil {
			return nil, fmt.Errorf("store: scan change: %w", err)
		}
		c.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, c)
	}
	if out == nil {
		out = []Change{}
	}
	return out, rows.Err()
}
