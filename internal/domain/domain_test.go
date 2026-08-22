package domain_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func openDomain(t *testing.T) (*domain.Service, *store.Store) {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return domain.New(st), st
}

func TestCreateRoundtripProvenance(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Ship causal", Confidence: 0.9})
	if err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	if g.SourceType != domain.DefaultSourceType || g.Status != store.StatusActive || g.Confidence != 0.9 {
		t.Fatalf("goal provenance: %+v", g)
	}
	got, err := svc.GetGoal(ctx, g.ID)
	if err != nil || got.Title != "Ship causal" {
		t.Fatalf("GetGoal: %+v err=%v", got, err)
	}

	if _, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "SQLite", Body: "local", Confidence: 1}); err != nil {
		t.Fatalf("CreateDecision: %v", err)
	}
	if _, err := svc.CreateAssumption(ctx, domain.AssumptionInput{Title: "CGO free store"}); err != nil {
		t.Fatalf("CreateAssumption: %v", err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Wire domain", GoalID: &g.ID})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if task.WorkState != store.WorkStatePending {
		t.Fatalf("default work_state: %q", task.WorkState)
	}
	if task.GoalID == nil || *task.GoalID != g.ID {
		t.Fatalf("task goal_id: %+v", task.GoalID)
	}
	if _, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "Schema gap"}); err != nil {
		t.Fatalf("CreateDiscovery: %v", err)
	}
	if _, err := svc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "Add work_state"}); err != nil {
		t.Fatalf("CreatePlanChange: %v", err)
	}

	created, err := st.ListEventsByEntity(domain.EntityGoal, g.ID)
	if err != nil || len(created) < 1 || created[0].Type != domain.EventEntityCreated {
		t.Fatalf("entity.created: %+v err=%v", created, err)
	}

	if _, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "  "}); err == nil {
		t.Fatal("expected empty title rejection")
	}
}

func TestGoalTaskLinkViaGoalID(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "T"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkGoalTask(ctx, g.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkGoalTask: %v", err)
	}
	got, err := svc.GetTask(ctx, task.ID)
	if err != nil || got.GoalID == nil || *got.GoalID != g.ID {
		t.Fatalf("goal_id not set: %+v err=%v", got, err)
	}
	links, err := st.ListLinksFrom(domain.EntityGoal, g.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 0 {
		t.Fatalf("Goal→Task must not use entity_links; got %+v", links)
	}
	evs, err := st.ListEventsByEntity(domain.EntityGoal, g.ID)
	if err != nil {
		t.Fatal(err)
	}
	var linked bool
	for _, e := range evs {
		if e.Type == domain.EventEntityLinked {
			linked = true
			var p map[string]any
			_ = json.Unmarshal([]byte(e.PayloadJSON), &p)
			if p["rel"] != domain.RelGoalHasTaskEvent {
				t.Fatalf("link event rel: %#v", p)
			}
		}
	}
	if !linked {
		t.Fatal("expected entity.linked for Goal→Task")
	}
}

func TestDecisionAffectsTaskLink(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	d, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "D"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "T"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkDecisionTask(ctx, d.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDecisionTask: %v", err)
	}
	links, err := svc.ListLinksFrom(ctx, domain.EntityDecision, d.ID)
	if err != nil || len(links) != 1 {
		t.Fatalf("links: %+v err=%v", links, err)
	}
	if links[0].Rel != domain.RelDecisionAffectsTask || links[0].ToID != task.ID {
		t.Fatalf("link: %+v", links[0])
	}
	evs, _ := st.ListEventsByEntity(domain.EntityDecision, d.ID)
	found := false
	for _, e := range evs {
		if e.Type == domain.EventEntityLinked {
			found = true
		}
	}
	if !found {
		t.Fatal("expected entity.linked")
	}
}

func TestDiscoveryCausesPlanChangeLink(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	disc, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "Found gap"})
	if err != nil {
		t.Fatal(err)
	}
	pc, err := svc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "Fix gap"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkDiscoveryPlanChange(ctx, disc.ID, pc.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDiscoveryPlanChange: %v", err)
	}
	links, err := svc.ListLinksFrom(ctx, domain.EntityDiscovery, disc.ID)
	if err != nil || len(links) != 1 || links[0].Rel != domain.RelDiscoveryCausesPlanChange {
		t.Fatalf("links: %+v err=%v", links, err)
	}
}

func TestLinkDiscoveryMentionsTask(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	disc, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "Mentions work"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Attributed task"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkDiscoveryMentionsTask(ctx, disc.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDiscoveryMentionsTask: %v", err)
	}
	links, err := svc.ListLinksFrom(ctx, domain.EntityDiscovery, disc.ID)
	if err != nil || len(links) != 1 {
		t.Fatalf("links: %+v err=%v", links, err)
	}
	if links[0].Rel != domain.RelDiscoveryMentionsTask || links[0].ToID != task.ID {
		t.Fatalf("link: %+v", links[0])
	}
	evs, _ := st.ListEventsByEntity(domain.EntityDiscovery, disc.ID)
	found := false
	for _, e := range evs {
		if e.Type == domain.EventEntityLinked {
			found = true
		}
	}
	if !found {
		t.Fatal("expected entity.linked")
	}
}

func TestTransitionLegalAndIllegal(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Work"})
	if err != nil {
		t.Fatal(err)
	}

	err = svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "alice", Reason: "too soon", AllowDoneWithoutReview: true,
	})
	var inv *domain.ErrInvalidTransition
	if !errors.As(err, &inv) {
		t.Fatalf("PENDING→DONE should fail closed: %v", err)
	}

	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "alice", Reason: "start",
	}); err != nil {
		t.Fatalf("PENDING→IN_PROGRESS: %v", err)
	}
	got, _ := svc.GetTask(ctx, task.ID)
	if got.WorkState != store.WorkStateInProgress {
		t.Fatalf("work_state: %q", got.WorkState)
	}

	evs, err := st.ListEventsByEntity(domain.EntityTask, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	var tr *store.Event
	for i := range evs {
		if evs[i].Type == domain.EventTaskTransition {
			tr = &evs[i]
		}
	}
	if tr == nil {
		t.Fatal("expected task.transition event")
	}
	var p map[string]any
	if err := json.Unmarshal([]byte(tr.PayloadJSON), &p); err != nil {
		t.Fatal(err)
	}
	if p["actor"] != "alice" || p["from"] != store.WorkStatePending || p["to"] != store.WorkStateInProgress || p["reason"] != "start" {
		t.Fatalf("payload: %#v", p)
	}
	if _, ok := p["evidence_ids"]; !ok {
		t.Fatal("evidence_ids required in payload")
	}
}

func TestDonePolicyStub(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Finish"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "go",
	}); err != nil {
		t.Fatal(err)
	}

	err = svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "done?",
	})
	if err == nil {
		t.Fatal("DONE without flag/review must fail")
	}
	if !strings.Contains(err.Error(), "DONE") {
		t.Fatalf("unexpected err: %v", err)
	}

	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "ship", AllowDoneWithoutReview: true,
	}); err != nil {
		t.Fatalf("AllowDoneWithoutReview: %v", err)
	}
	got, _ := svc.GetTask(ctx, task.ID)
	if got.WorkState != store.WorkStateDone {
		t.Fatalf("want DONE, got %q", got.WorkState)
	}

	// EvidenceIDs alone must NOT authorize DONE (replaced stub).
	task2, _ := svc.CreateTask(ctx, domain.TaskInput{Title: "With evidence only"})
	_ = svc.TransitionTask(ctx, task2.ID, store.WorkStateInProgress, domain.TransitionOptions{Actor: "a", Reason: "go"})
	err = svc.TransitionTask(ctx, task2.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "proven?", EvidenceIDs: []string{"ev-1"},
	})
	if err == nil {
		t.Fatal("DONE with EvidenceIDs alone must fail")
	}
	if !strings.Contains(err.Error(), "DONE") {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestDoneRequiresReviewPass(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	mkInProgress := func(title string) store.Task {
		t.Helper()
		task, err := svc.CreateTask(ctx, domain.TaskInput{Title: title})
		if err != nil {
			t.Fatal(err)
		}
		if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
			Actor: "impl", Reason: "start",
		}); err != nil {
			t.Fatal(err)
		}
		return task
	}

	// No review → reject
	t1 := mkInProgress("no review")
	if err := svc.TransitionTask(ctx, t1.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "impl", Reason: "self claim",
	}); err == nil {
		t.Fatal("DONE with no review must fail")
	}

	// FAIL review → reject
	t2 := mkInProgress("fail review")
	revFail, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "Reject"})
	if err != nil {
		t.Fatalf("CreateReview: %v", err)
	}
	if err := svc.LinkReviewTask(ctx, revFail.ID, t2.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkReviewTask: %v", err)
	}
	if err := svc.SetReviewResult(ctx, revFail.ID, store.ReviewResultFail, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "not proven",
	}); err != nil {
		t.Fatalf("SetReviewResult FAIL: %v", err)
	}
	if err := svc.TransitionTask(ctx, t2.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "impl", Reason: "try done",
	}); err == nil {
		t.Fatal("DONE with FAIL review must fail")
	}

	// UNCERTAIN → reject
	t3 := mkInProgress("uncertain")
	revU, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "Unsure"})
	if err != nil {
		t.Fatal(err)
	}
	_ = svc.LinkReviewTask(ctx, revU.ID, t3.ID, domain.LinkMeta{})
	_ = svc.SetReviewResult(ctx, revU.ID, store.ReviewResultUncertain, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "needs more",
	})
	if err := svc.TransitionTask(ctx, t3.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "impl", Reason: "try",
	}); err == nil {
		t.Fatal("DONE with UNCERTAIN review must fail")
	}

	// Linked review still open (no SetReviewResult) → reject
	t4 := mkInProgress("open review")
	revOpen, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "Pending"})
	if err != nil {
		t.Fatal(err)
	}
	_ = svc.LinkReviewTask(ctx, revOpen.ID, t4.ID, domain.LinkMeta{})
	if err := svc.TransitionTask(ctx, t4.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "impl", Reason: "too early",
	}); err == nil {
		t.Fatal("DONE with open review must fail")
	}

	// PASS linked → accept; task.transition may include review_id
	t5 := mkInProgress("pass path")
	revPass, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "Approve"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, revPass.ID, t5.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, revPass.ID, store.ReviewResultPass, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "looks good",
	}); err != nil {
		t.Fatalf("SetReviewResult PASS: %v", err)
	}
	gotRev, err := svc.GetReview(ctx, revPass.ID)
	if err != nil || gotRev.Result != store.ReviewResultPass {
		t.Fatalf("GetReview: %+v err=%v", gotRev, err)
	}
	if err := svc.TransitionTask(ctx, t5.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "impl", Reason: "promote", EvidenceIDs: []string{"ev-optional"},
		AllowOperatorDone: true,
	}); err != nil {
		t.Fatalf("DONE after PASS: %v", err)
	}
	got, _ := svc.GetTask(ctx, t5.ID)
	if got.WorkState != store.WorkStateDone {
		t.Fatalf("want DONE, got %q", got.WorkState)
	}

	evs, err := st.ListEventsByEntity(domain.EntityTask, t5.ID)
	if err != nil {
		t.Fatal(err)
	}
	var tr *store.Event
	for i := range evs {
		if evs[i].Type == domain.EventTaskTransition {
			var p map[string]any
			_ = json.Unmarshal([]byte(evs[i].PayloadJSON), &p)
			if p["to"] == store.WorkStateDone {
				tr = &evs[i]
			}
		}
	}
	if tr == nil {
		t.Fatal("expected DONE task.transition")
	}
	var p map[string]any
	_ = json.Unmarshal([]byte(tr.PayloadJSON), &p)
	if p["review_id"] != revPass.ID {
		t.Fatalf("expected review_id in payload, got %#v", p)
	}

	// review.result event on SetReviewResult
	revEvs, _ := st.ListEventsByEntity(domain.EntityReview, revPass.ID)
	foundResult := false
	for _, e := range revEvs {
		if e.Type == domain.EventReviewResult {
			foundResult = true
			var rp map[string]any
			_ = json.Unmarshal([]byte(e.PayloadJSON), &rp)
			if rp["result"] != store.ReviewResultPass || rp["actor"] != "reviewer" {
				t.Fatalf("review.result payload: %#v", rp)
			}
		}
	}
	if !foundResult {
		t.Fatal("expected review.result event")
	}

	links, err := svc.ListLinksFrom(ctx, domain.EntityReview, revPass.ID)
	if err != nil || len(links) != 1 || links[0].Rel != domain.RelReviewJudgesTask {
		t.Fatalf("review_judges_task links: %+v err=%v", links, err)
	}
}

func TestMarkStale(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Old"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.MarkStale(ctx, domain.EntityGoal, g.ID, "superseded by new plan"); err != nil {
		t.Fatalf("MarkStale: %v", err)
	}
	got, _ := svc.GetGoal(ctx, g.ID)
	if got.Status != store.StatusStale {
		t.Fatalf("status: %q", got.Status)
	}
	evs, _ := st.ListEventsByEntity(domain.EntityGoal, g.ID)
	found := false
	for _, e := range evs {
		if e.Type == "entity.stale" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected entity.stale event")
	}
}

func TestClaimEvidenceStub(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	c, err := svc.CreateClaim(ctx, domain.ClaimInput{Title: "Store is CGO-free"})
	if err != nil {
		t.Fatalf("CreateClaim: %v", err)
	}
	e, err := svc.CreateEvidence(ctx, domain.EvidenceInput{Title: "modernc.org/sqlite"})
	if err != nil {
		t.Fatalf("CreateEvidence: %v", err)
	}
	if err := svc.LinkClaimEvidence(ctx, c.ID, e.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkClaimEvidence: %v", err)
	}
	links, err := svc.ListLinksFrom(ctx, domain.EntityClaim, c.ID)
	if err != nil || len(links) != 1 || links[0].Rel != domain.RelClaimHasEvidence {
		t.Fatalf("links: %+v err=%v", links, err)
	}
}

func seedPlanScope(t *testing.T, st *store.Store, goalID string) store.PlanScope {
	t.Helper()
	ph, err := st.InsertPlanPhase(store.PlanPhase{GoalID: goalID, Title: "P1", Ord: 0})
	if err != nil {
		t.Fatalf("InsertPlanPhase: %v", err)
	}
	sc, err := st.InsertPlanScope(store.PlanScope{PhaseID: ph.ID, Title: "S1", Ord: 0})
	if err != nil {
		t.Fatalf("InsertPlanScope: %v", err)
	}
	return sc
}

func TestLinkReviewScopeAndResiduals(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Scope review goal"})
	if err != nil {
		t.Fatal(err)
	}
	sc := seedPlanScope(t, st, g.ID)

	rev, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "Scope check"})
	if err != nil {
		t.Fatalf("CreateReview: %v", err)
	}

	if err := svc.LinkReviewScope(ctx, rev.ID, "missing-scope", domain.LinkMeta{}); err == nil {
		t.Fatal("LinkReviewScope missing scope must fail")
	}
	if err := svc.LinkReviewScope(ctx, "missing-review", sc.ID, domain.LinkMeta{}); err == nil {
		t.Fatal("LinkReviewScope missing review must fail")
	}

	if err := svc.LinkReviewScope(ctx, rev.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkReviewScope: %v", err)
	}
	links, err := svc.ListLinksFrom(ctx, domain.EntityReview, rev.ID)
	if err != nil || len(links) != 1 || links[0].Rel != domain.RelReviewJudgesScope ||
		links[0].ToType != domain.EntityPlanScope || links[0].ToID != sc.ID {
		t.Fatalf("review_judges_scope links: %+v err=%v", links, err)
	}

	// Scope status must not change from linking.
	gotSc, err := st.GetPlanScope(sc.ID)
	if err != nil || gotSc.Status != sc.Status {
		t.Fatalf("plan_scopes.status mutated: before=%q after=%+v err=%v", sc.Status, gotSc, err)
	}

	if _, err := svc.AddResidual(ctx, rev.ID, domain.ResidualInput{Code: ""}); err == nil {
		t.Fatal("empty residual code must fail")
	}
	if _, err := svc.AddResidual(ctx, rev.ID, domain.ResidualInput{
		Code: domain.ResidualCodeOpenGap, Severity: "CRITICAL",
	}); err == nil {
		t.Fatal("bad residual severity must fail")
	}

	res, err := svc.AddResidual(ctx, rev.ID, domain.ResidualInput{
		Code: domain.ResidualCodeMissingEvidence, Body: "need logs",
	})
	if err != nil {
		t.Fatalf("AddResidual: %v", err)
	}
	if res.Severity != domain.ResidualSeverityINFO || res.Status != domain.ResidualStatusOpen {
		t.Fatalf("defaults: %+v", res)
	}

	byRev, err := svc.ListResidualsByReview(ctx, rev.ID)
	if err != nil || len(byRev) != 1 || byRev[0].ID != res.ID {
		t.Fatalf("ListResidualsByReview: %+v err=%v", byRev, err)
	}
	byScope, err := svc.ListResidualsByScope(ctx, sc.ID)
	if err != nil || len(byScope) != 1 || byScope[0].Code != domain.ResidualCodeMissingEvidence {
		t.Fatalf("ListResidualsByScope: %+v err=%v", byScope, err)
	}
	n, err := svc.CountOpenResidualsByScope(ctx, sc.ID)
	if err != nil || n != 1 {
		t.Fatalf("CountOpenResidualsByScope: n=%d err=%v", n, err)
	}

	if err := svc.SetResidualStatus(ctx, res.ID, domain.ResidualStatusAcked, domain.ResidualStatusOptions{}); err == nil {
		t.Fatal("SetResidualStatus without actor/reason must fail")
	}
	if err := svc.SetResidualStatus(ctx, res.ID, "NOPE", domain.ResidualStatusOptions{
		Actor: "r", Reason: "bad",
	}); err == nil {
		t.Fatal("unknown residual status must fail")
	}
	if err := svc.SetResidualStatus(ctx, res.ID, domain.ResidualStatusAcked, domain.ResidualStatusOptions{
		Actor: "reviewer", Reason: "noted",
	}); err != nil {
		t.Fatalf("SetResidualStatus ACKED: %v", err)
	}
	n, err = svc.CountOpenResidualsByScope(ctx, sc.ID)
	if err != nil || n != 0 {
		t.Fatalf("after ACKED open count want 0, got %d err=%v", n, err)
	}
	if err := svc.SetResidualStatus(ctx, res.ID, domain.ResidualStatusResolved, domain.ResidualStatusOptions{
		Actor: "reviewer", Reason: "closed",
	}); err != nil {
		t.Fatalf("SetResidualStatus RESOLVED: %v", err)
	}
}

func TestScopeReviewDoesNotWeakenDoneGate(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "DONE gate"})
	if err != nil {
		t.Fatal(err)
	}
	sc := seedPlanScope(t, st, g.ID)

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "impl"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}

	// EvidenceIDs alone still rejected.
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "claim", EvidenceIDs: []string{"ev-1"},
	}); err == nil {
		t.Fatal("EvidenceIDs alone must still fail DONE")
	}

	// Scope review PASS alone must not authorize task DONE.
	scopeRev, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "scope only"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewScope(ctx, scopeRev.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, scopeRev.ID, store.ReviewResultPass, domain.ReviewResultOptions{
		Actor: "r", Reason: "scope ok",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "try via scope",
	}); err == nil {
		t.Fatal("scope PASS must not authorize task DONE")
	}

	// Task PASS review_judges_task still works.
	taskRev, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "task pass"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, taskRev.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, taskRev.ID, store.ReviewResultPass, domain.ReviewResultOptions{
		Actor: "r", Reason: "ok",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "promote", AllowOperatorDone: true,
	}); err != nil {
		t.Fatalf("DONE after task PASS: %v", err)
	}

	// Escape hatch still explicit.
	task2, _ := svc.CreateTask(ctx, domain.TaskInput{Title: "escape"})
	_ = svc.TransitionTask(ctx, task2.ID, store.WorkStateInProgress, domain.TransitionOptions{Actor: "a", Reason: "go"})
	if err := svc.TransitionTask(ctx, task2.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "ship", AllowDoneWithoutReview: true,
	}); err != nil {
		t.Fatalf("AllowDoneWithoutReview: %v", err)
	}
}

func TestImpactFindingAddListAndReject(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	d, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Use SQLite"})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: "", Kind: domain.FindingKindAffectedWork,
	}); err == nil {
		t.Fatal("empty impact_class must fail")
	}
	if _, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: "MILD", Kind: domain.FindingKindAffectedWork,
	}); err == nil {
		t.Fatal("unknown impact_class must fail")
	}
	if _, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassSAFE, Kind: "",
	}); err == nil {
		t.Fatal("empty kind must fail")
	}
	if _, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassSAFE, Kind: "OTHER",
	}); err == nil {
		t.Fatal("unknown kind must fail")
	}
	if _, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassSAFE, Kind: domain.FindingKindAffectedWork, Uncertainty: "GUESS",
	}); err == nil {
		t.Fatal("unknown uncertainty must fail")
	}

	f, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: "caution", Kind: "affected_work", Body: "touch tasks",
	})
	if err != nil {
		t.Fatalf("AddImpactFinding: %v", err)
	}
	if f.ImpactClass != domain.ImpactClassCaution || f.Kind != domain.FindingKindAffectedWork {
		t.Fatalf("normalized: %+v", f)
	}
	if f.Uncertainty != domain.UncertaintyUNKNOWN {
		t.Fatalf("empty uncertainty → UNKNOWN, got %q", f.Uncertainty)
	}

	list, err := svc.ListImpactFindings(ctx, d.ID)
	if err != nil || len(list) != 1 || list[0].ID != f.ID {
		t.Fatalf("ListImpactFindings: %+v err=%v", list, err)
	}
}

func TestDecisionAlternativeRecommendExclusivity(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	d, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Pick store"})
	if err != nil {
		t.Fatal(err)
	}
	a1, err := svc.AddDecisionAlternative(ctx, d.ID, domain.AlternativeInput{
		Title: "SQLite", Recommended: true,
	})
	if err != nil {
		t.Fatalf("AddDecisionAlternative a1: %v", err)
	}
	if !a1.IsRecommended {
		t.Fatalf("a1 should be recommended: %+v", a1)
	}
	a2, err := svc.AddDecisionAlternative(ctx, d.ID, domain.AlternativeInput{
		Title: "Postgres", Recommended: true,
	})
	if err != nil {
		t.Fatalf("AddDecisionAlternative a2: %v", err)
	}
	list, err := svc.ListDecisionAlternatives(ctx, d.ID)
	if err != nil || len(list) != 2 {
		t.Fatalf("ListDecisionAlternatives: %+v err=%v", list, err)
	}
	var recCount int
	for _, a := range list {
		if a.IsRecommended {
			recCount++
			if a.ID != a2.ID {
				t.Fatalf("only a2 should be recommended, got %s", a.ID)
			}
		}
	}
	if recCount != 1 {
		t.Fatalf("want exactly one recommended, got %d", recCount)
	}

	if err := svc.SetRecommendedAlternative(ctx, d.ID, a1.ID); err != nil {
		t.Fatalf("SetRecommendedAlternative: %v", err)
	}
	list, _ = svc.ListDecisionAlternatives(ctx, d.ID)
	recCount = 0
	for _, a := range list {
		if a.IsRecommended {
			recCount++
			if a.ID != a1.ID {
				t.Fatalf("after set, only a1 recommended, got %s", a.ID)
			}
		}
	}
	if recCount != 1 {
		t.Fatalf("after set want 1 recommended, got %d", recCount)
	}

	other, _ := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Other"})
	if err := svc.SetRecommendedAlternative(ctx, other.ID, a1.ID); err == nil {
		t.Fatal("SetRecommendedAlternative wrong decision must fail")
	}
}

func TestImpactReportFailClosedAndRollup(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	d, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Impact report"})
	if err != nil {
		t.Fatal(err)
	}

	// Empty findings + no links → OverallClass "" + Incomplete.
	rep, err := svc.ImpactReport(ctx, d.ID)
	if err != nil {
		t.Fatalf("ImpactReport empty: %v", err)
	}
	if rep.OverallClass != "" || !rep.Incomplete || rep.HasUnknown {
		t.Fatalf("empty+no links: %+v", rep)
	}
	if rep.OverallUncertainty != domain.UncertaintyUNKNOWN {
		t.Fatalf("empty findings OverallUncertainty: %q", rep.OverallUncertainty)
	}

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "affected"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkDecisionTask(ctx, d.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkDecisionTask: %v", err)
	}

	// Empty findings + ≥1 task link → HasUnknown + Incomplete.
	rep, err = svc.ImpactReport(ctx, d.ID)
	if err != nil {
		t.Fatalf("ImpactReport linked: %v", err)
	}
	if !rep.HasUnknown || !rep.Incomplete || rep.OverallClass != "" {
		t.Fatalf("empty+link: %+v", rep)
	}
	if len(rep.AffectedTaskIDs) != 1 || rep.AffectedTaskIDs[0] != task.ID {
		t.Fatalf("AffectedTaskIDs: %+v", rep.AffectedTaskIDs)
	}
	if len(rep.AffectsTaskLinks) != 1 || rep.AffectsTaskLinks[0].Rel != domain.RelDecisionAffectsTask {
		t.Fatalf("AffectsTaskLinks: %+v", rep.AffectsTaskLinks)
	}

	// OverallClass rollup SAFE < CAUTION < HIGH < DESTRUCTIVE < REVERSAL.
	if _, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassSAFE, Kind: domain.FindingKindNewWork, Uncertainty: domain.UncertaintyKNOWN,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassHigh, Kind: domain.FindingKindWorkAtRisk, Uncertainty: domain.UncertaintyLIKELY,
	}); err != nil {
		t.Fatal(err)
	}
	unk, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassCaution, Kind: domain.FindingKindUnresolved, Uncertainty: domain.UncertaintyUNKNOWN,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassDestructive, Kind: domain.FindingKindInvalidatedAssumption, Uncertainty: domain.UncertaintyPOSSIBLE,
	}); err != nil {
		t.Fatal(err)
	}

	rep, err = svc.ImpactReport(ctx, d.ID)
	if err != nil {
		t.Fatalf("ImpactReport findings: %v", err)
	}
	if rep.OverallClass != domain.ImpactClassDestructive {
		t.Fatalf("OverallClass want DESTRUCTIVE got %q", rep.OverallClass)
	}
	if rep.OverallUncertainty != domain.UncertaintyUNKNOWN {
		t.Fatalf("OverallUncertainty want UNKNOWN got %q", rep.OverallUncertainty)
	}
	if !rep.HasUnknown || !rep.Incomplete {
		t.Fatalf("UNKNOWN present must HasUnknown+Incomplete: %+v", rep)
	}
	foundUnk := false
	for _, f := range rep.Findings {
		if f.ID == unk.ID && f.Uncertainty == domain.UncertaintyUNKNOWN {
			foundUnk = true
		}
	}
	if !foundUnk {
		t.Fatal("UNKNOWN finding must not be omitted from Findings")
	}

	// Alternatives without recommended → Incomplete.
	if _, err := svc.AddDecisionAlternative(ctx, d.ID, domain.AlternativeInput{Title: "A"}); err != nil {
		t.Fatal(err)
	}
	rep, err = svc.ImpactReport(ctx, d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !rep.Incomplete || len(rep.Alternatives) != 1 {
		t.Fatalf("alts without recommended: %+v", rep)
	}

	// REVERSAL beats DESTRUCTIVE; clear UNKNOWN path with known-only findings on fresh decision.
	d2, _ := svc.CreateDecision(ctx, domain.DecisionInput{Title: "rollup2"})
	_, _ = svc.AddImpactFinding(ctx, d2.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassDestructive, Kind: domain.FindingKindAffectedWork, Uncertainty: domain.UncertaintyKNOWN,
	})
	_, _ = svc.AddImpactFinding(ctx, d2.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassReversal, Kind: domain.FindingKindAffectedWork, Uncertainty: domain.UncertaintyLIKELY,
	})
	rep2, err := svc.ImpactReport(ctx, d2.ID)
	if err != nil {
		t.Fatal(err)
	}
	if rep2.OverallClass != domain.ImpactClassReversal {
		t.Fatalf("REVERSAL rollup: %q", rep2.OverallClass)
	}
	if rep2.OverallUncertainty != domain.UncertaintyLIKELY {
		t.Fatalf("OverallUncertainty LIKELY: %q", rep2.OverallUncertainty)
	}
	if rep2.HasUnknown {
		t.Fatalf("no UNKNOWN should HasUnknown=false: %+v", rep2)
	}
	// Still Incomplete: empty links + findings present is OK for HasUnknown, but
	// empty findings+no links is Incomplete; here we have findings so Incomplete only if alts.
	if rep2.Incomplete {
		t.Fatalf("known findings no alts should Incomplete=false: %+v", rep2)
	}
}

func TestImpactDoesNotAddEntityLinkRels(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	d, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "no new rels"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkDecisionTask(ctx, d.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AddImpactFinding(ctx, d.ID, domain.ImpactFindingInput{
		ImpactClass: domain.ImpactClassSAFE, Kind: domain.FindingKindAffectedWork,
		RelatedType: domain.EntityTask, RelatedID: task.ID,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AddDecisionAlternative(ctx, d.ID, domain.AlternativeInput{Title: "keep"}); err != nil {
		t.Fatal(err)
	}
	_, _ = svc.ImpactReport(ctx, d.ID)

	links, err := st.ListLinksFrom(domain.EntityDecision, d.ID)
	if err != nil || len(links) != 1 || links[0].Rel != domain.RelDecisionAffectsTask {
		t.Fatalf("only decision_affects_task expected: %+v err=%v", links, err)
	}
	byRel, err := st.ListLinksByRel(domain.RelDecisionAffectsTask)
	if err != nil || len(byRel) != 1 {
		t.Fatalf("decision_affects_task count: %+v err=%v", byRel, err)
	}
}
