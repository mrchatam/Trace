package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	tracemcp "github.com/mrchatam/Trace/internal/mcp"
)

func printMCPHelp(w io.Writer) {
	names := tracemcp.RegisteredToolNames()
	fmt.Fprintf(w, "usage: trace-mcp [-C|--project <dir>]\n")
	fmt.Fprintf(w, "  Thin MCP stdio server (official go-sdk). Tools: %s.\n", strings.Join(names, ", "))
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	root, err := parseProject(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "trace-mcp: %v\n", err)
		return 2
	}
	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: root})
	if err := srv.RunStdio(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "trace-mcp: %v\n", err)
		return 1
	}
	return 0
}

// parseProject accepts optional -C / --project like cmd/trace. Remaining args are ignored
// (stdio MCP has no subcommands). Empty root → cwd via internal/mcp.
func parseProject(args []string) (string, error) {
	root := ""
	i := 0
	for i < len(args) {
		a := args[i]
		switch {
		case a == "-C" || a == "--project":
			if i+1 >= len(args) {
				return "", fmt.Errorf("flag %s requires a directory argument", a)
			}
			root = args[i+1]
			i += 2
		case strings.HasPrefix(a, "-C="):
			root = strings.TrimPrefix(a, "-C=")
			i++
		case strings.HasPrefix(a, "--project="):
			root = strings.TrimPrefix(a, "--project=")
			i++
		case a == "-h" || a == "--help":
			printMCPHelp(os.Stdout)
			os.Exit(0)
			return "", nil
		default:
			return "", fmt.Errorf("unknown argument %q (want -C|--project only)", a)
		}
	}
	return root, nil
}
