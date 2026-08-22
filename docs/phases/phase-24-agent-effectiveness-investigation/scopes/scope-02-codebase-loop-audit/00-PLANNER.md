# P24-S02-00 — Codebase loop audit planner

## Objective
Lock audit map: which packages to read for FM-01..FM-10. Thicken `01-codebase-loop-audit.md`.

## Required inventory (minimum — verified live 2026-08-20)

| Area | Paths | Notes |
|------|-------|-------|
| Loop CLI | `cmd/trace/loop.go` | ✓ exists |
| Loop core | `internal/loop/{next,apply,policy,gate}.go` | ✓ all four; also `apply_writes.go`, `deliberation_packet.go` |
| Deliberation | `internal/deliberation/select.go` | ✓ + `types.go` |
| Task add | `cmd/trace/add.go`, `internal/domain/create.go` | store via `internal/store/entities.go` |
| MCP | `cmd/trace-mcp/main.go` (entry); **`internal/mcp/`** handlers | add→`tools_write.go`; loop→`tools_loop.go`; context→`tools_context.go`; tasks→`tools_parity.go`; registry→`server.go` |
| Install | `internal/install/enforcement.go` | ✓ + cursor/claude installers |
| Seed | `cmd/trace/seed.go` | ✓ + `internal/domain/seed_export.go`, `seed_import.go` |

## S01 residuals locked for S02-01

- SelectNext reason-code mapping (`p19_saturated` vs `hop_budget_exceeded`)
- Export vs SQLite sync (G1 `task_not_found` vs graph.json)
- Deliberation reset after gap pass
- FM→code paths for all FM-01..FM-10

## Deliverable (S02-01)

`CODEBASE-AUDIT.md` — table: FM-ID | mechanism | file:line | agent-visible symptom | change lever

## Next
**P24-S02-01**
