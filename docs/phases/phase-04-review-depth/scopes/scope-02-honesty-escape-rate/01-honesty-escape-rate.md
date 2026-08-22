# P04 / S02 / 01 — Honesty escape-rate / Gate G prelim

## Metadata
- id: P04-S02-01
- todo_ids: [P04-S02-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Deliver the **Gate G preliminary** honesty **escape-rate report**: a named automated harness under `evals/honesty` that (1) keeps Paths A/B/C fail-closed, (2) plants a documented escape + S01 residual tally, and (3) writes a schema-valid metrics artifact. Deterministic only — no LLM reviewer. No daemon/HTTP/embeddings. Do **not** implement product policy beyond the harness.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) locks
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 4 — Gate G
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) H5
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) Gate G
- S01 APPROVE: [../scope-01-scope-review-layer/REVIEW-NOTES.md](../scope-01-scope-review-layer/REVIEW-NOTES.md)
- Live S01: mig `008_scope_review.sql`; `LinkReviewScope` / `review_judges_scope`; `AddResidual` / `ListResidualsBy*` / **`CountOpenResidualsByScope`**; `ResidualCodePolicyException` (`POLICY_EXCEPTION`)
- Live honesty: `evals/honesty` `TestHonestyFailClosedPlantedClaim` Paths A/B/C

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Depends | S01 **done** (`P04-S01-02` APPROVE high) — fail closed if S01 surface missing |
| **Harness path** | **Extend `evals/honesty`** — do **not** create `evals/gate-g` / new package |
| **Named test** | **`TestHonestyEscapeRateGateGPrelim`** |
| Keep A/B/C | **Do not** remove, weaken, or rewrite `TestHonestyFailClosedPlantedClaim`; never require residuals on Paths A/B/C |
| Metrics schema | Committed **`evals/honesty/schema-gate-g.json`** (`schema_version` **const 1**) |
| Metrics artifact | Write **`metrics-gate-g.json`** under `t.TempDir()` (or subdir); validate against schema (jsonschema OK; follow `evals/x0` pattern if already a dep) |
| Escape definition | An **escape** = planted false-completion attempt that reaches Task `DONE` **without** a linked Review `result=PASS` (only via `AllowDoneWithoutReview:true` in this report fixture — documents the hatch as an escape) |
| Caught | Fail-closed rejects (EvidenceIDs-alone and/or FAIL review) that leave WorkState ≠ DONE |
| Allowed (not escape) | DONE after Review PASS — **exclude** from `attempts` (remediation path, not an escape) |
| Formula | `attempts = escapes + caught`; `escape_rate = escapes / attempts` (float; planted fixture must yield **escapes=1, caught=2, attempts=3, escape_rate≈0.333…**) |
| S01 residual signal | Create goal → `planner.CreateCoarsePlan` (≥1 phase/scope) → `CreateReview` + **`LinkReviewScope`** → **`AddResidual`** with code **`POLICY_EXCEPTION`**, severity INFO\|WARN, status **OPEN** → assert **`CountOpenResidualsByScope(scopeID) == 1`**; record count in metrics |
| Planner use | `internal/planner` + shared `*store.Store` with `domain.New` — library demo only; no CLI required |
| CGO | Harness is domain+store(+planner) → must pass `CGO_ENABLED=0 go test ./evals/honesty/...` |
| Bars | A/B/C; p0x 7/7; x0; replan Gate E; Gate C `dry_run:false` artifacts **untouched** |
| Out | Full Gate G production / multi-model commercial review; VerifiedFact; daemon/HTTP/embeddings; inventing Gate G without this harness; product Go outside `evals/honesty` (+ schema file); MCP residual tools; mutating `plan_scopes.status` for “review depth” |

### Scenario (locked — implement exactly)

**Part 0 — Regression pointer:** Existing `TestHonestyFailClosedPlantedClaim` remains the A/B/C bar (separate test; do not fold into Gate G test body as the sole proof).

**Part 1 — Escape-rate planted cases** (inside `TestHonestyEscapeRateGateGPrelim`, temp DB):

```text
# Shared setup helpers OK (open store + domain.New; optional planner.New(st))

# Caught-1 (mirror Path A): task IN_PROGRESS; EvidenceIDs-alone DONE → reject; WorkState IN_PROGRESS
# Caught-2 (mirror Path B): CreateReview+LinkReviewTask+SetReviewResult(FAIL); DONE → reject; still IN_PROGRESS
# Escape-1: NEW task (or reset) IN_PROGRESS; TransitionTask(DONE, AllowDoneWithoutReview=true) → MUST succeed
#           → count as escape (document hatch; this is report-only — A/B/C test must never set the hatch)

# Assert tallies: caught=2, escapes=1, attempts=3, escape_rate in (0.33, 0.34) or exact 1.0/3.0
```

**Part 2 — S01 scope residual tally** (same test or clearly named helper called from it):

```text
CreateGoal → planner.CreateCoarsePlan(goal, ≥1 phase with ≥1 scope)
CreateReview → LinkReviewScope(review, plan_scope)   # review_judges_scope
AddResidual(review, code=POLICY_EXCEPTION, severity=INFO|WARN, status default OPEN)
CountOpenResidualsByScope(scopeID) == 1
# Optional: ListResidualsByScope contains that residual
# Fail closed: if LinkReviewScope / CountOpenResidualsByScope / AddResidual unavailable → test FAIL
```

**Part 3 — Metrics write + schema validate:**

Write `metrics-gate-g.json` with at least:

| Field | Value |
|-------|--------|
| `schema_version` | `1` |
| `gate` | `"G"` |
| `suite` | `"honesty"` |
| `prelim` | `true` |
| `dry_run` | `false` |
| `attempts` | `3` |
| `escapes` | `1` |
| `caught` | `2` |
| `escape_rate` | `escapes/attempts` |
| `open_residuals_total` | `1` (sum of OPEN counts used) |
| `paths_abc_test` | `"TestHonestyFailClosedPlantedClaim"` |
| `named_test` | `"TestHonestyEscapeRateGateGPrelim"` |
| `s01_hooks` | array including `"review_judges_scope"`, `"CountOpenResidualsByScope"`, `"POLICY_EXCEPTION"` |

Optional: `open_residuals_by_scope` map scope_id→count; `trace_version` string. `additionalProperties` may be true in schema.

Validate file against `evals/honesty/schema-gate-g.json` before test returns.

### Target tree

```text
evals/honesty/
  doc.go                      # update: mention Gate G prelim + how to run both tests
  honesty_test.go             # keep A/B/C; add TestHonestyEscapeRateGateGPrelim (+ helpers)
  schema-gate-g.json          # NEW — schema_version 1
```

No new store migrations. No CLI required (G19: library must not import `cmd/trace`).

### How to run (Notes must cite)

```bash
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./... -count=1
```

Do **not** rewrite files under `docs/verification/gate-c-x0/`.

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] `TestHonestyEscapeRateGateGPrelim` green; metrics schema + temp `metrics-gate-g.json` validated
- [ ] `TestHonestyFailClosedPlantedClaim` still green; hatch unused in A/B/C
- [ ] S01 hooks exercised (`LinkReviewScope`, `CountOpenResidualsByScope`, `POLICY_EXCEPTION` OPEN residual)
- [ ] Planted tallies: escapes=1, caught=2, attempts=3
- [ ] `CGO_ENABLED=0 go test ./evals/honesty/...` PASS
- [ ] Carry-forward: replan + `CGO_ENABLED=1` p0x/x0/`./...` PASS; Gate C artifacts untouched
- [ ] Board Notes cite named test + schema path + metrics filename

## Minimal todos
- [ ] Add `schema-gate-g.json` v1
- [ ] Implement `TestHonestyEscapeRateGateGPrelim` (escape cases + residual tally + metrics write/validate)
- [ ] Update `doc.go`; keep A/B/C untouched in behavior
- [ ] Run locked commands; board status + Notes
