# Scope 04 — board map

**S04 Cross-matrix gap registry** — serial: **P38-S04-00 → P38-S04-01 → P38-S04-02**. Artifact: `GAP-REGISTRY.md`. **Gate before S05 saturation.** **Investigate only — no product code.**

| Order | Board ID | Prompt | Role | Artifact |
|------:|----------|--------|------|----------|
| 659 | P38-S04-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 — **done 2026-08-22** |
| 660 | P38-S04-01 | [01-investigate.md](01-investigate.md) | Investigate | Author `GAP-REGISTRY.md` |
| 661 | P38-S04-02 | [02-review.md](02-review.md) | Reviewer | APPROVE / REQUEST CHANGES / SPAWN |

## Planner locks (P38-S04-00)

| Lock | Value |
|------|-------|
| Hypotheses | **All H1–H11** aggregated; H11 doc read owned by S04 |
| Inputs | TRACE-AUDIT, PEER-CG, PEER-UA-GF (APPROVED — do not re-audit) |
| Method | Synthesis + H11 doc read — link S01–S03 evidence |
| Evidence | `experiments/runs/YYYY-MM-DD-p38-s04-660/evidence/` (H11 + matrix notes) |
| Matrix columns (LOCK) | **Trace \| CG \| UA \| GF \| MP \| moat row** |
| Gap IDs | **G-001…G-011** (1:1 H*); moat **M-001** |
| Verdicts | gap \| non-gap \| defer |
| Severity | Investigation confidence only — **not** build priority (S06) |
| Dual-side evidence | Required per **gap** row (Trace + peer cite) |
| Spawn | Incomplete registry → §6 + may spawn S01–S03 or S04-01a |
| Non-goals | No REMEDIATION-PLAN; no product code; no G1–Gn ranking |

## Gap ID registry (planner lock)

| Gap ID | H* | Theme | Upstream verdict seeds |
|--------|-----|-------|------------------------|
| G-001 | H1 | Unified query+task orient packet | S01 confirmed gap (partial); S02/S03 supported |
| G-002 | H2 | Compiler FTS title-only | S01 confirmed gap |
| G-003 | H3 | Layers 2–3 designed not shipped | S01 confirmed gap |
| G-004 | H4 | Semantic/concept retrieval | S03 supported; DR-NOSSEM vector **defer** |
| G-005 | H5 | Index langs + manual vs watcher | S01 + S02 supported |
| G-006 | H6 | MCP discovery (16 / 1 / 44) | S01 + S02 + S03 MP slice supported |
| G-007 | H7 | `trace_explore` unified read | S02 supported; P24 still deferred |
| G-008 | H8 | Graph-first onboarding hook | S01 inconclusive → S03 supported |
| G-009 | H9 | Intent pipeline doc-only | S01 + S03 MP contrast supported |
| G-010 | H10 | Moat under-promoted in install | S01 confirmed gap |
| G-011 | H11 | Trace+CG stack undocumented | S04 T8 doc read |
| M-001 | moat | Trace strengths peers lack | S01 §5 + S02 §5 + S03 §5 merge |

## Hypothesis → investigation todo map

| H / G-ID | Todo(s) | Primary evidence (link, don't duplicate) |
|----------|---------|---------------------------------------------|
| All ingest | T1, T2, T3 | Upstream artifacts + `t1-*` `t2-*` `t3-*` seeds |
| G-001–G-004 | T4, T5 | `t4-gap-id-registry.md`, `t5-matrix-h1-h4.md` |
| G-005–G-007 | T6 | `t6-matrix-h5-h7.md` |
| G-008–G-010 | T7 | `t7-matrix-h8-h10.md` |
| G-011 | T8 | `h11-stack-docs.md` |
| M-001 | T9 | `t9-moat-row-m001.md` |
| Non-gaps / defer | T10 | `t10-non-gaps-deferrals.md` |
| Spawn / S05 | T11 | `t11-spawn-triggers.md` |
| All | T0, T12 | Preflight + GAP-REGISTRY.md synthesis |

## Planner must-answer → GAP-REGISTRY section map

| # | Question | GAP-REGISTRY target |
|---|----------|---------------------|
| Q1 | Unified G-001…G-011 linked to H* | §2 preamble + main table |
| Q2 | Severity = investigation confidence only | §2 Severity column + §1 note |
| Q3 | Moat row Trace strengths | §3 M-001 |
| Q4 | Spawn triggers for S05 | §6 |

## Matrix column rules (review gate)

| Column | Source artifact | Required for H* |
|--------|-----------------|-----------------|
| **Trace** | TRACE-AUDIT | All G-001…G-011 |
| **CG** | PEER-CG | H1, H5, H6, H7; N/A justified elsewhere |
| **UA** | PEER-UA-GF §1 | H1, H8; optional H2 contrast |
| **GF** | PEER-UA-GF §2 | H4, H8; N/A justified elsewhere |
| **MP** | PEER-UA-GF §3 | H1, H4, H6 slice, H8, H9 contrast |
| **Moat** | §3 M-001 row | Once per registry (not per H*) |

## Dual-side evidence minimum (review gate)

Reviewer expects for each **gap** row:

```
Trace:  file:line OR experiments/runs/…-p38-s01-651/evidence/…
Peer:   file:line OR experiments/runs/…-p38-s02-654/… OR …-s03-657/…
```

H11 (G-011): Trace doc cites sufficient; CG complementary cite if claiming dual-stack gap.

## S05 spawn triggers (planner lock)

| Trigger | Action |
|---------|--------|
| Any gap row missing dual-side evidence | REQUEST CHANGES or S04-01a |
| G-011 inconclusive after T8 | S04-01a doc slice or defer in §6 with trigger |
| G-007 compose-equivalence untested | Spawn S01/S02 live compose OR §6 explicit open |
| Upstream verdict contradiction | Justify in registry or spawn back to owning scope |
| New H12+ uncovered | INVESTIGATION-INDEX §5 spawn rules |

## Out of scope (S04)

- Live re-audit of entire Trace stack (S01 done)
- Full CG/UA/GF/MP peer re-read (S02/S03 done)
- SATURATION-NOTES (S05), REMEDIATION-PLAN (S06)
- Any Go/TS product change
- Ranked build priority / G1–Gn themes

## Spawn rule

If T5–T11 reveal unbounded synthesis (re-litigating S01 live commands, full peer re-scan), reviewer inserts **P38-S04-01a/01b** below row 661 — do not block S05 indefinitely without SATURATION-NOTES defer.
