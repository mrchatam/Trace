# P29-S03-02 — HTTP API review

## Metadata
- id: P29-S03-02
- todo_ids: [P29-S03-02]
- role: reviewer
- skills: [code-review-and-quality, security-and-hardening, silent-failure-hunter]
- verification: automated
- hooks: []

## Objective

Independent review of `internal/httpapi` + `trace serve` vs OpenAPI, ADR, Law 19, and S03-00 locks. Fresh subagent. Small inline fixes OK; structural gaps spawn `P29-S03-02a` / `02b`.

## Session start

Follow agent-loop-protocol Session start. Do not share the implementer session. Do not start S04.

## References

- [00-PLANNER.md](00-PLANNER.md) locked defaults
- [01-implement.md](01-implement.md) wave matrix + library mapping
- [`api/openapi.yaml`](../../../../../api/openapi.yaml)
- [`docs/adr/ADR-HTTP-API-GUI.md`](../../../../../adr/ADR-HTTP-API-GUI.md)
- Implementer board Notes on **P29-S03-01**

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -d internal/httpapi
grep -q 'case "serve"' cmd/trace/root.go
go test ./internal/httpapi/...
```

If package/CLI missing, **fail** the row (or spawn 02a) — do not implement the whole API in review.

## Checklist

### Contract

- [ ] Every `x-trace-wave: p0` and `p1` operationId is a real library-backed handler (not silent 501)
- [ ] Every `defer` operation is **501** `NOT_IMPLEMENTED` envelope (not 404, not a fake success)
- [ ] Prefix `/v1`; health/version shapes match OpenAPI
- [ ] Loop status/next/apply/gate/reset bodies are library JSON (`trace.loop.*.v1` / gate envelope / deliberation reset) — **no invented fields**
- [ ] Plans/capability/impact/changes/regressions/agents/context/why match CLI/MCP JSON, not a parallel DTO
- [ ] `GET /v1/graph` requires `center` + `max_nodes`; no unbounded dump
- [ ] Seed HTTP returns status/summary; paths confined to project root
- [ ] No `POST /rpc` / MCP tools/call browser transport

### Law 19

- [ ] Handlers call `domain` / `loop` / `compiler` / `retrieval` / `planner` / `store` APIs
- [ ] **No** `database/sql` or SQLite DSN usage under `internal/httpapi`
- [ ] CLI `serve` is flags + listen + pass root; no duplicated transition/loop policy
- [ ] New graph/seed helpers (if any) live in library packages, not copied business rules in HTTP

### Security (human locks)

- [ ] Default listen `127.0.0.1:7432` (or documented equivalent default host loopback + that port)
- [ ] Non-loopback / `0.0.0.0` refused without `--allow-remote` (test or CLI evidence)
- [ ] Bearer required when bind is non-loopback; 401 envelope `UNAUTHORIZED`
- [ ] CORS: no `Access-Control-Allow-Origin: *`
- [ ] Tokens not logged; errors have no stack traces / accidental abs paths outside intentional fields
- [ ] Placeholder/static cannot list `.trace/` or sqlite

### Tests / UX-IA contract

- [ ] `go test ./internal/httpapi/...` PASS
- [ ] Floor: bind + health + ≥3 read + ≥2 write
- [ ] Help documents `--addr`, `--allow-remote`, `--token`
- [ ] Static placeholder (or `web/dist`) served; **no** accidental full SPA in this scope (S04 owns `web/`)
- [ ] GUI ship ≠ API wave: p1 HTTP exists even if UX-IA keeps reviews/graph-rich for S05

### Error model

- [ ] Envelope `{error:{code,message,details}}`
- [ ] Sensible mapping: 400 validation/budget, 401 auth, 403 library forbid, 404, 409 apply/state, 501 defer

## Findings

Classify: blocker | high | medium | low | nit. Cite files.

## Spawn policy

- blocker/high: inline fix if small; else insert **P29-S03-02a** (implement) + **P29-S03-02b** (review) immediately below this row with full protocol prompts
- medium: prefer spawn unless trivial
- Do not rewrite P29-S03-01 prompt; do not edit done S00–S02 history
- Upcoming only: if S04/S06 locks need a one-line correction (placeholder vs embed, Vite CORS), thicken **those** planners

## Exit criteria

- [ ] No open blocker/high without a pending follow-up row
- [ ] Confidence **medium** or **high** with evidence (commands, greps)
- [ ] Board Notes on **P29-S03-02**; next **P29-S04-00**

## Next

**P29-S04-00**
