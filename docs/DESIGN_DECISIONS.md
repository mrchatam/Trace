# Initial Design Decisions

This document records the decisions made before implementation so future contributors can distinguish settled direction from open research questions.

## D1 — Project identity

The product is a project knowledge/causal graph with progressive planning and decision impact analysis, not primarily an agent orchestrator.

## D2 — Git remains canonical for file/history content

The project does not recreate Git's content/history store.

## D3 — Graph is a logical model

The first implementation may use SQLite/relational tables. A graph database is not required.

## D4 — Retrieval is hybrid

Use exact, lexical, semantic, graph, and temporal retrieval together.

## D5 — Context is progressive

The agent gets minimal task-relevant context first and can request targeted expansion.

## D6 — Agent claims are not authoritative

Evidence and verification are required according to task policy.

## D7 — Reviews are multi-layered

Todo, scope, and phase reviews have different objectives and different review strategies.

## D8 — Discovery is first-class

Implementation gaps are expected and can modify future planning.

## D9 — User decisions are first-class

Mid-project decisions trigger impact analysis and may create alternative routes.

## D10 — The system advises; the human remains authoritative

High-impact decisions should warn, explain, and propose alternatives rather than silently block users.

## D11 — Environment is part of planning

Skills, rules, tools, MCPs, hooks, agents, and permissions are project-environment entities and can be task-selected.

## D12 — Worktree isolation for concurrency

When parallel execution is introduced, prefer Git worktrees over logical file locks alone.

## D13 — Open-source core

The current intended license is Apache-2.0. Commercial value should primarily come from services, hosting, enterprise features, support, or other offerings around the open core.

## D14 — Benchmark before scale

No sophisticated optimization or orchestration feature should be treated as successful without measurement.
