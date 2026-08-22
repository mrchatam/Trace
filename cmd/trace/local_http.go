package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/mrchatam/Trace/internal/httpapi"
)

// localHTTPMode selects CLI wording and post-listen browser behavior.
type localHTTPMode string

const (
	localHTTPServe localHTTPMode = "serve"
	localHTTPGUI   localHTTPMode = "gui"
)

// notifyContext is the serve/gui cancel source; injectable for short-lived tests.
var notifyContext = func() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

// guiListenHookMu guards guiListenHook (tests swap it around short-lived servers).
var guiListenHookMu sync.Mutex
var guiListenHook = func(addr string) {}

func setGUIListenHook(fn func(addr string)) (restore func()) {
	guiListenHookMu.Lock()
	prev := guiListenHook
	if fn == nil {
		guiListenHook = func(addr string) {}
	} else {
		guiListenHook = fn
	}
	guiListenHookMu.Unlock()
	return func() {
		guiListenHookMu.Lock()
		guiListenHook = prev
		guiListenHookMu.Unlock()
	}
}

func callGUIListenHook(addr string) {
	guiListenHookMu.Lock()
	fn := guiListenHook
	guiListenHookMu.Unlock()
	fn(addr)
}
func cmdLocalHTTP(root string, args []string, mode localHTTPMode) int {
	name := string(mode)
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	addr := fs.String("addr", httpapi.DefaultAddr, "listen address (host:port); default 127.0.0.1:7432")
	allowRemote := fs.Bool("allow-remote", false, "allow non-loopback bind (requires bearer token)")
	token := fs.String("token", "", "bearer token (required off-loopback; optional on loopback)")
	tokenFile := fs.String("token-file", "", "read bearer token from file")
	rootAlias := fs.String("root", "", "project root alias for -C/--project (one root only)")
	staticDir := fs.String("static-dir", "", "GUI dist directory (default <root>/web/dist); refused if equal to project root")
	corsOrigin := fs.String("cors-origin", "", "optional exact Origin to reflect for CORS (e.g. http://127.0.0.1:5173); never *")
	var noOpen *bool
	if mode == localHTTPGUI {
		noOpen = fs.Bool("no-open", false, "do not open the default browser after listen")
	}

	usageServe := `usage: trace serve [-C|--project DIR] [--addr host:port] [--allow-remote] [--token TOKEN] [--token-file PATH] [--root DIR] [--static-dir DIR] [--cors-origin URL]

Opt-in local HTTP API. Default bind 127.0.0.1:7432 (loopback-trust).
Project root: cwd, or -C/--project before or after serve, or --root (same as other trace commands).

Consumer workflow (no Trace checkout required):
  cd your-project && trace gui
  # uses .trace/ under that project; GUI ships inside the trace binary (go:embed)

Non-loopback requires --allow-remote and a bearer token (--token, --token-file, or auto-generated).
Default busy → next free loopback port (7432–7441). --addr pins a port and fails if busy.

CORS is deny-by-default (never Access-Control-Allow-Origin: *). Optional --cors-origin
reflects only that exact Origin (Vite DX). Prefer Vite proxy to :7432 for same-origin.

--static-dir defaults to <root>/web/dist when present (contributor override). Otherwise
the Explore SPA is served from the trace binary embed — consumer projects need only .trace/.

Seed export with strict/task_id over HTTP returns 501 — use CLI for honesty/gate export.

`
	usageGUI := `usage: trace gui [-C|--project DIR] [--addr host:port] [--allow-remote] [--token TOKEN] [--token-file PATH] [--root DIR] [--static-dir DIR] [--cors-origin URL] [--no-open]

Start the local GUI HTTP server and open the default browser to Explore (/).
Same bind/token/CORS/project-root policy as trace serve. Consumer: cd your-project && trace gui.
Default busy → next free loopback port (7432–7441). --addr pins a port and fails if busy.

`
	fs.Usage = func() {
		if mode == localHTTPGUI {
			fmt.Fprint(os.Stderr, usageGUI)
		} else {
			fmt.Fprint(os.Stderr, usageServe)
		}
		fs.PrintDefaults()
	}

	projectRoot, serveArgs, err := peelServeProjectFlags(root, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		return exitUsage
	}

	flagDefs := map[string]bool{
		"addr": true, "token": true, "token-file": true, "root": true, "static-dir": true, "cors-origin": true,
	}
	if err := fs.Parse(flagsFirst(serveArgs, flagDefs)); err != nil {
		if err == flag.ErrHelp {
			return exitOK
		}
		return exitUsage
	}
	if fs.NArg() != 0 {
		if mode == localHTTPGUI {
			fmt.Fprint(os.Stderr, "usage: trace gui [--addr host:port] [--allow-remote] [--token TOKEN] [--token-file PATH] [--root DIR] [--static-dir DIR] [--cors-origin URL] [--no-open]\n")
		} else {
			fmt.Fprint(os.Stderr, "usage: trace serve [--addr host:port] [--allow-remote] [--token TOKEN] [--token-file PATH] [--root DIR] [--static-dir DIR] [--cors-origin URL]\n")
		}
		return exitUsage
	}

	if *rootAlias != "" {
		if projectRoot != "" {
			absGlobal, _ := filepath.Abs(projectRoot)
			absAlias, err := filepath.Abs(*rootAlias)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
				return exitFail
			}
			cwd, _ := os.Getwd()
			absCWD, _ := filepath.Abs(cwd)
			if absGlobal != absCWD && absGlobal != absAlias {
				fmt.Fprintf(os.Stderr, "%s: conflicting -C/--project and --root\n", name)
				return exitUsage
			}
		}
		projectRoot = *rootAlias
	}

	abs, err := resolveRoot(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		return exitFail
	}

	tok := *token
	generated := false
	host, _, err := httpapi.ParseListenAddr(*addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		return exitUsage
	}
	if err := httpapi.RefuseRemote(host, *allowRemote); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		return exitFail
	}
	if !httpapi.IsLoopbackHost(host) && tok == "" && *tokenFile == "" {
		t, gerr := httpapi.GenerateToken()
		if gerr != nil {
			fmt.Fprintf(os.Stderr, "%s: generate token: %v\n", name, gerr)
			return exitFail
		}
		tok = t
		generated = true
	}

	skipOpen := mode == localHTTPGUI && noOpen != nil && *noOpen
	// AddrExplicit: stdlib flag has no Flag.Changed; Visit == "was set on cmdline"
	// (same intent as PLAN flag.Changed — never DefaultAddr string-equal).
	addrExplicit := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "addr" {
			addrExplicit = true
		}
	})
	opts := httpapi.Options{
		Root:         abs,
		Addr:         *addr,
		AllowRemote:  *allowRemote,
		Token:        tok,
		TokenFile:    *tokenFile,
		StaticDir:    *staticDir,
		CorsOrigin:   *corsOrigin,
		AddrExplicit: addrExplicit,
	}
	opts.OnListening = func(listenAddr string) {
		base := fmt.Sprintf("http://%s", listenAddr)
		if mode == localHTTPServe {
			fmt.Fprintf(os.Stderr, "trace serve: listening on %s (root %s)\n", base, abs)
			return
		}
		land := base + "/"
		fmt.Fprintf(os.Stderr, "trace gui: listening on %s (root %s)\n", base, abs)
		callGUIListenHook(listenAddr)
		if skipOpen {
			return
		}
		if err := openBrowserFn(land); err != nil {
			fmt.Fprintf(os.Stderr, "trace gui: open browser: %v\n", err)
			tipOpenManuallyFn(land)
		}
	}

	srv, err := httpapi.New(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		return exitFail
	}
	if generated {
		fmt.Fprintf(os.Stdout, "trace %s: generated bearer token (print once):\n%s\n", name, srv.Token())
	}

	ctx, stop := notifyContext()
	defer stop()
	if err := srv.ListenAndServe(ctx); err != nil && err != context.Canceled {
		var exhausted *httpapi.AutoPortExhaustedError
		if errors.As(err, &exhausted) {
			fmt.Fprint(os.Stderr, err.Error())
		} else if httpapi.IsAddrInUse(err) {
			fmt.Fprint(os.Stderr, httpapi.FormatAddrInUseMessage(srv.Addr()))
		} else {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		}
		return exitFail
	}
	return exitOK
}
