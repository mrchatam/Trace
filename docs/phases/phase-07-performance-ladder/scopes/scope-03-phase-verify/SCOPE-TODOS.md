# Scope S03 — Phase 07 VERIFY

**Depends-on:** S01 + S02 APPROVE (**S02-02 done 2026-08-16**). **S01-02:** T0 + isolation proofs required at VERIFY; **no** `evals/perf` fixtures from S01 (only T0 `t.Logf`). **S02-02:** Go adapter live (`tree-sitter-go` v0.25.0); optional `.go` ladder fixtures OK; **no** thresholds from S02.

**Gate H ownership (FINAL P07-S03-00):** **S03-01 creates** planted `evals/perf` harness as VERIFY work (measure-then-threshold). Do **not** block Gate H awaiting a prior seed.

**DR-HANDOFF:** On pass — VERIFY starts next-phase scaffold **`phase-08-ecosystem-hardening`** + board `P08-00`; S03-02 owns completion (or records `no successor`).

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P07-S03-00 | planner | done | 2026-08-16: locked Gate H names + measure-then-threshold + S03-01 creates harness; thickened 01+02; Phase 08 checklist |
| P07-S03-01 | verify | done | 2026-08-16: Gate H green; VERIFY-NOTES; Phase 08 scaffold + P08-00 started |
| P07-S03-02 | review | pending | Owns DR-HANDOFF completion |

## Checklist

- [x] P07-S03-00 planner — locks + thickened 01/02
- [x] P07-S03-01 verify — create Gate H harness; measure→threshold; VERIFY-NOTES; scaffold or spawn
- [ ] P07-S03-02 review — fresh re-check; complete Phase 08 handoff; mark Phase 07 complete

## Locked Gate H (FINAL P07-S03-00)

```bash
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
```

| Item | Value |
|------|-------|
| Package | `evals/perf` |
| Named test | `TestPlantedPerfLadderGateH` |
| Schema / metrics | `schema-gate-h.json` v1 + temp `metrics-gate-h.json` (`dry_run:false`) |
| Ownership | **S03-01 creates** (S01/S02 left none) |
| Rungs | smoke (~50–200 LOC) + ~1k + ~10k synthetic; **100k/1M deferred** |
| Thresholds | Measure-then-threshold — derive ceilings after first measurements; document in VERIFY-NOTES |
| Pass bar | Named test PASS + schema-valid metrics + structural T0/isolation (+ optional Go) + locked derived ceilings |

## Carry-forward (must stay green)

- S01 T0 + isolation (`cmd/trace` named tests)
- S02 Go golden (`TestIndexFileGoGolden` + DetectLanguage; tree-sitter-go v0.25.0)
- Honesty A/B/C + Gate G (`evals/honesty`)
- Gate E (`evals/replan` `TestPlantedDiscoveryReplan`)
- Gate F (`evals/impact` `TestPlantedImpactConflictsGateFPrelim`)
- Ablation (`evals/capability` `TestPlantedCapabilitySelectionAblation`)
- p0x 7/7; x0; domain/store/planner/compiler; Gate C `dry_run:false` intact
- Full `CGO_ENABLED=1` suite incl. `./evals/perf/...`
- Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H

## Phase 08 handoff

| Item | Value |
|------|-------|
| Folder | `docs/phases/phase-08-ecosystem-hardening/` |
| First board row | `P08-00` → `00-PHASE-PLANNER.md` |
| Who starts | P07-S03-01 on VERIFY pass |
| Who owns completion | P07-S03-02 |
