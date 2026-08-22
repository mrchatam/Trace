# P06-S03-02 — Phase review notes (capability-selection ablation close / DR-HANDOFF)

**Date:** 2026-08-16  
**Verdict:** APPROVE — Phase 06 complete; next runnable `P07-00`  
**Confidence:** **high**  
**Spawns:** none

Independent review of S03 VERIFY (`01-verify.md` + `VERIFY-NOTES.md` + board Notes for `P06-S03-01`). Fresh session ≠ S03-01.

**Explicit:** Capability-selection ablation = `evals/capability` **`TestPlantedCapabilitySelectionAblation`** (planted 4-probe P/R; not vibes / Notes-only / Gate C scores / commercial multi-model). Phase 01 dry-run ≠ Gate C pass, ≠ Gate F pass, ≠ Gate G evidence, and ≠ capability-ablation pass. Gate C **Go** re-confirmed from `dry_run:false` artifacts (mean G1 0.800 > B0 0.000). Mode-B packs remain historical. No new Gate C Go invented. Full commercial multi-model capability benchmark not claimed — ablation = planted automated P/R only. No ontology megastore. No `plan simulate`. DR-HANDOFF closed.

## Plan (executed)

1. Compare VERIFY claims to S01/S02 REVIEW-NOTES + live ablation harness + Gate C metrics
2. Fresh suite re-run: all locked VERIFY commands (ablation + honesty A/B/C + Gate G + Gate E + Gate F + domain/store/planner/compiler + full `./...`)
3. Spot-check planted tallies (TP=3/FN=0/FP=0/TN=1; P/R=1.0) + S01 hooks (`UpsertCapability` / `RequireCapability` / `MissingCapabilities` + packet required/missing)
4. Confirm DR-HANDOFF scaffold under `docs/phases/phase-07-performance-ladder/` + board `P07-00` first pending after Phase 06
5. Carry residuals; write these notes; mark Phase 06 complete

## Claims vs evidence

| Claim (VERIFY-NOTES / P06-S03-01 Notes) | Evidence |
|-----------------------------------------|----------|
| Capability-selection ablation green | Fresh `CGO_ENABLED=0 go test ./evals/capability/... -run TestPlantedCapabilitySelectionAblation` PASS |
| schema-capability.json v1 + temp metrics | Committed schema `schema_version` const 1; harness writes/validates temp `metrics-capability.json` |
| TP=3 / FN=0 / FP=0 / TN=1; precision=1.0; recall=1.0 | Harness `t.Fatalf` tallies + P/R asserts; Pos-1 UNAVAILABLE / Pos-2 UNKNOWN / Pos-3 selection-filter / Neg-1 clean AVAILABLE |
| S01 `UpsertCapability` + `RequireCapability` + `MissingCapabilities` + packet required/missing | Ablation harness plants via S01 APIs + packet fields; mig `010_capability_surface.sql` present; no `internal/capability` |
| Honesty H5 Paths A/B/C | Fresh `TestHonestyFailClosedPlantedClaim` + full `./evals/honesty/...` PASS |
| Gate G prelim | Fresh `TestHonestyEscapeRateGateGPrelim` PASS (escapes=1/caught=2/attempts=3 carry) |
| Gate E mini-eval | Fresh `TestPlantedDiscoveryReplan` PASS |
| Gate F prelim | Fresh `TestPlantedImpactConflictsGateFPrelim` PASS (TP=3/FN=0/FP=0/TN=1; P/R=1.0) |
| S01 mig 010 + capability surface | Fresh `./internal/domain/...` + store + planner + compiler PASS |
| P0-X 7/7 | Fresh `./evals/p0x/...` package PASS under full suite |
| X0 packages | Fresh `./evals/x0/...` PASS |
| Gate C `dry_run:false` intact | `metrics-{b0,g1}.json`: `dry_run=false`, 3 runs each; B0 tools `read_file`/`grep` (no why/context); G1 includes `trace why`/`trace context`; means 0.000 / 0.800; `GATE-C-NOTES.md` still **Go** |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation | VERIFY + this review reject Phase 01 `dry_run:true` as Gate C, Gate F, Gate G, or ablation |
| `go test ./...` | Fresh `CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./... -count=1` EXIT 0 |
| Fixture hash pin | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Law checks | No `ListenAndServe`/`http.Server` under product paths; 0 `.trace/` under fixtures/evals; G19 empty on `internal/` non-test; mig 010 present; no ontology megastore / `plan simulate` |
| Residuals non-blocking | DPC-global; non-tx Apply; UNIQUE re-link; MCP no severity; GC-03/04 deferred; S02 lows — carried |
| DR-HANDOFF complete | See checklist below — not README-only / blocked-until-noticed |

## Re-verification commands (2026-08-16, reviewer)

```text
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
# ok evals/capability 0.027s — EXIT:0

CGO_ENABLED=0 go test ./evals/capability/... -count=1
# ok capability 0.026s — EXIT:0

CGO_ENABLED=0 go test ./evals/honesty/... -count=1
# ok honesty 0.043s — EXIT:0

CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
# ok honesty 0.043s — EXIT:0

CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
# ok evals/replan 0.030s — EXIT:0

CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
# ok evals/impact 0.024s — EXIT:0

CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... -count=1
# ok domain; store; planner; compiler — EXIT:0

CGO_ENABLED=0 go test ./evals/capability/... -count=1 -v -run TestPlantedCapabilitySelectionAblation
# --- PASS: TestPlantedCapabilitySelectionAblation (0.02s)

CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./... -count=1
# ok honesty; p0x 1.872s; x0 1.932s; replan; impact; capability; cmd/trace; internal/* — EXIT:0

test -f evals/capability/schema-capability.json
# schema_version const 1

find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22  -
```

Gate C artifact inspect (no re-score): `dry_run: false` N=3; means match GATE-C-NOTES (B0 0.000 / G1 0.800); B0 no why/context; G1 includes why/context.

## DR-HANDOFF checklist (this row owns completion)

| Item | Status |
|------|--------|
| `docs/phases/phase-07-performance-ladder/README.md` (goal = Performance ladder & language plugins / A_PROJECT_PLAN Phase 7) | **ok** |
| `00-PHASE-PLANNER.md` runnable (session-start + exit criteria) | **ok** |
| Scope stubs S01–S03 each with `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md` | **ok** |
| `docs/TODO.md` Phase 07 section; first pending after Phase 06 last `done` = **`P07-00`** | **ok** (after this row marks `done`) |
| Not README-only / blocked-until-noticed | **ok** — full stub tree + board `P07-00` unblocked for orchestrator |

Do **not** execute `P07-00` in this review.

## Findings

| Severity | Location | Failure mode | Disposition |
|----------|----------|--------------|-------------|
| medium (residual) | retrieval task Expand DPC attach | Every task Expand attaches all `discovery_causes_plan_change` edges | Residual — fine for fixtures/x0; Phase 07+ if multi-goal measured |
| medium (residual) | `planner.ApplyDiscoveryReplan` | Non-tx multi-step window on mid-failure | Residual — Phase 03 S02; not ablation blocker |
| low | Discovery→PlanChange UNIQUE re-link | Re-apply may need idempotent handling | Residual — prior |
| low | MCP discovery path | No severity knob on MCP tools | Residual — CLI/domain path owns severity |
| low | S02 capability schema / probe notes | Schema looseness / probe scoring notes per S02 REVIEW-NOTES | Residual — prior low; not ablation fail |
| nit | GC-03/04 deferred | Recorded-operator-sim; N=3 identical grades | Residual — correct; do not promote without evidence |

No blocker/high. No open medium without follow-up. No spawn.

## Phase close declaration

- **Phase 06 / Environment / capability graph:** complete (S01 capability surface + S02 capability-selection ablation + VERIFY + DR-HANDOFF).  
- **Capability-selection ablation:** planted `evals/capability` `TestPlantedCapabilitySelectionAblation` — green (TP=3/FN=0/FP=0/TN=1; P/R=1.0).  
- **Phase 01 dry-run:** still **not** Gate C pass (and **not** Gate F / Gate G / ablation).  
- **Gate C artifacts:** intact (`dry_run:false`, Go).  
- **Gate F / Gate G / Gate E:** still green.  
- **Board:** all Phase 06 rows `done` after this review marks `P06-S03-02` done.  
- **Next runnable:** **`P07-00`** (`docs/phases/phase-07-performance-ladder/00-PHASE-PLANNER.md`) — pending; do not start until orchestrator launches a fresh subagent for that row.

## Residuals (explicit; do not undermine high confidence)

1. Global DPC attach on every task Expand (medium residual from Phase 02).  
2. Non-tx `ApplyDiscoveryReplan` (medium residual from Phase 03).  
3. UNIQUE re-link / MCP no severity / S02 lows (low).  
4. GC-03/04 deferred; soft p0x decision-constraint OR (prior).  
5. Ontology megastore / VerifiedFact / `plan simulate` remain rejected / out.

None undermine capability-selection ablation, honesty A/B/C, Gate G, Gate E, Gate F, p0x 7/7, x0, S01/S02 surfaces, Gate C artifact integrity, or `./...` on this run.
