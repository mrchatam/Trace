package analyzers

import (
	"context"
	"testing"
)

func TestIndexFileForceRebuildsDespiteHashMatch(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	content := []byte("package p\n\nfunc Hello() {}\n")
	if err := IndexFile(ctx, st, "pkg/hello.go", content, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile: %v", err)
	}
	edges, err := st.ListEdgesByFile("pkg/hello.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) == 0 {
		t.Fatal("expected outgoing edges after index")
	}
	if err := st.ReplaceFileEdges("pkg/hello.go", nil); err != nil {
		t.Fatalf("clear edges: %v", err)
	}
	cleared, err := st.ListEdgesByFile("pkg/hello.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(cleared) != 0 {
		t.Fatalf("edges not cleared: %d", len(cleared))
	}

	// Without Force, hash match skips — edges stay gone.
	var skipped bool
	if err := IndexFile(ctx, st, "pkg/hello.go", content, IndexOptions{HashSkipped: &skipped}); err != nil {
		t.Fatal(err)
	}
	if !skipped {
		t.Fatal("expected hash skip without Force")
	}
	still, err := st.ListEdgesByFile("pkg/hello.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(still) != 0 {
		t.Fatalf("hash skip should not rebuild edges, got %d", len(still))
	}

	// With Force, edges heal.
	skipped = false
	if err := IndexFile(ctx, st, "pkg/hello.go", content, IndexOptions{Force: true, HashSkipped: &skipped}); err != nil {
		t.Fatalf("Force IndexFile: %v", err)
	}
	if skipped {
		t.Fatal("Force must not report HashSkipped")
	}
	healed, err := st.ListEdgesByFile("pkg/hello.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(healed) == 0 {
		t.Fatal("Force should rebuild outgoing edges")
	}
}
