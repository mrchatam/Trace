# Phase 25 — Orchestrator + default gap pass (P25-C)

Human-promoted successor to Phase 24 close. This phase implements the first recommended theme: collapse E01 build-mode behavior toward directed-gap behavior by default.

## Goal

Make default build sessions run with Trace-first orchestration and mandatory post-build gap pass so discoveries/decisions happen without custom human gap prompts.

## Evidence basis

- Phase 24 matrix top ranks: [INTERVENTION-MATRIX.md](../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) (`INT-03`, `INT-04`, `INT-11`)
- Two-mode investigation: [POSTMORTEM.md](../phase-24-agent-effectiveness-investigation/scopes/scope-01-dogfood-postmortem/POSTMORTEM.md) §2
- Close decision: [DR-HANDOFF.md](../phase-24-agent-effectiveness-investigation/DR-HANDOFF.md)

## In scope

- `internal/install/` gap-pass install bundle and prompt wiring (`INT-03`)
- Orchestrator Trace-first failClosed gating and parent ownership (`INT-04`)
- Hook API drift check and maintenance spike (`INT-11`)
- Prompt/harness assets required to run P25-C safely

## Out of scope

- Product daemon/HTTP core-path work
- Discovery→task promotion implementation (`P25-A`)
- Full loop recalibration and deliberation reset (`P25-B`)
- Rewriting Phase 24 history or re-scoring E01

## Phase intent

This phase is **P25-C only**. Additional themes remain queued by promotion order: `P25-A`, then `P25-B`.
