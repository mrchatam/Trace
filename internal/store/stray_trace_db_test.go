package store

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func captureWarnWriter(t *testing.T) *bytes.Buffer {
	t.Helper()
	buf := &bytes.Buffer{}
	prev := warnWriter
	warnWriter = buf
	t.Cleanup(func() { warnWriter = prev })
	return buf
}

func writeRootStub(t *testing.T, root string, payload []byte) string {
	t.Helper()
	stub := filepath.Join(root, dbFileName)
	if err := os.WriteFile(stub, payload, 0o644); err != nil {
		t.Fatalf("write stub: %v", err)
	}
	return stub
}

func TestOpenWarnsWhenRootStubPresent(t *testing.T) {
	root := t.TempDir()
	stubPayload := []byte("not-a-trace-store")
	stub := writeRootStub(t, root, stubPayload)
	buf := captureWarnWriter(t)

	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	if !strings.Contains(buf.String(), "project-root trace.db exists") {
		t.Fatalf("expected stray warn, got %q", buf.String())
	}
	wantDB := filepath.Join(root, traceDirName, dbFileName)
	if s.DBPath() != wantDB {
		t.Fatalf("DBPath: got %q want %q", s.DBPath(), wantDB)
	}
	if _, err := os.Stat(wantDB); err != nil {
		t.Fatalf("live db missing: %v", err)
	}
	got, err := os.ReadFile(stub)
	if err != nil {
		t.Fatalf("read stub after open: %v", err)
	}
	if !bytes.Equal(got, stubPayload) {
		t.Fatalf("stub mutated: got %q want %q", got, stubPayload)
	}
}

func TestOpenExistingWarnsWhenRootStubPresent(t *testing.T) {
	root := t.TempDir()
	s0, err := Open(root)
	if err != nil {
		t.Fatalf("init Open: %v", err)
	}
	if err := s0.Close(); err != nil {
		t.Fatalf("close init: %v", err)
	}

	stubPayload := []byte{0}
	stub := writeRootStub(t, root, stubPayload)
	fiBefore, err := os.Stat(stub)
	if err != nil {
		t.Fatal(err)
	}
	buf := captureWarnWriter(t)

	s, err := OpenExisting(root)
	if err != nil {
		t.Fatalf("OpenExisting: %v", err)
	}
	defer s.Close()

	if !strings.Contains(buf.String(), "project-root trace.db exists") {
		t.Fatalf("expected stray warn, got %q", buf.String())
	}
	wantDB := filepath.Join(root, traceDirName, dbFileName)
	if s.DBPath() != wantDB {
		t.Fatalf("DBPath: got %q want %q", s.DBPath(), wantDB)
	}
	fiAfter, err := os.Stat(stub)
	if err != nil {
		t.Fatalf("stub missing after open: %v", err)
	}
	if fiAfter.Size() != fiBefore.Size() {
		t.Fatalf("stub size changed: %d → %d", fiBefore.Size(), fiAfter.Size())
	}
}

func TestOpenQuietWhenNoRootStub(t *testing.T) {
	root := t.TempDir()
	buf := captureWarnWriter(t)

	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	if buf.Len() != 0 {
		t.Fatalf("expected no warn, got %q", buf.String())
	}
	wantDB := filepath.Join(root, traceDirName, dbFileName)
	if s.DBPath() != wantDB {
		t.Fatalf("DBPath: got %q want %q", s.DBPath(), wantDB)
	}
}

func TestOpenLeavesRootStubUntouched(t *testing.T) {
	root := t.TempDir()
	stubPayload := []byte("keep-me")
	stub := writeRootStub(t, root, stubPayload)
	fiBefore, err := os.Stat(stub)
	if err != nil {
		t.Fatal(err)
	}
	_ = captureWarnWriter(t)

	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	wantDB := filepath.Join(root, traceDirName, dbFileName)
	if _, err := os.Stat(wantDB); err != nil {
		t.Fatalf("live db under .trace missing: %v", err)
	}
	fiAfter, err := os.Stat(stub)
	if err != nil {
		t.Fatalf("stub removed: %v", err)
	}
	if fiAfter.Size() != fiBefore.Size() {
		t.Fatalf("stub size changed: %d → %d", fiBefore.Size(), fiAfter.Size())
	}
	got, err := os.ReadFile(stub)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, stubPayload) {
		t.Fatalf("stub content changed: got %q want %q", got, stubPayload)
	}
}

// TestOpenQuietWhenRootStubIsDirectory: a directory named trace.db at root is
// not a regular file, so warnIfStrayRootTraceDB stays quiet (!IsRegular).
func TestOpenQuietWhenRootStubIsDirectory(t *testing.T) {
	root := t.TempDir()
	dirStub := filepath.Join(root, dbFileName)
	if err := os.Mkdir(dirStub, 0o755); err != nil {
		t.Fatalf("mkdir dir stub: %v", err)
	}
	buf := captureWarnWriter(t)

	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer s.Close()

	if buf.Len() != 0 {
		t.Fatalf("expected no warn for dir-named stub, got %q", buf.String())
	}
	wantDB := filepath.Join(root, traceDirName, dbFileName)
	if s.DBPath() != wantDB {
		t.Fatalf("DBPath: got %q want %q", s.DBPath(), wantDB)
	}
	if _, err := os.Stat(wantDB); err != nil {
		t.Fatalf("live db missing: %v", err)
	}
	fi, err := os.Stat(dirStub)
	if err != nil {
		t.Fatalf("dir stub removed/renamed: %v", err)
	}
	if !fi.IsDir() {
		t.Fatalf("dir stub is no longer a directory")
	}
}
