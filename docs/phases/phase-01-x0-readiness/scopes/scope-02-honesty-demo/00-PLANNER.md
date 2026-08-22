# P01 / S02 / 00-PLANNER — Honesty demo

## Metadata
- id: P01-S02-00
- todo_ids: [P01-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-honesty.md` for **Honesty demo** against live S01 promotion API. Lock deterministic fail-closed scenario. No product code beyond what demo/tests need (scope of implement row).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — H5 honesty
- [docs/init/H_VERIFICATION_STRATEGY.md](../../../../init/H_VERIFICATION_STRATEGY.md)
- Sibling S01 prompts + board Notes after S01 done
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Planner work
- Read S01 Notes + live Review/DONE API.
- Thicken `01-honesty.md`: exact scenario, package path (`evals/` or `internal/domain` test), success/fail asserts.
- Prefer `verification: automated`; mark `mixed` only if unavoidable.
- Sync SCOPE-TODOS.md.

## Depends-on
- **S01** review `done` (promotion rules must exist).
- **S01 surface (P01-S01-00 locks):** DONE iff linked Review PASS or explicit escape hatch; EvidenceIDs alone insufficient; APIs `CreateReview` / `SetReviewResult` / `LinkReviewTask`; honesty demo must prove FAIL/no-PASS path — not `--allow-done`.

## Exit criteria
- [x] `01-*` runnable without guessing
- [x] SCOPE-TODOS + TODO.md Notes updated
- [x] No product feature work in this planner row

## Minimal todos
- [x] Inspect S01 surface (REVIEW-NOTES + live review/task_state/claim APIs)
- [x] Thicken 01 + 02 prompts; light S03 package-boundary notes
- [x] Sync SCOPE-TODOS / board Notes
