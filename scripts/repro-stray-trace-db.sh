#!/usr/bin/env bash
# Durable dogfood: init → python root stub → CLI warn → live store under .trace/
# Does not delete/rename the root stub except via temp-dir cleanup on exit.
set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

pass() {
  echo "PASS: $*"
}

resolve_trace() {
  if [[ -n "${TRACE:-}" ]]; then
    echo "$TRACE"
    return
  fi
  if [[ -x "$REPO/bin/trace" ]]; then
    echo "$REPO/bin/trace"
    return
  fi
  mkdir -p "$REPO/bin"
  ( cd "$REPO" && go build -o "$REPO/bin/trace" ./cmd/trace )
  echo "$REPO/bin/trace"
}

TRACE_BIN="$(resolve_trace)"
[[ -x "$TRACE_BIN" ]] || fail "trace binary not executable: $TRACE_BIN"

echo "tmp=$TMP"
echo "trace=$TRACE_BIN"

"$TRACE_BIN" -C "$TMP" init

if [[ -e "$TMP/trace.db" ]]; then
  fail "init created root trace.db"
fi
pass "init did not create root trace.db"

if [[ ! -f "$TMP/.trace/trace.db" ]]; then
  fail "missing .trace/trace.db after init"
fi
pass ".trace/trace.db exists after init"

( cd "$TMP" && python3 -c "import sqlite3; sqlite3.connect('trace.db').close()" )
[[ -f "$TMP/trace.db" ]] || fail "python did not create root stub"

STAT_BEFORE="$(stat -c '%s %Y' "$TMP/trace.db")"

WARN_ERR="$(mktemp)"
trap 'rm -rf "$TMP"; rm -f "$WARN_ERR"' EXIT
set +e
"$TRACE_BIN" -C "$TMP" tasks >/dev/null 2>"$WARN_ERR"
set -e

STAT_AFTER="$(stat -c '%s %Y' "$TMP/trace.db")"
if [[ "$STAT_BEFORE" != "$STAT_AFTER" ]]; then
  fail "root stub size/mtime changed ($STAT_BEFORE → $STAT_AFTER)"
fi
pass "root stub untouched"

grep -Fq 'project-root trace.db exists but is not the Trace store' "$WARN_ERR" \
  || fail "warn missing substring: project-root trace.db exists but is not the Trace store"
grep -Fq '.trace/trace.db' "$WARN_ERR" \
  || fail "warn missing substring: .trace/trace.db"
grep -Fq 'agents: use CLI/MCP' "$WARN_ERR" \
  || fail "warn missing substring: agents: use CLI/MCP"
pass "warn observed on stderr"

if [[ ! -f "$TMP/.trace/trace.db" ]]; then
  fail "live store missing after tasks"
fi
pass "live store still .trace/trace.db"

echo "ALL PASS"
