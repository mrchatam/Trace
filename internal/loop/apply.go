package loop

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	ApplySchemaVersion  = "trace.loop.apply.v1"
	StatusSchemaVersion = "trace.loop.status.v1"
	EventLoopStep       = "loop.step.applied"
)

type ApplyEnvelope struct {
	SchemaVersion string      `json:"schema_version"`
	ApplyID       string      `json:"apply_id"`
	Seed          ApplySeed   `json:"seed"`
	Writes        ApplyWrites `json:"writes"`
}

type ApplySeed struct {
	TaskID string `json:"task_id"`
	GoalID string `json:"goal_id"`
}

type ApplyWrites struct {
	Discoveries   []ApplyDiscovery    `json:"discoveries"`
	PlanChanges   []ApplyPlanChange   `json:"plan_changes"`
	SpawnedTasks  []ApplySpawnedTask  `json:"spawned_tasks"`
	Stop          *ApplyStop          `json:"stop,omitempty"`
	Uncertainties []ApplyUncertainty  `json:"uncertainties"`
	Hypotheses    []ApplyHypothesis   `json:"hypotheses"`
	Changes       []ApplyChange       `json:"changes"`
	Effects       []ApplyEffect       `json:"effects"`
	TestResults   []ApplyTestResult   `json:"test_results"`
	Verifications []ApplyVerification `json:"verifications"`
	Evaluations   []ApplyEvaluation   `json:"evaluations"`
	Regressions   []ApplyRegression   `json:"regressions"`
	Reflections   []ApplyReflection   `json:"reflections"`
}

type ApplyUncertainty struct {
	ID       string      `json:"id"`
	Title    string      `json:"title"`
	Body     string      `json:"body,omitempty"`
	Severity string      `json:"severity,omitempty"`
	Status   string      `json:"status,omitempty"`
	Kind     string      `json:"kind,omitempty"`
	TaskID   string      `json:"task_id,omitempty"`
	Links    []ApplyLink `json:"links,omitempty"`
}

type ApplyHypothesis struct {
	ID            string      `json:"id"`
	Title         string      `json:"title"`
	Body          string      `json:"body,omitempty"`
	EvidenceIDs   []string    `json:"evidence_ids,omitempty"`
	UncertaintyID string      `json:"uncertainty_id,omitempty"`
	Links         []ApplyLink `json:"links,omitempty"`
}

type ApplyChange struct {
	ID             string          `json:"id"`
	GitCommit      string          `json:"git_commit,omitempty"`
	ParentChangeID string          `json:"parent_change_id,omitempty"`
	Actor          string          `json:"actor,omitempty"`
	Reason         string          `json:"reason,omitempty"`
	Paths          []string        `json:"paths"`
	Expected       []ApplyExpected `json:"expected,omitempty"`
	DecisionID     string          `json:"decision_id,omitempty"`
}

type ApplyExpected struct {
	Dimension string `json:"dimension"`
	Expected  string `json:"expected"`
}

type ApplyEffect struct {
	ChangeID     string   `json:"change_id"`
	Dimension    string   `json:"dimension"`
	Expected     string   `json:"expected,omitempty"`
	Actual       string   `json:"actual,omitempty"`
	Comparison   string   `json:"comparison,omitempty"`
	EvidenceIDs  []string `json:"evidence_ids,omitempty"`
	HypothesisID string   `json:"hypothesis_id,omitempty"`
}

type ApplyTestResult struct {
	ID          string   `json:"id"`
	TestName    string   `json:"test_name"`
	TestStatus  string   `json:"test_status"`
	Summary     string   `json:"summary,omitempty"`
	EvidenceIDs []string `json:"evidence_ids,omitempty"`
}

type ApplyVerification struct {
	ID                 string   `json:"id"`
	GoalID             string   `json:"goal_id"`
	VerificationStatus string   `json:"verification_status"`
	EvidenceIDs        []string `json:"evidence_ids"`
	Summary            string   `json:"summary,omitempty"`
}

type ApplyEvaluation struct {
	ID         string `json:"id"`
	BaselineID string `json:"baseline_id"`
	ScoresJSON string `json:"scores_json"`
}

type ApplyRegression struct {
	SourceKind string `json:"source_kind"`
	SourceID   string `json:"source_id"`
	TaskID     string `json:"task_id,omitempty"`
	Summary    string `json:"summary,omitempty"`
}

type ApplyReflection struct {
	Summary                  string            `json:"summary,omitempty"`
	InvalidatedAssumptionIDs []string          `json:"invalidated_assumption_ids,omitempty"`
	NewDependencies          []ApplyDependency `json:"new_dependencies,omitempty"`
	UsefulTests              []string          `json:"useful_tests,omitempty"`
	BroadenTestsNote         string            `json:"broaden_tests_note,omitempty"`
}

type ApplyDependency struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}

type ApplyDiscovery struct {
	ID       string      `json:"id"`
	Title    string      `json:"title"`
	Body     string      `json:"body,omitempty"`
	Severity string      `json:"severity,omitempty"`
	Links    []ApplyLink `json:"links,omitempty"`
}

type ApplyPlanChange struct {
	ID          string       `json:"id"`
	Title       string       `json:"title,omitempty"`
	Body        string       `json:"body,omitempty"`
	DiscoveryID string       `json:"discovery_id,omitempty"`
	Replan      *ApplyReplan `json:"replan,omitempty"`
}

type ApplyReplan struct {
	ScopeID          string             `json:"scope_id"`
	ExitCriteria     []string           `json:"exit_criteria,omitempty"`
	Constraints      []string           `json:"constraints,omitempty"`
	WorkItems        []planner.WorkItem `json:"work_items,omitempty"`
	LookaheadScopeID string             `json:"lookahead_scope_id,omitempty"`
	LookaheadSummary string             `json:"lookahead_summary,omitempty"`
	MaxAutoReplans   int                `json:"max_auto_replans,omitempty"`
}

type ApplySpawnedTask struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Body        string      `json:"body,omitempty"`
	GoalID      string      `json:"goal_id,omitempty"`
	DiscoveryID string      `json:"discovery_id,omitempty"`
	Links       []ApplyLink `json:"links,omitempty"`
}

type ApplyStop struct {
	Reason               string `json:"reason"`
	MaxIterationsReached bool   `json:"max_iterations_reached,omitempty"`
}

type ApplyLink struct {
	Rel    string `json:"rel"`
	ToType string `json:"to_type"`
	ToID   string `json:"to_id"`
}

type ApplyResult struct {
	NewDiscoveries  int      `json:"new_discoveries"`
	NewPlanChanges  int      `json:"new_plan_changes"`
	NewSpawnedTasks int      `json:"new_spawned_tasks"`
	SpawnedTaskIDs  []string `json:"spawned_task_ids"`
	Saturated       bool     `json:"saturated"`
	Replay          bool     `json:"replay"`
}

type Advisory struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type StatusResult struct {
	SchemaVersion               string                   `json:"schema_version"`
	Seed                        ApplySeed                `json:"seed"`
	LastApplyID                 string                   `json:"last_apply_id"`
	NewPlanChangesSinceLastStep int                      `json:"new_plan_changes_since_last_step"`
	NewTasksSinceLastStep       int                      `json:"new_tasks_since_last_step"`
	MaxIterationsReached        bool                     `json:"max_iterations_reached"`
	Saturated                   bool                     `json:"saturated"`
	Reason                      string                   `json:"reason"`
	Deliberation                *StatusDeliberation      `json:"deliberation,omitempty"`
	PromotionBlocked            *PromotionBlocked        `json:"promotion_blocked,omitempty"`
	VerificationCycle           *VerificationCycleStatus `json:"verification_cycle,omitempty"`
	Violations                  []Violation              `json:"violations"`
	Advisories                  []Advisory               `json:"advisories"`
}

// VerificationCycleStatus mirrors policy cycle flags plus human incomplete_reason (C11).
type VerificationCycleStatus struct {
	ExecutePending         bool   `json:"execute_pending"`
	TestPending            bool   `json:"test_pending"`
	VerificationIncomplete bool   `json:"verification_incomplete"`
	EvaluationPending      bool   `json:"evaluation_pending"`
	ReflectPending         bool   `json:"reflect_pending"`
	IncompleteReason       string `json:"incomplete_reason,omitempty"`
	RegressionDetected     bool   `json:"regression_detected,omitempty"`
	TestName               string `json:"test_name,omitempty"`
}

type PromotionBlocked struct {
	Present bool   `json:"present"`
	Reason  string `json:"reason,omitempty"`
}

type loopStepPayload struct {
	ApplyID                 string    `json:"apply_id"`
	Seed                    ApplySeed `json:"seed"`
	NewDiscoveries          int       `json:"new_discoveries"`
	NewPlanChanges          int       `json:"new_plan_changes"`
	NewSpawnedTasks         int       `json:"new_spawned_tasks"`
	SpawnedTaskIDs          []string  `json:"spawned_task_ids"`
	MaxIterationsReached    bool      `json:"max_iterations_reached"`
	Saturated               bool      `json:"saturated"`
	ConsecutiveEmptyApplies int       `json:"consecutive_empty_applies"`
}

func ParseApplyEnvelope(raw []byte) (ApplyEnvelope, error) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return ApplyEnvelope{}, fmt.Errorf("loop apply: parse envelope: %w", err)
	}
	required := []string{"schema_version", "apply_id", "seed", "writes"}
	for _, key := range required {
		if _, ok := top[key]; !ok {
			return ApplyEnvelope{}, fmt.Errorf("loop apply: missing required field %q", key)
		}
	}
	var env ApplyEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return ApplyEnvelope{}, fmt.Errorf("loop apply: parse envelope: %w", err)
	}
	if err := validateWritesKeys(top["writes"]); err != nil {
		return ApplyEnvelope{}, err
	}
	if err := ValidateApplyEnvelope(env); err != nil {
		return ApplyEnvelope{}, err
	}
	return env, nil
}

func ValidateApplyEnvelope(env ApplyEnvelope) error {
	if env.SchemaVersion != ApplySchemaVersion {
		return fmt.Errorf("loop apply: schema_version must be %q", ApplySchemaVersion)
	}
	if err := requireUUID("apply_id", env.ApplyID); err != nil {
		return fmt.Errorf("loop apply: %w", err)
	}
	if err := requireUUID("seed.task_id", env.Seed.TaskID); err != nil {
		return fmt.Errorf("loop apply: %w", err)
	}
	if err := requireUUID("seed.goal_id", env.Seed.GoalID); err != nil {
		return fmt.Errorf("loop apply: %w", err)
	}
	for i, d := range env.Writes.Discoveries {
		if err := requireUUID(fmt.Sprintf("writes.discoveries[%d].id", i), d.ID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
		if strings.TrimSpace(d.Title) == "" {
			return fmt.Errorf("loop apply: writes.discoveries[%d].title is required", i)
		}
		for j, l := range d.Links {
			if err := validateLink(fmt.Sprintf("writes.discoveries[%d].links[%d]", i, j), l); err != nil {
				return fmt.Errorf("loop apply: %w", err)
			}
		}
	}
	for i, p := range env.Writes.PlanChanges {
		if err := requireUUID(fmt.Sprintf("writes.plan_changes[%d].id", i), p.ID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
		if p.DiscoveryID != "" {
			if err := requireUUID(fmt.Sprintf("writes.plan_changes[%d].discovery_id", i), p.DiscoveryID); err != nil {
				return fmt.Errorf("loop apply: %w", err)
			}
		}
		if p.Replan != nil {
			if err := requireUUID(fmt.Sprintf("writes.plan_changes[%d].replan.scope_id", i), p.Replan.ScopeID); err != nil {
				return fmt.Errorf("loop apply: %w", err)
			}
		}
	}
	for i, t := range env.Writes.SpawnedTasks {
		if err := requireUUID(fmt.Sprintf("writes.spawned_tasks[%d].id", i), t.ID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
		if strings.TrimSpace(t.DiscoveryID) != "" {
			if err := requireUUID(fmt.Sprintf("writes.spawned_tasks[%d].discovery_id", i), t.DiscoveryID); err != nil {
				return fmt.Errorf("loop apply: %w", err)
			}
		}
		if strings.TrimSpace(t.Title) == "" && strings.TrimSpace(t.DiscoveryID) == "" {
			return fmt.Errorf("loop apply: writes.spawned_tasks[%d].title is required", i)
		}
		if strings.TrimSpace(t.GoalID) != "" && strings.TrimSpace(t.GoalID) != env.Seed.GoalID {
			return fmt.Errorf("loop apply: writes.spawned_tasks[%d].goal_id must match seed.goal_id", i)
		}
		for j, l := range t.Links {
			if err := validateLink(fmt.Sprintf("writes.spawned_tasks[%d].links[%d]", i, j), l); err != nil {
				return fmt.Errorf("loop apply: %w", err)
			}
		}
	}
	for i, u := range env.Writes.Uncertainties {
		if err := requireUUID(fmt.Sprintf("writes.uncertainties[%d].id", i), u.ID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
		if strings.TrimSpace(u.Title) == "" && strings.TrimSpace(u.Status) != store.UncertaintyStatusResolved {
			return fmt.Errorf("loop apply: writes.uncertainties[%d].title is required", i)
		}
	}
	for i, h := range env.Writes.Hypotheses {
		if err := requireUUID(fmt.Sprintf("writes.hypotheses[%d].id", i), h.ID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
		if strings.TrimSpace(h.Title) == "" {
			return fmt.Errorf("loop apply: writes.hypotheses[%d].title is required", i)
		}
	}
	for i, c := range env.Writes.Changes {
		if err := requireUUID(fmt.Sprintf("writes.changes[%d].id", i), c.ID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
		if len(c.Paths) < 1 {
			return fmt.Errorf("loop apply: writes.changes[%d].paths is required", i)
		}
	}
	for i, e := range env.Writes.Effects {
		if strings.TrimSpace(e.ChangeID) == "" {
			return fmt.Errorf("loop apply: writes.effects[%d].change_id is required", i)
		}
		if err := requireUUID(fmt.Sprintf("writes.effects[%d].change_id", i), e.ChangeID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
		if strings.TrimSpace(e.Dimension) == "" {
			return fmt.Errorf("loop apply: writes.effects[%d].dimension is required", i)
		}
	}
	for i, tr := range env.Writes.TestResults {
		if err := requireUUID(fmt.Sprintf("writes.test_results[%d].id", i), tr.ID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
		if strings.TrimSpace(tr.TestName) == "" {
			return fmt.Errorf("loop apply: writes.test_results[%d].test_name is required", i)
		}
	}
	for i, v := range env.Writes.Verifications {
		if err := requireUUID(fmt.Sprintf("writes.verifications[%d].id", i), v.ID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
		if len(v.EvidenceIDs) < 1 {
			return fmt.Errorf("loop apply: writes.verifications[%d].evidence_ids is required", i)
		}
	}
	for i, ev := range env.Writes.Evaluations {
		if err := requireUUID(fmt.Sprintf("writes.evaluations[%d].id", i), ev.ID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
		if err := requireUUID(fmt.Sprintf("writes.evaluations[%d].baseline_id", i), ev.BaselineID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
	}
	for i, r := range env.Writes.Regressions {
		if err := requireUUID(fmt.Sprintf("writes.regressions[%d].source_id", i), r.SourceID); err != nil {
			return fmt.Errorf("loop apply: %w", err)
		}
	}
	return nil
}

func Apply(ctx context.Context, st *store.Store, plan *planner.Service, env ApplyEnvelope) (ApplyResult, error) {
	if st == nil {
		return ApplyResult{}, fmt.Errorf("loop apply: store is required")
	}
	task, err := st.GetTask(env.Seed.TaskID)
	if err != nil {
		return ApplyResult{}, fmt.Errorf("loop apply: load seed task: %w", err)
	}
	if task.GoalID == nil || *task.GoalID != env.Seed.GoalID {
		return ApplyResult{}, fmt.Errorf("loop apply: seed goal mismatch for task %q", env.Seed.TaskID)
	}
	if prev, ok := findLoopStep(st, env.Seed, env.ApplyID); ok {
		return ApplyResult{
			NewDiscoveries:  prev.NewDiscoveries,
			NewPlanChanges:  prev.NewPlanChanges,
			NewSpawnedTasks: prev.NewSpawnedTasks,
			SpawnedTaskIDs:  append([]string(nil), prev.SpawnedTaskIDs...),
			Saturated:       prev.MaxIterationsReached || prev.Saturated,
			Replay:          true,
		}, nil
	}

	out := ApplyResult{SpawnedTaskIDs: []string{}}
	err = st.WithTx(func(txSt *store.Store) error {
		txDom := domain.New(txSt)
		var txPlan *planner.Service
		if plan != nil {
			txPlan = planner.New(txSt)
		}
		for _, d := range env.Writes.Discoveries {
			disc, inserted, err := txDom.ImportSeedDiscovery(ctx, domain.SeedEntity{
				ID: d.ID, Title: strings.TrimSpace(d.Title), Body: d.Body,
			})
			if err != nil {
				return err
			}
			if inserted {
				out.NewDiscoveries++
			}
			if strings.TrimSpace(d.Severity) != "" {
				if err := txDom.SetDiscoverySeverity(ctx, disc.ID, d.Severity); err != nil {
					return err
				}
			}
			for _, link := range d.Links {
				if err := importSupportedLink(ctx, txDom, domain.EntityDiscovery, d.ID, link); err != nil {
					return err
				}
			}
		}
		for _, p := range env.Writes.PlanChanges {
			title := strings.TrimSpace(p.Title)
			if title == "" {
				title = "Plan change " + p.ID
			}
			_, inserted, err := txDom.ImportSeedPlanChange(ctx, domain.SeedEntity{
				ID: p.ID, Title: title, Body: p.Body,
			})
			if err != nil {
				return err
			}
			if inserted {
				out.NewPlanChanges++
			}
			if strings.TrimSpace(p.DiscoveryID) != "" {
				if err := txDom.ImportSeedLink(ctx, domain.SeedLink{
					Rel:  "discovery-plan-change",
					From: p.DiscoveryID,
					To:   p.ID,
				}); err != nil {
					return err
				}
			}
			if p.Replan != nil {
				if txPlan == nil {
					return fmt.Errorf("loop apply: planner is required for replan")
				}
				if strings.TrimSpace(p.DiscoveryID) == "" {
					return fmt.Errorf("loop apply: plan_change %q replan requires discovery_id", p.ID)
				}
				_, err := txPlan.ApplyDiscoveryReplan(ctx, planner.ApplyDiscoveryReplanInput{
					DiscoveryID:      p.DiscoveryID,
					ScopeID:          p.Replan.ScopeID,
					PlanChangeID:     p.ID,
					PlanChangeTitle:  title,
					PlanChangeBody:   p.Body,
					ExitCriteria:     p.Replan.ExitCriteria,
					Constraints:      p.Replan.Constraints,
					WorkItems:        p.Replan.WorkItems,
					LookaheadScopeID: p.Replan.LookaheadScopeID,
					LookaheadSummary: p.Replan.LookaheadSummary,
					MaxAutoReplans:   p.Replan.MaxAutoReplans,
					Actor:            "loop.apply",
				})
				if err != nil {
					return err
				}
			}
		}
		for _, t := range env.Writes.SpawnedTasks {
			gid := env.Seed.GoalID
			if strings.TrimSpace(t.GoalID) != "" {
				gid = t.GoalID
			}
			if strings.TrimSpace(t.DiscoveryID) != "" {
				taskID, inserted, err := txDom.PromoteBlockingDiscovery(ctx, t.DiscoveryID, gid)
				if err != nil {
					return err
				}
				if inserted {
					out.NewSpawnedTasks++
				}
				out.SpawnedTaskIDs = append(out.SpawnedTaskIDs, taskID)
				for _, link := range t.Links {
					if err := importSupportedLink(ctx, txDom, domain.EntityTask, taskID, link); err != nil {
						return err
					}
				}
				continue
			}
			task, inserted, err := txDom.ImportSeedTask(ctx, domain.SeedTask{
				ID: t.ID, GoalID: gid, Title: strings.TrimSpace(t.Title), Body: t.Body,
			}, &gid)
			if err != nil {
				return err
			}
			if inserted {
				out.NewSpawnedTasks++
			}
			out.SpawnedTaskIDs = append(out.SpawnedTaskIDs, task.ID)
			for _, link := range t.Links {
				if err := importSupportedLink(ctx, txDom, domain.EntityTask, task.ID, link); err != nil {
					return err
				}
			}
		}
		maxReached := env.Writes.Stop != nil && env.Writes.Stop.MaxIterationsReached
		discoveryWrites := len(env.Writes.Discoveries)
		prevConsec := 0
		if existing, err := txSt.GetDeliberationState(env.Seed.TaskID); err == nil {
			prevConsec = existing.ConsecutiveEmptyApplies
		} else if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("loop apply: load deliberation state: %w", err)
		}
		consec := deliberation.NextConsecutiveEmptyApplies(
			prevConsec, out.NewPlanChanges, out.NewSpawnedTasks, discoveryWrites,
		)
		out.Saturated = deliberation.SaturatedFromCounter(maxReached, consec)
		// Discoveries-only applies are non-saturating for this transition (counter unchanged).
		if !maxReached && discoveryWrites > 0 && out.NewPlanChanges == 0 && out.NewSpawnedTasks == 0 {
			out.Saturated = false
		}
		out.Replay = false

		if err := applyCognitiveWrites(ctx, txDom, txSt, env.Seed, env.Writes); err != nil {
			return err
		}

		// Persist consecutive counter before transition so ApplyTransition preserves it.
		dRow, err := ensureDeliberationRow(txSt, env.Seed.TaskID, env.Seed.GoalID)
		if err != nil {
			return err
		}
		dRow.ConsecutiveEmptyApplies = consec
		if _, err := txSt.UpsertDeliberationState(dRow); err != nil {
			return fmt.Errorf("loop apply: persist consecutive empty: %w", err)
		}

		inputs, err := BuildPolicyInputs(ctx, txDom, txPlan, env.Seed.TaskID, env.Seed.GoalID, &env.Writes, out.Saturated)
		if err != nil {
			return err
		}
		if _, _, err := txDom.ApplyDeliberationTransition(ctx, env.Seed.TaskID, env.Seed.GoalID, inputs); err != nil {
			return fmt.Errorf("loop apply: deliberation transition: %w", err)
		}

		step := loopStepPayload{
			ApplyID:                 env.ApplyID,
			Seed:                    env.Seed,
			NewDiscoveries:          out.NewDiscoveries,
			NewPlanChanges:          out.NewPlanChanges,
			NewSpawnedTasks:         out.NewSpawnedTasks,
			SpawnedTaskIDs:          out.SpawnedTaskIDs,
			MaxIterationsReached:    maxReached,
			Saturated:               out.Saturated,
			ConsecutiveEmptyApplies: consec,
		}
		raw, err := json.Marshal(step)
		if err != nil {
			return fmt.Errorf("loop apply: encode step summary: %w", err)
		}
		if _, err := txSt.AppendEvent(store.Event{
			TS:          time.Now().UTC().Format(time.RFC3339Nano),
			Type:        EventLoopStep,
			EntityType:  domain.EntityTask,
			EntityID:    env.Seed.TaskID,
			PayloadJSON: string(raw),
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return ApplyResult{}, err
	}
	return out, nil
}

func Status(ctx context.Context, st *store.Store, plan *planner.Service, seed ApplySeed) (StatusResult, error) {
	if st == nil {
		return StatusResult{}, fmt.Errorf("loop status: store is required")
	}
	if plan == nil {
		return StatusResult{}, fmt.Errorf("loop status: planner is required")
	}
	if err := requireUUID("seed.task_id", seed.TaskID); err != nil {
		return StatusResult{}, fmt.Errorf("loop status: %w", err)
	}
	task, err := st.GetTask(seed.TaskID)
	if err != nil {
		return StatusResult{}, fmt.Errorf("loop status: load seed task: %w", err)
	}
	if task.GoalID == nil || *task.GoalID == "" {
		return StatusResult{}, fmt.Errorf("loop status: task %q has no goal_id", seed.TaskID)
	}
	if strings.TrimSpace(seed.GoalID) == "" {
		seed.GoalID = *task.GoalID
	}
	if *task.GoalID != seed.GoalID {
		return StatusResult{}, fmt.Errorf("loop status: seed goal mismatch for task %q", seed.TaskID)
	}

	step, ok := latestLoopStep(st, seed)
	promotionBlocked, err := buildPromotionBlocked(ctx, st, seed.TaskID)
	if err != nil {
		return StatusResult{}, err
	}
	if !ok {
		delib, err := buildStatusDeliberation(ctx, st, plan, seed, false)
		if err != nil {
			return StatusResult{}, err
		}
		dom := domain.New(st)
		cycle, err := buildVerificationCycle(ctx, dom, delib.PolicyInputs, seed.TaskID)
		if err != nil {
			return StatusResult{}, err
		}
		return attachStatusViolations(ctx, st, plan, seed, StatusResult{
			SchemaVersion:     StatusSchemaVersion,
			Seed:              seed,
			Reason:            "insufficient_history",
			Saturated:         false,
			Deliberation:      delib,
			PromotionBlocked:  promotionBlocked,
			VerificationCycle: cycle,
		})
	}
	maxReached := step.MaxIterationsReached
	consec := step.ConsecutiveEmptyApplies
	if existing, err := st.GetDeliberationState(seed.TaskID); err == nil {
		consec = existing.ConsecutiveEmptyApplies
	}
	saturated := deliberation.SaturatedFromCounter(maxReached, consec)
	reason := "tasks_and_plan_unchanged"
	if maxReached {
		reason = "max_iterations_reached"
	}
	if !saturated {
		reason = "insufficient_history"
	}
	delib, err := buildStatusDeliberation(ctx, st, plan, seed, saturated)
	if err != nil {
		return StatusResult{}, err
	}
	dom := domain.New(st)
	cycle, err := buildVerificationCycle(ctx, dom, delib.PolicyInputs, seed.TaskID)
	if err != nil {
		return StatusResult{}, err
	}
	return attachStatusViolations(ctx, st, plan, seed, StatusResult{
		SchemaVersion:               StatusSchemaVersion,
		Seed:                        seed,
		LastApplyID:                 step.ApplyID,
		NewPlanChangesSinceLastStep: step.NewPlanChanges,
		NewTasksSinceLastStep:       step.NewSpawnedTasks,
		MaxIterationsReached:        maxReached,
		Saturated:                   saturated,
		Reason:                      reason,
		Deliberation:                delib,
		PromotionBlocked:            promotionBlocked,
		VerificationCycle:           cycle,
	})
}

func attachStatusViolations(ctx context.Context, st *store.Store, plan *planner.Service, seed ApplySeed, res StatusResult) (StatusResult, error) {
	dom := domain.New(st)
	allowed, violations, err := EvaluateGate(ctx, dom, plan, st, seed.TaskID, GateForEdit)
	if err != nil {
		return StatusResult{}, fmt.Errorf("loop status: gate: %w", err)
	}
	if allowed || violations == nil {
		violations = []Violation{}
	}
	res.Violations = violations
	res.Advisories = []Advisory{}
	pAdvs, err := plan.StatusAdvisories(ctx, seed.GoalID)
	if err != nil {
		return StatusResult{}, fmt.Errorf("loop status: advisories: %w", err)
	}
	for _, a := range pAdvs {
		res.Advisories = append(res.Advisories, Advisory{Code: a.Code, Message: a.Message})
	}
	return res, nil
}

func buildPromotionBlocked(ctx context.Context, st *store.Store, taskID string) (*PromotionBlocked, error) {
	dom := domain.New(st)
	allowed, reason, err := dom.CheckPromotionGate(ctx, taskID, "")
	if err != nil {
		return nil, err
	}
	out := &PromotionBlocked{Present: !allowed}
	if !allowed {
		out.Reason = reason
	}
	return out, nil
}

func buildStatusDeliberation(ctx context.Context, st *store.Store, plan *planner.Service, seed ApplySeed, p19Sat bool) (*StatusDeliberation, error) {
	dom := domain.New(st)
	inputs, err := BuildPolicyInputs(ctx, dom, plan, seed.TaskID, seed.GoalID, nil, p19Sat)
	if err != nil {
		return nil, err
	}
	dState := loadDeliberationState(ctx, dom, seed.TaskID, seed.GoalID)
	recommended, reason, stopped := deliberation.SelectNext(dState, inputs)
	blocked := statusBlocked(inputs)
	needs := statusNeedsPhase(recommended, blocked, stopped || dState.Stopped, p19Sat)
	out := &StatusDeliberation{
		Phase:            dState.CurrentPhase,
		RecommendedPhase: recommended,
		WhySelected:      reason,
		HopCount:         dState.HopCount,
		Stopped:          stopped || dState.Stopped,
		Blocked:          blocked,
		NeedsPhase:       needs,
		PolicyInputs:     inputs,
	}
	if dState.Stopped {
		out.StopReason = dState.StopReason
	} else if stopped {
		out.StopReason = string(reason)
	}
	return out, nil
}

func latestLoopStep(st *store.Store, seed ApplySeed) (loopStepPayload, bool) {
	evs, err := st.ListEventsByEntity(domain.EntityTask, seed.TaskID)
	if err != nil {
		return loopStepPayload{}, false
	}
	for i := len(evs) - 1; i >= 0; i-- {
		ev := evs[i]
		if ev.Type != EventLoopStep {
			continue
		}
		var step loopStepPayload
		if json.Unmarshal([]byte(ev.PayloadJSON), &step) != nil {
			continue
		}
		if step.Seed.GoalID != seed.GoalID {
			continue
		}
		return step, true
	}
	return loopStepPayload{}, false
}

func ensureDeliberationRow(st *store.Store, taskID, goalID string) (store.DeliberationState, error) {
	existing, err := st.GetDeliberationState(taskID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return store.DeliberationState{}, err
	}
	return store.DeliberationState{
		TaskID:       taskID,
		GoalID:       goalID,
		CurrentPhase: string(deliberation.PhaseOrient),
	}, nil
}

func findLoopStep(st *store.Store, seed ApplySeed, applyID string) (loopStepPayload, bool) {
	evs, err := st.ListEventsByEntity(domain.EntityTask, seed.TaskID)
	if err != nil {
		return loopStepPayload{}, false
	}
	for _, ev := range evs {
		if ev.Type != EventLoopStep {
			continue
		}
		var step loopStepPayload
		if json.Unmarshal([]byte(ev.PayloadJSON), &step) != nil {
			continue
		}
		if step.Seed.GoalID == seed.GoalID && step.ApplyID == applyID {
			return step, true
		}
	}
	return loopStepPayload{}, false
}

func validateLink(prefix string, link ApplyLink) error {
	if strings.TrimSpace(link.Rel) == "" {
		return fmt.Errorf("%s.rel is required", prefix)
	}
	if strings.TrimSpace(link.ToType) == "" {
		return fmt.Errorf("%s.to_type is required", prefix)
	}
	if err := requireUUID(prefix+".to_id", link.ToID); err != nil {
		return err
	}
	return nil
}

func requireUUID(field, v string) error {
	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("%s is required", field)
	}
	if _, err := uuid.Parse(v); err != nil {
		return fmt.Errorf("%s must be UUID", field)
	}
	return nil
}

func importSupportedLink(ctx context.Context, dom *domain.Service, fromType, fromID string, link ApplyLink) error {
	switch fromType {
	case domain.EntityDiscovery:
		if link.Rel == domain.RelDiscoveryMentionsTask && link.ToType == domain.EntityTask {
			return dom.ImportSeedLink(ctx, domain.SeedLink{
				Rel:  "discovery-mentions-task",
				From: fromID,
				To:   link.ToID,
			})
		}
		if link.Rel == domain.RelDiscoveryCausesPlanChange && link.ToType == domain.EntityPlanChange {
			return dom.ImportSeedLink(ctx, domain.SeedLink{
				Rel:  "discovery-plan-change",
				From: fromID,
				To:   link.ToID,
			})
		}
	case domain.EntityTask:
		if link.Rel == domain.RelGoalHasTaskEvent && link.ToType == domain.EntityGoal {
			return dom.ImportSeedLink(ctx, domain.SeedLink{
				Rel:  "goal-task",
				From: link.ToID,
				To:   fromID,
			})
		}
	}
	return fmt.Errorf("loop apply: unsupported link from %s rel=%s to_type=%s", fromType, link.Rel, link.ToType)
}
