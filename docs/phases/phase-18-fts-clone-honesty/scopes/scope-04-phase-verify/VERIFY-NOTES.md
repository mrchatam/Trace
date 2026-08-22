# P18-S04-01 — Phase 18 VERIFY notes (fts-clone-honesty product evidence)

**Date:** 2026-08-18  
**Verifier:** independent re-run (does **not** trust S01–S03 Notes alone)  
**Verdict:** **Phase 18 VERIFY PASS / fts-clone-honesty product green**  
**Confidence:** high  
**Spawns:** none  

**Explicit claims:** S01–S03 named DF regressions (DF-87 / DF-88 / DF-89) green on live packages. DF-87 four named + two S01 keepers PASS on **CGO0 then CGO1** (both slash names kept). DF-88 re-prove is **document-only** (help named + omit/path keepers; no clone dirs; exclude unreversed). DF-89 named + golden keeper PASS on **CGO1**. P17 seed keepers (`TestSeedExportRoundTrip` + omit + `TestSeedExportWritesExportedAtCommit`) PASS — **not** two-clone. Carry-forward honesty Paths A/B/C + Gate G + Gate E + Gate F + capability ablation + Gate H + compat + p0x + x0 + product `./cmd|internal|evals` (CGO1, `GOMODCACHE`+`GOPROXY=off`) PASS. Gate C artifacts remain **Go** (`dry_run:false`, N=3). DF-86 git-hook **absent** (non-fail; grep recorded).

**Explicit non-claims:** Phase 18 is **not** complete. DR-HANDOFF is **started**, **not closed** (S05-02 owns close). Stale `bin/trace` / `bin/trace-mcp` not rebuilt (S05). Two-clone **shell recipe was not run** and is **not** a P18 fail bar. No dedicated `TestPortableGraphTwoCloneWhyContextPlan` `-run`. DF-86 hook not implemented. DF-67 / DF-22/37 not closed. CGO=0 `cmd/trace` / analyzers not used as fail bar. Phase 01 dry-run is **not** Gate C, **not** Gate F, **not** Gate G, **not** ablation, **not** Gate H, and **not** the compat checklist. No research S05 / `plan simulate` / D21+ / hosted MCP scaffold. Phase 17 historical `no successor` left intact.

**DR-HANDOFF started = `no successor`.** Close owned by **P18-S05-02**. Research S05 / `plan simulate` / D21+ / hosted MCP stay off-board.

## Environment

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` (`go.mod`) |
| `go version` | go1.24.2 linux/amd64 |
| DF-87 named + S01 keepers | `CGO_ENABLED=0` authoritative, then `CGO_ENABLED=1` corroboration |
| DF-88 / DF-89 / P17 seed keepers | `CGO_ENABLED=1` |
| Honesty / E / F / ablation | `CGO_ENABLED=0` |
| Gate H / compat / p0x / x0 / product bar | `CGO_ENABLED=1` |
| Product bar env | `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off` |
| Gate C metrics | `docs/verification/gate-c-x0/` (`dry_run:false`, N=3/condition) |
| Preflight | All imported `func Test*` names exist (2026-08-18 live grep) — no product Go |

## Evidence table (independent)

| Bucket / command | Result |
|------------------|--------|
| Preflight imported `func Test*` names (DF-87/88/89 + keepers + P17 seed keepers) | **PASS** — all present; no missing-name spawn |
| S01 DF-87 CGO0 `CGO_ENABLED=0 go test ./internal/store/... ./internal/compiler/... ./internal/retrieval/... -count=1 -run 'TestSanitizeFTSQueryPunctuationClass\|TestSearchFTSSlashInQuery\|TestFTSFindsEntityTitleAndPathSymbol\|TestTaskContextSlashTitle\|TestTaskContextContinuesWhenSearchErrors\|TestIncludeWhyFailClosed'` | **PASS** — store 0.067s; compiler 0.088s; retrieval `[no tests to run]` (OK) |
| S01 DF-87 CGO1 (same `-run`) | **PASS** — store 0.062s; compiler 0.089s; retrieval `[no tests to run]` (OK). Both slash names kept (`TestSearchFTSSlashInQuery` + `TestTaskContextSlashTitle`) |
| S02 DF-88 document-only CGO1 `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpCloneTasksImportPending\|TestHelpSeedExportPath\|TestSeedExportOmitsDeniedSurfaces'` | **PASS** — 0.049s. No clone dirs. Exclude not reversed. `SeedTask` JSON tags `id/title/body/goal_id` only |
| S03 DF-89 CGO1 `CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestIndexFileGoHandlerMethods\|TestIndexFileGoGolden'` | **PASS** — 0.048s |
| P17 seed keepers (exclude/round-trip — **NOT** two-clone) `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip\|TestSeedExportOmitsDeniedSurfaces\|TestSeedExportWritesExportedAtCommit'` | **PASS** — 0.196s. Dedicated two-clone `-run` **not** used. Shell recipe **not** used |
| Honesty full `CGO_ENABLED=0 go test ./evals/honesty/... -count=1` | **PASS** — 0.058s |
| Honesty A/B/C + Gate G `CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim\|TestHonestyEscapeRateGateGPrelim'` | **PASS** — 0.057s |
| Gate E `CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan` | **PASS** — 0.036s |
| Gate F `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` | **PASS** — 0.030s |
| Ablation `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` | **PASS** — 0.033s |
| Gate H `CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH` | **PASS** — 5.455s |
| Compat `CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist` | **PASS** — 0.934s |
| P0-X + X0 `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1` | **PASS** — p0x 1.444s; x0 1.504s |
| Product `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1` | **PASS** — all product pkgs ok (`cmd/trace` 2.654s; `cmd/trace-mcp` `[no test files]`). R3 graphify space FAIL did **not** occur. R4 CGO0 analyzers / CGO0 `cmd/trace` **not** used as fail bar |
| DF-86 grep `install git-hook` / `git-hook` in `*.go` | **PASS (absent)** — zero matches; non-fail per lock |
| Gate C artifacts inspect | **PASS** — `dry_run:false`; N=3 (`b0-gatec-1..3`, `g1-gatec-1..3`); **not** re-scored |
| G19 library packages do not import `cmd/trace` / `cmd/trace-mcp` | **PASS** — only test-string mentions in `internal/mcp/mcp_test.go` and `evals/compat/compat_test.go` |
| DF-88 docs (optional) | **PASS** — CONTRIBUTING bullet 7; README clone-PENDING sentence; `SeedTask` tags `id/title/body/goal_id` only; no `include-reviews` flags |
| `.gitignore` inspect | **PASS** — `.trace/` only for Trace store; `trace/graph.json` not ignored |
| P17 history intact | **PASS** — Phase 17 DR-HANDOFF still **CLOSED** / `no successor` |

## DF-88 document-only

Re-proved **docs/help/omit**, not clone dirs:

| Check | Hold? |
|-------|-------|
| Named `TestHelpCloneTasksImportPending` | Yes — CGO1 PASS |
| Keepers `TestSeedExportOmitsDeniedSurfaces` + `TestHelpSeedExportPath` | Yes — CGO1 PASS |
| P17 exclude reversed? | **No** — omit keeper unweakened; `SeedTask` has no `work_state` JSON tag |
| Include flags (`include-reviews`)? | **None** in product Go |
| Clone dirs / two-clone shell as DF-88 proof? | **Not used** |

Clone PENDING after import remains **expected** (`DF-88-DECISION.md`).

## Two-clone — not required

| Item | This VERIFY |
|------|-------------|
| Dedicated `-run TestPortableGraphTwoCloneWhyContextPlan` | **Not run** — **not** a P18 fail bar |
| P17 two-clone **shell recipe** | **Not run** — **not** a P18 fail bar |
| Product `./cmd/...` bar | May ride that test if present; product bar **PASS** — still **not** used as DF-88 proof |

## Law checks

| Check | Hold? |
|-------|-------|
| No daemon / always-on HTTP as primary surface | Yes — no `ListenAndServe` in `cmd/` / `internal/` |
| No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only) | Yes — zero matches |
| G19 — library packages do not import `cmd/trace` or `cmd/trace-mcp` | Yes |
| S01–S03 evidence is **named tests** — not Notes-only | Yes |
| **No** new MCP seed tool; local stdio MCP unchanged | Yes |
| Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat still green | Yes |
| Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone | Yes |
| **DF-86 hook absence does not fail VERIFY** | Yes — absent, grep recorded |
| **Stale binaries do not fail VERIFY** (S05) | Yes — not rebuilt this row |
| **No Phase 19 / research S05 / plan simulate / D21+ / hosted MCP scaffold** | Yes — not promoted |
| Phase 17 historical `no successor` left intact | Yes |
| Forward-only: do **not** rewrite Phase 00–17 `done` history | Yes |
| DF-88 re-prove is **document-only** (help + omit); exclude unreversed | Yes |
| Two-clone shell **not** used as fail bar; no dedicated two-clone `-run` | Yes |
| Phase 18 **not** marked complete; DR-HANDOFF **not** closed | Yes |
| CGO0 `cmd/trace` / analyzers **not** used as fail bar | Yes |

## Residuals / deferrals (non-blocking)

| Residual | Disposition | VERIFY note |
|----------|-------------|-------------|
| DF-86 git-hook | deferred | Absent — grep zero matches in `*.go` — **not** fail criteria |
| DF-67 symbol-entity staleness | deferred | Note only — **not** fail criteria |
| DF-22 / DF-37 MCP reload | deferred | Tip already shipped — **not** fail criteria |
| CGO=0 `cmd/trace` / analyzers | carry-forward | CGO=1 authoritative for those pkgs — **not** fail criteria |
| Harness rsync / MCP stdio EOF | harness-only | Product tests PASS — **not** fail criteria |
| Stale `bin/trace` / `bin/trace-mcp` | **S05** | **Not** fail criteria; next after S04-02 = **P18-S05-00** |
| Clone PENDING after import | expected (DF-88) | Document-only; exclude kept — **not** fail criteria |
| Goals #2–#4 | deferred | Research S05 / `plan simulate` / D21+ stay **off-board** |
| Hosted MCP / HTTP / OAuth | out | TODO Later developments — **not** successor |

## Dry-run ≠ gates

**Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist.** Gate C artifacts remain Mode-B `dry_run:false` (inspect only; not re-scored). Phase 01 dry-run is regression-only.

## Handoff

| Item | Value |
|------|-------|
| **DR-HANDOFF** | **started = `no successor`** (this row). Successor field left **TBD / default no successor**. **Not closed.** Close owned by **P18-S05-02** |
| Research S05 / plan simulate / D21+ / hosted MCP | **Do not scaffold** — no promotion |
| Phase 18 complete? | **No** — S05 rebuild still pending |
| Binaries rebuilt? | **No** — S05 |
| Spawns | **none** |
| Next board | **P18-S04-02** |
