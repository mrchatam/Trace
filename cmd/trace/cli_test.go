package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func TestDispatchHelpVersionUnknown(t *testing.T) {
	if code := run(nil); code != exitOK {
		t.Fatalf("empty args: exit %d", code)
	}
	if code := run([]string{"help"}); code != exitOK {
		t.Fatalf("help: exit %d", code)
	}
	if code := run([]string{"version"}); code != exitOK {
		t.Fatalf("version: exit %d", code)
	}
	if code := run([]string{"--version"}); code != exitOK {
		t.Fatalf("--version: exit %d", code)
	}
	if code := run([]string{"not-a-command"}); code != exitUsage {
		t.Fatalf("unknown: want exit %d got %d", exitUsage, code)
	}
}

// TestAsOperatorFlagIdentityDocs locks DF-44: help states conscious flag≠identity.
func TestAsOperatorFlagIdentityDocs(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	if !strings.Contains(out, "flag≠identity") && !strings.Contains(out, "not verified") {
		t.Fatalf("help must mention flag≠identity / not verified identity: %q", out)
	}
	if !strings.Contains(out, "conscious") {
		t.Fatalf("help must mention conscious claim: %q", out)
	}
	if !strings.Contains(out, "Actor string ≠ auth") && !strings.Contains(out, "Actor") {
		t.Fatalf("help must keep Actor ≠ auth honesty: %q", out)
	}
}

// TestHelpSerializeLockGuidance locks DF-47: Global help guides serialize CLI↔MCP / worktrees.
func TestHelpSerializeLockGuidance(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	lower := strings.ToLower(out)
	for _, want := range []string{"serialize", "cli", "mcp", "worktree"} {
		if !strings.Contains(lower, want) {
			t.Fatalf("help missing %q: %q", want, out)
		}
	}
}

// TestLinkDiscoveryMentionsTaskCLI covers DF-42 CLI alias → store rel.
func TestLinkDiscoveryMentionsTaskCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	discOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "discovery", "--title", "D"})
	})
	var discRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(discOut)), &discRes); err != nil {
		t.Fatalf("disc json: %v", err)
	}
	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--title", "T"})
	})
	var taskRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(taskOut)), &taskRes); err != nil {
		t.Fatalf("task json: %v", err)
	}
	if code := run([]string{
		"-C", dir, "link", "discovery-mentions-task",
		"--from", discRes["id"].(string), "--to", taskRes["id"].(string),
	}); code != exitOK {
		t.Fatalf("link discovery-mentions-task: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	links, err := st.ListLinksFrom("discovery", discRes["id"].(string))
	if err != nil || len(links) != 1 || links[0].Rel != "discovery_mentions_task" {
		t.Fatalf("store links: %+v err=%v", links, err)
	}
	help := captureStdout(t, func() int { return run([]string{"help"}) })
	if !strings.Contains(help, "discovery-mentions-task") {
		t.Fatalf("help missing discovery-mentions-task: %q", help)
	}
}

func TestAddTaskFromDiscoveryCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "goal", "--title", "Goal"})
	})
	var goalRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(goalOut)), &goalRes); err != nil {
		t.Fatalf("goal json: %v", err)
	}
	goalID := goalRes["id"].(string)

	discOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "discovery", "--title", "Gap", "--severity", "BLOCKING"})
	})
	var discRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(discOut)), &discRes); err != nil {
		t.Fatalf("discovery json: %v", err)
	}
	discID := discRes["id"].(string)

	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--from-discovery", discID, "--goal-id", goalID})
	})
	var taskRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(taskOut)), &taskRes); err != nil {
		t.Fatalf("task json: %v", err)
	}
	if taskRes["id"] != discID {
		t.Fatalf("promoted task id mismatch: %v", taskRes)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	links, err := st.ListLinksFrom("discovery", discID)
	if err != nil || len(links) != 1 || links[0].Rel != "discovery_mentions_task" || links[0].ToID != discID {
		t.Fatalf("promotion link missing: links=%+v err=%v", links, err)
	}
}

func TestInitCreatesDB(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init exit %d", code)
	}
	db := filepath.Join(dir, ".trace", "trace.db")
	if _, err := os.Stat(db); err != nil {
		t.Fatalf("expected db at %s: %v", db, err)
	}
}

func TestCausalWhyContextRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	goalOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "goal", "--title", "Ship CLI", "--id", "11111111-1111-1111-1111-111111111111"})
	})
	var goalRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(goalOut)), &goalRes); err != nil {
		t.Fatalf("goal json: %v (%q)", err, goalOut)
	}
	goalID := goalRes["id"].(string)

	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--title", "Wire commands", "--goal-id", goalID, "--id", "22222222-2222-2222-2222-222222222222"})
	})
	var taskRes map[string]any
	_ = json.Unmarshal([]byte(strings.TrimSpace(taskOut)), &taskRes)
	taskID := taskRes["id"].(string)

	decOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "decision", "--title", "Stdlib argv", "--id", "33333333-3333-3333-3333-333333333333"})
	})
	var decRes map[string]any
	_ = json.Unmarshal([]byte(strings.TrimSpace(decOut)), &decRes)
	decID := decRes["id"].(string)

	if code := run([]string{"-C", dir, "link", "decision-task", "--from", decID, "--to", taskID}); code != exitOK {
		t.Fatalf("link: %d", code)
	}

	whyOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "why", "task", taskID})
	})
	var why map[string]any
	if err := json.Unmarshal([]byte(whyOut), &why); err != nil {
		t.Fatalf("why json: %v\n%s", err, whyOut)
	}
	steps, ok := why["steps"].([]any)
	if !ok || len(steps) == 0 {
		t.Fatalf("why missing steps: %v", why)
	}
	for _, s := range steps {
		step := s.(map[string]any)
		if _, ok := step["reason_code"]; !ok {
			t.Fatalf("step missing reason_code: %v", step)
		}
	}

	ctxOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "context", taskID, "--format", "json"})
	})
	var pkt map[string]any
	if err := json.Unmarshal([]byte(ctxOut), &pkt); err != nil {
		t.Fatalf("context json: %v\n%s", err, ctxOut)
	}
	items, ok := pkt["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatalf("context missing items: %v", pkt)
	}
	budget, ok := pkt["budget"].(map[string]any)
	if !ok {
		t.Fatalf("context missing budget: %v", pkt)
	}
	if budget["token_limit"] == nil || budget["max_items"] == nil {
		t.Fatalf("budget incomplete: %v", budget)
	}
	foundTrust := false
	for _, it := range items {
		m := it.(map[string]any)
		if m["trust"] != nil && m["trust"] != "" {
			foundTrust = true
		}
		if m["excerpt"] != nil && m["excerpt"] != "" {
			if m["trust"] != "untrusted_data" && m["trust"] != "system" {
				t.Fatalf("excerpt item unexpected trust: %v", m)
			}
		}
	}
	if !foundTrust {
		t.Fatalf("expected trust labels on context items")
	}
}

func TestSeedImportAndWhy(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedPath := filepath.Join(dir, "seed.json")
	seed := `{
  "version": 1,
  "goals": [{"id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","title":"G","body":""}],
  "tasks": [{"id":"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb","title":"T","body":"task body","goal_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"}],
  "decisions": [{"id":"cccccccc-cccc-cccc-cccc-cccccccccccc","title":"D","body":""}],
  "assumptions": [],
  "discoveries": [],
  "plan_changes": [],
  "claims": [],
  "evidence": [],
  "links": [
    {"rel":"decision_affects_task","from":"cccccccc-cccc-cccc-cccc-cccccccccccc","to":"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"}
  ],
  "transitions": [
    {"task_id":"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb","to":"IN_PROGRESS","actor":"seed","reason":"start","allow_done":false}
  ]
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "seed", "import", seedPath})
	})
	var summary map[string]any
	if err := json.Unmarshal([]byte(out), &summary); err != nil {
		t.Fatalf("seed summary: %v (%s)", err, out)
	}
	if summary["ok"] != true {
		t.Fatalf("seed not ok: %v", summary)
	}

	taskID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	whyOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "why", "task", taskID})
	})
	if !strings.Contains(whyOut, "reason_code") {
		t.Fatalf("why missing reason_code: %s", whyOut)
	}
	ctxOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "context", taskID})
	})
	if !strings.Contains(ctxOut, `"items"`) {
		t.Fatalf("context missing items: %s", ctxOut)
	}
}

func TestIndexIncrementalIsolation(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	aPath := filepath.Join(dir, "a.js")
	bPath := filepath.Join(dir, "b.js")
	if err := os.WriteFile(aPath, []byte("export function alpha() { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bPath, []byte("export function beta() { return 2 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if code := run([]string{"-C", dir, "index", "a.js", "b.js"}); code != exitOK {
		t.Fatalf("index both: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	symsB1, err := st.ListSymbolsByPath("b.js")
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if len(symsB1) == 0 {
		st.Close()
		t.Fatal("expected symbols on b.js after index")
	}
	st.Close()

	if err := os.WriteFile(aPath, []byte("export function alpha() { return 99 }\nexport function alpha2() { return 0 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "index", "a.js"}); code != exitOK {
		t.Fatalf("reindex a: %d", code)
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	symsB2, err := st2.ListSymbolsByPath("b.js")
	if err != nil {
		t.Fatal(err)
	}
	if len(symsB2) != len(symsB1) {
		t.Fatalf("b.js symbols changed after indexing only a.js: before=%d after=%d", len(symsB1), len(symsB2))
	}
	for i := range symsB1 {
		if symsB1[i].Name != symsB2[i].Name || symsB1[i].Kind != symsB2[i].Kind {
			t.Fatalf("b.js symbol drift: %+v vs %+v", symsB1[i], symsB2[i])
		}
	}
	symsA, err := st2.ListSymbolsByPath("a.js")
	if err != nil {
		t.Fatal(err)
	}
	if len(symsA) < 2 {
		t.Fatalf("expected updated a.js symbols, got %d", len(symsA))
	}
}

func TestIndexGCAfterPathRename(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	aPath := filepath.Join(dir, "a.js")
	bPath := filepath.Join(dir, "b.js")
	if err := os.WriteFile(aPath, []byte("export function alpha() { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bPath, []byte("export function beta() { return 2 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	stderr1 := captureStderr(t, func() int {
		return run([]string{"-C", dir, "index"})
	})
	if !strings.Contains(stderr1, "removed") {
		t.Fatalf("expected removed count in stderr, got %q", stderr1)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetFileByPath("a.js"); err != nil {
		st.Close()
		t.Fatalf("a.js after index: %v", err)
	}
	symsB1, err := st.ListSymbolsByPath("b.js")
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if len(symsB1) == 0 {
		st.Close()
		t.Fatal("expected b.js symbols")
	}
	st.Close()

	cPath := filepath.Join(dir, "c.js")
	if err := os.Rename(aPath, cPath); err != nil {
		t.Fatal(err)
	}

	stderr2 := captureStderr(t, func() int {
		return run([]string{"-C", dir, "index"})
	})
	if !strings.Contains(stderr2, "removed 1") {
		t.Fatalf("expected removed 1 after rename, got %q", stderr2)
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()

	if _, err := st2.GetFileByPath("a.js"); err == nil {
		t.Fatal("ghost a.js still in store after full-tree GC")
	}
	if _, err := st2.ListSymbolsByPath("a.js"); err == nil {
		t.Fatal("ghost symbols for a.js still present")
	}
	ghostHits, err := st2.SearchFTS("alpha", 10)
	if err != nil {
		t.Fatalf("SearchFTS alpha: %v", err)
	}
	for _, h := range ghostHits {
		if h.Path == "a.js" {
			t.Fatalf("ghost FTS row for a.js: %+v", h)
		}
	}

	if _, err := st2.GetFileByPath("c.js"); err != nil {
		t.Fatalf("c.js missing after rename+index: %v", err)
	}
	symsC, err := st2.ListSymbolsByPath("c.js")
	if err != nil {
		t.Fatal(err)
	}
	foundAlpha := false
	for _, sym := range symsC {
		if sym.Name == "alpha" {
			foundAlpha = true
			break
		}
	}
	if !foundAlpha {
		t.Fatalf("expected alpha on c.js, got %+v", symsC)
	}

	symsB2, err := st2.ListSymbolsByPath("b.js")
	if err != nil {
		t.Fatal(err)
	}
	if len(symsB2) != len(symsB1) {
		t.Fatalf("b.js changed after GC: before=%d after=%d", len(symsB1), len(symsB2))
	}
	for i := range symsB1 {
		if symsB1[i].Name != symsB2[i].Name {
			t.Fatalf("b.js symbol drift: %+v vs %+v", symsB1[i], symsB2[i])
		}
	}
}

func TestIndexPartialArgvGCAfterRename(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	aPath := filepath.Join(dir, "a.js")
	bPath := filepath.Join(dir, "b.js")
	if err := os.WriteFile(aPath, []byte("export function alpha() { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bPath, []byte("export function beta() { return 2 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if code := run([]string{"-C", dir, "index", "a.js", "b.js"}); code != exitOK {
		t.Fatalf("index both: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	symsB1, err := st.ListSymbolsByPath("b.js")
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if len(symsB1) == 0 {
		st.Close()
		t.Fatal("expected b.js symbols")
	}
	st.Close()

	cPath := filepath.Join(dir, "c.js")
	if err := os.Rename(aPath, cPath); err != nil {
		t.Fatal(err)
	}

	stderr := captureStderr(t, func() int {
		return run([]string{"-C", dir, "index", "c.js"})
	})
	if !strings.Contains(stderr, "removed") {
		t.Fatalf("expected removed in stderr after partial rename GC, got %q", stderr)
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()

	if _, err := st2.GetFileByPath("a.js"); err == nil {
		t.Fatal("ghost a.js still in store after partial argv GC")
	}
	if _, err := st2.ListSymbolsByPath("a.js"); err == nil {
		t.Fatal("ghost symbols for a.js still present")
	}
	ghostHits, err := st2.SearchFTS("alpha", 10)
	if err != nil {
		t.Fatalf("SearchFTS alpha: %v", err)
	}
	for _, h := range ghostHits {
		if h.Path == "a.js" {
			t.Fatalf("ghost FTS row for a.js: %+v", h)
		}
	}

	if _, err := st2.GetFileByPath("c.js"); err != nil {
		t.Fatalf("c.js missing after rename+partial index: %v", err)
	}
	symsC, err := st2.ListSymbolsByPath("c.js")
	if err != nil {
		t.Fatal(err)
	}
	foundAlpha := false
	for _, sym := range symsC {
		if sym.Name == "alpha" {
			foundAlpha = true
			break
		}
	}
	if !foundAlpha {
		t.Fatalf("expected alpha on c.js, got %+v", symsC)
	}

	if _, err := st2.GetFileByPath("b.js"); err != nil {
		t.Fatalf("b.js must remain after partial argv: %v", err)
	}
	symsB2, err := st2.ListSymbolsByPath("b.js")
	if err != nil {
		t.Fatal(err)
	}
	if len(symsB2) != len(symsB1) {
		t.Fatalf("b.js changed after partial GC: before=%d after=%d", len(symsB1), len(symsB2))
	}
	for i := range symsB1 {
		if symsB1[i].Name != symsB2[i].Name {
			t.Fatalf("b.js symbol drift: %+v vs %+v", symsB1[i], symsB2[i])
		}
	}
}

func TestIndexArgvMissingPathDeletesOnlyThatPath(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	aPath := filepath.Join(dir, "a.js")
	bPath := filepath.Join(dir, "b.js")
	if err := os.WriteFile(aPath, []byte("export function alpha() { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bPath, []byte("export function beta() { return 2 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "index", "a.js", "b.js"}); code != exitOK {
		t.Fatalf("index both: %d", code)
	}
	if err := os.Remove(aPath); err != nil {
		t.Fatal(err)
	}

	stderr := captureStderr(t, func() int {
		return run([]string{"-C", dir, "index", "a.js"})
	})
	if !strings.Contains(stderr, "removed 1") {
		t.Fatalf("expected removed 1 for missing argv, got %q", stderr)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, err := st.GetFileByPath("a.js"); err == nil {
		t.Fatal("a.js should be deleted after missing argv index")
	}
	if _, err := st.GetFileByPath("b.js"); err != nil {
		t.Fatalf("b.js must remain (no project-wide GC): %v", err)
	}
}

func TestWalkIndexableT0AlwaysSkip(t *testing.T) {
	dir := t.TempDir()
	plant := []struct {
		rel  string
		body string
	}{
		{"src/ok.js", "export function ok() { return 1 }\n"},
		{"node_modules/pkg/index.js", "export function hidden() { return 0 }\n"},
		{"vendor/lib/util.js", "export function vendored() { return 0 }\n"},
		{"__pycache__/mod.py", "def cached():\n    return 0\n"},
		{"dist/bundle.js", "export function bundled() { return 0 }\n"},
		{".next/server.js", "export function nextish() { return 0 }\n"},
		{"target/out.js", "export function targeted() { return 0 }\n"},
		{"coverage/cov.js", "export function covered() { return 0 }\n"},
		{".venv/lib/site.py", "def venvish():\n    return 0\n"},
		{"venv/lib/site.py", "def venv2():\n    return 0\n"},
		{"foo.min.js", "export function minned() { return 0 }\n"},
		{"lib/bar.min.mjs", "export function minned2() { return 0 }\n"},
		{"lib/baz.min.cjs", "exports.minned3 = function() { return 0 }\n"},
	}
	for _, p := range plant {
		full := filepath.Join(dir, filepath.FromSlash(p.rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(p.body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	start := time.Now()
	got, err := walkIndexable(dir, false)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("walkIndexable T0 plant: %d paths in %s", len(got), elapsed)

	set := map[string]struct{}{}
	for _, p := range got {
		set[p] = struct{}{}
	}
	if _, ok := set["src/ok.js"]; !ok {
		t.Fatalf("expected control src/ok.js in walk; got %#v", got)
	}
	forbidden := []string{
		"node_modules/pkg/index.js",
		"vendor/lib/util.js",
		"__pycache__/mod.py",
		"dist/bundle.js",
		".next/server.js",
		"target/out.js",
		"coverage/cov.js",
		".venv/lib/site.py",
		"venv/lib/site.py",
		"foo.min.js",
		"lib/bar.min.mjs",
		"lib/baz.min.cjs",
	}
	for _, f := range forbidden {
		if _, ok := set[f]; ok {
			t.Errorf("T0 path %q must not appear in walkIndexable; got %#v", f, got)
		}
	}
	if len(got) != 1 {
		t.Fatalf("want exactly control path; got %#v", got)
	}
}

func TestIndexSkipsExplicitT0Path(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	nm := filepath.Join(dir, "node_modules", "x.js")
	if err := os.MkdirAll(filepath.Dir(nm), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nm, []byte("export function secret() { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Control: index a normal file so the store is known-good.
	okPath := filepath.Join(dir, "ok.js")
	if err := os.WriteFile(okPath, []byte("export function ok() { return 1 }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if code := run([]string{"-C", dir, "index", "node_modules/x.js", "ok.js"}); code != exitOK {
		t.Fatalf("index: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if _, err := st.GetFileByPath("node_modules/x.js"); err == nil {
		t.Fatal("explicit T0 path must not be upserted")
	}
	if _, err := st.GetFileByPath("ok.js"); err != nil {
		t.Fatalf("control ok.js should be indexed: %v", err)
	}
}

func TestReviewCreateSetDone(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--title", "Ship review", "--id", "dddddddd-dddd-dddd-dddd-dddddddddddd"})
	})
	var taskRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(taskOut)), &taskRes); err != nil {
		t.Fatalf("task json: %v", err)
	}
	taskID := taskRes["id"].(string)

	if code := run([]string{"-C", dir, "transition", "--task", taskID, "--to", "IN_PROGRESS", "--reason", "start"}); code != exitOK {
		t.Fatalf("to IN_PROGRESS: %d", code)
	}
	// Evidence alone must not unlock DONE
	if code := run([]string{"-C", dir, "transition", "--task", taskID, "--to", "DONE", "--reason", "try", "--evidence", "ev-x"}); code == exitOK {
		t.Fatal("DONE with evidence alone must fail")
	}

	revOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "review", "create", "--title", "Gate", "--task", taskID, "--id", "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"})
	})
	var revRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(revOut)), &revRes); err != nil {
		t.Fatalf("review create: %v (%s)", err, revOut)
	}
	if revRes["rel"] != "review_judges_task" {
		t.Fatalf("expected link: %v", revRes)
	}
	revID := revRes["id"].(string)

	if code := run([]string{"-C", dir, "review", "set", "--id", revID, "--result", "PASS", "--reason", "ok", "--actor", "rev"}); code != exitOK {
		t.Fatalf("review set: %d", code)
	}
	if code := run([]string{"-C", dir, "transition", "--task", taskID, "--to", "DONE", "--reason", "promote"}); code == exitOK {
		t.Fatal("DONE after PASS without --as-operator must fail")
	}
	if code := run([]string{"-C", dir, "transition", "--task", taskID, "--to", "DONE", "--reason", "promote", "--as-operator"}); code != exitOK {
		t.Fatalf("DONE after PASS + --as-operator: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task, err := st.GetTask(taskID)
	if err != nil || task.WorkState != store.WorkStateDone {
		t.Fatalf("want DONE: %+v err=%v", task, err)
	}
}

func TestInitFailClosedWhenStoreLocked(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	code := run([]string{"-C", dir, "init"})
	if code != exitFail {
		t.Fatalf("init while locked: want exit %d got %d", exitFail, code)
	}
}

func TestMigrateBackupAuthCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "migrate", "status"})
	})
	if !strings.Contains(out, "embed_expected=28") || !strings.Contains(out, "max_applied=28") {
		t.Fatalf("migrate status: %q", out)
	}

	bak := filepath.Join(t.TempDir(), "cli-snap.db")
	out = captureStdout(t, func() int {
		return run([]string{"-C", dir, "backup", "-o", bak})
	})
	if strings.TrimSpace(out) != bak {
		t.Fatalf("backup stdout: %q", out)
	}
	if _, err := os.Stat(bak); err != nil {
		t.Fatalf("backup file: %v", err)
	}

	out = captureStdout(t, func() int {
		return run([]string{"-C", dir, "auth", "set", "cli-secret"})
	})
	if strings.TrimSpace(out) != "enabled" {
		t.Fatalf("auth set: %q", out)
	}
	if code := run([]string{"-C", dir, "migrate", "status"}); code != exitFail {
		t.Fatalf("migrate without token: want exitFail got %d", code)
	}
	t.Setenv(store.AccessTokenEnv, "cli-secret")
	if code := run([]string{"-C", dir, "migrate", "status"}); code != exitOK {
		t.Fatalf("migrate with token: %d", code)
	}
	out = captureStdout(t, func() int {
		return run([]string{"-C", dir, "auth", "clear"})
	})
	if strings.TrimSpace(out) != "disabled" {
		t.Fatalf("auth clear: %q", out)
	}
	t.Setenv(store.AccessTokenEnv, "")
	if code := run([]string{"-C", dir, "migrate", "status"}); code != exitOK {
		t.Fatalf("migrate after clear: %d", code)
	}

	dirB := t.TempDir()
	absB, err := filepath.Abs(dirB)
	if err != nil {
		t.Fatal(err)
	}
	out = captureStdout(t, func() int {
		return run([]string{"-C", dirB, "restore", bak})
	})
	if strings.TrimSpace(out) != absB {
		t.Fatalf("restore stdout: %q want %q", out, absB)
	}
}

func TestTasksListAfterSeed(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	taskID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	seedPath := filepath.Join(dir, "seed.json")
	seed := `{
  "version": 1,
  "goals": [{"id":"` + goalID + `","title":"G","body":""}],
  "tasks": [{"id":"` + taskID + `","title":"Wire tasks","body":"","goal_id":"` + goalID + `"}],
  "decisions": [],
  "assumptions": [],
  "discoveries": [],
  "plan_changes": [],
  "claims": [],
  "evidence": [],
  "links": [],
  "transitions": [
    {"task_id":"` + taskID + `","to":"IN_PROGRESS","actor":"seed","reason":"start","allow_done":false}
  ]
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("seed: %d", code)
	}

	emptyDir := t.TempDir()
	if code := run([]string{"-C", emptyDir, "init"}); code != exitOK {
		t.Fatalf("empty init: %d", code)
	}
	emptyOut := captureStdout(t, func() int {
		return run([]string{"-C", emptyDir, "tasks"})
	})
	if strings.TrimSpace(emptyOut) != "[]" {
		t.Fatalf("empty tasks: %q", emptyOut)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "tasks"})
	})
	var rows []map[string]any
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("tasks json: %v (%s)", err, out)
	}
	if len(rows) != 1 {
		t.Fatalf("tasks len: %v", rows)
	}
	if rows[0]["id"] != taskID || rows[0]["title"] != "Wire tasks" {
		t.Fatalf("row identity: %v", rows[0])
	}
	if rows[0]["work_state"] != "IN_PROGRESS" {
		t.Fatalf("work_state: %v", rows[0]["work_state"])
	}
	if rows[0]["goal_id"] != goalID {
		t.Fatalf("goal_id: %v", rows[0]["goal_id"])
	}

	filtered := captureStdout(t, func() int {
		return run([]string{"-C", dir, "tasks", "--goal", goalID})
	})
	var filt []map[string]any
	if err := json.Unmarshal([]byte(filtered), &filt); err != nil {
		t.Fatalf("tasks --goal: %v (%s)", err, filtered)
	}
	if len(filt) != 1 || filt[0]["id"] != taskID {
		t.Fatalf("filtered: %v", filt)
	}
	other := captureStdout(t, func() int {
		return run([]string{"-C", dir, "tasks", "--goal", "cccccccc-cccc-cccc-cccc-cccccccccccc"})
	})
	if strings.TrimSpace(other) != "[]" {
		t.Fatalf("unknown goal filter: %q", other)
	}
}

func TestSeedImportRelativePathAgainstC(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedDir := filepath.Join(dir, "seed")
	if err := os.MkdirAll(seedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	rel := "seed/relative.json"
	seedPath := filepath.Join(dir, filepath.FromSlash(rel))
	seed := `{
  "version": 1,
  "goals": [],
  "tasks": [{"id":"dddddddd-dddd-dddd-dddd-dddddddddddd","title":"Rel","body":""}],
  "decisions": [],
  "assumptions": [],
  "discoveries": [],
  "plan_changes": [],
  "claims": [],
  "evidence": [],
  "links": [],
  "transitions": []
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}

	cwd := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "seed", "import", rel})
	})
	var summary map[string]any
	if err := json.Unmarshal([]byte(out), &summary); err != nil {
		t.Fatalf("seed summary: %v (%s)", err, out)
	}
	if summary["ok"] != true {
		t.Fatalf("seed not ok: %v", summary)
	}

	tasksOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "tasks"})
	})
	var rows []map[string]any
	if err := json.Unmarshal([]byte(tasksOut), &rows); err != nil {
		t.Fatalf("tasks: %v (%s)", err, tasksOut)
	}
	if len(rows) != 1 || rows[0]["id"] != "dddddddd-dddd-dddd-dddd-dddddddddddd" {
		t.Fatalf("tasks after relative seed: %v", rows)
	}
	if rows[0]["goal_id"] != nil {
		t.Fatalf("unset goal_id want null, got %v", rows[0]["goal_id"])
	}

	// Absolute path still works when cwd ≠ project.
	absSeed := filepath.Join(dir, "seed", "abs.json")
	if err := os.WriteFile(absSeed, []byte(`{
  "version": 1,
  "goals": [],
  "tasks": [{"id":"eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee","title":"Abs","body":""}],
  "decisions": [],
  "assumptions": [],
  "discoveries": [],
  "plan_changes": [],
  "claims": [],
  "evidence": [],
  "links": [],
  "transitions": []
}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "seed", "import", absSeed}); code != exitOK {
		t.Fatalf("abs seed: %d", code)
	}
}

func TestAllowDoneWarnsOnStderr(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--title", "Hatch", "--id", "f0000000-0000-4000-8000-0000000000f1"})
	})
	var taskRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(taskOut)), &taskRes); err != nil {
		t.Fatal(err)
	}
	taskID := taskRes["id"].(string)
	if code := run([]string{"-C", dir, "transition", "--task", taskID, "--to", "IN_PROGRESS", "--reason", "start"}); code != exitOK {
		t.Fatalf("IN_PROGRESS: %d", code)
	}
	stderr := captureStderr(t, func() int {
		return run([]string{"-C", dir, "transition", "--task", taskID, "--to", "DONE", "--reason", "escape", "--allow-done"})
	})
	if !strings.Contains(stderr, "WARNING") || !strings.Contains(stderr, "allow-done") {
		t.Fatalf("expected loud allow-done WARNING, got %q", stderr)
	}
	if !strings.Contains(stderr, "allow-missing-caps") && !strings.Contains(stderr, "missing capabilities") {
		t.Fatalf("WARNING must mention missing-caps independence, got %q", stderr)
	}
}

func TestCapabilityMissingRequiresTaskHint(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = w
	code := run([]string{"-C", dir, "capability", "missing"})
	_ = w.Close()
	os.Stderr = old
	if code != exitUsage {
		t.Fatalf("want exitUsage, got %d", code)
	}
	buf := make([]byte, 1<<16)
	n, _ := r.Read(buf)
	_ = r.Close()
	stderr := string(buf[:n])
	if !strings.Contains(stderr, "usage:") || !strings.Contains(stderr, "trace tasks") {
		t.Fatalf("expected usage + tasks hint, got %q", stderr)
	}
}

// captureStdout runs fn while redirecting os.Stdout to a pipe.
func captureStdout(t *testing.T, fn func() int) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	code := fn()
	_ = w.Close()
	os.Stdout = old
	if code != exitOK {
		t.Fatalf("command exit %d", code)
	}
	buf := make([]byte, 1<<20)
	n, _ := r.Read(buf)
	_ = r.Close()
	return string(buf[:n])
}

// TestSeedImportFromIDAliases locks DF-33: from_id/to_id alone import a goal-task link.
func TestSeedImportFromIDAliases(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	taskID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	seedPath := filepath.Join(dir, "seed.json")
	seed := `{
  "version": 1,
  "goals": [{"id":"` + goalID + `","title":"G","body":""}],
  "tasks": [{"id":"` + taskID + `","title":"T","body":"","goal_id":"` + goalID + `"}],
  "links": [
    {"rel":"goal-task","from_id":"` + goalID + `","to_id":"` + taskID + `"}
  ]
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "seed", "import", seedPath})
	})
	var summary map[string]any
	if err := json.Unmarshal([]byte(out), &summary); err != nil {
		t.Fatalf("seed summary: %v (%s)", err, out)
	}
	if summary["ok"] != true {
		t.Fatalf("seed not ok: %v", summary)
	}
	if links, _ := summary["links"].(float64); links != 1 {
		t.Fatalf("want 1 link, got %v", summary["links"])
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	task, err := st.GetTask(taskID)
	if err != nil || task.GoalID == nil || *task.GoalID != goalID {
		t.Fatalf("goal link via from_id/to_id: %+v err=%v", task, err)
	}
}

// TestSeedImportMissingEndpointsMessage locks DF-33 alias-aware empty error.
func TestSeedImportMissingEndpointsMessage(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedPath := filepath.Join(dir, "seed.json")
	seed := `{
  "version": 1,
  "goals": [{"id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","title":"G","body":""}],
  "tasks": [{"id":"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb","title":"T","body":""}],
  "links": [{"rel":"goal-task"}]
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = w
	code := run([]string{"-C", dir, "seed", "import", seedPath})
	_ = w.Close()
	os.Stderr = old
	if code != exitUsage {
		t.Fatalf("want exitUsage, got %d", code)
	}
	buf := make([]byte, 1<<16)
	n, _ := r.Read(buf)
	_ = r.Close()
	stderr := string(buf[:n])
	lower := strings.ToLower(stderr)
	for _, want := range []string{"from", "to", "from_id", "to_id"} {
		if !strings.Contains(lower, want) {
			t.Fatalf("stderr missing %q: %q", want, stderr)
		}
	}
}

// TestPlanShowSnakeCaseAndEmptyPhases locks DF-30/46: phases [] + tasks + snake_case.
func TestPlanShowSnakeCaseAndEmptyPhases(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	taskID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	if code := run([]string{"-C", dir, "add", "goal", "--title", "G", "--id", goalID}); code != exitOK {
		t.Fatalf("add goal: %d", code)
	}
	if code := run([]string{"-C", dir, "add", "task", "--title", "T", "--id", taskID, "--goal-id", goalID}); code != exitOK {
		t.Fatalf("add task: %d", code)
	}
	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "plan", "show", "--goal", goalID})
	})
	if strings.Contains(out, `"GoalID"`) || strings.Contains(out, `"Phases"`) {
		t.Fatalf("PascalCase keys leaked: %s", out)
	}
	var view map[string]any
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatalf("json: %v (%s)", err, out)
	}
	if view["goal_id"] != goalID {
		t.Fatalf("goal_id: %v", view["goal_id"])
	}
	phases, ok := view["phases"].([]any)
	if !ok {
		t.Fatalf("phases must be array, got %T %v", view["phases"], view["phases"])
	}
	if phases == nil || len(phases) != 0 {
		t.Fatalf("want empty phases [], got %#v", phases)
	}
	tasks, ok := view["tasks"].([]any)
	if !ok || len(tasks) == 0 {
		t.Fatalf("want nonempty tasks, got %#v", view["tasks"])
	}
	row, _ := tasks[0].(map[string]any)
	if row["id"] != taskID || row["work_state"] == nil {
		t.Fatalf("task row shape: %#v", row)
	}
}

// TestPlanShowWithPhasesSnakeCase locks DF-46 nested phase/scope snake_case.
func TestPlanShowWithPhasesSnakeCase(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	goalOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "goal", "--title", "G"})
	})
	var goalRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(goalOut)), &goalRes); err != nil {
		t.Fatalf("goal json: %v", err)
	}
	goalID := goalRes["id"].(string)
	if code := run([]string{
		"-C", dir, "plan", "create-coarse", "--goal", goalID,
		"--phase", "P1", "--scope", "S1",
	}); code != exitOK {
		t.Fatalf("create-coarse: %d", code)
	}
	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "plan", "show", "--goal", goalID})
	})
	var view map[string]any
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		t.Fatalf("json: %v (%s)", err, out)
	}
	phases, ok := view["phases"].([]any)
	if !ok || len(phases) != 1 {
		t.Fatalf("phases: %#v", view["phases"])
	}
	ph, _ := phases[0].(map[string]any)
	for _, key := range []string{"id", "title", "body", "ord", "status", "scopes"} {
		if _, ok := ph[key]; !ok {
			t.Fatalf("phase missing %q: %#v", key, ph)
		}
	}
	scopes, _ := ph["scopes"].([]any)
	if len(scopes) != 1 {
		t.Fatalf("scopes: %#v", ph["scopes"])
	}
	sc, _ := scopes[0].(map[string]any)
	for _, key := range []string{"id", "phase_id", "title", "ord", "status", "auto_replan_count"} {
		if _, ok := sc[key]; !ok {
			t.Fatalf("scope missing %q: %#v", key, sc)
		}
	}
	if strings.Contains(out, `"PhaseID"`) || strings.Contains(out, `"AutoReplanCount"`) {
		t.Fatalf("PascalCase nested keys: %s", out)
	}
}

// TestReviewGetShowList locks DF-45 review get|show|list.
func TestReviewGetShowList(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	emptyOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "review", "list"})
	})
	if strings.TrimSpace(emptyOut) != "[]" {
		t.Fatalf("empty list want [], got %q", emptyOut)
	}

	taskOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "task", "--title", "T"})
	})
	var taskRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(taskOut)), &taskRes); err != nil {
		t.Fatalf("task: %v", err)
	}
	taskID := taskRes["id"].(string)

	revOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "review", "create", "--title", "Gate", "--body", "body-x", "--task", taskID})
	})
	var revRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(revOut)), &revRes); err != nil {
		t.Fatalf("create: %v", err)
	}
	revID := revRes["id"].(string)

	getOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "review", "get", "--id", revID})
	})
	var got map[string]any
	if err := json.Unmarshal([]byte(getOut), &got); err != nil {
		t.Fatalf("get json: %v", err)
	}
	if got["id"] != revID || got["title"] != "Gate" || got["body"] != "body-x" {
		t.Fatalf("get: %#v", got)
	}
	if _, has := got["result"]; !has {
		t.Fatalf("get missing result: %#v", got)
	}
	if _, has := got["status"]; !has {
		t.Fatalf("get missing status: %#v", got)
	}

	showOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "review", "show", "--id", revID})
	})
	if strings.TrimSpace(showOut) != strings.TrimSpace(getOut) {
		t.Fatalf("show != get:\n%s\n%s", showOut, getOut)
	}

	listOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "review", "list"})
	})
	var listed []map[string]any
	if err := json.Unmarshal([]byte(listOut), &listed); err != nil {
		t.Fatalf("list: %v (%s)", err, listOut)
	}
	if len(listed) != 1 || listed[0]["id"] != revID {
		t.Fatalf("list: %#v", listed)
	}
	if _, hasBody := listed[0]["body"]; hasBody {
		t.Fatalf("list must omit body: %#v", listed[0])
	}

	byTask := captureStdout(t, func() int {
		return run([]string{"-C", dir, "review", "list", "--task", taskID})
	})
	var filtered []map[string]any
	if err := json.Unmarshal([]byte(byTask), &filtered); err != nil {
		t.Fatalf("list --task: %v", err)
	}
	if len(filtered) != 1 || filtered[0]["id"] != revID {
		t.Fatalf("list --task: %#v", filtered)
	}
}

// TestHelpSeedExportPath locks DF-82/85: help mentions recommended commit path and evidence-not-identity.
func TestHelpSeedExportPath(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	lower := strings.ToLower(out)
	if !strings.Contains(out, "trace/graph.json") {
		t.Fatalf("help missing trace/graph.json: %q", out)
	}
	if !strings.Contains(out, "exported_at_commit") {
		t.Fatalf("help missing exported_at_commit: %q", out)
	}
	hasEvidence := strings.Contains(lower, "not identity") ||
		(strings.Contains(lower, "evidence") && strings.Contains(lower, "export"))
	if !hasEvidence {
		t.Fatalf("help missing evidence-not-identity phrasing: %q", out)
	}
}

// TestHelpCloneTasksImportPending locks DF-88: clone tasks import PENDING;
// default export omits reviews/transitions/work_state (not bare work_state).
func TestHelpCloneTasksImportPending(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	lower := strings.ToLower(out)
	for _, want := range []string{
		"pending",
		"import",
		"omits reviews",
		"transitions",
		"work_state",
	} {
		if !strings.Contains(lower, want) {
			t.Fatalf("help missing %q (DF-88 clone PENDING honesty): %q", want, out)
		}
	}
}

// TestHelpHandoffSoT locks DF-28 thin handoff guidance in help.
func TestHelpHandoffSoT(t *testing.T) {
	out := captureStdout(t, func() int {
		return run([]string{"help"})
	})
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "handoff") {
		t.Fatalf("help missing handoff: %q", out)
	}
	if !strings.Contains(lower, "context") {
		t.Fatalf("help missing context: %q", out)
	}
	if !strings.Contains(lower, "why") {
		t.Fatalf("help missing why: %q", out)
	}
}

// TestImpactWalkCLI covers P14-S01 thin adapter: seeds excluded, loud totals JSON.
func TestImpactWalkCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	a, err := st.UpsertFile("a.go", "ha", nil)
	if err != nil {
		t.Fatalf("UpsertFile a: %v", err)
	}
	b, err := st.UpsertFile("b.go", "hb", nil)
	if err != nil {
		t.Fatalf("UpsertFile b: %v", err)
	}
	if err := st.ReplaceFileImports("b.go", []store.Import{{ImportedPath: "a.go"}}); err != nil {
		t.Fatalf("imports: %v", err)
	}
	_ = st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "impact", "walk", "--seed", "file:" + a.ID, "--depth", "1"})
	})
	var res map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &res); err != nil {
		t.Fatalf("json: %v (%s)", err, out)
	}
	if res["ok"] != true {
		t.Fatalf("ok: %#v", res)
	}
	if res["truncated"] != false {
		t.Fatalf("truncated: %#v", res)
	}
	blast, _ := res["blast"].([]any)
	if len(blast) != 1 {
		t.Fatalf("blast: %#v", blast)
	}
	hit := blast[0].(map[string]any)
	if hit["entity_id"] != b.ID {
		t.Fatalf("want importer B: %#v", hit)
	}
	if hit["entity_id"] == a.ID {
		t.Fatal("seed must be excluded")
	}
	help := captureStdout(t, func() int { return run([]string{"help"}) })
	if !strings.Contains(help, "walk") {
		t.Fatalf("help missing walk: %q", help)
	}
	if !strings.Contains(help, "predict") || !strings.Contains(help, "compare") {
		t.Fatalf("help missing predict/compare: %q", help)
	}
}

// TestImpactPredictCompareCLI covers C08 thin adapters.
func TestImpactPredictCompareCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	svc := domain.New(st)
	lib, err := st.UpsertFile("lib.go", "hlib", nil)
	if err != nil {
		t.Fatal(err)
	}
	imp, err := st.UpsertFile("imp.go", "himp", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileImports("imp.go", []store.Import{{ImportedPath: "lib.go"}}); err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(context.Background(), domain.TaskInput{Title: "impact-cli"})
	if err != nil {
		t.Fatal(err)
	}
	retrieval.WireDomainImpactWalker(svc, retrieval.New(st))
	change, err := svc.CreateChange(context.Background(), domain.ChangeInput{
		TaskID: task.ID,
		Reason: "touch lib",
		Paths:  []domain.ChangePathInput{{Path: "lib.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	predictOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "impact", "predict", "--change", change.ID, "--depth", "1"})
	})
	var predict map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(predictOut)), &predict); err != nil {
		t.Fatalf("predict json: %v (%s)", err, predictOut)
	}
	if predict["ok"] != true {
		t.Fatalf("predict ok: %#v", predict)
	}
	blastKeys, _ := predict["blast_keys"].([]any)
	if len(blastKeys) != 1 {
		t.Fatalf("predict blast_keys: %#v", blastKeys)
	}
	if blastKeys[0].(string) != "file:"+imp.ID {
		t.Fatalf("predict blast key: %#v", blastKeys[0])
	}
	_ = lib

	compareOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "impact", "compare", "--change", change.ID})
	})
	var compare map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(compareOut)), &compare); err != nil {
		t.Fatalf("compare json: %v (%s)", err, compareOut)
	}
	if compare["ok"] != true || compare["compared_at"] == "" {
		t.Fatalf("compare: %#v", compare)
	}
	matched, _ := compare["matched"].([]any)
	if len(matched) != 1 || matched[0].(string) != "file:"+imp.ID {
		t.Fatalf("compare matched: %#v", matched)
	}

	failCode := run([]string{"-C", dir, "impact", "compare", "--change", "00000000-0000-0000-0000-000000000099"})
	if failCode == exitOK {
		t.Fatal("compare without prediction must fail")
	}
}

// TestSeedImportDiscoveryMentionsTask locks DF-70: seed accepts underscore and hyphen
// mentions-task rels and writes store rel discovery_mentions_task.
func TestSeedImportDiscoveryMentionsTask(t *testing.T) {
	const (
		discID = "dddddddd-dddd-dddd-dddd-dddddddddddd"
		taskID = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	)
	for _, tc := range []struct {
		name string
		rel  string
	}{
		{name: "underscore", rel: "discovery_mentions_task"},
		{name: "hyphen", rel: "discovery-mentions-task"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if code := run([]string{"-C", dir, "init"}); code != exitOK {
				t.Fatalf("init: %d", code)
			}
			seedPath := filepath.Join(dir, "seed.json")
			seed := fmt.Sprintf(`{
  "version": 1,
  "discoveries": [{"id":%q,"title":"D","body":""}],
  "tasks": [{"id":%q,"title":"T","body":""}],
  "links": [{"rel":%q,"from":%q,"to":%q}]
}`, discID, taskID, tc.rel, discID, taskID)
			if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
				t.Fatal(err)
			}
			out := captureStdout(t, func() int {
				return run([]string{"-C", dir, "seed", "import", seedPath})
			})
			var summary map[string]any
			if err := json.Unmarshal([]byte(out), &summary); err != nil {
				t.Fatalf("seed summary: %v (%s)", err, out)
			}
			if summary["ok"] != true {
				t.Fatalf("seed not ok: %v", summary)
			}
			st, err := store.Open(dir)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = st.Close() })
			links, err := st.ListLinksFrom("discovery", discID)
			if err != nil || len(links) != 1 || links[0].Rel != "discovery_mentions_task" {
				t.Fatalf("store links: %+v err=%v", links, err)
			}
			if links[0].ToID != taskID {
				t.Fatalf("to_id: %q", links[0].ToID)
			}
		})
	}

	t.Run("unknown_rel", func(t *testing.T) {
		dir := t.TempDir()
		if code := run([]string{"-C", dir, "init"}); code != exitOK {
			t.Fatalf("init: %d", code)
		}
		seedPath := filepath.Join(dir, "seed.json")
		seed := `{
  "version": 1,
  "discoveries": [{"id":"dddddddd-dddd-dddd-dddd-dddddddddddd","title":"D","body":""}],
  "tasks": [{"id":"eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee","title":"T","body":""}],
  "links": [{"rel":"not-a-rel","from":"dddddddd-dddd-dddd-dddd-dddddddddddd","to":"eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"}]
}`
		if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
			t.Fatal(err)
		}
		code, _, stderr := runCapture(t, []string{"-C", dir, "seed", "import", seedPath})
		if code != exitUsage {
			t.Fatalf("unknown rel want exitUsage got %d stderr=%q", code, stderr)
		}
		if !strings.Contains(stderr, "unknown link rel") {
			t.Fatalf("stderr: %q", stderr)
		}
	})
}

// TestSeedImportImpactFindings locks DF-73: top-level findings/alternatives seed
// through domain APIs; stub key impact_findings stays unknown.
func TestSeedImportImpactFindings(t *testing.T) {
	const (
		taskID = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
		decID  = "cccccccc-cccc-cccc-cccc-cccccccccccc"
		findID = "ffffffff-ffff-ffff-ffff-ffffffffffff"
		altID  = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01"
	)
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedPath := filepath.Join(dir, "seed.json")
	seed := `{
  "version": 1,
  "tasks": [{"id":"` + taskID + `","title":"T","body":""}],
  "decisions": [{"id":"` + decID + `","title":"D","body":""}],
  "links": [{"rel":"decision_affects_task","from":"` + decID + `","to":"` + taskID + `"}],
  "findings": [{
    "id": "` + findID + `",
    "decision_id": "` + decID + `",
    "impact_class": "DESTRUCTIVE",
    "kind": "AFFECTED_WORK",
    "body": "wipes state"
  }],
  "alternatives": [{
    "id": "` + altID + `",
    "decision_id": "` + decID + `",
    "title": "safer path",
    "recommended": true
  }]
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "seed", "import", seedPath})
	})
	var summary map[string]any
	if err := json.Unmarshal([]byte(out), &summary); err != nil {
		t.Fatalf("seed summary: %v (%s)", err, out)
	}
	if summary["ok"] != true {
		t.Fatalf("seed not ok: %v", summary)
	}
	if summary["findings"] != float64(1) {
		t.Fatalf("summary findings count: %v", summary["findings"])
	}
	if summary["alternatives"] != float64(1) {
		t.Fatalf("summary alternatives count: %v", summary["alternatives"])
	}

	repOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "impact", "report", "--decision", decID})
	})
	if !strings.Contains(repOut, `"impact_class"`) || !strings.Contains(repOut, "DESTRUCTIVE") {
		t.Fatalf("report missing snake_case finding: %s", repOut)
	}
	if !strings.Contains(repOut, `"is_recommended"`) {
		t.Fatalf("report missing is_recommended: %s", repOut)
	}
	if strings.Contains(repOut, `"ImpactClass"`) || strings.Contains(repOut, `"IsRecommended"`) {
		t.Fatalf("report still PascalCase: %s", repOut)
	}

	badPath := filepath.Join(dir, "bad.json")
	bad := `{
  "version": 1,
  "decisions": [{"id":"cccccccc-cccc-cccc-cccc-cccccccccccc","title":"D","body":""}],
  "impact_findings": []
}`
	if err := os.WriteFile(badPath, []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, stderr := runCapture(t, []string{"-C", dir, "seed", "import", badPath})
	if code != exitUsage {
		t.Fatalf("impact_findings want exitUsage got %d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "unknown top-level key") || !strings.Contains(stderr, "impact_findings") {
		t.Fatalf("stderr: %q", stderr)
	}
}

// TestImpactReportJSONSnakeCase locks DF-74: nested report JSON uses snake_case keys.
func TestImpactReportJSONSnakeCase(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	decOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "add", "decision", "--title", "Ship"})
	})
	var decRes map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(decOut)), &decRes); err != nil {
		t.Fatalf("decision json: %v", err)
	}
	decID := decRes["id"].(string)

	if code := run([]string{
		"-C", dir, "impact", "finding", "add",
		"--decision", decID, "--class", "DESTRUCTIVE", "--kind", "AFFECTED_WORK",
		"--body", "breaks callers",
	}); code != exitOK {
		t.Fatalf("finding add: %d", code)
	}
	if code := run([]string{
		"-C", dir, "impact", "alternative", "add",
		"--decision", decID, "--title", "Keep old API", "--recommended",
	}); code != exitOK {
		t.Fatalf("alternative add: %d", code)
	}

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "impact", "report", "--decision", decID})
	})
	for _, want := range []string{`"impact_class"`, `"is_recommended"`, `"overall_class"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("report missing %s: %s", want, out)
		}
	}
	if !strings.Contains(out, "DESTRUCTIVE") {
		t.Fatalf("report missing DESTRUCTIVE: %s", out)
	}
	for _, bad := range []string{`"ImpactClass"`, `"IsRecommended"`, `"ID"`} {
		if strings.Contains(out, bad) {
			t.Fatalf("report must not contain object key %s: %s", bad, out)
		}
	}
}

// TestWhyIncludesImpactOverallClass locks DF-71: why JSON inherits packet-shaped impact.
func TestWhyIncludesImpactOverallClass(t *testing.T) {
	const (
		taskID = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
		decID  = "cccccccc-cccc-cccc-cccc-cccccccccccc"
		emptyD = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	)
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedPath := filepath.Join(dir, "seed.json")
	seed := `{
  "version": 1,
  "tasks": [{"id":"` + taskID + `","title":"T","body":""}],
  "decisions": [
    {"id":"` + decID + `","title":"hot","body":""},
    {"id":"` + emptyD + `","title":"cold","body":""}
  ],
  "links": [
    {"rel":"decision_affects_task","from":"` + decID + `","to":"` + taskID + `"},
    {"rel":"decision_affects_task","from":"` + emptyD + `","to":"` + taskID + `"}
  ]
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("seed: %d", code)
	}
	if code := run([]string{
		"-C", dir, "impact", "finding", "add",
		"--decision", decID, "--class", "DESTRUCTIVE", "--kind", "WORK_AT_RISK",
	}); code != exitOK {
		t.Fatalf("finding add: %d", code)
	}

	t.Run("task", func(t *testing.T) {
		out := captureStdout(t, func() int {
			return run([]string{"-C", dir, "why", "task", taskID})
		})
		assertWhyImpactDestructive(t, out, decID)
	})
	t.Run("decision", func(t *testing.T) {
		out := captureStdout(t, func() int {
			return run([]string{"-C", dir, "why", "decision", decID})
		})
		assertWhyImpactDestructive(t, out, decID)
	})
}

func assertWhyImpactDestructive(t *testing.T, out, decID string) {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("why json: %v\n%s", err, out)
	}
	impact, ok := payload["impact"].([]any)
	if !ok || len(impact) != 1 {
		t.Fatalf("want one impact row: %s", out)
	}
	row := impact[0].(map[string]any)
	if row["decision_id"] != decID {
		t.Fatalf("decision_id: %v", row["decision_id"])
	}
	if row["overall_class"] != "DESTRUCTIVE" {
		t.Fatalf("overall_class: %v", row["overall_class"])
	}
}

// TestSeedExportRoundTrip locks DF-80/84: import → export → fresh import preserves ids + plan tree.
func TestSeedExportRoundTrip(t *testing.T) {
	const (
		goalID    = "11111111-1111-1111-1111-111111111111"
		taskID    = "22222222-2222-2222-2222-222222222222"
		decID     = "33333333-3333-3333-3333-333333333333"
		phaseID   = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01"
		scopeID   = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1"
		deepID    = "cccccccc-cccc-cccc-cccc-cccccccccc01"
		deepOldID = "cccccccc-cccc-cccc-cccc-cccccccccc02"
		findID    = "dddddddd-dddd-dddd-dddd-dddddddddd01"
	)
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedPath := filepath.Join(dir, "seed.json")
	seed := `{
  "version": 1,
  "goals": [{"id":"` + goalID + `","title":"G","body":"goal body"}],
  "tasks": [{"id":"` + taskID + `","title":"T","body":"task body","goal_id":"` + goalID + `"}],
  "decisions": [{"id":"` + decID + `","title":"D","body":""}],
  "assumptions": [],
  "discoveries": [],
  "plan_changes": [],
  "claims": [],
  "evidence": [],
  "links": [{"rel":"decision_affects_task","from":"` + decID + `","to":"` + taskID + `"}],
  "findings": [{
    "id":"` + findID + `",
    "decision_id":"` + decID + `",
    "impact_class":"CAUTION",
    "kind":"AFFECTED_WORK",
    "body":"watch"
  }],
  "alternatives": [],
  "plan_phases": [{
    "id":"` + phaseID + `",
    "goal_id":"` + goalID + `",
    "title":"Phase 1",
    "body":"",
    "ord": 0,
    "status":"ACTIVE"
  }],
  "plan_scopes": [{
    "id":"` + scopeID + `",
    "phase_id":"` + phaseID + `",
    "title":"Scope 1",
    "body":"",
    "ord": 0,
    "status":"ACTIVE",
    "auto_replan_count": 2
  }],
  "scope_deep_plans": [
    {
      "id":"` + deepOldID + `",
      "scope_id":"` + scopeID + `",
      "content_json":"{\"rev\":1}",
      "status":"SUPERSEDED"
    },
    {
      "id":"` + deepID + `",
      "scope_id":"` + scopeID + `",
      "content_json":"{\"rev\":2}",
      "status":"ACTIVE"
    }
  ],
  "goal_plan_state": [{
    "goal_id":"` + goalID + `",
    "current_scope_id":"` + scopeID + `"
  }]
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("seed import: %d", code)
	}

	exported := captureStdout(t, func() int {
		return run([]string{"-C", dir, "seed", "export"})
	})

	dir2 := t.TempDir()
	if code := run([]string{"-C", dir2, "init"}); code != exitOK {
		t.Fatalf("init dir2: %d", code)
	}
	exportPath := filepath.Join(dir2, "roundtrip.json")
	if err := os.WriteFile(exportPath, []byte(exported), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir2, "seed", "import", exportPath}); code != exitOK {
		t.Fatalf("seed re-import: %d", code)
	}

	st, err := store.Open(dir2)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	if g, err := st.GetGoal(goalID); err != nil || g.Title != "G" {
		t.Fatalf("goal: %+v err=%v", g, err)
	}
	if task, err := st.GetTask(taskID); err != nil || task.Title != "T" {
		t.Fatalf("task: %+v err=%v", task, err)
	}
	links, err := st.ListLinksByRel("decision_affects_task")
	if err != nil || len(links) != 1 || links[0].FromID != decID || links[0].ToID != taskID {
		t.Fatalf("decision link: %+v err=%v", links, err)
	}
	findings, err := st.ListDecisionImpactFindingsByDecisionID(decID)
	if err != nil || len(findings) != 1 || findings[0].ID != findID {
		t.Fatalf("findings: %+v err=%v", findings, err)
	}
	ph, err := st.GetPlanPhase(phaseID)
	if err != nil || ph.GoalID != goalID || ph.Title != "Phase 1" {
		t.Fatalf("plan phase: %+v err=%v", ph, err)
	}
	sc, err := st.GetPlanScope(scopeID)
	if err != nil || sc.PhaseID != phaseID || sc.AutoReplanCount != 2 {
		t.Fatalf("plan scope: %+v err=%v", sc, err)
	}
	d1, err := st.GetScopeDeepPlan(deepID)
	if err != nil || d1.Status != store.StatusActive {
		t.Fatalf("active deep plan: %+v err=%v", d1, err)
	}
	d0, err := st.GetScopeDeepPlan(deepOldID)
	if err != nil || d0.Status != store.StatusSuperseded {
		t.Fatalf("superseded deep plan: %+v err=%v", d0, err)
	}
	gps, err := st.GetGoalPlanState(goalID)
	if err != nil || gps.CurrentScopeID == nil || *gps.CurrentScopeID != scopeID {
		t.Fatalf("goal plan state: %+v err=%v", gps, err)
	}
}

// TestSeedExportOmitsDeniedSurfaces locks DF-80: export excludes transitions, work_state, index, token, caps, reviews.
func TestSeedExportOmitsDeniedSurfaces(t *testing.T) {
	const taskID = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedPath := filepath.Join(dir, "seed.json")
	seed := `{
  "version": 1,
  "tasks": [{"id":"` + taskID + `","title":"T","body":""}],
  "transitions": [{"task_id":"` + taskID + `","to":"IN_PROGRESS","reason":"start"}]
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("seed import: %d", code)
	}

	okPath := filepath.Join(dir, "ok.js")
	if err := os.WriteFile(okPath, []byte("export function ok() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "index", "ok.js"}); code != exitOK {
		t.Fatalf("index: %d", code)
	}
	if code := run([]string{"-C", dir, "auth", "set", "secret-token"}); code != exitOK {
		t.Fatalf("auth set: %d", code)
	}
	t.Setenv(store.AccessTokenEnv, "secret-token")
	if code := run([]string{"-C", dir, "capability", "decide", "--slug", "cli:add", "--decision", "DENIED"}); code != exitOK {
		t.Fatalf("capability decide: %d", code)
	}
	if code := run([]string{"-C", dir, "review", "create", "--title", "R", "--task", taskID}); code != exitOK {
		t.Fatalf("review create: %d", code)
	}

	exported := captureStdout(t, func() int {
		return run([]string{"-C", dir, "seed", "export"})
	})
	if strings.Contains(exported, `"transitions"`) {
		t.Fatalf("export must omit transitions key: %s", exported)
	}
	if strings.Contains(exported, `"work_state"`) {
		t.Fatalf("export must omit task work_state: %s", exported)
	}
	for _, forbidden := range []string{
		"secret-token", "access.token", "capability_tool_decisions",
		"review_judges_task", "node_modules", `"files"`, `"symbols"`,
	} {
		if strings.Contains(exported, forbidden) {
			t.Fatalf("export leaked %q: %s", forbidden, exported)
		}
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal([]byte(exported), &doc); err != nil {
		t.Fatalf("export json: %v", err)
	}
	if _, ok := doc["transitions"]; ok {
		t.Fatal("transitions key present in parsed export")
	}
	var tasks []map[string]any
	if err := json.Unmarshal(doc["tasks"], &tasks); err != nil || len(tasks) != 1 {
		t.Fatalf("tasks: %v", err)
	}
	if _, ok := tasks[0]["work_state"]; ok {
		t.Fatal("task object includes work_state")
	}
}

// TestSeedExportWritesExportedAtCommit locks DF-85: git HEAD on export; omit outside git; import ignores field.
func TestSeedExportWritesExportedAtCommit(t *testing.T) {
	t.Run("git", func(t *testing.T) {
		if _, err := exec.LookPath("git"); err != nil {
			t.Skip("git not on PATH")
		}
		dir := t.TempDir()
		gitCmd := func(args ...string) string {
			t.Helper()
			cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
			cmd.Env = append(os.Environ(),
				"GIT_AUTHOR_NAME=Trace Test",
				"GIT_AUTHOR_EMAIL=trace@test.local",
				"GIT_COMMITTER_NAME=Trace Test",
				"GIT_COMMITTER_EMAIL=trace@test.local",
			)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("git %v: %v\n%s", args, err, out)
			}
			return strings.TrimSpace(string(out))
		}
		gitCmd("init")
		gitCmd("config", "user.email", "trace@test.local")
		gitCmd("config", "user.name", "Trace Test")
		if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("seed export test\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		gitCmd("add", "-A")
		gitCmd("commit", "-m", "init")
		head := gitCmd("rev-parse", "HEAD")

		if code := run([]string{"-C", dir, "init"}); code != exitOK {
			t.Fatalf("init: %d", code)
		}
		goalID := "11111111-1111-1111-1111-111111111111"
		seedPath := filepath.Join(dir, "seed.json")
		seed := `{"version":1,"goals":[{"id":"` + goalID + `","title":"G","body":""}]}`
		if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
			t.Fatal(err)
		}
		if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
			t.Fatalf("seed import: %d", code)
		}

		exported := captureStdout(t, func() int {
			return run([]string{"-C", dir, "seed", "export"})
		})
		var doc map[string]any
		if err := json.Unmarshal([]byte(exported), &doc); err != nil {
			t.Fatalf("export json: %v", err)
		}
		got, _ := doc["exported_at_commit"].(string)
		if got == "" {
			t.Fatalf("exported_at_commit missing in git repo: %s", exported)
		}
		if got != head {
			t.Fatalf("exported_at_commit=%q want HEAD %q", got, head)
		}
	})

	t.Run("non-git", func(t *testing.T) {
		dir := t.TempDir()
		if code := run([]string{"-C", dir, "init"}); code != exitOK {
			t.Fatalf("init: %d", code)
		}
		goalID := "11111111-1111-1111-1111-111111111111"
		seedPath := filepath.Join(dir, "seed.json")
		seed := `{"version":1,"goals":[{"id":"` + goalID + `","title":"G","body":""}]}`
		if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
			t.Fatal(err)
		}
		if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
			t.Fatalf("seed import: %d", code)
		}

		exported := captureStdout(t, func() int {
			return run([]string{"-C", dir, "seed", "export"})
		})
		if strings.Contains(exported, `"exported_at_commit"`) {
			t.Fatalf("non-git export should omit exported_at_commit: %s", exported)
		}

		withSHA := strings.TrimSpace(exported)
		withSHA = strings.TrimSuffix(withSHA, "}")
		withSHA += `,"exported_at_commit":"deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"}`
		reimportPath := filepath.Join(dir, "with-sha.json")
		if err := os.WriteFile(reimportPath, []byte(withSHA+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if code := run([]string{"-C", dir, "seed", "import", reimportPath}); code != exitOK {
			t.Fatalf("re-import with exported_at_commit: %d", code)
		}
	})
}

func seedImportRoundTripFixture() string {
	const (
		goalID    = "11111111-1111-1111-1111-111111111111"
		taskID    = "22222222-2222-2222-2222-222222222222"
		decID     = "33333333-3333-3333-3333-333333333333"
		phaseID   = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01"
		scopeID   = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1"
		deepID    = "cccccccc-cccc-cccc-cccc-cccccccccc01"
		deepOldID = "cccccccc-cccc-cccc-cccc-cccccccccc02"
		findID    = "dddddddd-dddd-dddd-dddd-dddddddddd01"
	)
	return `{
  "version": 1,
  "goals": [{"id":"` + goalID + `","title":"G","body":"goal body"}],
  "tasks": [{"id":"` + taskID + `","title":"T","body":"task body","goal_id":"` + goalID + `"}],
  "decisions": [{"id":"` + decID + `","title":"D","body":""}],
  "assumptions": [],
  "discoveries": [],
  "plan_changes": [],
  "claims": [],
  "evidence": [],
  "links": [{"rel":"decision_affects_task","from":"` + decID + `","to":"` + taskID + `"}],
  "findings": [{
    "id":"` + findID + `",
    "decision_id":"` + decID + `",
    "impact_class":"CAUTION",
    "kind":"AFFECTED_WORK",
    "body":"watch"
  }],
  "alternatives": [],
  "plan_phases": [{
    "id":"` + phaseID + `",
    "goal_id":"` + goalID + `",
    "title":"Phase 1",
    "body":"",
    "ord": 0,
    "status":"ACTIVE"
  }],
  "plan_scopes": [{
    "id":"` + scopeID + `",
    "phase_id":"` + phaseID + `",
    "title":"Scope 1",
    "body":"",
    "ord": 0,
    "status":"ACTIVE",
    "auto_replan_count": 2
  }],
  "scope_deep_plans": [
    {
      "id":"` + deepOldID + `",
      "scope_id":"` + scopeID + `",
      "content_json":"{\"rev\":1}",
      "status":"SUPERSEDED"
    },
    {
      "id":"` + deepID + `",
      "scope_id":"` + scopeID + `",
      "content_json":"{\"rev\":2}",
      "status":"ACTIVE"
    }
  ],
  "goal_plan_state": [{
    "goal_id":"` + goalID + `",
    "current_scope_id":"` + scopeID + `"
  }]
}`
}

// TestSeedImportIdempotent locks DF-81: second import of same file exits 0 with stable counts.
func TestSeedImportIdempotent(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedPath := filepath.Join(dir, "seed.json")
	if err := os.WriteFile(seedPath, []byte(seedImportRoundTripFixture()), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("first import: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	links1, err := st.ListLinksByRel("decision_affects_task")
	if err != nil {
		t.Fatal(err)
	}
	phases1, err := st.ListAllPlanPhases()
	if err != nil {
		t.Fatal(err)
	}
	findings1, err := st.ListAllDecisionImpactFindings()
	if err != nil {
		t.Fatal(err)
	}
	_ = st.Close()

	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("second import: %d", code)
	}
	st, err = store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	links2, err := st.ListLinksByRel("decision_affects_task")
	if err != nil {
		t.Fatal(err)
	}
	if len(links2) != len(links1) || len(links2) != 1 {
		t.Fatalf("link count changed: before=%d after=%d", len(links1), len(links2))
	}
	phases2, err := st.ListAllPlanPhases()
	if err != nil {
		t.Fatal(err)
	}
	if len(phases2) != len(phases1) {
		t.Fatalf("plan phase count changed: before=%d after=%d", len(phases1), len(phases2))
	}
	findings2, err := st.ListAllDecisionImpactFindings()
	if err != nil {
		t.Fatal(err)
	}
	if len(findings2) != len(findings1) {
		t.Fatalf("finding count changed: before=%d after=%d", len(findings1), len(findings2))
	}
}

// TestSeedImportDuplicateLinksNoOp locks DF-81: duplicate link endpoints do not error or multiply rows.
func TestSeedImportDuplicateLinksNoOp(t *testing.T) {
	const (
		decID  = "33333333-3333-3333-3333-333333333333"
		taskID = "22222222-2222-2222-2222-222222222222"
	)
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedPath := filepath.Join(dir, "seed.json")
	seed := `{
  "version": 1,
  "decisions": [{"id":"` + decID + `","title":"D","body":""}],
  "tasks": [{"id":"` + taskID + `","title":"T","body":""}],
  "links": [{"rel":"decision_affects_task","from":"` + decID + `","to":"` + taskID + `"}]
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("first import: %d", code)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("second import: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	links, err := st.ListLinksByRel("decision_affects_task")
	if err != nil || len(links) != 1 {
		t.Fatalf("expected 1 link, got %+v err=%v", links, err)
	}
}

// TestSeedImportSameIdLastWins locks DF-83: later import overwrites body/title/plan fields for same UUID.
func TestSeedImportSameIdLastWins(t *testing.T) {
	const (
		goalID  = "11111111-1111-1111-1111-111111111111"
		phaseID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01"
	)
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedA := filepath.Join(dir, "a.json")
	seedB := filepath.Join(dir, "b.json")
	fileA := `{
  "version": 1,
  "goals": [{"id":"` + goalID + `","title":"Title A","body":"body A"}],
  "plan_phases": [{
    "id":"` + phaseID + `",
    "goal_id":"` + goalID + `",
    "title":"Phase A",
    "body":"phase body A",
    "ord": 0,
    "status":"ACTIVE"
  }]
}`
	fileB := `{
  "version": 1,
  "goals": [{"id":"` + goalID + `","title":"Title B","body":"body B"}],
  "plan_phases": [{
    "id":"` + phaseID + `",
    "goal_id":"` + goalID + `",
    "title":"Phase B",
    "body":"phase body B",
    "ord": 1,
    "status":"ACTIVE"
  }]
}`
	if err := os.WriteFile(seedA, []byte(fileA), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(seedB, []byte(fileB), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedA}); code != exitOK {
		t.Fatalf("import A: %d", code)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedB}); code != exitOK {
		t.Fatalf("import B: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	g, err := st.GetGoal(goalID)
	if err != nil || g.Title != "Title B" || g.Body != "body B" {
		t.Fatalf("goal last-wins: %+v err=%v", g, err)
	}
	ph, err := st.GetPlanPhase(phaseID)
	if err != nil || ph.Title != "Phase B" || ph.Body != "phase body B" || ph.Ord != 1 {
		t.Fatalf("plan phase last-wins: %+v err=%v", ph, err)
	}
}

// TestSeedImportPlanTreeIdempotent locks DF-84: plan tree re-import preserves ids without PK errors.
func TestSeedImportPlanTreeIdempotent(t *testing.T) {
	const (
		goalID  = "11111111-1111-1111-1111-111111111111"
		phaseID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01"
		scopeID = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1"
		deepID  = "cccccccc-cccc-cccc-cccc-cccccccccc01"
	)
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	seedPath := filepath.Join(dir, "plan.json")
	seed := `{
  "version": 1,
  "goals": [{"id":"` + goalID + `","title":"G","body":""}],
  "plan_phases": [{
    "id":"` + phaseID + `",
    "goal_id":"` + goalID + `",
    "title":"Phase 1",
    "body":"",
    "ord": 0,
    "status":"ACTIVE"
  }],
  "plan_scopes": [{
    "id":"` + scopeID + `",
    "phase_id":"` + phaseID + `",
    "title":"Scope 1",
    "body":"",
    "ord": 0,
    "status":"ACTIVE",
    "auto_replan_count": 1
  }],
  "scope_deep_plans": [{
    "id":"` + deepID + `",
    "scope_id":"` + scopeID + `",
    "content_json":"{\"k\":1}",
    "status":"ACTIVE"
  }],
  "goal_plan_state": [{
    "goal_id":"` + goalID + `",
    "current_scope_id":"` + scopeID + `"
  }]
}`
	if err := os.WriteFile(seedPath, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("first import: %d", code)
	}
	if code := run([]string{"-C", dir, "seed", "import", seedPath}); code != exitOK {
		t.Fatalf("second import: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if ph, err := st.GetPlanPhase(phaseID); err != nil || ph.ID != phaseID {
		t.Fatalf("plan phase: %+v err=%v", ph, err)
	}
	if sc, err := st.GetPlanScope(scopeID); err != nil || sc.ID != scopeID || sc.AutoReplanCount != 1 {
		t.Fatalf("plan scope: %+v err=%v", sc, err)
	}
	if dp, err := st.GetScopeDeepPlan(deepID); err != nil || dp.ContentJSON != `{"k":1}` {
		t.Fatalf("deep plan: %+v err=%v", dp, err)
	}
	gps, err := st.GetGoalPlanState(goalID)
	if err != nil || gps.CurrentScopeID == nil || *gps.CurrentScopeID != scopeID {
		t.Fatalf("goal plan state: %+v err=%v", gps, err)
	}
}

// portableGraphSeedJSON is a committed-shape seed for trace/graph.json (plan tree + causal links).
func portableGraphSeedJSON() string {
	const (
		goalID  = "11111111-1111-1111-1111-111111111111"
		taskID  = "22222222-2222-2222-2222-222222222222"
		decID   = "33333333-3333-3333-3333-333333333333"
		phaseID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01"
		scopeID = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1"
		deepID  = "cccccccc-cccc-cccc-cccc-cccccccccc01"
	)
	return `{
  "version": 1,
  "goals": [{"id":"` + goalID + `","title":"Portable goal","body":"goal body"}],
  "tasks": [{"id":"` + taskID + `","title":"Portable task","body":"task body","goal_id":"` + goalID + `"}],
  "decisions": [{"id":"` + decID + `","title":"Portable decision","body":""}],
  "assumptions": [],
  "discoveries": [],
  "plan_changes": [],
  "claims": [],
  "evidence": [],
  "links": [{"rel":"decision_affects_task","from":"` + decID + `","to":"` + taskID + `"}],
  "findings": [],
  "alternatives": [],
  "plan_phases": [{
    "id":"` + phaseID + `",
    "goal_id":"` + goalID + `",
    "title":"Phase 1",
    "body":"",
    "ord": 0,
    "status":"ACTIVE"
  }],
  "plan_scopes": [{
    "id":"` + scopeID + `",
    "phase_id":"` + phaseID + `",
    "title":"Scope 1",
    "body":"",
    "ord": 0,
    "status":"ACTIVE",
    "auto_replan_count": 0
  }],
  "scope_deep_plans": [{
    "id":"` + deepID + `",
    "scope_id":"` + scopeID + `",
    "content_json":"{\"rev\":1}",
    "status":"ACTIVE"
  }],
  "goal_plan_state": [{
    "goal_id":"` + goalID + `",
    "current_scope_id":"` + scopeID + `"
  }]
}`
}

// TestPortableGraphTwoCloneWhyContextPlan locks Phase 17 two-clone git-JSON recipe:
// two independent dirs, no shared .trace/; init → import trace/graph.json → index → why + context + plan show.
func TestPortableGraphTwoCloneWhyContextPlan(t *testing.T) {
	const (
		goalID = "11111111-1111-1111-1111-111111111111"
		taskID = "22222222-2222-2222-2222-222222222222"
		decID  = "33333333-3333-3333-3333-333333333333"
	)

	sourceDir := t.TempDir()
	if code := run([]string{"-C", sourceDir, "init"}); code != exitOK {
		t.Fatalf("source init: %d", code)
	}
	sourceDB := filepath.Join(sourceDir, ".trace", "trace.db")
	sourceInfoBefore, err := os.Stat(sourceDB)
	if err != nil {
		t.Fatalf("source db stat: %v", err)
	}

	graphDir := filepath.Join(sourceDir, "trace")
	if err := os.MkdirAll(graphDir, 0o755); err != nil {
		t.Fatal(err)
	}
	graphPath := filepath.Join(graphDir, "graph.json")
	if err := os.WriteFile(graphPath, []byte(portableGraphSeedJSON()), 0o644); err != nil {
		t.Fatal(err)
	}
	graphBytes, err := os.ReadFile(graphPath)
	if err != nil {
		t.Fatal(err)
	}

	runClone := func(name string) string {
		t.Helper()
		clone := t.TempDir()
		if _, err := os.Stat(filepath.Join(clone, ".trace")); !os.IsNotExist(err) {
			t.Fatalf("%s: .trace must not exist before init", name)
		}
		cloneGraphDir := filepath.Join(clone, "trace")
		if err := os.MkdirAll(cloneGraphDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(cloneGraphDir, "graph.json"), graphBytes, 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(clone, "sample.js"), []byte("export function sample() { return 1 }\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if code := run([]string{"-C", clone, "init"}); code != exitOK {
			t.Fatalf("%s init: %d", name, code)
		}
		cloneDB := filepath.Join(clone, ".trace", "trace.db")
		if cloneDB == sourceDB {
			t.Fatalf("%s: clone must not share source .trace/", name)
		}
		if code := run([]string{"-C", clone, "seed", "import", "trace/graph.json"}); code != exitOK {
			t.Fatalf("%s seed import: %d", name, code)
		}
		if code := run([]string{"-C", clone, "index", "sample.js"}); code != exitOK {
			t.Fatalf("%s index: %d", name, code)
		}
		planOut := captureStdout(t, func() int {
			return run([]string{"-C", clone, "plan", "show", "--goal", goalID})
		})
		if !strings.Contains(planOut, goalID) || !strings.Contains(planOut, "Phase 1") {
			t.Fatalf("%s plan show: %s", name, planOut)
		}
		whyOut := captureStdout(t, func() int {
			return run([]string{"-C", clone, "why", "decision", decID})
		})
		var why map[string]any
		if err := json.Unmarshal([]byte(whyOut), &why); err != nil {
			t.Fatalf("%s why decision json: %v\n%s", name, err, whyOut)
		}
		if why["seed_id"] != decID {
			t.Fatalf("%s why decision missing seed_id: %v", name, why)
		}
		steps, ok := why["steps"].([]any)
		if !ok || len(steps) == 0 {
			t.Fatalf("%s why decision missing steps: %v", name, why)
		}
		ctxOut := captureStdout(t, func() int {
			return run([]string{"-C", clone, "context", taskID})
		})
		if !strings.Contains(ctxOut, `"items"`) {
			t.Fatalf("%s context: %s", name, ctxOut)
		}
		return clone
	}

	_ = runClone("cloneA")
	_ = runClone("cloneB")

	sourceInfoAfter, err := os.Stat(sourceDB)
	if err != nil {
		t.Fatalf("source db stat after clones: %v", err)
	}
	if !sourceInfoBefore.ModTime().Equal(sourceInfoAfter.ModTime()) || sourceInfoBefore.Size() != sourceInfoAfter.Size() {
		t.Fatal("source .trace/ must not be touched by clone workflow")
	}
}

func TestCLIWhyUncertainty(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "blocked work"})
	if err != nil {
		t.Fatal(err)
	}
	u, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "API shape unknown",
		Severity: store.UncertaintySeverityBlocking,
		TaskID:   task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	whyOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "why", "uncertainty", u.ID})
	})
	var why map[string]any
	if err := json.Unmarshal([]byte(whyOut), &why); err != nil {
		t.Fatalf("why json: %v\n%s", err, whyOut)
	}
	steps, ok := why["steps"].([]any)
	if !ok || len(steps) == 0 {
		t.Fatalf("why missing steps: %v", why)
	}
	step0, ok := steps[0].(map[string]any)
	if !ok {
		t.Fatalf("steps[0] shape: %#v", steps[0])
	}
	if step0["reason_code"] != retrieval.ReasonExactID {
		t.Fatalf("steps[0] reason_code: %#v", step0)
	}
	if step0["entity_type"] != domain.EntityUncertainty || step0["entity_id"] != u.ID {
		t.Fatalf("steps[0] seed: %#v", step0)
	}
}

func TestCLIWhyRegression(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "eval task"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit: "abc1234", ScoresJSON: `{"correctness":0.99}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: task.ID, BaselineID: b.ID, ScoresJSON: `{"correctness":0.50}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	reg, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: out.ID, TaskID: task.ID, Summary: "correctness dropped",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	whyOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "why", "regression", reg.ID})
	})
	var why map[string]any
	if err := json.Unmarshal([]byte(whyOut), &why); err != nil {
		t.Fatalf("why json: %v\n%s", err, whyOut)
	}
	steps, ok := why["steps"].([]any)
	if !ok || len(steps) == 0 {
		t.Fatalf("why missing steps: %v", why)
	}
	foundSource := false
	for _, s := range steps {
		step := s.(map[string]any)
		if step["entity_type"] == domain.EntityOutcomeResult && step["entity_id"] == out.ID {
			foundSource = true
			break
		}
	}
	if !foundSource {
		t.Fatalf("why missing linked outcome_result step: %#v", steps)
	}
}
