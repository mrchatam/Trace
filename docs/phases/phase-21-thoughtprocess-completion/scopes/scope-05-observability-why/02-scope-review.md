# P21-S05-02 — Review: observability + why

## Metadata
- id: P21-S05-02
- todo_ids: [P21-S05-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective
Independent review: why surfaces for P20 types, deliberation transition audit in task why, and `historical_relationships` loop section match S05-00 locks. Law 5/11 attribution unchanged. **No mig 021.**

## Session start
**Fresh subagent** (not S05-01 session). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: status + notes only; spawn Na/Nb if gap found.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-observability-why.md](01-observability-why.md) — implementer deliverable
- [DECISION-LOG.md](../../DECISION-LOG.md) D-11, D-12
- [WORK-MAP.md](../../WORK-MAP.md) W-07, W-08

## Review checklist

| Check | Evidence |
|-------|----------|
| P20 why CLI | `trace why uncertainty|regression <id>` exit 0 + JSON steps |
| Why types | Same 8 canonical P20 types as S02 — no duplicate registration |
| Transition step | Task why includes `deliberation.transition` with `reason_code: deliberation_transition` and payload `reason_code` in detail |
| Parse fallback | Malformed transition payload falls back to `recent_event` — Why does not fail |
| Impact — regression | Optional `impact[]` when source eval/effect link exists |
| Impact — uncertainty | Nil impact OK; graph expand provides neighbors |
| Historical section | `loop next` JSON has top-level `historical_relationships` |
| Item shape | Fields: `rel`, `from_type`, `from_id`, `to_type`, `to_id`, `confidence` only |
| Cap | ≤ **8** items |
| observed vs caused | `observed_relationship` always eligible; `caused_by` only with evidence links |
| caused_by ≠ attribution | Causal link does **not** auto-set `regressions.attribution=caused` (Law 5) |
| No Relationship table | Links from `entity_links` only — no new table |
| Loop schema string | `schema_version` still `trace.loop.next.v1` |
| No mig 021 | Schema dir has **20** files max; compat ceiling **20** |
| 6 named tests | All exist + PASS |
| P17/S02 keepers | `TestCausalWhyContextRoundTrip`, `TestWhyUncertaintyIncludesGraphSteps`, `TestSeedImportAndWhy` green |
| MCP | No new tools |

## D-11 / D-12 closure

- **D-11 promote:** `trace why` works for P20 nouns at CLI + retrieval; regression why reaches source eval/effect.
- **D-12 promote:** `historical_relationships` surfaced in loop next packet (§17 Should).

## Keeper command floor

```bash
go test ./internal/retrieval/... -count=1 -run 'TestWhyUncertaintyIncludesGraphSteps'
go test ./cmd/trace -count=1 -run 'TestCLIWhyUncertainty|TestCLIWhyRegression|TestWhyTaskIncludesDeliberationTransition|TestLoopNextHistoricalRelationshipsSection|TestHistoricalRelationshipsObservedVsCaused|TestCausalWhyContextRoundTrip|TestSeedImportAndWhy|TestLoopApplyDeliberationTransitionEvent'
go test ./internal/domain/... -count=1 -run 'TestCausalRelationshipFailClosedWithoutEvidence|TestRegressionAttribution'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Review focus

- **Transition inspectability:** Agent can answer "why INVESTIGATE?" from task why steps without reading raw events table.
- **Evidence gate:** Packet must not leak unproven `caused_by` edges (P20 S05 library invariant).
- **Orthogonal sections:** `historical_relationships` is separate from `why`, `deliberation`, `recent_changes` — no merge.
- **Blast radius:** S07 owns mig 021 — confirm S05 did not land experiment tables.
- **Attribution Law 5:** `caused_by` link presence ≠ regression `attribution=caused`.

## Spawn policy

- **Na (implement):** missing section/transition handling, caused_by filter wrong, named test missing/failing, compat/schema bump, new MCP tool
- **Nb (review):** re-review after Na
- Do **not** spawn for optional MCP why parity if CLI tests cover P20 types

## Exit criteria

- [ ] No blocker/high without spawn or inline fix
- [ ] Confidence **high** with test output pasted in board Notes
- [ ] D-11 + D-12 closure evidenced
- [ ] Spawn Na/Nb only if observability/historical gap found

## Next

**P21-S06-00** (unless Na spawned)
