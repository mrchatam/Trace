# S03 — idempotent import — scope todos

**Depends-on:** P17-S01-02 APPROVE + P17-S02-02 APPROVE (export + convention docs).

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **FINAL** — DF-81/83/84 import locks + named tests |
| 2 | 01-idempotent-import | implementer | pending — stop if 00 DRAFT |
| 3 | 02-scope-review | reviewer | pending |

## Locks (FINAL)
- DF-81: UUID entity upsert + duplicate-link no-op + insert-only events
- DF-83: last-import-wins + CONTRIBUTING union-by-id merge (human resolve; no driver)
- DF-84: plan-tree upsert on import
- Named tests: `TestSeedImportIdempotent`, `TestSeedImportDuplicateLinksNoOp`, `TestSeedImportSameIdLastWins`, `TestSeedImportPlanTreeIdempotent`
- SoT: [00-PLANNER.md](00-PLANNER.md) + [DF-84-FORWARD.md](../../DF-84-FORWARD.md)

## Next
**P17-S03-01** (implementer)
