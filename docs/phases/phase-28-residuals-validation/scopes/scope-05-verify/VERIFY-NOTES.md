# Phase 28 VERIFY notes

**Run:** 2026-08-20  
**Row:** P28-S05-01  
**Evidence dir:** experiments/runs/2026-08-20-p28-s05-01-verify/evidence/  
**Git SHA:** `unknown`  
**Harness:** score.sh G1 --p25 (dual-lane; no prepare.sh)

## Verdict

**PASS** — confidence **high**

## Per-block

| Block | Result | Notes |
|-------|--------|-------|
| 1 Build | PASS | `CGO_ENABLED=1 go build -o bin/trace ./cmd/trace` exit 0 |
| 2 Unit + cmd | PASS | `GOPROXY=direct go test ./internal/...` + `./cmd/trace/...` both exit 0 (`unit.txt`, `cmd.txt`) |
| 3 Install + hook smoke | PASS | Option A: `CursorLoopGateFailClosed` / `CursorLoopGateAllowNonStrict` / `HookDrift*` exit 0 (`install.txt`, `hook-smoke.txt`) |
| 4 Matrix M-16 | PASS | `evals/p28-regression/score_arm_labels_test.sh` PASS; TEST-MATRIX.md M-01..M-16 cited (M-16 = label smoke, not P25-3b dogfood) |
| 5 Score directed | PASS | `p25 arm: directed`; **P25-3b PASS** (disc=1 dec=1); P25-1/2 PASS; G2 `--strict --enforce` PASS; P25-4 PASS (`P25_ATTEST_DIRECTED=Y`); VERDICT PASS. Matches S02 `SESSION-B-SCORE.txt` |
| 6 Score build rich | PASS | `p25 arm: build`; **P25-3a PASS** (disc=1 dec=1) on rich graph; P25-4 PASS (`P25_ATTEST_BUILD=Y`); labeled **post-Session-B rich build-arm** |
| 6 Thin baseline docs | PASS | `SESSION-A-GRAPH-SNAPSHOT.json` discoveries=[] decisions=[]; historical thin P25-3a FAIL remains valid (P27 E02); no wipe |

## Dual-lane G1

| Lane | Result | Interpretation |
|------|--------|----------------|
| Directed live | P25-3b PASS (disc=1 dec=1) | Primary dogfood regression — no prepare.sh |
| Build live rich | P25-3a PASS (disc=1 dec=1) | post-Session-B — not thin Session-A |
| Build thin (snapshot) | disc=0/dec=0 | Historical thin baseline; do not wipe |

## Residual register R1–R8 (final)

| ID | Final status | Evidence |
|----|--------------|----------|
| R1 | closed | Session-B + directed re-score P25-3b PASS (`score-directed.txt`) |
| R2 | closed | Option A hook: strict+empty task deny (`hook-smoke.txt`) |
| R3 | closed | FM-05 script deny under Option A (install/hook smoke) |
| R4 | closed | honesty single source (S04); unit/cmd green |
| R5 | closed | `P25_ATTEST_DIRECTED=Y` / `P25_ATTEST_BUILD=Y` → P25-4 PASS both arms |
| R6 | partial/deferred | FM-01/02/04/07/08/09/10 remain behavioral/measurement gaps; FM-03/05/06 largely closed |
| R7 | closed | TEST-MATRIX M-01..M-16 + this VERIFY (unit/cmd/install/matrix/score) |
| R8 | closed | `hook_drift_test` via HookDrift filter smoke |

## vs Phase 27 VERIFY baseline

| Topic | P27 | P28 |
|-------|-----|-----|
| P25-3b | FAIL OK (no Session-B) | PASS (Session-B + re-score) |
| Hook failClosed | residual | Option A closed |
| Honesty dup | residual | R4 closed |
| P25-4 | skip/manual | env attestation PASS both arms |
| Build thin P25-3a | FAIL expected | snapshot docs + rich live PASS labeled |

## Gaps / spawn

(none) for VERIFY repair. Remaining FM measurement gaps (R6) are deferred — S05-02 decides successor / `no successor`.

## DR-HANDOFF status

**OPEN** — S05-02 closes with successor decision (never TBD).

## Residual-wave pointer (additive — P28-S06-07)

**FM-07 / FR-P28-04:** Remain **warn-only** (git-sparsity / post-hoc SPEC SHA drift). Decision: [`../scope-06-r6-fm-residuals/FM07-DECISION.md`](../scope-06-r6-fm-residuals/FM07-DECISION.md). Harness: PROTOCOL FM-07 + `score.sh` T03 — warn ≠ fail; do not treat FM-07 WARN as G2 hard-fail.
