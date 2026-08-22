# Scope S03 — Phase 06 VERIFY

**Depends-on:** S01 + S02 APPROVE (**P06-S02-02 APPROVE high** — 2026-08-16).

**DR-HANDOFF:** On pass — VERIFY starts next-phase scaffold **`phase-07-performance-ladder`** + board `P07-00`; S03-02 owns completion (or records `no successor`).

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P06-S03-00 | planner | done | 2026-08-16: locked VERIFY commands + evidence table + Phase 07 checklist; thickened 01+02 |
| P06-S03-01 | verify | done | 2026-08-16: **VERIFY PASS / capability-selection ablation green**; Phase 07 scaffold started; see [VERIFY-NOTES.md](VERIFY-NOTES.md); pending S03-02 |
| P06-S03-02 | review | done | 2026-08-16: APPROVE high; DR-HANDOFF complete; Phase 06 complete; next P07-00 — [REVIEW-NOTES.md](REVIEW-NOTES.md) |

## Checklist

- [x] P06-S03-00 planner — locks + thickened 01/02
- [x] P06-S03-01 verify — run locked commands; VERIFY-NOTES; scaffold or spawn
- [x] P06-S03-02 review — fresh re-check; complete Phase 07 handoff; mark Phase 06 complete

## Locked ablation re-prove (FINAL P06-S02-00 / S03-00)

```bash
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
```

| Item | Value |
|------|-------|
| Package | `evals/capability` |
| Named test | `TestPlantedCapabilitySelectionAblation` |
| Schema / metrics | `schema-capability.json` v1 + temp `metrics-capability.json` |
| Tallies | TP=3 FN=0 FP=0 TN=1; precision=1.0; recall=1.0 |
| Probes | Pos UNAVAILABLE / UNKNOWN / selection-filter; Neg AVAILABLE |

## Carry-forward (must stay green)

- Honesty A/B/C + Gate G (`evals/honesty`)
- Gate E (`evals/replan` `TestPlantedDiscoveryReplan`)
- Gate F (`evals/impact` `TestPlantedImpactConflictsGateFPrelim`)
- p0x 7/7; x0; domain/store/planner/compiler; Gate C `dry_run:false` intact
- Full `CGO_ENABLED=1` suite incl. `./evals/capability/...`
- Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation

## Phase 07 handoff

| Item | Value |
|------|-------|
| Folder | `docs/phases/phase-07-performance-ladder/` |
| First board row | `P07-00` → `00-PHASE-PLANNER.md` |
| Who starts | P06-S03-01 on VERIFY pass (**done** — scaffold present) |
| Who owns completion | P06-S03-02 |
