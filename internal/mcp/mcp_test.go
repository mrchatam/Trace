package mcp_test

import (
	"context"
	"encoding/json"
	"errors"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/install"
	"github.com/mrchatam/Trace/internal/loop"
	tracemcp "github.com/mrchatam/Trace/internal/mcp"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func TestTraceWhyAndContextCallLibrary(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}

	svc := domain.New(st)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "T", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkGoalTask(ctx, g.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})

	// Direct library baseline (close before MCP tools open the same root).
	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	libWhy, err := retrieval.New(st2).Why(ctx, domain.EntityTask, task.ID)
	if err != nil {
		st2.Close()
		t.Fatal(err)
	}
	libPkt, err := compiler.New(st2).WithRetrieval(retrieval.New(st2)).TaskContext(ctx, task.ID, compiler.ContextOptions{})
	if err != nil {
		st2.Close()
		t.Fatal(err)
	}
	libJSON, err := libPkt.JSON()
	if err != nil {
		st2.Close()
		t.Fatal(err)
	}
	st2.Close()

	whyRes, _, err := callWhy(srv, ctx, tracemcp.WhyInput{EntityType: domain.EntityTask, ID: task.ID})
	if err != nil {
		t.Fatalf("trace_why: %v", err)
	}
	whyText := mustText(t, whyRes)
	var mcpWhy retrieval.WhyResult
	if err := json.Unmarshal([]byte(whyText), &mcpWhy); err != nil {
		t.Fatalf("why json: %v\n%s", err, whyText)
	}
	if mcpWhy.SeedID != libWhy.SeedID || mcpWhy.SeedType != libWhy.SeedType {
		t.Fatalf("why seed mismatch: mcp=%+v lib=%+v", mcpWhy, libWhy)
	}
	if len(mcpWhy.Steps) == 0 {
		t.Fatal("why: expected steps from library")
	}

	ctxRes, _, err := callContext(srv, ctx, tracemcp.ContextInput{TaskID: task.ID, Depth: 1, Format: "json"})
	if err != nil {
		t.Fatalf("trace_context: %v", err)
	}
	ctxText := mustText(t, ctxRes)
	var mcpPkt map[string]any
	if err := json.Unmarshal([]byte(ctxText), &mcpPkt); err != nil {
		t.Fatalf("context json: %v\n%s", err, ctxText)
	}
	var libMap map[string]any
	if err := json.Unmarshal(libJSON, &libMap); err != nil {
		t.Fatal(err)
	}
	if mcpPkt["task_id"] != libMap["task_id"] {
		t.Fatalf("context task_id: mcp=%v lib=%v", mcpPkt["task_id"], libMap["task_id"])
	}
	if mcpPkt["schema_version"] != libMap["schema_version"] {
		t.Fatalf("context schema_version mismatch")
	}
}

func TestMCPContextQueryMerged(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "MCP query merge task"})
	if err != nil {
		t.Fatal(err)
	}
	const token = "mcpg1querytokenABC"
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Unlinked decision " + token})
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})

	without, _, err := callContext(srv, ctx, tracemcp.ContextInput{TaskID: task.ID, Depth: 1, Format: "json"})
	if err != nil {
		t.Fatalf("trace_context without query: %v", err)
	}
	var pktNoQuery compiler.Packet
	if err := json.Unmarshal([]byte(mustText(t, without)), &pktNoQuery); err != nil {
		t.Fatal(err)
	}
	for _, it := range pktNoQuery.Items {
		if it.EntityID == dec.ID {
			t.Fatalf("unlinked decision must not appear without query")
		}
	}

	with, _, err := callContext(srv, ctx, tracemcp.ContextInput{TaskID: task.ID, Depth: 1, Format: "json", Query: token})
	if err != nil {
		t.Fatalf("trace_context with query: %v", err)
	}
	var pktQuery compiler.Packet
	if err := json.Unmarshal([]byte(mustText(t, with)), &pktQuery); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, it := range pktQuery.Items {
		if it.EntityType == "decision" && it.EntityID == dec.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected query-merged decision in MCP packet; items=%+v", pktQuery.Items)
	}

	_, _, err = callContext(srv, ctx, tracemcp.ContextInput{Query: token})
	if err == nil || !strings.Contains(err.Error(), "task_id is required") {
		t.Fatalf("expected empty task_id rejection, got %v", err)
	}
}

func TestMCPContextMaxLayer2(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "MCP G8 goal"})
	if err != nil {
		t.Fatal(err)
	}
	taskA, err := svc.CreateTask(ctx, domain.TaskInput{Title: "MCP primary"})
	if err != nil {
		t.Fatal(err)
	}
	taskB, err := svc.CreateTask(ctx, domain.TaskInput{Title: "MCP sibling"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkGoalTask(ctx, goal.ID, taskA.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkGoalTask(ctx, goal.ID, taskB.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	res, _, err := callContext(srv, ctx, tracemcp.ContextInput{
		TaskID:   taskA.ID,
		Depth:    2,
		MaxLayer: 2,
		Format:   "json",
	})
	if err != nil {
		t.Fatalf("trace_context max_layer=2: %v", err)
	}
	var pkt compiler.Packet
	if err := json.Unmarshal([]byte(mustText(t, res)), &pkt); err != nil {
		t.Fatal(err)
	}
	var sawL2 bool
	for _, it := range pkt.Items {
		if it.Layer == 2 {
			sawL2 = true
			break
		}
	}
	if !sawL2 {
		t.Fatalf("MCP max_layer=2 expected layer-2 items; got %+v", pkt.Items)
	}
}

func TestParityWriteTools(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	ctx := context.Background()

	addRes, _, err := callAdd(srv, ctx, tracemcp.AddInput{Kind: "goal", Title: "Goal"})
	if err != nil {
		t.Fatal(err)
	}
	goalID := mustOKID(t, mustText(t, addRes))

	addRes, _, err = callAdd(srv, ctx, tracemcp.AddInput{Kind: "task", Title: "Task", GoalID: goalID})
	if err != nil {
		t.Fatal(err)
	}
	taskID := mustOKID(t, mustText(t, addRes))

	_, _, err = callLink(srv, ctx, tracemcp.LinkInput{Rel: "goal-task", From: goalID, To: taskID})
	if err != nil {
		t.Fatal(err)
	}

	addRes, _, err = callAdd(srv, ctx, tracemcp.AddInput{Kind: "discovery", Title: "Disc"})
	if err != nil {
		t.Fatal(err)
	}
	discID := mustOKID(t, mustText(t, addRes))
	_, _, err = callLink(srv, ctx, tracemcp.LinkInput{Rel: "discovery-mentions-task", From: discID, To: taskID})
	if err != nil {
		t.Fatalf("discovery-mentions-task: %v", err)
	}

	revRes, _, err := callReview(srv, ctx, tracemcp.ReviewInput{
		Action: "create", Title: "R", TaskID: taskID,
	})
	if err != nil {
		t.Fatal(err)
	}
	reviewID := mustOKID(t, mustText(t, revRes))

	_, _, err = callReview(srv, ctx, tracemcp.ReviewInput{
		Action: "set", ID: reviewID, Result: "PASS", Reason: "ok",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = callTransition(srv, ctx, tracemcp.TransitionInput{
		TaskID: taskID, ToState: store.WorkStateInProgress, Reason: "start",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = callTransition(srv, ctx, tracemcp.TransitionInput{
		TaskID: taskID, ToState: store.WorkStateDone, Reason: "done via review",
	})
	if err == nil {
		t.Fatal("DONE after PASS without as_operator must fail")
	}
	_, _, err = callTransition(srv, ctx, tracemcp.TransitionInput{
		TaskID: taskID, ToState: store.WorkStateDone, Reason: "done via review", AsOperator: true,
	})
	if err != nil {
		t.Fatalf("DONE after PASS + as_operator: %v", err)
	}
}

// TestAsOperatorSchemaIdentityDocs locks DF-44: as_operator schema is conscious flag≠identity.
func TestAsOperatorSchemaIdentityDocs(t *testing.T) {
	typ := reflect.TypeOf(tracemcp.TransitionInput{})
	f, ok := typ.FieldByName("AsOperator")
	if !ok {
		t.Fatal("TransitionInput.AsOperator missing")
	}
	tag := f.Tag.Get("jsonschema")
	if !strings.Contains(tag, "flag≠identity") && !strings.Contains(tag, "not verified") {
		t.Fatalf("as_operator jsonschema must mention flag≠identity / not verified: %q", tag)
	}
	if !strings.Contains(tag, "conscious") {
		t.Fatalf("as_operator jsonschema must mention conscious: %q", tag)
	}
}

func TestTransitionAllowDoneEmitsWarning(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	ctx := context.Background()

	addRes, _, err := callAdd(srv, ctx, tracemcp.AddInput{Kind: "task", Title: "hatch"})
	if err != nil {
		t.Fatal(err)
	}
	taskID := mustOKID(t, mustText(t, addRes))
	_, _, err = callTransition(srv, ctx, tracemcp.TransitionInput{
		TaskID: taskID, ToState: store.WorkStateInProgress, Reason: "start",
	})
	if err != nil {
		t.Fatal(err)
	}
	res, _, err := callTransition(srv, ctx, tracemcp.TransitionInput{
		TaskID: taskID, ToState: store.WorkStateDone, Reason: "escape", AllowDone: true,
	})
	if err != nil {
		t.Fatalf("allow_done: %v", err)
	}
	text := mustText(t, res)
	var env map[string]any
	if err := json.Unmarshal([]byte(text), &env); err != nil {
		t.Fatal(err)
	}
	warn, _ := env["warning"].(string)
	if warn == "" || !strings.Contains(warn, "allow_done") {
		t.Fatalf("expected non-empty allow_done warning, got %#v", env)
	}
	if !strings.Contains(warn, "allow_missing_caps") && !strings.Contains(warn, "missing capabilities") {
		t.Fatalf("warning must mention missing-caps independence, got %q", warn)
	}
}

func TestCapabilityMissingRequiresTaskParam(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	ctx := context.Background()
	_, _, err = callCapability(srv, ctx, tracemcp.CapabilityInput{Action: "missing"})
	if err == nil {
		t.Fatal("missing without task must fail")
	}
	msg := err.Error()
	if !strings.Contains(msg, "task") || !strings.Contains(msg, "trace_tasks") {
		t.Fatalf("error should name task/task_id and hint trace_tasks: %q", msg)
	}
}

func TestImportBoundaryNoLibraryImportsMCP(t *testing.T) {
	root := moduleRoot(t)
	forbiddenImporters := []string{
		"github.com/mrchatam/Trace/internal/store",
		"github.com/mrchatam/Trace/internal/vcs",
		"github.com/mrchatam/Trace/internal/gitcli",
		"github.com/mrchatam/Trace/internal/analyzers",
		"github.com/mrchatam/Trace/internal/domain",
		"github.com/mrchatam/Trace/internal/retrieval",
		"github.com/mrchatam/Trace/internal/compiler",
		"github.com/mrchatam/Trace/cmd/trace",
	}
	mcpImport := "github.com/mrchatam/Trace/internal/mcp"
	mcpCmd := "github.com/mrchatam/Trace/cmd/trace-mcp"

	for _, pkgPath := range forbiddenImporters {
		dir := filepath.Join(root, strings.TrimPrefix(pkgPath, "github.com/mrchatam/Trace/"))
		pkgs, err := build.ImportDir(dir, 0)
		if err != nil {
			// cmd/trace is package main; ImportDir still works
			t.Fatalf("import %s: %v", pkgPath, err)
		}
		for _, imp := range pkgs.Imports {
			if imp == mcpImport || imp == mcpCmd || strings.HasPrefix(imp, mcpImport+"/") {
				t.Fatalf("%s must not import %s (G19)", pkgPath, imp)
			}
		}
		for _, imp := range pkgs.TestImports {
			if imp == mcpImport || imp == mcpCmd {
				t.Fatalf("%s test must not import %s (G19)", pkgPath, imp)
			}
		}
	}
}

func TestMCPVirginProjectDoesNotMkdir(t *testing.T) {
	ctx := context.Background()

	t.Run("virgin defaultRoot", func(t *testing.T) {
		dir := t.TempDir()
		srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
		_, _, err := callAdd(srv, ctx, tracemcp.AddInput{Kind: "goal", Title: "G"})
		mustNotInitialized(t, err)
		_, _, err = callVersion(srv, ctx, tracemcp.VersionInput{})
		mustNotInitialized(t, err)
		assertNoTraceDir(t, dir)
	})

	t.Run("initialized defaultRoot plus virgin project override", func(t *testing.T) {
		bound := t.TempDir()
		st, err := store.Open(bound)
		if err != nil {
			t.Fatal(err)
		}
		st.Close()
		virgin := t.TempDir()
		srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: bound})
		_, _, err = callAdd(srv, ctx, tracemcp.AddInput{Project: virgin, Kind: "goal", Title: "G"})
		mustNotInitialized(t, err)
		assertNoTraceDir(t, virgin)
	})

	t.Run("empty trace dir no db", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, ".trace"), 0o755); err != nil {
			t.Fatal(err)
		}
		srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
		_, _, err := callAdd(srv, ctx, tracemcp.AddInput{Kind: "goal", Title: "G"})
		mustNotInitialized(t, err)
		assertNoTraceDB(t, dir)
	})
}

func TestMCPInitializedOtherRootIsolated(t *testing.T) {
	ctx := context.Background()
	rootA := t.TempDir()
	stA, err := store.Open(rootA)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := domain.New(stA).DecideTool(ctx, domain.DecideToolInput{
		Slug: "mcp:trace_add", Decision: "DENIED", Reason: "planted deny A",
	}); err != nil {
		stA.Close()
		t.Fatalf("DecideTool A: %v", err)
	}
	stA.Close()

	rootB := t.TempDir()
	stB, err := store.Open(rootB)
	if err != nil {
		t.Fatal(err)
	}
	stB.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: rootA})
	res, _, err := callAdd(srv, ctx, tracemcp.AddInput{Project: rootB, Kind: "goal", Title: "from B"})
	if err != nil {
		t.Fatalf("add on initialized B should succeed: %v", err)
	}
	if mustOKID(t, mustText(t, res)) == "" {
		t.Fatal("expected goal id from B")
	}

	_, _, err = callAdd(srv, ctx, tracemcp.AddInput{Kind: "goal", Title: "from A"})
	if err == nil {
		t.Fatal("unbound add on A must remain DENIED")
	}
	if !strings.Contains(err.Error(), "DENIED") {
		t.Fatalf("error should mention DENIED, got %q", err.Error())
	}

	stB2, err := store.Open(rootB)
	if err != nil {
		t.Fatal(err)
	}
	defer stB2.Close()
	row, err := stB2.GetCapabilityToolDecisionBySlug("mcp:trace_add")
	if err != nil {
		t.Fatalf("B should have own allowlist row, not inherit A DENIED: %v", err)
	}
	if row.Decision == store.ToolDecisionDenied {
		t.Fatalf("B inherited A's DENIED: %+v", row)
	}
}

func TestToolNamesRegistered(t *testing.T) {
	srv := tracemcp.NewServer(tracemcp.Options{})
	if srv.MCP() == nil {
		t.Fatal("nil MCP server")
	}
	names := tracemcp.RegisteredToolNames()
	want := []string{
		"trace_why", "trace_context", "trace_add",
		"trace_link", "trace_transition", "trace_review",
		"trace_tasks", "trace_capability", "trace_impact", "trace_version",
		"trace_search", "trace_changes", "trace_regressions", "trace_loop",
		"trace_agents", "trace_plan", "trace_explore",
	}
	if len(names) != 17 {
		t.Fatalf("want 17 tools, got %d: %v", len(names), names)
	}
	for i, n := range want {
		if names[i] != n {
			t.Fatalf("tool[%d]=%q want %q", i, names[i], n)
		}
	}
}

func TestServerInstructionsNonEmpty(t *testing.T) {
	instr := tracemcp.ServerInstructions()
	if strings.TrimSpace(instr) == "" {
		t.Fatal("ServerInstructions must be non-empty")
	}
	if !strings.Contains(instr, "moat-first") {
		t.Fatalf("expected moat-first heading, got %q", instr)
	}
}

func TestServerInstructionsMoatLead(t *testing.T) {
	instr := tracemcp.ServerInstructions()
	moat := []string{"trace_tasks", "trace_context", "trace_loop", "trace_review", "trace_plan"}
	prev := -1
	for _, tool := range moat {
		idx := strings.Index(instr, tool)
		if idx < 0 {
			t.Fatalf("missing moat tool %q in instructions", tool)
		}
		if idx <= prev {
			t.Fatalf("moat tool %q out of order (idx=%d prev=%d)", tool, idx, prev)
		}
		prev = idx
	}
	if !strings.Contains(instr, "query") {
		t.Fatal("instructions should mention optional trace_context query")
	}
}

func TestServerInstructionsComposeRecipe(t *testing.T) {
	instr := tracemcp.ServerInstructions()
	if !strings.Contains(strings.ToLower(instr), "compose") {
		t.Fatal("instructions should mention compose-first recipe")
	}
	compose := []string{"trace_search", "trace_why", "trace_impact", "trace_capability"}
	prev := -1
	for _, tool := range compose {
		idx := strings.Index(instr, tool)
		if idx < 0 {
			t.Fatalf("missing compose tool %q in instructions", tool)
		}
		if idx <= prev {
			t.Fatalf("compose tool %q out of order (idx=%d prev=%d)", tool, idx, prev)
		}
		prev = idx
	}
}

func TestServerInstructionsExploreOptional(t *testing.T) {
	instr := tracemcp.ServerInstructions()
	if !strings.Contains(instr, "trace_explore") {
		t.Fatal("instructions should mention optional trace_explore")
	}
	moatIdx := strings.Index(instr, "trace_tasks")
	exploreIdx := strings.Index(instr, "trace_explore")
	if moatIdx < 0 || exploreIdx < 0 || exploreIdx <= moatIdx {
		t.Fatalf("trace_explore should appear after moat lead (moat=%d explore=%d)", moatIdx, exploreIdx)
	}
	if !strings.Contains(strings.ToLower(instr), "optional") {
		t.Fatal("explore section should mark optional convenience")
	}
}

func TestServerInstructionsStaleHygiene(t *testing.T) {
	instr := tracemcp.ServerInstructions()
	for _, want := range []string{"trace_version", "9/17", "stale", "reload"} {
		if !strings.Contains(instr, want) {
			t.Fatalf("instructions missing stale hygiene token %q", want)
		}
	}
}

func TestMCPExploreTaskRequired(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	_, _, err = callExplore(srv, context.Background(), tracemcp.ExploreInput{})
	if err == nil || !strings.Contains(err.Error(), "task_id is required") {
		t.Fatalf("expected task_id required, got %v", err)
	}
}

func TestMCPExploreQueryMerged(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "MCP explore merge"})
	if err != nil {
		t.Fatal(err)
	}
	const token = "mcpexplorequeryABC"
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "MCP explore unlinked " + token})
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	res, _, err := callExplore(srv, ctx, tracemcp.ExploreInput{TaskID: task.ID, Query: token})
	if err != nil {
		t.Fatalf("trace_explore: %v", err)
	}
	var out compiler.ExploreResult
	if err := json.Unmarshal([]byte(mustText(t, res)), &out); err != nil {
		t.Fatal(err)
	}
	if out.TaskSummary.TaskID != task.ID {
		t.Fatalf("task moat lost: %+v", out.TaskSummary)
	}
	found := false
	for _, h := range out.SearchHits {
		if h.EntityType == "decision" && h.EntityID == dec.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected query-merged search hit; hits=%+v", out.SearchHits)
	}
}

func TestTraceAddDescriptionMentionsPromotionPath(t *testing.T) {
	desc := tracemcp.TraceAddDescriptionForTest()
	if !strings.Contains(desc, "(discovery|task|goal|decision|assumption|plan-change|claim|evidence)") {
		t.Fatalf("trace_add description should list discovery first: %q", desc)
	}
	for _, want := range []string{
		"BLOCKING discovery",
		"trace_add kind=task",
		"spawned_tasks",
		"discovery_id",
		"Prefer the task/promotion path",
		"discovery-only",
		"before product edits",
	} {
		if !strings.Contains(desc, want) {
			t.Fatalf("trace_add description missing %q: %q", want, desc)
		}
	}
}

func TestTraceAddKindSchemaListsDiscoveryFirst(t *testing.T) {
	f, ok := reflect.TypeOf(tracemcp.AddInput{}).FieldByName("Kind")
	if !ok {
		t.Fatal("AddInput.Kind missing")
	}
	tag := f.Tag.Get("jsonschema")
	want := "discovery|task|goal|decision|assumption|plan-change|claim|evidence"
	if tag != want {
		t.Fatalf("AddInput.Kind jsonschema=%q want %q", tag, want)
	}
}

func TestMCPSearchRegistered(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "MCP search task zephyrmcpsearch"})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "mcp zephyrmcpchange token",
		Paths:  []domain.ChangePathInput{{Path: "x.go", Status: "M"}},
	})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	names := tracemcp.RegisteredToolNames()
	found := false
	for _, n := range names {
		if n == "trace_search" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("trace_search not in catalog: %v", names)
	}

	res, _, err := callSearch(srv, ctx, tracemcp.SearchInput{Query: "zephyrmcpchange"})
	if err != nil {
		t.Fatalf("trace_search: %v", err)
	}
	text := mustText(t, res)
	var resp struct {
		OK    bool            `json:"ok"`
		Hits  []retrieval.Hit `json:"hits"`
		Count int             `json:"count"`
	}
	if err := json.Unmarshal([]byte(text), &resp); err != nil {
		t.Fatalf("json: %v\n%s", err, text)
	}
	if !resp.OK || resp.Count < 1 {
		t.Fatalf("expected hits: %+v", resp)
	}
	foundChange := false
	for _, h := range resp.Hits {
		if h.EntityType == "change" && h.EntityID == chg.ID {
			foundChange = true
			break
		}
	}
	if !foundChange {
		t.Fatalf("expected change %s in hits: %+v", chg.ID, resp.Hits)
	}
}

func TestMCPChangesList(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "MCP changes list"})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "listed",
		Paths:  []domain.ChangePathInput{{Path: "a.go", Status: "M"}},
	}); err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	res, _, err := callChanges(srv, ctx, tracemcp.ChangesInput{Action: "list", TaskID: task.ID})
	if err != nil {
		t.Fatalf("trace_changes list: %v", err)
	}
	text := mustText(t, res)
	var rows []map[string]any
	if err := json.Unmarshal([]byte(text), &rows); err != nil {
		t.Fatalf("json: %v\n%s", err, text)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 row, got %v", rows)
	}
	if rows[0]["task_id"] != task.ID {
		t.Fatalf("row: %v", rows[0])
	}
}

func TestMCPRegressionsRegistered(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "reg task", GoalID: &g.ID})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	bl, err := svc.CreateBaseline(ctx, domain.BaselineInput{
		GitCommit: "abc1234", ScoresJSON: `{"latency_ms": 10}`,
	})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	outEval, err := svc.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID: task.ID, BaselineID: bl.ID, ScoresJSON: `{"latency_ms": 50}`,
	})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	reg, err := svc.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
		OutcomeID: outEval.ID, TaskID: task.ID,
	})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	found := false
	for _, n := range tracemcp.RegisteredToolNames() {
		if n == "trace_regressions" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("trace_regressions not registered: %v", tracemcp.RegisteredToolNames())
	}

	res, _, err := callRegressions(srv, ctx, tracemcp.RegressionsInput{Action: "list", TaskID: task.ID})
	if err != nil {
		t.Fatalf("trace_regressions list: %v", err)
	}
	text := mustText(t, res)
	var payload struct {
		OK          bool `json:"ok"`
		Count       int  `json:"count"`
		Regressions []struct {
			ID string `json:"id"`
		} `json:"regressions"`
	}
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		t.Fatalf("json: %v\n%s", err, text)
	}
	if !payload.OK || payload.Count < 1 {
		t.Fatalf("payload: %#v", payload)
	}
	got := false
	for _, r := range payload.Regressions {
		if r.ID == reg.ID {
			got = true
		}
	}
	if !got {
		t.Fatalf("regression %s not in MCP list: %+v", reg.ID, payload.Regressions)
	}
}

func seedMCPGoalTaskPlan(t *testing.T, st *store.Store) (goalID, taskID string) {
	t.Helper()
	ctx := context.Background()
	dsvc := domain.New(st)
	psvc := planner.New(st)
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "MCP loop goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "MCP loop seed", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	cp, err := psvc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{
			Title:  "Phase 1",
			Scopes: []planner.ScopeInput{{Title: "Scope 1"}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	scopeID := cp.Phases[0].Scopes[0].ID
	if err := psvc.SetCurrentScope(ctx, g.ID, scopeID); err != nil {
		t.Fatal(err)
	}
	if _, err := psvc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID:      scopeID,
		ExitCriteria: []string{"packet ready"},
		WorkItems:    []planner.WorkItem{{Title: "emit loop packet"}},
	}); err != nil {
		t.Fatal(err)
	}
	return g.ID, task.ID
}

func TestMCPLoopNext(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	_, taskID := seedMCPGoalTaskPlan(t, st)
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	ctx := context.Background()
	res, _, err := callLoop(srv, ctx, tracemcp.LoopInput{Action: "next", TaskID: taskID})
	if err != nil {
		t.Fatalf("trace_loop next: %v", err)
	}
	text := mustText(t, res)
	var pkt map[string]any
	if err := json.Unmarshal([]byte(text), &pkt); err != nil {
		t.Fatalf("json: %v\n%s", err, text)
	}
	if pkt["schema_version"] != loop.NextSchemaVersion {
		t.Fatalf("schema_version: got %v want %s", pkt["schema_version"], loop.NextSchemaVersion)
	}
	seed, _ := pkt["seed"].(map[string]any)
	if seed == nil || seed["task_id"] != taskID {
		t.Fatalf("seed: %#v", pkt["seed"])
	}
}

func TestMCPLoopApply(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	goalID, taskID := seedMCPGoalTaskPlan(t, st)
	st.Close()

	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       "11111111-1111-4111-8111-111111111112",
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes:        loop.ApplyWrites{},
	}
	envRaw, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	ctx := context.Background()
	res, _, err := callLoop(srv, ctx, tracemcp.LoopInput{Action: "apply", Envelope: envRaw})
	if err != nil {
		t.Fatalf("trace_loop apply: %v", err)
	}
	text := mustText(t, res)
	var out loop.ApplyResult
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		t.Fatalf("json: %v\n%s", err, text)
	}
	if out.Saturated {
		t.Fatalf("first empty apply must not saturate: %+v", out)
	}

	env2 := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       "11111111-1111-4111-8111-111111111113",
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes:        loop.ApplyWrites{},
	}
	envRaw2, err := json.Marshal(env2)
	if err != nil {
		t.Fatal(err)
	}
	res2, _, err := callLoop(srv, ctx, tracemcp.LoopInput{Action: "apply", Envelope: envRaw2})
	if err != nil {
		t.Fatalf("trace_loop apply 2: %v", err)
	}
	text2 := mustText(t, res2)
	var out2 loop.ApplyResult
	if err := json.Unmarshal([]byte(text2), &out2); err != nil {
		t.Fatalf("json: %v\n%s", err, text2)
	}
	if !out2.Saturated {
		t.Fatalf("second empty apply should saturate: %+v", out2)
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	evs, err := st2.ListEventsByEntity(domain.EntityTask, taskID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range evs {
		if e.Type == domain.EventDeliberationTransition {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("apply must persist loop-step (deliberation.transition event)")
	}
}

func TestMCPLoopStatus(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	goalID, taskID := seedMCPGoalTaskPlan(t, st)
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	ctx := context.Background()
	res, _, err := callLoop(srv, ctx, tracemcp.LoopInput{Action: "status", TaskID: taskID})
	if err != nil {
		t.Fatalf("trace_loop status: %v", err)
	}
	text := mustText(t, res)
	var status map[string]any
	if err := json.Unmarshal([]byte(text), &status); err != nil {
		t.Fatalf("json: %v\n%s", err, text)
	}
	if status["schema_version"] != loop.StatusSchemaVersion {
		t.Fatalf("schema_version: got %v want %s", status["schema_version"], loop.StatusSchemaVersion)
	}
	seed, _ := status["seed"].(map[string]any)
	if seed == nil || seed["task_id"] != taskID || seed["goal_id"] != goalID {
		t.Fatalf("seed: %#v", status["seed"])
	}
}

func TestMCPLoopGate_MatchesCLI(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	dsvc := domain.New(st)
	psvc := planner.New(st)
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "gate goal"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "gate task", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	allowed, violations, err := loop.EvaluateGate(ctx, dsvc, psvc, st, task.ID, loop.GateForEdit)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	expected := map[string]any{
		"schema_version": "trace.loop.gate.v1",
		"task_id":        task.ID,
		"for":            "edit",
		"allowed":        allowed,
		"violations":     violations,
	}
	if !allowed && len(violations) == 1 {
		expected["recommended_phase"] = violations[0].RecommendedPhase
		expected["reason_code"] = violations[0].ReasonCode
	}
	expRaw, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	res, _, err := callLoop(srv, ctx, tracemcp.LoopInput{Action: "gate", TaskID: task.ID, For: "edit"})
	if err != nil {
		t.Fatalf("trace_loop gate: %v", err)
	}
	gotRaw := []byte(mustText(t, res))
	var got map[string]any
	if err := json.Unmarshal(gotRaw, &got); err != nil {
		t.Fatalf("mcp gate json: %v\n%s", err, gotRaw)
	}
	var exp map[string]any
	if err := json.Unmarshal(expRaw, &exp); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"schema_version", "task_id", "for", "allowed", "reason_code", "recommended_phase"} {
		if exp[key] != got[key] {
			t.Fatalf("gate field %q mismatch exp=%#v got=%#v", key, exp[key], got[key])
		}
	}
	expViol, _ := json.Marshal(exp["violations"])
	gotViol, _ := json.Marshal(got["violations"])
	if string(expViol) != string(gotViol) {
		t.Fatalf("violations mismatch exp=%s got=%s", expViol, gotViol)
	}
	if allowed {
		t.Fatal("expected blocked edit without plan")
	}
}

func TestMCPAgentsList(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	if err := installAgentCatalog(t, dir); err != nil {
		t.Fatal(err)
	}

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	ctx := context.Background()
	res, _, err := callAgents(srv, ctx, tracemcp.AgentsInput{Action: "list"})
	if err != nil {
		t.Fatalf("trace_agents list: %v", err)
	}
	text := mustText(t, res)
	var items []map[string]any
	if err := json.Unmarshal([]byte(text), &items); err != nil {
		t.Fatalf("json: %v\n%s", err, text)
	}
	if len(items) != 6 {
		t.Fatalf("want 6 agents, got %d: %s", len(items), text)
	}
}

func TestMCPAgentsRecommend(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	if err := installAgentCatalog(t, dir); err != nil {
		t.Fatal(err)
	}

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	ctx := context.Background()
	res, _, err := callAgents(srv, ctx, tracemcp.AgentsInput{Action: "recommend", Phase: "CRITIQUE"})
	if err != nil {
		t.Fatalf("trace_agents recommend: %v", err)
	}
	text := mustText(t, res)
	var recs []map[string]any
	if err := json.Unmarshal([]byte(text), &recs); err != nil {
		t.Fatalf("json: %v\n%s", err, text)
	}
	if len(recs) == 0 || recs[0]["agent_slug"] != "agent:code-reviewer" {
		t.Fatalf("recommendations: %s", text)
	}
}

func installAgentCatalog(t *testing.T, dir string) error {
	t.Helper()
	return install.InstallAgentDefaults(install.InstallOpts{
		Write:       true,
		ProjectRoot: dir,
		ErrOut:      io.Discard,
	})
}

// TestMCPUnprefixedDecideGatesCallTool: unprefixed DecideTool DENIED gates CallTool Name (DF-78).
func TestMCPUnprefixedDecideGatesCallTool(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if _, err := domain.New(st).DecideTool(ctx, domain.DecideToolInput{
		Slug: "trace_why", Decision: "DENIED", Reason: "unprefixed deny",
	}); err != nil {
		st.Close()
		t.Fatalf("DecideTool: %v", err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	_, _, err = callWhy(srv, ctx, tracemcp.WhyInput{
		EntityType: domain.EntityTask, ID: "00000000-0000-0000-0000-000000000001",
	})
	if err == nil {
		t.Fatal("CallTool must fail when unprefixed trace_why was DENIED")
	}
	if !strings.Contains(err.Error(), "DENIED") {
		t.Fatalf("error should mention DENIED, got %q", err.Error())
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	row, err := st2.GetCapabilityToolDecisionBySlug("mcp:trace_why")
	if err != nil {
		t.Fatalf("canonical mcp:trace_why row: %v", err)
	}
	if row.Decision == store.ToolDecisionAutoAllowed {
		t.Fatal("must not AUTO_ALLOW mcp:trace_why after unprefixed DENIED")
	}
	if row.Decision != store.ToolDecisionDenied {
		t.Fatalf("mcp:trace_why decision: got %q want DENIED", row.Decision)
	}
}

// TestCLIAddDeniedDoesNotBlockMCPAdd: cli:add DENIED must not gate CallTool trace_add.
func TestCLIAddDeniedDoesNotBlockMCPAdd(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if _, err := domain.New(st).DecideTool(ctx, domain.DecideToolInput{
		Slug: "cli:add", Decision: "DENIED", Reason: "shell lockdown",
	}); err != nil {
		st.Close()
		t.Fatalf("DecideTool cli:add: %v", err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	res, _, err := callAdd(srv, ctx, tracemcp.AddInput{Kind: "goal", Title: "from MCP"})
	if err != nil {
		t.Fatalf("CallTool trace_add must succeed when only cli:add is DENIED: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(mustText(t, res)), &m); err != nil {
		t.Fatal(err)
	}
	if m["ok"] != true {
		t.Fatalf("add payload: %v", m)
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	cliRow, err := st2.GetCapabilityToolDecisionBySlug("cli:add")
	if err != nil {
		t.Fatalf("cli:add row: %v", err)
	}
	if cliRow.Decision != store.ToolDecisionDenied {
		t.Fatalf("cli:add must stay DENIED, got %q", cliRow.Decision)
	}
	mcpRow, err := st2.GetCapabilityToolDecisionBySlug("mcp:trace_add")
	if err != nil {
		t.Fatalf("expected durable AUTO_ALLOWED mcp:trace_add: %v", err)
	}
	if mcpRow.Decision != store.ToolDecisionAutoAllowed {
		t.Fatalf("mcp:trace_add want AUTO_ALLOWED got %q", mcpRow.Decision)
	}
}

// TestMCPAssertDeniedBlocksCallTool: DecideTool DENIED on a builtin slug blocks CallTool.
func TestMCPAssertDeniedBlocksCallTool(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	slug := "mcp:trace_why"
	if _, err := domain.New(st).DecideTool(ctx, domain.DecideToolInput{
		Slug: slug, Decision: "DENIED", Reason: "planted deny",
	}); err != nil {
		st.Close()
		t.Fatalf("DecideTool: %v", err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	_, _, err = callWhy(srv, ctx, tracemcp.WhyInput{
		EntityType: domain.EntityTask, ID: "00000000-0000-0000-0000-000000000001",
	})
	if err == nil {
		t.Fatal("CallTool must fail when mcp:trace_why is DENIED")
	}
	if !strings.Contains(err.Error(), "DENIED") {
		t.Fatalf("error should mention DENIED, got %q", err.Error())
	}
}

// TestMCPAssertBuiltinAutoAllowedSucceeds: fresh project DB; builtin CallTool succeeds via AUTO_ALLOWED.
func TestMCPAssertBuiltinAutoAllowedSucceeds(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	ctx := context.Background()
	res, _, err := callVersion(srv, ctx, tracemcp.VersionInput{})
	if err != nil {
		t.Fatalf("builtin CallTool should AUTO_ALLOW: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(mustText(t, res)), &m); err != nil {
		t.Fatal(err)
	}
	if m["ok"] != true || m["name"] != "trace" {
		t.Fatalf("version payload: %v", m)
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	row, err := st2.GetCapabilityToolDecisionBySlug("mcp:trace_version")
	if err != nil {
		t.Fatalf("expected durable AUTO_ALLOWED row: %v", err)
	}
	if row.Decision != store.ToolDecisionAutoAllowed {
		t.Fatalf("want AUTO_ALLOWED got %q", row.Decision)
	}
}

func TestTraceTasksParity(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "T", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	res, _, err := callTasks(srv, ctx, tracemcp.TasksInput{})
	if err != nil {
		t.Fatalf("trace_tasks: %v", err)
	}
	text := mustText(t, res)
	var rows []map[string]any
	if err := json.Unmarshal([]byte(text), &rows); err != nil {
		t.Fatalf("tasks json: %v\n%s", err, text)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 task, got %s", text)
	}
	for _, k := range []string{"id", "title", "work_state", "goal_id"} {
		if _, ok := rows[0][k]; !ok {
			t.Fatalf("missing key %q in %v", k, rows[0])
		}
	}
	if rows[0]["id"] != task.ID || rows[0]["title"] != "T" {
		t.Fatalf("row mismatch: %v", rows[0])
	}

	filtered, _, err := callTasks(srv, ctx, tracemcp.TasksInput{GoalID: g.ID})
	if err != nil {
		t.Fatal(err)
	}
	var frows []map[string]any
	if err := json.Unmarshal([]byte(mustText(t, filtered)), &frows); err != nil {
		t.Fatal(err)
	}
	if len(frows) != 1 {
		t.Fatalf("goal filter: %v", frows)
	}
	empty, _, err := callTasks(srv, ctx, tracemcp.TasksInput{GoalID: "00000000-0000-0000-0000-000000000000"})
	if err != nil {
		t.Fatal(err)
	}
	var erows []map[string]any
	if err := json.Unmarshal([]byte(mustText(t, empty)), &erows); err != nil {
		t.Fatal(err)
	}
	if len(erows) != 0 {
		t.Fatalf("want empty filter result, got %v", erows)
	}
}

func TestTraceCapabilityActions(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Cap task"})
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})

	decl, _, err := callCapability(srv, ctx, tracemcp.CapabilityInput{
		Action: "declare", Kind: "MCP", Slug: "mcp:demo", Title: "Demo", Status: "UNAVAILABLE",
	})
	if err != nil {
		t.Fatalf("declare: %v", err)
	}
	capID := mustOKID(t, mustText(t, decl))

	listRes, _, err := callCapability(srv, ctx, tracemcp.CapabilityInput{Action: "list"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var listEnv map[string]any
	if err := json.Unmarshal([]byte(mustText(t, listRes)), &listEnv); err != nil {
		t.Fatal(err)
	}
	caps, ok := listEnv["capabilities"].([]any)
	if !ok || len(caps) != 1 {
		t.Fatalf("list: %v", listEnv)
	}
	row := caps[0].(map[string]any)
	for _, k := range []string{"id", "kind", "slug", "title", "status"} {
		if _, ok := row[k]; !ok {
			t.Fatalf("list row missing %q: %v", k, row)
		}
	}
	if _, ok := row["ID"]; ok {
		t.Fatalf("PascalCase leaked: %v", row)
	}

	_, _, err = callCapability(srv, ctx, tracemcp.CapabilityInput{
		Action: "require", TaskID: task.ID, Capability: "mcp:demo",
	})
	if err != nil {
		t.Fatalf("require: %v", err)
	}

	missRes, _, err := callCapability(srv, ctx, tracemcp.CapabilityInput{
		Action: "missing", TaskID: task.ID,
	})
	if err != nil {
		t.Fatalf("missing: %v", err)
	}
	var missEnv map[string]any
	if err := json.Unmarshal([]byte(mustText(t, missRes)), &missEnv); err != nil {
		t.Fatal(err)
	}
	missing, ok := missEnv["missing"].([]any)
	if !ok || len(missing) != 1 {
		t.Fatalf("missing: %v", missEnv)
	}
	mrow := missing[0].(map[string]any)
	if mrow["id"] != capID || mrow["slug"] != "mcp:demo" {
		t.Fatalf("missing row: %v", mrow)
	}

	_, _, err = callCapability(srv, ctx, tracemcp.CapabilityInput{
		Action: "unrequire", TaskID: task.ID, Capability: capID,
	})
	if err != nil {
		t.Fatalf("unrequire: %v", err)
	}
	missRes, _, err = callCapability(srv, ctx, tracemcp.CapabilityInput{
		Action: "missing", TaskID: task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(mustText(t, missRes)), &missEnv); err != nil {
		t.Fatal(err)
	}
	if n, _ := missEnv["count"].(float64); n != 0 {
		t.Fatalf("after unrequire count=%v", missEnv["count"])
	}
}

func TestTraceVersion(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	res, _, err := callVersion(srv, context.Background(), tracemcp.VersionInput{})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(mustText(t, res)), &m); err != nil {
		t.Fatal(err)
	}
	if m["ok"] != true || m["name"] != "trace" || m["version"] != "0.0.0-dev" {
		t.Fatalf("version payload: %v", m)
	}
}

func TestRegisteredToolNames_IncludesTracePlan(t *testing.T) {
	names := tracemcp.RegisteredToolNames()
	if len(names) != 17 {
		t.Fatalf("want 17 tools, got %d", len(names))
	}
	if names[len(names)-2] != "trace_plan" {
		t.Fatalf("penultimate tool=%q want trace_plan", names[len(names)-2])
	}
	if names[len(names)-1] != "trace_explore" {
		t.Fatalf("last tool=%q want trace_explore", names[len(names)-1])
	}
}

func TestGreenfield_MCPPlanBootstrap_EditGatePasses(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	ctx := context.Background()

	addRes, _, err := callAdd(srv, ctx, tracemcp.AddInput{Kind: "goal", Title: "Greenfield goal"})
	if err != nil {
		t.Fatal(err)
	}
	goalID := mustOKID(t, mustText(t, addRes))
	addRes, _, err = callAdd(srv, ctx, tracemcp.AddInput{Kind: "task", Title: "Greenfield task", GoalID: goalID})
	if err != nil {
		t.Fatal(err)
	}
	taskID := mustOKID(t, mustText(t, addRes))

	coarseRes, _, err := callPlan(srv, ctx, tracemcp.PlanInput{
		Action: "create-coarse", GoalID: goalID, Phase: "Phase 1", Scope: "Scope 1",
	})
	if err != nil {
		t.Fatalf("create-coarse: %v", err)
	}
	var coarse map[string]any
	if err := json.Unmarshal([]byte(mustText(t, coarseRes)), &coarse); err != nil {
		t.Fatal(err)
	}
	phases, _ := coarse["phases"].([]any)
	if len(phases) == 0 {
		t.Fatalf("no phases: %v", coarse)
	}
	phase0 := phases[0].(map[string]any)
	scopes, _ := phase0["scopes"].([]any)
	scopeID, _ := scopes[0].(map[string]any)["id"].(string)
	if scopeID == "" {
		t.Fatal("missing scope id")
	}

	if _, _, err := callPlan(srv, ctx, tracemcp.PlanInput{
		Action: "set-current", GoalID: goalID, ScopeID: scopeID,
	}); err != nil {
		t.Fatalf("set-current: %v", err)
	}
	if _, _, err := callPlan(srv, ctx, tracemcp.PlanInput{
		Action: "deep", ScopeID: scopeID, Exit: []string{"packet ready"}, Work: []string{"work"},
	}); err != nil {
		t.Fatalf("deep: %v", err)
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	if _, err := st2.UpsertDeliberationState(store.DeliberationState{
		TaskID: taskID, GoalID: goalID, CurrentPhase: "CRITIQUE", PlanCritiqued: true,
	}); err != nil {
		t.Fatal(err)
	}
	dom := domain.New(st2)
	psvc := planner.New(st2)
	allowed, violations, err := loop.EvaluateGate(ctx, dom, psvc, st2, taskID, loop.GateForEdit)
	if err != nil {
		t.Fatal(err)
	}
	if !allowed {
		t.Fatalf("edit gate blocked: %+v", violations)
	}
	for _, v := range violations {
		if v.ReasonCode == "plan_missing" {
			t.Fatalf("plan_missing after bootstrap: %+v", violations)
		}
	}
	status, err := loop.Status(ctx, st2, psvc, loop.ApplySeed{TaskID: taskID, GoalID: goalID})
	if err != nil {
		t.Fatal(err)
	}
	if status.Deliberation == nil || !status.Deliberation.PolicyInputs.PlanExists {
		t.Fatalf("status policy_inputs.plan_exists false: %+v", status.Deliberation)
	}
}

func TestImportBoundaryMCPNoPlanImpactIndexTools(t *testing.T) {
	hasImpact := false
	hasPlan := false
	for _, name := range tracemcp.RegisteredToolNames() {
		switch name {
		case "trace_plan":
			hasPlan = true
		case "trace_index":
			t.Fatalf("forbidden tool registered: %s", name)
		case "trace_install", "trace_decide":
			t.Fatalf("forbidden tool registered: %s", name)
		case "trace_impact":
			hasImpact = true
		}
		if strings.Contains(name, "install") || strings.Contains(name, "decide") {
			t.Fatalf("forbidden install/decide MCP name: %s", name)
		}
		if name == "trace_index" {
			t.Fatalf("forbidden tool registered: %s", name)
		}
	}
	if !hasImpact {
		t.Fatal("trace_impact must be registered")
	}
	if !hasPlan {
		t.Fatal("trace_plan must be registered")
	}
}

// TestMCPTraceImpactReport locks DF-72: action=report returns snake_case overall_class.
func TestMCPTraceImpactReport(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "break API"})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if _, err := svc.AddImpactFinding(ctx, dec.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassDestructive,
		Kind:        domain.FindingKindWorkAtRisk,
		Body:        "callers break",
	}); err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	res, _, err := callImpact(srv, ctx, tracemcp.ImpactInput{
		Action: "report", Decision: dec.ID,
	})
	if err != nil {
		t.Fatalf("trace_impact report: %v", err)
	}
	text := mustText(t, res)
	var m map[string]any
	if err := json.Unmarshal([]byte(text), &m); err != nil {
		t.Fatalf("json: %v\n%s", err, text)
	}
	if m["ok"] != true {
		t.Fatalf("ok: %v", m)
	}
	if m["overall_class"] != "DESTRUCTIVE" {
		t.Fatalf("overall_class: %v", m["overall_class"])
	}
	findings, ok := m["findings"].([]any)
	if !ok || len(findings) != 1 {
		t.Fatalf("findings: %v", m["findings"])
	}
	f := findings[0].(map[string]any)
	if f["impact_class"] != "DESTRUCTIVE" {
		t.Fatalf("findings[].impact_class: %v", f)
	}
	if strings.Contains(text, `"ImpactClass"`) {
		t.Fatalf("PascalCase ImpactClass: %s", text)
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	row, err := st2.GetCapabilityToolDecisionBySlug("mcp:trace_impact")
	if err != nil {
		t.Fatalf("expected AUTO_ALLOWED row after CallTool: %v", err)
	}
	if row.Decision != store.ToolDecisionAutoAllowed {
		t.Fatalf("want AUTO_ALLOWED got %q", row.Decision)
	}
}

// TestMCPImpactDeniedBlocksCallTool locks DF-72: DENIED mcp:trace_impact fails closed.
func TestMCPImpactDeniedBlocksCallTool(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	svc := domain.New(st)
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "gated"})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if _, err := svc.DecideTool(ctx, domain.DecideToolInput{
		Slug: "mcp:trace_impact", Decision: "DENIED", Reason: "planted deny",
	}); err != nil {
		st.Close()
		t.Fatalf("DecideTool: %v", err)
	}
	st.Close()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: dir})
	_, _, err = callImpact(srv, ctx, tracemcp.ImpactInput{
		Action: "report", Decision: dec.ID,
	})
	if err == nil {
		t.Fatal("CallTool must fail when mcp:trace_impact is DENIED")
	}
	if !strings.Contains(err.Error(), "DENIED") {
		t.Fatalf("error should mention DENIED, got %q", err.Error())
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	findings, err := domain.New(st2).ListImpactFindings(ctx, dec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("DENIED report must not write findings, got %+v", findings)
	}
}

// --- helpers that invoke unexported tool methods via exported input types ---
// Handlers are package-internal; tests live in mcp_test and call through a thin
// exported test seam in the mcp package (see export_test.go).

func callWhy(s *tracemcp.Server, ctx context.Context, in tracemcp.WhyInput) (any, any, error) {
	return tracemcp.CallWhy(s, ctx, in)
}
func callContext(s *tracemcp.Server, ctx context.Context, in tracemcp.ContextInput) (any, any, error) {
	return tracemcp.CallContext(s, ctx, in)
}
func callAdd(s *tracemcp.Server, ctx context.Context, in tracemcp.AddInput) (any, any, error) {
	return tracemcp.CallAdd(s, ctx, in)
}
func callLink(s *tracemcp.Server, ctx context.Context, in tracemcp.LinkInput) (any, any, error) {
	return tracemcp.CallLink(s, ctx, in)
}
func callTransition(s *tracemcp.Server, ctx context.Context, in tracemcp.TransitionInput) (any, any, error) {
	return tracemcp.CallTransition(s, ctx, in)
}
func callReview(s *tracemcp.Server, ctx context.Context, in tracemcp.ReviewInput) (any, any, error) {
	return tracemcp.CallReview(s, ctx, in)
}
func callTasks(s *tracemcp.Server, ctx context.Context, in tracemcp.TasksInput) (any, any, error) {
	return tracemcp.CallTasks(s, ctx, in)
}
func callCapability(s *tracemcp.Server, ctx context.Context, in tracemcp.CapabilityInput) (any, any, error) {
	return tracemcp.CallCapability(s, ctx, in)
}
func callVersion(s *tracemcp.Server, ctx context.Context, in tracemcp.VersionInput) (any, any, error) {
	return tracemcp.CallVersion(s, ctx, in)
}
func callImpact(s *tracemcp.Server, ctx context.Context, in tracemcp.ImpactInput) (any, any, error) {
	return tracemcp.CallImpact(s, ctx, in)
}
func callSearch(s *tracemcp.Server, ctx context.Context, in tracemcp.SearchInput) (any, any, error) {
	return tracemcp.CallSearch(s, ctx, in)
}
func callChanges(s *tracemcp.Server, ctx context.Context, in tracemcp.ChangesInput) (any, any, error) {
	return tracemcp.CallChanges(s, ctx, in)
}
func callRegressions(s *tracemcp.Server, ctx context.Context, in tracemcp.RegressionsInput) (any, any, error) {
	return tracemcp.CallRegressions(s, ctx, in)
}
func callLoop(s *tracemcp.Server, ctx context.Context, in tracemcp.LoopInput) (any, any, error) {
	return tracemcp.CallLoop(s, ctx, in)
}
func callAgents(s *tracemcp.Server, ctx context.Context, in tracemcp.AgentsInput) (any, any, error) {
	return tracemcp.CallAgents(s, ctx, in)
}
func callPlan(s *tracemcp.Server, ctx context.Context, in tracemcp.PlanInput) (any, any, error) {
	return tracemcp.CallPlan(s, ctx, in)
}
func callExplore(s *tracemcp.Server, ctx context.Context, in tracemcp.ExploreInput) (any, any, error) {
	return tracemcp.CallExplore(s, ctx, in)
}

func mustNotInitialized(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("want ErrNotInitialized, got nil")
	}
	if !errors.Is(err, store.ErrNotInitialized) && !strings.Contains(err.Error(), "not initialized") {
		t.Fatalf("want ErrNotInitialized (or wrap), got %v", err)
	}
}

func assertNoTraceDir(t *testing.T, root string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(root, ".trace")); !os.IsNotExist(err) {
		t.Fatalf("did not want .trace/ under %s: %v", root, err)
	}
}

func assertNoTraceDB(t *testing.T, root string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(root, ".trace", "trace.db")); !os.IsNotExist(err) {
		t.Fatalf("did not want trace.db under %s: %v", root, err)
	}
}

func mustText(t *testing.T, res any) string {
	t.Helper()
	return tracemcp.ResultText(t, res)
}

func mustOKID(t *testing.T, text string) string {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(text), &m); err != nil {
		t.Fatal(err)
	}
	if m["ok"] != true {
		t.Fatalf("expected ok: %s", text)
	}
	id, _ := m["id"].(string)
	if id == "" {
		t.Fatalf("missing id: %s", text)
	}
	return id
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}
