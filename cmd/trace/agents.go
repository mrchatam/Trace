package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/agents"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

func cmdAgents(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace agents list|describe|recommend\n")
		return exitUsage
	}
	switch args[0] {
	case "list":
		return cmdAgentsList(root, args[1:])
	case "describe":
		return cmdAgentsDescribe(root, args[1:])
	case "recommend":
		return cmdAgentsRecommend(root, args[1:])
	case "help", "-h", "--help":
		printAgentsHelp(os.Stdout)
		return exitOK
	default:
		fmt.Fprintf(os.Stderr, "agents: unknown subcommand %q\n", args[0])
		printAgentsHelp(os.Stderr)
		return exitUsage
	}
}

func printAgentsHelp(w *os.File) {
	fmt.Fprint(w, `trace agents — harness agent catalog (recommend-only; no spawn)

Subcommands:
  list
        JSON array of catalog profiles (slug, title, subagent_type,
        deliberation_phases, requirements). Empty catalog → [].
  describe <slug>
        Full profile + requirements + registry metadata for one agent slug.
  recommend (--task <id> | --phase <PHASE>) [--goal-keywords "kw ..."]
        Ranked recommendations (max 4); same shape as harness_recommendations[].
        --task resolves deliberation phase from persisted loop state when available.
`)
}

func cmdAgentsList(root string, args []string) int {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace agents list\n")
		return exitUsage
	}
	_, st, _, code := openAgentsStore(root, "agents list")
	if code != exitOK {
		return code
	}
	defer st.Close()

	items, err := agents.ListAgentSummaries(context.Background(), st)
	if err != nil {
		fmt.Fprintf(os.Stderr, "agents list: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(items); err != nil {
		fmt.Fprintf(os.Stderr, "agents list: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdAgentsDescribe(root string, args []string) int {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: trace agents describe <slug>\n")
		return exitUsage
	}
	_, st, _, code := openAgentsStore(root, "agents describe")
	if code != exitOK {
		return code
	}
	defer st.Close()

	profile, err := agents.DescribeAgent(context.Background(), st, args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(profile); err != nil {
		fmt.Fprintf(os.Stderr, "agents describe: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdAgentsRecommend(root string, args []string) int {
	fs := flag.NewFlagSet("agents recommend", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "seed task UUID")
	phase := fs.String("phase", "", "deliberation phase e.g. CRITIQUE")
	keywords := fs.String("goal-keywords", "", "optional keyword injection for routing")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true, "phase": true, "goal-keywords": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace agents recommend (--task <id> | --phase <PHASE>) [--goal-keywords \"kw ...\"]\n")
		return exitUsage
	}

	_, st, svc, code := openAgentsStore(root, "agents recommend")
	if code != exitOK {
		return code
	}
	defer st.Close()

	recs, err := loop.RecommendHarness(context.Background(), svc, st, planner.New(st), loop.RecommendHarnessInput{
		TaskID:   *taskID,
		Phase:    *phase,
		Keywords: *keywords,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return exitFail
	}
	if recs == nil {
		recs = []agents.Recommendation{}
	}
	if err := json.NewEncoder(os.Stdout).Encode(recs); err != nil {
		fmt.Fprintf(os.Stderr, "agents recommend: %v\n", err)
		return exitFail
	}
	return exitOK
}

func openAgentsStore(root, prefix string) (abs string, st *store.Store, svc *domain.Service, code int) {
	var err error
	abs, err = resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", prefix, err)
		return "", nil, nil, exitFail
	}
	st, err = store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", prefix, err)
		return "", nil, nil, exitFail
	}
	svc = domain.New(st)
	if code = failCLIDenied(svc, "agents", prefix); code != exitOK {
		st.Close()
		return "", nil, nil, code
	}
	return abs, st, svc, exitOK
}
