# P08 / 00-PHASE-PLANNER — Ecosystem & hardening phase scaffold

## Metadata
- id: P08-00
- todo_ids: [P08-00]
- role: planner
- skills: [planning-and-task-breakdown, grilling]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Light replan of **Phase 08 — Ecosystem & hardening** after Phase 07 Gate H VERIFYs green and S03-02 closes DR-HANDOFF. Confirm scope order, refresh stub prompts (00/01/02), lock phase defaults from live repo + `A_PROJECT_PLAN` Phase 8, and sync board rows. **Do not** implement product Go. **Do not** deep-finalize every implement prompt — each scope’s `00-PLANNER` does that.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [phase README](./README.md)
- [docs/init/A_PROJECT_PLAN.md](../../init/A_PROJECT_PLAN.md) Phase 8
- [docs/init/D_DECISION_REGISTER.md](../../init/D_DECISION_REGISTER.md) DR-HANDOFF
- [docs/TODO.md](../../TODO.md)
- Phase 07 VERIFY: [VERIFY-NOTES.md](../phase-07-performance-ladder/scopes/scope-03-phase-verify/VERIFY-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner only).

## Prior locks to respect (from Phase 07)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate H | Green via `evals/perf` `TestPlantedPerfLadderGateH` — do not weaken |
| S01 T0 / S02 Go | Keep green |
| Gate F / Gate G / Gate E / Gate C / ablation | Green / Go — do not silently reopen |
| Honesty / p0x / x0 | Keep green |
| 100k/1M ladders / GC-03/04 | Deferred unless promoted with measurement |
| Daemon/HTTP/embeddings | Still forbidden as primary unless explicitly promoted |
| VerifiedFact / `plan simulate` | Out unless explicitly promoted with Notes |
| Full-rebuild-on-any-change | Still forbidden |

## Planner work
1. Inspect live plugin/analyzer/CLI/MCP surfaces vs ecosystem + hardening needs.
2. Confirm or revise scope order (plugin APIs → worktrees / production concerns → VERIFY).
3. Ensure each scope folder has `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md`.
4. Patch **upcoming** stubs only; never edit Phase 07 `done` prompt bodies.
5. Unblock/sync `docs/TODO.md` Phase 08 section; set this row `done` with Notes.
6. Update phase README if order/locks change.
7. Lock compatibility + security checklist path (prefer automated checks; no commercial theater).

## Live inventory (confirmed 2026-08-16)
| Surface | Finding |
|---------|---------|
| Analyzers | `internal/analyzers` — JS/TS/TSX/Python/**Go**; extension via `DetectLanguage` + `extract` switch + `extract_*.go`; **no** versioned plugin iface / contributor registry |
| CLI | Thin `cmd/trace` argv — init/index/add/link/transition/review/impact/capability/plan/seed/why/context; `-C`/`--project` root bind |
| MCP | `internal/mcp` + `cmd/trace-mcp` stdio; six tools (`trace_why`/`context`/`add`/`link`/`transition`/`review`); G19 library-only |
| Store | Embed mig `001`…`010` via `schema/*.sql` + `schema_migrations`; Open→`.trace/trace.db`; **no** backup/export CLI; **no** auth surface |
| Worktrees | Single project-root bind; **no** multi-worktree / concurrent-agent isolation story |
| Evals | `evals/{capability,honesty,impact,p0x,perf,replan,x0}` — **no** `evals/compat` / security checklist harness |
| Gaps | Versioned analyzer contribution API; worktree-safe `.trace` bind; migrate/backup/auth hardening; automated compat+security gate |

## Locked phase defaults (do not weaken — set 2026-08-16)
| Item | Value |
|------|-------|
| Goal | Ecosystem & hardening (`A_PROJECT_PLAN` Phase 8) |
| Scope order | **S01** plugin APIs → **S02** multi-agent worktrees → **S03** production hardening (migrations/backup/auth) → **S04** VERIFY + compat/security checklist |
| Validation gate | Compatibility + security checklist — prefer **`evals/compat`** planted harness; preferred names `TestCompatibilitySecurityChecklist` + `schema-compat.json` v1 + temp `metrics-compat.json` (**S04-00 finalizes**) |
| Package hint | Prefer **versioned analyzer extension points** (document + stabilize DetectLanguage/extract contract) over premature plugin megastore / dynamic loader |
| Plugin risk | Do **not** lock a bad universal registry early — prefer adapter-shaped + `APIVersion` / contributor doc path (**S01-00 finalizes**) |
| Worktree hint | Prefer safe multi-root / worktree project bind + concurrent-agent fail-closed (**S02-00 finalizes**; not swarm orchestration) |
| Production hint | Prefer migrate hygiene + backup/restore + **local** auth/binding (not cloud SaaS OAuth) (**S03-00 finalizes**) |
| Migration hint | Additive `011_*`+ only when required; keep embed `schema/*.sql` pattern |
| Successor | `A_PROJECT_PLAN` ends at Phase 8 — VERIFY records **`no successor`** unless Notes promote follow-on |
| Review policy | Every scope gets `02-review` before next scope implement |
| Carry-forward bars | Gate H + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + Gate C `dry_run:false` |
| Still forbidden | Daemon/always-on HTTP/embeddings as primary; full-rebuild-on-any-change; VerifiedFact / `plan simulate` unless promoted |

## Cross-scope blast radius
- S01 versioned analyzer surface thickens contributor path; must not break Gate H Go golden / T0 / incremental.
- S02 worktree bind must not fork store schema carelessly or weaken G19 / CLI `-C` contracts.
- S03 migrations/backup/auth must preserve embed-mig idempotency and Gate C / honesty evidence integrity.
- S04 re-proves checklist + carry-forward; scaffolds **no successor** (or explicit follow-on only with Notes).

## Exit criteria
- [x] Scope order + stubs runnable; board Notes; next row ready
- [x] No product Go on this planner row

## Minimal todos
- [x] Re-read Phase 07 VERIFY-NOTES + live inventory
- [x] Confirm scopes S01→S04; thicken stubs; sync TODO.md
- [x] Mark `P08-00` done with Notes
