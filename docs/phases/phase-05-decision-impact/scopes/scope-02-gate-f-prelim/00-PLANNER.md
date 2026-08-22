# P05 / S02 / 00-PLANNER — Gate F prelim

## Metadata
- id: P05-S02-00
- todo_ids: [P05-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-gate-f-prelim.md` for **Gate F preliminary** (impact precision/recall on planted conflicts). Lock harness path under `evals/impact`, named test, schema/metrics, and exit criteria. No product code in this planner row.

## References
- [phase README](../../README.md) — Gate F = planted `evals/impact`
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) Gate F
- Prior planted patterns: `evals/replan` Gate E; `evals/honesty` Gate G
- S01 APPROVE surface (impact classes / report hooks)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Depends
S01 impact surface APPROVE (`P05-S01-02`).

### Expected S01 hooks (confirmed live after P05-S01-02 APPROVE — 2026-08-16)
| Hook | Use in Gate F plant |
|------|---------------------|
| `CreateDecision` / `LinkDecisionTask` | Seed decision ↔ task graph (`decision_affects_task` only) |
| `AddImpactFinding` | Plant class/uncertainty/kind rows (incl. UNKNOWN / conflicting classes) |
| `AddDecisionAlternative` / `SetRecommendedAlternative` | Optional alternative/recommended cases (not required in locked 4-probe fixture) |
| `ImpactReport` → `ImpactReportResult` | Score `OverallClass`, `HasUnknown`, `Incomplete`, Findings (UNKNOWN never omitted); do not trust `OverallClass` alone when `HasUnknown` |
| Mig `009_decision_impact` | Tables must exist; no S02 schema fork |
| CLI `trace impact` | Optional for humans; **harness must call domain APIs** (G19), not scrape CLI only |

## Phase Gate F lock (P05-S02-00 FINAL — 2026-08-16)
| Item | Value |
|------|-------|
| Package | **`evals/impact`** (new; do not overload honesty/replan/x0/p0x) |
| Named test | **`TestPlantedImpactConflictsGateFPrelim`** |
| Schema | **`evals/impact/schema-gate-f.json`** v1 (`schema_version` const 1) |
| Metrics artifact | Temp **`metrics-gate-f.json`** (schema-validated in test) |
| Evidence shape | 4 planted probes → TP=3 / FN=0 / FP=0 / TN=1 → **precision=1.0**, **recall=1.0** |
| Probes | Pos-1 UNKNOWN; Pos-2 SAFE+DESTRUCTIVE rollup; Pos-3 link/no findings; Neg-1 clean SAFE KNOWN |
| Non-claim | Not commercial multi-model Gate F; not Gate C / Gate G / Gate E substitute |
| Carry-forward | Honesty A/B/C; Gate G; Gate E; p0x; x0; Gate C `dry_run:false` intact |
| Re-prove (S03) | `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` |

## Planner work
- [x] Confirm S01 hooks available for planting (classes / links / report)
- [x] Lock Gate F named test + schema/metrics filenames
- [x] Thicken `01-gate-f-prelim.md`
- [x] Thicken `02-scope-review.md`
- [x] Light S03 Depends note (VERIFY re-prove command)
- [x] Sync SCOPE-TODOS + board

## Exit criteria
- [x] Gate F path locked with planted evidence shape
- [x] `01-gate-f-prelim.md` runnable alone
- [x] No product Go in this row

## Out of scope
- Commercial multi-model Gate F cost explosion; weakening Gate G/E/C bars; inventing Gate F from vibes/Notes-only
