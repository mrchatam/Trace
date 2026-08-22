# S04 — Capability upsert + hatch vs caps — scope todos

**Depends-on:** P11-S03-02 done. Owns DF-41, DF-51. (S03 FINAL: exclusive `trace.lock` + short retry + serialize UX — no capability/hatch coupling.)

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks 2026-08-16 |
| 2 | 01-capability-upsert-hatch | implement | pending — DF-41 slug upsert + DF-51 WARNING/docs |
| 3 | 02-scope-review | review | pending |

## Locked reminders (FINAL)
- **DF-41:** empty-ID `UpsertCapability` resolves by slug → stable id update; explicit different-id clash still fails; **no mig**
- **DF-51:** hatch ≠ missing-caps override (independence); thicken WARNING + help/MCP; keep DF-24 fail-closed + Gate G
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Forward-only board; implementers: status + Notes only
- Next after APPROVE: **P11-S05-00**
