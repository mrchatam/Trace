# P15-S02-01 — Phase VERIFY notes (P14 residual remediation closeout)

**Date:** 2026-08-17  
**Verifier:** independent re-run (does **not** trust S01 Notes alone)  
**Verdict:** **Phase 15 VERIFY PASS / P14 residual remediation green**  
**Confidence:** high  
**Spawns:** none  

**Explicit claims:** S01 MCP Assert named regressions green — `TestMCPAssertDeniedBlocksCallTool` (DENIED fail-closed), `TestMCPAssertBuiltinAutoAllowedSucceeds` (builtin AUTO_ALLOWED), `TestToolNamesRegistered` (exactly nine tools), `TestBuiltinMCPCapabilitySpecs` / capability decision package run. Wire-up: nine `assertMCPToolAllowed` call sites → `AssertToolAllowed(ctx, "mcp:"+name)` incl. `toolVersion` openStore+Assert. Carry-forward honesty Paths A/B/C + Gate G + Gate E + Gate F + capability ablation + Gate H + compat (mig ceiling **13**) + p0x 7/7 + x0 + product `./cmd|internal|evals` (CGO1) PASS. Gate C artifacts remain **Go** (`dry_run:false`, N=3; mean G1 **0.800** > B0 **0.000**).

**Explicit non-claims:** R2 / R3 / R4 **not** fixed. Phase 01 dry-run is **not** Gate C, **not** Gate F, **not** Gate G, **not** ablation, **not** Gate H, and **not** the compat checklist. VerifiedFact still out. No install/decide MCP dump. No product Go on this row. No Phase 16 / S05 / `plan simulate` / D21+ scaffold. Phase 15 not marked complete here — **P15-S02-02** owns handoff close + phase complete. Phase 14 historical `no successor` left intact as history.

**DR-HANDOFF = `no successor`.** Parallel dogfood / research FUTURE may continue under `experiments/` and research docs off-board. Do **not** scaffold Phase 16 / S05 / plan simulate / D21+ unless Notes explicitly promote.

## Environment

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` (`go.mod`) |
| `go version` | go1.24.2 linux/amd64 |
| Named MCP Assert / honesty / E / F / ablation | `CGO_ENABLED=0` |
| Gate H / compat / p0x / x0 / product bar | `CGO_ENABLED=1` |
| Product bar env | `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off` |
| Gate C metrics | `docs/verification/gate-c-x0/` (`dry_run:false`, N=3/condition) |
| Gate C means (inspect only) | B0 mean **0.000**; G1 mean **0.800** — **not** re-scored |
| Optional dogfood | `experiments/` **not** run this row (non-blocking) |

## Evidence table (independent)

| Bucket / command | Result |
|------------------|--------|
| S01 named `CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered\|TestMCPAssert\|TestBuiltinMCPCapabilitySpecs\|TestCapabilityDecision'` | **PASS** — mcp ok; domain ok; store [no tests to run] |
| Honesty full `CGO_ENABLED=0 go test ./evals/honesty/... -count=1` | **PASS** |
| Honesty A/B/C + G `CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim\|TestHonestyEscapeRateGateGPrelim'` | **PASS** |
| Gate E `CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan` | **PASS** |
| Gate F `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` | **PASS** |
| Ablation `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` | **PASS** |
| Gate H `CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH` | **PASS** (~7.3s) |
| Compat `CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist` | **PASS** (mig ceiling **13**; no 014+) |
| P0-X + X0 `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1` | **PASS** — p0x + x0 |
| Product `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1` | **PASS** — all product pkgs ok (incl. analyzers under CGO1) |
| S01 wire-up grep `assertMCPToolAllowed` | **PASS** — 9 call sites: why/context/add/link/transition/review/tasks/capability/version (+ helper in `assert.go`) |
| MCP nine tools + `trace_version` / no install-decide MCP | **PASS** — named registry tests green; no install/decide tools |
| Gate C artifacts inspect | **PASS** — `dry_run:false`; N=3; mean G1 0.800 > B0 0.000; **not** re-scored |
| Mig 013 / no VERIFY mig | **PASS** — no new migration from this row; ceiling **13** |
| No committed `.trace/` under `fixtures/` / `evals/` | **PASS** |
| G19 library packages do not import `cmd/trace` / `cmd/trace-mcp` | **PASS** |
| No Phase 16 scaffold | **PASS** — intentional absence |

## Law checks

| Check | Hold? |
|-------|-------|
| No daemon / always-on HTTP as primary surface | Yes |
| No committed `.trace/` under `fixtures/` or `evals/` | Yes |
| G19 — library packages do not import `cmd/trace` or `cmd/trace-mcp` | Yes |
| S01 evidence is **named tests** — not Notes-only | Yes |
| MCP remains **nine** tools + `trace_version`; no install/decide MCP | Yes |
| Assert **is** on MCP dispatch (R1) — DENIED + AUTO_ALLOWED named tests green | Yes |
| Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat green | Yes |
| Gate C evidence remains `dry_run:false` — not Phase 01 dry-run | Yes |
| Embeddings / VerifiedFact / Neo4j SoT still out | Yes |
| No full-rebuild-on-any-change indexer architecture | Yes |
| No new migration from VERIFY; mig 013 already in tree; compat ceiling **13** | Yes |
| No YOLO / AllowAll defaults | Yes |
| **R2/R3/R4 do not fail VERIFY** (defer / wontfix) | Yes |
| **No Phase 16 / S05 / plan simulate / D21+ scaffold** | Yes |
| Forward-only: do **not** rewrite Phase 00–14 `done` history; Phase 14 historical `no successor` left intact | Yes |

## Residuals / deferrals (non-blocking)

| Residual | Disposition | VERIFY note |
|----------|-------------|-------------|
| **R2** `allowContainsOut` late-upgrade | **defer** | Spot-checked still present in `internal/retrieval/impact_walk.go` — **not** fail criteria; **not** claimed fixed |
| **R3** graphify space-in-path on full `./...` | **wontfix** | Product bar is `./cmd\|internal\|evals` (PASS) — do **not** fail VERIFY |
| **R4** CGO0 analyzers FAIL | **wontfix** | Product bar is CGO1 (analyzers PASS under product suite) — do **not** fail VERIFY |
| Goals #2–#4 | deferred | S05 / `plan simulate` / D21+ stay **off-board** |
| Parallel dogfood | off-board | `experiments/` may continue; **not** board-blocking; **not** run this row |
| VerifiedFact / embeddings / daemon-HTTP primary | out | unchanged |

## Dry-run ≠ gates

**Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist.** Gate C artifacts remain Mode-B `dry_run:false` (inspect only; not re-scored). Phase 01 dry-run is regression-only.

## Handoff

| Item | Value |
|------|-------|
| **DR-HANDOFF** | **`no successor`** (**started** P15-S02-01; close owned by **P15-S02-02**) |
| Phase 16 / S05 / plan simulate / D21+ | **Do not scaffold** — intentional absence (no promotion) |
| Parallel dogfood / research FUTURE | May continue off-board under `experiments/` / research docs |
| Completion owner | **P15-S02-02** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff explicitly `no successor` |
| Spawns | **none** |
| Next board | **P15-S02-02** |
