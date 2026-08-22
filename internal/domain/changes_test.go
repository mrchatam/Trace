package domain_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
	_ "modernc.org/sqlite"
)

func mustChangeTask(t *testing.T, svc *domain.Service) store.Task {
	t.Helper()
	task, err := svc.CreateTask(context.Background(), domain.TaskInput{Title: "change-task"})
	if err != nil {
		t.Fatal(err)
	}
	return task
}

func mustIndexedCommit(t *testing.T, st *store.Store, oid, subject string, paths []store.IndexedPathChange) {
	t.Helper()
	seq, err := st.NextCommitSeq()
	if err != nil {
		t.Fatal(err)
	}
	if err := st.UpsertIndexedCommit(store.IndexedCommit{
		OID:         oid,
		CommittedAt: "2026-01-01T00:00:00Z",
		Subject:     subject,
		Seq:         seq,
	}, paths); err != nil {
		t.Fatal(err)
	}
}

func TestPromoteVCSCommitCreatesChangeIdempotent(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	const oid = "abc123def4567890abc123def4567890abc123de"

	mustIndexedCommit(t, st, oid, "add handler", []store.IndexedPathChange{
		{Path: "internal/handler.go", Status: "A"},
	})

	first, err := svc.PromoteVCSCommitToChange(ctx, oid)
	if err != nil {
		t.Fatalf("first promote: %v", err)
	}
	if first.ID == "" || first.GitCommit != oid || first.TaskID != domain.VCSCaptureTaskID {
		t.Fatalf("first change: %+v", first)
	}
	if first.SourceType != domain.SourceTypeVCS || first.Status != store.ChangeStatusRecorded {
		t.Fatalf("provenance/status: %+v", first)
	}
	if first.Reason != "add handler" {
		t.Fatalf("reason from subject: %q", first.Reason)
	}

	second, err := svc.PromoteVCSCommitToChange(ctx, oid)
	if err != nil {
		t.Fatalf("second promote: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("idempotent id: first=%q second=%q", first.ID, second.ID)
	}

	all, err := st.ListAllChanges()
	if err != nil {
		t.Fatal(err)
	}
	var matches int
	for _, c := range all {
		if c.GitCommit == oid {
			matches++
		}
	}
	if matches != 1 {
		t.Fatalf("want one change row for oid, got %d: %+v", matches, all)
	}
}

func TestPromoteVCSCommitRecordsPathsNoBlob(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	const oid = "1111111111111111111111111111111111111111"

	mustIndexedCommit(t, st, oid, "touch files", []store.IndexedPathChange{
		{Path: "pkg/foo.go", Status: "M"},
		{Path: "pkg/bar.go", Status: "A"},
	})

	c, err := svc.PromoteVCSCommitToChange(ctx, oid)
	if err != nil {
		t.Fatalf("promote: %v", err)
	}
	paths, err := svc.ListChangePaths(ctx, c.ID)
	if err != nil || len(paths) != 2 {
		t.Fatalf("paths: %+v err=%v", paths, err)
	}
	for _, p := range paths {
		if p.Status == "" {
			t.Fatalf("status copied: %+v", p)
		}
		if strings.Contains(strings.ToLower(p.Path), "blob") {
			t.Fatalf("path looks like blob: %+v", p)
		}
	}

	bad, where, err := st.HasBlobLikeColumns()
	if err != nil || bad {
		t.Fatalf("HasBlobLikeColumns: bad=%v where=%s err=%v", bad, where, err)
	}
}

func TestPromoteVCSCommitSkipsNonMeaningful(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	const oid = "2222222222222222222222222222222222222222"

	mustIndexedCommit(t, st, oid, "docs only", []store.IndexedPathChange{
		{Path: "docs/notes.md", Status: "M"},
		{Path: "CHANGELOG", Status: "M"},
	})

	skipped, err := svc.PromoteVCSCommitToChange(ctx, oid)
	if err != nil {
		t.Fatalf("promote meaningful filter: %v", err)
	}
	if skipped.ID != "" {
		t.Fatalf("expected skip, got change: %+v", skipped)
	}
	all, err := st.ListAllChanges()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 0 {
		t.Fatalf("docs-only commit must not create change: %+v", all)
	}

	captured, err := svc.PromoteVCSCommitToChangeOpts(ctx, oid, domain.PromoteVCSCommitOptions{AllPaths: true})
	if err != nil {
		t.Fatalf("promote --all: %v", err)
	}
	if captured.ID == "" {
		t.Fatal("expected change with --all")
	}
	paths, err := svc.ListChangePaths(ctx, captured.ID)
	if err != nil || len(paths) != 2 {
		t.Fatalf("all paths: %+v err=%v", paths, err)
	}
}

func TestCreateChangeWithGitSHAAndPathsNoBlob(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "ABCDEF0",
		Actor:     "agent",
		Reason:    "add store",
		Paths: []domain.ChangePathInput{
			{Path: "./internal/store/changes.go", Status: "A"},
			{Path: "internal/domain/changes.go", Status: "M", SymbolID: "sym-1"},
		},
	})
	if err != nil {
		t.Fatalf("CreateChange: %v", err)
	}
	if c.Status != store.ChangeStatusRecorded {
		t.Fatalf("status: %+v", c)
	}
	if c.GitCommit != "abcdef0" {
		t.Fatalf("git_commit stored lowercase: %q", c.GitCommit)
	}

	paths, err := svc.ListChangePaths(ctx, c.ID)
	if err != nil || len(paths) != 2 {
		t.Fatalf("paths: %+v err=%v", paths, err)
	}
	for _, p := range paths {
		if p.Path == "" || strings.Contains(strings.ToLower(p.Path), "blob") {
			t.Fatalf("path looks like content: %+v", p)
		}
		if p.Path != store.NormalizePath(p.Path) {
			t.Fatalf("path not normalized: %+v", p)
		}
	}

	bad, where, err := st.HasBlobLikeColumns()
	if err != nil || bad {
		t.Fatalf("HasBlobLikeColumns: bad=%v where=%s err=%v", bad, where, err)
	}

	evs, err := st.ListEventsByEntity(domain.EntityChange, c.ID)
	if err != nil || len(evs) < 1 || evs[0].Type != domain.EventEntityCreated {
		t.Fatalf("entity.created: %+v err=%v", evs, err)
	}

	listed, err := svc.ListChangesByTaskID(ctx, task.ID)
	if err != nil || len(listed) != 1 || listed[0].ID != c.ID {
		t.Fatalf("ListChangesByTaskID: %+v err=%v", listed, err)
	}
}

func TestCreateChangeRequiresTaskIDAndPath(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	_, err := svc.CreateChange(ctx, domain.ChangeInput{
		Paths: []domain.ChangePathInput{{Path: "a.go"}},
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("missing TaskID: want ErrValidation, got %v", err)
	}

	_, err = svc.CreateChange(ctx, domain.ChangeInput{TaskID: task.ID})
	if !errors.As(err, &ve) {
		t.Fatalf("missing paths: want ErrValidation, got %v", err)
	}

	_, err = svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Paths:  []domain.ChangePathInput{{Path: "  "}},
	})
	if !errors.As(err, &ve) {
		t.Fatalf("empty path: want ErrValidation, got %v", err)
	}
}

func TestRecordExpectedThenActualSupported(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abcdef0",
		Paths:     []domain.ChangePathInput{{Path: "a.go"}},
		Expected: []domain.ExpectedEffectInput{
			{Dimension: "latency", Expected: "p99 < 50ms"},
			{Dimension: "correctness", Expected: "tests pass"},
		},
	})
	if err != nil {
		t.Fatalf("CreateChange: %v", err)
	}
	if c.Status != store.ChangeStatusRecorded {
		t.Fatalf("before actuals: %+v", c)
	}

	one, disc, err := svc.RecordActualEffect(ctx, c.ID, domain.RecordActualEffectInput{
		Dimension:  "latency",
		Actual:     "p99 40ms",
		Comparison: store.EffectComparisonSupported,
	})
	if err != nil || disc != nil {
		t.Fatalf("first actual: %+v disc=%v err=%v", one, disc, err)
	}
	if one.Comparison != store.EffectComparisonSupported {
		t.Fatalf("comparison: %+v", one)
	}
	mid, err := svc.GetChange(ctx, c.ID)
	if err != nil || mid.Status != store.ChangeStatusRecorded {
		t.Fatalf("partial compare stays RECORDED: %+v err=%v", mid, err)
	}

	two, _, err := svc.RecordActualEffect(ctx, c.ID, domain.RecordActualEffectInput{
		Dimension:  "correctness",
		Actual:     "ok",
		Comparison: store.EffectComparisonSupported,
	})
	if err != nil {
		t.Fatalf("second actual: %v", err)
	}
	if two.Comparison != store.EffectComparisonSupported {
		t.Fatalf("second comparison: %+v", two)
	}
	got, err := svc.GetChange(ctx, c.ID)
	if err != nil || got.Status != store.ChangeStatusCompared {
		t.Fatalf("all compared → COMPARED: %+v err=%v", got, err)
	}

	evs, err := st.ListEventsByEntity(domain.EntityEffect, one.ID)
	if err != nil {
		t.Fatal(err)
	}
	var sawCompared bool
	for _, e := range evs {
		if e.Type == domain.EventEffectCompared {
			sawCompared = true
		}
	}
	if !sawCompared {
		t.Fatalf("missing effect.compared: %+v", evs)
	}
}

func TestRecordActualRequiresExpectedDimension(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Paths:  []domain.ChangePathInput{{Path: "a.go"}},
		Expected: []domain.ExpectedEffectInput{
			{Dimension: "correctness", Expected: "ok"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = svc.RecordActualEffect(ctx, c.ID, domain.RecordActualEffectInput{
		Dimension:  "latency",
		Actual:     "slow",
		Comparison: store.EffectComparisonSupported,
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("want ErrValidation, got %v", err)
	}
}

func TestUnknownEffectComparisonFailClosed(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Paths:  []domain.ChangePathInput{{Path: "a.go"}},
		Expected: []domain.ExpectedEffectInput{
			{Dimension: "correctness", Expected: "ok"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = svc.RecordActualEffect(ctx, c.ID, domain.RecordActualEffectInput{
		Dimension:  "correctness",
		Actual:     "maybe",
		Comparison: "unknown",
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("unknown comparison: want ErrValidation, got %v", err)
	}

	_, _, err = svc.RecordActualEffect(ctx, c.ID, domain.RecordActualEffectInput{
		Dimension:  "correctness",
		Actual:     "maybe",
		Comparison: "",
	})
	if !errors.As(err, &ve) {
		t.Fatalf("empty comparison: want ErrValidation, got %v", err)
	}
}

func TestRecordActualContradictedLinksHypothesisWithoutDiscoveryFork(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abcdef0",
		Paths:     []domain.ChangePathInput{{Path: "a.go"}},
		Expected: []domain.ExpectedEffectInput{
			{Dimension: "correctness", Expected: "green"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	eff, disc, err := svc.RecordActualEffect(ctx, c.ID, domain.RecordActualEffectInput{
		Dimension:        "correctness",
		Actual:           "red",
		Comparison:       store.EffectComparisonContradicted,
		CreateHypothesis: true,
		HypothesisTitle:  "flaky fixture",
		HypothesisBody:   "seed data",
	})
	if err != nil {
		t.Fatalf("RecordActualEffect: %v", err)
	}
	if disc != nil {
		t.Fatalf("must not emit discovery unless EmitDiscovery: %+v", disc)
	}

	discs, err := st.ListDiscoveries()
	if err != nil {
		t.Fatal(err)
	}
	if len(discs) != 0 {
		t.Fatalf("zero extra Discovery: %+v", discs)
	}

	links, err := st.ListLinksTo(domain.EntityEffect, eff.ID)
	if err != nil {
		t.Fatal(err)
	}
	var hypID string
	for _, l := range links {
		if l.Rel == domain.RelHypothesisExplainsEffect && l.FromType == domain.EntityHypothesis {
			hypID = l.FromID
		}
		if l.Rel == domain.RelDiscoveryFromContradictedEffect {
			t.Fatalf("must not fork discovery: %+v", l)
		}
	}
	if hypID == "" {
		t.Fatalf("missing hypothesis_explains_effect: %+v", links)
	}
	h, err := svc.GetHypothesis(ctx, hypID)
	if err != nil || h.Status != store.HypothesisStatusOpen || h.Title != "flaky fixture" {
		t.Fatalf("hypothesis: %+v err=%v", h, err)
	}
	if _, err := svc.GetHypothesis(ctx, hypID); err != nil {
		t.Fatal(err)
	}
	for _, d := range discs {
		if d.ID == h.ID || d.Title == h.Title {
			t.Fatalf("hypothesis must not be a discovery: %+v", d)
		}
	}
}

func TestRecordActualContradictedOptionalPlanAffectingDiscovery(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Paths:  []domain.ChangePathInput{{Path: "a.go"}},
		Expected: []domain.ExpectedEffectInput{
			{Dimension: "correctness", Expected: "green"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	eff, disc, err := svc.RecordActualEffect(ctx, c.ID, domain.RecordActualEffectInput{
		Dimension:      "correctness",
		Actual:         "red",
		Comparison:     store.EffectComparisonContradicted,
		EmitDiscovery:  true,
		DiscoveryTitle: "effect contradicted",
		DiscoveryBody:  "tests failed",
	})
	if err != nil {
		t.Fatalf("RecordActualEffect: %v", err)
	}
	if disc == nil || disc.Severity != domain.SeverityPlanAffecting || disc.Title != "effect contradicted" {
		t.Fatalf("discovery: %+v", disc)
	}
	links, err := svc.ListLinksFrom(ctx, domain.EntityDiscovery, disc.ID)
	if err != nil {
		t.Fatal(err)
	}
	var saw bool
	for _, l := range links {
		if l.Rel == domain.RelDiscoveryFromContradictedEffect && l.ToType == domain.EntityEffect && l.ToID == eff.ID {
			saw = true
		}
	}
	if !saw {
		t.Fatalf("discovery_from_contradicted_effect: %+v", links)
	}
}

func TestContradictedEffectFiresDecisionReconsideration(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Use SQLite"})
	if err != nil {
		t.Fatal(err)
	}

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:     task.ID,
		GitCommit:  "abcdef0",
		DecisionID: dec.ID,
		Paths:      []domain.ChangePathInput{{Path: "a.go"}},
		Expected: []domain.ExpectedEffectInput{
			{Dimension: "latency", Expected: "fast"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	eff, _, err := svc.RecordActualEffect(ctx, c.ID, domain.RecordActualEffectInput{
		Dimension:  "latency",
		Actual:     "slow",
		Comparison: store.EffectComparisonContradicted,
	})
	if err != nil {
		t.Fatalf("RecordActualEffect: %v", err)
	}

	rows, err := st.ListDecisionReconsiderationsByDecisionID(dec.ID)
	if err != nil || len(rows) != 1 {
		t.Fatalf("reconsiderations: %+v err=%v", rows, err)
	}
	if rows[0].Trigger != store.ReconsiderTriggerContradictedEffect || rows[0].Status != store.ReconsiderStatusFired {
		t.Fatalf("row: %+v", rows[0])
	}
	if rows[0].RelatedType != domain.EntityEffect || rows[0].RelatedID != eff.ID {
		t.Fatalf("related: %+v", rows[0])
	}
	got, err := svc.GetDecision(ctx, dec.ID)
	if err != nil || got.Title != "Use SQLite" {
		t.Fatalf("decision preserved: %+v err=%v", got, err)
	}
}

func TestContradictedEffectDoesNotCreateRegressionOrAutoReplan(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abcdef0",
		Paths:     []domain.ChangePathInput{{Path: "a.go"}},
		Expected: []domain.ExpectedEffectInput{
			{Dimension: "correctness", Expected: "green"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.RecordActualEffect(ctx, c.ID, domain.RecordActualEffectInput{
		Dimension:  "correctness",
		Actual:     "red",
		Comparison: store.EffectComparisonContradicted,
	}); err != nil {
		t.Fatal(err)
	}

	n, err := st.CountOpenRegressionsByTaskID(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatal("contradicted effect must not auto-create regression row")
	}

	evs, err := st.ListRecentEvents(200)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range evs {
		if e.Type == domain.EventDeliberationTransition {
			t.Fatalf("must not auto-replan: %+v", e)
		}
	}
	if _, err := svc.GetDeliberationState(ctx, task.ID); err == nil {
		t.Fatal("must not write deliberation_state")
	}
}

func TestParentChangeChain(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	parent, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abcdef0",
		Paths:     []domain.ChangePathInput{{Path: "a.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	child, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:         task.ID,
		ParentChangeID: parent.ID,
		GitCommit:      "1234567",
		Paths:          []domain.ChangePathInput{{Path: "a.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if child.ParentChangeID != parent.ID {
		t.Fatalf("parent_change_id: %+v", child)
	}
	got, err := svc.GetChange(ctx, parent.ID)
	if err != nil || got.Status != store.ChangeStatusRecorded {
		t.Fatalf("parent must not auto-SUPERSEDE: %+v err=%v", got, err)
	}
}

func TestResolveChangePathViaGitNotSQLite(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)
	const marker = "UNIQUE_CHANGE_BLOB_MARKER_P20S03"

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abcdef0",
		Paths:     []domain.ChangePathInput{{Path: "pkg/foo.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	repo := &vcs.Fake{Files: map[string][]byte{
		"abcdef0:pkg/foo.go": []byte(marker),
	}}
	got, err := svc.ResolveChangePath(ctx, c.ID, "pkg/foo.go", repo)
	if err != nil {
		t.Fatalf("ResolveChangePath: %v", err)
	}
	if string(got) != marker {
		t.Fatalf("bytes: %q", got)
	}

	bad, where, err := st.HasBlobLikeColumns()
	if err != nil || bad {
		t.Fatalf("HasBlobLikeColumns: bad=%v where=%s err=%v", bad, where, err)
	}
	paths, err := svc.ListChangePaths(ctx, c.ID)
	if err != nil || len(paths) != 1 || paths[0].Path != "pkg/foo.go" {
		t.Fatalf("paths: %+v err=%v", paths, err)
	}
	if strings.Contains(paths[0].Path, marker) {
		t.Fatal("path stored file bytes")
	}

	db, err := sql.Open("sqlite", st.DBPath())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var hits int
	if err := db.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM changes WHERE instr(reason, ?) > 0 OR instr(git_commit, ?) > 0) +
			(SELECT COUNT(*) FROM change_paths WHERE instr(path, ?) > 0) +
			(SELECT COUNT(*) FROM effects WHERE instr(expected, ?) > 0 OR instr(actual, ?) > 0) +
			(SELECT COUNT(*) FROM events WHERE instr(payload_json, ?) > 0)
	`, marker, marker, marker, marker, marker, marker).Scan(&hits); err != nil {
		t.Fatal(err)
	}
	if hits != 0 {
		t.Fatalf("SQLite copied ShowFile bytes (%d hits)", hits)
	}
}

func TestResolveChangePathFailsClosedWithoutCommit(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Paths:  []domain.ChangePathInput{{Path: "pkg/foo.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Status != store.ChangeStatusOpen || c.GitCommit != "" {
		t.Fatalf("OPEN without SHA: %+v", c)
	}

	repo := &vcs.Fake{Files: map[string][]byte{"HEAD:pkg/foo.go": []byte("x")}}
	_, err = svc.ResolveChangePath(ctx, c.ID, "pkg/foo.go", repo)
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("want ErrValidation, got %v", err)
	}

	_, err = svc.ResolveChangePath(ctx, c.ID, "pkg/foo.go", nil)
	if !errors.As(err, &ve) {
		t.Fatalf("nil repo: want ErrValidation, got %v", err)
	}

	c2, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    task.ID,
		GitCommit: "abcdef0",
		Paths:     []domain.ChangePathInput{{Path: "pkg/foo.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.ResolveChangePath(ctx, c2.ID, "missing.go", repo)
	if !errors.As(err, &ve) {
		t.Fatalf("unknown path: want ErrValidation, got %v", err)
	}
}

func TestOversizedEffectTextFailClosed(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)
	over := strings.Repeat("x", 8193)
	overDim := strings.Repeat("d", 65)

	_, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: over,
		Paths:  []domain.ChangePathInput{{Path: "a.go"}},
	})
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("oversized reason: want ErrValidation, got %v", err)
	}

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Paths:  []domain.ChangePathInput{{Path: "a.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.RecordExpectedEffect(ctx, c.ID, domain.ExpectedEffectInput{
		Dimension: "ok",
		Expected:  over,
	})
	if !errors.As(err, &ve) {
		t.Fatalf("oversized expected: want ErrValidation, got %v", err)
	}
	_, err = svc.RecordExpectedEffect(ctx, c.ID, domain.ExpectedEffectInput{
		Dimension: overDim,
		Expected:  "ok",
	})
	if !errors.As(err, &ve) {
		t.Fatalf("oversized dimension: want ErrValidation, got %v", err)
	}

	if _, err := svc.RecordExpectedEffect(ctx, c.ID, domain.ExpectedEffectInput{
		Dimension: "correctness",
		Expected:  "ok",
	}); err != nil {
		t.Fatal(err)
	}
	_, _, err = svc.RecordActualEffect(ctx, c.ID, domain.RecordActualEffectInput{
		Dimension:  "correctness",
		Actual:     over,
		Comparison: store.EffectComparisonSupported,
	})
	if !errors.As(err, &ve) {
		t.Fatalf("oversized actual: want ErrValidation, got %v", err)
	}
}

func TestRecordChangeCommitThenCompared(t *testing.T) {
	// Keeper for OPEN → RECORDED via RecordChangeCommit (not a planner-named test).
	svc, st := openDomain(t)
	ctx := context.Background()
	task := mustChangeTask(t, svc)

	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Paths:  []domain.ChangePathInput{{Path: "a.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := svc.RecordChangeCommit(ctx, c.ID, "ABCDEF012")
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != store.ChangeStatusRecorded || out.GitCommit != "abcdef012" {
		t.Fatalf("recorded: %+v", out)
	}
	evs, err := st.ListEventsByEntity(domain.EntityChange, c.ID)
	if err != nil {
		t.Fatal(err)
	}
	var saw bool
	for _, e := range evs {
		if e.Type == domain.EventChangeRecorded {
			saw = true
			var payload map[string]string
			_ = json.Unmarshal([]byte(e.PayloadJSON), &payload)
			if payload["git_commit"] != "abcdef012" {
				t.Fatalf("payload: %s", e.PayloadJSON)
			}
		}
	}
	if !saw {
		t.Fatalf("missing change.recorded: %+v", evs)
	}
	if _, err := svc.RecordChangeCommit(ctx, c.ID, "9999999"); err == nil {
		t.Fatal("replacing SHA must fail closed")
	}
}
