# P18-S02-02 — clone PENDING honesty scope review (DF-88)

**Date:** 2026-08-18  
**Reviewer:** independent (fresh session ≠ implementer)  
**Verdict:** **APPROVE** (confidence: high)  
**Spawn:** none — proceed **P18-S03-00**

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | CONTRIBUTING bullet **7** matches 00 locked clone-honesty text; bullets 1–6 still present | PASS | `CONTRIBUTING.md` Portable graph: bullets 1–6 unchanged (path/PR/clone/evidence/merge/hook). Bullet 7 exact vs 00 FINAL: omit reviews/transitions/`work_state`; import PENDING until clone operator transitions; live DONE/SKIPPED local; `why` / `plan show` without reviews/`work_state` |
| 2 | README clone recipe has the locked PENDING sentence | PASS | `README.md` `## Portable graph (clone recipe)`: sentence after the bash fence is exact 00 FINAL. Index/CONTRIBUTING sentence kept |
| 3 | Help seed export has locked omit + PENDING **and** `trace/graph.json` + `exported_at_commit` + `not identity` | PASS | `cmd/trace/help.go` seed-export block exact 00 FINAL. `TestHelpCloneTasksImportPending` + `TestHelpSeedExportPath` PASS. Import/handoff/build-note still present |
| 4 | Named test asserts `pending`, `import`, `omits reviews`, `transitions`, `work_state` (not bare `work_state`) | PASS | `TestHelpCloneTasksImportPending` in `cmd/trace/cli_test.go`: `captureStdout` + `run([]string{"help"})`; lowercased contains all five substrings. Comment states not bare `work_state` |
| 5 | `TestSeedExportOmitsDeniedSurfaces` PASS and assertions not weakened | PASS | Independent CGO1 re-run PASS. Still fails on `"transitions"` key, task `work_state`, and leaks (`secret-token`, `access.token`, `capability_tool_decisions`, `review_judges_task`, `node_modules`, `"files"`, `"symbols"`) |
| 6 | `SeedTask` still `id/title/body/goal_id` only; comments-only on `seed_export.go` | PASS | JSON tags `id`, `title`, `body`, `goal_id` only (no `work_state`). Comments on `SeedTask` + `BuildSeedDocument` match 00 locked text. Builder still copies `ID/Title/Body/GoalID` only; `Transitions` never filled (nil + `omitempty`) |
| 7 | No new export keys / include flags | PASS | No `--include-reviews` / `--include-work-state`. `ExportOpts` still `ProjectRoot` only. `cmd/trace/seed.go` mtime 2026-08-17 15:04 (P17); S02 files 2026-08-18 ~08:27 |
| 8 | P17 prompts not rewritten | PASS | No `.git` in this workspace. P17 mtimes 2026-08-17 (`00-PHASE-PLANNER.md` 06:50; S01 seed-export prompts 14:52; `DF-84-FORWARD.md` 06:50). S02 docs/help 2026-08-18 |
| 9 | S05 rows still pending after VERIFY; no binary rebuild in this scope | PASS | Board P18-S05-00/01/02 still `pending`. `bin/trace` + `bin/trace-mcp` mtime 2026-08-17 17:32. Help MCP build-note still `CGO_ENABLED=1` (S05 owns retarget; S02 did not touch) |

Reject-if (none tripped): export exclude not reversed; no include flags; omit keeper not weakened; clone PENDING documented as expected, not as a product bug to fix by exporting reviews.

## Landed `func Test*` names (S04 import)

| Test | File |
|------|------|
| `TestHelpCloneTasksImportPending` | `cmd/trace/cli_test.go` |

Keepers (unchanged names): `TestSeedExportOmitsDeniedSurfaces`; `TestHelpSeedExportPath`.

## Verify (independent re-run)

```text
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpCloneTasksImportPending|TestHelpSeedExportPath|TestSeedExportOmitsDeniedSurfaces'
→ PASS (0.049s)

CGO_ENABLED=1 same -run -v
→ PASS named:
  TestHelpSeedExportPath
  TestHelpCloneTasksImportPending
  TestSeedExportOmitsDeniedSurfaces
```

CGO0 `cmd/trace` is carry-forward non-fail (tree-sitter). Not used as this scope’s bar.

## Findings

| Severity | Location | Issue | Failure mode |
|----------|----------|-------|--------------|
| — | — | No blocker/high/medium issues | — |

### Residuals (non-fail, documented)

| Severity | Note |
|----------|------|
| low | `TestHelpCloneTasksImportPending` also matches `import` / `transitions` / `work_state` on other help lines (`seed import`, “operator transitions them”, `tasks` `work_state`). Unique fail bars are **`pending`** + **`omits reviews`** — planner-locked. Combined set meets DF-88. |
| low | No `.git` in this workspace; checklist 8 used mtimes instead of `git diff` on `docs/phases/phase-17-portable-graph-git`. |
| nit | `AGENTS.md` Current focus still says next **P18-S02-01** (01 was forbidden to rewrite it). S03-00 planner owns the next-runnable line. |
| nit | CGO0 `./cmd/trace` remains carry-forward non-fail (tree-sitter). |

## Architecture compliance

- 00-PLANNER FINAL locked strings satisfied (CONTRIBUTING bullet 7, README sentence, help omit+PENDING, `seed_export.go` comments).
- DF-88-DECISION option A: keep P17 exclude; clone PENDING expected; `why` / `plan show` without reviews/`work_state`.
- `SeedTask` JSON and `BuildSeedDocument` body unchanged. No seed-format change. No MCP/FTS/analyzer drive-by.
- Did not start S03-00 impl. Did not own S05 rebuild rows. Did not reverse export omit.

## Spawn decision

**No spawn.** Zero blocker/high findings. Next runnable: **P18-S03-00**.
