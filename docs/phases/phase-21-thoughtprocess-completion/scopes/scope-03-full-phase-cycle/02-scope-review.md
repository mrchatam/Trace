# P21-S03-02 — Review: full SelectNext phase cycle

## Metadata
- id: P21-S03-02
- todo_ids: [P21-S03-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective
Independent review: extended `SelectNext` matches S03-00 **14-row FINAL** table; EXECUTE gate preserved; EXPLORE optional row correct; D-18 floats observability-only; no ML; P19/P20 loop keepers green; `select.go` remains pure (no I/O).

## Session start
**Fresh subagent** (not S03-01 session). Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: status + notes only; spawn Na/Nb if gap found.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- [01-full-phase-cycle.md](01-full-phase-cycle.md) — implementer deliverable
- [DECISION-LOG.md](../../DECISION-LOG.md) D-03, D-04, D-18, D-21
- [WORK-MAP.md](../../WORK-MAP.md) W-04, W-15

## Review checklist

| Check | Evidence |
|-------|----------|
| 14-row table in `select.go` | Grep reason codes + phase returns; order matches S03-00 |
| Row 6 EXPLORE (D-04) | `open_decision_alternatives > 0` AND `!plan_critiqued` → EXPLORE; without alternatives → CRITIQUE (row 7) |
| Row 7 CRITIQUE | Requires `plan_exists && !plan_critiqued` |
| EXECUTE gate (D-03) | Row 3 blocks before row 8; `TestSelectNextNeverExecuteOnBlockingUncertainty` green |
| Row 8 EXECUTE | Only when `execute_pending` and upstream gates clear |
| Rows 9–13 cycle | TEST, VERIFY (10), EVALUATE, REFLECT, REPLAN fire on respective pending flags |
| Default ORIENT | Row 14 when all flags clear |
| PolicyInputs struct | 6 legacy + 8 new fields with correct JSON tags |
| Reason codes | 6 new constants match table `reason_code` strings |
| Float observability (D-18) | Vary `plan_confidence`/`requirement_coverage` in test — phase unchanged |
| No ML / scoring | No weighted sums, rand, or external models in `select.go` |
| Pure policy | `select.go` has no store/domain/planner imports |
| Transition payload | `ApplyTransition` embeds extended `policy_inputs`; optional floats when set |
| 9 named tests exist + PASS | All tests from S03-00 named table |
| Legacy tests green | `TestSelectNext` table cases still valid (CRITIQUE row renumbered) |
| No mig 020+ | Schema max **019**; compat ceiling **19** |
| Loop schema unchanged | `NextSchemaVersion` still `trace.loop.next.v1` |
| BuildPolicyInputs | May still stub new fields — document in Notes if unwired (expected until S06) |

## D-03 / D-04 / D-18 closure

- **D-03 promote:** EXECUTE/TEST/EVALUATE/REFLECT/REPLAN returned by SelectNext when pending flags set — not MVP 8-row only.
- **D-04 keep merged:** EXPLORE only as optional row 6 when alternatives exist; no separate DecisionAlternative entity work.
- **D-18 optional:** Float fields in JSON for observability; must not reorder table.

## Keeper command floor

```bash
go test ./internal/deliberation/... -count=1
go test ./cmd/trace -count=1 -run 'TestLoopNext|TestLoopApplyDeliberationTransitionEvent|TestLoopApplyUncertaintyWriteAffectsNextSelectNext|TestLoopApplyRegressionWriteAffectsPolicyInputs|TestLoopNextDeliberationSectionPresent|TestLoopNextExecuteEmphasizesContextAndRelated'
go test ./internal/domain/... -count=1 -run 'TestApplyDeliberationTransition|TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition|TestHasOpenRegressionFeedsApplyDeliberationTransition'
go test ./internal/loop/... -count=1
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Review focus

- **First-match ordering:** `TestSelectNextFullCycleOrdering` or manual trace proves VERIFY (row 10) beats EVALUATE when both flags set incorrectly — table order is law.
- **CRITIQUE vs EXPLORE:** With alternatives + `!plan_critiqued`, EXPLORE wins (row 6 before 7).
- **EXECUTE simulation:** Blocking uncertainty + `execute_pending=true` must still yield INVESTIGATE.
- **Blast radius:** S06 will wire `BuildPolicyInputs` for pending flags — confirm S03 did not partially wire incorrect semantics.
- **P20 keeper:** Hop budget, P19 saturated, ApplyTransition sticky INVESTIGATE tests unchanged.

## Spawn policy

- **Na (implement):** table row missing/wrong order, EXECUTE gate regression, floats affect priority, I/O in select.go, named test missing/failing
- **Nb (review):** re-review after Na
- Do **not** spawn for BuildPolicyInputs still stubbing new fields (expected; S06 scope)

## Exit criteria

- [ ] No blocker/high without spawn or inline fix
- [ ] Confidence **high** with test output pasted in board Notes
- [ ] D-03 closure evidenced (EXECUTE returned when pending)
- [ ] Spawn Na/Nb only if table/gate gap found

## Next

**P21-S04-00** (unless Na spawned)
