package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/httpapi"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

func TestBindRefuseRemote(t *testing.T) {
	if err := httpapi.RefuseRemote("127.0.0.1", false); err != nil {
		t.Fatalf("loopback should allow: %v", err)
	}
	if err := httpapi.RefuseRemote("::1", false); err != nil {
		t.Fatalf("::1 should allow: %v", err)
	}
	if err := httpapi.RefuseRemote("0.0.0.0", false); err == nil {
		t.Fatal("0.0.0.0 without AllowRemote must refuse")
	}
	if !strings.Contains(errString(httpapi.RefuseRemote("0.0.0.0", false)), "--allow-remote") {
		t.Fatal("refusal must name --allow-remote")
	}
	if err := httpapi.RefuseRemote("0.0.0.0", true); err != nil {
		t.Fatalf("AllowRemote should permit: %v", err)
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func TestNewRefuseWithoutAllowRemote(t *testing.T) {
	dir := t.TempDir()
	_, err := httpapi.New(httpapi.Options{Root: dir, Addr: "0.0.0.0:0", AllowRemote: false})
	if err == nil {
		t.Fatal("expected refuse")
	}
}

func TestNewRemoteRequiresToken(t *testing.T) {
	dir := t.TempDir()
	_, err := httpapi.New(httpapi.Options{Root: dir, Addr: "0.0.0.0:0", AllowRemote: true})
	if err == nil {
		t.Fatal("expected token required")
	}
	srv, err := httpapi.New(httpapi.Options{
		Root: dir, Addr: "0.0.0.0:0", AllowRemote: true, Token: "secret-token",
	})
	if err != nil {
		t.Fatal(err)
	}
	if srv.Token() != "secret-token" {
		t.Fatal("token not set")
	}
}

func openFixture(t *testing.T) (dir string, srv *httpapi.Server) {
	t.Helper()
	dir = t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	srv, err = httpapi.New(httpapi.Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatal(err)
	}
	return dir, srv
}

func TestHealthVersionStaticCORS(t *testing.T) {
	_, srv := openFixture(t)
	h := srv.Handler()

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
	if rr.Code != 200 {
		t.Fatalf("health: %d %s", rr.Code, rr.Body.String())
	}
	var health map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &health)
	if health["ok"] != true {
		t.Fatalf("health body: %v", health)
	}
	if acao := rr.Header().Get("Access-Control-Allow-Origin"); acao == "*" {
		t.Fatal("CORS must not be *")
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/version", nil))
	if rr.Code != 200 {
		t.Fatal(rr.Body.String())
	}
	var ver map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &ver)
	if ver["name"] != "trace" || ver["api_version"] != "1.0.0" {
		t.Fatalf("version: %v", ver)
	}

	rr = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://evil.example")
	h.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("GET /: %d", rr.Code)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Fatalf("Content-Type: %s", ct)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") == "*" {
		t.Fatal("CORS * forbidden")
	}
	body := rr.Body.String()
	if strings.Contains(body, "Embedded GUI stub") {
		t.Fatal("must not serve stub")
	}
	if !strings.Contains(body, `id="root"`) {
		t.Fatal("GET / should serve embedded Explore SPA (#root)")
	}
}

func TestRPCNotMCP(t *testing.T) {
	_, srv := openFixture(t)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/rpc", strings.NewReader(`{}`)))
	if rr.Code == 200 {
		t.Fatal("/rpc must not succeed as MCP")
	}
	if rr.Code != 404 && rr.Code != 501 {
		t.Fatalf("want 404/501 got %d", rr.Code)
	}
}

func TestAuthLoopbackAndRemote(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	// Loopback without token: health + data OK
	srv, err := httpapi.New(httpapi.Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/tasks", nil))
	if rr.Code != 200 {
		t.Fatalf("loopback tasks: %d %s", rr.Code, rr.Body.String())
	}

	// Simulated remote: RequireBearer via non-loopback options
	remote, err := httpapi.New(httpapi.Options{
		Root: dir, Addr: "0.0.0.0:0", AllowRemote: true, Token: "tok-abc",
	})
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	remote.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/tasks", nil))
	if rr.Code != 401 {
		t.Fatalf("want 401 got %d", rr.Code)
	}
	var env map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &env)
	errObj := env["error"].(map[string]any)
	if errObj["code"] != "UNAUTHORIZED" {
		t.Fatalf("code: %v", errObj)
	}

	rr = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	req.Header.Set("Authorization", "Bearer tok-abc")
	remote.Handler().ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("bearer health: %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	req.Header.Set("Authorization", "Bearer tok-abc")
	remote.Handler().ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("bearer tasks: %d %s", rr.Code, rr.Body.String())
	}
}

func TestReadsAndWrites(t *testing.T) {
	dir, srv := openFixture(t)
	h := srv.Handler()

	// create goal + task + link
	postJSON := func(path string, body any) *httptest.ResponseRecorder {
		t.Helper()
		b, _ := json.Marshal(body)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b)))
		return rr
	}

	rr := postJSON("/v1/entities", map[string]any{"kind": "goal", "title": "G1"})
	if rr.Code != 201 {
		t.Fatalf("create goal: %d %s", rr.Code, rr.Body.String())
	}
	var goal map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &goal)
	goalID := goal["id"].(string)

	rr = postJSON("/v1/entities", map[string]any{"kind": "task", "title": "T1", "goal_id": goalID})
	if rr.Code != 201 {
		t.Fatalf("create task: %d %s", rr.Code, rr.Body.String())
	}
	var task map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &task)
	taskID := task["id"].(string)

	rr = postJSON("/v1/links", map[string]any{"rel": "goal-task", "from": goalID, "to": taskID})
	if rr.Code != 201 {
		t.Fatalf("link: %d %s", rr.Code, rr.Body.String())
	}

	// reads
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/version", nil))
	if rr.Code != 200 {
		t.Fatal("version")
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/tasks", nil))
	if rr.Code != 200 {
		t.Fatalf("tasks: %s", rr.Body.String())
	}
	var list map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &list)
	items := list["items"].([]any)
	if len(items) < 1 {
		t.Fatal("expected tasks")
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/loop/status?task_id="+taskID, nil))
	if rr.Code != 200 {
		t.Fatalf("loop status: %d %s", rr.Code, rr.Body.String())
	}
	var status map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &status)
	if status["schema_version"] != "trace.loop.status.v1" {
		t.Fatalf("schema: %v", status["schema_version"])
	}

	// persist check via store
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	got, err := st.GetTask(taskID)
	if err != nil || got.Title != "T1" {
		t.Fatalf("persisted task: %v %v", got, err)
	}
}

func TestGraphBudgetAndDefer(t *testing.T) {
	dir, srv := openFixture(t)
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()
	g, err := svc.CreateGoal(ctx, domain.GoalInput{Title: "G"})
	if err != nil {
		t.Fatal(err)
	}
	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "T", GoalID: &g.ID})
	if err != nil {
		t.Fatal(err)
	}
	_ = svc.LinkGoalTask(ctx, g.ID, task.ID, domain.LinkMeta{})
	st.Close()

	h := srv.Handler()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/graph", nil))
	if rr.Code != 400 {
		t.Fatalf("missing params: %d", rr.Code)
	}
	var env map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &env)
	if env["error"].(map[string]any)["code"] != "VALIDATION_ERROR" {
		t.Fatalf("%v", env)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/graph?center="+task.ID+"&max_nodes=50", nil))
	if rr.Code != 200 {
		t.Fatalf("graph: %d %s", rr.Code, rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), `"edges":null`) {
		t.Fatalf("GET /v1/graph must not emit edges:null: %s", rr.Body.String())
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/graph?mode=project&max_nodes=50", nil))
	if rr.Code != 200 {
		t.Fatalf("project graph: %d %s", rr.Code, rr.Body.String())
	}
	var proj map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &proj)
	if proj["mode"] != "project" {
		t.Fatalf("mode: %#v", proj["mode"])
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/backup", nil))
	if rr.Code != 501 {
		t.Fatalf("defer: %d", rr.Code)
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &env)
	if env["error"].(map[string]any)["code"] != "NOT_IMPLEMENTED" {
		t.Fatalf("%v", env)
	}
}

func TestSeedPathConfinement(t *testing.T) {
	_, srv := openFixture(t)
	h := srv.Handler()
	body, _ := json.Marshal(map[string]any{"input_path": "../outside.json"})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/seed/import", bytes.NewReader(body)))
	if rr.Code != 400 {
		t.Fatalf("want 400 got %d %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "escapes") && !strings.Contains(rr.Body.String(), "VALIDATION") {
		t.Fatalf("body: %s", rr.Body.String())
	}
}

func TestListenLoopbackWithToken(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()

	srv, err := httpapi.New(httpapi.Options{
		Root: dir, Addr: addr, AllowRemote: true, Token: "listen-tok",
	})
	// Addr may still be 127.0.0.1 — AllowRemote with loopback is fine; use 0.0.0.0 for remote listen test
	if err != nil {
		// if addr is loopback, try 0.0.0.0:0 via ListenAndServe path differently
	}
	_ = srv

	remote, err := httpapi.New(httpapi.Options{
		Root: dir, Addr: "127.0.0.1:0", AllowRemote: false, Token: "x",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	go func() { errCh <- remote.ListenAndServe(ctx) }()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case <-errCh:
	case <-time.After(2 * time.Second):
		t.Fatal("listen did not stop")
	}
}

func TestNoSQLInPackage(t *testing.T) {
	root := filepath.Join("..", "..", "internal", "httpapi")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Skip(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		if strings.Contains(s, "database/sql") || strings.Contains(s, "modernc.org/sqlite") {
			t.Fatalf("%s imports SQL — Law 19 violation", e.Name())
		}
	}
}

func TestSeedExportSummaryOnly(t *testing.T) {
	dir, srv := openFixture(t)
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = domain.New(st).CreateGoal(context.Background(), domain.GoalInput{Title: "G"})
	st.Close()

	rr := httptest.NewRecorder()
	body := `{"output_path":"trace/graph.json"}`
	srv.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/seed/export", strings.NewReader(body)))
	if rr.Code != 200 {
		t.Fatalf("export: %d %s", rr.Code, rr.Body.String())
	}
	var job map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &job)
	if job["status"] != "ok" {
		t.Fatalf("%v", job)
	}
	if _, ok := job["goals"]; ok {
		t.Fatal("HTTP body must not embed full graph")
	}
	if _, err := os.Stat(filepath.Join(dir, "trace", "graph.json")); err != nil {
		t.Fatal("expected file write")
	}

	rr = httptest.NewRecorder()
	strictBody := `{"output_path":"trace/graph.json","strict":true}`
	srv.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/seed/export", strings.NewReader(strictBody)))
	if rr.Code != 501 {
		t.Fatalf("strict export want 501 got %d %s", rr.Code, rr.Body.String())
	}
	var env map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &env)
	if env["error"].(map[string]any)["code"] != "NOT_IMPLEMENTED" {
		t.Fatalf("%v", env)
	}
}

func TestAgentsListAndRecommend(t *testing.T) {
	_, srv := openFixture(t)
	h := srv.Handler()

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/agents", nil))
	if rr.Code != 200 {
		t.Fatalf("list: %d %s", rr.Code, rr.Body.String())
	}
	var list map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &list)
	if _, ok := list["items"]; !ok {
		t.Fatalf("want items: %v", list)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/agents?action=recommend", nil))
	if rr.Code != 400 {
		t.Fatalf("recommend without seed: %d %s", rr.Code, rr.Body.String())
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/agents?action=recommend&phase=CRITIQUE", nil))
	if rr.Code != 200 {
		t.Fatalf("recommend: %d %s", rr.Code, rr.Body.String())
	}
}

func TestCORSExactOriginReflect(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	allowed := "http://127.0.0.1:5173"
	srv, err := httpapi.New(httpapi.Options{
		Root: dir, Addr: "127.0.0.1:0", CorsOrigin: allowed,
	})
	if err != nil {
		t.Fatal(err)
	}
	h := srv.Handler()

	// Matching Origin → reflect
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	req.Header.Set("Origin", allowed)
	h.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("health: %d", rr.Code)
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != allowed {
		t.Fatalf("ACAO want %q got %q", allowed, got)
	}
	if rr.Header().Get("Vary") != "Origin" {
		t.Fatal("want Vary: Origin")
	}

	// Wrong Origin → no *
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	req.Header.Set("Origin", "https://evil.example")
	h.ServeHTTP(rr, req)
	if acao := rr.Header().Get("Access-Control-Allow-Origin"); acao == "*" || acao == "https://evil.example" {
		t.Fatalf("wrong origin must not be reflected: %q", acao)
	}

	// Preflight with matching Origin
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodOptions, "/v1/tasks", nil)
	req.Header.Set("Origin", allowed)
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Authorization")
	h.ServeHTTP(rr, req)
	if rr.Code != 204 {
		t.Fatalf("preflight: %d", rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != allowed {
		t.Fatal("preflight ACAO")
	}
	if !strings.Contains(rr.Header().Get("Access-Control-Allow-Headers"), "Authorization") {
		t.Fatal("preflight Allow-Headers")
	}

	// CorsOrigin * refused at New
	_, err = httpapi.New(httpapi.Options{Root: dir, Addr: "127.0.0.1:0", CorsOrigin: "*"})
	if err == nil {
		t.Fatal("CorsOrigin * must refuse")
	}
}

func TestStaticDirRefuseProjectRoot(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	_, err = httpapi.New(httpapi.Options{Root: dir, Addr: "127.0.0.1:0", StaticDir: dir})
	if err == nil {
		t.Fatal("StaticDir == root must refuse")
	}
	if !strings.Contains(err.Error(), ".trace") && !strings.Contains(err.Error(), "project root") {
		t.Fatalf("refusal should name footgun: %v", err)
	}
}

func TestMapDomainErrInvalidUUID(t *testing.T) {
	_, srv := openFixture(t)
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/loop/status?task_id=rl-not-a-uuid", nil))
	if rr.Code != 400 {
		t.Fatalf("want 400 got %d %s", rr.Code, rr.Body.String())
	}
	var env map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &env)
	errObj := env["error"].(map[string]any)
	if errObj["code"] != "VALIDATION_ERROR" {
		t.Fatalf("code: %v", errObj)
	}
	msg, _ := errObj["message"].(string)
	if !strings.Contains(msg, "must be UUID") {
		t.Fatalf("message: %q", msg)
	}
}

// T1 — consumer root with .trace/ only (no web/) serves embedded real SPA.
func TestStaticCSPAndEmbedFallback(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	// Default StaticDir is <root>/web/dist — missing → embedded real SPA
	srv, err := httpapi.New(httpapi.Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != 200 {
		t.Fatalf("GET /: %d", rr.Code)
	}
	csp := rr.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "default-src 'self'") || !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Fatalf("CSP: %q", csp)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") == "*" {
		t.Fatal("CORS * forbidden on static")
	}
	body := rr.Body.String()
	if strings.Contains(body, "Embedded GUI stub") {
		t.Fatal("embedded index must not be the stub phrase")
	}
	if !strings.Contains(body, `id="root"`) {
		t.Fatalf("want real SPA #root marker, body: %s", body)
	}
	if !strings.Contains(body, "/assets/") || !strings.Contains(body, `type="module"`) {
		t.Fatalf("want /assets/ module script, body: %s", body)
	}
}

// T2 — planted disk web/dist/index.html wins over embed.
func TestStaticDiskWinsOverEmbed(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()
	dist := filepath.Join(dir, "web", "dist")
	if err := os.MkdirAll(dist, 0o755); err != nil {
		t.Fatal(err)
	}
	const marker = "P34-DISK-WINS-MARKER-7f3a"
	diskHTML := `<!DOCTYPE html><html><body><p>` + marker + `</p></body></html>`
	if err := os.WriteFile(filepath.Join(dist, "index.html"), []byte(diskHTML), 0o644); err != nil {
		t.Fatal(err)
	}
	srv, err := httpapi.New(httpapi.Options{Root: dir, Addr: "127.0.0.1:0"})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != 200 {
		t.Fatalf("GET /: %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, marker) {
		t.Fatalf("disk marker missing: %s", body)
	}
	if strings.Contains(body, `id="root"`) {
		t.Fatal("must not serve embedded SPA when disk index exists")
	}
}

func readBody(t *testing.T, rr *httptest.ResponseRecorder) string {
	t.Helper()
	b, _ := io.ReadAll(rr.Body)
	return string(b)
}

func TestHTTPPlanBootstrap_CreatesPlannerRows(t *testing.T) {
	dir, srv := openFixture(t)
	h := srv.Handler()
	ctx := context.Background()

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	dsvc := domain.New(st)
	g, err := dsvc.CreateGoal(ctx, domain.GoalInput{Title: "HTTP bootstrap goal"})
	if err != nil {
		t.Fatal(err)
	}
	pc, err := dsvc.CreatePlanChange(ctx, domain.PlanChangeInput{Title: "bootstrap pc"})
	if err != nil {
		t.Fatal(err)
	}
	_ = pc
	st.Close()

	body, _ := json.Marshal(map[string]string{"goal_id": g.ID})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/plans/bootstrap", bytes.NewReader(body)))
	if rr.Code != http.StatusOK {
		t.Fatalf("bootstrap: %d %s", rr.Code, rr.Body.String())
	}
	var res map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res["goal_id"] != g.ID {
		t.Fatalf("goal_id: %#v", res["goal_id"])
	}
	if res["already_exists"] == true {
		t.Fatalf("expected fresh bootstrap, got already_exists")
	}
	if scopeID, _ := res["scope_id"].(string); scopeID == "" {
		t.Fatalf("scope_id missing: %#v", res)
	}

	st2, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	view, err := planner.New(st2).GetPlan(ctx, g.ID)
	if err != nil {
		t.Fatal(err)
	}
	if view.CurrentScopeID == nil || *view.CurrentScopeID == "" || view.CurrentDeepPlan == nil {
		t.Fatalf("planner rows not created: %+v", view)
	}
}

func TestHTTPIndexStatusLanguages(t *testing.T) {
	_, srv := openFixture(t)
	h := srv.Handler()

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/index", nil))
	if rr.Code != 200 {
		t.Fatalf("index status: %d %s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	raw, ok := out["supported_languages"].([]any)
	if !ok {
		t.Fatalf("supported_languages missing or wrong type: %#v", out["supported_languages"])
	}
	if len(raw) != 5 {
		t.Fatalf("len=%d want 5: %v", len(raw), raw)
	}
}
