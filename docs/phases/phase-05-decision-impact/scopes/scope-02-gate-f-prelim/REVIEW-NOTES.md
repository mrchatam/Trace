# P05-S02-02 — Scope review notes (Gate F prelim)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16

## Summary

Independent review of P05-S02-01 against S02-00 / `01-gate-f-prelim.md` locks. Claims match live repo: new package `evals/impact` with named `TestPlantedImpactConflictsGateFPrelim`; committed `schema-gate-f.json` v1 (`schema_version`/`gate`/`suite` consts); temp `metrics-gate-f.json` written + jsonschema-validated; planted tallies TP=3 / FN=0 / FP=0 / TN=1 → precision=1.0 / recall=1.0 / probes=4; Pos-1 scores HasUnknown+Incomplete+UNKNOWN finding (does **not** trust OverallClass alone); Pos-2 SAFE+DESTRUCTIVE → OverallClass=DESTRUCTIVE; Pos-3 link+empty findings → Incomplete; Neg-1 clean SAFE KNOWN → TN; S01 hooks only (`AddImpactFinding` / `LinkDecisionTask` / `ImpactReport` / `decision_affects_task`); mig **009** only (no S02 schema fork); no `cmd/trace` import, no `internal/impact`, no commercial multi-model; Gate G / Gate E / honesty A/B/C / p0x / x0 / `./...` green; Gate C `dry_run:false` N=3 G1 0.800 > B0 0.000 intact. No blocker/high; no spawns.

## Checklist (review focus)

| Focus | Result |
|-------|--------|
| Package `evals/impact` (not folded into honesty/replan/x0/p0x) | **Pass** — `doc.go` + `impact_test.go` + `schema-gate-f.json` |
| Named test `TestPlantedImpactConflictsGateFPrelim` PASS | **Pass** — fresh `CGO_ENABLED=0` below |
| `schema-gate-f.json` v1 + gate F + suite impact consts | **Pass** — `schema_version` const 1; `gate` const `F`; `suite` const `impact` |
| Temp `metrics-gate-f.json` schema-validated | **Pass** — write under `t.TempDir()` + `validateGateFMetricsFile` |
| Tallies TP=3 FN=0 FP=0 TN=1; P/R=1.0; probes=4 | **Pass** — asserted in test + metrics fields |
| Probes Pos-1 UNKNOWN / Pos-2 rollup / Pos-3 empty+link / Neg-1 SAFE | **Pass** — probe_ids match locks |
| Scoring does not trust OverallClass alone when HasUnknown required | **Pass** — Pos-1 asserts flags + UNKNOWN finding; comment forbids OverallClass==UNKNOWN |
| Evidence is harness not Notes-only | **Pass** — automated named test |
| S01 APIs; mig 009 only; no new entity_links rels | **Pass** — domain only; schema list ends at `009_decision_impact.sql` |
| No `cmd/trace` import; no `internal/impact`; no `plan simulate` | **Pass** |
| Carry-forward honesty / Gate G / Gate E / p0x / x0 / `./...` | **Pass** — fresh suites below |
| Gate C untouched; dry-run ≠ Gate C ≠ Gate F | **Pass** — `dry_run:false` N=3; means 0.000 / 0.800 |
| No commercial multi-model Gate F | **Pass** — deterministic domain plants only |

## Claims → evidence

| Claim (P05-S02-01 Notes) | Evidence |
|--------------------------|----------|
| `evals/impact` + `TestPlantedImpactConflictsGateFPrelim` | `evals/impact/impact_test.go` |
| `schema-gate-f.json` v1 | `evals/impact/schema-gate-f.json` (`const: 1` / `F` / `impact`) |
| Temp `metrics-gate-f.json` | Written under `t.TempDir()`; Draft2020 validate |
| TP=3 FN=0 FP=0 TN=1; precision=recall=1.0 | Tallies + metrics fields |
| S01 hooks | `AddImpactFinding` / `LinkDecisionTask` / `ImpactReport`; metrics `s01_hooks` |
| Mig 009 only | `internal/store/schema/009_decision_impact.sql`; no 010+ |
| Gate C untouched | `docs/verification/gate-c-x0/metrics-{b0,g1}.json` still `dry_run:false` N=3 |
| CGO bars | Fresh re-verify below |

## Required tests (fresh this review)

```text
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
  → PASS

CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
  → PASS

CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
  → PASS

CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... -count=1
  → PASS

CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./... -count=1
  → PASS (p0x, x0, honesty, replan, impact, domain, store, planner, mcp, analyzers, …)
```

Gate C spot-check: `dry_run:false`; N=3; mean understanding_accuracy G1 0.800 > B0 0.000 — packs not rewritten.

## Findings

### blocker
_None._

### high
_None._

### medium
_None open (no spawn)._

### low

1. **Schema does not const-lock planted tallies** — `true_positives` etc. are `minimum: 0` only; `s01_hooks` is `minItems: 1` without required string tokens. Live test asserts locked values; a future metrics writer could omit them and still validate. Acceptable for prelim; VERIFY can assert contents if desired (same residual shape as Gate G).

### nit

1. External test package (`package impact_test`) vs `doc.go` (`package impact`) — intentional; matches other eval packages’ doc/test split patterns.
2. Pos-2 does not assert `Incomplete==false` (lock only requires OverallClass=DESTRUCTIVE && !HasUnknown) — correct to lock; Incomplete would also be false for two KNOWN findings.

## Spawns
_None._

## Next board row
**P05-S03-00** (Phase 05 VERIFY planner). Do not invent commercial multi-model Gate F without this harness bar.
