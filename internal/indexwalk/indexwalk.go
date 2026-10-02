// Package indexwalk defines what `trace index` considers indexable: the T0
// always-skip rules, language detection, and best-effort gitignore filtering.
// One definition is shared by the CLI walk, `index status` drift, and the
// loop gate, so "the graph is smaller than the tree" means the same thing
// everywhere it is evaluated.
package indexwalk

import (
	"io/fs"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/store"
)

// t0SkipDirs are always-skip directory basenames (case-sensitive). Not build/bin.
var t0SkipDirs = map[string]struct{}{
	".git": {}, ".trace": {}, "node_modules": {}, "vendor": {},
	"__pycache__": {}, ".venv": {}, "venv": {}, "dist": {},
	".next": {}, "target": {}, "coverage": {},
}

// IsT0SkipDir reports whether a directory basename is always skipped.
func IsT0SkipDir(name string) bool {
	_, ok := t0SkipDirs[name]
	return ok
}

// IsT0SkipPath is true when any path component is a T0 dir, or the basename
// ends with a T0 minified suffix (.min.js / .min.mjs / .min.cjs).
func IsT0SkipPath(rel string) bool {
	rel = store.NormalizePath(rel)
	base := filepath.Base(rel)
	if strings.HasSuffix(base, ".min.js") ||
		strings.HasSuffix(base, ".min.mjs") ||
		strings.HasSuffix(base, ".min.cjs") {
		return true
	}
	for _, seg := range strings.Split(rel, "/") {
		if seg != "" && IsT0SkipDir(seg) {
			return true
		}
	}
	return false
}

// Walk returns the root-relative paths a whole-tree `trace index` would
// index: T0-skipped dirs/suffixes and unsupported languages are excluded,
// and gitignore filtering applies when useGitIgnore is set.
func Walk(root string, useGitIgnore bool) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			// 1. T0 always-skip dirs before descent
			if IsT0SkipDir(name) {
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
		if IsT0SkipPath(rel) {
			return nil
		}
		// 4. best-effort gitignore after T0
		if useGitIgnore && GitIgnored(root, rel) {
			return nil
		}
		out = append(out, rel)
		return nil
	})
	return out, err
}

// GitIgnored is best-effort: uses `git check-ignore` when available. Without
// a git repo (or git binary) it reports not-ignored, so gitignore filtering
// safely self-disables outside git work trees.
func GitIgnored(root, rel string) bool {
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

// MaxPathsSample caps the path samples in a Drift report; counts always carry
// the full number and PathsTruncated flags a capped sample.
const MaxPathsSample = 20

// Drift compares the on-disk walk (exactly what `trace index` would index)
// against the store. NotIndexed is the silent-truncation signal: a whole-tree
// index that aborted partway (or a store frozen by a constraint failure)
// leaves the graph a smaller world than the tree. NotOnDisk is the other
// direction: indexed paths that a full-tree index would garbage-collect
// (deleted, renamed, or no longer indexable).
type Drift struct {
	Walkable        int
	Indexed         int
	NotIndexed      int
	NotIndexedPaths []string
	NotOnDisk       int
	NotOnDiskPaths  []string
	PathsTruncated  bool
}

// ComputeDrift walks root (gitignore filtering on — it self-disables outside
// git work trees) and diffs against the store's indexed paths.
func ComputeDrift(root string, st *store.Store) (*Drift, error) {
	walkable, err := Walk(root, true)
	if err != nil {
		return nil, err
	}
	indexed, err := st.ListFilePaths()
	if err != nil {
		return nil, err
	}
	sort.Strings(walkable)
	sort.Strings(indexed)
	inStore := make(map[string]bool, len(indexed))
	for _, p := range indexed {
		inStore[p] = true
	}
	onDisk := make(map[string]bool, len(walkable))
	for _, p := range walkable {
		onDisk[p] = true
	}

	d := &Drift{Walkable: len(walkable), Indexed: len(indexed)}
	for _, p := range walkable {
		if !inStore[p] {
			d.NotIndexed++
			if len(d.NotIndexedPaths) < MaxPathsSample {
				d.NotIndexedPaths = append(d.NotIndexedPaths, p)
			} else {
				d.PathsTruncated = true
			}
		}
	}
	for _, p := range indexed {
		if !onDisk[p] {
			d.NotOnDisk++
			if len(d.NotOnDiskPaths) < MaxPathsSample {
				d.NotOnDiskPaths = append(d.NotOnDiskPaths, p)
			} else {
				d.PathsTruncated = true
			}
		}
	}
	return d, nil
}
