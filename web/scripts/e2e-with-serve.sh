#!/usr/bin/env bash
# Start trace serve with web/dist, run Playwright e2e, tear down.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
WEB="$ROOT/web"
PORT="${TRACE_E2E_PORT:-7439}"
BASE="http://127.0.0.1:${PORT}"

cd "$WEB"
npm run build

if [[ ! -x "$ROOT/bin/trace" ]]; then
  (cd "$ROOT" && go build -o bin/trace ./cmd/trace)
fi

# Serve from project root so .trace/ + web/dist resolve correctly.
(
  cd "$ROOT"
  exec "$ROOT/bin/trace" serve --addr "127.0.0.1:${PORT}"
) &
PID=$!
cleanup() { kill "$PID" 2>/dev/null || true; wait "$PID" 2>/dev/null || true; }
trap cleanup EXIT

for i in $(seq 1 60); do
  if curl -sf "$BASE/v1/health" >/dev/null; then
    break
  fi
  sleep 0.25
done
curl -sf "$BASE/v1/health" >/dev/null

export TRACE_E2E_BASE="$BASE"
cd "$WEB"
npx playwright test "$@"
