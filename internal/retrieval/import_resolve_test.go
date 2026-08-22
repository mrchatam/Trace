package retrieval

import (
	"testing"
)

func TestImportPathCandidates_subdirRelativeOrder(t *testing.T) {
	got := importPathCandidates("src/app.js", "./util.js")
	wantPrefix := []string{"util.js", "src/util.js"}
	if len(got) < 2 {
		t.Fatalf("candidates too short: %v", got)
	}
	for i, w := range wantPrefix {
		if got[i] != w {
			t.Fatalf("cand[%d]=%q want %q; full=%v", i, got[i], w, got)
		}
	}
}

func TestImportPathCandidates_extensionlessThenIndex(t *testing.T) {
	got := importPathCandidates("src/app.ts", "./util")
	// exact, joined, then exts on bases without suffix, then joined/index+exts
	mustContainInOrder(t, got, []string{
		"util",
		"src/util",
		"util.js",
		"src/util.ts",
		"src/util/index.js",
	})
	// dedupe: no duplicate paths
	seen := map[string]bool{}
	for _, p := range got {
		if seen[p] {
			t.Fatalf("duplicate candidate %q in %v", p, got)
		}
		seen[p] = true
	}
}

func TestImportPathCandidates_bareModuleExactOnly(t *testing.T) {
	got := importPathCandidates("src/app.go", "fmt")
	if len(got) != 1 || got[0] != "fmt" {
		t.Fatalf("bare want [fmt] got %v", got)
	}
	got = importPathCandidates("pkg/a.go", "github.com/example/whyfile/pkg/b")
	if len(got) != 1 || got[0] != "github.com/example/whyfile/pkg/b" {
		t.Fatalf("module path want exact only got %v", got)
	}
}

func TestImportPathCandidates_backslashRelative(t *testing.T) {
	// Single backslash form (Windows-style relative) before slash-normalize.
	got := importPathCandidates("src/app.js", `.\util.js`)
	if len(got) < 2 || got[0] != "util.js" || got[1] != "src/util.js" {
		t.Fatalf("backslash relative: %v", got)
	}
}

func mustContainInOrder(t *testing.T, got, want []string) {
	t.Helper()
	idx := 0
	for _, w := range want {
		for idx < len(got) && got[idx] != w {
			idx++
		}
		if idx >= len(got) {
			t.Fatalf("missing %q (in order) in %v", w, got)
		}
		idx++
	}
}
