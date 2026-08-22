#!/usr/bin/env bash
# P28-S01 optional P25-D smoke: score.sh arm labels (no harness execution).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SCORE="$ROOT/experiments/ab-p25-gap-pass-validation/score.sh"

if [[ ! -f "$SCORE" ]]; then
  echo "FAIL: missing $SCORE" >&2
  exit 1
fi

fail=0
for needle in 'P25-3a' 'P25-3b' '--arm'; do
  if grep -q -- "$needle" "$SCORE"; then
    echo "PASS: score.sh contains ${needle}"
  else
    echo "FAIL: score.sh missing ${needle}" >&2
    fail=1
  fi
done

exit "$fail"
