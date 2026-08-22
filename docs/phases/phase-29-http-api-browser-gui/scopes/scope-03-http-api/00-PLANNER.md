# P29-S03-00 — Scope planner (HTTP API)

## Metadata
- id: P29-S03-00
- todo_ids: [P29-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, security-and-hardening]
- verification: automated

## Objective

Finalize S03 implement/review prompts for OpenAPI-backed HTTP adapter + `trace serve`. Handlers call library only.

## Session start

Follow agent-loop-protocol Session start.

## Locked defaults (final — honor in 01/02)

| Item | Value |
|------|-------|
| Package | `internal/httpapi` |
| CLI | `cmd/trace/serve.go` + `case "serve"` in `root.go`; flags only — no domain logic |
| Router | stdlib `net/http` (Go 1.22+ method mux). **No** new HTTP frameworks |
| Handlers | Parse/validate → canonical library (`domain`, `loop`, `compiler`, `retrieval`, `planner`, `store` reads) → JSON. **No SQL in handlers** |
| Bind | Default `127.0.0.1:7432` (OpenAPI `servers` port). Non-loopback / `0.0.0.0` / `::` refused unless `--allow-remote` |
| Auth | Loopback-trust (no bearer). Non-loopback: bearer **required**. Generate-and-print if remote and no `--token`/`--token-file` |
| Prefix | `/v1` |
| Flags | Global `-C`/`--project`; serve `--addr`, `--allow-remote`, `--token`, `--token-file`; optional `--root` alias of project bind |
| CORS | Deny (no `Access-Control-Allow-Origin` for foreign origins; **never** `*`). Vite-dev origin reflect → **S06**, not S03 |
| MCP | **No** `POST /rpc`. 404 + envelope if hit |
| Static | Serve `web/dist` if `index.html` exists; else **inline placeholder HTML** at `/`. Do **not** scaffold the React SPA (`web/` is S04) |
| Waves | **Implement** all OpenAPI `x-trace-wave: p0` **and** `p1`. **Stub 501** all `defer` (including `POST /v1/index`) |
| Contract | `api/openapi.yaml` SoT for routes. Map loop/context/why/plans JSON from library — **do not invent fields** |
| Errors | ADR envelope `{error:{code,message,details}}`. Map library deny → 403; not found → 404; conflict → 409; graph over budget → 400 `BUDGET_EXCEEDED` |
| Graph | Budgeted neighborhood in `internal/retrieval` (new helper OK). Require `center` + `max_nodes`. No full dump |
| Seed paths | Confine `input_path`/`output_path` under project root |

## Exit criteria (scope — verified by S03-01/02)

- `go test ./internal/httpapi/...` PASS
- CLI `trace serve --help` documents flags
- Integration tests: health + ≥3 read + ≥2 write paths
- Bind + bearer-off-loopback + CORS deny + no `/rpc` tests

## Next

P29-S03-01 → P29-S03-02
