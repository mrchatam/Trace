package store

import "testing"

func TestReplaceFileEdgesIsFileLocal(t *testing.T) {
	s, _ := openTempStore(t)

	fa, err := s.UpsertFile("a.go", "hash-a", nil)
	if err != nil {
		t.Fatalf("UpsertFile a: %v", err)
	}
	fb, err := s.UpsertFile("b.go", "hash-b", nil)
	if err != nil {
		t.Fatalf("UpsertFile b: %v", err)
	}
	ft, err := s.UpsertFile("t.go", "hash-t", nil)
	if err != nil {
		t.Fatalf("UpsertFile t: %v", err)
	}

	if err := s.ReplaceFileEdges("a.go", []CodeEdge{{
		ToFileID:   ft.ID,
		Rel:        RelValidates,
		Provenance: ImportProvenanceExtracted,
	}}); err != nil {
		t.Fatalf("ReplaceFileEdges a: %v", err)
	}
	if err := s.ReplaceFileEdges("b.go", []CodeEdge{{
		ToFileID:   ft.ID,
		Rel:        RelValidates,
		Provenance: ImportProvenanceInferred,
	}}); err != nil {
		t.Fatalf("ReplaceFileEdges b: %v", err)
	}

	bBefore, err := s.ListEdgesByFile("b.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(bBefore) != 1 || bBefore[0].ToFileID != ft.ID || bBefore[0].Provenance != ImportProvenanceInferred {
		t.Fatalf("b edges before: %+v", bBefore)
	}
	bID := bBefore[0].ID

	if err := s.ReplaceFileEdges("a.go", []CodeEdge{{
		ToFileID:   fb.ID,
		Rel:        RelValidates,
		Provenance: ImportProvenanceExtracted,
	}}); err != nil {
		t.Fatalf("ReplaceFileEdges a again: %v", err)
	}

	aAfter, err := s.ListEdgesByFile("a.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(aAfter) != 1 || aAfter[0].ToFileID != fb.ID {
		t.Fatalf("a edges after replace: %+v", aAfter)
	}

	bAfter, err := s.ListEdgesByFile("b.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(bAfter) != 1 {
		t.Fatalf("b outgoing deleted by a reindex: %+v", bAfter)
	}
	if bAfter[0].ID != bID || bAfter[0].ToFileID != ft.ID || bAfter[0].FromFileID != fb.ID {
		t.Fatalf("b edge mutated: before id=%s %+v after %+v", bID, bBefore[0], bAfter[0])
	}

	_ = fa
}

func TestUpsertFilePairEdgesDoesNotDeleteOtherOutgoing(t *testing.T) {
	s, _ := openTempStore(t)
	ft, err := s.UpsertFile("t_test.go", "ht", nil)
	if err != nil {
		t.Fatal(err)
	}
	fa, err := s.UpsertFile("a.go", "ha", nil)
	if err != nil {
		t.Fatal(err)
	}
	fb, err := s.UpsertFile("b.go", "hb", nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = ft

	if err := s.ReplaceFileEdges("t_test.go", []CodeEdge{
		{ToFileID: fa.ID, Rel: RelValidates, Provenance: ImportProvenanceExtracted},
		{ToFileID: fb.ID, Rel: RelValidates, Provenance: ImportProvenanceExtracted},
	}); err != nil {
		t.Fatal(err)
	}

	if err := s.UpsertFilePairEdges("t_test.go", "a.go", RelValidates, []CodeEdge{
		{Rel: RelValidates, Provenance: ImportProvenanceInferred},
	}); err != nil {
		t.Fatal(err)
	}

	edges, err := s.ListEdgesByFile("t_test.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 2 {
		t.Fatalf("want 2 outgoing, got %+v", edges)
	}
	var sawA, sawB bool
	for _, e := range edges {
		switch e.ToFileID {
		case fa.ID:
			sawA = true
			if e.Provenance != ImportProvenanceInferred {
				t.Fatalf("pair upsert provenance: %q", e.Provenance)
			}
		case fb.ID:
			sawB = true
			if e.Provenance != ImportProvenanceExtracted {
				t.Fatalf("other outgoing provenance mutated: %q", e.Provenance)
			}
		}
	}
	if !sawA || !sawB {
		t.Fatalf("missing pair/other edges: %+v", edges)
	}
}

func TestListExportsAndModuleContents(t *testing.T) {
	s, _ := openTempStore(t)
	fa, err := s.UpsertFile("pkg/a.go", "ha", nil)
	if err != nil {
		t.Fatal(err)
	}
	fb, err := s.UpsertFile("pkg/b.go", "hb", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.ReplaceFileSymbols("pkg/a.go", []Symbol{{Name: "A", Kind: "function", StartLine: 1, EndLine: 1}}); err != nil {
		t.Fatal(err)
	}
	if err := s.ReplaceFileSymbols("pkg/b.go", []Symbol{{Name: "B", Kind: "function", StartLine: 1, EndLine: 1}}); err != nil {
		t.Fatal(err)
	}
	sa, err := s.ListSymbolsByPath("pkg/a.go")
	if err != nil || len(sa) != 1 {
		t.Fatalf("symbols a: %v %v", sa, err)
	}
	sb, err := s.ListSymbolsByPath("pkg/b.go")
	if err != nil || len(sb) != 1 {
		t.Fatalf("symbols b: %v %v", sb, err)
	}
	idA, idB := sa[0].ID, sb[0].ID

	if err := s.ReplaceFileEdges("pkg/a.go", []CodeEdge{
		{ToFileID: fa.ID, ToSymbolID: &idA, Rel: RelContainsModule, Provenance: ImportProvenanceExtracted},
		{ToFileID: fa.ID, ToSymbolID: &idA, Rel: RelExportsAPI, Provenance: ImportProvenanceExtracted},
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.ReplaceFileEdges("pkg/b.go", []CodeEdge{
		{ToFileID: fb.ID, ToSymbolID: &idB, Rel: RelContainsModule, Provenance: ImportProvenanceExtracted},
	}); err != nil {
		t.Fatal(err)
	}

	exports, err := s.ListExports("pkg/a.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(exports) != 1 || exports[0].Rel != RelExportsAPI || exports[0].ToSymbolID == nil || *exports[0].ToSymbolID != idA {
		t.Fatalf("ListExports: %+v", exports)
	}

	fileContents, err := s.ListModuleContents("pkg/a.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(fileContents) != 1 || fileContents[0].Rel != RelContainsModule || fileContents[0].ToSymbolID == nil || *fileContents[0].ToSymbolID != idA {
		t.Fatalf("ListModuleContents file: %+v", fileContents)
	}

	dirContents, err := s.ListModuleContents("pkg")
	if err != nil {
		t.Fatal(err)
	}
	if len(dirContents) != 2 {
		t.Fatalf("ListModuleContents dir: %+v", dirContents)
	}
	seen := map[string]bool{}
	for _, e := range dirContents {
		if e.Rel != RelContainsModule || e.ToSymbolID == nil {
			t.Fatalf("dir edge: %+v", e)
		}
		seen[*e.ToSymbolID] = true
	}
	if !seen[idA] || !seen[idB] {
		t.Fatalf("dir missing symbols: %+v", dirContents)
	}

	nested, err := s.UpsertFile("pkg/sub/c.go", "hc", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.ReplaceFileSymbols("pkg/sub/c.go", []Symbol{{Name: "C", Kind: "function", StartLine: 1, EndLine: 1}}); err != nil {
		t.Fatal(err)
	}
	sc, err := s.ListSymbolsByPath("pkg/sub/c.go")
	if err != nil || len(sc) != 1 {
		t.Fatalf("symbols c: %v %v", sc, err)
	}
	idC := sc[0].ID
	if err := s.ReplaceFileEdges("pkg/sub/c.go", []CodeEdge{
		{ToFileID: nested.ID, ToSymbolID: &idC, Rel: RelContainsModule, Provenance: ImportProvenanceExtracted},
	}); err != nil {
		t.Fatal(err)
	}
	dirAgain, err := s.ListModuleContents("pkg")
	if err != nil {
		t.Fatal(err)
	}
	if len(dirAgain) != 2 {
		t.Fatalf("nested dir leaked into pkg: %+v", dirAgain)
	}
}

func TestReplaceFileSymbolsKeepsIncomingValidatesOnStableIDs(t *testing.T) {
	s, _ := openTempStore(t)
	lib, err := s.UpsertFile("lib.js", "h1", nil)
	if err != nil {
		t.Fatal(err)
	}
	test, err := s.UpsertFile("lib.test.js", "ht", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.ReplaceFileSymbols("lib.js", []Symbol{
		{Name: "alpha", Kind: "function", StartLine: 1, EndLine: 1},
		{Name: "beta", Kind: "function", StartLine: 2, EndLine: 2},
	}); err != nil {
		t.Fatal(err)
	}
	syms, err := s.ListSymbolsByPath("lib.js")
	if err != nil || len(syms) != 2 {
		t.Fatalf("symbols: %v %v", syms, err)
	}
	idA, idB := syms[0].ID, syms[1].ID
	if err := s.ReplaceFileEdges("lib.test.js", []CodeEdge{
		{FromFileID: test.ID, ToFileID: lib.ID, ToSymbolID: &idA, Rel: RelValidates, Provenance: ImportProvenanceExtracted},
		{FromFileID: test.ID, ToFileID: lib.ID, ToSymbolID: &idB, Rel: RelValidates, Provenance: ImportProvenanceExtracted},
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.ReplaceFileSymbols("lib.js", []Symbol{
		{Name: "alpha", Kind: "function", StartLine: 1, EndLine: 1},
		{Name: "beta", Kind: "function", StartLine: 2, EndLine: 2},
	}); err != nil {
		t.Fatalf("stable reindex unique-collide: %v", err)
	}
	va, err := s.ListValidatesForSymbol(idA)
	if err != nil {
		t.Fatal(err)
	}
	vb, err := s.ListValidatesForSymbol(idB)
	if err != nil {
		t.Fatal(err)
	}
	if len(va) != 1 || va[0].ToSymbolID == nil || *va[0].ToSymbolID != idA {
		t.Fatalf("alpha incoming: %+v", va)
	}
	if len(vb) != 1 || vb[0].ToSymbolID == nil || *vb[0].ToSymbolID != idB {
		t.Fatalf("beta incoming: %+v", vb)
	}
}

func TestReplaceFileSymbolsCollapsesIncomingWhenSymbolDropped(t *testing.T) {
	s, _ := openTempStore(t)
	lib, err := s.UpsertFile("lib.js", "h1", nil)
	if err != nil {
		t.Fatal(err)
	}
	test, err := s.UpsertFile("lib.test.js", "ht", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.ReplaceFileSymbols("lib.js", []Symbol{
		{Name: "alpha", Kind: "function", StartLine: 1, EndLine: 1},
		{Name: "beta", Kind: "function", StartLine: 2, EndLine: 2},
	}); err != nil {
		t.Fatal(err)
	}
	syms, err := s.ListSymbolsByPath("lib.js")
	if err != nil || len(syms) != 2 {
		t.Fatalf("symbols: %v %v", syms, err)
	}
	var idA, idB string
	for _, sym := range syms {
		switch sym.Name {
		case "alpha":
			idA = sym.ID
		case "beta":
			idB = sym.ID
		}
	}
	if err := s.ReplaceFileEdges("lib.test.js", []CodeEdge{
		{FromFileID: test.ID, ToFileID: lib.ID, ToSymbolID: &idA, Rel: RelValidates, Provenance: ImportProvenanceExtracted},
		{FromFileID: test.ID, ToFileID: lib.ID, ToSymbolID: &idB, Rel: RelValidates, Provenance: ImportProvenanceExtracted},
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.ReplaceFileSymbols("lib.js", []Symbol{
		{Name: "gamma", Kind: "function", StartLine: 3, EndLine: 3},
	}); err != nil {
		t.Fatalf("drop alpha+beta unique-collide: %v", err)
	}
	edges, err := s.ListEdgesByFile("lib.test.js")
	if err != nil {
		t.Fatal(err)
	}
	if len(edges) != 1 {
		t.Fatalf("collapsed leftover incoming: %+v", edges)
	}
	if edges[0].ToSymbolID != nil {
		t.Fatalf("leftover incoming must SET NULL: %+v", edges[0])
	}
	if edges[0].ToFileID != lib.ID || edges[0].Rel != RelValidates {
		t.Fatalf("collapsed edge: %+v", edges[0])
	}
}
