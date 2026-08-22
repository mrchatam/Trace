# P09 / 00-PHASE-PLANNER — Dogfood hardening phase scaffold

## Metadata
- id: P09-00
- todo_ids: [P09-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Scaffold **Phase 09** after Phase 08 closed with `no successor`, driven by Cursor dogfood findings. Confirm scope order, stub prompts, lock defaults, sync `docs/TODO.md`. **Do not** implement product Go in this row.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [phase README](./README.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../experiments/DOGFOOD-FINDINGS.md)
- [experiments/LADDER.md](../../../experiments/LADDER.md)
- [experiments/RESULTS.md](../../../experiments/RESULTS.md)
- Phase 08 closeout: [REVIEW-NOTES.md](../phase-08-ecosystem-hardening/scopes/scope-04-phase-verify/REVIEW-NOTES.md)

## Prior locks to respect
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gates C/E/F/G/H + ablation + compat | Green — do not weaken |
| Daemon/HTTP/embeddings | Still forbidden as primary |
| Full-rebuild-on-any-change | Forbidden |
| Forward-only board | Do not rewrite Phase 08 `done` prompts; Phase 09 is new |

## Dogfood-driven backlog (import)
| DF | Priority | Phase home |
|----|----------|------------|
| DF-01 why/context `unknown entity type "review"` | P0 for S01 | S01 |
| DF-02 no task list / status CLI | P1 | S02 |
| DF-03 no `trace install` / Cursor wire | P1 | S03 |
| DF-04 seed path vs `-C` | P2 | S02 |
| DF-05 workspaceFolder footgun | P2 | S03 docs |
| DF-08 missing dogfood H4/H6/H7 | parallel experiments | not S01–S03 code unless promoted |

## Scope order (locked)
1. **S01 retrieval-review** — ExactLookup + why/context include `review`; regression test; fail-soft or full hit — prefer full hit with title/result excerpt  
2. **S02 discoverability** — `trace tasks` (or `trace status`) listing work_state + titles; document seed abs-path  
3. **S03 install-wire** — generate/merge Cursor MCP snippet (`trace-mcp -C ${workspaceFolder}`) from CLI or documented script; no daemon  
4. **S04 VERIFY** — carry-forward gates + DF-01 regression + DR-HANDOFF

## Live inventory (2026-08-16)
| Surface | Finding |
|---------|---------|
| `internal/retrieval/exact.go` | switch misses `review` (and possibly other domain types used in events) |
| CLI | no `tasks`/`status` list command |
| Install | binaries on PATH; MCP hand-edited in `~/.cursor/mcp.json` |
| Dogfood | 6/6 G1 wins; DF-01 blocks post-review why/context |

## Exit for this planner row
- [x] Phase folder + README  
- [x] Scope stubs S01–S04 with 00/01/02 + SCOPE-TODOS  
- [x] Board Phase 09 section + P09-00 done Notes  
- [x] AGENTS.md current focus updated  
- [ ] Product Go — **not** this row  

## Next
**P09-S01-00** (scope planner for retrieval-review).
