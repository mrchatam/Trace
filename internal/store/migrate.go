package store

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed schema/*.sql
var schemaFS embed.FS

const migrationTableSQL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY
);
`

// MigrationStatus is a read-only view of applied vs embedded migrations.
type MigrationStatus struct {
	AppliedVersions []int
	MaxApplied      int
	EmbedExpected   int
	PendingCount    int
}

// MigrationStatus reports applied schema versions and the embedded maximum.
func (s *Store) MigrationStatus() (MigrationStatus, error) {
	if s == nil || s.conn == nil {
		return MigrationStatus{}, fmt.Errorf("store: migration status: nil store")
	}
	return migrationStatus(s.conn)
}

type embedMig struct {
	version int
	name    string
}

func listEmbeddedMigrations() ([]embedMig, error) {
	entries, err := schemaFS.ReadDir("schema")
	if err != nil {
		return nil, fmt.Errorf("store: read schema dir: %w", err)
	}
	var migs []embedMig
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		ver, err := parseMigrationVersion(e.Name())
		if err != nil {
			return nil, err
		}
		migs = append(migs, embedMig{version: ver, name: e.Name()})
	}
	sort.Slice(migs, func(i, j int) bool { return migs[i].version < migs[j].version })
	return migs, nil
}

func embedExpectedMax() (int, error) {
	migs, err := listEmbeddedMigrations()
	if err != nil {
		return 0, err
	}
	if len(migs) == 0 {
		return 0, nil
	}
	return migs[len(migs)-1].version, nil
}

func migrationStatus(db *sql.DB) (MigrationStatus, error) {
	embedMax, err := embedExpectedMax()
	if err != nil {
		return MigrationStatus{}, err
	}
	rows, err := db.Query(`SELECT version FROM schema_migrations ORDER BY version`)
	if err != nil {
		return MigrationStatus{}, fmt.Errorf("store: list migrations: %w", err)
	}
	defer rows.Close()

	var applied []int
	maxApplied := 0
	appliedSet := map[int]struct{}{}
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return MigrationStatus{}, err
		}
		applied = append(applied, v)
		appliedSet[v] = struct{}{}
		if v > maxApplied {
			maxApplied = v
		}
	}
	if err := rows.Err(); err != nil {
		return MigrationStatus{}, err
	}

	migs, err := listEmbeddedMigrations()
	if err != nil {
		return MigrationStatus{}, err
	}
	pending := 0
	for _, m := range migs {
		if _, ok := appliedSet[m.version]; !ok {
			pending++
		}
	}
	return MigrationStatus{
		AppliedVersions: applied,
		MaxApplied:      maxApplied,
		EmbedExpected:   embedMax,
		PendingCount:    pending,
	}, nil
}

func migrate(db *sql.DB) error {
	if _, err := db.Exec(migrationTableSQL); err != nil {
		return fmt.Errorf("store: create schema_migrations: %w", err)
	}

	migs, err := listEmbeddedMigrations()
	if err != nil {
		return err
	}

	for _, m := range migs {
		var applied int
		err := db.QueryRow(`SELECT COUNT(1) FROM schema_migrations WHERE version = ?`, m.version).Scan(&applied)
		if err != nil {
			return fmt.Errorf("store: check migration %d: %w", m.version, err)
		}
		if applied > 0 {
			continue
		}

		body, err := schemaFS.ReadFile("schema/" + m.name)
		if err != nil {
			return fmt.Errorf("store: read migration %s: %w", m.name, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("store: begin migration %d: %w", m.version, err)
		}
		if _, err := tx.Exec(string(body)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("store: apply migration %d (%s): %w", m.version, m.name, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations(version) VALUES (?)`, m.version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("store: record migration %d: %w", m.version, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("store: commit migration %d: %w", m.version, err)
		}
	}
	return nil
}

func parseMigrationVersion(name string) (int, error) {
	// Expect NNN_name.sql (e.g. 001_init.sql)
	parts := strings.SplitN(name, "_", 2)
	if len(parts) < 1 {
		return 0, fmt.Errorf("store: bad migration name %q", name)
	}
	ver, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("store: bad migration version in %q: %w", name, err)
	}
	return ver, nil
}
