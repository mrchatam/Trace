#!/usr/bin/env bash
# build-with-proxy.sh — embed GUI and build trace binaries when proxy.golang.org returns 403.
#
# Usage (from repo root):
#   ./scripts/build-with-proxy.sh
#   SOCKS_PORT=10808 ./scripts/build-with-proxy.sh
#
# Sets ALL_PROXY / HTTPS_PROXY / HTTP_PROXY to socks5://127.0.0.1:${SOCKS_PORT:-10808}.
# Override host with SOCKS_HOST (default 127.0.0.1).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SOCKS_HOST="${SOCKS_HOST:-127.0.0.1}"
SOCKS_PORT="${SOCKS_PORT:-10808}"
PROXY_URL="socks5://${SOCKS_HOST}:${SOCKS_PORT}"

export ALL_PROXY="$PROXY_URL"
export HTTPS_PROXY="$PROXY_URL"
export HTTP_PROXY="$PROXY_URL"
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
export CGO_ENABLED=1

cd "$ROOT"
make embed-gui
go build -o bin/trace ./cmd/trace
go build -o bin/trace-mcp ./cmd/trace-mcp
echo "build-with-proxy: OK → $ROOT/bin/trace $ROOT/bin/trace-mcp"
