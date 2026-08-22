package store

import (
	"database/sql"
	"fmt"
)

const graphSyncRowID = 1

// GraphSyncState tracks the last successful symbol/file graph index at git HEAD.
// Separate from vcs_meta / MetaVCSWatermark (VCS history index).
type GraphSyncState struct {
	LastIndexedCommit string
	LastIndexedAt     string
	HookInstalled     bool
}

// GetGraphSyncState returns the singleton graph sync row (id=1).
func (s *Store) GetGraphSyncState() (GraphSyncState, error) {
	if s == nil || s.db == nil {
		return GraphSyncState{}, fmt.Errorf("store: graph sync state: nil store")
	}
	var hook int
	var st GraphSyncState
	err := s.db.QueryRow(`
		SELECT last_indexed_commit, last_indexed_at, hook_installed
		FROM graph_sync_state WHERE id = ?
	`, graphSyncRowID).Scan(&st.LastIndexedCommit, &st.LastIndexedAt, &hook)
	if err == sql.ErrNoRows {
		return GraphSyncState{}, nil
	}
	if err != nil {
		return GraphSyncState{}, fmt.Errorf("store: get graph sync state: %w", err)
	}
	st.HookInstalled = hook != 0
	return st, nil
}

// UpsertGraphSyncState updates the singleton graph sync row (id=1).
func (s *Store) UpsertGraphSyncState(st GraphSyncState) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store: upsert graph sync state: nil store")
	}
	hook := 0
	if st.HookInstalled {
		hook = 1
	}
	_, err := s.db.Exec(`
		INSERT INTO graph_sync_state(id, last_indexed_commit, last_indexed_at, hook_installed)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			last_indexed_commit = excluded.last_indexed_commit,
			last_indexed_at = excluded.last_indexed_at,
			hook_installed = excluded.hook_installed
	`, graphSyncRowID, st.LastIndexedCommit, st.LastIndexedAt, hook)
	if err != nil {
		return fmt.Errorf("store: upsert graph sync state: %w", err)
	}
	return nil
}
