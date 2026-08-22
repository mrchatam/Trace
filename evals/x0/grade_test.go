package x0

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGradeAnswerCriticalMissAndCorrect(t *testing.T) {
	bank, err := LoadQueryBank(filepath.Join(moduleRoot(t), "evals", "x0", "queries.json"))
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]Query{}
	for _, q := range bank.Queries {
		byID[q.ID] = q
	}

	if g := GradeAnswer(byID["q1"], PackAnswer{Assert: []string{"11111111-1111-1111-1111-111111111111", "Ship greeter + math demo"}}); g != GradeCorrect {
		t.Fatalf("q1 want correct got %s", g)
	}
	if g := GradeAnswer(byID["q1"], PackAnswer{Assert: []string{"decision is the parent goal"}}); g != GradeCriticalMiss {
		t.Fatalf("q1 want critical_miss got %s", g)
	}
	if g := GradeAnswer(byID["q1"], PackAnswer{Text: "no goal found", Assert: nil}); g != GradeIncorrect {
		t.Fatalf("empty assert want incorrect got %s", g)
	}
	if g := GradeAnswer(byID["q3"], PackAnswer{Assert: []string{
		"44444444-4444-4444-4444-444444444444",
		"55555555-5555-5555-5555-555555555555",
		"discovery_causes_plan_change",
	}}); g != GradeCorrect {
		t.Fatalf("q3 want correct got %s", g)
	}
	if g := GradeAnswer(byID["q3"], PackAnswer{Assert: []string{"55555555-5555-5555-5555-555555555555 caused 44444444-4444-4444-4444-444444444444"}}); g != GradeCriticalMiss {
		t.Fatalf("q3 swap want critical_miss got %s", g)
	}
}

func TestGradePackB0AndG1Sample(t *testing.T) {
	root := moduleRoot(t)
	bank, err := LoadQueryBank(filepath.Join(root, "evals", "x0", "queries.json"))
	if err != nil {
		t.Fatal(err)
	}
	b0, err := LoadAnswerPack(filepath.Join(root, "evals", "x0", "testdata", "gate-c", "b0-run1.json"))
	if err != nil {
		t.Fatal(err)
	}
	g1, err := LoadAnswerPack(filepath.Join(root, "evals", "x0", "testdata", "gate-c", "g1-run1.json"))
	if err != nil {
		t.Fatal(err)
	}
	qb0 := GradePack(bank, b0)
	qg1 := GradePack(bank, g1)
	if qb0.UnderstandingAccuracy >= qg1.UnderstandingAccuracy {
		t.Fatalf("expected G1 sample accuracy > B0; B0=%v G1=%v", qb0, qg1)
	}
	if qg1.Correct < 4 {
		t.Fatalf("G1 sample expected ≥4 correct, got %d (%v)", qg1.Correct, qg1.PerQuery)
	}
	if _, err := QualityJSON(qg1); err != nil {
		t.Fatal(err)
	}
}

func TestQueryBankHasAtLeastFive(t *testing.T) {
	bank, err := LoadQueryBank(filepath.Join(moduleRoot(t), "evals", "x0", "queries.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(bank.Queries) < 5 {
		t.Fatalf("want ≥5 got %d", len(bank.Queries))
	}
	raw, err := os.ReadFile(filepath.Join(moduleRoot(t), "evals", "x0", "queries.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
}
