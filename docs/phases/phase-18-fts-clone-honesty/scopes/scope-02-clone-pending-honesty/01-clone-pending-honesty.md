# P18 / S02 / 01 — clone PENDING honesty

## Metadata
- id: P18-S02-01
- todo_ids: [P18-S02-01]
- role: implementer
- skills: [incremental-implementation, writing-for-agents]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Document **DF-88** per sibling **00-PLANNER FINAL** + [DF-88-DECISION.md](../../DF-88-DECISION.md). Keep export omit. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT** (it is **FINAL**).

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (locked strings + test asserts)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [DF-88-DECISION.md](../../DF-88-DECISION.md)
- `CONTRIBUTING.md`; `README.md` Portable graph clone recipe; `cmd/trace/help.go`; `internal/domain/seed_export.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do **not** add reviews/`work_state` to export. **No board spawn.** Implementers: **status + Notes only**. Do **not** rewrite Current focus in `AGENTS.md` (planner already pointed next at this row; after land, Notes say **P18-S02-02**).

## Locked defaults (FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Product | Keep omit. No `--include-reviews`. No `--include-work-state`. No new seed keys |
| Named | `TestHelpCloneTasksImportPending` (new). Pattern: `captureStdout` + `run([]string{"help"})` beside `TestHelpSeedExportPath` |
| Fail bar | `TestSeedExportOmitsDeniedSurfaces` PASS (do not weaken) |
| Path keeper | `TestHelpSeedExportPath` unchanged purpose (DF-82/85) |
| CGO | `CGO_ENABLED=1` for `./cmd/trace/...` |
| S05 | Do **not** rebuild binaries; do **not** retarget MCP CGO build-note lines |
| Forbidden | Include flags; rewriting P17 prompts; changing `SeedTask` JSON; hijacking AGENTS current-focus into a novel product story |

## Files to touch

| File | Change |
|------|--------|
| `CONTRIBUTING.md` | Append locked bullet **7** (Clone honesty DF-88). Do not delete/renumber 1–6 |
| `README.md` | One locked sentence after the clone-recipe fence |
| `cmd/trace/help.go` | Replace `seed export` continuation with locked block (keep path + `exported_at_commit` + `not identity`) |
| `cmd/trace/cli_test.go` | Add `TestHelpCloneTasksImportPending` |
| `internal/domain/seed_export.go` | Comments only on `SeedTask` + `BuildSeedDocument` (00 locked text). **No** builder/JSON change |
| `AGENTS.md` | Optional append on Portable graph hard-boundary bullet only |

**Do not touch:** `cmd/trace/seed.go`, `BuildSeedDocument` body, MCP, analyzers, FTS, P17 `docs/phases/phase-17-*`, `bin/`.

## Role work
1. TDD: add red `TestHelpCloneTasksImportPending` with 00 locked asserts (`pending`, `import`, `omits reviews`, `transitions`, `work_state` on lowercased `trace help` stdout).
2. Help + CONTRIBUTING + README + comments per **exact** 00 locked strings.
3. Optional AGENTS portable-graph clause (append only).
4. Re-run locked verify. Confirm omit keeper still fails if reviews/`work_state` appeared (it must not).
5. Board Notes → **P18-S02-02**.

## Named test asserts (copy from 00)

`TestHelpCloneTasksImportPending`: lowercased help must contain:

- `pending`
- `import`
- `omits reviews`
- `transitions`
- `work_state`

Do **not** treat a lone `work_state` on the `transition`/`tasks` lines as sufficient.

## Locked verify

```text
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpCloneTasksImportPending|TestHelpSeedExportPath|TestSeedExportOmitsDeniedSurfaces'
```

CGO0 `cmd/trace` is carry-forward non-fail (tree-sitter). Do not use it as this scope’s bar.

## Todo updates
Board **status + Notes only**. Do not spawn. Do not reverse DF-88. Do not edit S05 rows.

## Exit criteria
- [ ] Docs/help/comments match DF-88-DECISION + 00 locked strings (clone PENDING expected; omit stays)
- [ ] `TestHelpCloneTasksImportPending` green
- [ ] `TestSeedExportOmitsDeniedSurfaces` and `TestHelpSeedExportPath` PASS
- [ ] No export-include flags; `SeedTask` still has no `work_state` JSON tag
- [ ] Board Notes; next **P18-S02-02**

## Minimal todos
- [ ] Red named help test
- [ ] CONTRIBUTING bullet 7 + README sentence + help block + comments
- [ ] Locked verify
- [ ] Board sync
