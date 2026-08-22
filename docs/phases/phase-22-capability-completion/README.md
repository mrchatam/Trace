# Phase 22 — Capability completion

**Status:** **complete** (closed 2026-08-18). Every remaining unchecked item in [`docs/CAPABILITIES_CHECKLIST.md`](../../CAPABILITIES_CHECKLIST.md) is `[x]` (141/141). **No post-MVP leftovers. No "later". No residuals at phase close** beyond hard-boundary outs.

Phase 21 DR-HANDOFF remains historically **`no successor`** at P21 close. Phase 22 is a **new forward queue** — not a rewrite of P21 board history.

## Why this phase exists

Phase 21 completed cognition **foundation** (portable P20 seed, retrieval/FTS, 14-row SelectNext, promotion, why, tx apply, thin experiments). Live audit of the 141-item checklist found **43 unchecked** capabilities: test→code graph, graph sync, change capture, verification gates, query surfaces, learning, extensible evaluation, MCP/CLI workflow parity.

**Who decided the P21 stops that left these open?** See [`DECISION-LOG.md`](DECISION-LOG.md) — P21 D-16 (no autonomous runner), P17 DF-86 (git-hook deferred), P20/P21 MCP catalog freeze at 10 tools, `BuildPolicyInputs` cycle-flag stubs.

## Scope order (locked at scaffold)

| Scope | Focus | Work map |
|-------|--------|----------|
| S01 | Test→code graph + artifact/architecture relationships + impact-on-tests | W-01…W-04 |
| S02 | Graph sync (git-hook + honesty) + meaningful change capture + state comparison | W-05…W-08 |
| S03 | Deliberation policy wiring + test invocation + verification cycle + scoring coordination | W-09…W-13 |
| S04 | Predicted vs actual impact + regression `caused` + improvements | W-14…W-16 |
| S05 | Query/search/history CLI+MCP + context evidence surfacing | W-17…W-22 |
| S06 | Patterns, project knowledge, tend-to-help/hurt, evidence-based feedback | W-23…W-27 |
| S07 | Extensible evaluation (mechanisms, project rules, additive contract) | W-28…W-31 |
| S08 | Agent workflow (MCP loop, conflict/redundant work, useful workflow) + phase VERIFY | W-32…W-35 |
| S09 | Harness agent catalog + routing recommendations (human-promoted post-scaffold) | W-36…W-40 |

## Hard boundaries (unchanged)

- No hosted MCP / daemon / HTTP on the core product path
- Local CLI + library + **stdio MCP** extensions **are in scope**
- Optional **local git hook** (`trace install git-hook`, DF-86) is in scope — **must not** wrap `git commit`
- Trace **recommends** harness agents/subagents (S09); Trace **never** spawns agents or acts as a harness
- Incremental localized reindex required (full-rebuild-on-any-change remains forbidden)
- Git remains canonical for diffs (Law 1); no source blobs in SQLite
- No ML / embeddings / graph DB (Law 13)
- Do **not** check boxes in `docs/CAPABILITIES_CHECKLIST.md` from planner rows — implementers check with evidence

## Coverage matrix (43 bullets → owner rows)

Every unchecked checklist line is listed. If a capability is split, both owners are named.

| ID | Checklist bullet | Scope | Owner implement rows | Split |
|----|------------------|-------|----------------------|-------|
| C01 | Track relationships between files, modules, components, functions, types, APIs, tests, and other relevant artifacts | S01 | **P22-S01-03** (artifacts); tests via **P22-S01-01** | S01-01 owns tests; S01-03 owns remaining artifact kinds |
| C02 | Track architectural relationships and boundaries | S01 | **P22-S01-05** | — |
| C03 | Represent tests and what they validate | S01 | **P22-S01-01** | — |
| C04 | Keep graph state synchronized with the actual project | S02 | **P22-S02-01** (hook) + **P22-S02-03** (honesty + incremental index) | Hook installs the path; honesty proves lag is visible |
| C05 | Record every meaningful change made to the project | S02 | **P22-S02-05** | VCS-promoted + apply path; not loop-apply-only |
| C06 | Allow comparison between project states | S02 | **P22-S02-07** | — |
| C07 | Identify potentially affected tests | S01 | **P22-S01-07** | Depends on S01-01 `validates` edges |
| C08 | Compare predicted impact with actual impact after implementation | S04 | **P22-S04-01** | — |
| C09 | Support a structured implementation/thought process rather than blindly executing instructions | S03 | **P22-S03-01** | Wire `BuildPolicyInputs` cycle flags |
| C10 | Allow agents to learn from previous engineering decisions | S06 | **P22-S06-03** | — |
| C11 | Require implementations to go through a test/verification cycle | S03 | **P22-S03-05** | Deliberation/status gate; not a rewrite of DONE/Review PASS |
| C12 | Automatically or explicitly run relevant tests after changes | S03 | **P22-S03-03** | Explicit `trace test run` and/or recorded invoke; no daemon |
| C13 | Verify that existing functionality has not regressed | S03 | **P22-S03-05** | Stored test+verification vs prior baseline |
| C14 | Verify relevant architectural/invariant constraints where possible | S03 | **P22-S03-07** | Uses S01 boundaries + S07 mechanisms |
| C15 | Compare results between iterations | S03 | **P22-S03-07** | — |
| C16 | Identify which change is associated with a regression | S04 | **P22-S04-03** | Evidence-backed `caused` when possible; else `correlated` |
| C17 | Make regression history queryable | S05 | **P22-S05-03** | Data from S04; query surface here |
| C18 | Record improvements | S04 | **P22-S04-05** | — |
| C19 | Identify recurring patterns between types of changes and outcomes | S06 | **P22-S06-01** | Deterministic aggregation; no ML |
| C20 | Allow agents to query historical evidence before making similar changes | S06 | **P22-S06-01** | Also exposed via S05-05 planning evidence |
| C21 | Gradually build project-specific engineering knowledge from historical changes | S06 | **P22-S06-03** | — |
| C22 | Help agents understand what tends to improve or damage a particular project | S06 | **P22-S06-05** | — |
| C23 | Surface successful approaches | S06 | **P22-S06-05** | Query CLI also S05-03 |
| C24 | Help agents make decisions based on accumulated project evidence | S06 | **P22-S06-05** | Context packet S05-05 |
| C25 | Continuously update Trace as the project changes | S02 | **P22-S02-01** | Local git-hook + documented `trace index`; no daemon |
| C26 | Accumulate engineering knowledge over time | S06 | **P22-S06-03** | — |
| C27 | Make the system progressively more knowledgeable about the specific project | S06 | **P22-S06-03** | — |
| C28 | Reduce conflicting or redundant work | S08 | **P22-S08-03** | — |
| C29 | Allow agents to ask questions about project history | S05 | **P22-S05-01** | — |
| C30 | Allow agents to ask what changed | S05 | **P22-S05-01** | Uses S02 capture |
| C31 | Allow agents to ask what tests verify something | S05 | **P22-S05-03** | Uses S01 `validates` graph |
| C32 | Allow agents to ask what previously failed | S05 | **P22-S05-03** | — |
| C33 | Allow agents to ask what approaches previously worked | S05 | **P22-S05-03** | — |
| C34 | Allow agents to ask about regressions | S05 | **P22-S05-03** | — |
| C35 | Allow agents to query accumulated evidence when planning new work | S05 | **P22-S05-05** | — |
| C36 | Coordinate testing, verification, and scoring | S03 | **P22-S03-05** | Orchestrates existing outcome kinds |
| C37 | Allow agents to inspect historical evidence | S05 | **P22-S05-01** | — |
| C38 | Allow agents to initiate or participate in verification loops | S03 + S08 | **P22-S03-03** (CLI test/verify) + **P22-S08-01** (MCP `trace_loop`) | Split: CLI invoke vs MCP loop parity |
| C39 | Make Trace useful as part of an agent's normal engineering workflow | S08 | **P22-S08-01** + **P22-S08-05** | MCP loop + help/docs/context completeness |
| C40 | Support multiple verification mechanisms | S07 | **P22-S07-01** | — |
| C41 | Allow project-specific evaluation rules | S07 | **P22-S07-03** | — |
| C42 | Make evaluation results available to future agents | S07 + S05 | **P22-S07-05** (library) + **P22-S05-05** (context/MCP) | Split: store/query vs agent surface |
| C43 | Allow the evaluation system to become more sophisticated without redesigning Trace's core model | S07 | **P22-S07-01** | Additive evaluator contract |

**Count:** 43 rows in this table. VERIFY (P22-S08-07/08, runs **after S09**) re-reads the checklist and fails the gate if any line is still `[ ]` without an in-phase spawned remediation row.

## Enhancement matrix (S09 — harness agents; not separate checklist bullets)

Human-promoted after initial scaffold. Strengthens partial checklist items; VERIFY should note evidence in Notes when these close supporting gaps.

| ID | Enhancement | Strengthens checklist | Scope | Owner rows |
|----|-------------|----------------------|-------|------------|
| E01 | Recommend fresh **subagent** for independent review when `harness:subagent` AVAILABLE | C09 (structured thought process), agent-loop-protocol | S09 | **P22-S09-05**, S09-06 |
| E02 | Route task/phase signals to catalog agents (e.g. perf task → `performance-reviewer`) | C09, C28, C39 | S09 | **P22-S09-01**, S09-05, S09-07 |
| E03 | Bundled default agent catalog mapped via skills/MCPs/hooks/tools | C39 | S09 | **P22-S09-01**, S09-03 |
| E04 | Extensible registry for future trace host (schema + docs; no network in P22) | C39 | S09 | **P22-S09-03**, S09-07 |

**S09 enhancements (closed at P22-S09-08):** [x] **E01** · [x] **E02** · [x] **E03** · [x] **E04**

**S09 rule:** recommendations only — harness executes delegation. Missing `harness:subagent` → honest inline fallback, not silent skip.

## Portable graph

Code graph (files/symbols/`validates`/architecture edges) stays **local** (`.trace/`). Clones rebuild it with `trace index` after `seed import`. Seed JSON continues to omit index blobs (P17 `TestSeedExportOmitsDeniedSurfaces`). This is **not** a residual: C04 closes via incremental index + git-hook + honesty + clone index recipe.

Semantic graph (entities, cognition, changes, outcomes, knowledge, eval rules pointer) remains portable via `trace seed export -o trace/graph.json`. S06/S07 additive seed keys are in-scope when those tables exist.

## Completion bar

After Phase 22, a fresh agent on a bound repo:

1. Indexes tests and can ask what they validate; impact walk lists affected tests.
2. Installs an optional local git-hook that incrementally indexes (and may export seed); without the hook, honesty banners + `trace index` keep the graph honest.
3. Records meaningful VCS-promoted changes; can compare two project states.
4. Gets store-driven EXECUTE/TEST/EVALUATE/REFLECT from `BuildPolicyInputs`; can `trace test run` relevant tests; cannot skip the verification cycle in deliberation.
5. Can query search/changes/regressions/outcomes/history via CLI and stdio MCP; context surfaces evaluations and reflections.
6. Sees similar-change evidence and project-specific tend-to-help/hurt knowledge.
7. Loads project eval rules without changing the core outcome model.
8. Uses `trace_loop` over MCP; is warned about overlapping/redundant open work.
9. Gets **harness agent recommendations** in loop next (e.g. performance-reviewer for perf tasks; fresh subagent hint for CRITIQUE when supported).
10. Can `trace install agents` + `trace agents recommend` to discover bundled catalog profiles.

Board: [`docs/TODO/phase-22.md`](../../TODO/phase-22.md)
