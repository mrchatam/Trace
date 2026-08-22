# P36-S02-02 — Review implement

## Metadata
- id: P36-S02-02
- todo_ids: [P36-S02-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter]
- verification: automated
- hooks: []

## Objective

Independent review of P36-S02-01 deliverables vs [PLAN.md](../scope-01-plan/PLAN.md) + [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md). Re-run tests. Small inline fixes OK; spawn `02a`/`02b` for structural gaps. **No new features.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [PLAN.md](../scope-01-plan/PLAN.md) — §2 fix set, §4 preserve, §6 acceptance tests, §7 non-goals
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md)
- [01-implement.md](01-implement.md)
- Implementer board Notes on P36-S02-01

## Session start

Fresh subagent. Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Confidence bar | **High** (or **Medium** with explicit residual risks listed) |
| Spawn threshold | Blocker/high without pending follow-up → spawn implement+review pair |
| Live dogfood | Feet-seller mutation is **S03** — S02 reviewer checks temp import tests only |
| PlanExists bridge | Must **not** appear in `policy.go` or deliberation |

## Preflight / Plan

1. Read implementer Notes + `git diff` / changed files list
2. Map each PLAN §2 accept item to repo evidence
3. Run test suite for touched packages
4. Spot-read Law 19 boundaries (MCP/CLI/GUI vs library)

## Checklist — fix set (PLAN §2)

### §2.1 MCP `trace_plan` — ACCEPT

- [ ] `internal/mcp/tools_plan.go` exists; thin adapter over `internal/planner.Service`
- [ ] Actions: `create-coarse`, `set-current`, `deep`, `show`, `bootstrap` (single tool, `action` param — mirrors `trace_loop`)
- [ ] `internal/mcp/server.go` registers handler; `RegisteredToolNames()` returns **16** tools including `trace_plan`
- [ ] `internal/mcp/mcp_test.go` locked name list updated (was 15)
- [ ] Error/JSON contract matches CLI (`cmd/trace/plan.go`) — no business logic fork in MCP
- [ ] No MCP `trace_loop action=gate` added (out of scope)

### §2.2 Bootstrap command — ACCEPT

- [ ] `internal/planner/bootstrap.go` exists with plan-change → progressive planner heuristic
- [ ] `trace plan bootstrap --goal <id>` in `cmd/trace/plan.go`
- [ ] MCP `trace_plan action=bootstrap` wired
- [ ] Idempotent when goal already has `current_scope_id` + `current_deep_plan`
- [ ] No LLM generation; minimal recovery documented in help text
- [ ] `TestPlanBootstrap_Idempotent` passes (or equivalent coverage)

### §2.3 Install contract — ACCEPT

- [ ] `internal/install/enforcement.go` — bootstrap step in `AgentsEnforcementBlock()`, `cursorRulesMDCContent()`, `EnforcementRulesMarkdown()`
- [ ] Step placed between TRACE_TASK_ID setup and pre-edit gate (mentions `create-coarse` / MCP `trace_plan`)
- [ ] `cmd/trace/install.go` success output hints bootstrap when goals exist, planner empty
- [ ] Default enforce mode **unchanged** (still opt-in off/warn/strict)
- [ ] Install tests updated if present (`enforcement_test.go` snapshot on `create-coarse` mention)

### §2.4 PlanExists bridge — DEFER

- [ ] **No** auto-satisfy `PlanExists` from plan-change density in `policy.go` or elsewhere
- [ ] No silent bridge heuristic shipped

### §2.5 Terminal gate honesty — ACCEPT

- [ ] `internal/loop/gate.go` — DONE/SKIPPED tasks with goal plan gap emit advisory, not actionable block
- [ ] `reason_code: goal_plan_gap_terminal_advisory`; `allowed: true` for `--for done` and `--for edit` on terminal tasks
- [ ] Verification/regression/uncertainty blocks preserved **before** terminal advisory (`gate.go:243–257`)
- [ ] Active non-terminal still blocked with `plan_missing`, `allowed: false`
- [ ] HTTP/GUI follow library JSON only (Law 19)

### §2.6 Enforce nudge — ACCEPT

- [ ] Stderr hint when `.trace/` exists without valid config (init and/or install surfaces)
- [ ] Suggests `enforce: warn` after install — **no** default mode change
- [ ] `LoadEnforceMode` behavior unchanged for missing config → `EnforceOff`

### §2.7 Goal structure warning — ACCEPT

- [ ] Threshold **N = 15**; condition `task_count(goal) > N && !PlanExists(goal)`
- [ ] Surfaces: plan show stderr, MCP show field, optional loop status advisory
- [ ] Warning only — does not weaken gate

### §2.8 Feet-seller recovery — S02 ship + S03 verify

- [ ] Bootstrap tool shipped (CLI + MCP)
- [ ] Temp import test only in S02 — **no** live mutation of dogfood fixture
- [ ] S03 live verify **not** claimed done in S02 Notes

## Checklist — acceptance tests (PLAN §6)

Run and verify:

```bash
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./cmd/trace/... -count=1 -run 'Greenfield|FeetSeller|ActiveWork|TerminalPlanGap|Bootstrap_Idempotent|GoalStructure|TracePlan|ToolNames'
```

### §6.1 `TestGreenfield_MCPPlanBootstrap_EditGatePasses`

- [ ] Temp dir init → goal + task → MCP `trace_plan` chain
- [ ] Edit gate: `allowed: true`, `reason_code` ≠ `plan_missing`
- [ ] Loop status: `policy_inputs.plan_exists: true`

### §6.2 `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap`

- [ ] Fixture `internal/loop/testdata/feet-export-min.json` (or equivalent) — 11 plan-changes, 0 planner
- [ ] Pre-bootstrap: DONE task `--for done` → `allowed: true`, advisory `goal_plan_gap_terminal_advisory` (not identical actionable premature_implementation alarm)
- [ ] Post-bootstrap: edit gate `allowed: true`, `plan_exists: true`

### §6.3 `TestActiveWork_PlanMissingStillBlocksEdit`

- [ ] Non-terminal task, no planner → edit `allowed: false`, `reason_code: plan_missing`, `recommended_phase: PLAN`
- [ ] Extends / guards `TestEvaluateGate_Edit_PlanMissing` pattern (`gate_test.go:198–205`)

### Recommended (spot-check)

- [ ] `TestEvaluateGate_Done_TerminalPlanGapAdvisory`
- [ ] `TestGoalStructureWarning_OverThresholdNoPlan`
- [ ] `TestRegisteredToolNames_IncludesTracePlan` / 16-tool lock

## Checklist — preserve invariants (PLAN §4)

- [ ] Non-terminal + `!PlanExists` → edit/done blocked, `reason_code: plan_missing`
- [ ] Deliberation order fail-closed: PLAN before CRITIQUE before EXECUTE (`select.go:28–38`)
- [ ] `PlanExists` requires progressive planner rows (`policy.go:45–49`) — unchanged
- [ ] `PlanCritiqued` via plan-changes does **not** substitute for `PlanExists`
- [ ] Verification/regression blocks on done preserved
- [ ] No global weaken of PLAN phase for active work

## Checklist — Law 19 & non-goals (PLAN §5, §7)

- [ ] Policy changes in `internal/loop/` and `internal/planner/` — not in `web/` or MCP business logic
- [ ] `web/src/components/GateStrip.tsx` — adapter only; maps library `allowed` + violations to warn/error
- [ ] `web/src/screens/TaskDetail.tsx` — no independent gate evaluation
- [ ] **No** HTTP POST plan mutation routes added
- [ ] **No** GUI-only patch hiding `plan_missing` without library terminal semantics
- [ ] **No** feet-seller history deletion
- [ ] **No** default enforce strict
- [ ] **No** unrelated refactors outside touch-list

## Checklist — touch-list completeness

Verify implementer touched (or explicitly skipped with reason) each PLAN §8 file:

| File | Expected |
|------|----------|
| `internal/planner/bootstrap.go` | Created |
| `internal/planner/advisory.go` | Created |
| `internal/loop/gate.go` | Edited |
| `internal/loop/gate_test.go` | Edited |
| `internal/loop/testdata/feet-export-min.json` | Created |
| `internal/mcp/tools_plan.go` | Created |
| `internal/mcp/server.go` | Edited |
| `internal/mcp/mcp_test.go` | Edited |
| `cmd/trace/plan.go` | Edited |
| `internal/install/enforcement.go` | Edited |
| `cmd/trace/install.go` | Edited |
| `internal/config/enforce.go` | Edited |
| `cmd/trace/init.go` | Edited |
| `web/src/components/GateStrip.tsx` | Edited or verified no change needed |
| `internal/httpapi/server.go` | **Not** mutated (defer) |
| `internal/loop/policy.go` | **Not** mutated (bridge defer) |

## Reviewer loop

```text
1. Compare deliverables to checklist above
2. Classify findings: blocker | high | medium | low | nit
3. blocker/high: small fix OR spawn 02a implement + 02b review immediately below this row
4. Re-run tests after fixes
5. UNTIL no open blocker/high without pending follow-up
   AND confidence medium or high with evidence
```

## Todo updates

Set own row `done` + Notes: PASS/FAIL, confidence, findings count, spawn IDs if any.

## Exit criteria

- [ ] All three primary acceptance tests verified passing
- [ ] PLAN §2 accept items evidenced; §2.4 defer confirmed absent
- [ ] PLAN §4 preserve invariants hold
- [ ] Law 19 — no business-logic fork in `web/` or MCP
- [ ] No open blocker/high without spawned follow-up row
- [ ] Confidence **high** (or **medium** with residual risks in Notes)
- [ ] Next: **P36-S03-00** on PASS

## Next

`P36-S03-00` (on PASS)
