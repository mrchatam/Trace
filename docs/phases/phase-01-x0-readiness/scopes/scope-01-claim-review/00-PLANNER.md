# P01 / S01 / 00-PLANNER — Claim/Evidence/Review promotion

## Metadata
- id: P01-S01-00
- todo_ids: [P01-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-claim-review.md` for **Claim/Evidence/Review promotion** against live repo + Phase 01 locks. Lock defaults so implementer does not re-debate DONE policy. No product code.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/H_VERIFICATION_STRATEGY.md](../../../../init/H_VERIFICATION_STRATEGY.md) — authority matrix
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G2/G3/G14
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md)
- Live: `internal/domain/claim_stub.go`, `task_state.go` DONE stub, `store/schema/001_init.sql` `reviews` table
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (what S01 must replace)
| Item | Today |
|------|-------|
| Claim/Evidence | `CreateClaim` / `CreateEvidence` / `LinkClaimEvidence` (`claim_has_evidence`) — stubs only |
| DONE gate | `AllowDoneWithoutReview \|\| len(EvidenceIDs)>0` — bypass via flag still legal |
| Review | Table + `EntityReview` const; **no** CreateReview / PASS|FAIL / promote API |
| CLI | `transition` wires domain; no review-specific commands yet |

## Planner work
- Re-read live domain + H authority matrix; lock promotion rules in `01-*`.
- Thicken `01-claim-review.md` exit criteria + paths + skills/MCPs (enough to run alone).
- Light-update **upcoming** S02/S03 stubs if public DONE/Review surface changes expectations.
- Sync SCOPE-TODOS.md + board Notes.

## Depends-on
- Phase 00 complete (P0-X 7/7). First product scope of Phase 01.

## Exit criteria
- [x] `01-*` runnable without guessing (locked DONE/Review policy)
- [x] SCOPE-TODOS + TODO.md Notes updated
- [x] No product code
- [x] Upcoming S02 honesty stub aware of promotion surface

## Minimal todos
- [x] Inspect claim_stub + TransitionTask + reviews schema
- [x] Thicken 01 prompt + light S02 Depends
- [x] Sync todos / board Notes
