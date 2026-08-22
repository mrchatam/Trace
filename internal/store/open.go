package store

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

const (
	traceDirName = ".trace"
	dbFileName   = "trace.db"

	// strayRootDBWarn is printed once per openStore when <root>/trace.db is a regular file.
	// That path is never the Trace store; live DB is always .trace/trace.db.
	strayRootDBWarn = "trace: warning: project-root trace.db exists but is not the Trace store; using .trace/trace.db. Do not open or create a root trace.db (agents: use CLI/MCP).\n"
)

// warnWriter receives stray-root-trace.db warnings. Default os.Stderr; tests may redirect.
var warnWriter io.Writer = os.Stderr

// sqlExecutor is the database/sql surface used by Store query methods.
// Both *sql.DB and *sql.Tx implement it so WithTx can run store methods in a transaction.
type sqlExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// Store is a per-project SQLite handle bound to projectRoot/.trace/trace.db.
// Open holds an exclusive advisory lock on .trace/trace.lock until Close.
type Store struct {
	conn        *sql.DB
	db          sqlExecutor
	projectRoot string
	projectID   string
	lock        *flockHandle
}

// Open creates <projectRoot>/.trace/ if needed, acquires exclusive trace.lock,
// checks optional .trace/access.token against TRACE_ACCESS_TOKEN, opens
// trace.db, applies migrations, and ensures a projects row for this root.
// A second Open on the same abs root fails with ErrLocked until Close.
//
// Auth order: mkdir → lock → token gate → db open → migrate → ensureProject.
func Open(projectRoot string) (*Store, error) {
	return openStore(projectRoot, false)
}

// OpenExisting opens an already-initialized store at
// <projectRoot>/.trace/trace.db. Unlike Open, it does not create .trace/ or
// trace.db. The database must already exist as a regular file (Abs only; no
// parent .trace walk-up). Missing file or an empty .trace/ directory returns
// ErrNotInitialized — never call Open on a Stat miss.
func OpenExisting(projectRoot string) (*Store, error) {
	absRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("store: resolve project root: %w", err)
	}
	dbPath := filepath.Join(absRoot, traceDirName, dbFileName)
	fi, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: missing %s", ErrNotInitialized, dbPath)
		}
		return nil, fmt.Errorf("store: stat db: %w", err)
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("%w: missing %s", ErrNotInitialized, dbPath)
	}
	return openStore(absRoot, false)
}

func openStore(projectRoot string, rebindRoot bool) (*Store, error) {
	absRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("store: resolve project root: %w", err)
	}

	warnIfStrayRootTraceDB(absRoot)

	traceDir := filepath.Join(absRoot, traceDirName)
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		return nil, fmt.Errorf("store: mkdir .trace: %w", err)
	}

	lock, err := acquireTraceLock(traceDir)
	if err != nil {
		return nil, err
	}

	if err := checkAccessToken(traceDir); err != nil {
		_ = lock.release()
		return nil, err
	}

	dbPath := filepath.Join(traceDir, dbFileName)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		_ = lock.release()
		return nil, fmt.Errorf("store: open db: %w", err)
	}
	// Single-writer local store; keep one connection to avoid lock surprises.
	db.SetMaxOpenConns(1)

	fail := func(err error) (*Store, error) {
		_ = db.Close()
		_ = lock.release()
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return fail(fmt.Errorf("store: enable foreign_keys: %w", err))
	}

	if err := migrate(db); err != nil {
		return fail(err)
	}

	s := &Store{conn: db, db: db, projectRoot: absRoot, lock: lock}
	if rebindRoot {
		if err := s.rebindProjectRoot(); err != nil {
			return fail(err)
		}
	}
	if err := s.ensureProject(); err != nil {
		return fail(err)
	}
	// After mig 004 (or a wiped index), Sync* only covers new Upserts — backfill
	// once when fts_docs is empty but entity/file/symbol rows already exist.
	if err := s.ensureFTSPopulated(); err != nil {
		return fail(err)
	}
	return s, nil
}

// warnIfStrayRootTraceDB writes a non-fatal warning when <absRoot>/trace.db exists
// as a regular file. It never opens, deletes, or renames that path.
func warnIfStrayRootTraceDB(absRoot string) {
	fi, err := os.Stat(filepath.Join(absRoot, dbFileName))
	if err != nil || !fi.Mode().IsRegular() {
		return
	}
	_, _ = io.WriteString(warnWriter, strayRootDBWarn)
}

// Close releases the database handle and the exclusive .trace/trace.lock.
func (s *Store) Close() error {
	if s == nil {
		return nil
	}
	var err error
	if s.conn != nil {
		err = s.conn.Close()
		s.conn = nil
		s.db = nil
	}
	if s.lock != nil {
		if lerr := s.lock.release(); err == nil {
			err = lerr
		}
		s.lock = nil
	}
	return err
}

// DBPath returns the absolute path to trace.db (for tests/diagnostics).
func (s *Store) DBPath() string {
	return filepath.Join(s.projectRoot, traceDirName, dbFileName)
}

// ProjectRoot returns the absolute project root this store is bound to.
func (s *Store) ProjectRoot() string {
	return s.projectRoot
}

// ProjectID returns the bound projects.id for this store.
func (s *Store) ProjectID() string {
	return s.projectID
}

// rebindProjectRoot updates the single projects.root_path to the current Abs
// root (restore path). Fail-closed if more than one projects row exists.
func (s *Store) rebindProjectRoot() error {
	var n int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM projects`).Scan(&n); err != nil {
		return fmt.Errorf("store: count projects for rebind: %w", err)
	}
	if n == 0 {
		return nil
	}
	if n > 1 {
		return fmt.Errorf("store: rebind: expected at most one projects row, found %d", n)
	}
	if _, err := s.db.Exec(`UPDATE projects SET root_path = ?`, s.projectRoot); err != nil {
		return fmt.Errorf("store: rebind root_path: %w", err)
	}
	return nil
}

func (s *Store) ensureProject() error {
	var id string
	err := s.db.QueryRow(`SELECT id FROM projects WHERE root_path = ?`, s.projectRoot).Scan(&id)
	if err == nil {
		s.projectID = id
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("store: lookup project: %w", err)
	}
	id = uuid.NewString()
	now := nowRFC3339()
	if _, err := s.db.Exec(
		`INSERT INTO projects(id, root_path, created_at) VALUES (?, ?, ?)`,
		id, s.projectRoot, now,
	); err != nil {
		return fmt.Errorf("store: insert project: %w", err)
	}
	s.projectID = id
	return nil
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func nowRFC3339Nano() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
