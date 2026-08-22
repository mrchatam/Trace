package testrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RunnerConfig is optional trace/test-runner.json override.
type RunnerConfig struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Cwd     string   `json:"cwd"`
}

func loadRunnerConfig(root string) (*RunnerConfig, error) {
	path := filepath.Join(root, "trace", "test-runner.json")
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
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

// resolveDefaultRunner returns custom config or go test default. Fail-closed when unknown stack.
func resolveDefaultRunner(root string) (RunSpec, bool, error) {
	if cfg, err := loadRunnerConfig(root); err != nil {
		return RunSpec{}, false, err
	} else if cfg != nil {
		return RunSpec{
			Command: cfg.Command,
			Args:    append([]string(nil), cfg.Args...),
			Cwd:     cfg.Cwd,
		}, true, nil
	}
	if goModPresent(root) {
		return RunSpec{
			Command: "go",
			Args:    []string{"test", "./..."},
			Cwd:     root,
		}, true, nil
	}
	return RunSpec{}, false, errors.New("testrun: no test-runner.json and no go.mod at project root")
}
