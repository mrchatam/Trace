# P14-S03-01 — Phase VERIFY notes (peer-impact-install-gates closeout)

**Date:** 2026-08-17  
**Verifier:** independent re-run (does **not** trust S01–S02 Notes alone)  
**Verdict:** **Phase 14 VERIFY PASS / peer-impact-install-gates green**  
**Confidence:** high  
**Spawns:** none  

**Explicit claims:** S01 ImpactWalk named regressions (multi-seed + seed exclusion; contains asymmetry; incoming import hop; loud truncation; hop_risk monotonic) + Gate F planted green; S02 install named (Cursor STABLE detect; uninstall idempotent; CONDITIONAL refuse/write) + `TestInstallCursor*` / usage / ImpactWalk CLI + capability decision named (AUTO_ALLOWED / PENDING fail-closed / ALLOWED persists / DENIED blocks) + capability ablation green. Honesty Paths A/B/C + Gate G + Gate E + Gate F + Gate H + compat (mig ceiling **13**) + p0x 7/7 + x0 + supporting domain/store/planner/compiler/retrieval/install/mcp + `cmd/trace` + product packages PASS. Gate C artifacts remain **Go** (`dry_run:false`, N=3; mean G1 **0.800** > B0 **0.000**).

**Explicit non-claims:** Phase 01 dry-run is **not** Gate C, **not** Gate F, **not** Gate G, **not** ablation, **not** Gate H, and **not** the compat checklist. VerifiedFact still out. No install/decide MCP dump. No product Go on this row. No Phase 15 / S05 / `plan simulate` / D21+ scaffold. Phase 14 not marked complete here — **P14-S03-02** owns handoff close + phase complete. S02 APPROVE ≠ “every MCP call gated.”

**DR-HANDOFF = `no successor`.** Parallel dogfood / research FUTURE may continue under `experiments/` and research docs off-board. Do **not** scaffold Phase 15 / S05 / plan simulate / D21+ unless Notes explicitly promote.

## Environment

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| `go version` | go1.24.2 linux/amd64 |
| Analyzers / Gate H / compat / p0x / x0 / `cmd/trace` / full product bar | `CGO_ENABLED=1` |
| Retrieval / install / domain / honesty / Gate E / Gate F / ablation / support pkgs | `CGO_ENABLED=0` where locked |
| Gate C metrics | `docs/verification/gate-c-x0/` (`dry_run:false`, N=3/condition) |
| Gate C means (inspect only) | B0 mean **0.000**; G1 mean **0.800** — **not** re-scored |
| Optional dogfood | `experiments/` **not** run this row (non-blocking) |

## Evidence table (independent)

| Bucket / command | Result |
|------------------|--------|
| S01 `CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestImpactWalk'` | **PASS** (multi-seed/exclusion, contains asymmetry, incoming import, loud truncation, hop_risk, CLI helpers as packaged) |
| S01/F `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` | **PASS** (Gate F planted) |
| S02 install `CGO_ENABLED=0 go test ./internal/install/... -count=1 -run 'TestInstallDetectListsCursorStable\|TestInstallCursorUninstallIdempotent\|TestInstallConditional'` | **PASS** |
| S02 decisions `CGO_ENABLED=0 go test ./internal/domain/... -count=1 -run 'TestCapabilityDecision'` | **PASS** (AUTO_ALLOWED / PENDING / ALLOWED / DENIED) |
| S02 CLI `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallCursor\|TestInstallUsage\|TestImpactWalkCLI'` | **PASS** |
| Ablation `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` | **PASS** |
| Honesty full `CGO_ENABLED=0 go test ./evals/honesty/... -count=1` | **PASS** |
| Honesty A/B/C + G `CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim\|TestHonestyEscapeRateGateGPrelim'` | **PASS** |
| Gate E `CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan` | **PASS** |
| Gate F (carry) `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` | **PASS** |
| Ablation (carry) `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` | **PASS** |
| Gate H `CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH` | **PASS** (~5.2s named) |
| Compat `CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist` | **PASS** (mig ceiling **13**; saw 013; no 014+) |
| P0-X + X0 `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1` | **PASS** — p0x 7/7 (`TestP0XAllCriteria`) + x0 |
| Support `CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... ./internal/retrieval/... ./internal/install/... -count=1` | **PASS** |
| MCP `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... -count=1` | **PASS** |
| `CGO_ENABLED=1 go test ./cmd/trace/... -count=1` | **PASS** |
| Product `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1` | **PASS** |
| `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./... -count=1` | Product pkgs **PASS**; known FAIL only `similar projects/graphify` space-in-path (non-product residual; exit 1) |
| MCP nine tools + `trace_version` / no install-decide MCP | **PASS** — `RegisteredToolNames` = nine + `trace_version`; `TestToolNamesRegistered` PASS; no install/decide tool Names |
| Gate C artifacts inspect | **PASS** — `dry_run:false`; N=3; mean G1 0.800 > B0 0.000; **not** re-scored |
| Mig 013 / no VERIFY mig | **PASS** — `013_capability_tool_decisions.sql` present; no `014_*`; VERIFY added no migration |
| No committed `.trace/` under `fixtures/` / `evals/` | **PASS** |
| G19 library packages do not import `cmd/trace` | **PASS** (`go list` / deps clean) |
| No Phase 15 scaffold | **PASS** — no `docs/phases/phase-15*` |

## Law checks

| Check | Hold? |
|-------|-------|
| No daemon / always-on HTTP as primary surface | Yes |
| No committed `.trace/` under `fixtures/` or `evals/` | Yes |
| G19 — library packages do not import `cmd/trace` or `cmd/trace-mcp` | Yes |
| S01–S02 evidence is **named tests** — not Notes-only | Yes |
| MCP remains **nine** tools + `trace_version`; no install/decide MCP | Yes |
| Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat green | Yes |
| Gate C evidence remains `dry_run:false` — not Phase 01 dry-run | Yes |
| Embeddings / VerifiedFact / Neo4j SoT still out | Yes |
| No full-rebuild-on-any-change indexer architecture | Yes |
| No new migration from VERIFY; mig 013 already in S02; compat ceiling **13** | Yes |
| No YOLO / AllowAll defaults | Yes |
| Assert ≠ MCP dispatch — honesty Note only (not a VERIFY fail) | Yes |
| **No Phase 15 / S05 / plan simulate / D21+ scaffold** | Yes |
| Forward-only: do **not** rewrite Phase 00–13 `done` history; Phase 13 historical `no successor` left intact | Yes |

## Residuals / deferrals

- **Assert ≠ MCP dispatch (by design):** `AssertToolAllowed` lives in `internal/domain` + CLI/tests only — **not** wired into `internal/mcp` request dispatch. S02 APPROVE proves audit/CLI fail-closed; it does **not** mean every MCP call is gated. Wiring Assert into MCP stays deferred unless separately promoted.
- **Optional `allowContainsOut` late-upgrade:** Spot-checked in `internal/retrieval/impact_walk.go` (file→symbols OK; `false` after contains-UP / no sibling climb; late-upgrade path at enqueue still present). Non-blocking residual from P14-S01-02; **not** a Phase 15 trigger.
- **Goals #2–#4** (S05 supersession / `plan simulate` / D21+) stay **off-board**.
- **Known `./...` nit:** `similar projects/graphify` space-in-path setup FAIL — pre-existing non-product; product pkgs PASS.
- **CGO0 analyzers FAIL OK** residual if present on zero-CGO analyzer path — product bar uses CGO1 for analyzers (PASS under product suite).
- Parallel dogfood under `experiments/` — **not** board-blocking; **not** run this row.
- VerifiedFact / embeddings / daemon-HTTP primary still out.

## Dry-run ≠ gates

**Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist.** Gate C artifacts remain Mode-B `dry_run:false` (inspect only; not re-scored). Phase 01 dry-run is regression-only.

## Handoff

| Item | Value |
|------|-------|
| **DR-HANDOFF** | **`no successor`** (**started** P14-S03-01; close owned by **P14-S03-02**) |
| Phase 15 / S05 / plan simulate / D21+ | **Do not scaffold** — intentional absence (no promotion) |
| Parallel dogfood / research FUTURE | May continue off-board under `experiments/` / research docs |
| Completion owner | **P14-S03-02** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff explicitly `no successor` |
