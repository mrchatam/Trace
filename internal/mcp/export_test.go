package mcp

import (
	"context"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Test seams: export tool handlers for mcp_test without widening the public API
// beyond what VERIFY needs.

func CallWhy(s *Server, ctx context.Context, in WhyInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolWhy(ctx, nil, in)
}

func CallContext(s *Server, ctx context.Context, in ContextInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolContext(ctx, nil, in)
}

func CallAdd(s *Server, ctx context.Context, in AddInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolAdd(ctx, nil, in)
}

func CallLink(s *Server, ctx context.Context, in LinkInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolLink(ctx, nil, in)
}

func CallTransition(s *Server, ctx context.Context, in TransitionInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolTransition(ctx, nil, in)
}

func CallReview(s *Server, ctx context.Context, in ReviewInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolReview(ctx, nil, in)
}

func CallTasks(s *Server, ctx context.Context, in TasksInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolTasks(ctx, nil, in)
}

func CallCapability(s *Server, ctx context.Context, in CapabilityInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolCapability(ctx, nil, in)
}

func CallVersion(s *Server, ctx context.Context, in VersionInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolVersion(ctx, nil, in)
}

func CallImpact(s *Server, ctx context.Context, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolImpact(ctx, nil, in)
}

func CallSearch(s *Server, ctx context.Context, in SearchInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolSearch(ctx, nil, in)
}

func CallChanges(s *Server, ctx context.Context, in ChangesInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolChanges(ctx, nil, in)
}

func CallRegressions(s *Server, ctx context.Context, in RegressionsInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolRegressions(ctx, nil, in)
}

func CallLoop(s *Server, ctx context.Context, in LoopInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolLoop(ctx, nil, in)
}

func CallAgents(s *Server, ctx context.Context, in AgentsInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolAgents(ctx, nil, in)
}

func CallPlan(s *Server, ctx context.Context, in PlanInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolPlan(ctx, nil, in)
}

func CallExplore(s *Server, ctx context.Context, in ExploreInput) (*sdkmcp.CallToolResult, any, error) {
	return s.toolExplore(ctx, nil, in)
}

// ResultText extracts the first text content block (test helper).
func ResultText(t testing.TB, res any) string {
	t.Helper()
	r, ok := res.(*sdkmcp.CallToolResult)
	if !ok || r == nil {
		t.Fatalf("expected *CallToolResult, got %T", res)
	}
	if len(r.Content) == 0 {
		t.Fatal("empty content")
	}
	tc, ok := r.Content[0].(*sdkmcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", r.Content[0])
	}
	return tc.Text
}

func TraceAddDescriptionForTest() string {
	return traceAddDescription
}
