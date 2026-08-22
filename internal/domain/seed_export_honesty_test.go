package domain_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
)

func loadSeedDocumentFixture(t *testing.T, name string) domain.SeedDocument {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "cmd", "trace", "testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	var doc domain.SeedDocument
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parse fixture %s: %v", name, err)
	}
	return doc
}

func TestSeedDocumentHonestyThinGraphP26Snippet(t *testing.T) {
	doc := loadSeedDocumentFixture(t, "p26-export-snippet.json")
	violations := domain.CollectSeedDocumentHonestyViolations(doc)
	if len(violations) == 0 {
		t.Fatal("expected min-count violation for P26 thin graph")
	}
	found := false
	for _, v := range violations {
		if strings.Contains(v.Message, "graph honesty") &&
			strings.Contains(v.Message, "discoveries=0") &&
			strings.Contains(v.Message, "decisions=0") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("missing thin-graph violation: %+v", violations)
	}
}

func TestSeedDocumentHonestyOrphanDiscovery(t *testing.T) {
	doc := domain.SeedDocument{
		Version: 1,
		Goals:   []domain.SeedEntity{{ID: "g1", Title: "G"}},
		Tasks:   []domain.SeedTask{{ID: "t1", Title: "T", GoalID: "g1"}},
		Discoveries: []domain.SeedEntity{
			{ID: "disc-1", Title: "Gap found"},
		},
		Links: []domain.SeedLink{
			{Rel: domain.RelGoalHasTaskEvent, From: "g1", To: "t1"},
		},
	}
	violations := domain.CollectSeedDocumentHonestyViolations(doc)
	if len(violations) != 1 {
		t.Fatalf("want 1 violation got %d: %+v", len(violations), violations)
	}
	if !strings.Contains(violations[0].Message, "disc-1") ||
		!strings.Contains(violations[0].Message, "discovery_mentions_task") {
		t.Fatalf("unexpected message: %q", violations[0].Message)
	}
}

func TestSeedDocumentHonestyOrphanDecision(t *testing.T) {
	doc := domain.SeedDocument{
		Version: 1,
		Goals:   []domain.SeedEntity{{ID: "g1", Title: "G"}},
		Tasks:   []domain.SeedTask{{ID: "t1", Title: "T", GoalID: "g1"}},
		Decisions: []domain.SeedEntity{
			{ID: "dec-1", Title: "Use SQLite"},
		},
		Links: []domain.SeedLink{
			{Rel: domain.RelGoalHasTaskEvent, From: "g1", To: "t1"},
		},
	}
	violations := domain.CollectSeedDocumentHonestyViolations(doc)
	if len(violations) != 1 {
		t.Fatalf("want 1 violation got %d: %+v", len(violations), violations)
	}
	if !strings.Contains(violations[0].Message, "dec-1") ||
		!strings.Contains(violations[0].Message, "decision_affects_task") {
		t.Fatalf("unexpected message: %q", violations[0].Message)
	}
}

func TestSeedDocumentHonestyCleanWithLinkedDecision(t *testing.T) {
	doc := domain.SeedDocument{
		Version: 1,
		Goals:   []domain.SeedEntity{{ID: "g1", Title: "G"}},
		Tasks:   []domain.SeedTask{{ID: "t1", Title: "T", GoalID: "g1"}},
		Decisions: []domain.SeedEntity{
			{ID: "dec-1", Title: "Use SQLite"},
		},
		Links: []domain.SeedLink{
			{Rel: domain.RelGoalHasTaskEvent, From: "g1", To: "t1"},
			{Rel: domain.RelDecisionAffectsTask, From: "dec-1", To: "t1"},
		},
	}
	violations := domain.CollectSeedDocumentHonestyViolations(doc)
	if len(violations) != 0 {
		t.Fatalf("expected no violations: %+v", violations)
	}
}
