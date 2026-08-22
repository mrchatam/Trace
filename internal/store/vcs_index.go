package store

import (
	"database/sql"
	"fmt"
	"strings"
)

// Meta keys for the thin VCS index.
const (
	MetaVCSWatermark = "vcs_index_watermark"
)

// IndexedCommit is a thin commit row (no patch body).
type IndexedCommit struct {
	OID         string
	ParentOIDs  []string
	CommittedAt string
	Subject     string
	Seq         int64 // refresh order (oldest→newest); used for stable history sort
}

// IndexedPathChange is a path touched by a commit.
type IndexedPathChange struct {
	Path   string
	Status string
}

// GetMeta returns a vcs_meta value, or "" if unset.
func (s *Store) GetMeta(key string) (string, error) {
	var v string
	err := s.db.QueryRow(`SELECT value FROM vcs_meta WHERE key = ?`, key).Scan(&v)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("store: get meta %q: %w", key, err)
	}
	return v, nil
}

// SetMeta upserts a vcs_meta key/value.
func (s *Store) SetMeta(key, value string) error {
	if key == "" {
		return fmt.Errorf("store: set meta: key required")
	}
	_, err := s.db.Exec(`
		INSERT INTO vcs_meta(key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, value)
	if err != nil {
		return fmt.Errorf("store: set meta %q: %w", key, err)
	}
	return nil
}

// UpsertIndexedCommit inserts or replaces a thin commit row and its path list.
// Does not store patch/diff bodies.
func (s *Store) UpsertIndexedCommit(c IndexedCommit, paths []IndexedPathChange) error {
	if c.OID == "" {
		return fmt.Errorf("store: upsert indexed commit: oid required")
	}
	parents := strings.Join(c.ParentOIDs, " ")

	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: begin upsert indexed commit: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`
		INSERT INTO vcs_commits(oid, parent_oids, committed_at, subject, seq)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(oid) DO UPDATE SET
			parent_oids = excluded.parent_oids,
			committed_at = excluded.committed_at,
			subject = excluded.subject
	`, c.OID, parents, c.CommittedAt, c.Subject, c.Seq); err != nil {
		return fmt.Errorf("store: upsert vcs_commits: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM vcs_commit_paths WHERE commit_oid = ?`, c.OID); err != nil {
		return fmt.Errorf("store: clear commit paths: %w", err)
	}
	for _, p := range paths {
		path := NormalizePath(p.Path)
		if path == "" {
			continue
		}
		if _, err := tx.Exec(`
			INSERT INTO vcs_commit_paths(commit_oid, path, status)
			VALUES (?, ?, ?)
		`, c.OID, path, p.Status); err != nil {
			return fmt.Errorf("store: insert commit path: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit upsert indexed commit: %w", err)
	}
	return nil
}

// CountIndexedCommits returns the number of rows in vcs_commits.
func (s *Store) CountIndexedCommits() (int, error) {
	var n int
	if err := s.db.QueryRow(`SELECT COUNT(1) FROM vcs_commits`).Scan(&n); err != nil {
		return 0, fmt.Errorf("store: count indexed commits: %w", err)
	}
	return n, nil
}

// GetIndexedCommit loads a thin commit by OID.
func (s *Store) GetIndexedCommit(oid string) (IndexedCommit, error) {
	var c IndexedCommit
	var parents string
	err := s.db.QueryRow(`
		SELECT oid, parent_oids, committed_at, subject, seq FROM vcs_commits WHERE oid = ?
	`, oid).Scan(&c.OID, &parents, &c.CommittedAt, &c.Subject, &c.Seq)
	if err == sql.ErrNoRows {
		return IndexedCommit{}, fmt.Errorf("store: indexed commit %q: %w", oid, err)
	}
	if err != nil {
		return IndexedCommit{}, fmt.Errorf("store: get indexed commit: %w", err)
	}
	c.ParentOIDs = splitParents(parents)
	return c, nil
}

// ListIndexedCommitPaths returns path changes for a commit.
func (s *Store) ListIndexedCommitPaths(oid string) ([]IndexedPathChange, error) {
	rows, err := s.db.Query(`
		SELECT path, status FROM vcs_commit_paths
		WHERE commit_oid = ?
		ORDER BY path ASC
	`, oid)
	if err != nil {
		return nil, fmt.Errorf("store: list commit paths: %w", err)
	}
	defer rows.Close()

	var out []IndexedPathChange
	for rows.Next() {
		var p IndexedPathChange
		if err := rows.Scan(&p.Path, &p.Status); err != nil {
			return nil, fmt.Errorf("store: scan commit path: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ListIndexedHistory returns commits that touched path, newest first.
func (s *Store) ListIndexedHistory(path string, limit int) ([]IndexedCommit, error) {
	path = NormalizePath(path)
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.db.Query(`
		SELECT c.oid, c.parent_oids, c.committed_at, c.subject, c.seq
		FROM vcs_commits c
		INNER JOIN vcs_commit_paths p ON p.commit_oid = c.oid
		WHERE p.path = ?
		ORDER BY c.seq DESC
		LIMIT ?
	`, path, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list indexed history: %w", err)
	}
	defer rows.Close()
	return scanIndexedCommits(rows)
}

// ListIndexedCommitsBetween returns indexed commits whose OIDs appear in oids
// (caller supplies the ordered OID list from git). Missing OIDs are skipped.
func (s *Store) ListIndexedCommitsByOIDs(oids []string) ([]IndexedCommit, error) {
	out := make([]IndexedCommit, 0, len(oids))
	for _, oid := range oids {
		c, err := s.GetIndexedCommit(oid)
		if err != nil {
			continue
		}
		out = append(out, c)
	}
	return out, nil
}

func scanIndexedCommits(rows *sql.Rows) ([]IndexedCommit, error) {
	var out []IndexedCommit
	for rows.Next() {
		var c IndexedCommit
		var parents string
		if err := rows.Scan(&c.OID, &parents, &c.CommittedAt, &c.Subject, &c.Seq); err != nil {
			return nil, fmt.Errorf("store: scan indexed commit: %w", err)
		}
		c.ParentOIDs = splitParents(parents)
		out = append(out, c)
	}
	return out, rows.Err()
}

// ListIndexedCommitsSince returns indexed commits with seq strictly after sinceOID's seq.
// Empty sinceOID returns all commits ordered oldest→newest.
func (s *Store) ListIndexedCommitsSince(sinceOID string) ([]IndexedCommit, error) {
	sinceOID = strings.TrimSpace(sinceOID)
	var afterSeq int64
	if sinceOID != "" {
		c, err := s.GetIndexedCommit(sinceOID)
		if err != nil {
			return nil, err
		}
		afterSeq = c.Seq
	}
	rows, err := s.db.Query(`
		SELECT oid, parent_oids, committed_at, subject, seq
		FROM vcs_commits
		WHERE seq > ?
		ORDER BY seq ASC
	`, afterSeq)
	if err != nil {
		return nil, fmt.Errorf("store: list indexed commits since: %w", err)
	}
	defer rows.Close()
	return scanIndexedCommits(rows)
}

// NextCommitSeq returns max(seq)+1 for assigning refresh order.
func (s *Store) NextCommitSeq() (int64, error) {
	var max sql.NullInt64
	if err := s.db.QueryRow(`SELECT MAX(seq) FROM vcs_commits`).Scan(&max); err != nil {
		return 0, fmt.Errorf("store: next commit seq: %w", err)
	}
	if !max.Valid {
		return 1, nil
	}
	return max.Int64 + 1, nil
}

func splitParents(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Fields(s)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// HasBlobLikeColumns reports whether any user table has BLOB columns or
// suspicious body/content/diff/patch columns on VCS index tables.
// Used in tests to guard G1 (no source content in SQLite).
//
// FTS5 shadow tables (fts_*_data etc.) store inverted-index blocks as BLOB;
// those are lexical index structures, not source-file content, and are skipped.
func (s *Store) HasBlobLikeColumns() (bool, string, error) {
	rows, err := s.db.Query(`
		SELECT m.name, p.name, p.type
		FROM sqlite_master m
		JOIN pragma_table_info(m.name) p
		WHERE m.type = 'table' AND m.name NOT LIKE 'sqlite_%'
	`)
	if err != nil {
		return false, "", fmt.Errorf("store: audit columns: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table, col, colType string
		if err := rows.Scan(&table, &col, &colType); err != nil {
			return false, "", err
		}
		if isFTS5ShadowTable(table) {
			continue
		}
		upperType := strings.ToUpper(colType)
		if strings.Contains(upperType, "BLOB") {
			return true, table + "." + col + " (" + colType + ")", nil
		}
		if strings.HasPrefix(table, "vcs_") {
			lower := strings.ToLower(col)
			switch {
			case strings.Contains(lower, "diff"),
				strings.Contains(lower, "patch"),
				strings.Contains(lower, "body"),
				strings.Contains(lower, "content"),
				strings.Contains(lower, "blob"):
				return true, table + "." + col, nil
			}
		}
	}
	return false, "", rows.Err()
}

// isFTS5ShadowTable reports FTS5 auxiliary tables that hold index blocks (not source).
func isFTS5ShadowTable(name string) bool {
	// e.g. fts_docs_data, fts_docs_idx, fts_docs_config, fts_docs_docsize, fts_docs_content
	for _, suf := range []string{"_data", "_idx", "_content", "_docsize", "_config"} {
		if strings.HasSuffix(name, suf) && strings.Contains(name, "fts") {
			return true
		}
	}
	return false
}
