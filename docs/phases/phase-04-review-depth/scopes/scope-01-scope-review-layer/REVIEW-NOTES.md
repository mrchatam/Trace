# P04-S01-02 — Scope review notes (scope review layer)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16

## Summary

Independent review of P04-S01-01 against S01-00 / `01-scope-review-layer.md` locks. Claims match live repo: mig `008_scope_review.sql` + `review_residuals`; `LinkReviewScope` (`review_judges_scope` → `plan_scope`); residual Add/List/CountOpen/SetStatus (INFO|WARN|BLOCKING; OPEN|ACKED|RESOLVED); reuse CreateReview/SetReviewResult; no second review stack / no planner fork; task DONE unchanged (EvidenceIDs-alone reject; scope PASS ≠ DONE; PASS `review_judges_task` + escape hatch); honesty Paths A/B/C fail-closed without `AllowDoneWithoutReview` or residual requirements; thin CLI `--scope` + `residual add|list` (G19); VerifiedFact absent; `plan_scopes.status` not mutated by scope reviews. Gate E / Gate C / p0x / x0 / honesty / `./...` green. S02 stubs already list S01 hooks. No blocker/high; no spawns.

## Checklist (review focus)

| Focus | Result |
|-------|--------|
| Mig `008_scope_review.sql` + `review_residuals` | **Pass** — schema matches lock; Open applies v8; `store_test` requires table + versions 1–8 |
| `LinkReviewScope` / `review_judges_scope` (to=`plan_scope`) | **Pass** — validates GetReview + GetPlanScope; entity_links + event; domain tests |
| Residual Add/List/CountOpen/SetStatus vocabulary | **Pass** — domain + store helpers; fail-closed severity/status; empty code rejected |
| Reuse CreateReview/SetReviewResult; no second stack; no planner fork | **Pass** — APIs in `internal/domain`; planner has zero residual/scope-review imports |
| Task DONE unchanged | **Pass** — `TestScopeReviewDoesNotWeakenDoneGate` + retained EvidenceIDs / PASS / escape tests |
| Honesty A/B/C fail-closed | **Pass** — proof never sets escape hatch; no residual requirement |
| Gate E / Gate C / p0x / x0 | **Pass** — fresh suites; Gate C metrics `dry_run:false` N=3 mean G1 0.800 > B0 0.000 |
| VerifiedFact out; `plan_scopes.status` not mutated | **Pass** — no VerifiedFact type; link test asserts status unchanged |
| Thin CLI G19; no daemon/HTTP/embeddings primary | **Pass** — `cmd/trace/review.go` adapters only; no MCP residual tools |
| S02 stubs list S01 hooks | **Pass** — `CountOpenResidualsByScope` / residual codes / `review_judges_scope` in S02-00/01/SCOPE-TODOS |

## Claims → evidence

| Claim (P04-S01-01 Notes) | Evidence |
|--------------------------|----------|
| Mig 008 `review_residuals` | `internal/store/schema/008_scope_review.sql`; embed via `schema/*.sql` |
| Store Insert/Get/Update/List/CountOpen | `internal/store/residuals.go` |
| `LinkReviewScope` (`review_judges_scope`→`plan_scope`) | `internal/domain/review.go`; consts in `service.go` |
| Residual domain APIs | `internal/domain/residual.go` — Add/SetStatus/List*/CountOpen |
| Severity/status/codes | `ResidualSeverity*` / `ResidualStatus*` / `ResidualCode*` in domain+store |
| Reuse CreateReview/SetReviewResult | Unchanged entrypoints in `review.go`; CLI create/set still call them |
| Task DONE unchanged | `TestScopeReviewDoesNotWeakenDoneGate`; EvidenceIDs-alone reject; scope PASS ≠ DONE |
| Thin CLI `--scope` + `residual add\|list` | `cmd/trace/review.go` + help line |
| VerifiedFact / MCP residual out | Grep: VerifiedFact only in comments; `internal/mcp` has no residual/scope review |
| Gate C intact | `docs/verification/gate-c-x0/metrics-{b0,g1}.json` `dry_run:false` N=3; means 0.000 / 0.800 |
| CGO bars | Fresh re-verify below |

## Required tests (fresh this review)

```text
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./evals/honesty/... ./evals/replan/... -count=1
  → PASS (domain, store, honesty, replan)

CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./... -count=1
  → PASS (p0x, x0, honesty, replan, domain, store, planner, mcp, analyzers, …)
```

Gate C artifacts: not rewritten; `dry_run:false` confirmed; Go inequality intact.

## Findings

### blocker
_None._

### high
_None._

### medium
_None open (no spawn)._

### low

1. **CLI has no `residual set-status`** — API `SetResidualStatus` is library-only; matches locked CLI surface (`add|list` only). S02/consumers use domain APIs.
2. **Store Insert does not re-validate severity vocabulary** — domain `NormalizeResidualSeverity` fail-closes; direct store callers could insert non-canonical severity. Acceptable while all product paths go through domain.

### nit

1. No dedicated `cmd/trace` integration test for `--scope` / residual subcommands (domain coverage is strong).
2. Residual severity/status consts duplicated store↔domain (re-export pattern matches prior migs).

## Spawns
_None._

## Next board row
**P04-S02-00** (honesty escape-rate / Gate G prelim planner). Do not start S02 implement until S02-00 locks.
