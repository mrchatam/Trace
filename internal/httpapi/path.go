package httpapi

import (
	"fmt"
	"path/filepath"
	"strings"
)

// confineUnderRoot resolves path under projectRoot and rejects traversal/escape.
// Absolute paths outside the project root are rejected (stricter than CLI).
func confineUnderRoot(projectRoot, path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	absRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", err
	}
	var candidate string
	if filepath.IsAbs(path) {
		candidate, err = filepath.Abs(path)
	} else {
		candidate, err = filepath.Abs(filepath.Join(absRoot, path))
	}
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(absRoot, candidate)
	if err != nil {
		return "", fmt.Errorf("path escapes project root")
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes project root")
	}
	return candidate, nil
}
