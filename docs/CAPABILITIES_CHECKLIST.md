# Trace — Final Capability Checklist

Use `- [x]` to mark items as implemented. Checkbox state is the only field intended to change during review.

**141 capabilities** across 17 sections · **8 non-goals** · [Jump to reference](#reference)

---

## Table of contents

| # | Section | Items |
|---|---------|------:|
| 1 | [Persistent Project Understanding](#1-persistent-project-understanding) | 8 |
| 2 | [Project Graph](#2-project-graph) | 9 |
| 3 | [Change Tracking](#3-change-tracking) | 10 |
| 4 | [Impact Analysis](#4-impact-analysis) | 8 |
| 5 | [Agent Thought Process / Engineering Reasoning](#5-agent-thought-process--engineering-reasoning) | 9 |
| 6 | [Test → Verify → Score Loop](#6-test--verify--score-loop) | 13 |
| 7 | [Regression Detection](#7-regression-detection) | 7 |
| 8 | [Change → Effect Learning](#8-change--effect-learning) | 11 |
| 9 | [Evidence-Based Agent Feedback](#9-evidence-based-agent-feedback) | 8 |
| 10 | [Continuous Project Intelligence](#10-continuous-project-intelligence) | 7 |
| 11 | [Agent Coordination Through Shared Understanding](#11-agent-coordination-through-shared-understanding) | 7 |
| 12 | [Observability of Engineering Work](#12-observability-of-engineering-work) | 8 |
| 13 | [Queryable Engineering Knowledge](#13-queryable-engineering-knowledge) | 10 |
| 14 | [Incremental / Non-Destructive Operation](#14-incremental--non-destructive-operation) | 5 |
| 15 | [Trace Core](#15-trace-core) | 7 |
| 16 | [Agent-Facing Interface](#16-agent-facing-interface) | 7 |
| 17 | [Extensible Evaluation System](#17-extensible-evaluation-system) | 7 |

---

## Capabilities

### Foundation — understanding, graph, changes, impact

#### 1. Persistent Project Understanding

- [x] Maintain a persistent representation of the project and its evolution.
- [x] Understand the structure of the codebase beyond individual files.
- [x] Track relationships between files, modules, components, functions, types, APIs, tests, and other relevant artifacts.
- [x] Track dependencies and dependency relationships.
- [x] Track architectural relationships and boundaries.
- [x] Allow agents to query the project model instead of repeatedly rediscovering the same context.
- [x] Preserve useful project knowledge across agent sessions.
- [x] Update project knowledge as the codebase changes.

#### 2. Project Graph

- [x] Maintain a graph representing important entities and relationships in the project.
- [x] Represent code entities and their relationships.
- [x] Represent dependencies and impact relationships.
- [x] Represent tests and what they validate.
- [x] Represent changes and what they affect.
- [x] Represent relevant architectural/project-level relationships.
- [x] Allow graph traversal and targeted queries.
- [x] Support incremental graph updates rather than rebuilding everything unnecessarily.
- [x] Keep graph state synchronized with the actual project.

#### 3. Change Tracking

- [x] Record every meaningful change made to the project.
- [x] Associate changes with the agent/task/action that produced them when possible.
- [x] Record what files/components/entities were changed.
- [x] Record why the change was made when that information is available.
- [x] Record the intended effect of a change.
- [x] Record the actual observed effect.
- [x] Preserve historical project state.
- [x] Allow agents to understand how the project evolved.
- [x] Allow comparison between project states.
- [x] Make changes and their consequences queryable.

#### 4. Impact Analysis

- [x] Determine what a proposed change may affect before implementation.
- [x] Identify directly affected components.
- [x] Identify transitively affected components.
- [x] Identify potentially affected tests.
- [x] Identify dependencies and consumers that may be affected.
- [x] Use the project graph to reason about change impact.
- [x] Help agents discover areas they might otherwise overlook.
- [x] Compare predicted impact with actual impact after implementation.

---

### Agent reasoning

#### 5. Agent Thought Process / Engineering Reasoning

- [x] Give agents structured project context for reasoning.
- [x] Help agents understand the current state before modifying it.
- [x] Support a structured implementation/thought process rather than blindly executing instructions.
- [x] Encourage agents to form an implementation plan before making changes.
- [x] Encourage agents to reason about affected components.
- [x] Encourage agents to identify risks and assumptions.
- [x] Preserve relevant reasoning/context needed to understand decisions.
- [x] Connect reasoning/decisions to resulting changes and outcomes.
- [x] Allow agents to learn from previous engineering decisions.

---

### Verification, regression, and learning

#### 6. Test → Verify → Score Loop

> This is a core part of the final direction.

- [x] Require implementations to go through a test/verification cycle.
- [x] Automatically or explicitly run relevant tests after changes.
- [x] Verify that the implementation satisfies its intended behavior.
- [x] Verify that existing functionality has not regressed.
- [x] Verify relevant architectural/invariant constraints where possible.
- [x] Evaluate implementation quality using a scoring mechanism.
- [x] Score outcomes based on evidence rather than simply accepting successful execution.
- [x] Detect failures and regressions.
- [x] Feed verification results back into the agent's reasoning loop.
- [x] Allow agents to iterate on implementations based on verification results.
- [x] Record test and verification results as part of project history.
- [x] Record scores and evaluation results.
- [x] Compare results between iterations.

#### 7. Regression Detection

- [x] Detect regressions caused by changes.
- [x] Identify which change is associated with a regression.
- [x] Track previously working behavior.
- [x] Detect deterioration across iterations.
- [x] Preserve evidence of regressions.
- [x] Make regression history queryable.
- [x] Help agents avoid repeating changes that previously caused problems.

#### 8. Change → Effect Learning

> This is another major distinction from the original control-plane concept.

- [x] Track the relationship between a change and its observed effects.
- [x] Record positive effects.
- [x] Record negative effects.
- [x] Record regressions.
- [x] Record improvements.
- [x] Record verification results.
- [x] Record performance/quality changes when measurable.
- [x] Identify recurring patterns between types of changes and outcomes.
- [x] Allow agents to query historical evidence before making similar changes.
- [x] Gradually build project-specific engineering knowledge from historical changes.
- [x] Help agents understand what tends to improve or damage a particular project.

---

### Intelligence, coordination, and observability

#### 9. Evidence-Based Agent Feedback

- [x] Give agents evidence from previous implementations.
- [x] Surface relevant historical changes.
- [x] Surface previous failures and regressions.
- [x] Surface successful approaches.
- [x] Surface verification results.
- [x] Surface relevant project relationships.
- [x] Distinguish observed facts from assumptions where possible.
- [x] Help agents make decisions based on accumulated project evidence.

#### 10. Continuous Project Intelligence

- [x] Continuously update Trace as the project changes.
- [x] Maintain an evolving understanding of the project.
- [x] Accumulate engineering knowledge over time.
- [x] Avoid treating every agent task as an isolated event.
- [x] Connect tasks, changes, tests, results, and project evolution.
- [x] Allow later agents to benefit from earlier agents' work.
- [x] Make the system progressively more knowledgeable about the specific project.

#### 11. Agent Coordination Through Shared Understanding

- [x] Provide a shared project state for multiple agents.
- [x] Allow agents to understand work performed by other agents.
- [x] Make previous changes and their effects visible to subsequent agents.
- [x] Reduce duplicated investigation.
- [x] Reduce conflicting or redundant work.
- [x] Allow agents to build on previous work instead of starting from zero.
- [x] Preserve continuity across independent agent sessions.

#### 12. Observability of Engineering Work

- [x] Record what happened during agent-driven engineering work.
- [x] Record actions relevant to understanding project evolution.
- [x] Connect actions to changes.
- [x] Connect changes to verification.
- [x] Connect verification to outcomes.
- [x] Make this history inspectable.
- [x] Make failures and successful iterations traceable.
- [x] Provide enough context to understand why the project reached its current state.

#### 13. Queryable Engineering Knowledge

- [x] Allow agents to ask questions about the current project.
- [x] Allow agents to ask questions about project history.
- [x] Allow agents to ask what changed.
- [x] Allow agents to ask why something changed.
- [x] Allow agents to ask what a change affected.
- [x] Allow agents to ask what tests verify something.
- [x] Allow agents to ask what previously failed.
- [x] Allow agents to ask what approaches previously worked.
- [x] Allow agents to ask about regressions.
- [x] Allow agents to query accumulated evidence when planning new work.

---

### Platform and operations

#### 14. Incremental / Non-Destructive Operation

- [x] Integrate with existing repositories rather than requiring Trace to own the project.
- [x] Observe and analyze project changes without unnecessarily interfering with development.
- [x] Update its knowledge incrementally.
- [x] Preserve historical information rather than replacing it with only the latest state.
- [x] Support normal developer/agent workflows.

#### 15. Trace Core

- [x] Provide a core engine responsible for maintaining Trace's project intelligence.
- [x] Maintain the project graph.
- [x] Maintain project/change history.
- [x] Maintain relationships between changes and effects.
- [x] Coordinate testing, verification, and scoring.
- [x] Provide the underlying query/reasoning capabilities to agents.
- [x] Provide a foundation that other agent-facing interfaces can use.

#### 16. Agent-Facing Interface

- [x] Expose Trace's knowledge to agents.
- [x] Allow agents to retrieve relevant context.
- [x] Allow agents to query project relationships.
- [x] Allow agents to inspect historical evidence.
- [x] Allow agents to record/associate work with Trace.
- [x] Allow agents to initiate or participate in verification loops.
- [x] Make Trace useful as part of an agent's normal engineering workflow.

#### 17. Extensible Evaluation System

- [x] Support multiple verification mechanisms.
- [x] Support multiple scoring criteria.
- [x] Allow project-specific evaluation rules.
- [x] Allow different types of tests/evidence to contribute to evaluation.
- [x] Preserve evaluation history.
- [x] Make evaluation results available to future agents.
- [x] Allow the evaluation system to become more sophisticated without redesigning Trace's core model.

---
