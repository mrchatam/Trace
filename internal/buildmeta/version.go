// Package buildmeta exposes Trace process identity for CLI, MCP, and HTTP.
// Version may be injected at link time:
//
//	go build -ldflags "-X github.com/mrchatam/Trace/internal/buildmeta.Version=v1.2.3"
//
// When Version is the default "0.0.0-dev" and the binary has VCS build info,
// String() appends a short commit so agents can detect stale servers (#87).
package buildmeta

import (
	"runtime/debug"
	"strings"
)

// Version is the marketing/semver identity. Overridable via -ldflags -X.
var Version = "0.0.0-dev"

// String returns the version agents should compare for stale-server detection.
// Prefer an explicit ldflags Version; otherwise append short VCS revision when
// available so tip builds are not stuck at bare "0.0.0-dev".
func String() string {
	v := strings.TrimSpace(Version)
	if v == "" {
		v = "0.0.0-dev"
	}
	if v != "0.0.0-dev" {
		return v
	}
	if sha := ShortCommit(); sha != "" {
		return v + "+" + sha
	}
	return v
}

// ShortCommit returns the short VCS revision from debug.ReadBuildInfo, or "".
func ShortCommit() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		return ""
	}
	var rev string
	var modified bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = strings.TrimSpace(s.Value)
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	if rev == "" {
		return ""
	}
	if len(rev) > 7 {
		rev = rev[:7]
	}
	if modified {
		return rev + "-dirty"
	}
	return rev
}

// Payload is the shared JSON identity for trace_version / GET /v1/version / CLI.
func Payload(extra map[string]any) map[string]any {
	m := map[string]any{
		"ok":      true,
		"name":    "trace",
		"version": String(),
	}
	if sha := ShortCommit(); sha != "" {
		m["commit"] = sha
	}
	for k, v := range extra {
		m[k] = v
	}
	return m
}
