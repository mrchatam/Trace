# Storage and Performance

## 1. Core principle

The system should store rich project knowledge without repeatedly reconstructing it or sending it all to agents.

## 2. Git remains the content/history substrate

Do not duplicate:

- blobs;
- complete file snapshots;
- ordinary commit history.

Store references such as:

- repository ID;
- commit OID;
- tree OID;
- path;
- symbol identity.

Retrieve actual source/history from Git when necessary.

## 3. Event history

Maintain append-only project events for:

- task changes;
- plan changes;
- discoveries;
- decisions;
- reviews;
- evidence;
- semantic updates;
- environment changes.

Current state can be materialized from events plus Git.

## 4. Hot/warm/cold model

### Hot

- current plan;
- active tasks;
- current decisions;
- current dependencies;
- active files;
- recent discoveries.

### Warm

- recent history;
- recent review records;
- related semantic summaries.

### Cold

- old semantic inferences;
- superseded planning state;
- old intermediate artifacts.

Cold data should remain queryable but should not be loaded by default.

## 5. Incremental indexing

A change should trigger a localized pipeline:

```text
changed files
   ↓
update direct code graph
   ↓
invalidate affected semantic facts
   ↓
mark dependent summaries stale
   ↓
update task/file history
```

Do not rebuild the entire repository graph for ordinary edits.

## 6. Lazy semantic analysis

Cheap deterministic information should be available continuously.

Expensive semantic analysis should run:

- on demand;
- when a task requires it;
- when a major decision changes;
- when a summary becomes stale.

## 7. Caching

Semantic artifacts should carry enough input identity to be cacheable:

- source hash;
- graph-context hash;
- analysis version;
- prompt version;
- model;
- tool version.

If the effective inputs have not changed, reuse the cached result.

## 8. Query locality

Most agent questions are local.

Prefer bounded traversals:

- direct;
- module;
- scope;
- phase;
- project.

Do not assume a project-wide traversal is necessary.

## 9. Repository tiers

Suggested repository analysis tiers:

```text
T0 ignored/generated/vendor/binary
T1 structural only
T2 semantic summaries
T3 deep semantic and historical reasoning
```

Promote active areas dynamically.

## 10. Performance metrics

Measure:

- initial indexing;
- incremental indexing;
- database size;
- memory;
- query latency;
- graph traversal cost;
- semantic cache hit rate;
- LLM calls;
- input/output tokens;
- context tokens;
- time-to-first-context.

## 11. Scaling strategy

Start with SQLite and normal indexed relational tables.

Treat the graph as a logical model, not necessarily a specialized graph database.

Do not introduce a graph database or distributed architecture until benchmark evidence requires it.
