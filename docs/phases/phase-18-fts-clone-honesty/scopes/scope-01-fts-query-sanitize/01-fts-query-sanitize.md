# P18 / S01 / 01 — FTS query sanitize

## Metadata
- id: P18-S01-01
- todo_ids: [P18-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-87** per sibling **00-PLANNER FINAL**. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT.** It is FINAL — implement.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **SoT** (algorithm + MATCH table + named tests)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- Live: `internal/store/fts.go`; `internal/compiler/compiler.go` (`compileAtDepth` Search)
- Pattern: `failWhyRetriever` in `internal/compiler/compiler_test.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do **not** re-debate FINAL locks. Do **not** implement S02 docs or S03 analyzer. **No board spawn.** Implementers: **status + Notes only**.

## Locked defaults (FINAL — copy)

| Item | Value |
|------|-------|
| Sanitizer | Keep rune iff `unicode.IsLetter \|\| unicode.IsNumber`; else space. `Fields`. Quote each token (`"` … `"`; interior `"` doubled). Join ` AND `. Empty → `""` |
| MATCH examples | `GET /notes` → `"GET" AND "notes"`; `GET /notes/search` → `"GET" AND "notes" AND "search"` — full table in 00-PLANNER |
| SQL | Keep `WHERE fts_docs MATCH ?`. Do not add parameters as the “fix” |
| Packet | Slash titles succeed; Search error → `fts = nil`, continue (Expand-only) |
| Swallow | **Search errors only.** Expand and IncludeWhy fail-closed stay |
| Forbidden | Slash-only `/`→space; phrase-quote whole title as the only strategy; export `sanitizeFTSQuery`; MCP/CLI/seed/analyzer edits |

## Named tests (must exist)

| Test | File |
|------|------|
| `TestSanitizeFTSQueryPunctuationClass` | `internal/store/fts_test.go` |
| `TestSearchFTSSlashInQuery` | `internal/store/fts_test.go` |
| `TestTaskContextSlashTitle` | `internal/compiler/compiler_test.go` |
| `TestTaskContextContinuesWhenSearchErrors` | `internal/compiler/compiler_test.go` |

Keepers must stay green: `TestFTSFindsEntityTitleAndPathSymbol`, `TestIncludeWhyFailClosed`, `TestSeedExportOmitsDeniedSurfaces`.

## Extension points / files

- `internal/store/fts.go` (`sanitizeFTSQuery`) + `fts_test.go`
- `internal/compiler/compiler.go` (Search error fallback in `compileAtDepth`) + `compiler_test.go`

## Role work (TDD)

1. Write the four named tests **red** (`GET /notes` MATCH errors today; TaskContext aborts).
2. Replace `sanitizeFTSQuery` with the charset + quote algorithm. Do not only add `/` to the old `Replacer`.
3. In `compileAtDepth`, if `Search` errors, set `fts = nil` and continue (do not `return Packet{}, err`).
4. Self-check exit criteria; board **status + Notes only** → **P18-S01-02**.

`TestTaskContextContinuesWhenSearchErrors`: stub `Retriever` like `failWhyRetriever` but Search returns a forced error; Expand/Why delegate to `*retrieval.Engine`. Use a slash-free title so the test proves fallback, not the sanitizer.

## Locked verify

```text
CGO_ENABLED=0 go test ./internal/store/... ./internal/compiler/... ./internal/retrieval/... -count=1 -run 'TestSanitizeFTSQueryPunctuationClass|TestSearchFTSSlashInQuery|TestFTSFindsEntityTitleAndPathSymbol|TestTaskContextSlashTitle|TestTaskContextContinuesWhenSearchErrors|TestIncludeWhyFailClosed'
CGO_ENABLED=1 go test ./internal/store/... ./internal/compiler/... ./internal/retrieval/... -count=1 -run 'TestSanitizeFTSQueryPunctuationClass|TestSearchFTSSlashInQuery|TestFTSFindsEntityTitleAndPathSymbol|TestTaskContextSlashTitle|TestTaskContextContinuesWhenSearchErrors|TestIncludeWhyFailClosed'
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSearchFTS|TestTaskContext|TestSeedExportOmitsDeniedSurfaces'
```

## Todo updates
Board **status + Notes only**. Do not spawn rows. Do not thicken later prompts.

## Exit criteria
- [ ] Four named DF-87 tests green (CGO0 + CGO1 store/compiler/retrieval `-run` above)
- [ ] Slash title context packet (`err==nil`; no `syntax error near "/"`)
- [ ] SearchFTS slash queries return the seeded task (not merely no error)
- [ ] Keepers green; no seed/MCP/analyzer edits
- [ ] Board Notes; next **P18-S01-02**

## Minimal todos
- [ ] Red named tests
- [ ] Sanitizer + compiler fallback
- [ ] Locked verify cmds
- [ ] Board sync
