# P22-S08-00 — Planner: agent workflow + phase VERIFY

## Metadata
- id: P22-S08-00
- todo_ids: [P22-S08-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep]
- verification: automated

## Objective

Lock S08. Owned: **C28, C38-MCP, C39**. Then **VERIFY all 43** (README matrix C01–C43). **No product Go.**

## Live inventory (2026-08-18, post-S07)

| Surface | Live state |
|---------|------------|
| Schema max | **026** (`026_eval_rules.sql`; 26 embed sql files) |
| Compat ceiling | **26** (`evals/compat/compat_test.go` — no 027+) |
| Checklist open | **3** bullets: C28 (§11), C38 MCP half (§16), C39 workflow (§16); **138/141** already `[x]` |
| MCP catalog | **13** tools — `RegisteredToolNames`: why, context, add, link, transition, review, tasks, capability, impact, version, search, changes, regressions |
| MCP loop | **Absent** — no `trace_loop`, no `tools_loop.go` |
| CLI loop | Live — `trace loop next\|apply\|status` (`cmd/trace/loop.go` → `internal/loop`) |
| Loop next packet | `trace.loop.next.v1` — planning_evidence, tendencies, similar_changes, risk_hints; **no** `work_conflicts[]` |
| Conflict detection | **Absent** — no `internal/domain/conflicts.go`, no `trace tasks conflicts` |
| Help | `cmd/trace/help.go` lists search, test run, verify, changes, knowledge, loop, git-hook — **partial C39** |
| trace-mcp `-h` | **Stale** — lists only through `trace_changes`; omits `trace_regressions` (and future loop/agents) |
| CONTRIBUTING | **No** “agent workflow” paragraph (index → loop → test → apply → search/why) |
| C38 CLI half | **Closed** (S03) — `trace test run`, `trace verify run`; checklist C38 still `[ ]` until MCP half lands |
| S09 dependency | Mig **027**, `trace_agents`, `harness_recommendations[]` — **must complete before S08-07 VERIFY** |
| Binaries | Rebuild `bin/trace-mcp` after MCP catalog changes (P18 lesson); note in S08-05 evidence, not a residual |

S01–S07 closed — do not reopen schema, eval contract, query tools, or verification cycle in S08 implement rows (except additive loop packet field + MCP mirror).

## References

- [DECISION-LOG.md](../../DECISION-LOG.md) — D-22-05 (MCP may grow), D-22-13 (advisory overlap), D-22-24 (spawn policy), D-22-21 (schema sequence)
- [WORK-MAP.md](../../WORK-MAP.md) — W-32…W-35
- [README.md](../../README.md) — C28, C38, C39 rows; completion bar §8–9
- [VERIFY.md](../../VERIFY.md) — gate index; S09 before S08-07
- Live: `internal/loop/{next,apply}.go`, `internal/mcp/server.go`, `cmd/trace/loop.go`, `cmd/trace-mcp/main.go`

## FINAL locked defaults

| Item | Value |
|------|-------|
| MCP loop | `trace_loop` action=`next\|apply\|status` mirroring CLI (G19 → `internal/loop`). **`apply`**: `DestructiveHint=true`, `ReadOnlyHint=false`; **`next`/`status`**: read-only |
| MCP catalog | S08-01 → **14** tools (`trace_loop` after `trace_regressions`); S09 → **15** (`trace_agents`); update `RegisteredToolNames`, `BuiltinMCPCapabilitySpecs`, `TestToolNamesRegistered`, `cmd/trace-mcp` `-h` each row |
| Conflict | **Advisory only** — overlapping active tasks sharing `change_paths` or impact seeds; same `goal_id` + normalized similar title. CLI `trace tasks conflicts [--task <id>]`. Loop next `work_conflicts[]` cap **8**. **No** distributed locks, no auto-block |
| Active task set | `work_state` ∈ {PENDING, IN_PROGRESS, AWAITING_REVIEW, BLOCKED} — excludes DONE, SKIPPED, STALE, FAILED |
| Path overlap | Repo-relative prefix match on union of paths from each task’s changes (`ListChangesByTaskID` → `ListChangePaths`); OPEN change status preferred when present |
| Impact seed overlap | Intersect `ImpactSeedRef` sets from `seedsFromChangePaths` (file/symbol entity_id match) |
| Title redundancy | Same `goal_id`; normalized titles equal **or** mutual substring (min rune length **8** after normalize) |
| Workflow | Help lists search/test/verify/changes/knowledge/loop/git-hook; CONTRIBUTING “agent workflow” paragraph; `trace_version` still required after MCP rebuild |
| VERIFY | Re-read checklist + README C01–C43; every `[x]` with evidence or in-phase spawn; **E01–E04** evidence in Notes; DR-HANDOFF `no successor` only if fully `[x]` |
| Schema | Max **027** after S09 (022–027); **no 028+** at VERIFY |
| Board order | S08-01…S08-06, then **S09 entire**, then S08-07 VERIFY + S08-08 review |
| Rebuild | After MCP changes, note `bin/trace-mcp` rebuild in board Notes — S08-05 may rebuild binaries in evidence |

## Work map (owned)

| Work ID | Capability | Board rows |
|---------|------------|------------|
| W-32 | C38 MCP `trace_loop` + C39 loop half | S08-01, S08-02 |
| W-33 | C28 conflict / redundant work | S08-03, S08-04 |
| W-34 | C39 help/docs/catalog completeness | S08-05, S08-06 |
| W-35 | VERIFY C01–C43 + DR-HANDOFF | S08-07, S08-08 |

## Named tests (product rows)

| Test | Row |
|------|-----|
| `TestMCPLoopNext` | S08-01 |
| `TestMCPLoopApply` | S08-01 |
| `TestMCPLoopStatus` | S08-01 |
| `TestDetectOverlappingOpenTasks` | S08-03 |
| `TestRedundantSimilarTitleSameGoal` | S08-03 |
| `TestNoConflictWhenTasksDisjoint` | S08-03 |
| `TestTasksConflictsCLI` | S08-03 |
| `TestLoopNextIncludesWorkConflicts` | S08-03 |
| `TestHelpIncludesSearchTestVerify` | S08-05 |

## Spawn / VERIFY policy (locked)

- Implementers: board **status + notes only**; may unbox checklist caps when evidence exists.
- Reviewers: unmet owned bullet → spawn `P22-Sxx-0Na` + `0Nb` immediately below; **do not** close with “later”.
- S08-07: any C01–C43 or E01–E04 without evidence and without runnable in-phase spawn → **FAIL** + spawn `07a`/`07b`.
- S08-08: only row that may set DR-HANDOFF **CLOSED** + **`no successor`** when all 43 matrix rows and checklist are `[x]`.
- Forbidden close residual: checklist `[ ]` for C28/C38/C39 or any README matrix row still open.

## Residual risks for S08-01 (carry forward)

| Risk | Mitigation |
|------|------------|
| MCP handler duplicates loop logic | G19: call `loop.BuildNextPacket`, `loop.Apply`, `loop.Status` — grep handlers for SQL/domain forks |
| Apply destructive annotation wrong | Match `trace_transition`/`trace_review` pattern: `DestructiveHint: true` on apply action only |
| Catalog drift (14 vs specs vs `-h`) | Single edit pass: `server.go`, `capability.go`, `mcp_test.go`, `cmd/trace-mcp/main.go` |
| Assert/gating missing on new tool | `assertMCPToolAllowed(ctx, st, "trace_loop")` on every action entry |
| C38 boxed before MCP review | S08-01 unboxes only on test PASS; S08-02 boxes C38 when both halves evidenced |
| Binary stale after register | Board Notes: rebuild command; S08-05 verifies live `-h` |

## Exit criteria

- [x] 01–08 thickened
- [x] VERIFY spawn/residual policy explicit
- [x] No product Go

## Next

**P22-S08-01**
