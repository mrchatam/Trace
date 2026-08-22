# P22-S09-00 — Planner: harness agent catalog + routing recommendations

## Metadata
- id: P22-S09-00
- todo_ids: [P22-S09-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents, brainstorming]
- mcps: [Read, Glob, Grep]
- verification: automated

## Objective

Lock S09 (human-promoted post-scaffold). Trace **recommends** harness agents/subagents; Trace **never** spawns or orchestrates agents (**E01–E04**). Board: status + notes only. **No product Go this row.**

## References

- [README.md](../../README.md) — enhancement matrix E01–E04
- [DECISION-LOG.md](../../DECISION-LOG.md) — D-22-25…D-22-30
- [WORK-MAP.md](../../WORK-MAP.md) — W-36…W-40
- Existing capability graph: `internal/domain/capability.go`, mig **010**
- Agent-loop protocol fresh-subagent review rule: `docs/rules/agent-loop-protocol.md`

## FINAL locked defaults

| Item | Value |
|------|-------|
| Trace role | **Recommend only** — emit `harness_recommendations[]` in loop next / context; optional `trace agents recommend`. **No** Task tool invocation, no subprocess agent runner |
| New kind | Extend `capabilities.kind` with **`AGENT`** (harness agent profile slug e.g. `agent:performance-reviewer`) |
| Schema | Mig **027** `harness_agents.sql`: `harness_agents`, `harness_agent_requirements` (agent → required SKILL/MCP/TOOL/HOOK capability slugs), optional `harness_agent_routing` (phase/tag keywords). Compat ceiling **27** after S09-01 |
| Bundled catalog | Committed `trace/agents/default.json` (versioned); `trace install agents` upserts into `.trace/`; idempotent with `trace init` optional `--agents-defaults` |
| Routing signals | **Deliberation phase** (CRITIQUE→code-reviewer, VERIFY→performance-reviewer when task tags match, etc.), **task title/tags**, **required capabilities** on task, **harness:subagent** availability (declare via `trace_capability` or bundled HOOK) |
| Subagent hint | When phase is CRITIQUE or post-EXECUTE review and `harness:subagent` is AVAILABLE: recommend **fresh subagent** with named profile + prompt stub — harness executes |
| Default profiles (minimum) | `code-reviewer`, `performance-reviewer`, `security-reviewer`, `nested-reviewer`, `explore` (investigation), `generalPurpose` (fallback). Map each to required skills/MCPs/hooks as metadata — not enforced at runtime by Trace |
| MCP | Add `trace_agents` actions `list`, `recommend`, `describe` (stdio only). Update `TestToolNamesRegistered` |
| Future host | `registry_source`, `registry_version`, `external_url` columns nullable — **no network fetch in P22**; document extension point in `trace/agents/README.md` |
| Hard out | Trace as harness, hosted agent runner, daemon, spawning subagents, ML routing |

## Named tests (product rows)

| Test | Row |
|------|-----|
| `TestHarnessAgentCatalogMigrate027` | S09-01 |
| `TestRecommendAgentForPhaseCritique` | S09-05 |
| `TestRecommendPerformanceReviewerForPerfTask` | S09-05 |
| `TestRecommendSubagentWhenAvailable` | S09-05 |
| `TestInstallAgentsSeedsDefaults` | S09-03 |
| `TestMCPAgentsRecommend` | S09-07 |

## Exit criteria

- [x] 01–08 thickened with exit criteria mapping E01–E04
- [x] VERIFY (S08-07) updated to run after S09 (board order)
- [x] No product Go this row

## Next

**P22-S09-01**
