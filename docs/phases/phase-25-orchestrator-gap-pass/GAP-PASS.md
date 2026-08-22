# GAP-PASS design SoT

Phase 25 codifies default gap-pass behavior and Trace-first orchestration to reduce mode dependence observed in E01.

## Problem statement

Phase 24 showed two modes:
- Build mode converges to thin graph / seed anchoring.
- Directed gap mode captures discoveries and fixes but required explicit human steering.

## Target outcomes

1. Default build flow ends with a mandatory gap pass (`INT-03`).
2. Parent orchestrator owns Trace task context and fails closed on non-compliant edits (`INT-04`).
3. Hook integration remains stable across Cursor API drift (`INT-11`).

## In scope

- Install-time prompt bundle and gap-pass wiring
- Orchestrator guardrails for `TRACE_TASK_ID` ownership
- Hook drift verification and maintenance checks

## Out of scope

- Autonomous discovery→task spawning policy decisions
- Full P19/hop threshold recalibration and reset semantics
- Hosted services, daemon architectures, or non-local-first control planes
