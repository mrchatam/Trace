#!/usr/bin/env bash
# E03 ab-library-hold-desk — harness scoring + Phase 25 checks.
# Strict enforce (T02): --strict --enforce on both arms; thin graph fails G2 (FM-07 stays warn-only).
set -euo pipefail

EXP="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$EXP/../.." && pwd)"
TRACE_BIN="${TRACE_BIN:-$REPO/bin/trace}"

ARM=""
RUN_TESTS=0
CHECK_P25=0
P25_ARM="build"
ARGS=("$@")
i=0
while [[ $i -lt ${#ARGS[@]} ]]; do
  arg="${ARGS[$i]}"
  case "$arg" in
    B0|G1) ARM="$arg" ;;
    --test) RUN_TESTS=1 ;;
    --p25) CHECK_P25=1 ;;
    --arm)
      i=$((i + 1))
      if [[ $i -ge ${#ARGS[@]} ]]; then
        echo "error: --arm requires build or directed" >&2
        exit 2
      fi
      P25_ARM="${ARGS[$i]}"
      ;;
    -h|--help)
      echo "usage: $0 [B0|G1] [--test] [--p25] [--arm build|directed]"
      exit 0
      ;;
    *)
      echo "unknown argument: $arg" >&2
      exit 2
      ;;
  esac
  i=$((i + 1))
done

[[ -n "$ARM" ]] || { echo "usage: $0 [B0|G1] [--test] [--p25] [--arm build|directed]" >&2; exit 2; }

if [[ "$CHECK_P25" -eq 1 ]]; then
  case "$P25_ARM" in
    build|directed) ;;
    *)
      echo "error: --arm must be build or directed (got: $P25_ARM)" >&2
      exit 2
      ;;
  esac
fi

WS="$EXP/runs/$ARM"
GT="$EXP/seed/gt.json"
[[ -d "$WS" ]] || { echo "missing workspace: $WS (run ./prepare.sh)" >&2; exit 2; }

FAILS=0
pass() { echo "PASS  $*"; }
fail() { echo "FAIL  $*"; FAILS=$((FAILS + 1)); }
skip() { echo "SKIP  $*"; }

count_kind() {
  local file="$1" kind="$2"
  python3 - "$file" "$kind" <<'PY' 2>/dev/null || echo 0
import json, sys
path, kind = sys.argv[1], sys.argv[2]
with open(path) as f:
    data = json.load(f)
# Support both export shapes:
# 1) legacy/graph-store shape: entities[] with {"kind": "..."}
# 2) current seed export shape: top-level arrays (goals/tasks/decisions/discoveries/...)
entities = data.get("entities")
if isinstance(entities, list):
    print(sum(1 for e in entities if isinstance(e, dict) and e.get("kind") == kind))
else:
    plural_map = {
        "discovery": "discoveries",
    }
    bucket = data.get(plural_map.get(kind, kind + "s"))
    print(len(bucket) if isinstance(bucket, list) else 0)
PY
}

echo "=== E03 score: arm $ARM ==="
echo "workspace: $WS"
[[ "$CHECK_P25" -eq 1 ]] && echo "p25 arm: $P25_ARM"
echo

if [[ "$ARM" == "B0" ]]; then
  if [[ -d "$WS/.trace" ]]; then fail "B0-1 no .trace/"; else pass "B0-1 no .trace/"; fi
  if [[ -f "$WS/cmd/holddeskd/main.go" ]] && [[ $(find "$WS/internal" -name '*.go' 2>/dev/null | wc -l) -gt 0 ]]; then
    pass "B0-2 deliverable beyond starter"
  elif [[ -f "$WS/cmd/holddeskd/main.go" ]] && [[ -f "$WS/SPEC.md" ]]; then
    pass "B0-2 deliverable (SPEC + cmd)"
  else
    fail "B0-2 deliverable beyond starter"
  fi
  if [[ "$RUN_TESTS" -eq 1 ]]; then
    if (cd "$WS" && CGO_ENABLED=1 go test ./...); then pass "B0-3 go test ./..."; else fail "B0-3 go test ./..."; fi
  else
    skip "B0-3 go test ./... — pass --test"
  fi
else
  [[ -d "$WS/.trace" ]] && pass "G1 .trace/ present" || fail "G1 .trace/ present"

  GRAPH="$WS/trace/graph.json"

  # T01: preflight export when graph.json missing (respect existing agent/operator export)
  if [[ ! -f "$GRAPH" ]]; then
    if ! "$TRACE_BIN" -C "$WS" seed export -o trace/graph.json 2>/dev/null; then
      fail "G2 graph export — preflight export failed (run: $TRACE_BIN -C $WS seed export -o trace/graph.json)"
    else
      pass "G2 graph export — preflight export created"
    fi
  else
    pass "G2 graph export"
  fi

  # T02: --strict --enforce on both arms; honesty failures fail G2 (not WARN-only)
  if [[ -f "$GRAPH" ]]; then
    enforce_err=""
    enforce_rc=0
    enforce_err="$("$TRACE_BIN" -C "$WS" seed export -o trace/graph.json --strict --enforce 2>&1 >/dev/null)" || enforce_rc=$?
    if [[ "$enforce_rc" -ne 0 ]]; then
      fail "G2 graph honesty --strict --enforce (exit $enforce_rc)"
      while IFS= read -r line; do
        [[ -n "$line" ]] && echo "      enforce: $line"
      done <<< "$enforce_err"
    else
      pass "G2 graph honesty --strict --enforce"
    fi
  fi

  # T03: FM-07 git-sparsity warn-only (never fail G2 on drift alone)
  if [[ -f "$GRAPH" ]]; then
    export_sha=$(python3 - "$GRAPH" <<'PY' 2>/dev/null || echo ""
import json, sys
try:
    with open(sys.argv[1]) as f:
        d = json.load(f)
    val = d.get("exported_at_commit", "")
    print(val if val else "")
except Exception:
    print("")
PY
)
    if [[ -n "$export_sha" ]] && git -C "$WS" rev-parse HEAD >/dev/null 2>&1; then
      head_sha=$(git -C "$WS" rev-parse HEAD)
      if [[ "$export_sha" != "$head_sha" ]]; then
        echo "WARN  FM-07 exported_at_commit ($export_sha) behind HEAD ($head_sha) — re-export recommended"
      else
        pass "FM-07 git-sparsity — export SHA matches HEAD"
      fi
    else
      skip "FM-07 git-sparsity — no git or exported_at_commit"
    fi
  fi

  if [[ -f "$GRAPH" ]]; then
    goals=$(count_kind "$GRAPH" goal)
    tasks=$(count_kind "$GRAPH" task)
    decisions=$(count_kind "$GRAPH" decision)
    discoveries=$(count_kind "$GRAPH" discovery)
    if [[ "$goals" -ge 1 && "$tasks" -ge 3 ]]; then
      pass "G3 graph has goal + tasks (goals=$goals tasks=$tasks)"
    else
      fail "G3 graph has goal + tasks — goals=$goals tasks=$tasks"
    fi
    echo "      graph counts: decisions=$decisions discoveries=$discoveries"
  fi

  [[ -f "$WS/.cursor/rules/trace-enforcement.mdc" ]] && pass "E1 cursor rules" || fail "E1 cursor rules"
  [[ -f "$WS/.cursor/hooks/trace-loop-gate.sh" ]] && pass "E2 cursor-hook" || fail "E2 cursor-hook"
  if grep -q '"enforce"[[:space:]]*:[[:space:]]*"strict"' "$WS/.trace/config.json" 2>/dev/null; then
    pass "E3 config enforce strict"
  else
    fail "E3 config enforce strict"
  fi

  if [[ "$CHECK_P25" -eq 1 ]]; then
    RULES="$WS/.cursor/rules/trace-enforcement.mdc"
    if grep -qi "mandatory gap pass" "$RULES" 2>/dev/null; then
      pass "P25-1 GapPassPrompt in installed rules"
    else
      fail "P25-1 GapPassPrompt in installed rules"
    fi
    if grep -qi "Parent orchestrator" "$RULES" 2>/dev/null; then
      pass "P25-2 Parent orchestrator rule in installed rules"
    else
      fail "P25-2 Parent orchestrator rule in installed rules"
    fi
    if [[ -f "$GRAPH" ]]; then
      disc=$(count_kind "$GRAPH" discovery)
      dec=$(count_kind "$GRAPH" decision)
      if [[ "$P25_ARM" == "build" ]]; then
        # P25-3a: build-only baseline — FAIL expected on thin graph
        if [[ "$disc" -ge 1 || "$dec" -ge 1 ]]; then
          pass "P25-3a graph richness (discoveries=$disc decisions=$dec)"
        else
          fail "P25-3a graph richness — need >=1 discovery OR decision (build-only baseline; thin graph expected until directed gap pass)"
        fi
      else
        # P25-3b: directed-gap — PASS required for P25-C validation
        if [[ "$disc" -ge 1 || "$dec" -ge 1 ]]; then
          pass "P25-3b graph richness (discoveries=$disc decisions=$dec)"
        else
          fail "P25-3b graph richness — need >=1 discovery OR decision (directed gap pass required)"
        fi
      fi
    else
      if [[ "$P25_ARM" == "build" ]]; then
        fail "P25-3a graph richness — no export"
      else
        fail "P25-3b graph richness — no export"
      fi
    fi
    # P25-4: arm-matched env attestation (RESULTS.md remains human narrative)
    if [[ "$P25_ARM" == "build" && "${P25_ATTEST_BUILD:-}" == "Y" ]]; then
      pass "P25-4 operator attestation (arm=build)"
    elif [[ "$P25_ARM" == "directed" && "${P25_ATTEST_DIRECTED:-}" == "Y" ]]; then
      pass "P25-4 operator attestation (arm=directed)"
    else
      skip "P25-4 operator attestation — set P25_ATTEST_BUILD=Y / P25_ATTEST_DIRECTED=Y before score (arm=$P25_ARM)"
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
