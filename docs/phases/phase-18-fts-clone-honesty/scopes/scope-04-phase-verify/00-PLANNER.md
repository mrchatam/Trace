# P18-S04-00 — Phase VERIFY / fts-clone-honesty (FINAL)

## Metadata
- id: P18-S04-00
- todo_ids: [P18-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock Phase 18 VERIFY: **re-prove named DF-87 / DF-88 / DF-89 tests** + keepers + carry-forward. Default **DR-HANDOFF = `no successor`**. DF-86 / DF-67 / harness rsync / stdio EOF / stale binaries are **non-fail**. **Two-clone recipe is not required.** DF-88 re-prove is **document-only** (help + omit keepers — no clone dirs). **No product Go.** **S05 rebuild still follows this scope** — VERIFY does **not** mark Phase 18 complete and does **not** close DR-HANDOFF.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- [../../DF-88-DECISION.md](../../DF-88-DECISION.md)
- Pattern (shape only — **not** two-clone fail bar): [../../../phase-17-portable-graph-git/scopes/scope-04-phase-verify/00-PLANNER.md](../../../phase-17-portable-graph-git/scopes/scope-04-phase-verify/00-PLANNER.md)
- S01 named SoT (landed): [../scope-01-fts-query-sanitize/REVIEW-NOTES.md](../scope-01-fts-query-sanitize/REVIEW-NOTES.md)
- S02 named SoT (landed): [../scope-02-clone-pending-honesty/REVIEW-NOTES.md](../scope-02-clone-pending-honesty/REVIEW-NOTES.md)
- S03 named SoT (landed): [../scope-03-go-method-extract/REVIEW-NOTES.md](../scope-03-go-method-extract/REVIEW-NOTES.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Depends-on: **P18-S01-02 / P18-S02-02 / P18-S03-02 APPROVE**. Import **live** `func Test*` names from S01–S03 REVIEW-NOTES (do not invent). If a REVIEW-NOTES name differs, **live name wins**. Unattended: no Plan-mode switch; no product Go; **do not run S04-01**.

## Depends-on (S01–S03 — landed; live named tests confirmed 2026-08-18)

Grep of `func Test*` in repo matches REVIEW-NOTES **Landed `func Test*` names (S04 import)** tables. Live names win.

| Scope | Board | Named tests imported (live `func Test*` exists) | Keepers |
|-------|-------|--------------------------------------------------|---------|
| S01 DF-87 | **APPROVE high** (P18-S01-02) | `TestSanitizeFTSQueryPunctuationClass` (`internal/store/fts_test.go`); `TestSearchFTSSlashInQuery` (`internal/store/fts_test.go`); `TestTaskContextSlashTitle` (`internal/compiler/compiler_test.go`); `TestTaskContextContinuesWhenSearchErrors` (`internal/compiler/compiler_test.go`) | `TestFTSFindsEntityTitleAndPathSymbol`; `TestIncludeWhyFailClosed`. Omit keeper `TestSeedExportOmitsDeniedSurfaces` lives on the S02 `cmd/trace` line |
| S02 DF-88 | **APPROVE high** (P18-S02-02) | `TestHelpCloneTasksImportPending` (`cmd/trace/cli_test.go`) | `TestSeedExportOmitsDeniedSurfaces`; `TestHelpSeedExportPath`. Does **not** reverse P17 exclude |
| S03 DF-89 | **APPROVE high** (P18-S03-02) | `TestIndexFileGoHandlerMethods` (`internal/analyzers/analyzers_test.go`) — exact `kind:name` `method:Search` + `method:SearchCursor` + `type:Memory` + `type:Notes` on `testdata/handler_methods.go` | `TestIndexFileGoGolden`. `extract_go.go` untouched |

**S05** rebuild still follows this VERIFY scope — do not drop those rows. This VERIFY does **not** complete the phase.

### CGO matrix (FINAL)

| Bucket | Packages | `CGO_ENABLED` | Rule |
|--------|----------|---------------|------|
| DF-87 named + S01 keepers | `./internal/store/...` `./internal/compiler/...` `./internal/retrieval/...` | **0 authoritative** | modernc; retrieval may report `no tests to run` |
| DF-87 same `-run` | same | **1 corroboration** | required second pass (S01-02 ran both) |
| DF-88 named + S02 keepers + P17 seed keepers | `./cmd/trace/...` | **1 required** | tree-sitter |
| DF-89 named + S03 keeper | `./internal/analyzers/...` | **1 required** | tree-sitter |
| Honesty / E / F / ablation | `./evals/honesty\|replan\|impact\|capability/...` | **0** | |
| H / compat / p0x / x0 / product bar | evals + `./cmd\|internal\|evals` | **1** (+ `GOMODCACHE` on product bar) | |
| CGO0 `./cmd/trace/...` / `./internal/analyzers/...` | — | **non-fail** (R4) | do not use as fail bar |

Keep **both** DF-87 slash names: `TestTaskContextSlashTitle` is packet/Layer-0; MATCH hit coverage is `TestSearchFTSSlashInQuery`.

### Two-clone — not required

`TestPortableGraphTwoCloneWhyContextPlan` is P17 history. It may ride the product `./cmd/...` bar. S04 **must not** re-run the P17 two-clone shell recipe, **must not** add a dedicated `-run` for that test, **must not** implement that test if absent, **must not** spawn two-clone remediation, **must not** treat two-clone as DF-88 proof.

### DF-88 — document-only re-prove

Re-prove **docs/help/omit**, not clone dirs: named `TestHelpCloneTasksImportPending` + keepers `TestSeedExportOmitsDeniedSurfaces` / `TestHelpSeedExportPath`. Optional grep: CONTRIBUTING bullet 7; README clone-PENDING sentence; `SeedTask` tags `id/title/body/goal_id` only. Clone PENDING is **expected** ([DF-88-DECISION.md](../../DF-88-DECISION.md)). Do **not** reverse exclude.

## Live residuals → DR-HANDOFF decision (2026-08-18)

| Bucket | Items | Phase implication |
|--------|-------|-------------------|
| Product gaps scheduled in Phase 18 | DF-87/88/89 | Closed by S01–S03 APPROVE — VERIFY must **re-prove named tests** |
| Explicit residual OK into VERIFY | **DF-86** hook absent; **DF-67**; DF-22/37 tip; CGO0 `cmd/trace` / analyzers; harness rsync; MCP stdio EOF; stale `bin/trace` / `bin/trace-mcp` | **Do not fail VERIFY** |
| In-phase after VERIFY | S05 workspace binary rebuild (`bin/trace` CGO1 + `bin/trace-mcp` CGO0) | **Not** a VERIFY fail; **not** research S05; next after S04-02 = **P18-S05-00** |
| Goals sequence #2–#4 | Research S05 / `plan simulate` / D21+ | Stay off-board |
| Hosted MCP / HTTP / OAuth | TODO Later developments | **Not** a VERIFY successor |

**DR-HANDOFF = `no successor` (default).** S04-01 **starts** Notes. S04-02 re-verifies product only — **does not close**. **S05-02 owns close.** Phase 17 historical `no successor` left intact.

## Planner work
1. [x] Import named `-run` filters from S01–S03 REVIEW-NOTES (live grep confirmed)
2. [x] Lock verify command set + CGO matrix + two-clone **not required** + DF-88 document-only
3. [x] Thicken 01-verify / 02-scope-review / SCOPE-TODOS + DR-HANDOFF intent (close = S05-02)
4. [x] Light S05 Depends (upcoming only); leave S05 board rows in place
5. [x] Stamp this prompt **FINAL**; next **P18-S04-01**

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Re-prove named DF-87/88/89 + keepers + carry-forward — **not** a new planted eval; **not** phase complete |
| Named | Exact `-run` from landed REVIEW-NOTES. **2026-08-18:** all imported names exist. Absent at S04-01 preflight → FAIL + spawn; **no product Go** |
| Two-clone | **Not required.** No shell recipe. No dedicated `-run`. Do not implement `TestPortableGraphTwoCloneWhyContextPlan` |
| DF-88 | **Document-only** (help named + omit/path keepers). Keep exclude |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat; p0x; x0; Gate C `dry_run:false`; product `./cmd\|internal\|evals`; P17 `TestSeedExportRoundTrip` + omit + `TestSeedExportWritesExportedAtCommit` (**not** two-clone) |
| Residuals non-fail | DF-86; DF-67; DF-22/37; CGO0 `cmd/trace` / analyzers; harness rsync/stdio EOF; stale binaries (**S05**) |
| DR-HANDOFF | default **`no successor`** — **S04-01 starts** Notes; **S05-02 owns close**. S04-02 product only |
| Next after S04-02 | **P18-S05-00** (rebuild). Not “phase complete” |
| Dry-run ≠ | Gate C / F / G / ablation / H / checklist |
| Allowed Go on VERIFY | **None.** Missing named test → FAIL + spawn |
| Evidence | **`VERIFY-NOTES.md`** in this folder |
| Spawn | On fail: `P18-S04-01a` / `01b` (+`01c`) immediately below |
| Forbidden | Phase 18 complete; closing DR-HANDOFF; two-clone shell or dedicated two-clone `-run` as fail bar; hosted MCP; reversing DF-88 exclude; rewriting P17 history |

### Locked verify command set (FINAL)

**CGO0 `./cmd/trace/...` and CGO0 `./internal/analyzers/...` are R4 — do not use as fail bar.** `GOMODCACHE`+`GOPROXY=off` on full product bar.

```bash
# --- DF-87 S01 (CGO0 store/compiler authoritative) ---
CGO_ENABLED=0 go test ./internal/store/... ./internal/compiler/... ./internal/retrieval/... -count=1 -run 'TestSanitizeFTSQueryPunctuationClass|TestSearchFTSSlashInQuery|TestFTSFindsEntityTitleAndPathSymbol|TestTaskContextSlashTitle|TestTaskContextContinuesWhenSearchErrors|TestIncludeWhyFailClosed'

# --- DF-87 same -run (CGO1 corroboration) ---
CGO_ENABLED=1 go test ./internal/store/... ./internal/compiler/... ./internal/retrieval/... -count=1 -run 'TestSanitizeFTSQueryPunctuationClass|TestSearchFTSSlashInQuery|TestFTSFindsEntityTitleAndPathSymbol|TestTaskContextSlashTitle|TestTaskContextContinuesWhenSearchErrors|TestIncludeWhyFailClosed'

# --- DF-88 S02 document-only (CGO1 cmd/trace) ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpCloneTasksImportPending|TestHelpSeedExportPath|TestSeedExportOmitsDeniedSurfaces'

# --- DF-89 S03 (CGO1 analyzers) ---
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestIndexFileGoHandlerMethods|TestIndexFileGoGolden'

# --- P17 portable-graph keepers (exclude/round-trip — NOT two-clone recipe) ---
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit'

# --- Honesty: Paths A/B/C + Gate G ---
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# --- Gate E / F / capability ablation ---
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# --- Gate H + compat ---
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# --- P0-X + X0 ---
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1

# --- Product regression bar ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Optional (strong evidence, **not** substitutes for named PASS):

```bash
# DF-86: grep — no install git-hook in product Go (non-fail if absent)
# Gate C inspect: docs/verification/gate-c-x0/ dry_run:false N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# DF-88 docs: CONTRIBUTING bullet 7; README clone-PENDING sentence; SeedTask tags id/title/body/goal_id only
# Do NOT fail for R3 graphify space FAIL or R4 CGO0 analyzers / CGO0 cmd/trace FAIL
# Do NOT fail for stale bin/trace or bin/trace-mcp (S05)
# Do NOT run P17 two-clone shell recipe
# retrieval/... may report "no tests to run" — not a fail
```

## Exit criteria
- [x] FINAL command set from live REVIEW-NOTES
- [x] CGO matrix + keepers + two-clone **not required** + DF-88 document-only locked
- [x] 01/02/SCOPE-TODOS thickened; next after review = **P18-S05-00**
- [x] Board Notes; next **P18-S04-01**
- [x] No product Go this row; S04-01 verify **not** started

## Out of scope
- Running VERIFY (S04-01)
- Product Go / new MCP tools / daemon / DF-86 hook
- Two-clone shell recipe / dedicated two-clone `-run` / implementing `TestPortableGraphTwoCloneWhyContextPlan`
- Rebuilding `bin/trace` / `bin/trace-mcp` (S05)
- Closing DR-HANDOFF / marking Phase 18 complete
- Rewriting Phase 00–17 `done` history
- Auto-boarding Phase 19 / research S05 / plan simulate / D21+ / hosted MCP

## Next
**P18-S04-01** (independent VERIFY run — named DF-87/88/89 + carry-forward; does not close the phase).
