# P04 / S01 / 00-PLANNER — Scope review layer

## Metadata
- id: P04-S01-00
- todo_ids: [P04-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-scope-review-layer.md` for **richer scope-level review / evidence policies** against live domain review surface. Lock package paths, APIs, persistence, CLI surface, and exit criteria. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 4
- [docs/init/H_VERIFICATION_STRATEGY.md](../../../../init/H_VERIFICATION_STRATEGY.md) — scope review layer
- Live: `internal/domain` CreateReview/SetReviewResult/LinkReviewTask; DONE=PASS\|escape hatch; `internal/planner` `plan_scopes`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (gaps — lock paths in this planner)
| Item | Today (post–Phase 03 / P04-00) | S01 need (sketch) |
|------|-------------------------------|-------------------|
| Review | Task-level `review_judges_task` PASS / FAIL + `AllowDoneWithoutReview` | **Scope-level** review against `plan_scopes` (rel / API locked here) |
| Residuals | Free-text review body only | Structured **residual tracking** hooks S02 can count |
| Honesty | Paths A/B/C planted claim | Keep green; do not weaken fail-closed; never use `AllowDoneWithoutReview` in honesty |
| Planner | Coarse + discovery replan live | Do not fork planner; review layer only |
| VerifiedFact | Absent (Phase 01 deferred) | **Out** — residuals only, not promotion engine |

## Phase defaults already locked (respect)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Honesty / p0x / x0 / replan / Gate C / Gate E | Keep green / intact |
| Daemon/HTTP/embeddings | Forbidden as primary |
| MCP | Optional; CLI primary |
| VerifiedFact | Out this phase |

## Locked by this planner (2026-08-16)

| Item | Value |
|------|-------|
| Package | **`internal/domain`** + store helpers — **no** second review stack; **no** planner fork |
| Migration | Additive **`008_scope_review.sql`**: table `review_residuals` |
| Scope link | **`review_judges_scope`** (`entity_links`: from=`review`, to=`plan_scope`; entity type **`plan_scope`**) via `LinkReviewScope`; validate scope with `store.GetPlanScope` |
| Reuse | Keep `CreateReview` / `SetReviewResult` / `LinkReviewTask`; same PASS\|FAIL\|UNCERTAIN |
| Residuals | `AddResidual` / `SetResidualStatus` / `ListResidualsByReview` / `ListResidualsByScope` / `CountOpenResidualsByScope` — severity INFO\|WARN\|BLOCKING; status OPEN\|ACKED\|RESOLVED; codes include MISSING_EVIDENCE, OPEN_GAP, CONTRACT_GAP, POLICY_EXCEPTION |
| Scope policy | Recording only — **do not** mutate `plan_scopes.status` or task DONE from scope reviews |
| Task DONE | Unchanged (PASS `review_judges_task` \| `AllowDoneWithoutReview`); honesty never uses escape hatch |
| CLI | Thin `trace review` — `--scope` on create; `residual add\|list` |
| S02 hooks | Count/list OPEN residuals by scope + residual codes; `review_judges_scope` links |
| Out | VerifiedFact; Gate G harness (S02); VERIFY (S03); daemon/HTTP/embeddings |

## Planner work
- [x] Lock review-depth surface and persistence model (prefer extend `internal/domain` + store mig; do not invent a second review stack).
- [x] Thicken `01-scope-review-layer.md` exit criteria enough to run alone.
- [x] Light-update **upcoming** S02 stubs with expected escape-rate hooks from S01 surface.
- [x] Sync SCOPE-TODOS + board Notes; mark this row done.

## Exit criteria
- [x] `01-scope-review-layer.md` runnable alone
- [x] Package path + scope-review + residual model locked
- [x] Light S02 Depends note updated
- [x] No product Go in this row

## Minimal todos
- [x] Inventory live Review / plan_scopes APIs vs S01 needs
- [x] Thicken 01 + 02 + light S02 Depends
- [x] Mark P04-S01-00 done

## Out of scope
- Product Go; Gate G harness (S02); phase VERIFY (S03); VerifiedFact promotion; weakening DONE fail-closed
