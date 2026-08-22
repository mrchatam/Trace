package loop_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

func openLoopTestStore(t *testing.T) (*store.Store, *planner.Service, *domain.Service) {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st, planner.New(st), domain.New(st)
}

func seedGoalTaskPlan(t *testing.T, psvc *planner.Service, dsvc *domain.Service) (goalID, taskID, scopeID string) {
	t.Helper()
	ctx := context.Background()
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "loop goal"})
	if err != nil {
		t.Fatalf("CreateGoal: %v", err)
	}
	task, err := dsvc.CreateTask(ctx, domain.TaskInput{Title: "seed", GoalID: &g.ID})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	cp, err := psvc.CreateCoarsePlan(ctx, planner.CoarsePlanInput{
		GoalID: g.ID,
		Phases: []planner.PhaseInput{{
			Title:  "Phase 1",
			Scopes: []planner.ScopeInput{{Title: "Scope 1"}},
		}},
	})
	if err != nil {
		t.Fatalf("CreateCoarsePlan: %v", err)
	}
	scopeID = cp.Phases[0].Scopes[0].ID
	if err := psvc.SetCurrentScope(ctx, g.ID, scopeID); err != nil {
		t.Fatalf("SetCurrentScope: %v", err)
	}
	if _, err := psvc.DeepPlan(ctx, planner.DeepPlanInput{
		ScopeID:      scopeID,
		ExitCriteria: []string{"packet ready"},
		WorkItems:    []planner.WorkItem{{Title: "emit loop packet"}},
	}); err != nil {
		t.Fatalf("DeepPlan: %v", err)
	}
	return g.ID, task.ID, scopeID
}

func countTransitionEvents(t *testing.T, st *store.Store, taskID string) int {
	t.Helper()
	evs, err := st.ListEventsByEntity(domain.EntityTask, taskID)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, e := range evs {
		if e.Type == domain.EventDeliberationTransition {
			n++
		}
	}
	return n
}

func TestLoopApplyTransactionalRollbackOnFailure(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, scopeID := seedGoalTaskPlan(t, psvc, dsvc)

	discID := uuid.NewString()
	pcID := uuid.NewString()
	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes: loop.ApplyWrites{
			Discoveries: []loop.ApplyDiscovery{{
				ID: discID, Title: "found gap",
			}},
			PlanChanges: []loop.ApplyPlanChange{{
				ID: pcID, Title: "replan attempt",
				Replan: &loop.ApplyReplan{
					ScopeID:      scopeID,
					ExitCriteria: []string{"should not commit"},
				},
			}},
		},
	}

	_, err := loop.Apply(ctx, st, psvc, env)
	if err == nil || !strings.Contains(err.Error(), "replan requires discovery_id") {
		t.Fatalf("want replan discovery_id error, got %v", err)
	}

	discoveries, err := st.ListDiscoveries()
	if err != nil {
		t.Fatal(err)
	}
	if len(discoveries) != 0 {
		t.Fatalf("rollback must leave zero discoveries, got %+v", discoveries)
	}
	if n := countTransitionEvents(t, st, taskID); n != 0 {
		t.Fatalf("rollback must leave zero transitions, got %d", n)
	}
}

func TestLoopApplyGoalIDMismatchFailsClosed(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	otherGoal, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "other"})
	if err != nil {
		t.Fatal(err)
	}

	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: otherGoal.ID},
		Writes:        loop.ApplyWrites{},
	}
	_, err = loop.Apply(ctx, st, psvc, env)
	if err == nil || !strings.Contains(err.Error(), "seed goal mismatch") {
		t.Fatalf("want seed goal mismatch, got %v", err)
	}
	_ = goalID

	if n := countTransitionEvents(t, st, taskID); n != 0 {
		t.Fatalf("must not write on goal mismatch, transitions=%d", n)
	}
}

func TestLoopApplyDeliberationTransitionEvent(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       "11111111-1111-4111-8111-111111111112",
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes:        loop.ApplyWrites{},
	}
	if _, err := loop.Apply(ctx, st, psvc, env); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	evs, err := st.ListEventsByEntity(domain.EntityTask, taskID)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, e := range evs {
		if e.Type != domain.EventDeliberationTransition {
			continue
		}
		found = true
		var payload map[string]any
		if err := json.Unmarshal([]byte(e.PayloadJSON), &payload); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"to_phase", "reason_code", "policy_inputs"} {
			if payload[key] == nil {
				t.Fatalf("transition payload missing %q: %#v", key, payload)
			}
		}
	}
	if !found {
		t.Fatal("missing deliberation.transition event")
	}
}

func TestLoopApplyNoPartialWritesOnValidationFailure(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       "33333333-3333-4333-8333-333333333334",
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes: loop.ApplyWrites{
			Discoveries: []loop.ApplyDiscovery{{
				ID: uuid.NewString(), Title: "ok",
			}},
			Uncertainties: []loop.ApplyUncertainty{{
				ID: uuid.NewString(), Title: "",
			}},
		},
	}
	if err := loop.ValidateApplyEnvelope(env); err == nil || !strings.Contains(err.Error(), "uncertainties[0].title is required") {
		t.Fatalf("want validation error, got %v", err)
	}

	discoveries, err := st.ListDiscoveries()
	if err != nil {
		t.Fatal(err)
	}
	if len(discoveries) != 0 {
		t.Fatalf("partial write occurred: %+v", discoveries)
	}
	uncCount, err := st.CountOpenBlockingUncertaintiesByTaskID(taskID)
	if err != nil {
		t.Fatal(err)
	}
	if uncCount != 0 {
		t.Fatalf("partial uncertainty write: count=%d", uncCount)
	}
	_ = ctx
	_ = psvc
}

func TestLoopApplyReplaySkipsDuplicateTransition(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       "22222222-2222-4222-8222-222222222223",
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes:        loop.ApplyWrites{},
	}
	if _, err := loop.Apply(ctx, st, psvc, env); err != nil {
		t.Fatalf("first apply: %v", err)
	}
	if n := countTransitionEvents(t, st, taskID); n != 1 {
		t.Fatalf("want 1 transition after first apply, got %d", n)
	}

	res, err := loop.Apply(ctx, st, psvc, env)
	if err != nil {
		t.Fatalf("replay apply: %v", err)
	}
	if !res.Replay {
		t.Fatalf("want replay=true, got %+v", res)
	}
	if n := countTransitionEvents(t, st, taskID); n != 1 {
		t.Fatalf("replay must not duplicate transition, got %d", n)
	}
}

func TestLoopApplyUnknownWriteKeyFailsClosed(t *testing.T) {
	raw := []byte(`{"schema_version":"trace.loop.apply.v1","apply_id":"cccccccc-cccc-4ccc-8ccc-cccccccccccc","seed":{"task_id":"dddddddd-dddd-4ddd-8ddd-dddddddddddd","goal_id":"eeeeeeee-eeee-4eee-8eee-eeeeeeeeeeee"},"writes":{"discoveries":[],"plan_changes":[],"spawned_tasks":[],"unknown_key":[]}}`)
	_, err := loop.ParseApplyEnvelope(raw)
	if err == nil || !strings.Contains(err.Error(), "unknown writes key") {
		t.Fatalf("want unknown key fail, got %v", err)
	}
}

func TestValidateApplyEnvelopeSpawnedTaskGoalMismatch(t *testing.T) {
	seedGoal := uuid.NewString()
	otherGoal := uuid.NewString()
	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: uuid.NewString(), GoalID: seedGoal},
		Writes: loop.ApplyWrites{
			SpawnedTasks: []loop.ApplySpawnedTask{{
				ID: uuid.NewString(), Title: "child", GoalID: otherGoal,
			}},
		},
	}
	err := loop.ValidateApplyEnvelope(env)
	if err == nil || !strings.Contains(err.Error(), "spawned_tasks[0].goal_id must match seed.goal_id") {
		t.Fatalf("want spawned goal mismatch, got %v", err)
	}
}

func TestLoopApplySuccessPersistsLoopStepEvent(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	applyID := uuid.NewString()
	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       applyID,
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes:        loop.ApplyWrites{},
	}
	if _, err := loop.Apply(ctx, st, psvc, env); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	evs, err := st.ListEventsByEntity(domain.EntityTask, taskID)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, e := range evs {
		if e.Type != loop.EventLoopStep {
			continue
		}
		found = true
		var step struct {
			ApplyID string         `json:"apply_id"`
			Seed    loop.ApplySeed `json:"seed"`
		}
		if err := json.Unmarshal([]byte(e.PayloadJSON), &step); err != nil {
			t.Fatal(err)
		}
		if step.ApplyID != applyID || step.Seed.TaskID != taskID || step.Seed.GoalID != goalID {
			t.Fatalf("loop step payload: %+v", step)
		}
	}
	if !found {
		t.Fatal("missing loop.step.applied event")
	}
}

func TestLoopApplyPromotesBlockingDiscoveryViaSpawnedTask(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	disc, err := dsvc.CreateDiscovery(ctx, domain.DiscoveryInput{
		Title:    "Blocking gap",
		Severity: domain.SeverityBlocking,
	})
	if err != nil {
		t.Fatal(err)
	}

	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes: loop.ApplyWrites{
			SpawnedTasks: []loop.ApplySpawnedTask{{
				ID:          uuid.NewString(),
				Title:       "placeholder",
				DiscoveryID: disc.ID,
			}},
		},
	}
	res, err := loop.Apply(ctx, st, psvc, env)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if res.NewSpawnedTasks != 1 {
		t.Fatalf("NewSpawnedTasks=%d want 1", res.NewSpawnedTasks)
	}
	if len(res.SpawnedTaskIDs) != 1 || res.SpawnedTaskIDs[0] != disc.ID {
		t.Fatalf("spawned_task_ids mismatch: %+v", res.SpawnedTaskIDs)
	}
	tasks, err := st.ListTasksByGoalID(goalID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, task := range tasks {
		if task.ID == disc.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("promoted task %q not found in goal roster", disc.ID)
	}
	links, err := st.ListLinksFrom(domain.EntityDiscovery, disc.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) == 0 || links[0].Rel != domain.RelDiscoveryMentionsTask || links[0].ToID != disc.ID {
		t.Fatalf("discovery link missing: %+v", links)
	}
}

func TestLoopApplyDiscoveriesOnlyDoNotCreateTasks(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)

	before, err := st.ListTasksByGoalID(goalID)
	if err != nil {
		t.Fatal(err)
	}
	env := loop.ApplyEnvelope{
		SchemaVersion: loop.ApplySchemaVersion,
		ApplyID:       uuid.NewString(),
		Seed:          loop.ApplySeed{TaskID: taskID, GoalID: goalID},
		Writes: loop.ApplyWrites{
			Discoveries: []loop.ApplyDiscovery{{
				ID:    uuid.NewString(),
				Title: "Only discovery",
			}},
		},
	}
	res, err := loop.Apply(ctx, st, psvc, env)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if res.NewSpawnedTasks != 0 || len(res.SpawnedTaskIDs) != 0 {
		t.Fatalf("discoveries-only should not spawn tasks: %+v", res)
	}
	after, err := st.ListTasksByGoalID(goalID)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Fatalf("task roster changed from discoveries-only apply: before=%d after=%d", len(before), len(after))
	}
}
