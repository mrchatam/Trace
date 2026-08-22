# P21-S01-02 — Review: portable P20 seed

## Metadata
- id: P21-S01-02
- todo_ids: [P21-S01-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective
Independent review: P20 seed round-trip + P17 keeper regression. Confirm D-05 omit policy **retired** (portable clone includes cognition).

## Session start
**Fresh subagent** (not S01-01 session). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: status + notes only; spawn Na/Nb if gap found.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-portable-p20-seed.md](01-portable-p20-seed.md) — implementer deliverable
- P20 superseded policy: [01-verify.md](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/01-verify.md) seed omit bar

## Review checklist

| Check | Evidence |
|-------|----------|
| 11 tables exported/imported | Grep `SeedDocument` fields + `BuildSeedDocument` / `ImportSeedDocument` paths |
| 10 JSON keys (paths nested) | `deliberation_states`, `uncertainties`, `hypotheses`, `decision_reconsiderations`, `changes`+`paths`, `effects`, `outcome_results`, `baselines`, `regressions`, `reflections` |
| Denied surfaces still omitted | `TestSeedExportOmitsDeniedSurfaces` — no transitions, work_state, reviews, caps, tokens, index |
| Idempotent re-import | `TestSeedImportP20RoundTrip` — duplicate import no extra rows |
| Import order FK-safe | baselines → outcomes → changes/effects → cognitive → regressions/reflections → deliberation_state |
| No mig 020+ | `evals/compat` ceiling **19** |
| No blob/path violations | changes export git SHA + path refs only (Law 1) |
| Old seeds still import | version 1 doc without P20 keys imports cleanly |
| ListAll helpers only | No new SQL migrations under `internal/store/schema/` |

## Keeper command floor

```bash
go test ./internal/domain/... -count=1 -run 'TestSeedExportIncludesP20Cognition|TestSeedImportP20RoundTrip'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./cmd/trace -count=1 -run 'TestLoopApplyDeliberationTransitionEvent|TestLoopApplyUncertaintyWriteAffectsNextSelectNext'
```

## Spawn policy

- **Na (implement):** export/import gap for any of 11 tables, idempotency break, or denied surface leak
- **Nb (review):** re-review after Na
- Do **not** spawn for style-only nits

## Exit criteria

- [ ] No blocker/high without spawn or inline fix
- [ ] Confidence **high** with test output pasted in board Notes
- [ ] D-05 closed: portable P20 cognition in seed JSON confirmed
- [ ] Spawn Na/Nb only if export/import gap found
