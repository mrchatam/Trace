package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// taskListRow is the DF-02 cold-start list shape (id/title/work_state/goal_id).
type taskListRow struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	WorkState string  `json:"work_state"`
	GoalID    *string `json:"goal_id"`
}

type tasksConflictsResponse struct {
	OK        bool                  `json:"ok"`
	Conflicts []domain.WorkConflict `json:"conflicts"`
}

func cmdTasks(root string, args []string) int {
	if len(args) > 0 && args[0] == "conflicts" {
		return cmdTasksConflicts(root, args[1:])
	}
	return cmdTasksList(root, args)
}

func cmdTasksList(root string, args []string) int {
	fs := flag.NewFlagSet("tasks", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	goalID := fs.String("goal", "", "optional goal UUID filter")
	limit := fs.Int("limit", store.DefaultTaskListLimit, "max rows (default 50, max 500); ignored with --all")
	cursor := fs.String("cursor", "", "opaque pagination cursor from a prior next_cursor")
	all := fs.Bool("all", false, "return all matching tasks (no limit; prefer --limit for agents)")
	var workStates multiFlag
	fs.Var(&workStates, "work-state", "filter by work_state (repeatable); default any")
	if err := fs.Parse(flagsFirst(args, map[string]bool{
		"goal": true, "limit": true, "cursor": true, "work-state": true, "all": false,
	})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace tasks [--goal <id>] [--work-state STATE]... [--limit N] [--cursor C] [--all]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tasks: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tasks: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "tasks", "tasks"); code != exitOK {
		return code
	}

	filt := store.TaskListFilter{
		GoalID:     *goalID,
		WorkStates: []string(workStates),
		Cursor:     *cursor,
	}
	if *all {
		filt.Limit = -1
	} else {
		filt.Limit = *limit
	}
	res, err := st.ListTasksFiltered(filt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tasks: %v\n", err)
		return exitFail
	}

	outRows := make([]taskListRow, 0, len(res.Tasks))
	for _, t := range res.Tasks {
		outRows = append(outRows, taskListRow{
			ID: t.ID, Title: t.Title, WorkState: t.WorkState, GoalID: t.GoalID,
		})
	}
	payload := map[string]any{
		"items":     outRows,
		"count":     len(outRows),
		"truncated": res.Truncated,
	}
	if res.NextCursor != "" {
		payload["next_cursor"] = res.NextCursor
	}
	if err := json.NewEncoder(os.Stdout).Encode(payload); err != nil {
		fmt.Fprintf(os.Stderr, "tasks: %v\n", err)
		return exitFail
	}
	return exitOK
}

// multiFlag collects repeatable string flags.
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

func cmdTasksConflicts(root string, args []string) int {
	fs := flag.NewFlagSet("tasks conflicts", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "optional task UUID filter")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace tasks conflicts [--task <id>]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tasks conflicts: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tasks conflicts: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "tasks", "tasks conflicts"); code != exitOK {
		return code
	}

	conflicts, err := svc.DetectWorkConflicts(context.Background(), domain.DetectWorkConflictsOpts{
		TaskID: *taskID,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "tasks conflicts: %v\n", err)
		return exitFail
	}
	if conflicts == nil {
		conflicts = []domain.WorkConflict{}
	}
	resp := tasksConflictsResponse{OK: true, Conflicts: conflicts}
	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		fmt.Fprintf(os.Stderr, "tasks conflicts: %v\n", err)
		return exitFail
	}
	return exitOK
}
