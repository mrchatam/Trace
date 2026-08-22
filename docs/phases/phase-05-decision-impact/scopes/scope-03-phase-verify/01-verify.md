# P05 / S03 / 01 — Phase 05 VERIFY

## Metadata
- id: P05-S03-01
- todo_ids: [P05-S03-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently re-prove Phase 05 closeout — Gate F prelim + S01/S02 decision-impact surfaces + honesty/Gate G/Gate E/p0x/x0/Gate C carry-forward bars — against live packages. Do **not** trust S01/S02 Notes alone. Do **not** reopen Gate C, invent commercial multi-model Gate F, claim `plan simulate`, or declare commercial A1 / product thesis.

Write durable evidence, then either:

1. **Pass** → declare **Phase 05 VERIFY PASS / Gate F prelim green** + **start DR-HANDOFF** Phase 06 (`phase-06-environment-capability` + board `P06-00`), or
2. **Fail** → **spawn forward-only remediations** (01a/01b/+01c).

No product features.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF, DR-AGENT, DR-NOIMP
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 5 Gate F; Phase 6 Environment/capability graph
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) — Gate F: impact precision/recall on planted conflicts
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — Gate F
- Sibling: [S01 REVIEW-NOTES.md](../scope-01-impact-classes/REVIEW-NOTES.md), [S02 REVIEW-NOTES.md](../scope-02-gate-f-prelim/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 04 VERIFY [`../../../phase-04-review-depth/scopes/scope-03-phase-verify/01-verify.md`](../../../phase-04-review-depth/scopes/scope-03-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate F (definition) | Impact analysis precision/recall on planted conflicts (`PROJECT_MODEL_SNAPSHOT`) |
| Gate F path / harness | **`evals/impact`** package; named test **`TestPlantedImpactConflictsGateFPrelim`** |
| Gate F pass bar | Planted 4-probe P/R: **TP=3 FN=0 FP=0 TN=1**; **precision=1.0**; **recall=1.0**; schema-valid temp **`metrics-gate-f.json`** vs committed **`schema-gate-f.json`** v1 |
| S01 surface (must stay green) | mig **`009_decision_impact.sql`** (`decision_impact_findings`+`decision_alternatives`); Add/List findings+alternatives; fail-closed **`ImpactReport`**; CLI `trace impact`; **no new** entity_links rels (keep `decision_affects_task`) |
| S02 surface | New package `evals/impact` only; S01 hooks only (`AddImpactFinding` / `LinkDecisionTask` / `ImpactReport`); mig 009 only (no S02 schema fork) |
| Honesty Paths A/B/C | Fail-closed; no `AllowDoneWithoutReview` in A/B/C proof — `TestHonestyFailClosedPlantedClaim` |
| Gate G | **Green** — `evals/honesty` **`TestHonestyEscapeRateGateGPrelim`** (escapes=1/caught=2/attempts=3) |
| Gate E | **Green** — `evals/replan` **`TestPlantedDiscoveryReplan`** (PLAN_AFFECTING+; churn N=5 + ack) |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** artifacts under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; means G1 0.800 > B0 0.000; **do not invent new Go**; Mode-B packs historical |
| Dry-run ≠ Gate C / ≠ Gate F | Phase 01 dry-run is regression-only — **not** Gate C pass; **not** Gate F pass; also ≠ Gate G |
| Fixture hash pin (carry) | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Deferrals (carry) | GC-03/04 still deferred; **`plan simulate`** still out (P13) |
| Residuals (non-blocking) | DPC-global attach; non-tx Apply; UNIQUE re-link; MCP no severity; s01_hooks schema looseness; Pos-1 does not trust OverallClass alone when HasUnknown required |
| VerifiedFact | Still **out** — not a promotion engine |
| DR-NOIMP | Still PROVISIONAL — Phase 05 = planted Gate F prelim, **not** commercial multi-model impact engine |
| Product Go | **Forbidden** — verify + evidence + handoff scaffold + spawn prompts only |
| MCP | Optional; Gate F is impact harness — do not require MCP for VERIFY pass |
| Daemon / HTTP / embeddings | Still forbidden as primary |
| Phase 06 folder name | **`phase-06-environment-capability`** (A_PROJECT_PLAN Phase 6 — Environment/capability graph; refine only with Notes) |

### Locked verify commands

```bash
# Gate F prelim (primary — named planted P/R harness)
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim

# Impact package (CGO-free)
CGO_ENABLED=0 go test ./evals/impact/... -count=1

# Honesty: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E carry-forward
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan

# S01/S02 supporting surfaces (domain+store+planner; impact findings / report)
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... -count=1

# Full regression bar (includes p0x 7/7, x0, replan, honesty, impact, analyzers)
CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./evals/impact/... ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -v -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -v -run TestHonestyEscapeRateGateGPrelim
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -v -run TestHonestyFailClosedPlantedClaim
CGO_ENABLED=1 go test ./evals/p0x/... -count=1 -v -run TestP0XAllCriteria
# Confirm schema file present (metrics are temp — written by Gate F test)
test -f evals/impact/schema-gate-f.json
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3, means match GATE-C-NOTES — do not re-score
# Fixture content hash (must still match pin unless Notes promote change)
find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# expect: 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/`
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] Gate F evidence is **`evals/impact` `TestPlantedImpactConflictsGateFPrelim`** + schema/metrics — not vibes / Notes-only
- [ ] Planted tallies: TP=3, FN=0, FP=0, TN=1; precision=1.0; recall=1.0
- [ ] S01 hooks remain green: `AddImpactFinding` / `LinkDecisionTask` / `ImpactReport` / `decision_affects_task`
- [ ] Mig 009 present; no commercial `internal/impact` stack; no `plan simulate`
- [ ] Paths A/B/C still fail-closed without hatch in that proof
- [ ] Gate G = `TestHonestyEscapeRateGateGPrelim` still green
- [ ] Gate E = `TestPlantedDiscoveryReplan` still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone; no silent Go flip
- [ ] Mode-B packs not falsified
- [ ] Embeddings still absent; VerifiedFact absent as promotion engine
- [ ] DR-NOIMP not violated (no commercial multi-model Gate F claim)
- [ ] GC-03/04 remain deferred unless explicitly promoted in VERIFY-NOTES

### DR-HANDOFF duties (this row + S03-02)

Per protocol Phase handoff + **DR-HANDOFF**. On Gate F + regression green → scaffold Phase 06. Do **not** leave Phase 06 as README-only / blocked-until-noticed.

| Who | Duty |
|-----|------|
| **P05-S03-01 (this VERIFY)** | On pass: **Create** Phase 06 scaffold under `docs/phases/phase-06-environment-capability/` (checklist below) **and** append Phase 06 board section to `docs/TODO.md` with `P06-00` as first pending row after Phase 05’s last row. Record handoff progress in `VERIFY-NOTES.md`. |
| **P05-S03-02 (final review)** | **Owns completion check** — refuse `done` until scaffold is runnable (planner + stubs + board). May finish missing pieces. Marks Phase 05 complete only when handoff + VERIFY evidence agree. |

#### Phase 06 scaffold checklist (minimum)

Create / ensure:

```text
docs/phases/phase-06-environment-capability/
  README.md                 # goal = Environment/capability graph (A_PROJECT_PLAN Phase 6)
  00-PHASE-PLANNER.md       # runnable (Agent→clarify→Plan→execute); light locks OK
  scopes/                   # ≥1 scope stub folder recommended; minimal OK
    scope-01-…/
      00-PLANNER.md
      01-*.md
      02-scope-review.md
      SCOPE-TODOS.md
```

Board (`docs/TODO.md`): new **Phase 06** section after Phase 05 rows; first pending ID **`P06-00`** → `phase-06-environment-capability/00-PHASE-PLANNER.md`. Do **not** start executing Phase 06 until Phase 05 S03-02 is `done`.

Deep tasking stays with Phase 06’s own phase/scope planners. This row delivers a **runnable handoff**, not finished implement prompts.

**Counterfactual:** If Gate F / primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** scaffold Phase 06 until VERIFY finally greens (or user records `no successor`).

## Board rights
Verify: **status + notes** on `P05-S03-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails; **may create upcoming Phase 06 scaffold** (DR-HANDOFF). Do **not** rewrite Phase 05 `done` history. Do **not** mark `P05-S03-02` or Phase 06 rows `done`.

## Preflight / Plan
1. Re-read this prompt + board row + S01/S02 REVIEW-NOTES + Gate F locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`, packages `evals/{honesty,x0,p0x,replan,impact}`, mig `009`, `schema-gate-f.json`, Gate C metrics under `docs/verification/gate-c-x0/`.
3. Plan: Gate F harness check → run locked commands → fill evidence → pass→VERIFY-NOTES + Phase 06 scaffold **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. Gate F prelim re-check (required)

Confirm all of:

| Check | Expect |
|-------|--------|
| Harness path | `evals/impact/` exists; `doc.go` names Gate F prelim |
| Named test | `TestPlantedImpactConflictsGateFPrelim` PASS under locked command |
| Schema / metrics | Committed `schema-gate-f.json` v1; test writes schema-valid temp `metrics-gate-f.json` |
| Tallies | TP=3, FN=0, FP=0, TN=1; precision=1.0; recall=1.0 |
| Probes | Pos-1 UNKNOWN / Pos-2 DESTRUCTIVE rollup / Pos-3 link+empty findings / Neg-1 clean SAFE |
| S01 consume | `AddImpactFinding` + `LinkDecisionTask` + `ImpactReport` + `decision_affects_task` exercised |
| Honesty of claim | Notes map Gate F to this harness — **not** “vibes”, Gate C scores, or commercial multi-model |

Record pass/fail in VERIFY-NOTES evidence table.

### B. Carry-forward bars (required)

| Check | Expect |
|-------|--------|
| Honesty H5 Paths A/B/C | `TestHonestyFailClosedPlantedClaim` PASS |
| Gate G | `TestHonestyEscapeRateGateGPrelim` PASS (escapes=1/caught=2/attempts=3) |
| Gate E | `TestPlantedDiscoveryReplan` PASS |
| P0-X 7/7 | `evals/p0x` PASS (`TestP0XAllCriteria` or package) |
| X0 packages | `./evals/x0/...` PASS |
| S01 domain/store | mig 009 + findings/alternatives/`ImpactReport` surfaces PASS |
| Gate C artifacts | `GATE-C-NOTES.md` still **Go**; metrics `dry_run:false`, N=3; means G1 0.800 > B0 0.000 — **re-check files, do not invent new Go** |
| Dry-run ≠ Gate C / ≠ Gate F | Explicit in VERIFY-NOTES (also: dry-run ≠ Gate G) |
| Full suite | `CGO_ENABLED=1 go test ./... -count=1` PASS (via locked full command) |

### C. Independent re-run (required)

From module root `/home/ali/Desktop/Trace`, run the **Locked verify commands**. All required commands must exit 0 for a green gate. Capture command + PASS/FAIL in `VERIFY-NOTES.md`.

### D. Evidence table (required in `VERIFY-NOTES.md`)

Copy and fill:

| Gate | Result (pass/fail) | Evidence (test / log gist) |
|------|--------------------|----------------------------|
| Gate F prelim | | `evals/impact` `TestPlantedImpactConflictsGateFPrelim` |
| schema-gate-f.json v1 + temp metrics | | file present; test validation |
| TP=3 FN=0 FP=0 TN=1; P/R=1.0 | | metrics fields / asserts |
| S01 `AddImpactFinding` + `LinkDecisionTask` + `ImpactReport` | | Gate F harness +/or domain tests |
| Honesty H5 Paths A/B/C | | `TestHonestyFailClosedPlantedClaim` |
| Gate G prelim | | `evals/honesty` `TestHonestyEscapeRateGateGPrelim` |
| Gate E mini-eval | | `evals/replan` `TestPlantedDiscoveryReplan` |
| S01 mig 009 + impact surface | | domain/store package PASS; schema file present |
| P0-X 7/7 | | `TestP0XAllCriteria` (or package PASS + note) |
| X0 packages | | `./evals/x0/...` PASS |
| Gate C `dry_run:false` intact | | paths under `docs/verification/gate-c-x0/`; GATE-C-NOTES Go |
| Dry-run ≠ Gate C / ≠ Gate F | | explicit: Phase 01 dry-run not used as Gate C / Gate F / Gate G |
| `go test ./...` | | PASS/FAIL |
| Law checks | | no daemon/HTTP primary; no committed `.trace/`; G19; VerifiedFact out; DR-NOIMP; no `plan simulate` |
| Residuals (non-blocking) | | DPC-global; non-tx Apply; UNIQUE re-link; MCP no severity; s01_hooks; GC-03/04 deferred |
| DR-HANDOFF | | Phase 06 path created? board `P06-00` present? |

Also record: date, go/CGO note, residuals carried forward.

### E. On **all gates pass** + law checks

1. Write `VERIFY-NOTES.md` with verdict **Phase 05 VERIFY PASS / Gate F prelim green**, confidence, evidence table, residuals as **non-blocking** if primary paths green. Explicitly: Gate F = `evals/impact` named planted test; Gate G + Gate E + Gate C artifacts intact; Phase 01 dry-run ≠ Gate C / ≠ Gate F.
2. **Start DR-HANDOFF:** create `docs/phases/phase-06-environment-capability/` per checklist + append Phase 06 section / `P06-00` to `docs/TODO.md`.
3. Board Notes: short “Gate F + honesty/Gate G/Gate E/p0x/x0/S01–S02 PASS; Gate C intact; Phase 06 scaffold started; see VERIFY-NOTES.md; pending P05-S03-02 handoff close”.
4. Mark `P05-S03-01` **done**. Do **not** mark S03-02 done.

### F. On **any fail**

1. `VERIFY-NOTES.md`: verdict **FAIL**, which gate/law, minimal reproduction.
2. Insert `P05-S03-01a` (implement) + `P05-S03-01b` (review) immediately below this row; set this row **`blocked`** (or `failed` + plan `01c` re-VERIFY) with reason. Prefer new `01c` verify row if this row already closed as failed/blocked (forward-only history).
3. Spawn prompts must be **full** protocol skeletons. Scope remediations to the failing layer — **do not** weaken bars, invent Gate F pass from Notes, or rewrite Gate C Mode-B packs.
4. **Forbidden “fixes”:** claiming Gate F without `TestPlantedImpactConflictsGateFPrelim`; claiming Gate C from dry-run; requiring MCP for Gate F; adding daemon/HTTP; weakening honesty Paths A/B/C; rewriting `done` S01–S02 prompts; promoting GC-03/04 without evidence; scaffolding Phase 06 on a red VERIFY; inventing VerifiedFact; shipping `plan simulate` / commercial multi-model Gate F.

### Spawn ID convention

```text
… P05-S03-01  (this VERIFY)
… P05-S03-01a (remediation implement)   ← insert immediately below
… P05-S03-01b (remediation review)
… P05-S03-01c (re-VERIFY)               ← if original VERIFY closed as failed/blocked
… P05-S03-02  (phase review — only after VERIFY finally done; owns DR-HANDOFF completion)
```

Update `SCOPE-TODOS.md` when spawning.

## Out of scope
- Implementing Phase 06 environment/capability product features
- Re-running live multi-model Gate C / flipping Go without contradicting evidence
- Rewriting Mode-B packs
- Promoting GC-03/04 without measurement need
- Declaring A1 / product thesis commercially validated
- Starting Phase 06 implement wave before S03-02 closes
- Expanding Gate F into a multi-model scored commercial impact benchmark (prelim = planted automated P/R only)
- Introducing VerifiedFact promotion engine
- Implementing `plan simulate` / PlanVersion branches

## Todo updates
Status + Notes on `P05-S03-01`; spawn rows if needed; Phase 06 upcoming artifacts on pass; `SCOPE-TODOS.md` checkboxes.

## Exit criteria
- [ ] Independent Gate F (`TestPlantedImpactConflictsGateFPrelim` + schema/metrics + tallies + S01 hooks) + honesty A/B/C + Gate G + Gate E + p0x + x0 + S01 surfaces + Gate C artifact integrity + `./...` recorded in `VERIFY-NOTES.md`
- [ ] Evidence table includes Gate F path + dry-run≠Gate C/≠Gate F + residuals + handoff
- [ ] Law checks recorded; Mode-B packs not falsified; DR-NOIMP respected
- [ ] Either gates green **and** Phase 06 scaffold started (folder + `P06-00` board) **or** remediations spawned with full prompts + this row blocked/failed honestly
- [ ] TODO.md status + Notes updated; SCOPE-TODOS synced
- [ ] No product feature Go on this row

## Minimal todos
- [ ] Preflight: S01/S02 REVIEW-NOTES + mig 009 + schema-gate-f.json + Gate C metrics + Gate F harness path
- [ ] Re-check Gate F (`TestPlantedImpactConflictsGateFPrelim` + tallies + S01 hooks)
- [ ] Run honesty A/B/C / Gate G / Gate E / p0x / x0 / domain+store+planner / Gate C artifact inspect / `./...` independently
- [ ] Write `VERIFY-NOTES.md` evidence table
- [ ] Pass → scaffold `phase-06-environment-capability` + board `P06-00`; Fail → spawn 01a/01b (+01c plan)
- [ ] Board status + SCOPE-TODOS
