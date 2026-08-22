# P11 / 00-PHASE-PLANNER — Residual surfaces phase scaffold

## Metadata
- id: P11-00
- todo_ids: [P11-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Scaffold **Phase 11** after Phase 10 closed with `no successor`, driven by **all** still-open / ops / deferred DFs from post–P10 dogfood + MCP + adversarial hunts (severity-agnostic). Confirm scope order, stub prompts, lock defaults, sync `docs/TODO.md` + `AGENTS.md` + findings schedule. **Do not** implement product Go in this row.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md)
- [phase README](./README.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../experiments/DOGFOOD-FINDINGS.md) — Still open / deferred
- [experiments/POST-P10-DOGFOOD.md](../../../experiments/POST-P10-DOGFOOD.md)
- [experiments/POST-P10-MCP.md](../../../experiments/POST-P10-MCP.md)
- [experiments/_post_p10/DOGFOOD.md](../../../experiments/_post_p10/DOGFOOD.md)
- [experiments/_post_p10/MCP-PARITY.md](../../../experiments/_post_p10/MCP-PARITY.md)
- [experiments/_post_p10/BUGHUNT.md](../../../experiments/_post_p10/BUGHUNT.md)
- [experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md](../../../experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md)
- Phase 10 closeout: [REVIEW-NOTES.md](../phase-10-integrity-surfaces/scopes/scope-05-phase-verify/REVIEW-NOTES.md)

## Prior locks to respect
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gates C/E/F/G/H + ablation + compat | Green — do not weaken |
| Daemon/HTTP/embeddings | Still forbidden as primary |
| Full-rebuild-on-any-change | Forbidden (DR-INCREMENTAL) |
| Forward-only board | Do not rewrite Phase 10 `done` prompts; Phase 11 is new |
| G19 | CLI/MCP adapters never fork domain logic |
| Law 4 / 9 | Retrieved text ≠ policy; user decisions authoritative — DF-48 must reconcile dogfood “binding” vs `untrusted_data` without elevating blobs to system |
| Phase 10 product | DF-17…32 stay fixed; do not regress |

## Problem statement (post–P10)

Phase 10 integrity holds under full-tree index / operator DONE / capability gates, but agents still hit:

1. **Partial-path index ghosts** — `trace index <new-path>` after rename leaves old file/symbols (DF-40).
2. **Review / operator residuals** — any PASS authorizes DONE despite sibling FAIL (DF-43); `--as-operator` is a freestanding flag, not identity (DF-44).
3. **Concurrency** — exclusive `trace.lock` makes CLI↔MCP natural mix brittle (DF-47).
4. **Capability ergonomics** — slug re-declare UNIQUE without `--id` (DF-41); `--allow-done` does not bypass missing-caps (DF-51).
5. **Retrieval / trust** — `why symbol` unknown (DF-49); depth≥2 sibling body leak (DF-35); binding vs untrusted conflict (DF-48); no CLI/MCP `discovery_mentions_task` (DF-42).
6. **MCP/install ops** — stale Cursor catalog / reload tips incomplete (DF-22/37/50).
7. **Seed/plan/review polish** — handoff SoT (DF-28), empty plan phases (DF-30), seed `from_id`/`to_id` (DF-33), review list/show (DF-45), plan JSON PascalCase (DF-46).

## Assumptions locked (grill deferred)

| # | Assumption |
|---|------------|
| A1 | **All 18 open/ops/deferred product DFs** are in-scope this phase (severity-agnostic); experiment-only DF-06/07/13/34/36 stay off-board |
| A2 | Scope planners pick exact APIs; phase planner only locks **homes** + order + non-goals |
| A3 | DF-28 handoff SoT may be thin pointer / docs / first-class entity — S07 planner decides; must not invent daemon |
| A4 | DF-47 may be soft UX (clearer ErrLocked / serialize guidance) and/or safer multi-open — **not** drop exclusive lock without security review |
| A5 | DF-44 may stay “conscious flag” with louder docs/warnings **or** tighten identity — S02 planner picks; Gate G hatch retained |
| A6 | Thin MCP expansion only when a DF requires parity (prefer CLI); no plan/impact/index MCP dump unless promoted |
| A7 | VERIFY default DR-HANDOFF = **`no successor`** |

## Canonical open DF list (2026-08-16)

| ID | Sev | Status | Phase home |
|----|-----|--------|------------|
| DF-22 | high | open/ops | S06 |
| DF-28 | med | deferred | S07 |
| DF-30 | low | open | S07 |
| DF-33 | low | deferred | S07 |
| DF-35 | low | open | S05 |
| DF-37 | med | open/ops | S06 |
| DF-40 | med | open | S01 |
| DF-41 | med | open | S04 |
| DF-42 | med | open | S05 |
| DF-43 | med | open | S02 |
| DF-44 | low | open | S02 |
| DF-45 | low | open | S07 |
| DF-46 | low | open | S07 |
| DF-47 | med | open | S03 |
| DF-48 | med | open | S05 |
| DF-49 | low | open | S05 |
| DF-50 | low | open | S06 |
| DF-51 | med | open | S04 |

**Count: 18** (deduped; DF-35 listed once). Collisions already resolved in findings collision map.

## Scope order (locked)
1. **S01 index-partial-path-gc** — DF-40  
2. **S02 review-pass-fail-operator** — DF-43, DF-44  
3. **S03 store-lock-concurrency** — DF-47  
4. **S04 capability-upsert-hatch** — DF-41, DF-51  
5. **S05 retrieval-why-depth-trust** — DF-49, DF-35, DF-48, DF-42  
6. **S06 mcp-install-reload** — DF-22, DF-37, DF-50  
7. **S07 seed-plan-review-polish** — DF-28, DF-30, DF-33, DF-45, DF-46  
8. **S08 VERIFY** — regressions + carry-forward gates + DR-HANDOFF

## Non-goals
- Product Go in this planner row  
- Reopening fixed DF-17…32 as new work (regression only via VERIFY)  
- Rewriting Mode-B Gate C packs / claiming new Gate C  
- Daemon / HTTP / embeddings / full-rebuild indexer  

## Acceptance (this row)
- Phase folder + README + DR-HANDOFF stub  
- Scope stubs S01–S08 with 00/01/02 + SCOPE-TODOS  
- Board Phase 11 section; **P11-00 done**; next **P11-S01-00**  
- AGENTS.md + PROJECT_DOCS_INDEX + DOGFOOD-FINDINGS scheduled → Phase 11  

## Exit for this planner row
- [x] Phase folder + README + DR-HANDOFF stub  
- [x] Scope stubs S01–S08 with 00/01/02 + SCOPE-TODOS  
- [x] Board Phase 11 section + P11-00 done Notes  
- [x] AGENTS.md + PROJECT_DOCS_INDEX + DOGFOOD-FINDINGS schedule note  
- [ ] Product Go — **not** this row  

## Next
**P11-S01-00** (scope planner for index-partial-path-gc).
