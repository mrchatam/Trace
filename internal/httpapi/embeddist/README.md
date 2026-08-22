# Embedded Explore SPA

This directory is the `go:embed` source for the Trace browser GUI
(`internal/httpapi/embed.go`). Release and everyday consumer binaries serve
**this** tree when disk `<root>/web/dist/index.html` is absent.

Consumer projects need only `.trace/` — they must not require a project `web/`.

## Refresh pipeline

From the Trace repo root:

```bash
./scripts/embed-gui.sh
# or: go generate ./internal/httpapi
# or: make embed-gui
```

That builds `web/` (`npm ci && npm run build`), syncs `web/dist/**` here, and
rewrites this README. Then `go build` so `//go:embed` picks up the assets.

## Disk override (contributor DX)

If `<root>/web/dist/index.html` exists, `trace serve` / `trace gui` prefer disk
over embed (Vite contributor path). Optional `--static-dir` still applies;
StaticDir equal to the project root is refused.

## Last resort

If embed is empty or broken (pipeline not run), the HTTP adapter falls back to
an inline placeholder page. That is a packaging mistake for releases — not the
consumer story.
