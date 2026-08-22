# P10 / 00-PHASE-PLANNER — Integrity surfaces phase scaffold

## Metadata
- id: P10-00
- todo_ids: [P10-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Scaffold **Phase 10** after Phase 09 closed with `no successor`, driven by **deduped** DF-17+ from natural re-run, D09b/D30b/D34, and product bughunt. Confirm scope order, stub prompts, lock defaults, sync `docs/TODO.md` + `AGENTS.md` + findings. **Do not** implement product Go in this row.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md)
- [phase README](./README.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../experiments/DOGFOOD-FINDINGS.md) — canonical DF-17+ + collision map
- [experiments/NATURAL-RERUN.md](../../../experiments/NATURAL-RERUN.md)
- [experiments/BATCH-D09b-D30b-D34.md](../../../experiments/BATCH-D09b-D30b-D34.md)
- Phase 09 closeout: [REVIEW-NOTES.md](../phase-09-dogfood-hardening/scopes/scope-04-phase-verify/REVIEW-NOTES.md)

## Prior locks to respect
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gates C/E/F/G/H + ablation + compat | Green — do not weaken |
| Daemon/HTTP/embeddings | Still forbidden as primary |
| Full-rebuild-on-any-change | Forbidden (DR-INCREMENTAL) |
| Forward-only board | Do not rewrite Phase 09 `done` prompts; Phase 10 is new |
| G19 | CLI/MCP adapters never fork domain logic |
| Law 4 / 9 | Retrieved text ≠ policy; user decisions authoritative — DF-27 must reconcile labeling without elevating untrusted blobs |

## Problem statement (dogfood)

Agents win G1 when Trace is discoverable, but product integrity fails in four clusters:

1. **Retrieval lies / floods** — every task pack gets global discovery↔plan_change; entity type aliases break why; decisions look “not policy” while tasks treat them as binding; IncludeWhy hides failures.
2. **MCP cold-start gap** — CLI has `tasks`/`capability`; MCP does not; Cursor may keep a stale `trace-mcp` after Phase 09 DF-01 library fix.
3. **Index ghosts** — rename leaves old path/symbols, undermining stale-index dogfood (DF-14 closed as experiment still needs product GC).
4. **Gates are advisory** — PASS authorizes any actor to DONE; sticky PASS after reopen; capability missing never blocks transition; `--allow-done` is quiet.

## Assumptions locked (grill deferred)

| # | Assumption |
|---|------------|
| A1 | **Thin MCP tools** `trace_tasks` + capability mirror (list/missing/require as needed) are in-scope; **not** full plan/impact/index MCP this phase |
| A2 | Operator DONE (DF-17): prefer **explicit operator/allow flag or actor attribute** over soft docs; scope planner picks exact API |
| A3 | Sticky PASS (DF-18): reopen **invalidates** prior PASS for DONE (require new review) — not “keep forever” |
| A4 | Capability gate (DF-24): **fail-closed block** (or hard warn + opt-in override) on transition when required cap is UNAVAILABLE/UNKNOWN — S04 planner picks; default block |
| A5 | DPC attach (DF-19): attach only edges **causally reachable** from the queried entity (or goal-scoped) — not “all project DPC” |
| A6 | DF-28/30/33/34/35/36 stay out of S01–S04 unless a scope planner promotes with board spawn |
| A7 | VERIFY default DR-HANDOFF = **`no successor`** |

## Dogfood-driven backlog (import)
| DF | Priority | Phase home |
|----|----------|------------|
| DF-19 DPC pollution | P0 | S01 |
| DF-23/25/27/29 why fidelity | P0–P1 | S01 |
| DF-21/22 MCP + stale process | P0 | S02 |
| DF-32 JSON case | P2 | S02 |
| DF-20 index GC | P0 | S03 |
| DF-17/18/24/26/31 gates | P0–P1 | S04 |

## Scope order (locked)
1. **S01 retrieval-why-fidelity** — DF-19, DF-23, DF-25, DF-27, DF-29  
2. **S02 mcp-parity-install** — DF-21, DF-22, DF-32  
3. **S03 index-gc** — DF-20  
4. **S04 operator-capability-gates** — DF-17, DF-18, DF-24, DF-26, DF-31  
5. **S05 VERIFY** — regressions + carry-forward gates + DR-HANDOFF

## Live inventory (2026-08-16)
| Surface | Finding |
|---------|---------|
| `internal/retrieval/expand.go` | `discoveryPlanChangeHits` global attach residual |
| `internal/retrieval/exact.go` | `plan_change` underscore; no `capability` case |
| `internal/domain/task_state.go` | PASS→DONE any actor; sticky PASS; no cap check |
| `internal/mcp` | Six tools; no tasks/capability |
| `cmd/trace/index.go` | No delete-on-missing path after rename |
| `trace install cursor` | Writes mcp.json; does not force MCP process restart |

## Exit for this planner row
- [x] Phase folder + README + DR-HANDOFF stub  
- [x] Scope stubs S01–S05 with 00/01/02 + SCOPE-TODOS  
- [x] Board Phase 10 section + P10-00 done Notes  
- [x] AGENTS.md + PROJECT_DOCS_INDEX + DOGFOOD-FINDINGS canonical merge  
- [ ] Product Go — **not** this row  

## Next
**P10-S01-00** (scope planner for retrieval-why-fidelity).
