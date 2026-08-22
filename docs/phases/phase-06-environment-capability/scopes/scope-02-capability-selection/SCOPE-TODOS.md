# Scope S02 — Capability selection

**Depends-on:** S01 APPROVE (`P06-S01-02`).

**Out:** Phase VERIFY (S03); commercial multi-model theater; ontology expansion beyond S01.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P06-S02-00 | planner | done | 2026-08-16: FINAL locks — `evals/capability` / `TestPlantedCapabilitySelectionAblation` / `schema-capability.json` v1 + temp `metrics-capability.json`; P/R TP=3/FN=0/FP=0/TN=1; no product Go |
| P06-S02-01 | implement | done | 2026-08-16: ablation harness shipped; TP=3/FN=0/FP=0/TN=1 P/R=1.0 |
| P06-S02-02 | review | done | 2026-08-16: APPROVE high; no spawns — [REVIEW-NOTES.md](REVIEW-NOTES.md); next P06-S03-00 |

## Checklist

- [x] P06-S02-00 planner
- [x] P06-S02-01 implement
- [x] P06-S02-02 review

## Ablation path (P06-S02-00 FINAL)

| Item | Value |
|------|-------|
| Package | `evals/capability` |
| Named test | `TestPlantedCapabilitySelectionAblation` |
| Schema | `evals/capability/schema-capability.json` v1 |
| Metrics | Temp `metrics-capability.json` |
| Evidence | TP=3 / FN=0 / FP=0 / TN=1 → P/R=1.0 |
| Probes | Pos-1 UNAVAILABLE; Pos-2 UNKNOWN; Pos-3 no catalog dump; Neg-1 clean AVAILABLE |
| Re-prove | `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` |

## S01 hooks to consume (live-confirmed)

- Domain: `UpsertCapability`, `RequireCapability`, `MissingCapabilities`, `ListRequiredCapabilities`
- Compiler packet fields: `required_capabilities` / `missing_capabilities`
- Mig `010_capability_surface` only (no S02 schema fork)
- Optional: `BuiltinMCPCapabilitySpecs()` for six MCP slugs
- Harness calls library APIs (G19) — CLI optional
- Product Go: **none outside `evals/capability`** (+ schema)
