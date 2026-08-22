package domain_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestSeedImportP20RoundTrip(t *testing.T) {
	_, st := openDomain(t)
	seedP20Cognition(t, st)

	doc1, err := domain.BuildSeedDocument(context.Background(), st, domain.ExportOpts{})
	if err != nil {
		t.Fatalf("export1: %v", err)
	}

	svc2, st2 := openDomain(t)
	if _, err := svc2.ImportSeedDocument(context.Background(), doc1); err != nil {
		t.Fatalf("import1: %v", err)
	}

	doc2, err := domain.BuildSeedDocument(context.Background(), st2, domain.ExportOpts{})
	if err != nil {
		t.Fatalf("export2: %v", err)
	}

	svc3, st3 := openDomain(t)
	if _, err := svc3.ImportSeedDocument(context.Background(), doc2); err != nil {
		t.Fatalf("import2: %v", err)
	}

	assertP20IDs(t, st3)

	ds, err := st3.GetDeliberationState(seedTaskID)
	if err != nil {
		t.Fatalf("GetDeliberationState: %v", err)
	}
	if ds.CurrentPhase != string(deliberation.PhaseInvestigate) {
		t.Fatalf("current_phase: got %q want %q", ds.CurrentPhase, deliberation.PhaseInvestigate)
	}
	if ds.HopCount != 2 || !ds.PlanCritiqued {
		t.Fatalf("deliberation fields: %+v", ds)
	}

	// Idempotent re-import must not duplicate rows.
	before, err := countP20Rows(st3)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc3.ImportSeedDocument(context.Background(), doc2); err != nil {
		t.Fatalf("re-import: %v", err)
	}
	after, err := countP20Rows(st3)
	if err != nil {
		t.Fatal(err)
	}
	if before != after {
		t.Fatalf("re-import duplicated rows: before=%+v after=%+v", before, after)
	}
}

func assertP20IDs(t *testing.T, st *store.Store) {
	t.Helper()
	checks := []struct {
		name string
		fn   func() error
	}{
		{"baseline", func() error { _, err := st.GetBaseline(seedBaseline); return err }},
		{"outcome", func() error { _, err := st.GetOutcomeResult(seedOutcome); return err }},
		{"change", func() error { _, err := st.GetChange(seedChange); return err }},
		{"effect", func() error { _, err := st.GetEffect(seedEffect); return err }},
		{"uncertainty", func() error { _, err := st.GetUncertainty(seedUncert); return err }},
		{"hypothesis", func() error { _, err := st.GetHypothesis(seedHypo); return err }},
		{"reconsideration", func() error { _, err := st.GetDecisionReconsideration(seedRecon); return err }},
		{"regression", func() error { _, err := st.GetRegression(seedRegr); return err }},
		{"reflection", func() error { _, err := st.GetReflection(seedReflect); return err }},
		{"deliberation", func() error { _, err := st.GetDeliberationState(seedTaskID); return err }},
	}
	for _, c := range checks {
		if err := c.fn(); err != nil {
			t.Fatalf("%s id missing: %v", c.name, err)
		}
	}
	paths, err := st.ListChangePaths(seedChange)
	if err != nil || len(paths) != 1 || paths[0].Path != "internal/foo.go" {
		t.Fatalf("change paths: %+v err=%v", paths, err)
	}
}

type p20RowCounts struct {
	baselines, outcomes, changes, effects, uncertainties, hypotheses,
	reconsiderations, regressions, reflections, deliberation int
}

func countP20Rows(st *store.Store) (p20RowCounts, error) {
	var c p20RowCounts
	b, err := st.ListAllBaselines()
	if err != nil {
		return c, err
	}
	c.baselines = len(b)
	o, err := st.ListAllOutcomeResults()
	if err != nil {
		return c, err
	}
	c.outcomes = len(o)
	ch, err := st.ListAllChanges()
	if err != nil {
		return c, err
	}
	c.changes = len(ch)
	e, err := st.ListAllEffects()
	if err != nil {
		return c, err
	}
	c.effects = len(e)
	u, err := st.ListAllUncertainties()
	if err != nil {
		return c, err
	}
	c.uncertainties = len(u)
	h, err := st.ListAllHypotheses()
	if err != nil {
		return c, err
	}
	c.hypotheses = len(h)
	r, err := st.ListAllDecisionReconsiderations()
	if err != nil {
		return c, err
	}
	c.reconsiderations = len(r)
	reg, err := st.ListAllRegressions()
	if err != nil {
		return c, err
	}
	c.regressions = len(reg)
	ref, err := st.ListAllReflections()
	if err != nil {
		return c, err
	}
	c.reflections = len(ref)
	d, err := st.ListAllDeliberationStates()
	if err != nil {
		return c, err
	}
	c.deliberation = len(d)
	return c, nil
}

func TestSeedImportP20RoundTripJSON(t *testing.T) {
	// Ensure JSON marshal/unmarshal round-trip preserves P20 fields.
	_, st := openDomain(t)
	seedP20Cognition(t, st)
	doc, err := domain.BuildSeedDocument(context.Background(), st, domain.ExportOpts{})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var decoded domain.SeedDocument
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.DeliberationStates) != 1 || len(decoded.Changes) != 1 {
		t.Fatalf("decoded P20 slices: ds=%d changes=%d", len(decoded.DeliberationStates), len(decoded.Changes))
	}
}
