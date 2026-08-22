package store

import (
	"database/sql"
	"fmt"
)

// WithTx runs fn inside a single database transaction. All store methods invoked on
// the Store passed to fn use the same tx; Commit runs only when fn returns nil.
func (s *Store) WithTx(fn func(*Store) error) error {
	if s == nil || s.conn == nil {
		return fmt.Errorf("store: WithTx requires an open root store")
	}
	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	txStore := &Store{
		db:          tx,
		projectRoot: s.projectRoot,
		projectID:   s.projectID,
	}
	if err := fn(txStore); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit tx: %w", err)
	}
	return nil
}

// Conn returns the underlying *sql.DB for diagnostics and tests. Nil on tx-scoped stores.
func (s *Store) Conn() *sql.DB {
	if s == nil {
		return nil
	}
	return s.conn
}

// runInTx executes fn in the active transaction when s.db is *sql.Tx; otherwise it
// begins/commits its own transaction on s.conn.
func (s *Store) runInTx(fn func(*sql.Tx) error) error {
	if tx, ok := s.db.(*sql.Tx); ok {
		return fn(tx)
	}
	if s.conn == nil {
		return fmt.Errorf("store: runInTx: no database connection")
	}
	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("store: begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := fn(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: commit tx: %w", err)
	}
	return nil
}
