# P05 / 00-PHASE-PLANNER — Decision impact phase scaffold

## Metadata
- id: P05-00
- todo_ids: [P05-00]
- role: planner
- skills: [planning-and-task-breakdown, grilling]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Light replan of **Phase 05 — Decision impact & simulation** after Phase 04 Gate G VERIFYs green and S03-02 closes DR-HANDOFF. Confirm scope order, refresh stub prompts (00/01/02), lock phase defaults from live repo + `A_PROJECT_PLAN` Phase 5, and sync board rows. **Do not** implement product Go. **Do not** deep-finalize every implement prompt — each scope’s `00-PLANNER` does that.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [phase README](./README.md)
- [docs/init/A_PROJECT_PLAN.md](../../init/A_PROJECT_PLAN.md) Phase 5
- [docs/init/D_DECISION_REGISTER.md](../../init/D_DECISION_REGISTER.md) DR-HANDOFF, DR-NOIMP
- [docs/TODO.md](../../TODO.md)
- Phase 04 VERIFY: [VERIFY-NOTES.md](../phase-04-review-depth/scopes/scope-03-phase-verify/VERIFY-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner only).

## Prior locks to respect (from Phase 04)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate G prelim | Green via `evals/honesty` `TestHonestyEscapeRateGateGPrelim` — do not weaken |
| Gate E | Green via `evals/replan` `TestPlantedDiscoveryReplan` — do not weaken |
| Gate C | **Go** — do not silently reopen |
| Honesty / p0x / x0 / replan | Keep green |
| GC-03/04 | Deferred unless promoted with measurement |
| Daemon/HTTP/embeddings | Still forbidden as primary |
| VerifiedFact | Out unless explicitly promoted with Notes |

## Planner work
1. Inspect live Decision / entity_links / planner surfaces vs impact-class needs.
2. Confirm or revise S01→S03 order (impact classes → Gate F prelim → VERIFY/handoff).
3. Ensure each scope folder has `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md`.
4. Patch **upcoming** stubs only; never edit Phase 04 `done` prompt bodies.
5. Unblock/sync `docs/TODO.md` Phase 05 section; set this row `done` with Notes.
6. Update phase README if order/locks change.
7. Lock Gate F harness path (prefer planted eval under `evals/` — do not invent commercial multi-model Gate F in P05-00 without evidence).

## Locked phase defaults (do not weaken — set 2026-08-16)
| Item | Value |
|------|-------|
| Goal | Decision impact & simulation (`A_PROJECT_PLAN` Phase 5) |
| Scope order | S01 impact classes → S02 Gate F prelim → S03 VERIFY + Phase 06 handoff |
| Validation gate | Gate F preliminary — planted precision/recall under **`evals/impact`** (prefer `TestPlantedImpactConflictsGateFPrelim` + `schema-gate-f.json` v1; S02-00 finalizes names) |
| Package / mig hint | Prefer `internal/domain` + store additive **`009_*`**; no second impact stack; no planner fork |
| Impact band sketch | safe \| caution \| high \| destructive \| reversal; uncertainty KNOWN \| LIKELY \| POSSIBLE \| UNKNOWN (S01-00 finalizes) |
| `plan simulate` | Out this phase (roadmap P13) |
| DR-NOIMP | Still bars commercial automated impact engine |
| Phase 06 folder | `phase-06-environment-capability` |
| Review policy | Every scope gets `02-review` before next scope implement |
| Carry-forward bars | Honesty A/B/C; Gate G; Gate E; p0x 7/7; Gate C artifacts intact |
| MCP | Optional; CLI path primary |

## Cross-scope blast radius
- S01 impact surface thickens S02 Gate F planted harness.
- S02 must not weaken Gate G / Gate E / Gate C evidence integrity.
- S03 scaffolds Phase 06 `phase-06-environment-capability` (or records explicit `no successor`).

## Exit criteria
- [x] Phase README accurate vs live Phase 04 outcomes
- [x] Scopes S01–S03 each have 00/01/02 stubs (+ SCOPE-TODOS.md)
- [x] `docs/TODO.md` Phase 05 lists planner + scope rows; `P05-00` done after Phase 04 S03-02
- [x] No product Go in this row
- [x] Board Notes record locks + next runnable (`P05-S01-00`)

## Minimal todos
- [x] Inventory live Decision APIs vs impact gaps
- [x] Write/refresh S01–S03 stub prompts
- [x] Sync TODO.md Phase 05 section
- [x] Update README + mark P05-00 done

## Out of scope
- Product feature implementation
- Re-scoring Gate C / inventing Gate G without named honesty escape-rate test
- Starting Phase 06 before S03 VERIFY
