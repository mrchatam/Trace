package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Code-graph relationship kinds persisted on code_edges.rel.
const (
	RelValidates             = "validates"
	RelContainsModule        = "contains_module"
	RelExportsAPI            = "exports_api"
	RelArchitecturalBoundary = "architectural_boundary"
	RelDependsOn             = "depends_on"
)

// Namespaces for deterministic graph ids (reindex must not churn FKs).
var (
	symbolIDNamespace = uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	edgeIDNamespace   = uuid.MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
)

// DeterministicSymbolID is file_id + name + kind + start_line (UUIDv5-style).
func DeterministicSymbolID(fileID, name, kind string, startLine int) string {
	return uuid.NewSHA1(symbolIDNamespace, []byte(fmt.Sprintf("%s\n%s\n%s\n%d", fileID, name, kind, startLine))).String()
}

func deterministicEdgeID(fromFileID, fromSymbolID, toFileID, toSymbolID, rel string) string {
	return uuid.NewSHA1(edgeIDNamespace, []byte(fmt.Sprintf("%s\n%s\n%s\n%s\n%s", fromFileID, fromSymbolID, toFileID, toSymbolID, rel))).String()
}

func validateCodeRel(rel string) error {
	switch rel {
	case RelValidates, RelContainsModule, RelExportsAPI, RelArchitecturalBoundary, RelDependsOn:
		return nil
	default:
		return fmt.Errorf("store: invalid code edge rel %q", rel)
	}
}

// FileRecord is a code-graph file stub (path identity + content hash; no source body).
type FileRecord struct {
	ID             string
	Path           string
	ContentHash    string
	GitOID         *string
	Language       *string
	IndexedAt      string
	Status         string
	SourceType     string
	Confidence     float64
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// Symbol is a symbol stub scoped to a file.
type Symbol struct {
	ID        string
	FileID    string
	Name      string
	Kind      string
	StartLine int
	EndLine   int
}

// Structural import edge provenance (Law 5 — inferred ≠ extracted). Exact strings.
const (
	ImportProvenanceExtracted = "EXTRACTED"
	ImportProvenanceInferred  = "INFERRED"
	ImportProvenanceAmbiguous = "AMBIGUOUS"
)

// Import is an import stub scoped to a file.
type Import struct {
	ID           string
	FileID       string
	ImportedPath string
	Symbol       *string
	// Provenance is EXTRACTED|INFERRED|AMBIGUOUS. Empty on write/read defaults to EXTRACTED.
	// Garbage values are rejected on write (DF-64); not silently coerced.
	Provenance string
}

// validateImportProvenance normalizes empty → EXTRACTED and rejects unknown values.
func validateImportProvenance(p string) (string, error) {
	if p == "" {
		return ImportProvenanceExtracted, nil
	}
	switch p {
	case ImportProvenanceExtracted, ImportProvenanceInferred, ImportProvenanceAmbiguous:
		return p, nil
	default:
		return "", fmt.Errorf("store: invalid import provenance %q (want EXTRACTED|INFERRED|AMBIGUOUS)", p)
	}
}

// normalizeImportProvenance maps empty DB values to EXTRACTED (pre-012 defense).
func normalizeImportProvenance(p string) string {
	if p == "" {
		return ImportProvenanceExtracted
	}
	return p
}

// UpsertFile upserts a file by repo-relative path (UNIQUE).
// gitOID may be nil. Returns the stored row.
func (s *Store) UpsertFile(path, contentHash string, gitOID *string) (FileRecord, error) {
	path = NormalizePath(path)
	if path == "" {
		return FileRecord{}, fmt.Errorf("store: upsert file: path required")
	}
	if contentHash == "" {
		return FileRecord{}, fmt.Errorf("store: upsert file: content_hash required")
	}

	now := nowRFC3339()
	var existingID string
	var createdAt string
	err := s.db.QueryRow(`SELECT id, created_at FROM files WHERE path = ?`, path).Scan(&existingID, &createdAt)
	switch {
	case err == sql.ErrNoRows:
		id := uuid.NewString()
		_, err = s.db.Exec(`
			INSERT INTO files(id, path, content_hash, git_oid, language, indexed_at, status, source_type, confidence, created_at, updated_at, last_verified_at)
			VALUES (?, ?, ?, ?, NULL, ?, ?, 'DETERMINISTIC', 1, ?, ?, NULL)
		`, id, path, contentHash, nullStr(gitOID), now, StatusActive, now, now)
		if err != nil {
			return FileRecord{}, fmt.Errorf("store: insert file: %w", err)
		}
		out, err := s.GetFileByPath(path)
		if err != nil {
			return FileRecord{}, err
		}
		if err := s.SyncFileFTS(path); err != nil {
			return FileRecord{}, err
		}
		return out, nil
	case err != nil:
		return FileRecord{}, fmt.Errorf("store: lookup file: %w", err)
	default:
		_, err = s.db.Exec(`
			UPDATE files SET content_hash = ?, git_oid = ?, indexed_at = ?, updated_at = ?
			WHERE id = ?
		`, contentHash, nullStr(gitOID), now, now, existingID)
		if err != nil {
			return FileRecord{}, fmt.Errorf("store: update file: %w", err)
		}
		out, err := s.GetFileByPath(path)
		if err != nil {
			return FileRecord{}, err
		}
		if err := s.SyncFileFTS(path); err != nil {
			return FileRecord{}, err
		}
		return out, nil
	}
}

// SetFileLanguage sets files.language for an existing path (no migration required).
func (s *Store) SetFileLanguage(path, language string) error {
	path = NormalizePath(path)
	if path == "" {
		return fmt.Errorf("store: set file language: path required")
	}
	if language == "" {
		return fmt.Errorf("store: set file language: language required")
	}
	res, err := s.db.Exec(`UPDATE files SET language = ?, updated_at = ? WHERE path = ?`, language, nowRFC3339(), path)
	if err != nil {
		return fmt.Errorf("store: set file language: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: set file language rows: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("store: set file language: file %q not found", path)
	}
	return nil
}

// GetFileByPath loads a file row by path.
func (s *Store) GetFileByPath(path string) (FileRecord, error) {
	path = NormalizePath(path)
	var f FileRecord
	var gitOID, language, lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, path, content_hash, git_oid, language, indexed_at, status, source_type, confidence, created_at, updated_at, last_verified_at
		FROM files WHERE path = ?
	`, path).Scan(
		&f.ID, &f.Path, &f.ContentHash, &gitOID, &language, &f.IndexedAt, &f.Status,
		&f.SourceType, &f.Confidence, &f.CreatedAt, &f.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return FileRecord{}, fmt.Errorf("store: file %q: %w", path, err)
	}
	if err != nil {
		return FileRecord{}, fmt.Errorf("store: get file: %w", err)
	}
	f.GitOID = nullStrPtr(gitOID)
	f.Language = nullStrPtr(language)
	f.LastVerifiedAt = nullStrPtr(lastVerified)
	return f, nil
}

// ListFilePaths returns all repo-relative paths currently in files, ordered by path.
func (s *Store) ListFilePaths() ([]string, error) {
	rows, err := s.db.Query(`SELECT path FROM files ORDER BY path ASC`)
	if err != nil {
		return nil, fmt.Errorf("store: list file paths: %w", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, fmt.Errorf("store: scan file path: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ListFilePathsByContentHash returns paths with the given content_hash, ordered by path.
func (s *Store) ListFilePathsByContentHash(hash string) ([]string, error) {
	if hash == "" {
		return nil, nil
	}
	rows, err := s.db.Query(`SELECT path FROM files WHERE content_hash = ? ORDER BY path ASC`, hash)
	if err != nil {
		return nil, fmt.Errorf("store: list file paths by content hash: %w", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, fmt.Errorf("store: scan file path by hash: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// DeleteFileByPath removes a file row and its FTS file/symbol docs.
// Symbols and imports CASCADE from files. Idempotent: missing path is success.
func (s *Store) DeleteFileByPath(path string) error {
	path = NormalizePath(path)
	if path == "" {
		return fmt.Errorf("store: delete file: path required")
	}

	f, err := s.GetFileByPath(path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: begin delete file: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM fts_docs WHERE entity_type = 'file' AND entity_id = ?`, f.ID); err != nil {
		return fmt.Errorf("store: delete file fts: %w", err)
	}
	if _, err := tx.Exec(`DELETE FROM fts_docs WHERE entity_type = 'symbol' AND path = ?`, path); err != nil {
		return fmt.Errorf("store: delete symbol fts: %w", err)
	}
	if _, err := tx.Exec(`DELETE FROM files WHERE path = ?`, path); err != nil {
		return fmt.Errorf("store: delete file: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit delete file: %w", err)
	}
	return nil
}

// ReplaceFileSymbols deletes symbols for the given path's file_id only, then inserts.
// Other files' symbols are untouched (DR-INCREMENTAL substrate).
func (s *Store) ReplaceFileSymbols(path string, symbols []Symbol) error {
	path = NormalizePath(path)
	f, err := s.GetFileByPath(path)
	if err != nil {
		return err
	}

	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: begin replace symbols: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Upsert first (deterministic ids). A blanket DELETE would ON DELETE SET NULL
	// every incoming to_symbol_id; two incoming validates from the same file then
	// collide on idx_code_edges_unique (NULL to_symbol_id collapses).
	keep := make([]string, 0, len(symbols))
	for _, sym := range symbols {
		id := sym.ID
		if id == "" {
			id = DeterministicSymbolID(f.ID, sym.Name, sym.Kind, sym.StartLine)
		}
		keep = append(keep, id)
		if _, err := tx.Exec(`
			INSERT INTO symbols(id, file_id, name, kind, start_line, end_line)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				file_id = excluded.file_id,
				name = excluded.name,
				kind = excluded.kind,
				start_line = excluded.start_line,
				end_line = excluded.end_line
		`, id, f.ID, sym.Name, sym.Kind, sym.StartLine, sym.EndLine); err != nil {
			return fmt.Errorf("store: insert symbol: %w", err)
		}
	}

	leftover, err := leftoverSymbolIDs(tx, f.ID, keep)
	if err != nil {
		return err
	}
	if err := collapseEdgesTargetingSymbols(tx, leftover); err != nil {
		return err
	}

	if len(keep) == 0 {
		if _, err := tx.Exec(`DELETE FROM symbols WHERE file_id = ?`, f.ID); err != nil {
			return fmt.Errorf("store: delete symbols for file: %w", err)
		}
	} else if len(leftover) > 0 {
		args := make([]any, 0, 1+len(keep))
		args = append(args, f.ID)
		for _, id := range keep {
			args = append(args, id)
		}
		q := `DELETE FROM symbols WHERE file_id = ? AND id NOT IN (` + sqlPlaceholders(len(keep)) + `)`
		if _, err := tx.Exec(q, args...); err != nil {
			return fmt.Errorf("store: delete leftover symbols: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit replace symbols: %w", err)
	}
	return s.SyncFileFTS(path)
}

func leftoverSymbolIDs(tx *sql.Tx, fileID string, keep []string) ([]string, error) {
	var rows *sql.Rows
	var err error
	if len(keep) == 0 {
		rows, err = tx.Query(`SELECT id FROM symbols WHERE file_id = ?`, fileID)
	} else {
		args := make([]any, 0, 1+len(keep))
		args = append(args, fileID)
		for _, id := range keep {
			args = append(args, id)
		}
		q := `SELECT id FROM symbols WHERE file_id = ? AND id NOT IN (` + sqlPlaceholders(len(keep)) + `)`
		rows, err = tx.Query(q, args...)
	}
	if err != nil {
		return nil, fmt.Errorf("store: leftover symbols: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("store: scan leftover symbol: %w", err)
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

// collapseEdgesTargetingSymbols drops extras that would share
// (from_file_id, from_symbol_id, to_file_id, rel) after ON DELETE SET NULL.
func collapseEdgesTargetingSymbols(tx *sql.Tx, leftover []string) error {
	if len(leftover) == 0 {
		return nil
	}
	in := sqlPlaceholders(len(leftover))
	args := make([]any, 0, 2*len(leftover))
	for _, id := range leftover {
		args = append(args, id)
	}
	for _, id := range leftover {
		args = append(args, id)
	}
	q := `
		DELETE FROM code_edges WHERE id IN (
			SELECT id FROM (
				SELECT id FROM code_edges
				WHERE to_symbol_id IN (` + in + `)
				  AND id NOT IN (
					SELECT MIN(id) FROM code_edges
					WHERE to_symbol_id IN (` + in + `)
					GROUP BY from_file_id, IFNULL(from_symbol_id, ''), to_file_id, rel
				  )
			)
		)`
	if _, err := tx.Exec(q, args...); err != nil {
		return fmt.Errorf("store: collapse edges before symbol delete: %w", err)
	}
	return nil
}

func sqlPlaceholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.TrimRight(strings.Repeat("?,", n), ",")
}

// GetSymbolByID loads a symbol by id and returns its owning file path (JOIN files).
// Misses return sql.ErrNoRows wrapped like other Get* helpers (DF-49).
func (s *Store) GetSymbolByID(id string) (Symbol, string, error) {
	if id == "" {
		return Symbol{}, "", fmt.Errorf("store: get symbol by id: id required")
	}
	var sym Symbol
	var path string
	err := s.db.QueryRow(`
		SELECT s.id, s.file_id, s.name, s.kind, s.start_line, s.end_line, f.path
		FROM symbols s
		JOIN files f ON f.id = s.file_id
		WHERE s.id = ?
	`, id).Scan(&sym.ID, &sym.FileID, &sym.Name, &sym.Kind, &sym.StartLine, &sym.EndLine, &path)
	if err == sql.ErrNoRows {
		return Symbol{}, "", fmt.Errorf("store: symbol id %q: %w", id, err)
	}
	if err != nil {
		return Symbol{}, "", fmt.Errorf("store: get symbol by id: %w", err)
	}
	return sym, path, nil
}

// ListSymbolsByPath returns symbols for a file path.
func (s *Store) ListSymbolsByPath(path string) ([]Symbol, error) {
	path = NormalizePath(path)
	f, err := s.GetFileByPath(path)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`
		SELECT id, file_id, name, kind, start_line, end_line
		FROM symbols WHERE file_id = ?
		ORDER BY start_line ASC, name ASC
	`, f.ID)
	if err != nil {
		return nil, fmt.Errorf("store: list symbols: %w", err)
	}
	defer rows.Close()

	var out []Symbol
	for rows.Next() {
		var sym Symbol
		if err := rows.Scan(&sym.ID, &sym.FileID, &sym.Name, &sym.Kind, &sym.StartLine, &sym.EndLine); err != nil {
			return nil, fmt.Errorf("store: scan symbol: %w", err)
		}
		out = append(out, sym)
	}
	return out, rows.Err()
}

// ReplaceFileImports deletes imports for the given path's file_id only, then inserts.
func (s *Store) ReplaceFileImports(path string, imports []Import) error {
	path = NormalizePath(path)
	f, err := s.GetFileByPath(path)
	if err != nil {
		return err
	}

	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: begin replace imports: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM imports WHERE file_id = ?`, f.ID); err != nil {
		return fmt.Errorf("store: delete imports for file: %w", err)
	}

	for _, imp := range imports {
		id := imp.ID
		if id == "" {
			id = uuid.NewString()
		}
		prov, err := validateImportProvenance(imp.Provenance)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(`
			INSERT INTO imports(id, file_id, imported_path, symbol, provenance)
			VALUES (?, ?, ?, ?, ?)
		`, id, f.ID, imp.ImportedPath, nullStr(imp.Symbol), prov); err != nil {
			return fmt.Errorf("store: insert import: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit replace imports: %w", err)
	}
	return nil
}

// ListImportsByPath returns imports for a file path.
func (s *Store) ListImportsByPath(path string) ([]Import, error) {
	path = NormalizePath(path)
	f, err := s.GetFileByPath(path)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`
		SELECT id, file_id, imported_path, symbol, provenance
		FROM imports WHERE file_id = ?
		ORDER BY imported_path ASC
	`, f.ID)
	if err != nil {
		return nil, fmt.Errorf("store: list imports: %w", err)
	}
	defer rows.Close()

	var out []Import
	for rows.Next() {
		var imp Import
		var sym sql.NullString
		if err := rows.Scan(&imp.ID, &imp.FileID, &imp.ImportedPath, &sym, &imp.Provenance); err != nil {
			return nil, fmt.Errorf("store: scan import: %w", err)
		}
		imp.Symbol = nullStrPtr(sym)
		imp.Provenance = normalizeImportProvenance(imp.Provenance)
		out = append(out, imp)
	}
	return out, rows.Err()
}

// ImportEdge is one imports row joined to its importer file path (no migration).
type ImportEdge struct {
	ImporterPath string
	ImportedPath string
	Provenance   string
}

// ListImportEdges returns all structural import edges with importer paths,
// ordered by importer path then imported_path. Used for reverse (incoming) walks.
func (s *Store) ListImportEdges() ([]ImportEdge, error) {
	rows, err := s.db.Query(`
		SELECT f.path, i.imported_path, i.provenance
		FROM imports i
		JOIN files f ON f.id = i.file_id
		ORDER BY f.path ASC, i.imported_path ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list import edges: %w", err)
	}
	defer rows.Close()

	var out []ImportEdge
	for rows.Next() {
		var e ImportEdge
		if err := rows.Scan(&e.ImporterPath, &e.ImportedPath, &e.Provenance); err != nil {
			return nil, fmt.Errorf("store: scan import edge: %w", err)
		}
		e.Provenance = normalizeImportProvenance(e.Provenance)
		out = append(out, e)
	}
	return out, rows.Err()
}

// CodeEdge is a structural relationship between files/symbols (no source body).
type CodeEdge struct {
	ID           string
	FromFileID   string
	FromSymbolID *string
	ToFileID     string
	ToSymbolID   *string
	Rel          string
	Provenance   string
}

func derefID(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func (s *Store) insertCodeEdge(tx interface {
	Exec(query string, args ...any) (sql.Result, error)
}, e CodeEdge) error {
	if err := validateCodeRel(e.Rel); err != nil {
		return err
	}
	prov, err := validateImportProvenance(e.Provenance)
	if err != nil {
		return err
	}
	id := e.ID
	if id == "" {
		id = deterministicEdgeID(e.FromFileID, derefID(e.FromSymbolID), e.ToFileID, derefID(e.ToSymbolID), e.Rel)
	}
	_, err = tx.Exec(`
		INSERT INTO code_edges(id, from_file_id, from_symbol_id, to_file_id, to_symbol_id, rel, provenance)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, e.FromFileID, nullStr(e.FromSymbolID), e.ToFileID, nullStr(e.ToSymbolID), e.Rel, prov)
	if err != nil {
		return fmt.Errorf("store: insert code edge: %w", err)
	}
	return nil
}

// ReplaceFileEdges deletes outgoing edges (from_file_id of path) only, then inserts.
// Incoming edges and other files' outgoing edges are untouched (DR-INCREMENTAL).
func (s *Store) ReplaceFileEdges(path string, edges []CodeEdge) error {
	path = NormalizePath(path)
	f, err := s.GetFileByPath(path)
	if err != nil {
		return err
	}

	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: begin replace edges: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM code_edges WHERE from_file_id = ?`, f.ID); err != nil {
		return fmt.Errorf("store: delete edges for file: %w", err)
	}

	for _, e := range edges {
		e.FromFileID = f.ID
		if err := s.insertCodeEdge(tx, e); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit replace edges: %w", err)
	}
	return nil
}

// UpsertFilePairEdges replaces edges of rel from fromPath to toPath only.
// Other outgoing edges of fromPath are kept (target IndexFile incoming validates).
func (s *Store) UpsertFilePairEdges(fromPath, toPath, rel string, edges []CodeEdge) error {
	fromPath = NormalizePath(fromPath)
	toPath = NormalizePath(toPath)
	if err := validateCodeRel(rel); err != nil {
		return err
	}
	fromFile, err := s.GetFileByPath(fromPath)
	if err != nil {
		return err
	}
	toFile, err := s.GetFileByPath(toPath)
	if err != nil {
		return err
	}

	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: begin upsert pair edges: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(
		`DELETE FROM code_edges WHERE from_file_id = ? AND to_file_id = ? AND rel = ?`,
		fromFile.ID, toFile.ID, rel,
	); err != nil {
		return fmt.Errorf("store: delete pair edges: %w", err)
	}

	for _, e := range edges {
		e.FromFileID = fromFile.ID
		e.ToFileID = toFile.ID
		e.Rel = rel
		if err := s.insertCodeEdge(tx, e); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit upsert pair edges: %w", err)
	}
	return nil
}

// ListEdgesByFile returns outgoing code_edges for path (from_file_id).
func (s *Store) ListEdgesByFile(path string) ([]CodeEdge, error) {
	path = NormalizePath(path)
	f, err := s.GetFileByPath(path)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`
		SELECT id, from_file_id, from_symbol_id, to_file_id, to_symbol_id, rel, provenance
		FROM code_edges WHERE from_file_id = ?
		ORDER BY rel ASC, to_file_id ASC, id ASC
	`, f.ID)
	if err != nil {
		return nil, fmt.Errorf("store: list edges: %w", err)
	}
	defer rows.Close()
	return scanCodeEdges(rows)
}

// ListValidatesForSymbol returns validates edges whose to_symbol_id matches.
func (s *Store) ListValidatesForSymbol(symbolID string) ([]CodeEdge, error) {
	if symbolID == "" {
		return nil, fmt.Errorf("store: list validates: symbol id required")
	}
	rows, err := s.db.Query(`
		SELECT id, from_file_id, from_symbol_id, to_file_id, to_symbol_id, rel, provenance
		FROM code_edges WHERE rel = ? AND to_symbol_id = ?
		ORDER BY from_file_id ASC, id ASC
	`, RelValidates, symbolID)
	if err != nil {
		return nil, fmt.Errorf("store: list validates: %w", err)
	}
	defer rows.Close()
	return scanCodeEdges(rows)
}

// ListValidatesForFile returns validates edges whose to_file_id matches (any to_symbol_id).
func (s *Store) ListValidatesForFile(fileID string) ([]CodeEdge, error) {
	if fileID == "" {
		return nil, fmt.Errorf("store: list validates: file id required")
	}
	rows, err := s.db.Query(`
		SELECT id, from_file_id, from_symbol_id, to_file_id, to_symbol_id, rel, provenance
		FROM code_edges WHERE rel = ? AND to_file_id = ?
		ORDER BY to_symbol_id ASC, from_file_id ASC, id ASC
	`, RelValidates, fileID)
	if err != nil {
		return nil, fmt.Errorf("store: list validates: %w", err)
	}
	defer rows.Close()
	return scanCodeEdges(rows)
}

func scanCodeEdges(rows *sql.Rows) ([]CodeEdge, error) {
	var out []CodeEdge
	for rows.Next() {
		var e CodeEdge
		var fromSym, toSym sql.NullString
		if err := rows.Scan(&e.ID, &e.FromFileID, &fromSym, &e.ToFileID, &toSym, &e.Rel, &e.Provenance); err != nil {
			return nil, fmt.Errorf("store: scan code edge: %w", err)
		}
		e.FromSymbolID = nullStrPtr(fromSym)
		e.ToSymbolID = nullStrPtr(toSym)
		e.Provenance = normalizeImportProvenance(e.Provenance)
		out = append(out, e)
	}
	return out, rows.Err()
}

// ListExports returns outgoing exports_api edges for path.
func (s *Store) ListExports(path string) ([]CodeEdge, error) {
	path = NormalizePath(path)
	f, err := s.GetFileByPath(path)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`
		SELECT id, from_file_id, from_symbol_id, to_file_id, to_symbol_id, rel, provenance
		FROM code_edges WHERE from_file_id = ? AND rel = ?
		ORDER BY to_symbol_id ASC, id ASC
	`, f.ID, RelExportsAPI)
	if err != nil {
		return nil, fmt.Errorf("store: list exports: %w", err)
	}
	defer rows.Close()
	return scanCodeEdges(rows)
}

// ListModuleContents returns contains_module edges for a file path, or for every
// file whose directory is dirOrPath (directory is module identity; no modules table).
func (s *Store) ListModuleContents(dirOrPath string) ([]CodeEdge, error) {
	dirOrPath = strings.TrimSuffix(NormalizePath(dirOrPath), "/")
	if dirOrPath == "" {
		dirOrPath = "."
	}
	if _, err := s.GetFileByPath(dirOrPath); err == nil {
		return s.listContainsByFile(dirOrPath)
	}
	return s.listContainsByDir(dirOrPath)
}

func (s *Store) listContainsByFile(path string) ([]CodeEdge, error) {
	f, err := s.GetFileByPath(path)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`
		SELECT id, from_file_id, from_symbol_id, to_file_id, to_symbol_id, rel, provenance
		FROM code_edges WHERE from_file_id = ? AND rel = ?
		ORDER BY to_symbol_id ASC, id ASC
	`, f.ID, RelContainsModule)
	if err != nil {
		return nil, fmt.Errorf("store: list module contents: %w", err)
	}
	defer rows.Close()
	return scanCodeEdges(rows)
}

func (s *Store) listContainsByDir(dir string) ([]CodeEdge, error) {
	rows, err := s.db.Query(`
		SELECT e.id, e.from_file_id, e.from_symbol_id, e.to_file_id, e.to_symbol_id, e.rel, e.provenance, f.path
		FROM code_edges e
		JOIN files f ON f.id = e.from_file_id
		WHERE e.rel = ?
		ORDER BY f.path ASC, e.id ASC
	`, RelContainsModule)
	if err != nil {
		return nil, fmt.Errorf("store: list module contents: %w", err)
	}
	defer rows.Close()

	var out []CodeEdge
	for rows.Next() {
		var e CodeEdge
		var fromSym, toSym sql.NullString
		var p string
		if err := rows.Scan(&e.ID, &e.FromFileID, &fromSym, &e.ToFileID, &toSym, &e.Rel, &e.Provenance, &p); err != nil {
			return nil, fmt.Errorf("store: scan module contents: %w", err)
		}
		if fileDir(p) != dir {
			continue
		}
		e.FromSymbolID = nullStrPtr(fromSym)
		e.ToSymbolID = nullStrPtr(toSym)
		e.Provenance = normalizeImportProvenance(e.Provenance)
		out = append(out, e)
	}
	return out, rows.Err()
}

func fileDir(p string) string {
	p = NormalizePath(p)
	i := strings.LastIndex(p, "/")
	if i < 0 {
		return "."
	}
	return p[:i]
}

// ArchitectureLayerPrefix is the TO-file path prefix for layer identity stubs.
// Outgoing architectural_boundary always starts at the indexed source file
// (from_file_id); these stubs are never IndexFile'd and never carry outgoing rels.
const ArchitectureLayerPrefix = "architecture/"

// ArchitectureLayerPath is the repo-relative identity path for a named layer.
func ArchitectureLayerPath(layer string) string {
	layer = strings.Trim(NormalizePath(layer), "/")
	if layer == "" {
		return ""
	}
	return ArchitectureLayerPrefix + layer
}

// LayerNameFromIdentityPath extracts the layer name from an architecture/ TO path.
func LayerNameFromIdentityPath(path string) string {
	path = NormalizePath(path)
	if !strings.HasPrefix(path, ArchitectureLayerPrefix) {
		return ""
	}
	return strings.TrimPrefix(path, ArchitectureLayerPrefix)
}

// ListArchitecturalBoundaries returns outgoing architectural_boundary edges for path.
func (s *Store) ListArchitecturalBoundaries(path string) ([]CodeEdge, error) {
	path = NormalizePath(path)
	f, err := s.GetFileByPath(path)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`
		SELECT id, from_file_id, from_symbol_id, to_file_id, to_symbol_id, rel, provenance
		FROM code_edges WHERE from_file_id = ? AND rel = ?
		ORDER BY to_file_id ASC, id ASC
	`, f.ID, RelArchitecturalBoundary)
	if err != nil {
		return nil, fmt.Errorf("store: list architectural boundaries: %w", err)
	}
	defer rows.Close()
	return scanCodeEdges(rows)
}

// FileLayer returns the layer name and provenance for a source file path
// (package→layer membership on that file's outgoing architectural_boundary).
func (s *Store) FileLayer(path string) (layer, provenance string, err error) {
	edges, err := s.ListArchitecturalBoundaries(path)
	if err != nil {
		return "", "", err
	}
	if len(edges) == 0 {
		return "", "", fmt.Errorf("store: no architectural_boundary for %s", path)
	}
	tf, err := s.GetFileByID(edges[0].ToFileID)
	if err != nil {
		return "", "", err
	}
	layer = LayerNameFromIdentityPath(tf.Path)
	if layer == "" {
		return "", "", fmt.Errorf("store: architectural_boundary to_file %q is not a layer identity", tf.Path)
	}
	return layer, edges[0].Provenance, nil
}

// CrossLayerImport is an observed imports row whose importer and resolved
// target live in different architectural layers. The import graph is not rewritten.
type CrossLayerImport struct {
	ImporterPath string
	ImportedPath string
	FromLayer    string
	ToLayer      string
	Provenance   string
}

// ListCrossLayerImports joins stored imports with FileLayer membership.
// Unresolved import targets are omitted (still visible via ListImportEdges).
func (s *Store) ListCrossLayerImports() ([]CrossLayerImport, error) {
	layers, err := s.fileLayerMap()
	if err != nil {
		return nil, err
	}
	imps, err := s.ListImportEdges()
	if err != nil {
		return nil, err
	}
	var out []CrossLayerImport
	for _, im := range imps {
		fromLayer, ok := layers[im.ImporterPath]
		if !ok {
			continue
		}
		toLayer, ok := resolveImportedLayer(im.ImportedPath, layers)
		if !ok || toLayer == fromLayer {
			continue
		}
		out = append(out, CrossLayerImport{
			ImporterPath: im.ImporterPath,
			ImportedPath: im.ImportedPath,
			FromLayer:    fromLayer,
			ToLayer:      toLayer,
			Provenance:   im.Provenance,
		})
	}
	return out, nil
}

func (s *Store) fileLayerMap() (map[string]string, error) {
	rows, err := s.db.Query(`
		SELECT f.path, tf.path
		FROM code_edges e
		JOIN files f ON f.id = e.from_file_id
		JOIN files tf ON tf.id = e.to_file_id
		WHERE e.rel = ?
		ORDER BY f.path ASC, e.id ASC
	`, RelArchitecturalBoundary)
	if err != nil {
		return nil, fmt.Errorf("store: list file layers: %w", err)
	}
	defer rows.Close()
	out := make(map[string]string)
	for rows.Next() {
		var fromPath, toPath string
		if err := rows.Scan(&fromPath, &toPath); err != nil {
			return nil, fmt.Errorf("store: scan file layer: %w", err)
		}
		if _, exists := out[fromPath]; exists {
			continue
		}
		if layer := LayerNameFromIdentityPath(toPath); layer != "" {
			out[fromPath] = layer
		}
	}
	return out, rows.Err()
}

func resolveImportedLayer(imported string, layers map[string]string) (string, bool) {
	imported = strings.Trim(NormalizePath(imported), "/")
	if imported == "" {
		return "", false
	}
	for srcPath, layer := range layers {
		dir := fileDir(srcPath)
		if dir == "." {
			continue
		}
		if imported == dir || strings.HasSuffix(imported, "/"+dir) {
			return layer, true
		}
	}
	return "", false
}
