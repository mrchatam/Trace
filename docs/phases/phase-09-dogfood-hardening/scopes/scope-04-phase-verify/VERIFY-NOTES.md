# P09-S04-01 — Phase VERIFY notes (dogfood hardening closeout)

**Date:** 2026-08-16  
**Verifier:** independent re-run (does **not** trust S01/S02/S03 Notes alone)  
**Verdict:** **Phase 09 VERIFY PASS / dogfood hardening green**  
**Confidence:** high  
**Spawns:** none  

**Explicit claims:** DF-01 regression green via live **`TestWhyAndContextWithLinkedReview`** (`lookupEntity` `case "review"`). S02 discoverability green via **`TestTasksListAfterSeed`** + **`TestSeedImportRelativePathAgainstC`** (+ optional **`TestListTasks`**). S03 install-wire green via **`TestInstallCursor*`** (PrintSnippet / PrintBin / WriteMergeBackup / WriteInvalidJSON) + DF-05 docs (README + `experiments/ab-simple/PROTOCOL.md` run-folder / `${workspaceFolder}` footgun). MCP still **six** tools (`trace_why`/`trace_context`/`trace_add`/`trace_link`/`trace_transition`/`trace_review`) — **no** list-tasks. Honesty Paths A/B/C + Gate G + Gate E + Gate F + capability ablation + Gate H + compat checklist + p0x + x0 + domain/store/planner/compiler/mcp + full `./...` PASS. Gate C artifacts remain **Go** (`dry_run:false`, N=3; mean G1 **0.800** > B0 **0.000**).  

**Explicit non-claims:** Phase 01 dry-run is **not** Gate C, **not** Gate F, **not** Gate G, **not** ablation, **not** Gate H, and **not** the compat checklist. Mode-B packs remain historical. GC-03/04 stay deferred. **`plan simulate`** still out. **100k / 1M** planted CI ladders deferred. VerifiedFact still out. No product Go on this row. Phase 09 not marked complete here — **P09-S04-02** owns handoff close + phase complete.

**DR-HANDOFF = `no successor`.** Ladder gaps (D08/D09/combos/multi-agent; DF-11/12 tighten) stay on the parallel `experiments/` dogfood track. Do **not** scaffold Phase 10 unless Notes explicitly promote.

## Environment

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| `go version` | go1.24.2 linux/amd64 |
| S02/S03 CLI / Gate H / compat / p0x / x0 / full suite | `CGO_ENABLED=1` |
| DF-01 / store ListTasks / Honesty / Gate E / Gate F / ablation / domain/store/planner/compiler/mcp | `CGO_ENABLED=0` where locked |
| Fixture content hash (live) | `2d1ac2a7f142fb715a6b138be5064fc1877674105f273db89e0e5782851d2e3a` (drift from S02 `fixtures/x0/README.md` DF-04 note — see residuals) |
| Gate C recorded pin / metrics `git_sha` | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` (historical; artifacts intact) |
| Gate C metrics | `docs/verification/gate-c-x0/` (`dry_run:false`, N=3/condition) |

## Evidence table (independent)

| Command / check | Result |
|-----------------|--------|
| `CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run TestWhyAndContextWithLinkedReview` | **PASS** — `ok retrieval` EXIT:0 (~0.02s) |
| `CGO_ENABLED=0 go test ./internal/store/... -count=1 -run TestListTasks` | **PASS** — `ok store` EXIT:0 |
| `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestTasksListAfterSeed\|TestSeedImportRelativePathAgainstC'` | **PASS** — `ok cmd/trace` EXIT:0 |
| `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallCursor'` | **PASS** — `ok cmd/trace` EXIT:0 (PrintSnippet + PrintBin + WriteMergeBackup + WriteCreateMissing + WriteInvalidJSON) |
| `CGO_ENABLED=0 go test ./evals/honesty/... -count=1` | **PASS** |
| `CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim\|TestHonestyEscapeRateGateGPrelim'` | **PASS** (A/B/C + Gate G) |
| `CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan` | **PASS** (Gate E) |
| `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` | **PASS** (Gate F) |
| `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` | **PASS** |
| `CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH` | **PASS** (~5.1s) |
| `CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist` | **PASS** |
| `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1` | **PASS** — p0x + x0 |
| `CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... ./internal/mcp/... -count=1` | **PASS** — all five |
| `CGO_ENABLED=1 go test ./... -count=1` | **PASS** — full regression bar |
| Install snippet spot-check (`trace install cursor`) | **PASS** — `mcpServers.trace` stdio + `-C` `${workspaceFolder}` |
| DF-05 docs (README + `experiments/ab-simple/PROTOCOL.md`) | **PASS** — run-folder / workspaceFolder footgun documented |
| MCP six tools / no list-tasks | **PASS** — `internal/mcp/server.go` + `BuiltinMCPCapabilitySpecs` = six; no `trace_tasks` in `internal/mcp` |
| Gate C artifacts inspect | **PASS** — `dry_run:false`; N=3; mean G1 0.800 > B0 0.000; **not** re-scored |
| Fixture content hash vs historical pin | **NOTE** — live `2d1ac2a7…` ≠ historical `15fe50a1…` (S02 README); Gate C Go metrics unchanged |

## Law checks

| Check | Hold? |
|-------|-------|
| No daemon / always-on HTTP as primary surface | Yes |
| No committed `.trace/` under `fixtures/` or `evals/` | Yes |
| G19 — library packages do not import `cmd/trace` or `cmd/trace-mcp` | Yes (`go list` imports clean; mcp_test boundary strings only) |
| DF-01 evidence is live named test (not Notes-only) | Yes — `TestWhyAndContextWithLinkedReview` + `case "review"` |
| S02/S03 proven by named tests + DF-05 docs; no MCP list-tasks | Yes |
| Honesty A/B/C + Gate G + E + F + ablation + H + compat green | Yes |
| Gate C evidence remains `dry_run:false` — not Phase 01 dry-run | Yes |
| Mode-B packs not falsified | Yes (historical) |
| Embeddings / VerifiedFact / `plan simulate` still out | Yes |
| No full-rebuild-on-any-change indexer architecture | Yes |
| No new migration `011_*` from Phase 09 | Yes (schema through `010_capability_surface.sql`) |
| No Phase 10 scaffold | Yes |

## Residuals / deferrals

- **Fixture content hash drift:** live `2d1ac2a7…` after S02 updated `fixtures/x0/README.md` (DF-04 harness note). Gate C recorded pin / metrics `git_sha` remain `15fe50a1…`. Non-blocking; do not re-score Gate C.
- S01 residuals (carry): `plan_scope` ExactLookup out; scope-only review expand untested.
- S03 residual (carry): degenerate non-object `mcpServers`.
- GC-03/04 deferred; `plan simulate` out; 100k/1M planted CI ladders deferred.
- Parallel dogfood ladder gaps (D08/D09/combos/multi-agent; DF-11/12) — **not** board-blocking; stay on `experiments/`.

## Handoff

| Item | Value |
|------|-------|
| **DR-HANDOFF** | **`no successor`** |
| Phase 10 | **Do not scaffold** |
| Completion owner | **P09-S04-02** — closed 2026-08-16 (fresh suite re-check + APPROVE high; see [REVIEW-NOTES.md](REVIEW-NOTES.md)) |
| Next board row | **none** (roadmap closed; parallel dogfood off-board) |
