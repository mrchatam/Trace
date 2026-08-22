# P20-S03-02 — Review change + effects

## Metadata
- id: P20-S03-02
- todo_ids: [P20-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: automated

## Objective
Independent fresh-session review: no-blob discipline (Law 1); effect comparison honesty; contradiction must not auto-promote to causal regression (S05 owns attribution). Paths live in `change_paths`, not JSON.

## Session start
Follow agent-loop-protocol. Unattended after S03-01 `done`. Reviewer ≠ implementer session. Board: status/notes; spawn forward on blocker/high. Do not rewrite `done` prompts.

## Keeper tests (must re-run)

```bash
go test ./internal/domain/ -count=1 -run 'TestCreateChange|TestRecordExpected|TestRecordActual|TestUnknownEffect|TestContradicted|TestParentChange|TestResolveChange|TestOversizedEffect'
go test ./internal/store/ -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestChange|TestNoSourceContentColumns'
go test ./internal/deliberation/...
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Review checklist

Compare S03-01 Notes + repo to [00-PLANNER.md](00-PLANNER.md) FINAL locks.

- [ ] `017_changes_effects.sql` exists; no 018+; embed/compat ceiling **17**
- [ ] Tables are `changes` + `change_paths` + `effects` only — **no** `paths_json` / `files_json` / patch/diff/blob/content columns
- [ ] `git_commit` is OID reference only; content bytes come from `vcs.Repository.ShowFile`, not SQLite
- [ ] Change APIs do **not** write `vcs_commits` / `vcs_commit_paths` (no Git index duplication)
- [ ] Comparison enum fail-closed on unknown values; actual requires prior expected dimension
- [ ] Contradicted path: optional PLAN_AFFECTING Discovery **or** Hypothesis link (`hypothesis_explains_effect`); Hypothesis is not a Discovery row
- [ ] Linked decision gets FIRED `contradicted_effect` reconsideration; Decision survives
- [ ] **No** Regression row; **no** `attribution=caused`; **no** auto-replan / SelectNext hop
- [ ] No tests/baseline/score columns on `changes` (S04)
- [ ] Parent chain does not auto-SUPERSEDE parent
- [ ] No P19 loop/CLI/MCP edits; no CoT/blobs
- [ ] Law 19: library-only this scope (`vcs` iface, not `gitcli`)
- [ ] Residuals listed if seed/FTS omitted (expected per S03-00)

## Spawn rule

blocker/high: small inline fix **or** insert `P20-S03-02a` (implement) + `P20-S03-02b` (review) immediately below this row with full prompts. Medium: prefer spawn unless trivial.

## Exit criteria

- blocker/high fixed or spawned forward
- confidence medium or high with evidence
- residuals listed explicitly if medium (never silent)
- Next runnable after APPROVE: **P20-S04-00**
