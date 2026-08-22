# P21-S02-02 — Review: retrieval + FTS

## Metadata
- id: P21-S02-02
- todo_ids: [P21-S02-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective
Independent review: P20 types registered in Exact/Why/Expand + FTS; INVESTIGATE stderr residual (D-06) **closed**; D-07 FTS sync **closed**; no regression on P17/P19/P20 keepers or S01 seed.

## Session start
**Fresh subagent** (not S02-01 session). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: status + notes only; spawn Na/Nb if gap found.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-retrieval-fts.md](01-retrieval-fts.md) — implementer deliverable
- [DECISION-LOG.md](../../DECISION-LOG.md) D-06, D-07
- P20 residual: [VERIFY-NOTES.md](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/VERIFY-NOTES.md) item 4

## Review checklist

| Check | Evidence |
|-------|----------|
| 8 P20 types in `lookupEntity` | Grep `case "uncertainty"`…`outcome_result` in `exact.go` |
| Aliases in `NormalizeEntityType` | `outcome`, `plan-change` + canonical pass-through |
| Unknown type fail-closed | Bad type to Exact/Why returns error (not empty hit + log) |
| FTS Must tables in `RebuildFTS` | uncertainties, hypotheses, changes, regressions, reflections |
| FTS Should (if claimed) | baselines, outcome_results — or documented defer in Notes |
| `effect` Exact-only | No FTS insert for `effect` entity_type |
| Cognitive upsert sync | `UpsertUncertainty`/`UpsertHypothesis` call `SyncEntityFTS`; comments removed |
| INVESTIGATE stderr clean | `TestLoopNextInvestigateNoRetrievalStderr` + manual spot-check |
| Expand through uncertainty link | Task with `uncertainty_blocks_task` — `TaskContext` succeeds without `minimalTaskContextPacket` fallback |
| Law 1 FTS bodies | No file blobs/paths in indexed body text |
| S01 seed untouched | No edits to `seed_export.go` / `seed_import.go` |
| No mig 020+ | Schema max **019**; compat ceiling **19** |

## Keeper command floor

```bash
go test ./internal/retrieval/... -count=1 -run 'TestExactLookupUncertainty|TestExactLookupHypothesis|TestWhyUncertaintyIncludesGraphSteps|TestNormalizeEntityTypeP20Aliases'
go test ./internal/store/... -count=1 -run 'TestSyncEntityFTSUncertainty|TestRebuildFTSIncludesP20Types|TestOpenCreatesDBAndMigratesIdempotent'
go test ./cmd/trace -count=1 -run 'TestLoopNextInvestigateNoRetrievalStderr|TestCausalWhyContextRoundTrip|TestLoopNextInvestigateEmphasizesUncertainties|TestLoopApplyUncertaintyWriteAffectsNextSelectNext'
go test ./internal/deliberation/... -count=1
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Review focus

- **D-06 closed:** `retrieval: unknown entity type "uncertainty"` cannot appear on INVESTIGATE `loop next` stderr when blocking uncertainty present.
- **D-07 closed:** P20 cognitive rows appear in FTS after upsert (uncertainty minimum; hypothesis if wired).
- **Fail-closed boundary:** Malformed CLI/MCP type strings error at API boundary — not silent success with degraded logging elsewhere.
- **Blast radius:** S05 `trace why` P20 nouns still pending scope S05 — confirm S02 only registers retrieval; CLI why for P20 types should work as side effect (spot-check `trace why uncertainty <id>` if fixture exists).

## Spawn policy

- **Na (implement):** missing P20 type in lookup/FTS, stderr regression, FTS indexes blobs, or seed touched
- **Nb (review):** re-review after Na
- Do **not** spawn for optional baseline/outcome FTS deferral if documented in implementer Notes

## Exit criteria

- [ ] No blocker/high without spawn or inline fix
- [ ] Confidence **high** with test output pasted in board Notes
- [ ] D-06 + D-07 closed with evidence
- [ ] Spawn Na/Nb only if retrieval/FTS gap found
