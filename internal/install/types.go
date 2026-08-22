package install

import "io"

// Tier vocabulary (exact Trace strings).
const (
	TierStable      = "STABLE"
	TierConditional = "CONDITIONAL"
	TierOptIn       = "OPT_IN"
)

// Target IDs.
const (
	TargetCursor     = "cursor"
	TargetClaude     = "claude"
	TargetGitHook    = "git-hook"
	TargetCursorHook = "cursor-hook"
)

// DetectResult is the outcome of Target.Detect.
type DetectResult struct {
	Detected bool
	Reason   string
}

// DetectInfo is the JSON shape for `trace install detect`.
type DetectInfo struct {
	ID       string `json:"id"`
	Tier     string `json:"tier"`
	Detected bool   `json:"detected"`
	Reason   string `json:"reason"`
}

// InstallOpts configures Install / Uninstall / Detect for a target.
type InstallOpts struct {
	// Write upserts config; when false, Install may print a snippet to Out.
	Write bool
	// Bin is the trace-mcp binary path/name (Cursor).
	Bin string
	// MCPJSON overrides the Cursor mcp.json path (tests / --mcp-json).
	MCPJSON string
	// ProjectRoot is the project directory for marker checks (claude). Empty → ".".
	ProjectRoot string
	// HomeDir overrides os.UserHomeDir (tests). Empty → real home.
	HomeDir string
	// CatalogPath overrides bundled trace/agents/default.json (tests).
	CatalogPath string
	// Out receives print-only JSON (default os.Stdout when nil at CLI).
	Out io.Writer
	// ErrOut receives tips / backup lines (default os.Stderr when nil at CLI).
	ErrOut io.Writer
}

// Target is one installable client/config surface.
type Target interface {
	ID() string
	Tier() string
	Detect(opts InstallOpts) DetectResult
	Install(opts InstallOpts) error
	Uninstall(opts InstallOpts) error
}
