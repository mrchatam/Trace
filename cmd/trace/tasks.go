package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

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
	if err := fs.Parse(flagsFirst(args, map[string]bool{"goal": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace tasks [--goal <id>]\n")
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

	var tasks []store.Task
	if *goalID != "" {
		tasks, err = st.ListTasksByGoalID(*goalID)
	} else {
		tasks, err = st.ListTasks()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "tasks: %v\n", err)
		return exitFail
	}

	out := make([]taskListRow, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, taskListRow{
			ID:        t.ID,
			Title:     t.Title,
			WorkState: t.WorkState,
			GoalID:    t.GoalID,
		})
	}
	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "tasks: %v\n", err)
		return exitFail
	}
	return exitOK
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
