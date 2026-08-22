# S02 — Loop gate CLI — scope todos

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | P23-S02-00 | scope planner | **done** — `trace.loop.gate.v1` schema + exit 0/1/2 locked; 01/02 thickened |
| 2 | P23-S02-01 | implementer | **next** — `trace loop gate` CLI + 14 named tests |
| 3 | P23-S02-02 | reviewer | pending — exit code + JSON + thin-adapter review |

**Depends on:** P23-S01-02 done. **Blocks:** S05 harness hooks (call gate CLI).

## Locked artifacts (S02-00)

- Command: `trace loop gate --task <uuid> [--for orient|edit|execute|done|export]`
- Default `--for`: `edit`
- Exit: **0** allowed, **1** blocked, **2** usage/internal (stdout JSON only on 0/1)
- Schema: `trace.loop.gate.v1` — `{schema_version, task_id, for, allowed, violations[], recommended_phase?, reason_code?}`
- Top-level lift: `recommended_phase` + `reason_code` from `violations[0]` when blocked
- Evaluator: S01 `EvaluateGate` only — no policy in cmd
- Help: gate subcommand + exit semantics in `printLoopHelp`
- Tests: 14 named CLI tests in `01-loop-gate-cli.md`
- No MCP, no SQL
