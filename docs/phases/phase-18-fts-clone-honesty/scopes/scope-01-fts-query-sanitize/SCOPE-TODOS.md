# S01 — FTS query sanitize — scope todos

**Depends-on:** P18-00 done. **SoT:** [00-PLANNER.md](00-PLANNER.md) (**FINAL**).

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **FINAL** — done |
| 2 | 01-fts-query-sanitize | implementer | **done** (board) |
| 3 | 02-scope-review | reviewer | **APPROVE** — next P18-S02-00 |

## Phase locks (FINAL)

- DF-87: `sanitizeFTSQuery` — keep `unicode.IsLetter \|\| unicode.IsNumber`; other runes (incl. `/`) are separators; quote tokens; `AND`-join; `MATCH ?` stays
- Compiler: Search error → Expand-only (`fts = nil`); do not abort packet. Expand / IncludeWhy fail-closed unchanged
- Named: `TestSanitizeFTSQueryPunctuationClass`, `TestSearchFTSSlashInQuery`, `TestTaskContextSlashTitle`, `TestTaskContextContinuesWhenSearchErrors`
- Keepers: `TestFTSFindsEntityTitleAndPathSymbol`, `TestIncludeWhyFailClosed`, `TestSeedExportOmitsDeniedSurfaces`
- No MCP; no DF-88 export reverse; no analyzer; no product Go on 00

## Reminders
- Not “add SQL params”. Not slash-only. Not phrase-quote the whole title as the only strategy.
- S02/S03 do not need S01 seed or analyzer changes.
- S05 rebuild is after VERIFY — not this scope.
