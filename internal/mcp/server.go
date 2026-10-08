package mcp

import (
	"context"
	"fmt"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/buildmeta"
)

const serverName = "trace"

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
		Version: buildmeta.String(),
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
		Description: descWhy,
		InputSchema: schemaFor[WhyInput](nil),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Why",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolWhy)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_context",
		Description: descContext,
		InputSchema: schemaFor[ContextInput](map[string]propConstraint{
			"depth":     rangeOf(1, 2),
			"max_layer": rangeOf(1, 3),
			"format":    enumOf(contextFormats...),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Context",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolContext)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_add",
		Description: descAdd,
		InputSchema: schemaFor[AddInput](map[string]propConstraint{
			"kind":   enumOf(addKinds...),
			"status": enumOf(provenanceStatus...),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Add",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolAdd)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_link",
		Description: descLink,
		InputSchema: schemaFor[LinkInput](map[string]propConstraint{
			"rel": enumOf(linkRels...),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Link",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolLink)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_transition",
		Description: descTransition,
		InputSchema: schemaFor[TransitionInput](map[string]propConstraint{
			"to_state": enumOf(workStates...),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Transition",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveTrue,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolTransition)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_review",
		Description: descReview,
		InputSchema: schemaFor[ReviewInput](map[string]propConstraint{
			"action": enumOf(reviewActions...),
			"result": enumOf(reviewResults...),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Review",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveTrue,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolReview)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_tasks",
		Description: descTasks,
		InputSchema: schemaFor[TasksInput](map[string]propConstraint{
			"limit":       rangeOf(1, 500),
			"work_state":  enumOf(workStates...),
			"work_states": {enum: workStates},
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Tasks",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolTasks)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_capability",
		Description: descCapability,
		InputSchema: schemaFor[CapabilityInput](map[string]propConstraint{
			"action": enumOf(capabilityActions...),
			"kind":   enumOf(capabilityKinds...),
			"status": enumOf(capabilityStatus...),
			"limit":  rangeOf(1, 500),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Capability",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolCapability)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_impact",
		Description: descImpact,
		InputSchema: schemaFor[ImpactInput](map[string]propConstraint{
			"action": enumOf(impactActions...),
			"op":     enumOf(impactOps...),
			"depth":  rangeOf(1, 8),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Impact",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolImpact)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_version",
		Description: descVersion,
		InputSchema: schemaFor[VersionInput](nil),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Version",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolVersion)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_search",
		Description: descSearch,
		InputSchema: schemaFor[SearchInput](map[string]propConstraint{
			"limit":   rangeOf(1, 64),
			"matcher": enumOf(matchers...),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Search",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolSearch)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_changes",
		Description: descChanges,
		InputSchema: schemaFor[ChangesInput](map[string]propConstraint{
			"action": enumOf(changeActions...),
			"limit":  rangeOf(1, 64),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Changes",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolChanges)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_regressions",
		Description: descRegressions,
		InputSchema: schemaFor[RegressionsInput](map[string]propConstraint{
			"action": enumOf(regressionAction...),
			"limit":  rangeOf(1, 64),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:         "Regressions",
			ReadOnlyHint:  true,
			OpenWorldHint: openWorldFalse,
		},
	}, s.toolRegressions)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_loop",
		Description: descLoop,
		InputSchema: schemaFor[LoopInput](map[string]propConstraint{
			"action": enumOf(loopActions...),
			"for":    enumOf(loopGateFors...),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Loop",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveTrue,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolLoop)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_agents",
		Description: descAgents,
		InputSchema: schemaFor[AgentsInput](map[string]propConstraint{
			"action": enumOf(agentActions...),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Agents",
			ReadOnlyHint:    true,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolAgents)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_plan",
		Description: descPlan,
		InputSchema: schemaFor[PlanInput](map[string]propConstraint{
			"action": enumOf(planActions...),
		}),
		Annotations: &sdkmcp.ToolAnnotations{
			Title:           "Plan",
			ReadOnlyHint:    false,
			DestructiveHint: destructiveFalse,
			OpenWorldHint:   openWorldFalse,
		},
	}, s.toolPlan)

	sdkmcp.AddTool(s.mcp, &sdkmcp.Tool{
		Name:        "trace_explore",
		Description: descExplore,
		InputSchema: schemaFor[ExploreInput](map[string]propConstraint{
			"limit":     rangeOf(1, 64),
			"max_nodes": rangeOf(1, 500),
			"depth":     rangeOf(1, 8),
			"matcher":   enumOf(matchers...),
		}),
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
