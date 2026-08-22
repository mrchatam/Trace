# Scope 02 — board map

**S02 Codegraph peer** — serial: **P38-S02-00 → P38-S02-01 → P38-S02-02**. Artifact: `PEER-CG.md`. **Investigate only — no product code.**

| Order | Board ID | Prompt | Role | Artifact |
|------:|----------|--------|------|----------|
| 653 | P38-S02-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 — **done 2026-08-22** |
| 654 | P38-S02-01 | [01-investigate.md](01-investigate.md) | Investigate | Author `PEER-CG.md` |
| 655 | P38-S02-02 | [02-review.md](02-review.md) | Reviewer | APPROVE / REQUEST CHANGES / SPAWN |

## Planner locks (P38-S02-00)

| Lock | Value |
|------|-------|
| Hypotheses | H1 (partial), H5, H6, H7 |
| Peer root | `similar projects/codegraph/` (read-only) |
| Method | Peer file:line + optional CG MCP — not README-only for mechanisms |
| Evidence | `experiments/runs/YYYY-MM-DD-p38-s02-654/evidence/` |
| Verdicts | supported \| weakened \| rejected \| inconclusive |
| Dogfood | Trace repo has no `.codegraph/` — optional MCP on any indexed path |
| P24 | Extend EXTERNAL-RESEARCH CG row; state deferral of `trace_explore` / consolidation |
| Spawn | Unbounded slice → S02-01a row before S03 |
| Non-goals | No product code; no implement `trace_explore`; no REMEDIATION-PLAN |

## Hypothesis → investigation todo map

| H | Todo(s) | Primary evidence files |
|---|---------|------------------------|
| H7 | T1, T5 | `h7-explore-mechanism.md`, `h7-explore-gap.md`, `h7-p24-deferred-evidence.md` |
| H6 | T2, T6 | `h6-single-tool-ux.md`, `h6-benchmark-claims.md` |
| H5 | T3 | `h5-index-watch-contrast.md` |
| H1 partial | T4 | `h1-cg-partial.md` |
| Anti-patterns / Q5 | T7 | `anti-patterns-not-for-trace.md` |
| Optional live | T8 | `h8-live-explore-sample.txt` (or skip note) |
| All | T0, T9 | Preflight + PEER-CG.md synthesis |

## Planner must-answer → PEER-CG section map

| # | Question | PEER-CG target |
|---|----------|----------------|
| Q1 | `codegraph_explore` mechanism | §1 + §2 explore row + §6 T1 cites |
| Q2 | Index/watch vs Trace index | §2 + §3 H5 |
| Q3 | Single-tool vs 16-tool discovery | §2 + §3 H6 |
| Q4 | P24 transfer still deferred? | §3 H7 + §6 T5 grep/matrix cites |
| Q5 | What Trace must NOT adopt | §4 |

## Optional tools (S02)

| H | Peer read | CG MCP | Trace MCP |
|---|-----------|--------|-----------|
| H7 | ✓ required | optional T8 | contrast only (TRACE-AUDIT) |
| H6 | ✓ required | optional | defer live re-run |
| H5 | ✓ required | — | TRACE-AUDIT cite |
| H1 partial | ✓ required | optional T8 | TRACE-AUDIT cite |

## Peer mechanism minimum (review gate)

Reviewer expects file:line cites from peer repo for at least:

```
similar projects/codegraph/src/mcp/tools.ts       — explore schema, handleExplore, DEFAULT_MCP_TOOLS, blast radius
similar projects/codegraph/src/sync/watcher.ts    — auto-sync / debounce
similar projects/codegraph/src/mcp/daemon.ts      — anti-pattern contrast (optional but recommended for §4)
```

Optional live:

```
MCP codegraph_explore — only if .codegraph/ exists on target projectPath
```

## Out of scope (S02)

- H2, H3, H4, H8, H9, H10, H11 (other scopes / S01)
- Full H1 matrix closure (S04)
- UA/Graphify peer files (S03)
- GAP-REGISTRY, SATURATION-NOTES, REMEDIATION-PLAN
- Any Go/TS product change in Trace

## Spawn rule

If T1–T8 reveal unbounded CG slice (full daemon architecture audit, entire extraction pipeline), reviewer inserts **P38-S02-01a/01b** below row 655 — do not block S03 indefinitely without S05 note.
