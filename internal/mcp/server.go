package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName          = "trace"
	serverVersion       = "0.0.0-dev"
	traceAddDescription = "Create a domain entity (discovery|task|goal|decision|assumption|plan-change|claim|evidence). Prefer the task/promotion path over discovery-only edits: after a BLOCKING discovery, promote with trace_add kind=task or loop apply spawned_tasks with discovery_id before product edits — do not discovery-only then edit. Mirrors `trace add`."
)

// Options configure the Trace MCP server.
type Options struct {
	// ProjectRoot is the default project directory (-C / --project). Empty = cwd.
	ProjectRoot string
}

// Server is a thin MCP adapter over Trace libraries (G19: no domain fork).
type Server struct {
	defaultRoot string
	mcp         *sdkmcp.Server
}

// NewServer builds an MCP server with Trace tools registered.
func NewServer(opts Options) *Server {
	s := &Server{defaultRoot: opts.ProjectRoot}
	s.mcp = sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, &sdkmcp.ServerOptions{
		Instructions: ServerInstructions(),
	})
	s.registerTools()
	return s
}

// MCP returns the underlying SDK server (tests / advanced wiring).
func (s *Server) MCP() *sdkmcp.Server { return s.mcp }

// RunStdio serves MCP over stdin/stdout until the client disconnects.
func (s *Server) RunStdio(ctx context.Context) error {
	if s == nil || s.mcp == nil {
		return fmt.Errorf("mcp: nil server")
	}
	return s.mcp.Run(ctx, &sdkmcp.StdioTransport{})
}

func boolPtr(v bool) *bool { return &v }

func (s *Server) registerTools() {
	openWorldFalse := boolPtr(false)
	destructiveFalse := boolPtr(false)
	destructiveTrue := boolPtr(true)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_why",
		Description: "Explain why an entity exists via causal neighborhood (retrieval.Why). Mirrors `trace why <type> <id>`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Why",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolWhy)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_context",
		Description: "Compile a task context packet (compiler.TaskContext / ExpandContext). Mirrors `trace context`. depth 1|2; format json|markdown|both; optional query merges agent FTS hits (task_id required).",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Context",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolContext)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_add",
		Description: traceAddDescription,
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Add",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolAdd)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_link",
		Description: "Link two entities (goal-task|decision-task|discovery-plan-change|discovery-mentions-task|claim-evidence). Mirrors `trace link`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Link",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolLink)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_transition",
		Description: "Transition a task work_state. DONE requires linked Review PASS + as_operator with no linked FAIL, or allow_done hatch (returns warning); as_operator is a conscious claim (flag≠identity / not verified operator identity); evidence_ids alone do not authorize DONE; actor is not authorization; missing caps fail-closed unless allow_missing_caps; allow_done does not bypass missing caps. Mirrors `trace transition`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Transition",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveTrue,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolTransition)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_review",
		Description: "Create a review (optional task link) or set result PASS|FAIL|UNCERTAIN. action=create|set. Mirrors `trace review`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Review",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveTrue,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolReview)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_tasks",
		Description: "List tasks as JSON array (id, title, work_state, goal_id). Optional goal_id filter. Mirrors `trace tasks`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Tasks",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolTasks)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_capability",
		Description: "Capability catalog: action=declare|list|require|unrequire|missing. missing requires task or task_id (list via trace_tasks). Mirrors `trace capability`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Capability",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolCapability)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_impact",
		Description: "Decision impact: action=finding|alternative|report|walk. finding op=add|list; alternative op=add|list|recommend. Mirrors `trace impact`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Impact",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolImpact)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_version",
		Description: "Return live MCP process identity {ok,name,version} so agents can detect a stale stdio server after rebuild/install.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Version",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolVersion)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_search",
		Description: "Full-text search over indexed entities. Args: query (required), limit (optional). Mirrors `trace search`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Search",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolSearch)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_changes",
		Description: "Change history: action=list|show|compare. list: optional task_id, limit; show: change_id; compare: from, to. Mirrors `trace changes`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Changes",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolChanges)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_regressions",
		Description: "Regression history: action=list. Optional task_id, change_id, limit. Mirrors `trace regressions list`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Regressions",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolRegressions)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_loop",
		Description: "Deliberation loop: action=next|apply|status. next/status require task_id; apply requires envelope (trace.loop.apply.v1). Mirrors `trace loop`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Loop",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveTrue,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolLoop)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_agents",
		Description: "Harness agent catalog: action=list|describe|recommend. describe requires slug; recommend requires task_id or phase (+ optional keywords). Mirrors `trace agents`. Recommend-only — no spawn.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Agents",
			ReadOnlyHint:    true,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolAgents)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_plan",
		Description: "Progressive planner: action=create-coarse|set-current|deep|show|bootstrap. Mirrors `trace plan`. Bootstrap recovers PlanExists from plan_changes.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Plan",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolPlan)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_explore",
		Description: "Unified task-aware capped explore: task context packet + search + why on top hits + neighborhood. task_id required; optional query merges via G1. Mirrors `trace explore`.",
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Explore",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolExplore)
}

// RegisteredToolNames returns the locked seventeen MCP tool names (registration order).
func RegisteredToolNames() []string {
	return []string{
		"trace_why", "trace_context", "trace_add",
		"trace_link", "trace_transition", "trace_review",
		"trace_tasks", "trace_capability", "trace_impact", "trace_version",
		"trace_search", "trace_changes", "trace_regressions", "trace_loop",
		"trace_agents", "trace_plan", "trace_explore",
	}
}

func textResult(text string) *sdkmcp.CallToolResult {
	return &sdkmcp.CallToolResult{
		Content: []sdkmcp.Content{&sdkmcp.TextContent{Text: text}},
	}
}
