package store

import (
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestProjectBindPathLocalIsolation(t *testing.T) {
	a := t.TempDir()
	b := t.TempDir()

	sa, err := Open(a)
	if err != nil {
		t.Fatalf("Open(A): %v", err)
	}
	defer sa.Close()

	sb, err := Open(b)
	if err != nil {
		t.Fatalf("Open(B): %v", err)
	}
	defer sb.Close()

	dbA := sa.DBPath()
	dbB := sb.DBPath()
	if dbA == dbB {
		t.Fatalf("expected distinct DB paths, both %q", dbA)
	}
	wantA := filepath.Join(a, traceDirName, dbFileName)
	wantB := filepath.Join(b, traceDirName, dbFileName)
	if dbA != wantA {
		t.Fatalf("A DBPath: got %q want %q", dbA, wantA)
	}
	if dbB != wantB {
		t.Fatalf("B DBPath: got %q want %q", dbB, wantB)
	}

	absA, err := filepath.Abs(a)
	if err != nil {
		t.Fatal(err)
	}
	absB, err := filepath.Abs(b)
	if err != nil {
		t.Fatal(err)
	}
	var rootA, rootB string
	if err := sa.db.QueryRow(`SELECT root_path FROM projects WHERE id = ?`, sa.ProjectID()).Scan(&rootA); err != nil {
		t.Fatalf("A root_path: %v", err)
	}
	if err := sb.db.QueryRow(`SELECT root_path FROM projects WHERE id = ?`, sb.ProjectID()).Scan(&rootB); err != nil {
		t.Fatalf("B root_path: %v", err)
	}
	if rootA != absA || rootB != absB {
		t.Fatalf("root_path: A=%q want %q; B=%q want %q", rootA, absA, rootB, absB)
	}
	if sa.ProjectID() == sb.ProjectID() {
		t.Fatalf("expected distinct project IDs")
	}

	goal, err := sa.UpsertGoal(Goal{
		Title:      "only-in-A",
		SourceType: "USER_ASSERTED",
		Confidence: 1,
		Status:     StatusActive,
	})
	if err != nil {
		t.Fatalf("UpsertGoal A: %v", err)
	}
	if _, err := sb.GetGoal(goal.ID); err == nil {
		t.Fatal("B must not see A's goal")
	}
}

func TestConcurrentStoreOpenFailClosed(t *testing.T) {
	root := t.TempDir()

	s1, err := Open(root)
	if err != nil {
		t.Fatalf("first Open: %v", err)
	}

	var secondErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, secondErr = Open(root)
	}()
	wg.Wait()
	if !errors.Is(secondErr, ErrLocked) {
		s1.Close()
		t.Fatalf("second Open: want ErrLocked, got %v", secondErr)
	}

	if _, err := s1.UpsertGoal(Goal{
		Title:      "while-locked",
		SourceType: "USER_ASSERTED",
		Confidence: 1,
		Status:     StatusActive,
	}); err != nil {
		s1.Close()
		t.Fatalf("write under first Open: %v", err)
	}

	if err := s1.Close(); err != nil {
		t.Fatalf("Close first: %v", err)
	}

	s2, err := Open(root)
	if err != nil {
		t.Fatalf("Open after Close: %v", err)
	}
	defer s2.Close()

	var n int
	if err := s2.db.QueryRow(`SELECT COUNT(*) FROM goals WHERE title = ?`, "while-locked").Scan(&n); err != nil || n != 1 {
		t.Fatalf("goal preserved after re-open: n=%d err=%v", n, err)
	}
}

// TestOpenRetrySucceedsWhenLockReleasedSoon locks DF-47: brief contention during
// Open wait budget recovers when the holder Close()s soon (soft race).
func TestOpenRetrySucceedsWhenLockReleasedSoon(t *testing.T) {
	root := t.TempDir()
	s1, err := Open(root)
	if err != nil {
		t.Fatalf("first Open: %v", err)
	}

	started := make(chan struct{})
	var s2 *Store
	var secondErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		close(started)
		s2, secondErr = Open(root)
	}()
	<-started
	time.Sleep(50 * time.Millisecond)
	if err := s1.Close(); err != nil {
		t.Fatalf("Close first: %v", err)
	}
	wg.Wait()
	if secondErr != nil {
		t.Fatalf("second Open after brief release: %v", secondErr)
	}
	if s2 == nil {
		t.Fatal("second Open: nil store")
	}
	defer s2.Close()
}

// TestErrLockedSerializeGuidance locks DF-47: ErrLocked text guides serialize
// CLI↔MCP or separate -C / worktree roots (errors.Is identity unchanged).
func TestErrLockedSerializeGuidance(t *testing.T) {
	msg := ErrLocked.Error()
	for _, want := range []string{"serialize", "CLI", "MCP", "worktree"} {
		if !strings.Contains(msg, want) {
			t.Errorf("ErrLocked missing %q: %s", want, msg)
		}
	}
	if !errors.Is(ErrLocked, ErrLocked) {
		t.Fatal("ErrLocked must remain errors.Is identity")
	}
}
