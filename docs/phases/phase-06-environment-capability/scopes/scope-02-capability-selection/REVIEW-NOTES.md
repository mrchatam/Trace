# P06-S02-02 — Scope review notes (capability selection ablation)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16

## Summary

Independent review of P06-S02-01 against S02-00 / `01-capability-selection.md` locks. Claims match live repo: new package `evals/capability` with named `TestPlantedCapabilitySelectionAblation`; committed `schema-capability.json` v1 (`schema_version`/`gate`/`suite`/`ablation` consts); temp `metrics-capability.json` written under `t.TempDir()` + jsonschema-validated; planted tallies TP=3 / FN=0 / FP=0 / TN=1 → precision=1.0 / recall=1.0 / probes=4; Pos-1 UNAVAILABLE missing; Pos-2 UNKNOWN missing; Pos-3 selection filter (exactly 2 required; no catalog dump; BuiltinMCP specs seeded AVAILABLE); Neg-1 clean AVAILABLE → TN; S01 hooks only (`UpsertCapability` / `RequireCapability` / `MissingCapabilities` + packet `required_capabilities`/`missing_capabilities`); mig **010** only (no S02 schema fork); G19 library APIs only (no `cmd/trace` import, no `internal/capability`); no commercial multi-model; honesty A/B/C + Gate G / Gate E / Gate F / p0x / x0 / `./...` green; Gate C `dry_run:false` N=3 G1 0.800 > B0 0.000 intact. No blocker/high; no spawns.

## Checklist (review focus)

| Focus | Result |
|-------|--------|
| Package `evals/capability` (not folded into honesty/replan/impact/x0/p0x) | **Pass** — `doc.go` + `capability_test.go` + `schema-capability.json` |
| Named test `TestPlantedCapabilitySelectionAblation` PASS | **Pass** — fresh `CGO_ENABLED=0` below |
| `schema-capability.json` v1 + gate `capability-selection` + suite `capability` + ablation true | **Pass** — consts match locks |
| Temp `metrics-capability.json` schema-validated | **Pass** — write under `t.TempDir()` + `validateCapabilityMetricsFile` |
| Tallies TP=3 FN=0 FP=0 TN=1; P/R=1.0; probes=4 | **Pass** — asserted in test + metrics fields |
| Probes Pos-1 UNAVAILABLE / Pos-2 UNKNOWN / Pos-3 no catalog dump / Neg-1 AVAILABLE | **Pass** — probe_ids match locks |
| Packet lists only required (no catalog dump) | **Pass** — Pos-3 asserts len==2 + non-required slugs absent |
| Evidence is harness not Notes-only | **Pass** — automated named test |
| S01 APIs; mig 010 only; no product Go outside evals | **Pass** — domain+compiler only; schema list ends at `010_capability_surface.sql` |
| G19: no `cmd/trace` scrape; no `internal/capability` | **Pass** |
| Carry-forward honesty / Gate G / Gate E / Gate F / p0x / x0 / `./...` | **Pass** — fresh suites below |
| Gate C untouched; dry-run ≠ Gate C ≠ capability ablation | **Pass** — Gate C `dry_run:false` N=3; means 0.000 / 0.800 |
| No commercial multi-model theater | **Pass** — deterministic domain/compiler plants only |
| Light S03 re-prove command confirmed | **Pass** — already stamped on S03-00/01 |

## Claims → evidence

| Claim (P06-S02-01 Notes) | Evidence |
|--------------------------|----------|
| `evals/capability` + `TestPlantedCapabilitySelectionAblation` | `evals/capability/capability_test.go` |
| `schema-capability.json` v1 | `evals/capability/schema-capability.json` (`const: 1` / `capability-selection` / `capability` / ablation `true`) |
| Temp `metrics-capability.json` | Written under `t.TempDir()`; Draft2020 validate; `mig`=`010_capability_surface`; `dry_run:false` |
| TP=3 FN=0 FP=0 TN=1; precision=recall=1.0 | Tallies + metrics fields |
| S01 hooks | `UpsertCapability` / `RequireCapability` / `MissingCapabilities` + packet fields; metrics `s01_hooks` |
| Mig 010 only | `internal/store/schema/010_capability_surface.sql`; no 011+ |
| Gate C untouched | `docs/verification/gate-c-x0/metrics-{b0,g1}.json` still `dry_run:false` N=3 |
| CGO bars | Fresh re-verify below |

## Required tests (fresh this review)

```text
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
  → PASS

CGO_ENABLED=0 go test ./evals/honesty/... ./evals/replan/... ./evals/impact/... -count=1
  → PASS (honesty A/B/C + Gate G; Gate E; Gate F)

CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./... -count=1
  → PASS (p0x, x0, honesty, replan, impact, capability, domain, store, planner, mcp, analyzers, …)
```

Gate C spot-check: `dry_run:false`; N=3; mean understanding_accuracy G1 0.800 > B0 0.000 — packs not rewritten.

## Light S03 note

Upcoming VERIFY re-prove command already locked (P06-S02-00 / S03-00 Depends):

```bash
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
```

No spawn required to thicken S03.

## Findings

### blocker
_None._

### high
_None._

### medium
_None open (no spawn)._

### low

1. **Schema does not const-lock planted tallies** — `true_positives` etc. are `minimum: 0` only; `s01_hooks` is `minItems: 1` without required string tokens. Live test asserts locked values; a future metrics writer could omit them and still validate. Acceptable for ablation; VERIFY can assert contents if desired (same residual shape as Gate F/G).

### nit

1. External test package (`package capability_test`) vs `doc.go` (`package capability`) — intentional; matches other eval packages’ doc/test split patterns.
2. Pos-3 seeds BuiltinMCP specs plus extras (N≫3 AVAILABLE) then requires exactly `tool:alpha`/`tool:beta` — stronger than the lock’s N≥3 floor; correct and within optional BuiltinMCP plant.

## Spawns
_None._

## Next board row
**P06-S03-00** (Phase 06 VERIFY planner). Do not invent commercial multi-model capability theater without this harness bar.
