#!/usr/bin/env bash
# E01 ab-incident-tracker — harness scoring (B0 / G1) + enforcement checks.
set -euo pipefail

EXP="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$EXP/../.." && pwd)"
TRACE_BIN="${TRACE_BIN:-$REPO/bin/trace}"

ARM=""
RUN_TESTS=0
CHECK_GATE=0
for arg in "$@"; do
  case "$arg" in
    B0|G1) ARM="$arg" ;;
    --test) RUN_TESTS=1 ;;
    --gate) CHECK_GATE=1 ;;
    -h|--help)
      echo "usage: $0 [B0|G1] [--test] [--gate]"
      exit 0
      ;;
    *)
      echo "unknown argument: $arg" >&2
      exit 2
      ;;
  esac
done

[[ -n "$ARM" ]] || { echo "usage: $0 [B0|G1] [--test] [--gate]" >&2; exit 2; }

WS="$EXP/runs/$ARM"
GT="$EXP/seed/gt.json"
[[ -d "$WS" ]] || { echo "missing workspace: $WS (run ./prepare.sh)" >&2; exit 2; }

FAILS=0
pass() { echo "PASS  $*"; }
fail() { echo "FAIL  $*"; FAILS=$((FAILS + 1)); }
skip() { echo "SKIP  $*"; }

echo "=== E01 score: arm $ARM ==="
echo "workspace: $WS"
echo

if [[ "$ARM" == "B0" ]]; then
  if [[ -d "$WS/.trace" ]]; then fail "B0-1 no .trace/"; else pass "B0-1 no .trace/"; fi
  if [[ -f "$WS/cmd/incidentd/main.go" ]] || [[ -x "$WS/bin/incidentd" ]]; then
    pass "B0-2 deliverable (cmd/incidentd or binary)"
  else
    fail "B0-2 deliverable (cmd/incidentd or binary)"
  fi
  if [[ "$RUN_TESTS" -eq 1 ]]; then
    if (cd "$WS" && CGO_ENABLED=1 go test ./...); then pass "B0-3 go test ./..."; else fail "B0-3 go test ./..."; fi
  else
    skip "B0-3 go test ./... — pass --test to run"
  fi
else
  # G1 harness
  [[ -d "$WS/.trace" ]] && pass "G1 .trace/ present" || fail "G1 .trace/ present"

  GRAPH="$WS/trace/graph.json"
  if [[ -f "$GRAPH" ]]; then pass "G2 graph export trace/graph.json"; else fail "G2 graph export trace/graph.json"; fi

  if [[ -f "$GRAPH" ]]; then
    goals=$(grep -c '"kind": "goal"' "$GRAPH" 2>/dev/null || echo 0)
    tasks=$(grep -c '"kind": "task"' "$GRAPH" 2>/dev/null || echo 0)
    decisions=$(grep -c '"kind": "decision"' "$GRAPH" 2>/dev/null || echo 0)
    if [[ "$goals" -ge 1 && "$tasks" -ge 3 && "$decisions" -ge 3 ]]; then
      pass "G3 graph non-trivial (goals>=1 tasks>=3 decisions>=3)"
    else
      fail "G3 graph non-trivial — goals=$goals tasks=$tasks decisions=$decisions"
    fi
    missing=0
    while read -r id; do
      grep -q "$id" "$GRAPH" || missing=$((missing + 1))
    done < <(grep -o '"id": "[^"]*"' "$GT" | grep '0000000000[0-9a-f][0-9a-f]' | cut -d'"' -f4)
    if [[ "$missing" -eq 0 ]]; then pass "G-seed task IDs match gt.json"; else fail "G-seed task IDs match gt.json — $missing missing"; fi
  else
    fail "G3 graph non-trivial — no graph.json"
    fail "G-seed task IDs match gt.json — no graph.json"
  fi

  ev=$(find "$WS" -maxdepth 1 -name 'TRACE-EVIDENCE*.md' 2>/dev/null | head -1)
  if [[ -n "$ev" ]] && grep -qi deliberation "$ev"; then
    pass "G4 TRACE-EVIDENCE* contains deliberation"
  else
    fail "G4 TRACE-EVIDENCE* contains deliberation — ${ev:-no TRACE-EVIDENCE*.md found}"
  fi

  task_rows=$("$TRACE_BIN" -C "$WS" tasks list 2>/dev/null | grep -c . || echo 0)
  if [[ "$task_rows" -ge 5 ]]; then pass "G5 trace tasks >= 5 rows"; else fail "G5 trace tasks >= 5 rows — count=$task_rows"; fi

  unc=$("$TRACE_BIN" -C "$WS" query 'SELECT COUNT(*) FROM uncertainties WHERE resolved_at IS NOT NULL' 2>/dev/null | tail -1 || echo 0)
  if [[ "${unc:-0}" -ge 1 ]]; then pass "G6 uncertainty with resolution in store"; else fail "G6 uncertainty with resolution — count=${unc:-0}"; fi

  if [[ -n "$ev" ]]; then
    if grep -q 'wave-0\|Wave 0\|WAVE-0' "$ev" 2>/dev/null; then pass "G7 wave-0 before wave-2"; else skip "G7 wave-0 before wave-2 — not found in evidence"; fi
  else
    skip "G7 wave-0 before wave-2 — no TRACE-EVIDENCE file"
  fi

  if "$TRACE_BIN" -C "$WS" loop status --task "e0100000-0000-4000-8000-000000000010" 2>&1 | grep -q deliberation; then
    pass "G-opt loop status deliberation block"
  else
    skip "G-opt loop status deliberation block"
  fi

  [[ -f "$WS/.cursor/rules/trace-enforcement.mdc" ]] && pass "E1 cursor rules installed" || fail "E1 cursor rules installed"
  [[ -f "$WS/.cursor/hooks/trace-loop-gate.sh" ]] && pass "E2 cursor-hook installed" || fail "E2 cursor-hook installed"
  if grep -q '"enforce"[[:space:]]*:[[:space:]]*"strict"' "$WS/.trace/config.json" 2>/dev/null; then
    pass "E3 config enforce strict"
  else
    fail "E3 config enforce strict"
  fi

  if [[ "$CHECK_GATE" -eq 1 ]]; then
    VTASK="e0100000-0000-4000-8000-000000000050"
    if "$TRACE_BIN" -C "$WS" loop gate --task "$VTASK" --for edit >/dev/null 2>&1; then
      pass "E4 gate allows edit on verify task"
    else
      fail "E4 gate allows edit on verify task — task=$VTASK still blocked"
    fi
  fi

  if [[ "$RUN_TESTS" -eq 1 ]]; then
    if (cd "$WS" && CGO_ENABLED=1 go test ./...); then pass "G-opt go test ./..."; else fail "G-opt go test ./..."; fi
  fi
fi

echo
if [[ "$FAILS" -eq 0 ]]; then
  echo "VERDICT: PASS"
  exit 0
else
  echo "VERDICT: FAIL ($FAILS check(s) failed — see RUBRIC.md)"
  exit 1
fi
