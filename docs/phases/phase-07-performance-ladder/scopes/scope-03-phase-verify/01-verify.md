# P07 / S03 / 01 — Phase 07 VERIFY

## Metadata
- id: P07-S03-01
- todo_ids: [P07-S03-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently close Phase 07 — **Gate H** planted performance ladder + S01 T0/isolation + S02 Go adapter + carry-forward honesty/Gates/ablation/p0x/x0/Gate C — against live packages.

**Harness ownership (FINAL):** S01/S02 left **no** `evals/perf`. This VERIFY row **creates** the planted Gate H harness under `evals/perf` (like Gate F/G planted packages), then measures and locks thresholds. Do **not** mark Gate H blocked-until-harness; do **not** invent commercial theater numbers before measuring.

Do **not** trust S01/S02 Notes alone. Do **not** reopen Gate C, invent commercial multi-model perf theater, claim 1M-LOC CI theater, or declare commercial A1 / product thesis.

Write durable evidence, then either:

1. **Pass** → declare **Phase 07 VERIFY PASS / Gate H green** + **start DR-HANDOFF** Phase 08 (`phase-08-ecosystem-hardening` + board `P08-00`), or
2. **Fail** → **spawn forward-only remediations** (01a/01b/+01c).

No product features outside `evals/perf` harness creation (+ spawn remediations if needed).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF, DR-AGENT
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 7 Gate H; Phase 8 Ecosystem & hardening
- [docs/STORAGE_AND_PERFORMANCE.md](../../../../STORAGE_AND_PERFORMANCE.md) §10 metrics
- [docs/init/E_ASSUMPTION_REGISTER.md](../../../../init/E_ASSUMPTION_REGISTER.md) — A5 (SQLite ceiling)
- Sibling: [S01 REVIEW-NOTES.md](../scope-01-incremental-indexing/REVIEW-NOTES.md), [S02 REVIEW-NOTES.md](../scope-02-language-plugins/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 06 VERIFY [`../../../phase-06-environment-capability/scopes/scope-03-phase-verify/01-verify.md`](../../../phase-06-environment-capability/scopes/scope-03-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults (FINAL — P07-S03-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate H (definition) | Practical indexing on planted size ladders (`A_PROJECT_PLAN` Gate H / Phase 7) — **not** commercial multi-model theater |
| Gate H path / harness | **`evals/perf`** package; named test **`TestPlantedPerfLadderGateH`** |
| Schema / metrics | Committed **`schema-gate-h.json`** **v1**; temp **`metrics-gate-h.json`** under `t.TempDir()` (schema-validated) |
| Harness ownership | **S03-01 creates** `evals/perf` as VERIFY work (S01/S02 seeded none). Pattern: Gate F/G planted packages. **Do not** leave Gate H blocked awaiting a prior seed. |
| Allowed Go on this row | **`evals/perf/**` only** (doc.go, schema, planted fixtures, named test). No opportunistic `cmd/` / `internal/` product rewrites unless spawn remediation requires them. |
| Size ladder (synthetic planted) | **`smoke`** (~50–200 LOC) + **`rung-1k`** (~1k LOC) + **`rung-10k`** (~10k LOC). Toward A_PROJECT_PLAN 10k–1M tables; **100k / 1M CI plants deferred** (residual — not Phase 07 Gate H pass bar). |
| Optional Go fixtures | Tiny `.go` sources OK on smoke (and optionally higher rungs) to exercise S02 adapter — **not** required on every rung |
| Thresholds | **Measure-then-threshold** (below) — **forbidden** to invent commercial theater ms/LOC pass numbers before first measurement |
| Gate H pass bar | Named test PASS + schema-valid temp metrics + structural rung asserts (index counts / T0 skip / isolation) + **locked regression ceilings derived from first measurements** (documented in VERIFY-NOTES) |
| S01 surface (must stay green) | T0 always-skip (`isT0SkipDir`/`isT0SkipPath`); walk order T0→lang→T0 file→gitignore; `TestWalkIndexableT0AlwaysSkip` + `TestIndexSkipsExplicitT0Path` + `TestIndexIncrementalIsolation`; **no** `011_*` |
| S02 surface (must stay green) | `LangGo` + `.go`; `extract_go.go`; `tree-sitter-go` **v0.25.0**; `TestIndexFileGoGolden` + DetectLanguage `.go`; CGO analyzers-only; **no** plugin registry; **no** `011_*` |
| Honesty Paths A/B/C | Fail-closed — `TestHonestyFailClosedPlantedClaim` |
| Gate G | **Green** — `TestHonestyEscapeRateGateGPrelim` (escapes=1/caught=2/attempts=3) |
| Gate E | **Green** — `TestPlantedDiscoveryReplan` |
| Gate F | **Green** — `TestPlantedImpactConflictsGateFPrelim` (TP=3/FN=0/FP=0/TN=1; P/R=1.0) |
| Ablation | **Green** — `TestPlantedCapabilitySelectionAblation` (TP=3/FN=0/FP=0/TN=1; P/R=1.0) |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; means G1 0.800 > B0 0.000; **do not invent new Go**; Mode-B packs historical |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / **Gate H** — Phase 01 dry-run is regression-only |
| Fixture hash pin (carry) | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Deferrals (carry) | GC-03/04 still deferred; **`plan simulate`** still out; **100k/1M planted CI ladders** deferred |
| Residuals (non-blocking) | DPC-global; non-tx Apply; UNIQUE re-link; MCP no severity; S01 min argv nit; S02 blank/dot import nit; A5 SQLite ceiling still ACCEPTED_RISK until larger plants |
| VerifiedFact | Still **out** |
| Product Go (non-harness) | **Forbidden** — no indexer rewrite / graph DB / language mega-registry |
| MCP / daemon / HTTP / embeddings | Still forbidden as primary; Gate H does not require MCP |
| Phase 08 folder name | **`phase-08-ecosystem-hardening`** (A_PROJECT_PLAN Phase 8) |

### Measure-then-threshold protocol (FINAL)

```text
1. Create evals/perf (doc.go, schema-gate-h.json v1, planted synthetic ladders smoke/1k/10k, TestPlantedPerfLadderGateH).
2. First measurement pass: run ladder indexing; record per-rung file_count, approx_loc, initial_index_ms, incremental_index_ms (or re-index unchanged), db_bytes into VERIFY-NOTES. Do not invent ceilings beforehand.
3. Derive regression ceilings from measured values with generous headroom, e.g.:
     ceiling_ms  = max(measured_ms * 5, 2000)   # per rung metric
     ceiling_db  = measured_db_bytes * 3
   Encode ceilings in the harness (constants or metrics fields the test asserts). Document derivation in VERIFY-NOTES.
4. Re-run TestPlantedPerfLadderGateH — must PASS with ceilings locked + structural asserts.
5. Forbidden: commercial-repo theater numbers; claiming Gate H from S01 t.Logf alone; claiming 1M LOC CI pass; using Phase 01 dry-run as Gate H.
```

### Metrics schema shape (v1 — lock fields)

`schema-gate-h.json` **required** keys (extend with `additionalProperties: true` OK):

| Field | Notes |
|-------|-------|
| `schema_version` | const `1` |
| `gate` | const `"H"` |
| `suite` | const `"perf"` |
| `dry_run` | boolean — must be **`false`** for Gate H pass artifact |
| `named_test` | `"TestPlantedPerfLadderGateH"` |
| `rungs` | array of rung objects (id, approx_loc, file_count, initial_index_ms, incremental_index_ms, db_bytes, …) |
| `thresholds` | object documenting derived ceilings (ms/db) used by asserts |
| `t0_skip_ok` | boolean |
| `incremental_isolation_ok` | boolean |
| `go_adapter_exercised` | boolean (true if any `.go` plant indexed) |
| `s01_hooks` / `s02_hooks` | string arrays naming consumed surfaces |

### Locked verify commands

```bash
# --- Gate H: create harness if absent, then primary named test (CGO — analyzers) ---
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH

# Perf package
CGO_ENABLED=1 go test ./evals/perf/... -count=1

# S01 T0 + isolation (cmd/trace)
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestWalkIndexableT0AlwaysSkip|TestIndexIncrementalIsolation|TestIndexSkipsExplicitT0Path'

# S02 Go adapter
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1 -run 'TestIndexFileGoGolden|TestDetectLanguage'

# Honesty: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E / F / capability ablation carry-forward
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# Supporting surfaces
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... -count=1

# Full regression bar (includes perf + prior evals)
CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./evals/perf/... ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -v -run TestPlantedPerfLadderGateH
test -f evals/perf/schema-gate-h.json
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3, means match GATE-C-NOTES — do not re-score
find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# expect: 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/` (temp dirs only)
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19) — evals may drive CLI/analyzers as tests do
- [ ] Gate H evidence is **`evals/perf` `TestPlantedPerfLadderGateH`** + schema/metrics — not vibes / Notes-only / S01 `t.Logf` alone
- [ ] Thresholds derived from measurements (VERIFY-NOTES documents derivation)
- [ ] Ladder is **synthetic planted** (smoke/1k/10k) — not commercial multi-model theater; 100k/1M deferred
- [ ] S01 T0 + isolation still green; S02 Go golden still green
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone
- [ ] Mode-B packs not falsified
- [ ] Embeddings / VerifiedFact / `plan simulate` still out
- [ ] No full-rebuild-on-any-change indexer architecture
- [ ] GC-03/04 remain deferred unless explicitly promoted in VERIFY-NOTES

### DR-HANDOFF duties (this row + S03-02)

Per protocol Phase handoff + **DR-HANDOFF**. On Gate H + regression green → scaffold Phase 08. Do **not** leave Phase 08 as README-only / blocked-until-noticed.

| Who | Duty |
|-----|------|
| **P07-S03-01 (this VERIFY)** | On pass: **Create** Phase 08 scaffold under `docs/phases/phase-08-ecosystem-hardening/` (checklist below) **and** append Phase 08 board section to `docs/TODO.md` with `P08-00` as first pending row after Phase 07’s last row. Record handoff progress in `VERIFY-NOTES.md`. |
| **P07-S03-02 (final review)** | **Owns completion check** — refuse `done` until scaffold is runnable (planner + stubs + board). May finish missing pieces. Marks Phase 07 complete only when handoff + VERIFY evidence agree. |

#### Phase 08 scaffold checklist (minimum)

Create / ensure:

```text
docs/phases/phase-08-ecosystem-hardening/
  README.md                 # goal = Ecosystem & hardening (A_PROJECT_PLAN Phase 8)
  00-PHASE-PLANNER.md       # runnable (Agent→clarify→Plan→execute); light locks OK
  scopes/                   # ≥1 scope stub folder recommended; minimal OK
    scope-01-…/
      00-PLANNER.md
      01-*.md
      02-scope-review.md
      SCOPE-TODOS.md
```

Board (`docs/TODO.md`): new **Phase 08** section after Phase 07 rows; first pending ID **`P08-00`** → `phase-08-ecosystem-hardening/00-PHASE-PLANNER.md`. Do **not** start executing Phase 08 until Phase 07 S03-02 is `done`.

Deep tasking stays with Phase 08’s own phase/scope planners. This row delivers a **runnable handoff**, not finished implement prompts.

**Counterfactual:** If Gate H / primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** scaffold Phase 08 until VERIFY finally greens (or user records `no successor`).

## Board rights
Verify: **status + notes** on `P07-S03-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails; **may create** `evals/perf` harness + upcoming Phase 08 scaffold (DR-HANDOFF). Do **not** rewrite Phase 07 `done` history. Do **not** mark `P07-S03-02` or Phase 08 rows `done`.

## Preflight / Plan
1. Re-read this prompt + board row + S01/S02 REVIEW-NOTES + Gate H locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`; packages `evals/{honesty,x0,p0x,replan,impact,capability}`; **no** prior `evals/perf` expected (create it).
3. Plan: create Gate H harness → measure → lock thresholds → run locked commands → fill evidence → pass→VERIFY-NOTES + Phase 08 scaffold **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. Gate H harness create + re-check (required)

Confirm all of:

| Check | Expect |
|-------|--------|
| Harness path | `evals/perf/` created; `doc.go` names Gate H planted ladder |
| Named test | `TestPlantedPerfLadderGateH` PASS under locked command |
| Schema / metrics | Committed `schema-gate-h.json` v1; test writes schema-valid temp `metrics-gate-h.json` (`dry_run:false`) |
| Rungs | smoke + ~1k + ~10k synthetic plants present and measured |
| Thresholds | Derived from first measurements; ceilings encoded; derivation in VERIFY-NOTES |
| Structural | T0 skip OK under plant; incremental isolation OK; optional `.go` exercised if planted |
| Honesty of claim | Notes map Gate H to this harness — **not** S01 `t.Logf`, Gate C scores, or commercial theater |

Record pass/fail in VERIFY-NOTES evidence table.

### B. S01 / S02 surfaces (required)

| Check | Expect |
|-------|--------|
| S01 T0 + isolation | `TestWalkIndexableT0AlwaysSkip`, `TestIndexSkipsExplicitT0Path`, `TestIndexIncrementalIsolation` PASS |
| S02 Go adapter | `TestIndexFileGoGolden` (+ DetectLanguage `.go`) PASS; `tree-sitter-go` v0.25.0 still required |

### C. Carry-forward bars (required)

| Check | Expect |
|-------|--------|
| Honesty H5 Paths A/B/C | `TestHonestyFailClosedPlantedClaim` PASS |
| Gate G | `TestHonestyEscapeRateGateGPrelim` PASS |
| Gate E | `TestPlantedDiscoveryReplan` PASS |
| Gate F | `TestPlantedImpactConflictsGateFPrelim` PASS |
| Ablation | `TestPlantedCapabilitySelectionAblation` PASS |
| P0-X 7/7 | `evals/p0x` PASS |
| X0 packages | `./evals/x0/...` PASS |
| Gate C artifacts | `GATE-C-NOTES.md` still **Go**; metrics `dry_run:false`, N=3; means G1 0.800 > B0 0.000 — **re-check files, do not invent new Go** |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H | Explicit in VERIFY-NOTES |
| Full suite | Locked full command PASS |

### D. Independent re-run (required)

From module root `/home/ali/Desktop/Trace`, run the **Locked verify commands**. All required commands must exit 0 for a green gate. Capture command + PASS/FAIL in `VERIFY-NOTES.md`.

### E. Evidence table (required in `VERIFY-NOTES.md`)

Copy and fill:

| Gate | Result (pass/fail) | Evidence (test / log gist) |
|------|--------------------|----------------------------|
| Gate H harness created (`evals/perf`) | | path + files |
| `TestPlantedPerfLadderGateH` | | named test PASS |
| `schema-gate-h.json` v1 + temp `metrics-gate-h.json` | | file present; `dry_run:false`; validation |
| Measure-then-threshold derivation | | measured → ceilings (VERIFY-NOTES prose) |
| Rungs smoke / ~1k / ~10k | | counts + ms + db_bytes |
| T0 skip under plant | | `t0_skip_ok` |
| Incremental isolation under plant | | `incremental_isolation_ok` |
| Optional Go fixtures / `go_adapter_exercised` | | if planted |
| S01 `TestWalkIndexableT0AlwaysSkip` | | cmd/trace |
| S01 `TestIndexIncrementalIsolation` | | cmd/trace |
| S01 `TestIndexSkipsExplicitT0Path` | | cmd/trace |
| S02 `TestIndexFileGoGolden` / DetectLanguage | | analyzers + tree-sitter-go v0.25.0 |
| Honesty H5 Paths A/B/C | | `TestHonestyFailClosedPlantedClaim` |
| Gate G prelim | | `TestHonestyEscapeRateGateGPrelim` |
| Gate E mini-eval | | `TestPlantedDiscoveryReplan` |
| Gate F prelim | | `TestPlantedImpactConflictsGateFPrelim` |
| Capability ablation | | `TestPlantedCapabilitySelectionAblation` |
| P0-X 7/7 | | `TestP0XAllCriteria` (or package PASS + note) |
| X0 packages | | `./evals/x0/...` PASS |
| Gate C `dry_run:false` intact | | `docs/verification/gate-c-x0/`; GATE-C-NOTES Go |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H | | explicit |
| `go test ./...` (+ perf) | | PASS/FAIL |
| Law checks | | no daemon/HTTP primary; no committed `.trace/`; G19; no full-rebuild; no commercial perf theater; 100k/1M deferred |
| Residuals (non-blocking) | | DPC-global; GC-03/04; A5; deferred ladders; S01/S02 lows |
| DR-HANDOFF | | Phase 08 path created? board `P08-00` present? |

Also record: date, go/CGO note, threshold derivation, residuals carried forward.

### F. On **all gates pass** + law checks

1. Write `VERIFY-NOTES.md` with verdict **Phase 07 VERIFY PASS / Gate H green**, confidence, evidence table, residuals as **non-blocking** if primary paths green. Explicitly: Gate H = `evals/perf` named planted ladder; thresholds measured-then-locked; Gate C artifacts intact; Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H.
2. **Start DR-HANDOFF:** create `docs/phases/phase-08-ecosystem-hardening/` per checklist + append Phase 08 section / `P08-00` to `docs/TODO.md`.
3. Board Notes: short “Gate H + S01/S02 + honesty/Gates/ablation/p0x/x0 PASS; Gate C intact; Phase 08 scaffold started; see VERIFY-NOTES.md; pending P07-S03-02 handoff close”.
4. Mark `P07-S03-01` **done**. Do **not** mark S03-02 done.

### G. On **any fail**

1. `VERIFY-NOTES.md`: verdict **FAIL**, which gate/law, minimal reproduction.
2. Insert `P07-S03-01a` (implement) + `P07-S03-01b` (review) immediately below this row; set this row **`blocked`** (or `failed` + plan `01c` re-VERIFY) with reason. Prefer new `01c` verify row if this row already closed as failed/blocked (forward-only history).
3. Spawn prompts must be **full** protocol skeletons. Scope remediations to the failing layer — **do not** weaken bars, invent Gate H pass from Notes, or rewrite Gate C Mode-B packs.
4. **Forbidden “fixes”:** claiming Gate H without `TestPlantedPerfLadderGateH`; inventing theater thresholds without measurement; claiming Gate C from dry-run; requiring MCP for Gate H; adding daemon/HTTP; weakening honesty Paths A/B/C; rewriting `done` S01–S02 prompts; promoting GC-03/04 without evidence; scaffolding Phase 08 on a red VERIFY; inventing VerifiedFact; shipping commercial multi-model perf theater / full-rebuild indexer.

### Spawn ID convention

```text
… P07-S03-01  (this VERIFY)
… P07-S03-01a (remediation implement)   ← insert immediately below
… P07-S03-01b (remediation review)
… P07-S03-01c (re-VERIFY)               ← if original VERIFY closed as failed/blocked
… P07-S03-02  (phase review — only after VERIFY finally done; owns DR-HANDOFF completion)
```

Update `SCOPE-TODOS.md` when spawning.

## Out of scope
- Implementing Phase 08 ecosystem/hardening product features
- Re-running live multi-model Gate C / flipping Go without contradicting evidence
- Rewriting Mode-B packs
- Promoting GC-03/04 without measurement need
- Declaring A1 / product thesis commercially validated
- Starting Phase 08 implement wave before S03-02 closes
- Expanding Gate H into commercial multi-model / 1M-LOC CI theater (100k/1M deferred)
- Introducing VerifiedFact promotion engine / `plan simulate` / graph DB
- Inventing Gate H thresholds before first measurements

## Todo updates
Status + Notes on `P07-S03-01`; spawn rows if needed; Phase 08 upcoming artifacts on pass; `SCOPE-TODOS.md` checkboxes.

## Exit criteria
- [ ] `evals/perf` Gate H harness created + `TestPlantedPerfLadderGateH` + schema/metrics + measure-then-threshold documented
- [ ] S01 T0/isolation + S02 Go golden re-proved
- [ ] Honesty A/B/C + Gate G + Gate E + Gate F + ablation + p0x + x0 + Gate C artifact integrity + `./...` recorded in `VERIFY-NOTES.md`
- [ ] Evidence table includes Gate H path + dry-run≠… + residuals + handoff
- [ ] Law checks recorded; Mode-B packs not falsified; no commercial perf theater
- [ ] Either gates green **and** Phase 08 scaffold started (folder + `P08-00` board) **or** remediations spawned with full prompts + this row blocked/failed honestly
- [ ] TODO.md status + Notes updated; SCOPE-TODOS synced
- [ ] No non-harness product feature Go on this row

## Minimal todos
- [ ] Preflight: S01/S02 REVIEW-NOTES + Gate C metrics + confirm `evals/perf` absent → create
- [ ] Create planted Gate H harness (smoke/1k/10k) + schema-gate-h.json v1
- [ ] Measure → lock thresholds → re-run `TestPlantedPerfLadderGateH`
- [ ] Re-prove S01 T0/isolation + S02 Go + honesty/Gates/ablation/p0x/x0/Gate C/`./...`
- [ ] Write `VERIFY-NOTES.md` evidence table (incl. threshold derivation)
- [ ] Pass → scaffold `phase-08-ecosystem-hardening` + board `P08-00`; Fail → spawn 01a/01b (+01c plan)
- [ ] Board status + SCOPE-TODOS
