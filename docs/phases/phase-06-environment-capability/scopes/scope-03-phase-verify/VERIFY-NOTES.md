# P06-S03-01 — Phase VERIFY notes (capability-selection ablation closeout)

**Date:** 2026-08-16  
**Verifier:** independent re-run (does **not** trust S01/S02 Notes alone)  
**Verdict:** **Phase 06 VERIFY PASS / capability-selection ablation green**  
**Confidence:** high  
**Spawns:** none  

**Explicit claims:** Capability-selection ablation green via `evals/capability` **`TestPlantedCapabilitySelectionAblation`** (planted 4-probe P/R: TP=3 / FN=0 / FP=0 / TN=1; precision=1.0; recall=1.0; schema-valid temp `metrics-capability.json` vs committed `schema-capability.json` v1; S01 hooks `UpsertCapability` + `RequireCapability` + `MissingCapabilities` + compiler packet `required_capabilities`/`missing_capabilities`). Honesty Paths A/B/C (`TestHonestyFailClosedPlantedClaim`) still fail-closed without hatch in that proof. Gate G (`TestHonestyEscapeRateGateGPrelim` escapes=1/caught=2/attempts=3) + Gate E (`TestPlantedDiscoveryReplan`) + Gate F (`TestPlantedImpactConflictsGateFPrelim` TP=3/FN=0/FP=0/TN=1; P/R=1.0) + p0x 7/7 + x0 + domain/store/planner/compiler + full `./...` (incl. capability) PASS. Gate C artifacts remain **Go** (`dry_run:false`, N=3; G1 0.800 > B0 0.000).  
**Explicit non-claims:** Phase 01 dry-run is **not** Gate C pass, **not** Gate F pass, **not** Gate G evidence, and **not** capability-ablation pass. Mode-B packs remain historical. GC-03/04 stay deferred. A1 / product thesis not commercially validated. Full commercial multi-model capability benchmark not claimed — ablation = planted automated P/R only. No ontology megastore. No `plan simulate`. No product feature Go on this row. Phase 06 not marked complete here — **P06-S03-02** owns handoff close + phase complete.

## Environment

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| `go version` | go1.24.2 linux/amd64 |
| Ablation / honesty / Gate E / Gate F / planner / store / domain / compiler | `CGO_ENABLED=0` where locked |
| Full suite / p0x / x0 / analyzers | `CGO_ENABLED=1` |
| Fixture hash (carry) | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Gate C metrics | `docs/verification/gate-c-x0/` (`dry_run:false`, N=3/condition) |
| Migrations | `internal/store/schema/010_capability_surface.sql` (+ 006–009 carry) |
| Ablation schema | `evals/capability/schema-capability.json` (`schema_version` const **1**) |

## Commands (independent)

```text
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
# ok evals/capability 0.032s — EXIT:0

CGO_ENABLED=0 go test ./evals/capability/... -count=1
# ok capability 0.033s — EXIT:0

CGO_ENABLED=0 go test ./evals/honesty/... -count=1
# ok honesty 0.061s — EXIT:0

CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
# ok honesty 0.045s — EXIT:0

CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
# ok evals/replan 0.036s — EXIT:0

CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
# ok evals/impact 0.026s — EXIT:0

CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... -count=1
# ok domain; store; planner; compiler — EXIT:0

CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./... -count=1
# ok honesty; p0x 1.813s; x0 1.888s; replan; impact; capability; cmd/trace; internal/* — EXIT:0

# Optional (strong evidence)
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -v -run TestPlantedCapabilitySelectionAblation
# --- PASS: TestPlantedCapabilitySelectionAblation (0.03s)

test -f evals/capability/schema-capability.json
# present; properties.schema_version const 1

find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22  -
```

No CGO/binary skip treated as pass: harnesses built and ran.

## Evidence table

| Gate | Result (pass/fail) | Evidence (test / log gist) |
|------|--------------------|----------------------------|
| Capability-selection ablation | **pass** | `evals/capability` `TestPlantedCapabilitySelectionAblation` PASS (0.032s) |
| schema-capability.json v1 + temp metrics | **pass** | committed schema `schema_version` const 1; test writes/validates temp `metrics-capability.json` |
| TP=3 FN=0 FP=0 TN=1; P/R=1.0 | **pass** | harness `t.Fatalf` tallies + precision/recall asserts; Pos-1 UNAVAILABLE / Pos-2 UNKNOWN / Pos-3 selection-filter / Neg-1 clean AVAILABLE |
| S01 Upsert/Require/Missing + packet required/missing | **pass** | ablation harness plants via S01 APIs + packet fields; domain/store/compiler package PASS |
| Honesty H5 Paths A/B/C | **pass** | `TestHonestyFailClosedPlantedClaim` PASS; hatch not used in A/B/C proof |
| Gate G prelim | **pass** | `evals/honesty` `TestHonestyEscapeRateGateGPrelim` PASS (escapes=1/caught=2/attempts=3) |
| Gate E mini-eval | **pass** | `evals/replan` `TestPlantedDiscoveryReplan` PASS (0.036s) |
| Gate F prelim | **pass** | `evals/impact` `TestPlantedImpactConflictsGateFPrelim` PASS (TP=3/FN=0/FP=0/TN=1; P/R=1.0) |
| S01 mig 010 + capability surface | **pass** | `./internal/domain/...` + store + compiler PASS; `010_capability_surface.sql` present |
| P0-X 7/7 | **pass** | `evals/p0x` package PASS under full suite |
| X0 packages | **pass** | `./evals/x0/...` PASS |
| Gate C `dry_run:false` intact | **pass** | `metrics-b0.json` + `metrics-g1.json`: `dry_run: false`, 3 runs each; `GATE-C-NOTES.md` still **Go** (G1 0.800 > B0 0.000) |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation | **pass** | Explicit: Phase 01 dry-run **not** used as Gate C, Gate F, Gate G, or capability-ablation pass |
| `go test ./...` | **pass** | Full module `CGO_ENABLED=1` EXIT:0 (incl. capability) |
| Law checks | **pass** | See table below |
| Residuals (non-blocking) | **noted** | DPC-global; non-tx Apply; UNIQUE re-link; MCP no severity; GC-03/04 deferred; S02 lows; ontology megastore rejected |
| DR-HANDOFF | **pass (started)** | `docs/phases/phase-07-performance-ladder/` created; board `P07-00` appended. **P06-S03-02 owns completion check.** |

## Law / architecture checks

| Check | Result | Evidence |
|-------|--------|----------|
| No daemon / always-on HTTP as primary | **pass** | No product `ListenAndServe` / daemon primary; MCP stdio only |
| No committed `.trace/` under `fixtures/` or `evals/` | **pass** | `find` empty |
| G19: libraries do not import `cmd/trace` or `cmd/trace-mcp` | **pass** | `rg` on `internal/` non-test sources empty for cmd imports |
| Ablation evidence is `evals/capability` `TestPlantedCapabilitySelectionAblation` | **pass** | Named test re-run; not Notes-only / vibes / Gate C scores / commercial multi-model |
| Planted tallies TP=3 FN=0 FP=0 TN=1; P/R=1.0 | **pass** | Asserts in harness; Pos UNAVAILABLE/UNKNOWN/selection-filter; Neg AVAILABLE |
| S01 hooks remain green; mig 010; no ontology megastore / `internal/capability` | **pass** | Ablation + domain/store/compiler green; no `internal/capability` package |
| Paths A/B/C fail-closed without hatch in that proof | **pass** | `TestHonestyFailClosedPlantedClaim` PASS |
| Gate G / Gate E / Gate F still green | **pass** | Named tests re-run EXIT:0 |
| Gate C evidence remains `dry_run:false` — not Phase 01 dry-run | **pass** | Metrics + GATE-C-NOTES inspected; no new Go invented |
| Mode-B packs not falsified | **pass** | VERIFY did not rewrite packs; historical |
| Embeddings still absent; VerifiedFact out as promotion engine | **pass** | retrieval doc forbids embeddings; residuals/docs say not VerifiedFact |
| No commercial multi-model capability theater | **pass** | Ablation = planted P/R only |
| GC-03/04 remain deferred | **pass** | Not promoted in these Notes |

## Residuals (non-blocking; carried forward)

1. **Global DPC attach** on every task Expand (Phase 02 residual) — fine for `fixtures/x0`; multi-goal stores may over-attach.  
2. **Non-tx `ApplyDiscoveryReplan`** — partial-failure window if mid-steps fail (Phase 03 S02).  
3. **UNIQUE re-link** on Discovery→PlanChange — re-apply may need idempotent handling.  
4. **MCP no severity** — discovery severity is CLI/`domain` path; MCP tools do not expose severity knobs.  
5. **S02 low residuals** per S02 REVIEW-NOTES (schema looseness / probe scoring notes).  
6. **GC-03/04 deferred** — recorded-operator-sim model; N=3 identical grades.  
7. Soft `decision-constraint` OR / unchecked JSON asserts in p0x (prior phases).  
8. Ontology megastore / VerifiedFact / `plan simulate` remain rejected / out.

None undermine capability-selection ablation, honesty A/B/C, Gate G, Gate E, Gate F, p0x 7/7, x0, S01/S02 surfaces, Gate C artifact integrity, or `./...` on this run.

## DR-HANDOFF progress

Created under `docs/phases/phase-07-performance-ladder/`:

- `README.md` — goal = Performance ladder & language plugins (`A_PROJECT_PLAN` Phase 7)
- `00-PHASE-PLANNER.md` — runnable (Agent→clarify→Plan→execute); light locks OK
- `scopes/scope-01-incremental-indexing/` — 00/01/02 + SCOPE-TODOS
- `scopes/scope-02-language-plugins/` — 00/01/02 + SCOPE-TODOS
- `scopes/scope-03-phase-verify/` — 00/01/02 + SCOPE-TODOS

Board: Phase 07 section with first pending row **`P07-00`**. Do **not** execute Phase 07 until `P06-S03-02` is `done`.

## Board pointer

`P06-S03-01` Notes: ablation + honesty/Gate G/E/F/p0x/x0/S01–S02 PASS; Gate C intact; Phase 07 scaffold started; see this file; pending P06-S03-02 handoff close.
