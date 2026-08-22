# S03 — Index GC — scope todos

**Depends-on:** P10-S02-02 done. Owns DF-20.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks 2026-08-16; 01+02 thickened |
| 2 | 01-index-gc | implement | **done** — ListFilePaths + DeleteFileByPath + FTS; full-tree set-diff; argv missing single-delete |
| 3 | 02-scope-review | review | **done** — APPROVE high; no spawns — [REVIEW-NOTES.md](REVIEW-NOTES.md) |

## Reminders
- DR-INCREMENTAL: **no** full-rebuild-on-any-change — set-diff delete-on-missing only
- Full-tree GC when `len(args)==0`; argv must **not** GC siblings (`TestIndexIncrementalIsolation`)
- Delete must clear FTS (not only CASCADE symbols/imports)
- Fuel for experiment DF-14 residual — product fix lives here
- **S02:** MCP tasks/capability/version unrelated — do not touch MCP tool surface
- **S04:** serial after S03-02; no index coupling
- **S05 VERIFY:** re-prove `TestIndexGCAfterPathRename` + `TestIndexIncrementalIsolation`
