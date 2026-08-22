package domain_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestRecordImprovementQueryable(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	chg := mustRecordedChange(t, svc, task.ID)
	ev := mustEvidence(t, svc, "bench run")

	row, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
		ChangeID:    chg.ID,
		TaskID:      task.ID,
		Dimension:   "latency",
		Summary:     "p99 dropped after cache layer",
		EvidenceIDs: []string{ev.ID},
		Confidence:  0.9,
	})
	if err != nil {
		t.Fatalf("RecordImprovement: %v", err)
	}
	if row.ChangeID != chg.ID || row.TaskID != task.ID {
		t.Fatalf("row: %+v", row)
	}

	byChange, err := svc.ListImprovementsByChangeID(ctx, chg.ID)
	if err != nil || len(byChange) != 1 || byChange[0].ID != row.ID {
		t.Fatalf("ListImprovementsByChangeID: %+v err=%v", byChange, err)
	}
	byTask, err := svc.ListImprovementsByTaskID(ctx, task.ID)
	if err != nil || len(byTask) != 1 || byTask[0].ID != row.ID {
		t.Fatalf("ListImprovementsByTaskID: %+v err=%v", byTask, err)
	}
	got, err := svc.GetImprovement(ctx, row.ID)
	if err != nil || got.Summary != row.Summary {
		t.Fatalf("GetImprovement: %+v err=%v", got, err)
	}
}

func TestRecordImprovementFailClosedEmptySummary(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	_, task := mustGoalTask(t, svc)
	chg := mustRecordedChange(t, svc, task.ID)

	_, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
		ChangeID: chg.ID,
		TaskID:   task.ID,
		Summary:  "   ",
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("want ErrValidation, got %v", err)
	}
}

func TestSeedExportIncludesImprovements(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	seedP20Cognition(t, st)

	imID := "dddddddd-dddd-dddd-dddd-dddddddddd01"
	if _, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
		ID:        imID,
		ChangeID:  seedChange,
		TaskID:    seedTaskID,
		Dimension: "latency",
		Summary:   "seed improvement",
	}); err != nil {
		t.Fatalf("RecordImprovement: %v", err)
	}

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{})
	if err != nil {
		t.Fatalf("BuildSeedDocument: %v", err)
	}
	if len(doc.Improvements) != 1 || doc.Improvements[0].ID != imID {
		t.Fatalf("export improvements: %+v", doc.Improvements)
	}

	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatal(err)
	}
	if _, ok := parsed["improvements"]; !ok {
		t.Fatal("export missing improvements key")
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
	got, err := svc2.GetImprovement(ctx, imID)
	if err != nil || got.Summary != "seed improvement" {
		t.Fatalf("round-trip: %+v err=%v", got, err)
	}
}
