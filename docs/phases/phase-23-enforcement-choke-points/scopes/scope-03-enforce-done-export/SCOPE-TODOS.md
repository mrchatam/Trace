# S03 — Enforce DONE + strict export — scope todos

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | P23-S03-00 | scope planner | **done** 2026-08-20 — locked `--enforce` on DONE + export `--strict`/`--enforce`; two-layer DONE model; 17 named CLI tests; touch points in 01/02 |
| 2 | P23-S03-01 | implementer | pending — `transition.go` + `seed.go` + help + tests |
| 3 | P23-S03-02 | reviewer | pending — backward compat + no-write-on-enforce-fail |

**Depends on:** P23-S02-02 done (`trace loop gate` CLI). **Blocks:** S05 harness rules referencing export enforce + transition `--enforce`.

## Locked command shapes (S03-00)

```bash
trace transition --task <id> --to DONE --reason <text> [--as-operator] [--allow-done] [--enforce]
trace seed export [-o path] [--strict] [--enforce] [--task <id>]
```

- Transition `--enforce`: `EvaluateGate(..., GateForDone)` before `TransitionTask`; exit **1** on block.
- Export `--strict`: scan active tasks (or `--task` only) with `GateForExport`; stderr warnings; exit **0**.
- Export `--strict --enforce`: exit **1** + no write on any violation.
- Without new flags: behavior identical to pre-S03.
