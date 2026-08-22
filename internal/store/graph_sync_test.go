package store

import (
	"testing"
	"time"
)

func TestGraphSyncStateRoundTrip(t *testing.T) {
	s, _ := openTempStore(t)
	st, err := s.GetGraphSyncState()
	if err != nil {
		t.Fatalf("GetGraphSyncState: %v", err)
	}
	if st.LastIndexedCommit != "" || st.LastIndexedAt != "" || st.HookInstalled {
		t.Fatalf("initial state: %+v", st)
	}

	st.LastIndexedCommit = "abc123"
	st.LastIndexedAt = time.Now().UTC().Format(time.RFC3339)
	st.HookInstalled = true
	if err := s.UpsertGraphSyncState(st); err != nil {
		t.Fatalf("UpsertGraphSyncState: %v", err)
	}

	got, err := s.GetGraphSyncState()
	if err != nil {
		t.Fatalf("GetGraphSyncState after upsert: %v", err)
	}
	if got.LastIndexedCommit != st.LastIndexedCommit {
		t.Fatalf("commit: got %q want %q", got.LastIndexedCommit, st.LastIndexedCommit)
	}
	if got.LastIndexedAt != st.LastIndexedAt {
		t.Fatalf("at: got %q want %q", got.LastIndexedAt, st.LastIndexedAt)
	}
	if !got.HookInstalled {
		t.Fatal("hook_installed want true")
	}
}
