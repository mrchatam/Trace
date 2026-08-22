# S02 — MCP parity + install — scope todos

**Depends-on:** P10-S01-02 done (serial). Owns DF-21, DF-22, DF-32.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | FINAL locks; thicken 01+02 — **done** |
| 2 | 01-mcp-parity-install | implement | **done** — DF-21/22/32 shipped |
| 3 | 02-scope-review | review | **done** — APPROVE high; no spawns; next P10-S03-00 |

## Reminders
- G19: no domain fork in `internal/mcp`
- No daemon/HTTP; **no** plan/impact/index MCP this phase
- Update `BuiltinMCPCapabilitySpecs` when tools added (**nine** tools; no auto-seed)
- **S01 inherit:** why/Exact `plan-change` alias + `capability` Exact; IncludeWhy fail-closed — do not re-open DF-19/23/25/27/29
- DF-22: docs + stderr tip + `trace_version` — Trace does **not** kill Cursor MCP processes
- DF-32: capability list/missing snake_case rows (`id,kind,slug,title,status`)
