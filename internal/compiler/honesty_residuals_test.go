package compiler_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

// DF-61: >8 disk-stale files → cap 8 paths + stale_total + stale_truncated + MD total.
func TestIndexHonestyStaleTotalTruncated(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	c := compiler.New(st).WithRetrieval(retrieval.New(st))
	ctx := context.Background()

	const token = "stalecaptoken"
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: token})
	if err != nil {
		t.Fatalf("task: %v", err)
	}

	writeIndexed := func(rel, body string) {
		t.Helper()
		abs := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(abs, []byte(body), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		sum := sha256.Sum256([]byte(body))
		if _, err := st.UpsertFile(rel, hex.EncodeToString(sum[:]), nil); err != nil {
			t.Fatalf("UpsertFile %s: %v", rel, err)
		}
	}

	const n = 13
	var rels []string
	for i := 0; i < n; i++ {
		rel := fmt.Sprintf("pkg/aaa_%s_%02d.go", token, i)
		writeIndexed(rel, fmt.Sprintf("package pkg\n// body %d\n", i))
		rels = append(rels, rel)
	}
	for _, rel := range rels {
		abs := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.WriteFile(abs, []byte("package pkg\n// mutated\n"), 0o644); err != nil {
			t.Fatalf("mutate %s: %v", rel, err)
		}
	}

	pkt, err := c.TaskContext(ctx, task.ID, compiler.ContextOptions{
		MaxItems:        32,
		IncludeMarkdown: true,
	})
	if err != nil {
		t.Fatalf("TaskContext: %v", err)
	}
	if pkt.IndexHonesty == nil {
		t.Fatal("expected index_honesty")
	}
	if got := len(pkt.IndexHonesty.StalePaths); got != 8 {
		t.Fatalf("stale_paths len=%d want 8; %+v", got, pkt.IndexHonesty.StalePaths)
	}
	if pkt.IndexHonesty.StaleTotal <= 8 {
		t.Fatalf("stale_total=%d want >8", pkt.IndexHonesty.StaleTotal)
	}
	if !pkt.IndexHonesty.StaleTruncated {
		t.Fatal("expected stale_truncated=true")
	}
	md := pkt.Markdown()
	if !strings.Contains(md, fmt.Sprintf("stale_total=%d", pkt.IndexHonesty.StaleTotal)) {
		t.Fatalf("markdown missing stale_total signal:\n%s", md)
	}
}

// DF-62: MaxItems trim drops a disk-stale file → honesty still non-null (pre-trim universe).
func TestIndexHonestyPreTrimUniverse(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	c := compiler.New(st).WithRetrieval(retrieval.New(st))
	ctx := context.Background()

	const token = "staledroptoken"
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: token})
	if err != nil {
		t.Fatalf("task: %v", err)
	}
	// Fill MaxItems with higher-priority decisions so the FTS file is trimmed out.
	for i := 0; i < 6; i++ {
		dec, err := svc.CreateDecision(ctx, domain.DecisionInput{
			Title: fmt.Sprintf("staledropdec%d", i),
		})
		if err != nil {
			t.Fatalf("dec: %v", err)
		}
		if err := svc.LinkDecisionTask(ctx, dec.ID, task.ID, domain.LinkMeta{}); err != nil {
			t.Fatalf("link: %v", err)
		}
	}

	rel := "pkg/zzz_" + token + ".go"
	abs := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	body := "package pkg\n// original\n"
	if err := os.WriteFile(abs, []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	sum := sha256.Sum256([]byte(body))
	rec, err := st.UpsertFile(rel, hex.EncodeToString(sum[:]), nil)
	if err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if err := os.WriteFile(abs, []byte("package pkg\n// mutated stale\n"), 0o644); err != nil {
		t.Fatalf("mutate: %v", err)
	}

	pkt, err := c.TaskContext(ctx, task.ID, compiler.ContextOptions{
		MaxItems:        3, // task + task_state + 1 decision; file dropped by trim
		IncludeMarkdown: true,
	})
	if err != nil {
		t.Fatalf("TaskContext: %v", err)
	}
	for _, it := range pkt.Items {
		if it.EntityType == "file" && it.EntityID == rec.ID {
			t.Fatalf("precondition failed: stale file still kept; items=%+v", pkt.Items)
		}
	}
	if pkt.IndexHonesty == nil || pkt.IndexHonesty.StaleTotal == 0 {
		t.Fatalf("expected pre-trim index_honesty for dropped stale file; got %+v", pkt.IndexHonesty)
	}
	found := false
	for _, p := range pkt.IndexHonesty.StalePaths {
		if p == rel {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("stale_paths missing %q: %+v", rel, pkt.IndexHonesty.StalePaths)
	}
}

// DF-63: candidates capped → items_total reflects full admit universe (≫64), not post-cap len.
func TestCandidateCapAdmitUniverseTotal(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	c := compiler.New(st).WithRetrieval(retrieval.New(st))
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Admit universe cap"})
	if err != nil {
		t.Fatalf("task: %v", err)
	}
	const n = 80
	for i := 0; i < n; i++ {
		dec, err := svc.CreateDecision(ctx, domain.DecisionInput{
			Title: fmt.Sprintf("admituniv%d", i),
		})
		if err != nil {
			t.Fatalf("dec: %v", err)
		}
		if err := svc.LinkDecisionTask(ctx, dec.ID, task.ID, domain.LinkMeta{}); err != nil {
			t.Fatalf("link: %v", err)
		}
	}

	pkt, err := c.TaskContext(ctx, task.ID, compiler.ContextOptions{
		MaxItems:        compiler.DefaultMaxItems,
		IncludeMarkdown: true,
	})
	if err != nil {
		t.Fatalf("TaskContext: %v", err)
	}
	if !pkt.Budget.CandidatesCapped {
		t.Fatal("expected candidates_capped")
	}
	// layer0 (task+task_state) + 80 decisions ≈ 82; must not be post-cap 64.
	if pkt.Budget.ItemsTotal < n+2 {
		t.Fatalf("items_total=%d want >= %d (admit universe)", pkt.Budget.ItemsTotal, n+2)
	}
	if pkt.Budget.ItemsTotal <= compiler.MaxCandidateHits {
		t.Fatalf("items_total=%d still looks post-cap (<=%d)", pkt.Budget.ItemsTotal, compiler.MaxCandidateHits)
	}
	want := fmt.Sprintf("items=%d/%d", pkt.Budget.ItemsKept, pkt.Budget.ItemsTotal)
	md := pkt.Markdown()
	if !strings.Contains(md, want) {
		t.Fatalf("markdown missing %q:\n%s", want, md)
	}
}

// DF-65: FTS-admitted importer file → Expand file seeds → neighbor Item carries edge_provenance.
func TestContextImportHopEdgeProvenance(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	c := compiler.New(st).WithRetrieval(retrieval.New(st))
	ctx := context.Background()

	const token = "ctxprovhop"
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: token})
	if err != nil {
		t.Fatalf("task: %v", err)
	}

	writeFile := func(rel, body string) store.FileRecord {
		t.Helper()
		abs := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(abs, []byte(body), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		sum := sha256.Sum256([]byte(body))
		rec, err := st.UpsertFile(rel, hex.EncodeToString(sum[:]), nil)
		if err != nil {
			t.Fatalf("UpsertFile %s: %v", rel, err)
		}
		return rec
	}

	importer := writeFile("src/app_"+token+".js", "import './util'\n")
	util := writeFile("src/util.js", "export const x = 1\n")
	if err := st.ReplaceFileImports(importer.Path, []store.Import{
		{ImportedPath: "./util", Provenance: store.ImportProvenanceExtracted},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	pkt, err := c.ExpandContext(ctx, task.ID, 2, compiler.ContextOptions{
		MaxItems:        32,
		IncludeMarkdown: true,
	})
	if err != nil {
		t.Fatalf("ExpandContext: %v", err)
	}
	var saw bool
	for _, it := range pkt.Items {
		if it.EntityType == "file" && it.EntityID == util.ID && it.EdgeProvenance == store.ImportProvenanceExtracted {
			saw = true
			break
		}
	}
	if !saw {
		t.Fatalf("expected util neighbor with edge_provenance=EXTRACTED; items=%+v", pkt.Items)
	}
	md := pkt.Markdown()
	if !strings.Contains(md, "edge_provenance: `EXTRACTED`") {
		t.Fatalf("markdown missing edge_provenance EXTRACTED:\n%s", md)
	}
}
