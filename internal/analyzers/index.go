package analyzers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

// IndexOptions controls IndexFile / IndexFileAtRev.
type IndexOptions struct {
	// GitOID is passed through to UpsertFile when known (optional).
	GitOID *string
}

const binaryProbeBytes = 8 * 1024

// IndexFile indexes one path from content bytes (unit of file-local incremental update).
// It upserts the file stub, sets language, replaces symbols and imports for that path only,
// classifies test symbols, writes outgoing code_edges (validates + contains_module +
// exports_api + architectural_boundary) in one ReplaceFileEdges batch, and upserts incoming validates.
func IndexFile(ctx context.Context, st *store.Store, path string, content []byte, opts IndexOptions) error {
	_ = ctx
	if st == nil {
		return fmt.Errorf("analyzers: IndexFile: store is nil")
	}
	path = store.NormalizePath(path)

	lang, ok := DetectLanguage(path)
	if !ok {
		return unsupportedExt(path)
	}
	if looksBinary(content) {
		return binaryFile(path)
	}

	contentHash := sha256Hex(content)
	if _, err := st.UpsertFile(path, contentHash, opts.GitOID); err != nil {
		return fmt.Errorf("analyzers: upsert file: %w", err)
	}
	if err := st.SetFileLanguage(path, lang); err != nil {
		return fmt.Errorf("analyzers: set language: %w", err)
	}

	symbols, imports, err := extract(lang, content)
	if err != nil {
		return fmt.Errorf("analyzers: extract %s: %w", lang, err)
	}
	if symbols == nil {
		symbols = []store.Symbol{}
	}
	if imports == nil {
		imports = []store.Import{}
	}

	if isJSTSTestPath(path) {
		extra, err := extractJSTSTestCalls(lang, content)
		if err != nil {
			return fmt.Errorf("analyzers: extract jsts test calls: %w", err)
		}
		symbols = append(symbols, extra...)
	}
	symbols = classifyTestSymbols(path, symbols)
	if lang == LangGo && strings.HasSuffix(path, "_test.go") {
		if pkg := extractGoPackageName(content); pkg != "" {
			// Persist package clause so incoming validates (no source blob)
			// can honor package foo_test vs same-package inferred heuristic.
			symbols = append(symbols, store.Symbol{Name: pkg, Kind: "package", StartLine: 1, EndLine: 1})
		}
	}

	// Clear outgoing edges before symbol replace. Leftover-symbol DELETE would
	// SET NULL to_symbol_id on outgoing contains_module/exports_api and collide
	// idx_code_edges_unique. Incoming validates on stable ids are kept by
	// ReplaceFileSymbols upsert-first; leftover incoming is collapsed there.
	if err := st.ReplaceFileEdges(path, nil); err != nil {
		return fmt.Errorf("analyzers: clear edges: %w", err)
	}
	if err := st.ReplaceFileSymbols(path, symbols); err != nil {
		return fmt.Errorf("analyzers: replace symbols: %w", err)
	}
	if err := st.ReplaceFileImports(path, imports); err != nil {
		return fmt.Errorf("analyzers: replace imports: %w", err)
	}
	if err := indexCodeEdges(st, path, content, lang); err != nil {
		return err
	}
	return nil
}

// IndexFileAtRev loads path at rev via vcs.Repository.ShowFile, then IndexFile.
func IndexFileAtRev(ctx context.Context, st *store.Store, repo vcs.Repository, rev, path string, opts IndexOptions) error {
	if repo == nil {
		return fmt.Errorf("analyzers: IndexFileAtRev: repository is nil")
	}
	content, err := repo.ShowFile(ctx, rev, store.NormalizePath(path))
	if err != nil {
		return fmt.Errorf("analyzers: show file: %w", err)
	}
	if opts.GitOID == nil && rev != "" {
		oid := rev
		opts.GitOID = &oid
	}
	return IndexFile(ctx, st, path, content, opts)
}

func sha256Hex(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func looksBinary(content []byte) bool {
	n := len(content)
	if n > binaryProbeBytes {
		n = binaryProbeBytes
	}
	for i := 0; i < n; i++ {
		if content[i] == 0 {
			return true
		}
	}
	return false
}

func extract(lang string, content []byte) ([]store.Symbol, []store.Import, error) {
	a, ok := adapterByID(lang)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported language %q", lang)
	}
	return a.Extract(content)
}
