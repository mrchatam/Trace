package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/indexwalk"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

func cmdIndex(root string, args []string, command string) int {
	if len(args) > 0 && args[0] == "status" {
		return cmdIndexStatus(root, args[1:])
	}
	if len(args) > 0 && args[0] == "watch" {
		return cmdIndexWatch(root, args[1:])
	}
	force := false
	filtered := make([]string, 0, len(args))
	for _, a := range args {
		switch a {
		case "--force", "-f":
			force = true
		default:
			filtered = append(filtered, a)
		}
	}
	args = filtered
	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "index: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "index: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), command, "index"); code != exitOK {
		return code
	}

	ctx := context.Background()
	var repo vcs.Repository
	if r, rerr := tryOpenGit(abs, st); rerr == nil {
		repo = r
		defer repo.Close()
		if _, err := repo.Refresh(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "index: git refresh: %v\n", err)
			return exitFail
		}
	}

	fullTree := len(args) == 0
	paths := args
	if fullTree {
		paths, err = walkIndexable(abs, repo != nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "index: walk: %v\n", err)
			return exitFail
		}
	}

	var indexed, hashSkipped, skipped, removed int
	type indexFailure struct {
		path string
		err  error
	}
	var failures []indexFailure
	for _, p := range paths {
		rel, absPath, err := normalizeProjectPath(abs, p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "index: %v\n", err)
			return exitFail
		}
		// Explicit argv under T0 dirs/suffixes: count as skipped (same as SkipError), not hard fail.
		if isT0SkipPath(rel) {
			skipped++
			continue
		}
		if !fullTree {
			if _, sterr := os.Stat(absPath); sterr != nil && errors.Is(sterr, fs.ErrNotExist) {
				if derr := st.DeleteFileByPath(rel); derr != nil {
					fmt.Fprintf(os.Stderr, "index: %v\n", derr)
					return exitFail
				}
				removed++
				continue
			}
		}
		wasHashSkip, err := indexOne(ctx, st, repo, abs, rel, absPath, force)
		if err != nil {
			var skip *analyzers.SkipError
			if errors.As(err, &skip) {
				skipped++
				continue
			}
			if fullTree {
				// One bad record must not truncate the rest of the tree: a
				// whole-tree walk that aborted on the first failing file left
				// every later file unindexed with only a one-line error to show
				// for it (20 files went missing this way on a consumer repo).
				// Continue, report every failure by path below, keep the sync
				// watermark where it is so `index status` stays stale, and exit
				// non-zero. Explicit argv keeps hard-fail semantics: the caller
				// asked for exactly that file.
				failures = append(failures, indexFailure{path: rel, err: err})
				continue
			}
			fmt.Fprintf(os.Stderr, "index: %v\n", err)
			return exitFail
		}
		if wasHashSkip {
			hashSkipped++
		} else {
			indexed++
		}
		// DF-40: after successful partial argv index, drop same-hash orphans missing on disk.
		if !fullTree {
			n, gerr := gcContentHashOrphans(st, abs, rel)
			if gerr != nil {
				fmt.Fprintf(os.Stderr, "index: hash orphan gc: %v\n", gerr)
				return exitFail
			}
			removed += n
		}
	}

	if fullTree {
		live := make(map[string]struct{}, len(paths))
		for _, p := range paths {
			live[store.NormalizePath(p)] = struct{}{}
		}
		existing, err := st.ListFilePaths()
		if err != nil {
			fmt.Fprintf(os.Stderr, "index: list files: %v\n", err)
			return exitFail
		}
		for _, dbPath := range existing {
			if _, ok := live[dbPath]; ok {
				continue
			}
			if derr := st.DeleteFileByPath(dbPath); derr != nil {
				fmt.Fprintf(os.Stderr, "index: gc %s: %v\n", dbPath, derr)
				return exitFail
			}
			removed++
		}
	}

	fmt.Fprintf(os.Stderr, "indexed %d, hash_skipped %d, skipped %d, removed %d\n", indexed, hashSkipped, skipped, removed)
	if len(failures) > 0 {
		for _, f := range failures {
			fmt.Fprintf(os.Stderr, "index: FAILED %s: %v\n", f.path, f.err)
		}
		fmt.Fprintf(os.Stderr, "index: %d of %d files failed; store is incomplete until these index cleanly (re-run `trace index` after fixing)\n", len(failures), len(paths))
		return exitFail
	}
	if repo != nil {
		if err := updateGraphSyncWatermark(ctx, st, repo); err != nil {
			fmt.Fprintf(os.Stderr, "index: graph sync watermark: %v\n", err)
			return exitFail
		}
		promoteVCSCommitsAfterIndex(ctx, domain.New(st))
	}
	return exitOK
}

// gcContentHashOrphans deletes other DB paths that share indexedRel's content_hash
// and are missing on disk. On-disk siblings (duplicate content) are left alone.
func gcContentHashOrphans(st *store.Store, root, indexedRel string) (int, error) {
	f, err := st.GetFileByPath(indexedRel)
	if err != nil {
		return 0, err
	}
	candidates, err := st.ListFilePathsByContentHash(f.ContentHash)
	if err != nil {
		return 0, err
	}
	var removed int
	for _, q := range candidates {
		if q == indexedRel {
			continue
		}
		absQ := filepath.Join(root, filepath.FromSlash(q))
		if _, sterr := os.Stat(absQ); sterr == nil || !errors.Is(sterr, fs.ErrNotExist) {
			continue
		}
		if derr := st.DeleteFileByPath(q); derr != nil {
			return removed, derr
		}
		removed++
	}
	return removed, nil
}

func indexOne(ctx context.Context, st *store.Store, repo vcs.Repository, root, rel, absPath string, force bool) (hashSkipped bool, err error) {
	var skipped bool
	opts := analyzers.IndexOptions{Force: force, HashSkipped: &skipped}
	if repo != nil {
		head, herr := repo.Head(ctx)
		if herr == nil {
			err = analyzers.IndexFileAtRev(ctx, st, repo, head, rel, opts)
			if err == nil {
				return skipped, nil
			}
			// Untracked / missing at HEAD → fall through to working-tree bytes.
			// Other errors (SkipError, store failures, git failures) must not be masked.
			if !errors.Is(err, vcs.ErrNotFound) {
				return false, err
			}
		}
	}
	content, err := os.ReadFile(absPath)
	if err != nil {
		return false, err
	}
	err = analyzers.IndexFile(ctx, st, rel, content, opts)
	return skipped, err
}

func normalizeProjectPath(root, p string) (rel, absPath string, err error) {
	if filepath.IsAbs(p) {
		absPath = filepath.Clean(p)
	} else {
		absPath = filepath.Join(root, p)
	}
	absPath, err = filepath.Abs(absPath)
	if err != nil {
		return "", "", err
	}
	rel, err = filepath.Rel(root, absPath)
	if err != nil {
		return "", "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("path %q is outside project root", p)
	}
	rel = store.NormalizePath(rel)
	return rel, absPath, nil
}

// Walk/T0/gitignore rules live in internal/indexwalk so `index status` drift
// and the loop gate evaluate the identical definition of "what trace index
// would index" (D1: one smaller-world definition hid a 20-file truncation).
func isT0SkipDir(name string) bool { return indexwalk.IsT0SkipDir(name) }

func isT0SkipPath(rel string) bool { return indexwalk.IsT0SkipPath(rel) }

func walkIndexable(root string, useGitIgnore bool) ([]string, error) {
	return indexwalk.Walk(root, useGitIgnore)
}
