# P18-S03-00 — Go method extract golden (FINAL)

## Metadata
- id: P18-S03-00
- todo_ids: [P18-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock **FINAL** for **DF-89**: Go tree-sitter must extract handler-shaped methods `Search` and `SearchCursor`. Live `extract_go.go` already has `method_declaration`. This scope is a **named golden** (fix query only if red). **No why-by-name CLI. No indexer architecture change.** No product Go in this planner row. Stop if [../../00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) is not FINAL (it is).

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/rules/skills-map.md](../../../../rules/skills-map.md)
- [phase README](../../README.md)
- Live: `internal/analyzers/extract_go.go`; `internal/analyzers/testdata/sample.go`; `TestIndexFileGoGolden`
- D40 observation: G2 `Search` / `SearchCursor` missing from symbol table after reindex
- D40 shapes: `experiments/runs/2026-08-17-ab-compare/workspaces/G2/internal/handlers/notes.go` (`func (h *Notes) Search(w http.ResponseWriter, r *http.Request)`); `.../internal/store/memory.go` (`func (m *Memory) SearchCursor(...)`)
- [docs/TODO.md](../../../../TODO.md)
- S02 landed (no seed coupling): [../scope-02-clone-pending-honesty/REVIEW-NOTES.md](../scope-02-clone-pending-honesty/REVIEW-NOTES.md)

## Session start
Follow agent-loop-protocol. Depends-on: **P18-S02-02 APPROVE** (board order) — no S02 code coupling. **S02 landed:** DF-88 docs/help/comments + `TestHelpCloneTasksImportPending` (keepers `TestSeedExportOmitsDeniedSurfaces`, `TestHelpSeedExportPath`). Seed exclude **unchanged** (do not reverse; do not add reviews/`work_state` to export). **CGO_ENABLED=1 required** for analyzer tests (tree-sitter). CGO0 analyzers / `cmd/trace` remain carry-forward non-fail.

## Live inventory (2026-08-18)

| Area | Present? | Gap vs DF-89 |
|------|----------|--------------|
| `goSymbolQuery` `(method_declaration name: (field_identifier) @name) @method` | **Yes** | Name-agnostic. Pointer-receiver `func (w *Worker) Run()` already extracts as `method:Run` (`TestIndexFileGoGolden`) |
| Extract post-filter on names | **No** | `extractGo` maps capture kind + name text; no skip list |
| `IndexFile` `.go` → `goAdapter.Extract` → `extractGo` | **Yes** | File-local `ReplaceFileSymbols` (DR-INCREMENTAL). Not a full-rebuild indexer |
| Fixture `testdata/handler_methods.go` | **No** | New in 01. Do **not** mutate `sample.go` |
| Named `TestIndexFileGoHandlerMethods` | **No** | New in `analyzers_test.go` beside the keeper |
| D40 G2 miss `Search` / `SearchCursor` | Observation **disproven as a query hole** | Live `IndexFile` of the D40 G2 sources already emits `method:Search` / `method:SearchCursor` (and `function:NewNotes` / `function:paginateNotes`). Golden is the regression lock. D40 SCORE was operator / stale workspace index / `why <type> <uuid>` (Exact by id, not by name) — **not** a missing `method_declaration` capture |

### Live extract vs A/B (read-only — no product Go)

`CGO_ENABLED=1` `IndexFile` of D40 G2 files (deleted scratch; not in repo) + keeper:

| Input | `kind:name` includes |
|-------|----------------------|
| `workspaces/G2/internal/handlers/notes.go` | **`method:Search`**, `function:NewNotes`, `method:List`, `method:Create`, `type:Notes`, … |
| `workspaces/G2/internal/store/memory.go` | **`method:SearchCursor`**, `function:paginateNotes`, `method:ListCursor`, `type:Memory`, `type:Note`, … |
| Mini fixture (locked source below) | **`method:Search`**, **`method:SearchCursor`**, `type:Memory`, `type:Notes` |
| `TestIndexFileGoGolden` | **PASS** (`method:Run` on `sample.go`) |

A/B SCORE also claimed `paginateNotes` / `NewNotes` missing. Those extract today as `function:*`. That is the same class of miss as the methods: **not** `goSymbolQuery`.

**Lock:** golden-only if green (expected). Touch `extract_go.go` **only if** `TestIndexFileGoHandlerMethods` is red.

Keeper analog: `testdata/sample.go` `func (w *Worker) Run()` → `method:Run`. Receiver identifier is not captured (`n` vs G2 `h` is not the fail bar).

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Fixture path | `internal/analyzers/testdata/handler_methods.go` (**new**; do not mutate `sample.go` / `TestIndexFileGoGolden` expectations) |
| Fixture package | `package testdata` |
| Fixture shape | Pointer-receiver methods: `func (n *Notes) Search(w http.ResponseWriter, r *http.Request)` and `func (m *Memory) SearchCursor(query, cursor string, limit int) (items []string, next string)`. Types `Notes` and `Memory` only (no extra `type Note` — keeps the exact list at four). Import `"net/http"` |
| Receiver names | Lock `n` / `m` as sketched. G2 used `h *Notes`; identifier is **not** in `kind:name` |
| Index path in test | Virtual `pkg/handler_methods.go` (mirrors keeper `pkg/sample.go`). Read bytes via `readTestdata(t, "handler_methods.go")` |
| Named | **`TestIndexFileGoHandlerMethods`** in `internal/analyzers/analyzers_test.go` |
| Exact `kind:name` list | Sorted equality (same helper as keeper `symNamesKinds`): **`method:Search`**, **`method:SearchCursor`**, **`type:Memory`**, **`type:Notes`** |
| Fail bar | Must include **`method:Search`** and **`method:SearchCursor`**. The exact four-item list is the named-test SoT |
| Keeper | `TestIndexFileGoGolden` still PASS (`method:Run` on `sample.go`). Do **not** change its `wantSym` |
| If green on current extract | Golden is the deliverable (regression lock). **Expected:** green without `extract_go.go` edits |
| If red | Fix `goSymbolQuery` / extract **only** — file-local incremental stays (DR-INCREMENTAL). Do not add a full-rebuild indexer |
| CGO | Verify with **`CGO_ENABLED=1`**. Do not use CGO0 as this scope’s bar |
| S05 | **Leave** P18-S05-00/01/02 after VERIFY. This scope does **not** rebuild binaries |
| Out | Why-by-name CLI; FTS; seed; MCP; full-rebuild index; reversing DF-88 exclude |

### Locked fixture source (FINAL)

01 copies this into `internal/analyzers/testdata/handler_methods.go` (bodies may be empty; do not add extra types/funcs that would grow the `kind:name` list):

```go
package testdata

import "net/http"

type Notes struct{}

type Memory struct{}

func (n *Notes) Search(w http.ResponseWriter, r *http.Request) {}

func (m *Memory) SearchCursor(query, cursor string, limit int) (items []string, next string) {
	return nil, ""
}
```

D40 `SearchCursor` returned `[]Note`; fixture uses `[]string` so extract does not also emit `type:Note`. Tree-sitter still captures `method:SearchCursor`.

### Named test (FINAL)

`TestIndexFileGoHandlerMethods` pattern (copy keeper shape, do **not** overload `TestIndexFileGoGolden`):

1. `openTemp` + `readTestdata(t, "handler_methods.go")`
2. `IndexFile(..., "pkg/handler_methods.go", content, IndexOptions{})`
3. Optional: language is `LangGo`
4. `ListSymbolsByPath("pkg/handler_methods.go")` → `symNamesKinds` **equals** the exact four-item list above (order via `sort` in the helper)
5. Optional line-range sanity (`StartLine >= 1`) like the keeper

Do **not** use a bare `-run 'TestIndexFileGo'` as the named-DF bar (too broad).

## Files likely touched (implementer)

| File | Change |
|------|--------|
| `internal/analyzers/testdata/handler_methods.go` | **New** — locked source above |
| `internal/analyzers/analyzers_test.go` | Add `TestIndexFileGoHandlerMethods` beside `TestIndexFileGoGolden` |
| `internal/analyzers/extract_go.go` | **Only if** named test red |

**Do not touch:** `testdata/sample.go`; keeper `wantSym`; FTS; seed; MCP; `cmd/trace`; S05 `bin/`; P17 prompts; DF-88 docs.

## Locked verify

```text
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestIndexFileGoHandlerMethods|TestIndexFileGoGolden'
```

CGO0 analyzers are carry-forward non-fail (tree-sitter). Do not use CGO0 as this scope’s bar.

## Blast / later scopes (upcoming only)

- **S04 VERIFY:** import **`TestIndexFileGoHandlerMethods`** + keeper `TestIndexFileGoGolden` (names locked here; live name from S03 REVIEW-NOTES wins if 01 must rename — it must not). CGO1 analyzers. Do not treat CGO0 analyzers as fail.
- **S05:** Unchanged board rows after VERIFY. Golden is not a binary rebuild.
- **S02:** No coupling. Do not reverse DF-88 exclude.

## Non-goals
- Product Go on **this** planner row
- Why-by-name CLI; `why symbol Search` product surface
- Goldening D40 `paginateNotes` / `NewNotes` (functions; not DF-89)
- Mutating `sample.go` to add Search/SearchCursor
- FTS / seed / MCP / full-rebuild-on-any-change indexer
- Reversing DF-88 exclude; hosted MCP; DF-86 hook

## Planner work (this row)
1. [x] Live `goSymbolQuery` / extract vs scratch handler snippet (overlay; no product Go)
2. [x] Lock fixture path + exact `kind:name` list as **FINAL**
3. [x] Thicken 01/02/SCOPE-TODOS; light S04 Depends; S05 rows left in place
4. [x] Mark this prompt **FINAL**; board Notes; next **P18-S03-01**
5. [ ] Product Go — **not** this row

## Exit criteria
- [x] This prompt **FINAL** with fixture path + test name + exact `kind:name` list
- [x] 01/02 thickened enough to run alone
- [x] Board Notes; next **P18-S03-01**
- [x] No product Go this row

## Minimal todos
- [x] Live extract_go + golden pattern + scratch probe
- [x] FINAL + thicken
- [x] Board sync

## Next
Orchestrator: **P18-S03-01**. Do **not** start P18-S03-02 until S03-01 is `done`.
