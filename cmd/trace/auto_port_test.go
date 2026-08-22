package main

import (
	"context"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mrchatam/Trace/internal/httpapi"
)

// T6 — explicit --addr equal to DefaultAddr string must fail when busy (flag.Changed).
func TestGuiExplicitDefaultAddrBusyNoHop(t *testing.T) {
	ln, err := net.Listen("tcp", httpapi.DefaultAddr)
	if err != nil {
		t.Skipf("cannot bind DefaultAddr (busy externally): %v", err)
	}
	defer ln.Close()

	dir := t.TempDir()
	code, _, stderr := runCapture(t, []string{"-C", dir, "gui", "--addr", httpapi.DefaultAddr, "--no-open"})
	if code == exitOK {
		t.Fatal("gui --addr DefaultAddr on occupied port must fail (no hop)")
	}
	for _, want := range []string{"address already in use", "--addr"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr missing %q (exit %d):\n%s", want, code, stderr)
		}
	}
	if strings.Contains(stderr, "no free port in auto range") {
		t.Fatalf("explicit path must not use auto-exhausted message:\n%s", stderr)
	}
}

func TestServeExplicitDefaultAddrBusyNoHop(t *testing.T) {
	ln, err := net.Listen("tcp", httpapi.DefaultAddr)
	if err != nil {
		t.Skipf("cannot bind DefaultAddr (busy externally): %v", err)
	}
	defer ln.Close()

	dir := t.TempDir()
	code, _, stderr := runCapture(t, []string{"-C", dir, "serve", "--addr", httpapi.DefaultAddr})
	if code == exitOK {
		t.Fatal("serve --addr DefaultAddr on occupied port must fail (no hop)")
	}
	for _, want := range []string{"address already in use", "--addr"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr missing %q (exit %d):\n%s", want, code, stderr)
		}
	}
}

// T5 — two concurrent default binds → distinct ports; gui hook/open sees chosen addr.
func TestGuiServeConcurrentDefaultDistinctPorts(t *testing.T) {
	var (
		mu      sync.Mutex
		cancels []context.CancelFunc
	)
	defer func() {
		mu.Lock()
		for _, c := range cancels {
			c()
		}
		mu.Unlock()
	}()

	prevNotify := notifyContext
	notifyContext = func() (context.Context, context.CancelFunc) {
		ctx, cancel := context.WithCancel(context.Background())
		mu.Lock()
		cancels = append(cancels, cancel)
		mu.Unlock()
		return ctx, cancel
	}
	defer func() { notifyContext = prevNotify }()

	guiHeard := make(chan string, 1)
	restoreHook := setGUIListenHook(func(addr string) {
		select {
		case guiHeard <- addr:
		default:
		}
	})
	defer restoreHook()

	opened := make(chan string, 1)
	prevOpen := openBrowserFn
	openBrowserFn = func(url string) error {
		opened <- url
		return nil
	}
	defer func() { openBrowserFn = prevOpen }()

	dir1 := t.TempDir()
	dir2 := t.TempDir()

	serveDone := make(chan int, 1)
	go func() {
		serveDone <- run([]string{"-C", dir1, "serve"})
	}()

	serveAddr, ok := waitDefaultRangeBound(t, 5*time.Second)
	if !ok {
		t.Fatal("timeout waiting for first serve bind")
	}

	guiDone := make(chan int, 1)
	go func() {
		guiDone <- run([]string{"-C", dir2, "gui"})
	}()

	var guiAddr string
	select {
	case guiAddr = <-guiHeard:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for gui listen hook")
	}

	var openURL string
	select {
	case openURL = <-opened:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for gui open URL")
	}

	if guiAddr == serveAddr {
		t.Fatalf("concurrent defaults must use distinct ports; both %q", guiAddr)
	}
	wantOpen := "http://" + guiAddr + "/"
	if openURL != wantOpen {
		t.Fatalf("open URL=%q want %q", openURL, wantOpen)
	}

	mu.Lock()
	for _, c := range cancels {
		c()
	}
	mu.Unlock()
	<-serveDone
	<-guiDone
}

func waitDefaultRangeBound(t *testing.T, timeout time.Duration) (string, bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		for p := 7432; p <= 7441; p++ {
			addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(p))
			c, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
			if err == nil {
				_ = c.Close()
				return addr, true
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	return "", false
}
