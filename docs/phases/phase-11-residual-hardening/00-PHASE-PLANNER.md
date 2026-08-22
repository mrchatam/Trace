# P11 / 00-PHASE-PLANNER — Residual hardening phase scaffold

## Metadata
- id: P11-00
- todo_ids: [P11-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Scaffold **Phase 11** after Phase 10 closed with `no successor`, driven by **post–P10 residual** DF-40+ (not the fixed DF-17…32 cluster). Confirm scope order, stub prompts, lock defaults, sync `docs/TODO.md` + `AGENTS.md` + findings. **Do not** implement product Go in this row.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md)
- [phase README](./README.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../experiments/DOGFOOD-FINDINGS.md) — open residuals + collision map
- [experiments/POST-P10-DOGFOOD.md](../../../experiments/POST-P10-DOGFOOD.md)
- [experiments/POST-P10-MCP.md](../../../experiments/POST-P10-MCP.md)
- [experiments/_post_p10/DOGFOOD.md](../../../experiments/_post_p10/DOGFOOD.md)
- [experiments/_post_p10/BUGHUNT.md](../../../experiments/_post_p10/BUGHUNT.md)
- [experiments/_post_p10/MCP-PARITY.md](../../../experiments/_post_p10/MCP-PARITY.md)
- [experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md](../../../experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md)
- Phase 10 template: [phase-10-integrity-surfaces/](../phase-10-integrity-surfaces/)
- Phase 10 closeout: [REVIEW-NOTES.md](../phase-10-integrity-surfaces/scopes/scope-05-phase-verify/REVIEW-NOTES.md)

## Prior locks to respect
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gates C/E/F/G/H + ablation + compat | Green — do not weaken |
| Phase 10 DF-17…32 | Stay fixed — do not regress |
| Daemon/HTTP/embeddings | Still forbidden as primary |
| Full-rebuild-on-any-change | Forbidden (DR-INCREMENTAL) |
| Forward-only board | Do not rewrite Phase 10 `done` prompts; Phase 11 is new |
| G19 | CLI/MCP adapters never fork domain logic |
| Law 4 / 9 | Retrieved text ≠ policy; user decisions authoritative — DF-48 must reconcile dogfood wording without elevating untrusted blobs |
| P08 worktree lock | Path-local `.trace` + exclusive `trace.lock` — S03 may soften **concurrency UX**, not abandon exclusivity without explicit lock |

## Problem statement (post–P10)

Phase 10 integrity holds under fresh CLI/MCP, but residuals remain in six clusters:

1. **Honesty hole** — any linked PASS authorizes DONE even with sibling FAIL (DF-43); hatch/`--allow-done` does not bypass missing caps and copy omits that (DF-51); `--as-operator` is still flag≠identity (DF-44).
2. **Partial index GC** — full-tree `trace index` GC works (DF-20); `index <new-path>` after rename leaves ghosts (DF-40).
3. **Store lock** — exclusive `.trace/trace.lock` fails closed under parallel CLI↔CLI or CLI↔MCP (DF-47).
4. **Capability / CLI UX** — declare cannot upsert by slug without `--id` (DF-41); `why symbol` unknown (DF-49); print-only install omits reload tip (DF-50; pairs DF-22/37 ops).
5. **Context / trust / attribution** — depth≥2 sibling body leak (DF-35); TASK “binding” vs MD “not authority” (DF-48); no discovery→task attribution CLI/MCP (DF-42).
6. **Deferred lows** — DF-28/30/33/45/46 stay non-goals unless a scope planner promotes with board spawn.

## Assumptions locked (grill deferred)

| # | Assumption |
|---|------------|
| A1 | **DF-43:** DONE requires **no linked FAIL** among active reviews that could authorize PASS (exact rule finalized at S01-00; default = ignore superseded/UNCERTAIN; FAIL blocks even if another PASS exists) |
| A2 | **DF-51:** keep missing-caps gate; hatch WARNING must mention `--allow-missing-caps`; do **not** silently make `--allow-done` bypass caps unless S01-00 explicitly locks that (default: **document + warn**, two flags stay independent) |
| A3 | **DF-44:** no real auth this phase — clarify Actor≠auth in help/MCP schema; optional louder warning; **not** OAuth/token identity |
| A4 | **DF-40:** path-scoped index must GC ghosts for paths **replaced by rename within the indexed set** (or document argv semantics + implement delete-old when new path indexes after move); prefer incremental, **no** full-rebuild |
| A5 | **DF-47:** prefer **short retry / clearer ErrLocked UX** and/or documented single-writer; do **not** invent multi-writer SQLite or drop exclusivity without S03-00 evidence |
| A6 | **DF-41:** upsert-by-slug (or resolve `--id` from slug) on declare; keep UNIQUE slug |
| A7 | **DF-35/48:** depth sibling leak = fail-closed or task-scoped expand; trust wording = dogfood/docs + MD copy consistency with Law 4/9 (no `system` elevate) |
| A8 | **DF-42:** thin CLI/MCP link rel for discovery→task **or** documented store-only path — S05-00 picks; no full planner redesign |
| A9 | DF-28/30/33/45/46 stay out of S01–S05 unless promoted |
| A10 | VERIFY default DR-HANDOFF = **`no successor`** |

## Dogfood-driven backlog (import)
| DF | Priority | Phase home |
|----|----------|------------|
| DF-43 PASS+FAIL | P0 | S01 |
| DF-51 hatch vs caps | P1 | S01 |
| DF-44 flag≠identity | P2 | S01 |
| DF-40 partial-index ghosts | P0 | S02 |
| DF-47 store lock | P0 | S03 |
| DF-41 slug upsert | P1 | S04 |
| DF-49 why symbol | P2 | S04 |
| DF-50 / DF-22/37 install tip | P2 | S04 (docs/tip only for 22/37) |
| DF-35 depth leak | P1 | S05 |
| DF-48 binding vs trust | P1 | S05 |
| DF-42 discovery→task | P1 | S05 |
| DF-28/30/33/45/46 | — | **non-goal** |

## Scope order (locked)
1. **S01 honesty-review-gates** — DF-43, DF-51, DF-44  
2. **S02 index-partial-gc** — DF-40  
3. **S03 store-lock-concurrency** — DF-47  
4. **S04 capability-cli-ux** — DF-41, DF-49, DF-50 (+ DF-22/37 tip)  
5. **S05 context-trust-attribution** — DF-35, DF-48, DF-42  
6. **S06 VERIFY** — regressions + carry-forward gates + DR-HANDOFF

## Live inventory (2026-08-16)
| Surface | Finding |
|---------|---------|
| `findPassReviewID` / DONE policy | Any PASS authorizes; sibling FAIL ignored → DF-43 |
| Transition hatch | `--allow-done` loud WARNING; missing-caps still blocks → DF-51 |
| `cmd/trace/index.go` | Full-tree GC yes; argv path-only no project GC → DF-40 |
| `internal/store` Open | Exclusive `trace.lock` / `ErrLocked` → DF-47 |
| `capability declare` | UNIQUE slug without upsert-by-slug → DF-41 |
| Exact/Why | `file` OK; `symbol` unknown → DF-49 |
| `install cursor` print | Omits stderr reload tip; `--write` has tip → DF-50 |
| compiler Expand depth | Sibling task body at depth≥2 → DF-35 |
| Decision MD + dogfood TASK | Binding vs untrusted_data conflict → DF-48 |
| `trace link` / MCP | No discovery→task rel → DF-42 |

## Exit for this planner row
- [x] Phase folder + README + DR-HANDOFF stub  
- [x] Scope stubs S01–S06 with 00/01/02 + SCOPE-TODOS  
- [x] Board Phase 11 section + P11-00 done Notes  
- [x] AGENTS.md + PROJECT_DOCS_INDEX + DOGFOOD-FINDINGS schedule merge  
- [ ] Product Go — **not** this row  

## Next
**P11-S01-00** (scope planner for honesty-review-gates).
