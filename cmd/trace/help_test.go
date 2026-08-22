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
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Fatalf("help missing %q", want)
		}
	}
}
