package httpapi

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func occupyTCP(t *testing.T, addr string) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Skipf("cannot occupy %s (busy externally): %v", addr, err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	return ln
}

func newAutoPortServer(t *testing.T, addr string, explicit bool, onListening func(string)) *Server {
	t.Helper()
	srv, err := New(Options{
		Root:         t.TempDir(),
		Addr:         addr,
		AddrExplicit: explicit,
		OnListening:  onListening,
	})
	if err != nil {
		t.Fatal(err)
	}
	return srv
}

func TestListenAutoPort_freeDefaultStays7432(t *testing.T) {
	// T11 — no hop when DefaultAddr is free. Do not Parallel: binds 7432–7441.
	occupyTCP(t, "127.0.0.1:7433") // unrelated busy must not force hop from free 7432

	var heard string
	srv := newAutoPortServer(t, DefaultAddr, false, func(a string) { heard = a })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe(ctx) }()

	deadline := time.After(3 * time.Second)
	for heard == "" {
		select {
		case err := <-errCh:
			t.Fatalf("ListenAndServe failed early: %v", err)
		case <-deadline:
			t.Fatal("timeout waiting for OnListening")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	cancel()
	<-errCh

	if heard != DefaultAddr {
		t.Fatalf("OnListening addr=%q want %q", heard, DefaultAddr)
	}
	if srv.Addr() != DefaultAddr {
		t.Fatalf("Addr()=%q want %q", srv.Addr(), DefaultAddr)
	}
}

func TestListenAutoPort_busyDefaultHopsNext(t *testing.T) {
	// T4 — occupy :7432 → bind :7433.
	occupyTCP(t, DefaultAddr)

	var heard string
	srv := newAutoPortServer(t, DefaultAddr, false, func(a string) { heard = a })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe(ctx) }()

	deadline := time.After(3 * time.Second)
	for heard == "" {
		select {
		case err := <-errCh:
			t.Fatalf("ListenAndServe failed early: %v", err)
		case <-deadline:
			t.Fatal("timeout waiting for OnListening")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	cancel()
	<-errCh

	want := "127.0.0.1:7433"
	if heard != want {
		t.Fatalf("OnListening addr=%q want %q", heard, want)
	}
	if srv.Addr() != want {
		t.Fatalf("Addr()=%q want %q", srv.Addr(), want)
	}
}

func TestListenAutoPort_explicitBusyNoHop(t *testing.T) {
	occupyTCP(t, DefaultAddr)

	srv := newAutoPortServer(t, DefaultAddr, true, nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := srv.ListenAndServe(ctx)
	if err == nil {
		t.Fatal("explicit --addr on busy DefaultAddr must fail")
	}
	if !IsAddrInUse(err) {
		t.Fatalf("want IsAddrInUse, got %v", err)
	}
	var exhausted *AutoPortExhaustedError
	if errors.As(err, &exhausted) {
		t.Fatal("explicit path must not return AutoPortExhaustedError")
	}
	if srv.Addr() != DefaultAddr {
		t.Fatalf("Addr() hopped to %q; explicit must stay on first addr", srv.Addr())
	}
}

func TestListenAutoPort_exhaustedRange(t *testing.T) {
	// T7 — occupy 7432–7441 → fail with range + --addr.
	for p := 7432; p <= 7441; p++ {
		occupyTCP(t, net.JoinHostPort("127.0.0.1", strconv.Itoa(p)))
	}

	srv := newAutoPortServer(t, DefaultAddr, false, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := srv.ListenAndServe(ctx)
	if err == nil {
		t.Fatal("expected auto-port exhausted")
	}
	var exhausted *AutoPortExhaustedError
	if !errors.As(err, &exhausted) {
		t.Fatalf("want AutoPortExhaustedError, got %T %v", err, err)
	}
	msg := err.Error()
	for _, want := range []string{"7432", "7441", "10", "--addr"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("exhausted message missing %q:\n%s", want, msg)
		}
	}
}

func TestFormatAutoPortExhaustedMessage(t *testing.T) {
	msg := FormatAutoPortExhaustedMessage(DefaultAddr, MaxAutoPortAttempts)
	for _, want := range []string{"7432", "7441", "10", "--addr", "gui|serve:"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("message missing %q:\n%s", want, msg)
		}
	}
}

func TestIncrementListenPort(t *testing.T) {
	got, err := IncrementListenPort("127.0.0.1:7432")
	if err != nil {
		t.Fatal(err)
	}
	if got != "127.0.0.1:7433" {
		t.Fatalf("got %q", got)
	}
}
