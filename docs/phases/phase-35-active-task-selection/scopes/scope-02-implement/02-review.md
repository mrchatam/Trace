# P35-S02-02 — Implement review (active task pick)

## Metadata
- id: P35-S02-02
- todo_ids: [P35-S02-02]
- role: reviewer
- skills: [code-review-and-quality, security-and-hardening]
- verification: automated
- mcps: []

## Objective

Independent review of **P35-S02-01** vs [`PLAN.md`](../scope-01-plan/PLAN.md) + DESIGN-LOCKS + Law 19. Fresh context (do not share implementer session). Spawn remediation forward on blocker/high FAIL.

## References

- [`PLAN.md`](../scope-01-plan/PLAN.md) — normative policy
- [`01-implement.md`](01-implement.md) — claimed deliverables + locked defaults
- [`INVESTIGATION.md`](../scope-00-investigate/INVESTIGATION.md)
- [`DESIGN-LOCKS.md`](../../DESIGN-LOCKS.md)
- `docs/rules/agent-loop-protocol.md` (reviewer loop)
- Law 19 / Laws 6–7

## Session start

Follow agent-loop-protocol Session start. Re-verify against **repo evidence**, not Notes alone.

## Locked review anchors (from PLAN + S02-00)

| Anchor | Expect |
|--------|--------|
| Placement | Shared `web/src/lib/pickActiveTask.ts` only — **no** Go/HTTP current-work API |
| Semantics | P1→P2→P3a; **not** `tasks[0]` / `items[0]` as all-DONE fallback |
| Screens | Overview + Loop both import helper; no divergent pick forks |
| Overrides | `?task_id=` not overwritten; TRACE_TASK_ID not sole fix |
| Limit honesty | Client fetch-for-pick (`listTasksForPick` or equivalent); **no** S02 OpenAPI/`handlers_tasks.go` pagination |
| Tests | `node --experimental-strip-types --test src/lib/pickActiveTask.test.ts` green; all-DONE ≠ Step1; >100 honesty covered |
| Dogfood | Feet-seller not mutated |
| Out | `plan_missing` not weakened; no full-graph dump; Graph/`overviewCompose` smell not required |

## Checklist

### Policy + Law 19
- [ ] Matches PLAN P1→P2→P3a (not a silent alternate heuristic)
- [ ] Single pick policy; Overview/Loop do not diverge
- [ ] Adapters stay thin; no second business-logic fork in screens
- [ ] Placement A (Go) not sneakily half-implemented

### Wiring
- [ ] Local Overview `pickActiveTask` removed; imports shared helper
- [ ] Loop default uses helper when `!taskId`; `items[0]` auto-pick gone
- [ ] Explicit `?task_id=` preserved
- [ ] Auto-pick list comes from fetch-for-pick (page-until-exhausted / equivalent), not bare truncated `limit: 100` as sole P3a input

### Tests
- [ ] Acceptance #1–6 from PLAN covered (synthetic fixtures)
- [ ] Red→green seed: all-DONE oldest-first → pick ≠ Step1 (`33247e2d-…`)
- [ ] Locked test cmd exit 0 (re-run yourself)
- [ ] >100 / truncated-page honesty satisfied per PLAN prefer-(a)

### Safety / scope
- [ ] No feet-seller data destruction
- [ ] `plan_missing` / PLAN-phase gates not weakened
- [ ] No full-graph dump / Laws 6–7 regression
- [ ] No unauthorized `handlers_tasks.go` / `openapi.yaml` pagination “drive-by”
- [ ] Optional AGENTS.md only if present and accurate

### Findings disposition
- [ ] blocker/high: inline small fix **or** spawn `P35-S02-02a` / `02b` immediately below this row
- [ ] medium: prefer spawn unless trivial
- [ ] Confidence **medium** or **high** with residual risks listed (never silent)

## Exit criteria

- [ ] PASS/FAIL with Notes (cite files + test re-run)
- [ ] No open blocker/high without pending follow-up spawn
- [ ] On PASS, next **P35-S03-00**

## Minimal todos

- [ ] Diff S02-01 claims vs tree (`pickActiveTask.ts`, screens, ops fetch helper, tests)
- [ ] Re-run `cd web && node --experimental-strip-types --test src/lib/pickActiveTask.test.ts`
- [ ] Walk checklist; severity-rank findings
- [ ] Fix or spawn; re-verify
- [ ] Board Notes + status; on PASS hand off to S03

## Next

`P35-S03-00`
