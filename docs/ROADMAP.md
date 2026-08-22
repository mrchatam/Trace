# Project Roadmap

## 0. Purpose

This roadmap turns the project concept into an implementation sequence ordered by **risk reduction and learning value**, not by feature count.

The central hypothesis is:

> A persistent causal/semantic model of a software project can materially improve agent understanding, planning, verification, and decision-making compared with repository contents + Git history alone.

A second hypothesis is:

> Progressive planning — coarse planning first, deep planning near execution, continual replanning from discoveries — produces better long-running development outcomes than attempting to enumerate the entire project in advance.

A third hypothesis is:

> Evidence-driven, independent review loops can reduce false completion claims and regressions enough to justify their additional cost.

## 1. Non-goals for early development

Do not start by building:

- a cloud SaaS product;
- a large dashboard;
- a swarm framework;
- a proprietary graph database;
- support for every programming language;
- automatic project-wide LLM summaries;
- unrestricted autonomous execution;
- a generic “AI project manager”.

The project first needs to prove that its knowledge model and retrieval strategy are useful.

## 2. Phase sequence

### P0 — Specification and benchmark design

**Goal:** freeze the conceptual model before implementation.

Deliver:

- entity and relationship specification;
- provenance model;
- state-transition model;
- progressive-planning model;
- verification/evidence model;
- decision-impact model;
- environment/capability model;
- benchmark repository set;
- baseline agent tasks.

Exit criteria:

- every core concept has a defined owner and lifecycle;
- unknown/uncertain information is explicitly representable;
- benchmark tasks can be executed with a raw repository baseline.

---

### P1 — Git history substrate

**Goal:** make repository history queryable without recreating Git.

Deliver:

- repository binding;
- commit index;
- file history;
- changed-file index;
- branch/ref metadata;
- Git object references;
- incremental history refresh.

Core queries:

- `history(file)`
- `commits_between(a, b)`
- `changes(commit)`
- `origin(file)`
- `last_changed(file)`

Exit criteria:

- history can be queried from indexed metadata;
- file content remains delegated to Git;
- incremental updates do not require a full rebuild.

---

### P2 — Structural code graph

**Goal:** understand the repository’s structure.

Deliver:

- directories;
- files;
- symbols;
- imports/exports;
- references;
- calls where available;
- module relationships;
- language analyzer adapter interface.

Start with a small set of languages and expand through plugins.

Exit criteria:

- direct structural relationships can be queried quickly;
- a changed file causes only localized graph updates.

---

### P3 — Work and causal model

**Goal:** connect software evolution to project intent.

Entities:

- goals;
- requirements;
- decisions;
- assumptions;
- constraints;
- phases;
- scopes;
- tasks;
- discoveries;
- reviews;
- evidence.

Relations:

- goal → requirement;
- requirement → task;
- decision → task;
- decision → assumption;
- task → file;
- task → commit;
- task → discovery;
- discovery → task;
- decision → file;
- task → task.

Exit criteria:

- “why does this file exist?” can be answered with a causal chain where evidence exists;
- every semantic relationship has provenance.

---

### P4 — Retrieval engine

**Goal:** retrieve task-relevant knowledge without dumping the whole project into an agent context.

Implement a hybrid retrieval pipeline:

```text
exact lookup
+
lexical/BM25
+
semantic similarity
+
graph traversal
+
temporal filtering
        ↓
candidate set
        ↓
relevance/risk ranking
        ↓
context compiler
```

Retrieval should support:

- direct entity lookup;
- relevant semantic facts;
- direct graph neighbors;
- dependency expansion;
- historical evidence;
- task-linked decisions;
- recent discoveries.

Exit criteria:

- the same question answered with substantially less irrelevant context than a full-project dump;
- retrieval records why each item was selected.

---

### P5 — Context compiler

**Goal:** produce the smallest useful agent context.

Context layers:

1. task packet;
2. direct code/task context;
3. direct decisions/assumptions;
4. dependent/related work;
5. recent discoveries;
6. targeted historical context;
7. deep architecture context only when justified.

Features:

- token budgets;
- relevance scoring;
- confidence;
- provenance;
- progressive expansion;
- explicit “why retrieved” metadata.

Exit criteria:

- an agent can request more context without receiving the whole graph;
- context cost is measurable;
- cached context avoids repeated expensive work.

---

### P6 — Agent-facing interfaces

Build one canonical daemon/API with adapters:

```text
HTTP API
   ├── CLI
   ├── MCP
   ├── UI
   └── future plugins
```

Initial CLI:

- `planner explain`
- `planner why`
- `planner impact`
- `planner history`
- `planner task`
- `planner context`
- `planner discover`

MCP should expose semantic operations, not raw database primitives.

Exit criteria:

- Claude Code/Cursor/Codex-style agents can query the graph;
- the same operation has identical semantics through HTTP, CLI, and MCP.

---

### P7 — Progressive planner

**Goal:** implement the user’s actual planning methodology.

Algorithm:

```text
GOAL
  ↓
coarse project scaffold
  ↓
PHASE
  ↓
scope begins
  ↓
deep scope planning
  ↓
minimal runnable tasks
  ↓
task implementation
  ↓
review
  ↓
discoveries / evidence
  ↓
update graph
  ↓
re-evaluate future work
```

Key rule:

> The future is intentionally low-resolution until the project reaches it.

The planner must explicitly represent:

- known;
- assumed;
- unresolved;
- expected discovery areas;
- planning confidence.

Exit criteria:

- scope planning can revise downstream tasks without executing them;
- planning changes have provenance.

---

### P8 — Implementation + independent review loop

**Goal:** make implementation and verification first-class.

Flow:

```text
task
 ↓
verification plan
 ↓
implementation
 ↓
claims
 ↓
deterministic evidence
 ↓
fresh reviewer
 ↓
PASS / FAIL / UNCERTAIN
 ↓
fix / cause investigation / accept
```

Implement separate reviewer identities/contexts.

Exit criteria:

- implementation agent cannot promote its own claims directly to verified facts;
- every mandatory exit criterion has evidence;
- repeated failures trigger replanning/cause investigation.

---

### P9 — Multi-level reviews

Implement three primary review layers.

**Todo review:** correctness of the bounded task.

**Scope review:** cross-task consistency, architecture, shared assumptions, downstream effects.

**Phase review:** goal alignment, architectural trajectory, assumptions, plan integrity, and future phases.

Each level gets a different review strategy.

Exit criteria:

- all review layers can produce structured evidence;
- disagreement triggers investigation rather than a vote;
- review findings can modify future work with explicit authority.

---

### P10 — Discovery engine

**Goal:** turn implementation gaps into durable project knowledge.

Implement:

- discovery records;
- evidence links;
- affected entities;
- invalidated assumptions;
- triggered plan changes;
- discovery resolution.

Exit criteria:

- implementation can discover an unexpected constraint;
- planner can propagate that discovery into future work;
- discoveries have explicit provenance.

---

### P11 — Decision system

**Goal:** model user decisions during development.

Decision records contain:

- question;
- chosen option;
- alternatives;
- rationale;
- evidence;
- affected requirements;
- assumptions;
- confidence;
- approver;
- superseded decisions.

Exit criteria:

- the system can explain why an architecture or task direction was chosen.

---

### P12 — Decision impact analysis

**Goal:** warn about destructive or high-impact changes without taking control away from the user.

Classify impact:

- safe;
- caution;
- high impact;
- destructive;
- reversal.

Classify knowledge:

- known;
- likely;
- possible;
- unknown.

Return:

- directly affected work;
- invalidated assumptions;
- completed work at risk;
- new work required;
- unresolved consequences;
- recommended route;
- alternative routes.

Exit criteria:

- a mid-project user change produces a useful impact report before it changes the real plan;
- the user can proceed anyway;
- the system never silently changes direction.

---

### P13 — Hypothetical plan branches

Implement:

- `plan simulate`;
- alternative plan branches;
- compare impact/rework;
- adopt/discard simulation.

This is “Git branching for plan state”.

Exit criteria:

- users can compare architectural choices without modifying the active project state.

---

### P14 — Forward-state and reversal model

Explicitly model:

- project states;
- plan states;
- decision supersession;
- correction;
- reversal;
- rollback.

A reversal is a new recorded state transition, not deletion of history.

Exit criteria:

- the graph can distinguish forward progress, correction, and explicit reversal.

---

### P15 — Environment/capability graph

Represent:

- skills;
- rules;
- MCPs;
- tools;
- hooks;
- commands;
- agents;
- models;
- permissions.

Tasks can require capabilities.

Planner selects only relevant environment components.

Exit criteria:

- task planning includes environment availability;
- missing capabilities produce actionable warnings;
- agents do not receive every skill/tool/MCP by default.

---

### P16 — Performance and scale

Implement:

- incremental indexing;
- current-state materialization;
- lazy semantic analysis;
- semantic caches;
- hot/warm/cold knowledge tiers;
- bounded graph traversal;
- content-addressed analysis artifacts;
- Git delegation;
- background indexing.

Benchmark:

- 10k LOC;
- 50k LOC;
- 200k LOC;
- 1M+ LOC;
- monorepos.

Track:

- index time;
- incremental update time;
- database size;
- memory;
- query latency;
- LLM calls;
- tokens;
- cache hit rate.

Exit criteria are benchmark-driven.

---

### P17 — Multi-agent concurrency

Only after the single-agent workflow is proven.

Prefer:

- Git worktrees;
- task ownership;
- explicit integration;
- graph-aware conflict detection.

Logical path locks are advisory; worktree isolation is the stronger physical boundary.

---

### P18 — Ecosystem

Add plugin interfaces for:

- language analyzers;
- code-graph providers;
- VCS providers;
- agent adapters;
- LLM providers;
- visualization;
- planner strategies;
- impact analyzers;
- storage backends.

---

### P19 — Evaluation and research release

Publish benchmarks for:

1. project understanding;
2. planning quality;
3. discovery-driven replanning;
4. decision impact prediction;
5. review/verification quality;
6. context/token efficiency.

Compare:

```text
agent
agent + repository tools
agent + code graph
agent + project graph
agent + project graph + progressive planning
agent + graph + planning + environment selection
```

---

### P20 — Production hardening

Only after the core is validated.

Add:

- schema migrations;
- backup/recovery;
- corruption handling;
- security hardening;
- local authentication;
- observability;
- compatibility guarantees;
- stable plugin APIs.

## 3. Milestone gates

Do not advance only because features are implemented.

Advance when:

**Gate A:** graph can reconstruct useful project history.

**Gate B:** graph can answer “why” questions.

**Gate C:** agents perform better with graph context than without it.

**Gate D:** progressive planning outperforms static upfront planning on selected tasks.

**Gate E:** implementation discoveries improve future plans.

**Gate F:** decision impact analysis predicts meaningful consequences.

**Gate G:** review/evidence reduces false completion and regressions enough to justify cost.

**Gate H:** all of the above remain practical on large repositories.

## 4. Immediate implementation target

The first experimental system should be much smaller than the final roadmap:

```text
Git index
+
small code graph
+
Goal/Task/Decision/Discovery model
+
hybrid retrieval
+
context compiler
+
CLI
+
one agent integration
+
fresh review loop
```

The purpose of this first slice is to falsify or validate the core hypothesis, not to build the finished product.
