# P03-S01-02 — Scope review notes (coarse planner)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16

## Summary

Independent review of P03-S01-01. Claims match live repo: `internal/planner` + mig `006_plan_hierarchy.sql`, Goal→phase→scope hierarchy, supersede-not-delete deep plans, `goal_plan_state` current pointer, progressive DeepPlan (current-only + one shallow lookahead), S02 hooks (`SupersedeDeepPlan`, List/GetCurrent, `auto_replan_count`), thin `trace plan` CLI (G19), no MCP plan tools / daemon / HTTP / embeddings. Honesty, p0x, x0, and `./...` bars green. Gate C packs untouched. No blocker/high; no spawns.

## Checklist (review focus)

| Focus | Result |
|-------|--------|
| Package `internal/planner` (not domain Phase/Scope dump) | **Pass** — no CreatePhase/DeepPlan in `internal/domain`; planner uses `store.GetGoal` |
| Mig `006_plan_hierarchy.sql` embedded | **Pass** — `schema/006_plan_hierarchy.sql` via `//go:embed schema/*.sql`; store tests list new tables |
| Hierarchy Goal→phase→scope + ord | **Pass** — `CreateCoarsePlan` ord=index; lists `ORDER BY phase.ord, scope.ord` |
| `scope_deep_plans` supersede-not-delete | **Pass** — `SupersedeActiveScopeDeepPlans` UPDATE→SUPERSEDED; no DELETE; re-deep keeps prior rows |
| `goal_plan_state` current pointer | **Pass** — Ensure on create (NULL); `SetCurrentScope` / `GetCurrentScope` |
| DeepPlan fail-closed unless current | **Pass** — `ErrNotCurrent`; tested before set-current and for non-current |
| One lookahead shallow only | **Pass** — next scope id+summary; may update lookahead `plan_scopes.body`; no deep revisions for non-current |
| No whole-backlog auto-gen | **Pass** — caller-supplied work items; no LLM/heuristics expanding past current+lookahead |
| S02 hooks | **Pass** — `SupersedeDeepPlan` (not current-gated), `ListScopes`/`GetCurrentScope`/`GetPlan`, column `auto_replan_count` |
| Thin `trace plan` CLI (G19) | **Pass** — `cmd/trace/plan.go` adapter only; help wired; create-coarse ordered `--phase`/`--scope` argv (not `--from` JSON) |
| No MCP plan tools | **Pass** — MCP still six tools (`trace_why`…`trace_review`); no `trace_plan*` |
| Laws: no daemon/HTTP/embeddings primary | **Pass** — no ListenAndServe / embeddings in planner path |
| Honesty / p0x / x0 still green | **Pass** — fresh test run below |
| Gate C packs untouched | **Pass** — `docs/verification/gate-c-x0/` still Mode-B `dry_run:false`; scores not rewritten |
| S02 stubs compatible | **Pass** — consume surface matches; S02 must add count mutator (see residuals) |

## Claims → evidence

| Claim (P03-S01-01 Notes) | Evidence |
|--------------------------|----------|
| `internal/planner` + APIs | `service.go`: New, CreateCoarsePlan, SetCurrentScope, DeepPlan, SupersedeDeepPlan, GetPlan, ListScopes, GetCurrentScope |
| Mig 006 tables | `plan_phases` / `plan_scopes` (+`auto_replan_count`) / `scope_deep_plans` / `goal_plan_state` |
| Store helpers | `internal/store/plan_hierarchy.go` |
| Thin CLI | `cmd/trace/plan.go` + `root.go` `case "plan"` + help line |
| Tests | `planner_test.go`: ord/rejects, set-current foreign reject, deep current+lookahead+supersede, SupersedeDeepPlan smoke, GetPlan, auto_replan_count default 0 |

### Implementer residual notes (verified)

| Note | Verified |
|------|----------|
| Create-coarse CLI uses ordered `--phase`/`--scope` argv (not `--from` JSON) | Yes — `parseCreateCoarseArgs`; help documents same |
| DeepPlan may update lookahead `plan_scopes.body`; no deep revision for lookahead in that call | Yes — `UpdatePlanScopeBody` + `CountScopeDeepPlans`==0 for non-current |
| `auto_replan_count` init/preserve only; S02 owns budget | Yes — DEFAULT 0 on insert; `UpdatePlanScopeBody` does not touch count; no enforce |
| Gate C packs untouched | Yes |

## Required tests (fresh this review)

```text
CGO_ENABLED=0 go test ./internal/planner/... ./internal/store/...   PASS
CGO_ENABLED=0 go test ./evals/honesty/...                           PASS
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/...                PASS
CGO_ENABLED=1 go test ./...                                         PASS
```

## Findings

### blocker
_None._

### high
_None._

### medium
_None open (no spawn)._

### low

1. **`CreateCoarsePlan` non-transactional** — mid-failure can leave partial phases/scopes without `goal_plan_state`. Acceptable for S01 CLI/library callers; S02+ may wrap in a store tx if multi-call durability matters.
2. **DeepPlan body-then-revision order** — lookahead `body` update precedes deep-plan insert; rare failure window leaves body updated without new revision.
3. **No `IncrementAutoReplanCount` (or equivalent) store helper yet** — column + read path exist (`ScopeView.AutoReplanCount`); S02 must add mutator + budget enforce (expected by design).

### nit

1. Re-`DeepPlan` that supersedes emits `plan.deep_superseded` (not `plan.deep_planned`) when prior ACTIVE existed — fine; document for consumers if event semantics matter later.
2. Library `DeepPlan` allows empty `ExitCriteria`; CLI requires `--exit` — intentional thin CLI stricter than library.

## Spawns

_None._

## Upcoming prompt light-touch

S02 Depends note thickened: S01 ships column + read surface only; S02 owns increment/check/`ack` on `auto_replan_count` (add store helper as needed). Do not invent a second planner package.

## Residual risks (explicit)

- Partial coarse-plan writes without transaction (low).
- S02 must implement churn budget on `auto_replan_count` (not started; next scope).
- DPC-global attach residual from Phase 02 unchanged (S02 may scope if measured).

## Board

Mark **P03-S01-02** `done`. Next runnable: **P03-S02-00**.
