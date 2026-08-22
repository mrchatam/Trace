package install

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// CursorReloadTip is the shared DF-22/50 stderr line (print and --write).
const CursorReloadTip = "install: tip: rebuild trace-mcp, prefer absolute --bin, then reload/restart Cursor MCP (or reload window) so the stdio process is not stale; partial tool list (9/17) means stale stdio — call trace_version after reload"

// cursorMCPServer is the Cursor mcp.json entry for Trace (stdio + workspace -C).
type cursorMCPServer struct {
	Type    string   `json:"type"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type cursorTarget struct{}

func (cursorTarget) ID() string   { return TargetCursor }
func (cursorTarget) Tier() string { return TierStable }

func (c cursorTarget) Detect(opts InstallOpts) DetectResult {
	path := cursorMCPPath(opts)
	cursorDir := filepath.Dir(path)
	if st, err := os.Stat(cursorDir); err == nil && st.IsDir() {
		return DetectResult{Detected: true, Reason: "found Cursor dir " + cursorDir}
	}
	if _, err := os.Stat(path); err == nil {
		return DetectResult{Detected: true, Reason: "found mcp.json at " + path}
	}
	return DetectResult{Detected: false, Reason: "no Cursor dir or mcp.json at " + path}
}

func (c cursorTarget) Install(opts InstallOpts) error {
	bin := opts.Bin
	if bin == "" {
		bin = "trace-mcp"
	}
	entry := cursorMCPServer{
		Type:    "stdio",
		Command: bin,
		Args:    []string{"-C", "${workspaceFolder}"},
	}
	snippet := map[string]any{
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
		if err := enc.Encode(snippet); err != nil {
			return err
		}
		if opts.ProjectRoot != "" {
			enforcementPrintHint(errOut)
		}
		fmt.Fprintln(errOut, CursorReloadTip)
		return nil
	}

	path := cursorMCPPath(opts)
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if err := upsertCursorMCP(abs, entry, errOut); err != nil {
		return err
	}
	if opts.ProjectRoot != "" {
		if err := WriteCursorRulesMDC(opts.ProjectRoot); err != nil {
			return err
		}
		if err := UpsertAgentsMD(opts.ProjectRoot); err != nil {
			return err
		}
		fmt.Fprintf(errOut, "install: wrote %s\n", filepath.Join(opts.ProjectRoot, cursorRulesRelPath))
	}
	fmt.Fprintln(errOut, CursorReloadTip)
	return nil
}

func (c cursorTarget) Uninstall(opts InstallOpts) error {
	path := cursorMCPPath(opts)
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	errOut := opts.ErrOut
	if errOut == nil {
		errOut = io.Discard
	}
	if err := removeCursorTrace(abs, errOut); err != nil {
		return err
	}
	if opts.ProjectRoot != "" {
		if err := RemoveCursorRulesMDC(opts.ProjectRoot); err != nil {
			return err
		}
		return StripAgentsMD(opts.ProjectRoot)
	}
	return nil
}

func cursorMCPPath(opts InstallOpts) string {
	if opts.MCPJSON != "" {
		return opts.MCPJSON
	}
	home := opts.HomeDir
	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			home = ""
		}
	}
	return filepath.Join(home, ".cursor", "mcp.json")
}

func upsertCursorMCP(path string, entry cursorMCPServer, errOut io.Writer) error {
	root, exists, fi, err := readMCPRoot(path)
	if err != nil {
		return err
	}

	servers, _ := root["mcpServers"].(map[string]any)
	if servers == nil {
		servers = map[string]any{}
		root["mcpServers"] = servers
	}
	servers["trace"] = entry

	return writeMCPRoot(path, root, exists, fi, errOut)
}

func removeCursorTrace(path string, errOut io.Writer) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // idempotent
		}
		return err
	}
	if fi.Size() == 0 {
		return nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var root map[string]any
	if err := json.Unmarshal(raw, &root); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	if root == nil {
		return nil
	}
	servers, _ := root["mcpServers"].(map[string]any)
	if servers == nil {
		return nil
	}
	if _, ok := servers["trace"]; !ok {
		return nil // already absent — idempotent
	}
	delete(servers, "trace")

	return writeMCPRoot(path, root, true, fi, errOut)
}

func readMCPRoot(path string) (root map[string]any, exists bool, fi os.FileInfo, err error) {
	fi, err = os.Stat(path)
	exists = err == nil
	if err != nil && !os.IsNotExist(err) {
		return nil, false, nil, err
	}
	if exists && fi.Size() > 0 {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, exists, fi, err
		}
		if err := json.Unmarshal(raw, &root); err != nil {
			return nil, exists, fi, fmt.Errorf("invalid JSON in %s: %w", path, err)
		}
		if root == nil {
			root = map[string]any{}
		}
		return root, exists, fi, nil
	}
	return map[string]any{}, exists, fi, nil
}

func writeMCPRoot(path string, root map[string]any, exists bool, fi os.FileInfo, errOut io.Writer) error {
	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	if exists && fi != nil && fi.Size() > 0 {
		stamp := time.Now().UTC().Format("20060102T150405Z")
		bak := path + ".bak." + stamp
		if err := copyFile(path, bak); err != nil {
			return fmt.Errorf("backup: %w", err)
		}
		if errOut != nil {
			fmt.Fprintf(errOut, "install: backup %s\n", bak)
		}
	}

	tmp := path + ".tmp." + fmt.Sprintf("%d", time.Now().UnixNano())
	if err := os.WriteFile(tmp, out, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func copyFile(src, dst string) error {
	raw, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, raw, 0o644)
}
