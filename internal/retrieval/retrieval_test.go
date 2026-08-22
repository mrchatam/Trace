package retrieval_test

import (
	"context"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

func openEngine(t *testing.T) (*retrieval.Engine, *store.Store, *domain.Service) {
	t.Helper()
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return retrieval.New(st), st, domain.New(st)
}

func TestExactSearchExpandWhy(t *testing.T) {
	eng, st, svc := openEngine(t)
	ctx := context.Background()

	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Objective alpha"})
	if err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Task beta", Body: "exit: tests pass"})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if err := svc.LinkGoalTask(ctx, goal.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkGoalTask: %v", err)
	}
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Use FTS5 lexicolight"})
	if err != nil {
		t.Fatalf("CreateDecision: %v", err)
	}
	if err := svc.LinkDecisionTask(ctx, dec.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDecisionTask: %v", err)
	}
	disc, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "Found gap zephyrwhy"})
	if err != nil {
		t.Fatalf("CreateDiscovery: %v", err)
	}
	pc, err := svc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "Adjust plan"})
	if err != nil {
		t.Fatalf("CreatePlanChange: %v", err)
	}
	if err := svc.LinkDiscoveryPlanChange(ctx, disc.ID, pc.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDiscoveryPlanChange: %v", err)
	}

	_, err = st.UpsertFile("src/main.go", "hash1", nil)
	if err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if err := st.ReplaceFileSymbols("src/main.go", []store.Symbol{
		{Name: "RunMain", Kind: "function", StartLine: 1, EndLine: 10},
	}); err != nil {
		t.Fatalf("symbols: %v", err)
	}
	_, err = st.UpsertFile("src/util.go", "hash2", nil)
	if err != nil {
		t.Fatalf("UpsertFile util: %v", err)
	}
	if err := st.ReplaceFileImports("src/main.go", []store.Import{
		{ImportedPath: "src/util.go"},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	hits, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "task", EntityID: task.ID})
	if err != nil || len(hits) != 1 || hits[0].ReasonCode != retrieval.ReasonExactID {
		t.Fatalf("Exact task: %v %+v", err, hits)
	}
	miss, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "task", EntityID: "00000000-0000-0000-0000-000000000000"})
	if err != nil || len(miss) != 0 {
		t.Fatalf("Exact miss should be empty: %v %+v", err, miss)
	}

	pathHits, err := eng.Exact(ctx, retrieval.ExactQuery{Path: "src/main.go", SymbolName: "RunMain"})
	if err != nil {
		t.Fatalf("Exact path+symbol: %v", err)
	}
	var sawPath, sawSym bool
	for _, h := range pathHits {
		if h.ReasonCode == retrieval.ReasonExactPath {
			sawPath = true
		}
		if h.ReasonCode == retrieval.ReasonExactSymbol {
			sawSym = true
		}
	}
	if !sawPath || !sawSym {
		t.Fatalf("path/symbol exact: %+v", pathHits)
	}

	fts, err := eng.Search(ctx, "lexicolight", retrieval.SearchOptions{})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	foundDec := false
	for _, h := range fts {
		if h.ReasonCode != retrieval.ReasonFTSMatch {
			t.Fatalf("fts reason: %q", h.ReasonCode)
		}
		if h.EntityID == dec.ID {
			foundDec = true
		}
	}
	if !foundDec {
		t.Fatalf("expected decision in FTS: %+v", fts)
	}

	seed := hits[0]
	exp1, err := eng.Expand(ctx, []retrieval.Hit{seed}, 1)
	if err != nil {
		t.Fatalf("Expand1: %v", err)
	}
	var sawGoal, sawDec bool
	for _, h := range exp1 {
		if h.ReasonCode == "" {
			t.Fatal("missing reason_code on expand hit")
		}
		if h.EntityID == goal.ID && h.ReasonCode == retrieval.ReasonGoalHasTask {
			sawGoal = true
		}
		if h.EntityID == dec.ID && h.ReasonCode == retrieval.ReasonDecisionAffectsTask {
			sawDec = true
		}
	}
	if !sawGoal || !sawDec {
		t.Fatalf("expand1 goal/decision: sawGoal=%v sawDec=%v hits=%+v", sawGoal, sawDec, exp1)
	}

	fileHits, err := eng.Exact(ctx, retrieval.ExactQuery{Path: "src/main.go"})
	if err != nil || len(fileHits) == 0 {
		t.Fatalf("file exact: %v %+v", err, fileHits)
	}
	sexp, err := eng.Expand(ctx, fileHits, 1)
	if err != nil {
		t.Fatalf("struct expand: %v", err)
	}
	var sawSymbol, sawImport bool
	for _, h := range sexp {
		if h.EntityType == "symbol" && h.Title == "RunMain" {
			sawSymbol = true
		}
		if h.EntityType == "file" && h.Path == "src/util.go" {
			sawImport = true
		}
	}
	if !sawSymbol || !sawImport {
		t.Fatalf("structural: symbol=%v import=%v %+v", sawSymbol, sawImport, sexp)
	}

	discHit, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "discovery", EntityID: disc.ID})
	if err != nil || len(discHit) == 0 {
		t.Fatalf("disc exact: %v", err)
	}
	exp2, err := eng.Expand(ctx, discHit, 2)
	if err != nil {
		t.Fatalf("Expand2: %v", err)
	}
	sawPC := false
	for _, h := range exp2 {
		if h.EntityID == pc.ID {
			sawPC = true
		}
	}
	if !sawPC {
		t.Fatalf("depth2 should include plan_change: %+v", exp2)
	}
	if _, err := eng.Expand(ctx, discHit, 3); err == nil {
		t.Fatal("expected depth>2 rejection")
	}
	if _, err := eng.Expand(ctx, discHit, 0); err == nil {
		t.Fatal("expected depth<1 rejection")
	}

	why, err := eng.Why(ctx, "task", task.ID)
	if err != nil {
		t.Fatalf("Why: %v", err)
	}
	if len(why.Steps) == 0 {
		t.Fatal("why empty")
	}
	for _, step := range why.Steps {
		if step.ReasonCode == "" {
			t.Fatalf("why step missing reason: %+v", step)
		}
	}
	joined := ""
	for _, s := range why.Steps {
		joined += s.ReasonCode + " "
	}
	if !strings.Contains(joined, retrieval.ReasonExactID) {
		t.Fatalf("why missing exact_id: %s", joined)
	}
}

// TestWhyTaskDPCGoalScoped covers DF-19 single-goal / x0-shaped: Why + Expand
// from a task still surfaces unattributed discovery↔plan_change when the store
// has exactly one goal (supersedes global GC-01 attach).
func TestWhyTaskDPCGoalScoped(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Ship greeter + math demo"})
	if err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Wire greeting to arithmetic helpers"})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if err := svc.LinkGoalTask(ctx, goal.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkGoalTask: %v", err)
	}
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Prefer TypeScript greeter surface"})
	if err != nil {
		t.Fatalf("CreateDecision: %v", err)
	}
	if err := svc.LinkDecisionTask(ctx, dec.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDecisionTask: %v", err)
	}
	disc, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "math_util lacks a clamp helper"})
	if err != nil {
		t.Fatalf("CreateDiscovery: %v", err)
	}
	pc, err := svc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "Add clamp helper to math_util"})
	if err != nil {
		t.Fatalf("CreatePlanChange: %v", err)
	}
	if err := svc.LinkDiscoveryPlanChange(ctx, disc.ID, pc.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDiscoveryPlanChange: %v", err)
	}

	why, err := eng.Why(ctx, "task", task.ID)
	if err != nil {
		t.Fatalf("Why: %v", err)
	}
	var sawDisc, sawPC bool
	for _, s := range why.Steps {
		if s.EntityID == disc.ID && s.ReasonCode == retrieval.ReasonDiscoveryCausesPlanChg {
			sawDisc = true
		}
		if s.EntityID == pc.ID && s.ReasonCode == retrieval.ReasonDiscoveryCausesPlanChg {
			sawPC = true
		}
	}
	if !sawDisc || !sawPC {
		t.Fatalf("Why(task) missing discovery/plan_change with reason %q: sawDisc=%v sawPC=%v steps=%+v",
			retrieval.ReasonDiscoveryCausesPlanChg, sawDisc, sawPC, why.Steps)
	}

	seed, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "task", EntityID: task.ID})
	if err != nil || len(seed) == 0 {
		t.Fatalf("Exact task: %v %+v", err, seed)
	}
	exp1, err := eng.Expand(ctx, seed, 1)
	if err != nil {
		t.Fatalf("Expand1: %v", err)
	}
	sawDisc, sawPC = false, false
	for _, h := range exp1 {
		if h.EntityID == disc.ID && h.ReasonCode == retrieval.ReasonDiscoveryCausesPlanChg {
			sawDisc = true
		}
		if h.EntityID == pc.ID && h.ReasonCode == retrieval.ReasonDiscoveryCausesPlanChg {
			sawPC = true
		}
	}
	if !sawDisc || !sawPC {
		t.Fatalf("Expand(task,1) missing discovery/plan_change: sawDisc=%v sawPC=%v hits=%+v", sawDisc, sawPC, exp1)
	}
}

// TestWhyTaskDPCMultiGoalNoForeignPollution covers DF-19: DPC attributed to
// goal B must not appear on Why for a task under goal A.
func TestWhyTaskDPCMultiGoalNoForeignPollution(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	goalA, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Goal A"})
	if err != nil {
		t.Fatalf("goalA: %v", err)
	}
	taskA, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Task A"})
	if err != nil {
		t.Fatalf("taskA: %v", err)
	}
	if err := svc.LinkGoalTask(ctx, goalA.ID, taskA.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("link A: %v", err)
	}

	goalB, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Goal B"})
	if err != nil {
		t.Fatalf("goalB: %v", err)
	}
	taskB, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Task B"})
	if err != nil {
		t.Fatalf("taskB: %v", err)
	}
	if err := svc.LinkGoalTask(ctx, goalB.ID, taskB.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("link B: %v", err)
	}

	discB, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "Foreign discovery"})
	if err != nil {
		t.Fatalf("discB: %v", err)
	}
	pcB, err := svc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "Foreign plan change"})
	if err != nil {
		t.Fatalf("pcB: %v", err)
	}
	if err := svc.LinkDiscoveryPlanChange(ctx, discB.ID, pcB.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDiscoveryPlanChange: %v", err)
	}
	// Attribute DPC to goal B via domain discovery_mentions_task (DF-42).
	if err := svc.LinkDiscoveryMentionsTask(ctx, discB.ID, taskB.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDiscoveryMentionsTask: %v", err)
	}

	// Unattributed DPC under multi-goal must also be omitted from task A.
	discU, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "Unattributed discovery"})
	if err != nil {
		t.Fatalf("discU: %v", err)
	}
	pcU, err := svc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "Unattributed plan change"})
	if err != nil {
		t.Fatalf("pcU: %v", err)
	}
	if err := svc.LinkDiscoveryPlanChange(ctx, discU.ID, pcU.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDiscoveryPlanChange U: %v", err)
	}

	whyA, err := eng.Why(ctx, "task", taskA.ID)
	if err != nil {
		t.Fatalf("Why A: %v", err)
	}
	for _, s := range whyA.Steps {
		if s.EntityID == discB.ID || s.EntityID == pcB.ID || s.EntityID == discU.ID || s.EntityID == pcU.ID {
			t.Fatalf("Why(task A) leaked foreign/unattributed DPC step: %+v", s)
		}
	}

	whyB, err := eng.Why(ctx, "task", taskB.ID)
	if err != nil {
		t.Fatalf("Why B: %v", err)
	}
	var sawDisc, sawPC bool
	for _, s := range whyB.Steps {
		if s.EntityID == discB.ID && s.ReasonCode == retrieval.ReasonDiscoveryCausesPlanChg {
			sawDisc = true
		}
		if s.EntityID == pcB.ID && s.ReasonCode == retrieval.ReasonDiscoveryCausesPlanChg {
			sawPC = true
		}
		if s.EntityID == discU.ID || s.EntityID == pcU.ID {
			t.Fatalf("Why(task B) must omit unattributed multi-goal DPC: %+v", s)
		}
	}
	if !sawDisc || !sawPC {
		t.Fatalf("Why(task B) missing in-goal DPC: sawDisc=%v sawPC=%v steps=%+v", sawDisc, sawPC, whyB.Steps)
	}
}

// TestWhySymbolExact covers DF-49: Exact/Why by symbol id via GetSymbolByID.
func TestWhySymbolExact(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	if _, err := st.UpsertFile("pkg/sym.go", "hash-sym", nil); err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if err := st.ReplaceFileSymbols("pkg/sym.go", []store.Symbol{
		{Name: "ExactSym", Kind: "function", StartLine: 1, EndLine: 4},
	}); err != nil {
		t.Fatalf("ReplaceFileSymbols: %v", err)
	}
	syms, err := st.ListSymbolsByPath("pkg/sym.go")
	if err != nil || len(syms) != 1 {
		t.Fatalf("ListSymbolsByPath: %v %+v", err, syms)
	}
	symID := syms[0].ID

	hits, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "symbol", EntityID: symID})
	if err != nil || len(hits) != 1 {
		t.Fatalf("Exact symbol: %v %+v", err, hits)
	}
	if hits[0].EntityType != "symbol" || hits[0].Title != "ExactSym" || hits[0].Excerpt != "function" {
		t.Fatalf("Exact hit: %+v", hits[0])
	}
	if hits[0].Path != "pkg/sym.go" {
		t.Fatalf("Exact path: %q", hits[0].Path)
	}

	miss, err := eng.Exact(ctx, retrieval.ExactQuery{
		EntityType: "symbol",
		EntityID:   "00000000-0000-0000-0000-000000000000",
	})
	if err != nil || len(miss) != 0 {
		t.Fatalf("Exact miss want empty: %v %+v", err, miss)
	}

	why, err := eng.Why(ctx, "symbol", symID)
	if err != nil {
		t.Fatalf("Why symbol: %v", err)
	}
	if why.SeedType != "symbol" || why.SeedID != symID {
		t.Fatalf("Why seed: %+v", why)
	}
	if len(why.Steps) == 0 || why.Steps[0].EntityType != "symbol" {
		t.Fatalf("Why steps: %+v", why.Steps)
	}
}

// TestExpandDepth2NoSiblingTaskBody covers DF-35 at retrieval Expand.
func TestExpandDepth2NoSiblingTaskBody(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Shared goal"})
	if err != nil {
		t.Fatalf("goal: %v", err)
	}
	taskA, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Planner task", Body: "plan only"})
	if err != nil {
		t.Fatalf("taskA: %v", err)
	}
	taskB, err := svc.CreateTask(ctx, domain.TaskInput{
		Title: "Implementer sibling",
		Body:  "SECRET_HANDOFF do not leak via depth-2",
	})
	if err != nil {
		t.Fatalf("taskB: %v", err)
	}
	if err := svc.LinkGoalTask(ctx, goal.ID, taskA.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("link A: %v", err)
	}
	if err := svc.LinkGoalTask(ctx, goal.ID, taskB.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("link B: %v", err)
	}

	seed, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "task", EntityID: taskA.ID})
	if err != nil || len(seed) != 1 {
		t.Fatalf("Exact A: %v %+v", err, seed)
	}
	exp, err := eng.Expand(ctx, seed, 2)
	if err != nil {
		t.Fatalf("Expand2: %v", err)
	}
	var sawSiblingTitle bool
	for _, h := range exp {
		if strings.Contains(h.Excerpt, "SECRET_HANDOFF") {
			t.Fatalf("Expand depth-2 leaked sibling body: %+v", h)
		}
		if h.EntityID == taskB.ID {
			sawSiblingTitle = true
			if h.Title != "Implementer sibling" {
				t.Fatalf("sibling title: %+v", h)
			}
			if h.Excerpt != "" {
				t.Fatalf("sibling Excerpt must be empty, got %q", h.Excerpt)
			}
		}
	}
	if !sawSiblingTitle {
		t.Fatalf("expected sibling title in Expand depth-2: %+v", exp)
	}
}

// TestExactWhyPlanChangeAlias covers DF-23: plan-change alias → plan_change.
func TestExactWhyPlanChangeAlias(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()
	pc, err := svc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "Rename alias"})
	if err != nil {
		t.Fatalf("CreatePlanChange: %v", err)
	}
	hits, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "plan-change", EntityID: pc.ID})
	if err != nil || len(hits) != 1 {
		t.Fatalf("Exact plan-change: %v %+v", err, hits)
	}
	if hits[0].EntityType != "plan_change" {
		t.Fatalf("Exact emitted type %q want plan_change", hits[0].EntityType)
	}
	why, err := eng.Why(ctx, "plan-change", pc.ID)
	if err != nil {
		t.Fatalf("Why plan-change: %v", err)
	}
	if why.SeedType != "plan_change" {
		t.Fatalf("Why SeedType %q want plan_change", why.SeedType)
	}
	if len(why.Steps) == 0 || why.Steps[0].EntityType != "plan_change" {
		t.Fatalf("Why seed step type: %+v", why.Steps)
	}
}

// TestExactWhyCapability covers DF-25: Exact/Why for capability via GetCapability.
func TestExactWhyCapability(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()
	cap, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind:   "SKILL",
		Slug:   "review-notes",
		Title:  "Review notes skill",
		Status: "AVAILABLE",
	})
	if err != nil {
		t.Fatalf("UpsertCapability: %v", err)
	}
	hits, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "capability", EntityID: cap.ID})
	if err != nil || len(hits) != 1 {
		t.Fatalf("Exact capability: %v %+v", err, hits)
	}
	if hits[0].EntityType != "capability" || hits[0].Title != "Review notes skill" {
		t.Fatalf("Exact hit: %+v", hits[0])
	}
	if hits[0].Excerpt != "AVAILABLE" {
		t.Fatalf("Exact excerpt want status, got %q", hits[0].Excerpt)
	}
	why, err := eng.Why(ctx, "capability", cap.ID)
	if err != nil {
		t.Fatalf("Why capability: %v", err)
	}
	if why.SeedType != "capability" || why.SeedID != cap.ID {
		t.Fatalf("Why seed: %+v", why)
	}
}

func TestWhyWithVCS(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()
	_, err := st.UpsertFile("tracked.go", "h", nil)
	if err != nil {
		t.Fatalf("file: %v", err)
	}
	fake := &vcs.Fake{
		PathHistory: map[string][]vcs.CommitMeta{
			"tracked.go": {{OID: "abc123", Subject: "touch tracked"}},
		},
	}
	eng.WithVCS(fake)

	hits, err := eng.Exact(ctx, retrieval.ExactQuery{Path: "tracked.go"})
	if err != nil || len(hits) == 0 {
		t.Fatalf("exact: %v", err)
	}
	why, err := eng.Why(ctx, "file", hits[0].EntityID)
	if err != nil {
		t.Fatalf("Why: %v", err)
	}
	saw := false
	for _, s := range why.Steps {
		if s.ReasonCode == retrieval.ReasonHistoricalVCS {
			saw = true
		}
	}
	if !saw {
		t.Fatalf("expected historical_vcs step: %+v", why.Steps)
	}
}

// TestWhyAndContextWithLinkedReview covers DF-01: Why + Expand succeed when a
// review is linked to a task (honesty/D07 plant shape). Full review Hit required.
func TestWhyAndContextWithLinkedReview(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{
		Title: "Implement retrieval review support",
		Body:  "exit: Why/Expand include linked review",
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	rev, err := svc.CreateReview(ctx, domain.ReviewInput{
		Title: "DF-01 linked review",
		Body:  "judge task completion against exit criteria",
	})
	if err != nil {
		t.Fatalf("CreateReview: %v", err)
	}
	if err := svc.LinkReviewTask(ctx, rev.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkReviewTask: %v", err)
	}
	if err := svc.SetReviewResult(ctx, rev.ID, store.ReviewResultPass, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "exit criteria met",
	}); err != nil {
		t.Fatalf("SetReviewResult PASS: %v", err)
	}

	exact, err := eng.Exact(ctx, retrieval.ExactQuery{EntityType: "review", EntityID: rev.ID})
	if err != nil || len(exact) != 1 {
		t.Fatalf("Exact review: err=%v hits=%+v", err, exact)
	}
	if exact[0].EntityType != "review" || exact[0].Title != "DF-01 linked review" {
		t.Fatalf("Exact review hit shape: %+v", exact[0])
	}
	if exact[0].Excerpt != store.ReviewResultPass {
		t.Fatalf("Exact review Excerpt want PASS got %q", exact[0].Excerpt)
	}

	why, err := eng.Why(ctx, "task", task.ID)
	if err != nil {
		t.Fatalf("Why: %v", err)
	}
	var sawWhyReview bool
	for _, s := range why.Steps {
		if s.EntityType == "review" && s.EntityID == rev.ID && s.ReasonCode == retrieval.ReasonReviewJudgesTask {
			sawWhyReview = true
		}
	}
	if !sawWhyReview {
		t.Fatalf("Why missing review neighbor: %+v", why.Steps)
	}

	seed := retrieval.Hit{EntityType: "task", EntityID: task.ID, Title: task.Title, ReasonCode: retrieval.ReasonExactID}
	exp, err := eng.Expand(ctx, []retrieval.Hit{seed}, 1)
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	var sawExpandReview bool
	for _, h := range exp {
		if h.EntityType == "review" && h.EntityID == rev.ID && h.ReasonCode == retrieval.ReasonReviewJudgesTask {
			sawExpandReview = true
			if h.Title != "DF-01 linked review" {
				t.Fatalf("Expand review Title: %+v", h)
			}
			if h.Excerpt != store.ReviewResultPass {
				t.Fatalf("Expand review Excerpt want PASS got %q", h.Excerpt)
			}
		}
	}
	if !sawExpandReview {
		t.Fatalf("Expand missing review neighbor: %+v", exp)
	}
}

// TestExpandImportEdgeProvenance + TestWhySurfacesEdgeProvenance: structural import
// hops carry edge_provenance; INFERRED via store fixture (not analyzer).
func TestExpandImportEdgeProvenance(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	src, err := st.UpsertFile("src/main.go", "h1", nil)
	if err != nil {
		t.Fatalf("UpsertFile src: %v", err)
	}
	if _, err := st.UpsertFile("src/util.go", "h2", nil); err != nil {
		t.Fatalf("UpsertFile util: %v", err)
	}
	if _, err := st.UpsertFile("src/inferred.go", "h3", nil); err != nil {
		t.Fatalf("UpsertFile inferred: %v", err)
	}
	if err := st.ReplaceFileImports("src/main.go", []store.Import{
		{ImportedPath: "src/util.go", Provenance: store.ImportProvenanceExtracted},
		{ImportedPath: "src/inferred.go", Provenance: store.ImportProvenanceInferred},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	seed := retrieval.Hit{EntityType: "file", EntityID: src.ID, Path: "src/main.go", Title: "src/main.go", ReasonCode: retrieval.ReasonExactPath}
	exp, err := eng.Expand(ctx, []retrieval.Hit{seed}, 1)
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	got := map[string]string{}
	for _, h := range exp {
		if h.EntityType == "file" && h.ReasonCode == retrieval.ReasonGraphNeighbor {
			got[h.Path] = h.EdgeProvenance
		}
	}
	if got["src/util.go"] != store.ImportProvenanceExtracted {
		t.Fatalf("util edge_provenance: got %q want EXTRACTED; hits=%+v", got["src/util.go"], exp)
	}
	if got["src/inferred.go"] != store.ImportProvenanceInferred {
		t.Fatalf("inferred edge_provenance: got %q want INFERRED (must not silent EXTRACTED); hits=%+v", got["src/inferred.go"], exp)
	}
}

func TestWhySurfacesEdgeProvenance(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	src, err := st.UpsertFile("why/a.go", "ha", nil)
	if err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if _, err := st.UpsertFile("why/b.go", "hb", nil); err != nil {
		t.Fatalf("UpsertFile b: %v", err)
	}
	if err := st.ReplaceFileImports("why/a.go", []store.Import{
		{ImportedPath: "why/b.go", Provenance: store.ImportProvenanceInferred},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	why, err := eng.Why(ctx, "file", src.ID)
	if err != nil {
		t.Fatalf("Why: %v", err)
	}
	var saw bool
	for _, s := range why.Steps {
		if s.EntityType == "file" && s.Title == "why/b.go" && s.ReasonCode == retrieval.ReasonGraphNeighbor {
			saw = true
			if s.EdgeProvenance != store.ImportProvenanceInferred {
				t.Fatalf("Why step edge_provenance=%q want INFERRED", s.EdgeProvenance)
			}
		}
	}
	if !saw {
		t.Fatalf("Why missing import neighbor: %+v", why.Steps)
	}
}

// DF-64: empty Provenance on write surfaces as EXTRACTED on Expand (not omitempty-hidden).
func TestExpandEmptyProvenanceSurfacesExtracted(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	src, err := st.UpsertFile("empty/main.go", "h1", nil)
	if err != nil {
		t.Fatalf("UpsertFile src: %v", err)
	}
	if _, err := st.UpsertFile("empty/util.go", "h2", nil); err != nil {
		t.Fatalf("UpsertFile util: %v", err)
	}
	if err := st.ReplaceFileImports("empty/main.go", []store.Import{
		{ImportedPath: "empty/util.go"}, // empty Provenance → EXTRACTED
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	seed := retrieval.Hit{EntityType: "file", EntityID: src.ID, Path: "empty/main.go", Title: "empty/main.go", ReasonCode: retrieval.ReasonExactPath}
	exp, err := eng.Expand(ctx, []retrieval.Hit{seed}, 1)
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	var saw string
	for _, h := range exp {
		if h.EntityType == "file" && h.Path == "empty/util.go" && h.ReasonCode == retrieval.ReasonGraphNeighbor {
			saw = h.EdgeProvenance
		}
	}
	if saw != store.ImportProvenanceExtracted {
		t.Fatalf("edge_provenance: got %q want EXTRACTED (must not be empty/hidden); hits=%+v", saw, exp)
	}
}

// DF-60: analyzer-shaped relative imports under a subdirectory resolve via Expand.
func TestExpandSubdirRelativeImportJS(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	src, err := st.UpsertFile("src/app.js", "h1", nil)
	if err != nil {
		t.Fatalf("UpsertFile app: %v", err)
	}
	if _, err := st.UpsertFile("src/util.js", "h2", nil); err != nil {
		t.Fatalf("UpsertFile util: %v", err)
	}
	if err := st.ReplaceFileImports("src/app.js", []store.Import{
		{ImportedPath: "./util.js", Provenance: store.ImportProvenanceExtracted},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	seed := retrieval.Hit{EntityType: "file", EntityID: src.ID, Path: "src/app.js", Title: "src/app.js", ReasonCode: retrieval.ReasonExactPath}
	exp, err := eng.Expand(ctx, []retrieval.Hit{seed}, 1)
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	var saw bool
	for _, h := range exp {
		if h.EntityType == "file" && h.Path == "src/util.js" && h.ReasonCode == retrieval.ReasonGraphNeighbor {
			saw = true
			if h.EdgeProvenance != store.ImportProvenanceExtracted {
				t.Fatalf("edge_provenance=%q want EXTRACTED", h.EdgeProvenance)
			}
		}
	}
	if !saw {
		t.Fatalf("Expand missing subdir neighbor src/util.js: %+v", exp)
	}
}

func TestExpandSubdirExtensionlessImport(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	src, err := st.UpsertFile("src/app.js", "h1", nil)
	if err != nil {
		t.Fatalf("UpsertFile app: %v", err)
	}
	if _, err := st.UpsertFile("src/util.ts", "h2", nil); err != nil {
		t.Fatalf("UpsertFile util: %v", err)
	}
	if err := st.ReplaceFileImports("src/app.js", []store.Import{
		{ImportedPath: "./util", Provenance: store.ImportProvenanceExtracted},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	seed := retrieval.Hit{EntityType: "file", EntityID: src.ID, Path: "src/app.js", Title: "src/app.js", ReasonCode: retrieval.ReasonExactPath}
	exp, err := eng.Expand(ctx, []retrieval.Hit{seed}, 1)
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	var saw bool
	for _, h := range exp {
		if h.EntityType == "file" && h.Path == "src/util.ts" && h.ReasonCode == retrieval.ReasonGraphNeighbor {
			saw = true
			if h.EdgeProvenance != store.ImportProvenanceExtracted {
				t.Fatalf("edge_provenance=%q want EXTRACTED", h.EdgeProvenance)
			}
		}
	}
	if !saw {
		t.Fatalf("Expand missing extensionless neighbor src/util.ts: %+v", exp)
	}
}

func TestExpandRootRelativeImportPositive(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	src, err := st.UpsertFile("app.js", "h1", nil)
	if err != nil {
		t.Fatalf("UpsertFile app: %v", err)
	}
	if _, err := st.UpsertFile("util.js", "h2", nil); err != nil {
		t.Fatalf("UpsertFile util: %v", err)
	}
	if err := st.ReplaceFileImports("app.js", []store.Import{
		{ImportedPath: "./util.js", Provenance: store.ImportProvenanceExtracted},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	seed := retrieval.Hit{EntityType: "file", EntityID: src.ID, Path: "app.js", Title: "app.js", ReasonCode: retrieval.ReasonExactPath}
	exp, err := eng.Expand(ctx, []retrieval.Hit{seed}, 1)
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	var saw bool
	for _, h := range exp {
		if h.EntityType == "file" && h.Path == "util.js" && h.ReasonCode == retrieval.ReasonGraphNeighbor {
			saw = true
			if h.EdgeProvenance != store.ImportProvenanceExtracted {
				t.Fatalf("edge_provenance=%q want EXTRACTED", h.EdgeProvenance)
			}
		}
	}
	if !saw {
		t.Fatalf("Expand missing root neighbor util.js: %+v", exp)
	}
}

// DF-60: parent-relative ../ under a nested importer joins Dir(importer).
func TestExpandParentRelativeImport(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	src, err := st.UpsertFile("src/nested/app.js", "h1", nil)
	if err != nil {
		t.Fatalf("UpsertFile app: %v", err)
	}
	if _, err := st.UpsertFile("src/util.js", "h2", nil); err != nil {
		t.Fatalf("UpsertFile util: %v", err)
	}
	if err := st.ReplaceFileImports("src/nested/app.js", []store.Import{
		{ImportedPath: "../util.js", Provenance: store.ImportProvenanceExtracted},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	seed := retrieval.Hit{EntityType: "file", EntityID: src.ID, Path: "src/nested/app.js", Title: "src/nested/app.js", ReasonCode: retrieval.ReasonExactPath}
	exp, err := eng.Expand(ctx, []retrieval.Hit{seed}, 1)
	if err != nil {
		t.Fatalf("Expand: %v", err)
	}
	var saw bool
	for _, h := range exp {
		if h.EntityType == "file" && h.Path == "src/util.js" && h.ReasonCode == retrieval.ReasonGraphNeighbor {
			saw = true
			if h.EdgeProvenance != store.ImportProvenanceExtracted {
				t.Fatalf("edge_provenance=%q want EXTRACTED", h.EdgeProvenance)
			}
		}
	}
	if !saw {
		t.Fatalf("Expand missing parent-relative neighbor src/util.js: %+v", exp)
	}
}

func TestWhySurfacesSubdirRelativeImportProvenance(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	src, err := st.UpsertFile("src/app.js", "ha", nil)
	if err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if _, err := st.UpsertFile("src/util.js", "hb", nil); err != nil {
		t.Fatalf("UpsertFile util: %v", err)
	}
	if err := st.ReplaceFileImports("src/app.js", []store.Import{
		{ImportedPath: "./util.js", Provenance: store.ImportProvenanceExtracted},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	why, err := eng.Why(ctx, "file", src.ID)
	if err != nil {
		t.Fatalf("Why: %v", err)
	}
	var saw bool
	for _, s := range why.Steps {
		if s.EntityType == "file" && s.Title == "src/util.js" && s.ReasonCode == retrieval.ReasonGraphNeighbor {
			saw = true
			if s.EdgeProvenance != store.ImportProvenanceExtracted {
				t.Fatalf("Why step edge_provenance=%q want EXTRACTED", s.EdgeProvenance)
			}
		}
	}
	if !saw {
		t.Fatalf("Why missing subdir import neighbor: %+v", why.Steps)
	}
}
