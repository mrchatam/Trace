# Architecture

## 1. Architectural goal

The system is a local-first project intelligence and planning layer between a software repository and AI agents.

It is not:

- Git;
- an IDE;
- a coding model;
- a generic project manager;
- a cloud requirement.

## 2. High-level topology

```text
                 HUMAN / IDE / AGENT
                         │
              ┌──────────┴──────────┐
              │                     │
             CLI                   MCP
              │                     │
              └──────────┬──────────┘
                         ▼
                 Planner Daemon
                         │
          ┌──────────────┼───────────────┐
          ▼              ▼               ▼
      Project DB      Retrieval       Git Adapter
          │              │               │
          │              ├── exact       ├── commits
          │              ├── lexical     ├── trees
          │              ├── semantic    ├── diffs
          │              ├── graph       └── content
          │              └── temporal
          │
          └──────────────┬────────────────
                         ▼
                 Context Compiler
                         │
                         ▼
                      Agent
```

## 3. Core separation

### Git

Source of truth for:

- file contents;
- commits;
- branches;
- diffs;
- object identities;
- repository history.

### Project database

Source of truth for:

- goals;
- decisions;
- tasks;
- discoveries;
- semantic facts;
- provenance;
- plan state;
- evidence;
- environment metadata.

### Code analysis index

Source of truth for derived structural relationships:

- symbols;
- imports;
- references;
- calls;
- module structure.

### Retrieval layer

Chooses relevant information.

### Context compiler

Converts internal knowledge into a compact agent-facing representation.

### Planner

Decides what work should happen next and how newly discovered information modifies future work.

### Review system

Validates claims and creates evidence.

## 4. One canonical service

The daemon owns state transitions.

The CLI, MCP, and UI must not independently implement business logic.

```text
CLI ───────┐
MCP ───────┤
UI ────────┤──> canonical API/state machine
Plugins ───┘
```

This prevents diverging behavior.

## 5. Local-first deployment

Initial deployment:

- single user;
- loopback daemon;
- one daemon, many projects;
- repository data remains local;
- optional local token for API access.

A remote/hosted service is a later deployment target, not a core assumption.

## 6. Query architecture

Every query should follow roughly:

```text
intent
  ↓
query classifier
  ↓
cheap exact/lexical retrieval
  ↓
semantic retrieval where useful
  ↓
graph expansion
  ↓
temporal filtering
  ↓
relevance/risk ranking
  ↓
context compiler
```

No layer should assume that the whole project can be loaded.

## 7. Event and state model

State changes produce append-only events.

Derived indexes/materialized state are allowed to be rebuilt from events + Git.

This is intended to make:

- audit;
- debugging;
- recovery;
- provenance;
- historical reasoning

reliable.

## 8. Adapter boundaries

### Agent adapter

Responsibilities:

- start/stop agent;
- provide selected context/capabilities;
- capture output;
- collect evidence;
- report completion/failure.

### Code analyzer adapter

Responsibilities:

- parse source;
- emit structural entities/edges;
- support incremental updates.

### VCS adapter

Responsibilities:

- read repository history;
- resolve paths;
- read diffs;
- provide commit/tree references.

### Storage adapter

The first implementation may be SQLite. The data model should not depend on SQLite-specific behavior.

## 9. Physical isolation

When multiple agents work concurrently:

- prefer Git worktrees;
- each work unit gets an explicit integration boundary;
- path locks are advisory conflict detection, not a security boundary.

The system must never claim that a logical path lock alone prevents an agent from modifying an undeclared file.
