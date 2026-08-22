# P06 / 00-PHASE-PLANNER — Environment/capability phase scaffold

## Metadata
- id: P06-00
- todo_ids: [P06-00]
- role: planner
- skills: [planning-and-task-breakdown, grilling]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Light replan of **Phase 06 — Environment / capability graph** after Phase 05 Gate F VERIFYs green and S03-02 closes DR-HANDOFF. Confirm scope order, refresh stub prompts (00/01/02), lock phase defaults from live repo + `A_PROJECT_PLAN` Phase 6, and sync board rows. **Do not** implement product Go. **Do not** deep-finalize every implement prompt — each scope’s `00-PLANNER` does that.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [phase README](./README.md)
- [docs/init/A_PROJECT_PLAN.md](../../init/A_PROJECT_PLAN.md) Phase 6
- [docs/AGENT_ENVIRONMENT.md](../../AGENT_ENVIRONMENT.md)
- [docs/init/D_DECISION_REGISTER.md](../../init/D_DECISION_REGISTER.md) DR-HANDOFF
- [docs/TODO.md](../../TODO.md)
- Phase 05 VERIFY: [VERIFY-NOTES.md](../phase-05-decision-impact/scopes/scope-03-phase-verify/VERIFY-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner only).

## Prior locks to respect (from Phase 05)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate F prelim | Green via `evals/impact` `TestPlantedImpactConflictsGateFPrelim` — do not weaken |
| Gate G / Gate E / Gate C | Green / Go — do not silently reopen |
| Honesty / p0x / x0 / replan / impact | Keep green |
| GC-03/04 | Deferred unless promoted with measurement |
| Daemon/HTTP/embeddings | Still forbidden as primary |
| VerifiedFact / `plan simulate` | Out unless explicitly promoted with Notes |

## Planner work
1. Inspect live MCP/CLI/task packet surfaces vs capability-selection needs.
2. Confirm or revise S01→S03 order (capability surface → selection/ablation → VERIFY/handoff).
3. Ensure each scope folder has `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md`.
4. Patch **upcoming** stubs only; never edit Phase 05 `done` prompt bodies.
5. Unblock/sync `docs/TODO.md` Phase 06 section; set this row `done` with Notes.
6. Update phase README if order/locks change.
7. Lock ablation harness path (prefer planted eval under `evals/` — do not invent commercial multi-model capability theater without evidence).

## Locked phase defaults (do not weaken — set 2026-08-16)
| Item | Value |
|------|-------|
| Goal | Environment/capability graph (`A_PROJECT_PLAN` Phase 6 / H7) |
| Scope order | S01 capability surface → S02 selection/ablation → S03 VERIFY + Phase 07 handoff |
| Validation gate | Capability-selection ablation — planted under **`evals/capability`** (prefer `TestPlantedCapabilitySelectionAblation` + `schema-capability.json` v1; **S02-00 finalizes** names/path) |
| Package / mig hint | Prefer `internal/domain` + store + `internal/compiler` packet attach; optional MCP id mirror — avoid ontology megastore |
| Migration hint | Additive **`010_*`** (S01-00 locks filename) |
| Phase 07 folder | **`phase-07-performance-ladder`** |
| Review policy | Every scope gets `02-review` before next scope implement |
| Carry-forward bars | Honesty A/B/C; Gate G; Gate E; Gate F; p0x 7/7; Gate C artifacts intact |
| MCP | In-scope as capability **ids** / thin G19 adapter — not primary daemon/HTTP |
| Ablation policy | Planted eval only — no commercial multi-model theater |

## Cross-scope blast radius
- S01 surface thickens S02 ablation harness (catalog + packet attach hooks).
- S02 must not weaken Gate F / Gate G / Gate E / Gate C evidence integrity.
- S03 scaffolds Phase 07 `phase-07-performance-ladder` (or records explicit `no successor`).

## Exit criteria
- [x] Phase README accurate vs live Phase 05 outcomes
- [x] Scopes S01–S03 each have 00/01/02 stubs (+ SCOPE-TODOS.md)
- [x] `docs/TODO.md` Phase 06 lists planner + scope rows; `P06-00` done after Phase 05 S03-02
- [x] No product Go in this row
- [x] Board Notes record locks + next runnable (`P06-S01-00`)

## Minimal todos
- [x] Inventory live MCP/CLI/task-context vs capability gaps
- [x] Write/refresh S01–S03 stub prompts
- [x] Sync TODO.md Phase 06 section
- [x] Update README + mark P06-00 done

## Out of scope
- Product feature implementation
- Re-scoring Gate C / inventing Gate F/G without named tests
- Starting Phase 07 before S03 VERIFY
- Deep-finalizing every implement prompt (scope planners)
