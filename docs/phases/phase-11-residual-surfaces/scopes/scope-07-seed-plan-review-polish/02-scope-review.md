# P11 / S07 / 02 — Scope review (Seed / plan / review show polish)

## Metadata
- id: P11-S07-02
- todo_ids: [P11-S07-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S07 (**DF-28, DF-30, DF-33, DF-45, DF-46**). Fresh subagent. Compare claims + locks to live code/tests. Spawn `02a`/`02b` for blocker/high. Do not rewrite prior `done` history.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- Sibling [01-seed-plan-review-polish.md](01-seed-plan-review-polish.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-28, DF-30, DF-33, DF-45, DF-46
- Live: `cmd/trace/{seed,plan,review,help}.go`; `internal/planner`; domain/store review list helpers
- Prior: P11-S06 no seed/plan/review coupling; P10 DF-32 snake_case pattern

## Checklist (must all pass for APPROVE)

| # | Check |
|---|--------|
| 1 | DF-33: seed accepts `from_id`/`to_id`; empty-endpoint error mentions alias keys |
| 2 | DF-30: `plan show` with tasks-only goal → `phases` is `[]` not null; `tasks` present |
| 3 | DF-46: plan show JSON snake_case (`goal_id`, not `GoalID`); nested phase/scope keys snake_case |
| 4 | DF-45: `review get`/`show`/`list` (+ optional `--task`); empty list `[]`; create/set/residual retained |
| 5 | DF-28: help (and README if touched) declares thin handoff SoT (task body + Trace-pull); **no** new entity/mig/MCP |
| 6 | G19 — no domain fork in adapters; migration none; no new MCP tools |
| 7 | No forbidden architecture (daemon/HTTP/full-rebuild) |
| 8 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + P11-S01–S06 + Gate C `dry_run:false` |
| 9 | Board Notes accurate; planner row had no product Go |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Prefer named asserts: `TestSeedImportFromIDAliases`, `TestPlanShowSnakeCaseAndEmptyPhases`, `TestPlanShowWithPhasesSnakeCase`, `TestReviewGetShowList`, `TestHelpHandoffSoT` (or equiv).

## Exit criteria
- [x] Checklist evidenced; confidence high (or medium with residuals)
- [x] Board status + Notes; next **P11-S08-00** (unless spawn)
- [x] Write [REVIEW-NOTES.md](REVIEW-NOTES.md) on APPROVE / spawn
