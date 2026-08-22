# S05 — rebuild CLI + MCP binaries — scope todos

**Depends-on:** P18-S04-02 APPROVE + P18-S05-00 **FINAL**. Last Phase 18 scope. Owns **DR-HANDOFF close**. **S04-00 FINAL:** VERIFY is named DF-87/88/89 + carry-forward only (two-clone **not required**; DF-88 document-only); stale binaries remain this scope’s job — do not start from S04-01.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **FINAL** — next **P18-S05-01** |
| 2 | 01-rebuild-binaries | implementer | **done** on board — next **P18-S05-02** |
| 3 | 02-scope-review | reviewer / handoff close | **done** — DR-HANDOFF **CLOSED** `no successor`; Phase 18 complete |

## Phase locks (00 FINAL)

- Rebuild `bin/trace` with **`CGO_ENABLED=1`**; `bin/trace-mcp` with **`CGO_ENABLED=0`**
- Preferred env: `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off` (sandbox 403 `segmentio/encoding` class — retry is not a product defect)
- `trace-mcp -h` lists **10** tools including `trace_impact` (fail if any name missing)
- Optional thin DF-87: `trace context` on title `GET /notes/search` — **skip = non-fail** (Notes required); **run red = FAIL**
- No product feature Go; do not reverse DF-88; do not re-run S04 named suite as fail bar
- DR-HANDOFF default **`no successor`** — 02 row closes it (S04 only starts Notes; 00/01 must not close)
- Planner row does **not** rebuild

## Reminders
- Stale MCP binary is the defect this scope exists to prevent — rebuild even if live `-h` already lists 10
- Not a successor phase; not hosted MCP; not research S05
