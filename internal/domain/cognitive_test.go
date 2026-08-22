package domain_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestCreateUncertaintyDefaultsOpenInfo(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	u, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{Title: "What is N?"})
	if err != nil {
		t.Fatalf("CreateUncertainty: %v", err)
	}
	if u.Status != store.UncertaintyStatusOpen || u.Severity != store.UncertaintySeverityINFO {
		t.Fatalf("defaults: %+v", u)
	}
	got, err := svc.GetUncertainty(ctx, u.ID)
	if err != nil || got.ID != u.ID {
		t.Fatalf("GetUncertainty: %+v err=%v", got, err)
	}
	evs, err := st.ListEventsByEntity(domain.EntityUncertainty, u.ID)
	if err != nil || len(evs) < 1 || evs[0].Type != domain.EventEntityCreated {
		t.Fatalf("entity.created: %+v err=%v", evs, err)
	}
}

func TestBlockingUncertaintyRequiresTaskID(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	_, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "blocked?",
		Severity: store.UncertaintySeverityBlocking,
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("want ErrValidation, got %v", err)
	}

	_, err = svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "blocked?",
		Severity: store.UncertaintySeverityBlocking,
		TaskID:   "missing-task",
	})
	if err == nil {
		t.Fatal("missing task must fail closed")
	}
}

func TestBlockingUncertaintyIncrementsCountForTask(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	taskA, err := svc.CreateTask(ctx, domain.TaskInput{Title: "A"})
	if err != nil {
		t.Fatal(err)
	}
	taskB, err := svc.CreateTask(ctx, domain.TaskInput{Title: "B"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "blocking A?",
		Severity: store.UncertaintySeverityBlocking,
		TaskID:   taskA.ID,
		GoalID:   g.ID,
	})
	if err != nil {
		t.Fatalf("CreateUncertainty: %v", err)
	}

	nA, err := svc.CountBlockingUncertainties(ctx, taskA.ID)
	if err != nil || nA != 1 {
		t.Fatalf("task A count=%d err=%v want 1", nA, err)
	}
	nB, err := svc.CountBlockingUncertainties(ctx, taskB.ID)
	if err != nil || nB != 0 {
		t.Fatalf("task B count=%d err=%v want 0 (goal link must not increment)", nB, err)
	}
}

func TestInfoUncertaintyDoesNotIncrementBlockingCount(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "info?",
		Severity: store.UncertaintySeverityINFO,
		TaskID:   task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	n, err := svc.CountBlockingUncertainties(ctx, task.ID)
	if err != nil || n != 0 {
		t.Fatalf("INFO must not count: n=%d err=%v", n, err)
	}
}

func TestResolveUncertaintyClearsBlockingCount(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t"})
	if err != nil {
		t.Fatal(err)
	}
	u, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "blocking?",
		Severity: store.UncertaintySeverityBlocking,
		TaskID:   task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	n, err := svc.CountBlockingUncertainties(ctx, task.ID)
	if err != nil || n != 1 {
		t.Fatalf("before resolve: n=%d err=%v", n, err)
	}

	out, err := svc.ResolveUncertainty(ctx, u.ID, "answered")
	if err != nil {
		t.Fatalf("ResolveUncertainty: %v", err)
	}
	if out.Status != store.UncertaintyStatusResolved || out.Resolution != "answered" {
		t.Fatalf("resolved: %+v", out)
	}
	n, err = svc.CountBlockingUncertainties(ctx, task.ID)
	if err != nil || n != 0 {
		t.Fatalf("after resolve: n=%d err=%v", n, err)
	}
	evs, err := st.ListEventsByEntity(domain.EntityUncertainty, u.ID)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, e := range evs {
		if e.Type == domain.EventUncertaintyResolved {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing uncertainty.resolved: %+v", evs)
	}
	if _, err := svc.ResolveUncertainty(ctx, u.ID, "again"); err == nil {
		t.Fatal("RESOLVED is terminal")
	}
}

func TestSupersedeUncertaintyClearsBlockingCount(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t"})
	if err != nil {
		t.Fatal(err)
	}
	u, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "blocking?",
		Severity: store.UncertaintySeverityBlocking,
		TaskID:   task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	out, err := svc.SupersedeUncertainty(ctx, u.ID, "replaced")
	if err != nil {
		t.Fatalf("SupersedeUncertainty: %v", err)
	}
	if out.Status != store.UncertaintyStatusSuperseded {
		t.Fatalf("superseded: %+v", out)
	}
	n, err := svc.CountBlockingUncertainties(ctx, task.ID)
	if err != nil || n != 0 {
		t.Fatalf("after supersede: n=%d err=%v", n, err)
	}
	evs, err := st.ListEventsByEntity(domain.EntityUncertainty, u.ID)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, e := range evs {
		if e.Type == domain.EventUncertaintySuperseded {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing uncertainty.superseded: %+v", evs)
	}
}

func TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "seed", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.CreateUncertainty(ctx, domain.UncertaintyInput{
		Title:    "blocking?",
		Severity: store.UncertaintySeverityBlocking,
		TaskID:   task.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	n, err := svc.CountBlockingUncertainties(ctx, task.ID)
	if err != nil || n != 1 {
		t.Fatalf("count=%d err=%v want 1", n, err)
	}

	inputs := deliberation.PolicyInputs{
		BlockingUncertaintyCount: n,
		PlanExists:               true,
		PlanCritiqued:            true,
	}
	next, ev, err := svc.ApplyDeliberationTransition(ctx, task.ID, g.ID, inputs)
	if err != nil {
		t.Fatalf("ApplyDeliberationTransition: %v", err)
	}
	if next.CurrentPhase != deliberation.PhaseInvestigate {
		t.Fatalf("phase %s want INVESTIGATE", next.CurrentPhase)
	}
	if next.CurrentPhase == deliberation.PhaseExecute {
		t.Fatal("EXECUTE forbidden while blocking uncertainty")
	}
	var payload deliberation.TransitionPayload
	if err := json.Unmarshal([]byte(ev.PayloadJSON), &payload); err != nil {
		t.Fatalf("payload: %v", err)
	}
	if payload.ReasonCode != deliberation.ReasonBlockingUncertainty {
		t.Fatalf("reason %s want blocking_uncertainty", payload.ReasonCode)
	}
	if payload.ToPhase == deliberation.PhaseExecute {
		t.Fatal("EXECUTE in payload")
	}
	if payload.PolicyInputs.BlockingUncertaintyCount != 1 {
		t.Fatalf("policy count=%d", payload.PolicyInputs.BlockingUncertaintyCount)
	}
}

func TestInvalidateAssumptionSetsStaleAndKeepsRow(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	a, err := svc.CreateAssumption(ctx, domain.AssumptionInput{Title: "CGO free"})
	if err != nil {
		t.Fatal(err)
	}
	out, disc, err := svc.InvalidateAssumption(ctx, a.ID, domain.InvalidateAssumptionInput{
		Status: store.StatusStale,
		Reason: "modernc still required",
	})
	if err != nil {
		t.Fatalf("InvalidateAssumption: %v", err)
	}
	if disc != nil {
		t.Fatal("no discovery expected")
	}
	if out.Status != store.StatusStale {
		t.Fatalf("status: %+v", out)
	}
	got, err := svc.GetAssumption(ctx, a.ID)
	if err != nil || got.Status != store.StatusStale || got.Title != "CGO free" {
		t.Fatalf("row kept: %+v err=%v", got, err)
	}
	evs, err := st.ListEventsByEntity(domain.EntityAssumption, a.ID)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, e := range evs {
		if e.Type == domain.EventAssumptionInvalidated {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing assumption.invalidated: %+v", evs)
	}

	if _, _, err := svc.InvalidateAssumption(ctx, a.ID, domain.InvalidateAssumptionInput{
		Status: store.StatusActive,
		Reason: "revive",
	}); err == nil {
		t.Fatal("revive to ACTIVE must fail")
	}
}

func TestInvalidateAssumptionSupersededNoDelete(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	a, err := svc.CreateAssumption(ctx, domain.AssumptionInput{Title: "old"})
	if err != nil {
		t.Fatal(err)
	}
	repl, err := svc.CreateAssumption(ctx, domain.AssumptionInput{Title: "new"})
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := svc.InvalidateAssumption(ctx, a.ID, domain.InvalidateAssumptionInput{
		Status:       store.StatusSuperseded,
		Reason:       "replaced",
		SupersededBy: repl.ID,
	})
	if err != nil {
		t.Fatalf("InvalidateAssumption: %v", err)
	}
	if out.Status != store.StatusSuperseded {
		t.Fatalf("status: %+v", out)
	}
	got, err := svc.GetAssumption(ctx, a.ID)
	if err != nil || got.ID != a.ID {
		t.Fatalf("must not delete: %+v err=%v", got, err)
	}
	if _, _, err := svc.InvalidateAssumption(ctx, a.ID, domain.InvalidateAssumptionInput{
		Status: store.StatusStale,
		Reason: "back",
	}); err == nil {
		t.Fatal("SUPERSEDED is terminal")
	}
}

func TestInvalidateAssumptionEmitsImpactFindingOnLinkedDecision(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "SQLite"})
	if err != nil {
		t.Fatal(err)
	}
	a, err := svc.CreateAssumption(ctx, domain.AssumptionInput{
		Title:       "local-first",
		DecisionIDs: []string{dec.ID},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = svc.InvalidateAssumption(ctx, a.ID, domain.InvalidateAssumptionInput{
		Status: store.StatusStale,
		Reason: "hosting considered",
	})
	if err != nil {
		t.Fatalf("InvalidateAssumption: %v", err)
	}

	findings, err := svc.ListImpactFindings(ctx, dec.ID)
	if err != nil || len(findings) != 1 {
		t.Fatalf("findings: %+v err=%v", findings, err)
	}
	f := findings[0]
	if f.Kind != domain.FindingKindInvalidatedAssumption || f.ImpactClass != domain.ImpactClassCaution {
		t.Fatalf("finding: %+v", f)
	}
	if f.Uncertainty != domain.UncertaintyUNKNOWN || f.RelatedType != domain.EntityAssumption || f.RelatedID != a.ID {
		t.Fatalf("finding related: %+v", f)
	}

	recs, err := st.ListDecisionReconsiderationsByDecisionID(dec.ID)
	if err != nil || len(recs) != 1 {
		t.Fatalf("reconsiderations: %+v err=%v", recs, err)
	}
	if recs[0].Trigger != store.ReconsiderTriggerInvalidatedAssumption || recs[0].Status != store.ReconsiderStatusFired {
		t.Fatalf("reconsider: %+v", recs[0])
	}

	gotDec, err := svc.GetDecision(ctx, dec.ID)
	if err != nil || gotDec.Status != store.StatusActive {
		t.Fatalf("decision not auto-stale: %+v err=%v", gotDec, err)
	}
}

func TestInvalidateAssumptionOptionalPlanAffectingDiscovery(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "t"})
	if err != nil {
		t.Fatal(err)
	}
	a, err := svc.CreateAssumption(ctx, domain.AssumptionInput{Title: "assume"})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = svc.InvalidateAssumption(ctx, a.ID, domain.InvalidateAssumptionInput{
		Status:         store.StatusStale,
		Reason:         "disproved",
		EmitDiscovery:  true,
		DiscoveryTitle: "",
	})
	if err == nil {
		t.Fatal("EmitDiscovery without title must fail")
	}

	out, disc, err := svc.InvalidateAssumption(ctx, a.ID, domain.InvalidateAssumptionInput{
		Status:         store.StatusStale,
		Reason:         "disproved",
		EmitDiscovery:  true,
		DiscoveryTitle: "assumption failed",
		DiscoveryBody:  "bench",
		TaskIDs:        []string{task.ID},
	})
	if err != nil {
		t.Fatalf("InvalidateAssumption: %v", err)
	}
	if disc == nil || disc.Severity != domain.SeverityPlanAffecting || disc.Title != "assumption failed" {
		t.Fatalf("discovery: %+v", disc)
	}
	if out.Status != store.StatusStale {
		t.Fatalf("assumption: %+v", out)
	}
	links, err := svc.ListLinksFrom(ctx, domain.EntityDiscovery, disc.ID)
	if err != nil {
		t.Fatal(err)
	}
	var sawInvalidate, sawMentions bool
	for _, l := range links {
		if l.Rel == domain.RelDiscoveryInvalidatesAssumption && l.ToID == a.ID {
			sawInvalidate = true
		}
		if l.Rel == domain.RelDiscoveryMentionsTask && l.ToID == task.ID {
			sawMentions = true
		}
	}
	if !sawInvalidate || !sawMentions {
		t.Fatalf("discovery links: %+v", links)
	}
}

func TestHypothesisLinksEvidenceWithoutDiscoveryTable(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	ev, err := svc.CreateEvidence(ctx, domain.EvidenceInput{Title: "bench.json"})
	if err != nil {
		t.Fatal(err)
	}
	unc, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{Title: "is sqlite enough?"})
	if err != nil {
		t.Fatal(err)
	}
	h, err := svc.CreateHypothesis(ctx, domain.HypothesisInput{
		Title:         "SQLite is enough",
		EvidenceIDs:   []string{ev.ID},
		UncertaintyID: unc.ID,
	})
	if err != nil {
		t.Fatalf("CreateHypothesis: %v", err)
	}
	if h.Status != store.HypothesisStatusOpen {
		t.Fatalf("status: %+v", h)
	}

	links, err := svc.ListLinksFrom(ctx, domain.EntityHypothesis, h.ID)
	if err != nil {
		t.Fatal(err)
	}
	var sawEvidence, sawUnc bool
	for _, l := range links {
		if l.Rel == domain.RelHypothesisSupportedBy && l.ToType == domain.EntityEvidence && l.ToID == ev.ID {
			sawEvidence = true
		}
		if l.Rel == domain.RelHypothesisAddressesUncertainty && l.ToID == unc.ID {
			sawUnc = true
		}
	}
	if !sawEvidence || !sawUnc {
		t.Fatalf("links: %+v", links)
	}

	discs, err := st.ListDiscoveries()
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range discs {
		if d.ID == h.ID || d.Title == h.Title {
			t.Fatalf("hypothesis must not be a discovery: %+v", d)
		}
	}
	if _, err := svc.GetHypothesis(ctx, h.ID); err != nil {
		t.Fatalf("hypothesis row missing: %v", err)
	}

	confirmed, err := svc.ConfirmHypothesis(ctx, h.ID, "bench passed")
	if err != nil || confirmed.Status != store.HypothesisStatusConfirmed {
		t.Fatalf("confirm: %+v err=%v", confirmed, err)
	}
}

func TestDecisionReconsiderPreservesDecisionAndAlternatives(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "SQLite"})
	if err != nil {
		t.Fatal(err)
	}
	alt, err := svc.AddDecisionAlternative(ctx, dec.ID, domain.AlternativeInput{
		Title: "Postgres", Body: "hosted",
	})
	if err != nil {
		t.Fatal(err)
	}

	row, err := svc.RecordDecisionReconsideration(ctx, dec.ID, domain.ReconsiderationInput{
		Trigger:     store.ReconsiderTriggerNewEvidence,
		Reason:      "new latency numbers",
		RelatedType: domain.EntityEvidence,
		RelatedID:   "e1",
	})
	if err != nil {
		t.Fatalf("RecordDecisionReconsideration: %v", err)
	}
	if row.Status != store.ReconsiderStatusFired || row.DecisionID != dec.ID {
		t.Fatalf("row: %+v", row)
	}

	got, err := svc.GetDecision(ctx, dec.ID)
	if err != nil || got.Title != "SQLite" {
		t.Fatalf("decision preserved: %+v err=%v", got, err)
	}
	alts, err := svc.ListDecisionAlternatives(ctx, dec.ID)
	if err != nil || len(alts) != 1 || alts[0].ID != alt.ID {
		t.Fatalf("alternatives preserved: %+v err=%v", alts, err)
	}
	listed, err := st.ListDecisionReconsiderationsByDecisionID(dec.ID)
	if err != nil || len(listed) != 1 {
		t.Fatalf("reconsiderations: %+v err=%v", listed, err)
	}
	evs, err := st.ListEventsByEntity(domain.EntityDecision, dec.ID)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, e := range evs {
		if e.Type == domain.EventDecisionReconsider {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing decision.reconsider: %+v", evs)
	}
}

func TestUnknownUncertaintySeverityFailClosed(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	for _, sev := range []string{"PLAN_AFFECTING", "WEIRD", "blocking"} {
		_, err := svc.CreateUncertainty(ctx, domain.UncertaintyInput{
			Title:    "q?",
			Severity: sev,
		})
		var ve *domain.ErrValidation
		if !errors.As(err, &ve) {
			t.Fatalf("severity %q: want ErrValidation, got %v", sev, err)
		}
	}
}
