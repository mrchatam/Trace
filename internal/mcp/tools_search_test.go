package mcp_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	tracemcp "github.com/mrchatam/Trace/internal/mcp"
	"github.com/mrchatam/Trace/internal/store"
)

// searchFixture creates a project with two rows a boolean query can separate.
func searchFixture(t *testing.T) (string, context.Context) {
	t.Helper()
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	if _, err := svc.CreateTask(ctx, domain.TaskInput{Title: "zeppelin cache layer"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateTask(ctx, domain.TaskInput{Title: "unrelated widget"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateTask(ctx, domain.TaskInput{Title: "cache eviction policy"}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	return dir, ctx
}

func decodeSearch(t *testing.T, raw string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatalf("search json: %v\n%s", err, raw)
	}
	return m
}

// TestSearchAndModeNotesOperatorRewrite is the #149 proof at the MCP surface:
// the default mode still returns zero hits, but it now says why and names the
// opt-in instead of failing silently.
func TestSearchAndModeNotesOperatorRewrite(t *testing.T) {
	dir, ctx := searchFixture(t)
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})

	res, _, err := tracemcp.CallSearch(srv, ctx, tracemcp.SearchInput{Query: "zeppelin OR cache"})
	if err != nil {
		t.Fatalf("trace_search: %v", err)
	}
	got := decodeSearch(t, tracemcp.ResultText(t, res))
	if int(got["count"].(float64)) != 0 {
		t.Fatalf("AND mode should still return 0 hits, got %v", got)
	}
	note, _ := got["query_note"].(string)
	if !strings.Contains(note, "matcher=ops") {
		t.Fatalf("query_note must name the opt-in, got %q", note)
	}
	if got["matcher"] != "and" {
		t.Fatalf("matcher = %v, want and", got["matcher"])
	}
}

// TestSearchOpsModeReturnsUnion proves matcher=ops turns the same query into a
// union instead of an empty result.
func TestSearchOpsModeReturnsUnion(t *testing.T) {
	dir, ctx := searchFixture(t)
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})

	res, _, err := tracemcp.CallSearch(srv, ctx, tracemcp.SearchInput{Query: "zeppelin OR cache", Matcher: "ops"})
	if err != nil {
		t.Fatalf("trace_search: %v", err)
	}
	got := decodeSearch(t, tracemcp.ResultText(t, res))
	if int(got["count"].(float64)) != 2 {
		t.Fatalf("ops mode want 2 hits, got %v", got)
	}
	if note, _ := got["query_note"].(string); note != "" {
		t.Fatalf("a honored boolean query must not warn, got %q", note)
	}
	if got["matcher"] != "ops" {
		t.Fatalf("matcher = %v, want ops", got["matcher"])
	}
}

// TestSearchRejectsUnknownMatcher keeps the wire vocabulary closed.
func TestSearchRejectsUnknownMatcher(t *testing.T) {
	dir, ctx := searchFixture(t)
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	_, _, err := tracemcp.CallSearch(srv, ctx, tracemcp.SearchInput{Query: "cache", Matcher: "regex"})
	if err == nil || !strings.Contains(err.Error(), "matcher must be and or ops") {
		t.Fatalf("want a closed-vocabulary error, got %v", err)
	}
}

// TestSearchResponseEchoesEffectiveLimit proves a clamped request is visible.
func TestSearchResponseEchoesEffectiveLimit(t *testing.T) {
	dir, ctx := searchFixture(t)
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	res, _, err := tracemcp.CallSearch(srv, ctx, tracemcp.SearchInput{Query: "cache", Limit: 5000})
	if err != nil {
		t.Fatalf("trace_search: %v", err)
	}
	got := decodeSearch(t, tracemcp.ResultText(t, res))
	if int(got["limit"].(float64)) != 64 {
		t.Fatalf("limit = %v, want the effective 64", got["limit"])
	}
}

// TestSearchResponseCarriesGraphSyncHonesty is the #153 proof: a read surface
// that does not compile a packet still reports index staleness.
func TestSearchResponseCarriesGraphSyncHonesty(t *testing.T) {
	dir, ctx := searchFixture(t)
	// A non-git project has no HEAD, so the block is legitimately absent.
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	res, _, err := tracemcp.CallSearch(srv, ctx, tracemcp.SearchInput{Query: "cache"})
	if err != nil {
		t.Fatalf("trace_search: %v", err)
	}
	got := decodeSearch(t, tracemcp.ResultText(t, res))
	if _, present := got["graph_sync_honesty"]; present {
		t.Fatalf("non-git project must not claim staleness: %v", got)
	}
}

// TestExploreRejectsOutOfRange mirrors the trace_context validation model that
// #158 asks trace_explore to copy.
func TestExploreRejectsOutOfRange(t *testing.T) {
	dir, ctx := searchFixture(t)
	svc := domain.New(mustOpenStore(t, dir))
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "explore bounds task"})
	if err != nil {
		t.Fatal(err)
	}
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})

	for _, tc := range []struct {
		name string
		in   tracemcp.ExploreInput
		want string
	}{
		{"negative limit", tracemcp.ExploreInput{TaskID: task.ID, Limit: -1}, "limit must be 1..64"},
		{"negative max_nodes", tracemcp.ExploreInput{TaskID: task.ID, MaxNodes: -5}, "max_nodes must be 1.."},
		{"negative depth", tracemcp.ExploreInput{TaskID: task.ID, Depth: -2}, "depth must be 1..8"},
		{"unknown matcher", tracemcp.ExploreInput{TaskID: task.ID, Matcher: "regex"}, "matcher must be and or ops"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.in.Project = dir
			_, _, err := tracemcp.CallExplore(srv, ctx, tc.in)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want %q, got %v", tc.want, err)
			}
		})
	}
}

func mustOpenStore(t *testing.T, dir string) *store.Store {
	t.Helper()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}
