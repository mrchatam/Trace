# P17-00 — Plan portable graph via git (FINAL)

> **Forward correction (2026-08-17, after this FINAL):** upcoming S01–S04 **supersede** this prompt’s exclude of planner tables, missing `exported_at_commit`, and “no git hook this phase.” See [`DF-84-FORWARD.md`](DF-84-FORWARD.md) (DF-84…86). This prompt body remains historical FINAL. Do not rewrite the locks below as if they always included the plan tree.

## Metadata
- id: P17-00
- todo_ids: [P17-00]
- role: planner
- skills: [planning-and-task-breakdown, research, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live repo + four investigation loops, lock a **FINAL** thin Phase 17 for clone-readable semantic collaboration (seed export + git convention + idempotent import). **No product Go.** **Do not** edit Phase 16 product Go or rewrite P16 `done` rows. **Do not** steal P16 DFs 68, 70–78.

## References
- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md) — Laws 1, 4, 6, 9, 13, 18, 19
- [phase README](README.md)
- Research: [docs/research/PORTABLE-GRAPH-GIT-2026-08-17.md](../../research/PORTABLE-GRAPH-GIT-2026-08-17.md)
- Cross-check: [docs/research/TRACE-GOALS-PROGRESS-2026-08-17.md](../../research/TRACE-GOALS-PROGRESS-2026-08-17.md), [docs/research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md](../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md)
- Findings SoT: [experiments/DOGFOOD-FINDINGS.md](../../../experiments/DOGFOOD-FINDINGS.md)
- Live: `cmd/trace/seed.go` (import only); `fixtures/x0/seed/gt.json`; `internal/store` Upsert* vs `InsertLink`; `.gitignore` `.trace/`
- P16 S05 (depend, do not duplicate): [../phase-16-assert-root-and-surfaces/scopes/scope-05-seed-impact-packet/00-PLANNER.md](../phase-16-assert-root-and-surfaces/scopes/scope-05-seed-impact-packet/00-PLANNER.md)
- [docs/TODO.md](../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). **Unattended:** human already requested investigation + scaffold; defaults below are FINAL. Phase 16 remains the running phase.

## Live confirmation (2026-08-17)

| Claim | Still true? | Evidence |
|-------|-------------|----------|
| Seed import only | **Yes** | `cmdSeed` requires `args[0] == "import"`; help has no export |
| `.trace/` gitignored | **Yes** | `.gitignore`; project-rules store row |
| Entity create upserts by UUID | **Yes** | `UpsertGoal` ON CONFLICT; `CreateGoal` still appends `entity.created` |
| Link re-import fails | **Yes** | `InsertLink` plain INSERT; UNIQUE on `entity_links` |
| P16 S05 seed completeness pending | **Yes** | DF-70/73 boarded on P16-S05-*; P17 must depend |
| Next free DF before this row | **DF-79 unused** | DOGFOOD-FINDINGS; this phase uses **DF-80…83** |

**Count boarded fix: 4** (DF-80, 81, 82, 83). Next free after this phase: **DF-84**.

## Disposition matrix (FINAL)

See [README.md](README.md). **Boarded:** S01→S03 fix + S04 VERIFY. **Not boarded:** encryption, `.trace/` commit, daemon, P16 S05 work, reviews/DONE in default export.

## Locked defaults (FINAL — phase)

| Item | Value |
|------|-------|
| Phase | Thin portable graph via git (seed JSON v1 snapshot) |
| History | Do not rewrite Phase 00–16 `done` prompts; P16 stays current runnable |
| Product Go | **Forbidden** on P17-00 |
| Format | Seed JSON **v1** round-trip; unknown keys still rejected on import |
| Path | Recommended commit **`trace/graph.json`**; CLI `-o` or stdout |
| `.gitignore` | **Unchanged** (`.trace/` stays ignored) |
| Export include | goals, tasks (id/title/body/goal_id), decisions, assumptions, discoveries, plan_changes, claims, evidence, links (post–P16 S05 rels), findings/alternatives **if** S05 landed keys |
| Export exclude | index, tokens, lock, tool decisions, capabilities, events, reviews, transitions/work_state (default) |
| Import | UUID upsert; duplicate links **no-op**; do not replay default transitions |
| MCP | **No** `trace_seed` / new tools |
| G19 | Thin CLI adapter; library owns export/import |
| Depends-on | **P16 S05** for DF-70/73 keys — do not duplicate |
| DR-HANDOFF intent | After VERIFY: default **`no successor`** |
| Forbidden | Daemon; encryption-as-git; committing `.trace/`; git merge driver; NDJSON split this phase |

## Scope order (locked)
1. **S01 seed-export** — DF-80 (`seed export` + import↔export tests)
2. **S02 commit-convention** — DF-82 (path, gitignore stays, AGENTS/help/CONTRIBUTING)
3. **S03 idempotent-import** — DF-81, DF-83 (upsert + conflict docs)
4. **S04 VERIFY** — named S01–S03 + clone recipe + carry-forward; DR-HANDOFF `no successor`

Board order is sequential (protocol) even though S02 is docs-heavy.

## Non-goals
- Product Go in this planner row
- Committing `.trace/`; encryption; daemon/sync; MCP seed
- Duplicating P16 S05 seed rels / impact / `trace_impact`
- Default export of reviews/DONE/allowlist
- Claiming current focus is Phase 17 while P16 is incomplete
- Auto-boarding research S05 / `plan simulate` / D21+

## Planner work (this row)
1. [x] Four investigation loops → research doc
2. [x] Disposition matrix FINAL (GO)
3. [x] Create scope folders + `00-PLANNER`/`01`/`02` stubs + SCOPE-TODOS + DR-HANDOFF stub
4. [x] Board P17-* **after** P16 rows; P17-00 done; next **runnable remains P16-S02-01**
5. [x] AGENTS.md queued one-liner (P16 current focus); PROJECT_DOCS_INDEX; DOGFOOD-FINDINGS DF-80…83
6. [x] Light P16 S06 note: P17 independently queued; VERIFY still `no successor`

## Exit criteria
- [x] Research doc + GO
- [x] Disposition matrix FINAL for DF-80…83
- [x] Board updated after P16 block
- [x] 00-PHASE-PLANNER marked **FINAL**
- [x] Notes: next **runnable** = **P16-S02-01**; first P17 row after P16 VERIFY = **P17-S01-00**
- [ ] Product Go — **not** this row

## Minimal todos
- [x] Investigate loops 1–4
- [x] Write research doc
- [x] Spawn S01–S04 stubs
- [x] Board + README + index + findings + AGENTS queued note

## Next
Orchestrator: continue **P16-S02-01**. Do **not** start P17-S01-00 until every P16 row is `done`/`skipped`.
