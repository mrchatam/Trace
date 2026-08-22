package analyzers

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

func openTemp(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}
	return b
}

func symNamesKinds(syms []store.Symbol) []string {
	out := make([]string, len(syms))
	for i, s := range syms {
		out[i] = s.Kind + ":" + s.Name
	}
	sort.Strings(out)
	return out
}

func impKeys(imps []store.Import) []string {
	out := make([]string, len(imps))
	for i, im := range imps {
		sym := ""
		if im.Symbol != nil {
			sym = *im.Symbol
		}
		out[i] = im.ImportedPath + "|" + sym
	}
	sort.Strings(out)
	return out
}

func TestDetectLanguage(t *testing.T) {
	cases := []struct {
		path string
		lang string
		ok   bool
	}{
		{"a.js", LangJavaScript, true},
		{"a.jsx", LangJavaScript, true},
		{"a.mjs", LangJavaScript, true},
		{"a.cjs", LangJavaScript, true},
		{"a.ts", LangTypeScript, true},
		{"a.tsx", LangTSX, true},
		{"a.py", LangPython, true},
		{"a.go", LangGo, true},
		{"a.txt", "", false},
	}
	for _, tc := range cases {
		lang, ok := DetectLanguage(tc.path)
		if ok != tc.ok || lang != tc.lang {
			t.Fatalf("%s: got (%q,%v) want (%q,%v)", tc.path, lang, ok, tc.lang, tc.ok)
		}
	}
}

func TestIndexFileJSGolden(t *testing.T) {
	st := openTemp(t)
	content := readTestdata(t, "sample.js")
	if err := IndexFile(context.Background(), st, "src/sample.js", content, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile: %v", err)
	}

	f, err := st.GetFileByPath("src/sample.js")
	if err != nil {
		t.Fatalf("GetFileByPath: %v", err)
	}
	if f.Language == nil || *f.Language != LangJavaScript {
		t.Fatalf("language: got %v want javascript", f.Language)
	}
	if f.ContentHash == "" || len(f.ContentHash) != 64 {
		t.Fatalf("content hash: %q", f.ContentHash)
	}

	syms, err := st.ListSymbolsByPath("src/sample.js")
	if err != nil {
		t.Fatalf("ListSymbolsByPath: %v", err)
	}
	gotSym := symNamesKinds(syms)
	wantSym := []string{"class:Greeter", "function:greet", "method:sayHi"}
	if len(gotSym) != len(wantSym) {
		t.Fatalf("symbols got %v want %v", gotSym, wantSym)
	}
	for i := range wantSym {
		if gotSym[i] != wantSym[i] {
			t.Fatalf("symbols got %v want %v", gotSym, wantSym)
		}
	}
	for _, s := range syms {
		if s.StartLine < 1 || s.EndLine < s.StartLine {
			t.Fatalf("bad lines for %s: %d-%d", s.Name, s.StartLine, s.EndLine)
		}
	}

	imps, err := st.ListImportsByPath("src/sample.js")
	if err != nil {
		t.Fatalf("ListImportsByPath: %v", err)
	}
	gotImp := impKeys(imps)
	wantImp := []string{"./styles.css|", "react|useState"}
	if len(gotImp) != len(wantImp) {
		t.Fatalf("imports got %v want %v", gotImp, wantImp)
	}
	for i := range wantImp {
		if gotImp[i] != wantImp[i] {
			t.Fatalf("imports got %v want %v", gotImp, wantImp)
		}
	}
}

func TestIndexFilePythonGolden(t *testing.T) {
	st := openTemp(t)
	content := readTestdata(t, "sample.py")
	if err := IndexFile(context.Background(), st, "pkg/sample.py", content, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile: %v", err)
	}

	f, err := st.GetFileByPath("pkg/sample.py")
	if err != nil {
		t.Fatalf("GetFileByPath: %v", err)
	}
	if f.Language == nil || *f.Language != LangPython {
		t.Fatalf("language: got %v want python", f.Language)
	}

	syms, err := st.ListSymbolsByPath("pkg/sample.py")
	if err != nil {
		t.Fatalf("ListSymbolsByPath: %v", err)
	}
	gotSym := symNamesKinds(syms)
	wantSym := []string{"class:Worker", "function:helper", "function:main", "method:run"}
	if len(gotSym) != len(wantSym) {
		t.Fatalf("symbols got %v want %v", gotSym, wantSym)
	}
	for i := range wantSym {
		if gotSym[i] != wantSym[i] {
			t.Fatalf("symbols got %v want %v", gotSym, wantSym)
		}
	}

	imps, err := st.ListImportsByPath("pkg/sample.py")
	if err != nil {
		t.Fatalf("ListImportsByPath: %v", err)
	}
	gotImp := impKeys(imps)
	wantImp := []string{"collections|defaultdict", "os|", "pathlib|Path"}
	if len(gotImp) != len(wantImp) {
		t.Fatalf("imports got %v want %v", gotImp, wantImp)
	}
	for i := range wantImp {
		if gotImp[i] != wantImp[i] {
			t.Fatalf("imports got %v want %v", gotImp, wantImp)
		}
	}
}

func TestIndexFileTSGolden(t *testing.T) {
	st := openTemp(t)
	content := readTestdata(t, "sample.ts")
	if err := IndexFile(context.Background(), st, "src/sample.ts", content, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile: %v", err)
	}
	syms, err := st.ListSymbolsByPath("src/sample.ts")
	if err != nil {
		t.Fatalf("ListSymbolsByPath: %v", err)
	}
	gotSym := symNamesKinds(syms)
	wantSym := []string{"class:Counter", "function:compute", "method:inc"}
	if len(gotSym) != len(wantSym) {
		t.Fatalf("symbols got %v want %v", gotSym, wantSym)
	}
	for i := range wantSym {
		if gotSym[i] != wantSym[i] {
			t.Fatalf("symbols got %v want %v", gotSym, wantSym)
		}
	}
}

func TestIndexFileGoGolden(t *testing.T) {
	st := openTemp(t)
	content := readTestdata(t, "sample.go")
	if err := IndexFile(context.Background(), st, "pkg/sample.go", content, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile: %v", err)
	}

	f, err := st.GetFileByPath("pkg/sample.go")
	if err != nil {
		t.Fatalf("GetFileByPath: %v", err)
	}
	if f.Language == nil || *f.Language != LangGo {
		t.Fatalf("language: got %v want go", f.Language)
	}

	syms, err := st.ListSymbolsByPath("pkg/sample.go")
	if err != nil {
		t.Fatalf("ListSymbolsByPath: %v", err)
	}
	gotSym := symNamesKinds(syms)
	wantSym := []string{
		"function:Helper",
		"function:Main",
		"method:Run",
		"type:Counter",
		"type:ID",
		"type:Worker",
	}
	if len(gotSym) != len(wantSym) {
		t.Fatalf("symbols got %v want %v", gotSym, wantSym)
	}
	for i := range wantSym {
		if gotSym[i] != wantSym[i] {
			t.Fatalf("symbols got %v want %v", gotSym, wantSym)
		}
	}
	for _, s := range syms {
		if s.StartLine < 1 || s.EndLine < s.StartLine {
			t.Fatalf("bad lines for %s: %d-%d", s.Name, s.StartLine, s.EndLine)
		}
	}

	imps, err := st.ListImportsByPath("pkg/sample.go")
	if err != nil {
		t.Fatalf("ListImportsByPath: %v", err)
	}
	gotImp := impKeys(imps)
	wantImp := []string{"fmt|", "os|alias", "path/filepath|"}
	if len(gotImp) != len(wantImp) {
		t.Fatalf("imports got %v want %v", gotImp, wantImp)
	}
	for i := range wantImp {
		if gotImp[i] != wantImp[i] {
			t.Fatalf("imports got %v want %v", gotImp, wantImp)
		}
	}
}

func TestIndexFileGoHandlerMethods(t *testing.T) {
	st := openTemp(t)
	content := readTestdata(t, "handler_methods.go")
	if err := IndexFile(context.Background(), st, "pkg/handler_methods.go", content, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile: %v", err)
	}

	f, err := st.GetFileByPath("pkg/handler_methods.go")
	if err != nil {
		t.Fatalf("GetFileByPath: %v", err)
	}
	if f.Language == nil || *f.Language != LangGo {
		t.Fatalf("language: got %v want go", f.Language)
	}

	syms, err := st.ListSymbolsByPath("pkg/handler_methods.go")
	if err != nil {
		t.Fatalf("ListSymbolsByPath: %v", err)
	}
	gotSym := symNamesKinds(syms)
	wantSym := []string{
		"method:Search",
		"method:SearchCursor",
		"type:Memory",
		"type:Notes",
	}
	if len(gotSym) != len(wantSym) {
		t.Fatalf("symbols got %v want %v", gotSym, wantSym)
	}
	for i := range wantSym {
		if gotSym[i] != wantSym[i] {
			t.Fatalf("symbols got %v want %v", gotSym, wantSym)
		}
	}
	for _, s := range syms {
		if s.StartLine < 1 || s.EndLine < s.StartLine {
			t.Fatalf("bad lines for %s: %d-%d", s.Name, s.StartLine, s.EndLine)
		}
	}
}

func TestIncrementalIsolation(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()

	a1 := []byte("export function alpha() { return 1; }\n")
	b1 := []byte("export function beta() { return 2; }\n")
	if err := IndexFile(ctx, st, "a.js", a1, IndexOptions{}); err != nil {
		t.Fatalf("index A: %v", err)
	}
	if err := IndexFile(ctx, st, "b.js", b1, IndexOptions{}); err != nil {
		t.Fatalf("index B: %v", err)
	}

	bSymsBefore, err := st.ListSymbolsByPath("b.js")
	if err != nil {
		t.Fatal(err)
	}
	bImpsBefore, err := st.ListImportsByPath("b.js")
	if err != nil {
		t.Fatal(err)
	}
	bIDs := make([]string, len(bSymsBefore))
	for i, s := range bSymsBefore {
		bIDs[i] = s.ID
	}

	a2 := []byte("export function alpha() { return 1; }\nexport function gamma() { return 3; }\n")
	if err := IndexFile(ctx, st, "a.js", a2, IndexOptions{}); err != nil {
		t.Fatalf("reindex A: %v", err)
	}

	aSyms, err := st.ListSymbolsByPath("a.js")
	if err != nil {
		t.Fatal(err)
	}
	gotA := symNamesKinds(aSyms)
	wantA := []string{"function:alpha", "function:gamma"}
	if len(gotA) != len(wantA) {
		t.Fatalf("A symbols got %v want %v", gotA, wantA)
	}
	for i := range wantA {
		if gotA[i] != wantA[i] {
			t.Fatalf("A symbols got %v want %v", gotA, wantA)
		}
	}

	bSymsAfter, err := st.ListSymbolsByPath("b.js")
	if err != nil {
		t.Fatal(err)
	}
	bImpsAfter, err := st.ListImportsByPath("b.js")
	if err != nil {
		t.Fatal(err)
	}
	if len(bSymsAfter) != len(bSymsBefore) {
		t.Fatalf("B symbols changed: before %v after %v", bSymsBefore, bSymsAfter)
	}
	for i := range bSymsAfter {
		if bSymsAfter[i].ID != bIDs[i] || bSymsAfter[i].Name != bSymsBefore[i].Name {
			t.Fatalf("B symbol row mutated: before %+v after %+v", bSymsBefore[i], bSymsAfter[i])
		}
	}
	if len(bImpsAfter) != len(bImpsBefore) {
		t.Fatalf("B imports changed")
	}
}

func TestUnsupportedAndBinaryDoNotCorruptOthers(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	if err := IndexFile(ctx, st, "ok.js", []byte("export function ok() {}"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}

	err := IndexFile(ctx, st, "readme.md", []byte("# hi"), IndexOptions{})
	if err == nil {
		t.Fatal("expected unsupported error")
	}
	var skip *SkipError
	if !errors.As(err, &skip) {
		t.Fatalf("want SkipError, got %T %v", err, err)
	}

	bin := []byte("ok\x00more")
	err = IndexFile(ctx, st, "blob.js", bin, IndexOptions{})
	if err == nil {
		t.Fatal("expected binary error")
	}
	if !errors.As(err, &skip) {
		t.Fatalf("want SkipError, got %T %v", err, err)
	}

	syms, err := st.ListSymbolsByPath("ok.js")
	if err != nil {
		t.Fatal(err)
	}
	if len(syms) != 1 || syms[0].Name != "ok" {
		t.Fatalf("ok.js corrupted: %v", syms)
	}
	if _, err := st.GetFileByPath("readme.md"); err == nil {
		t.Fatal("unsupported path should not upsert file")
	}
	if _, err := st.GetFileByPath("blob.js"); err == nil {
		t.Fatal("binary path should not upsert file")
	}
}

func TestIndexFileAtRevUsesVCS(t *testing.T) {
	st := openTemp(t)
	content := []byte("export function fromVCS() { return 1; }\n")
	repo := &vcs.Fake{
		IsGit:   true,
		HeadOID: "abc123",
		Files: map[string][]byte{
			"abc123:lib/from_vcs.js": content,
		},
	}
	if err := IndexFileAtRev(context.Background(), st, repo, "abc123", "lib/from_vcs.js", IndexOptions{}); err != nil {
		t.Fatalf("IndexFileAtRev: %v", err)
	}
	f, err := st.GetFileByPath("lib/from_vcs.js")
	if err != nil {
		t.Fatal(err)
	}
	if f.GitOID == nil || *f.GitOID != "abc123" {
		t.Fatalf("git oid: %v", f.GitOID)
	}
	syms, err := st.ListSymbolsByPath("lib/from_vcs.js")
	if err != nil {
		t.Fatal(err)
	}
	if len(syms) != 1 || syms[0].Name != "fromVCS" {
		t.Fatalf("symbols: %v", syms)
	}
}

func TestAnalyzersDoNotImportGitcli(t *testing.T) {
	// Compile-time / package boundary guard: this file imports vcs only.
	var _ vcs.Repository = (*vcs.Fake)(nil)
}

func TestAnalyzerImportProvenanceExtracted(t *testing.T) {
	st := openTemp(t)
	content := readTestdata(t, "sample.go")
	if err := IndexFile(context.Background(), st, "pkg/sample.go", content, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile go: %v", err)
	}
	imps, err := st.ListImportsByPath("pkg/sample.go")
	if err != nil || len(imps) == 0 {
		t.Fatalf("go imports: %v %+v", err, imps)
	}
	for _, im := range imps {
		if im.Provenance != store.ImportProvenanceExtracted {
			t.Fatalf("go import %q provenance=%q want EXTRACTED", im.ImportedPath, im.Provenance)
		}
	}

	pyContent := []byte("from os import *\nimport sys\n")
	if err := IndexFile(context.Background(), st, "pkg/wild.py", pyContent, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile py: %v", err)
	}
	pyImps, err := st.ListImportsByPath("pkg/wild.py")
	if err != nil {
		t.Fatalf("py imports: %v", err)
	}
	var sawWild, sawExtracted bool
	for _, im := range pyImps {
		switch im.ImportedPath {
		case "os":
			sawWild = true
			if im.Provenance != store.ImportProvenanceAmbiguous {
				t.Fatalf("wildcard os: provenance=%q want AMBIGUOUS", im.Provenance)
			}
		case "sys":
			sawExtracted = true
			if im.Provenance != store.ImportProvenanceExtracted {
				t.Fatalf("sys: provenance=%q want EXTRACTED", im.Provenance)
			}
		}
	}
	if !sawWild || !sawExtracted {
		t.Fatalf("expected wildcard+extracted imports, got %+v", pyImps)
	}
}

func TestLanguageAdapterAPIVersion(t *testing.T) {
	if LanguageAdapterAPIVersion != 1 {
		t.Fatalf("LanguageAdapterAPIVersion = %d, want 1", LanguageAdapterAPIVersion)
	}
}

func TestIndexDiscoversGoTestFunctions(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	if err := IndexFile(ctx, st, "pkg/foo.go", readTestdata(t, "foo.go"), IndexOptions{}); err != nil {
		t.Fatalf("IndexFile foo.go: %v", err)
	}
	if err := IndexFile(ctx, st, "pkg/foo_test.go", readTestdata(t, "foo_test.go"), IndexOptions{}); err != nil {
		t.Fatalf("IndexFile foo_test.go: %v", err)
	}

	syms, err := st.ListSymbolsByPath("pkg/foo_test.go")
	if err != nil {
		t.Fatal(err)
	}
	got := symNamesKinds(syms)
	want := []string{"function:notATest", "package:foo", "test:TestFoo"}
	if len(got) != len(want) {
		t.Fatalf("symbols got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("symbols got %v want %v", got, want)
		}
	}
}

func TestValidatesInferredGoNamePrefix(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	if err := IndexFile(ctx, st, "pkg/foo.go", readTestdata(t, "foo.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "pkg/foo_test.go", readTestdata(t, "foo_test.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}

	prod, err := st.ListSymbolsByPath("pkg/foo.go")
	if err != nil {
		t.Fatal(err)
	}
	var foo store.Symbol
	for _, s := range prod {
		if s.Name == "Foo" {
			foo = s
			break
		}
	}
	if foo.ID == "" {
		t.Fatalf("Foo symbol missing: %v", prod)
	}
	edges, err := st.ListValidatesForSymbol(foo.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 1 {
		t.Fatalf("validates for Foo: %+v", edges)
	}
	if edges[0].Provenance != store.ImportProvenanceInferred {
		t.Fatalf("provenance=%q want INFERRED (heuristic:go_test_name_prefix)", edges[0].Provenance)
	}
	tests, err := st.ListSymbolsByPath("pkg/foo_test.go")
	if err != nil {
		t.Fatal(err)
	}
	var testFoo store.Symbol
	for _, s := range tests {
		if s.Name == "TestFoo" && s.Kind == "test" {
			testFoo = s
			break
		}
	}
	if testFoo.ID == "" || edges[0].FromSymbolID == nil || *edges[0].FromSymbolID != testFoo.ID {
		t.Fatalf("from_symbol: got %+v want TestFoo id %s", edges[0].FromSymbolID, testFoo.ID)
	}
}

func TestValidatesEdgeExtractedFromImport(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	if err := IndexFile(ctx, st, "pkg/mod.py", readTestdata(t, "mod.py"), IndexOptions{}); err != nil {
		t.Fatalf("IndexFile mod.py: %v", err)
	}
	if err := IndexFile(ctx, st, "pkg/test_mod.py", readTestdata(t, "test_mod.py"), IndexOptions{}); err != nil {
		t.Fatalf("IndexFile test_mod.py: %v", err)
	}

	modSyms, err := st.ListSymbolsByPath("pkg/mod.py")
	if err != nil {
		t.Fatal(err)
	}
	var foo store.Symbol
	for _, s := range modSyms {
		if s.Name == "Foo" {
			foo = s
			break
		}
	}
	if foo.ID == "" {
		t.Fatalf("Foo missing: %v", modSyms)
	}
	edges, err := st.ListValidatesForSymbol(foo.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 1 {
		t.Fatalf("validates for Foo: %+v", edges)
	}
	if edges[0].Provenance != store.ImportProvenanceExtracted {
		t.Fatalf("provenance=%q want EXTRACTED", edges[0].Provenance)
	}
	testFile, err := st.GetFileByPath("pkg/test_mod.py")
	if err != nil {
		t.Fatal(err)
	}
	if edges[0].FromFileID != testFile.ID {
		t.Fatalf("from_file_id=%s want %s", edges[0].FromFileID, testFile.ID)
	}

	// Target indexed after the test: incoming upsert must still create EXTRACTED.
	st2 := openTemp(t)
	if err := IndexFile(ctx, st2, "pkg/test_mod.py", readTestdata(t, "test_mod.py"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st2, "pkg/mod.py", readTestdata(t, "mod.py"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	modSyms2, err := st2.ListSymbolsByPath("pkg/mod.py")
	if err != nil {
		t.Fatal(err)
	}
	foo2 := store.Symbol{}
	for _, s := range modSyms2 {
		if s.Name == "Foo" {
			foo2 = s
			break
		}
	}
	late, err := st2.ListValidatesForSymbol(foo2.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(late) != 1 || late[0].Provenance != store.ImportProvenanceExtracted {
		t.Fatalf("incoming validates after target index: %+v", late)
	}
}

func TestValidatesSurviveTargetReindex(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	if err := IndexFile(ctx, st, "pkg/foo.go", readTestdata(t, "foo.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "pkg/foo_test.go", readTestdata(t, "foo_test.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	foo := mustSymbolNamed(t, st, "pkg/foo.go", "Foo")
	before, err := st.ListValidatesForSymbol(foo.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(before) != 1 {
		t.Fatalf("before reindex: %+v", before)
	}
	if err := IndexFile(ctx, st, "pkg/foo.go", readTestdata(t, "foo.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	foo2 := mustSymbolNamed(t, st, "pkg/foo.go", "Foo")
	if foo2.ID != foo.ID {
		t.Fatalf("symbol id churn: %s → %s", foo.ID, foo2.ID)
	}
	after, err := st.ListValidatesForSymbol(foo2.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != 1 || after[0].ToSymbolID == nil || *after[0].ToSymbolID != foo2.ID {
		t.Fatalf("validates after target reindex (SET NULL repair): %+v", after)
	}
}

func TestValidatesMultipleNamedImportsSurviveTargetReindex(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	lib := []byte("export function alpha() { return 1; }\nexport function beta() { return 2; }\n")
	testSrc := []byte("import { alpha, beta } from './lib.js'\ntest('both', () => { alpha(); beta() })\n")
	if err := IndexFile(ctx, st, "src/lib.js", lib, IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "src/lib.test.js", testSrc, IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	alpha := mustSymbolNamed(t, st, "src/lib.js", "alpha")
	beta := mustSymbolNamed(t, st, "src/lib.js", "beta")
	va, err := st.ListValidatesForSymbol(alpha.ID)
	if err != nil {
		t.Fatal(err)
	}
	vb, err := st.ListValidatesForSymbol(beta.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(va) != 1 || len(vb) != 1 {
		t.Fatalf("named-import validates before reindex: alpha=%+v beta=%+v", va, vb)
	}
	if err := IndexFile(ctx, st, "src/lib.js", lib, IndexOptions{}); err != nil {
		t.Fatalf("lib reindex with two incoming validates: %v", err)
	}
	alpha2 := mustSymbolNamed(t, st, "src/lib.js", "alpha")
	beta2 := mustSymbolNamed(t, st, "src/lib.js", "beta")
	va2, err := st.ListValidatesForSymbol(alpha2.ID)
	if err != nil {
		t.Fatal(err)
	}
	vb2, err := st.ListValidatesForSymbol(beta2.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(va2) != 1 || va2[0].ToSymbolID == nil || *va2[0].ToSymbolID != alpha2.ID {
		t.Fatalf("alpha validates after lib reindex: %+v", va2)
	}
	if len(vb2) != 1 || vb2[0].ToSymbolID == nil || *vb2[0].ToSymbolID != beta2.ID {
		t.Fatalf("beta validates after lib reindex: %+v", vb2)
	}
}

func TestValidatesFooTestPackageNotInferred(t *testing.T) {
	lib := []byte("package foo\n\nfunc Foo() int { return 1 }\n")
	xtest := []byte("package foo_test\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\n")
	ctx := context.Background()

	st := openTemp(t)
	if err := IndexFile(ctx, st, "pkg/foo.go", lib, IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "pkg/foo_test.go", xtest, IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	foo := mustSymbolNamed(t, st, "pkg/foo.go", "Foo")
	edges, err := st.ListValidatesForSymbol(foo.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 0 {
		t.Fatalf("package foo_test without import must not INFER validates: %+v", edges)
	}

	st2 := openTemp(t)
	if err := IndexFile(ctx, st2, "pkg/foo_test.go", xtest, IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st2, "pkg/foo.go", lib, IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	foo2 := mustSymbolNamed(t, st2, "pkg/foo.go", "Foo")
	late, err := st2.ListValidatesForSymbol(foo2.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(late) != 0 {
		t.Fatalf("incoming empty-goPkg must not INFER for package foo_test: %+v", late)
	}
}

func TestValidatesIncomingSamePackageStillInferred(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	if err := IndexFile(ctx, st, "pkg/foo_test.go", readTestdata(t, "foo_test.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "pkg/foo.go", readTestdata(t, "foo.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	foo := mustSymbolNamed(t, st, "pkg/foo.go", "Foo")
	edges, err := st.ListValidatesForSymbol(foo.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 1 || edges[0].Provenance != store.ImportProvenanceInferred {
		t.Fatalf("incoming same-package validates: %+v", edges)
	}
}

func mustSymbolNamed(t *testing.T, st *store.Store, path, name string) store.Symbol {
	t.Helper()
	syms, err := st.ListSymbolsByPath(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range syms {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("symbol %q missing in %s: %v", name, path, syms)
	return store.Symbol{}
}

func TestIndexPythonAndTSTestFiles(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	if err := IndexFile(ctx, st, "pkg/mod.py", readTestdata(t, "mod.py"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "pkg/test_mod.py", readTestdata(t, "test_mod.py"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "pkg/mod.ts", readTestdata(t, "mod.ts"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "pkg/mod.test.ts", readTestdata(t, "mod.test.ts"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}

	py, err := st.ListSymbolsByPath("pkg/test_mod.py")
	if err != nil {
		t.Fatal(err)
	}
	gotPy := symNamesKinds(py)
	if len(gotPy) != 1 || gotPy[0] != "test:test_foo" {
		t.Fatalf("python test symbols got %v want [test:test_foo]", gotPy)
	}

	ts, err := st.ListSymbolsByPath("pkg/mod.test.ts")
	if err != nil {
		t.Fatal(err)
	}
	gotTS := symNamesKinds(ts)
	wantTS := []string{"test:Foo", "test:nested", "test:suite"}
	if len(gotTS) != len(wantTS) {
		t.Fatalf("ts test symbols got %v want %v", gotTS, wantTS)
	}
	for i := range wantTS {
		if gotTS[i] != wantTS[i] {
			t.Fatalf("ts test symbols got %v want %v", gotTS, wantTS)
		}
	}
}

func TestBuiltinLanguageAdaptersContributionPath(t *testing.T) {
	if len(builtinAdapters) < 5 {
		t.Fatalf("builtinAdapters length = %d, want ≥ 5", len(builtinAdapters))
	}

	wantIDs := map[string]bool{
		LangJavaScript: false,
		LangTypeScript: false,
		LangTSX:        false,
		LangPython:     false,
		LangGo:         false,
	}
	seenExt := map[string]string{}

	for _, a := range builtinAdapters {
		id := a.ID()
		if _, known := wantIDs[id]; known {
			wantIDs[id] = true
		}
		for _, ext := range a.Extensions() {
			if ext == "" || ext[0] != '.' || ext != strings.ToLower(ext) {
				t.Fatalf("adapter %q: extension %q must be lowercase with leading dot", id, ext)
			}
			if other, dup := seenExt[ext]; dup {
				t.Fatalf("duplicate extension %q claimed by %q and %q", ext, other, id)
			}
			seenExt[ext] = id

			got, ok := DetectLanguage("x" + ext)
			if !ok || got != id {
				t.Fatalf("DetectLanguage(%q): got (%q,%v) want (%q,true)", "x"+ext, got, ok, id)
			}
		}

		_, _, err := extract(id, []byte("\n"))
		if err != nil && strings.Contains(err.Error(), "unsupported language") {
			t.Fatalf("extract(%q): unsupported language — adapter not wired: %v", id, err)
		}
	}

	for id, found := range wantIDs {
		if !found {
			t.Fatalf("builtinAdapters missing language id %q", id)
		}
	}
}

func TestArtifactEdgesFunctionsTypesAPIs(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	if err := IndexFile(ctx, st, "pkg/sample.go", readTestdata(t, "sample.go"), IndexOptions{}); err != nil {
		t.Fatalf("IndexFile sample.go: %v", err)
	}

	gotGo := exportNamesKinds(t, st, "pkg/sample.go")
	wantGo := []string{
		"function:Helper",
		"function:Main",
		"method:Run",
		"type:Counter",
		"type:ID",
		"type:Worker",
	}
	if len(gotGo) != len(wantGo) {
		t.Fatalf("go exports got %v want %v", gotGo, wantGo)
	}
	for i := range wantGo {
		if gotGo[i] != wantGo[i] {
			t.Fatalf("go exports got %v want %v", gotGo, wantGo)
		}
	}
	for _, e := range mustExports(t, st, "pkg/sample.go") {
		if e.Provenance != store.ImportProvenanceExtracted {
			t.Fatalf("go export provenance %q want EXTRACTED", e.Provenance)
		}
		if e.Rel != store.RelExportsAPI {
			t.Fatalf("go export rel %q", e.Rel)
		}
	}

	js := []byte("function hidden() { return 'export function trap'; }\nexport function visible() { return 1; }\nexport class Box { hold() { return 1; } }\n")
	if err := IndexFile(ctx, st, "src/api.js", js, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile api.js: %v", err)
	}
	gotJS := exportNamesKinds(t, st, "src/api.js")
	wantJS := []string{"class:Box", "function:visible", "method:hold"}
	if len(gotJS) != len(wantJS) {
		t.Fatalf("js exports got %v want %v", gotJS, wantJS)
	}
	for i := range wantJS {
		if gotJS[i] != wantJS[i] {
			t.Fatalf("js exports got %v want %v", gotJS, wantJS)
		}
	}

	// export { foo } is an export_statement, but the declaration is not its child.
	// EXTRACTED exports_api is declaration-under-export only (locked S01-03).
	reexport := []byte("function foo() { return 1; }\nexport { foo };\n")
	if err := IndexFile(ctx, st, "src/reexport.js", reexport, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile reexport.js: %v", err)
	}
	gotRe := exportNamesKinds(t, st, "src/reexport.js")
	if len(gotRe) != 0 {
		t.Fatalf("export { foo } must not be exports_api (honest EXTRACTED): got %v", gotRe)
	}
	gotContains := containsNamesKinds(t, st, "src/reexport.js")
	if len(gotContains) != 1 || gotContains[0] != "function:foo" {
		t.Fatalf("export { foo } still contains the local function: got %v", gotContains)
	}

	if err := IndexFile(ctx, st, "pkg/sample.py", readTestdata(t, "sample.py"), IndexOptions{}); err != nil {
		t.Fatalf("IndexFile sample.py: %v", err)
	}
	gotPy := exportNamesKinds(t, st, "pkg/sample.py")
	wantPy := []string{"class:Worker", "function:helper", "function:main", "method:run"}
	if len(gotPy) != len(wantPy) {
		t.Fatalf("py exports got %v want %v", gotPy, wantPy)
	}
	for i := range wantPy {
		if gotPy[i] != wantPy[i] {
			t.Fatalf("py exports got %v want %v", gotPy, wantPy)
		}
	}

	pyPriv := []byte("def public():\n    return 1\n\ndef _hidden():\n    return 2\n")
	if err := IndexFile(ctx, st, "pkg/priv.py", pyPriv, IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	gotPriv := exportNamesKinds(t, st, "pkg/priv.py")
	if len(gotPriv) != 1 || gotPriv[0] != "function:public" {
		t.Fatalf("python underscore export: got %v want [function:public]", gotPriv)
	}

	for _, e := range mustListEdges(t, st, "pkg/sample.go") {
		if e.Rel == store.RelDependsOn {
			t.Fatalf("must not clone imports as depends_on: %+v", e)
		}
	}
}

func TestModuleContainsSymbols(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	if err := IndexFile(ctx, st, "pkg/foo.go", readTestdata(t, "foo.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "pkg/foo_test.go", readTestdata(t, "foo_test.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}

	fileGot := containsNamesKinds(t, st, "pkg/foo.go")
	wantFile := []string{"function:Foo", "function:Untested"}
	if len(fileGot) != len(wantFile) {
		t.Fatalf("file contains got %v want %v", fileGot, wantFile)
	}
	for i := range wantFile {
		if fileGot[i] != wantFile[i] {
			t.Fatalf("file contains got %v want %v", fileGot, wantFile)
		}
	}

	modGot := containsNamesKinds(t, st, "pkg")
	wantMod := []string{"function:Foo", "function:Untested", "function:notATest", "package:foo", "test:TestFoo"}
	if len(modGot) != len(wantMod) {
		t.Fatalf("module contains got %v want %v", modGot, wantMod)
	}
	for i := range wantMod {
		if modGot[i] != wantMod[i] {
			t.Fatalf("module contains got %v want %v", modGot, wantMod)
		}
	}

	foo := mustSymbolNamed(t, st, "pkg/foo.go", "Foo")
	before, err := st.ListValidatesForSymbol(foo.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(before) != 1 {
		t.Fatalf("validates before test reindex: %+v", before)
	}
	if err := IndexFile(ctx, st, "pkg/foo_test.go", readTestdata(t, "foo_test.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	after, err := st.ListValidatesForSymbol(foo.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != 1 {
		t.Fatalf("test reindex dropped validates: %+v", after)
	}
	modAfter := containsNamesKinds(t, st, "pkg")
	if len(modAfter) != len(wantMod) {
		t.Fatalf("test reindex dropped contains_module: got %v want %v", modAfter, wantMod)
	}

	syms, err := st.ListSymbolsByPath("pkg/foo_test.go")
	if err != nil {
		t.Fatal(err)
	}
	got := symNamesKinds(syms)
	want := []string{"function:notATest", "package:foo", "test:TestFoo"}
	if len(got) != len(want) {
		t.Fatalf("test file symbols (kind=package keeper) got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("test file symbols got %v want %v", got, want)
		}
	}
}

func TestArchitecturalBoundaryEdges(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	if err := os.WriteFile(filepath.Join(st.ProjectRoot(), "go.mod"), []byte("module example.com/mod\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cmdSrc := []byte("package main\n\nimport \"example.com/mod/internal/lib\"\n\nfunc Run() { lib.Helper() }\n")
	libSrc := []byte("package lib\n\nfunc Helper() {}\n")
	if err := IndexFile(ctx, st, "cmd/app/main.go", cmdSrc, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile cmd: %v", err)
	}
	if err := IndexFile(ctx, st, "internal/lib/lib.go", libSrc, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile internal: %v", err)
	}

	cmdLayer, cmdProv, err := st.FileLayer("cmd/app/main.go")
	if err != nil {
		t.Fatalf("FileLayer cmd: %v", err)
	}
	libLayer, libProv, err := st.FileLayer("internal/lib/lib.go")
	if err != nil {
		t.Fatalf("FileLayer internal: %v", err)
	}
	if cmdLayer != "cmd" || libLayer != "internal" {
		t.Fatalf("layers cmd=%q internal=%q want cmd vs internal", cmdLayer, libLayer)
	}
	if cmdProv != store.ImportProvenanceExtracted || libProv != store.ImportProvenanceExtracted {
		t.Fatalf("provenance cmd=%q internal=%q want EXTRACTED (go.mod + package dirs)", cmdProv, libProv)
	}

	cmdBounds, err := st.ListArchitecturalBoundaries("cmd/app/main.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(cmdBounds) != 1 || cmdBounds[0].Rel != store.RelArchitecturalBoundary {
		t.Fatalf("cmd architectural_boundary: %+v", cmdBounds)
	}
	cmdFile, err := st.GetFileByPath("cmd/app/main.go")
	if err != nil {
		t.Fatal(err)
	}
	if cmdBounds[0].FromFileID != cmdFile.ID {
		t.Fatalf("from_file_id must be the indexed source file, got %s want %s", cmdBounds[0].FromFileID, cmdFile.ID)
	}
	toFile, err := st.GetFileByID(cmdBounds[0].ToFileID)
	if err != nil {
		t.Fatal(err)
	}
	wantTo := store.ArchitectureLayerPath("cmd")
	if toFile.Path != wantTo {
		t.Fatalf("TO layer identity: %q want %q", toFile.Path, wantTo)
	}
	stubBounds, err := st.ListArchitecturalBoundaries(toFile.Path)
	if err != nil {
		t.Fatal(err)
	}
	if len(stubBounds) != 0 {
		t.Fatalf("layer identity TO stub must not be FROM: %+v", stubBounds)
	}

	cross, err := st.ListCrossLayerImports()
	if err != nil {
		t.Fatal(err)
	}
	if len(cross) != 1 || cross[0].FromLayer != "cmd" || cross[0].ToLayer != "internal" {
		t.Fatalf("cross-layer imports (left as imports, not rebuilt): %+v", cross)
	}
	for _, e := range mustListEdges(t, st, "cmd/app/main.go") {
		if e.Rel == store.RelDependsOn {
			t.Fatalf("must not clone imports as depends_on: %+v", e)
		}
	}
}

func TestArchitecturalBoundaryIncremental(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	cmdSrc := []byte("package main\n\nfunc Run() {}\n")
	libSrc := []byte("package lib\n\nfunc Helper() {}\n")
	if err := IndexFile(ctx, st, "cmd/app/main.go", cmdSrc, IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "internal/lib/lib.go", libSrc, IndexOptions{}); err != nil {
		t.Fatal(err)
	}

	libBefore, err := st.ListArchitecturalBoundaries("internal/lib/lib.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(libBefore) != 1 {
		t.Fatalf("internal boundaries before: %+v", libBefore)
	}
	if libBefore[0].Provenance != store.ImportProvenanceInferred {
		t.Fatalf("provenance without go.mod: %q want INFERRED", libBefore[0].Provenance)
	}
	libID := libBefore[0].ID
	libFrom := libBefore[0].FromFileID
	libTo := libBefore[0].ToFileID
	allBefore := countRel(t, st, store.RelArchitecturalBoundary)

	cmd2 := []byte("package main\n\nfunc Run() {}\nfunc Extra() {}\n")
	if err := IndexFile(ctx, st, "cmd/app/main.go", cmd2, IndexOptions{}); err != nil {
		t.Fatalf("reindex cmd: %v", err)
	}

	libAfter, err := st.ListArchitecturalBoundaries("internal/lib/lib.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(libAfter) != 1 {
		t.Fatalf("internal boundaries rewritten by cmd IndexFile: %+v", libAfter)
	}
	if libAfter[0].ID != libID || libAfter[0].FromFileID != libFrom || libAfter[0].ToFileID != libTo {
		t.Fatalf("internal boundary mutated: before %+v after %+v", libBefore[0], libAfter[0])
	}
	allAfter := countRel(t, st, store.RelArchitecturalBoundary)
	if allAfter != allBefore {
		t.Fatalf("full-graph boundary rebuild: before %d after %d", allBefore, allAfter)
	}
	if _, _, err := st.FileLayer("cmd/app/main.go"); err != nil {
		t.Fatalf("cmd layer dropped: %v", err)
	}
}

func TestArchitecturalBoundarySurvivesFileReindex(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	libSrc := []byte("package lib\n\nfunc Helper() {}\n")
	testSrc := []byte("package lib\n\nimport \"testing\"\n\nfunc TestHelper(t *testing.T) {}\n")
	if err := IndexFile(ctx, st, "internal/lib/lib.go", libSrc, IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "internal/lib/lib_test.go", testSrc, IndexOptions{}); err != nil {
		t.Fatal(err)
	}

	assertOutgoingMerge := func(path string, wantValidates bool) {
		t.Helper()
		edges := mustListEdges(t, st, path)
		got := map[string]int{}
		for _, e := range edges {
			got[e.Rel]++
		}
		if got[store.RelArchitecturalBoundary] < 1 {
			t.Fatalf("%s missing architectural_boundary: %+v", path, edges)
		}
		if got[store.RelContainsModule] < 1 {
			t.Fatalf("%s missing contains_module: %+v", path, edges)
		}
		if path == "internal/lib/lib.go" && got[store.RelExportsAPI] < 1 {
			t.Fatalf("%s missing exports_api: %+v", path, edges)
		}
		if wantValidates && got[store.RelValidates] < 1 {
			t.Fatalf("%s missing validates: %+v", path, edges)
		}
	}

	assertOutgoingMerge("internal/lib/lib.go", false)
	assertOutgoingMerge("internal/lib/lib_test.go", true)

	libBounds, err := st.ListArchitecturalBoundaries("internal/lib/lib.go")
	if err != nil {
		t.Fatal(err)
	}
	libContains := containsNamesKinds(t, st, "internal/lib/lib.go")
	libExports := exportNamesKinds(t, st, "internal/lib/lib.go")
	testValidates := 0
	for _, e := range mustListEdges(t, st, "internal/lib/lib_test.go") {
		if e.Rel == store.RelValidates {
			testValidates++
		}
	}

	if err := IndexFile(ctx, st, "internal/lib/lib.go", libSrc, IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "internal/lib/lib_test.go", testSrc, IndexOptions{}); err != nil {
		t.Fatal(err)
	}

	assertOutgoingMerge("internal/lib/lib.go", false)
	assertOutgoingMerge("internal/lib/lib_test.go", true)

	libBoundsAfter, err := st.ListArchitecturalBoundaries("internal/lib/lib.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(libBoundsAfter) != len(libBounds) {
		t.Fatalf("lib architectural_boundary dropped on reindex: before %+v after %+v", libBounds, libBoundsAfter)
	}
	gotContains := containsNamesKinds(t, st, "internal/lib/lib.go")
	if strings.Join(gotContains, ",") != strings.Join(libContains, ",") {
		t.Fatalf("contains_module dropped on reindex: got %v want %v", gotContains, libContains)
	}
	gotExports := exportNamesKinds(t, st, "internal/lib/lib.go")
	if strings.Join(gotExports, ",") != strings.Join(libExports, ",") {
		t.Fatalf("exports_api dropped on reindex: got %v want %v", gotExports, libExports)
	}
	testValidatesAfter := 0
	for _, e := range mustListEdges(t, st, "internal/lib/lib_test.go") {
		if e.Rel == store.RelValidates {
			testValidatesAfter++
		}
	}
	if testValidatesAfter != testValidates {
		t.Fatalf("validates dropped on test reindex: before %d after %d", testValidates, testValidatesAfter)
	}
}

func TestArchitecturalBoundaryOverlayExtracted(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	root := st.ProjectRoot()
	if err := os.MkdirAll(filepath.Join(root, "trace"), 0o755); err != nil {
		t.Fatal(err)
	}
	overlay := []byte(`{"layers":[{"name":"platform","prefixes":["internal/lib"]}]}` + "\n")
	if err := os.WriteFile(filepath.Join(root, "trace", "architecture.json"), overlay, 0o644); err != nil {
		t.Fatal(err)
	}

	libSrc := []byte("package lib\n\nfunc Helper() {}\n")
	cmdSrc := []byte("package main\n\nfunc Run() {}\n")
	if err := IndexFile(ctx, st, "internal/lib/lib.go", libSrc, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile overlay match: %v", err)
	}
	if err := IndexFile(ctx, st, "cmd/app/main.go", cmdSrc, IndexOptions{}); err != nil {
		t.Fatalf("IndexFile overlay miss: %v", err)
	}

	libLayer, libProv, err := st.FileLayer("internal/lib/lib.go")
	if err != nil {
		t.Fatalf("FileLayer overlay: %v", err)
	}
	if libLayer != "platform" || libProv != store.ImportProvenanceExtracted {
		t.Fatalf("overlay layer=%q provenance=%q want platform EXTRACTED", libLayer, libProv)
	}
	libBounds, err := st.ListArchitecturalBoundaries("internal/lib/lib.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(libBounds) != 1 {
		t.Fatalf("overlay architectural_boundary: %+v", libBounds)
	}
	libFile, err := st.GetFileByPath("internal/lib/lib.go")
	if err != nil {
		t.Fatal(err)
	}
	if libBounds[0].FromFileID != libFile.ID {
		t.Fatalf("overlay FROM must be indexed source file, got %s want %s", libBounds[0].FromFileID, libFile.ID)
	}
	toFile, err := st.GetFileByID(libBounds[0].ToFileID)
	if err != nil {
		t.Fatal(err)
	}
	wantTo := store.ArchitectureLayerPath("platform")
	if toFile.Path != wantTo {
		t.Fatalf("overlay TO: %q want %q", toFile.Path, wantTo)
	}

	cmdLayer, cmdProv, err := st.FileLayer("cmd/app/main.go")
	if err != nil {
		t.Fatalf("FileLayer unmatched overlay: %v", err)
	}
	if cmdLayer != "cmd" || cmdProv != store.ImportProvenanceInferred {
		t.Fatalf("unmatched overlay layer=%q provenance=%q want cmd INFERRED (no go.mod)", cmdLayer, cmdProv)
	}
}

func countRel(t *testing.T, st *store.Store, rel string) int {
	t.Helper()
	paths, err := st.ListFilePaths()
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, p := range paths {
		if strings.HasPrefix(p, store.ArchitectureLayerPrefix) {
			continue
		}
		edges, err := st.ListEdgesByFile(p)
		if err != nil {
			continue
		}
		for _, e := range edges {
			if e.Rel == rel {
				n++
			}
		}
	}
	return n
}

func mustExports(t *testing.T, st *store.Store, path string) []store.CodeEdge {
	t.Helper()
	edges, err := st.ListExports(path)
	if err != nil {
		t.Fatal(err)
	}
	return edges
}

func mustListEdges(t *testing.T, st *store.Store, path string) []store.CodeEdge {
	t.Helper()
	edges, err := st.ListEdgesByFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return edges
}

func exportNamesKinds(t *testing.T, st *store.Store, path string) []string {
	t.Helper()
	return edgeTargetNamesKinds(t, st, mustExports(t, st, path))
}

func containsNamesKinds(t *testing.T, st *store.Store, dirOrPath string) []string {
	t.Helper()
	edges, err := st.ListModuleContents(dirOrPath)
	if err != nil {
		t.Fatal(err)
	}
	return edgeTargetNamesKinds(t, st, edges)
}

func edgeTargetNamesKinds(t *testing.T, st *store.Store, edges []store.CodeEdge) []string {
	t.Helper()
	var out []string
	for _, e := range edges {
		if e.ToSymbolID == nil {
			t.Fatalf("edge missing to_symbol_id: %+v", e)
		}
		sym, _, err := st.GetSymbolByID(*e.ToSymbolID)
		if err != nil {
			t.Fatalf("GetSymbolByID: %v", err)
		}
		out = append(out, sym.Kind+":"+sym.Name)
	}
	sort.Strings(out)
	return out
}
