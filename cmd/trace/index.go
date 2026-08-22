package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/domain"
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

	var indexed, skipped, removed int
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
		if err := indexOne(ctx, st, repo, abs, rel, absPath); err != nil {
			var skip *analyzers.SkipError
			if errors.As(err, &skip) {
				skipped++
				continue
			}
			fmt.Fprintf(os.Stderr, "index: %v\n", err)
			return exitFail
		}
		indexed++
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

	fmt.Fprintf(os.Stderr, "indexed %d, skipped %d, removed %d\n", indexed, skipped, removed)
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

func indexOne(ctx context.Context, st *store.Store, repo vcs.Repository, root, rel, absPath string) error {
	if repo != nil {
		head, err := repo.Head(ctx)
		if err == nil {
			err = analyzers.IndexFileAtRev(ctx, st, repo, head, rel, analyzers.IndexOptions{})
			if err == nil {
				return nil
			}
			// Untracked / missing at HEAD → fall through to working-tree bytes.
			// Other errors (SkipError, store failures, git failures) must not be masked.
			if !errors.Is(err, vcs.ErrNotFound) {
				return err
			}
		}
	}
	content, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}
	return analyzers.IndexFile(ctx, st, rel, content, analyzers.IndexOptions{})
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

// t0SkipDirs are always-skip directory basenames (case-sensitive). Not build/bin.
var t0SkipDirs = map[string]struct{}{
	".git": {}, ".trace": {}, "node_modules": {}, "vendor": {},
	"__pycache__": {}, ".venv": {}, "venv": {}, "dist": {},
	".next": {}, "target": {}, "coverage": {},
}

func isT0SkipDir(name string) bool {
	_, ok := t0SkipDirs[name]
	return ok
}

// isT0SkipPath is true when any path component is a T0 dir, or the basename
// ends with a T0 minified suffix (.min.js / .min.mjs / .min.cjs).
func isT0SkipPath(rel string) bool {
	rel = store.NormalizePath(rel)
	base := filepath.Base(rel)
	if strings.HasSuffix(base, ".min.js") ||
		strings.HasSuffix(base, ".min.mjs") ||
		strings.HasSuffix(base, ".min.cjs") {
		return true
	}
	for _, seg := range strings.Split(rel, "/") {
		if seg != "" && isT0SkipDir(seg) {
			return true
		}
	}
	return false
}

func walkIndexable(root string, useGitIgnore bool) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			// 1. T0 always-skip dirs before descent
			if isT0SkipDir(name) {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = store.NormalizePath(rel)
		// 2. unsupported language → skip
		if _, ok := analyzers.DetectLanguage(rel); !ok {
			return nil
		}
		// 3. T0 file suffix or path-segment
		if isT0SkipPath(rel) {
			return nil
		}
		// 4. best-effort gitignore after T0
		if useGitIgnore && gitIgnored(root, rel) {
			return nil
		}
		out = append(out, rel)
		return nil
	})
	return out, err
}

// gitIgnored is best-effort: uses `git check-ignore` when available.
func gitIgnored(root, rel string) bool {
	cmd := exec.Command("git", "-C", root, "check-ignore", "-q", "--", rel)
	err := cmd.Run()
	if err == nil {
		return true // exit 0 → ignored
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 1 {
		return false // not ignored
	}
	return false
}
