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

// capabilityListRow is the DF-32 list/missing shape (snake_case; body/timestamps omitted).
type capabilityListRow struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Slug   string `json:"slug"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func capabilityListRows(list []store.Capability) []capabilityListRow {
	out := make([]capabilityListRow, 0, len(list))
	for _, c := range list {
		out = append(out, capabilityListRow{
			ID: c.ID, Kind: c.Kind, Slug: c.Slug, Title: c.Title, Status: c.Status,
		})
	}
	return out
}

// cmdCapability is a thin G19 adapter: declare|list|require|unrequire|missing|decide|decisions via domain only.
func cmdCapability(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace capability declare|list|require|unrequire|missing|decide|decisions …\n")
		return exitUsage
	}
	switch args[0] {
	case "declare":
		return cmdCapabilityDeclare(root, args[1:])
	case "list":
		return cmdCapabilityList(root, args[1:])
	case "require":
		return cmdCapabilityRequire(root, args[1:])
	case "unrequire":
		return cmdCapabilityUnrequire(root, args[1:])
	case "missing":
		return cmdCapabilityMissing(root, args[1:])
	case "decide":
		return cmdCapabilityDecide(root, args[1:])
	case "decisions":
		return cmdCapabilityDecisions(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "capability: unknown subcommand %q (want declare|list|require|unrequire|missing|decide|decisions)\n", args[0])
		return exitUsage
	}
}

func cmdCapabilityDeclare(root string, args []string) int {
	fs := flag.NewFlagSet("capability declare", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	kind := fs.String("kind", "", "SKILL|RULE|MCP|TOOL|HOOK")
	slug := fs.String("slug", "", "unique slug (e.g. mcp:trace_why)")
	title := fs.String("title", "", "optional title")
	status := fs.String("status", "", "AVAILABLE|UNAVAILABLE|UNKNOWN (default UNKNOWN)")
	body := fs.String("body", "", "optional body")
	id := fs.String("id", "", "optional UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if strings.TrimSpace(*kind) == "" || strings.TrimSpace(*slug) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace capability declare --kind SKILL|RULE|MCP|TOOL|HOOK --slug <slug> [--title <t>] [--status AVAILABLE|UNAVAILABLE|UNKNOWN] [--body <text>] [--id <uuid>]\n")
		fmt.Fprintf(os.Stderr, "note: omit --id to update an existing capability by slug (stable id)\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	c, err := svc.UpsertCapability(context.Background(), domain.CapabilityInput{
		ID: *id, Kind: *kind, Slug: *slug, Title: *title, Status: *status, Body: *body,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "capability: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "id": c.ID, "kind": c.Kind, "slug": c.Slug,
		"title": c.Title, "status": c.Status,
	})
	return exitOK
}

func cmdCapabilityList(root string, args []string) int {
	fs := flag.NewFlagSet("capability list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	kind := fs.String("kind", "", "optional SKILL|RULE|MCP|TOOL|HOOK")
	status := fs.String("status", "", "optional AVAILABLE|UNAVAILABLE|UNKNOWN")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	list, err := svc.ListCapabilities(context.Background(), domain.ListCapabilitiesFilter{
		Kind: *kind, Status: *status,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "capability: %v\n", err)
		return exitFail
	}
	if list == nil {
		list = []store.Capability{}
	}
	rows := capabilityListRows(list)
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "capabilities": rows, "count": len(rows),
	})
	return exitOK
}

func cmdCapabilityRequire(root string, args []string) int {
	fs := flag.NewFlagSet("capability require", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "task UUID")
	capRef := fs.String("capability", "", "capability id or slug")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *taskID == "" || strings.TrimSpace(*capRef) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace capability require --task <id> --capability <id|slug>\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	cap, err := svc.ResolveCapabilityIDOrSlug(context.Background(), *capRef)
	if err != nil {
		fmt.Fprintf(os.Stderr, "capability: %v\n", err)
		return exitFail
	}
	r, err := svc.RequireCapability(context.Background(), *taskID, cap.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "capability: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "id": r.ID, "task": r.TaskID, "capability": r.CapabilityID, "slug": cap.Slug,
	})
	return exitOK
}

func cmdCapabilityUnrequire(root string, args []string) int {
	fs := flag.NewFlagSet("capability unrequire", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "task UUID")
	capRef := fs.String("capability", "", "capability id or slug")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *taskID == "" || strings.TrimSpace(*capRef) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace capability unrequire --task <id> --capability <id|slug>\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	cap, err := svc.ResolveCapabilityIDOrSlug(context.Background(), *capRef)
	if err != nil {
		fmt.Fprintf(os.Stderr, "capability: %v\n", err)
		return exitFail
	}
	if err := svc.UnrequireCapability(context.Background(), *taskID, cap.ID); err != nil {
		fmt.Fprintf(os.Stderr, "capability: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "task": *taskID, "capability": cap.ID, "slug": cap.Slug,
	})
	return exitOK
}

func cmdCapabilityMissing(root string, args []string) int {
	fs := flag.NewFlagSet("capability missing", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "task UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *taskID == "" {
		fmt.Fprintf(os.Stderr, "usage: trace capability missing --task <id>\n")
		fmt.Fprintf(os.Stderr, "task id required; list tasks: trace tasks\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	list, err := svc.MissingCapabilities(context.Background(), *taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "capability: %v\n", err)
		return exitFail
	}
	if list == nil {
		list = []store.Capability{}
	}
	rows := capabilityListRows(list)
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "task": *taskID, "missing": rows, "count": len(rows),
	})
	return exitOK
}

func cmdCapabilityDecide(root string, args []string) int {
	fs := flag.NewFlagSet("capability decide", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	slug := fs.String("slug", "", "tool slug (e.g. mcp:trace_why)")
	decision := fs.String("decision", "", "ALLOWED|DENIED")
	reason := fs.String("reason", "", "optional reason")
	actor := fs.String("actor", "cli", "actor label (not authorization)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if strings.TrimSpace(*slug) == "" || strings.TrimSpace(*decision) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace capability decide --slug <slug> --decision ALLOWED|DENIED [--reason <text>] [--actor cli]\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	row, err := svc.DecideTool(context.Background(), domain.DecideToolInput{
		Slug: *slug, Decision: *decision, Reason: *reason, Actor: *actor,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "capability: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "id": row.ID, "slug": row.Slug, "decision": row.Decision,
		"reason": row.Reason, "actor": row.Actor,
		"created_at": row.CreatedAt, "updated_at": row.UpdatedAt,
	})
	return exitOK
}

func cmdCapabilityDecisions(root string, args []string) int {
	fs := flag.NewFlagSet("capability decisions", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	list, err := svc.ListToolDecisions(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "capability: %v\n", err)
		return exitFail
	}
	if list == nil {
		list = []store.CapabilityToolDecision{}
	}
	type row struct {
		ID        string `json:"id"`
		Slug      string `json:"slug"`
		Decision  string `json:"decision"`
		Reason    string `json:"reason"`
		Actor     string `json:"actor"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	out := make([]row, 0, len(list))
	for _, d := range list {
		out = append(out, row{
			ID: d.ID, Slug: d.Slug, Decision: d.Decision, Reason: d.Reason,
			Actor: d.Actor, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
		})
	}
	_ = json.NewEncoder(os.Stdout).Encode(out)
	return exitOK
}
