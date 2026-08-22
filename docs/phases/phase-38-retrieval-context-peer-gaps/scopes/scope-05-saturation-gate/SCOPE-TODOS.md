# Scope 05 — board map

**Saturation gate** — exit investigation loops here or spawn back to S01–S04. **Blocks S06** until S05-02 APPROVE.

Serial: **P38-S05-00 → P38-S05-01 → P38-S05-02**. Artifact: `SATURATION-NOTES.md`. **Investigate only — no product code.**

| Order | Board ID | Prompt | Role | Artifact |
|------:|----------|--------|------|----------|
| 662 | P38-S05-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 — **done 2026-08-22** |
| 663 | P38-S05-01 | [01-investigate.md](01-investigate.md) | Investigate | Author `SATURATION-NOTES.md` |
| 664 | P38-S05-02 | [02-review.md](02-review.md) | Reviewer | APPROVE (saturated) / SPAWN / REQUEST CHANGES |

## Planner locks (P38-S05-00)

| Lock | Value |
|------|-------|
| Inputs | GAP-REGISTRY (661 APPROVE), INVESTIGATION-INDEX, TRACE-AUDIT / PEER-CG / PEER-UA-GF summaries |
| Method | DESIGN-LOCKS checklist walk — link S01–S04 evidence, no full re-audit |
| Evidence | `experiments/runs/YYYY-MM-DD-p38-s05-663/evidence/` |
| H7 compose-equivalence | **Defer-with-trigger** → S06; T4 desk-check in S05-01; **no spawn row** |
| Expected exit | `PROCEED_TO_S06` + `ready_for_REMEDIATION_PLAN: true` |
| Product edits | **Forbidden** |
| S06 gate | Blocked until row 664 APPROVE (saturated) |

## DESIGN-LOCKS checklist registry (S05-01 §1)

| # | Criterion | Upstream evidence |
|---|-----------|-------------------|
| 1 | H1–H11 verified/rejected/deferred | GAP-REGISTRY §2 G-001…G-011; §5 defer/reject |
| 2 | Live Trace per major gap | GAP-REGISTRY §7 → S01 `$EV/` |
| 3 | Peer mechanism CG, UA, GF | GAP-REGISTRY §2 dual-side; S02/S03 |
| 4 | Moat row | GAP-REGISTRY §3 M-001 |
| 5 | Spawn empty or deferred w/ trigger | §3–§4; H7 → §4 |
| 6 | High confidence — rows duplicate | §2 + §5 reject list |

## Investigation todo map (S05-01)

| Todo | Purpose | Evidence target |
|------|---------|-----------------|
| T0 | Preflight + `$EV/` | folder |
| T1 | DESIGN-LOCKS checklist walk | `t1-saturation-checklist.md` |
| T2 | Live Trace coverage confirm (link) | `t2-live-coverage-confirm.md` |
| T3 | Peer mechanism confirm (link) | `t3-peer-mechanism-confirm.md` |
| T4 | H7 compose desk-check | `h7-compose-desk-check.md` |
| T5 | Residual “what if X?” triage | `t5-residual-triage.md` |
| T6 | Duplication confidence + rejects | `t6-duplication-confidence.md` |
| T7 | Author SATURATION-NOTES.md | artifact |
| T8 | Self-check | — |

## H7 compose-equivalence decision (planner lock)

| Path | Verdict |
|------|---------|
| Spawn S01-01a / S02-01a live test | **Reject** — duplicate structural evidence |
| T4 desk-check + §4 defer → S06 | **Accept** — closes GAP-REGISTRY §6 trigger |
| Silent close | **Reject** — violates INVESTIGATION-INDEX H7 rejected-if path |

**Defer trigger owner:** S06 REMEDIATION-PLAN — rank G-007 unified tool vs compose-first; optional Phase 39 live spike.

## Rejected spawn seeds (planner — S05-01 §5 must cover)

| Idea | Reason |
|------|--------|
| H7 live compose-equivalence row | Duplicate |
| S01-01a symbol richness | Defer (GAP-REGISTRY §6) |
| H12+ Mempalace | Duplicate (S03 mapped) |
| S01–S03 re-audit | Duplicate |
| CG live explore on Trace | No `.codegraph/` |
| MCP 9/16 exposure | Fold → S06 |
| Embedding spike | DR-NOSSEM law |
| P38 implement trace_explore | Phase law |
| S04 registry rework | Duplicate (661 APPROVE) |

## Planner must-answer → SATURATION-NOTES section map

| # | Question | SATURATION-NOTES target |
|---|----------|-------------------------|
| Q1 | DESIGN-LOCKS checklist all boxes | §1 |
| Q2 | Confidence + rationale | §2 |
| Q3 | Rejected spawns | §5 |
| Q4 | ready_for_REMEDIATION_PLAN boolean | §7 |

## Review gate (S05-02)

Checklists A–E in [02-review.md](02-review.md):

- A: six DESIGN-LOCKS criteria
- B: planner must-answer Q1–Q4
- C: H7 defer-with-trigger + T4 desk-check
- D: spawn/defer hygiene
- E: recommendation ↔ S06 gate consistency

## Out of scope (S05)

- REMEDIATION-PLAN content (S06)
- Product Go/TS changes
- Full S01–S03 re-audit
- Live `codegraph_explore` unless T4 forces spawn
- Ranked G1–Gn build themes

## Spawn rule

If T4 or Checklist A reveals high-value room: reviewer inserts **P38-S05-01a/01b** or back-spawns to **S01–S04** below row 664 — re-enter S05 after `done`. Do **not** enter S06 until saturated APPROVE.
