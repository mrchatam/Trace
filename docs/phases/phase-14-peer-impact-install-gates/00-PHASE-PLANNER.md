# P14 / 00-PHASE-PLANNER — Peer impact + install gates (thin) phase scaffold

## Metadata
- id: P14-00
- todo_ids: [P14-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Scaffold **thin Phase 14** after Phase 13 closed with `no successor`, driven by human-approved promotion of goals-gap recommended step **#1** ([`TRACE-GOALS-PROGRESS-2026-08-17.md`](../../research/TRACE-GOALS-PROGRESS-2026-08-17.md) §4): **peer impact + install gates** = research FUTURE **S03 + S04** (ranks **6** + **4–5**). Confirm scope order, stub prompts, lock defaults, sync `docs/TODO.md` + `AGENTS.md` + `PROJECT_DOCS_INDEX.md` + research schedule notes. **Do not** implement product Go in this row. **Do not** board supersession (S05), `plan simulate`, or dogfood D21+.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md) — Laws 5–7, 9, 12–13, 17–19; DR-INCREMENTAL; DR-NOIMP commercial block context
- [phase README](./README.md)
- [TRACE-GOALS-PROGRESS-2026-08-17.md](../../research/TRACE-GOALS-PROGRESS-2026-08-17.md) — §4 recommended sequence #1
- [SIMILAR-PROJECTS-REVIEW-2026-08-16.md](../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md) — ranks 4–6 accept; FUTURE S03–S04; §D rejects
- Phase 13 closeout: [DR-HANDOFF.md](../phase-13-import-resolve-honesty/DR-HANDOFF.md) — historical `no successor` (intact)
- Design: [DECISION_IMPACT.md](../../DECISION_IMPACT.md), [AGENT_ENVIRONMENT.md](../../AGENT_ENVIRONMENT.md), [EVALUATION.md](../../EVALUATION.md) H4/H7

## Prior locks to respect
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gates C/E/F/G/H + ablation + compat | Green — do not weaken |
| Daemon/HTTP/embeddings/Neo4j | Forbidden as primary / SoT |
| Full-rebuild-on-any-change | Forbidden (DR-INCREMENTAL) |
| Forward-only board | Do not rewrite Phase 13 `done` prompts; Phase 14 is new |
| G19 | CLI/MCP adapters never fork domain logic |
| Law 5 / 9 / 17 | No YOLO/AllowAll defaults; inferred ≠ verified |
| Phase 05 impact | Gate F prelim + `trace impact` stay green; extend walks — do not replace commercial-block honesty |
| Phase 06/09/11 capability + install | Catalog + `trace install cursor` stay green; harden matrix/audit — do not dump MCP tools |
| Phase 12–13 honesty | Edge provenance + packet honesty + DF-60…67 dispositions stay green; **do not re-open** closed DFs |

## Problem statement (post–P13 product gap)

Phases 00–13 closed core graph, progressive plan/replan, review honesty, import-resolve, and peer honesty ranks 1–3. Remaining high-value peer methods still deferred:

1. **Impact walks shallow** — Phase 05 impact exists, but multi-seed BFS + contains-asymmetric radius (research rank **6** / CBM + codegraph) are not wired — H4 depth incomplete.
2. **Install / capability gates thin** — Cursor install + capability catalog exist, but marker-gated multi-client detect/install/uninstall registry (rank **4**) and graduated allowlist + durable tool-decision audit (rank **5**) are not productized — H7 depth incomplete.

Thin Phase 14 closes those two surfaces only.

## Assumptions locked (grill deferred — human already cut scope)

| # | Assumption |
|---|------------|
| A1 | **Board only S01 + S02 + VERIFY (S03)** — research FUTURE S05 supersession, ranks 7–20, `plan simulate`, D21+ **not** boarded |
| A2 | Scope planners pick exact APIs/migrations; phase planner locks **homes** + order + non-goals + research citations |
| A3 | Impact upgrades wire into **existing** impact domain / `trace impact` — no new MCP impact tool menu |
| A4 | Install matrix stays CLI/rules-first; optional thin MCP only if library already exposes parity (G19); **no** YOLO / AllowAll |
| A5 | Contains asymmetry: walk deps in; expand containers out via contains; never climb contains-up into siblings (codegraph method) |
| A6 | Multi-seed BFS: one walk, seed exclusion, hop risk, loud truncation (CBM method) — align with Phase 12 packet honesty loudness |
| A7 | VERIFY default DR-HANDOFF = **`no successor`** (VERIFY may promote S05 only with explicit Notes + human) |
| A8 | DF-60…67 remain closed — no re-litigation |

## Research citations (in-scope steals)

| Rank | Technique | Source | Trace home |
|-----:|-----------|--------|------------|
| 6 | Multi-seed impact BFS + depth-bounded contains asymmetry | codebase-memory-mcp, codegraph | **S01** |
| 4 | Conditional / marker-gated install matrix + detect→install→uninstall | codebase-memory-mcp, codegraph | **S02** |
| 5 | Graduated capability allowlist + durable tool-decision audit (≠ chat) | agentrq | **S02** |

## Explicit deferrals (not this phase)

| Research / goals item | Why deferred |
|-----------------------|--------------|
| Ranks 9–10 / FUTURE S05 supersession + episodes | Goals sequence **#2** — later thin phase |
| `plan simulate` / adopt-discard branches | Goals sequence **#3** — after impact walks trustworthy |
| Dogfood D21–D23 / D31 / D33 | Goals sequence **#4** — off-board unless DF-* forward |
| Ranks 7–8, 11–20; Gate D harness as product | Medium / dogfood-first; not #1 cut |
| §D rejects (MCP/daemon/HTTP P0, Neo4j, embeddings default, Hybrid LSP, …) | Hard boundaries |
| Closed DF-60…67 | Forward-only; do not reopen |

## Scope order (locked)
1. **S01 impact-walks** — multi-seed + contains asymmetry on existing impact surface  
2. **S02 install-capability-gates** — marker install registry + graduated allowlist audit  
3. **S03 VERIFY** — S01/S02 regressions + carry-forward gates + DR-HANDOFF  

## Non-goals
- Product Go in this planner row  
- Boarding S05 / ranks 7+ / `plan simulate` / entire D21+ ladder  
- MCP/daemon/HTTP as P0; embeddings / Neo4j SoT  
- Full-rebuild indexer; rewriting Phase 00–13 `done` history  
- Re-opening DF-60…67  
- Claiming Mode-B Gate C / new Gate C from this phase  

## Acceptance (this row)
- Phase folder + README + DR-HANDOFF stub  
- Scope stubs S01–S03 with 00/01/02 + SCOPE-TODOS  
- Board Phase 14 section; **P14-00 done**; next **P14-S01-00**  
- AGENTS.md + PROJECT_DOCS_INDEX + research schedule notes  

## Exit for this planner row
- [x] Phase folder + README + DR-HANDOFF stub  
- [x] Scope stubs S01–S03 with 00/01/02 + SCOPE-TODOS  
- [x] Board Phase 14 section + P14-00 done Notes  
- [x] AGENTS.md + PROJECT_DOCS_INDEX + research schedule notes  
- [ ] Product Go — **not** this row  

## Next
**P14-S01-00** (scope planner for impact-walks / research rank 6).
