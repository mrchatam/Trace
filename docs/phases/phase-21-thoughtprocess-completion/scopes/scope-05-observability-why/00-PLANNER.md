# P21-S05-00 — Planner: observability + trace why for P20

## Metadata
- id: P21-S05-00
- todo_ids: [P21-S05-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock `trace why` for P20 entity types, `deliberation.transition` in task why chains, and `historical_relationships` loop next section (§17). **No product Go this row.**

## References
- [DECISION-LOG.md](../../DECISION-LOG.md) D-11, D-12
- [WORK-MAP.md](../../WORK-MAP.md) W-07, W-08
- P20 relationships: [scope-05/00-PLANNER.md](../../../phase-20-cognitive-deliberation/scopes/scope-05-regression-reflect/00-PLANNER.md) — `observed_relationship` / `caused_by` / evidence policy
- Live: `cmd/trace/why.go`, `internal/retrieval/why.go`, `internal/loop/next.go`, `internal/domain/{impact.go,links.go,regressions.go}`, `internal/store/links.go`

## Live inventory (confirmed 2026-08-18)

| Surface | Location | Today (live read) | S05 action |
|---------|----------|-------------------|------------|
| `trace why` CLI | `cmd/trace/why.go` L14–54 | Calls `eng.Why` + `ImpactSummariesForWhySeed`; JSON `{steps, impact?}` | **Verify** P20 types work post-S02; **extend** impact for `uncertainty` / `regression` seeds |
| `lookupEntity` P20 | `retrieval/exact.go` L208–314 | All 8 P20 types registered (S02 landed) | **Reuse** — no new types |
| `Why` graph expand | `retrieval/why.go` L36–53 | Expand from seed; all link rels including `regression_from_*` | **Keep** — regression why reaches source eval/effect via Expand |
| Task events in Why | `retrieval/why.go` L55–74 | Last 8 events; `ReasonRecentEvent`; Detail = payload excerpt | **Special-case** `deliberation.transition` → `ReasonDeliberationTransition` + `reason_code` in Detail |
| `ImpactSummariesForWhySeed` | `domain/impact.go` L373–391 | Only `task` + `decision` cases; P20 → `nil` | **Add** `uncertainty` (optional task-neighborhood) + `regression` (source eval/effect summary) |
| `NextPacket` | `loop/next.go` L47–61 | 11 sections; **no** `historical_relationships` | **Add** top-level `historical_relationships` section |
| `entity_links` rels | `domain/service.go` L58–60 | `observed_relationship`, `caused_by`, `relationship_supported_by` | **Reuse** — no Relationship table, no mig |
| `RecordObservedRelationship` | `regressions.go` L747–777 | No evidence required | Input to historical section |
| `RecordCausalRelationship` | `regressions.go` L779–835 | Requires evidence IDs → `relationship_supported_by` links | Packet **includes** caused_by only when evidence link exists |
| Schema / compat | `internal/store/schema/` | **20** files, max **020** (S04 landed); compat ceiling **20** | **No new migration** in S05; compat **20** unchanged |
| S02 keeper | `retrieval/p20_test.go` | `TestWhyUncertaintyIncludesGraphSteps` at retrieval layer | **Keep** green; S05 adds CLI-level tests |
| P17 keeper | `cmd/trace/cli_test.go` L115 | `TestCausalWhyContextRoundTrip` | **Keep** green |

## W-07 / W-08 locks

| Work ID | TRACE § | Lock |
|---------|---------|------|
| **W-07** | §25 Observability | `trace why <p20-type> <id>` exits 0 with ordered `steps[]` + optional `impact[]`; task why surfaces `deliberation.transition` events with transition `reason_code` inspectable |
| **W-08** | §17 Historical knowledge | Loop `next` packet adds bounded `historical_relationships` from `entity_links` (`observed_relationship` + evidence-backed `caused_by` only); max **8** rows; **no** new tables |

## FINAL locked defaults (S05-01 must not re-debate)

| Item | Value |
|------|-------|
| Why entity types | Same canonical **8** as S02: `uncertainty`, `hypothesis`, `change`, `effect`, `regression`, `reflection`, `baseline`, `outcome_result` |
| CLI aliases | Reuse S02 `NormalizeEntityType` (`outcome`→`outcome_result`, `plan-change`→`plan_change`) |
| Transition events | Event type `deliberation.transition` on task seed → Why step with `reason_code: deliberation_transition`; `title` = `to_phase`; `detail` = transition `reason_code` string from payload |
| New reason constant | `ReasonDeliberationTransition = "deliberation_transition"` in `retrieval/types.go` |
| Impact overlay — uncertainty | Return `nil` impact **or** lightweight task-neighborhood summary; graph expand is sufficient for steps — **must not** fail CLI |
| Impact overlay — regression | When `regression_from_evaluation` / `regression_from_effect` link exists, include one `DecisionImpact`-shaped summary referencing source kind + id (optional `impact[]` entry) |
| `historical_relationships` | Top-level `NextPacket` field; schema string **`trace.loop.next.v1` unchanged** |
| Section shape | `{freshness, items[]}` — same freshness pattern as `recent_changes` (seed task freshness) |
| Item fields | `rel`, `from_type`, `from_id`, `to_type`, `to_id`, `confidence` — **omit** link id / evidence ids in packet |
| Link rels | **`observed_relationship`** + **`caused_by`** only — no `relationships` table |
| Selection | Merge `ListLinksByRel(observed_relationship)` + filtered `ListLinksByRel(caused_by)`; sort `created_at DESC, id DESC`; cap **8** |
| caused_by filter | Include only when ≥1 `relationship_supported_by` link from same `(from_type, from_id)` to evidence (matches `RecordCausalRelationship` invariant) |
| Scope | **Repo-wide** historical edges (§17 cross-cutting knowledge) — not task-filtered |
| Empty section | Valid packet when no links: `{freshness, items:[]}` |
| Schema / compat | Max mig **020**; compat ceiling **20** — **no mig 021** in S05 |
| MCP | No new tools |
| Seed | No seed changes |

### Why transition step shape (FINAL)

When `ev.Type == domain.EventDeliberationTransition` (`"deliberation.transition"`):

```go
WhyStep{
    EntityType: ev.EntityType,  // task
    EntityID:   ev.EntityID,
    Title:      string(payload.ToPhase),
    ReasonCode: ReasonDeliberationTransition,
    Detail:     string(payload.ReasonCode),
    Distance:   0,
}
```

Parse failures on transition payload → fall back to existing `ReasonRecentEvent` + excerpt (do not fail Why).

### Historical relationships structs (FINAL)

```go
type HistoricalRelationshipsSection struct {
    Freshness string                       `json:"freshness"`
    Items     []HistoricalRelationshipItem `json:"items"`
}

type HistoricalRelationshipItem struct {
    Rel        string  `json:"rel"`
    FromType   string  `json:"from_type"`
    FromID     string  `json:"from_id"`
    ToType     string  `json:"to_type"`
    ToID       string  `json:"to_id"`
    Confidence float64 `json:"confidence"`
}
```

Wire in `BuildNextPacket` after `recent_changes` assembly; helper `buildHistoricalRelationshipsSection(st *store.Store, freshness string)`.

### Impact extension sketch (FINAL)

```go
// impact.go — extend ImpactSummariesForWhySeed switch:
case EntityRegression:
    // ListLinksFrom(regression_id); find regression_from_evaluation | regression_from_effect
    // return []DecisionImpact{{Kind: sourceKind, EntityID: sourceID, Summary: ...}} or nil
case EntityUncertainty:
    return nil, nil // graph expand covers uncertainty_blocks_task neighbors
```

## Named tests (S05-01)

| Test | Location | Proves |
|------|----------|--------|
| `TestCLIWhyUncertainty` | `cmd/trace/cli_test.go` | `trace why uncertainty <id>` exit 0; JSON has `steps` with seed step |
| `TestCLIWhyRegression` | `cmd/trace/cli_test.go` | `trace why regression <id>` exit 0; steps include linked `outcome_result` or `effect` via Expand |
| `TestWhyTaskIncludesDeliberationTransition` | `cmd/trace/loop_test.go` or `cli_test.go` | After `loop apply` hop: task why lists step with `reason_code` `deliberation_transition` and transition `reason_code` in detail |
| `TestLoopNextHistoricalRelationshipsSection` | `cmd/trace/loop_test.go` | `loop next` JSON has `historical_relationships.items` (≤8); observed link present |
| `TestHistoricalRelationshipsObservedVsCaused` | `internal/domain/regressions_test.go` or `loop_test.go` | `caused_by` without evidence excluded from packet; with evidence included |
| `TestCausalWhyContextRoundTrip` | `cmd/trace/cli_test.go` | P17 keeper — unchanged behavior |

**Keep** S02 `TestWhyUncertaintyIncludesGraphSteps`, `TestSeedImportAndWhy`, `TestLoopApplyDeliberationTransitionEvent`.

## Touch files

- `internal/retrieval/types.go` — `ReasonDeliberationTransition`
- `internal/retrieval/why.go` — transition event parsing + ordering
- `internal/domain/impact.go` — `ImpactSummariesForWhySeed` regression (± uncertainty) cases
- `internal/loop/next.go` — `HistoricalRelationshipsSection` + wire in `BuildNextPacket`
- `internal/loop/deliberation_packet.go` — optional: `buildHistoricalRelationshipsSection` helper + cap const
- `internal/store/links.go` — optional: `ListHistoricalRelationshipLinks(limit int)` if cleaner than inline merge
- `cmd/trace/why.go` — no logic change expected (verify only)
- `cmd/trace/cli_test.go` — `TestCLIWhyUncertainty`, `TestCLIWhyRegression`
- `cmd/trace/loop_test.go` — `TestWhyTaskIncludesDeliberationTransition`, `TestLoopNextHistoricalRelationshipsSection`, `TestHistoricalRelationshipsObservedVsCaused`

## Planner work

1. [x] Live inventory `why.go` / `impact.go` / `next.go` / links / schema max **020** / compat **20**.
2. [x] Lock **W-07** (trace why P20 + deliberation audit in task why).
3. [x] Lock **W-08** (`historical_relationships` loop section — §17 Should).
4. [x] Lock transition step shape + caused_by evidence filter + repo-wide cap **8**.
5. [x] Thicken `01-observability-why.md` + `02-scope-review.md` with before-state, 6 tests, keeper floor.
6. [x] Update `SCOPE-TODOS.md`.

## Exit criteria

- [x] 6 named tests locked
- [x] W-07/W-08 locks explicit
- [x] 01/02 thickened enough to implement alone
- [x] No product Go
- [x] No mig 021

## Next

**P21-S05-01**
