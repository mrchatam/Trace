package domain_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestSeedExportSkipsDanglingScopeMemberFromReview(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()

	sc, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth", Title: "Auth", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Task in scope"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, task.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	rev, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "Review"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, rev.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{})
	if err != nil {
		t.Fatalf("BuildSeedDocument: %v", err)
	}

	var scopeMemberLinks []domain.SeedLink
	for _, l := range doc.Links {
		if l.Rel == domain.RelScopeMember {
			scopeMemberLinks = append(scopeMemberLinks, l)
		}
	}

	foundTaskLink := false
	foundReviewLink := false
	for _, l := range scopeMemberLinks {
		if l.From == task.ID && l.To == sc.ID {
			foundTaskLink = true
		}
		if l.From == rev.ID && l.To == sc.ID {
			foundReviewLink = true
		}
	}

	if !foundTaskLink {
		t.Fatal("export missing scope_member link from task to scope")
	}
	if foundReviewLink {
		t.Fatal("export must not include scope_member link from review (not exported entity)")
	}

	if len(doc.Scopes) != 1 || doc.Scopes[0].ID != sc.ID {
		t.Fatalf("export missing scope: %#v", doc.Scopes)
	}

	dir2 := t.TempDir()
	st2, err := store.Open(dir2)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	svc2 := domain.New(st2)

	summary, err := svc2.ImportSeedDocument(ctx, doc)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if !summary.OK {
		t.Fatalf("import not OK: %+v", summary)
	}

	gotScope, err := st2.GetScopeBySlug("auth")
	if err != nil {
		t.Fatalf("scope not imported: %v", err)
	}
	if gotScope.Title != "Auth" || gotScope.Kind != "feature" {
		t.Fatalf("scope mismatch: %+v", gotScope)
	}

	links, err := st2.ListLinksByRel(domain.RelScopeMember)
	if err != nil {
		t.Fatalf("ListLinksByRel: %v", err)
	}
	foundImportedTaskLink := false
	for _, l := range links {
		if l.FromID == task.ID && l.ToID == sc.ID {
			foundImportedTaskLink = true
		}
	}
	if !foundImportedTaskLink {
		t.Fatal("task scope_member link not round-tripped")
	}

	for _, l := range links {
		if l.FromID == rev.ID {
			t.Fatal("review scope_member link should not be imported")
		}
	}
}

func TestSeedExportSkipsDanglingScopeMemberFromChange(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()

	sc, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "api", Title: "API", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "API task"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, task.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Paths:  []domain.ChangePathInput{{Path: "internal/api.go"}},
		Reason: "Add endpoint",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, chg.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{})
	if err != nil {
		t.Fatalf("BuildSeedDocument: %v", err)
	}

	var scopeMemberLinks []domain.SeedLink
	for _, l := range doc.Links {
		if l.Rel == domain.RelScopeMember {
			scopeMemberLinks = append(scopeMemberLinks, l)
		}
	}

	foundTaskLink := false
	foundChangeLink := false
	for _, l := range scopeMemberLinks {
		if l.From == task.ID && l.To == sc.ID {
			foundTaskLink = true
		}
		if l.From == chg.ID && l.To == sc.ID {
			foundChangeLink = true
		}
	}

	if !foundTaskLink {
		t.Fatal("export missing scope_member link from task to scope")
	}
	if foundChangeLink {
		t.Fatal("export must not include scope_member link from change (not exported entity)")
	}

	dir2 := t.TempDir()
	st2, err := store.Open(dir2)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	svc2 := domain.New(st2)

	if _, err := svc2.ImportSeedDocument(ctx, doc); err != nil {
		t.Fatalf("import failed: %v", err)
	}

	links, err := st2.ListLinksByRel(domain.RelScopeMember)
	if err != nil {
		t.Fatalf("ListLinksByRel: %v", err)
	}
	foundImportedTaskLink := false
	for _, l := range links {
		if l.FromID == task.ID && l.ToID == sc.ID {
			foundImportedTaskLink = true
		}
	}
	if !foundImportedTaskLink {
		t.Fatal("task scope_member link not round-tripped")
	}

	for _, l := range links {
		if l.FromID == chg.ID {
			t.Fatal("change scope_member link should not be imported")
		}
	}
}

func TestSeedExportIncludesCausalLinks(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()

	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Use SQLite"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Implement storage"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkDecisionTask(ctx, dec.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	sc, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "storage", Title: "Storage", Kind: "layer",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, task.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{})
	if err != nil {
		t.Fatalf("BuildSeedDocument: %v", err)
	}

	hasDecisionAffectsTask := false
	hasScopeMember := false
	for _, l := range doc.Links {
		if l.Rel == domain.RelDecisionAffectsTask && l.From == dec.ID && l.To == task.ID {
			hasDecisionAffectsTask = true
		}
		if l.Rel == domain.RelScopeMember && l.From == task.ID && l.To == sc.ID {
			hasScopeMember = true
		}
	}

	if !hasDecisionAffectsTask {
		t.Fatal("export missing decision_affects_task causal link")
	}
	if !hasScopeMember {
		t.Fatal("export missing scope_member link from task to scope")
	}

	dir2 := t.TempDir()
	st2, err := store.Open(dir2)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	svc2 := domain.New(st2)

	if _, err := svc2.ImportSeedDocument(ctx, doc); err != nil {
		t.Fatalf("import failed: %v", err)
	}

	decLinks, err := st2.ListLinksByRel(domain.RelDecisionAffectsTask)
	if err != nil || len(decLinks) != 1 || decLinks[0].FromID != dec.ID || decLinks[0].ToID != task.ID {
		t.Fatalf("decision_affects_task not round-tripped: %+v err=%v", decLinks, err)
	}

	scopeLinks, err := st2.ListLinksByRel(domain.RelScopeMember)
	if err != nil || len(scopeLinks) != 1 || scopeLinks[0].FromID != task.ID || scopeLinks[0].ToID != sc.ID {
		t.Fatalf("scope_member not round-tripped: %+v err=%v", scopeLinks, err)
	}
}
