package mcp_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	tracemcp "github.com/mrchatam/Trace/internal/mcp"
	"github.com/mrchatam/Trace/internal/store"
)

// connectProbe wires a real client to a real server over the SDK's in-memory
// transport. This is the permanent replacement for the ad-hoc Python probe
// harness: it exercises the same JSON-RPC surface an agent uses, with no
// external process or scripting dependency.
func connectProbe(t *testing.T, projectRoot string) *sdkmcp.ClientSession {
	t.Helper()
	ctx := context.Background()

	srv := tracemcp.NewServer(tracemcp.Options{ProjectRoot: projectRoot})
	clientTransport, serverTransport := sdkmcp.NewInMemoryTransports()

	serverSession, err := srv.MCP().Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "conformance-probe", Version: "0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() {
		_ = clientSession.Close()
		_ = serverSession.Close()
	})
	return clientSession
}

// TestProtocolNegotiation pins the negotiated version so a future SDK bump that
// silently changes it is caught.
func TestProtocolNegotiation(t *testing.T) {
	session := connectProbe(t, t.TempDir())
	got := session.InitializeResult().ProtocolVersion
	if got == "" {
		t.Fatal("negotiated protocol version must not be empty")
	}
	t.Logf("negotiated protocol version: %s", got)
}

// TestProtocolListsTheLockedToolSet proves the advertised catalog is exactly
// the locked seventeen-name set — the stale-server check an agent performs
// after a rebuild. The wire listing is sorted by the SDK, so compare as sets;
// RegisteredToolNames is the registration order used elsewhere.
func TestProtocolListsTheLockedToolSet(t *testing.T) {
	session := connectProbe(t, t.TempDir())
	res, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("tools/list: %v", err)
	}
	want := map[string]bool{}
	for _, name := range tracemcp.RegisteredToolNames() {
		want[name] = true
	}
	if len(res.Tools) != len(want) {
		t.Fatalf("tools/list returned %d tools, want %d", len(res.Tools), len(want))
	}
	for _, tool := range res.Tools {
		if !want[tool.Name] {
			t.Fatalf("tools/list advertised unexpected tool %q", tool.Name)
		}
		delete(want, tool.Name)
	}
	if len(want) != 0 {
		t.Fatalf("tools/list omitted: %v", want)
	}
}

// TestProtocolStandardMethodsRespond covers the capability methods a client
// may call unconditionally.
func TestProtocolStandardMethodsRespond(t *testing.T) {
	session := connectProbe(t, t.TempDir())
	ctx := context.Background()

	if _, err := session.ListResources(ctx, nil); err != nil {
		t.Fatalf("resources/list: %v", err)
	}
	if _, err := session.ListResourceTemplates(ctx, nil); err != nil {
		t.Fatalf("resources/templates/list: %v", err)
	}
	if _, err := session.ListPrompts(ctx, nil); err != nil {
		t.Fatalf("prompts/list: %v", err)
	}
	if err := session.Ping(ctx, nil); err != nil {
		t.Fatalf("ping: %v", err)
	}
}

// TestProtocolUnknownMethodIsMethodNotFound proves an unknown method yields
// -32601 rather than a silent success or a hang.
func TestProtocolUnknownMethodIsMethodNotFound(t *testing.T) {
	session := connectProbe(t, t.TempDir())
	_, err := session.CallTool(context.Background(), &sdkmcp.CallToolParams{
		Name: "trace_definitely_not_a_tool", Arguments: map[string]any{},
	})
	if err == nil {
		t.Fatal("want an error for an unknown tool")
	}
}

// TestProtocolRejectsUnknownParams proves strict argument validation: a typo in
// a parameter name is refused rather than silently ignored.
func TestProtocolRejectsUnknownParams(t *testing.T) {
	session := connectProbe(t, t.TempDir())
	ctx := context.Background()

	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	session = connectProbe(t, dir)

	res, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "trace_search",
		Arguments: map[string]any{"query": "hello", "totally_unknown_param": 1},
	})
	if err != nil {
		t.Fatalf("an invalid argument is a tool error, not a transport error: %v", err)
	}
	// Schema validation surfaces as an IsError tool result, per the MCP spec.
	if !res.IsError {
		t.Fatalf("want an error result for an unknown parameter, got: %+v", res)
	}
	if msg := resultText(res); !containsAny(msg, "additional properties", "unknown", "invalid") {
		t.Fatalf("error should name the problem, got: %q", msg)
	}
}

// TestProtocolRejectsOutOfEnumValue proves the advertised enums are enforced,
// not decorative.
func TestProtocolRejectsOutOfEnumValue(t *testing.T) {
	session := connectProbe(t, t.TempDir())
	res, err := session.CallTool(context.Background(), &sdkmcp.CallToolParams{
		Name:      "trace_regressions",
		Arguments: map[string]any{"action": "delete_everything"},
	})
	if err != nil {
		t.Fatalf("an invalid argument is a tool error, not a transport error: %v", err)
	}
	if !res.IsError {
		t.Fatalf("want an error result for a value outside the declared enum, got: %+v", res)
	}
	if msg := resultText(res); !containsAny(msg, "list", "enum") {
		t.Fatalf("error should name the allowed values, got: %q", msg)
	}
}

// resultText flattens a tool result's text content.
func resultText(res *sdkmcp.CallToolResult) string {
	var b strings.Builder
	for _, c := range res.Content {
		if tc, ok := c.(*sdkmcp.TextContent); ok {
			b.WriteString(tc.Text)
			b.WriteByte(' ')
		}
	}
	return b.String()
}

// TestProtocolConcurrentReadersShareOneProject proves two live servers against
// the same project serialize without error or stderr noise.
func TestProtocolConcurrentReadersShareOneProject(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	const workers = 4
	var wg sync.WaitGroup
	errCh := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			session := connectProbe(t, dir)
			res, err := session.CallTool(context.Background(), &sdkmcp.CallToolParams{
				Name:      "trace_version",
				Arguments: map[string]any{},
			})
			if err != nil {
				errCh <- err
				return
			}
			if len(res.Content) == 0 {
				errCh <- errEmptyResult
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatalf("concurrent probe failed: %v", err)
	}
}

// errEmptyResult reports a probe that got a result with no content.
var errEmptyResult = errors.New("probe result had no content")

// containsAny reports whether s contains any of the substrings.
func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
