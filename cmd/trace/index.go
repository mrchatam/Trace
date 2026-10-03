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
	var failures []pathFailure
	if fullTree {
		// Whole-tree walks reconcile with continue-past-failure semantics: one
		// bad record must not truncate the rest of the tree — an abort left
		// every later file unindexed with only a one-line error to show for it
		// (20 files went missing this way on a consumer repo). Failures print
		// as they happen; the sync watermark below stays untouched so `index
		// status` stays stale.
		var rels []string
		for _, p := range paths {
			rel, _, nerr := normalizeProjectPath(abs, p)
			if nerr != nil {
				fmt.Fprintf(os.Stderr, "index: %v\n", nerr)
				return exitFail
			}
			if isT0SkipPath(rel) {
				skipped++
				continue
			}
			rels = append(rels, rel)
		}
		out := reconcilePaths(ctx, st, repo, abs, rels, force, "index: ", nil)
		indexed = out.indexed
		hashSkipped = out.hashSkipped
		skipped += out.skipped
		failures = out.failures
	} else {
		// Explicit argv reconciles with the same continue-past-failure
		// semantics as the whole-tree walk: one bad file must not stop the
		// rest of the caller's list from indexing (an abort turned a partial
		// success into no work at all). Missing paths still delete-only;
		// T0 paths still count as skipped; failures print as they happen and
		// the summary below exits non-zero.
		var rels []string
		for _, p := range paths {
			rel, absPath, err := normalizeProjectPath(abs, p)
			if err != nil {
				failures = append(failures, pathFailure{path: p, err: err})
				fmt.Fprintf(os.Stderr, "index: FAILED %s: %v\n", p, err)
				continue
			}
			// Explicit argv under T0 dirs/suffixes: count as skipped (same as SkipError), not a failure.
			if isT0SkipPath(rel) {
				skipped++
				continue
			}
			if _, sterr := os.Stat(absPath); sterr != nil && errors.Is(sterr, fs.ErrNotExist) {
				if derr := st.DeleteFileByPath(rel); derr != nil {
					failures = append(failures, pathFailure{path: rel, err: derr})
					fmt.Fprintf(os.Stderr, "index: FAILED %s: %v\n", rel, derr)
					continue
				}
				removed++
				continue
			}
			rels = append(rels, rel)
		}
		out := reconcilePaths(ctx, st, repo, abs, rels, force, "index: ", func(rel string) error {
			// DF-40: after a successful partial argv index, drop same-hash orphans missing on disk.
			n, gerr := gcContentHashOrphans(st, abs, rel)
			if gerr != nil {
				return fmt.Errorf("hash orphan gc: %w", gerr)
			}
			removed += n
			return nil
		})
		indexed += out.indexed
		hashSkipped += out.hashSkipped
		skipped += out.skipped
		failures = append(failures, out.failures...)
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
		// Per-file FAILED lines were printed by reconcilePaths as they
		// happened; the summary carries the count and the remedy.
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

// pathFailure is one file that failed to reconcile.
type pathFailure struct {
	path string
	err  error
}

// reconcileOutcome is the tally of indexing a list of paths.
type reconcileOutcome struct {
	indexed     int
	hashSkipped int
	skipped     int
	failures    []pathFailure
}

// reconcilePaths indexes each rel path against root, continuing past
// per-file failures so one bad record cannot truncate the rest of the set
// (the whole-tree lesson: an abort converts loud failure into silent
// truncation). SkipError paths count as skipped and print nothing; every
// other failure prints "<logPrefix>FAILED <rel>: <err>" as it happens and is
// returned in outcome.failures for caller policy (exit codes, summaries).
// onIndexed, when non-nil, fires per successful path and its error, if any,
// is recorded as that path's failure (watch mode logs "indexed <rel>" and
// returns nil; the CLI argv path runs its post-index content-hash orphan GC
// here). Shared by the whole-tree walk, explicit argv, and watch-mode drift
// healing so all surfaces keep identical reconcile semantics.
func reconcilePaths(
	ctx context.Context,
	st *store.Store,
	repo vcs.Repository,
	root string,
	rels []string,
	force bool,
	logPrefix string,
	onIndexed func(rel string) error,
) reconcileOutcome {
	var out reconcileOutcome
	for _, rel := range rels {
		absPath := filepath.Join(root, filepath.FromSlash(rel))
		wasHashSkip, err := indexOne(ctx, st, repo, root, rel, absPath, force)
		if err != nil {
			var skip *analyzers.SkipError
			if errors.As(err, &skip) {
				out.skipped++
				continue
			}
			out.failures = append(out.failures, pathFailure{path: rel, err: err})
			fmt.Fprintf(os.Stderr, "%sFAILED %s: %v\n", logPrefix, rel, err)
			continue
		}
		if wasHashSkip {
			out.hashSkipped++
		} else {
			out.indexed++
		}
		if onIndexed != nil {
			if herr := onIndexed(rel); herr != nil {
				out.failures = append(out.failures, pathFailure{path: rel, err: herr})
				fmt.Fprintf(os.Stderr, "%sFAILED %s: %v\n", logPrefix, rel, herr)
				continue
			}
		}
	}
	return out
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
