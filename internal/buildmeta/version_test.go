package buildmeta

import (
	"strings"
	"testing"
)

func TestStringDefaultOrVCS(t *testing.T) {
	got := String()
	if got == "" {
		t.Fatal("empty version")
	}
	// Either bare default (no VCS in this build) or default+shortsha.
	if got != "0.0.0-dev" && !strings.HasPrefix(got, "0.0.0-dev+") {
		t.Fatalf("unexpected version %q", got)
	}
	if sha := ShortCommit(); sha != "" {
		if !strings.Contains(got, sha) && !strings.HasPrefix(sha, strings.TrimPrefix(got, "0.0.0-dev+")) {
			// String embeds ShortCommit; accept dirty suffix mismatch only if equal base
			base := strings.TrimSuffix(sha, "-dirty")
			if !strings.Contains(got, base) {
				t.Fatalf("version %q should include commit %q", got, sha)
			}
		}
		if got == "0.0.0-dev" {
			t.Fatalf("VCS commit %q present but String still bare 0.0.0-dev", sha)
		}
	}
}

func TestStringHonorsLdflagsOverride(t *testing.T) {
	prev := Version
	t.Cleanup(func() { Version = prev })
	Version = "v9.9.9"
	if got := String(); got != "v9.9.9" {
		t.Fatalf("got %q want v9.9.9", got)
	}
}

func TestPayloadIncludesCommitWhenPresent(t *testing.T) {
	m := Payload(map[string]any{"api_version": "1.0.0"})
	if m["ok"] != true || m["name"] != "trace" {
		t.Fatalf("%v", m)
	}
	if m["api_version"] != "1.0.0" {
		t.Fatalf("extra lost: %v", m)
	}
	if _, ok := m["version"].(string); !ok {
		t.Fatalf("version missing: %v", m)
	}
	if sha := ShortCommit(); sha != "" {
		if m["commit"] != sha {
			t.Fatalf("commit=%v want %q", m["commit"], sha)
		}
	}
}
