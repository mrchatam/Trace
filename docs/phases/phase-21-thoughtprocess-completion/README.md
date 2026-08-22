# Phase 21 — TRACE thoughtprocess completion

**Status:** scaffolded (human-promoted 2026-08-18). Completes Phase 20 residuals, documented gaps vs [`TRACE_THOUGHTPROCESS.md`](../../TRACE_THOUGHTPROCESS.md), and **promoted** MVP cuts that P20 planners intentionally deferred.

Phase 20 DR-HANDOFF remains historically **`no successor`** at P20 close. Phase 21 is a **new forward queue** — not a rewrite of P20 board history.

## Why this phase exists

Phase 20 shipped the cognition **foundation** (controller, artifacts, loop integration, schema 015–019). It deliberately stopped short of:

- portable P20 cognition in seed JSON
- full deliberation phase cycle in `SelectNext`
- retrieval/why/FTS for new entity types
- baseline promotion and eval-driven promotion blocks
- §16/§18 minimal product surfaces

**Who decided those stops?** See [`DECISION-LOG.md`](DECISION-LOG.md) — P20 planners executing doc §26 MVP cuts, plus explicit S01/S07 locks, under human-promoted Phase 20 scope.

## Scope order (locked at scaffold)

| Scope | Focus | Work map |
|-------|--------|----------|
| S01 | P20 seed export + idempotent import | W-01 |
| S02 | Retrieval Exact/Why/compiler + FTS for P20 types | W-02, W-03 |
| S03 | Full SelectNext phase cycle + optional PolicyInputs enrichment | W-04, W-15 |
| S04 | Baseline promotion chain + eval regression promotion block | W-05, W-06 |
| S05 | `trace why` for P20 + historical relationship packet section | W-07, W-08 |
| S06 | Transactional apply + goal_id validation + verify-floor hygiene | W-09–W-11 |
| S07 | Promoted MVP cuts: thin §16 experiments + minimal §18 risk hints | W-12, W-13 |
| S08 | VERIFY + §31 live mini-eval + DR-HANDOFF | W-14 |

## Hard boundaries (unchanged)

- No hosted MCP / daemon / HTTP on core path
- No autonomous test runner — harness records results via `loop apply`
- No ML phase scoring; deterministic SelectNext extensions only
- No Requirement table fork (stay merged in Goal)
- Git remains canonical for diffs (Law 1)

## Completion bar

After Phase 21, a fresh agent on a **cloned repo** with imported `trace/graph.json`:

1. Sees portable deliberation state, uncertainties, changes, outcomes, regressions.
2. Gets INVESTIGATE context without retrieval stderr for `uncertainty`.
3. Can traverse ORIENT→…→REFLECT/REPLAN via `SelectNext` when preconditions met.
4. Can promote baseline after accepted evaluation; blocked when regression policy exceeded.
5. Can `trace why uncertainty <id>` and see deliberation transitions on task.
6. §31 user story proven at live (non-fixture) scale in verify evidence.

Board: [`docs/TODO/phase-21.md`](../../TODO/phase-21.md)
