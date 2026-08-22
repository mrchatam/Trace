# P06 / S03 / 01 — Phase 06 VERIFY

## Metadata
- id: P06-S03-01
- todo_ids: [P06-S03-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently re-prove Phase 06 closeout — capability-selection ablation + S01/S02 capability surfaces + Gate F/G/E/honesty/p0x/x0/Gate C carry-forward bars — against live packages. Do **not** trust S01/S02 Notes alone. Do **not** reopen Gate C, invent commercial multi-model capability theater, claim ontology megastore success, or declare commercial A1 / product thesis.

Write durable evidence, then either:

1. **Pass** → declare **Phase 06 VERIFY PASS / capability-selection ablation green** + **start DR-HANDOFF** Phase 07 (`phase-07-performance-ladder` + board `P07-00`), or
2. **Fail** → **spawn forward-only remediations** (01a/01b/+01c).

No product features.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF, DR-AGENT
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 6 Environment/capability; Phase 7 Performance ladder & language plugins
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md)
- [docs/AGENT_ENVIRONMENT.md](../../../../AGENT_ENVIRONMENT.md)
- Sibling: [S01 REVIEW-NOTES.md](../scope-01-capability-surface/REVIEW-NOTES.md), [S02 REVIEW-NOTES.md](../scope-02-capability-selection/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 05 VERIFY [`../../../phase-05-decision-impact/scopes/scope-03-phase-verify/01-verify.md`](../../../phase-05-decision-impact/scopes/scope-03-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Ablation (definition) | Capability-selection precision/recall on planted probes (`A_PROJECT_PLAN` Phase 6 / H7) |
| Ablation path / harness | **`evals/capability`** package; named test **`TestPlantedCapabilitySelectionAblation`** |
| Ablation pass bar | Planted 4-probe P/R: **TP=3 FN=0 FP=0 TN=1**; **precision=1.0**; **recall=1.0**; schema-valid temp **`metrics-capability.json`** vs committed **`schema-capability.json`** v1 |
| Ablation probes | Pos-1 UNAVAILABLE / Pos-2 UNKNOWN / Pos-3 selection-filter (no catalog dump) / Neg-1 clean AVAILABLE |
| S01 surface (must stay green) | mig **`010_capability_surface.sql`** (`capabilities`+`task_capability_requirements`); kinds SKILL\|RULE\|MCP\|TOOL\|HOOK; domain Upsert/Require/Missing + `BuiltinMCPCapabilitySpecs` (no auto-seed); compiler packet `required_capabilities`+`missing_capabilities` (SchemaVersion `0.1`); thin `trace capability`; **no new** MCP tools / entity_links / `internal/capability` |
| S02 surface | New package `evals/capability` only; S01 hooks only; mig 010 only (no S02 schema fork) |
| Honesty Paths A/B/C | Fail-closed; no `AllowDoneWithoutReview` in A/B/C proof — `TestHonestyFailClosedPlantedClaim` |
| Gate G | **Green** — `evals/honesty` **`TestHonestyEscapeRateGateGPrelim`** (escapes=1/caught=2/attempts=3) |
| Gate E | **Green** — `evals/replan` **`TestPlantedDiscoveryReplan`** (PLAN_AFFECTING+; churn N=5 + ack) |
| Gate F | **Green** — `evals/impact` **`TestPlantedImpactConflictsGateFPrelim`** (TP=3/FN=0/FP=0/TN=1; P/R=1.0) |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** artifacts under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; means G1 0.800 > B0 0.000; **do not invent new Go**; Mode-B packs historical |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation | Phase 01 dry-run is regression-only — **not** Gate C / Gate F / Gate G / capability-ablation pass |
| Fixture hash pin (carry) | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Deferrals (carry) | GC-03/04 still deferred; **`plan simulate`** still out (P13) |
| Residuals (non-blocking) | DPC-global attach; non-tx Apply; UNIQUE re-link; MCP no severity; S02 low residuals per REVIEW-NOTES; ontology megastore still rejected |
| VerifiedFact | Still **out** — not a promotion engine |
| Product Go | **Forbidden** — verify + evidence + handoff scaffold + spawn prompts only |
| MCP | Optional; ablation is capability harness — do not require MCP for VERIFY pass; no new MCP write tools required |
| Daemon / HTTP / embeddings | Still forbidden as primary |
| Phase 07 folder name | **`phase-07-performance-ladder`** (A_PROJECT_PLAN Phase 7 — Performance ladder & language plugins; refine only with Notes) |

### Locked verify commands

```bash
# Capability-selection ablation (primary — named planted P/R harness)
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# Capability package (CGO-free)
CGO_ENABLED=0 go test ./evals/capability/... -count=1

# Honesty: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E carry-forward
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan

# Gate F carry-forward
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim

# S01/S02 supporting surfaces (domain+store+planner+compiler; capability surface / packet)
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... -count=1

# Full regression bar (includes p0x 7/7, x0, replan, honesty, impact, capability, analyzers)
CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -v -run TestPlantedCapabilitySelectionAblation
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -v -run TestHonestyEscapeRateGateGPrelim
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -v -run TestHonestyFailClosedPlantedClaim
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -v -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=1 go test ./evals/p0x/... -count=1 -v -run TestP0XAllCriteria
# Confirm schema file present (metrics are temp — written by ablation test)
test -f evals/capability/schema-capability.json
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3, means match GATE-C-NOTES — do not re-score
# Fixture content hash (must still match pin unless Notes promote change)
find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# expect: 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/`
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] Ablation evidence is **`evals/capability` `TestPlantedCapabilitySelectionAblation`** + schema/metrics — not vibes / Notes-only
- [ ] Planted tallies: TP=3, FN=0, FP=0, TN=1; precision=1.0; recall=1.0
- [ ] S01 hooks remain green: `UpsertCapability` / `RequireCapability` / `MissingCapabilities` + packet required/missing
- [ ] Mig 010 present; no ontology megastore / `internal/capability` stack; no new MCP write tools required for VERIFY
- [ ] Paths A/B/C still fail-closed without hatch in that proof
- [ ] Gate G = `TestHonestyEscapeRateGateGPrelim` still green
- [ ] Gate E = `TestPlantedDiscoveryReplan` still green
- [ ] Gate F = `TestPlantedImpactConflictsGateFPrelim` still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone; no silent Go flip
- [ ] Mode-B packs not falsified
- [ ] Embeddings still absent; VerifiedFact absent as promotion engine
- [ ] No commercial multi-model capability theater claim
- [ ] GC-03/04 remain deferred unless explicitly promoted in VERIFY-NOTES

### DR-HANDOFF duties (this row + S03-02)

Per protocol Phase handoff + **DR-HANDOFF**. On ablation + regression green → scaffold Phase 07. Do **not** leave Phase 07 as README-only / blocked-until-noticed.

| Who | Duty |
|-----|------|
| **P06-S03-01 (this VERIFY)** | On pass: **Create** Phase 07 scaffold under `docs/phases/phase-07-performance-ladder/` (checklist below) **and** append Phase 07 board section to `docs/TODO.md` with `P07-00` as first pending row after Phase 06’s last row. Record handoff progress in `VERIFY-NOTES.md`. |
| **P06-S03-02 (final review)** | **Owns completion check** — refuse `done` until scaffold is runnable (planner + stubs + board). May finish missing pieces. Marks Phase 06 complete only when handoff + VERIFY evidence agree. |

#### Phase 07 scaffold checklist (minimum)

Create / ensure:

```text
docs/phases/phase-07-performance-ladder/
  README.md                 # goal = Performance ladder & language plugins (A_PROJECT_PLAN Phase 7)
  00-PHASE-PLANNER.md       # runnable (Agent→clarify→Plan→execute); light locks OK
  scopes/                   # ≥1 scope stub folder recommended; minimal OK
    scope-01-…/
      00-PLANNER.md
      01-*.md
      02-scope-review.md
      SCOPE-TODOS.md
```

Board (`docs/TODO.md`): new **Phase 07** section after Phase 06 rows; first pending ID **`P07-00`** → `phase-07-performance-ladder/00-PHASE-PLANNER.md`. Do **not** start executing Phase 07 until Phase 06 S03-02 is `done`.

Deep tasking stays with Phase 07’s own phase/scope planners. This row delivers a **runnable handoff**, not finished implement prompts.

**Counterfactual:** If ablation / primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** scaffold Phase 07 until VERIFY finally greens (or user records `no successor`).

## Board rights
Verify: **status + notes** on `P06-S03-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails; **may create upcoming Phase 07 scaffold** (DR-HANDOFF). Do **not** rewrite Phase 06 `done` history. Do **not** mark `P06-S03-02` or Phase 07 rows `done`.

## Preflight / Plan
1. Re-read this prompt + board row + S01/S02 REVIEW-NOTES + ablation locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`, packages `evals/{honesty,x0,p0x,replan,impact,capability}`, mig `010`, `schema-capability.json`, Gate C metrics under `docs/verification/gate-c-x0/`.
3. Plan: ablation harness check → run locked commands → fill evidence → pass→VERIFY-NOTES + Phase 07 scaffold **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. Capability-selection ablation re-check (required)

Confirm all of:

| Check | Expect |
|-------|--------|
| Harness path | `evals/capability/` exists; `doc.go` names capability-selection ablation |
| Named test | `TestPlantedCapabilitySelectionAblation` PASS under locked command |
| Schema / metrics | Committed `schema-capability.json` v1; test writes schema-valid temp `metrics-capability.json` |
| Tallies | TP=3, FN=0, FP=0, TN=1; precision=1.0; recall=1.0 |
| Probes | Pos-1 UNAVAILABLE / Pos-2 UNKNOWN / Pos-3 selection-filter (no catalog dump) / Neg-1 clean AVAILABLE |
| S01 consume | `UpsertCapability` + `RequireCapability` + `MissingCapabilities` + packet required/missing exercised |
| Honesty of claim | Notes map ablation to this harness — **not** “vibes”, Gate C scores, or commercial multi-model |

Record pass/fail in VERIFY-NOTES evidence table.

### B. Carry-forward bars (required)

| Check | Expect |
|-------|--------|
| Honesty H5 Paths A/B/C | `TestHonestyFailClosedPlantedClaim` PASS |
| Gate G | `TestHonestyEscapeRateGateGPrelim` PASS (escapes=1/caught=2/attempts=3) |
| Gate E | `TestPlantedDiscoveryReplan` PASS |
| Gate F | `TestPlantedImpactConflictsGateFPrelim` PASS (TP=3/FN=0/FP=0/TN=1; P/R=1.0) |
| P0-X 7/7 | `evals/p0x` PASS (`TestP0XAllCriteria` or package) |
| X0 packages | `./evals/x0/...` PASS |
| S01 domain/store/compiler | mig 010 + capability surface / packet PASS |
| Gate C artifacts | `GATE-C-NOTES.md` still **Go**; metrics `dry_run:false`, N=3; means G1 0.800 > B0 0.000 — **re-check files, do not invent new Go** |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation | Explicit in VERIFY-NOTES |
| Full suite | `CGO_ENABLED=1 go test ./... -count=1` PASS (via locked full command) |

### C. Independent re-run (required)

From module root `/home/ali/Desktop/Trace`, run the **Locked verify commands**. All required commands must exit 0 for a green gate. Capture command + PASS/FAIL in `VERIFY-NOTES.md`.

### D. Evidence table (required in `VERIFY-NOTES.md`)

Copy and fill:

| Gate | Result (pass/fail) | Evidence (test / log gist) |
|------|--------------------|----------------------------|
| Capability-selection ablation | | `evals/capability` `TestPlantedCapabilitySelectionAblation` |
| schema-capability.json v1 + temp metrics | | file present; test validation |
| TP=3 FN=0 FP=0 TN=1; P/R=1.0 | | metrics fields / asserts |
| S01 Upsert/Require/Missing + packet required/missing | | ablation harness +/or domain/compiler tests |
| Honesty H5 Paths A/B/C | | `TestHonestyFailClosedPlantedClaim` |
| Gate G prelim | | `evals/honesty` `TestHonestyEscapeRateGateGPrelim` |
| Gate E mini-eval | | `evals/replan` `TestPlantedDiscoveryReplan` |
| Gate F prelim | | `evals/impact` `TestPlantedImpactConflictsGateFPrelim` |
| S01 mig 010 + capability surface | | domain/store/compiler package PASS; schema file present |
| P0-X 7/7 | | `TestP0XAllCriteria` (or package PASS + note) |
| X0 packages | | `./evals/x0/...` PASS |
| Gate C `dry_run:false` intact | | paths under `docs/verification/gate-c-x0/`; GATE-C-NOTES Go |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation | | explicit: Phase 01 dry-run not used as those passes |
| `go test ./...` | | PASS/FAIL |
| Law checks | | no daemon/HTTP primary; no committed `.trace/`; G19; VerifiedFact out; no ontology megastore; no commercial capability theater |
| Residuals (non-blocking) | | DPC-global; non-tx Apply; UNIQUE re-link; MCP no severity; GC-03/04 deferred; S02 lows |
| DR-HANDOFF | | Phase 07 path created? board `P07-00` present? |

Also record: date, go/CGO note, residuals carried forward.

### E. On **all gates pass** + law checks

1. Write `VERIFY-NOTES.md` with verdict **Phase 06 VERIFY PASS / capability-selection ablation green**, confidence, evidence table, residuals as **non-blocking** if primary paths green. Explicitly: ablation = `evals/capability` named planted test; Gate F + Gate G + Gate E + Gate C artifacts intact; Phase 01 dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation.
2. **Start DR-HANDOFF:** create `docs/phases/phase-07-performance-ladder/` per checklist + append Phase 07 section / `P07-00` to `docs/TODO.md`.
3. Board Notes: short “ablation + honesty/Gate G/E/F/p0x/x0/S01–S02 PASS; Gate C intact; Phase 07 scaffold started; see VERIFY-NOTES.md; pending P06-S03-02 handoff close”.
4. Mark `P06-S03-01` **done**. Do **not** mark S03-02 done.

### F. On **any fail**

1. `VERIFY-NOTES.md`: verdict **FAIL**, which gate/law, minimal reproduction.
2. Insert `P06-S03-01a` (implement) + `P06-S03-01b` (review) immediately below this row; set this row **`blocked`** (or `failed` + plan `01c` re-VERIFY) with reason. Prefer new `01c` verify row if this row already closed as failed/blocked (forward-only history).
3. Spawn prompts must be **full** protocol skeletons. Scope remediations to the failing layer — **do not** weaken bars, invent ablation pass from Notes, or rewrite Gate C Mode-B packs.
4. **Forbidden “fixes”:** claiming ablation without `TestPlantedCapabilitySelectionAblation`; claiming Gate C from dry-run; requiring MCP for ablation; adding daemon/HTTP; weakening honesty Paths A/B/C; rewriting `done` S01–S02 prompts; promoting GC-03/04 without evidence; scaffolding Phase 07 on a red VERIFY; inventing VerifiedFact; shipping ontology megastore / commercial multi-model capability theater.

### Spawn ID convention

```text
… P06-S03-01  (this VERIFY)
… P06-S03-01a (remediation implement)   ← insert immediately below
… P06-S03-01b (remediation review)
… P06-S03-01c (re-VERIFY)               ← if original VERIFY closed as failed/blocked
… P06-S03-02  (phase review — only after VERIFY finally done; owns DR-HANDOFF completion)
```

Update `SCOPE-TODOS.md` when spawning.

## Out of scope
- Implementing Phase 07 performance-ladder / language-plugin product features
- Re-running live multi-model Gate C / flipping Go without contradicting evidence
- Rewriting Mode-B packs
- Promoting GC-03/04 without measurement need
- Declaring A1 / product thesis commercially validated
- Starting Phase 07 implement wave before S03-02 closes
- Expanding ablation into a multi-model scored commercial capability benchmark (prelim = planted automated P/R only)
- Introducing VerifiedFact promotion engine
- Implementing `plan simulate` / ontology megastore

## Todo updates
Status + Notes on `P06-S03-01`; spawn rows if needed; Phase 07 upcoming artifacts on pass; `SCOPE-TODOS.md` checkboxes.

## Exit criteria
- [ ] Independent ablation (`TestPlantedCapabilitySelectionAblation` + schema/metrics + tallies + S01 hooks) + honesty A/B/C + Gate G + Gate E + Gate F + p0x + x0 + S01 surfaces + Gate C artifact integrity + `./...` recorded in `VERIFY-NOTES.md`
- [ ] Evidence table includes ablation path + dry-run≠Gate C/≠Gate F/≠Gate G/≠ablation + residuals + handoff
- [ ] Law checks recorded; Mode-B packs not falsified; no commercial capability theater
- [ ] Either gates green **and** Phase 07 scaffold started (folder + `P07-00` board) **or** remediations spawned with full prompts + this row blocked/failed honestly
- [ ] TODO.md status + Notes updated; SCOPE-TODOS synced
- [ ] No product feature Go on this row

## Minimal todos
- [ ] Preflight: S01/S02 REVIEW-NOTES + mig 010 + schema-capability.json + Gate C metrics + ablation harness path
- [ ] Re-check ablation (`TestPlantedCapabilitySelectionAblation` + tallies + S01 hooks)
- [ ] Run honesty A/B/C / Gate G / Gate E / Gate F / p0x / x0 / domain+store+planner+compiler / Gate C artifact inspect / `./...` independently
- [ ] Write `VERIFY-NOTES.md` evidence table
- [ ] Pass → scaffold `phase-07-performance-ladder` + board `P07-00`; Fail → spawn 01a/01b (+01c plan)
- [ ] Board status + SCOPE-TODOS
