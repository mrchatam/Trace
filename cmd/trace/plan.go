package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

// cmdPlan is a thin adapter over internal/planner (G19).
func cmdPlan(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace plan create-coarse|set-current|deep|show|bootstrap|apply-discovery|ack-replan …\n")
		return exitUsage
	}
	sub := args[0]
	switch sub {
	case "create-coarse":
		return cmdPlanCreateCoarse(root, args[1:])
	case "set-current":
		return cmdPlanSetCurrent(root, args[1:])
	case "deep":
		return cmdPlanDeep(root, args[1:])
	case "show":
		return cmdPlanShow(root, args[1:])
	case "bootstrap":
		return cmdPlanBootstrap(root, args[1:])
	case "apply-discovery":
		return cmdPlanApplyDiscovery(root, args[1:])
	case "ack-replan":
		return cmdPlanAckReplan(root, args[1:])
	case "help", "-h", "--help":
		printPlanHelp(os.Stdout)
		return exitOK
	default:
		fmt.Fprintf(os.Stderr, "plan: unknown subcommand %q\n", sub)
		printPlanHelp(os.Stderr)
		return exitUsage
	}
}

func printPlanHelp(w *os.File) {
	fmt.Fprint(w, `trace plan — progressive coarse planner (library: internal/planner)

Subcommands:
  create-coarse --goal <id> --phase <title> [--scope <title> ...]
                Repeatable --phase; each --phase consumes following --scope
                flags until the next --phase. Caller-supplied structure only.
  set-current   --goal <id> --scope <id>
  deep          --scope <id> --exit <text> [--constraint <text>] [--work <title>]
                [--lookahead <summary>]
                Requires scope == current (fail closed).
  show          --goal <id>
                JSON plan view on stdout (snake_case; phases always []; includes goal tasks);
                progress on stderr; goal_structure_warning when task count > 15 without plan.
  bootstrap     --goal <id>
                Recover progressive plan from plan_changes (idempotent when plan exists).
                Bootstrap yields a minimal plan; refine with create-coarse / deep afterward.
  apply-discovery --discovery <id> --scope <id>
                [--plan-change <id>] [--pc-title <t>] [--pc-body <b>]
                [--exit <text>] [--constraint <text>] [--work <title>] …
                Maps to ApplyDiscoveryReplan (severity + churn budget in library).
  ack-replan    --scope <id>
                Reset auto_replan_count to 0 (human ack).

MCP: trace_plan action=create-coarse|set-current|deep|show|bootstrap. No LLM backlog generation.
`)
}

func cmdPlanCreateCoarse(root string, args []string) int {
	goalID, phases, err := parseCreateCoarseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan create-coarse: %v\n", err)
		return exitUsage
	}
	svc, st, code := openPlanner(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	fmt.Fprintf(os.Stderr, "plan: creating coarse plan for goal %s (%d phases)\n", goalID, len(phases))
	cp, err := svc.CreateCoarsePlan(context.Background(), planner.CoarsePlanInput{
		GoalID: goalID, Phases: phases, Actor: "cli",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan create-coarse: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(cp)
	return exitOK
}

// parseCreateCoarseArgs parses --goal and ordered --phase/--scope pairs.
func parseCreateCoarseArgs(args []string) (goalID string, phases []planner.PhaseInput, err error) {
	var cur *planner.PhaseInput
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--goal" || a == "-goal":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--goal requires a value")
			}
			goalID = args[i+1]
			i++
		case strings.HasPrefix(a, "--goal="):
			goalID = strings.TrimPrefix(a, "--goal=")
		case a == "--phase" || a == "-phase":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--phase requires a title")
			}
			if cur != nil {
				phases = append(phases, *cur)
			}
			cur = &planner.PhaseInput{Title: args[i+1]}
			i++
		case strings.HasPrefix(a, "--phase="):
			if cur != nil {
				phases = append(phases, *cur)
			}
			cur = &planner.PhaseInput{Title: strings.TrimPrefix(a, "--phase=")}
		case a == "--scope" || a == "-scope":
			if cur == nil {
				return "", nil, fmt.Errorf("--scope before --phase")
			}
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--scope requires a title")
			}
			cur.Scopes = append(cur.Scopes, planner.ScopeInput{Title: args[i+1]})
			i++
		case strings.HasPrefix(a, "--scope="):
			if cur == nil {
				return "", nil, fmt.Errorf("--scope before --phase")
			}
			cur.Scopes = append(cur.Scopes, planner.ScopeInput{Title: strings.TrimPrefix(a, "--scope=")})
		case a == "-h" || a == "--help":
			return "", nil, fmt.Errorf("help")
		default:
			return "", nil, fmt.Errorf("unknown arg %q", a)
		}
	}
	if cur != nil {
		phases = append(phases, *cur)
	}
	if strings.TrimSpace(goalID) == "" {
		return "", nil, fmt.Errorf("usage: trace plan create-coarse --goal <id> --phase <title> [--scope <title> ...]")
	}
	if len(phases) == 0 {
		return "", nil, fmt.Errorf("at least one --phase is required")
	}
	return goalID, phases, nil
}

func cmdPlanSetCurrent(root string, args []string) int {
	fs := flag.NewFlagSet("plan set-current", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	goal := fs.String("goal", "", "goal UUID")
	scope := fs.String("scope", "", "scope UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *goal == "" || *scope == "" {
		fmt.Fprintf(os.Stderr, "usage: trace plan set-current --goal <id> --scope <id>\n")
		return exitUsage
	}
	svc, st, code := openPlanner(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	fmt.Fprintf(os.Stderr, "plan: setting current scope %s for goal %s\n", *scope, *goal)
	if err := svc.SetCurrentScope(context.Background(), *goal, *scope); err != nil {
		fmt.Fprintf(os.Stderr, "plan set-current: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": true, "goal": *goal, "scope": *scope})
	return exitOK
}

func cmdPlanDeep(root string, args []string) int {
	scopeID, exits, constraints, works, lookahead, err := parseDeepArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan deep: %v\n", err)
		return exitUsage
	}
	svc, st, code := openPlanner(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	items := make([]planner.WorkItem, 0, len(works))
	for _, w := range works {
		items = append(items, planner.WorkItem{Title: w})
	}
	fmt.Fprintf(os.Stderr, "plan: deep-planning scope %s\n", scopeID)
	res, err := svc.DeepPlan(context.Background(), planner.DeepPlanInput{
		ScopeID:          scopeID,
		ExitCriteria:     exits,
		Constraints:      constraints,
		WorkItems:        items,
		LookaheadSummary: lookahead,
		Actor:            "cli",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan deep: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(res)
	return exitOK
}

func parseDeepArgs(args []string) (scope string, exits, constraints, works []string, lookahead string, err error) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--scope" || a == "-scope":
			if i+1 >= len(args) {
				return "", nil, nil, nil, "", fmt.Errorf("--scope requires a value")
			}
			scope = args[i+1]
			i++
		case strings.HasPrefix(a, "--scope="):
			scope = strings.TrimPrefix(a, "--scope=")
		case a == "--exit" || a == "-exit":
			if i+1 >= len(args) {
				return "", nil, nil, nil, "", fmt.Errorf("--exit requires a value")
			}
			exits = append(exits, args[i+1])
			i++
		case strings.HasPrefix(a, "--exit="):
			exits = append(exits, strings.TrimPrefix(a, "--exit="))
		case a == "--constraint" || a == "-constraint":
			if i+1 >= len(args) {
				return "", nil, nil, nil, "", fmt.Errorf("--constraint requires a value")
			}
			constraints = append(constraints, args[i+1])
			i++
		case strings.HasPrefix(a, "--constraint="):
			constraints = append(constraints, strings.TrimPrefix(a, "--constraint="))
		case a == "--work" || a == "-work":
			if i+1 >= len(args) {
				return "", nil, nil, nil, "", fmt.Errorf("--work requires a value")
			}
			works = append(works, args[i+1])
			i++
		case strings.HasPrefix(a, "--work="):
			works = append(works, strings.TrimPrefix(a, "--work="))
		case a == "--lookahead" || a == "-lookahead":
			if i+1 >= len(args) {
				return "", nil, nil, nil, "", fmt.Errorf("--lookahead requires a value")
			}
			lookahead = args[i+1]
			i++
		case strings.HasPrefix(a, "--lookahead="):
			lookahead = strings.TrimPrefix(a, "--lookahead=")
		default:
			return "", nil, nil, nil, "", fmt.Errorf("unknown arg %q", a)
		}
	}
	if scope == "" || len(exits) == 0 {
		return "", nil, nil, nil, "", fmt.Errorf("usage: trace plan deep --scope <id> --exit <text> [--constraint <text>] [--work <title>] [--lookahead <summary>]")
	}
	return scope, exits, constraints, works, lookahead, nil
}

func cmdPlanShow(root string, args []string) int {
	fs := flag.NewFlagSet("plan show", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	goal := fs.String("goal", "", "goal UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *goal == "" {
		fmt.Fprintf(os.Stderr, "usage: trace plan show --goal <id>\n")
		return exitUsage
	}
	svc, st, code := openPlanner(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	fmt.Fprintf(os.Stderr, "plan: showing plan for goal %s\n", *goal)
	view, err := svc.GetPlan(context.Background(), *goal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan show: %v\n", err)
		return exitFail
	}
	if view.Phases == nil {
		view.Phases = []planner.PhaseView{}
	}
	if warn, err := svc.GoalStructureWarning(context.Background(), *goal); err == nil && warn != "" {
		fmt.Fprintf(os.Stderr, "plan show advisory: %s\n", warn)
	}
	tasks, err := st.ListTasksByGoalID(*goal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan show: %v\n", err)
		return exitFail
	}
	taskRows := make([]taskListRow, 0, len(tasks))
	for _, t := range tasks {
		taskRows = append(taskRows, taskListRow{
			ID: t.ID, Title: t.Title, WorkState: t.WorkState, GoalID: t.GoalID,
		})
	}
	out := planShowDTO{
		GoalID:           view.GoalID,
		CurrentScopeID:   view.CurrentScopeID,
		Phases:           view.Phases,
		CurrentDeepPlan:  view.CurrentDeepPlan,
		LookaheadScopeID: view.LookaheadScopeID,
		LookaheadSummary: view.LookaheadSummary,
		Tasks:            taskRows,
	}
	_ = json.NewEncoder(os.Stdout).Encode(out)
	return exitOK
}

type planShowDTO struct {
	GoalID           string                    `json:"goal_id"`
	CurrentScopeID   *string                   `json:"current_scope_id"`
	Phases           []planner.PhaseView       `json:"phases"`
	CurrentDeepPlan  *planner.DeepPlanDocument `json:"current_deep_plan"`
	LookaheadScopeID string                    `json:"lookahead_scope_id"`
	LookaheadSummary string                    `json:"lookahead_summary"`
	Tasks            []taskListRow             `json:"tasks"`
}

func cmdPlanBootstrap(root string, args []string) int {
	fs := flag.NewFlagSet("plan bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	goal := fs.String("goal", "", "goal UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *goal == "" {
		fmt.Fprintf(os.Stderr, "usage: trace plan bootstrap --goal <id>\n")
		return exitUsage
	}
	svc, st, code := openPlanner(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	fmt.Fprintf(os.Stderr, "plan: bootstrapping progressive plan for goal %s\n", *goal)
	fmt.Fprintln(os.Stderr, "plan bootstrap: yields minimal plan — refine with create-coarse / deep afterward")
	res, err := svc.BootstrapFromPlanChanges(context.Background(), *goal, "cli")
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan bootstrap: %v\n", err)
		return exitFail
	}
	if res.AlreadyExists {
		fmt.Fprintf(os.Stderr, "plan bootstrap: %s\n", res.Note)
	}
	_ = json.NewEncoder(os.Stdout).Encode(res)
	return exitOK
}

func cmdPlanApplyDiscovery(root string, args []string) int {
	discID, scopeID, pcID, pcTitle, pcBody, exits, constraints, works, err := parseApplyDiscoveryArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan apply-discovery: %v\n", err)
		return exitUsage
	}
	svc, st, code := openPlanner(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	items := make([]planner.WorkItem, 0, len(works))
	for _, w := range works {
		items = append(items, planner.WorkItem{Title: w})
	}
	fmt.Fprintf(os.Stderr, "plan: apply-discovery %s → scope %s\n", discID, scopeID)
	res, err := svc.ApplyDiscoveryReplan(context.Background(), planner.ApplyDiscoveryReplanInput{
		DiscoveryID:     discID,
		ScopeID:         scopeID,
		PlanChangeID:    pcID,
		PlanChangeTitle: pcTitle,
		PlanChangeBody:  pcBody,
		ExitCriteria:    exits,
		Constraints:     constraints,
		WorkItems:       items,
		Actor:           "cli",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan apply-discovery: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(res)
	return exitOK
}

func parseApplyDiscoveryArgs(args []string) (disc, scope, pcID, pcTitle, pcBody string, exits, constraints, works []string, err error) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--discovery" || a == "-discovery":
			if i+1 >= len(args) {
				return "", "", "", "", "", nil, nil, nil, fmt.Errorf("--discovery requires a value")
			}
			disc = args[i+1]
			i++
		case strings.HasPrefix(a, "--discovery="):
			disc = strings.TrimPrefix(a, "--discovery=")
		case a == "--scope" || a == "-scope":
			if i+1 >= len(args) {
				return "", "", "", "", "", nil, nil, nil, fmt.Errorf("--scope requires a value")
			}
			scope = args[i+1]
			i++
		case strings.HasPrefix(a, "--scope="):
			scope = strings.TrimPrefix(a, "--scope=")
		case a == "--plan-change" || a == "-plan-change":
			if i+1 >= len(args) {
				return "", "", "", "", "", nil, nil, nil, fmt.Errorf("--plan-change requires a value")
			}
			pcID = args[i+1]
			i++
		case strings.HasPrefix(a, "--plan-change="):
			pcID = strings.TrimPrefix(a, "--plan-change=")
		case a == "--pc-title" || a == "-pc-title":
			if i+1 >= len(args) {
				return "", "", "", "", "", nil, nil, nil, fmt.Errorf("--pc-title requires a value")
			}
			pcTitle = args[i+1]
			i++
		case strings.HasPrefix(a, "--pc-title="):
			pcTitle = strings.TrimPrefix(a, "--pc-title=")
		case a == "--pc-body" || a == "-pc-body":
			if i+1 >= len(args) {
				return "", "", "", "", "", nil, nil, nil, fmt.Errorf("--pc-body requires a value")
			}
			pcBody = args[i+1]
			i++
		case strings.HasPrefix(a, "--pc-body="):
			pcBody = strings.TrimPrefix(a, "--pc-body=")
		case a == "--exit" || a == "-exit":
			if i+1 >= len(args) {
				return "", "", "", "", "", nil, nil, nil, fmt.Errorf("--exit requires a value")
			}
			exits = append(exits, args[i+1])
			i++
		case strings.HasPrefix(a, "--exit="):
			exits = append(exits, strings.TrimPrefix(a, "--exit="))
		case a == "--constraint" || a == "-constraint":
			if i+1 >= len(args) {
				return "", "", "", "", "", nil, nil, nil, fmt.Errorf("--constraint requires a value")
			}
			constraints = append(constraints, args[i+1])
			i++
		case strings.HasPrefix(a, "--constraint="):
			constraints = append(constraints, strings.TrimPrefix(a, "--constraint="))
		case a == "--work" || a == "-work":
			if i+1 >= len(args) {
				return "", "", "", "", "", nil, nil, nil, fmt.Errorf("--work requires a value")
			}
			works = append(works, args[i+1])
			i++
		case strings.HasPrefix(a, "--work="):
			works = append(works, strings.TrimPrefix(a, "--work="))
		default:
			return "", "", "", "", "", nil, nil, nil, fmt.Errorf("unknown arg %q", a)
		}
	}
	if disc == "" || scope == "" {
		return "", "", "", "", "", nil, nil, nil, fmt.Errorf("usage: trace plan apply-discovery --discovery <id> --scope <id> […]")
	}
	return disc, scope, pcID, pcTitle, pcBody, exits, constraints, works, nil
}

func cmdPlanAckReplan(root string, args []string) int {
	fs := flag.NewFlagSet("plan ack-replan", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	scope := fs.String("scope", "", "scope UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *scope == "" {
		fmt.Fprintf(os.Stderr, "usage: trace plan ack-replan --scope <id>\n")
		return exitUsage
	}
	svc, st, code := openPlanner(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	fmt.Fprintf(os.Stderr, "plan: ack-replan scope %s\n", *scope)
	if err := svc.AckReplan(context.Background(), *scope); err != nil {
		fmt.Fprintf(os.Stderr, "plan ack-replan: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"ok": true, "scope": *scope, "auto_replan_count": 0})
	return exitOK
}

func openPlanner(root string) (*planner.Service, *store.Store, int) {
	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan: %v\n", err)
		return nil, nil, exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "plan: %v\n", err)
		return nil, nil, exitFail
	}
	if code := failCLIDenied(domain.New(st), "plan", "plan"); code != exitOK {
		st.Close()
		return nil, nil, code
	}
	return planner.New(st), st, exitOK
}
