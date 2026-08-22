# Trace browser GUI (`web/`) — contributor DX

**Audience:** Trace **checkout** contributors. App/consumer repos do **not** need this tree — `trace gui` serves the Explore SPA from the Trace **binary** embed; consumer footprint is `.trace/` only. See [`docs/gui-quickstart.md`](../docs/gui-quickstart.md).

Operator SPA for `trace gui` (and `trace serve` headless/DX). Home is **Explore** (`/`): graph canvas + inspector. Overview and other ops routes are secondary (`/overview`, …).

## Stack

TypeScript + Vite + React · `react-router-dom` BrowserRouter · plain CSS variables · Context + `fetch` (no Redux / TanStack Query / Tailwind / MUI).

## Dev

```bash
# terminal 1 — primary launch opens the browser; use --no-open for Vite DX
trace gui --no-open --addr 127.0.0.1:7432
# or headless twin: trace serve --addr 127.0.0.1:7432

# terminal 2
cd web && npm install && npm run dev
# Vite proxies /v1 → http://127.0.0.1:7432
```

Pin `--addr` when proxying so Vite and the API share one host:port. Without `--addr`, default bind auto-picks the next free loopback port in **`7432`–`7441`** if `:7432` is busy — then update the Vite `/v1` proxy to match the printed URL (or keep pinning `:7432` for DX).

One `trace gui` / `trace serve` = one project root. A second root on default bind hops automatically; use `--addr` only when you want a **strict** pin (fail if busy).

## Production build (embed pipeline / disk override)

```bash
cd web && npm run build
# writes web/dist — disk override when present; also feeds scripts/embed-gui.sh
```

Consumers do not run this. Contributors refresh the binary embed with `./scripts/embed-gui.sh` (or `make embed-gui` / `go generate ./internal/httpapi`) then `go build` — see [`internal/httpapi/embeddist/README.md`](../internal/httpapi/embeddist/README.md).

## API types

```bash
npm run gen:api
# openapi-typescript → src/api/schema.d.ts
```

Hand wrappers live in `src/api/ops.ts` by OpenAPI `operationId`. Browser calls relative `/v1` only (Law 19).
