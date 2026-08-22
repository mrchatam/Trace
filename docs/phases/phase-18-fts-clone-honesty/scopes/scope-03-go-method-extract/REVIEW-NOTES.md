# P18-S03-02 — Go method extract scope review (DF-89)

**Date:** 2026-08-18  
**Reviewer:** independent (fresh session ≠ implementer)  
**Verdict:** **APPROVE** (confidence: high)  
**Spawn:** none — proceed **P18-S04-00**

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | Named golden exact list is `method:Search`, `method:SearchCursor`, `type:Memory`, `type:Notes` | PASS | `TestIndexFileGoHandlerMethods` `wantSym` is those four strings. `symNamesKinds` sorts. Fail path is length + pairwise equality (not a subset/contains-only check) |
| 2 | Fixture is 00 locked handler shape | PASS | `internal/analyzers/testdata/handler_methods.go` byte-matches 00 FINAL: `package testdata`; `import "net/http"`; `type Notes`/`type Memory` only; `func (n *Notes) Search(w http.ResponseWriter, r *http.Request)`; `func (m *Memory) SearchCursor(query, cursor string, limit int) (items []string, next string)` |
| 3 | Fixture path `internal/analyzers/testdata/handler_methods.go`; virtual index `pkg/handler_methods.go` | PASS | File exists. Test: `readTestdata(t, "handler_methods.go")` + `IndexFile(..., "pkg/handler_methods.go", ...)` + `ListSymbolsByPath("pkg/handler_methods.go")` |
| 4 | Existing `TestIndexFileGoGolden` still PASS; `sample.go` / keeper `wantSym` untouched | PASS | Independent CGO1 re-run PASS. `sample.go` mtime 2026-08-16 03:25 (pre-S03). Keeper `wantSym` still `function:Helper`/`function:Main`/`method:Run`/`type:Counter`/`type:ID`/`type:Worker` — no Search/SearchCursor. DF-89 vehicle is `handler_methods.go`, not `sample.go` |
| 5 | `extract_go.go` unchanged unless named test was red; if changed, only query/extract | PASS | **Not touched.** mtime 2026-08-17 01:07 vs fixture+test 2026-08-18 08:41. `goSymbolQuery` still has `(method_declaration name: (field_identifier) @name) @method`. No extract drive-by on an already-green golden |
| 6 | No why-by-name CLI; no FTS/seed/MCP; no full-rebuild indexer; DF-88 exclude not reversed | PASS | No `why-by-name` / `why symbol` in `cmd/`. `IndexFile` still path-local `ReplaceFileSymbols` (`index.go` mtime 2026-08-16). FTS `internal/store/fts.go` 2026-08-17 18:49; MCP 2026-08-17 07:10; `cmd/trace/seed.go` 2026-08-17 15:04 — all before S03-01. `SeedTask` still `id/title/body/goal_id` only (no `work_state`). No `--include-reviews` / `--include-work-state`. `TestHelpCloneTasksImportPending` + `TestSeedExportOmitsDeniedSurfaces` still present |
| 7 | CGO analyzer tests green on locked verify | PASS | Independent re-run below. Not CGO0. Not bare `-run 'TestIndexFileGo'` |
| 8 | S05 rows still pending after VERIFY; no binary rebuild in this scope | PASS | Board P18-S05-00/01/02 still `pending`. `bin/trace` + `bin/trace-mcp` mtime 2026-08-17 17:32 (not 2026-08-18) |

Reject-if (none tripped): named test missing the two methods; fixture not handler-shaped; keeper mutated; `sample.go` used as DF-89 vehicle; extract drive-by while golden was already green; full-rebuild indexer; why-by-name; reversing DF-88.

## Landed `func Test*` names (S04 import)

| Test | File |
|------|------|
| `TestIndexFileGoHandlerMethods` | `internal/analyzers/analyzers_test.go` |

Keeper (unchanged name): `TestIndexFileGoGolden`.

Exact `kind:name` SoT for the named test: `method:Search`, `method:SearchCursor`, `type:Memory`, `type:Notes`.

## Verify (independent re-run)

```text
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestIndexFileGoHandlerMethods|TestIndexFileGoGolden'
→ PASS (0.047s)

CGO_ENABLED=1 same -run -v
→ PASS named:
  TestIndexFileGoGolden (0.02s)
  TestIndexFileGoHandlerMethods (0.02s)
```

CGO0 analyzers are carry-forward non-fail (tree-sitter). Not used as this scope’s bar.

## Findings

| Severity | Location | Issue | Failure mode |
|----------|----------|-------|--------------|
| — | — | No blocker/high/medium issues | — |

### Residuals (non-fail, documented)

| Severity | Note |
|----------|------|
| nit | No `.git` in this workspace; checklist 4/5/6/8 used mtimes instead of `git diff`. |
| nit | CGO0 `./internal/analyzers` remains carry-forward non-fail (tree-sitter). |
| nit | `AGENTS.md` Current focus still says next **P18-S03-01** (01 was forbidden to rewrite it). S04-00 planner owns the next-runnable line. |

## Architecture compliance

- 00-PLANNER FINAL: golden-only-if-green. Named test + fixture landed; `extract_go.go` untouched (query already captured `method_declaration`).
- Keeper analog intact: `sample.go` `func (w *Worker) Run()` → `method:Run`. Receiver identifier not in `kind:name`.
- `IndexFile` remains file-local incremental (`ReplaceFileSymbols` / `ReplaceFileImports` per path). Not a full-rebuild indexer.
- Did not goldene `paginateNotes` / `NewNotes`. Did not add why-by-name CLI. Did not reverse DF-88 exclude. Did not own S05 rebuild.

## Five-axis (code-review-and-quality)

| Axis | Result |
|------|--------|
| Correctness | Named test asserts exact four-item sorted `kind:name` including both handler methods; keeper unchanged and green |
| Readability | Copies keeper `openTemp` / `readTestdata` / `symNamesKinds` / line-range pattern; no extra abstraction |
| Architecture | Regression lock only; extract/query/indexer untouched; fixture isolated from `sample.go` |
| Security | Testdata fixture; no secrets, no new input surface |
| Performance | No hot-path change |

Change size: fixture (~13 lines) + one test (~40 lines). Healthy.

## Spawn decision

**No spawn.** Zero blocker/high findings. Next runnable: **P18-S04-00**.
