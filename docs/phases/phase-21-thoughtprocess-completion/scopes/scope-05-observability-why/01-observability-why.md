# P21-S05-01 — Implement: observability + why

## Metadata
- id: P21-S05-01
- todo_ids: [P21-S05-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective
Implement S05-00: CLI why for P20 types (verify + impact overlay), `deliberation.transition` in task why chains, and `historical_relationships` loop next section. **No new migration.** Compat ceiling **20** unchanged.

## Session start
Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: **status + notes only**.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- [DECISION-LOG.md](../../DECISION-LOG.md) D-11, D-12
- [WORK-MAP.md](../../WORK-MAP.md) W-07, W-08
- P20 relationships: [scope-05/00-PLANNER.md](../../../phase-20-cognitive-deliberation/scopes/scope-05-regression-reflect/00-PLANNER.md)
- Live: `cmd/trace/why.go`, `internal/retrieval/{why.go,types.go}`, `internal/domain/impact.go`, `internal/loop/next.go`, `internal/store/links.go`

## Prerequisites
**S02 done** — P20 types registered in `lookupEntity` / `NormalizeEntityType`.

## Locked defaults (from S05-00 — do not re-debate)

| Item | Value |
|------|-------|
| Why P20 types | Same 8 as S02 — no new entity types |
| Transition Why step | `reason_code: deliberation_transition`; `title` = `to_phase`; `detail` = payload `reason_code` |
| Historical section | Top-level `historical_relationships`; max **8** items; repo-wide links |
| Link rels | `observed_relationship` + evidence-backed `caused_by` only |
| caused_by filter | Require `relationship_supported_by` from same from-entity |
| Loop schema | `trace.loop.next.v1` string **unchanged** (additive JSON field only) |
| Schema / compat | Max mig **020**; ceiling **20** — no **021+** |
| MCP | No new tools |

## Live inventory (before — confirmed S05-00)

| Surface | Location | Today |
|---------|----------|-------|
| CLI why | `cmd/trace/why.go` | Why + impact JSON; P20 lookup works via S02 |
| Why events | `retrieval/why.go` L55–74 | Generic `recent_event` for all event types |
| Impact | `domain/impact.go` L375–391 | `task` + `decision` only; P20 → nil |
| Next packet | `loop/next.go` L47–61 | No `historical_relationships` field |
| Observed/causal API | `domain/regressions.go` L747–835 | Library exists; not surfaced in loop next |
| Schema | 20 files, max **020** | S04 promotion mig landed |

## Requirements

### 1. Deliberation transition in task Why (W-07)

1. Add `ReasonDeliberationTransition = "deliberation_transition"` to `internal/retrieval/types.go`.
2. In `retrieval/why.go` event loop (task seed events):
   - When `ev.Type == "deliberation.transition"`, unmarshal `deliberation.TransitionPayload`.
   - Emit `WhyStep` per FINAL shape in 00-PLANNER (title=`to_phase`, detail=`reason_code`).
   - On parse error, fall back to existing `ReasonRecentEvent` behavior.
3. **Do not** change event ordering policy (still last ≤8 events on seed entity, appended after graph expand).

### 2. Impact overlay for P20 why seeds (W-07)

Extend `ImpactSummariesForWhySeed` in `internal/domain/impact.go`:

1. **`regression`**: Follow `ListLinksFrom(regression_id)` for `regression_from_evaluation` or `regression_from_effect`. Return one lightweight `DecisionImpact` summarizing source (kind + id + regression dimension/summary). Return `nil` when no source link.
2. **`uncertainty`**: Return `nil, nil` — graph expand covers `uncertainty_blocks_task` neighbors (S02).
3. **Other P20 types**: Return `nil, nil` (no impact array in CLI JSON).
4. **Keep** existing `task` + `decision` behavior unchanged.

`cmd/trace/why.go` should need **no** logic changes (already calls `ImpactSummariesForWhySeed`).

### 3. Historical relationships loop section (W-08)

1. Add types to `internal/loop/next.go` (or `deliberation_packet.go` if helpers live there):

```go
const maxHistoricalRelationshipsCap = 8

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

2. Implement `buildHistoricalRelationshipsSection(st *store.Store, freshness string) (HistoricalRelationshipsSection, error)`:
   - Load `ListLinksByRel("observed_relationship")`.
   - Load `ListLinksByRel("caused_by")`; for each, verify ≥1 `relationship_supported_by` link from same `(from_type, from_id)` via `ListLinksFrom`.
   - Merge, sort `created_at DESC, id DESC`, cap at **8**.
   - Map to items (omit link id / evidence ids).
3. Add `HistoricalRelationships HistoricalRelationshipsSection` field on `NextPacket`.
4. Wire in `BuildNextPacket` — use same freshness as seed task section.
5. Empty DB → `{freshness, items:[]}` — valid packet, exit 0.

### 4. Tests (6 named + keepers)

| Test | Assert |
|------|--------|
| `TestCLIWhyUncertainty` | init + create uncertainty via domain or apply; `trace why uncertainty <id>` exit 0; JSON `steps[0]` matches seed |
| `TestCLIWhyRegression` | regression with `regression_from_evaluation` or `regression_from_effect` link; why steps include `outcome_result` or `effect` entity |
| `TestWhyTaskIncludesDeliberationTransition` | loop apply triggers transition; `trace why task <id>` has step with `reason_code==deliberation_transition` and detail matching transition reason |
| `TestLoopNextHistoricalRelationshipsSection` | insert observed_relationship; loop next JSON has `historical_relationships.items` length ≥1, ≤8 |
| `TestHistoricalRelationshipsObservedVsCaused` | caused_by without evidence → absent from packet; with evidence → present |
| `TestCausalWhyContextRoundTrip` | P17 keeper — must stay green |

Also keep green: `TestWhyUncertaintyIncludesGraphSteps`, `TestSeedImportAndWhy`, `TestLoopApplyDeliberationTransitionEvent`.

## Touch files

- `internal/retrieval/types.go`
- `internal/retrieval/why.go`
- `internal/domain/impact.go`
- `internal/loop/next.go`
- `internal/loop/deliberation_packet.go` (helper + cap const if placed here)
- `internal/store/links.go` (optional list helper)
- `cmd/trace/cli_test.go`
- `cmd/trace/loop_test.go`

## Keeper floor

```bash
go test ./internal/retrieval/... -count=1 -run 'TestWhyUncertaintyIncludesGraphSteps'
go test ./cmd/trace -count=1 -run 'TestCLIWhyUncertainty|TestCLIWhyRegression|TestWhyTaskIncludesDeliberationTransition|TestLoopNextHistoricalRelationshipsSection|TestHistoricalRelationshipsObservedVsCaused|TestCausalWhyContextRoundTrip|TestSeedImportAndWhy|TestLoopApplyDeliberationTransitionEvent'
go test ./internal/domain/... -count=1 -run 'TestCausalRelationshipFailClosedWithoutEvidence'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [ ] 6 named tests PASS
- [ ] `trace.loop.next.v1` schema version string unchanged
- [ ] No mig **021+**; compat ceiling **20**
- [ ] D-11 + D-12 behaviors evidenced in tests

## Next

**P21-S05-02**
