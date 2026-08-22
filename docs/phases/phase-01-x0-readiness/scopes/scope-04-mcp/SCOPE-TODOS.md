# Scope S04 — MCP adapter

**Depends-on:** S01–S03 all `done` (DR-AGENT: MCP after CLI context validated).

**Locks (P01-S04-00):**
| Item | Value |
|------|-------|
| Layout | `internal/mcp` + thin `cmd/trace-mcp` |
| SDK | Official `github.com/modelcontextprotocol/go-sdk` (stdio only) |
| Required tools | `trace_why`, `trace_context` |
| Parity tools | `trace_add`, `trace_link`, `trace_transition`, `trace_review` |
| G19 | Library-only; no domain fork; no reverse imports into MCP |
| Out | Daemon/HTTP primary; embeddings; raw SQL; init/index/seed as MCP tools |
| X0 | Remains CLI-capable without MCP (`evals/x0` B0/G1) |

**S03 note:** `evals/x0` owns B0/G1 dry-run metrics. MCP must not gate or replace that path; regression keeps `./evals/x0/...` green.

- [x] P01-S04-00 planner
- [x] P01-S04-01 implement
- [x] P01-S04-02 review (+ spawns as needed)
