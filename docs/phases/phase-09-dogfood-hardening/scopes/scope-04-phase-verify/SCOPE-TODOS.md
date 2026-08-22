# P09-S04 — SCOPE-TODOS

| Order | ID | Role | Status |
|------:|----|------|--------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks DF-01 + S02/S03 spot-checks + carry-forward + `./...`; DR-HANDOFF=`no successor` |
| 2 | 01-verify | verify | pending — run locked commands; write VERIFY-NOTES; start handoff |
| 3 | 02-scope-review | review / handoff | pending — owns DR-HANDOFF completion |

## Locked command summary (see 01-verify.md)

1. DF-01: `TestWhyAndContextWithLinkedReview`
2. S02: `TestTasksListAfterSeed` + `TestSeedImportRelativePathAgainstC`
3. S03: `TestInstallCursor*` + DF-05 docs; **no** MCP list-tasks
4. Carry-forward: honesty G, replan E, impact F, capability ablation, perf H, compat checklist, p0x, x0, Gate C `dry_run:false`
5. Full: `CGO_ENABLED=1 go test ./...`

## DR-HANDOFF

**`no successor`** — remaining ladder gaps (D08/D09/combos/multi-agent; DF-11/12) stay parallel `experiments/`.
