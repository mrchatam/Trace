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

func TestLoadEnforceModeMalformedJSONTreatAsWarn(t *testing.T) {
	dir := t.TempDir()
	writeTraceConfig(t, dir, "{not json")
	load := config.LoadEnforceModeDetail(dir)
	if load.Mode != config.EnforceWarn || !load.Invalid {
		t.Fatalf("LoadEnforceModeDetail() = %+v want warn+invalid", load)
	}
	if got := config.LoadEnforceMode(dir); got != config.EnforceWarn {
		t.Fatalf("LoadEnforceMode() = %q want %q", got, config.EnforceWarn)
	}
}

func TestLoadEnforceModeUnknownValueTreatAsWarn(t *testing.T) {
	dir := t.TempDir()
	writeTraceConfig(t, dir, `{"enforce":"yolo"}`)
	load := config.LoadEnforceModeDetail(dir)
	if load.Mode != config.EnforceWarn || !load.Invalid {
		t.Fatalf("LoadEnforceModeDetail() = %+v want warn+invalid", load)
	}
}

func TestDoneGateAction(t *testing.T) {
	cases := []struct {
		explicit bool
		mode     config.EnforceMode
		run      bool
		reject   bool
	}{
		{false, config.EnforceOff, false, false},
		{false, config.EnforceWarn, true, false},
		{false, config.EnforceStrict, true, true},
		{true, config.EnforceOff, true, true},
		{true, config.EnforceWarn, true, true},
	}
	for _, tc := range cases {
		run, reject := config.DoneGateAction(tc.explicit, tc.mode)
		if run != tc.run || reject != tc.reject {
			t.Fatalf("DoneGateAction(%v,%q)=(%v,%v) want (%v,%v)",
				tc.explicit, tc.mode, run, reject, tc.run, tc.reject)
		}
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
