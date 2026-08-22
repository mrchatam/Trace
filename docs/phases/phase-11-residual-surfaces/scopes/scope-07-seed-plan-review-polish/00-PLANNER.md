# P11-S07-00 — Seed / plan / review show polish (FINAL)

## Metadata
- id: P11-S07-00
- todo_ids: [P11-S07-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S07 implement/review prompts for **DF-28, DF-30, DF-33, DF-45, DF-46**. Confirm live inventory; lock APIs/tests. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A3 DF-28 thin OK; A6 CLI-primary
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-28, DF-30, DF-33, DF-45, DF-46
- [experiments/_post_p10/BUGHUNT.md](../../../../../experiments/_post_p10/BUGHUNT.md) — DF-30/45/46/33/28
- [experiments/_post_p10/bughunt/df33_seed/RESULTS.txt](../../../../../experiments/_post_p10/bughunt/df33_seed/RESULTS.txt) — `from_id` → empty endpoints → opaque fail
- Phase 10 DF-32 pattern: [../../../phase-10-integrity-surfaces/scopes/scope-02-mcp-parity-install/](../../../phase-10-integrity-surfaces/scopes/scope-02-mcp-parity-install/) — snake_case CLI DTOs
- Live: `cmd/trace/{seed,plan,review,help,tasks}.go`; `internal/planner/{types,service}.go` (`GetPlan` / `PlanView`); `internal/domain/review.go`; `internal/store` reviews + `ListLinksTo`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no grill (A3 thin handoff + live inventory + prefer-no-mig; A1–A7 do not conflict).

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `seedLink` | JSON tags only `from`/`to` — `from_id`/`to_id` silently empty → `goalID and taskID are required` (DF-33; df33_seed exit 2) |
| `cmdPlanShow` | Encodes raw `planner.PlanView` (no json tags) → PascalCase `GoalID`/`Phases` (DF-46) |
| `GetPlan` empty | No phases → `Phases` nil → JSON `null` when only tasks seeded (DF-30) |
| `cmdReview` | Subcommands `create\|set\|residual` only — no get/list/show (DF-45) |
| Domain/store | `GetReview` exists; no CLI list; reviews via `ListLinksTo(task)` + GetReview |
| Handoff SoT | No first-class entity; D30b residual still OK via pointer+context (DF-28) |
| Migration | Prefer **none** — aliases/DTO/help/list need no schema |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-33 | Seed `from_id`/`to_id` opaque fail | Accept aliases; clear error if still empty |
| DF-30 | `plan show` empty/`null` Phases when only tasks seeded | `phases: []` never null; include goal `tasks` list |
| DF-46 | Plan show PascalCase JSON | Snake_case plan show JSON (DF-32 parity) |
| DF-45 | No review get\|list\|show | Thin CLI get/show + list (optional `--task`) |
| DF-28 | No first-class handoff SoT | **Thin** help/docs SoT — task body + Trace-pull; no new entity/mig |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| DF home | **DF-28, DF-30, DF-33, DF-45, DF-46** only |
| Packages | Thin **`cmd/trace`** (`seed.go`, `plan.go`, `review.go`, `help.go`, tests); **`internal/planner`** (GetPlan / view encoding or empty-slice + tasks populate); thin **`internal/domain`** / **`internal/store`** only if list-reviews helper needed. **G19** — no business logic fork in adapters |
| Migration | **None** |
| DF-33 seed aliases | `seedLink` accepts **`from`/`to` or `from_id`/`to_id`** (and hyphen variants only if already used elsewhere — prefer underscore aliases matching finding). If both forms present and differ → usage error. After resolve, empty endpoints → stderr must mention accepted keys (`from`/`to` / `from_id`/`to_id`), not only "goalID and taskID are required" |
| DF-33 non-goals | Do not change link `rel` vocabulary; do not invent new seed top-level keys; discovery-mentions-task stays S05 |
| DF-30 phases | `plan show` JSON **`phases` is always an array** (`[]` when none) — never `null` |
| DF-30 tasks | Same payload includes **`tasks`**: goal-linked tasks via `ListTasksByGoalID` as snake_case rows `{id,title,work_state,goal_id}` (mirror `trace tasks` shape). Present even when phases empty so seeded-only goals are useful |
| DF-30 non-goals | Do not auto-create coarse phases from tasks; do not invent plan MCP |
| DF-46 snake_case | All `plan show` object keys snake_case: `goal_id`, `current_scope_id`, `phases` (nested `id`/`title`/`body`/`ord`/`status`/`scopes`), `current_deep_plan`, `lookahead_scope_id`, `lookahead_summary`, `tasks`. Prefer json tags on planner view types **or** CLI DTO mapping (DF-32 `capabilityListRow` pattern). Nested deep-plan fields snake_case too |
| DF-46 non-goals | Do not reopen DF-32 capability paths; do not change non-show plan create/deep stdout unless same types are shared (if shared, keep tests green) |
| DF-45 CLI | Add **`review get --id`**, **`review show --id`** (alias of get), **`review list [--task <id>]`**. Keep `create\|set\|residual`. Usage/help updated |
| DF-45 JSON | Snake_case rows: at least `id`, `title`, `result`, `status`; include `body` on get/show; list may omit body. Empty list → `[]` |
| DF-45 list semantics | No `--task`: list reviews (store/domain list — add thin `ListReviews` if missing). With `--task`: reviews linked via `review_judges_task` to that task (`ListLinksTo` + GetReview). Order stable (created_at) |
| DF-45 MCP | **CLI-only** this scope (A6) — no new MCP tool required |
| DF-45 non-goals | Do not change DONE/PASS/FAIL gates (S02); residual subcommands unchanged |
| DF-28 posture | **Thin SoT (A3)** — no handoff entity, no mig, no daemon. Product close = **help (+ optional light README)** declaring: handoff source of truth is the **predecessor task body** (plus linked decisions via retrieval); successors **must** Trace-pull with `trace context` / `trace why` (and `trace tasks` as needed) |
| DF-28 test | Assert help (and README if touched) contains assertable handoff guidance (`handoff` + `context` and/or `why`) |
| DF-28 non-goals | No new entity type; no `trace handoff` command required; no MCP tool; do not re-open DF-35 body-redact |
| Tests (required) | (1) **`TestSeedImportFromIDAliases`** (or equiv): link with only `from_id`/`to_id` succeeds. (2) **`TestSeedImportMissingEndpointsMessage`**: empty endpoints mention alias keys. (3) **`TestPlanShowSnakeCaseAndEmptyPhases`**: seeded goal+tasks, no coarse plan → `phases` is `[]` not null; `tasks` nonempty; keys snake_case (`goal_id`, not `GoalID`). (4) **`TestPlanShowWithPhasesSnakeCase`**: after coarse plan, phases/scopes snake_case. (5) **`TestReviewGetShowList`**: create+get/show; list / list `--task`; empty `[]`. (6) **`TestHelpHandoffSoT`** (or help substring in existing help test): DF-28 wording. (7) Carry-forward suites |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` untouched; P10 DF-32 snake_case pattern; P11-S01…S06 |
| Forbidden | New mig; handoff entity/daemon; auto-synthesize phases; new MCP tools; daemon/HTTP/embeddings; full-rebuild indexer; rewriting Phase 00–10 / P11-S01–S06 `done` history; S08 product work beyond Depends note |

## Effects on later scopes
- **S08 VERIFY:** include DF-33 seed aliases, DF-30 empty `phases`+`tasks`, DF-46 snake_case plan show, DF-45 review get/list/show, DF-28 help handoff SoT in evidence table.
- Light Depends note on S08 stubs only (forward).

## Exit
- [x] Thicken `01-seed-plan-review-polish.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light upcoming Depends note (S08)
- [x] Board Notes; next **P11-S07-01**
- [x] Product Go — **not** this row
