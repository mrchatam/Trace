package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// run dispatches CLI args. Global -C/--project may precede the command.
func run(args []string) int {
	root, rest, err := parseGlobal(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return exitUsage
	}
	if len(rest) == 0 {
		printHelp(os.Stdout)
		return exitOK
	}

	cmd := rest[0]
	cmdArgs := rest[1:]

	switch cmd {
	case "help", "-h", "--help":
		printHelp(os.Stdout)
		return exitOK
	case "version", "--version", "-version":
		fmt.Println(version)
		return exitOK
	case "init":
		return cmdInit(root, cmdArgs)
	case "index":
		return cmdIndex(root, cmdArgs, "index")
	case "reindex":
		return cmdIndex(root, cmdArgs, "reindex")
	case "add":
		return cmdAdd(root, cmdArgs)
	case "link":
		return cmdLink(root, cmdArgs)
	case "transition":
		return cmdTransition(root, cmdArgs)
	case "review":
		return cmdReview(root, cmdArgs)
	case "impact":
		return cmdImpact(root, cmdArgs)
	case "capability":
		return cmdCapability(root, cmdArgs)
	case "plan":
		return cmdPlan(root, cmdArgs)
	case "seed":
		return cmdSeed(root, cmdArgs)
	case "tasks":
		return cmdTasks(root, cmdArgs)
	case "why":
		return cmdWhy(root, cmdArgs)
	case "context":
		return cmdContext(root, cmdArgs)
	case "explore":
		return cmdExplore(root, cmdArgs)
	case "loop":
		return cmdLoop(root, cmdArgs)
	case "agents":
		return cmdAgents(root, cmdArgs)
	case "migrate":
		return cmdMigrate(root, cmdArgs)
	case "backup":
		return cmdBackup(root, cmdArgs)
	case "restore":
		return cmdRestore(root, cmdArgs)
	case "auth":
		return cmdAuth(root, cmdArgs)
	case "install":
		return cmdInstall(root, cmdArgs)
	case "changes":
		return cmdChanges(root, cmdArgs)
	case "patterns":
		return cmdPatterns(root, cmdArgs)
	case "knowledge":
		return cmdKnowledge(root, cmdArgs)
	case "search":
		return cmdSearch(root, cmdArgs)
	case "test":
		return cmdTest(root, cmdArgs)
	case "tests":
		return cmdTests(root, cmdArgs)
	case "verify":
		return cmdVerify(root, cmdArgs)
	case "eval":
		return cmdEval(root, cmdArgs)
	case "outcomes":
		return cmdOutcomes(root, cmdArgs)
	case "regressions":
		return cmdRegressions(root, cmdArgs)
	case "serve":
		return cmdServe(root, cmdArgs)
	case "gui":
		return cmdGui(root, cmdArgs)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printHelp(os.Stderr)
		return exitUsage
	}
}

func parseGlobal(args []string) (root string, rest []string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", nil, err
	}
	root = cwd
	i := 0
	for i < len(args) {
		a := args[i]
		switch {
		case a == "-C" || a == "--project":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("flag %s requires a directory argument", a)
			}
			root = args[i+1]
			i += 2
		case strings.HasPrefix(a, "-C="):
			root = strings.TrimPrefix(a, "-C=")
			i++
		case strings.HasPrefix(a, "--project="):
			root = strings.TrimPrefix(a, "--project=")
			i++
		default:
			return root, args[i:], nil
		}
	}
	return root, nil, nil
}

func resolveRoot(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	return abs, nil
}
