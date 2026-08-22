# P18-S01-00 — FTS query sanitize (FINAL)

## Metadata
- id: P18-S01-00
- todo_ids: [P18-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Replan against **live** `sanitizeFTSQuery` + `TaskContext`. Lock **FINAL** defaults for **DF-87**. Thicken sibling `01`/`02`/SCOPE-TODOS. **No product Go in this row.** [../../00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) is FINAL — proceed.

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Laws 6, 7, 15, 19
- [phase README](../../README.md)
- [../../00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — phase locks (do not re-debate)
- Live: `internal/store/fts.go` (`sanitizeFTSQuery`, `SearchFTS` `MATCH ?`); `internal/store/schema/004_fts.sql` (`tokenize = 'unicode61'`); `internal/compiler/compiler.go` `compileAtDepth` Search on `task.Title`; `internal/retrieval/search.go`; `internal/compiler/compiler_test.go` `failWhyRetriever`
- Evidence: `experiments/runs/2026-08-17-ab-compare/verify/G1-context.err` (`fts5: syntax error near "/"`)
- Findings: [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) DF-87
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Depends-on: **P18-00 done**. **No product Go.** Do not re-open DF-88 exclude or DF-89 golden.

## Live inventory (2026-08-17)

| Area | Present? | Gap vs DF-87 |
|------|----------|--------------|
| `SearchFTS` `WHERE fts_docs MATCH ?` | **Yes** | Not an SQL-injection / “add params” bug |
| `sanitizeFTSQuery` strips `" ' * ( ) { } [ ] ^ ~ : - +` then `Fields` + `AND`-join | **Yes** | Replacer **omits `/`**. `strings.Fields("GET /notes")` → `GET`, `/notes` → MATCH `GET AND /notes` → FTS5 syntax error |
| unicode61 indexes `/` as a separator | **Yes** | Query language still rejects `/` before tokenization of MATCH |
| `compileAtDepth` Search on `task.Title` | **Yes** | `if err != nil { return Packet{}, err }` **aborts** the whole packet |
| Why / Exact by UUID | **Yes** | Does not MATCH the title — why still works in D40 |
| MCP/CLI context | Library `TaskContext` | **Zero** adapter edits (G19) |
| Seed export / DF-88 | Untouched | S01 must not change seed format |

### Repro (locked)

D40 G1: task titles `GET /notes`, `GET /notes/search`. `trace context` → `store: search fts: SQL logic error: fts5: syntax error near "/"`.

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Bug class | MATCH **expression** syntax, not SQL injection. Keep `MATCH ?` |
| Home | `sanitizeFTSQuery` in `internal/store/fts.go`; Search-error fallback in `internal/compiler/compiler.go` `compileAtDepth` |
| G19 | Library only. CLI/`trace context` and MCP `context` inherit. **No** new flags, tools, or MCP schema |
| Migration | **None** (004_fts unicode61 stays) |
| Packet schema | Unchanged (`0.2`). No new honesty/FTS-error fields |
| DF-88 | Untouched — no seed export/import edits |
| DF-89 | Untouched — no analyzer edits |
| Forbidden | “Add SQL parameters” as the fix; **slash-only** escape; phrase-quoting the **whole title** as the only strategy; swallowing **Expand** or **IncludeWhy** errors; reversing P17 exclude; daemon/hosted MCP |

### Sanitizer algorithm (FINAL)

Replace the current explicit-replacer list. Do **not** only add `/` to that list.

1. Walk runes. **Keep** a rune iff `unicode.IsLetter(r) || unicode.IsNumber(r)`. **Every other rune is a separator** (write ASCII space). This class **includes `/`** and also `. : - _ + * " ' ( )` and FTS operator punctuation — aligned with FTS5 **unicode61** default tokenchars (Unicode `L*` / `N*`).
2. `strings.Fields` the result. Drop empty tokens.
3. **Quote** each remaining token as an FTS5 phrase: wrap in `"…"`; double any interior `"` (`"` → `""`) even though the charset should never leave quotes in a token.
4. Join with ` AND ` (spaces required). The `AND` between phrases is **ours**; user words `AND`/`OR`/`NOT`/`NEAR` become quoted literals.
5. If no tokens remain (`///`, empty, punctuation-only), return `""`. `SearchFTS` already no-ops empty match (`nil, nil`).

Locked expected MATCH strings (table test SoT):

| Input | `sanitizeFTSQuery` output |
|-------|---------------------------|
| `GET /notes` | `"GET" AND "notes"` |
| `GET /notes/search` | `"GET" AND "notes" AND "search"` |
| `zephyrunique` | `"zephyrunique"` |
| `foo AND bar` | `"foo" AND "AND" AND "bar"` |
| `title:secret` | `"title" AND "secret"` |
| `NEAR foo` | `"NEAR" AND "foo"` |
| `work_state` | `"work" AND "state"` |
| `  GET   /notes  ` | `"GET" AND "notes"` |
| `///` | `` (empty) |
| `C++` | `"C"` |

Keepers: existing `TestFTSFindsEntityTitleAndPathSymbol` / `TestFTSBackfillOnOpenWhenIndexEmpty` must stay green (`"zephyrunique"` phrase ≡ token match under unicode61).

### Compiler fallback (FINAL)

In `compileAtDepth`, after Expand (Expand errors **still abort**):

```text
fts, err := c.retr.Search(ctx, task.Title, retrieval.SearchOptions{Limit: 16})
if err != nil {
    fts = nil  // Expand-only; do not return Packet{}
}
```

- Silent: no new packet field, no log requirement.
- File-seed Expand (DF-65) still runs on Expand∪FTS; FTS empty on Search error is OK.
- `TestIncludeWhyFailClosed` (DF-29) **unchanged**: Why errors still fail-close when `IncludeWhy=true`.
- Swallow **Search only**. Do not change retrieval.Search / SearchFTS to hide MATCH errors from other callers.

### Context success bar (FINAL)

Task titled `GET /notes` and `GET /notes/search`: `TaskContext` returns a packet (Layer-0 task present). Error must not contain `syntax error near "/"`. Same `compileAtDepth` path covers `ExpandContext`.

`SearchFTS` of those query strings: `err == nil` **and** the seeded task id is among hits (`"GET" AND "notes"` matches unicode61-indexed title tokens).

## Named tests (FINAL)

| Test | File | Assert |
|------|------|--------|
| `TestSanitizeFTSQueryPunctuationClass` | `internal/store/fts_test.go` | Table above (exact strings). Same-package; do **not** export `sanitizeFTSQuery` |
| `TestSearchFTSSlashInQuery` | `internal/store/fts_test.go` | Subtests `GET /notes` and `GET /notes/search`: seed task with that title; `SearchFTS(title)` `err==nil`; hit `entity_type=task` + that id |
| `TestTaskContextSlashTitle` | `internal/compiler/compiler_test.go` | Same two titles via `TaskContext`; `err==nil`; Layer-0 task item present |
| `TestTaskContextContinuesWhenSearchErrors` | `internal/compiler/compiler_test.go` | Stub `Retriever` (mirror `failWhyRetriever`): Search returns error; Expand/Why delegate. Plain title. `TaskContext` `err==nil`; Layer-0 task present |
| Keepers | existing files | `TestFTSFindsEntityTitleAndPathSymbol`; `TestIncludeWhyFailClosed`; `TestSeedExportOmitsDeniedSurfaces` (untouched) |

TDD: red tests first (`GET /notes` MATCH errors today), then sanitizer, then compiler fallback.

## Files likely touched (implementer)

- `internal/store/fts.go` — `sanitizeFTSQuery` only (SearchFTS bind stays `?`)
- `internal/store/fts_test.go` — two named tests
- `internal/compiler/compiler.go` — Search error → `fts = nil`
- `internal/compiler/compiler_test.go` — two named tests + small Search-fail stub

**Do not touch:** seed export/import, analyzers, MCP tools, CLI flags, `004_fts.sql`, DF-88 docs.

## Locked verify

```text
CGO_ENABLED=0 go test ./internal/store/... ./internal/compiler/... ./internal/retrieval/... -count=1 -run 'TestSanitizeFTSQueryPunctuationClass|TestSearchFTSSlashInQuery|TestFTSFindsEntityTitleAndPathSymbol|TestTaskContextSlashTitle|TestTaskContextContinuesWhenSearchErrors|TestIncludeWhyFailClosed'
CGO_ENABLED=1 go test ./internal/store/... ./internal/compiler/... ./internal/retrieval/... -count=1 -run 'TestSanitizeFTSQueryPunctuationClass|TestSearchFTSSlashInQuery|TestFTSFindsEntityTitleAndPathSymbol|TestTaskContextSlashTitle|TestTaskContextContinuesWhenSearchErrors|TestIncludeWhyFailClosed'
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1 -run 'TestSearchFTS|TestTaskContext|TestSeedExportOmitsDeniedSurfaces'
```

modernc.org/sqlite: CGO0 store/compiler tests are expected to run. Product `cmd` CGO0 remain carry-forward non-fail (not this scope’s bar).

## Blast / later scopes (upcoming only)

- **S02:** S01 does not change seed format. DF-88 stays wontfix+docs. Slash titles are context/FTS only.
- **S03:** No analyzer coupling.
- **S04 VERIFY:** import the four named tests above from S01 REVIEW-NOTES after land; keepers listed.

## Non-goals
- Product Go on **this** planner row
- Why-by-name CLI; embeddings; FTS prefix/NEAR features; column filters
- Kitchen-sink harness (rsync / stdio EOF)

## Planner work (this row)
1. [x] Live re-read `fts.go` + `compileAtDepth` Search abort + G1-context.err
2. [x] Lock token charset + quote + AND-join + compiler Search fallback + exact test names
3. [x] Thicken `01-fts-query-sanitize.md` / `02-scope-review.md` / SCOPE-TODOS; light S02/S04 Depends
4. [x] Mark this prompt **FINAL**; board Notes; next **P18-S01-01**
5. [ ] Product Go — **not** this row

## Exit criteria
- [x] This prompt **FINAL** with locked algorithm, MATCH table, test names, files, verify cmds
- [x] Sibling 01/02/SCOPE-TODOS thickened enough to run alone
- [x] Board status + Notes; next **P18-S01-01**
- [x] No product Go

## Minimal todos
- [x] Live re-read fts.go + compiler Search call
- [x] FINAL locks + thicken 01/02
- [x] Board sync

## Next
Orchestrator: **P18-S01-01**. Do **not** start P18-S01-02 until S01-01 is `done`.
