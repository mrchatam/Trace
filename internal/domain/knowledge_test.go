package domain_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestUpsertEngineeringKnowledge(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	ev := mustEvidence(t, svc, "load test run")

	row, err := svc.UpsertEngineeringKnowledge(ctx, domain.EngineeringKnowledgeInput{
		Title: "Cache tuning",
		Body: domain.KnowledgeBody{
			Summary:          "Use bounded LRU for hot paths",
			SourceEntityType: "manual",
			SourceEntityID:   "note-1",
		},
		Topic:       "improvement",
		EvidenceIDs: []string{ev.ID},
		Confidence:  0.85,
		SourceType:  "USER_ASSERTED",
	})
	if err != nil {
		t.Fatalf("UpsertEngineeringKnowledge: %v", err)
	}
	if row.Topic != "improvement" || row.Status != store.KnowledgeStatusActive {
		t.Fatalf("row: %+v", row)
	}

	got, err := svc.GetEngineeringKnowledge(ctx, row.ID)
	if err != nil || got.Title != row.Title {
		t.Fatalf("GetEngineeringKnowledge: %+v err=%v", got, err)
	}

	listed, err := svc.ListEngineeringKnowledge(ctx, domain.ListEngineeringKnowledgeOpts{Topic: "improvement"})
	if err != nil || len(listed) != 1 || listed[0].ID != row.ID {
		t.Fatalf("ListEngineeringKnowledge: %+v err=%v", listed, err)
	}

	_, err = svc.UpsertEngineeringKnowledge(ctx, domain.EngineeringKnowledgeInput{
		Title: "too big",
		Body: domain.KnowledgeBody{
			Summary:          strings.Repeat("x", 9000),
			SourceEntityType: "manual",
			SourceEntityID:   "note-2",
		},
		Topic: "improvement",
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("want body size validation, got %v", err)
	}

	_, err = svc.UpsertEngineeringKnowledge(ctx, domain.EngineeringKnowledgeInput{
		Title: "bad evidence",
		Body: domain.KnowledgeBody{
			Summary:          "x",
			SourceEntityType: "manual",
			SourceEntityID:   "note-3",
		},
		Topic:       "improvement",
		EvidenceIDs: []string{"00000000-0000-0000-0000-000000000000"},
	})
	if err == nil {
		t.Fatal("expected missing evidence error")
	}
}

func TestSynthesizeKnowledgeFromPatterns(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	for i := 0; i < 2; i++ {
		chg, err := svc.CreateChange(ctx, domain.ChangeInput{
			TaskID: task.ID,
			Reason: "internal improvement",
			Paths:  []domain.ChangePathInput{{Path: "internal/a.go", Status: "M"}},
		})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
			ChangeID: chg.ID,
			TaskID:   task.ID,
			Summary:  "faster path",
		}); err != nil {
			t.Fatal(err)
		}
	}

	first, err := svc.SynthesizeKnowledge(ctx)
	if err != nil {
		t.Fatalf("SynthesizeKnowledge: %v", err)
	}
	if first.Created < 1 {
		t.Fatalf("expected created rows: %+v", first)
	}

	rows, err := svc.ListEngineeringKnowledge(ctx, domain.ListEngineeringKnowledgeOpts{Topic: "pattern"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 {
		t.Fatal("expected pattern knowledge row")
	}

	time.Sleep(1100 * time.Millisecond)
	second, err := svc.SynthesizeKnowledge(ctx)
	if err != nil {
		t.Fatalf("second SynthesizeKnowledge: %v", err)
	}
	if second.Updated < 1 || second.Created != 0 {
		t.Fatalf("second run should update not duplicate: %+v", second)
	}
	rowsAfter, err := svc.ListEngineeringKnowledge(ctx, domain.ListEngineeringKnowledgeOpts{Topic: "pattern", Limit: 64})
	if err != nil {
		t.Fatal(err)
	}
	if len(rowsAfter) != len(rows) {
		t.Fatalf("duplicate rows: before=%d after=%d", len(rows), len(rowsAfter))
	}
	if rowsAfter[0].UpdatedAt <= rows[0].UpdatedAt {
		t.Fatalf("updated_at should advance: before=%q after=%q", rows[0].UpdatedAt, rowsAfter[0].UpdatedAt)
	}
}

func TestKnowledgeLinksDecision(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	seedP20Cognition(t, st)

	result, err := svc.SynthesizeKnowledge(ctx)
	if err != nil {
		t.Fatalf("SynthesizeKnowledge: %v", err)
	}
	if result.Created < 1 {
		t.Fatalf("expected knowledge rows: %+v", result)
	}

	rows, err := svc.ListEngineeringKnowledge(ctx, domain.ListEngineeringKnowledgeOpts{Topic: "decision"})
	if err != nil || len(rows) == 0 {
		t.Fatalf("decision knowledge: %+v err=%v", rows, err)
	}

	var body domain.KnowledgeBody
	if err := json.Unmarshal([]byte(rows[0].BodyJSON), &body); err != nil {
		t.Fatalf("body json: %v", err)
	}
	if body.DecisionID != seedDecID {
		t.Fatalf("decision_id=%q want %q body=%+v", body.DecisionID, seedDecID, body)
	}
}

func TestSeedExportIncludesKnowledge(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	seedP20Cognition(t, st)

	if _, err := svc.RefreshChangePatterns(ctx); err != nil {
		t.Fatalf("RefreshChangePatterns: %v", err)
	}
	knID := "eeeeeeee-eeee-eeee-eeee-eeeeeeeeee01"
	if _, err := svc.UpsertEngineeringKnowledge(ctx, domain.EngineeringKnowledgeInput{
		ID:    knID,
		Title: "seed knowledge",
		Body: domain.KnowledgeBody{
			Summary:          "portable row",
			SourceEntityType: "manual",
			SourceEntityID:   "seed-note",
		},
		Topic: "improvement",
	}); err != nil {
		t.Fatalf("UpsertEngineeringKnowledge: %v", err)
	}

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{})
	if err != nil {
		t.Fatalf("BuildSeedDocument: %v", err)
	}
	if len(doc.EngineeringKnowledge) != 1 || doc.EngineeringKnowledge[0].ID != knID {
		t.Fatalf("export knowledge: %+v", doc.EngineeringKnowledge)
	}
	if len(doc.ChangePatterns) == 0 {
		t.Fatalf("export change_patterns empty")
	}

	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"change_patterns", "engineering_knowledge"} {
		if _, ok := parsed[key]; !ok {
			t.Fatalf("export missing %s key", key)
		}
	}

	dir2 := t.TempDir()
	st2, err := store.Open(dir2)
	if err != nil {
		t.Fatal(err)
	}
	svc2 := domain.New(st2)
	if _, err := svc2.ImportSeedDocument(ctx, doc); err != nil {
		t.Fatalf("ImportSeedDocument: %v", err)
	}
	got, err := svc2.GetEngineeringKnowledge(ctx, knID)
	if err != nil || got.Title != "seed knowledge" {
		t.Fatalf("round-trip knowledge: %+v err=%v", got, err)
	}
	patterns, err := st2.ListAllChangePatterns()
	if err != nil || len(patterns) == 0 {
		t.Fatalf("round-trip patterns: %+v err=%v", patterns, err)
	}
}
