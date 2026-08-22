# FM-07 / FR-P28-04 — decision (P28-S06-07)

**Date:** 2026-08-20  
**Status:** Accepted — **remain warn-only**  
**Row:** P28-S06-07

## Context

FM-07 (git-sparsity) detects when `trace/graph.json`’s `exported_at_commit` lags `git HEAD` — typically after post-hoc SPEC / product commits without a re-export. Operators may also see this after Session-A/B work that commits without refreshing the portable graph.

FR-P28-04 required an **explicit** decision: keep the harness check warn-only, **or** ship a plan-before-edit fail-closed gate with tests. Silent policy flip to hard-fail is forbidden.

## Decision

**Remain warn-only.** No human approval was given for plan-before-edit fail-closed; product does not gain a commit-ordering gate in this row.

| Outcome | Meaning |
|---------|---------|
| SHA match | `pass` FM-07 |
| SHA mismatch | `WARN` + re-export recommended — **does not** fail G2 alone |
| No git / no `exported_at_commit` | `skip` |

**warn ≠ fail.** Thin-graph / honesty failures remain T02 (`seed export --strict --enforce`). FM-07 drift never increments FAIL by itself.

## Evidence (live harness)

[`experiments/ab-p25-gap-pass-validation/score.sh`](../../../../../experiments/ab-p25-gap-pass-validation/score.sh) T03 (~L135–157):

- Header comment: FM-07 stays warn-only
- Drift → `echo "WARN  FM-07 … — re-export recommended"`
- Match → `pass`; missing → `skip`

Planner lock: [`00-PLANNER.md`](00-PLANNER.md) — default stay warn-only unless explicit human decision to ship plan-before-edit.

Harness operator text: [`PROTOCOL.md`](../../../../../experiments/ab-p25-gap-pass-validation/PROTOCOL.md) FM-07 paragraph.

## Alternatives considered

### Silent hard-fail on SHA drift
- Rejected: contradicts FR acceptance (“do not silently turn FM-07 into hard fail”); breaks existing PASS semantics when export lags HEAD after legitimate commits.

### Plan-before-edit fail-closed (product gate + tests)
- Deferred: material product decision; requires human approval (none given this wave). Out of scope for P28-S06-07.

## How to upgrade later

1. Human explicitly approves plan-before-edit / fail-closed commit ordering.
2. Ship product + harness tests proving drift fails only under that mode.
3. Supersede this decision with a forward residual-wave / phase note (do not rewrite this file’s Accepted status — add supersession pointer).

## Cross-links

- S05 VERIFY: [`../scope-05-verify/VERIFY-NOTES.md`](../scope-05-verify/VERIFY-NOTES.md) (additive residual-wave pointer)
- Review row: P28-S06-08
