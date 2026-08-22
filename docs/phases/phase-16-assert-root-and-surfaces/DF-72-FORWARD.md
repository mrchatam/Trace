# Forward correction — DF-72 (thin MCP `trace_impact`)

P16-00 `00-PHASE-PLANNER.md` is **done history** and recorded DF-72 as P14 A3 defer. A subsequent human/planner cut (user-required **all** post-P15 findings) **supersedes that defer for upcoming scopes only**. Do not rewrite the P16-00 prompt body; do not claim P14/P15 close was wrong.

## P14 A3 supersession (Option A — one-paragraph lock)

**P14 A3** (“Impact upgrades wire into **existing** impact domain / `trace impact` — no new MCP impact tool menu”) is **superseded for one G19 adapter only**. Phase 16 S05 will register MCP tool `trace_impact` that wraps the same impact library the CLI already uses (`finding` / `alternative` / `report` / `walk`) — the same adapter pattern as `trace_tasks` / `trace_capability`. This is not a product daemon, not HTTP, not an install/decide/plan/index MCP dump, and not “MCP on the P0-X critical path” (P0-X already closed; thin stdio MCP adapters exist post-P10). Law 19 forbids forking business logic in the adapter; Law 13 forbids new infrastructure — a stdio tool calling existing domain code satisfies both. The import-boundary keeper must **allow** `trace_impact` and still **forbid** `trace_plan` / `trace_index`.

## Live lock (upcoming S05+)

| Item | Value |
|------|-------|
| DF-72 | **fix** — thin G19 MCP tool `trace_impact` (CLI `finding`/`alternative`/`report`/`walk` subset) |
| P14 A3 | **Superseded for this tool only** — not an install/decide/plan/index MCP dump; not a daemon |
| Home | **S05** (`scope-05-seed-impact-packet`); VERIFY (S06) named test is a **fail bar** |
| Keeper | Update `TestImportBoundaryMCPNoPlanImpactIndexTools` to **allow** `trace_impact`; still **forbid** `trace_plan` / `trace_index` |
| Catalog | Ten tools + `trace_version`; slug `mcp:trace_impact`; P15 Assert at tool entry |
| Do not | Rewrite P16-00 prompt body; claim P14/P15 history was wrong; add daemon/HTTP; add install/decide MCP |

Scope planners for **P16-S05-00** must treat this file + phase README as SoT over the historical DF-72 defer line in P16-00.
