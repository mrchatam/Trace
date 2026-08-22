# P06 / S02 / 01 — Capability selection

## Metadata
- id: P06-S02-01
- todo_ids: [P06-S02-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Deliver the **capability-selection ablation** planted harness: a named automated package under **`evals/capability`** that (1) plants selection / missing probes via live S01 capability APIs + compiler packet attach, (2) scores **precision/recall** against locked ground-truth assertions, and (3) writes a schema-valid temp metrics artifact. Deterministic only — no LLM, no commercial multi-model theater. Do **not** invent product selection policy beyond the harness; S01 already attaches required+missing to packets.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) locks (finalized 2026-08-16)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 6
- [docs/EVALUATION.md](../../../../EVALUATION.md) H7
- S01 APPROVE: [../scope-01-capability-surface/REVIEW-NOTES.md](../scope-01-capability-surface/REVIEW-NOTES.md)
- Live S01: mig `010_capability_surface.sql`; `UpsertCapability` / `RequireCapability` / `MissingCapabilities` / `ListRequiredCapabilities`; compiler `TaskContext` → packet `required_capabilities`+`missing_capabilities`; optional `BuiltinMCPCapabilitySpecs`
- Prior planted patterns: `evals/impact` Gate F (`schema-gate-f.json` + temp `metrics-gate-f.json`); `evals/honesty` Gate G; `evals/replan` Gate E

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate — P06-S02-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Depends | S01 **done** (`P06-S01-02` APPROVE high) — fail closed if S01 surface missing |
| **Harness package** | **New `evals/capability`** — do **not** overload `evals/honesty` / `replan` / `impact` / `x0` / `p0x` |
| **Named test** | **`TestPlantedCapabilitySelectionAblation`** |
| Metrics schema | Committed **`evals/capability/schema-capability.json`** (`schema_version` **const 1**) |
| Metrics artifact | Write **`metrics-capability.json`** under `t.TempDir()` (or subdir); validate against schema (jsonschema OK; same dep as `evals/x0` / Gate F / Gate G) |
| S01 plant surface | Domain + compiler library APIs only (G19) — **not** CLI-scrape; optional CLI is human-only |
| Mig | Consume **`010_capability_surface`** via `store.Open` — **no** S02 schema fork / new migration |
| Scoring rule | Score `MissingCapabilities` + packet `required_capabilities` / `missing_capabilities` per probe GT; packet must list **only** required caps (no catalog dump) |
| CGO | Harness is domain+store+compiler → must pass `CGO_ENABLED=0 go test ./evals/capability/...` |
| Carry-forward | Honesty A/B/C; Gate G `TestHonestyEscapeRateGateGPrelim`; Gate E `TestPlantedDiscoveryReplan`; Gate F `TestPlantedImpactConflictsGateFPrelim`; p0x 7/7; x0; Gate C `docs/verification/gate-c-x0/` **untouched** (`dry_run:false` intact) |
| Out | Commercial multi-model theater; product Go outside `evals/capability` (+ schema); inventing capability APIs / `internal/capability`; daemon/HTTP/embeddings; rewriting Mode-B / Gate C packs; weakening Gate F/G/E/C/honesty/p0x |

### Precision / recall planting protocol (locked — implement exactly)

**Positive class** = planted selection / missing behavior that **must** hold.  
**Negative class** = planted clean AVAILABLE case that must **not** false-alarm missing.

Score each probe against its GT assertion set:

| Outcome | When |
|---------|------|
| **TP** | Positive probe: all GT assertions hold |
| **FN** | Positive probe: any GT assertion fails |
| **TN** | Negative probe: all GT assertions hold |
| **FP** | Negative probe: any GT assertion fails (false alarm / wrong selection) |

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

**Part 0 — Open:** `store.Open(t.TempDir())` + `domain.New(st)` + compiler on same store. Shared helpers OK. Separate task per probe (isolation).

**Part 1 — Pos-1 — UNAVAILABLE missing warning**

```text
CreateTask
UpsertCapability(slug=tool:down, status=UNAVAILABLE, kind=TOOL)
RequireCapability(task, capID)
miss := MissingCapabilities(taskID)
pkt := compiler.TaskContext(taskID, …)
GT: miss contains tool:down (by slug or id)
     && pkt.missing_capabilities contains tool:down
     && pkt.required_capabilities contains tool:down
→ TP else FN
```

**Part 2 — Pos-2 — UNKNOWN missing warning**

```text
CreateTask
UpsertCapability(slug=skill:maybe, status=UNKNOWN, kind=SKILL)
RequireCapability(task, capID)
miss := MissingCapabilities(taskID)
pkt := compiler.TaskContext(taskID, …)
GT: miss contains skill:maybe
     && pkt.missing_capabilities contains skill:maybe
     && pkt.required_capabilities contains skill:maybe
→ TP else FN
```

**Part 3 — Pos-3 — Selection filter (no catalog dump)**

```text
CreateTask
UpsertCapability × N≥3 all AVAILABLE (distinct slugs; optional: seed BuiltinMCPCapabilitySpecs via Upsert)
RequireCapability exactly 2 of them
miss := MissingCapabilities(taskID)
pkt := compiler.TaskContext(taskID, …)
GT: len(pkt.required_capabilities)==2 && only those two slugs
     && len(pkt.missing_capabilities)==0 && len(miss)==0
     && no non-required catalog slug appears in pkt.required_capabilities
→ TP else FN
```

**Part 4 — Neg-1 — Clean AVAILABLE (no false missing)**

```text
CreateTask
UpsertCapability(slug=tool:ok, status=AVAILABLE, kind=TOOL)
RequireCapability(task, capID)
miss := MissingCapabilities(taskID)
pkt := compiler.TaskContext(taskID, …)
GT: len(miss)==0 && len(pkt.missing_capabilities)==0
     && pkt.required_capabilities == required set (exactly tool:ok)
→ TN else FP
```

**Part 5 — Metrics write + schema validate**

Write `metrics-capability.json` with at least:

| Field | Value |
|-------|--------|
| `schema_version` | `1` |
| `gate` | `"capability-selection"` |
| `suite` | `"capability"` |
| `ablation` | `true` |
| `dry_run` | `false` |
| `true_positives` | `3` |
| `false_positives` | `0` |
| `false_negatives` | `0` |
| `true_negatives` | `1` |
| `precision` | `1.0` (or exact `TP/(TP+FP)`) |
| `recall` | `1.0` (or exact `TP/(TP+FN)`) |
| `probes` | `4` |
| `named_test` | `"TestPlantedCapabilitySelectionAblation"` |
| `mig` | `"010_capability_surface"` |
| `s01_hooks` | array including `"UpsertCapability"`, `"RequireCapability"`, `"MissingCapabilities"`, `"required_capabilities"`, `"missing_capabilities"` |

Optional: `probe_ids` string array; `trace_version` string. Schema may allow `additionalProperties: true`.

Validate file against `evals/capability/schema-capability.json` before test returns.

### Target tree

```text
evals/capability/
  doc.go                       # NEW — package purpose + how to run named test
  capability_test.go           # NEW — TestPlantedCapabilitySelectionAblation (+ helpers)
  schema-capability.json       # NEW — schema_version 1
```

No new store migrations. No CLI required (G19: eval package must not import `cmd/trace`). No `internal/capability` package.

### How to run (Notes must cite)

```bash
CGO_ENABLED=0 go test ./evals/capability/... -count=1
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./... -count=1
```

Do **not** rewrite files under `docs/verification/gate-c-x0/`.

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] `TestPlantedCapabilitySelectionAblation` green; schema + temp `metrics-capability.json` validated
- [ ] Planted tallies: TP=3, FN=0, FP=0, TN=1; precision=1.0; recall=1.0
- [ ] S01 hooks exercised (`UpsertCapability`, `RequireCapability`, `MissingCapabilities`, packet required+missing)
- [ ] `CGO_ENABLED=0 go test ./evals/capability/...` PASS
- [ ] Carry-forward: honesty A/B/C + Gate G + Gate E + Gate F + `CGO_ENABLED=1` p0x/x0/`./...` PASS; Gate C artifacts untouched
- [ ] Board Notes cite named test + schema path + metrics filename + tallies
- [ ] No product Go outside `evals/capability` (+ schema)

## Minimal todos
- [ ] Create `evals/capability` package (`doc.go`, `capability_test.go`, `schema-capability.json` v1)
- [ ] Implement `TestPlantedCapabilitySelectionAblation` (Pos-1..3 + Neg-1 + metrics write/validate)
- [ ] Run locked commands; board status + Notes only

## Out of scope
- Phase VERIFY (S03); daemon/HTTP/embeddings; rewriting Gate C Mode-B packs; product selection API beyond S01
