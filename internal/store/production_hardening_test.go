package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestMigrationStatusReportsEmbedMax(t *testing.T) {
	s, _ := openTempStore(t)
	st, err := s.MigrationStatus()
	if err != nil {
		t.Fatalf("MigrationStatus: %v", err)
	}
	if st.EmbedExpected != 28 {
		t.Fatalf("EmbedExpected: got %d want 28", st.EmbedExpected)
	}
	if st.MaxApplied != 28 {
		t.Fatalf("MaxApplied: got %d want 28", st.MaxApplied)
	}
	if st.PendingCount != 0 {
		t.Fatalf("PendingCount: got %d want 0", st.PendingCount)
	}
	if len(st.AppliedVersions) != 28 {
		t.Fatalf("AppliedVersions len: got %d want 28", len(st.AppliedVersions))
	}
	for i, v := range st.AppliedVersions {
		if v != i+1 {
			t.Fatalf("AppliedVersions[%d]=%d want %d", i, v, i+1)
		}
	}
}

func TestBackupRestoreRoundTrip(t *testing.T) {
	rootA := t.TempDir()
	sa, err := Open(rootA)
	if err != nil {
		t.Fatalf("Open A: %v", err)
	}
	goal, err := sa.UpsertGoal(Goal{
		Title:      "backup-marker-goal",
		Body:       "distinguishable",
		SourceType: "USER_ASSERTED",
		Confidence: 1,
		Status:     StatusActive,
	})
	if err != nil {
		t.Fatalf("UpsertGoal: %v", err)
	}
	goalID := goal.ID

	// Token must not ride along by default.
	if err := SetAccessToken(rootA, "secret-should-not-copy"); err != nil {
		t.Fatalf("SetAccessToken: %v", err)
	}

	bak := filepath.Join(t.TempDir(), "snap.db")
	if err := sa.BackupTo(bak); err != nil {
		sa.Close()
		t.Fatalf("BackupTo: %v", err)
	}
	if err := sa.Close(); err != nil {
		t.Fatalf("Close A: %v", err)
	}

	rootB := t.TempDir()
	if err := Restore(rootB, bak, false); err != nil {
		t.Fatalf("Restore B: %v", err)
	}

	// Default restore must not install access.token.
	if _, err := os.Stat(filepath.Join(rootB, traceDirName, accessTokenFileName)); !os.IsNotExist(err) {
		t.Fatalf("access.token must be absent on B after default restore; err=%v", err)
	}

	sb, err := Open(rootB)
	if err != nil {
		t.Fatalf("Open B: %v", err)
	}
	defer sb.Close()

	got, err := sb.GetGoal(goalID)
	if err != nil {
		t.Fatalf("GetGoal on B: %v", err)
	}
	if got.Title != "backup-marker-goal" {
		t.Fatalf("goal title: got %q", got.Title)
	}

	absB, err := filepath.Abs(rootB)
	if err != nil {
		t.Fatal(err)
	}
	var rootPath string
	if err := sb.db.QueryRow(`SELECT root_path FROM projects WHERE id = ?`, sb.ProjectID()).Scan(&rootPath); err != nil {
		t.Fatalf("root_path: %v", err)
	}
	if rootPath != absB {
		t.Fatalf("rebind: root_path=%q want %q", rootPath, absB)
	}

	bad, where, err := sb.HasBlobLikeColumns()
	if err != nil {
		t.Fatalf("HasBlobLikeColumns: %v", err)
	}
	if bad {
		t.Fatalf("unexpected BLOB-like column: %s", where)
	}
}

func TestRestoreForceAndExists(t *testing.T) {
	root := t.TempDir()
	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	bak := filepath.Join(t.TempDir(), "snap.db")
	if err := s.BackupTo(bak); err != nil {
		s.Close()
		t.Fatalf("BackupTo: %v", err)
	}
	s.Close()

	if err := Restore(root, bak, false); !errors.Is(err, ErrRestoreExists) {
		t.Fatalf("restore without force: want ErrRestoreExists, got %v", err)
	}
	if err := Restore(root, bak, true); err != nil {
		t.Fatalf("restore --force: %v", err)
	}
}

func TestLocalAccessTokenFailClosed(t *testing.T) {
	root := t.TempDir()
	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open cold: %v", err)
	}
	s.Close()

	if err := SetAccessToken(root, "correct-token"); err != nil {
		t.Fatalf("SetAccessToken: %v", err)
	}

	t.Setenv(AccessTokenEnv, "")
	if _, err := Open(root); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("Open without env: want ErrUnauthorized, got %v", err)
	}

	t.Setenv(AccessTokenEnv, "wrong-token")
	if _, err := Open(root); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("Open wrong env: want ErrUnauthorized, got %v", err)
	}

	t.Setenv(AccessTokenEnv, "correct-token")
	s2, err := Open(root)
	if err != nil {
		t.Fatalf("Open matching token: %v", err)
	}
	s2.Close()

	if err := ClearAccessToken(root); err != nil {
		t.Fatalf("ClearAccessToken: %v", err)
	}
	t.Setenv(AccessTokenEnv, "")
	s3, err := Open(root)
	if err != nil {
		t.Fatalf("Open after clear: %v", err)
	}
	s3.Close()

	enabled, err := AccessTokenEnabled(root)
	if err != nil || enabled {
		t.Fatalf("AccessTokenEnabled after clear: enabled=%v err=%v", enabled, err)
	}
}

func TestBackupFailsWhenLocked(t *testing.T) {
	root := t.TempDir()
	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	bak := filepath.Join(t.TempDir(), "snap.db")
	err = Backup(root, bak, BackupOptions{})
	if !errors.Is(err, ErrLocked) {
		t.Fatalf("Backup while locked: want ErrLocked, got %v", err)
	}
	if _, err := os.Stat(bak); !os.IsNotExist(err) {
		t.Fatalf("backup file must not be created on lock failure; err=%v", err)
	}

	// Holding Open must still allow in-process BackupTo.
	if err := s.BackupTo(bak); err != nil {
		t.Fatalf("BackupTo under held Open: %v", err)
	}
}

func TestRestoreFailsWhenLocked(t *testing.T) {
	rootA := t.TempDir()
	sa, err := Open(rootA)
	if err != nil {
		t.Fatalf("Open A: %v", err)
	}
	bak := filepath.Join(t.TempDir(), "snap.db")
	if err := sa.BackupTo(bak); err != nil {
		sa.Close()
		t.Fatalf("BackupTo: %v", err)
	}
	sa.Close()

	rootB := t.TempDir()
	sb, err := Open(rootB)
	if err != nil {
		t.Fatalf("Open B: %v", err)
	}
	defer sb.Close()

	if err := Restore(rootB, bak, true); !errors.Is(err, ErrLocked) {
		t.Fatalf("Restore while locked: want ErrLocked, got %v", err)
	}
}
