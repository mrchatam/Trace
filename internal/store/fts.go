package store

import (
	"database/sql"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// FTSHit is one row from the fts_docs lexical index.
// Rank is SQLite FTS5 bm25 (lower is better). Score is -Rank for higher-is-better ranking.
type FTSHit struct {
	EntityType string
	EntityID   string
	Title      string
	Body       string
	Path       string
	SymbolName string
	SymbolKind string
	Rank       float64
	Score      float64
}

// SearchFTS runs an FTS5 MATCH query against fts_docs.
// Tokenizer is unicode61 (see schema/004_fts.sql). limit <= 0 defaults to 32; hard-capped at 64.
func (s *Store) SearchFTS(query string, limit int) ([]FTSHit, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = 32
	}
	if limit > 64 {
		limit = 64
	}
	match := sanitizeFTSQuery(query)
	if match == "" {
		return nil, nil
	}

	rows, err := s.db.Query(`
		SELECT entity_type, entity_id, title, body, path, symbol_name, symbol_kind,
		       bm25(fts_docs) AS rank
		FROM fts_docs
		WHERE fts_docs MATCH ?
		ORDER BY rank ASC, entity_type ASC, entity_id ASC
		LIMIT ?
	`, match, limit)
	if err != nil {
		return nil, fmt.Errorf("store: search fts: %w", err)
	}
	defer rows.Close()

	var out []FTSHit
	for rows.Next() {
		var h FTSHit
		if err := rows.Scan(
			&h.EntityType, &h.EntityID, &h.Title, &h.Body, &h.Path,
			&h.SymbolName, &h.SymbolKind, &h.Rank,
		); err != nil {
			return nil, fmt.Errorf("store: scan fts hit: %w", err)
		}
		h.Score = -h.Rank
		out = append(out, h)
	}
	return out, rows.Err()
}

// sanitizeFTSQuery turns free text into a safe FTS5 MATCH expression.
// Keep Unicode letters/numbers; every other rune (including `/`) is a separator.
// Remaining tokens are FTS5-quoted and joined with AND so operators cannot inject MATCH syntax.
func sanitizeFTSQuery(q string) string {
	var b strings.Builder
	b.Grow(len(q))
	for _, r := range q {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
		}
	}
	parts := strings.Fields(b.String())
	if len(parts) == 0 {
		return ""
	}
	quoted := make([]string, len(parts))
	for i, p := range parts {
		quoted[i] = `"` + strings.ReplaceAll(p, `"`, `""`) + `"`
	}
	return strings.Join(quoted, " AND ")
}

// ensureFTSPopulated runs RebuildFTS when the lexical index is empty but
// content tables already have rows (typical after applying 004_fts onto an
// existing DB). No-op when fts_docs already has rows or the DB is empty.
func (s *Store) ensureFTSPopulated() error {
	var ftsCount int
	if err := s.db.QueryRow(`SELECT COUNT(1) FROM fts_docs`).Scan(&ftsCount); err != nil {
		return fmt.Errorf("store: count fts_docs: %w", err)
	}
	if ftsCount > 0 {
		return nil
	}
	var content int
	if err := s.db.QueryRow(`
		SELECT
			(SELECT COUNT(1) FROM goals) +
			(SELECT COUNT(1) FROM tasks) +
			(SELECT COUNT(1) FROM decisions) +
			(SELECT COUNT(1) FROM assumptions) +
			(SELECT COUNT(1) FROM discoveries) +
			(SELECT COUNT(1) FROM plan_changes) +
			(SELECT COUNT(1) FROM claims) +
			(SELECT COUNT(1) FROM evidence) +
			(SELECT COUNT(1) FROM reviews) +
			(SELECT COUNT(1) FROM files) +
			(SELECT COUNT(1) FROM symbols) +
			(SELECT COUNT(1) FROM uncertainties) +
			(SELECT COUNT(1) FROM hypotheses) +
			(SELECT COUNT(1) FROM changes) +
			(SELECT COUNT(1) FROM regressions) +
			(SELECT COUNT(1) FROM reflections) +
			(SELECT COUNT(1) FROM baselines) +
			(SELECT COUNT(1) FROM outcome_results)
	`).Scan(&content); err != nil {
		return fmt.Errorf("store: count content for fts backfill: %w", err)
	}
	if content == 0 {
		return nil
	}
	return s.RebuildFTS()
}

// RebuildFTS clears and reindexes all FTS documents from entity/file/symbol tables.
func (s *Store) RebuildFTS() error {
	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: begin rebuild fts: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM fts_docs`); err != nil {
		return fmt.Errorf("store: clear fts: %w", err)
	}

	type textEntity struct {
		table, typ string
	}
	for _, e := range []textEntity{
		{"goals", "goal"},
		{"tasks", "task"},
		{"decisions", "decision"},
		{"assumptions", "assumption"},
		{"discoveries", "discovery"},
		{"plan_changes", "plan_change"},
		{"claims", "claim"},
		{"evidence", "evidence"},
		{"reviews", "review"},
	} {
		if err := rebuildTextTable(tx, e.table, e.typ); err != nil {
			return err
		}
	}

	rows, err := tx.Query(`SELECT id, path FROM files`)
	if err != nil {
		return fmt.Errorf("store: list files for fts: %w", err)
	}
	for rows.Next() {
		var id, path string
		if err := rows.Scan(&id, &path); err != nil {
			rows.Close()
			return fmt.Errorf("store: scan file for fts: %w", err)
		}
		if err := insertFTS(tx, "file", id, path, "", path, "", ""); err != nil {
			rows.Close()
			return err
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	symRows, err := tx.Query(`
		SELECT s.id, s.name, s.kind, f.path
		FROM symbols s
		JOIN files f ON f.id = s.file_id
	`)
	if err != nil {
		return fmt.Errorf("store: list symbols for fts: %w", err)
	}
	for symRows.Next() {
		var id, name, kind, path string
		if err := symRows.Scan(&id, &name, &kind, &path); err != nil {
			symRows.Close()
			return fmt.Errorf("store: scan symbol for fts: %w", err)
		}
		if err := insertFTS(tx, "symbol", id, name, "", path, name, kind); err != nil {
			symRows.Close()
			return err
		}
	}
	symRows.Close()
	if err := symRows.Err(); err != nil {
		return err
	}

	if err := rebuildUncertaintyTable(tx); err != nil {
		return err
	}
	if err := rebuildHypothesisTable(tx); err != nil {
		return err
	}
	if err := rebuildChangeTable(tx); err != nil {
		return err
	}
	if err := rebuildRegressionTable(tx); err != nil {
		return err
	}
	if err := rebuildReflectionTable(tx); err != nil {
		return err
	}
	if err := rebuildBaselineTable(tx); err != nil {
		return err
	}
	if err := rebuildOutcomeResultTable(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit rebuild fts: %w", err)
	}
	return nil
}

func rebuildTextTable(tx *sql.Tx, table, entityType string) error {
	q := fmt.Sprintf(`SELECT id, title, body FROM %s`, table)
	rows, err := tx.Query(q)
	if err != nil {
		return fmt.Errorf("store: list %s for fts: %w", table, err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, title, body string
		if err := rows.Scan(&id, &title, &body); err != nil {
			return fmt.Errorf("store: scan %s for fts: %w", table, err)
		}
		if err := insertFTS(tx, entityType, id, title, body, "", "", ""); err != nil {
			return err
		}
	}
	return rows.Err()
}

func insertFTS(tx *sql.Tx, entityType, entityID, title, body, path, symbolName, symbolKind string) error {
	_, err := tx.Exec(`
		INSERT INTO fts_docs(entity_type, entity_id, title, body, path, symbol_name, symbol_kind)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, entityType, entityID, title, body, path, symbolName, symbolKind)
	if err != nil {
		return fmt.Errorf("store: insert fts %s/%s: %w", entityType, entityID, err)
	}
	return nil
}

// SyncEntityFTS refreshes FTS rows for one causal/text entity (goal, task, …).
func (s *Store) SyncEntityFTS(entityType, entityID string) error {
	if entityType == "" || entityID == "" {
		return fmt.Errorf("store: sync fts: entity_type and entity_id required")
	}
	return s.runInTx(func(tx *sql.Tx) error {
		if _, err := tx.Exec(`DELETE FROM fts_docs WHERE entity_type = ? AND entity_id = ?`, entityType, entityID); err != nil {
			return fmt.Errorf("store: delete fts row: %w", err)
		}
		title, body, ok, err := loadEntityText(tx, entityType, entityID)
		if err != nil {
			return err
		}
		if ok {
			if err := insertFTS(tx, entityType, entityID, title, body, "", "", ""); err != nil {
				return err
			}
		}
		return nil
	})
}

func loadEntityText(tx *sql.Tx, entityType, entityID string) (title, body string, ok bool, err error) {
	switch entityType {
	case "goal", "task", "decision", "assumption", "discovery", "plan_change", "claim", "evidence", "review":
		table := map[string]string{
			"goal":        "goals",
			"task":        "tasks",
			"decision":    "decisions",
			"assumption":  "assumptions",
			"discovery":   "discoveries",
			"plan_change": "plan_changes",
			"claim":       "claims",
			"evidence":    "evidence",
			"review":      "reviews",
		}[entityType]
		err = tx.QueryRow(
			fmt.Sprintf(`SELECT title, body FROM %s WHERE id = ?`, table), entityID,
		).Scan(&title, &body)
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		if err != nil {
			return "", "", false, fmt.Errorf("store: load %s for fts: %w", entityType, err)
		}
		return title, body, true, nil
	case "uncertainty":
		var severity, kind string
		err = tx.QueryRow(`
			SELECT title, body, severity, kind FROM uncertainties WHERE id = ?
		`, entityID).Scan(&title, &body, &severity, &kind)
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		if err != nil {
			return "", "", false, fmt.Errorf("store: load uncertainty for fts: %w", err)
		}
		body = joinNonEmpty(body, severity, kind)
		return title, body, true, nil
	case "hypothesis":
		var status string
		err = tx.QueryRow(`
			SELECT title, body, status FROM hypotheses WHERE id = ?
		`, entityID).Scan(&title, &body, &status)
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		if err != nil {
			return "", "", false, fmt.Errorf("store: load hypothesis for fts: %w", err)
		}
		body = joinNonEmpty(body, status)
		return title, body, true, nil
	case "change":
		var reason, status string
		err = tx.QueryRow(`
			SELECT git_commit, reason, status FROM changes WHERE id = ?
		`, entityID).Scan(&title, &reason, &status)
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		if err != nil {
			return "", "", false, fmt.Errorf("store: load change for fts: %w", err)
		}
		body = joinNonEmpty(reason, status)
		return title, body, true, nil
	case "regression":
		var summary, attribution string
		err = tx.QueryRow(`
			SELECT dimension, summary, attribution FROM regressions WHERE id = ?
		`, entityID).Scan(&title, &summary, &attribution)
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		if err != nil {
			return "", "", false, fmt.Errorf("store: load regression for fts: %w", err)
		}
		body = joinNonEmpty(summary, attribution)
		return title, body, true, nil
	case "reflection":
		var assumptions string
		err = tx.QueryRow(`
			SELECT summary, invalidated_assumptions_json FROM reflections WHERE id = ?
		`, entityID).Scan(&title, &assumptions)
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		if err != nil {
			return "", "", false, fmt.Errorf("store: load reflection for fts: %w", err)
		}
		body = truncateText(assumptions, 512)
		return title, body, true, nil
	case "baseline":
		var label, gitCommit, scores string
		err = tx.QueryRow(`
			SELECT label, git_commit, scores_json FROM baselines WHERE id = ?
		`, entityID).Scan(&label, &gitCommit, &scores)
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		if err != nil {
			return "", "", false, fmt.Errorf("store: load baseline for fts: %w", err)
		}
		title = label
		if title == "" {
			title = gitCommit
		}
		body = truncateText(scores, 512)
		return title, body, true, nil
	case "outcome_result":
		var testName, kind, summary string
		err = tx.QueryRow(`
			SELECT test_name, kind, summary FROM outcome_results WHERE id = ?
		`, entityID).Scan(&testName, &kind, &summary)
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		if err != nil {
			return "", "", false, fmt.Errorf("store: load outcome_result for fts: %w", err)
		}
		title = testName
		if title == "" {
			title = kind
		}
		body = summary
		return title, body, true, nil
	case "effect":
		// Exact/Why only — no FTS row.
		return "", "", false, nil
	default:
		return "", "", false, fmt.Errorf("store: sync fts: unknown entity type %q", entityType)
	}
}

func truncateText(s string, maxRunes int) string {
	s = strings.TrimSpace(s)
	if maxRunes <= 0 || utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes])
}

func joinNonEmpty(parts ...string) string {
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, " ")
}

func rebuildUncertaintyTable(tx *sql.Tx) error {
	rows, err := tx.Query(`SELECT id, title, body, severity, kind FROM uncertainties`)
	if err != nil {
		return fmt.Errorf("store: list uncertainties for fts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, title, body, severity, kind string
		if err := rows.Scan(&id, &title, &body, &severity, &kind); err != nil {
			return fmt.Errorf("store: scan uncertainty for fts: %w", err)
		}
		body = joinNonEmpty(body, severity, kind)
		if err := insertFTS(tx, "uncertainty", id, title, body, "", "", ""); err != nil {
			return err
		}
	}
	return rows.Err()
}

func rebuildHypothesisTable(tx *sql.Tx) error {
	rows, err := tx.Query(`SELECT id, title, body, status FROM hypotheses`)
	if err != nil {
		return fmt.Errorf("store: list hypotheses for fts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, title, body, status string
		if err := rows.Scan(&id, &title, &body, &status); err != nil {
			return fmt.Errorf("store: scan hypothesis for fts: %w", err)
		}
		body = joinNonEmpty(body, status)
		if err := insertFTS(tx, "hypothesis", id, title, body, "", "", ""); err != nil {
			return err
		}
	}
	return rows.Err()
}

func rebuildChangeTable(tx *sql.Tx) error {
	rows, err := tx.Query(`SELECT id, git_commit, reason, status FROM changes`)
	if err != nil {
		return fmt.Errorf("store: list changes for fts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, gitCommit, reason, status string
		if err := rows.Scan(&id, &gitCommit, &reason, &status); err != nil {
			return fmt.Errorf("store: scan change for fts: %w", err)
		}
		if err := insertFTS(tx, "change", id, gitCommit, joinNonEmpty(reason, status), "", "", ""); err != nil {
			return err
		}
	}
	return rows.Err()
}

func rebuildRegressionTable(tx *sql.Tx) error {
	rows, err := tx.Query(`SELECT id, dimension, summary, attribution FROM regressions`)
	if err != nil {
		return fmt.Errorf("store: list regressions for fts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, dimension, summary, attribution string
		if err := rows.Scan(&id, &dimension, &summary, &attribution); err != nil {
			return fmt.Errorf("store: scan regression for fts: %w", err)
		}
		if err := insertFTS(tx, "regression", id, dimension, joinNonEmpty(summary, attribution), "", "", ""); err != nil {
			return err
		}
	}
	return rows.Err()
}

func rebuildReflectionTable(tx *sql.Tx) error {
	rows, err := tx.Query(`SELECT id, summary, invalidated_assumptions_json FROM reflections`)
	if err != nil {
		return fmt.Errorf("store: list reflections for fts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, summary, assumptions string
		if err := rows.Scan(&id, &summary, &assumptions); err != nil {
			return fmt.Errorf("store: scan reflection for fts: %w", err)
		}
		if err := insertFTS(tx, "reflection", id, summary, truncateText(assumptions, 512), "", "", ""); err != nil {
			return err
		}
	}
	return rows.Err()
}

func rebuildBaselineTable(tx *sql.Tx) error {
	rows, err := tx.Query(`SELECT id, label, git_commit, scores_json FROM baselines`)
	if err != nil {
		return fmt.Errorf("store: list baselines for fts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, label, gitCommit, scores string
		if err := rows.Scan(&id, &label, &gitCommit, &scores); err != nil {
			return fmt.Errorf("store: scan baseline for fts: %w", err)
		}
		title := label
		if title == "" {
			title = gitCommit
		}
		if err := insertFTS(tx, "baseline", id, title, truncateText(scores, 512), "", "", ""); err != nil {
			return err
		}
	}
	return rows.Err()
}

func rebuildOutcomeResultTable(tx *sql.Tx) error {
	rows, err := tx.Query(`SELECT id, test_name, kind, summary FROM outcome_results`)
	if err != nil {
		return fmt.Errorf("store: list outcome_results for fts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, testName, kind, summary string
		if err := rows.Scan(&id, &testName, &kind, &summary); err != nil {
			return fmt.Errorf("store: scan outcome_result for fts: %w", err)
		}
		title := testName
		if title == "" {
			title = kind
		}
		if err := insertFTS(tx, "outcome_result", id, title, summary, "", "", ""); err != nil {
			return err
		}
	}
	return rows.Err()
}

// SyncFileFTS refreshes FTS rows for a file path and its symbols.
func (s *Store) SyncFileFTS(path string) error {
	path = NormalizePath(path)
	if path == "" {
		return fmt.Errorf("store: sync file fts: path required")
	}
	f, err := s.GetFileByPath(path)
	if err != nil {
		return err
	}

	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: begin sync file fts: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM fts_docs WHERE entity_type = 'file' AND entity_id = ?`, f.ID); err != nil {
		return fmt.Errorf("store: delete file fts: %w", err)
	}
	// Drop symbols previously indexed for this path (by path column + type).
	if _, err := tx.Exec(`DELETE FROM fts_docs WHERE entity_type = 'symbol' AND path = ?`, path); err != nil {
		return fmt.Errorf("store: delete symbol fts: %w", err)
	}

	if err := insertFTS(tx, "file", f.ID, f.Path, "", f.Path, "", ""); err != nil {
		return err
	}

	symRows, err := tx.Query(`
		SELECT id, name, kind FROM symbols WHERE file_id = ?
	`, f.ID)
	if err != nil {
		return fmt.Errorf("store: list symbols for file fts: %w", err)
	}
	for symRows.Next() {
		var id, name, kind string
		if err := symRows.Scan(&id, &name, &kind); err != nil {
			symRows.Close()
			return fmt.Errorf("store: scan symbol for file fts: %w", err)
		}
		if err := insertFTS(tx, "symbol", id, name, "", path, name, kind); err != nil {
			symRows.Close()
			return err
		}
	}
	symRows.Close()
	if err := symRows.Err(); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit sync file fts: %w", err)
	}
	return nil
}
