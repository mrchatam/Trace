package install

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Claude marker: project-relative .claude/ directory or CLAUDE.md file.
// On Install --write with marker: writes .claude/trace-mcp.json (stdio snippet).

const claudeConfigName = "trace-mcp.json"

type claudeTarget struct{}

func (claudeTarget) ID() string   { return TargetClaude }
func (claudeTarget) Tier() string { return TierConditional }

func (c claudeTarget) Detect(opts InstallOpts) DetectResult {
	root := projectRoot(opts)
	if hasClaudeMarker(root) {
		return DetectResult{Detected: true, Reason: "found Claude marker (.claude/ or CLAUDE.md) under " + root}
	}
	return DetectResult{Detected: false, Reason: "no Claude marker (.claude/ or CLAUDE.md) under " + root}
}

func (c claudeTarget) Install(opts InstallOpts) error {
	root := projectRoot(opts)
	if !hasClaudeMarker(root) {
		return fmt.Errorf("install: claude is CONDITIONAL — require .claude/ or CLAUDE.md under %s", root)
	}

	bin := opts.Bin
	if bin == "" {
		bin = "trace-mcp"
	}
	entry := map[string]any{
		"type":    "stdio",
		"command": bin,
		"args":    []string{"-C", "${workspaceFolder}"},
	}
	doc := map[string]any{
		"mcpServers": map[string]any{
			"trace": entry,
		},
	}

	out := opts.Out
	if out == nil {
		out = io.Discard
	}
	errOut := opts.ErrOut
	if errOut == nil {
		errOut = io.Discard
	}

	if !opts.Write {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(doc); err != nil {
			return err
		}
		enforcementPrintHint(errOut)
		return nil
	}

	claudeDir := filepath.Join(root, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(claudeDir, claudeConfigName)
	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(errOut, "install: wrote %s\n", path)
	if err := UpsertClaudeRules(root); err != nil {
		return err
	}
	return UpsertAgentsMD(root)
}

func (c claudeTarget) Uninstall(opts InstallOpts) error {
	root := projectRoot(opts)
	path := filepath.Join(root, ".claude", claudeConfigName)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := StripClaudeRules(root); err != nil {
		return err
	}
	return StripAgentsMD(root)
}

func projectRoot(opts InstallOpts) string {
	if opts.ProjectRoot != "" {
		return opts.ProjectRoot
	}
	return "."
}

func hasClaudeMarker(root string) bool {
	claudeDir := filepath.Join(root, ".claude")
	if st, err := os.Stat(claudeDir); err == nil && st.IsDir() {
		return true
	}
	if _, err := os.Stat(filepath.Join(root, "CLAUDE.md")); err == nil {
		return true
	}
	return false
}
