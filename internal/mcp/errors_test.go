package mcp_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	tracemcp "github.com/mrchatam/Trace/internal/mcp"
	"github.com/mrchatam/Trace/internal/store"
)

// TestLinkErrorsNameRelVocabulary is the #155 proof. Previously
// `rel=goal-task from=<task-id>` reported the goal lookup as the failure,
// which blames the wrong parameter and leaked the driver text.
func TestLinkErrorsNameRelVocabulary(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)

	goalRes, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{Kind: "goal", Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	taskRes, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{Kind: "task", Title: "t"})
	if err != nil {
		t.Fatal(err)
	}
	var goal, task struct {
		ID string `json:"id"`
	}
	mustDecodeID(t, tracemcp.ResultText(t, goalRes), &goal)
	mustDecodeID(t, tracemcp.ResultText(t, taskRes), &task)

	// Wrong type in the `from` position: must name `from`, not the goal lookup.
	_, _, err = tracemcp.CallLink(srv, ctx, tracemcp.LinkInput{
		Rel: "goal-task", From: task.ID, To: goal.ID,
	})
	assertNoDriverText(t, "wrong-type from", err)
	for _, want := range []string{"from", "rel goal-task", "expects from=goal to=task"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("wrong-type from: want %q in %v", want, err)
		}
	}

	// Missing id in the `to` position: must name `to`.
	_, _, err = tracemcp.CallLink(srv, ctx, tracemcp.LinkInput{
		Rel: "goal-task", From: goal.ID, To: "00000000-0000-0000-0000-0000000000ab",
	})
	assertNoDriverText(t, "missing to", err)
	if err == nil || !strings.Contains(err.Error(), "to") {
		t.Fatalf("missing to: error must name the to parameter, got %v", err)
	}
}

// TestToolErrorsNeverLeakDriverText sweeps the tools named in #155.
func TestToolErrorsNeverLeakDriverText(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	const missing = "00000000-0000-0000-0000-0000000000cd"

	cases := []struct {
		name string
		call func() error
	}{
		{"trace_why", func() error {
			_, _, err := tracemcp.CallWhy(srv, ctx, tracemcp.WhyInput{EntityType: "task", ID: missing})
			return err
		}},
		{"trace_changes show", func() error {
			_, _, err := tracemcp.CallChanges(srv, ctx, tracemcp.ChangesInput{
				Action: "show", ChangeID: missing,
			})
			return err
		}},
		{"trace_review", func() error {
			_, _, err := tracemcp.CallReview(srv, ctx, tracemcp.ReviewInput{
				Action: "set", ID: missing, Result: "PASS", Reason: "checking",
			})
			return err
		}},
		{"trace_transition", func() error {
			_, _, err := tracemcp.CallTransition(srv, ctx, tracemcp.TransitionInput{
				TaskID: missing, ToState: "IN_PROGRESS", Reason: "starting",
			})
			return err
		}},
		{"trace_link", func() error {
			_, _, err := tracemcp.CallLink(srv, ctx, tracemcp.LinkInput{
				Rel: "blocks", From: missing, To: "00000000-0000-0000-0000-0000000000ce",
			})
			return err
		}},
		// S02 typed errors made most tools honest at the source; these cases
		// pin the rest of the surface so a future unconverted getter leaks
		// here instead of to a model.
		{"trace_context", func() error {
			_, _, err := tracemcp.CallContext(srv, ctx, tracemcp.ContextInput{TaskID: missing})
			return err
		}},
		{"trace_capability missing", func() error {
			_, _, err := tracemcp.CallCapability(srv, ctx, tracemcp.CapabilityInput{
				Action: "missing", TaskID: missing,
			})
			return err
		}},
		{"trace_capability require", func() error {
			_, _, err := tracemcp.CallCapability(srv, ctx, tracemcp.CapabilityInput{
				Action: "require", TaskID: missing, Capability: "some-slug",
			})
			return err
		}},
		{"trace_plan deep", func() error {
			_, _, err := tracemcp.CallPlan(srv, ctx, tracemcp.PlanInput{
				Action: "deep", GoalID: missing,
			})
			return err
		}},
		{"trace_impact finding list", func() error {
			_, _, err := tracemcp.CallImpact(srv, ctx, tracemcp.ImpactInput{
				Action: "finding", Op: "list", Decision: missing,
			})
			return err
		}},
		{"trace_impact walk", func() error {
			_, _, err := tracemcp.CallImpact(srv, ctx, tracemcp.ImpactInput{
				Action: "walk", Decision: missing, Depth: 2,
			})
			return err
		}},
		{"trace_loop status", func() error {
			_, _, err := tracemcp.CallLoop(srv, ctx, tracemcp.LoopInput{
				Action: "status", TaskID: missing,
			})
			return err
		}},
		// trace_loop gate on a missing task intentionally returns an
		// allowed:false gate envelope (orientReasonTaskNotFound), not an error —
		// the envelope IS the honest answer. Covered by the loop tests.
		{"trace_regressions list", func() error {
			_, _, err := tracemcp.CallRegressions(srv, ctx, tracemcp.RegressionsInput{
				Action: "list", ChangeID: missing,
			})
			return err
		}},
		{"trace_changes compare", func() error {
			_, _, err := tracemcp.CallChanges(srv, ctx, tracemcp.ChangesInput{
				Action: "compare", From: "0000000000000000000000000000000000000000", To: "0000000000000000000000000000000000000001",
			})
			return err
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			if err == nil {
				t.Fatal("want an error for a missing entity")
			}
			assertNoDriverText(t, tc.name, err)
			if !strings.Contains(err.Error(), tc.name) {
				t.Fatalf("error must name the tool, got %v", err)
			}
		})
	}
}

// TestLoopNextReturnsGateEnvelopeOnPrecondition is the #156 proof: a blocked
// action=next answers with the same structured envelope as action=gate.
func TestLoopNextReturnsGateEnvelopeOnPrecondition(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)

	// A task with no goal cannot produce a packet.
	taskRes, _, err := tracemcp.CallAdd(srv, ctx, tracemcp.AddInput{Kind: "task", Title: "unlinked"})
	if err != nil {
		t.Fatal(err)
	}
	var task struct {
		ID string `json:"id"`
	}
	mustDecodeID(t, tracemcp.ResultText(t, taskRes), &task)

	res, _, err := tracemcp.CallLoop(srv, ctx, tracemcp.LoopInput{Action: "next", TaskID: task.ID})
	if err != nil {
		t.Fatalf("a precondition must not be a protocol error: %v", err)
	}
	var env struct {
		SchemaVersion string `json:"schema_version"`
		For           string `json:"for"`
		Allowed       bool   `json:"allowed"`
		ReasonCode    string `json:"reason_code"`
		Remedy        string `json:"remedy"`
		Violations    []struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Remedy  string `json:"remedy"`
		} `json:"violations"`
	}
	if err := json.Unmarshal([]byte(tracemcp.ResultText(t, res)), &env); err != nil {
		t.Fatalf("next must return the gate envelope as JSON: %v", err)
	}
	if env.SchemaVersion != "trace.loop.gate.v1" {
		t.Fatalf("schema_version = %q, want the gate envelope version", env.SchemaVersion)
	}
	if env.For != "next" || env.Allowed {
		t.Fatalf("want a blocked next envelope, got %+v", env)
	}
	if env.ReasonCode != "goal_link_missing" {
		t.Fatalf("reason_code = %q", env.ReasonCode)
	}
	if env.Remedy == "" || env.Violations[0].Remedy == "" {
		t.Fatalf("a blocked envelope must carry a remedy: %+v", env)
	}
}

// TestLoopNextUnknownTaskReturnsGateEnvelope covers the not-found precondition.
func TestLoopNextUnknownTaskReturnsGateEnvelope(t *testing.T) {
	dir, ctx := writeFixture(t)
	srv := addServer(t, dir)
	res, _, err := tracemcp.CallLoop(srv, ctx, tracemcp.LoopInput{
		Action: "next", TaskID: "00000000-0000-0000-0000-0000000000ff",
	})
	if err != nil {
		t.Fatalf("a missing seed must not be a protocol error: %v", err)
	}
	var env struct {
		ReasonCode string `json:"reason_code"`
		Remedy     string `json:"remedy"`
	}
	if err := json.Unmarshal([]byte(tracemcp.ResultText(t, res)), &env); err != nil {
		t.Fatalf("next must return the gate envelope as JSON: %v", err)
	}
	if env.ReasonCode != "seed_task_not_found" || env.Remedy == "" {
		t.Fatalf("want a seed_task_not_found envelope with a remedy, got %+v", env)
	}
}

func mustDecodeID(t *testing.T, raw string, out any) {
	t.Helper()
	if err := json.Unmarshal([]byte(raw), out); err != nil {
		t.Fatalf("decode id: %v\n%s", err, raw)
	}
}

// assertNoDriverText fails when raw SQLite text leaks into a user-facing error.
func assertNoDriverText(t *testing.T, label string, err error) {
	t.Helper()
	if err == nil {
		return
	}
	for _, leak := range []string{"sql: no rows", "no rows in result set", "sql.ErrNoRows"} {
		if strings.Contains(err.Error(), leak) {
			t.Fatalf("%s leaked driver text %q: %v", label, leak, err)
		}
	}
}

// seedLinkedTask creates a goal and a task linked to it, then releases the
// store lock so the MCP server can open the same project.
func seedLinkedTask(t *testing.T, dir string) string {
	t.Helper()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	goal, err := svc.CreateGoal(context.Background(), domain.GoalInput{Title: "g"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(context.Background(), domain.TaskInput{Title: "t", GoalID: &goal.ID})
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	return task.ID
}
