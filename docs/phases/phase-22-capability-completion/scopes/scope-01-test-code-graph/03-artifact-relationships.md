# P22-S01-03 — Implement: artifact relationships

## Metadata
- id: P22-S01-03
- todo_ids: [P22-S01-03]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Record graph relationships among **files, modules, components, functions, types, and APIs** (tests already in S01-01). Closes **C01** together with S01-01. Board: **status + notes only**.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md)
- Live symbols: `kind` is free text on `symbols` (`001_init.sql` — no CHECK). Extractors emit `function` / `method` / `type` (Go) / `class` (Python/JS/TS). Tests are `kind=test`. Go `_test.go` also persists `kind=package` (package clause) so incoming `validates` can distinguish `package foo_test` from same-package tests. Do not drop `kind=package` on test files.
- Live deps: `imports` table + `ReplaceFileImports` / `ListImportEdges` — C01 “dependencies” is **already** `[x]`.
- Store helpers to add on `file_graph.go`: `ListExports`, `ListModuleContents` (names locked).
- **Live IndexFile edges (S01-01):** `indexValidates` calls `ReplaceFileEdges` **only** when `isTestFile(path)`. Non-test files only run `upsertIncomingValidates` (pair upsert of others’ `validates` onto this target). `ReplaceFileSymbols` DELETE-all then insert; `code_edges` symbol FKs are `ON DELETE SET NULL` — incoming upsert repairs `validates.to_symbol_id` for targets. **This row must call `ReplaceFileEdges` for every indexed file** (test and non-test), merging existing `validates` (tests) with `contains_module` / `exports_api` in one outgoing batch. If you only write contains/exports on a lib file, that is a new `ReplaceFileEdges` (today it is a no-op) — still merge, because a later test-file reindex must not drop contains/exports, and a lib reindex must not drop them either. Do not call `ReplaceFileEdges` with only the new rels on a test file (that would delete `validates`).

## Locked defaults

| Item | Value |
|------|-------|
| SQL | **No new migration** — reuse `code_edges` rels `contains_module`, `exports_api` |
| `depends_on` | **Do not duplicate** `imports` into `code_edges`. Imports **are** the dependency relationship. Tests may JOIN `imports` to prove C01 dependencies still hold |
| Compat | Stays **22** (forbid 023+) |
| Incremental | Same `ReplaceFileEdges(path)` outgoing-only; include **all** outgoing rels for that file (`validates` + `contains_module` + `exports_api`). Today only test files call `ReplaceFileEdges` — **extend to every IndexFile path** |
| Modules | Directory of the file is the module identity (Go package dir, Python package dir, TS/JS file-as-module). **No** `modules` table |
| Components | Language-idiomatic: Go package, Python package, TS/JS module file — do not invent a new ontology |
| `contains_module` | file → its symbols (`from_file_id` = file, `to_symbol_id` = each symbol, `to_file_id` = same file) |
| `exports_api` | public/exported symbols only (provenance EXTRACTED when the grammar shows export) |
| Export rules | Go: `unicode.IsUpper` first rune of `function`/`method`/`type`. Python: name not starting with `_` (optional EXTRACTED `__all__`). JS/TS: EXTRACTED when the declaration is under `export_statement` / `export` keyword; do not regex the source blob |

## Requirements

1. On `IndexFile`, emit `contains_module` + `exports_api` in the same `ReplaceFileEdges` batch as `validates` (S01-01) so reindex of a test file does not drop validates, and reindex of a lib file does not drop its contains/exports.
2. `ListExports(path)` / `ListModuleContents(dirOrPath)` on `*store.Store`.
3. Named tests. Keep S01-01 keepers green.

## Touch files

- `internal/analyzers/index.go` (and extractors only if you add export captures — prefer IndexFile post-pass using existing symbol kinds + a small export query in `extract_javascript.go` / `extract_go.go`)
- `internal/store/file_graph.go`, `file_graph_test.go`
- `internal/analyzers/analyzers_test.go` + testdata

## Named tests

| Test | Proves |
|------|--------|
| `TestArtifactEdgesFunctionsTypesAPIs` | function + type + exported API present as `exports_api` |
| `TestModuleContainsSymbols` | file/module `contains_module` its symbols |
| `TestReplaceFileEdgesIsFileLocal` | keeper |
| `TestIndexDiscoversGoTestFunctions` | keeper — tests still `kind=test`; Go `_test.go` still has `kind=package` |
| `TestValidatesFooTestPackageNotInferred` | keeper — do not drop package-clause persist (incoming `foo_test` honesty) |

```bash
go test ./internal/analyzers/... ./internal/store/... -count=1 -run 'TestArtifactEdgesFunctionsTypesAPIs|TestModuleContainsSymbols|TestReplaceFileEdgesIsFileLocal|TestIndexDiscoversGoTestFunctions|TestValidatesFooTestPackageNotInferred'
```

## Exit criteria

- [ ] Named tests PASS
- [ ] C01 true for files, modules, components, functions, types, APIs, tests (tests via S01-01)
- [ ] No 023+; no `depends_on` clone of `imports`
- [ ] Board Notes with evidence

## Minimal todos

- [ ] Emit contains_module + exports_api on IndexFile (same ReplaceFileEdges batch)
- [ ] ListExports / ListModuleContents + tests
- [ ] Board notes
