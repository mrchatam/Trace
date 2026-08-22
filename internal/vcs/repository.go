package vcs

import "context"

// Repository is the VCS adapter surface used by the rest of Trace.
//
// ShowFile always returns bytes from Git (canonical content). History-style
// methods may be accelerated by a thin SQLite index after Refresh; they may
// also fall back to Git. Refresh indexes only commits since the durable
// watermark and is a noop when HEAD has not advanced.
type Repository interface {
	// Head returns the OID of HEAD.
	Head(ctx context.Context) (string, error)

	// IsRepo reports whether the bound root is inside a Git work tree.
	IsRepo(ctx context.Context) (bool, error)

	// ShowFile returns file bytes at rev:path (equivalent to git show).
	ShowFile(ctx context.Context, rev, path string) ([]byte, error)

	// History returns commits that touched path, newest first, up to limit.
	// limit <= 0 means an implementation-defined default.
	History(ctx context.Context, path string, limit int) ([]CommitMeta, error)

	// CommitsBetween returns commits reachable from toRev but not fromRev,
	// oldest first (toRev exclusive of fromRev ancestry).
	CommitsBetween(ctx context.Context, fromRev, toRev string) ([]CommitMeta, error)

	// Changes returns paths changed in commitOID (with status when available).
	Changes(ctx context.Context, commitOID string) ([]PathChange, error)

	// DiffNameStatus returns path changes between fromRev and toRev using
	// git diff --name-status fromRev..toRev (no patch bodies).
	DiffNameStatus(ctx context.Context, fromRev, toRev string) ([]PathChange, error)

	// LastChanged returns the most recent commit that touched path.
	LastChanged(ctx context.Context, path string) (CommitMeta, error)

	// Refresh indexes new commits since the durable watermark.
	Refresh(ctx context.Context) (RefreshResult, error)

	// Close releases resources (e.g. SQLite handle). Idempotent.
	Close() error
}
