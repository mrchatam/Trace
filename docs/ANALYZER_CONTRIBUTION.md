# Analyzer contribution

How to add a language to Trace’s indexer without rewriting `IndexFile` or introducing a plugin megastore.

## Contract version

- Package: `internal/analyzers`
- Const: `LanguageAdapterAPIVersion` (currently **1**)
- Bump the const only when the `LanguageAdapter` interface or dispatch semantics break contributors

## `LanguageAdapter` (v1)

```go
type LanguageAdapter interface {
    ID() string
    Extensions() []string // lowercase, leading dot, e.g. ".go"
    Extract(content []byte) ([]store.Symbol, []store.Import, error)
}
```

- `ID()` must match the persisted `files.language` value (see existing `Lang*` consts).
- `Extensions()` must not overlap other adapters’ claims.
- `Extract` may return empty slices for minimal/invalid content; missing adapters fail with `unsupported language`.

## Built-in registration

Adapters live in a **compile-time static table** (`builtinAdapters` in `language_adapter.go`).

- `DetectLanguage` resolves extension → adapter `ID()`
- Unexported `extract(lang, …)` dispatches adapter `ID()` → `Extract`
- `IndexFile` / `IndexFileAtRev` orchestration stays Upsert → SetFileLanguage → extract → Replace*

There is **no** public mutable `Register` API and **no** dynamic `.so` / plugin loader.

## Add a language (checklist)

Tier-2 languages are **deferred** until a human-promoted board row. See [INDEX_LANG_POLICY.md](INDEX_LANG_POLICY.md) for tier tables and the contribution gate.

1. Add `extract_<lang>.go` using an official tree-sitter grammar (`bindings/go`). Avoid `*_js.go` GOOS traps.
2. Implement a small `LanguageAdapter` (struct wrapping the extractor).
3. Append the adapter to `builtinAdapters`.
4. Add a golden `IndexFile` test under `internal/analyzers` (+ `testdata/` sample).
5. Update this doc’s supported-language list and [INDEX_LANG_POLICY.md](INDEX_LANG_POLICY.md) Tier-1 table if promoted from Tier-2.
6. Run: `CGO_ENABLED=1 go test ./internal/analyzers/... -count=1`

Do **not** change `IndexFile` body orchestration for a new language.

## Forbidden

- Dynamic `.so` / plugin megastore / hosted marketplace loaders
- Public mutable global Register for product paths
- Store migrations for language version metadata (version is a Go const)
- Full-rebuild-on-any-change indexer architecture
- New CLI/MCP surface solely to register languages

## Supported today (v1 table)

| ID | Extensions |
|----|------------|
| `javascript` | `.js` `.jsx` `.mjs` `.cjs` |
| `typescript` | `.ts` |
| `tsx` | `.tsx` |
| `python` | `.py` |
| `go` | `.go` |

## Import edge provenance (Law 5)

Analyzers emit **`EXTRACTED`** for concrete import edges and **`AMBIGUOUS`** only where the extractor cannot pin a single target (e.g. Python star/wildcard). They do **not** emit **`INFERRED`**.

`INFERRED` is reserved for store/fixture paths that prove Law 5 (inferred ≠ extracted) on Expand/Why/compiler surfaces until a future call-graph phase. There is **no** `trace` CLI (or MCP) command to set import provenance — write validation rejects unknown strings; empty defaults to `EXTRACTED`.
