# P28-S02-00 — Scope planner (Session-B dogfood)

## Metadata
- id: P28-S02-00
- todo_ids: [P28-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: mixed
- hooks: []

## Objective

Lock Session-B protocol to validate **P25-3b** (directed gap richness) on E02 G1 using Phase 27 harness. Finalize `01-run-and-score.md` + `02-review.md` without invalidating Session-A build evidence.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R1 status
- [Phase 28 README](../../README.md)
- `experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-DIRECTED-GAP.md`
- [Phase 27 PROTOCOL](../../../phase-27-protocol-measurement-graph-honesty/scopes/scope-01-protocol-v2/) (via harness PROTOCOL.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked protocol

| Step | Action |
|------|--------|
| 1 | Operator opens `runs/G1` — **no** `./prepare.sh` wipe |
| 2 | Paste `PROMPT-G1-DIRECTED-GAP.md` |
| 3 | Agent runs mandatory gap pass per installed rules |
| 4 | Export: `trace seed export -o trace/graph.json --strict --enforce` |
| 5 | Score: `./score.sh G1 --p25 --arm directed --test` |
| 6 | Append separate row in `experiments/RESULTS.md` (Session-B) |

## Preconditions

- E02 Session A (build-only) already scored — do not wipe G1 product work
- `score.sh` supports `--arm directed` (Phase 27 S01)
- `TRACE_TASK_ID=e0200000-0000-4000-8000-000000000010` for G1

## S01 automated baseline (do not re-inventory)

Automated P25-A/B/C/D/E coverage is locked in [`../scope-01-integration-tests/TEST-MATRIX.md`](../scope-01-integration-tests/TEST-MATRIX.md) (M-01..M-16). S02 validates **live Session-B P25-3b richness** only.

- M-16 (`evals/p28-regression/score_arm_labels_test.sh`) is grep-only for `P25-3a`/`P25-3b`/`--arm` — **not** a P25-3b PASS.
- Do not add `apply_promotion_test.go` (M-01 already covers apply E2E).
- Do not treat S01 as dogfood evidence. Do not implement hook deny (S03) or honesty-dup/attestation (S04) here.

## Success criteria (P25-3b)

- discoveries ≥ 1 OR decisions ≥ 1
- P25-1/2 still PASS
- Optional: promotion path if agent uses `loop apply` / `--from-discovery`
- P25-4 operator attestation recorded (manual until S04)

## Planner gate

- [x] `RESIDUAL-AUDIT.md` confirms R1 open
- [x] `01-run-and-score.md` + `02-review.md` runnable
- [x] Arm isolation documented (build vs directed)

## Exit criteria

- [x] Dogfood operator prompt locked for fresh subagent
- [x] Board row P28-S02-00 Notes cite locks
- [x] Next runnable **P28-S02-01**

## Todo updates

Status + notes on **P28-S02-00** only.

## Live locks (planner 2026-08-20)

Preflight verified in-repo (no product code):

| Check | Result |
|-------|--------|
| R1 open | RESIDUAL-AUDIT.md: Session-B not run; P25-3b unvalidated |
| G1 present | `runs/G1/.trace/` + `trace/graph.json` — goals=1 tasks=5 disc=0 dec=0 |
| Session-A scored | `experiments/RESULTS.md` E02 P27 row; VERIFY-NOTES P25-3a FAIL expected |
| `score.sh --arm directed` | `P25_ARM` parse L22–28; P25-3b labels L204–208 |
| `TRACE_TASK_ID` | `e0200000-0000-4000-8000-000000000010` (PROTOCOL + PROMPT) |
| `PROMPT-G1-DIRECTED-GAP.md` | exists |
| `trace gap` | **not a command** — S02-01 uses `loop status` |
| Live gate | `loop status` → edit blocked `plan_missing` (record as gap, do not wipe) |

**S02-01 mode:** agent-executable with `TRACE_BIN -C runs/G1`. Human paste optional. Blocked only if G1 missing (then `SESSION-B-BLOCKED.md`, no `prepare.sh`).

**Arm isolation:** snapshot graph before writes; never `./score.sh G1 --p25` after mutation (defaults to build).

## Next

`P28-S02-01`
