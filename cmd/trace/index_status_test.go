package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/analyzers"
)

func statusJSON(t *testing.T, dir string) indexStatusJSON {
	t.Helper()
	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "index", "status"})
	})
	var status indexStatusJSON
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatalf("status json: %v\n%s", err, out)
	}
	if status.Drift == nil {
		t.Fatalf("drift section missing in: %s", out)
	}
	return status
}

func TestIndexStatusSupportedLanguages(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	statusJSON := captureStdout(t, func() int {
		return run([]string{"-C", dir, "index", "status"})
	})
	var status indexStatusJSON
	if err := json.Unmarshal([]byte(statusJSON), &status); err != nil {
		t.Fatalf("status json: %v\n%s", err, statusJSON)
	}
	want := analyzers.SupportedLanguages()
	if !reflect.DeepEqual(status.SupportedLanguages, want) {
		t.Fatalf("supported_languages=%v want %v", status.SupportedLanguages, want)
	}
	if len(status.SupportedLanguages) != 5 {
		t.Fatalf("len=%d want 5", len(status.SupportedLanguages))
	}
}

// TestIndexStatusReportsNotIndexedDrift pins the silent-truncation signal:
// files on disk that the store never got (a whole-tree index aborted partway,
// or a store frozen by a constraint failure) must be visible from
// `index status` without a full re-index.
func TestIndexStatusReportsNotIndexedDrift(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	for _, name := range []string{"a.js", "b.js", "c.js"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("export function f"+name+"() { return 1 }\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Index only two of the three files.
	if code := run([]string{"-C", dir, "index", "a.js", "b.js"}); code != exitOK {
		t.Fatalf("partial index: %d", code)
	}

	status := statusJSON(t, dir)
	d := status.Drift
	if d.Walkable != 3 || d.Indexed != 2 || d.NotIndexed != 1 {
		t.Fatalf("drift counts: %+v", d)
	}
	if len(d.NotIndexedPaths) != 1 || d.NotIndexedPaths[0] != "c.js" {
		t.Fatalf("not_indexed_paths: %v", d.NotIndexedPaths)
	}
	if d.NotOnDisk != 0 || len(d.NotOnDiskPaths) != 0 || d.PathsTruncated {
		t.Fatalf("unexpected other-direction drift: %+v", d)
	}
}

// TestIndexStatusReportsNotOnDiskDrift pins the other direction: indexed
// paths whose file disappeared are reported so a caller knows a full index
// will shrink the store.
func TestIndexStatusReportsNotOnDiskDrift(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.js"), []byte("export function fa() { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "index", "a.js"}); code != exitOK {
		t.Fatalf("index: %d", code)
	}
	if err := os.Remove(filepath.Join(dir, "a.js")); err != nil {
		t.Fatal(err)
	}

	status := statusJSON(t, dir)
	d := status.Drift
	if d.NotOnDisk != 1 || len(d.NotOnDiskPaths) != 1 || d.NotOnDiskPaths[0] != "a.js" {
		t.Fatalf("not_on_disk drift: %+v", d)
	}
	if d.NotIndexed != 0 {
		t.Fatalf("unexpected not_indexed drift: %+v", d)
	}
}

// TestIndexStatusNoDriftAfterFullIndex pins the healthy state: a complete
// whole-tree index leaves no drift in either direction.
func TestIndexStatusNoDriftAfterFullIndex(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	for _, name := range []string{"a.js", "b.ts"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("export function f() { return 1 }\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if code := run([]string{"-C", dir, "index"}); code != exitOK {
		t.Fatalf("full index: %d", code)
	}

	status := statusJSON(t, dir)
	d := status.Drift
	if d.Walkable != 2 || d.Indexed != 2 || d.NotIndexed != 0 || d.NotOnDisk != 0 {
		t.Fatalf("drift after full index: %+v", d)
	}
	if d.NotIndexedPaths != nil || d.NotOnDiskPaths != nil || d.PathsTruncated {
		t.Fatalf("drift lists should be absent: %+v", d)
	}
}

// TestIndexStatusDriftStderrNamesSample checks the operator-facing stderr
// summary names the missing paths and the remedy.
func TestIndexStatusDriftStderrNamesSample(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.js"), []byte("export function fa() { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	stderr := captureStderr(t, func() int {
		return run([]string{"-C", dir, "index", "status"})
	})
	if !strings.Contains(stderr, "not in the store") || !strings.Contains(stderr, `"a.js"`) || !strings.Contains(stderr, "trace index") {
		t.Fatalf("stderr drift summary: %q", stderr)
	}
}
