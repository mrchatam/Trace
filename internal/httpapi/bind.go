package httpapi

import (
	"fmt"
	"net"
	"strings"
)

const (
	DefaultAddr = "127.0.0.1:7432"
	// MaxAutoPortAttempts is the UA-increment window for default (non-explicit) bind:
	// ports DefaultAddr … DefaultAddr+9 inclusive (7432–7441).
	MaxAutoPortAttempts = 10
)

// ParseListenAddr splits host:port. Empty host defaults to 127.0.0.1.
func ParseListenAddr(addr string) (host, port string, err error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		addr = DefaultAddr
	}
	h, p, err := net.SplitHostPort(addr)
	if err != nil {
		// Allow ":7432" form.
		if strings.HasPrefix(addr, ":") {
			return "0.0.0.0", strings.TrimPrefix(addr, ":"), nil
		}
		return "", "", fmt.Errorf("httpapi: invalid --addr %q: %w", addr, err)
	}
	if h == "" {
		h = "0.0.0.0"
	}
	return h, p, nil
}

// IsLoopbackHost reports whether host is loopback (127.0.0.1, ::1, localhost, IPv4-mapped).
func IsLoopbackHost(host string) bool {
	host = strings.TrimSpace(host)
	if host == "" {
		return false
	}
	lower := strings.ToLower(host)
	if lower == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}

// RefuseRemote returns an error when host is non-loopback and allowRemote is false.
func RefuseRemote(host string, allowRemote bool) error {
	if IsLoopbackHost(host) {
		return nil
	}
	if allowRemote {
		return nil
	}
	return fmt.Errorf("httpapi: refusing non-loopback bind %q without --allow-remote", host)
}

// NormalizeListenAddr returns host:port suitable for net.Listen.
func NormalizeListenAddr(addr string) (string, string, error) {
	host, port, err := ParseListenAddr(addr)
	if err != nil {
		return "", "", err
	}
	return host, net.JoinHostPort(host, port), nil
}
