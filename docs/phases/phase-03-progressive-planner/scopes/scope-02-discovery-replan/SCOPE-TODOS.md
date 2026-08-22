# Scope S02 — Discovery→PlanChange replan

**Depends-on:** S01 done (`P03-S01-02`). Consume **`internal/planner`** (`SupersedeDeepPlan`, current/list scopes, `GetPlan`, `auto_replan_count` column) — no second planning stack. Keep `LinkDiscoveryPlanChange`.

**Locks (P03-S02-00):**
- Mig `007_discovery_severity.sql` — severity `INFO` \| `PLAN_AFFECTING` \| `BLOCKING`
- Only `PLAN_AFFECTING`+ auto-opens replan
- Churn N=5 (`DefaultMaxAutoReplans`); store `IncrementAutoReplanCount` + `AckAutoReplan` (reset 0)
- `planner.ApplyDiscoveryReplan` + `AckReplan`
- Demo: `evals/replan` / `TestPlantedDiscoveryReplan`
- CLI: `add discovery --severity`, `plan apply-discovery`, `plan ack-replan`
- DPC attach scoping **out**; honesty / p0x / Gate C intact

**Out:** Phase VERIFY / Gate E closeout (S03).

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P03-S02-00 | planner | done | 2026-08-16: locked severity/churn/demo/CLI; 01 runnable alone |
| P03-S02-01 | implement | done | 2026-08-16: mig 007 + ApplyDiscoveryReplan + evals/replan + CLI |
| P03-S02-02 | review | done | 2026-08-16: APPROVE high; no spawns — REVIEW-NOTES.md |

## Checklist

- [x] P03-S02-00 planner
- [x] P03-S02-01 implement
- [x] P03-S02-02 review
