# Scope 00 — board map (G1)

**S00 context orient merge** — G-001, G-002. Serial: **P39-S00-00 → P39-S00-01 → P39-S00-02**.

| Order | Board ID | Prompt | Role | Status |
|------:|----------|--------|------|--------|
| 672 | P39-S00-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | **done** |
| 673 | P39-S00-01 | [01-implement.md](01-implement.md) | Implementer | pending |
| 674 | P39-S00-02 | [02-review.md](02-review.md) | Reviewer | pending |

## Planner locks (P39-S00-00 — verified 2026-08-22)

| Lock | Value |
|------|-------|
| Theme | G1 query+task orient merge |
| Verdict | **Accept** — optional query merges into packet |
| Moat | M-001 — task_id required; Layer 0 preserved |
| Query semantics | Adds FTS hits; title FTS still runs when query set |
| Merge point | `compileAtDepth` after title FTS (~154), before file-seed expand (~156) |
| Touch | `compiler.go`, `context.go`, `tools_context.go`, `server.go` desc, tests |
| Out | G2 unified explore; query-only; cap raise; new MCP tool |
| MCP tools | **16** unchanged (`RegisteredToolNames()` verified) |
| Tests | T1–T6 + T1-MCP in `01-implement.md`; regression list in `02-review.md` |

## Live-repo spot-check (P39-S00-00)

| Anchor | Status |
|--------|--------|
| `ContextOptions` no Query (`compiler.go:14–19`) | confirmed |
| Title FTS only (`compiler.go:146–154`, Limit 16) | confirmed (line +3 vs P39-00) |
| Caps 4096/32/64 (`packet.go:18–20`) | confirmed |
| Reason codes (`types.go:13–14`) | confirmed |
| CLI no `--query` (`context.go:18–68`) | confirmed |
| MCP no `query` (`tools_context.go:14–21`; live schema) | confirmed |
| 16 MCP tools (`server.go` 16× `AddTool`) | confirmed |
| Codegraph | no `.codegraph/` — used Read/Grep |

## Next

**P39-S00-01** — implement G1 per thickened [01-implement.md](01-implement.md).
