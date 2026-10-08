package mcp_test

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	tracemcp "github.com/mrchatam/Trace/internal/mcp"
	"github.com/mrchatam/Trace/internal/store"
)

// assertLowerKeys walks a decoded JSON value and fails on any object key that
// contains an uppercase rune.
//
// This is the guard for the response-key-shape class that produced closed #116,
// closed #105 and open #161: an untagged struct slipped into a response and
// marshalled its Go field names. One test now covers every tool.
func assertLowerKeys(t *testing.T, label string, v any) {
	t.Helper()
	walkJSON(t, label, "$", v)
}

func walkJSON(t *testing.T, label, path string, v any) {
	t.Helper()
	switch node := v.(type) {
	case map[string]any:
		for k, child := range node {
			for _, r := range k {
				if r >= 'A' && r <= 'Z' {
					t.Errorf("%s: key %q at %s contains an uppercase rune %q", label, k, path, r)
				}
			}
			walkJSON(t, label, path+"."+k, child)
		}
	case []any:
		for i, child := range node {
			walkJSON(t, label, path+"["+strconv.Itoa(i)+"]", child)
		}
	}
}

// keyShapeFixture holds the ids the sweep needs after seeding.
type keyShapeFixture struct {
	taskIDs  []string
	goalID   string
	changeID string
	scopeID  string
}

// TestNoToolResponseKeyContainsUppercase sweeps every tool reachable on a
// seeded project and fails if any response leaks a Go field name as a key.
func TestNoToolResponseKeyContainsUppercase(t *testing.T) {
	dir, ctx := writeFixture(t)
	fx := seedKeyShapeFixture(t, dir)
	srv := addServer(t, dir)

	calls := []struct {
		tool string
		call func() (any, any, error)
	}{
		{"trace_version", func() (any, any, error) {
			return tracemcp.CallVersion(srv, ctx, tracemcp.VersionInput{})
		}},
		{"trace_tasks", func() (any, any, error) {
			return tracemcp.CallTasks(srv, ctx, tracemcp.TasksInput{})
		}},
		{"trace_tasks filtered", func() (any, any, error) {
			return tracemcp.CallTasks(srv, ctx, tracemcp.TasksInput{WorkState: "PENDING"})
		}},
		{"trace_capability list", func() (any, any, error) {
			return tracemcp.CallCapability(srv, ctx, tracemcp.CapabilityInput{Action: "list"})
		}},
		{"trace_capability missing", func() (any, any, error) {
			return tracemcp.CallCapability(srv, ctx, tracemcp.CapabilityInput{
				Action: "missing", TaskID: fx.taskIDs[0],
			})
		}},
		{"trace_search", func() (any, any, error) {
			return tracemcp.CallSearch(srv, ctx, tracemcp.SearchInput{Query: "keyshape"})
		}},
		{"trace_search ops", func() (any, any, error) {
			return tracemcp.CallSearch(srv, ctx, tracemcp.SearchInput{
				Query: "keyshape OR plan", Matcher: "ops",
			})
		}},
		{"trace_changes list", func() (any, any, error) {
			return tracemcp.CallChanges(srv, ctx, tracemcp.ChangesInput{Action: "list"})
		}},
		{"trace_changes show", func() (any, any, error) {
			return tracemcp.CallChanges(srv, ctx, tracemcp.ChangesInput{
				Action: "show", ChangeID: fx.changeID,
			})
		}},
		{"trace_regressions", func() (any, any, error) {
			return tracemcp.CallRegressions(srv, ctx, tracemcp.RegressionsInput{Action: "list"})
		}},
		{"trace_agents list", func() (any, any, error) {
			return tracemcp.CallAgents(srv, ctx, tracemcp.AgentsInput{Action: "list"})
		}},
		{"trace_agents describe", func() (any, any, error) {
			return tracemcp.CallAgents(srv, ctx, tracemcp.AgentsInput{
				Action: "describe", Slug: "agent:code-reviewer",
			})
		}},
		{"trace_why", func() (any, any, error) {
			return tracemcp.CallWhy(srv, ctx, tracemcp.WhyInput{
				EntityType: "task", ID: fx.taskIDs[0],
			})
		}},
		{"trace_impact walk", func() (any, any, error) {
			return tracemcp.CallImpact(srv, ctx, tracemcp.ImpactInput{
				Action: "walk", Seeds: []string{"file:pkg/keyshape.go"},
			})
		}},
		{"trace_link api-contract", func() (any, any, error) {
			return tracemcp.CallLink(srv, ctx, tracemcp.LinkInput{
				Rel: "api-contract", From: fx.taskIDs[0], To: fx.taskIDs[1],
			})
		}},
		{"trace_link scope-member", func() (any, any, error) {
			return tracemcp.CallLink(srv, ctx, tracemcp.LinkInput{
				Rel: "scope-member", From: fx.taskIDs[0], To: fx.scopeID,
			})
		}},
		{"trace_loop gate", func() (any, any, error) {
			return tracemcp.CallLoop(srv, ctx, tracemcp.LoopInput{
				Action: "gate", TaskID: fx.taskIDs[0], For: "edit",
			})
		}},
		{"trace_loop next", func() (any, any, error) {
			return tracemcp.CallLoop(srv, ctx, tracemcp.LoopInput{
				Action: "next", TaskID: fx.taskIDs[0],
			})
		}},
		{"trace_plan show", func() (any, any, error) {
			return tracemcp.CallPlan(srv, ctx, tracemcp.PlanInput{
				Action: "show", GoalID: fx.goalID,
			})
		}},
		{"trace_explore", func() (any, any, error) {
			return tracemcp.CallExplore(srv, ctx, tracemcp.ExploreInput{
				TaskID: fx.taskIDs[0], Query: "keyshape",
			})
		}},
	}

	for _, c := range calls {
		t.Run(c.tool, func(t *testing.T) {
			res, _, err := c.call()
			if err != nil {
				t.Fatalf("%s: %v", c.tool, err)
			}
			raw := tracemcp.ResultText(t, res)
			var v any
			if err := json.Unmarshal([]byte(raw), &v); err != nil {
				t.Fatalf("%s: response is not JSON: %v\n%s", c.tool, err, raw)
			}
			assertLowerKeys(t, c.tool, v)
		})
	}
}

// seedKeyShapeFixture writes the rows the sweep needs so every tool returns a
// real response rather than an error path, then releases the store lock.
func seedKeyShapeFixture(t *testing.T, dir string) keyShapeFixture {
	t.Helper()

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()

	// A goal with a plan tree, so trace_loop next and trace_plan show build.
	goalID, taskID := seedMCPGoalTaskPlan(t, st)

	// A second task, so api-contract (task -> task) has two real endpoints.
	second, err := svc.CreateTask(ctx, domain.TaskInput{
		Title: "keyshape second task", GoalID: &goalID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// A change with a path, for trace_changes show and search hits.
	change, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: taskID,
		Reason: "keyshape change body",
		Paths:  []domain.ChangePathInput{{Path: "pkg/keyshape.go", Status: "M"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	// A regression, so trace_regressions is non-empty.
	if _, err := st.UpsertRegression(store.Regression{
		TaskID:      second.ID,
		SourceKind:  store.RegressionSourceEvaluation,
		SourceID:    "keyshape-outcome",
		Dimension:   "overall",
		Attribution: store.RegressionAttributionCorrelated,
		Summary:     "keyshape regression",
	}); err != nil {
		t.Fatal(err)
	}

	// A capability, so trace_capability has a row.
	if _, err := st.UpsertCapability(store.Capability{
		Kind: "MCP", Slug: "mcp:demo", Title: "Demo", Status: store.CapabilityStatusAvailable,
	}); err != nil {
		t.Fatal(err)
	}

	// An indexed file, so trace_impact walk resolves its seed.
	if _, err := st.UpsertFile("pkg/keyshape.go", "hash-keyshape", nil); err != nil {
		t.Fatal(err)
	}
	if err := st.SyncFileFTS("pkg/keyshape.go"); err != nil {
		t.Fatal(err)
	}

	// A scope, so scope-member (any entity -> scope) has a real endpoint.
	scope, err := svc.CreateScope(ctx, domain.ScopeInput{Slug: "keyshape-scope", Title: "keyshape scope", Kind: "feature"})
	if err != nil {
		t.Fatal(err)
	}

	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	// The harness agent catalog, so trace_agents describe resolves a slug.
	if err := installAgentCatalog(t, dir); err != nil {
		t.Fatal(err)
	}

	return keyShapeFixture{
		taskIDs:  []string{taskID, second.ID},
		goalID:   goalID,
		changeID: change.ID,
		scopeID:  scope.ID,
	}
}

// TestKeyShapeGuardCatchesPascalCase is the guard's own keeper: it proves the
// walker actually fails on a PascalCase key, so a future regression cannot pass
// because the guard silently does nothing.
func TestKeyShapeGuardCatchesPascalCase(t *testing.T) {
	fake := &testing.T{}
	assertLowerKeys(fake, "probe", map[string]any{"Actor": "x"})
	if !fake.Failed() {
		t.Fatal("the key-shape guard must fail on a PascalCase key")
	}
}
