# P23-S03-00 — Enforce DONE + strict export planner

## Metadata
- id: P23-S03-00
- todo_ids: [P23-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock `--enforce` on `trace transition … DONE` and `trace seed export --strict`. **No product Go this row.**

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- Live: `cmd/trace/transition.go`, `cmd/trace/seed.go`, `internal/loop/gate.go`, existing DONE escape hatches (`--allow-done`, `--as-operator`)

## Live inventory (P23-00 + S02 verified)

**Transition today** (`cmd/trace/transition.go`):

- Flags: `--task`, `--to`, `--reason`, `--actor`, `--allow-done`, `--as-operator`, `--allow-missing-caps`, `--evidence`
- DONE requires Review PASS + `--as-operator`, or `--allow-done` hatch
- **No `--enforce`** — S03 adds opt-in only
- Exit codes today: `exitOK=0`, `exitUsage=1`, `exitFail=2`

**Seed export today** (`cmd/trace/seed.go`):

- `trace seed export [-o <file>]` only — **no `--strict` or `--enforce`**
- S03 adds `--strict`, `--enforce`, optional `--task`
- `--strict` alone may warn; `--enforce` fail-closed + no write

**Gate library today** (`internal/loop/gate.go`):

- `EvaluateGate(..., GateForDone | GateForExport)` — export mirrors done policy (verification debt, regression, deliberation incomplete)
- S01-02 residual: export-honesty **extensions** beyond GateForExport deferred; S03 wires CLI only

## Locked defaults (S03-01 must not re-debate)

| Item | Value |
|------|-------|
| Transition | `trace transition --task <id> --to DONE … **--enforce**` calls `EvaluateGate(..., GateForDone)` **before** `TransitionTask` |
| Transition scope | `--enforce` only when `--to DONE` (case-insensitive); **ignored** for other target states |
| Without `--enforce` | **Unchanged** current behavior (all existing flags work) |
| Two-layer DONE | **Gate** = deliberation policy (opt-in); **Domain** = review PASS/FAIL/caps (always) — gate does **not** check review PASS |
| `--enforce` + `--allow-done` | Gate runs first; blocks even if `--allow-done` would bypass review |
| Transition block exit | **1** (policy block, match `exitGateBlocked` in `loop.go`); stderr = violation message; no stdout success JSON |
| Export | `trace seed export [-o path] **--strict**` validates; **`--enforce`** requires `--strict` |
| Export gate scope | All tasks with `work_state ∉ {DONE, SKIPPED, STALE}`; optional **`--task <id>`** for single-task scan |
| Export doc keys | `version == 1`; `goals`/`tasks` non-nil after `BuildSeedDocument` |
| `--strict` alone | Print violations to stderr; exit **0** unless `--enforce` |
| `--strict --enforce` | Exit **1** on violation; **no write** (file or stdout blocked payload) |
| Evaluator | S01 `EvaluateGate` with `GateForDone` / `GateForExport` — no policy changes in S03 |
| Config | S04 `enforce` mode may affect stderr verbosity only — **not** auto `--enforce` in S03 |

## Planner work

1. [x] Map existing DONE transition checks — gate must align, not weaken; domain review/caps stay in `TransitionTask`.
2. [x] Lock export strict checks (GateForExport per active task; doc version; optional `--task` filter).
3. [x] Thicken 01/02 with named CLI tests (17: 9 transition + 8 export).

## Exit criteria

- [x] Flag semantics locked
- [x] Escape hatches preserved when `--enforce` absent
- [x] No product Go

## Next

**P23-S03-01** after this row is `done`.
