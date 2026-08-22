package main

import (
	"strings"
	"testing"
)

func TestPlanHelp_MentionsRefinement(t *testing.T) {
	help := captureStdout(t, func() int {
		return run([]string{"plan", "help"})
	})
	if !strings.Contains(help, "create-coarse") || !strings.Contains(help, "deep") {
		t.Fatalf("plan help missing refinement note: %q", help)
	}
	if !strings.Contains(help, "minimal plan") {
		t.Fatalf("plan help missing minimal plan honesty: %q", help)
	}
}
