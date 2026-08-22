# P17 / S03 / 02 — idempotent import review

## Metadata
- id: P17-S03-02
- todo_ids: [P17-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of **DF-81/83/84** import implementation vs **00-PLANNER FINAL**. Spawn 02a/02b on blocker/high. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Sibling **FINAL:** [00-PLANNER.md](00-PLANNER.md)
- Implementer row: [01-idempotent-import.md](01-idempotent-import.md)

## Session start
Follow agent-loop-protocol. Fresh subagent. Compare deliverables to repo evidence.

## Review checklist

1. **Second import** — same seed file exits 0; no UNIQUE/PK failures on links, plan rows, findings, alternatives
2. **Duplicate links** — `InsertLinkOrIgnore` or equivalent; duplicate endpoints no-op; link count stable
3. **Entity upsert** — UUID last-wins; **`entity.created` insert-only** (no duplicate events on re-import)
4. **Task `work_state`** — not overwritten by seed upsert
5. **Plan tree** — `UpsertPlanPhase/Scope/DeepPlan` by id; `goal_plan_state` last-wins on `goal_id`; SUPERSEDED deep plans round-trip
6. **Last-wins** — `TestSeedImportSameIdLastWins` covers entity + plan field overwrite
7. **Findings/alternatives** — upsert by id (round-trip fixture re-import clean)
8. **Transitions** — skip when already at target state (if implemented in import path)
9. **Merge docs** — CONTRIBUTING union-by-id for entities **and plan arrays**; no merge driver mention added
10. **`exported_at_commit`** — still ignored for identity
11. **G19** — domain seed-import; thin CLI; no MCP
12. **Keepers** — `TestSeedExportRoundTrip`, S02 `TestHelpSeedExportPath`, P16 seed import tests green
13. **Scope boundary** — no export/help/gitignore/hook changes

## Locked verify (reviewer runs)

```text
CGO_ENABLED=0 go test ./cmd/trace/... ./internal/domain/... ./internal/store/... -count=1 \
  -run 'TestSeedImportIdempotent|TestSeedImportDuplicateLinksNoOp|TestSeedImportSameIdLastWins|TestSeedImportPlanTreeIdempotent|TestSeedExportRoundTrip|TestSeedImport|TestHelpSeedExportPath'
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Record results in **REVIEW-NOTES.md** (create if missing).

## Exit criteria
- [ ] APPROVE (medium+ confidence) + REVIEW-NOTES.md **or** spawn 02a/02b
- [ ] No open blocker/high without pending follow-up
- [ ] Board Notes → **P17-S04-00**
