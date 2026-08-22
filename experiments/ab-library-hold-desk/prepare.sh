#!/usr/bin/env bash
# Prepare E03 workspaces. Never wipe an arm unless that arm is requested.
#
#   ./prepare.sh          # both, but refuse if B0 already has product work
#   ./prepare.sh G1       # G1 only (use this after B0 is done)
#   ./prepare.sh B0       # B0 only
#   PREPARE_FORCE=1 ./prepare.sh [both|B0|G1]
set -euo pipefail

EXP="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO="$(cd "$EXP/../.." && pwd)"
TRACE_BIN="${TRACE_BIN:-$REPO/bin/trace}"
FORCE="${PREPARE_FORCE:-0}"
TARGET="${1:-both}"

usage() {
  echo "usage: $0 [both|B0|G1]" >&2
  echo "  default both — will not destroy a populated B0 unless PREPARE_FORCE=1" >&2
  echo "  after B0 completes, run: $0 G1" >&2
  exit 2
}

b0_has_product_work() {
  local ws="$EXP/runs/B0"
  [[ -d "$ws" ]] || return 1
  local extra
  extra=$(find "$ws" \( -name '*.go' -o -name '*.html' -o -name 'SPEC.md' -o -name 'VERIFY.md' \) \
    ! -path "$ws/.git/*" 2>/dev/null | wc -l)
  extra=${extra// /}
  [[ "$extra" -gt 1 ]]
}

trace_supports_cursor_hook() {
  [[ -x "$TRACE_BIN" ]] && "$TRACE_BIN" help 2>&1 | grep -q 'cursor-hook'
}

if [[ ! -x "$TRACE_BIN" ]]; then
  echo "error: trace binary not found at $TRACE_BIN" >&2
  echo "  rebuild: cd $REPO && CGO_ENABLED=1 go build -o bin/trace ./cmd/trace" >&2
  exit 2
fi
if ! trace_supports_cursor_hook; then
  echo "error: $TRACE_BIN is too old (missing cursor-hook)" >&2
  exit 2
fi

case "$TARGET" in
  both|B0|G1) ;;
  -h|--help) usage ;;
  *) usage ;;
esac

if [[ "$TARGET" == "both" || "$TARGET" == "B0" ]]; then
  if b0_has_product_work && [[ "$FORCE" != "1" ]]; then
    echo "error: runs/B0 already has product work — refusing to wipe B0" >&2
    echo "  to prepare G1 only (keeps B0): $0 G1" >&2
    echo "  to wipe B0 anyway: PREPARE_FORCE=1 $0 $TARGET" >&2
    exit 3
  fi
fi

mkdir -p "$EXP/runs"

prep_ws() {
  local arm="$1"
  local ws="$EXP/runs/$arm"
  rm -rf "$ws"
  cp -a "$EXP/project" "$ws"
  (
    cd "$ws"
    git init -q
    git config user.email "relay@trace.local"
    git config user.name "relay"
    git add -A
    git commit -q -m "init $arm starter (E03 library hold desk)"
  )
}

do_b0=0
do_g1=0
[[ "$TARGET" == "both" || "$TARGET" == "B0" ]] && do_b0=1
[[ "$TARGET" == "both" || "$TARGET" == "G1" ]] && do_g1=1

[[ "$do_b0" -eq 1 ]] && prep_ws B0
[[ "$do_g1" -eq 1 ]] && prep_ws G1

install_workspace_guard() {
  local ws="$1"
  local rule_src="$EXP/.cursor/rules/e03-workspace-only.mdc"
  if [[ -f "$rule_src" ]]; then
    mkdir -p "$ws/.cursor/rules"
    cp -a "$rule_src" "$ws/.cursor/rules/e03-workspace-only.mdc"
  fi
}

if [[ "$do_g1" -eq 1 ]]; then
  cp -a "$EXP/seed" "$EXP/runs/G1/seed"
  (
    cd "$EXP/runs/G1"
    "$TRACE_BIN" init
    CGO_ENABLED=1 "$TRACE_BIN" index
    "$TRACE_BIN" seed import seed/gt.json
    echo '{"enforce":"strict"}' > .trace/config.json
    "$TRACE_BIN" install cursor --write
    "$TRACE_BIN" install cursor-hook --write
  )
  install_workspace_guard "$EXP/runs/G1"
  echo "note: score.sh exports trace/graph.json from DB if missing before scoring"
fi

if [[ "$do_b0" -eq 1 ]]; then
  install_workspace_guard "$EXP/runs/B0"
fi

echo "prepared E03: target=$TARGET TRACE_BIN=$TRACE_BIN"

if [[ "$do_g1" -eq 1 ]]; then
  WS="$EXP/runs/G1"
  RULES="$WS/.cursor/rules/trace-enforcement.mdc"
  if grep -q "mandatory gap pass" "$RULES" 2>/dev/null; then
    echo "verify: GapPassPrompt present in installed rules — OK"
  else
    echo "warning: GapPassPrompt missing — rebuild bin/trace from repo HEAD" >&2
  fi
  if grep -q "Parent orchestrator" "$RULES" 2>/dev/null; then
    echo "verify: ParentOrchestratorRule present — OK"
  else
    echo "warning: Parent orchestrator missing from rules" >&2
  fi
  TASK="e0300000-0000-4000-8000-000000000010"
  if ! "$TRACE_BIN" -C "$WS" loop gate --task "$TASK" --for edit >/tmp/e03-gate.json 2>/dev/null; then
    if grep -q plan_missing /tmp/e03-gate.json 2>/dev/null; then
      echo "verify: gate blocks edit on fresh seed (plan_missing) — OK"
    else
      echo "warning: gate blocked but reason unexpected — see /tmp/e03-gate.json" >&2
    fi
  else
    echo "warning: gate should block fresh seed" >&2
  fi
fi
