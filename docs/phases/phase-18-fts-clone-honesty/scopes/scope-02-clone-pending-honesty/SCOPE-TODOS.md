# S02 — clone PENDING honesty — scope todos

**Depends-on:** P18-S01-02 APPROVE (board order). SoT: [DF-88-DECISION.md](../../DF-88-DECISION.md) + sibling [00-PLANNER.md](00-PLANNER.md) **FINAL**. S01 does not change seed format (slash titles are FTS/context only).

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **FINAL** — next **P18-S02-01** |
| 2 | 01-clone-pending-honesty | implementer | **done** — DF-88 docs/help/comments + named test |
| 3 | 02-scope-review | reviewer | **APPROVE** high — next **P18-S03-00** |

## Phase locks (FINAL)

- DF-88 **wontfix** product: keep omit reviews/`work_state`/`transitions` (and P17 denied surfaces)
- Docs/help/comments: clone PENDING expected (not a product bug)
- CONTRIBUTING bullet **7** + README one sentence + help omit/PENDING lines (00 exact strings)
- Named: `TestHelpCloneTasksImportPending` (new; do not overload `TestHelpSeedExportPath`)
- Keepers: `TestSeedExportOmitsDeniedSurfaces`, `TestHelpSeedExportPath`
- CGO: `CGO_ENABLED=1` for `./cmd/trace/...`
- S05 rebuild stays after VERIFY — **not** this scope (leave P18-S05-* rows)

## Reminders
- Do not reverse P17 export exclude
- S04 VERIFY imports the named help test + omit keeper + path keeper
- S03 has no seed coupling
