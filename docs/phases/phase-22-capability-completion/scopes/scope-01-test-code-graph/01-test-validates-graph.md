# P22-S01-01 — Implement: test discovery + `validates` edges

## Metadata
- id: P22-S01-01
- todo_ids: [P22-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Index **tests** as first-class graph artifacts and record **`validates`** edges to the code they exercise. Closes **C03**. Unlocks C07 (S01-07). Board edits: **status + notes only**.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- Live: `internal/analyzers/index.go` (`IndexFile`), `extract.go`, `extract_go.go`, `extract_py.go`, `extract_javascript.go` (not `extract_js.go`), `extract_ts.go`, `language_adapter.go`
- Live store: `internal/store/file_graph.go` (`UpsertFile`, `ReplaceFileSymbols`, `ReplaceFileImports` only — no edges yet)
- Schema max **before** this row: **021**. This row adds **022**.
- Compat ceiling **before** this row: **21** (`evals/compat/compat_test.go`, `internal/store/production_hardening_test.go`, `internal/store/deliberation_test.go`)

## Locked defaults

| Item | Value |
|------|-------|
| Migration | **`internal/store/schema/022_code_relationships.sql`** only; **forbid 023+** |
| Embed | `//go:embed schema/*.sql` in `internal/store/migrate.go` — **adding the SQL file is enough**; do not edit migrate.go unless version parse fails |
| Compat | Ceiling **22** |
| Table | `code_edges`: `id`, `from_file_id`, `from_symbol_id` (nullable), `to_file_id`, `to_symbol_id` (nullable), `rel`, `provenance` |
| `rel` CHECK | `validates` \| `contains_module` \| `exports_api` \| `architectural_boundary` \| `depends_on` (this row **writes** `validates` only; other values exist so S01-03/05 need no 023) |
| Provenance | `EXTRACTED` \| `INFERRED` \| `AMBIGUOUS` — same honesty as imports (`validateImportProvenance` / shared helper) |
| FKs | file ids `ON DELETE CASCADE`; symbol ids `ON DELETE SET NULL` (ReplaceFileSymbols must not wipe other files' incoming `validates`) |
| Unique | `(from_file_id, IFNULL(from_symbol_id,''), to_file_id, IFNULL(to_symbol_id,''), rel)` |
| Incremental | `ReplaceFileEdges(path)` deletes **outgoing** (`from_file_id` of that path) only — mirror `ReplaceFileSymbols` |
| IndexFile | After symbols+imports: classify tests, then `ReplaceFileEdges`. `LanguageAdapter.Extract` signature **unchanged**; `LanguageAdapterAPIVersion` stays **1** |
| Seed | Do **not** add `code_edges` / files / symbols to `SeedDocument` (D-22-06). Keeper `TestSeedExportOmitsDeniedSurfaces` |
| Blobs | No source bodies (Law 1). `TestNoSourceContentColumns` must stay green |
| MCP/CLI query | **S05-03** — do not add `trace tests` or a new MCP tool this row |

## Locked test heuristics (must match live queries)

Classify in **`IndexFile`** (it has `path`). Testdata lives under `internal/analyzers/testdata/` (today: `sample.go`, `sample.py`, `sample.ts`, `sample.js`, `handler_methods.go` — add new fixtures; do not break goldens).

| Lang | File gate | Symbol | Query fact |
|------|-----------|--------|------------|
| Go | suffix `_test.go` | `function` named `Test*` / `Benchmark*` / `Example*` / `Fuzz*` → persist **`kind=test`** | `goSymbolQuery` already captures `function_declaration`. No `t.Run` subtests |
| Python | basename `test_*.py` or `*_test.py` | `function`/`method` named `test_*`, or `class` named `Test*` → **`kind=test`** | `pySymbolQuery` is module-level + class methods only. Nested `def` **out of scope** |
| JS/TS/TSX | basename `*.test.*` / `*.spec.*` **or** parent dir `__tests__` | New **tree-sitter** `call_expression` (not regex-only) on test files: `test`/`it`/`xtest`/`fit` and `test.only`/`test.skip`/`it.only`/`it.skip`; `describe*` also `kind=test` | `jsSymbolQuery`/`tsSymbolQuery` are **declarations only**. Add `extractJSTSTestCalls` (name OK) gated by the file gate; call from `IndexFile`, not from `Extract` unless you keep API version 1 |

**`validates`:**

1. **EXTRACTED** when the test file’s `imports` row matches a symbol (or file, if `Import.Symbol` is nil) in an **already-indexed** resolved target. Resolve relative JS/TS/Python paths against the test file directory. Go `package foo_test` import of the package under test → sibling non-test `.go` files already in the store.
2. **INFERRED** only with name `heuristic:go_test_name_prefix`: same-directory `_test.go` whose package is **not** `*_test`; `TestFoo` → sibling symbol `Foo` (strip `Test`/`Benchmark`/`Example`/`Fuzz`). Comment in code + a named test.
3. If the test is indexed **before** the target: `IndexFile` of the target **upserts** incoming `validates` from already-indexed test files (unique key). Do not `DELETE` other files’ outgoing edges.
4. Prefer deterministic symbol IDs (`file_id + name + kind + start_line`) so reindex does not churn `to_symbol_id`.

## Requirements

1. Add `code_edges` + store helpers `ReplaceFileEdges`, `ListEdgesByFile`, `ListValidatesForSymbol` on `*store.Store` in `file_graph.go`.
2. Wire `IndexFile` as above. Keep `IndexFileAtRev` as thin caller of `IndexFile`.
3. Named tests below. Keep `TestIncrementalIsolation`, language goldens, `TestAnalyzerImportProvenanceExtracted`.
4. Bump **every** embed-max assertion from **21 → 22**, including `evals/compat/compat_test.go` (`checkMigrateStatus`: `saw022`, forbid **023+**, `EmbedExpected == 22`), `evals/compat/doc.go`, `internal/store/production_hardening_test.go`, `internal/store/deliberation_test.go`. Add `code_edges` to `TestOpenCreatesDBAndMigratesIdempotent` required tables if that list is exhaustive.

## Touch files

- `internal/store/schema/022_code_relationships.sql` (**new**)
- `internal/store/file_graph.go`
- `internal/store/file_graph_test.go` (**new** — none exists today) for `TestReplaceFileEdgesIsFileLocal`
- `internal/analyzers/index.go`, `extract_javascript.go` (test-call query), optionally `extract_go.go` only if you add a package_clause helper
- `internal/analyzers/analyzers_test.go` + testdata
- `evals/compat/compat_test.go`, `evals/compat/doc.go`
- `internal/store/production_hardening_test.go`, `internal/store/deliberation_test.go`, `internal/store/store_test.go` as needed for ceiling/tables

## Named tests

| Test | Proves |
|------|--------|
| `TestIndexDiscoversGoTestFunctions` | `foo_test.go` `TestFoo` stored with `kind=test` |
| `TestValidatesEdgeExtractedFromImport` | test importing a symbol yields EXTRACTED `validates` |
| `TestReplaceFileEdgesIsFileLocal` | reindex path A does not delete edges of path B |
| `TestIndexPythonAndTSTestFiles` | `test_mod.py` + `mod.test.ts` discovered (tree-sitter calls for TS) |
| `TestOpenCreatesDBAndMigratesIdempotent` | keeper — embed max **22** |
| `TestMigrationStatusReportsEmbedMax` | keeper — **22** |
| `TestNoSourceContentColumns` | keeper |
| `TestCompatibilitySecurityChecklist` | ceiling **22**, no 023+ |
| `TestIncrementalIsolation` | keeper — symbols/imports of other files still isolated |
| `TestIndexFileGoGolden` | keeper — non-test `kind` unchanged |

```bash
go test ./internal/analyzers/... ./internal/store/... -count=1 -run 'TestIndexDiscoversGoTestFunctions|TestValidatesEdgeExtractedFromImport|TestReplaceFileEdgesIsFileLocal|TestIndexPythonAndTSTestFiles|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestNoSourceContentColumns|TestIncrementalIsolation|TestIndexFileGoGolden'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [ ] Named tests PASS
- [ ] C03: tests and what they validate are stored as graph edges
- [ ] Compat **22**; no 023+
- [ ] Checklist box **not** checked unless C03 is fully true (this row owns C03)
- [ ] Board Notes: test output

## Minimal todos

- [ ] Mig 022 + `ReplaceFileEdges` / list helpers + file-local test
- [ ] IndexFile test discovery + validates (Go/Python/JS-TS heuristics)
- [ ] Named tests + bump every embed-max **21 → 22**
- [ ] Board status + notes
