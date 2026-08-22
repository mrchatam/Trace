# Phase 28 — Residuals closure + implementation validation

**Phase planner.** Runs as row `P28-00` on the board.

## Metadata
- id: P28-00
- todo_ids: [P28-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified]
- verification: automated

## Mission

Close all open residuals from Phases 24–27 and **validate shipped implementations** with automated tests + live dogfood.

| Residual source | Theme | Scope |
|-----------------|-------|-------|
| P27 VERIFY | Session-B unvalidated (P25-3b) | S02 |
| P26/P27 DR-HANDOFF | Hook failClosed beyond install text (FM-04/05) | S03 |
| P27 S02-02 | BLOCKING duplicate honesty msg | S04 |
| P27 VERIFY | P25-4 attestation manual | S04 |
| FM matrix §3 | Residual gaps after INT-01..11 | S00 → S01 |
| All P25 themes | No consolidated regression matrix | S01 + S05 |

## Scope sequence

```
S00 Residual audit → S01 Integration tests → S02 Session-B dogfood → S03 Hook failClosed → S04 Polish → S05 VERIFY
```

S02 (live agent dogfood) may run in parallel with S01 after S00 audit locks protocol — default serial.

## Hard constraints

- No daemon, HTTP server, or hosted service
- No rewriting done board history (Phases 24–27 rows stay closed)
- Session-B dogfood must not invalidate build-only Session-A evidence (separate score row / `--arm directed`)
- Auto-spawn from discoveries remains human-gated (Phase 24 §4)
- `go test ./internal/...` green after product-touching scopes

## Deliverables expected by VERIFY (S05)

| ID | Deliverable | Scope |
|----|-------------|-------|
| D1 | `RESIDUAL-AUDIT.md` — every open FM/residual mapped to closure task or explicit defer | S00 |
| D2 | Integration test matrix covering promotion, reset, saturation, honesty, install | S01 |
| D3 | E02 Session-B scored: `./score.sh G1 --p25 --arm directed` → P25-3b PASS | S02 |
| D4 | Hook deny when TRACE_TASK_ID absent under strict config (or documented failClosed path) | S03 |
| D5 | BLOCKING dup msg fixed; P25-4 attestation in score.sh or PROTOCOL | S04 |
| D6 | Full regression: unit + cmd + score build + score directed PASS/FAIL documented | S05 |

## Planner gate (P28-00)

Verify all scope scaffold paths exist before closing this row:

- `docs/phases/phase-28-residuals-validation/` — this file, README.md, DR-HANDOFF.md OPEN
- `scopes/scope-00-residual-audit/` — 00-PLANNER, 01-residual-audit, 02-review
- `scopes/scope-01-integration-tests/` — 00-PLANNER, 01-implement, 02-review
- `scopes/scope-02-session-b-dogfood/` — 00-PLANNER, 01-run-and-score, 02-review
- `scopes/scope-03-hook-failclosed/` — 00-PLANNER, 01-implement, 02-review
- `scopes/scope-04-product-polish/` — 00-PLANNER, 01-implement, 02-review
- `scopes/scope-05-verify/` — 00-PLANNER, 01-verify, 02-dr-handoff
- Board row in `docs/TODO/phase-28.md` ✓

## Next

`P28-S00-00`
