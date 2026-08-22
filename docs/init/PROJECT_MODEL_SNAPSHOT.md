# Project Model Snapshot

Internal model established before the implementation plan. Derived from repository documentation as of initialization; not a substitute for `docs/PROJECT_MODEL.md`.

## 1. Actual problem

Coding agents can read repositories and Git history, but they lack a durable, provenance-bearing model of:

- why code exists;
- which goals/decisions/tasks caused it;
- how discoveries should change unfinished work;
- what evidence justifies “done”;
- which capabilities a task needs.

Without that layer, agents repeatedly rediscover context, invent false causal stories, and promote claims without evidence.

## 2. Who it is for

**Primary (near-term):** individual developers and researchers running local AI coding agents who need better long-horizon project memory and planning.

**Secondary (later):** agent-harness vendors integrating via MCP/HTTP; open-source contributors building analyzers/plugins; research groups evaluating agent planning systems.

**Not primary early:** enterprise SaaS buyers, non-technical PMs, or teams wanting a generic AI project manager.

## 3. Core product

A **local-first project intelligence and progressive-planning layer**: a versioned knowledge graph over Git + code structure + work/causal/environment semantics, with hybrid retrieval, context compilation, evidence-backed review, and decision impact analysis.

It is **not** Git, an IDE, a coding model, a swarm framework, or (initially) a cloud product.

## 4. Explicitly out of scope (early)

From `docs/ROADMAP.md` non-goals, reinforced here:

- cloud SaaS / large dashboard;
- multi-agent swarm control plane;
- proprietary graph database;
- every language;
- automatic full-project LLM summaries;
- unrestricted autonomous execution;
- generic AI project manager UX;
- rebuilding Git content storage;
- environment-graph completeness before Gate C;
- semantic embeddings before cheaper retrieval proves insufficient.

## 5. Core entities

| Layer | Entities |
|-------|----------|
| Intent | Goal, Requirement, Constraint, Decision, Assumption |
| Work | Project, Phase, Scope, Task, Review, Evidence, Discovery, PlanVersion, PlanChange |
| Code | Repository, Directory, File, Symbol, Module |
| History | Commit, Change, ProjectState (refs into Git) |
| Environment | Skill, Rule, Tool, MCP, Hook, Agent, Model, Capability |

## 6. Core relationships (must exist in v0 conceptual schema)

- Goal → Requirement → Task
- Phase → Scope → Task
- Decision → {Task, Assumption, File}
- Task → {File, Commit, Discovery, Evidence, Decision, Task}
- Discovery → {Task, Assumption, PlanChange}
- File → Symbol; Symbol → Symbol; File → File (imports)
- Review → Evidence; Claim → Evidence → Verification → VerifiedFact

## 7. Critical human workflows

1. Bind a repository and initialize project state.
2. State a goal; receive coarse phases/scopes (not exhaustive todos).
3. Start a scope; receive deep plan + minimal tasks + exit criteria.
4. Propose a mid-project decision; see impact + alternatives; accept/reject.
5. Review discoveries and approve plan changes that affect future work.
6. Query “why does X exist?” / “what does changing X affect?”

## 8. Critical agent workflows

1. Fetch bounded task packet + progressive context expansion.
2. Implement within declared scope; emit claims + discoveries (not self-verified facts).
3. Independent reviewer reconstructs requirements, gathers evidence, PASS/FAIL/UNCERTAIN.
4. Planner propagates discoveries to affected future work with provenance.
5. Capability selection: only required skills/tools/MCPs for the task (later phases).

## 9. Architectural boundaries

| Boundary | Owns | Must not own |
|----------|------|--------------|
| Git | content, commits, diffs, OIDs | goals, decisions, evidence semantics |
| Project DB | work/causal/provenance/plan/events | file blobs / full history bodies |
| Code analysis index | structural entities/edges | intent/causal meaning |
| Retrieval | candidate selection + reasons | authority of facts |
| Context compiler | agent-facing packets | storage schema |
| Planner daemon/API | state transitions | per-adapter divergent logic |
| Review system | evidence promotion rules | implementer narrative as truth |

## 10. Unknowns (material)

- Implementation language/runtime (OPEN — blocks scaffolding).
- First agent adapter target (OPEN — blocks integration task).
- Benchmark corpus selection and licensing (OPEN — blocks Gate C).
- Whether causal edges can be seeded cheaply enough for cold-start value (EXPERIMENT).
- Whether hybrid retrieval without embeddings is enough for first Gate C (EXPERIMENT).
- Event-sourcing depth vs simple mutable tables + audit log (PROVISIONAL lean: thin events).
- Embedding/vector store choice (DEFERRED until needed).
- Plugin ABI stability timeline (DEFERRED until P18).

## 11. Most dangerous assumptions

1. Agents (or humans) will create enough high-quality causal links for “why” queries to beat Git blame + README reading.
2. Independent LLM reviewers catch false completions often enough to justify cost (H5).
3. Progressive planning reduces rework more than it adds bureaucracy (H2).
4. Graph retrieval improves outcomes vs strong baseline tooling (H1/H6) — **must be falsifiable early**.
5. SQLite + incremental indexing scales to the first interesting monorepo sizes without a graph DB.

## 12. Hardest technical problems

1. **Cold-start causality:** structural graph ≠ why-graph.
2. **Incremental correctness:** stale semantic facts after edits.
3. **Context selection:** usefulness under token budgets without missing critical constraints.
4. **Anti-hallucination promotion:** preventing claims from becoming VerifiedFacts.
5. **Replan without churn:** discovery propagation that does not thrash the board.
6. **Prompt injection via retrieved project text.**

## 13. Must be deterministic

- Git OID resolution, diffs, path existence.
- Structural parse edges (for supported languages), within analyzer version.
- Task/plan state machine transitions and append-only event append.
- Evidence artifact capture (command, exit code, stdout hash, paths).
- Retrieval “why selected” metadata for exact/lexical/graph paths.
- Scope of writes claimed vs files actually touched (diff-based).

## 14. May rely on LLM reasoning

- Goal interpretation and coarse phase proposals.
- Scope deep-planning narratives and task splits.
- Relevance ranking beyond deterministic scores (when used).
- Discovery interpretation and suggested plan changes (as proposals).
- Review narratives and hypothesis generation (never sole evidence).
- Alternative route suggestions for decisions.

## 15. Must be experimentally validated

- Gate C: graph context beats repository-only agent.
- Gate D: progressive planning beats static upfront plans on selected tasks.
- Gate E: discovery propagation reduces downstream rework.
- Gate F: impact analysis precision/recall on planted conflicts.
- Gate G: review escape rate under honesty tests.
- Gate H: performance on sized repos (10k → 1M+ LOC).

See `I_BENCHMARK_PLAN.md`.
