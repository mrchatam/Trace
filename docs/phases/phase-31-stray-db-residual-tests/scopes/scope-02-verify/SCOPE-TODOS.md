# Scope 02 — board map

VERIFY + DR-HANDOFF. Serial: **S02-00 → S02-01 → S02-02**. Do not close phase until S02-02. Successor on green: **Phase 32** (`P32-00`) — never TBD.

| Board ID | Row | Prompt | Role | Status (planner) |
|----------|-----|--------|------|------------------|
| 545 | P31-S02-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner (lock VERIFY floor; no product code) | **done** — floor locked 2026-08-21 |
| 546 | P31-S02-01 | [01-verify.md](01-verify.md) | Verifier → VERIFY-NOTES + evidence | pending — run next |
| 547 | P31-S02-02 | [02-dr-handoff.md](02-dr-handoff.md) | Reviewer → close DR-HANDOFF; successor Phase 32 | pending — after S02-01 |

## Locked VERIFY floor (S02-00)

| Block | Owner | Action |
|------:|-------|--------|
| 0 | S02-01 | `experiments/runs/YYYY-MM-DD-p31-s02-01-verify/evidence/` + metadata (git SHA; cite S01-02 PASS) |
| 1 | S02-01 | Five store tests in `stray_trace_db_test.go` (incl. G1 dir-stub quiet) |
| 2 | S02-01 | `go test ./internal/...` |
| 3 | S02-01 | `bash scripts/repro-stray-trace-db.sh` (G5 shipped; missing → FAIL) |
| 4 | S02-01 | gitignore + open.go join/Stat-only + G6 CONTRIBUTING/AGENTS |
| 5 | S02-01 | Residuals listed; overall PASS/FAIL; DR-HANDOFF stays OPEN |

## Fail vs residual (summary)

- **FAIL:** stray tests / internal bar / repro / missing G1|G5|G6 / path redesign / silent delete / init creates root `trace.db`
- **Residual OK:** G2 nice skip; G3/G4 deferred; multi-open once-per-`openStore`; agent stubs; optional delete future-only

## DR-HANDOFF close (S02-02 only)

- Evidence gatherer: **S02-01**
- Close owner: **S02-02**
- On APPROVE: CLOSED + successor **Phase 32** / **P32-00**
- On regression: keep OPEN; spawn `P31-S02-02a`/`02b`; successor = `pending repair spawn` (not TBD)
- Confirm Phase 32 scaffold exists; do not implement GUI here

DR-HANDOFF stays **OPEN** until S02-02.
