# P01-S04-02 — Scope review notes (MCP adapter)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Reviewer:** independent (fresh context; claims checked against repo + re-run tests)

## Claims vs evidence

| Claim (P01-S04-01 Notes / 01 prompt) | Evidence |
|--------------------------------------|----------|
| Layout `internal/mcp` + thin `cmd/trace-mcp` | Tree matches; `main.go` only parses `-C`/`--project` then `RunStdio` |
| Official `github.com/modelcontextprotocol/go-sdk` v1.4.0 | `go.mod` require; no mark3labs/mcp-go |
| Transport stdio only | `RunStdio` → `&sdkmcp.StdioTransport{}`; no ListenAndServe/SSE/HTTP in mcp packages |
| Server identity `trace` / `0.0.0-dev` | `server.go` constants match CLI version string |
| Tools: `trace_why` `trace_context` `trace_add` `trace_link` `trace_transition` `trace_review` | All registered in `registerTools`; help text lists same six; **no deferral** |
| Required why/context → library | `tools_why.go` → `retrieval.Why`; `tools_context.go` → `compiler.TaskContext`/`ExpandContext` (depth 1\|2; format json\|markdown\|both) |
| Parity writes → domain | `tools_write.go` → Create*/Link*/TransitionTask/CreateReview/SetReviewResult/LinkReviewTask; DONE gate unchanged |
| G19: no reverse imports | `rg` empty under `internal/{store,vcs,gitcli,analyzers,domain,retrieval,compiler}` + `cmd/trace`; `TestImportBoundaryNoLibraryImportsMCP` PASS |
| Project root `-C`/`--project` + per-tool `project` | `cmd/trace-mcp` flags + `Project` on each input struct; `resolveProject` Abs |
| X0 without MCP (DR-AGENT) | `evals/x0` has no mcp import; dry-run still CLI why/context |
| No daemon/HTTP/SQL/embeddings | Grep clean on mcp packages; store opens via `store.Open` only |
| Tests / build / regression | `CGO_ENABLED=0 go test ./internal/mcp/...` PASS; `go build -o bin/trace-mcp` PASS; `CGO_ENABLED=1 go test ./evals/x0/... ./evals/p0x/... ./evals/honesty/... ./...` PASS |

## Checklist (02-scope-review)

### Thin adapter / G19 — PASS
- Handlers open store and call domain/retrieval/compiler; no duplicated DONE/why/context policy
- CLI semantics mirrored (kinds, rels, depth/format, allow_done, review create|set)
- Import edges: no library → mcp reverse deps

### Tool list vs CLI — PASS
- Required + all four parity tools shipped (snake_case args documented in jsonschema tags)
- Out of scope tools absent (no init/index/seed/SQL)

### Transport / surfaces — PASS
- Stdio primary; no HTTP/SSE daemon product surface
- MCP not required for X0

### Regression — PASS
- Independent re-run: x0 + p0x + honesty + `./...` green (CGO=1)

### Cross-scope — PASS
- S05 VERIFY MCP checklist thickened with live six-tool list + binary/SDK pin
- Did not rewrite S04-01 `done` prompt body

## Findings

| Sev | Finding | Disposition |
|-----|---------|-------------|
| nit | `TestToolNamesRegistered` is constructor smoke only (does not assert the six names) | Accept residual — names covered by help + registerTools + handler tests |
| nit | `tryOpenGit`/`isNotRepo` dead-branch pattern copied from CLI | Accept residual (identical to `cmd/trace`) |
| nit | `SCOPE-TODOS.md` implement/review boxes lagged | **Fixed inline** — synced |

**Blocker/high:** none  
**Spawns:** none

## Residuals (none material)

- Proxy note: first `go build` hit `403` on `proxy.golang.org` for `segmentio/encoding`; succeeded with `GOPROXY=direct` (env/network, not product defect).
- `export_test.go` test seams (`CallWhy`, …) are package-test-only; public surface remains Options/NewServer/RunStdio + input structs.

## Next

**P01-S05-00** (Phase 01 VERIFY planner)
