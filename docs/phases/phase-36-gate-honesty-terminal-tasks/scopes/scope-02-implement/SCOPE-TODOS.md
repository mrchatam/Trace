# Scope 02 — board map

**S02 implement** — library + MCP + install + GUI adapter per [PLAN.md](../scope-01-plan/PLAN.md). Serial: **P36-S02-00 → P36-S02-01 → P36-S02-02**. PLAN.md locked by P36-S01-01.

| Order | Board ID | Prompt | Role | Artifact |
|------:|----------|--------|------|----------|
| 628 | P36-S02-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Thickened implement + review prompts |
| 629 | P36-S02-01 | [01-implement.md](01-implement.md) | Implementer | Code + tests |
| 630 | P36-S02-02 | [02-review.md](02-review.md) | Reviewer | Independent review |

## Locked fix set (PLAN.md §2 — do not re-debate in S02)

| § | Fix | S02 action |
|---|-----|------------|
| 2.1 | MCP `trace_plan` (16 tools) | Implement |
| 2.2 | Bootstrap CLI + MCP | Implement |
| 2.3 | Install contract bootstrap step | Implement |
| 2.4 | PlanExists bridge | **Defer** — no code |
| 2.5 | Terminal gate `goal_plan_gap_terminal_advisory` | Implement |
| 2.6 | Enforce document nudge | Implement |
| 2.7 | Goal structure warning N=15 | Implement |
| 2.8 | Feet-seller recovery | Ship tool + temp test; live verify S03 |

## Touch-list order (PLAN.md §5 — verified live 2026-08-22)

```text
library (planner bootstrap/advisory)
  → library (loop gate terminal honesty)
  → MCP (tools_plan.go + server.go → 16 tools)
  → CLI (plan.go bootstrap + show advisory)
  → install (enforcement.go + install.go)
  → config (enforce.go + init.go nudge)
  → GUI (GateStrip/TaskDetail adapter — secondary)
  → HTTP POST plan routes DEFER
```

### File targets

| File | Action | PLAN § |
|------|--------|--------|
| `internal/planner/bootstrap.go` | Create | 2.2 |
| `internal/planner/bootstrap_test.go` | Create | 6.2 |
| `internal/planner/advisory.go` | Create | 2.7 |
| `internal/planner/advisory_test.go` | Create | 2.7 |
| `internal/planner/service.go` | Edit if needed | 2.2, 2.7 |
| `internal/loop/gate.go` | Edit | 2.5 |
| `internal/loop/gate_test.go` | Edit | 6.2, 6.3 |
| `internal/loop/testdata/feet-export-min.json` | Create | 6.2 |
| `internal/mcp/tools_plan.go` | Create | 2.1 |
| `internal/mcp/server.go` | Edit | 2.1 |
| `internal/mcp/mcp_test.go` | Edit | 6.1 |
| `cmd/trace/plan.go` | Edit | 2.2, 2.7 |
| `internal/install/enforcement.go` | Edit | 2.3 |
| `cmd/trace/install.go` | Edit | 2.3 |
| `internal/config/enforce.go` | Edit | 2.6 |
| `cmd/trace/init.go` | Edit | 2.6 |
| `web/src/components/GateStrip.tsx` | Edit/verify | 2.5 |
| `web/src/screens/TaskDetail.tsx` | Edit optional | 2.5 |

**Non-touch:** `internal/loop/policy.go`, `internal/deliberation/select.go`, `internal/httpapi/server.go` (POST defer)

## Acceptance tests (PLAN.md §6 — S02 must pass)

| Name | File | Purpose |
|------|------|---------|
| `TestGreenfield_MCPPlanBootstrap_EditGatePasses` | `internal/mcp/mcp_test.go` | MCP-first agent path |
| `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap` | `internal/loop/gate_test.go` | Terminal honesty + bootstrap recovery |
| `TestActiveWork_PlanMissingStillBlocksEdit` | `internal/loop/gate_test.go` | Regression guard |

Recommended: `TestEvaluateGate_Done_TerminalPlanGapAdvisory`, `TestPlanBootstrap_Idempotent`, `TestGoalStructureWarning_OverThresholdNoPlan`, 16-tool registration lock.

## Feet-seller recovery scope

- **S02:** Read-only export import in tests; bootstrap command shipped
- **S03:** Live bootstrap on `/home/ali/Desktop/feet seller telegram app` with human approval

## Out of scope

- VERIFY live dogfood (S03)
- PlanExists bridge (§2.4 defer)
- HTTP POST plan mutation routes
- MCP `trace_loop action=gate`
- Default enforce strict
- Scope creep beyond PLAN.md
