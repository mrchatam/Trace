# P21-S07-02 — Review: promoted MVP cuts

## Metadata
- id: P21-S07-02
- todo_ids: [P21-S07-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective
Independent review: thin §16 experiments + minimal §18 risk hints match S07-00 locks. No bake-off engine, no ML, no autonomous test runner. **Mig 021 only.**

## Session start
**Fresh subagent** (not S07-01 session). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: status + notes only; spawn Na/Nb if gap found.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-promoted-mvp-cuts.md](01-promoted-mvp-cuts.md) — implementer deliverable
- [DECISION-LOG.md](../../DECISION-LOG.md) D-01, D-02
- [WORK-MAP.md](../../WORK-MAP.md) W-12, W-13

## Review checklist

| Check | Evidence |
|-------|----------|
| Migration | **`021_experiments.sql`** only; **21** schema files; **no 022+** |
| Experiment table minimal | Columns match lock — no candidate arrays, no runner state, no blob fields |
| Status CHECK | `planned`\|`running`\|`completed` only |
| Outcome link | `outcome_result_id` on experiment row; **no** new `outcome_results.kind` |
| No bake-off runner | No multi-candidate compare logic; no agent orchestration |
| No subprocess runner | `TestNoExperimentRunnerInvoked` green; grep experiment path for `os/exec` |
| Risk hints section | `loop next` JSON has top-level `risk_hints` |
| Hint shape | `{code, severity, detail}` only |
| Hint codes | Subset of locked four — no ad-hoc codes |
| Cap | `len(risk_hints.items) <= 4` always |
| Priority order | When all fire, `blocking_uncertainty` first, then `missing_verification`, `many_paths`, `high_churn_path` |
| Deterministic | Same DB state → same hints (no randomness / ML) |
| Thresholds | many_paths **>8** on latest change; churn **≥3** changes/path |
| Advisory only | Hints do **not** change SelectNext rows or auto-execute tests |
| Loop schema string | `schema_version` still `trace.loop.next.v1` |
| Compat ceiling **21** | `TestCompatibilitySecurityChecklist` + embed max **21** |
| Seed export | **Unchanged** — no experiments key added |
| MCP | Still **10** tools |
| 6 named tests | All exist + PASS |
| S06 blast radius | Apply tx tests still green — S07 did not revert WithTx |

## D-01 / D-02 closure

- **D-01 promote:** Thin experiment records exist — table + domain CRUD + optional outcome link; **not** full §16 bake-off engine.
- **D-02 promote:** Deterministic risk hints in loop next — **not** ML risk-adaptive test-selection engine.

## Keeper command floor

```bash
go test ./internal/domain/... -count=1 -run 'TestCreateExperimentLinksOutcome|TestExperimentStatusLifecycle|TestNoExperimentRunnerInvoked'
go test ./internal/loop/... -count=1 -run 'TestRiskHintsManyPaths|TestRiskHintsBlockingUncertainty|TestLoopNextRiskHintsBounded'
go test ./internal/store/... -count=1 -run 'TestMigrationStatusReportsEmbedMax|TestOpenCreatesDBAndMigratesIdempotent'
go test ./cmd/trace -count=1 -run 'TestLoopNextPacketShape|TestLoopNextHistoricalRelationshipsSection'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

Document schema file count in board Notes (expect **21**).

## Review focus

- **Thin §16:** Table is metadata + link — agents still run comparisons externally (D-16).
- **§18 minimal:** Hints surface verification risk; they do **not** select or run tests.
- **Orthogonal sections:** `risk_hints` separate from `deliberation`, `verification_debt`, `recent_changes`.
- **Compat bump:** Only S07 adds **021** — confirm S04–S06 did not pre-land experiment tables.
- **Hard boundary grep:** No `exec.Command("go", "test"` or similar in `experiments.go` / `risk_hints.go`.

## Spawn policy

- **Na (implement):** missing migration/table, runner invoked, hint cap wrong, compat ≠ 21, named test missing/failing, new MCP tool, SelectNext behavior changed by hints
- **Nb (review):** re-review after Na
- Do **not** spawn for optional seed export of experiments (explicitly out of S07 scope)

## Exit criteria

- [ ] No blocker/high without spawn or inline fix
- [ ] Confidence **high** with test output pasted in board Notes
- [ ] D-01 + D-02 closure evidenced
- [ ] Spawn Na/Nb only if experiment/risk-hint gap found

## Next

**P21-S08-00** (unless Na spawned)
