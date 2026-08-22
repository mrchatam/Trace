# P22-S01-02 — Review: test discovery + `validates` edges

## Metadata
- id: P22-S01-02
- todo_ids: [P22-S01-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Independent review of S01-01. Confirm **C03** is fully met (tests represented; `validates` edges exist with provenance).

## Session start

**Fresh subagent** (not S01-01). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — live inventory + FINAL locks
- [01-test-validates-graph.md](01-test-validates-graph.md)
- [README.md](../../README.md) C03
- Live before S01-01 (for drift): `IndexFile` had no edges; schema max 021; compat 21; extractors as named in 00-PLANNER

## Review checklist

| Check | Evidence |
|-------|----------|
| Mig **`022_code_relationships.sql`** only | `internal/store/schema/` — 22 files; **no** `023_*` |
| Compat **22** | `TestCompatibilitySecurityChecklist`; `TestMigrationStatusReportsEmbedMax`; `deliberation_test.go` embed max |
| File-local incremental | `TestReplaceFileEdgesIsFileLocal` + `TestIncrementalIsolation` |
| EXTRACTED vs INFERRED | grep `heuristic:go_test_name_prefix`; `TestValidatesEdgeExtractedFromImport` |
| JS/TS tests are tree-sitter calls | not regex-only on source text; `extract_javascript.go` / IndexFile |
| Adapter API | `LanguageAdapter.Extract` still `(content []byte)`; `LanguageAdapterAPIVersion == 1` |
| No blobs | `TestNoSourceContentColumns` |
| Seed omit | `SeedDocument` has no `code_edges`; `TestSeedExportOmitsDeniedSurfaces` still forbids `"files"`/`"symbols"` |
| C03 complete | tests indexed (`kind=test`) + `ListValidatesForSymbol` queryable |
| Symbol FK | `ON DELETE SET NULL` (not CASCADE) so `ReplaceFileSymbols` does not drop other files’ incoming validates |

## Keeper commands

```bash
go test ./internal/analyzers/... ./internal/store/... -count=1 -run 'TestIndexDiscoversGoTestFunctions|TestValidatesEdgeExtractedFromImport|TestReplaceFileEdgesIsFileLocal|TestIndexPythonAndTSTestFiles|TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestNoSourceContentColumns|TestIncrementalIsolation|TestIndexFileGoGolden'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Spawn policy

If **C03** is unmet (no test symbols, no validates, full rebuild, blobs, 023+ leak, regex-only JS tests, or adapter API bump without cause): spawn **`P22-S01-02a` implement** + **`P22-S01-02b` review** immediately below this row with full prompts. **Do not** mark this review `done` while leaving residuals for a later phase.

## Exit criteria

- [ ] No blocker/high without spawn or inline fix
- [ ] Confidence **high** with re-run test output in Notes
- [ ] C03 closed or spawned
