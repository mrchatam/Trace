package vcs

import (
	"context"
	"sync"
)

// Fake is an in-memory Repository for unit tests outside gitcli.
type Fake struct {
	mu sync.Mutex

	HeadOID string
	IsGit   bool

	// Files maps "rev:path" → content.
	Files map[string][]byte

	// Commits is the ordered history (oldest first) used by CommitsBetween.
	Commits []CommitMeta

	// PathsByCommit maps commit OID → path changes.
	PathsByCommit map[string][]PathChange

	// DiffByRange maps "from..to" → path changes for DiffNameStatus.
	DiffByRange map[string][]PathChange

	// PathHistory maps path → commits newest-first for History/LastChanged.
	PathHistory map[string][]CommitMeta

	RefreshFn func(ctx context.Context) (RefreshResult, error)

	Closed bool
}

// Head implements Repository.
func (f *Fake) Head(ctx context.Context) (string, error) {
	_ = ctx
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.HeadOID == "" {
		return "", &Error{Op: "Head", Err: ErrNotFound}
	}
	return f.HeadOID, nil
}

// IsRepo implements Repository.
func (f *Fake) IsRepo(ctx context.Context) (bool, error) {
	_ = ctx
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.IsGit, nil
}

// ShowFile implements Repository.
func (f *Fake) ShowFile(ctx context.Context, rev, path string) ([]byte, error) {
	_ = ctx
	f.mu.Lock()
	defer f.mu.Unlock()
	key := rev + ":" + path
	b, ok := f.Files[key]
	if !ok {
		return nil, &Error{Op: "ShowFile", Err: ErrNotFound}
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out, nil
}

// History implements Repository.
func (f *Fake) History(ctx context.Context, path string, limit int) ([]CommitMeta, error) {
	_ = ctx
	f.mu.Lock()
	defer f.mu.Unlock()
	h := append([]CommitMeta(nil), f.PathHistory[path]...)
	if limit > 0 && len(h) > limit {
		h = h[:limit]
	}
	return h, nil
}

// CommitsBetween implements Repository.
func (f *Fake) CommitsBetween(ctx context.Context, fromRev, toRev string) ([]CommitMeta, error) {
	_ = ctx
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []CommitMeta
	inRange := fromRev == ""
	for _, c := range f.Commits {
		if !inRange {
			if c.OID == fromRev {
				inRange = true
			}
			continue
		}
		out = append(out, c)
		if c.OID == toRev {
			break
		}
	}
	return out, nil
}

// Changes implements Repository.
func (f *Fake) Changes(ctx context.Context, commitOID string) ([]PathChange, error) {
	_ = ctx
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]PathChange(nil), f.PathsByCommit[commitOID]...), nil
}

// DiffNameStatus implements Repository.
func (f *Fake) DiffNameStatus(ctx context.Context, fromRev, toRev string) ([]PathChange, error) {
	_ = ctx
	f.mu.Lock()
	defer f.mu.Unlock()
	key := toRev
	if fromRev != "" {
		key = fromRev + ".." + toRev
	}
	if f.DiffByRange == nil {
		return nil, &Error{Op: "DiffNameStatus", Err: ErrNotFound}
	}
	ch, ok := f.DiffByRange[key]
	if !ok {
		return nil, &Error{Op: "DiffNameStatus", Err: ErrNotFound}
	}
	return append([]PathChange(nil), ch...), nil
}

// LastChanged implements Repository.
func (f *Fake) LastChanged(ctx context.Context, path string) (CommitMeta, error) {
	_ = ctx
	f.mu.Lock()
	defer f.mu.Unlock()
	h := f.PathHistory[path]
	if len(h) == 0 {
		return CommitMeta{}, &Error{Op: "LastChanged", Err: ErrNotFound}
	}
	return h[0], nil
}

// Refresh implements Repository.
func (f *Fake) Refresh(ctx context.Context) (RefreshResult, error) {
	f.mu.Lock()
	fn := f.RefreshFn
	f.mu.Unlock()
	if fn != nil {
		return fn(ctx)
	}
	return RefreshResult{NewCommits: 0}, nil
}

// Close implements Repository.
func (f *Fake) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Closed = true
	return nil
}

// Ensure Fake satisfies Repository.
var _ Repository = (*Fake)(nil)
