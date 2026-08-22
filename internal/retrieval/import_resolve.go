package retrieval

import (
	"path"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// Source extensions tried when resolving extensionless relative imports (DF-60).
var importResolveExts = []string{
	".js", ".jsx", ".mjs", ".cjs", ".ts", ".tsx", ".go", ".py",
}

// isRelativeImport reports whether raw import text is a relative specifier
// (./ or ../), accepting backslash forms before slash-normalize.
func isRelativeImport(raw string) bool {
	n := strings.ReplaceAll(raw, `\`, "/")
	return strings.HasPrefix(n, "./") || strings.HasPrefix(n, "../")
}

func hasImportResolveExt(p string) bool {
	for _, ext := range importResolveExts {
		if strings.HasSuffix(p, ext) {
			return true
		}
	}
	return false
}

// importPathCandidates returns ordered, deduped lookup paths for an import string.
// Bare modules: NormalizePath exact only. Relative: exact, importer-dir join+Clean,
// then extension / index.* variants; first GetFileByPath hit wins at call site.
func importPathCandidates(importerPath, importedPath string) []string {
	var out []string
	seen := map[string]struct{}{}
	add := func(p string) {
		if p == "" || p == "." {
			return
		}
		if _, ok := seen[p]; ok {
			return
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}

	// Slash-normalize then Clean so `.\\util` / `.//util` become stable bases.
	slashImported := strings.ReplaceAll(importedPath, `\`, "/")
	exact := store.NormalizePath(path.Clean(slashImported))
	add(exact)

	if !isRelativeImport(importedPath) {
		return out
	}

	impDir := path.Dir(store.NormalizePath(importerPath))
	joined := store.NormalizePath(path.Clean(path.Join(impDir, slashImported)))
	add(joined)

	bases := append([]string(nil), out...)
	for _, base := range bases {
		if hasImportResolveExt(base) {
			continue
		}
		for _, ext := range importResolveExts {
			add(base + ext)
		}
	}

	if !hasImportResolveExt(joined) {
		indexBase := path.Join(joined, "index")
		for _, ext := range importResolveExts {
			add(indexBase + ext)
		}
	}

	return out
}

// resolveImportedFile returns the first indexed file matching import candidates,
// or a not-found error (isNotFound) when none match.
func (e *Engine) resolveImportedFile(importerPath, importedPath string) (store.FileRecord, error) {
	var lastNotFound error
	for _, cand := range importPathCandidates(importerPath, importedPath) {
		f, err := e.store.GetFileByPath(cand)
		if err != nil {
			if isNotFound(err) {
				lastNotFound = err
				continue
			}
			return store.FileRecord{}, err
		}
		return f, nil
	}
	if lastNotFound != nil {
		return store.FileRecord{}, lastNotFound
	}
	// Empty candidate list (should not happen for non-empty import).
	f, err := e.store.GetFileByPath(store.NormalizePath(importedPath))
	return f, err
}
