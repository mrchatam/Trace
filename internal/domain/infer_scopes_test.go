package domain_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
)

func TestInferScopes_TitleRuleHappyPath(t *testing.T) {
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

	links, err := st.ListLinksFrom(domain.EntityTask, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, l := range links {
		if l.Rel != domain.RelScopeMember {
			continue
		}
		found = true
		if l.ToID != auth.ID {
			t.Fatalf("to_id=%q want %q", l.ToID, auth.ID)
		}
		if l.SourceType != domain.SourceTypeInferred {
			t.Fatalf("source_type=%q want INFERRED", l.SourceType)
		}
		if l.Confidence != domain.InferredLinkConfidence {
			t.Fatalf("confidence=%v want %v", l.Confidence, domain.InferredLinkConfidence)
		}
	}
	if !found {
		t.Fatal("expected INFERRED scope_member")
	}
}

func TestInferScopes_PathRuleHappyPath(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	be, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "backend", Title: "Backend layer", Kind: "layer",
	})
	if err != nil {
		t.Fatal(err)
	}
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Wire store", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	mustInProgress(t, svc, task.ID)
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Paths:  []domain.ChangePathInput{{Path: "internal/store/links.go"}},
	}); err != nil {
		t.Fatal(err)
	}

	rep, err := svc.InferScopes(ctx, domain.InferOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Inserted < 1 {
		t.Fatalf("inserted=%d want ≥1; %+v", rep.Inserted, rep)
	}
	links, err := st.ListLinksByRel(domain.RelScopeMember)
	if err != nil {
		t.Fatal(err)
	}
	ok := false
	for _, l := range links {
		if l.FromID == task.ID && l.ToID == be.ID && l.SourceType == domain.SourceTypeInferred {
			ok = true
			if l.Confidence != domain.InferredLinkConfidence {
				t.Fatalf("confidence=%v", l.Confidence)
			}
		}
	}
	if !ok {
		t.Fatalf("path rule did not link task→backend: %+v", links)
	}
}

func TestInferScopes_ExplicitWins(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	auth, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth", Title: "Auth", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	other, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "billing", Title: "Billing", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = other
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "auth login work", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkScopeMember(ctx, task.ID, auth.ID, domain.LinkMeta{
		SourceType: domain.DefaultSourceType, Confidence: 1.0,
	}); err != nil {
		t.Fatal(err)
	}
	before, err := st.ListLinksByRel(domain.RelScopeMember)
	if err != nil {
		t.Fatal(err)
	}

	rep, err := svc.InferScopes(ctx, domain.InferOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Inserted != 0 {
		t.Fatalf("inserted=%d want 0 (explicit wins)", rep.Inserted)
	}
	if rep.SkippedExplicit < 1 {
		t.Fatalf("skipped_explicit=%d want ≥1", rep.SkippedExplicit)
	}
	after, err := st.ListLinksByRel(domain.RelScopeMember)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Fatalf("link count changed %d→%d", len(before), len(after))
	}
	for _, l := range after {
		if l.FromID == task.ID {
			if l.SourceType != domain.DefaultSourceType {
				t.Fatalf("explicit source_type mutated to %q", l.SourceType)
			}
			if l.Confidence != 1.0 {
				t.Fatalf("explicit confidence mutated to %v", l.Confidence)
			}
		}
	}
}

func TestInferScopes_ConflictFailClosed(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	fe, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "frontend", Title: "Frontend", Kind: "layer",
	})
	if err != nil {
		t.Fatal(err)
	}
	auth, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth", Title: "Auth", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = fe
	_ = auth

	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	// Title → auth; path → frontend — different scopes → conflict.
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "auth login", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	mustInProgress(t, svc, task.ID)
	if _, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Paths:  []domain.ChangePathInput{{Path: "web/src/pages/home.tsx"}},
	}); err != nil {
		t.Fatal(err)
	}

	rep, err := svc.InferScopes(ctx, domain.InferOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Inserted != 0 {
		t.Fatalf("inserted=%d want 0 on conflict; %+v", rep.Inserted, rep)
	}
	if rep.SkippedConflict < 1 {
		t.Fatalf("skipped_conflict=%d want ≥1; %+v", rep.SkippedConflict, rep)
	}
	links, err := st.ListLinksFrom(domain.EntityTask, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range links {
		if l.Rel == domain.RelScopeMember {
			t.Fatalf("unexpected scope_member on conflict: %+v", l)
		}
	}
}

func TestInferScopes_IdempotentAndDryRun(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	_, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth", Title: "Auth", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.CreateTask(ctx, domain.TaskInput{Title: "session auth", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}

	dry, err := svc.InferScopes(ctx, domain.InferOptions{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if dry.Inserted != 0 {
		t.Fatalf("dry-run inserted=%d", dry.Inserted)
	}
	if len(dry.Candidates) < 1 {
		t.Fatalf("dry-run candidates empty: %+v", dry)
	}
	links, err := st.ListLinksByRel(domain.RelScopeMember)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 0 {
		t.Fatalf("dry-run wrote links: %+v", links)
	}

	first, err := svc.InferScopes(ctx, domain.InferOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if first.Inserted < 1 {
		t.Fatalf("first inserted=%d", first.Inserted)
	}
	second, err := svc.InferScopes(ctx, domain.InferOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if second.Inserted != 0 {
		t.Fatalf("second inserted=%d want 0", second.Inserted)
	}
	if second.SkippedExisting < 1 {
		t.Fatalf("second skipped_existing=%d", second.SkippedExisting)
	}
}

func TestInferScopes_NoGateRels(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	_, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug: "auth", Title: "Auth", Kind: "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.CreateTask(ctx, domain.TaskInput{Title: "auth", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.InferScopes(ctx, domain.InferOptions{}); err != nil {
		t.Fatal(err)
	}
	forbidden := []string{
		"review_judges_task", "decision_affects_task", "uncertainty_blocks_task",
		"change_implements_decision", "goal_has_task",
	}
	for _, rel := range forbidden {
		rows, err := st.ListLinksByRel(rel)
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != 0 {
			t.Fatalf("infer must not write %s: %+v", rel, rows)
		}
	}
}

func TestInferScopes_NoIndexHookRegistration(t *testing.T) {
	// CI-style: index/watch sources must not call InferScopes.
	root := filepath.Join("..", "..", "cmd", "trace")
	for _, name := range []string{"index.go", "index_watch.go"} {
		b, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		src := string(b)
		for _, needle := range []string{"InferScopes", "graph infer", "cmdGraph", "infer_scopes"} {
			if strings.Contains(src, needle) {
				t.Fatalf("%s must not reference %q (D4 CLI-only)", name, needle)
			}
		}
	}
}

func TestInferLinkConfidenceLowerThanExplicit(t *testing.T) {
	if domain.InferredLinkConfidence >= 0.9 {
		t.Fatalf("InferredLinkConfidence=%v must be < typical explicit 0.9", domain.InferredLinkConfidence)
	}
	if domain.SourceTypeInferred != "INFERRED" {
		t.Fatalf("SourceTypeInferred=%q", domain.SourceTypeInferred)
	}
}
