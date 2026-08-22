# S07 — Phase verify — scope todos

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | P20-S07-00 | scope planner | **done** — FINAL evidence bar locked |
| 2 | P20-S07-01 | verify | pending — run locked floor + archive CLI evidence |
| 3 | P20-S07-02 | reviewer | pending — re-verify + **close DR-HANDOFF = `no successor`** |

**Depends on:** S01–S06 complete.

## Locked verify floor (import from 00-PLANNER / 01-verify)

- `./internal/deliberation/...` SelectNext + hop budget
- S02–S05 domain 14+14+14+14 named tests
- `./internal/loop/...` 14 S06 named tests
- P19 six loop keepers in `./cmd/trace`
- `./internal/store/...` migration embed max + P20 table roundtrips + `TestNoSourceContentColumns`
- `CGO_ENABLED=1 ./evals/compat/... TestCompatibilitySecurityChecklist` (ceiling **19**)
- P17 seed keepers (P20 export **omitted** by design)
- CLI evidence + §31 mini-eval (fixture-scale OK with residual)

## DR-HANDOFF

- S07-01: gather only
- S07-02: **CLOSED**, successor **`no successor`**
