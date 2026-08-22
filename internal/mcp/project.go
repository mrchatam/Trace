package mcp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mrchatam/Trace/internal/gitcli"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

// resolveProject returns an absolute project root. Per-tool project overrides
// the server default; empty both defaults to the process cwd.
// Path-local only: filepath.Abs with no parent .trace walk-up.
func (s *Server) resolveProject(projectOverride string) (string, error) {
	root := projectOverride
	if root == "" {
		root = s.defaultRoot
	}
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("mcp: project: %w", err)
		}
		root = cwd
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("mcp: project: %w", err)
	}
	return abs, nil
}

func (s *Server) openStore(projectOverride string) (abs string, st *store.Store, err error) {
	abs, err = s.resolveProject(projectOverride)
	if err != nil {
		return "", nil, err
	}
	st, err = store.OpenExisting(abs)
	if err != nil {
		// ErrLocked already guides serialize CLI↔MCP / worktrees (DF-47).
		// ErrNotInitialized is fail-closed missing store (DF-76: no auto-init).
		if errors.Is(err, store.ErrLocked) || errors.Is(err, store.ErrUnauthorized) || errors.Is(err, store.ErrNotInitialized) {
			return "", nil, fmt.Errorf("mcp: %w", err)
		}
		return "", nil, fmt.Errorf("mcp: open store: %w", err)
	}
	return abs, st, nil
}

func tryOpenGit(absRoot string, st *store.Store) (vcs.Repository, error) {
	repo, err := gitcli.OpenWithStore(absRoot, st)
	if err != nil {
		if isNotRepo(err) {
			return nil, err
		}
		return nil, err
	}
	return repo, nil
}

func isNotRepo(err error) bool {
	var ve *vcs.Error
	if errors.As(err, &ve) {
		return errors.Is(ve.Err, vcs.ErrNotRepo)
	}
	return errors.Is(err, vcs.ErrNotRepo)
}
