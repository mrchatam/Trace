package domain_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestSeedExportImportScopesRoundTrip(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	t1, err := svc.CreateTask(ctx, domain.TaskInput{Title: "FE", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	t2, err := svc.CreateTask(ctx, domain.TaskInput{Title: "BE", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	sc, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth", Title: "Auth", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, t1.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkAPIContract(ctx, t1.ID, t2.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkBlocks(ctx, t1.ID, t2.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkImplements(ctx, t1.ID, t2.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Scopes) != 1 || doc.Scopes[0].Slug != "auth" {
		t.Fatalf("scopes export: %#v", doc.Scopes)
	}
	wantRels := map[string]bool{
		domain.RelScopeMember: false,
		domain.RelAPIContract: false,
		domain.RelImplements:  false,
		domain.RelBlocks:      false,
	}
	for _, l := range doc.Links {
		if _, ok := wantRels[l.Rel]; ok {
			wantRels[l.Rel] = true
		}
	}
	for rel, ok := range wantRels {
		if !ok {
			t.Fatalf("export missing rel %s", rel)
		}
	}

	dir2 := t.TempDir()
	st2, err := store.Open(dir2)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	svc2 := domain.New(st2)
	if _, err := svc2.ImportSeedDocument(ctx, doc); err != nil {
		t.Fatalf("import: %v", err)
	}
	// Idempotent second import
	if _, err := svc2.ImportSeedDocument(ctx, doc); err != nil {
		t.Fatalf("re-import: %v", err)
	}

	got, err := st2.GetScopeBySlug("auth")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Auth" || got.Kind != "feature" {
		t.Fatalf("%+v", got)
	}
	for _, rel := range []string{domain.RelScopeMember, domain.RelAPIContract, domain.RelImplements, domain.RelBlocks} {
		rows, err := st2.ListLinksByRel(rel)
		if err != nil || len(rows) < 1 {
			t.Fatalf("imported %s: %v %#v", rel, err, rows)
		}
	}
}

func TestINFERREDScopeMemberSeedRoundTrip(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	auth, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth-fe", Title: "Auth frontend", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Ship"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Fix login session", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}

	rep, err := svc.InferScopes(ctx, domain.InferOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Inserted != 1 {
		t.Fatalf("inserted=%d want 1; report=%+v", rep.Inserted, rep)
	}

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{})
	if err != nil {
		t.Fatal(err)
	}
	var exported *domain.SeedLink
	for i := range doc.Links {
		l := &doc.Links[i]
		if l.Rel == domain.RelScopeMember && l.From == task.ID && l.To == auth.ID {
			exported = l
			break
		}
	}
	if exported == nil {
		t.Fatal("export missing scope_member for inferred task")
	}
	if exported.SourceType != domain.SourceTypeInferred {
		t.Fatalf("export source_type=%q want INFERRED", exported.SourceType)
	}
	if exported.Confidence != domain.InferredLinkConfidence {
		t.Fatalf("export confidence=%v want %v", exported.Confidence, domain.InferredLinkConfidence)
	}

	dir2 := t.TempDir()
	st2, err := store.Open(dir2)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	svc2 := domain.New(st2)
	if _, err := svc2.ImportSeedDocument(ctx, doc); err != nil {
		t.Fatalf("import: %v", err)
	}

	links, err := st2.ListLinksByRel(domain.RelScopeMember)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, l := range links {
		if l.FromID != task.ID || l.ToID != auth.ID {
			continue
		}
		found = true
		if l.SourceType != domain.SourceTypeInferred {
			t.Fatalf("import source_type=%q want INFERRED (H1 honesty)", l.SourceType)
		}
		if l.Confidence != domain.InferredLinkConfidence {
			t.Fatalf("import confidence=%v want %v", l.Confidence, domain.InferredLinkConfidence)
		}
	}
	if !found {
		t.Fatal("imported scope_member missing")
	}
}

func TestNoINFERREDFromScopePaths(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	sc, err := svc.CreateScope(ctx, domain.ScopeInput{Slug: "s", Title: "S", Kind: "layer"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "T"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, task.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	links, err := st.ListLinksByRel(domain.RelScopeMember)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range links {
		if l.SourceType == "INFERRED" {
			t.Fatal("S02 must not create INFERRED rows")
		}
	}
}
