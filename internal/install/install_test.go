package install_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/install"
)

func TestInstallDetectListsCursorStable(t *testing.T) {
	home := t.TempDir()
	cursorDir := filepath.Join(home, ".cursor")
	if err := os.MkdirAll(cursorDir, 0o755); err != nil {
		t.Fatal(err)
	}
	infos := install.ListTargets(install.InstallOpts{HomeDir: home})
	var cursor *install.DetectInfo
	for i := range infos {
		if infos[i].ID == install.TargetCursor {
			cursor = &infos[i]
			break
		}
	}
	if cursor == nil {
		t.Fatalf("detect missing cursor: %+v", infos)
	}
	if cursor.Tier != install.TierStable {
		t.Fatalf("tier: want STABLE got %q", cursor.Tier)
	}
	if !cursor.Detected {
		t.Fatalf("want detected with .cursor dir, got reason %q", cursor.Reason)
	}
	if strings.TrimSpace(cursor.Reason) == "" {
		t.Fatal("reason must be non-empty")
	}
}

func TestInstallCursorUninstallIdempotent(t *testing.T) {
	dir := t.TempDir()
	mcpPath := filepath.Join(dir, "mcp.json")
	prior := []byte(`{
  "mcpServers": {
    "other": {
      "type": "stdio",
      "command": "other-mcp",
      "args": []
    },
    "trace": {
      "type": "stdio",
      "command": "old-trace-mcp",
      "args": ["-C", "/old"]
    }
  }
}
`)
	if err := os.WriteFile(mcpPath, prior, 0o644); err != nil {
		t.Fatal(err)
	}

	tgt, err := install.Lookup(install.TargetCursor)
	if err != nil {
		t.Fatal(err)
	}
	opts := install.InstallOpts{
		Write:   true,
		MCPJSON: mcpPath,
		Bin:     "/tmp/trace-mcp",
		ErrOut:  os.Stderr,
	}
	if err := tgt.Install(opts); err != nil {
		t.Fatalf("Install: %v", err)
	}

	if err := tgt.Uninstall(install.InstallOpts{MCPJSON: mcpPath, ErrOut: os.Stderr}); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	raw, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	servers := doc["mcpServers"].(map[string]any)
	if _, ok := servers["trace"]; ok {
		t.Fatal("trace must be removed")
	}
	if _, ok := servers["other"]; !ok {
		t.Fatal("sibling other must be preserved")
	}

	// Second uninstall is idempotent.
	if err := tgt.Uninstall(install.InstallOpts{MCPJSON: mcpPath}); err != nil {
		t.Fatalf("second Uninstall: %v", err)
	}
	// Missing file path is also OK.
	missing := filepath.Join(dir, "missing", "mcp.json")
	if err := tgt.Uninstall(install.InstallOpts{MCPJSON: missing}); err != nil {
		t.Fatalf("uninstall missing file: %v", err)
	}
}

func TestInstallConditionalRefusesWithoutMarker(t *testing.T) {
	root := t.TempDir()
	tgt, err := install.Lookup(install.TargetClaude)
	if err != nil {
		t.Fatal(err)
	}
	if tgt.Tier() != install.TierConditional {
		t.Fatalf("tier: want CONDITIONAL got %q", tgt.Tier())
	}
	err = tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		Bin:         "trace-mcp",
	})
	if err == nil {
		t.Fatal("want refuse without marker")
	}
	if !strings.Contains(err.Error(), "CONDITIONAL") {
		t.Fatalf("error should mention CONDITIONAL: %v", err)
	}
	path := filepath.Join(root, ".claude", "trace-mcp.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("must not write without marker: %v", err)
	}
}

func TestInstallConditionalWritesWithMarker(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	tgt, err := install.Lookup(install.TargetClaude)
	if err != nil {
		t.Fatal(err)
	}
	if err := tgt.Install(install.InstallOpts{
		Write:       true,
		ProjectRoot: root,
		Bin:         "/opt/trace-mcp",
		ErrOut:      os.Stderr,
	}); err != nil {
		t.Fatalf("Install: %v", err)
	}
	path := filepath.Join(root, ".claude", "trace-mcp.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	trace := doc["mcpServers"].(map[string]any)["trace"].(map[string]any)
	if cmd, _ := trace["command"].(string); cmd != "/opt/trace-mcp" {
		t.Fatalf("command: %#v", trace["command"])
	}
}
