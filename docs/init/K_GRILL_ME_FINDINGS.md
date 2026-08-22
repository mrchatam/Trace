# K — Grill-Me Findings

Adversarial interrogation of the product thesis, architecture, and plan. Each finding records challenge → response → disposition.

Disposition values: **JUSTIFIED** | **NARROWED** | **EXPERIMENT** | **ACCEPTED_RISK** | **CHANGED** | **OPEN**

---

## Product thesis

### G1. Is this a real problem or graph fascination?

**Challenge:** Strong agents + ripgrep + Git + AGENTS.md may already be “good enough.”

**Response:** The distinctive gap is durable causal/planning state with provenance and discovery propagation—not another code search. Value appears only if agents/humans populate and use that state.

**Disposition:** **NARROWED** — first undeniable value defined as: *bounded task context + why-chain for seeded entities beats repo-only agent on planted understanding tasks.* If Gate C fails, stop expanding features.

### G2. Smallest experiment that could kill the project?

**Challenge:** What falsifies the thesis cheaply?

**Response:** On a fixture repo with known causal ground truth, agent+repo-tools vs agent+project-graph on understanding + one implementation task. If graph does not improve accuracy or reduce critical misses after fair seeding cost, thesis fails.

**Disposition:** **EXPERIMENT** — see `I_BENCHMARK_PLAN.md` Experiment X0.

### G3. First undeniable capability?

**Challenge:** Roadmap lists many capabilities; which one sells the rest?

**Response:** Not orchestration. **Task-bounded context packets with provenance + “why/impact” over seeded decisions/tasks.** Review loops and env graphs are amplifiers after that.

**Disposition:** **JUSTIFIED** / **CHANGED** relative to building P15–P17 early — deferred.

---

## Architecture

### G4. Duplicating Git?

**Challenge:** Commit index may reinvent `git log`/`git blame`.

**Response:** Store references and query accelerators only; content always from Git. Justification is joinability to Task/Decision/Discovery—not history UI.

**Disposition:** **JUSTIFIED** with law: no blob duplication.

### G5. Duplicating code graphs?

**Challenge:** SCIP/Graphify/stack-graphs already exist.

**Response:** Use a minimal analyzer first; treat external providers as adapters later. Structural graph is necessary but insufficient; product value is causal/work layer.

**Disposition:** **JUSTIFIED** / **ACCEPTED_RISK** of temporary inferior structure quality.

### G6. Overusing a graph database?

**Challenge:** Neo4j/etc. premature.

**Response:** Docs already say logical graph on SQLite. Keep it.

**Disposition:** **JUSTIFIED** (SETTLED D3).

### G7. Event model justified?

**Challenge:** Full event sourcing may be heavy for v0.

**Response:** Need audit/provenance/rebuild. **Thin append-only event log + materialized tables** is enough; not a distributed event bus.

**Disposition:** **NARROWED** — provisional thin events (see Decision DR-EVT).

### G8. Coupling to one LLM/agent?

**Challenge:** Planner prompts tied to one vendor.

**Response:** Canonical daemon API; agent adapters at edge; no vendor in core types.

**Disposition:** **JUSTIFIED**; first adapter still **OPEN** (Q-AGENT).

### G9. Retrieval coupled to planning?

**Challenge:** One blob of “smart context” mixes concerns.

**Response:** Keep retrieval + compiler as pure query path; planner consumes them but does not own ranking internals.

**Disposition:** **JUSTIFIED** — enforce in first scope module boundaries.

---

## Planning

### G10. Progressive planning as bureaucracy?

**Challenge:** Constant replans thrash agents.

**Response:** Severity tiers + churn budget + human ack for PLAN_AFFECTING changes outside current task.

**Disposition:** **CHANGED** — add churn controls to laws and first planner design (even if planner is minimal).

### G11. Endless divergence of plan vs reality?

**Challenge:** Graph becomes fanfiction about the repo.

**Response:** STALE markers on semantic facts when files change; todo review compares claims to actual Git state; dogfood will surface this.

**Disposition:** **EXPERIMENT** + **ACCEPTED_RISK** early.

### G12. Agents modifying future work recklessly?

**Challenge:** Implementer rewrites the roadmap.

**Response:** Implementer may *propose* PlanChange; adoption requires planner rules + review authority levels.

**Disposition:** **JUSTIFIED** — encode in verification/planner permissions matrix.

---

## Retrieval

### G13. Irrelevant expansion / context explosion?

**Response:** Depth caps, budgets, expand_context, ban full dumps.

**Disposition:** **JUSTIFIED**.

### G14. Semantic retrieval plausible-but-wrong?

**Response:** Do not ship embeddings in X0. Prefer exact/lexical/graph. When added, label SEMANTIC matches with lower authority than DETERMINISTIC edges.

**Disposition:** **CHANGED** — embeddings deferred.

### G15. How prove graph retrieval helps?

**Response:** Gate C / X0 with logged retrieval traces and token counts.

**Disposition:** **EXPERIMENT**.

---

## Verification

### G16. Implementer lies; reviewer shares blind spot?

**Response:** Deterministic evidence required for automatable criteria; honesty suite with planted failures; human gates for subjective claims; optional later multi-model review.

**Disposition:** **JUSTIFIED** strategy; **ACCEPTED_RISK** that two LLMs can agree on falsehood without deterministic checks.

### G17. Tests pass but claim still wrong?

**Response:** Verification policy by task class; residual risks recorded; not all truths are automatable.

**Disposition:** **JUSTIFIED**.

---

## Decision analysis

### G18. Impact engine misses or overstates?

**Response:** Confidence bands; user override; measure precision/recall in later Gate F. Not in first slice beyond Decision CRUD + manual affected-entity links.

**Disposition:** **NARROWED** — full impact engine deferred; schema ready.

---

## Scaling

### G19. Worst-case storage / monorepo?

**Response:** Git delegation + ignore tiers + incremental updates. Gate H benchmarks required before optimization theater.

**Disposition:** **JUSTIFIED** direction; numbers **UNKNOWN** until measured.

### G20. Stale semantic info / rebuild safety?

**Response:** content hashes; STALE; rebuild materialized views from events+Git; version analyzer outputs.

**Disposition:** **JUSTIFIED** design intent; implementation **EXPECTED_DISCOVERY**.

---

## Adoption

### G21. Why not Graphify/Graphiti/scripts?

**Response:** Those are partial layers (code graph or memory). Trace’s bet is progressive planning + evidence + decision causality integrated. If users only need code graphs, they should use those tools—Trace must not pretend to replace them on day one.

**Disposition:** **JUSTIFIED** positioning.

### G22. Minimum useful experience?

**Response:** CLI: `init`, `task context`, `why <entity>`, `discover`, `review submit` on one repo.

**Disposition:** **JUSTIFIED** — defines first scope UX surface.

---

## Final adversarial pass over *this* plan

| Question | Answer |
|----------|--------|
| Building too much? | First scope intentionally slims roadmap P1–P8 into a vertical slice without embeddings, env graph, impact engine, multi-agent, UI. |
| Building too little? | Slice still includes work/causal entities + retrieval + context + CLI + one review honesty path—enough to run X0. |
| Duplicating Git? | Forbidden by law; indexes only. |
| LLM where deterministic possible? | Forbidden for structure/evidence capture. |
| Can discoveries revise future tasks? | Yes via PlanChange proposals; adoption rules prevent silent rewrite. |
| Can users change direction? | Decision records + override; full simulate deferred. |
| False success claims? | Evidence policy + independent review on honesty tasks. |
| Endless fix loops? | K=3 then cause investigation. |
| Kill experiment defined? | X0 in benchmark plan. |
| Language/agent still open? | Yes — Round 1 blockers; recommended defaults recorded. |

---

## Unresolved after grilling (need user)

**Round 1–2:** resolved. Implementation of `T001+` is unblocked.

**Round 3:** no blocking questions. Soft defaults: ≥5 understanding queries (Q-UNDERSTAND-N); TS+Python fixture (Q-FIXTURE-LANG).

### G23. P0 closure criterion

**CHANGED** — P0 requires P0-X (DR-P0).

### G24. P0-X strength / incremental architecture

**Challenge:** A CRUD-only or full-rebuild foundation would validate the wrong product.

**Response:** User raised P0-X to 7 points including tree-sitter structure and incremental file update; DR-RISK / DR-INCREMENTAL / Law 12+21.

**Disposition:** **CHANGED** — settled.
