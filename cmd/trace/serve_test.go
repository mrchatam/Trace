package main

import (
	"net"
	"strings"
	"testing"
)

func TestServeHelpFlags(t *testing.T) {
	out := captureStdout(t, func() int { return run([]string{"help"}) })
	for _, want := range []string{
		"serve",
		"--addr",
		"--allow-remote",
		"--token",
		"--cors-origin",
		"--static-dir",
		"7432–7441",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("help missing %q", want)
		}
	}
	if strings.Contains(out, "no auto free-port") {
		t.Fatalf("help must not claim no auto free-port:\n%s", out)
	}
}

func TestServeRefuseRemoteCLI(t *testing.T) {
	dir := t.TempDir()
	code := run([]string{"-C", dir, "serve", "--addr", "0.0.0.0:17999"})
	if code == exitOK {
		t.Fatal("serve without --allow-remote must fail")
	}
}

func TestServeAddrInUseFriendlyMessage(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	addr := ln.Addr().String()

	dir := t.TempDir()
	code, _, stderr := runCapture(t, []string{"-C", dir, "serve", "--addr", addr})
	if code == exitOK {
		t.Fatal("serve on occupied addr must fail")
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
