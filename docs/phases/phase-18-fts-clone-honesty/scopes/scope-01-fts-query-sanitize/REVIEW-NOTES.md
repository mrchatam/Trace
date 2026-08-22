# P18-S01-02 — FTS query sanitize scope review (DF-87)

**Date:** 2026-08-17  
**Reviewer:** independent (fresh session ≠ implementer)  
**Verdict:** **APPROVE** (confidence: high)  
**Spawn:** none — proceed **P18-S02-00**

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | MATCH still parameterized (`?`) | PASS | `internal/store/fts.go` `SearchFTS` `WHERE fts_docs MATCH ?` then `LIMIT ?`; bind is `match, limit` |
| 2 | Punctuation-as-separators (letter/number charset), not slash-only | PASS | `sanitizeFTSQuery` keeps `unicode.IsLetter \|\| unicode.IsNumber`, else ASCII space. No `strings.NewReplacer`. Table includes `/`, `:`, `AND`, `NEAR`, `work_state` |
| 3 | Remaining tokens FTS5-quoted; join ` AND ` | PASS | `GET /notes` → `"GET" AND "notes"`; `GET /notes/search` → `"GET" AND "notes" AND "search"` (`TestSanitizeFTSQueryPunctuationClass`) |
| 4 | `SearchFTS` slash queries `err==nil` and hit seeded task | PASS | `TestSearchFTSSlashInQuery` subtests `GET /notes` and `GET /notes/search` |
| 5 | `TaskContext` on those titles returns a packet | PASS | `TestTaskContextSlashTitle` `err==nil` + Layer-0 task item |
| 6 | Search error does not abort packet (Expand-only) | PASS | `compileAtDepth` Search → `fts = nil`; Expand still `return Packet{}, err`; IncludeWhy still fail-closed. `TestTaskContextContinuesWhenSearchErrors` + keeper `TestIncludeWhyFailClosed` |
| 7 | Story is MATCH syntax, not “add SQL params” | PASS | Sanitizer builds MATCH expression; SQL bind stays `?` |
| 8 | Not phrase-quote of the whole title as the only strategy | PASS | AND of per-token quotes, not `"GET /notes"` as a single phrase |
| 9 | No MCP / no seed-exclude reverse / no analyzer | PASS | Workspace has no `.git` (diff N/A). Mtimes: S01 files ~18:49; MCP morning; seed ~15:00 (P17); `004_fts.sql` 08-15; analyzers unchanged. Catalog still 10 tools. `TestSeedExportOmitsDeniedSurfaces` still omits `work_state` / transitions / reviews |
| 10 | Keepers green | PASS | `TestFTSFindsEntityTitleAndPathSymbol`; `TestIncludeWhyFailClosed`; `TestSeedExportOmitsDeniedSurfaces` |

Reject-if (none tripped): `/` was not bolted onto the old Replacer; Search no longer `return Packet{}, err`; `sanitizeFTSQuery` stays unexported; packet schema still `0.2`; no seed/MCP/analyzer drive-by.

## Landed `func Test*` names (S04 import)

| Test | File |
|------|------|
| `TestSanitizeFTSQueryPunctuationClass` | `internal/store/fts_test.go` |
| `TestSearchFTSSlashInQuery` | `internal/store/fts_test.go` |
| `TestTaskContextSlashTitle` | `internal/compiler/compiler_test.go` |
| `TestTaskContextContinuesWhenSearchErrors` | `internal/compiler/compiler_test.go` |

Keepers (unchanged names): `TestFTSFindsEntityTitleAndPathSymbol`; `TestIncludeWhyFailClosed`; `TestSeedExportOmitsDeniedSurfaces`.

## Verify (independent re-run)

```text
CGO_ENABLED=0 go test ./internal/store/... ./internal/compiler/... ./internal/retrieval/... -count=1 -run 'TestSanitizeFTSQueryPunctuationClass|TestSearchFTSSlashInQuery|TestFTSFindsEntityTitleAndPathSymbol|TestTaskContextSlashTitle|TestTaskContextContinuesWhenSearchErrors|TestIncludeWhyFailClosed'
→ PASS (store 0.081s; compiler 0.099s; retrieval no tests to run)

CGO_ENABLED=1 same -run -v
→ PASS named:
  TestFTSFindsEntityTitleAndPathSymbol
  TestSanitizeFTSQueryPunctuationClass
  TestSearchFTSSlashInQuery/{GET /notes, GET /notes/search}
  TestIncludeWhyFailClosed
  TestTaskContextSlashTitle/{GET /notes, GET /notes/search}
  TestTaskContextContinuesWhenSearchErrors

CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSearchFTS|TestTaskContext|TestSeedExportOmitsDeniedSurfaces'
→ PASS (GOMODCACHE=$HOME/go/pkg/mod). cmd/trace includes TestSeedExportOmitsDeniedSurfaces. internal/mcp compiled; no matching MCP tests (expected).
```

Sandbox default `GOMODCACHE` lacked `segmentio/encoding` (proxy 403) on the first `./cmd/...` attempt; retry with user module cache is authoritative. Not a product defect.

## Findings

| Severity | Location | Issue | Failure mode |
|----------|----------|-------|--------------|
| — | — | No blocker/high/medium issues | — |

### Residuals (non-fail, documented)

| Severity | Note |
|----------|------|
| low | `TestTaskContextSlashTitle` asserts packet + Layer-0 only. After Search→`fts = nil`, that test would still pass if MATCH syntax failed. MATCH hit coverage is `TestSearchFTSSlashInQuery`. Combined bar meets DF-87; S04 should keep **both** names. |
| low | No `.git` in this workspace; checklist 9 used mtimes + grep + seed keeper instead of `git diff`. |
| nit | First CGO1 `./cmd/trace-mcp`/`./internal/mcp` setup failed on sandbox module cache; user `GOMODCACHE` PASS. |

## Architecture compliance

- FINAL locks in `00-PLANNER.md` satisfied (charset + quote + `AND`-join + `MATCH ?` + compiler Search-only swallow).
- G19: library `store` + `compiler` only; CLI/`trace context` and MCP `context` inherit. Packet schema `0.2` unchanged.
- `retrieval.Search` still propagates `SearchFTS` errors (not hidden for other callers).
- DF-88 exclude untouched; DF-89 analyzers untouched. Did not start S02/S03. Did not own sibling S05 rebuild rows.

## Spawn decision

**No spawn.** Zero blocker/high findings. Next runnable: **P18-S02-00**.
