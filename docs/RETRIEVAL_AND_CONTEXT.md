# Retrieval and Context

## 1. Objective

The system should store a rich project graph but provide agents only the smallest useful subset.

The core question is:

> What does this agent need to know right now to make the next correct decision?

## 2. Retrieval methods

Use a hybrid system.

### Exact

For:

- file paths;
- task IDs;
- commit IDs;
- symbols;
- decisions;
- requirements.

### Lexical

For exact terminology and project vocabulary.

### Semantic

For conceptual similarity.

**Shipped (Phase 42+):** graph-label concept retrieval (G6) — bounded FTS over concept entity types (`discovery`, `assumption`, `decision`, `goal`, `claim`) with reason_code `graph_label_match`. Non-semantic lexical match on graph-adjacent labels/summaries; G9 intent feeds search terms.

**Deferred (DR-NOSSEM):** embedding/vector semantic retrieval is forbidden in Trace core. See [§3 Retrieval pipeline](#3-retrieval-pipeline) — the semantic leg is not shipped; rule-based intent (G9), graph-label channel (G6), and lexical/graph channels cover the shipped path.

### Graph

For explicit relationships and impact.

### Temporal

For historical state and “why did this exist then?” queries.

## 3. Retrieval pipeline

**Shipped (Phase 41+):** bounded rule-based intent extraction (G9) precedes lexical lookup on the compile/explore path. **Shipped (Phase 42+):** graph-label concept channel (G6) merges `graph_label_match` hits after intent FTS. **Deferred (DR-NOSSEM):** semantic/embedding retrieval — not implemented; do not emit `semantic_match`.

```text
task (+ optional agent query)
  ↓
intent extraction (G9 — rule-based keywords/entity hints → FTS query)
  ↓
exact lookup
  +
lexical retrieval (FTS5 — reason_code fts_match)
  +
graph-label concept retrieval (G6 — reason_code graph_label_match; concept entity types only)
  ↓
candidate entities
  ↓
graph expansion
  ↓
temporal filtering (optional VCS enrich)
  ↓
relevance/risk ranking + progressive layer admission (G8 max_layer)
  ↓
context budget selection
  ↓
context compiler (+ G1 query-hit merge into packet)
```

**Not shipped:** semantic retrieval leg (DR-NOSSEM). MCP direct `trace search` with raw query and no `Intent` input preserves legacy lexical-only behavior.

## 4. Progressive context

### Layer 0

- task;
- objective;
- exit criteria;
- current state.

### Layer 1

- directly affected files;
- direct symbols;
- immediate task dependencies;
- direct decisions;
- relevant assumptions.

### Layer 2

- direct dependents;
- recent discoveries;
- related future tasks;
- architectural neighbors.

### Layer 3

- deeper historical decisions;
- cross-module impact;
- alternative architecture;
- older evidence.

Deeper layers should be requested or justified rather than loaded automatically.

## 5. Context budgeting

Each candidate should have:

- relevance;
- confidence;
- freshness;
- graph distance;
- importance;
- estimated token cost.

Context selection is an optimization problem:

```text
maximize expected usefulness
subject to token/cost budget
```

## 6. Retrieval provenance

Every included item should carry a reason such as:

```text
direct task scope
dependency
created_by_task
decision dependency
recent discovery
impact edge
historical evidence
semantic match
```

This metadata may be visible to agents and developers.

## 7. Internal knowledge vs agent prompt

Never couple the database representation to the prompt representation.

Internal representation should be structured and typed.

The context compiler may render:

- concise Markdown;
- structured JSON;
- tool results;
- task packets.

## 8. Avoid full graph dumps

No general-purpose tool should return unlimited project state.

Provide bounded operations such as:

- `get_task_context`;
- `explain_entity`;
- `expand_context`;
- `get_impact`;
- `get_history`.

The agent should be able to request more context in stages.
