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

type changeCaptureRow struct {
	ID        string `json:"id"`
	GitCommit string `json:"git_commit"`
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
}

type changeListRow struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	GitCommit  string `json:"git_commit"`
	Status     string `json:"status"`
	Reason     string `json:"reason"`
	SourceType string `json:"source_type"`
	CreatedAt  string `json:"created_at"`
}

type changeShowPathRow struct {
	Path     string `json:"path"`
	Status   string `json:"status"`
	SymbolID string `json:"symbol_id"`
}

type changeShowResponse struct {
	OK     bool                `json:"ok"`
	Change store.Change        `json:"change"`
	Paths  []changeShowPathRow `json:"paths"`
}

func cmdChanges(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace changes capture|compare|list|show|similar\n")
		return exitUsage
	}
	switch args[0] {
	case "capture":
		return cmdChangesCapture(root, args[1:])
	case "compare":
		return cmdChangesCompare(root, args[1:])
	case "list":
		return cmdChangesList(root, args[1:])
	case "show":
		return cmdChangesShow(root, args[1:])
	case "similar":
		return cmdChangesSimilar(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown changes subcommand: %s\n", args[0])
		return exitUsage
	}
}

func cmdChangesCapture(root string, args []string) int {
	fs := flag.NewFlagSet("changes capture", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	since := fs.String("since", "", "promote indexed commits after this OID (exclusive)")
	all := fs.Bool("all", false, "include commits with only non-meaningful paths")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"since": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace changes capture [--since <oid>] [--all]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes capture: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes capture: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "changes", "changes capture"); code != exitOK {
		return code
	}

	opts := domain.PromoteVCSCommitOptions{AllPaths: *all}
	promoted, err := svc.PromoteRecentVCSCommits(context.Background(), *since, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes capture: %v\n", err)
		return exitFail
	}

	out := make([]changeCaptureRow, 0, len(promoted))
	for _, c := range promoted {
		out = append(out, changeCaptureRow{
			ID:        c.ID,
			GitCommit: c.GitCommit,
			TaskID:    c.TaskID,
			Status:    c.Status,
			Reason:    c.Reason,
		})
	}
	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "changes capture: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdChangesCompare(root string, args []string) int {
	fs := flag.NewFlagSet("changes compare", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	from := fs.String("from", "", "start git OID (exclusive ancestry baseline)")
	to := fs.String("to", "", "end git OID")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"from": true, "to": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 || strings.TrimSpace(*from) == "" || strings.TrimSpace(*to) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace changes compare --from <oid> --to <oid>\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes compare: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes compare: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "changes", "changes compare"); code != exitOK {
		return code
	}

	result, err := svc.CompareStates(context.Background(), *from, *to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes compare: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "changes compare: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdChangesList(root string, args []string) int {
	fs := flag.NewFlagSet("changes list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "optional task UUID filter")
	limit := fs.Int("limit", 0, "max rows (default 32, cap 64)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true, "limit": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace changes list [--task <uuid>] [--limit N]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes list: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes list: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "changes", "changes list"); code != exitOK {
		return code
	}

	changes, err := svc.ListChanges(context.Background(), *limit, *taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes list: %v\n", err)
		return exitFail
	}

	out := make([]changeListRow, 0, len(changes))
	for _, c := range changes {
		out = append(out, changeListRow{
			ID:         c.ID,
			TaskID:     c.TaskID,
			GitCommit:  c.GitCommit,
			Status:     c.Status,
			Reason:     c.Reason,
			SourceType: c.SourceType,
			CreatedAt:  c.CreatedAt,
		})
	}
	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "changes list: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdChangesShow(root string, args []string) int {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: trace changes show <change-id>\n")
		return exitUsage
	}
	changeID := strings.TrimSpace(args[0])
	if changeID == "" {
		fmt.Fprintf(os.Stderr, "usage: trace changes show <change-id>\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes show: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes show: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "changes", "changes show"); code != exitOK {
		return code
	}

	ctx := context.Background()
	c, err := svc.GetChange(ctx, changeID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes show: %v\n", err)
		return exitFail
	}
	paths, err := svc.ListChangePaths(ctx, changeID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes show: %v\n", err)
		return exitFail
	}
	pathRows := make([]changeShowPathRow, 0, len(paths))
	for _, p := range paths {
		pathRows = append(pathRows, changeShowPathRow{
			Path:     p.Path,
			Status:   p.Status,
			SymbolID: p.SymbolID,
		})
	}
	resp := changeShowResponse{OK: true, Change: c, Paths: pathRows}
	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		fmt.Fprintf(os.Stderr, "changes show: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdChangesSimilar(root string, args []string) int {
	fs := flag.NewFlagSet("changes similar", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	pathPrefix := fs.String("path", "", "repo-relative path prefix filter")
	changeKind := fs.String("kind", "", "change kind filter (e.g. seg:internal)")
	limit := fs.Int("limit", 0, "max rows (default 32, cap 64)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"path": true, "kind": true, "limit": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace changes similar --path <prefix> | --kind <kind> [--limit N]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes similar: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes similar: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "changes", "changes similar"); code != exitOK {
		return code
	}

	result, err := svc.QuerySimilarChanges(context.Background(), domain.SimilarChangesOpts{
		PathPrefix: *pathPrefix,
		ChangeKind: *changeKind,
		Limit:      *limit,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "changes similar: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "changes similar: %v\n", err)
		return exitFail
	}
	return exitOK
}

// promoteVCSCommitsAfterIndex best-effort promotes new indexed commits after index.
func promoteVCSCommitsAfterIndex(ctx context.Context, svc *domain.Service) {
	if _, err := svc.PromoteRecentVCSCommits(ctx, "", domain.PromoteVCSCommitOptions{}); err != nil {
		fmt.Fprintf(os.Stderr, "index: vcs change capture: %v\n", err)
	}
}
