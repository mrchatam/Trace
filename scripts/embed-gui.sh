#!/usr/bin/env bash
# embed-gui.sh — build the Explore SPA and sync it into internal/httpapi/embeddist
# for //go:embed (consumer binary serves real UI without project web/).
#
# Usage (from Trace repo root or via go:generate from internal/httpapi):
#   ./scripts/embed-gui.sh
#   go generate ./internal/httpapi
#
# The script always resolves the Trace repo root from its own path, so
# //go:generate ../../scripts/embed-gui.sh works when cwd is internal/httpapi/.
#
# Requires: Node.js + npm on PATH.
# Idempotent: re-running replaces SPA assets and regenerates embeddist/README.md.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WEB="$ROOT/web"
EMBED="$ROOT/internal/httpapi/embeddist"

fail() {
  echo "embed-gui: $*" >&2
  exit 1
}

command -v npm >/dev/null 2>&1 || fail "npm not found on PATH (install Node.js)"
command -v node >/dev/null 2>&1 || fail "node not found on PATH (install Node.js)"
[[ -f "$WEB/package.json" ]] || fail "missing $WEB/package.json"

echo "embed-gui: building web/ (npm ci && npm run build)…"
(
  cd "$WEB"
  npm ci
  npm run build
)

[[ -f "$WEB/dist/index.html" ]] || fail "web/dist/index.html missing after build"
if grep -q 'Embedded GUI stub' "$WEB/dist/index.html" 2>/dev/null; then
  fail "web/dist/index.html looks like a stub — refusing to embed"
fi
grep -q 'id="root"' "$WEB/dist/index.html" || fail "web/dist/index.html missing id=\"root\""
grep -q '/assets/' "$WEB/dist/index.html" || fail "web/dist/index.html missing /assets/"

echo "embed-gui: syncing web/dist → internal/httpapi/embeddist…"
mkdir -p "$EMBED"
# Remove previous SPA assets; keep directory. README is rewritten below.
find "$EMBED" -mindepth 1 -maxdepth 1 ! -name 'README.md' -exec rm -rf {} +
# Copy dist contents (not the dist directory itself).
cp -a "$WEB/dist"/. "$EMBED"/

# Regenerate README after sync so it is never overwritten by dist (and not left as two-artifact teaching).
cat >"$EMBED/README.md" <<'EOF'
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
EOF

echo "embed-gui: done → $EMBED"
ls -la "$EMBED" | head -20
