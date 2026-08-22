# Trace GUI quickstart (`trace gui`)

Local-first **opt-in** HTTP API and browser GUI. Not a daemon, not a hosted SaaS product.

## Launch (primary)

From any Trace-initialized **app** repo (consumer footprint is **`.trace/` only** — no project `web/` required):

```bash
trace gui
```

Opens the browser to Explore home. Default bind is `http://127.0.0.1:7432/`. Same loopback HTTP server as `trace serve`; `gui` adds open-browser after listen. The Explore SPA is served from the Trace **binary** embed when disk `web/dist` is absent.

## Install CLI on PATH

This module is not published to the Go module proxy yet, so `go install …@latest` fails with **no matching versions**. From a Trace checkout:

```bash
cd /path/to/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
cp -f bin/trace ~/.local/bin/trace   # or any dir on your PATH
# then, from a Trace-initialized repo:
trace gui
```

If `~/.local/bin/trace` is root-owned (common after an earlier sudo install), overwrite with your user: `cp -f bin/trace ~/.local/bin/trace` (may need to fix ownership once).

`trace install …` configures agents/MCP/hooks only — it does **not** put the `trace` binary on PATH.

## Headless / scripting (`serve`)

Use when you need headless bind, CI, or no browser helper. `trace serve` is the twin of `trace gui --no-open`.

```bash
trace serve
# listening on http://127.0.0.1:7432  (or next free in 7432–7441 if busy)
```

Open the printed URL yourself. No app-repo SPA build step.

## What you see (Explore home)

Home is **Explore** (`/`): a canvas-first graph shell plus an inspector (summary → why → context → impact → reviews → links). Select a node to load detail; double-click or “Use as center” to re-center. Overview and other ops screens live under `/overview` and sibling routes — useful, but secondary to the graph home.

Craft cues (live CSS): `.graph-shell` with a taller canvas (`--graph-canvas-height`), inspector roughly `minmax(18rem, 26rem)`, structured `PacketView` fields with raw JSON in `<details>`, calm node chrome (no glow), and short settle/node/chrome motions that honor `prefers-reduced-motion`. Kind literacy: each node’s left chroma strip plus text labels (kind/state) — not color alone.

## Multi-project / ports

One `trace gui` / `trace serve` process binds **one** project root (`--root` / `-C`). Default listen address is `127.0.0.1:7432`. If that port is busy, default bind **auto-picks** the next free loopback port in **`7432`–`7441`**, prints the chosen URL, and (`gui`) opens that URL. You do not need `--addr` for a second project.

```bash
# Project A
trace gui
# → http://127.0.0.1:7432/  (or next free if busy)

# Project B (another cwd / -C) — auto-hops; opens the chosen URL
trace gui -C /path/to/other
# → http://127.0.0.1:7433/ … up to :7441

# Pin (strict): set --addr on the cmdline → fail if busy (no hop)
trace gui --addr 127.0.0.1:7433

# Headless twin
trace serve -C /path/to/other
# or: trace gui --no-open …
```

If you **pin** with `--addr` and that address is taken:

```text
gui|serve: address already in use (127.0.0.1:7432)
hint: another process (often trace gui or trace serve) is bound there.
  Default bind auto-hops to the next free port (7432–7441); --addr pins and fails if busy.
  To pin a free port, e.g.:
    trace gui -C /path/to/other --addr 127.0.0.1:7433
```

If the whole auto range is full:

```text
gui|serve: no free port in auto range 127.0.0.1:7432–7441 (10 ports tried)
hint: stop other trace gui/serve processes, or pin a free port with --addr:
  trace gui --addr 127.0.0.1:7450
```

Vite DX (Trace checkout) still proxies `/v1` to `127.0.0.1:7432` by default — if you pin a different bind, point the proxy (or `--addr`) at the same address. Historical design notes: [OPEN-PORT-MULTI.md](phases/phase-32-graph-first-gui/OPEN-PORT-MULTI.md) (superseded for default bind by Phase 34 auto free-port).

## Static assets (resolution)

| Source | When |
|--------|------|
| Disk `<root>/web/dist` | Prefer when `index.html` is present (contributor / Trace checkout override) |
| Embedded SPA | Everyday consumer path — baked into the Trace binary |
| Placeholder HTML | Packaging mistake (empty/broken embed) — not the consumer story |

Default `--static-dir` string is still `<root>/web/dist`, but consumers rarely need the flag: resolution is disk-if-present → embed → placeholder. **Do not** set `--static-dir` to the project root — refused (would expose `.trace/` and source).

## Contributor: build SPA / refresh embed

Only for Trace **checkout** work (not app repos):

```bash
cd web && npm ci && npm run build   # disk override while iterating
./scripts/embed-gui.sh              # or: make embed-gui / go generate ./internal/httpapi
go build -o bin/trace ./cmd/trace   # pick up go:embed
```

See [`web/README.md`](../web/README.md) and [`internal/httpapi/embeddist/README.md`](../internal/httpapi/embeddist/README.md).

## Security defaults

| Control | Behavior |
|---------|----------|
| Bind | Default `127.0.0.1:7432` (loopback-trust); busy → next free in `7432`–`7441` |
| `--addr` | Pin: fail if busy (no auto-hop) |
| Remote | `--allow-remote` plus bearer (`--token` / `--token-file` / one-time generated) |
| CORS | Deny by default; **never** `Access-Control-Allow-Origin: *` |
| Vite DX | Prefer Vite **proxy** `/v1` → `:7432`. Optional `--cors-origin http://127.0.0.1:5173` reflects **only** that exact Origin |
| Static dir | Default `<root>/web/dist` path string; resolution disk → embed → placeholder. Never project root |
| CSP | Set on static/SPA responses (`default-src 'self'`, `frame-ancestors 'none'`, …) |
| MCP | No browser `/rpc` transport — use local stdio MCP |

Example remote (explicit):

```bash
trace serve --addr 0.0.0.0:7432 --allow-remote --token "$TRACE_TOKEN"
```

## Seed honesty

`POST /v1/seed/export` with `strict` or `task_id` returns **501** `NOT_IMPLEMENTED`. Use the CLI for honesty/gate export (`trace seed export …`). HTTP export writes a summary + file under the project root only (path confined).

## Graph budgets

`GET /v1/graph` requires `center` and `max_nodes` — there is no unbounded full-graph dump. Search / changes / regressions use `limit` caps. Task lists are project-local (library has no paging yet); keep projects sized for local use.

## Cloud path

Same OpenAPI contract may later back a hosted product. Phase 29 does **not** ship multi-tenant hosting, OAuth, or public-internet defaults. See [CLOUD-APPENDIX.md](phases/phase-29-http-api-browser-gui/CLOUD-APPENDIX.md) (design only).

## Related

- ADR: [`docs/adr/ADR-HTTP-API-GUI.md`](adr/ADR-HTTP-API-GUI.md)
- OpenAPI: [`api/openapi.yaml`](../api/openapi.yaml)
- Ports design (historical): [`OPEN-PORT-MULTI.md`](phases/phase-32-graph-first-gui/OPEN-PORT-MULTI.md)
- `trace gui --help` — launch + open-browser flags (`--no-open`, shared serve flags)
- `trace serve --help` — same bind/CORS/static flags for headless / scripting
