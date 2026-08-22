package install

import (
	"context"
	"fmt"
	"io"

	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

// PrintBootstrapHintIfNeeded warns when goals exist but no progressive planner rows are present.
func PrintBootstrapHintIfNeeded(projectRoot string, w io.Writer) {
	if w == nil {
		return
	}
	st, err := store.Open(projectRoot)
	if err != nil {
		return
	}
	defer st.Close()

	goals, err := st.ListGoals()
	if err != nil || len(goals) == 0 {
		return
	}
	psvc := planner.New(st)
	ctx := context.Background()
	for _, g := range goals {
		ok, err := psvc.PlanExists(ctx, g.ID)
		if err != nil || ok {
			continue
		}
		fmt.Fprintf(w, "install: moat-first: trace_tasks → trace_context → trace_loop before product edits\n")
		fmt.Fprintf(w, "install: graph-first GUI: trace serve → open / (Explore) for orient + bounded neighborhood\n")
		fmt.Fprintf(w, "install: goal %q has no progressive plan — run trace plan create-coarse or trace plan bootstrap --goal %s (MCP trace_plan)\n", g.Title, g.ID)
		return
	}
}
