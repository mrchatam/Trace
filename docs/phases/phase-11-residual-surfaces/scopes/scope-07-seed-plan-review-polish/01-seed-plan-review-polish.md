# P11 / S07 / 01 — Seed / plan / review show polish

## Metadata
- id: P11-S07-01
- todo_ids: [P11-S07-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-28, DF-30, DF-33, DF-45, DF-46** per sibling **00-PLANNER** FINAL locks (2026-08-16). Seed link aliases; plan show snake_case + empty phases/`tasks`; review get|show|list; thin handoff help SoT. **No new migration. No new MCP tools. Gate C `dry_run:false` untouched.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-28, DF-30, DF-33, DF-45, DF-46
- [experiments/_post_p10/bughunt/df33_seed/](../../../../../experiments/_post_p10/bughunt/df33_seed/)
- [experiments/_post_p10/BUGHUNT.md](../../../../../experiments/_post_p10/BUGHUNT.md)
- Live: `cmd/trace/{seed,plan,review,help,tasks}.go`; `internal/planner/{types,service}.go`; `internal/domain/review.go`; `internal/store` reviews + links
- Prior: P11-S06 no seed/plan/review coupling
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate. If 00-PLANNER is still DRAFT, stop and return to planner.

## Locked defaults (FINAL — P11-S07-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `cmd/trace` (seed/plan/review/help); `internal/planner` as needed; thin domain/store list helpers |
| Migration | **None** |
| DF-33 | `from`/`to` **or** `from_id`/`to_id`; clear alias-aware empty error |
| DF-30 | `phases` always `[]`; include goal `tasks` rows |
| DF-46 | Plan show snake_case keys (DF-32 parity) |
| DF-45 | CLI `review get\|show\|list [--task]`; snake_case; no MCP |
| DF-28 | Thin help/docs SoT — task body + Trace-pull; no entity |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false`; P11-S01…S06 |
| Forbidden | Mig; handoff entity; auto-phases; new MCP tools; daemon/HTTP; rewrite `done` history |

## Extension points (exact files)

| File | Work |
|------|------|
| `cmd/trace/seed.go` | Alias resolve for link endpoints; clearer empty error |
| `cmd/trace/seed` tests / `cli_test.go` | DF-33 alias + message tests |
| `cmd/trace/plan.go` (+ planner types/service) | Snake_case show DTO or tags; empty `phases`; populate `tasks` |
| `cmd/trace/review.go` | `get` / `show` / `list [--task]`; help usage |
| `cmd/trace/help.go` (+ optional README) | Review subcommands; DF-28 handoff SoT |
| `internal/domain` / `internal/store` | Thin `ListReviews` / list-by-task if missing (G19 stays in domain) |

## Role work

1. TDD: seed JSON with only `from_id`/`to_id` imports a goal_has_task link.
2. Plan show on goal with tasks and no coarse plan → `phases:[]`, nonempty `tasks`, snake_case keys.
3. Plan show after coarse plan → nested phases/scopes snake_case.
4. Review get/show by id; list and list `--task`; empty list `[]`.
5. Help asserts handoff + context/why guidance (DF-28).
6. Run locked verify suite; board **status + Notes only** (cite test names + DF ids).

## Algorithm sketch (non-normative — locks above win)

```text
seed link endpoints:
  from = from OR from_id; to = to OR to_id
  if both forms set and differ → usage error
  if empty → stderr mentions from/to and from_id/to_id

plan show --goal G:
  view = GetPlan(G)  # phases never nil slice
  tasks = ListTasksByGoalID(G) → snake_case rows
  encode snake_case JSON {goal_id, phases, tasks, …}

review get|show --id:
  GetReview → snake_case object
review list [--task T]:
  all reviews OR reviews linked review_judges_task → T
  encode []
```

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Prefer named asserts: `TestSeedImportFromIDAliases`, `TestPlanShowSnakeCaseAndEmptyPhases`, `TestPlanShowWithPhasesSnakeCase`, `TestReviewGetShowList`, `TestHelpHandoffSoT` (or equiv).

## Exit criteria
- [ ] DF-28, DF-30, DF-33, DF-45, DF-46 green per FINAL locks
- [ ] Carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P11-S07-02**

## Out of scope
- Other Phase 11 scopes; daemon/HTTP/embeddings; handoff entity; rewriting `done` history
