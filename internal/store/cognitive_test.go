package store

import (
	"testing"
)

func TestCreateUncertaintyStoreRoundtrip(t *testing.T) {
	s, _ := openTempStore(t)

	u, err := s.UpsertUncertainty(Uncertainty{
		Title:      "What is the hop budget?",
		Body:       "N=?",
		SourceType: "USER_ASSERTED",
	})
	if err != nil {
		t.Fatalf("UpsertUncertainty: %v", err)
	}
	if u.ID == "" || u.CreatedAt == "" || u.UpdatedAt == "" {
		t.Fatalf("expected ids/timestamps: %+v", u)
	}
	if u.Status != UncertaintyStatusOpen || u.Severity != UncertaintySeverityINFO {
		t.Fatalf("defaults: %+v", u)
	}
	if u.Kind != "" || u.Resolution != "" {
		t.Fatalf("empty kind/resolution: %+v", u)
	}

	got, err := s.GetUncertainty(u.ID)
	if err != nil {
		t.Fatalf("GetUncertainty: %v", err)
	}
	if got.Title != "What is the hop budget?" || got.Body != "N=?" {
		t.Fatalf("roundtrip: %+v", got)
	}

	got.Status = UncertaintyStatusResolved
	got.Resolution = "N=12"
	got.Severity = UncertaintySeverityBlocking
	got.Kind = UncertaintyKindGap
	updated, err := s.UpsertUncertainty(got)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Status != UncertaintyStatusResolved || updated.Resolution != "N=12" {
		t.Fatalf("updated: %+v", updated)
	}
	if updated.CreatedAt != u.CreatedAt {
		t.Fatalf("created_at changed: %q vs %q", updated.CreatedAt, u.CreatedAt)
	}

	if _, err := s.UpsertUncertainty(Uncertainty{Title: ""}); err == nil {
		t.Fatal("expected title required")
	}
	if _, err := s.GetUncertainty(""); err == nil {
		t.Fatal("expected id required")
	}
	if _, err := s.GetUncertainty("missing"); err == nil {
		t.Fatal("expected missing uncertainty")
	}
}

func TestHypothesisStoreRoundtrip(t *testing.T) {
	s, _ := openTempStore(t)

	h, err := s.UpsertHypothesis(Hypothesis{Title: "SQLite is enough"})
	if err != nil {
		t.Fatalf("UpsertHypothesis: %v", err)
	}
	if h.Status != HypothesisStatusOpen {
		t.Fatalf("default status: %+v", h)
	}
	got, err := s.GetHypothesis(h.ID)
	if err != nil || got.Title != "SQLite is enough" {
		t.Fatalf("GetHypothesis: %+v err=%v", got, err)
	}
	got.Status = HypothesisStatusConfirmed
	updated, err := s.UpsertHypothesis(got)
	if err != nil || updated.Status != HypothesisStatusConfirmed {
		t.Fatalf("confirm: %+v err=%v", updated, err)
	}
	if _, err := s.UpsertHypothesis(Hypothesis{}); err == nil {
		t.Fatal("expected title required")
	}
}

func TestBlockingCountSQL(t *testing.T) {
	s, _ := openTempStore(t)

	if _, err := s.CountOpenBlockingUncertaintiesByTaskID(""); err == nil {
		t.Fatal("empty taskID must fail closed")
	}

	taskA, err := s.UpsertTask(Task{Title: "A"})
	if err != nil {
		t.Fatal(err)
	}
	taskB, err := s.UpsertTask(Task{Title: "B"})
	if err != nil {
		t.Fatal(err)
	}

	n, err := s.CountOpenBlockingUncertaintiesByTaskID(taskA.ID)
	if err != nil || n != 0 {
		t.Fatalf("empty count: n=%d err=%v", n, err)
	}

	blocking, err := s.UpsertUncertainty(Uncertainty{
		Title:    "blocking?",
		Severity: UncertaintySeverityBlocking,
		Status:   UncertaintyStatusOpen,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.InsertLink(EntityLink{
		FromType: "uncertainty", FromID: blocking.ID,
		Rel: "uncertainty_blocks_task", ToType: "task", ToID: taskA.ID,
	}); err != nil {
		t.Fatal(err)
	}

	info, err := s.UpsertUncertainty(Uncertainty{
		Title:    "info?",
		Severity: UncertaintySeverityINFO,
		Status:   UncertaintyStatusOpen,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.InsertLink(EntityLink{
		FromType: "uncertainty", FromID: info.ID,
		Rel: "uncertainty_blocks_task", ToType: "task", ToID: taskA.ID,
	}); err != nil {
		t.Fatal(err)
	}

	resolved, err := s.UpsertUncertainty(Uncertainty{
		Title:    "resolved blocking?",
		Severity: UncertaintySeverityBlocking,
		Status:   UncertaintyStatusResolved,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.InsertLink(EntityLink{
		FromType: "uncertainty", FromID: resolved.ID,
		Rel: "uncertainty_blocks_task", ToType: "task", ToID: taskA.ID,
	}); err != nil {
		t.Fatal(err)
	}

	nA, err := s.CountOpenBlockingUncertaintiesByTaskID(taskA.ID)
	if err != nil || nA != 1 {
		t.Fatalf("task A count: n=%d err=%v want 1", nA, err)
	}
	nB, err := s.CountOpenBlockingUncertaintiesByTaskID(taskB.ID)
	if err != nil || nB != 0 {
		t.Fatalf("task B count: n=%d err=%v want 0", nB, err)
	}

	blocking.Status = UncertaintyStatusSuperseded
	if _, err := s.UpsertUncertainty(blocking); err != nil {
		t.Fatal(err)
	}
	nA, err = s.CountOpenBlockingUncertaintiesByTaskID(taskA.ID)
	if err != nil || nA != 0 {
		t.Fatalf("after supersede: n=%d err=%v", nA, err)
	}
}

func TestDecisionReconsiderStoreRoundtrip(t *testing.T) {
	s, _ := openTempStore(t)

	dec, err := s.UpsertDecision(Decision{Title: "Use SQLite"})
	if err != nil {
		t.Fatal(err)
	}
	alt, err := s.InsertDecisionAlternative(DecisionAlternative{
		DecisionID: dec.ID, Title: "Postgres",
	})
	if err != nil {
		t.Fatal(err)
	}

	row, err := s.InsertDecisionReconsideration(DecisionReconsideration{
		DecisionID:  dec.ID,
		Trigger:     ReconsiderTriggerNewEvidence,
		Reason:      "new bench",
		RelatedType: "evidence",
		RelatedID:   "e1",
	})
	if err != nil {
		t.Fatalf("InsertDecisionReconsideration: %v", err)
	}
	if row.Status != ReconsiderStatusFired || row.ReconsiderAt == "" || row.CreatedAt == "" {
		t.Fatalf("defaults: %+v", row)
	}

	listed, err := s.ListDecisionReconsiderationsByDecisionID(dec.ID)
	if err != nil || len(listed) != 1 || listed[0].ID != row.ID {
		t.Fatalf("list: %+v err=%v", listed, err)
	}

	gotDec, err := s.GetDecision(dec.ID)
	if err != nil || gotDec.Title != "Use SQLite" {
		t.Fatalf("decision preserved: %+v err=%v", gotDec, err)
	}
	alts, err := s.ListDecisionAlternativesByDecisionID(dec.ID)
	if err != nil || len(alts) != 1 || alts[0].ID != alt.ID {
		t.Fatalf("alternatives preserved: %+v err=%v", alts, err)
	}

	if _, err := s.InsertDecisionReconsideration(DecisionReconsideration{Trigger: ReconsiderTriggerNewEvidence}); err == nil {
		t.Fatal("expected decision_id required")
	}
}
