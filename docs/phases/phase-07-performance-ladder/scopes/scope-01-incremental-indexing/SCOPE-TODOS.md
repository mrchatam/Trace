# Scope S01 — Incremental indexing / ignore tiers

**Depends-on:** `P07-00` done (phase locks).

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P07-S01-00 | planner | done | 2026-08-16: T0 dirs/suffixes/path-segment locked; no `011_*`; measurement = CLI tests + optional evals/perf seed; `01` runnable |
| P07-S01-01 | implement | done | 2026-08-16: T0 helpers + walk/explicit argv; isolation+T0 tests |
| P07-S01-02 | review | done | 2026-08-16: APPROVE high — [REVIEW-NOTES.md](REVIEW-NOTES.md) |

## Checklist

- [x] P07-S01-00 planner
- [x] P07-S01-01 implement
- [x] P07-S01-02 review

## Locked (S01-00)

- Package: `cmd/trace` walk/index; analyzers binary NUL unchanged; no `internal/indexer`
- Migration: **none** (path filter only)
- T0 dirs: `.git`/`.trace`/`node_modules`/`vendor`/`__pycache__`/`.venv`/`venv`/`dist`/`.next`/`target`/`coverage`
- T0 files: `.min.js`/`.min.mjs`/`.min.cjs` + path-segment rule
- Measurement: `TestIndexIncrementalIsolation` + `TestWalkIndexableT0AlwaysSkip` + explicit T0 argv; optional `evals/perf` seed — **no** Gate H pass
- T1–T3: notes only

## Phase locks (from P07-00)

- File-local incremental only; no full-rebuild-on-any-change
- Prefer T0 always-skip + measurable quality; Gate H thresholds later (`evals/perf`)
- Carry-forward honesty / Gates E/F/G / ablation / p0x / x0 / Gate C intact
