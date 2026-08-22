package httpapi

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/mrchatam/Trace/internal/store"
)

const (
	APIVersion   = "1.0.0"
	TraceVersion = "0.0.0-dev"
)

// Options configures the HTTP adapter.
type Options struct {
	Root        string // absolute project root
	Addr        string // host:port; default 127.0.0.1:7432
	AllowRemote bool
	Token       string
	TokenFile   string
	StaticDir   string // default <root>/web/dist; must not equal project root
	CorsOrigin  string // optional exact Origin to reflect; never "*"; empty = deny
	// AddrExplicit is true when the operator passed --addr (flag.Changed).
	// When true, ListenAndServe fails on EADDRINUSE instead of UA-increment hopping.
	AddrExplicit bool
	// OnListening runs after TCP bind succeeds and before Serve (CLI open-browser, etc.).
	// The addr argument is the chosen host:port after any auto-port hop.
	OnListening func(addr string)
}

// Server is a thin Law-19 HTTP adapter over Trace libraries.
type Server struct {
	root         string
	addr         string // host:port for Listen
	host         string
	allowRemote  bool
	token        string
	requireToken bool
	staticDir    string
	corsOrigin   string
	addrExplicit bool
	onListening  func(addr string)
	mux          *http.ServeMux
	handler      http.Handler
}

// New validates options, resolves bind policy, and builds the handler tree.
// It does not listen. Call ListenAndServe or Handler() from tests.
func New(opts Options) (*Server, error) {
	root := strings.TrimSpace(opts.Root)
	if root == "" {
		return nil, fmt.Errorf("httpapi: Root is required")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("httpapi: resolve root: %w", err)
	}

	addr := strings.TrimSpace(opts.Addr)
	if addr == "" {
		addr = DefaultAddr
	}
	host, listenAddr, err := NormalizeListenAddr(addr)
	if err != nil {
		return nil, err
	}
	if err := RefuseRemote(host, opts.AllowRemote); err != nil {
		return nil, err
	}

	token := strings.TrimSpace(opts.Token)
	if token == "" && strings.TrimSpace(opts.TokenFile) != "" {
		token, err = LoadTokenFile(opts.TokenFile)
		if err != nil {
			return nil, err
		}
	}

	loopback := IsLoopbackHost(host)
	requireToken := false
	if !loopback {
		if token == "" {
			return nil, fmt.Errorf("httpapi: non-loopback bind requires --token, --token-file, or a generated token")
		}
		requireToken = true
	} else if token != "" {
		requireToken = true // operator opted in
	}

	staticDir := strings.TrimSpace(opts.StaticDir)
	if staticDir == "" {
		staticDir = filepath.Join(abs, "web", "dist")
	} else {
		staticDir, err = filepath.Abs(staticDir)
		if err != nil {
			return nil, fmt.Errorf("httpapi: resolve StaticDir: %w", err)
		}
	}
	// Refuse serving the project root as static (would expose .trace/ and source).
	if staticDir == abs {
		return nil, fmt.Errorf("httpapi: --static-dir must not be the project root (would expose .trace/ and source); use <root>/web/dist")
	}

	corsOrigin := strings.TrimSpace(opts.CorsOrigin)
	if corsOrigin == "*" {
		return nil, fmt.Errorf("httpapi: --cors-origin must be an exact Origin URL, not *")
	}

	s := &Server{
		root:         abs,
		addr:         listenAddr,
		host:         host,
		allowRemote:  opts.AllowRemote,
		token:        token,
		requireToken: requireToken,
		staticDir:    staticDir,
		corsOrigin:   corsOrigin,
		addrExplicit: opts.AddrExplicit,
		onListening:  opts.OnListening,
		mux:          http.NewServeMux(),
	}
	s.registerRoutes()
	s.handler = applyCORS(s.corsOrigin, s.authMiddleware(http.HandlerFunc(s.dispatch)))
	return s, nil
}

// dispatch routes /v1 to the mux; static/placeholder and /rpc stay outside mux patterns.
func (s *Server) dispatch(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" {
		path = "/"
	}
	if path == "/rpc" || strings.HasPrefix(path, "/rpc/") {
		s.handleRPCForbidden(w, r)
		return
	}
	if strings.HasPrefix(path, "/v1/") || path == "/v1" {
		s.mux.ServeHTTP(w, r)
		return
	}
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		s.serveStaticOrPlaceholder(w, r)
		return
	}
	writeEnvelope(w, http.StatusNotFound, "NOT_FOUND", "not found", nil)
}

// Handler returns the HTTP handler (for httptest).
func (s *Server) Handler() http.Handler { return s.handler }

// Addr returns the listen address (host:port).
func (s *Server) Addr() string { return s.addr }

// Token returns the configured bearer token (may be empty on loopback-trust).
func (s *Server) Token() string { return s.token }

// SetToken updates the bearer token used by auth middleware (e.g. POST /v1/auth/token).
func (s *Server) SetToken(tok string) {
	s.token = tok
	if tok != "" {
		s.requireToken = true
	}
}

// Root returns the absolute project root.
func (s *Server) Root() string { return s.root }

// ListenAndServe binds and serves until ctx is cancelled or an error occurs.
// When AddrExplicit is false, EADDRINUSE triggers UA-increment hops up to
// MaxAutoPortAttempts (same host, port+1). Explicit --addr fails on first busy.
func (s *Server) ListenAndServe(ctx context.Context) error {
	ln, err := s.listenTCP()
	if err != nil {
		return err
	}
	if s.onListening != nil {
		s.onListening(s.addr)
	}
	srv := &http.Server{
		Handler:           s.handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(ln)
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		<-errCh
		return ctx.Err()
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

// listenTCP binds s.addr, hopping on EADDRINUSE when !addrExplicit.
// On success, s.addr is the chosen host:port.
func (s *Server) listenTCP() (net.Listener, error) {
	maxAttempts := 1
	if !s.addrExplicit {
		maxAttempts = MaxAutoPortAttempts
	}
	start := s.addr
	for attempt := 0; attempt < maxAttempts; attempt++ {
		ln, err := net.Listen("tcp", s.addr)
		if err == nil {
			return ln, nil
		}
		if s.addrExplicit || !IsAddrInUse(err) {
			return nil, err
		}
		if attempt == maxAttempts-1 {
			break
		}
		next, ierr := IncrementListenPort(s.addr)
		if ierr != nil {
			return nil, ierr
		}
		s.addr = next
	}
	return nil, &AutoPortExhaustedError{Start: start, Attempts: maxAttempts}
}

func (s *Server) openStore() (*store.Store, error) {
	st, err := store.Open(s.root)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Server) registerRoutes() {
	// Meta
	s.mux.HandleFunc("GET /v1/health", s.handleHealth)
	s.mux.HandleFunc("GET /v1/version", s.handleVersion)
	s.mux.HandleFunc("GET /v1/project", s.handleProject)

	// Tasks
	s.mux.HandleFunc("GET /v1/tasks", s.handleListTasks)
	s.mux.HandleFunc("GET /v1/tasks/{task_id}", s.handleGetTask)

	// Loop
	s.mux.HandleFunc("GET /v1/loop/status", s.handleLoopStatus)
	s.mux.HandleFunc("GET /v1/loop/next", s.handleLoopNext)
	s.mux.HandleFunc("POST /v1/loop/apply", s.handleLoopApply)
	s.mux.HandleFunc("GET /v1/loop/gate", s.handleLoopGate)
	s.mux.HandleFunc("POST /v1/loop/reset", s.handleLoopReset)

	// Entities / links / transitions
	s.mux.HandleFunc("POST /v1/entities", s.handleCreateEntity)
	s.mux.HandleFunc("GET /v1/entities/{entity_id}", s.handleGetEntity)
	s.mux.HandleFunc("POST /v1/links", s.handleCreateLink)
	s.mux.HandleFunc("POST /v1/transitions", s.handleCreateTransition)

	// Retrieval
	s.mux.HandleFunc("GET /v1/context", s.handleContext)
	s.mux.HandleFunc("GET /v1/why", s.handleWhy)
	s.mux.HandleFunc("GET /v1/search", s.handleSearch)
	s.mux.HandleFunc("GET /v1/graph", s.handleGraph)

	// Seed
	s.mux.HandleFunc("GET /v1/seed/status", s.handleSeedStatus)
	s.mux.HandleFunc("POST /v1/seed/export", s.handleSeedExport)
	s.mux.HandleFunc("POST /v1/seed/import", s.handleSeedImport)

	// P1
	s.mux.HandleFunc("GET /v1/reviews", s.handleListReviews)
	s.mux.HandleFunc("POST /v1/reviews", s.handleCreateReview)
	s.mux.HandleFunc("GET /v1/reviews/{review_id}", s.handleGetReview)
	s.mux.HandleFunc("GET /v1/plans", s.handleListPlans)
	s.mux.HandleFunc("POST /v1/plans/bootstrap", s.handlePlanBootstrap)
	s.mux.HandleFunc("GET /v1/capability", s.handleListCapability)
	s.mux.HandleFunc("GET /v1/impact", s.handleGetImpact)
	s.mux.HandleFunc("GET /v1/changes", s.handleListChanges)
	s.mux.HandleFunc("GET /v1/regressions", s.handleListRegressions)
	s.mux.HandleFunc("GET /v1/agents", s.handleListAgents)
	s.mux.HandleFunc("GET /v1/index", s.handleIndexStatus)
	s.mux.HandleFunc("POST /v1/auth/token", s.handleAuthToken)

	// Defer → 501
	s.mux.HandleFunc("POST /v1/index", s.handleNotImplemented)
	s.mux.HandleFunc("POST /v1/backup", s.handleNotImplemented)
	s.mux.HandleFunc("POST /v1/restore", s.handleNotImplemented)
	s.mux.HandleFunc("POST /v1/migrate", s.handleNotImplemented)
	s.mux.HandleFunc("POST /v1/install", s.handleNotImplemented)
	s.mux.HandleFunc("GET /v1/patterns", s.handleNotImplemented)
	s.mux.HandleFunc("GET /v1/knowledge", s.handleNotImplemented)
	s.mux.HandleFunc("GET /v1/eval", s.handleNotImplemented)
	s.mux.HandleFunc("GET /v1/events", s.handleNotImplemented)
}
