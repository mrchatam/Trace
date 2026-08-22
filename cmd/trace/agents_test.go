package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCLIAgentsList(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "install", "agents"}); code != exitOK {
		t.Fatalf("install agents: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "agents", "list"})
	})
	var items []map[string]any
	if err := json.Unmarshal([]byte(out), &items); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if len(items) != 6 {
		t.Fatalf("want 6 agents, got %d: %s", len(items), out)
	}
	found := false
	for _, item := range items {
		if item["slug"] == "agent:code-reviewer" {
			found = true
			if item["subagent_type"] != "code-reviewer" {
				t.Fatalf("subagent_type: %+v", item)
			}
			reqs, ok := item["requirements"].([]any)
			if !ok || len(reqs) == 0 {
				t.Fatalf("requirements: %+v", item["requirements"])
			}
			phases, ok := item["deliberation_phases"].([]any)
			if !ok || len(phases) == 0 {
				t.Fatalf("deliberation_phases: %+v", item["deliberation_phases"])
			}
		}
	}
	if !found {
		t.Fatalf("agent:code-reviewer missing: %s", out)
	}
}

func TestCLIAgentsDescribe(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "install", "agents"}); code != exitOK {
		t.Fatalf("install agents: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "agents", "describe", "agent:code-reviewer"})
	})
	var profile map[string]any
	if err := json.Unmarshal([]byte(out), &profile); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if profile["slug"] != "agent:code-reviewer" || profile["registry_source"] != "bundled" {
		t.Fatalf("profile: %+v", profile)
	}
	if profile["description"] == "" {
		t.Fatal("expected description")
	}

	code, _, stderr := runCapture(t, []string{"-C", dir, "agents", "describe", "agent:missing"})
	if code != exitFail {
		t.Fatalf("want exitFail got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "unknown slug") {
		t.Fatalf("stderr: %q", stderr)
	}
}

func TestCLIAgentsRecommend(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	if code := run([]string{"-C", dir, "install", "agents"}); code != exitOK {
		t.Fatalf("install agents: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "agents", "recommend", "--phase", "CRITIQUE"})
	})
	var recs []map[string]any
	if err := json.Unmarshal([]byte(out), &recs); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if len(recs) == 0 {
		t.Fatalf("expected recommendations: %s", out)
	}
	if recs[0]["agent_slug"] != "agent:code-reviewer" {
		t.Fatalf("first rec: %+v", recs[0])
	}
}

func TestHelpIncludesAgents(t *testing.T) {
	out := captureStdout(t, func() int { return run([]string{"help"}) })
	checks := []string{
		"agents list",
		"agents describe",
		"agents recommend",
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Fatalf("help missing %q", want)
		}
	}
}
