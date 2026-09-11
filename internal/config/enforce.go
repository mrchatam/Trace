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

// EnforceLoad describes LoadEnforceModeDetail outcome.
type EnforceLoad struct {
	Mode EnforceMode
	// Invalid is true when .trace/config.json exists but is malformed or has an unknown enforce value.
	Invalid bool
	// Warning is a human-readable honesty note when Invalid (or missing config with .trace/).
	Warning string
}

// LoadEnforceMode reads <root>/.trace/config.json.
// Missing file → EnforceOff. Malformed JSON or unknown enforce value with config present → EnforceWarn
// (not silent Off) so agents cannot skip deliberation via a broken config.
func LoadEnforceMode(projectRoot string) EnforceMode {
	return LoadEnforceModeDetail(projectRoot).Mode
}

// LoadEnforceModeDetail returns mode plus invalid-config honesty metadata.
func LoadEnforceModeDetail(projectRoot string) EnforceLoad {
	path := filepath.Join(projectRoot, ".trace", "config.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return EnforceLoad{Mode: EnforceOff}
	}
	var cfg traceConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return EnforceLoad{
			Mode:    EnforceWarn,
			Invalid: true,
			Warning: "trace: .trace/config.json is malformed — treating enforce as warn (fix JSON or set {\"enforce\":\"off\"|\"warn\"|\"strict\"})",
		}
	}
	switch EnforceMode(cfg.Enforce) {
	case EnforceOff, EnforceWarn, EnforceStrict:
		return EnforceLoad{Mode: EnforceMode(cfg.Enforce)}
	default:
		return EnforceLoad{
			Mode:    EnforceWarn,
			Invalid: true,
			Warning: fmt.Sprintf("trace: .trace/config.json unknown enforce %q — treating as warn (use off|warn|strict)", cfg.Enforce),
		}
	}
}

// WarnInvalidEnforce writes load.Warning when Invalid.
func WarnInvalidEnforce(load EnforceLoad, w io.Writer) {
	if w == nil || !load.Invalid || load.Warning == "" {
		return
	}
	fmt.Fprintln(w, load.Warning)
}

// DoneGateAction decides whether to run GateForDone and whether failure rejects the transition.
// explicitEnforce mirrors CLI --enforce / MCP enforce=true (always reject on fail).
func DoneGateAction(explicitEnforce bool, mode EnforceMode) (run bool, rejectOnFail bool) {
	if explicitEnforce {
		return true, true
	}
	switch mode {
	case EnforceStrict:
		return true, true
	case EnforceWarn:
		return true, false
	default:
		return false, false
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
