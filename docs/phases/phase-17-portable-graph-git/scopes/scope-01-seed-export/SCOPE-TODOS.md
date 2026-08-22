# S01 — seed export — scope todos

**Depends-on:** Phase 16 complete (P16-S06-02 done). P16 S05 seed keys (`findings`/`alternatives`/`discovery_mentions_task`) already live — export only.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **FINAL** — 2026-08-17 |
| 2 | 01-seed-export | implementer | pending — stop if 00 DRAFT |
| 3 | 02-scope-review | reviewer | pending |

## Locked (00 FINAL)
- CLI: `trace seed export [-o <file>]` | import unchanged
- Keys: P16 causal + `findings`/`alternatives` + **`plan_phases`/`plan_scopes`/`scope_deep_plans`/`goal_plan_state`** + **`exported_at_commit`** (import ignore)
- Named: `TestSeedExportRoundTrip`, `TestSeedExportOmitsDeniedSurfaces`, `TestSeedExportWritesExportedAtCommit`
- G19: `internal/domain` export builder + thin CLI
- Idempotent upsert → **S03**; docs/convention → **S02**

## Reminders
- DF-80 + **DF-84** + **DF-85** here; no MCP; no `.trace/` commit
- SoT: [00-PLANNER.md](00-PLANNER.md) + [DF-84-FORWARD.md](../../DF-84-FORWARD.md)
