// Package analyzers extracts File + minimal Symbol/Import edges via tree-sitter.
//
// Supported languages (DR-ANLANG): JavaScript, TypeScript/TSX, Python, and Go.
// Each language is a LanguageAdapter registered in a compile-time static table
// (LanguageAdapterAPIVersion). Contributor steps live in docs/ANALYZER_CONTRIBUTION.md.
//
// Persistence goes only through the store UpsertFile / ReplaceFileSymbols /
// ReplaceFileImports / ReplaceFileEdges APIs (plus SetFileLanguage). Source bodies
// are never stored; only a SHA-256 content hash is persisted (G1).
//
// Incremental updates are file-local (DR-INCREMENTAL / P0-X): reindexing path A
// replaces symbols, imports, and outgoing code_edges for A only. Incoming
// validates from already-indexed tests may be upserted onto a newly indexed
// target without deleting other files' outgoing edges. There is no cascade
// reindex of importers or dependents in this package, and no full-project rebuild path.
//
// CGO is required here for official tree-sitter bindings. Other Trace packages
// (store, vcs, gitcli) stay CGO-free.
package analyzers
