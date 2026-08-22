# S04 — Phase VERIFY — scope todos

**Depends-on:** P13-S03-02 APPROVE. Owns phase gate + DR-HANDOFF.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** (FINAL 2026-08-17) |
| 2 | 01-verify | verify | pending — next |
| 3 | 02-scope-review | review | pending |

## FINAL locks (P13-S04-00)

| Item | Value |
|------|-------|
| Gate | Phase 13 closeout — DF-60…67 named + P12 keepers + carry-forward |
| DR-HANDOFF | **`no successor`** — S04-01 starts; S04-02 owns completion |
| S01 named | Import-path candidates + Expand/Why subdir/parent/extensionless/root (+ P12 Expand/Why) |
| S02 named | `TestIndexHonestyStaleTotalTruncated` / `TestIndexHonestyPreTrimUniverse` / `TestCandidateCapAdmitUniverseTotal` / `TestContextImportHopEdgeProvenance` (+ P12 Budget/cap/stale/WhyTrace) |
| S03 named | Garbage reject / empty normalize / Expand empty→EXTRACTED / round-trip / analyzer EXTRACTED; mig **012**; compat **12** |
| DF-66 | Documented **wontfix** — confirm docs + Law 5 fixtures in VERIFY-NOTES |
| DF-67 | **Out-of-bar** residual must appear in VERIFY-NOTES (`symstale/`); file-hash only |
| SchemaVersion | Remains **`0.2`** |
| Optional dogfood | `experiments/ab-import-resolve/` prepare + probe — **non-blocking** (≠ Gate C) |
| Spawn | 01a/01b/01c on fail |
| Forbidden | Phase 14 scaffold without promotion; inventing INFERRED CLI/analyzer; inventing symbol honesty; product Go on VERIFY |

## Depends (from S03 APPROVE)
- Named **DF-64** regressions: `TestReplaceFileImportsRejectsGarbageProvenance` / `TestImportProvenanceEmptyWriteAndReadNormalize` / `TestExpandEmptyProvenanceSurfacesExtracted` + P12 keepers; mig **012** / compat ceiling **12** (no 013+).
- **DF-66** wontfix evidence: `ANALYZER_CONTRIBUTION.md` § Import edge provenance + Law 5 store-fixture Expand/Why/compiler tests green.
- **DF-67** residual **must** appear in VERIFY-NOTES: symbol-entity staleness still out of `index_honesty` bar (`experiments/_bughunt/post-p12/symstale/`).
- Packet SchemaVersion stays **`0.2`**.
- S03 REVIEW-NOTES: [../scope-03-provenance-schema/REVIEW-NOTES.md](../scope-03-provenance-schema/REVIEW-NOTES.md)

## Reminders
- Default DR-HANDOFF = `no successor` (**FINAL** — no Phase 14 unless Notes promote)
- Optional ab-import-resolve prepare is dogfood hook, not Gate C
- S04-01 starts handoff; S04-02 closes
- Dry-run ≠ Gate C / ≠ Gate H / ≠ checklist
