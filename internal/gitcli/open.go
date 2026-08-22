package gitcli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

// Repo is the git-CLI-backed vcs.Repository.
type Repo struct {
	root      string
	run       *runner
	st        *store.Store
	ownsStore bool
}

// Open binds repoRoot (absolute) as a Git work tree and opens the project store
// for the thin VCS index. Fails if git is missing or root is not a repository.
func Open(repoRoot string) (vcs.Repository, error) {
	abs, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("gitcli: resolve root: %w", err)
	}

	r := &Repo{
		root: abs,
		run:  newRunner(abs),
	}

	ctx := context.Background()
	ok, err := r.IsRepo(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &vcs.Error{Op: "Open", Err: vcs.ErrNotRepo}
	}

	st, err := store.Open(abs)
	if err != nil {
		return nil, fmt.Errorf("gitcli: open store: %w", err)
	}
	r.st = st
	r.ownsStore = true
	return r, nil
}

// OpenWithStore binds repoRoot as a Git work tree using an already-open store
// for the same abs root. Close does not close the store (caller owns it).
// Use this when the CLI/MCP already holds store.Open to avoid ErrLocked.
func OpenWithStore(repoRoot string, st *store.Store) (vcs.Repository, error) {
	if st == nil {
		return nil, fmt.Errorf("gitcli: OpenWithStore: nil store")
	}
	abs, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("gitcli: resolve root: %w", err)
	}

	r := &Repo{
		root:      abs,
		run:       newRunner(abs),
		st:        st,
		ownsStore: false,
	}

	ctx := context.Background()
	ok, err := r.IsRepo(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &vcs.Error{Op: "Open", Err: vcs.ErrNotRepo}
	}
	return r, nil
}

// Close releases the store handle when this Repo owns it (Open, not OpenWithStore).
func (r *Repo) Close() error {
	if r == nil || r.st == nil {
		return nil
	}
	if !r.ownsStore {
		r.st = nil
		return nil
	}
	err := r.st.Close()
	r.st = nil
	return err
}

// IsRepo implements vcs.Repository.
func (r *Repo) IsRepo(ctx context.Context) (bool, error) {
	out, err := r.run.run(ctx, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		if strings.Contains(err.Error(), "not a git repository") {
			return false, nil
		}
		// Missing git binary
		if _, ok := err.(*vcs.Error); ok {
			return false, err
		}
		// Other failures: treat as not a repo when stderr says so; else surface.
		low := strings.ToLower(err.Error())
		if strings.Contains(low, "not a git repository") || strings.Contains(low, "fatal:") {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(out) == "true", nil
}

// Head implements vcs.Repository.
func (r *Repo) Head(ctx context.Context) (string, error) {
	out, err := r.run.run(ctx, "rev-parse", "HEAD")
	if err != nil {
		return "", &vcs.Error{Op: "Head", Err: err}
	}
	oid := strings.TrimSpace(out)
	if oid == "" {
		return "", &vcs.Error{Op: "Head", Err: vcs.ErrNotFound}
	}
	return oid, nil
}

var _ vcs.Repository = (*Repo)(nil)
