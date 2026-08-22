# S04 — Violations + config — scope todos

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | P23-S04-00 | scope planner | **done** 2026-08-20 — locked status `violations[]` + `.trace/config.json` enforce modes; `internal/config` loader; GateForEdit parity; no auto-enforce; 18 named tests in 01/02 |
| 2 | P23-S04-01 | implementer | **done** 2026-08-20 — violations + config loader + status stderr hints |
| 3 | P23-S04-02 | reviewer | **done** 2026-08-20 — APPROVE high; default off; gate parity; no auto-enforce |

**Depends on:** P23-S01-02 done. **Parallel with:** S02/S03 after S01 (S03-02 review pending — not blocking S04).

## Locked artifacts (S04-00)

- Status: additive `violations[]` on `trace.loop.status.v1` (schema string unchanged)
- Violation shape: same as S01 `loop.Violation` / gate JSON
- Status gate: `EvaluateGate(..., GateForEdit)` — parity with `loop gate --for edit`
- Config: `.trace/config.json` → `{ "enforce": "off"|"warn"|"strict" }`
- Default: missing/malformed/invalid → `off`
- Loader: `internal/config.LoadEnforceMode(root)` — file only, no env
- Modes: `off` = JSON only; `warn`/`strict` = identical (+ stderr hints on status); **no auto `--enforce`**
- Init: does not require writing config
- Tests: 14 CLI + 4 unit minimum in `01-violations-config.md`
- No SQL, no `gate.go` changes

## Command shapes (unchanged)

```bash
trace loop status --task <id> [--goal <id>]
# stdout: trace.loop.status.v1 + violations[]
# stderr: hints when config enforce is warn|strict and violations non-empty

# manual config (optional, local):
# .trace/config.json
```

## Blocks

- S05 harness install (status violations + config pointer in AGENTS.md / rules)
