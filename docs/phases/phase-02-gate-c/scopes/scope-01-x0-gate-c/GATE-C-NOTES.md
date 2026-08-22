# GATE-C-NOTES — Experiment X0 (P02-S01-01)

## Verdict

**Go**

Mean G1 `understanding_accuracy` (0.800) **>** mean B0 (0.000) with complete evidence table, schema-valid `dry_run:false` metrics (N=3/condition), and non-trivial human GT seed. Kill criteria did **not** fire.

## Evidence table

| Condition | N | Mean understanding_accuracy | Mean critical_misses | Mean latency_ms (pack) | Mean tokens (pack) | Model pin |
|-----------|---|----------------------------:|---------------------:|-----------------------:|-------------------:|-----------|
| B0 | 3 | 0.000 | 3.000 | 1217 | 907 | `recorded-operator-sim/v1` |
| G1 | 3 | 0.800 | 0.000 | 450 | 1160 | `recorded-operator-sim/v1` |

Per-run accuracies:

| Run | Accuracy | Critical misses | Notes |
|-----|---------:|----------------:|-------|
| b0-gatec-1 | 0.0 | 4 | src-only; invented wrong causal roles |
| b0-gatec-2 | 0.0 | 3 | src-only |
| b0-gatec-3 | 0.0 | 2 | src-only |
| g1-gatec-1 | 0.8 | 0 | why+context; miss q3 only |
| g1-gatec-2 | 0.8 | 0 | why+context; miss q3 only |
| g1-gatec-3 | 0.8 | 0 | why+context; miss q3 only |

Artifacts:

- Metrics: [`docs/verification/gate-c-x0/metrics-b0.json`](../../../../verification/gate-c-x0/metrics-b0.json), [`metrics-g1.json`](../../../../verification/gate-c-x0/metrics-g1.json)
- Pins: [`docs/verification/gate-c-x0/pins.md`](../../../../verification/gate-c-x0/pins.md)
- Queries: [`evals/x0/queries.json`](../../../../../evals/x0/queries.json)
- Packs: [`evals/x0/testdata/gate-c/`](../../../../../evals/x0/testdata/gate-c/)

Fixture content hash: `bcc50f8e3b027c111a0fe1db251a440d99af1f7f29abbf16e93c267b2cf2074c`  
Trace version: `0.0.0-dev`  
Mode: **B — recorded** operator-sim packs (CLI path; MCP not used).

## Kill criteria check

| Conjunct | Result |
|----------|--------|
| mean G1 understanding_accuracy ≤ mean B0 | **No** (0.800 > 0.000) |
| Seeding non-trivial (multi-entity causal graph + links in human GT) | **Yes** (satisfied for `fixtures/x0`) |
| Thesis endangered? | **No** |

## Issue list (S02 / GC-NN)

```text
- id: GC-01
  severity: medium
  metric: G1 understanding_accuracy capped at 0.8; q3 incorrect on all 3 runs (0 critical_miss)
  evidence: docs/verification/gate-c-x0/metrics-g1.json per_query q3; GATE-C-NOTES per-run table; live why/context omit discovery→plan_change
  proposed_fix_surface: internal/retrieval Why/Expand neighbors + compiler TaskContext linkage for discovery↔plan_change in task neighborhood (or document secondary query path)
  defer: false

- id: GC-02
  severity: medium
  metric: fairness residual — B0 could ace GT if seed/README oracle readable (packs used oracle-exclusion brief only)
  evidence: fixtures/x0/seed/gt.json + fixtures/x0/README.md UUID table; pins.md Agent brief; evals/x0/testdata/gate-c/README.md
  proposed_fix_surface: fixtures/x0 layout / eval harness agent-hide of GT for live Gate C (or documented oracle policy)
  defer: false

- id: GC-03
  severity: low
  metric: model pin is recorded-operator-sim/v1 (not production coding model)
  evidence: docs/verification/gate-c-x0/pins.md; metrics-*.json model field
  proposed_fix_surface: evals/x0/testdata/gate-c pack refresh under pinned live model
  defer: true

- id: GC-04
  severity: low
  metric: within-condition variance = 0 (identical grade patterns across N=3)
  evidence: GATE-C-NOTES per-run table; metrics-b0.json / metrics-g1.json quality objects
  proposed_fix_surface: defer — increase N or live-model stochasticity only if significance claims needed
  defer: true
```

## Honesty

Phase 01 `TestX0DryRunMetricsB0AndG1` (`dry_run:true`) was **not** treated as a Gate C pass. Gate C evidence is the graded Mode-B packs + `dry_run:false` metrics above.

## Residuals

- Sample size N=3; no statistical significance test (per locked protocol).
- Unfairness risk if future operators allow B0 to read `seed/gt.json` / README GT tables without documenting it (see GC-02).
- G1 may still use repo tools; packs did not give G1 the seed oracle either.
- Secondary task families (`implementation`, `honesty`) not used for Gate C pass/fail; H5 remains in `evals/honesty`.

## How to re-verify

```bash
CGO_ENABLED=1 go test ./evals/x0/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/honesty/... ./evals/x0/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```
