package main

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestGuiHelpFlags(t *testing.T) {
	out := captureStdout(t, func() int { return run([]string{"help"}) })
	for _, want := range []string{
		"gui",
		"--no-open",
		"go build -o bin/trace ./cmd/trace",
		"127.0.0.1:7432",
		"7432–7441",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("help missing %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "does not put the") {
		t.Fatalf("help should clarify trace install ≠ PATH:\n%s", out)
	}
	if strings.Contains(out, "no auto free-port") {
		t.Fatalf("help must not claim no auto free-port:\n%s", out)
	}

	code, _, stderr := runCapture(t, []string{"gui", "--help"})
	if code != exitOK {
		t.Fatalf("gui --help exit=%d stderr=%s", code, stderr)
	}
	for _, want := range []string{"--no-open", "127.0.0.1:7432", "gui", "7432–7441"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("gui --help missing %q:\n%s", want, stderr)
		}
	}
	if strings.Contains(stderr, "does not auto-pick") {
		t.Fatalf("gui --help must not claim no auto-pick:\n%s", stderr)
	}
}

func TestGuiRefuseRemoteCLI(t *testing.T) {
	dir := t.TempDir()
	code := run([]string{"-C", dir, "gui", "--addr", "0.0.0.0:17998"})
	if code == exitOK {
		t.Fatal("gui without --allow-remote must fail")
	}
}

func TestGuiAddrInUseFriendlyMessage(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	addr := ln.Addr().String()

	dir := t.TempDir()
	code, _, stderr := runCapture(t, []string{"-C", dir, "gui", "--addr", addr, "--no-open"})
	if code == exitOK {
		t.Fatal("gui on occupied addr must fail")
	}
	for _, want := range []string{
		"address already in use",
		"--addr",
		"127.0.0.1:7433",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr missing %q (exit %d):\n%s", want, code, stderr)
		}
	}
}

func freeLoopbackAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	return addr
}

func withGUITestServer(t *testing.T) (cancel context.CancelFunc, restoreNotify func()) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	prevNotify := notifyContext
	notifyContext = func() (context.Context, context.CancelFunc) {
		return ctx, cancel
	}
	return cancel, func() { notifyContext = prevNotify }
}

func TestGuiNoOpenDoesNotCallOpener(t *testing.T) {
	addr := freeLoopbackAddr(t)

	var calls atomic.Int32
	prevOpen := openBrowserFn
	openBrowserFn = func(url string) error {
		calls.Add(1)
		return nil
	}
	defer func() { openBrowserFn = prevOpen }()

	cancel, restoreNotify := withGUITestServer(t)
	defer restoreNotify()

	listening := make(chan struct{}, 1)
	restoreHook := setGUIListenHook(func(string) {
		select {
		case listening <- struct{}{}:
		default:
		}
	})
	defer restoreHook()

	dir := t.TempDir()
	done := make(chan int, 1)
	go func() {
		done <- run([]string{"-C", dir, "gui", "--addr", addr, "--no-open"})
	}()

	select {
	case <-listening:
	case <-time.After(5 * time.Second):
		cancel()
		t.Fatal("timeout waiting for listen")
	}
	cancel()
	code := <-done
	if code != exitOK {
		t.Fatalf("gui --no-open exit=%d", code)
	}
	if calls.Load() != 0 {
		t.Fatalf("opener called %d times with --no-open", calls.Load())
	}
}

func TestGuiOpenURLLandsOnExploreRoot(t *testing.T) {
	addr := freeLoopbackAddr(t)

	opened := make(chan string, 1)
	prevOpen := openBrowserFn
	openBrowserFn = func(url string) error {
		opened <- url
		return nil
	}
	defer func() { openBrowserFn = prevOpen }()

	cancel, restoreNotify := withGUITestServer(t)
	defer restoreNotify()

	dir := t.TempDir()
	done := make(chan int, 1)
	go func() {
		done <- run([]string{"-C", dir, "gui", "--addr", addr})
	}()

	var url string
	select {
	case url = <-opened:
	case <-time.After(5 * time.Second):
		cancel()
		t.Fatal("timeout waiting for openBrowser")
	}
	cancel()
	code := <-done
	if code != exitOK {
		t.Fatalf("gui exit=%d", code)
	}
	want := "http://" + addr + "/"
	if url != want {
		t.Fatalf("open URL=%q want %q", url, want)
	}
	if strings.Contains(url, "overview") {
		t.Fatalf("must not land on /overview: %q", url)
	}
}

func TestGuiOpenFailStillListens(t *testing.T) {
	addr := freeLoopbackAddr(t)

	prevOpen := openBrowserFn
	openBrowserFn = func(url string) error {
		return errors.New("mock open fail")
	}
	defer func() { openBrowserFn = prevOpen }()

	tipped := make(chan string, 1)
	prevTip := tipOpenManuallyFn
	tipOpenManuallyFn = func(url string) { tipped <- url }
	defer func() { tipOpenManuallyFn = prevTip }()

	cancel, restoreNotify := withGUITestServer(t)
	defer restoreNotify()

	listening := make(chan struct{}, 1)
	restoreHook := setGUIListenHook(func(string) {
		select {
		case listening <- struct{}{}:
		default:
		}
	})
	defer restoreHook()

	dir := t.TempDir()
	done := make(chan int, 1)
	go func() {
		done <- run([]string{"-C", dir, "gui", "--addr", addr})
	}()

	select {
	case <-listening:
	case <-time.After(5 * time.Second):
		cancel()
		t.Fatal("timeout waiting for listen after open fail")
	}
	var tipURL string
	select {
	case tipURL = <-tipped:
	case <-time.After(2 * time.Second):
		cancel()
		t.Fatal("expected open-fail tip")
	}
	cancel()
	code := <-done
	if code != exitOK {
		t.Fatalf("open fail must not fail listen; exit=%d", code)
	}
	want := "http://" + addr + "/"
	if tipURL != want {
		t.Fatalf("tip URL=%q want %q", tipURL, want)
	}
}
