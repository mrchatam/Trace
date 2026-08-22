# FM-09 / FR-P28-06 — notes (P28-S06-11)

**Date:** 2026-08-20  
**Status:** mode-collapse dual-lane proof closed (dogfood; no product code)  
**Harness:** `experiments/ab-p25-gap-pass-validation/score.sh` G1 `--p25`  
**Critical:** `prepare.sh` **NOT RUN** (G1 Session-B rich graph preserved)

## Objective

Prove **build ≠ directed richness** beyond single Session-B: thin Session-A baseline remains documented; directed arm still PASSes P25-3b; live `--arm build` on the same rich graph is labeled **post-directed** and must **not** be read as Session-A thin FAIL.

## Dual-lane results (live re-score)

| Lane | Label | Source | Result | Interpretation |
|------|-------|--------|--------|----------------|
| Thin build baseline | **Session-A thin** | [`../scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json`](../scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json) | discoveries=`[]` decisions=`[]` (disc=0/dec=0) | Historical thin P25-3a FAIL (P27 E02) remains valid; docs-only — no wipe |
| Directed | **Directed P25-3b** | `score-directed.txt` | **PASS** P25-3b (disc=1 dec=1); P25-1/2 PASS; G2 `--strict --enforce` PASS; P25-4 PASS (`P25_ATTEST_DIRECTED=Y`); VERDICT PASS | Primary dogfood regression; matches S02 `SESSION-B-SCORE.txt` |
| Build rich | **post-directed rich build-arm** | `score-build-rich.txt` | **PASS** P25-3a (disc=1 dec=1); P25-4 PASS (`P25_ATTEST_BUILD=Y`); VERDICT PASS | Same on-disk rich G1 scored as build arm — **not** Session-A thin FAIL |

### Anti-conflation (FM-09)

- Rich `--arm build` **P25-3a PASS** does **not** invalidate the Session-A thin baseline.
- Do **not** cite `score-build-rich.txt` as Session-A thin FAIL evidence.
- Thin FAIL evidence stays: snapshot disc=0/dec=0 + historical P27 E02 build score.
- Mode collapse stays closed: directed richness and build-arm richness are separately labeled; build arm on thin graph ≠ build arm on post-directed graph.

## Commands run (no prepare)

```bash
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace

cd experiments/ab-p25-gap-pass-validation
P25_ATTEST_DIRECTED=Y ./score.sh G1 --p25 --arm directed --test
P25_ATTEST_BUILD=Y    ./score.sh G1 --p25 --arm build --test
# NEVER ./prepare.sh
```

## Evidence paths

| Artifact | Path |
|----------|------|
| Run metadata | `experiments/runs/2026-08-20-p28-s06-11-fm09/evidence/run-metadata.txt` |
| Directed score | `experiments/runs/2026-08-20-p28-s06-11-fm09/evidence/score-directed.txt` |
| Rich build score | `experiments/runs/2026-08-20-p28-s06-11-fm09/evidence/score-build-rich.txt` |
| Thin snapshot | `docs/phases/phase-28-residuals-validation/scopes/scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json` |
| Prior S02 directed | `docs/phases/phase-28-residuals-validation/scopes/scope-02-session-b-dogfood/SESSION-B-SCORE.txt` |

## Optional second directed fixture

Not required — live directed re-score on existing Session-B G1 already **P25-3b PASS** (disc=1 dec=1).

## Exit

Acceptance hint met. Next runnable: **P28-S06-12** (FM-09 review).
