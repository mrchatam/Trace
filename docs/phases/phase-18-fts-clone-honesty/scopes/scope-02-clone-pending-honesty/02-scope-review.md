# P18 / S02 / 02 — clone PENDING honesty review

## Metadata
- id: P18-S02-02
- todo_ids: [P18-S02-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of **DF-88** honesty docs vs [DF-88-DECISION.md](../../DF-88-DECISION.md) + sibling [00-PLANNER.md](00-PLANNER.md) **FINAL** locked strings. Confirm export exclude **unchanged**. Write `REVIEW-NOTES.md` (landed `func Test*` names for S04). Next **P18-S03-00** unless spawn.

**Stop if sibling `00-PLANNER.md` is still DRAFT** (it is **FINAL**).

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-clone-pending-honesty.md](01-clone-pending-honesty.md)
- [DF-88-DECISION.md](../../DF-88-DECISION.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Re-run verify; do not trust Notes. Do not reverse exclude. Do not re-open S01 FTS or S03 analyzer. Do not edit S05 board rows. Fresh session ≠ implementer.

## Checklist

| # | Check | How |
|---|--------|-----|
| 1 | CONTRIBUTING bullet **7** matches 00 locked clone-honesty text (omit reviews/transitions/`work_state`; import PENDING; local DONE/SKIPPED; why/plan without reviews) | Read Portable graph; bullets 1–6 still present |
| 2 | README clone recipe has the locked PENDING sentence | Read `## Portable graph (clone recipe)` |
| 3 | Help seed export has locked omit + PENDING lines **and** still has `trace/graph.json` + `exported_at_commit` + `not identity` | `TestHelpCloneTasksImportPending` + `TestHelpSeedExportPath` |
| 4 | Named test asserts `pending`, `import`, `omits reviews`, `transitions`, `work_state` (not bare `work_state`) | Read `TestHelpCloneTasksImportPending` |
| 5 | `TestSeedExportOmitsDeniedSurfaces` PASS and assertions not weakened | Re-run + read test body vs P17 (no `"transitions"` key; no task `work_state`; no review/token leak) |
| 6 | `SeedTask` still `id/title/body/goal_id` only; comments-only on `seed_export.go` | Diff `internal/domain/seed_export.go` |
| 7 | No new export keys / include flags | Diff `cmd/trace/seed.go`, `BuildSeedDocument` body, help flags |
| 8 | P17 prompts not rewritten | `git diff` `docs/phases/phase-17-portable-graph-git` |
| 9 | S05 rows still pending after VERIFY; no binary rebuild in this scope | Board table P18-S05-00/01/02 |

Reject (blocker/high) if export exclude reversed, include flags added, omit keeper weakened, or clone PENDING documented as a product bug to “fix” by exporting reviews.

## Locked verify (re-run)

```text
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpCloneTasksImportPending|TestHelpSeedExportPath|TestSeedExportOmitsDeniedSurfaces'
```

## REVIEW-NOTES.md (required)

Record: verdict; checklist; landed `func Test*` names (must include `TestHelpCloneTasksImportPending`); keepers; CGO; residuals (CGO0 `cmd/trace` carry-forward is non-fail). S04 imports names from this file.

## Exit criteria
- [ ] REVIEW-NOTES.md; confidence high or medium with residuals listed
- [ ] Landed test names recorded for S04 import
- [ ] No open blocker/high without a pending spawn
- [ ] Board Notes; next **P18-S03-00** unless spawn (S05 still after S04)

## Minimal todos
- [ ] Verify + checklist
- [ ] REVIEW-NOTES.md
- [ ] Board sync
