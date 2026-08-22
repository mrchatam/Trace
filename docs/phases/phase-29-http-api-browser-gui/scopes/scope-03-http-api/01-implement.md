# P29-S03-01 — HTTP API implementer

## Metadata
- id: P29-S03-01
- todo_ids: [P29-S03-01]
- role: implementer
- skills: [incremental-implementation, tdd, security-and-hardening]
- mcps: [user-context7]
- verification: automated
- hooks: []

## Objective

Implement **opt-in** `trace serve` + thin HTTP adapter `internal/httpapi` against [`api/openapi.yaml`](../../../../../api/openapi.yaml) and [`docs/adr/ADR-HTTP-API-GUI.md`](../../../../../adr/ADR-HTTP-API-GUI.md).

GUI ship (S04 vs S05) is **not** this row. S03 still implements the **HTTP contract** needed for both: **p0 + p1 live**, **defer = 501 stub**.

**No** React/`web/` SPA. **No** SQL in handlers (Law 19). **Do not** start P29-S03-02.

## References

- [00-PLANNER.md](00-PLANNER.md) — **final locked defaults**
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) §19
- [UX-IA.md](../scope-02-ux-ia/UX-IA.md) — which screens consume which ops (do not change IA)
- Adapter precedent: [`internal/mcp/`](../../../../../internal/mcp/) (open store → library → JSON; no domain fork)
- CLI flag/help patterns: [`cmd/trace/root.go`](../../../../../cmd/trace/root.go), [`cmd/trace/help.go`](../../../../../cmd/trace/help.go), [`cmd/trace/loop.go`](../../../../../cmd/trace/loop.go)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Human locks and S01–S02 locks are **settled** — do not re-grill bind/auth/CORS/MCP-RPC/package name.

## Locked defaults

| Item | Value |
|------|-------|
| Package | `internal/httpapi` |
| Type | `httpapi.Server` + `httpapi.Options` `{Root, Addr, AllowRemote, Token, TokenFile, StaticDir}` |
| CLI file | `cmd/trace/serve.go` |
| HTTP stack | stdlib `net/http` only (Go 1.22+ `ServeMux` method+path). **Do not** add chi/gin/echo/gorilla |
| Default listen | `127.0.0.1:7432` |
| Loopback | `127.0.0.1` and `::1` (and IPv4-mapped loopback). `0.0.0.0`, `::`, LAN/public IPs need `--allow-remote` |
| Project bind | Same as CLI: global `-C`/`--project` via existing `parseGlobal`. Serve may accept `--root` as alias; one root only |
| Static | `Options.StaticDir` default `<root>/web/dist`. If `index.html` missing, serve a **small placeholder** page at `GET /` (and SPA fallback for non-`/v1` GET) stating GUI ships in S04 |
| OpenAPI waves | **Implement** `x-trace-wave: p0` and `p1`. **501** `NOT_IMPLEMENTED` for every `defer` path (register the route; do not 404) |
| Schema residuals | Encode **library JSON as-is**. Optional: tighten `additionalProperties` schemas in `api/openapi.yaml` to match structs — **never invent HTTP fields** |
| CORS | No ACAO `*`. S03 default: omit CORS allow headers (same-origin). Do not implement Vite reflect here |
| Cloud headers | Ignore `X-Trace-Tenant` / `X-Trace-Request-Id` (no behavior) |
| Secrets | Never log bearer tokens or write them to error `details` |

### Package layout (implementer may split further, not fewer conceptual units)

```
internal/httpapi/
  server.go      // New, Handler(), ListenAndServe, mux register
  bind.go        // parse addr; IsLoopbackHost; RefuseRemote
  auth.go        // bearer middleware when !loopback
  cors.go        // explicit deny / no wildcard (testable)
  errors.go      // writeEnvelope(status, code, message, details)
  static.go      // dist or placeholder
  handlers_*.go  // by family is OK (meta, tasks, loop, entities, retrieval, seed, graph, p1, defer)
  *_test.go
cmd/trace/serve.go
cmd/trace/serve_test.go   // help flags; optional bind-refuse CLI
```

CLI `serve` **only** parses flags, enforces bind policy **before listen**, opens nothing domain-specific beyond passing `root` into `httpapi`. Token generate-and-print lives in CLI or `httpapi` helper — print **once** to stdout on remote serve when token was generated.

### Wave matrix (do not stub p0/p1)

**Implement (library-backed):**

| Wave | operationId |
|------|-------------|
| p0 | `getHealth`, `getVersion`, `getProject`, `listTasks`, `getTask`, `getLoopStatus`, `getLoopNext`, `postLoopApply`, `getLoopGate`, `postLoopReset`, `createEntity`, `getEntity`, `createLink`, `createTransition`, `getContext`, `getWhy`, `search`, `getSeedStatus`, `postSeedExport`, `postSeedImport`, `getGraph` |
| p1 | `listReviews`, `createReview`, `getReview`, `listPlans`, `listCapability`, `getImpact`, `listChanges`, `listRegressions`, `listAgents`, `getIndexStatus`, `postAuthToken` |

**Stub 501 `NOT_IMPLEMENTED`:** `postIndexReindex`, `postBackup`, `postRestore`, `postMigrate`, `postInstall`, `listPatterns`, `listKnowledge`, `getEval`, `getEventsSSE`.

### Library mapping (do not invent JSON)

Mirror MCP/CLI. Open store with `store.Open(absRoot)`; call packages below; `json.Marshal` the library types (or the same DTO the CLI already prints).

| HTTP | Call (canonical) |
|------|------------------|
| `GET /v1/health` | `{ok:true}` — process liveness, no DB required |
| `GET /v1/version` | `{ok, name:"trace", api_version:"1.0.0", trace_version}` — same binary version as `trace version` |
| `GET /v1/project` | Bound abs root + store open/readiness flags (no secrets, no extra filesystem dump) |
| `GET /v1/tasks` | `store.ListTasks` / `ListTasksByGoalID` → CLI `taskListRow` `{id,title,work_state,goal_id}` |
| `GET /v1/tasks/{id}` | Store/domain get task; 404 if missing |
| `GET /v1/loop/status` | `loop` status builder used by `cmdLoopStatus` / MCP `trace_loop` status — `schema_version` `trace.loop.status.v1` including `violations[]` |
| `GET /v1/loop/next` | `loop.BuildNextPacket` — `trace.loop.next.v1` |
| `POST /v1/loop/apply` | `loop.ParseApplyEnvelope` + `loop.Apply` — marshal apply result |
| `GET /v1/loop/gate` | `loop.EvaluateGate` + CLI `gateEnvelope` shape (`trace.loop.gate.v1`, `allowed`, `violations`) |
| `POST /v1/loop/reset` | `domain.Service.ResetDeliberationState` — marshal returned deliberation state (do not invent a parallel reset DTO) |
| `POST /v1/entities` | `domain.Service` create* by `kind` (same kinds as `trace add`) |
| `GET /v1/entities/{id}` | Store get by id/type as OpenAPI requires |
| `POST /v1/links` | `domain` link helpers (same rel vocabulary as `trace link` / MCP) |
| `POST /v1/transitions` | `domain.TransitionTask` — map gate/cap/review denials to 403/409 per ADR, not a second policy |
| `GET /v1/context` | `compiler` TaskContext / ExpandContext (MCP `trace_context`: depth 1\|2) |
| `GET /v1/why` | `retrieval.Engine.Why` |
| `GET /v1/search` | Same engine as `trace search` |
| seed * | `domain` seed import/export; HTTP body = **status/summary** (`SeedJobStatus`), **never** default full `graph.json`. Confine paths under project root (reject `..` / escape) — same honesty as CLI relative-path resolve |
| `GET /v1/graph` | **New** `internal/retrieval` (or store-backed) hop-limited neighborhood: require `center` + `max_nodes` (cap ≤ 5000). Missing params → 400 `VALIDATION_ERROR`. Over budget → 400 `BUDGET_EXCEEDED`. **No SQL in httpapi**. If a helper does not exist, add it in retrieval/store — not in the handler |
| reviews / capability / impact / changes / regressions / agents / plans / index status | CLI/MCP equivalents (`trace review`, `capability`, `impact`, `changes`, `regressions`, `agents`, `plan`, `index status`). Plans: `planner.Service` tree JSON the CLI already emits — do not invent a new plan document |
| `POST /v1/auth/token` | Local serve token rotate/issue (httpapi concern). Treat response `token` as secret |

If OpenAPI query/body names differ from library, **adapt at the edge only** (rename), keep value semantics identical.

### Bind / auth / CORS behavior

1. Parse `--addr`. If host is non-loopback and `--allow-remote` is unset → **refuse to listen** (CLI exit ≠ 0; message names the flag).
2. If listening non-loopback: require a token (`--token`, `--token-file`, or generate). Missing/invalid `Authorization: Bearer` → **401** `UNAUTHORIZED`.
3. Loopback: token optional; if `--token` was set, **do** enforce it (operator opted in).
4. CORS: tests assert response does **not** contain `Access-Control-Allow-Origin: *`.
5. `POST /rpc` and `/rpc` are **not** registered as MCP; if requested, 404 envelope (or 501) — never tools/call.

### Test strategy (TDD-friendly order)

Tests live primarily in `internal/httpapi`. Use `httptest` + `t.TempDir()` + `store.Open` like `internal/mcp/mcp_test.go`. Listen tests may use `127.0.0.1:0`.

**Required:**

| Area | Assertions |
|------|------------|
| Bind | Default/options treat loopback as allowed; `0.0.0.0` without `AllowRemote` fails **before** serve; with `AllowRemote` + token, listen OK |
| Auth | Loopback without token: `/v1/health` and a data GET succeed. Simulated remote (`RequireBearer` / non-loopback option): data GET without bearer → 401; with bearer → 200 |
| CORS | No `*` |
| RPC | `/rpc` not a successful MCP endpoint |
| Health | `GET /v1/health` → 200 `{ok:true}` |
| Reads (≥3) | e.g. `getVersion`, `listTasks`, `getLoopStatus` (fixture task) |
| Writes (≥2) | e.g. `createEntity` + `createLink` **or** `createTransition` — persist via library, visible on subsequent GET |
| Graph | Missing `max_nodes` or `center` → 400; bounded response when valid |
| Defer | One defer path returns 501 envelope `NOT_IMPLEMENTED` |
| Static | Placeholder (or dist) `GET /` is 200 text/html |
| CLI | `run([]string{"serve","--help"})` or `help` contains `serve`, `--addr`, `--allow-remote`, `--token` |
| Law 19 | Grep-level self-check: `internal/httpapi` has no `database/sql` / raw `sqlite` |

Seed path confinement: at least one test that `../` escape is rejected.

Do **not** require a live browser or Vite.

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f docs/adr/ADR-HTTP-API-GUI.md
test -f api/openapi.yaml
test ! -d internal/httpapi   # expected absent at start; if present, extend this contract — do not fork a second package
! grep -q 'case "serve"' cmd/trace/root.go || true
```

If `internal/httpapi` already exists from a partial attempt, **complete it** against this prompt; do not create `internal/http` / `internal/api` in parallel.

## Preflight / Plan

Short plan then execute. Slice vertically (each slice green tests):

1. Bind policy + Server mux + health/version + placeholder static + `trace serve` flags/help
2. Auth + CORS tests + project + tasks read
3. Entity/link/transition writes
4. Loop five ops (status/next/apply/gate/reset) mapped to library JSON
5. context/why/search + budgeted graph + seed status/export/import confinement
6. P1 families + auth token
7. Defer 501 stubs
8. `go test ./internal/httpapi/...` and `go test ./cmd/trace/...` (help)

## Role work

1. Implement package + CLI wiring per locked layout.
2. Register every OpenAPI path (p0/p1 live, defer 501).
3. Enforce bind/auth/CORS/no-rpc.
4. Tests as table above.
5. Board **P29-S03-01** status + Notes only (commands + key paths). Do not edit 02-review substance.

## Todo updates

Status + notes on **P29-S03-01** only.

## Exit criteria

- [ ] `go test ./internal/httpapi/...` PASS
- [ ] `go test ./cmd/trace/...` still PASS (help includes serve flags)
- [ ] `trace serve --help` (or `trace help`) documents `--addr`, `--allow-remote`, `--token`
- [ ] Default bind loopback `:7432`; remote refused without flag
- [ ] Bearer required off-loopback; loopback-trust unless token set
- [ ] No CORS `*`; no MCP `/rpc`
- [ ] p0 + p1 implemented via library; defer → 501 envelope
- [ ] Loop/gate/reset/next/apply JSON matches library (no invented fields)
- [ ] Graph requires center+max_nodes; seed paths confined; seed HTTP body not full graph
- [ ] No SQL in `internal/httpapi`
- [ ] Board Notes with test commands + files

## Minimal todos

- [ ] `internal/httpapi` server, bind, auth, errors, static placeholder
- [ ] `cmd/trace` `serve` + help + tests
- [ ] p0 routes + tests (health, ≥3 read, ≥2 write, bind, graph budget)
- [ ] p1 routes (reviews, plans, capability, impact, changes, regressions, agents, index GET, auth token)
- [ ] defer 501 stubs
- [ ] Notes on P29-S03-01

## Next

**P29-S03-02**
