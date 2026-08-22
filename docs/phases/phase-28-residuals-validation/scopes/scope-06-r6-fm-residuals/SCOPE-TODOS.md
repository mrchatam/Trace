# Scope 06 — R6 FM residual wave (locked)

**Planner:** P28-S06-00 (scaffold 2026-08-20)  
**Source FR register:** superseded [`docs/TODO/forward-p28-residuals.md`](../../../../TODO/forward-p28-residuals.md) → run board [`docs/TODO/phase-28.md`](../../../../TODO/phase-28.md)

## Actionable (board)

| Order | ID | FR | FM | Status target |
|------:|----|----|-----|---------------|
| 1 | P28-S06-01 / 02 | FR-P28-01 | FM-01 | implement → review |
| 2 | P28-S06-03 / 04 | FR-P28-02 | FM-02 | implement → review |
| 3 | P28-S06-05 / 06 | FR-P28-03 | FM-04 | implement → review |
| 4 | P28-S06-07 / 08 | FR-P28-04 | FM-07 | implement → review (decision/doc) |
| 5 | P28-S06-09 / 10 | FR-P28-05 | FM-08 | implement → review |
| 6 | P28-S06-11 / 12 | FR-P28-06 | FM-09 | implement → review |
| 7 | P28-S06-13 / 14 | FR-P28-07 | FM-10 | implement → review |

## Explicit defers (not board rows)

| ID | Topic | Why |
|----|-------|-----|
| FR-P28-D1 | Autonomous discovery→task spawn | Human gate required |
| FR-P28-D2 | Full Graphiti / temporal invalidation DB | Spike only after human promote |
| FR-P28-D3 | RESULTS.md parser for P25-4 | Env attestation closed R5 |
| FR-P28-D4 | Hook Option B | S03 locked Option A |
| FR-P28-X1 | Daemon / HTTP / hosted MCP | Project laws / wontfix |

## After S06

→ **S07** residual-wave VERIFY (`P28-S07-00`…) then close Residual wave section in `DR-HANDOFF.md`.
