# P22-S01-05 — Implement: architectural boundaries

## Metadata
- id: P22-S01-05
- todo_ids: [P22-S01-05]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Track **architectural relationships and boundaries** as graph edges. Closes **C02**.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md)
- Live: no architecture overlay file in-repo today; `IndexFile` is per-path; `cmd/trace/index.go` walks indexable files and calls `analyzers.IndexFile` (or `IndexFileAtRev`)
- Store: `code_edges` from S01-01; `ListImportEdges` for observed cross-layer imports
- **Live IndexFile outgoing edges (S01-03, reviewed S01-04):** `IndexFile` calls `ReplaceFileEdges(path, nil)` to pre-clear **all** outgoing rels, then `indexCodeEdges` writes **one** outgoing batch = `validates` + `contains_module` + `exports_api`. `ReplaceFileEdges` deletes `WHERE from_file_id = this file`. Any `architectural_boundary` whose `from_file_id` is that file is gone unless it is in that same rewrite batch.

## Locked defaults

| Item | Value |
|------|-------|
| Rel | `architectural_boundary` |
| Provenance | EXTRACTED from `go.mod` module path / package dirs; INFERRED from conventional layer folder names (`internal/`, `cmd/`, `pkg/`) with tests |
| SQL | No 023+; compat stays **22** |
| Optional overlay | If committed `trace/architecture.json` exists, EXTRACTED boundaries from that file — **do not require it** (absent in this repo today) |
| Incremental | Dirty file + **its package directory** only. Do not scan the whole repo on every `IndexFile` |
| Graph | SQLite `code_edges` only — not a new graph database; not LLM-classified layers |
| Query | Store enough for “what layer is this file in” and “what does this package depend on across layers” (helpers on `file_graph.go`) |
| Outgoing merge | **Same `ReplaceFileEdges` batch** as S01-03: `validates` + `contains_module` + `exports_api` + `architectural_boundary`. Compute boundaries in `indexCodeEdges` (or a helper it calls) and `append` them. One `ReplaceFileEdges` call per IndexFile path |
| Second replace | Forbidden. A follow-up `ReplaceFileEdges` with only `architectural_boundary` deletes validates/contains/exports. A batch that omits `architectural_boundary` deletes prior boundaries on that file |
| from_file_id | Persist package→layer on the **indexed source file** (file/package-local). Do not introduce a synthetic architecture file that IndexFile never rewrites — that breaks the incremental unit |

## Requirements

1. Emit `architectural_boundary` edges (package→layer membership at minimum: e.g. `cmd/trace` vs `internal/store`). Observed cross-layer imports may be recorded as INFERRED membership or left as `imports` — do not rebuild the import graph.
2. Incremental: `TestArchitecturalBoundaryIncremental` — editing one file does not rewrite all boundaries in the store.
3. Reindex merge: `TestArchitecturalBoundarySurvivesFileReindex` — IndexFile of the same path twice keeps `architectural_boundary` **and** that file’s `contains_module` / `exports_api` (and `validates` if it is a test file).
4. Named tests + small testdata repo layout (temp dir with `cmd/` + `internal/` is enough).

## Touch files

- `internal/analyzers/index.go` / `test_graph.go` `indexCodeEdges` (append `architectural_boundary` into the existing outgoing batch; still file/package-local)
- `internal/store/file_graph.go` (query helpers)
- `internal/analyzers/analyzers_test.go` + testdata

## Named tests

| Test | Proves |
|------|--------|
| `TestArchitecturalBoundaryEdges` | `internal/` vs `cmd/` membership edges |
| `TestArchitecturalBoundaryIncremental` | editing one file does not rebuild all boundaries |
| `TestArchitecturalBoundarySurvivesFileReindex` | same-path IndexFile keeps boundary + artifact rels (outgoing merge) |
| `TestReplaceFileEdgesIsFileLocal` | keeper |

```bash
go test ./internal/analyzers/... ./internal/store/... -count=1 -run 'TestArchitecturalBoundaryEdges|TestArchitecturalBoundaryIncremental|TestArchitecturalBoundarySurvivesFileReindex|TestReplaceFileEdgesIsFileLocal'
```

## Exit criteria

- [ ] C02 true with stored edges
- [ ] Incremental (Law 12) — no full-graph rebuild
- [ ] Named tests PASS
- [ ] Board Notes

## Minimal todos

- [ ] Boundary extraction + store (package-local), appended onto the existing `indexCodeEdges` outgoing batch
- [ ] Tests
- [ ] Board notes
