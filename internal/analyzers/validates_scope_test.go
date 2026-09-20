package analyzers

import (
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/mrchatam/Trace/internal/store"
)

func TestCandidateTestPathsForIncomingValidatesBounded(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()

	decoySrc := "package other\n\nfunc Decoy() {}\n"
	decoyTest := "package other\n\nfunc TestDecoy(t *testing.T) { Decoy() }\n"
	if err := IndexFile(ctx, st, "other/decoy.go", []byte(decoySrc), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "other/decoy_test.go", []byte(decoyTest), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "pkg/foo.go", readTestdata(t, "foo.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "pkg/foo_test.go", readTestdata(t, "foo_test.go"), IndexOptions{}); err != nil {
		t.Fatal(err)
	}

	tf, err := st.GetFileByPath("pkg/foo.go")
	if err != nil {
		t.Fatal(err)
	}
	cands, err := candidateTestPathsForIncomingValidates(st, "pkg/foo.go", tf.ID)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, p := range cands {
		seen[p] = true
	}
	if !seen["pkg/foo_test.go"] {
		t.Fatalf("expected same-dir test in candidates: %v", cands)
	}
	if seen["other/decoy_test.go"] {
		t.Fatalf("decoy other-package test must not be a candidate: %v", cands)
	}
}

func TestIncomingValidatesSkipsUnrelatedPackages(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		dir := path.Join("decoy", fmt.Sprintf("pkg%d", i))
		src := "package p\n\nfunc X() {}\n"
		tst := "package p\n\nfunc TestX(t *testing.T) {}\n"
		if err := IndexFile(ctx, st, path.Join(dir, "x.go"), []byte(src), IndexOptions{}); err != nil {
			t.Fatal(err)
		}
		if err := IndexFile(ctx, st, path.Join(dir, "x_test.go"), []byte(tst), IndexOptions{}); err != nil {
			t.Fatal(err)
		}
	}

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
	if len(edges) != 1 {
		t.Fatalf("incoming validates with decoys present: %+v", edges)
	}
}

func TestCandidateIncludesJSTSTestsDir(t *testing.T) {
	st := openTemp(t)
	ctx := context.Background()
	mod := "export function Foo() { return 1 }\n"
	tst := "import { Foo } from '../mod';\ntest('Foo', () => { Foo(); });\n"
	if err := IndexFile(ctx, st, "ui/mod.ts", []byte(mod), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := IndexFile(ctx, st, "ui/__tests__/mod.test.ts", []byte(tst), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	tf, err := st.GetFileByPath("ui/mod.ts")
	if err != nil {
		t.Fatal(err)
	}
	cands, err := candidateTestPathsForIncomingValidates(st, "ui/mod.ts", tf.ID)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, p := range cands {
		if p == "ui/__tests__/mod.test.ts" {
			found = true
		}
	}
	if !found {
		t.Fatalf("__tests__ candidate missing: %v", cands)
	}
	if err := IndexFile(ctx, st, "ui/mod.ts", []byte(mod), IndexOptions{}); err != nil {
		t.Fatal(err)
	}
	tf2, err := st.GetFileByPath("ui/mod.ts")
	if err != nil {
		t.Fatal(err)
	}
	edges, err := st.ListEdgesByFile("ui/__tests__/mod.test.ts")
	if err != nil {
		t.Fatal(err)
	}
	validates := 0
	for _, e := range edges {
		if e.Rel == store.RelValidates && e.ToFileID == tf2.ID {
			validates++
		}
	}
	if validates < 1 {
		t.Fatalf("expected validates from __tests__ after target reindex, edges=%+v", edges)
	}
}
