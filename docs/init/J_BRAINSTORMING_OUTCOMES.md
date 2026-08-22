# J — Brainstorming Outcomes

Deliberate exploration before locking the first implementation scope. Items are tagged: **ADOPT** (into near-term plan), **DEFER**, **REJECT**, or **RESEARCH**.

## Product

| Idea | Tag | Notes |
|------|-----|-------|
| “Why does this file exist?” as the first undeniable human+agent value | **ADOPT** | Requires seeded causal links; define seeding path in first slice |
| Self-hosting Trace on Trace as dogfood | **ADOPT** | After vertical slice exists; strongest continuous eval |
| PR annotation bot that posts impact of a change | **DEFER** | Needs Gate C first |
| IDE sidebar (VS Code/Cursor) | **DEFER** | CLI+MCP first |
| Team multiplayer sync | **REJECT** early | Conflicts with local-first; later enterprise |
| Import from Linear/Jira/GitHub Issues as Goals/Tasks | **DEFER** | Useful adoption wedge; not needed to test H1 |
| “Stop using” risks: high setup cost, stale graph, noisy warnings, bureaucracy | **ADOPT** as design pressures | Laws + churn controls must address these |
| Simplify v0 to: index + tasks + context packets only (no planner daemon completeness) | **ADOPT** | Vertical slice over roadmap completeness |
| Agent vendor SDK with stable MCP ops | **DEFER** | After semantics stable |

**Missing use cases worth tracking:** onboarding a new human to a legacy codebase; audit trail for regulated changes; explaining regressions via decision+discovery chains; capability gap detection before a UI task starts.

## Graph

| Idea | Tag | Notes |
|------|-----|-------|
| Entity set in PROJECT_MODEL is sufficient for v0 if Environment is deferred | **ADOPT** | Slim schema for experimental slice |
| Missing relation: `Evidence verifies Claim` and `Review judges Task` | **ADOPT** | Make claim explicit in schema |
| Missing: `PlanChange caused_by Discovery\|Decision\|Review` | **ADOPT** | |
| Missing: `Task invalidates Assumption` | **ADOPT** | |
| Temporal modeling: valid_from/valid_to on semantic facts + event log | **ADOPT** thin version | Avoid full bitemporal until needed |
| Keep in Git only: blobs, commit messages, tags, ordinary history | **ADOPT** | |
| Never store: secrets, credentials, raw PII from env dumps, full embeddings of ignored/vendor trees | **ADOPT** | |
| Deterministic edges: imports, containment, commit-touches-file | **ADOPT** | |
| Inference edges: “implements requirement”, “violates constraint” — always AGENT_INFERRED until verified | **ADOPT** | |
| Optional link to external code-graph tools (Graphify, SCIP, stack-graphs) via adapter | **DEFER** | Build minimal analyzer first |

## Planning

| Idea | Tag | Notes |
|------|-----|-------|
| Plan churn budget: max N auto-replans per scope without human ack | **ADOPT** | Prevents bureaucracy loops |
| Discovery severity tiers: INFO / PLAN_AFFECTING / BLOCKING | **ADOPT** | Only PLAN_AFFECTING+ auto-opens replan |
| Contradictory future scopes: mark `CONFLICT` edge; block deep-plan until resolved | **ADOPT** | |
| “Wrong plan” recovery: supersede PlanVersion; never delete | **ADOPT** | |
| Deep-plan only current scope + one lookahead horizon | **ADOPT** | Matches progressive philosophy |
| Auto-generate entire backlog from goal via LLM | **REJECT** | Violates progressive planning |

## Retrieval

| Idea | Tag | Notes |
|------|-----|-------|
| Exact → lexical → graph before any embedding | **ADOPT** | Semantic search EXPERIMENTAL later |
| Progressive layers 0–3 as in RETRIEVAL_AND_CONTEXT | **ADOPT** | |
| Hard default expansion depth = 1; require `expand_context` | **ADOPT** | |
| Always attach retrieval reason codes | **ADOPT** | |
| History queries via Git adapter, not duplicated bodies | **ADOPT** | |
| Agent needs: task packet, exit criteria, constraints, direct files, open assumptions, recent discoveries | **ADOPT** as Layer 0–1 content | |

## Verification

| Idea | Tag | Notes |
|------|-----|-------|
| Deterministic checks first: files exist, commands exit 0, tests, diff path allowlist | **ADOPT** | |
| Independent reviewer agent for non-trivial tasks | **ADOPT** for honesty tests; optional early otherwise | |
| Human gate for subjective UX / security-sensitive | **ADOPT** | |
| Reviewer fooling modes: shared blind spot, fake logs, green tests that miss behavior | **ADOPT** into honesty suite | |
| After K failures → cause investigation required (no more patch loops) | **ADOPT** K=3 provisional | |
| Different model for high-risk review | **DEFER** | Measure first with identity/context separation |

## Decisions

| Idea | Tag | Notes |
|------|-----|-------|
| Decision as first-class with SUPERSEDED chain | **ADOPT** | |
| Destructive = invalidates completed verified work or core architecture constraint | **ADOPT** provisional definition | |
| Impact miss → user override always allowed; record OVERRIDE | **ADOPT** | |
| `plan simulate` as separate PlanVersion branch | **DEFER** to P13 | Not in first slice |
| Uncertainty on impact findings: KNOWN/LIKELY/POSSIBLE/UNKNOWN | **ADOPT** | |

## Agent environment

| Idea | Tag | Notes |
|------|-----|-------|
| Skills = procedures; Rules = persistent constraints; Hooks = triggered policies | **ADOPT** definitions | |
| Represent MCP as Capability with permissions metadata | **ADOPT** when env graph lands | |
| Env change → mark dependent tasks STALE | **ADOPT** | |
| Defer full env graph until after Gate C | **ADOPT** | Reduces surface before proof |

## Performance

| Idea | Tag | Notes |
|------|-----|-------|
| Hot: active plan/tasks/decisions; cold: superseded plans | **ADOPT** | |
| Incremental file-level reindex mandatory from day one | **ADOPT** | |
| Never LLM for structural indexing | **ADOPT** | |
| Million-LOC: tier T0 ignore vendor/generated; T1 structural only by default | **ADOPT** | |
| Years of history: index commit metadata only; content via Git | **ADOPT** | |
| Cache semantic artifacts by input hash | **DEFER** until semantic analysis exists | |

## Security

| Idea | Tag | Notes |
|------|-----|-------|
| Loopback + optional token | **ADOPT** | |
| Retrieved content = untrusted data channel | **ADOPT** | Label in context packets |
| Human confirmation: destructive decisions, capability grants, reversals | **ADOPT** | |
| Trust boundary: OS user + local daemon; agents untrusted for truth | **ADOPT** | |
| Secret redaction on evidence capture | **ADOPT** minimal patterns early | |

## Open-source / community

| Idea | Tag | Notes |
|------|-----|-------|
| Stable plugin surfaces later: language analyzers, VCS, agent adapters | **ADOPT** as future | |
| Research platform: publish benchmarks + honesty suites | **ADOPT** | |
| Keep core Apache-2.0; commercial around hosting/support | **SETTLED** direction | |
| Avoid unstable public API until Gate C | **ADOPT** | |

## Ideas that changed the plan

1. **First slice drops semantic embeddings and environment graph** to isolate the causal/work graph hypothesis.
2. **Explicit claim entity** and evidence promotion path must exist before any “done” UX.
3. **Cold-start seeding** (human/agent bootstrap of goals/decisions for a fixture repo) is part of the experiment, not an afterthought.
4. **Plan churn controls** are first-class, not a later polish item.
5. **Dogfooding Trace on itself** is the preferred long-running benchmark once the CLI exists.
