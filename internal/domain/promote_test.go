package domain_test

import (
	"context"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestPromoteBlockingDiscoveryCreatesAndLinksTask(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	goal, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	disc, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title:    "Blocking gap",
		Severity: domain.SeverityBlocking,
	})
	if err != nil {
		t.Fatal(err)
	}

	taskID, inserted, err := svc.PromoteBlockingDiscovery(ctx, disc.ID, goal.ID)
	if err != nil {
		t.Fatalf("PromoteBlockingDiscovery: %v", err)
	}
	if !inserted {
		t.Fatal("expected inserted=true on first promote")
	}
	if taskID != disc.ID {
		t.Fatalf("task id must match discovery id: got %q want %q", taskID, disc.ID)
	}

	task, err := st.GetTask(taskID)
	if err != nil {
		t.Fatal(err)
	}
	if task.GoalID == nil || *task.GoalID != goal.ID {
		t.Fatalf("task goal mismatch: %+v", task)
	}
	links, err := st.ListLinksFrom(domain.EntityDiscovery, disc.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 1 || links[0].Rel != domain.RelDiscoveryMentionsTask || links[0].ToID != taskID {
		t.Fatalf("discovery links mismatch: %+v", links)
	}
}

func TestPromoteBlockingDiscoveryIsIdempotent(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	goal, _ := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	disc, _ := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "Blocking gap", Severity: domain.SeverityBlocking})

	taskID1, inserted1, err := svc.PromoteBlockingDiscovery(ctx, disc.ID, goal.ID)
	if err != nil {
		t.Fatal(err)
	}
	taskID2, inserted2, err := svc.PromoteBlockingDiscovery(ctx, disc.ID, goal.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !inserted1 || inserted2 {
		t.Fatalf("inserted flags mismatch: first=%v second=%v", inserted1, inserted2)
	}
	if taskID1 != taskID2 {
		t.Fatalf("idempotent task id mismatch: %q vs %q", taskID1, taskID2)
	}
}

func TestPromoteBlockingDiscoveryFailsClosed(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	goal, _ := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	infoDisc, _ := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "Info gap", Severity: domain.SeverityINFO})

	if _, _, err := svc.PromoteBlockingDiscovery(ctx, "", goal.ID); err == nil {
		t.Fatal("expected discovery_id required error")
	}
	if _, _, err := svc.PromoteBlockingDiscovery(ctx, infoDisc.ID, goal.ID); err == nil || !strings.Contains(err.Error(), "BLOCKING") {
		t.Fatalf("expected BLOCKING validation error, got %v", err)
	}
	blockingDisc, _ := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "Blocking", Severity: domain.SeverityBlocking})
	if _, _, err := svc.PromoteBlockingDiscovery(ctx, blockingDisc.ID, ""); err == nil || !strings.Contains(err.Error(), "goal_id is required") {
		t.Fatalf("expected goal_id validation error, got %v", err)
	}
	if _, _, err := svc.PromoteBlockingDiscovery(ctx, "00000000-0000-4000-8000-000000000000", goal.ID); err == nil {
		t.Fatal("expected missing discovery error")
	}
}

func TestPromoteBlockingDiscoveryPreservesSeedWorkState(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	goal, _ := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	disc, _ := svc.CreateDiscovery(ctx, domain.DiscoveryInput{Title: "Blocking", Severity: domain.SeverityBlocking})

	_, _, err := svc.ImportSeedTask(ctx, domain.SeedTask{
		ID:     disc.ID,
		GoalID: goal.ID,
		Title:  "Preseeded",
	}, &goal.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, disc.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "test", Reason: "preserve",
	}); err != nil {
		t.Fatal(err)
	}

	taskID, inserted, err := svc.PromoteBlockingDiscovery(ctx, disc.ID, goal.ID)
	if err != nil {
		t.Fatal(err)
	}
	if inserted {
		t.Fatal("expected no insert when seed task already exists")
	}
	task, err := st.GetTask(taskID)
	if err != nil {
		t.Fatal(err)
	}
	if task.WorkState != store.WorkStateInProgress {
		t.Fatalf("work_state should be preserved, got %q", task.WorkState)
	}
}

func TestImportSeedDocumentSurfacesPromotionCandidates(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	goalID := "a0000000-0000-4000-8000-0000000000a1"
	taskID := "a0000000-0000-4000-8000-0000000000a2"
	discID := "a0000000-0000-4000-8000-0000000000a3"

	summary, err := svc.ImportSeedDocument(ctx, domain.SeedDocument{
		Version: 1,
		Goals:   []domain.SeedEntity{{ID: goalID, Title: "Seed goal"}},
		Tasks:   []domain.SeedTask{{ID: taskID, GoalID: goalID, Title: "Seed task"}},
	})
	if err != nil {
		t.Fatalf("seed import roster: %v", err)
	}
	if len(summary.PromotionCandidates) != 0 {
		t.Fatalf("expected no candidates before BLOCKING discovery, got %+v", summary.PromotionCandidates)
	}

	if _, _, err := svc.ImportSeedDiscovery(ctx, domain.SeedEntity{
		ID:    discID,
		Title: "Orphan blocking after import",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetDiscoverySeverity(ctx, discID, domain.SeverityBlocking); err != nil {
		t.Fatal(err)
	}

	summary2, err := svc.ImportSeedDocument(ctx, domain.SeedDocument{
		Version: 1,
		Goals:   []domain.SeedEntity{{ID: goalID, Title: "Seed goal"}},
		Tasks:   []domain.SeedTask{{ID: taskID, GoalID: goalID, Title: "Seed task"}},
	})
	if err != nil {
		t.Fatalf("re-import: %v", err)
	}
	if len(summary2.PromotionCandidates) != 1 {
		t.Fatalf("promotion_candidates len=%d want 1: %+v", len(summary2.PromotionCandidates), summary2.PromotionCandidates)
	}
	if summary2.PromotionCandidates[0].DiscoveryID != discID {
		t.Fatalf("candidate id=%q want %q", summary2.PromotionCandidates[0].DiscoveryID, discID)
	}
	if summary2.PromotionHint == "" || !strings.Contains(summary2.PromotionHint, "--from-discovery") {
		t.Fatalf("promotion_hint missing guided path: %q", summary2.PromotionHint)
	}

	tasks, err := st.ListTasksByGoalID(goalID)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 {
		t.Fatalf("import must not auto-spawn tasks, got %d", len(tasks))
	}

	promotedID, inserted, err := svc.PromoteBlockingDiscovery(ctx, summary2.PromotionCandidates[0].DiscoveryID, goalID)
	if err != nil {
		t.Fatal(err)
	}
	if !inserted || promotedID != discID {
		t.Fatalf("promote mismatch: id=%q inserted=%v want id=%q", promotedID, inserted, discID)
	}

	summary3, err := svc.ImportSeedDocument(ctx, domain.SeedDocument{Version: 1})
	if err != nil {
		t.Fatalf("post-promote import: %v", err)
	}
	if len(summary3.PromotionCandidates) != 0 {
		t.Fatalf("after promote expect empty candidates, got %+v", summary3.PromotionCandidates)
	}
}

func TestPromoteBlockingDiscoveryAfterImport(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	goal, _ := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	importedID := "f0000000-0000-4000-8000-0000000000f1"

	if _, _, err := svc.ImportSeedDiscovery(ctx, domain.SeedEntity{
		ID:    importedID,
		Title: "Imported blocking discovery",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetDiscoverySeverity(ctx, importedID, domain.SeverityBlocking); err != nil {
		t.Fatal(err)
	}
	if tasks, err := st.ListTasksByGoalID(goal.ID); err != nil {
		t.Fatal(err)
	} else if len(tasks) != 0 {
		t.Fatalf("import must not auto-spawn tasks, got %d", len(tasks))
	}

	taskID, inserted, err := svc.PromoteBlockingDiscovery(ctx, importedID, goal.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !inserted || taskID != importedID {
		t.Fatalf("promotion mismatch: taskID=%q inserted=%v", taskID, inserted)
	}
	taskID2, inserted2, err := svc.PromoteBlockingDiscovery(ctx, importedID, goal.ID)
	if err != nil {
		t.Fatal(err)
	}
	if inserted2 || taskID2 != taskID {
		t.Fatalf("idempotent mismatch after import: id2=%q inserted2=%v", taskID2, inserted2)
	}
}
