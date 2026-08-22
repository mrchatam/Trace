# P18 / S01 / 02 — FTS query sanitize review

## Metadata
- id: P18-S01-02
- todo_ids: [P18-S01-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of **DF-87** vs FINAL locks. Fresh subagent ≠ implementer. Spawn `P18-S01-02a`/`02b` on blocker/high. Write `REVIEW-NOTES.md` in this folder. Next **P18-S02-00** unless spawn.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL SoT**
- Sibling [01-fts-query-sanitize.md](01-fts-query-sanitize.md)
- [phase README](../../README.md)
- Live: `internal/store/fts.go`; `internal/compiler/compiler.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Re-run verify; do not trust Notes alone. Do not re-open S02/S03. Do not reverse DF-88 exclude.

## Checklist

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | MATCH still parameterized (`?`) | Read `SearchFTS` |
| 2 | Sanitizer is punctuation-as-separators (letter/number charset), **not** slash-only | Diff `sanitizeFTSQuery`; `TestSanitizeFTSQueryPunctuationClass` table includes `/` **and** `:`, `AND`, `NEAR`, `work_state` |
| 3 | Remaining tokens are FTS5-quoted; join is ` AND ` | Same table: `GET /notes` → `"GET" AND "notes"` |
| 4 | `SearchFTS("GET /notes")` and `…/search` `err==nil` and hit the seeded task | `TestSearchFTSSlashInQuery` |
| 5 | `TaskContext` on those titles returns a packet | `TestTaskContextSlashTitle` |
| 6 | Search error does not abort packet (Expand-only) | `TestTaskContextContinuesWhenSearchErrors`; Expand/Why still fail-closed |
| 7 | Story is MATCH syntax, not “add SQL params” | Sanitizer + bind still `?` |
| 8 | Not phrase-quote of the **whole title** as the only strategy | MATCH is AND of quoted tokens |
| 9 | No MCP / no seed-exclude reverse / no analyzer | `git diff` vs S01 files; grep |
| 10 | Keepers green | `TestFTSFindsEntityTitleAndPathSymbol`; `TestIncludeWhyFailClosed`; `TestSeedExportOmitsDeniedSurfaces` |

Reject (blocker/high) if: `/` only added to the old `Replacer`; Search still `return Packet{}, err`; `sanitizeFTSQuery` exported without need; packet schema bump; seed/MCP/analyzer drive-by.

## Locked verify (re-run)

```text
CGO_ENABLED=0 go test ./internal/store/... ./internal/compiler/... ./internal/retrieval/... -count=1 -run 'TestSanitizeFTSQueryPunctuationClass|TestSearchFTSSlashInQuery|TestFTSFindsEntityTitleAndPathSymbol|TestTaskContextSlashTitle|TestTaskContextContinuesWhenSearchErrors|TestIncludeWhyFailClosed'
CGO_ENABLED=1 go test ./internal/store/... ./internal/compiler/... ./internal/retrieval/... -count=1 -run 'TestSanitizeFTSQueryPunctuationClass|TestSearchFTSSlashInQuery|TestFTSFindsEntityTitleAndPathSymbol|TestTaskContextSlashTitle|TestTaskContextContinuesWhenSearchErrors|TestIncludeWhyFailClosed'
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSearchFTS|TestTaskContext|TestSeedExportOmitsDeniedSurfaces'
```

## Role work
1. Re-run locked verify. Compare diff to 00-PLANNER MATCH table + compiler fallback.
2. Findings by severity. Blocker/high: small inline fix **or** spawn `P18-S01-02a`/`02b` immediately below this row.
3. Write `REVIEW-NOTES.md` (verdict, checklist evidence, residuals). Board Notes → **P18-S02-00** unless spawn.

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] REVIEW-NOTES.md written (landed `func Test*` names for S04 import)
- [ ] Board status + Notes; next **P18-S02-00** (unless spawn)

## Minimal todos
- [ ] Independent verify + checklist
- [ ] REVIEW-NOTES.md
- [ ] Board sync
