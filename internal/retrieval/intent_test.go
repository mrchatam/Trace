package retrieval_test

import (
	"context"
	"encoding/json"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
)

func TestExtractIntentFromTask(t *testing.T) {
	intent := retrieval.ExtractIntent(retrieval.IntentInput{
		TaskTitle: "Ship auth",
		TaskBody:  "Add JWT middleware",
	})
	if len(intent.Keywords) == 0 {
		t.Fatal("expected non-empty keywords")
	}
	joined := strings.Join(intent.Keywords, " ")
	for _, want := range []string{"ship", "auth", "add", "jwt", "middleware"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("keywords %q missing %q", joined, want)
		}
	}
	if intent.Scope != "task" || intent.Source != "task" {
		t.Fatalf("scope/source: scope=%q source=%q", intent.Scope, intent.Source)
	}
}

func TestExtractIntentEntityHints(t *testing.T) {
	const uuid = "550e8400-e29b-41d4-a716-446655440000"
	body := uuid + " touch internal/foo/bar.go and RunHandler"
	intent := retrieval.ExtractIntent(retrieval.IntentInput{
		TaskTitle: "Entity hints",
		TaskBody:  body,
	})
	if len(intent.EntityHints) == 0 {
		t.Fatal("expected entity hints")
	}
	var sawUUID, sawPath, sawSymbol bool
	for _, h := range intent.EntityHints {
		switch h.Kind {
		case "uuid":
			if strings.EqualFold(h.Value, uuid) {
				sawUUID = true
			}
		case "path":
			if h.Value == "internal/foo/bar.go" {
				sawPath = true
			}
		case "symbol":
			if h.Value == "RunHandler" {
				sawSymbol = true
			}
		}
	}
	if !sawUUID || !sawPath || !sawSymbol {
		t.Fatalf("hints incomplete uuid=%v path=%v symbol=%v: %+v", sawUUID, sawPath, sawSymbol, intent.EntityHints)
	}
}

func TestExtractIntentQueryMerge(t *testing.T) {
	intent := retrieval.ExtractIntent(retrieval.IntentInput{
		TaskTitle: "Ship auth",
		TaskBody:  "Add JWT middleware",
		Query:     "extra token",
	})
	joined := strings.Join(intent.Keywords, " ")
	if !strings.Contains(joined, "extra") || !strings.Contains(joined, "token") {
		t.Fatalf("query tokens missing: %q", joined)
	}
	if !strings.Contains(joined, "ship") {
		t.Fatalf("task tokens missing: %q", joined)
	}
	count := strings.Count(joined, "ship")
	if count != 1 {
		t.Fatalf("expected deduped ship once, got %d in %q", count, joined)
	}
	if intent.Scope != "task+query" || intent.Source != "task,query" {
		t.Fatalf("scope/source: scope=%q source=%q", intent.Scope, intent.Source)
	}
}

func TestSearchUsesIntent(t *testing.T) {
	eng, _, svc := openEngine(t)
	ctx := context.Background()

	const unique = "g9intentsearchtokenXYZ"
	task, err := svc.CreateTask(ctx, domain.TaskInput{
		Title: "G9 intent search",
		Body:  "Find " + unique + " in corpus",
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	dec, err := svc.CreateDecision(ctx, domain.DecisionInput{
		Title: "Decision about " + unique,
	})
	if err != nil {
		t.Fatalf("CreateDecision: %v", err)
	}
	_ = task

	in := retrieval.IntentInput{
		TaskTitle: "G9 intent search",
		TaskBody:  "Find " + unique + " in corpus",
	}
	rawHits, err := eng.Search(ctx, "", retrieval.SearchOptions{Limit: 16})
	if err != nil {
		t.Fatalf("Search raw empty: %v", err)
	}
	intentHits, err := eng.Search(ctx, "", retrieval.SearchOptions{Limit: 16, Intent: &in})
	if err != nil {
		t.Fatalf("Search intent: %v", err)
	}
	if len(intentHits) == 0 {
		t.Fatal("expected intent-enriched hits")
	}
	found := false
	for _, h := range intentHits {
		if h.EntityType == "decision" && h.EntityID == dec.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("intent search missed decision %s; hits=%+v raw=%+v", dec.ID, intentHits, rawHits)
	}
	if intentHits[0].ReasonCode != retrieval.ReasonFTSMatch {
		t.Fatalf("expected fts_match, got %q", intentHits[0].ReasonCode)
	}
}

func TestIntentNoSemantic(t *testing.T) {
	dir := filepath.Join("..", "..", "internal", "retrieval")
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") || name == "doc.go" {
			continue
		}
		files = append(files, filepath.Join(dir, name))
	}
	for _, path := range files {
		fset := token.NewFileSet()
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile %s: %v", path, err)
		}
		f, err := parser.ParseFile(fset, path, src, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		for _, imp := range f.Imports {
			impPath := strings.Trim(imp.Path.Value, `"`)
			if strings.Contains(impPath, "embedding") || strings.Contains(impPath, "vector") || strings.Contains(impPath, "openai") {
				t.Fatalf("%s imports forbidden %q", filepath.Base(path), impPath)
			}
		}
	}
	intent := retrieval.ExtractIntent(retrieval.IntentInput{TaskTitle: "lexical only"})
	if strings.Contains(intent.FTSQuery(), "semantic_match") {
		t.Fatal("FTSQuery must not emit semantic_match")
	}
}

func TestExtractIntentDeterministic(t *testing.T) {
	in := retrieval.IntentInput{
		TaskTitle: "Ship auth",
		TaskBody:  "Add JWT middleware internal/foo/bar.go",
		Query:     "extra token",
	}
	var baseline []byte
	for i := 0; i < 10; i++ {
		intent := retrieval.ExtractIntent(in)
		b, err := json.Marshal(intent)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if i == 0 {
			baseline = b
			continue
		}
		if string(b) != string(baseline) {
			t.Fatalf("call %d differed:\n%s\nvs\n%s", i, b, baseline)
		}
	}
}
