package retrieval_test

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func intentForTask(title, body string) retrieval.Intent {
	return retrieval.ExtractIntent(retrieval.IntentInput{TaskTitle: title, TaskBody: body})
}

// G6-C1: intent keyword matches discovery → graph_label_match.
func TestSearchGraphLabelsDiscovery(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	const token = "g6discoverylabeltokenXYZ"
	task, err := svc.CreateTask(ctx, domain.TaskInput{
		Title: "Find graph labels",
		Body:  "Search for " + token,
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	disc, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title: "Unlinked discovery",
		Body:  "Notes mention " + token + " in summary",
	})
	if err != nil {
		t.Fatalf("CreateDiscovery: %v", err)
	}
	_ = task

	hits, err := eng.SearchGraphLabels(ctx, intentForTask("Find graph labels", "Search for "+token), retrieval.SearchOptions{Limit: 16})
	if err != nil {
		t.Fatalf("SearchGraphLabels: %v", err)
	}
	var found *retrieval.Hit
	for i := range hits {
		if hits[i].EntityType == "discovery" && hits[i].EntityID == disc.ID {
			found = &hits[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("expected discovery hit; got %+v", hits)
	}
	if found.ReasonCode != retrieval.ReasonGraphLabelMatch {
		t.Fatalf("reason_code=%q want %q", found.ReasonCode, retrieval.ReasonGraphLabelMatch)
	}
}

// G6-C2: file/symbol hits excluded from concept channel.
func TestSearchGraphLabelsEntityFilter(t *testing.T) {
	eng, st, svc := openEngine(t)
	ctx := context.Background()

	const token = "g6entityfiltertokenXYZ"
	_, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title: "Concept " + token,
	})
	if err != nil {
		t.Fatalf("CreateDiscovery: %v", err)
	}
	_, err = st.UpsertFile("src/"+token+".go", "hash1", nil)
	if err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if err := st.ReplaceFileSymbols("src/"+token+".go", []store.Symbol{
		{Name: token + "Func", Kind: "function", StartLine: 1, EndLine: 5},
	}); err != nil {
		t.Fatalf("ReplaceFileSymbols: %v", err)
	}

	hits, err := eng.SearchGraphLabels(ctx, intentForTask("filter", token), retrieval.SearchOptions{Limit: 16})
	if err != nil {
		t.Fatalf("SearchGraphLabels: %v", err)
	}
	for _, h := range hits {
		if h.EntityType == "file" || h.EntityType == "symbol" {
			t.Fatalf("concept channel must exclude file/symbol: %+v", h)
		}
		if h.ReasonCode != retrieval.ReasonGraphLabelMatch {
			t.Fatalf("unexpected reason_code %q on %+v", h.ReasonCode, h)
		}
	}
}

// G6-C4: limit honored; no unbounded results.
func TestSearchGraphLabelsCap(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	const prefix = "g6caplabeltoken"
	for i := 0; i < 20; i++ {
		title := prefix + string(rune('A'+i%26))
		if _, err := svc.CreateAssumption(ctx, domain.AssumptionInput{Title: title}); err != nil {
			t.Fatalf("CreateAssumption %d: %v", i, err)
		}
	}

	hits, err := eng.SearchGraphLabels(ctx, intentForTask(prefix, prefix), retrieval.SearchOptions{Limit: 4})
	if err != nil {
		t.Fatalf("SearchGraphLabels: %v", err)
	}
	if len(hits) > 4 {
		t.Fatalf("expected at most 4 hits, got %d", len(hits))
	}

	hits64, err := eng.SearchGraphLabels(ctx, intentForTask(prefix, prefix), retrieval.SearchOptions{Limit: 128})
	if err != nil {
		t.Fatalf("SearchGraphLabels cap64: %v", err)
	}
	if len(hits64) > 64 {
		t.Fatalf("hard cap 64 exceeded: %d hits", len(hits64))
	}
}

// G6-C5: no vector deps; no semantic_match reason_code.
func TestSearchGraphLabelsNoSemantic(t *testing.T) {
	path := filepath.Join("..", "..", "internal", "retrieval", "concept.go")
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	for _, forbidden := range []string{"embedding", "vector", "openai", "semantic_match"} {
		if strings.Contains(string(src), forbidden) {
			t.Fatalf("concept.go must not reference %q", forbidden)
		}
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	for _, imp := range f.Imports {
		impPath := strings.Trim(imp.Path.Value, `"`)
		if strings.Contains(impPath, "embedding") || strings.Contains(impPath, "vector") {
			t.Fatalf("forbidden import %q", impPath)
		}
	}

	eng, _, svc := openEngine(t)
	ctx := context.Background()
	const token = "g6nosemantictokenXYZ"
	if _, err := svc.CreateClaim(ctx, domain.ClaimInput{Title: "Claim " + token}); err != nil {
		t.Fatalf("CreateClaim: %v", err)
	}
	hits, err := eng.SearchGraphLabels(ctx, intentForTask("claim", token), retrieval.SearchOptions{Limit: 8})
	if err != nil {
		t.Fatalf("SearchGraphLabels: %v", err)
	}
	for _, h := range hits {
		if h.ReasonCode == "semantic_match" {
			t.Fatalf("semantic_match forbidden: %+v", h)
		}
	}
}

// G6-C6: same input → same hits order.
func TestSearchGraphLabelsDeterministic(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	const token = "g6deterministictokenXYZ"
	if _, err := svc.CreateDecision(ctx, domain.DecisionInput{Title: "Decision " + token}); err != nil {
		t.Fatalf("CreateDecision: %v", err)
	}
	if _, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "Goal " + token}); err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	intent := intentForTask("decide", token)
	opts := retrieval.SearchOptions{Limit: 16}

	var baseline []retrieval.Hit
	for i := 0; i < 5; i++ {
		hits, err := eng.SearchGraphLabels(ctx, intent, opts)
		if err != nil {
			t.Fatalf("SearchGraphLabels: %v", err)
		}
		if i == 0 {
			baseline = hits
			continue
		}
		if len(hits) != len(baseline) {
			t.Fatalf("run %d len=%d baseline=%d", i, len(hits), len(baseline))
		}
		for j := range hits {
			if hits[j].EntityType != baseline[j].EntityType || hits[j].EntityID != baseline[j].EntityID {
				t.Fatalf("run %d differ at %d: %+v vs %+v", i, j, hits[j], baseline[j])
			}
		}
	}
}

type failGraphLabelsEngine struct {
	inner *retrieval.Engine
}

func (f failGraphLabelsEngine) Expand(ctx context.Context, seeds []retrieval.Hit, depth int) ([]retrieval.Hit, error) {
	return f.inner.Expand(ctx, seeds, depth)
}
func (f failGraphLabelsEngine) Search(ctx context.Context, query string, opts retrieval.SearchOptions) ([]retrieval.Hit, error) {
	return f.inner.Search(ctx, query, opts)
}
func (f failGraphLabelsEngine) SearchGraphLabels(ctx context.Context, intent retrieval.Intent, opts retrieval.SearchOptions) ([]retrieval.Hit, error) {
	return nil, fmt.Errorf("forced graph-label failure")
}
func (f failGraphLabelsEngine) Why(ctx context.Context, entityType, entityID string) (retrieval.WhyResult, error) {
	return f.inner.Why(ctx, entityType, entityID)
}

// G6-C7: concept path error → Context still returns packet with task seed.
func TestSearchGraphLabelsFailOpen(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Fail-open graph labels", Body: "task seed must survive"})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	inner := retrieval.New(st)
	c := compiler.New(st).WithRetrieval(failGraphLabelsEngine{inner: inner})

	pkt, err := c.TaskContext(ctx, task.ID, compiler.ContextOptions{})
	if err != nil {
		t.Fatalf("TaskContext: %v", err)
	}
	var sawTask bool
	for _, it := range pkt.Items {
		if it.EntityType == "task" && it.EntityID == task.ID {
			sawTask = true
			break
		}
	}
	if !sawTask {
		t.Fatalf("expected task seed in packet: %+v", pkt.Items)
	}
}

func TestMergeConceptHitsDedupe(t *testing.T) {
	base := []retrieval.Hit{{
		EntityType: "decision", EntityID: "d1",
		ReasonCode: retrieval.ReasonFTSMatch,
	}}
	concept := []retrieval.Hit{{
		EntityType: "decision", EntityID: "d1",
		ReasonCode: retrieval.ReasonGraphLabelMatch,
	}}
	merged := retrieval.MergeConceptHits(base, concept)
	if len(merged) != 1 {
		t.Fatalf("expected 1 hit, got %d", len(merged))
	}
	if merged[0].ReasonCode != retrieval.ReasonGraphLabelMatch {
		t.Fatalf("graph_label_match should win dedupe, got %q", merged[0].ReasonCode)
	}
}
