# S05 — Phase VERIFY — scope todos

**Depends-on:** P10-S04-02 done.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | VERIFY planner | **done** — evidence + DR-HANDOFF locked |
| 2 | 01-verify | verify | **done** — PASS; VERIFY-NOTES; DR-HANDOFF started=`no successor` |
| 3 | 02-scope-review | review | close phase; **own** DR-HANDOFF completion |

## Locked evidence imports (P10-S05-00)

| Scope | Named regressions (must re-prove) |
|-------|-----------------------------------|
| S01 | DF-19 `TestWhyTaskDPCGoalScoped` + `TestWhyTaskDPCMultiGoalNoForeignPollution` (+ compiler DPC scoped); DF-23 `TestExactWhyPlanChangeAlias`; DF-25 `TestExactWhyCapability`; DF-27 `TestDecisionMarkdownTrustLabels`; DF-29 `TestIncludeWhyFailClosed` |
| S02 | Nine MCP tools; `TestToolNamesRegistered` / BuiltinMCP×9; DF-21/22/32 tasks/capability/version + snake_case + install tip; G19; no plan/impact/index MCP |
| S03 | DF-20 `TestIndexGCAfterPathRename` + argv-only delete + `TestIndexIncrementalIsolation` |
| S04 | DF-17/18/24/26/31 operator/reopen/missing-cap/hatch UX tests; honesty Path C `AllowOperatorDone`; Gate G hatch retained; MCP `trace_transition` same gates |

## Carry-forward
Honesty A/B/C+G; Gate E; Gate F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` N=3; product `./...` (graphify space FAIL OK).

## Reminders
- Dry-run ≠ Gate C / F / G / ablation / H / checklist
- Forward-only: do not rewrite Phase 09 history
- Spawn on fail: `P10-S05-01a` / `01b` / (`01c`) immediately below
- **DR-HANDOFF:** default **`no successor`** — S05-01 starts Notes; **S05-02 owns completion**
- Residual OK: `plan_scope` Exact out; Mode-B historical; Cursor MCP reload; optional ab-* re-runs
