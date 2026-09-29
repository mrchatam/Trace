package testrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// RunnerConfig is the optional test-runner override. Locations, most specific
// first: trace/test-runner.json (canonical), then test-runner.json at project
// root. Schema: {"command":string, "args":[string...], "cwd":string} — cwd is
// relative to project root when not absolute. The literal {path} placeholder in
// any arg is replaced with the selected test file path (root-relative), enabling
// per-target runs (e.g. "args":["exec","vitest","run","{path}"]). Without {path},
// the exact command runs once per selected target. Custom config always wins
// over the go.mod → `go test` fallback.
type RunnerConfig struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Cwd     string   `json:"cwd"`
}

// runnerConfigPaths lists candidate locations for test-runner.json, most specific first:
// trace/test-runner.json (canonical, P22) then project-root test-runner.json (convenience).
var runnerConfigPaths = []string{
	filepath.Join("trace", "test-runner.json"),
	"test-runner.json",
}

func loadRunnerConfig(root string) (*RunnerConfig, error) {
	for _, rel := range runnerConfigPaths {
		cfg, err := loadRunnerConfigAt(root, filepath.Join(root, rel))
		if err != nil {
			return nil, err
		}
		if cfg != nil {
			return cfg, nil
		}
	}
	return nil, nil
}

func loadRunnerConfigAt(root, path string) (*RunnerConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		// NotExist and ENOTDIR (e.g. trace/ exists as a regular file) mean "no
		// config here" — callers fall through to the next candidate location.
		if os.IsNotExist(err) || errors.Is(err, syscall.ENOTDIR) {
			return nil, nil
		}
		return nil, err
	}
	var cfg RunnerConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, fmt.Errorf("testrun: parse test-runner.json: %w", err)
	}
	if strings.TrimSpace(cfg.Command) == "" {
		return nil, errors.New("testrun: test-runner.json command is required")
	}
	if cfg.Cwd == "" {
		cfg.Cwd = root
	} else if !filepath.IsAbs(cfg.Cwd) {
		cfg.Cwd = filepath.Join(root, cfg.Cwd)
	}
	return &cfg, nil
}

func goModPresent(root string) bool {
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "module ") {
			return true
		}
	}
	return false
}

func goModulePath(root string) (string, bool) {
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", false
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), true
		}
	}
	return "", false
}

// runnerSource reports which runner resolution produced a spec.
type runnerSource int

const (
	runnerSourceNone   runnerSource = iota // no config, no go.mod
	runnerSourceCustom                     // test-runner.json
	runnerSourceGo                         // go.mod fallback
)

// resolveDefaultRunner returns custom config or go test default. Fail-closed when unknown stack.
func resolveDefaultRunner(root string) (RunSpec, runnerSource, error) {
	if cfg, err := loadRunnerConfig(root); err != nil {
		return RunSpec{}, runnerSourceNone, err
	} else if cfg != nil {
		return RunSpec{
			Command: cfg.Command,
			Args:    append([]string(nil), cfg.Args...),
			Cwd:     cfg.Cwd,
		}, runnerSourceCustom, nil
	}
	if goModPresent(root) {
		return RunSpec{
			Command: "go",
			Args:    []string{"test", "./..."},
			Cwd:     root,
		}, runnerSourceGo, nil
	}
	return RunSpec{}, runnerSourceNone, errors.New("testrun: no test-runner.json and no go.mod at project root")
}
