# P07-S03-01 — Phase VERIFY notes (Gate H / performance ladder closeout)

**Date:** 2026-08-16  
**Verifier:** independent re-run (does **not** trust S01/S02 Notes alone)  
**Verdict:** **Phase 07 VERIFY PASS / Gate H green**  
**Confidence:** high  
**Spawns:** none  

**Explicit claims:** Gate H green via `evals/perf` **`TestPlantedPerfLadderGateH`** (planted synthetic ladders smoke/~1k/~10k; schema-valid temp `metrics-gate-h.json` vs committed `schema-gate-h.json` v1; `dry_run:false`; structural `t0_skip_ok` + `incremental_isolation_ok` + `go_adapter_exercised`; regression ceilings **measure-then-threshold** locked). S01 T0 + isolation (`TestWalkIndexableT0AlwaysSkip` / `TestIndexSkipsExplicitT0Path` / `TestIndexIncrementalIsolation`) + S02 Go golden (`TestIndexFileGoGolden` + DetectLanguage; `tree-sitter-go` v0.25.0) re-proved. Honesty Paths A/B/C (`TestHonestyFailClosedPlantedClaim`) + Gate G (`TestHonestyEscapeRateGateGPrelim` escapes=1/caught=2/attempts=3) + Gate E (`TestPlantedDiscoveryReplan`) + Gate F (`TestPlantedImpactConflictsGateFPrelim` TP=3/FN=0/FP=0/TN=1; P/R=1.0) + capability ablation (`TestPlantedCapabilitySelectionAblation` TP=3/FN=0/FP=0/TN=1; P/R=1.0) + p0x 7/7 + x0 + domain/store/planner/compiler + full `./...` (incl. perf) PASS. Gate C artifacts remain **Go** (`dry_run:false`, N=3; G1 0.800 > B0 0.000).  

**Explicit non-claims:** Phase 01 dry-run is **not** Gate C pass, **not** Gate F, **not** Gate G, **not** ablation, and **not** Gate H. Mode-B packs remain historical. GC-03/04 stay deferred. **100k / 1M** planted CI ladders deferred. A1 / product thesis not commercially validated. No commercial multi-model perf theater. No product feature Go outside `evals/perf/**` on this row. Phase 07 not marked complete here — **P07-S03-02** owns handoff close + phase complete.

## Environment

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| `go version` | go1.24.2 linux/amd64 |
| Gate H / S01 / S02 / p0x / x0 / full suite | `CGO_ENABLED=1` |
| Honesty / Gate E / Gate F / ablation / domain/store/planner/compiler | `CGO_ENABLED=0` where locked |
| Fixture hash (carry) | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Gate C metrics | `docs/verification/gate-c-x0/` (`dry_run:false`, N=3/condition) |
| Gate H schema | `evals/perf/schema-gate-h.json` (`schema_version` const **1**) |
| `tree-sitter-go` | v0.25.0 (go.mod) |

## Measure-then-threshold derivation

**Protocol:** first measure → `ceiling_ms = max(measured_ms * 5, 2000)` / `ceiling_db = measured_db_bytes * 3` → encode in harness → re-run PASS.

### First measurement (2026-08-16, before ceiling lock)

| Rung | approx_loc | file_count | initial_index_ms | incremental_index_ms | db_bytes |
|------|------------|------------|------------------|----------------------|----------|
| smoke | 129 | 6 | 38 | 42 | 454656 |
| rung-1k | 1002 | 12 | 113 | 166 | 864256 |
| rung-10k | 10002 | 52 | 1271 | 2158 | 4902912 |

### Locked ceilings (encoded in `evals/perf/perf_test.go`)

| Rung | initial_ms ceil | incr_ms ceil | db_bytes ceil | derivation |
|------|-----------------|--------------|---------------|------------|
| smoke | 2000 | 2000 | 1363968 | max(38×5,2000); max(42×5,2000); 454656×3 |
| rung-1k | 2000 | 2000 | 2592768 | max(113×5,2000); max(166×5,2000); 864256×3 |
| rung-10k | 6355 | 10790 | 14708736 | max(1271×5,2000); max(2158×5,2000); 4902912×3 |

**Re-run with ceilings locked:** PASS (`TestPlantedPerfLadderGateH` ~3.9–7.7s across runs). Structural: `t0_skip_ok=true`, `incremental_isolation_ok=true`, `go_adapter_exercised=true` (smoke `.go`).

## Commands (independent)

```text
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
# ok evals/perf — EXIT:0 (ceilings locked)

CGO_ENABLED=1 go test ./evals/perf/... -count=1
# ok evals/perf — EXIT:0

CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestWalkIndexableT0AlwaysSkip|TestIndexIncrementalIsolation|TestIndexSkipsExplicitT0Path'
# ok cmd/trace — EXIT:0

CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestIndexFileGoGolden|TestDetectLanguage'
# ok analyzers — EXIT:0

CGO_ENABLED=0 go test ./evals/honesty/... -count=1
# ok honesty — EXIT:0

CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
# ok honesty — EXIT:0

CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
# ok replan — EXIT:0

CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
# ok impact — EXIT:0

CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
# ok capability — EXIT:0

CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... -count=1
# ok domain; store; planner; compiler — EXIT:0

CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./evals/perf/... ./... -count=1
# ok all packages incl. perf/p0x/x0/cmd/trace/analyzers — EXIT:0

# Optional
test -f evals/perf/schema-gate-h.json
# present; schema_version const 1

find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22  -
```

No CGO/binary skip treated as pass: harnesses built and ran.

## Evidence table

| Gate | Result (pass/fail) | Evidence (test / log gist) |
|------|--------------------|----------------------------|
| Gate H harness created (`evals/perf`) | **pass** | `doc.go` + `schema-gate-h.json` + `perf_test.go`; synthetic plants in-test |
| `TestPlantedPerfLadderGateH` | **pass** | named test PASS (~4–8s) |
| `schema-gate-h.json` v1 + temp `metrics-gate-h.json` | **pass** | committed schema; temp metrics schema-validated; `dry_run:false` |
| Measure-then-threshold derivation | **pass** | first measure → ceilings above; encoded in harness; re-run PASS |
| Rungs smoke / ~1k / ~10k | **pass** | 129 / 1002 / 10002 approx_loc; ms+db recorded |
| T0 skip under plant | **pass** | `t0_skip_ok`; node_modules/vendor/`.min.js` not indexed |
| Incremental isolation under plant | **pass** | `incremental_isolation_ok`; iso/b.js stable after iso/a.js reindex |
| Optional Go fixtures / `go_adapter_exercised` | **pass** | smoke `pkg/plant.go` indexed + symbols |
| S01 `TestWalkIndexableT0AlwaysSkip` | **pass** | cmd/trace |
| S01 `TestIndexIncrementalIsolation` | **pass** | cmd/trace |
| S01 `TestIndexSkipsExplicitT0Path` | **pass** | cmd/trace |
| S02 `TestIndexFileGoGolden` / DetectLanguage | **pass** | analyzers; tree-sitter-go v0.25.0 |
| Honesty H5 Paths A/B/C | **pass** | `TestHonestyFailClosedPlantedClaim` |
| Gate G prelim | **pass** | `TestHonestyEscapeRateGateGPrelim` (escapes=1/caught=2/attempts=3) |
| Gate E mini-eval | **pass** | `TestPlantedDiscoveryReplan` |
| Gate F prelim | **pass** | `TestPlantedImpactConflictsGateFPrelim` (TP=3/FN=0/FP=0/TN=1; P/R=1.0) |
| Capability ablation | **pass** | `TestPlantedCapabilitySelectionAblation` (TP=3/FN=0/FP=0/TN=1; P/R=1.0) |
| P0-X 7/7 | **pass** | `evals/p0x` package PASS under full suite |
| X0 packages | **pass** | `./evals/x0/...` PASS |
| Gate C `dry_run:false` intact | **pass** | metrics-b0/g1: `dry_run:false`, 3 runs; GATE-C-NOTES still **Go** (G1 0.800 > B0 0.000) |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H | **pass** | Explicit: Phase 01 dry-run **not** used as any of these |
| `go test ./...` (+ perf) | **pass** | Full module `CGO_ENABLED=1` EXIT:0 |
| Law checks | **pass** | See table below |
| Residuals (non-blocking) | **noted** | DPC-global; GC-03/04; A5; 100k/1M deferred; S01/S02 lows |
| DR-HANDOFF | **pass (started)** | `docs/phases/phase-08-ecosystem-hardening/` created; board `P08-00` appended. **P07-S03-02 owns completion check.** |

## Law / architecture checks

| Check | Result | Evidence |
|-------|--------|----------|
| No daemon / always-on HTTP as primary | **pass** | No product `ListenAndServe` / daemon primary; MCP stdio only |
| No committed `.trace/` under `fixtures/` or `evals/` | **pass** | `find` empty |
| G19: libraries do not import `cmd/trace` or `cmd/trace-mcp` | **pass** | `rg` on `internal/` non-test sources empty |
| Gate H evidence is `evals/perf` `TestPlantedPerfLadderGateH` | **pass** | Named test + schema/metrics; not Notes-only / S01 `t.Logf` / commercial theater |
| Thresholds derived from measurements | **pass** | Documented above; constants in harness |
| Ladder synthetic planted smoke/1k/10k; 100k/1M deferred | **pass** | Harness rungs only those three |
| S01 T0 + isolation / S02 Go golden green | **pass** | Named tests re-run |
| Honesty A/B/C + Gate G/E/F + ablation green | **pass** | Named tests re-run |
| Gate C `dry_run:false` — not Phase 01 dry-run | **pass** | Artifacts + GATE-C-NOTES inspected; no new Go invented |
| Mode-B packs not falsified | **pass** | VERIFY did not rewrite packs |
| Embeddings / VerifiedFact / `plan simulate` still out | **pass** | No promotion this row |
| No full-rebuild-on-any-change | **pass** | File-local IndexFile path; isolation asserts |
| GC-03/04 remain deferred | **pass** | Not promoted |

## Residuals (non-blocking; carried forward)

1. **Global DPC attach** on every task Expand (Phase 02).  
2. **Non-tx `ApplyDiscoveryReplan`** partial-failure window (Phase 03).  
3. **UNIQUE re-link** on Discovery→PlanChange.  
4. **MCP no severity** for discoveries.  
5. **GC-03/04 deferred**.  
6. **A5 SQLite ceiling** still ACCEPTED_RISK until larger plants.  
7. **100k / 1M planted CI ladders** deferred (not Gate H pass bar).  
8. S01 low: explicit `.min.*` argv plant coverage soft.  
9. S02 low: blank/dot import nit; no CLI walk golden for `.go` (smoke Gate H exercises Go).  

None undermine Gate H, S01/S02 surfaces, honesty A/B/C, Gate G/E/F, ablation, p0x 7/7, x0, Gate C integrity, or `./...` on this run.

## DR-HANDOFF progress

Created under `docs/phases/phase-08-ecosystem-hardening/`:

- `README.md` — goal = Ecosystem & hardening (`A_PROJECT_PLAN` Phase 8)
- `00-PHASE-PLANNER.md` — runnable (Agent→clarify→Plan→execute); light locks OK
- `scopes/scope-01-plugin-apis/` — 00/01/02 + SCOPE-TODOS (minimal stub)

Board: Phase 08 section with first pending row **`P08-00`**. Do **not** execute Phase 08 until `P07-S03-02` is `done`.

## Board pointer

`P07-S03-01` Notes: Gate H + S01/S02 + honesty/Gates/ablation/p0x/x0 PASS; Gate C intact; Phase 08 scaffold started; see this file; pending P07-S03-02 handoff close.
