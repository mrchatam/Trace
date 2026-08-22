package httpapi

import (
	"fmt"
	"net"
	"strconv"
)

// AutoPortExhaustedError is returned when default bind hops MaxAutoPortAttempts
// times and every candidate is busy.
type AutoPortExhaustedError struct {
	Start    string // first listen addr tried (host:port)
	Attempts int
}

func (e *AutoPortExhaustedError) Error() string {
	return FormatAutoPortExhaustedMessage(e.Start, e.Attempts)
}

// IncrementListenPort returns host:port with port+1 (same host).
func IncrementListenPort(addr string) (string, error) {
	host, portStr, err := ParseListenAddr(addr)
	if err != nil {
		return "", err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", fmt.Errorf("httpapi: invalid port in %q: %w", addr, err)
	}
	if port < 0 || port >= 65535 {
		return "", fmt.Errorf("httpapi: cannot increment port %d", port)
	}
	return net.JoinHostPort(host, strconv.Itoa(port+1)), nil
}

// FormatAutoPortExhaustedMessage is CLI stderr copy when the default auto range is full.
func FormatAutoPortExhaustedMessage(startAddr string, attempts int) string {
	if attempts <= 0 {
		attempts = MaxAutoPortAttempts
	}
	host, portStr, err := ParseListenAddr(startAddr)
	if err != nil || host == "" {
		host = "127.0.0.1"
		portStr = "7432"
	}
	startPort, _ := strconv.Atoi(portStr)
	endPort := startPort + attempts - 1
	if startPort <= 0 {
		startPort = 7432
		endPort = startPort + attempts - 1
	}
	return fmt.Sprintf(`gui|serve: no free port in auto range %s:%d–%d (%d ports tried)
hint: stop other trace gui/serve processes, or pin a free port with --addr:
  trace gui --addr 127.0.0.1:7450
`, host, startPort, endPort, attempts)
}
