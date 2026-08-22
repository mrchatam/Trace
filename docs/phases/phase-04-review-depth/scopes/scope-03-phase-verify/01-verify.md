# P04 / S03 / 01 — Phase 04 VERIFY

## Metadata
- id: P04-S03-01
- todo_ids: [P04-S03-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently re-prove Phase 04 closeout — Gate G prelim + S01/S02 review-depth surfaces + honesty/p0x/x0/replan/Gate E/Gate C carry-forward bars — against live packages. Do **not** trust S01/S02 Notes alone. Do **not** reopen Gate C, invent full production Gate G without the named harness, or claim commercial A1 / product thesis.

Write durable evidence, then either:

1. **Pass** → declare **Phase 04 VERIFY PASS / Gate G prelim green** + **start DR-HANDOFF** Phase 05 (`phase-05-decision-impact` + board `P05-00`), or
2. **Fail** → **spawn forward-only remediations** (01a/01b/+01c).

No product features.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF, DR-AGENT
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 4 Gate G; Phase 5 Decision impact & simulation
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) — Gate G: review reduces false completion
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — H5 / Gate G
- Sibling: [S01 REVIEW-NOTES.md](../scope-01-scope-review-layer/REVIEW-NOTES.md), [S02 REVIEW-NOTES.md](../scope-02-honesty-escape-rate/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 03 VERIFY [`../../../phase-03-progressive-planner/scopes/scope-03-phase-verify/01-verify.md`](../../../phase-03-progressive-planner/scopes/scope-03-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate G (definition) | Review reduces false completion (`PROJECT_MODEL_SNAPSHOT` / H5 depth) |
| Gate G path / harness | **`evals/honesty`** package; named test **`TestHonestyEscapeRateGateGPrelim`** |
| Gate G pass bar | Planted escape-rate report: **escapes=1 / caught=2 / attempts=3** (escape_rate≈1/3); hatch counted as escape **only** in this report; schema-valid temp **`metrics-gate-g.json`** vs committed **`schema-gate-g.json`** v1 |
| S01 surface (must stay green) | mig **`008_scope_review.sql`** `review_residuals`; `LinkReviewScope` / `review_judges_scope`→`plan_scope`; residual Add/List/CountOpen/SetStatus; **`CountOpenResidualsByScope`**; OPEN **`POLICY_EXCEPTION`** exercised by Gate G harness |
| S02 surface | Extend `evals/honesty` only (no new package); keep Paths A/B/C (`TestHonestyFailClosedPlantedClaim`) untouched in fail-closed proof |
| Honesty Paths A/B/C | Fail-closed; no `AllowDoneWithoutReview` in A/B/C proof |
| Gate E | **Green** — `evals/replan` **`TestPlantedDiscoveryReplan`** (PLAN_AFFECTING+; churn N=5 + ack) |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** artifacts under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; means G1 0.800 > B0 0.000; **do not invent new Go**; Mode-B packs historical |
| Dry-run ≠ Gate C | Phase 01 dry-run is regression-only — **not** Gate C pass; **not** Gate G pass |
| Fixture hash pin (carry) | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Deferrals (carry) | GC-03/04 still deferred unless VERIFY Notes promote with evidence |
| Residuals (non-blocking) | DPC-global attach; non-tx Apply; UNIQUE re-link; MCP no severity; schema `s01_hooks` minItems-only (S02 low) |
| VerifiedFact | Still **out** — residuals only, not promotion engine |
| Product Go | **Forbidden** — verify + evidence + handoff scaffold + spawn prompts only |
| MCP | Optional; Gate G is honesty harness — do not require MCP for VERIFY pass |
| Daemon / HTTP / embeddings | Still forbidden as primary |
| Phase 05 folder name | **`phase-05-decision-impact`** (A_PROJECT_PLAN Phase 5 — Decision impact & simulation; refine only with Notes) |

### Locked verify commands

```bash
# Gate G prelim (primary — named escape-rate report)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run TestHonestyEscapeRateGateGPrelim

# Honesty package: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E carry-forward
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan

# S01/S02 supporting surfaces (domain+store+planner; residuals / scope review)
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... -count=1

# Full regression bar (includes p0x 7/7, x0, replan, honesty, analyzers)
CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -v -run TestHonestyEscapeRateGateGPrelim
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -v -run TestHonestyFailClosedPlantedClaim
CGO_ENABLED=1 go test ./evals/p0x/... -count=1 -v -run TestP0XAllCriteria
# Confirm schema file present (metrics are temp — written by Gate G test)
test -f evals/honesty/schema-gate-g.json
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3, means match GATE-C-NOTES — do not re-score
# Fixture content hash (must still match pin unless Notes promote change)
find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# expect: 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/`
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] Gate G evidence is **`evals/honesty` `TestHonestyEscapeRateGateGPrelim`** + schema/metrics — not vibes / Notes-only
- [ ] Planted tallies: escapes=1, caught=2, attempts=3; hatch=escape only in Gate G report
- [ ] S01 hooks remain green: `LinkReviewScope` / `CountOpenResidualsByScope` / OPEN `POLICY_EXCEPTION`
- [ ] Paths A/B/C still fail-closed without hatch in that proof
- [ ] Gate E = `TestPlantedDiscoveryReplan` still green
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone; no silent Go flip
- [ ] Mode-B packs not falsified
- [ ] Embeddings still absent; VerifiedFact absent as promotion engine
- [ ] Mig 008 present and exercised by green domain/store (+ Gate G) tests
- [ ] GC-03/04 remain deferred unless explicitly promoted in VERIFY-NOTES

### DR-HANDOFF duties (this row + S03-02)

Per protocol Phase handoff + **DR-HANDOFF**. On Gate G + regression green → scaffold Phase 05. Do **not** leave Phase 05 as README-only / blocked-until-noticed.

| Who | Duty |
|-----|------|
| **P04-S03-01 (this VERIFY)** | On pass: **Create** Phase 05 scaffold under `docs/phases/phase-05-decision-impact/` (checklist below) **and** append Phase 05 board section to `docs/TODO.md` with `P05-00` as first pending row after Phase 04’s last row. Record handoff progress in `VERIFY-NOTES.md`. |
| **P04-S03-02 (final review)** | **Owns completion check** — refuse `done` until scaffold is runnable (planner + stubs + board). May finish missing pieces. Marks Phase 04 complete only when handoff + VERIFY evidence agree. |

#### Phase 05 scaffold checklist (minimum)

Create / ensure:

```text
docs/phases/phase-05-decision-impact/
  README.md                 # goal = Decision impact & simulation (A_PROJECT_PLAN Phase 5)
  00-PHASE-PLANNER.md       # runnable (Agent→clarify→Plan→execute); light locks OK
  scopes/                   # ≥1 scope stub folder recommended; minimal OK
    scope-01-…/
      00-PLANNER.md
      01-*.md
      02-scope-review.md
      SCOPE-TODOS.md
```

Board (`docs/TODO.md`): new **Phase 05** section after Phase 04 rows; first pending ID **`P05-00`** → `phase-05-decision-impact/00-PHASE-PLANNER.md`. Do **not** start executing Phase 05 until Phase 04 S03-02 is `done`.

Deep tasking stays with Phase 05’s own phase/scope planners. This row delivers a **runnable handoff**, not finished implement prompts.

**Counterfactual:** If Gate G / primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** scaffold Phase 05 until VERIFY finally greens (or user records `no successor`).

## Board rights
Verify: **status + notes** on `P04-S03-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails; **may create upcoming Phase 05 scaffold** (DR-HANDOFF). Do **not** rewrite Phase 04 `done` history. Do **not** mark `P04-S03-02` or Phase 05 rows `done`.

## Preflight / Plan
1. Re-read this prompt + board row + S01/S02 REVIEW-NOTES + Gate G locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`, packages `evals/{honesty,x0,p0x,replan}`, mig `008`, `schema-gate-g.json`, Gate C metrics under `docs/verification/gate-c-x0/`.
3. Plan: Gate G harness check → run locked commands → fill evidence → pass→VERIFY-NOTES + Phase 05 scaffold **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. Gate G prelim re-check (required)

Confirm all of:

| Check | Expect |
|-------|--------|
| Harness path | `evals/honesty/` exists; `doc.go` names Gate G prelim |
| Named test | `TestHonestyEscapeRateGateGPrelim` PASS under locked command |
| Schema / metrics | Committed `schema-gate-g.json` v1; test writes schema-valid temp `metrics-gate-g.json` |
| Tallies | escapes=1, caught=2, attempts=3 (escape_rate≈1/3); hatch=escape only |
| S01 consume | `LinkReviewScope` + OPEN `POLICY_EXCEPTION` + `CountOpenResidualsByScope` exercised |
| Paths A/B/C | `TestHonestyFailClosedPlantedClaim` still PASS; hatch not used in that proof |
| Honesty of claim | Notes map Gate G to this harness — **not** “vibes” or Gate C scores |

Record pass/fail in VERIFY-NOTES evidence table.

### B. Carry-forward bars (required)

| Check | Expect |
|-------|--------|
| Honesty H5 Paths A/B/C | `TestHonestyFailClosedPlantedClaim` PASS |
| Gate E | `TestPlantedDiscoveryReplan` PASS |
| P0-X 7/7 | `evals/p0x` PASS (`TestP0XAllCriteria` or package) |
| X0 packages | `./evals/x0/...` PASS |
| S01 domain/store | mig 008 + residual/`LinkReviewScope` surfaces PASS |
| Gate C artifacts | `GATE-C-NOTES.md` still **Go**; metrics `dry_run:false`, N=3; means G1 0.800 > B0 0.000 — **re-check files, do not invent new Go** |
| Dry-run ≠ Gate C | Explicit in VERIFY-NOTES (also: dry-run ≠ Gate G) |
| Full suite | `CGO_ENABLED=1 go test ./... -count=1` PASS (via locked full command) |

### C. Independent re-run (required)

From module root `/home/ali/Desktop/Trace`, run the **Locked verify commands**. All required commands must exit 0 for a green gate. Capture command + PASS/FAIL in `VERIFY-NOTES.md`.

### D. Evidence table (required in `VERIFY-NOTES.md`)

Copy and fill:

| Gate | Result (pass/fail) | Evidence (test / log gist) |
|------|--------------------|----------------------------|
| Gate G prelim | | `evals/honesty` `TestHonestyEscapeRateGateGPrelim` |
| schema-gate-g.json v1 + temp metrics | | file present; test validation |
| escapes=1 / caught=2 / attempts=3 | | metrics fields / asserts |
| S01 `LinkReviewScope` + OPEN `POLICY_EXCEPTION` + `CountOpenResidualsByScope` | | Gate G harness +/or domain tests |
| Honesty H5 Paths A/B/C | | `TestHonestyFailClosedPlantedClaim` |
| Gate E mini-eval | | `evals/replan` `TestPlantedDiscoveryReplan` |
| S01 mig 008 + residuals surface | | domain/store package PASS; schema file present |
| P0-X 7/7 | | `TestP0XAllCriteria` (or package PASS + note) |
| X0 packages | | `./evals/x0/...` PASS |
| Gate C `dry_run:false` intact | | paths under `docs/verification/gate-c-x0/`; GATE-C-NOTES Go |
| Dry-run ≠ Gate C | | explicit: Phase 01 dry-run not used as Gate C / Gate G |
| `go test ./...` | | PASS/FAIL |
| Law checks | | no daemon/HTTP primary; no committed `.trace/`; G19; VerifiedFact out |
| Residuals (non-blocking) | | DPC-global; non-tx Apply; UNIQUE re-link; MCP no severity; s01_hooks schema looseness |
| DR-HANDOFF | | Phase 05 path created? board `P05-00` present? |

Also record: date, go/CGO note, residuals carried forward.

### E. On **all gates pass** + law checks

1. Write `VERIFY-NOTES.md` with verdict **Phase 04 VERIFY PASS / Gate G prelim green**, confidence, evidence table, residuals as **non-blocking** if primary paths green. Explicitly: Gate G = `evals/honesty` named escape-rate test; Gate E + Gate C artifacts intact; Phase 01 dry-run ≠ Gate C.
2. **Start DR-HANDOFF:** create `docs/phases/phase-05-decision-impact/` per checklist + append Phase 05 section / `P05-00` to `docs/TODO.md`.
3. Board Notes: short “Gate G + honesty/p0x/x0/replan/S01–S02 PASS; Gate C intact; Phase 05 scaffold started; see VERIFY-NOTES.md; pending P04-S03-02 handoff close”.
4. Mark `P04-S03-01` **done**. Do **not** mark S03-02 done.

### F. On **any fail**

1. `VERIFY-NOTES.md`: verdict **FAIL**, which gate/law, minimal reproduction.
2. Insert `P04-S03-01a` (implement) + `P04-S03-01b` (review) immediately below this row; set this row **`blocked`** (or `failed` + plan `01c` re-VERIFY) with reason. Prefer new `01c` verify row if this row already closed as failed/blocked (forward-only history).
3. Spawn prompts must be **full** protocol skeletons. Scope remediations to the failing layer — **do not** weaken bars, invent Gate G pass from Notes, or rewrite Gate C Mode-B packs.
4. **Forbidden “fixes”:** claiming Gate G without `TestHonestyEscapeRateGateGPrelim`; claiming Gate C from dry-run; requiring MCP for Gate G; adding daemon/HTTP; weakening honesty Paths A/B/C; rewriting `done` S01–S02 prompts; promoting GC-03/04 without evidence; scaffolding Phase 05 on a red VERIFY; inventing VerifiedFact.

### Spawn ID convention

```text
… P04-S03-01  (this VERIFY)
… P04-S03-01a (remediation implement)   ← insert immediately below
… P04-S03-01b (remediation review)
… P04-S03-01c (re-VERIFY)               ← if original VERIFY closed as failed/blocked
… P04-S03-02  (phase review — only after VERIFY finally done; owns DR-HANDOFF completion)
```

Update `SCOPE-TODOS.md` when spawning.

## Out of scope
- Implementing Phase 05 decision-impact / simulate product features
- Re-running live multi-model Gate C / flipping Go without contradicting evidence
- Rewriting Mode-B packs
- Promoting GC-03/04 without measurement need
- Declaring A1 / product thesis commercially validated
- Starting Phase 05 implement wave before S03-02 closes
- Expanding Gate G into a multi-model scored commercial review benchmark (prelim = planted automated escape-rate report only)
- Introducing VerifiedFact promotion engine

## Todo updates
Status + Notes on `P04-S03-01`; spawn rows if needed; Phase 05 upcoming artifacts on pass; `SCOPE-TODOS.md` checkboxes.

## Exit criteria
- [ ] Independent Gate G (`TestHonestyEscapeRateGateGPrelim` + schema/metrics + tallies + S01 hooks) + honesty A/B/C + Gate E + p0x + x0 + S01 surfaces + Gate C artifact integrity + `./...` recorded in `VERIFY-NOTES.md`
- [ ] Evidence table includes Gate G path + dry-run≠Gate C + residuals + handoff
- [ ] Law checks recorded; Mode-B packs not falsified
- [ ] Either gates green **and** Phase 05 scaffold started (folder + `P05-00` board) **or** remediations spawned with full prompts + this row blocked/failed honestly
- [ ] TODO.md status + Notes updated; SCOPE-TODOS synced
- [ ] No product feature Go on this row

## Minimal todos
- [ ] Preflight: S01/S02 REVIEW-NOTES + mig 008 + schema-gate-g.json + Gate C metrics + Gate G harness path
- [ ] Re-check Gate G (`TestHonestyEscapeRateGateGPrelim` + tallies + S01 hooks)
- [ ] Run honesty A/B/C / Gate E / p0x / x0 / domain+store+planner / Gate C artifact inspect / `./...` independently
- [ ] Write `VERIFY-NOTES.md` evidence table
- [ ] Pass → scaffold `phase-05-decision-impact` + board `P05-00`; Fail → spawn 01a/01b (+01c plan)
- [ ] Board status + SCOPE-TODOS
