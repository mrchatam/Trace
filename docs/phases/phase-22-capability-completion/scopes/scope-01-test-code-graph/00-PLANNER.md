# P22-S01-00 — Planner: test→code graph + architecture

## Metadata
- id: P22-S01-00
- todo_ids: [P22-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep]
- verification: automated

## Objective

Lock S01 against **live** analyzers + impact walk. Owned capabilities: **C01, C02, C03, C07**. **No product Go this row.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DECISION-LOG.md](../../DECISION-LOG.md) D-22-06, D-22-21
- [WORK-MAP.md](../../WORK-MAP.md) W-01…W-04
- Live: `internal/analyzers/index.go`, `extract.go`, `extract_go.go`, `extract_ts.go`, `extract_py.go`, `internal/store/file_graph.go`, `internal/retrieval/impact_walk.go`, schema max **021**

## Live inventory (confirmed 2026-08-18 — before S01-01)

| Surface | Live |
|---------|------|
| `IndexFile` | `internal/analyzers/index.go` — `UpsertFile` + `SetFileLanguage` + `ReplaceFileSymbols` + `ReplaceFileImports`; **no** edges / test classification |
| Extract | `extract.go` `runQuery`; `extract_go.go` `goSymbolQuery`/`goImportQuery`; `extract_py.go` `pySymbolQuery`/`pyImportQuery`; `extract_javascript.go` `jsSymbolQuery` + `tsSymbolQuery` + `extractJSTS` (file is **`extract_javascript.go`**, not `extract_js.go` — GOOS=js); `extract_ts.go` delegates to `extractJSTS` |
| Adapters | `language_adapter.go` `LanguageAdapter.Extract(content []byte)` — **no path**; `LanguageAdapterAPIVersion = 1` (do not bump unless iface changes) |
| Store | `internal/store/file_graph.go`: `UpsertFile`, `ReplaceFileSymbols`, `ReplaceFileImports`, `ListSymbolsByPath`, `ListImportsByPath`, `ListImportEdges` (path+provenance only; **no** import symbol on the edge DTO). **No** `ReplaceFileEdges`. **No** `file_graph_test.go` |
| ImpactWalk | `internal/retrieval/impact_walk.go` — seeds `file\|symbol` only; neighbors = contains-OUT + incoming imports + contains-UP; `MaxImpactBlast=64`; depth **1..2**; no `validates` |
| Schema | **21** files; max `021_experiments.sql`; **no** `code_edges` |
| Compat | Ceiling **21** (`evals/compat/compat_test.go` `checkMigrateStatus`; `internal/store/production_hardening_test.go`; `internal/store/deliberation_test.go`) |
| Seed | `SeedDocument` has no files/symbols/code_edges; keeper `TestSeedExportOmitsDeniedSurfaces` already forbids `"files"` / `"symbols"` |
| MCP | `trace_impact` action `walk` already encodes `ImpactWalkResult` (`internal/mcp/tools_impact.go`); catalog **10** tools (`TestToolNamesRegistered`) |
| CLI | `cmd/trace/impact.go` `cmdImpactWalk` encodes seeds/blast/depth; `cmd/trace/index.go` calls `analyzers.IndexFile` per path |

## FINAL locked defaults (implementers must not re-debate)

| Item | Value |
|------|-------|
| SQL migration | **`022_code_relationships.sql`** in S01-01 only |
| Compat after S01-01 | **22** (forbid 023+) |
| Table | `code_edges`: ids, from/to file_id + optional symbol_id, `rel`, `provenance` |
| `rel` enum | `validates` \| `contains_module` \| `exports_api` \| `architectural_boundary` \| `depends_on` |
| Provenance | `EXTRACTED` \| `INFERRED` \| `AMBIGUOUS` — same honesty as imports (DF-64) |
| Incremental | `ReplaceFileEdges(path)` on IndexFile — **no** full-graph rebuild |
| Test discovery | See **Locked test heuristics** below (must match live queries) |
| `validates` | EXTRACTED when test file imports/references the symbol; INFERRED only with documented heuristic; never silent |
| Seed JSON | **Still omit** code graph (D-22-06) — no `code_edges` on `SeedDocument` |
| MCP/CLI query of tests | **S05-03** (C31) — this scope ships **library + index**; S01-07 may extend existing `trace impact walk` / `trace_impact` walk JSON via the library DTO. **No new MCP tool** (catalog stays 10) |

## Locked test heuristics (live tree-sitter)

Classify in **`IndexFile`** (has path). Do **not** change `LanguageAdapter.Extract` signature.

| Lang | Test **file** gate | Test **symbol** | Live query |
|------|--------------------|-----------------|------------|
| Go | `strings.HasSuffix(path, "_test.go")` | `function` named `Test*` / `Benchmark*` / `Example*` / `Fuzz*` → persist `kind=test` | `goSymbolQuery` already captures `function_declaration`. **No** `t.Run` subtests |
| Python | basename `test_*.py` or `*_test.py` | `function`/`method` named `test_*`, or `class` named `Test*` → `kind=test` | `pySymbolQuery` is **module-level** functions + class methods. Nested `def` inside a test is **out of scope** |
| JS/TS/TSX | basename `*.test.*` / `*.spec.*` **or** parent dir `__tests__` | **New** tree-sitter `call_expression` on test files only: callees `test`/`it`/`xtest`/`fit` and members `test.only`/`test.skip`/`it.only`/`it.skip`; `describe`/`describe.only`/`describe.skip` also `kind=test` | `jsSymbolQuery` / `tsSymbolQuery` capture **declarations only** — they will **not** see Jest/Vitest `test()` today. Add `extractJSTSTestCalls` (or equivalent) gated by the file gate |

**`validates` provenance (DF-64 honesty):**

- **EXTRACTED:** test-file `imports` row whose `Symbol` matches a symbol in a **already-indexed** resolved target file; or (JS/TS/Python) relative `ImportedPath` resolves to that file (`to_symbol_id` NULL = file-level). Go `package foo_test` import of the package under test → sibling non-test `.go` files.
- **INFERRED** named `heuristic:go_test_name_prefix` only: same-directory `_test.go` whose package is **not** `*_test`, test name `TestFoo`/`BenchmarkFoo`/`ExampleFoo`/`FuzzFoo` → sibling symbol `Foo`. Comment + test required.
- Never write an edge with empty/garbage provenance (mirror `validateImportProvenance`).

**Index order / incremental:** `ReplaceFileEdges(path)` deletes **outgoing** edges (`from_file_id` of `path`) only. If a test file is indexed before its target, `IndexFile` of the **target** may upsert incoming `validates` from already-indexed test files (unique key) — still not a full-graph rebuild. Symbol FKs: `ON DELETE SET NULL`; file FKs: `ON DELETE CASCADE`. Prefer deterministic symbol IDs (`file_id+name+kind+start_line`) so reindex does not churn incoming `to_symbol_id`.

## Named tests (must exist by end of S01)

| Test | Row | Proves |
|------|-----|--------|
| `TestIndexDiscoversGoTestFunctions` | S01-01 | test symbols indexed |
| `TestValidatesEdgeExtractedFromImport` | S01-01 | EXTRACTED `validates` |
| `TestReplaceFileEdgesIsFileLocal` | S01-01 | incremental; other files untouched |
| `TestIndexPythonAndTSTestFiles` | S01-01 | py + ts/js test discovery |
| `TestArtifactEdgesFunctionsTypesAPIs` | S01-03 | C01 non-test artifacts |
| `TestModuleContainsSymbols` | S01-03 | module/file contains |
| `TestArchitecturalBoundaryEdges` | S01-05 | C02 boundaries |
| `TestImpactWalkIncludesAffectedTests` | S01-07 | C07 |
| `TestImpactWalkDepthStillCapped` | S01-07 | existing depth 1–2 keepers |

P21 keepers: `TestNoSourceContentColumns`; `TestIncrementalIsolation`; `TestIndexFileGoGolden` (and JS/Python/TS goldens); `go test ./internal/analyzers/...`; `go test ./internal/retrieval/... -run TestImpact`.

## Planner confirmation

- No **023+** this scope. S01-03/05/07 reuse `code_edges` (compat stays **22**).
- No product Go this row.

## Touch files (scope)

- `internal/store/schema/022_code_relationships.sql`
- `internal/store/file_graph.go` (+ tests)
- `internal/analyzers/index.go`, `extract_*.go`
- `internal/retrieval/impact_walk.go` (S01-07)
- `evals/compat/compat_test.go` ceiling **22** after S01-01
- `cmd/trace/impact.go` / `internal/mcp/tools_impact.go` only to encode library `affected_tests` (S01-07; no new MCP tool)

## Planner work

1. Re-read live extractors; lock test heuristics that match real tree-sitter queries.
2. Thicken 01–08 if live names drifted.
3. Confirm no 023+ this scope.

## Exit criteria

- [x] 01–08 prompts match live files
- [x] Mig **022** / compat **22** locked
- [x] No product Go

## Next

**P22-S01-01**
