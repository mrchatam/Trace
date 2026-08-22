# P03 / S03 / 01 — Phase 03 VERIFY

## Metadata
- id: P03-S03-01
- todo_ids: [P03-S03-01]
- role: verify
- skills: [incremental-implementation, debugging-and-error-recovery]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Phase gate (not a feature row):** independently re-prove Phase 03 closeout — Gate E mini-eval + S01/S02 planner surfaces + honesty/p0x/x0/Gate C carry-forward bars — against live packages. Do **not** trust S01/S02 Notes alone. Do **not** reopen Gate C or claim commercial A1 / product thesis.

Write durable evidence, then either:

1. **Pass** → declare **Phase 03 VERIFY PASS / Gate E mini-eval green** + **start DR-HANDOFF** Phase 04 (`phase-04-review-depth` + board `P04-00`), or
2. **Fail** → **spawn forward-only remediations** (01a/01b/+01c).

No product features.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff; VERIFY may spawn
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-HANDOFF, DR-AGENT
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 3 Gate E; Phase 4 Review depth & evidence policies
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) — Gate E: discovery propagation reduces downstream rework
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G16 churn
- Sibling: [S01 REVIEW-NOTES.md](../scope-01-coarse-planner/REVIEW-NOTES.md), [S02 REVIEW-NOTES.md](../scope-02-discovery-replan/REVIEW-NOTES.md)
- Gate C artifacts (carry-forward): [`docs/verification/gate-c-x0/`](../../../../verification/gate-c-x0/), [GATE-C-NOTES.md](../../../phase-02-gate-c/scopes/scope-01-x0-gate-c/GATE-C-NOTES.md)
- Phase README: [../../README.md](../../README.md)
- Pattern: Phase 02 VERIFY [`../../../phase-02-gate-c/scopes/scope-03-phase-verify/01-verify.md`](../../../phase-02-gate-c/scopes/scope-03-phase-verify/01-verify.md)
- Sibling locks: [00-PLANNER.md](00-PLANNER.md)

## Session start
Follow agent-loop-protocol: Agent → clarify if needed → Plan → execute (verify).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate E (definition) | Discovery propagation reduces downstream rework (`PROJECT_MODEL_SNAPSHOT` §15) |
| Gate E path / harness | **`evals/replan`** package; named test **`TestPlantedDiscoveryReplan`** |
| Gate E pass bar | Planted Goal→coarse→deep → **PLAN_AFFECTING+** supersedes deep plan + links Discovery→PlanChange; **INFO** does **not** auto-replan; churn **N=5** fail-closed then **AckReplan** recovers |
| Severity gate | Only `PLAN_AFFECTING` and `BLOCKING` auto-replan; `INFO` may link but never supersedes / increments |
| Churn | `DefaultMaxAutoReplans=5`; `ErrReplanBudgetExceeded` when count ≥ N; `AckAutoReplan` / `AckReplan` resets 0 |
| S01 surface | `internal/planner` + mig **`006_plan_hierarchy.sql`** (`plan_phases`/`plan_scopes`/`scope_deep_plans`/`goal_plan_state`); `SupersedeDeepPlan` + deep-plan current+1 lookahead |
| S02 surface | mig **`007_discovery_severity.sql`**; store `IncrementAutoReplanCount`/`AckAutoReplan`; `planner.ApplyDiscoveryReplan` + `AckReplan` |
| Supporting unit (strong) | `./internal/planner/...` — `TestApplyDiscoveryReplanINFONoSupersede`, `…PlanAffectingSupersedes`, `…BlockingLikePlanAffecting`, `…BudgetAndAck` |
| Honesty | `evals/honesty` Paths A/B/C; no `AllowDoneWithoutReview` in proof |
| P0-X | `evals/p0x` 7/7 — keep green |
| X0 | `evals/x0` packages green (dry-run harness intact) |
| Gate C integrity | Recorded **Go** artifacts under `docs/verification/gate-c-x0/` remain `dry_run:false`, N=3; **do not invent new Go**; Mode-B packs historical |
| Dry-run ≠ Gate C | Phase 01 dry-run is regression-only — not Gate C pass |
| Fixture hash pin (carry) | `15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22` |
| Deferrals (carry) | GC-03/04 still deferred unless VERIFY Notes promote with evidence |
| Residual (non-blocking) | Global DPC attach on every task Expand (S02 REVIEW-NOTES); non-tx Apply / UNIQUE re-link / MCP no severity |
| Product Go | **Forbidden** — verify + evidence + handoff scaffold + spawn prompts only |
| MCP | Optional; Gate E is library demo (`evals/replan`) — do not require MCP for VERIFY pass |
| Daemon / HTTP / embeddings | Still forbidden as primary |
| Phase 04 folder name | **`phase-04-review-depth`** (A_PROJECT_PLAN Phase 4 — Review depth & evidence policies; refine only with Notes) |

### Locked verify commands

```bash
# Gate E mini-eval (primary — named planted demo)
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan

# Full Gate E package + planner/store/domain supporting surfaces
CGO_ENABLED=0 go test ./evals/replan/... ./internal/planner/... ./internal/store/... ./internal/domain/... -count=1

# Honesty (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1

# Full regression bar (includes p0x 7/7, x0, replan, honesty, analyzers)
CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./... -count=1
```

Optional (strong evidence, not substitutes for package PASS):

```bash
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -v -run TestHonestyFailClosedPlantedClaim
CGO_ENABLED=1 go test ./evals/p0x/... -count=1 -v -run TestP0XAllCriteria
CGO_ENABLED=0 go test ./internal/planner/... -count=1 -run 'TestApplyDiscoveryReplan'
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3, means match GATE-C-NOTES — do not re-score
# Fixture content hash (must still match pin unless Notes promote change)
find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# expect: 15fe50a18dd4c970202822591d9d0759a43d4227bdfadf5ec612270338b06d22
```

### Architecture / law checks (must all hold)

- [ ] No daemon / always-on HTTP as primary surface
- [ ] No committed `.trace/` under `fixtures/` or `evals/`
- [ ] Library packages still do not import `cmd/trace` or `cmd/trace-mcp` (G19)
- [ ] Gate E evidence is **`evals/replan` `TestPlantedDiscoveryReplan`** — not vibes / Notes-only
- [ ] Severity: INFO does not auto-replan; PLAN_AFFECTING+ does
- [ ] Churn N=5 fail-closed + ack proven (demo and/or unit)
- [ ] Gate C evidence remains `dry_run:false` — **not** Phase 01 dry-run alone; no silent Go flip
- [ ] Mode-B packs not falsified
- [ ] Embeddings still absent
- [ ] Mig 006 + 007 present and exercised by green planner/store tests
- [ ] GC-03/04 remain deferred unless explicitly promoted in VERIFY-NOTES

### DR-HANDOFF duties (this row + S03-02)

Per protocol Phase handoff + **DR-HANDOFF**. On Gate E + regression green → scaffold Phase 04. Do **not** leave Phase 04 as README-only / blocked-until-noticed.

| Who | Duty |
|-----|------|
| **P03-S03-01 (this VERIFY)** | On pass: **Create** Phase 04 scaffold under `docs/phases/phase-04-review-depth/` (checklist below) **and** append Phase 04 board section to `docs/TODO.md` with `P04-00` as first pending row after Phase 03’s last row. Record handoff progress in `VERIFY-NOTES.md`. |
| **P03-S03-02 (final review)** | **Owns completion check** — refuse `done` until scaffold is runnable (planner + stubs + board). May finish missing pieces. Marks Phase 03 complete only when handoff + VERIFY evidence agree. |

#### Phase 04 scaffold checklist (minimum)

Create / ensure:

```text
docs/phases/phase-04-review-depth/
  README.md                 # goal = Review depth & evidence policies (A_PROJECT_PLAN Phase 4)
  00-PHASE-PLANNER.md       # runnable (Agent→clarify→Plan→execute); light locks OK
  scopes/                   # ≥1 scope stub folder recommended; minimal OK
    scope-01-…/
      00-PLANNER.md
      01-*.md
      02-scope-review.md
      SCOPE-TODOS.md
```

Board (`docs/TODO.md`): new **Phase 04** section after Phase 03 rows; first pending ID **`P04-00`** → `phase-04-review-depth/00-PHASE-PLANNER.md`. Do **not** start executing Phase 04 until Phase 03 S03-02 is `done`.

Deep tasking stays with Phase 04’s own phase/scope planners. This row delivers a **runnable handoff**, not finished implement prompts.

**Counterfactual:** If Gate E / primary bars fail and cannot be remediated in-wave → record honest FAIL + spawn; do **not** scaffold Phase 04 until VERIFY finally greens (or user records `no successor`).

## Board rights
Verify: **status + notes** on `P03-S03-01`; **may spawn** remediation implement+review pairs **immediately below this row** if any gate fails; **may create upcoming Phase 04 scaffold** (DR-HANDOFF). Do **not** rewrite Phase 03 `done` history. Do **not** mark `P03-S03-02` or Phase 04 rows `done`.

## Preflight / Plan
1. Re-read this prompt + board row + S01/S02 REVIEW-NOTES + Gate E locks above.
2. Confirm module root has `go.mod` = `github.com/mrchatam/Trace`, packages `evals/{honesty,x0,p0x,replan}`, `internal/planner`, mig `006`/`007`, Gate C metrics under `docs/verification/gate-c-x0/`.
3. Plan: Gate E harness check → run locked commands → fill evidence → pass→VERIFY-NOTES + Phase 04 scaffold **or** fail→spawn → board update.

## Role work (VERIFY procedure)

### A. Gate E mini-eval re-check (required)

Confirm all of:

| Check | Expect |
|-------|--------|
| Harness path | `evals/replan/` exists; `doc.go` names planted demo |
| Named test | `TestPlantedDiscoveryReplan` PASS under locked command |
| Severity | PLAN_AFFECTING+ triggers supersede+link; INFO does not auto-replan (asserted in demo and/or planner units) |
| Churn | N=5 fail-closed + ack recovery proven |
| Surfaces | Consumes S01 `SupersedeDeepPlan` / coarse plan; S02 `ApplyDiscoveryReplan` — no second planner package |
| Honesty of claim | Notes map Gate E to this harness — **not** “vibes” or Gate C scores |

Record pass/fail in VERIFY-NOTES evidence table.

### B. Carry-forward bars (required)

| Check | Expect |
|-------|--------|
| Honesty H5 Paths A/B/C | `TestHonestyFailClosedPlantedClaim` PASS |
| P0-X 7/7 | `evals/p0x` PASS (`TestP0XAllCriteria` or package) |
| X0 packages | `./evals/x0/...` PASS |
| Gate C artifacts | `GATE-C-NOTES.md` still **Go**; metrics `dry_run:false`, N=3; means G1 0.800 > B0 0.000 — **re-check files, do not invent new Go** |
| Dry-run ≠ Gate C | Explicit in VERIFY-NOTES |
| S01/S02 packages | `./internal/planner/...` (+ store/domain as in locked cmds) PASS |
| Full suite | `CGO_ENABLED=1 go test ./... -count=1` PASS |

### C. Independent re-run (required)

From module root `/home/ali/Desktop/Trace`, run the **Locked verify commands**. All required commands must exit 0 for a green gate. Capture command + PASS/FAIL in `VERIFY-NOTES.md`.

### D. Evidence table (required in `VERIFY-NOTES.md`)

Copy and fill:

| Gate | Result (pass/fail) | Evidence (test / log gist) |
|------|--------------------|----------------------------|
| Gate E mini-eval | | `evals/replan` `TestPlantedDiscoveryReplan` |
| Severity PLAN_AFFECTING+ only | | demo +/or `TestApplyDiscoveryReplanINFONoSupersede` / PA / BLOCKING |
| Churn N=5 fail-closed + ack | | demo budget/ack and/or `TestApplyDiscoveryReplanBudgetAndAck` |
| S01 `internal/planner` + mig 006 | | package PASS; schema file present |
| S02 mig 007 + ApplyDiscoveryReplan | | package PASS; schema file present |
| Honesty H5 Paths A/B/C | | `TestHonestyFailClosedPlantedClaim` |
| P0-X 7/7 | | `TestP0XAllCriteria` (or package PASS + note) |
| X0 packages | | `./evals/x0/...` PASS |
| Gate C `dry_run:false` intact | | paths under `docs/verification/gate-c-x0/`; GATE-C-NOTES Go |
| Dry-run ≠ Gate C | | explicit: Phase 01 dry-run not used as Gate C / Gate E |
| `go test ./...` | | PASS/FAIL |
| Law checks | | no daemon/HTTP primary; no committed `.trace/`; G19 |
| Residuals (non-blocking) | | DPC-global; non-tx Apply; UNIQUE re-link; MCP no severity |
| DR-HANDOFF | | Phase 04 path created? board `P04-00` present? |

Also record: date, go/CGO note, residuals carried forward.

### E. On **all gates pass** + law checks

1. Write `VERIFY-NOTES.md` with verdict **Phase 03 VERIFY PASS / Gate E mini-eval green**, confidence, evidence table, residuals as **non-blocking** if primary paths green. Explicitly: Gate E = `evals/replan` planted demo; Gate C artifacts intact; Phase 01 dry-run ≠ Gate C.
2. **Start DR-HANDOFF:** create `docs/phases/phase-04-review-depth/` per checklist + append Phase 04 section / `P04-00` to `docs/TODO.md`.
3. Board Notes: short “Gate E + honesty/p0x/x0/replan/S01–S02 PASS; Gate C intact; Phase 04 scaffold started; see VERIFY-NOTES.md; pending P03-S03-02 handoff close”.
4. Mark `P03-S03-01` **done**. Do **not** mark S03-02 done.

### F. On **any fail**

1. `VERIFY-NOTES.md`: verdict **FAIL**, which gate/law, minimal reproduction.
2. Insert `P03-S03-01a` (implement) + `P03-S03-01b` (review) immediately below this row; set this row **`blocked`** (or `failed` + plan `01c` re-VERIFY) with reason. Prefer new `01c` verify row if this row already closed as failed/blocked (forward-only history).
3. Spawn prompts must be **full** protocol skeletons. Scope remediations to the failing layer — **do not** weaken bars, invent Gate E pass from Notes, or rewrite Gate C Mode-B packs.
4. **Forbidden “fixes”:** claiming Gate E without `TestPlantedDiscoveryReplan`; claiming Gate C from dry-run; requiring MCP for Gate E; adding daemon/HTTP; weakening honesty Paths A/B/C; rewriting `done` S01–S02 prompts; promoting GC-03/04 without evidence; scaffolding Phase 04 on a red VERIFY.

### Spawn ID convention

```text
… P03-S03-01  (this VERIFY)
… P03-S03-01a (remediation implement)   ← insert immediately below
… P03-S03-01b (remediation review)
… P03-S03-01c (re-VERIFY)               ← if original VERIFY closed as failed/blocked
… P03-S03-02  (phase review — only after VERIFY finally done; owns DR-HANDOFF completion)
```

Update `SCOPE-TODOS.md` when spawning.

## Out of scope
- Implementing Phase 04 review-depth product features
- Re-running live multi-model Gate C / flipping Go without contradicting evidence
- Rewriting Mode-B packs
- Promoting GC-03/04 without measurement need
- Declaring A1 / product thesis commercially validated
- Starting Phase 04 implement wave before S03-02 closes
- Expanding Gate E into a multi-model scored benchmark (mini-eval = planted automated demo only)

## Todo updates
Status + Notes on `P03-S03-01`; spawn rows if needed; Phase 04 upcoming artifacts on pass; `SCOPE-TODOS.md` checkboxes.

## Exit criteria
- [ ] Independent Gate E (`TestPlantedDiscoveryReplan` + severity + churn) + honesty + p0x + x0 + S01/S02 surfaces + Gate C artifact integrity + `./...` recorded in `VERIFY-NOTES.md`
- [ ] Evidence table includes Gate E path + dry-run≠Gate C + residuals + handoff
- [ ] Law checks recorded; Mode-B packs not falsified
- [ ] Either gates green **and** Phase 04 scaffold started (folder + `P04-00` board) **or** remediations spawned with full prompts + this row blocked/failed honestly
- [ ] TODO.md status + Notes updated; SCOPE-TODOS synced
- [ ] No product feature Go on this row

## Minimal todos
- [ ] Preflight: S01/S02 REVIEW-NOTES + mig 006/007 + Gate C metrics + Gate E harness path
- [ ] Re-check Gate E (`TestPlantedDiscoveryReplan` + severity + N=5/ack)
- [ ] Run honesty / p0x / x0 / planner+replan / Gate C artifact inspect / `./...` independently
- [ ] Write `VERIFY-NOTES.md` evidence table
- [ ] Pass → scaffold `phase-04-review-depth` + board `P04-00`; Fail → spawn 01a/01b (+01c plan)
- [ ] Board status + SCOPE-TODOS
