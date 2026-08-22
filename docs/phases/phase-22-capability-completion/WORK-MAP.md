# Phase 22 — Work map

Source: 43 unchecked bullets in [`docs/CAPABILITIES_CHECKLIST.md`](../../CAPABILITIES_CHECKLIST.md) + audit findings in the human promotion prompt. Attribution: [`DECISION-LOG.md`](DECISION-LOG.md). Coverage: [`README.md`](README.md) matrix.

| Work ID | Audit / capability | Scope | Board rows |
|---------|--------------------|-------|------------|
| W-01 | C03 tests + `validates` edges; tree-sitter test discovery | S01 | S01-01, S01-02 |
| W-02 | C01 artifact relationships (files, modules, components, functions, types, APIs) | S01 | S01-03, S01-04 |
| W-03 | C02 architectural relationships and boundaries | S01 | S01-05, S01-06 |
| W-04 | C07 impact walk surfaces affected tests | S01 | S01-07, S01-08 |
| W-05 | C25 + DF-86 `trace install git-hook` (index, not wrap commit) | S02 | S02-01, S02-02 |
| W-06 | C04 graph sync honesty + incremental index-on-hook | S02 | S02-03, S02-04 |
| W-07 | C05 record every meaningful change (VCS-promoted) | S02 | S02-05, S02-06 |
| W-08 | C06 compare project states | S02 | S02-07, S02-08 |
| W-09 | C09 wire `BuildPolicyInputs` Execute/Test/Evaluation/Reflect pending | S03 | S03-01, S03-02 |
| W-10 | C12 + C38 CLI: run/record relevant tests | S03 | S03-03, S03-04 |
| W-11 | C11 + C13 + C36 verification cycle gate + scoring coordination | S03 | S03-05, S03-06 |
| W-12 | C14 invariant constraints + C15 iteration comparison | S03 | S03-07, S03-08 |
| W-13 | *(reserved — S03 planner may thicken; do not leave C11–C15 unowned)* | S03 | — |
| W-14 | C08 predicted vs actual impact | S04 | S04-01, S04-02 |
| W-15 | C16 regression↔change association; evidence-backed `caused` | S04 | S04-03, S04-04 |
| W-16 | C18 record improvements | S04 | S04-05, S04-06 |
| W-17 | C29 + C30 + C37 search/changes/history CLI+MCP | S05 | S05-01, S05-02 |
| W-18 | C17 + C31–C34 query tests/failures/successes/regressions | S05 | S05-03, S05-04 |
| W-19 | C35 + C42-surface context/MCP evidence (evals, reflections) | S05 | S05-05, S05-06 |
| W-20 | *(C38 MCP half is W-32)* | S08 | — |
| W-21 | *(FTS already in store — S05 wires CLI/MCP)* | S05 | S05-01 |
| W-22 | MCP parity for query tools (G19) | S05 | S05-01, S05-03 |
| W-23 | C19 + C20 patterns + similar-change query | S06 | S06-01, S06-02 |
| W-24 | C10 + C21 + C26 + C27 project knowledge synthesis | S06 | S06-03, S06-04 |
| W-25 | C22 + C23 + C24 tend-help/hurt + successful approaches + evidence decisions | S06 | S06-05, S06-06 |
| W-26 | Seed export of knowledge rows | S06 | S06-03 (additive keys) |
| W-27 | No ML (D-22-11) | S06 | all S06 reviews |
| W-28 | C40 + C43 evaluator contract (multiple mechanisms, additive) | S07 | S07-01, S07-02 |
| W-29 | C41 project-specific evaluation rules | S07 | S07-03, S07-04 |
| W-30 | C42 library: eval results queryable for future agents | S07 | S07-05, S07-06 |
| W-31 | Seed pointer / committed `trace/eval-rules.json` | S07 | S07-03 |
| W-32 | C38 MCP `trace_loop` + C39 workflow (loop half) | S08 | S08-01, S08-02 |
| W-33 | C28 conflict / redundant work | S08 | S08-03, S08-04 |
| W-34 | C39 remaining workflow usefulness (help, docs, packet completeness) | S08 | S08-05, S08-06 |
| W-35 | VERIFY all 43 `[x]` or in-phase spawn; DR-HANDOFF | S08 | S08-07, S08-08 |
| W-36 | E01 recommend subagents for independent review (supports C09) | S09 | S09-05, S09-06 |
| W-37 | E02 phase/task agent routing (performance-reviewer, etc.) | S09 | S09-01, S09-05, S09-07 |
| W-38 | E03 bundled agent catalog + install | S09 | S09-01, S09-03 |
| W-39 | E04 extensible registry hook (future host; no network P22) | S09 | S09-03, S09-07 |
| W-40 | MCP `trace_agents` + CLI (supports C39) | S09 | S09-07, S09-08 |

**Still out of Phase 22 (hard boundary — not leftovers):** hosted MCP, daemon, HTTP, wrap-`git commit`, ML policies, graph DB, Requirement table, committing `.trace/`, rewriting DONE/Review PASS.

**Schema note:** see DECISION-LOG D-22-21. S01 **022**; S02 **023**; S04 **024**; S06 **025**; S07 **026**; S09 **027**. S03/S05/S08 protocol/query unless a scope planner proves a table is required (then they take the next free number **only if** later scopes have not started — prefer no steal: use the reserved slot or spawn).
