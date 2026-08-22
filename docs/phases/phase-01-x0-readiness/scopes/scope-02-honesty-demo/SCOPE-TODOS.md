# Scope S02 — Honesty demo

**Depends-on:** S01 Claim/Evidence/Review promotion `done` (APPROVE high, 2026-08-16).

**Goal:** Deterministic H5-partial fail-closed proof that a planted false/incomplete claim cannot reach DONE.

## Locks (from P01-S02-00)

| Item | Value |
|------|-------|
| Package | `evals/honesty` (named artifact; keep separate from `evals/p0x` / future `evals/x0`) |
| Verification | `automated` |
| S01 surface | `CreateReview` / `SetReviewResult` / `LinkReviewTask` → DONE only on linked PASS |
| Scenario | Path A EvidenceIDs-alone reject; Path B FAIL reject; Path C second-review PASS → DONE |
| Escape hatch | Do **not** prove honesty via `AllowDoneWithoutReview` / `--allow-done` |
| CGO | `CGO_ENABLED=0` for `./evals/honesty/...`; regression `CGO_ENABLED=1 ./...` |

## Board

- [x] P01-S02-00 planner
- [x] P01-S02-01 implement
- [x] P01-S02-02 review (+ spawns as needed) — APPROVE high; no spawns; [REVIEW-NOTES.md](REVIEW-NOTES.md)
