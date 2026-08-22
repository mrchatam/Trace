package httpapi

import (
	"errors"
	"net"
	"strings"
	"testing"
)

func TestIsAddrInUse_realConflict(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	_, err = net.Listen("tcp", ln.Addr().String())
	if err == nil {
		t.Fatal("expected second listen to fail")
	}
	if !IsAddrInUse(err) {
		t.Fatalf("IsAddrInUse=false for %v", err)
	}
}

func TestIsAddrInUse_messageFallback(t *testing.T) {
	if !IsAddrInUse(errors.New("listen tcp 127.0.0.1:7432: bind: address already in use")) {
		t.Fatal("expected message fallback to match")
	}
	if IsAddrInUse(errors.New("connection refused")) {
		t.Fatal("non-in-use error must be false")
	}
	if IsAddrInUse(nil) {
		t.Fatal("nil must be false")
	}
}

func TestFormatAddrInUseMessage(t *testing.T) {
	msg := FormatAddrInUseMessage("127.0.0.1:7432")
	for _, want := range []string{
		"gui|serve:",
		"address already in use (127.0.0.1:7432)",
		"trace gui",
		"trace serve",
		"--addr",
		"127.0.0.1:7433",
		"auto-hops",
	} {
		if !strings.Contains(msg, want) {
			t.Fatalf("message missing %q:\n%s", want, msg)
		}
	}
}
