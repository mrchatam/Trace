package main

import (
	"strings"
	"testing"
)

func TestHelpIncludesSearchTestVerify(t *testing.T) {
	out := captureStdout(t, func() int { return run([]string{"help"}) })
	checks := []string{
		"search <query>",
		"test run",
		"verify run",
		"changes capture",
		"knowledge list",
		"install git-hook",
		"loop next --task <id>",
		"loop apply [--in <path>]",
		"add task --from-discovery <id>",
		"writes.spawned_tasks[].discovery_id",
		"loop status --task <id>",
		"loop gate --task",
		"agents list",
		"agents recommend",
		"graph infer",
		"Does not run on trace index",
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Fatalf("help missing %q", want)
		}
	}
}

// TestHelpBuildNoteCGOForTraceMCP locks the build note to scripts/build-with-proxy.sh:
// both binaries are built with CGO_ENABLED=1. The analyzers' tree-sitter bindings are
// cgo packages, so a CGO_ENABLED=0 trace-mcp build fails with "build constraints
// exclude all Go files" — the stale note sent rebuilders into that wall.
func TestHelpBuildNoteCGOForTraceMCP(t *testing.T) {
	out := captureStdout(t, func() int { return run([]string{"help"}) })
	if !strings.Contains(out, "CGO_ENABLED=1 go build -o bin/trace-mcp ./cmd/trace-mcp") {
		t.Fatalf("help build note must build trace-mcp with CGO_ENABLED=1 (see scripts/build-with-proxy.sh); got:\n%s", out)
	}
	if strings.Contains(out, "CGO_ENABLED=0") {
		t.Fatalf("help must not advertise any CGO_ENABLED=0 build; got:\n%s", out)
	}
}
