package httpapi

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
)

// IsAddrInUse reports whether err indicates the listen address is already bound.
// Prefers typed unwrap (net.OpError / syscall.EADDRINUSE); falls back to message match.
func IsAddrInUse(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, syscall.EADDRINUSE) {
		return true
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if errors.Is(opErr.Err, syscall.EADDRINUSE) {
			return true
		}
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) && errors.Is(sysErr.Err, syscall.EADDRINUSE) {
			return true
		}
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "address already in use")
}

// FormatAddrInUseMessage returns CLI stderr copy for an explicit --addr bind conflict.
// Wording covers both `trace gui` and `trace serve` (shared local HTTP).
// Default (non-explicit) bind auto-hops; this path is for pinned --addr only.
func FormatAddrInUseMessage(addr string) string {
	if strings.TrimSpace(addr) == "" {
		addr = DefaultAddr
	}
	return fmt.Sprintf(`gui|serve: address already in use (%s)
hint: another process (often trace gui or trace serve) is bound there.
  Default bind auto-hops to the next free port (7432–7441); --addr pins and fails if busy.
  To pin a free port, e.g.:
    trace gui -C /path/to/other --addr 127.0.0.1:7433
`, addr)
}
