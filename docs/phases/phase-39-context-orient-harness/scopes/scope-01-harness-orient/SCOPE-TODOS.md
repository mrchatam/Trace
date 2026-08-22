# Scope 01 — board map (G3)

**S01 harness orient** — G-006, G-010. Serial: **P39-S01-00 → P39-S01-01 → P39-S01-02**.

| Order | Board ID | Prompt | Role | Status |
|------:|----------|--------|------|--------|
| 675 | P39-S01-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | **done** |
| 676 | P39-S01-01 | [01-implement.md](01-implement.md) | Implementer | pending |
| 677 | P39-S01-02 | [02-review.md](02-review.md) | Reviewer | pending |

## Planner locks (P39-S01-00 — verified 2026-08-22)

| Lock | Value |
|------|-------|
| Theme | G3 MCP/harness orient |
| Verdict | **Accept** — Instructions + docs hygiene |
| Moat lead | `trace_tasks` → `trace_context`(+`query`) → `trace_loop` → `trace_review` → `trace_plan` |
| 9/16 | Docs + Instructions + `trace_version` — **not** tool reduction |
| G2 compose-first | Orient recipe in Instructions (no explore implement) |
| G1 | **Shipped S00** — `ContextInput.Query` + compiler merge; reflect in Instructions |
| Touch | `internal/mcp/instructions.go` (new), `server.go` ServerOptions, CONTRIBUTING, install optional |
| Tests | `TestServerInstructions*` + `TestToolNamesRegistered` |
| Out | 1-tool facade; bundled CG MCP; `trace_explore`; tool count change |

## Live-repo spot-check (P39-S01-00)

| Anchor | Status |
|--------|--------|
| `NewServer` no Instructions (`server.go:31–34`, `nil` ServerOptions) | confirmed — **implement wires here** |
| 16 MCP tools (`server.go` 16× `AddTool`; `RegisteredToolNames()`) | confirmed |
| `trace_version` tool (`server.go:153–161`) | confirmed |
| `trace_context` + optional `query` (`tools_context.go:21`; G1 S00) | confirmed shipped |
| `trace_loop` gate action (`tools_loop.go:18–24`) | confirmed |
| `CursorReloadTip` (`cursor.go:12–13`) | confirmed — optional strengthen |
| `PrintBootstrapHintIfNeeded` plan-only (`bootstrap_hint.go:12–36`) | confirmed — optional moat line |
| `InstallAgentDefaults` (`agents.go:17–57`) | confirmed |
| CONTRIBUTING MCP reload (`CONTRIBUTING.md:64–70`) | confirmed — expand orient |
| go-sdk v1.4.0 `ServerOptions.Instructions` | confirmed |
| `instructions.go` | **does not exist yet** — S01-01 creates |

## Next

**P39-S01-01** — implement G3 per thickened [01-implement.md](01-implement.md).
