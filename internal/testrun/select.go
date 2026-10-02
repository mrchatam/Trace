package testrun

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

// TestTarget is one executed test invocation recorded as kind=test.
type TestTarget struct {
	Name       string
	Package    string // go test package arg, e.g. ./internal/foo
	RunPattern string // optional -run value
	Path       string // root-relative test file path; enables {path} substitution for custom runners (may be empty for package fallback)
}

type targetKey struct {
	name string
	pkg  string
}

// SelectTestTargets picks relevant tests from seed paths via impact walk +
// validates; Go package fallback; JS/TS sibling test-dir fallback. When
// nothing is selected it fails closed with the reason: which stage produced
// zero targets, and which seed paths are not in the index at all.
func SelectTestTargets(
	ctx context.Context,
	st *store.Store,
	dom *domain.Service,
	taskID string,
	seedPaths []string,
) ([]TestTarget, error) {
	if st == nil {
		return nil, fmt.Errorf("testrun: store is required")
	}
	if dom == nil {
		return nil, fmt.Errorf("testrun: domain service is required")
	}

	paths, err := resolveSeedPaths(ctx, dom, taskID, seedPaths)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, nil
	}

	targets, err := targetsFromImpactAndValidates(ctx, st, paths)
	if err != nil {
		return nil, err
	}
	if len(targets) > 0 {
		return targets, nil
	}
	fallback, ferr := packageFallbackTargets(st.ProjectRoot(), paths)
	if ferr != nil {
		return nil, ferr
	}
	if len(fallback) > 0 {
		return fallback, nil
	}
	if fallback := jstsFallbackTargets(st, paths); len(fallback) > 0 {
		return fallback, nil
	}
	return nil, noTargetsDiagnosis(st, paths)
}

func resolveSeedPaths(ctx context.Context, dom *domain.Service, taskID string, explicit []string) ([]string, error) {
	if len(explicit) > 0 {
		out := make([]string, 0, len(explicit))
		for _, p := range explicit {
			p = store.NormalizePath(strings.TrimSpace(p))
			if p != "" {
				out = append(out, p)
			}
		}
		sort.Strings(out)
		return out, nil
	}

	changes, err := dom.ListChangesByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	latest, ok := latestRecordedOrComparedChange(changes)
	if !ok {
		return nil, nil
	}
	cpaths, err := dom.ListChangePaths(ctx, latest.ID)
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	var out []string
	for _, cp := range cpaths {
		p := store.NormalizePath(cp.Path)
		if p == "" {
			continue
		}
		if _, dup := seen[p]; dup {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	sort.Strings(out)
	return out, nil
}

func latestRecordedOrComparedChange(changes []store.Change) (store.Change, bool) {
	var latest store.Change
	var found bool
	for _, c := range changes {
		if c.Status != store.ChangeStatusRecorded && c.Status != store.ChangeStatusCompared {
			continue
		}
		if !found || c.CreatedAt > latest.CreatedAt || (c.CreatedAt == latest.CreatedAt && c.ID > latest.ID) {
			latest = c
			found = true
		}
	}
	return latest, found
}

func targetsFromImpactAndValidates(ctx context.Context, st *store.Store, paths []string) ([]TestTarget, error) {
	eng := retrieval.New(st)
	byKey := map[targetKey]TestTarget{}

	addTarget := func(t TestTarget) {
		if t.Name == "" {
			return
		}
		k := targetKey{name: t.Name, pkg: t.Package}
		if _, ok := byKey[k]; ok {
			return
		}
		byKey[k] = t
	}

	var seeds []retrieval.ImpactSeed
	for _, p := range paths {
		f, err := st.GetFileByPath(p)
		if err != nil {
			continue
		}
		seeds = append(seeds, retrieval.ImpactSeed{EntityType: "file", EntityID: f.ID})
		addValidatesTargets(st, p, addTarget)
	}
	if len(seeds) > 0 {
		res, err := eng.ImpactWalk(ctx, seeds, 2)
		if err != nil {
			return nil, err
		}
		for _, hit := range res.AffectedTests {
			t, ok := targetFromBlastHit(st, hit)
			if ok {
				addTarget(t)
			}
		}
	}

	out := make([]TestTarget, 0, len(byKey))
	for _, t := range byKey {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].Package < out[j].Package
	})
	return out, nil
}

func addValidatesTargets(st *store.Store, changedPath string, add func(TestTarget)) {
	f, err := st.GetFileByPath(changedPath)
	if err != nil {
		return
	}
	edges, err := st.ListValidatesForFile(f.ID)
	if err != nil {
		return
	}
	for _, e := range edges {
		fromFile, err := st.GetFileByID(e.FromFileID)
		if err != nil {
			continue
		}
		name := filepath.Base(fromFile.Path)
		if e.FromSymbolID != nil && *e.FromSymbolID != "" {
			if sym, _, err := st.GetSymbolByID(*e.FromSymbolID); err == nil && sym.Name != "" {
				name = sym.Name
			}
		}
		pkg, ok := goPackageArg(st.ProjectRoot(), fromFile.Path)
		if !ok {
			continue
		}
		t := TestTarget{Name: name, Package: pkg, Path: fromFile.Path}
		if strings.HasPrefix(name, "Test") {
			t.RunPattern = "^" + name + "$"
		}
		add(t)
	}
}

func targetFromBlastHit(st *store.Store, hit retrieval.BlastHit) (TestTarget, bool) {
	path := hit.Path
	name := hit.Title
	if hit.EntityType == "symbol" {
		if sym, symPath, err := st.GetSymbolByID(hit.EntityID); err == nil {
			if sym.Name != "" {
				name = sym.Name
			}
			if path == "" {
				path = symPath
			}
		}
	} else if hit.EntityType == "file" {
		if f, err := st.GetFileByID(hit.EntityID); err == nil {
			path = f.Path
			if name == "" {
				name = filepath.Base(path)
			}
		}
	}
	if name == "" || path == "" {
		return TestTarget{}, false
	}
	pkg, ok := goPackageArg(st.ProjectRoot(), path)
	if !ok {
		return TestTarget{}, false
	}
	t := TestTarget{Name: name, Package: pkg, Path: path}
	if strings.HasPrefix(name, "Test") {
		t.RunPattern = "^" + name + "$"
	}
	return t, true
}

func packageFallbackTargets(root string, paths []string) ([]TestTarget, error) {
	mod, ok := goModulePath(root)
	if !ok {
		return nil, nil
	}
	seen := map[string]TestTarget{}
	for _, p := range paths {
		if !strings.HasSuffix(p, ".go") {
			continue
		}
		pkgDir, ok := goPackageArg(root, p)
		if !ok {
			continue
		}
		dir := strings.TrimPrefix(strings.TrimPrefix(pkgDir, "./"), "./")
		importPath := mod
		if dir != "" && dir != "." {
			importPath = mod + "/" + filepath.ToSlash(dir)
		}
		name := "package:" + importPath
		seen[name] = TestTarget{Name: name, Package: pkgDir}
	}
	out := make([]TestTarget, 0, len(seen))
	for _, t := range seen {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// jstsExts are the extensions the tree-sitter analyzers index for JS/TS.
var jstsExts = map[string]bool{
	".ts": true, ".tsx": true, ".mts": true, ".cts": true,
	".js": true, ".jsx": true, ".mjs": true, ".cjs": true,
}

// jstsFallbackTargets mirrors the Go package fallback for JS/TS repos: for a
// changed source file, the indexed test files in its sibling test/ or
// __tests__/ directory are the same "run this file's package tests" idea. Only
// already-indexed test files become targets, so the fallback is honest about
// what the graph knows (mirrors analyzers.isTestFile).
func jstsFallbackTargets(st *store.Store, paths []string) []TestTarget {
	seen := map[string]bool{}
	var out []TestTarget
	for _, p := range paths {
		p = store.NormalizePath(p)
		if !jstsExts[path.Ext(p)] {
			continue
		}
		dir := path.Dir(p)
		for _, testDir := range []string{
			path.Join(dir, "test"),
			path.Join(dir, "__tests__"),
			path.Join(path.Dir(dir), "test"),
			path.Join(path.Dir(dir), "__tests__"),
		} {
			files, err := st.ListFilePathsInDir(testDir)
			if err != nil {
				continue
			}
			for _, f := range files {
				if !jstsTestPath(f) || seen[f] {
					continue
				}
				seen[f] = true
				out = append(out, TestTarget{Name: path.Base(f), Path: f})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// jstsTestPath mirrors analyzers.isTestFile for JS/TS layouts.
func jstsTestPath(p string) bool {
	base := path.Base(p)
	stem := strings.TrimSuffix(base, path.Ext(base))
	return strings.HasSuffix(stem, ".test") || strings.HasSuffix(stem, ".spec") || path.Base(path.Dir(p)) == "__tests__"
}

// maxNamedPathsInError caps path lists inside the diagnosis so a 9-path
// failure does not produce a 9-line error (progressive disclosure).
const maxNamedPathsInError = 5

// noTargetsDiagnosis explains WHY selection produced zero targets instead of a
// bare "no relevant tests selected". Consumers could not previously tell "the
// graph has no validates edges for these files" from "the index is stale" from
// "this repo has no tests for that package" — on the reporting repo all nine
// changed paths were silently missing from a truncated index while `index
// status` claimed the index was fine.
func noTargetsDiagnosis(st *store.Store, paths []string) error {
	var notIndexed, notSource []string
	for _, p := range paths {
		p = store.NormalizePath(p)
		if _, err := st.GetFileByPath(p); err == nil {
			continue
		}
		if _, ok := analyzers.DetectLanguage(p); ok {
			notIndexed = append(notIndexed, p)
		} else {
			notSource = append(notSource, p)
		}
	}
	parts := make([]string, 0, 3)
	if len(notIndexed) > 0 {
		parts = append(parts, fmt.Sprintf("%d changed path(s) not in the index — run `trace index` to add them: %s",
			len(notIndexed), quotePathList(notIndexed)))
	}
	if len(notSource) > 0 {
		parts = append(parts, fmt.Sprintf("ignored non-source path(s): %s", quotePathList(notSource)))
	}
	parts = append(parts, "no incoming validates edges and no impact-walk tests for the indexed changed paths")
	return fmt.Errorf("testrun: no relevant tests selected: %s", strings.Join(parts, "; "))
}

func quotePathList(paths []string) string {
	named := paths
	more := 0
	if len(paths) > maxNamedPathsInError {
		named = paths[:maxNamedPathsInError]
		more = len(paths) - maxNamedPathsInError
	}
	quoted := make([]string, 0, len(named))
	for _, p := range named {
		quoted = append(quoted, strconv.Quote(p))
	}
	if more > 0 {
		quoted = append(quoted, fmt.Sprintf("…and %d more", more))
	}
	return strings.Join(quoted, ", ")
}

func goPackageArg(root, filePath string) (string, bool) {
	dir := filepath.Dir(store.NormalizePath(filePath))
	if dir == "." {
		return "./.", true
	}
	return "./" + filepath.ToSlash(dir), true
}
