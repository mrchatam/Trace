package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/config"
)

func writeTraceConfig(t *testing.T, root, content string) {
	t.Helper()
	dir := filepath.Join(root, ".trace")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadEnforceModeMissingFile(t *testing.T) {
	dir := t.TempDir()
	if got := config.LoadEnforceMode(dir); got != config.EnforceOff {
		t.Fatalf("LoadEnforceMode() = %q want %q", got, config.EnforceOff)
	}
}

func TestLoadEnforceModeValidValues(t *testing.T) {
	dir := t.TempDir()
	cases := []struct {
		content string
		want    config.EnforceMode
	}{
		{`{"enforce":"off"}`, config.EnforceOff},
		{`{"enforce":"warn"}`, config.EnforceWarn},
		{`{"enforce":"strict"}`, config.EnforceStrict},
	}
	for _, tc := range cases {
		writeTraceConfig(t, dir, tc.content)
		if got := config.LoadEnforceMode(dir); got != tc.want {
			t.Fatalf("LoadEnforceMode(%q) = %q want %q", tc.content, got, tc.want)
		}
	}
}

func TestLoadEnforceModeMalformedJSON(t *testing.T) {
	dir := t.TempDir()
	writeTraceConfig(t, dir, "{not json")
	if got := config.LoadEnforceMode(dir); got != config.EnforceOff {
		t.Fatalf("LoadEnforceMode() = %q want %q", got, config.EnforceOff)
	}
}

func TestLoadEnforceModeUnknownValue(t *testing.T) {
	dir := t.TempDir()
	writeTraceConfig(t, dir, `{"enforce":"yolo"}`)
	if got := config.LoadEnforceMode(dir); got != config.EnforceOff {
		t.Fatalf("LoadEnforceMode() = %q want %q", got, config.EnforceOff)
	}
}

func TestWarnIfTraceDirWithoutConfig(t *testing.T) {
	dir := t.TempDir()
	traceDir := filepath.Join(dir, ".trace")
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	var buf strings.Builder
	config.WarnIfTraceDirWithoutConfig(dir, &buf)
	out := buf.String()
	if !strings.Contains(out, "config.json") || !strings.Contains(out, `"enforce": "warn"`) {
		t.Fatalf("stderr nudge missing expected substring: %q", out)
	}
}
