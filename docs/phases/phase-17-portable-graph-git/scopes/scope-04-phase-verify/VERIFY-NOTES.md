# P17-S04-01 — Phase VERIFY notes (portable-graph-git closeout)

**Date:** 2026-08-17  
**Verifier:** independent re-run (does **not** trust S01–S03 Notes alone)  
**Verdict:** **Phase 17 VERIFY PASS / portable-graph-git green**  
**Confidence:** high  
**Spawns:** none  

**Explicit claims:** S01–S03 named DF regressions green on live packages. Two-clone git-JSON recipe **`TestPortableGraphTwoCloneWhyContextPlan`** PASS (implemented at VERIFY preflight — was absent). Carry-forward honesty Paths A/B/C + Gate G + Gate E + Gate F + capability ablation + Gate H + compat + p0x + x0 + product `./cmd|internal|evals` (CGO1) PASS. Gate C artifacts remain **Go** (`dry_run:false`, N=3). `.gitignore` lists `.trace/` only; `trace/graph.json` not ignored. DF-86 git-hook **absent** (non-fail; grep recorded).

**Explicit non-claims:** Encryption-as-git wontfix; reviews omitted from default export; DF-86 hook not implemented; CGO=0 `cmd/trace` tree-sitter FAIL; S03 `work_state` preservation has no dedicated named re-import test (SQL-only). Phase 01 dry-run is **not** Gate C, **not** Gate F, **not** Gate G, **not** ablation, **not** Gate H, and **not** the compat checklist. No research S05 / `plan simulate` / D21+ / hosted MCP scaffold. Phase 16 historical `no successor` left intact. Phase 17 not marked complete here — **P17-S04-02** owns handoff close + phase complete.

**DR-HANDOFF = `no successor`.** Research S05 / `plan simulate` / D21+ / hosted MCP stay off-board unless Notes explicitly promote (they do not).

## Environment

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` (`go.mod`) |
| `go version` | go1.24.2 linux/amd64 |
| S01–S03 / two-clone named tests | `CGO_ENABLED=1` |
| Honesty / E / F / ablation | `CGO_ENABLED=0` |
| Gate H / compat / p0x / x0 / product bar | `CGO_ENABLED=1` |
| Product bar env | `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off` |
| Gate C metrics | `docs/verification/gate-c-x0/` (`dry_run:false`, N=3/condition) |

## Evidence table (independent)

| Bucket / command | Result |
|------------------|--------|
| S01 DF-80/84/85 + P16 keepers `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip\|TestSeedExportOmitsDeniedSurfaces\|TestSeedExportWritesExportedAtCommit\|TestSeedImportAndWhy\|TestSeedImportDiscoveryMentionsTask\|TestSeedImportImpactFindings\|TestSeedImportFromIDAliases\|TestSeedImportRelativePathAgainstC\|TestSeedImportMissingEndpointsMessage'` | **PASS** — 0.449s |
| S02 DF-82/85 `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpSeedExportPath\|TestHelpHandoffSoT\|TestAsOperatorFlagIdentityDocs\|TestSeedExport'` | **PASS** — 0.226s |
| S03 DF-81/83/84 `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportIdempotent\|TestSeedImportDuplicateLinksNoOp\|TestSeedImportSameIdLastWins\|TestSeedImportPlanTreeIdempotent\|TestSeedExportRoundTrip\|TestHelpSeedExportPath'` | **PASS** — 0.187s |
| Two-clone recipe `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run TestPortableGraphTwoCloneWhyContextPlan` | **PASS** — 0.113s (implemented VERIFY remediation in `cmd/trace/cli_test.go`) |
| Honesty full `CGO_ENABLED=0 go test ./evals/honesty/... -count=1` | **PASS** |
| Honesty A/B/C + G `CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim\|TestHonestyEscapeRateGateGPrelim'` | **PASS** |
| Gate E `CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan` | **PASS** |
| Gate F `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` | **PASS** |
| Ablation `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` | **PASS** |
| Gate H `CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH` | **PASS** (~5.4s) |
| Compat `CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist` | **PASS** |
| P0-X + X0 `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1` | **PASS** |
| Product `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1` | **PASS** — all product pkgs ok (cmd/trace 2.568s) |
| `.gitignore` inspect | **PASS** — `.trace/` only for Trace store; `trace/graph.json` not ignored |
| DF-86 grep `install git-hook` / `git-hook` in `*.go` | **PASS (absent)** — zero matches; non-fail per lock |
| Gate C artifacts inspect | **PASS** — `dry_run:false`; N=3 (b0-gatec-1..3, g1-gatec-1..3); **not** re-scored |
| G19 library packages do not import `cmd/trace` / `cmd/trace-mcp` | **PASS** — only test-string mentions in `internal/mcp/mcp_test.go` and `evals/compat/compat_test.go` |
| P16 history intact | **PASS** — Phase 16 DR-HANDOFF not rewritten |

## Two-clone recipe

**Named test (primary fail bar):** `TestPortableGraphTwoCloneWhyContextPlan` in `cmd/trace/cli_test.go`.

| Field | Value |
|-------|-------|
| Isolation | Two `t.TempDir()` clones; source dir has its own `.trace/`; clones never share source DB path |
| Workflow per clone | `init` → `seed import trace/graph.json` → `index sample.js` → `plan show --goal` → `why decision` → `context` |
| Seeded ids | goal `11111111-1111-1111-1111-111111111111`; task `22222222-2222-2222-2222-222222222222`; decision `33333333-3333-3333-3333-333333333333` |
| Offline | No HTTP, no MCP server, no account |
| VERIFY remediation | Test was **absent at preflight** — implemented on this row (allowed) |

**Shell corroboration (secondary):** per `00-PLANNER.md` locked recipe — rsync repo minus `.trace`/`.git`, then init/import/index/plan/why/context in two temp dirs. Named test is authoritative evidence this VERIFY run.

## Law checks

| Check | Hold? |
|-------|-------|
| No daemon / always-on HTTP as primary surface | Yes |
| No committed `.trace/` under `fixtures/` or `evals/` (temp dirs / gitignore only) | Yes |
| G19 — library packages do not import `cmd/trace` or `cmd/trace-mcp` | Yes |
| S01–S03 evidence is **named tests** — not Notes-only | Yes |
| **No** new MCP seed tool; local stdio MCP unchanged | Yes |
| Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat still green | Yes |
| Gate C evidence remains `dry_run:false` — not Phase 01 dry-run | Yes |
| `.gitignore` lists `.trace/` only — `trace/graph.json` **not** ignored | Yes |
| **DF-86 hook absence does not fail VERIFY** | Yes — absent, grep recorded |
| **No research S05 / plan simulate / D21+ scaffold** | Yes — not promoted |
| Phase 16 historical `no successor` left intact | Yes |
| Forward-only: do **not** rewrite Phase 00–16 `done` history | Yes |

## Residuals / deferrals (non-blocking)

| Residual | Disposition | VERIFY note |
|----------|-------------|-------------|
| Encryption-as-git | wontfix | Note only — **not** fail criteria |
| Reviews omitted from default export | out | Note only — **not** fail criteria |
| DF-86 git-hook | deferred | Absent — grep zero matches — **not** fail criteria |
| CGO=0 `cmd/trace` | carry-forward | CGO=1 authoritative — **not** fail criteria |
| S03 `work_state` preservation | SQL-only gap | No dedicated named re-import test — **not** fail criteria |
| Goals #2–#4 | deferred | Research S05 / `plan simulate` / D21+ stay **off-board** |
| Hosted MCP / HTTP / OAuth | out | TODO Later developments — **not** successor |

## Dry-run ≠ gates

**Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist.** Gate C artifacts remain Mode-B `dry_run:false` (inspect only; not re-scored). Phase 01 dry-run is regression-only.

## Handoff

| Item | Value |
|------|-------|
| **DR-HANDOFF** | **`no successor`** (**started** P17-S04-01; close owned by **P17-S04-02**) |
| Research S05 / plan simulate / D21+ / hosted MCP | **Do not scaffold** — no promotion |
| Completion owner | **P17-S04-02** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` |
| Spawns | **none** |
| Next board | **P17-S04-02** |
