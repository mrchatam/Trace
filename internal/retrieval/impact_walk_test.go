package retrieval_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func TestImpactWalkMultiSeedExcludeSeeds(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	a, err := st.UpsertFile("a.go", "ha", nil)
	if err != nil {
		t.Fatalf("UpsertFile a: %v", err)
	}
	b, err := st.UpsertFile("b.go", "hb", nil)
	if err != nil {
		t.Fatalf("UpsertFile b: %v", err)
	}
	c, err := st.UpsertFile("c.go", "hc", nil)
	if err != nil {
		t.Fatalf("UpsertFile c: %v", err)
	}
	d, err := st.UpsertFile("d.go", "hd", nil)
	if err != nil {
		t.Fatalf("UpsertFile d: %v", err)
	}
	// C imports A; D imports B — shared walk discovers both importers once.
	if err := st.ReplaceFileImports("c.go", []store.Import{{ImportedPath: "a.go"}}); err != nil {
		t.Fatalf("imports c: %v", err)
	}
	if err := st.ReplaceFileImports("d.go", []store.Import{{ImportedPath: "b.go"}}); err != nil {
		t.Fatalf("imports d: %v", err)
	}

	res, err := eng.ImpactWalk(ctx, []retrieval.ImpactSeed{
		{EntityType: "file", EntityID: a.ID},
		{EntityType: "file", EntityID: b.ID},
	}, 2)
	if err != nil {
		t.Fatalf("ImpactWalk: %v", err)
	}
	if len(res.Seeds) != 2 {
		t.Fatalf("seeds: %+v", res.Seeds)
	}
	ids := map[string]int{}
	for _, h := range res.Blast {
		if h.EntityID == a.ID || h.EntityID == b.ID {
			t.Fatalf("seed must be excluded from blast: %+v", h)
		}
		ids[h.EntityID]++
	}
	if ids[c.ID] != 1 || ids[d.ID] != 1 {
		t.Fatalf("want C and D once each; blast=%+v counts=%v", res.Blast, ids)
	}
}

func TestImpactWalkContainsAsymmetryNoSiblings(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	f, err := st.UpsertFile("pkg/mod.go", "hm", nil)
	if err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if err := st.ReplaceFileSymbols("pkg/mod.go", []store.Symbol{
		{Name: "Alpha", Kind: "function", StartLine: 1, EndLine: 5},
		{Name: "Beta", Kind: "function", StartLine: 6, EndLine: 10},
	}); err != nil {
		t.Fatalf("symbols: %v", err)
	}
	syms, err := st.ListSymbolsByPath("pkg/mod.go")
	if err != nil || len(syms) != 2 {
		t.Fatalf("ListSymbolsByPath: %v %+v", err, syms)
	}
	var seedSym, sibling store.Symbol
	for _, s := range syms {
		if s.Name == "Alpha" {
			seedSym = s
		} else {
			sibling = s
		}
	}
	imp, err := st.UpsertFile("pkg/user.go", "hu", nil)
	if err != nil {
		t.Fatalf("UpsertFile user: %v", err)
	}
	if err := st.ReplaceFileImports("pkg/user.go", []store.Import{
		{ImportedPath: "pkg/mod.go", Provenance: store.ImportProvenanceExtracted},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	res, err := eng.ImpactWalk(ctx, []retrieval.ImpactSeed{
		{EntityType: "symbol", EntityID: seedSym.ID},
	}, 2)
	if err != nil {
		t.Fatalf("ImpactWalk: %v", err)
	}
	sawFile, sawImporter, sawSibling := false, false, false
	for _, h := range res.Blast {
		if h.EntityID == seedSym.ID {
			t.Fatalf("seed symbol in blast: %+v", h)
		}
		if h.EntityID == sibling.ID {
			sawSibling = true
		}
		if h.EntityID == f.ID {
			sawFile = true
		}
		if h.EntityID == imp.ID {
			sawImporter = true
		}
	}
	if !sawFile {
		t.Fatalf("expected containing file in blast: %+v", res.Blast)
	}
	if !sawImporter {
		t.Fatalf("expected importer in blast: %+v", res.Blast)
	}
	if sawSibling {
		t.Fatalf("sibling symbol must not appear via contains climb: %+v", res.Blast)
	}
}

func TestImpactWalkIncomingImportHop(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	a, err := st.UpsertFile("lib/a.go", "ha", nil)
	if err != nil {
		t.Fatalf("UpsertFile a: %v", err)
	}
	b, err := st.UpsertFile("lib/b.go", "hb", nil)
	if err != nil {
		t.Fatalf("UpsertFile b: %v", err)
	}
	if err := st.ReplaceFileImports("lib/b.go", []store.Import{
		{ImportedPath: "lib/a.go", Provenance: store.ImportProvenanceInferred},
	}); err != nil {
		t.Fatalf("imports: %v", err)
	}

	res, err := eng.ImpactWalk(ctx, []retrieval.ImpactSeed{
		{EntityType: "file", EntityID: a.ID},
	}, 1)
	if err != nil {
		t.Fatalf("ImpactWalk: %v", err)
	}
	var found *retrieval.BlastHit
	for i := range res.Blast {
		if res.Blast[i].EntityID == b.ID {
			found = &res.Blast[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("expected importer B in blast: %+v", res.Blast)
	}
	if found.Hop < 1 {
		t.Fatalf("importer hop: %+v", found)
	}
	if found.EdgeProvenance != store.ImportProvenanceInferred {
		t.Fatalf("edge_provenance: got %q want INFERRED", found.EdgeProvenance)
	}
}

func TestImpactWalkLoudTruncation(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	seed, err := st.UpsertFile("core.go", "hc", nil)
	if err != nil {
		t.Fatalf("UpsertFile seed: %v", err)
	}
	const n = 70
	for i := 0; i < n; i++ {
		path := fmt.Sprintf("imp/%03d.go", i)
		if _, err := st.UpsertFile(path, fmt.Sprintf("h%d", i), nil); err != nil {
			t.Fatalf("UpsertFile %s: %v", path, err)
		}
		if err := st.ReplaceFileImports(path, []store.Import{{ImportedPath: "core.go"}}); err != nil {
			t.Fatalf("imports %s: %v", path, err)
		}
	}

	res, err := eng.ImpactWalk(ctx, []retrieval.ImpactSeed{
		{EntityType: "file", EntityID: seed.ID},
	}, 1)
	if err != nil {
		t.Fatalf("ImpactWalk: %v", err)
	}
	if !res.Truncated {
		t.Fatalf("expected truncated: total=%d kept=%d", res.BlastTotal, res.BlastKept)
	}
	if res.BlastKept != retrieval.MaxImpactBlast {
		t.Fatalf("blast_kept=%d want %d", res.BlastKept, retrieval.MaxImpactBlast)
	}
	if res.BlastTotal <= res.BlastKept {
		t.Fatalf("blast_total=%d should exceed kept=%d", res.BlastTotal, res.BlastKept)
	}
	if len(res.Blast) != res.BlastKept {
		t.Fatalf("len(blast)=%d kept=%d", len(res.Blast), res.BlastKept)
	}
}

func TestImpactWalkHopRiskIncreases(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	// A ← B ← C (incoming chain); seed A → B hop1, C hop2
	a, err := st.UpsertFile("chain/a.go", "ha", nil)
	if err != nil {
		t.Fatalf("UpsertFile a: %v", err)
	}
	b, err := st.UpsertFile("chain/b.go", "hb", nil)
	if err != nil {
		t.Fatalf("UpsertFile b: %v", err)
	}
	c, err := st.UpsertFile("chain/c.go", "hc", nil)
	if err != nil {
		t.Fatalf("UpsertFile c: %v", err)
	}
	if err := st.ReplaceFileImports("chain/b.go", []store.Import{{ImportedPath: "chain/a.go"}}); err != nil {
		t.Fatalf("imports b: %v", err)
	}
	if err := st.ReplaceFileImports("chain/c.go", []store.Import{{ImportedPath: "chain/b.go"}}); err != nil {
		t.Fatalf("imports c: %v", err)
	}

	res, err := eng.ImpactWalk(ctx, []retrieval.ImpactSeed{
		{EntityType: "file", EntityID: a.ID},
	}, 2)
	if err != nil {
		t.Fatalf("ImpactWalk: %v", err)
	}
	var hop1, hop2 *retrieval.BlastHit
	for i := range res.Blast {
		h := &res.Blast[i]
		if h.EntityID == b.ID {
			hop1 = h
		}
		if h.EntityID == c.ID {
			hop2 = h
		}
	}
	if hop1 == nil || hop2 == nil {
		t.Fatalf("need B and C: %+v", res.Blast)
	}
	if hop1.Hop != 1 || hop2.Hop != 2 {
		t.Fatalf("hops: B=%d C=%d", hop1.Hop, hop2.Hop)
	}
	if hop1.HopRisk != float64(hop1.Hop) {
		t.Fatalf("hop_risk B: %v", hop1.HopRisk)
	}
	if hop2.HopRisk < hop1.HopRisk {
		t.Fatalf("hop_risk must be non-decreasing: hop1=%v hop2=%v", hop1.HopRisk, hop2.HopRisk)
	}
}

func TestImpactWalkIncludesAffectedTests(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	lib, err := st.UpsertFile("pkg/foo.go", "hlib", nil)
	if err != nil {
		t.Fatalf("UpsertFile lib: %v", err)
	}
	testFile, err := st.UpsertFile("pkg/foo_test.go", "htest", nil)
	if err != nil {
		t.Fatalf("UpsertFile test: %v", err)
	}
	if err := st.ReplaceFileSymbols("pkg/foo.go", []store.Symbol{
		{Name: "Foo", Kind: "function", StartLine: 1, EndLine: 5},
	}); err != nil {
		t.Fatalf("symbols foo: %v", err)
	}
	if err := st.ReplaceFileSymbols("pkg/foo_test.go", []store.Symbol{
		{Name: "TestFoo", Kind: "test", StartLine: 3, EndLine: 10},
	}); err != nil {
		t.Fatalf("symbols test: %v", err)
	}
	prod, err := st.ListSymbolsByPath("pkg/foo.go")
	if err != nil || len(prod) != 1 {
		t.Fatalf("prod symbols: %v %v", prod, err)
	}
	tests, err := st.ListSymbolsByPath("pkg/foo_test.go")
	if err != nil || len(tests) != 1 {
		t.Fatalf("test symbols: %v %v", tests, err)
	}
	fromID := tests[0].ID
	toID := prod[0].ID
	if err := st.ReplaceFileEdges("pkg/foo_test.go", []store.CodeEdge{
		{
			FromFileID: testFile.ID, FromSymbolID: &fromID,
			ToFileID: lib.ID, ToSymbolID: &toID,
			Rel: store.RelValidates, Provenance: store.ImportProvenanceInferred,
		},
	}); err != nil {
		t.Fatalf("validates edge: %v", err)
	}
	_ = testFile

	res, err := eng.ImpactWalk(ctx, []retrieval.ImpactSeed{
		{EntityType: "symbol", EntityID: toID},
	}, 1)
	if err != nil {
		t.Fatalf("ImpactWalk: %v", err)
	}

	var testHit *retrieval.BlastHit
	for i := range res.Blast {
		if res.Blast[i].EntityID == fromID {
			testHit = &res.Blast[i]
			break
		}
	}
	if testHit == nil {
		t.Fatalf("TestFoo missing from blast: %+v", res.Blast)
	}
	if testHit.EntityType != "symbol" {
		t.Fatalf("test blast entity_type=%q want symbol", testHit.EntityType)
	}
	if testHit.EdgeProvenance != store.ImportProvenanceInferred {
		t.Fatalf("edge_provenance=%q want INFERRED", testHit.EdgeProvenance)
	}
	if len(res.AffectedTests) != 1 {
		t.Fatalf("affected_tests=%+v want one TestFoo", res.AffectedTests)
	}
	if res.AffectedTests[0].EntityID != fromID {
		t.Fatalf("affected_tests[0]=%+v want TestFoo id %s", res.AffectedTests[0], fromID)
	}
}

func TestImpactWalkDepthStillCapped(t *testing.T) {
	eng, st, _ := openEngine(t)
	ctx := context.Background()

	f, err := st.UpsertFile("x.go", "hx", nil)
	if err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	_, err = eng.ImpactWalk(ctx, []retrieval.ImpactSeed{
		{EntityType: "file", EntityID: f.ID},
	}, 3)
	if err == nil {
		t.Fatal("expected depth 3 to fail closed")
	}
}
