package domain_test

import (
	"context"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestScopeMVPRelsAndCRUD(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()

	// Collision consts must stay distinct from causal names.
	if domain.RelImplements == domain.RelChangeImplementsDecision {
		t.Fatal("implements must ≠ change_implements_decision")
	}
	if domain.RelBlocks == domain.RelUncertaintyBlocksTask {
		t.Fatal("blocks must ≠ uncertainty_blocks_task")
	}
	if domain.RelImplements != "implements" || domain.RelBlocks != "blocks" ||
		domain.RelScopeMember != "scope_member" || domain.RelAPIContract != "api_contract" {
		t.Fatalf("MVP rel const mismatch")
	}
	if domain.EntityScope == domain.EntityPlanScope {
		t.Fatal("EntityScope must ≠ EntityPlanScope")
	}

	sc, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth-be", Title: "Auth backend", Kind: "feature",
	})
	if err != nil {
		t.Fatalf("CreateScope: %v", err)
	}
	_, err = svc.CreateScope(ctx, domain.ScopeInput{Slug: "x", Title: "y", Kind: "nope"})
	if err == nil {
		t.Fatal("invalid kind must fail")
	}

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	t1, err := svc.CreateTask(ctx, domain.TaskInput{Title: "FE login", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	t2, err := svc.CreateTask(ctx, domain.TaskInput{Title: "BE login", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	d, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Use JWT"})
	if err != nil {
		t.Fatal(err)
	}

	if err := svc.LinkScopeMember(ctx, t1.ID, sc.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkScopeMember: %v", err)
	}
	links, err := st.ListLinksFrom(domain.EntityTask, t1.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, l := range links {
		if l.Rel == domain.RelScopeMember && l.ToType == domain.EntityScope && l.ToID == sc.ID {
			found = true
			if l.SourceType != domain.DefaultSourceType {
				t.Fatalf("source_type=%q", l.SourceType)
			}
		}
	}
	if !found {
		t.Fatal("scope_member link missing")
	}

	if err := svc.LinkAPIContract(ctx, t1.ID, t2.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkAPIContract: %v", err)
	}
	if err := svc.LinkImplements(ctx, d.ID, t2.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkImplements: %v", err)
	}
	if err := svc.LinkBlocks(ctx, t1.ID, t2.ID, domain.LinkMeta{}); err != nil {
		t.Fatalf("LinkBlocks: %v", err)
	}

	for _, want := range []string{domain.RelAPIContract, domain.RelImplements, domain.RelBlocks} {
		rows, err := st.ListLinksByRel(want)
		if err != nil || len(rows) < 1 {
			t.Fatalf("rel %s: %v %#v", want, err, rows)
		}
		if rows[0].Rel != want {
			t.Fatalf("stored rel %q want %q", rows[0].Rel, want)
		}
	}

	// Duplicate UNIQUE — InsertLink errors (match existing link style).
	err = svc.LinkScopeMember(ctx, t1.ID, sc.ID, domain.LinkMeta{})
	if err == nil {
		t.Fatal("duplicate scope_member should fail UNIQUE")
	}
	if !strings.Contains(err.Error(), "UNIQUE") && !strings.Contains(err.Error(), "constraint") &&
		!strings.Contains(err.Error(), "unique") {
		// sqlite may wrap; still must error
		t.Logf("duplicate err (ok if constraint-related): %v", err)
	}
}
