# Scope S02 — Gate F prelim

**Depends-on:** S01 APPROVE (`P05-S01-02`).

**Out:** Phase VERIFY (S03); commercial multi-model Gate F; inventing Gate F from vibes; inventing impact APIs (consume S01).

**Gate F lock (P05-S02-00 FINAL — 2026-08-16):**
| Item | Value |
|------|-------|
| Package | `evals/impact` |
| Named test | `TestPlantedImpactConflictsGateFPrelim` |
| Schema | `evals/impact/schema-gate-f.json` v1 |
| Metrics | temp `metrics-gate-f.json` |
| Tallies | TP=3 FN=0 FP=0 TN=1; precision=1.0; recall=1.0; probes=4 |
| Re-prove | `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` |

**S01 hooks (confirmed P05-S01-02 APPROVE):** plant via `AddImpactFinding` + `LinkDecisionTask` + `ImpactReport` (`HasUnknown`/`Incomplete`/`OverallClass`; do not trust OverallClass alone when HasUnknown); mig `009_decision_impact`; harness calls domain APIs (not CLI-only).

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P05-S02-00 | planner | done | 2026-08-16: locked Gate F names + 4-probe P/R protocol; thickened 01+02; light S03 Depends; no product Go |
| P05-S02-01 | implement | done | 2026-08-16: `evals/impact` harness + schema v1 + temp metrics; TP=3/FN=0/FP=0/TN=1 |
| P05-S02-02 | review | done | 2026-08-16: APPROVE **high**; no spawns — [REVIEW-NOTES.md](REVIEW-NOTES.md); next **P05-S03-00** |

## Checklist

- [x] P05-S02-00 planner
- [x] P05-S02-01 implement
- [x] P05-S02-02 review
