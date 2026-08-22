# P18 / S03 / 01 — Go method extract golden

## Metadata
- id: P18-S03-01
- todo_ids: [P18-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, grinding-until-pass]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated
- hooks: []

## Objective
Implement **DF-89** per sibling **00-PLANNER FINAL**: named golden so Go extract keeps handler-shaped `Search` / `SearchCursor` as `method:*`. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT** (it is **FINAL**).

Live extract already matches this shape (00 overlay probe). **Expected path:** add fixture + test, stay green, **do not** edit `extract_go.go`. Fix query/extract **only if** the named test is red.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (locked fixture source + exact `kind:name` list)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live: `internal/analyzers/extract_go.go`; `internal/analyzers/analyzers_test.go` `TestIndexFileGoGolden`; `internal/analyzers/testdata/sample.go`
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do **not** add why-by-name. **No board spawn.** Implementers: **status + Notes only**. Do **not** reverse DF-88. Do **not** rebuild S05 binaries. Do **not** rewrite Current focus in `AGENTS.md` (planner already pointed next at this row; after land, Notes say **P18-S03-02**).

## Locked defaults (FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Fixture | `internal/analyzers/testdata/handler_methods.go` — copy **exact** 00 locked source (`package testdata`; `Notes.Search` + `Memory.SearchCursor` pointer receivers; `net/http`; no extra types) |
| Named | `TestIndexFileGoHandlerMethods` in `analyzers_test.go` |
| Index path | Virtual `pkg/handler_methods.go`; bytes from `readTestdata(t, "handler_methods.go")` |
| Exact `kind:name` | `symNamesKinds` **equals** `method:Search`, `method:SearchCursor`, `type:Memory`, `type:Notes` (helper sorts) |
| Keeper | `TestIndexFileGoGolden` unchanged (`method:Run` on `sample.go`) |
| Fix extract | **Only if** named golden red. File-local incremental stays |
| CGO | `CGO_ENABLED=1` required |
| Forbidden | Mutating `sample.go` / keeper `wantSym`; full-rebuild indexer; FTS/seed/MCP; why-by-name; reversing DF-88; goldening `paginateNotes`/`NewNotes` |

## Files to touch

| File | Change |
|------|--------|
| `internal/analyzers/testdata/handler_methods.go` | **New** — 00 locked source |
| `internal/analyzers/analyzers_test.go` | Add `TestIndexFileGoHandlerMethods` beside `TestIndexFileGoGolden` (same `openTemp` / `readTestdata` / `symNamesKinds` pattern) |
| `internal/analyzers/extract_go.go` | **Only if red** — `goSymbolQuery` / extract only |

**Do not touch:** `testdata/sample.go`; FTS; seed; MCP; `cmd/trace`; `bin/`; P17 prompts; DF-88 docs.

## Role work
1. TDD: add locked fixture + `TestIndexFileGoHandlerMethods` asserting the exact four-item `kind:name` list. Do **not** overload the keeper.
2. Run locked verify. **If green:** stop — do not “improve” `extract_go.go`. **If red:** fix `goSymbolQuery` / `extractGo` only; re-run until green (grinding-until-pass). Keep `IndexFile` file-local replace.
3. Confirm keeper `TestIndexFileGoGolden` still PASS with the same `wantSym`.
4. Board Notes → **P18-S03-02**.

## Locked verify

```text
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestIndexFileGoHandlerMethods|TestIndexFileGoGolden'
```

Do **not** use a bare `-run 'TestIndexFileGo'` as the named-DF bar (too broad). CGO0 analyzers are carry-forward non-fail.

## Todo updates
Board **status + Notes only**. Do not spawn. Do not edit S04/S05 rows. Do not reverse DF-88.

## Exit criteria
- [ ] `TestIndexFileGoHandlerMethods` green with exact list `method:Search`, `method:SearchCursor`, `type:Memory`, `type:Notes`
- [ ] `TestIndexFileGoGolden` PASS (unchanged expectations)
- [ ] `extract_go.go` untouched unless the named test was red
- [ ] Board Notes; next **P18-S03-02**

## Minimal todos
- [ ] Fixture + named test (+ extract fix if red)
- [ ] Locked verify
- [ ] Board sync
