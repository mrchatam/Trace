# S06 — MCP / install reload UX — scope todos

**Depends-on:** P11-S05-02 done. Owns DF-22, DF-37, DF-50. (S05 FINAL: why-symbol + depth-2 body redact + trust MD + discovery-mentions-task — no install/reload coupling.)

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks 2026-08-16 |
| 2 | 01-mcp-install-reload | implement | **done** |
| 3 | 02-scope-review | review | **done** — APPROVE high; next **P11-S07-00** |

## FINAL locks (summary)
- **DF-50:** print-only `install cursor` emits **same** stderr tip as `--write`; stdout JSON-only
- **DF-22:** tip on all successful install paths; help/README print+write; keep `trace_version`
- **DF-37:** tip + docs only — cannot force Cursor reload; no new tools/daemon
- Migration **none**; packages thin `cmd/trace` (+ optional README/PROTOCOL)

## Reminders
- Shared tip must contain `reload` + `trace-mcp`
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Forward-only board; implementers: status + Notes only
- Next after APPROVE: **P11-S07-00**
