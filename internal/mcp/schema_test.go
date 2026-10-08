package mcp_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/domain"

	tracemcp "github.com/mrchatam/Trace/internal/mcp"
	"github.com/mrchatam/Trace/internal/store"
)

// toolListSchemas connects an in-memory client to the server and returns the
// advertised input schema per tool name, exactly as an agent would see it.
func toolListSchemas(t *testing.T) map[string]map[string]any {
	t.Helper()
	ctx := context.Background()
	srv := tracemcp.NewServer(tracemcp.Options{})
	ct, st := sdkmcp.NewInMemoryTransports()

	serverSession, err := srv.MCP().Connect(ctx, st, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "probe", Version: "0"}, nil)
	clientSession, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()
	defer serverSession.Close()

	res, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("tools/list: %v", err)
	}
	out := map[string]map[string]any{}
	for _, tool := range res.Tools {
		raw, err := json.Marshal(tool.InputSchema)
		if err != nil {
			t.Fatalf("marshal schema for %s: %v", tool.Name, err)
		}
		var m map[string]any
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatalf("decode schema for %s: %v", tool.Name, err)
		}
		out[tool.Name] = m
	}
	return out
}

// TestToolSchemasCarryEnumsAndRanges is the #158 proof: the advertised schemas
// must carry real enum/range keywords, not only free-text descriptions.
func TestToolSchemasCarryEnumsAndRanges(t *testing.T) {
	schemas := toolListSchemas(t)
	if len(schemas) != 17 {
		t.Fatalf("want 17 tools, got %d", len(schemas))
	}

	enumCases := []struct {
		tool, prop string
		want       []string
	}{
		{"trace_loop", "action", []string{"next", "apply", "status", "gate"}},
		{"trace_capability", "kind", []string{"SKILL", "RULE", "MCP", "TOOL", "HOOK", "AGENT"}},
		{"trace_capability", "action", []string{"declare", "list", "require", "unrequire", "missing"}},
		{"trace_review", "result", []string{"PASS", "FAIL", "UNCERTAIN"}},
		{"trace_search", "matcher", []string{"and", "ops"}},
		{"trace_explore", "matcher", []string{"and", "ops"}},
		{"trace_tasks", "work_state", []string{"PENDING", "IN_PROGRESS", "AWAITING_REVIEW", "BLOCKED", "FAILED", "DONE", "STALE", "SKIPPED"}},
		{"trace_transition", "to_state", []string{"DONE", "BLOCKED", "FAILED"}},
		{"trace_context", "format", []string{"json", "markdown", "both"}},
	}
	for _, tc := range enumCases {
		t.Run(tc.tool+"."+tc.prop, func(t *testing.T) {
			got := propSchema(t, schemas[tc.tool], tc.prop)
			rawEnum, ok := got["enum"].([]any)
			if !ok {
				t.Fatalf("%s.%s has no enum: %v", tc.tool, tc.prop, got)
			}
			for _, want := range tc.want {
				found := false
				for _, v := range rawEnum {
					if v == want {
						found = true
					}
				}
				if !found {
					t.Fatalf("%s.%s enum missing %q: %v", tc.tool, tc.prop, want, rawEnum)
				}
			}
		})
	}

	rangeCases := []struct {
		tool, prop string
		min, max   float64
	}{
		{"trace_search", "limit", 1, 64},
		{"trace_tasks", "limit", 1, 500},
		{"trace_explore", "limit", 1, 64},
		{"trace_context", "depth", 1, 2},
		{"trace_context", "max_layer", 1, 3},
	}
	for _, tc := range rangeCases {
		t.Run(tc.tool+"."+tc.prop+" range", func(t *testing.T) {
			got := propSchema(t, schemas[tc.tool], tc.prop)
			if v, ok := got["minimum"].(float64); !ok || v != tc.min {
				t.Fatalf("%s.%s minimum = %v, want %v", tc.tool, tc.prop, got["minimum"], tc.min)
			}
			if v, ok := got["maximum"].(float64); !ok || v != tc.max {
				t.Fatalf("%s.%s maximum = %v, want %v", tc.tool, tc.prop, got["maximum"], tc.max)
			}
		})
	}
}

// TestToolDescriptionsCoverSchemaEnums is the #159 drift guard: every enum
// value in a tool's schema must appear in that tool's description.
func TestToolDescriptionsCoverSchemaEnums(t *testing.T) {
	schemas := toolListSchemas(t)
	descriptions := tracemcp.ToolDescriptions()

	for tool, schema := range schemas {
		desc, ok := descriptions[tool]
		if !ok {
			t.Fatalf("no generated description for %s", tool)
		}
		props, _ := schema["properties"].(map[string]any)
		for prop, raw := range props {
			p, _ := raw.(map[string]any)
			rawEnum, ok := p["enum"].([]any)
			if !ok {
				continue
			}
			for _, v := range rawEnum {
				s, ok := v.(string)
				if !ok || !strings.Contains(desc, s) {
					t.Fatalf("%s.%s enum value %v is missing from the description: %q", tool, prop, v, desc)
				}
			}
		}
	}
}

// TestServerExposesEveryToolDescription guards that the registry and the
// generated descriptions stay in lockstep.
func TestServerExposesEveryToolDescription(t *testing.T) {
	schemas := toolListSchemas(t)
	for _, name := range tracemcp.RegisteredToolNames() {
		if _, ok := schemas[name]; !ok {
			t.Fatalf("registered tool %q is not advertised", name)
		}
	}
	if len(tracemcp.ToolDescriptions()) != len(tracemcp.RegisteredToolNames()) {
		t.Fatal("description count must match the registered tool count")
	}
}

// TestChangeShowRowIsSnakeCase is the #161 proof: the same document must not mix
// snake_case and PascalCase, and RowID must not be on the wire.
func TestChangeShowRowIsSnakeCase(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	taskID := seedEnvelopeFixture(t, dir)

	recRes, _, err := tracemcp.CallChanges(srv, ctx, tracemcp.ChangesInput{
		Action: "list", TaskID: taskID,
	})
	if err != nil {
		t.Fatal(err)
	}
	var listEnv struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(tracemcp.ResultText(t, recRes)), &listEnv); err != nil {
		t.Fatal(err)
	}
	if len(listEnv.Items) == 0 {
		t.Fatal("expected a change row")
	}
	changeID := listEnv.Items[0].ID

	showRes, _, err := tracemcp.CallChanges(srv, ctx, tracemcp.ChangesInput{
		Action: "show", ChangeID: changeID,
	})
	if err != nil {
		t.Fatal(err)
	}
	raw := tracemcp.ResultText(t, showRes)

	var show struct {
		Change map[string]any   `json:"change"`
		Paths  []map[string]any `json:"paths"`
	}
	if err := json.Unmarshal([]byte(raw), &show); err != nil {
		t.Fatalf("show json: %v\n%s", err, raw)
	}
	for _, want := range []string{"id", "task_id", "git_commit", "status", "created_at", "updated_at"} {
		if _, ok := show.Change[want]; !ok {
			t.Fatalf("change row missing snake_case key %q: %v", want, show.Change)
		}
	}
	for key := range show.Change {
		for _, r := range key {
			if r >= 'A' && r <= 'Z' {
				t.Fatalf("PascalCase key %q leaked into change row: %v", key, show.Change)
			}
		}
		if key == "RowID" {
			t.Fatalf("internal rowid must not be on the wire: %v", show.Change)
		}
	}
}

// TestListEnvelopeIsUniform proves #157: the list-shaped responses all carry
// ok/items/count/truncated.
func TestListEnvelopeIsUniform(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	seedEnvelopeFixture(t, dir)

	cases := []struct {
		name string
		call func() (any, any, error)
	}{
		{"trace_changes list", func() (any, any, error) {
			return tracemcp.CallChanges(srv, ctx, tracemcp.ChangesInput{Action: "list"})
		}},
		{"trace_regressions list", func() (any, any, error) {
			return tracemcp.CallRegressions(srv, ctx, tracemcp.RegressionsInput{Action: "list"})
		}},
		{"trace_agents list", func() (any, any, error) {
			return tracemcp.CallAgents(srv, ctx, tracemcp.AgentsInput{Action: "list"})
		}},
		{"trace_agents recommend", func() (any, any, error) {
			return tracemcp.CallAgents(srv, ctx, tracemcp.AgentsInput{Action: "recommend", Phase: "CRITIQUE"})
		}},
		{"trace_tasks", func() (any, any, error) {
			return tracemcp.CallTasks(srv, ctx, tracemcp.TasksInput{})
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, _, err := tc.call()
			if err != nil {
				t.Fatal(err)
			}
			raw := tracemcp.ResultText(t, res)
			var env map[string]any
			if err := json.Unmarshal([]byte(raw), &env); err != nil {
				t.Fatalf("must be an object envelope: %v\n%s", err, raw)
			}
			for _, key := range []string{"ok", "items", "count", "truncated"} {
				if _, ok := env[key]; !ok {
					t.Fatalf("envelope missing %q: %s", key, raw)
				}
			}
			items, ok := env["items"].([]any)
			if !ok {
				t.Fatalf("items must be an array (empty list must be [] not null): %s", raw)
			}
			if int(env["count"].(float64)) != len(items) {
				t.Fatalf("count %v must equal len(items) %d: %s", env["count"], len(items), raw)
			}
		})
	}
}

// TestAgentsRecommendEmptyCarriesHint proves the recommend action surfaces the
// same empty-catalog remedy the list action already did.
func TestAgentsRecommendEmptyCarriesHint(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	res, _, err := tracemcp.CallAgents(srv, ctx, tracemcp.AgentsInput{
		Action: "recommend", Phase: "CRITIQUE",
	})
	if err != nil {
		t.Fatal(err)
	}
	var env struct {
		Hint string `json:"hint"`
	}
	if err := json.Unmarshal([]byte(tracemcp.ResultText(t, res)), &env); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(env.Hint, "trace install agents") {
		t.Fatalf("empty recommendation must carry the install hint, got %q", env.Hint)
	}
}

func propSchema(t *testing.T, schema map[string]any, prop string) map[string]any {
	t.Helper()
	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema has no properties: %v", schema)
	}
	raw, ok := props[prop].(map[string]any)
	if !ok {
		t.Fatalf("schema has no property %q: %v", prop, props)
	}
	return raw
}

// seedEnvelopeFixture writes a goal, a linked task, one change and one
// regression, then releases the store lock so the MCP server can open the same
// project. It returns the task id.
func seedEnvelopeFixture(t *testing.T, dir string) string {
	t.Helper()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()

	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "envelope goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "envelope task", GoalID: &goal.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "envelope change",
		Paths:  []domain.ChangePathInput{{Path: "pkg/envelope.go", Status: "M"}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := st.UpsertRegression(store.Regression{
		TaskID:      task.ID,
		SourceKind:  store.RegressionSourceEvaluation,
		SourceID:    "envelope-outcome",
		Dimension:   "overall",
		Attribution: store.RegressionAttributionCorrelated,
		Summary:     "envelope regression",
	}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	return task.ID
}
