# P05 / S02 / 01 — Gate F prelim

## Metadata
- id: P05-S02-01
- todo_ids: [P05-S02-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Deliver the **Gate F preliminary** planted harness: a named automated package under **`evals/impact`** that (1) plants conflict / clean probes via live S01 impact APIs, (2) scores **precision/recall** against locked ground-truth assertions on `ImpactReport`, and (3) writes a schema-valid temp metrics artifact. Deterministic only — no LLM, no commercial multi-model Gate F. Do **not** invent or extend product impact policy beyond the harness.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) locks (finalized 2026-08-16)
- [phase README](../../README.md)
- [docs/init/PROJECT_MODEL_SNAPSHOT.md](../../../../init/PROJECT_MODEL_SNAPSHOT.md) Gate F
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 5
- S01 APPROVE: [../scope-01-impact-classes/REVIEW-NOTES.md](../scope-01-impact-classes/REVIEW-NOTES.md)
- Live S01: mig `009_decision_impact.sql`; `AddImpactFinding` / `ListImpactFindings`; `AddDecisionAlternative` / `SetRecommendedAlternative`; **`ImpactReport`** → `ImpactReportResult` (`OverallClass`, `HasUnknown`, `Incomplete`, Findings); `CreateDecision` / `LinkDecisionTask` (`decision_affects_task` only — no new rels)
- Prior planted patterns: `evals/honesty` Gate G (`schema-gate-g.json` + temp `metrics-gate-g.json`); `evals/replan` `TestPlantedDiscoveryReplan`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate — P05-S02-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Depends | S01 **done** (`P05-S01-02` APPROVE high) — fail closed if S01 surface missing |
| **Harness package** | **New `evals/impact`** — do **not** overload `evals/honesty` / `replan` / `x0` / `p0x` |
| **Named test** | **`TestPlantedImpactConflictsGateFPrelim`** |
| Metrics schema | Committed **`evals/impact/schema-gate-f.json`** (`schema_version` **const 1**) |
| Metrics artifact | Write **`metrics-gate-f.json`** under `t.TempDir()` (or subdir); validate against schema (jsonschema OK; same dep as `evals/x0` / Gate G) |
| S01 plant surface | Domain APIs only (G19) — **not** CLI-scrape; optional CLI is human-only |
| entity_links | **`decision_affects_task` only** — harness must not invent new rels |
| Mig | Consume **`009_decision_impact`** via `store.Open` — **no** S02 schema fork / new migration |
| Scoring rule | Score `HasUnknown` / `Incomplete` / Findings presence / `OverallClass` per probe GT; **never** trust `OverallClass` alone when GT requires `HasUnknown` (S01 residual: SAFE+UNKNOWN → OverallClass may be SAFE with flags true) |
| CGO | Harness is domain+store → must pass `CGO_ENABLED=0 go test ./evals/impact/...` |
| Carry-forward | Honesty A/B/C; Gate G `TestHonestyEscapeRateGateGPrelim`; Gate E `TestPlantedDiscoveryReplan`; p0x 7/7; x0; Gate C `docs/verification/gate-c-x0/` **untouched** (`dry_run:false` intact) |
| Out | Commercial multi-model Gate F; product Go outside `evals/impact` (+ schema); inventing impact APIs; `plan simulate`; daemon/HTTP/embeddings; rewriting Mode-B / Gate C packs; weakening Gate G/E/C/honesty/p0x |

### Precision / recall planting protocol (locked — implement exactly)

**Positive class** = planted conflict / incompleteness that **must** be surfaced by `ImpactReport`.  
**Negative class** = planted clean known impact that must **not** false-alarm `HasUnknown`.

Score each probe against its GT assertion set:

| Outcome | When |
|---------|------|
| **TP** | Positive probe: all GT assertions hold |
| **FN** | Positive probe: any GT assertion fails |
| **TN** | Negative probe: all GT assertions hold |
| **FP** | Negative probe: any GT assertion fails (false alarm / wrong class) |

**Planted fixture tallies (must assert in test + metrics):**

| Metric | Locked value |
|--------|----------------|
| `true_positives` | **3** |
| `false_negatives` | **0** |
| `false_positives` | **0** |
| `true_negatives` | **1** |
| `probes` | **4** |
| `precision` | `TP/(TP+FP)` → **1.0** |
| `recall` | `TP/(TP+FN)` → **1.0** |

**Part 0 — Open:** `store.Open(t.TempDir())` + `domain.New(st)`. Shared helpers OK. Separate decision per probe (isolation).

**Part 1 — Pos-1 — UNKNOWN uncertainty conflict**

```text
CreateDecision → CreateTask → LinkDecisionTask
AddImpactFinding(ImpactClass=CAUTION|SAFE, Uncertainty=UNKNOWN, Kind=UNRESOLVED|AFFECTED_WORK)
rep := ImpactReport(decisionID)
GT: HasUnknown==true && Incomplete==true
     && ≥1 finding with Uncertainty==UNKNOWN present in rep.Findings (never omitted)
# Do NOT assert OverallClass == UNKNOWN (OverallClass is max severity band, not uncertainty)
→ TP else FN
```

**Part 2 — Pos-2 — Multi-class severity conflict (rollup)**

```text
CreateDecision → CreateTask → LinkDecisionTask
AddImpactFinding(SAFE, KNOWN, …)
AddImpactFinding(DESTRUCTIVE, KNOWN, …)
rep := ImpactReport(decisionID)
GT: OverallClass == DESTRUCTIVE && HasUnknown == false
→ TP else FN
```

**Part 3 — Pos-3 — Linked tasks, empty findings (incomplete conflict)**

```text
CreateDecision → CreateTask → LinkDecisionTask
# no AddImpactFinding
rep := ImpactReport(decisionID)
GT: HasUnknown==true && Incomplete==true && OverallClass==""
→ TP else FN
```

**Part 4 — Neg-1 — Clean known SAFE (no false alarm)**

```text
CreateDecision → CreateTask → LinkDecisionTask
AddImpactFinding(SAFE, KNOWN, AFFECTED_WORK)  # no alternatives
rep := ImpactReport(decisionID)
GT: HasUnknown==false && Incomplete==false && OverallClass==SAFE
→ TN else FP
```

**Part 5 — Metrics write + schema validate**

Write `metrics-gate-f.json` with at least:

| Field | Value |
|-------|--------|
| `schema_version` | `1` |
| `gate` | `"F"` |
| `suite` | `"impact"` |
| `prelim` | `true` |
| `dry_run` | `false` |
| `true_positives` | `3` |
| `false_positives` | `0` |
| `false_negatives` | `0` |
| `true_negatives` | `1` |
| `precision` | `1.0` (or exact `TP/(TP+FP)`) |
| `recall` | `1.0` (or exact `TP/(TP+FN)`) |
| `probes` | `4` |
| `named_test` | `"TestPlantedImpactConflictsGateFPrelim"` |
| `mig` | `"009_decision_impact"` |
| `s01_hooks` | array including `"AddImpactFinding"`, `"LinkDecisionTask"`, `"ImpactReport"`, `"decision_affects_task"` |

Optional: `probe_ids` string array; `trace_version` string. Schema may allow `additionalProperties: true`.

Validate file against `evals/impact/schema-gate-f.json` before test returns.

### Target tree

```text
evals/impact/
  doc.go                      # NEW — package purpose + how to run named test
  impact_test.go              # NEW — TestPlantedImpactConflictsGateFPrelim (+ helpers)
  schema-gate-f.json          # NEW — schema_version 1
```

No new store migrations. No CLI required (G19: eval package must not import `cmd/trace`). No `internal/impact` package.

### How to run (Notes must cite)

```bash
CGO_ENABLED=0 go test ./evals/impact/... -count=1
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./... -count=1
```

Do **not** rewrite files under `docs/verification/gate-c-x0/`.

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] `TestPlantedImpactConflictsGateFPrelim` green; schema + temp `metrics-gate-f.json` validated
- [ ] Planted tallies: TP=3, FN=0, FP=0, TN=1; precision=1.0; recall=1.0
- [ ] S01 hooks exercised (`AddImpactFinding`, `LinkDecisionTask`, `ImpactReport` fields incl. HasUnknown/Incomplete)
- [ ] `CGO_ENABLED=0 go test ./evals/impact/...` PASS
- [ ] Carry-forward: honesty A/B/C + Gate G + Gate E + `CGO_ENABLED=1` p0x/x0/`./...` PASS; Gate C artifacts untouched
- [ ] Board Notes cite named test + schema path + metrics filename + tallies
- [ ] No product Go outside `evals/impact` (+ schema)

## Minimal todos
- [ ] Create `evals/impact` package (`doc.go`, `impact_test.go`, `schema-gate-f.json` v1)
- [ ] Implement `TestPlantedImpactConflictsGateFPrelim` (Pos-1..3 + Neg-1 + metrics write/validate)
- [ ] Run locked commands; board status + Notes only
