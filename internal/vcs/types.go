package vcs

import (
	"errors"
	"fmt"
)

// Sentinel errors for adapter callers.
var (
	ErrNotRepo     = errors.New("vcs: not a git repository")
	ErrNotFound    = errors.New("vcs: not found")
	ErrGitMissing  = errors.New("vcs: git binary not available")
	ErrInvalidPath = errors.New("vcs: invalid path")
)

// CommitMeta is thin commit metadata (no patch/diff body).
type CommitMeta struct {
	OID         string
	ParentOIDs  []string
	CommittedAt string // RFC3339 or git %cI
	Subject     string // short subject summary
}

// PathChange is a path touched by a commit (optional status A/M/D/R…).
type PathChange struct {
	Path   string
	Status string // A, M, D, R, C, T, … (empty if unknown)
}

// RefreshResult reports how many commits were newly indexed.
type RefreshResult struct {
	NewCommits int
}

// Error wraps a sentinel with context while preserving errors.Is.
type Error struct {
	Op  string
	Err error
}

func (e *Error) Error() string {
	if e.Op == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("vcs %s: %v", e.Op, e.Err)
}

func (e *Error) Unwrap() error { return e.Err }
