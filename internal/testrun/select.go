package testrun

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

// TestTarget is one executed test invocation recorded as kind=test.
type TestTarget struct {
	Name       string
	Package    string // go test package arg, e.g. ./internal/foo
	RunPattern string // optional -run value
}

type targetKey struct {
	name string
	pkg  string
}

// SelectTestTargets picks relevant tests from seed paths via impact walk + validates; Go package fallback.
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
	return packageFallbackTargets(st.ProjectRoot(), paths)
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
		t := TestTarget{Name: name, Package: pkg}
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
	t := TestTarget{Name: name, Package: pkg}
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

func goPackageArg(root, filePath string) (string, bool) {
	dir := filepath.Dir(store.NormalizePath(filePath))
	if dir == "." {
		return "./.", true
	}
	return "./" + filepath.ToSlash(dir), true
}
