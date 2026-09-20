package main

import (
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/store"
)

func TestIndexForceHealsCorruptEdgesAndReportsHashSkipped(t *testing.T) {
	dir := t.TempDir()
	git := gitTestHelper(t, dir)
	git("init")
	git("config", "user.email", "trace@test.local")
	git("config", "user.name", "Trace Test")

	writeGoFile(t, dir, "hello.go", "package main\nfunc Hello() {}\n")
	writeGoFile(t, dir, "hello_test.go", "package main\nimport \"testing\"\nfunc TestHello(t *testing.T) {}\n")
	git("add", "-A")
	git("commit", "-m", "add")

	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	stderr := captureStderr(t, func() int {
		return run([]string{"-C", dir, "index"})
	})
	if !strings.Contains(stderr, "indexed") {
		t.Fatalf("stderr missing indexed: %q", stderr)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	edges, err := st.ListEdgesByFile("hello.go")
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if len(edges) == 0 {
		st.Close()
		t.Fatal("expected edges after index")
	}
	if err := st.ReplaceFileEdges("hello.go", nil); err != nil {
		st.Close()
		t.Fatal(err)
	}
	if err := st.ReplaceFileEdges("hello_test.go", nil); err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	stderr = captureStderr(t, func() int {
		return run([]string{"-C", dir, "index"})
	})
	if !strings.Contains(stderr, "hash_skipped") {
		t.Fatalf("want hash_skipped in stderr: %q", stderr)
	}

	st, err = store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	edges, _ = st.ListEdgesByFile("hello.go")
	if len(edges) != 0 {
		st.Close()
		t.Fatalf("without --force edges should stay empty, got %d", len(edges))
	}
	st.Close()

	stderr = captureStderr(t, func() int {
		return run([]string{"-C", dir, "index", "--force"})
	})
	if !strings.Contains(stderr, "indexed") {
		t.Fatalf("--force stderr: %q", stderr)
	}

	st, err = store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	edges, err = st.ListEdgesByFile("hello.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) == 0 {
		t.Fatal("--force should heal edges")
	}
}
