#!/usr/bin/env bash
# E01 — Phase 23 enforcement mechanics (no CMS build).
set -euo pipefail

EXP="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$EXP/../.." && pwd)"
TRACE_BIN="${TRACE_BIN:-$REPO/bin/trace}"
WS="$EXP/runs/G1"
TASK="e0100000-0000-4000-8000-000000000010"

[[ -d "$WS/.trace" ]] || { echo "run ./prepare.sh first" >&2; exit 2; }

FAILS=0
pass() { echo "PASS  $*"; }
fail() { echo "FAIL  $*"; FAILS=$((FAILS + 1)); }

echo "=== E01 enforcement demo ==="
echo "workspace: $WS"
echo "task: $TASK"
echo

[[ -f "$WS/.cursor/rules/trace-enforcement.mdc" ]] && pass "cursor rules installed" || fail "cursor rules installed"
[[ -f "$WS/.cursor/hooks/trace-loop-gate.sh" ]] && pass "cursor-hook installed" || fail "cursor-hook installed"
grep -q '"enforce"[[:space:]]*:[[:space:]]*"strict"' "$WS/.trace/config.json" && pass "config enforce strict" || fail "config enforce strict"

if ! "$TRACE_BIN" -C "$WS" loop gate --task "$TASK" --for edit >/tmp/e01-demo-gate.json 2>/dev/null; then
  grep -q plan_missing /tmp/e01-demo-gate.json && pass "gate --for edit blocks (plan_missing)" || fail "gate --for edit blocks (plan_missing)"
else
  fail "gate --for edit blocks (plan_missing)"
fi

if "$TRACE_BIN" -C "$WS" loop status --task "$TASK" 2>&1 | grep -q '"violations"'; then
  "$TRACE_BIN" -C "$WS" loop status --task "$TASK" 2>&1 | grep -q '"violations":\[' && pass "status violations[] populated" || fail "status violations[] populated"
else
  fail "status violations[] populated"
fi

if ! "$TRACE_BIN" -C "$WS" seed export -o /tmp/e01-export.json --strict --enforce 2>/dev/null; then
  pass "export --strict --enforce blocks write"
else
  fail "export --strict --enforce blocks write"
fi

if ! "$TRACE_BIN" -C "$WS" transition --task "$TASK" --to DONE --enforce 2>/dev/null; then
  pass "transition --enforce blocks DONE"
else
  fail "transition --enforce blocks DONE"
fi

echo
if [[ "$FAILS" -eq 0 ]]; then
  echo "VERDICT: PASS (E01 enforcement mechanics demo)"
  exit 0
else
  echo "VERDICT: FAIL ($FAILS)"
  exit 1
fi
