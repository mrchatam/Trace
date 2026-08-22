package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestKnowledgeListAndSynthesize(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	task, err := svc.CreateTask(context.Background(), domain.TaskInput{Title: "knowledge CLI"})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	dec, err := svc.CreateDecision(context.Background(), domain.DecisionInput{Title: "use sqlite"})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if _, err := svc.RecordDecisionReconsideration(context.Background(), dec.ID, domain.ReconsiderationInput{
		Trigger: store.ReconsiderTriggerNewEvidence,
		Status:  store.ReconsiderStatusOpen,
		Reason:  "new benchmark data",
	}); err != nil {
		st.Close()
		t.Fatal(err)
	}
	chg, err := svc.CreateChange(context.Background(), domain.ChangeInput{
		TaskID: task.ID,
		Reason: "internal tweak",
		Paths:  []domain.ChangePathInput{{Path: "internal/x.go", Status: "M"}},
	})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	if _, err := svc.RecordImprovement(context.Background(), domain.ImprovementInput{
		ChangeID: chg.ID,
		TaskID:   task.ID,
		Summary:  "cache layer helped",
	}); err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	outSynth := captureStdout(t, func() int {
		return run([]string{"-C", dir, "knowledge", "synthesize"})
	})
	var synthResp knowledgeSynthesizeResponse
	if err := json.Unmarshal([]byte(strings.TrimSpace(outSynth)), &synthResp); err != nil {
		t.Fatalf("synthesize json: %v body=%q", err, outSynth)
	}
	if !synthResp.OK || synthResp.Created < 1 {
		t.Fatalf("synthesize: %+v", synthResp)
	}

	outList := captureStdout(t, func() int {
		return run([]string{"-C", dir, "knowledge", "list"})
	})
	var rows []store.EngineeringKnowledge
	if err := json.Unmarshal([]byte(strings.TrimSpace(outList)), &rows); err != nil {
		t.Fatalf("list json: %v body=%q", err, outList)
	}
	if len(rows) == 0 {
		t.Fatalf("knowledge list empty: %q", outList)
	}

	outTopic := captureStdout(t, func() int {
		return run([]string{"-C", dir, "knowledge", "list", "--topic", "decision"})
	})
	var decisionRows []store.EngineeringKnowledge
	if err := json.Unmarshal([]byte(strings.TrimSpace(outTopic)), &decisionRows); err != nil {
		t.Fatalf("topic json: %v body=%q", err, outTopic)
	}
	if len(decisionRows) == 0 {
		t.Fatalf("decision topic empty: %q", outTopic)
	}
}

func TestHelpIncludesKnowledge(t *testing.T) {
	out := captureStdout(t, func() int { return run([]string{"help"}) })
	if !strings.Contains(out, "knowledge list") || !strings.Contains(out, "knowledge synthesize") {
		t.Fatalf("help missing knowledge commands: %q", out)
	}
	if !strings.Contains(out, "knowledge tendencies") {
		t.Fatalf("help missing knowledge tendencies: %q", out)
	}
}

func TestKnowledgeTendenciesCLI(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	task, err := svc.CreateTask(context.Background(), domain.TaskInput{Title: "tend CLI"})
	if err != nil {
		st.Close()
		t.Fatal(err)
	}
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		chg, err := svc.CreateChange(ctx, domain.ChangeInput{
			TaskID: task.ID,
			Reason: "internal improvement",
			Paths:  []domain.ChangePathInput{{Path: "internal/a.go", Status: "M"}},
		})
		if err != nil {
			st.Close()
			t.Fatal(err)
		}
		if _, err := svc.RecordImprovement(ctx, domain.ImprovementInput{
			ChangeID: chg.ID, TaskID: task.ID, Summary: "faster",
		}); err != nil {
			st.Close()
			t.Fatal(err)
		}
	}
	if _, err := svc.RefreshChangePatterns(ctx); err != nil {
		st.Close()
		t.Fatal(err)
	}
	st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "knowledge", "tendencies"})
	})
	var rows []domain.TendencyRow
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &rows); err != nil {
		t.Fatalf("tendencies json: %v body=%q", err, out)
	}
	if len(rows) == 0 {
		t.Fatalf("tendencies empty: %q", out)
	}
}
