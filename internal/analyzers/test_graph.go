package analyzers

import (
	"fmt"
	"path"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

func isTestFile(p string) bool {
	p = store.NormalizePath(p)
	base := path.Base(p)
	lang, ok := DetectLanguage(p)
	if !ok {
		return false
	}
	switch lang {
	case LangGo:
		return strings.HasSuffix(p, "_test.go")
	case LangPython:
		return (strings.HasPrefix(base, "test_") && strings.HasSuffix(base, ".py")) || strings.HasSuffix(base, "_test.py")
	case LangJavaScript, LangTypeScript, LangTSX:
		return isJSTSTestPath(p)
	default:
		return false
	}
}

func isJSTSTestPath(p string) bool {
	p = store.NormalizePath(p)
	base := path.Base(p)
	ext := path.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	if strings.HasSuffix(stem, ".test") || strings.HasSuffix(stem, ".spec") {
		return true
	}
	return path.Base(path.Dir(p)) == "__tests__"
}

func isGoTestName(name string) bool {
	for _, prefix := range []string{"Test", "Benchmark", "Example", "Fuzz"} {
		if name == prefix {
			return true
		}
		if strings.HasPrefix(name, prefix) && len(name) > len(prefix) {
			c := name[len(prefix)]
			if c >= 'A' && c <= 'Z' {
				return true
			}
		}
	}
	return false
}

func classifyTestSymbols(p string, symbols []store.Symbol) []store.Symbol {
	if !isTestFile(p) {
		return symbols
	}
	lang, _ := DetectLanguage(p)
	out := append([]store.Symbol(nil), symbols...)
	for i, s := range out {
		switch lang {
		case LangGo:
			if s.Kind == "function" && isGoTestName(s.Name) {
				out[i].Kind = "test"
			}
		case LangPython:
			if (s.Kind == "function" || s.Kind == "method") && strings.HasPrefix(s.Name, "test_") {
				out[i].Kind = "test"
			}
			if s.Kind == "class" && strings.HasPrefix(s.Name, "Test") {
				out[i].Kind = "test"
			}
		}
	}
	return out
}

func indexCodeEdges(st *store.Store, p string, content []byte, lang string) error {
	var goPkg string
	if lang == LangGo {
		goPkg = extractGoPackageName(content)
	}
	var edges []store.CodeEdge
	if isTestFile(p) {
		v, err := computeValidatesEdges(st, p, goPkg)
		if err != nil {
			return err
		}
		edges = append(edges, v...)
	}
	art, err := computeArtifactEdges(st, p, content, lang)
	if err != nil {
		return err
	}
	edges = append(edges, art...)
	bounds, err := computeArchitecturalBoundaryEdges(st, p)
	if err != nil {
		return err
	}
	edges = append(edges, bounds...)
	// One ReplaceFileEdges after pre-clear: validates + contains_module +
	// exports_api + architectural_boundary. A second replace deletes the rest.
	if edges == nil {
		edges = []store.CodeEdge{}
	}
	if err := st.ReplaceFileEdges(p, edges); err != nil {
		return fmt.Errorf("analyzers: replace edges: %w", err)
	}
	return upsertIncomingValidates(st, p)
}

func upsertIncomingValidates(st *store.Store, targetPath string) error {
	targetPath = store.NormalizePath(targetPath)
	tf, err := st.GetFileByPath(targetPath)
	if err != nil {
		return err
	}
	paths, err := st.ListFilePaths()
	if err != nil {
		return fmt.Errorf("analyzers: list files for incoming validates: %w", err)
	}
	for _, p := range paths {
		if p == targetPath || !isTestFile(p) {
			continue
		}
		edges, err := computeValidatesEdges(st, p, "")
		if err != nil {
			return err
		}
		var toTarget []store.CodeEdge
		for _, e := range edges {
			if e.ToFileID == tf.ID {
				toTarget = append(toTarget, e)
			}
		}
		if len(toTarget) == 0 {
			continue
		}
		if err := st.UpsertFilePairEdges(p, targetPath, store.RelValidates, toTarget); err != nil {
			return fmt.Errorf("analyzers: upsert incoming validates: %w", err)
		}
	}
	return nil
}

func computeValidatesEdges(st *store.Store, testPath, goPkg string) ([]store.CodeEdge, error) {
	testPath = store.NormalizePath(testPath)
	testFile, err := st.GetFileByPath(testPath)
	if err != nil {
		return nil, err
	}
	lang, _ := DetectLanguage(testPath)
	imps, err := st.ListImportsByPath(testPath)
	if err != nil {
		return nil, err
	}
	extracted, err := extractedValidates(st, testPath, testFile, lang, imps)
	if err != nil {
		return nil, err
	}
	var edges []store.CodeEdge
	edges = append(edges, extracted...)

	if goPkg == "" && lang == LangGo {
		goPkg = goPackageNameFromSymbols(st, testPath)
	}
	inferredOK := lang == LangGo && strings.HasSuffix(testPath, "_test.go")
	if goPkg != "" && strings.HasSuffix(goPkg, "_test") {
		inferredOK = false
	}
	if len(extracted) > 0 {
		inferredOK = false
	}
	if inferredOK {
		testSyms, err := st.ListSymbolsByPath(testPath)
		if err != nil {
			return nil, err
		}
		edges = append(edges, inferredGoNamePrefix(st, testPath, testFile, testSyms)...)
	}
	return edges, nil
}

func extractedValidates(st *store.Store, testPath string, testFile store.FileRecord, lang string, imps []store.Import) ([]store.CodeEdge, error) {
	var out []store.CodeEdge
	for _, im := range imps {
		targets := resolveImportTargets(st, testPath, lang, im.ImportedPath)
		for _, tpath := range targets {
			tf, err := st.GetFileByPath(tpath)
			if err != nil {
				continue
			}
			var toSym *string
			if im.Symbol != nil && *im.Symbol != "" {
				syms, err := st.ListSymbolsByPath(tpath)
				if err != nil {
					return nil, err
				}
				for _, s := range syms {
					if s.Name == *im.Symbol && s.Kind != "test" {
						id := s.ID
						toSym = &id
						break
					}
				}
			}
			out = append(out, store.CodeEdge{
				FromFileID: testFile.ID,
				ToFileID:   tf.ID,
				ToSymbolID: toSym,
				Rel:        store.RelValidates,
				Provenance: store.ImportProvenanceExtracted,
			})
		}
	}
	return out, nil
}

func resolveImportTargets(st *store.Store, fromPath, lang, imported string) []string {
	switch lang {
	case LangJavaScript, LangTypeScript, LangTSX:
		return existingIndexedPaths(st, jsRelativeCandidates(fromPath, imported))
	case LangPython:
		return existingIndexedPaths(st, pythonModuleCandidates(fromPath, imported))
	case LangGo:
		return goPackageSiblings(st, fromPath, imported)
	default:
		return nil
	}
}

func existingIndexedPaths(st *store.Store, candidates []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, c := range candidates {
		c = store.NormalizePath(c)
		if c == "" || seen[c] {
			continue
		}
		if _, err := st.GetFileByPath(c); err != nil {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}

func jsRelativeCandidates(fromPath, imported string) []string {
	if !strings.HasPrefix(imported, ".") {
		return nil
	}
	joined := store.NormalizePath(path.Clean(path.Join(path.Dir(fromPath), imported)))
	cands := []string{joined}
	if path.Ext(joined) == "" {
		for _, e := range []string{".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs"} {
			cands = append(cands, joined+e)
		}
	}
	return cands
}

func pythonModuleCandidates(fromPath, imported string) []string {
	dir := path.Dir(fromPath)
	if strings.HasPrefix(imported, ".") {
		i := 0
		for i < len(imported) && imported[i] == '.' {
			i++
		}
		d := dir
		for n := 1; n < i; n++ {
			d = path.Dir(d)
		}
		rest := strings.ReplaceAll(imported[i:], ".", "/")
		if rest == "" {
			return []string{store.NormalizePath(path.Join(d, "__init__.py"))}
		}
		return []string{
			store.NormalizePath(path.Join(d, rest+".py")),
			store.NormalizePath(path.Join(d, rest, "__init__.py")),
		}
	}
	rel := strings.ReplaceAll(imported, ".", "/")
	return []string{
		store.NormalizePath(path.Join(dir, rel+".py")),
		store.NormalizePath(path.Join(dir, rel, "__init__.py")),
	}
}

func goPackageSiblings(st *store.Store, fromPath, imported string) []string {
	last := imported
	if i := strings.LastIndex(imported, "/"); i >= 0 {
		last = imported[i+1:]
	}
	dir := path.Dir(fromPath)
	if last == "" || last != path.Base(dir) {
		return nil
	}
	paths, err := st.ListFilePaths()
	if err != nil {
		return nil
	}
	var out []string
	for _, p := range paths {
		if path.Dir(p) != dir {
			continue
		}
		if strings.HasSuffix(p, "_test.go") || !strings.HasSuffix(p, ".go") {
			continue
		}
		out = append(out, p)
	}
	return out
}

// inferredGoNamePrefix implements heuristic:go_test_name_prefix (INFERRED only).
// Same-directory _test.go whose package is not *_test: TestFoo → sibling symbol Foo
// (strip Test / Benchmark / Example / Fuzz).
func inferredGoNamePrefix(st *store.Store, testPath string, testFile store.FileRecord, testSyms []store.Symbol) []store.CodeEdge {
	dir := path.Dir(testPath)
	paths, err := st.ListFilePaths()
	if err != nil {
		return nil
	}
	type hit struct {
		fileID string
		sym    store.Symbol
	}
	var siblings []hit
	for _, p := range paths {
		if p == testPath || path.Dir(p) != dir {
			continue
		}
		if strings.HasSuffix(p, "_test.go") || !strings.HasSuffix(p, ".go") {
			continue
		}
		syms, err := st.ListSymbolsByPath(p)
		if err != nil {
			continue
		}
		for _, s := range syms {
			if s.Kind == "test" {
				continue
			}
			siblings = append(siblings, hit{fileID: s.FileID, sym: s})
		}
	}
	var out []store.CodeEdge
	for _, ts := range testSyms {
		if ts.Kind != "test" {
			continue
		}
		target := stripGoTestPrefix(ts.Name)
		if target == "" {
			continue
		}
		for _, sib := range siblings {
			if sib.sym.Name != target {
				continue
			}
			fromID := ts.ID
			toID := sib.sym.ID
			out = append(out, store.CodeEdge{
				FromFileID:   testFile.ID,
				FromSymbolID: &fromID,
				ToFileID:     sib.fileID,
				ToSymbolID:   &toID,
				Rel:          store.RelValidates,
				Provenance:   store.ImportProvenanceInferred,
			})
			break
		}
	}
	return out
}

func goPackageNameFromSymbols(st *store.Store, testPath string) string {
	syms, err := st.ListSymbolsByPath(testPath)
	if err != nil {
		return ""
	}
	for _, s := range syms {
		if s.Kind == "package" && s.Name != "" {
			return s.Name
		}
	}
	return ""
}

func stripGoTestPrefix(name string) string {
	for _, prefix := range []string{"Test", "Benchmark", "Example", "Fuzz"} {
		if !strings.HasPrefix(name, prefix) || len(name) <= len(prefix) {
			continue
		}
		rest := name[len(prefix):]
		if rest[0] >= 'A' && rest[0] <= 'Z' {
			return rest
		}
	}
	return ""
}
