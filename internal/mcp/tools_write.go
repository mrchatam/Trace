package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/config"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

// AddInput mirrors `trace add <kind> --title …`.
type AddInput struct {
	Project    string  `json:"project,omitempty" jsonschema:"optional project root override"`
	Kind       string  `json:"kind" jsonschema:"discovery|task|goal|decision|assumption|plan-change|claim|evidence"`
	Title      string  `json:"title" jsonschema:"entity title (required)"`
	Body       string  `json:"body,omitempty" jsonschema:"entity body"`
	ID         string  `json:"id,omitempty" jsonschema:"optional UUID; must not already exist unless replace=true"`
	Replace    bool    `json:"replace,omitempty" jsonschema:"allow overwriting an existing entity with the same id (default false; refuse rather than silently overwrite)"`
	SourceType string  `json:"source_type,omitempty" jsonschema:"provenance source_type"`
	Confidence float64 `json:"confidence,omitempty" jsonschema:"provenance confidence"`
	Status     string  `json:"status,omitempty" jsonschema:"provenance status"`
	GoalID     string  `json:"goal_id,omitempty" jsonschema:"goal id (task only); must resolve to an existing goal"`
}

// LinkInput mirrors `trace link <rel> --from --to`.
type LinkInput struct {
	Project    string `json:"project,omitempty" jsonschema:"optional project root override"`
	Rel        string `json:"rel" jsonschema:"goal-task|decision-task|discovery-plan-change|discovery-mentions-task|claim-evidence|scope-member|scope_member|api-contract|api_contract|implements|blocks"`
	From       string `json:"from" jsonschema:"from entity UUID"`
	To         string `json:"to" jsonschema:"to entity UUID"`
	SourceType string `json:"source_type,omitempty" jsonschema:"optional source_type"`
}

// TransitionInput mirrors `trace transition --task --to --reason`.
type TransitionInput struct {
	Project          string   `json:"project,omitempty" jsonschema:"optional project root override"`
	TaskID           string   `json:"task_id" jsonschema:"task UUID"`
	ToState          string   `json:"to_state" jsonschema:"target work_state"`
	Actor            string   `json:"actor,omitempty" jsonschema:"actor; default mcp (not authorization)"`
	Reason           string   `json:"reason" jsonschema:"reason (required non-empty)"`
	AllowDone        bool     `json:"allow_done,omitempty" jsonschema:"AllowDoneWithoutReview escape hatch (emits warning); does not bypass missing caps (use allow_missing_caps)"`
	AsOperator       bool     `json:"as_operator,omitempty" jsonschema:"AllowOperatorDone; conscious claim flag≠identity / not verified operator identity; required with Review PASS; Actor ≠ auth"`
	AllowMissingCaps bool     `json:"allow_missing_caps,omitempty" jsonschema:"AllowMissingCapabilities override"`
	EvidenceIDs      []string `json:"evidence_ids,omitempty" jsonschema:"evidence ids (do not alone authorize DONE)"`
	Enforce          bool     `json:"enforce,omitempty" jsonschema:"run GateForDone before DONE (also honored via .trace/config.json enforce warn|strict)"`
}

// ReviewInput mirrors `trace review create|set`.
type ReviewInput struct {
	Project    string  `json:"project,omitempty" jsonschema:"optional project root override"`
	Action     string  `json:"action" jsonschema:"create|set"`
	Title      string  `json:"title,omitempty" jsonschema:"review title (create)"`
	Body       string  `json:"body,omitempty" jsonschema:"review body (create)"`
	ID         string  `json:"id,omitempty" jsonschema:"review UUID (create optional / set required)"`
	TaskID     string  `json:"task_id,omitempty" jsonschema:"optional task to link on create (review_judges_task)"`
	SourceType string  `json:"source_type,omitempty" jsonschema:"provenance source_type (create)"`
	Confidence float64 `json:"confidence,omitempty" jsonschema:"provenance confidence (create)"`
	Result     string  `json:"result,omitempty" jsonschema:"PASS|FAIL|UNCERTAIN (set)"`
	Actor      string  `json:"actor,omitempty" jsonschema:"actor for set; default mcp"`
	Reason     string  `json:"reason,omitempty" jsonschema:"reason for set (required)"`
}

func (s *Server) toolAdd(ctx context.Context, _ *sdkmcp.CallToolRequest, in AddInput) (*sdkmcp.CallToolResult, any, error) {
	if strings.TrimSpace(in.Title) == "" {
		return nil, nil, fmt.Errorf("trace_add: title is required")
	}
	kind, err := domain.NormalizeAddKind(in.Kind)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_add: %w", err)
	}
	entityID, err := domain.ValidateEntityID("id", in.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_add: %w", err)
	}

	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_add"); err != nil {
		return nil, nil, err
	}
	svc := domain.New(st)

	// Fail closed on an existing id: an explicit id is a claim about identity,
	// and the upsert underneath would silently destroy the row (#150).
	created := true
	if entityID != "" {
		exists, err := svc.EntityExists(ctx, kind, entityID)
		if err != nil {
			return nil, nil, fmt.Errorf("trace_add: %w", err)
		}
		if exists {
			if !in.Replace {
				return nil, nil, fmt.Errorf(
					"trace_add: %s %q already exists; pass replace=true to overwrite it (or omit id to create a new entity)",
					kind, entityID)
			}
			created = false
		}
	}

	// A task's goal_id has no FK, so it must be resolved before the write (#151).
	goalID := strings.TrimSpace(in.GoalID)
	if kind == domain.AddKindTask && goalID != "" {
		if _, err := domain.ValidateEntityID("goal_id", goalID); err != nil {
			return nil, nil, fmt.Errorf("trace_add: %w", err)
		}
		if _, err := svc.GetGoal(ctx, goalID); err != nil {
			return nil, nil, fmt.Errorf("trace_add: goal_id %q: %w", goalID, err)
		}
	}

	var outID, outType string
	switch kind {
	case domain.AddKindGoal:
		g, err := svc.CreateGoal(ctx, domain.GoalInput{
			ID: entityID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("trace_add: %w", err)
		}
		outID, outType = g.ID, domain.EntityGoal
	case domain.AddKindTask:
		var gid *string
		if goalID != "" {
			g := goalID
			gid = &g
		}
		t, err := svc.CreateTask(ctx, domain.TaskInput{
			ID: entityID, GoalID: gid, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("trace_add: %w", err)
		}
		outID, outType = t.ID, domain.EntityTask
	case domain.AddKindDecision:
		d, err := svc.CreateDecision(ctx, domain.DecisionInput{
			ID: in.ID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("trace_add: %w", err)
		}
		outID, outType = d.ID, domain.EntityDecision
	case domain.AddKindAssumption:
		a, err := svc.CreateAssumption(ctx, domain.AssumptionInput{
			ID: entityID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("trace_add: %w", err)
		}
		outID, outType = a.ID, domain.EntityAssumption
	case domain.AddKindDiscovery:
		d, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{
			ID: entityID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("trace_add: %w", err)
		}
		outID, outType = d.ID, domain.EntityDiscovery
	case domain.AddKindPlanChange:
		p, err := svc.CreatePlanChange(ctx, domain.PlanChangeInput{
			ID: entityID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("trace_add: %w", err)
		}
		outID, outType = p.ID, domain.EntityPlanChange
	case domain.AddKindClaim:
		c, err := svc.CreateClaim(ctx, domain.ClaimInput{
			ID: entityID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("trace_add: %w", err)
		}
		outID, outType = c.ID, domain.EntityClaim
	case domain.AddKindEvidence:
		e, err := svc.CreateEvidence(ctx, domain.EvidenceInput{
			ID: entityID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("trace_add: %w", err)
		}
		outID, outType = e.ID, domain.EntityEvidence
	default:
		return nil, nil, fmt.Errorf("trace_add: unknown kind %q", in.Kind)
	}

	b, err := json.Marshal(map[string]any{"ok": true, "type": outType, "id": outID, "created": created})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) toolLink(ctx context.Context, _ *sdkmcp.CallToolRequest, in LinkInput) (*sdkmcp.CallToolResult, any, error) {
	if in.Rel == "" || in.From == "" || in.To == "" {
		return nil, nil, fmt.Errorf("trace_link: rel, from, and to are required")
	}
	// A self-link is never meaningful and silently produces a cycle in every
	// rel that is read back as a graph edge (#162).
	if strings.TrimSpace(in.From) == strings.TrimSpace(in.To) {
		return nil, nil, fmt.Errorf(
			"trace_link: rel %s cannot link an entity to itself (from == to == %q)", in.Rel, strings.TrimSpace(in.From))
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_link"); err != nil {
		return nil, nil, err
	}
	svc := domain.New(st)
	meta := domain.LinkMeta{SourceType: in.SourceType}

	switch in.Rel {
	case "goal-task":
		err = svc.LinkGoalTask(ctx, in.From, in.To, meta)
	case "decision-task":
		err = svc.LinkDecisionTask(ctx, in.From, in.To, meta)
	case "discovery-plan-change":
		err = svc.LinkDiscoveryPlanChange(ctx, in.From, in.To, meta)
	case "discovery-mentions-task":
		err = svc.LinkDiscoveryMentionsTask(ctx, in.From, in.To, meta)
	case "claim-evidence":
		err = svc.LinkClaimEvidence(ctx, in.From, in.To, meta)
	case "scope-member", "scope_member":
		err = svc.LinkScopeMember(ctx, in.From, in.To, meta)
	case "api-contract", "api_contract":
		err = svc.LinkAPIContract(ctx, in.From, in.To, meta)
	case "implements":
		err = svc.LinkImplements(ctx, in.From, in.To, meta)
	case "blocks":
		err = svc.LinkBlocks(ctx, in.From, in.To, meta)
	default:
		return nil, nil, fmt.Errorf("trace_link: unknown rel %q", in.Rel)
	}
	if err != nil {
		return nil, nil, wrapLinkErr(in.Rel, in.From, in.To, err)
	}
	b, err := json.Marshal(map[string]any{"ok": true, "rel": in.Rel, "from": in.From, "to": in.To})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) toolTransition(ctx context.Context, _ *sdkmcp.CallToolRequest, in TransitionInput) (*sdkmcp.CallToolResult, any, error) {
	if in.TaskID == "" || in.ToState == "" || strings.TrimSpace(in.Reason) == "" {
		return nil, nil, fmt.Errorf("trace_transition: task_id, to_state, and reason are required")
	}
	actor := in.Actor
	if actor == "" {
		actor = "mcp"
	}
	abs, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_transition"); err != nil {
		return nil, nil, err
	}
	svc := domain.New(st)

	load := config.LoadEnforceModeDetail(abs)
	var warnings []string
	if load.Invalid && load.Warning != "" {
		warnings = append(warnings, load.Warning)
	}
	toDone := strings.EqualFold(strings.TrimSpace(in.ToState), store.WorkStateDone)
	runGate, rejectOnFail := config.DoneGateAction(in.Enforce, load.Mode)
	if toDone && runGate {
		plan := planner.New(st)
		allowed, violations, gateErr := loop.EvaluateGate(ctx, svc, plan, st, in.TaskID, loop.GateForDone)
		if gateErr != nil {
			return nil, nil, fmt.Errorf("trace_transition: gate: %w", gateErr)
		}
		if !allowed {
			if rejectOnFail {
				msg := "gate blocked"
				if len(violations) > 0 {
					msg = violations[0].Message
				}
				return nil, nil, fmt.Errorf("trace_transition: %s", msg)
			}
			for _, v := range violations {
				warnings = append(warnings, v.Message)
			}
		}
	}

	err = svc.TransitionTask(ctx, in.TaskID, in.ToState, domain.TransitionOptions{
		Actor:                    actor,
		Reason:                   in.Reason,
		EvidenceIDs:              in.EvidenceIDs,
		AllowDoneWithoutReview:   in.AllowDone,
		AllowOperatorDone:        in.AsOperator,
		AllowMissingCapabilities: in.AllowMissingCaps,
	})
	if err != nil {
		return nil, nil, wrapToolErr("trace_transition", err)
	}
	out := map[string]any{"ok": true, "task": in.TaskID, "to": in.ToState}
	if in.AllowDone {
		out["warning"] = "allow_done escape hatch used; Review PASS and as_operator were bypassed; missing capabilities still need allow_missing_caps"
	}
	if len(warnings) > 0 {
		out["warnings"] = warnings
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) toolReview(ctx context.Context, _ *sdkmcp.CallToolRequest, in ReviewInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	if err := assertMCPToolAllowed(ctx, st, "trace_review"); err != nil {
		_ = st.Close()
		return nil, nil, err
	}
	if err := st.Close(); err != nil {
		return nil, nil, err
	}
	switch strings.ToLower(in.Action) {
	case "create":
		return s.reviewCreate(ctx, in)
	case "set":
		return s.reviewSet(ctx, in)
	default:
		return nil, nil, fmt.Errorf("trace_review: action must be create|set")
	}
}

func (s *Server) reviewCreate(ctx context.Context, in ReviewInput) (*sdkmcp.CallToolResult, any, error) {
	if strings.TrimSpace(in.Title) == "" {
		return nil, nil, fmt.Errorf("trace_review create: title is required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	// Validate the task link target before creating anything, so a bogus
	// task_id cannot leave an orphaned review behind while the tool reports
	// an error (the caller would retry and create duplicates).
	if taskID := strings.TrimSpace(in.TaskID); taskID != "" {
		if _, err := svc.GetTask(ctx, taskID); err != nil {
			return nil, nil, wrapToolErr("trace_review create", err)
		}
	}

	r, err := svc.CreateReview(ctx, domain.ReviewInput{
		ID: in.ID, Title: in.Title, Body: in.Body,
		SourceType: in.SourceType, Confidence: in.Confidence,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_review: %w", err)
	}
	if in.TaskID != "" {
		if err := svc.LinkReviewTask(ctx, r.ID, in.TaskID, domain.LinkMeta{}); err != nil {
			return nil, nil, fmt.Errorf("trace_review: link: %w", err)
		}
	}
	out := map[string]any{"ok": true, "id": r.ID, "type": domain.EntityReview, "result": r.Result}
	if in.TaskID != "" {
		out["task"] = in.TaskID
		out["rel"] = domain.RelReviewJudgesTask
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) reviewSet(ctx context.Context, in ReviewInput) (*sdkmcp.CallToolResult, any, error) {
	if in.ID == "" || in.Result == "" || strings.TrimSpace(in.Reason) == "" {
		return nil, nil, fmt.Errorf("trace_review set: id, result, and reason are required")
	}
	actor := in.Actor
	if actor == "" {
		actor = "mcp"
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	result := strings.ToUpper(in.Result)
	err = svc.SetReviewResult(ctx, in.ID, result, domain.ReviewResultOptions{
		Actor: actor, Reason: in.Reason,
	})
	if err != nil {
		return nil, nil, wrapToolErr("trace_review", err)
	}
	b, err := json.Marshal(map[string]any{"ok": true, "id": in.ID, "result": result})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}
