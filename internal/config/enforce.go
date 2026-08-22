package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type EnforceMode string

const (
	EnforceOff    EnforceMode = "off"
	EnforceWarn   EnforceMode = "warn"
	EnforceStrict EnforceMode = "strict"
)

type traceConfig struct {
	Enforce string `json:"enforce"`
}

// LoadEnforceMode reads <root>/.trace/config.json.
// Returns EnforceOff when the file is missing or any parse/validation error occurs.
func LoadEnforceMode(projectRoot string) EnforceMode {
	path := filepath.Join(projectRoot, ".trace", "config.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return EnforceOff
	}
	var cfg traceConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return EnforceOff
	}
	switch EnforceMode(cfg.Enforce) {
	case EnforceOff, EnforceWarn, EnforceStrict:
		return EnforceMode(cfg.Enforce)
	default:
		return EnforceOff
	}
}

// WarnIfTraceDirWithoutConfig emits a one-time stderr nudge when .trace/ exists but config is missing/invalid.
func WarnIfTraceDirWithoutConfig(projectRoot string, w io.Writer) {
	if w == nil {
		return
	}
	traceDir := filepath.Join(projectRoot, ".trace")
	if st, err := os.Stat(traceDir); err != nil || !st.IsDir() {
		return
	}
	path := filepath.Join(traceDir, "config.json")
	if _, err := os.Stat(path); err == nil {
		return
	}
	fmt.Fprintf(w, "trace: .trace/ exists without config.json — consider `.trace/config.json` with {\"enforce\": \"warn\"} after `trace install`\n")
}
