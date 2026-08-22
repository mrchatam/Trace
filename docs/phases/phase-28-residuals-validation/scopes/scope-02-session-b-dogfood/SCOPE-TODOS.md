# S02 scope todos — Session-B dogfood

| ID | Status | Prompt | Lock |
|----|--------|--------|------|
| P28-S02-00 | done (after planner) | [00-PLANNER.md](00-PLANNER.md) | Protocol + live preflight |
| P28-S02-01 | pending | [01-run-and-score.md](01-run-and-score.md) | **Agent-executable** gap in `runs/G1` via `TRACE_BIN -C`; evidence in this folder |
| P28-S02-02 | pending | [02-review.md](02-review.md) | Independent; no prepare; no `--arm build` |

## Evidence artifacts (S02-01 writes)

| File | Role |
|------|------|
| `SESSION-A-GRAPH-SNAPSHOT.json` | Thin Session-A graph before mutation |
| `SESSION-B-SCORE.txt` | `./score.sh G1 --p25 --arm directed --test` |
| `SESSION-B-NOTES.md` | Operator notes + P25-4 |
| `SESSION-B-BLOCKED.md` | Only if G1/preflight missing |
| `experiments/RESULTS.md` | New row `E02-SB` |

## Explicit non-goals

- `./prepare.sh` wipe
- Re-score `--arm build` after directed writes
- S01 re-inventory / `apply_promotion_test.go`
- S03 hook deny / S04 honesty-dup + attestation automation
- `trace gap` CLI (does not exist — use `loop status`)
