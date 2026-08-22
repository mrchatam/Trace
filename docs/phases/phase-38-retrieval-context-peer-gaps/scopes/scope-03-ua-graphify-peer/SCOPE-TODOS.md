# Scope 03 — board map

**S03 UA + Graphify + Mempalace peers** — serial: **P38-S03-00 → P38-S03-01 → P38-S03-02**. Artifact: `PEER-UA-GF.md` (§1 UA · §2 GF · §3 MP). **Investigate only — no product code.**

| Order | Board ID | Prompt | Role | Artifact |
|------:|----------|--------|------|----------|
| 656 | P38-S03-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 — **done 2026-08-22** |
| 657 | P38-S03-01 | [01-investigate.md](01-investigate.md) | Investigate | Author `PEER-UA-GF.md` |
| 658 | P38-S03-02 | [02-review.md](02-review.md) | Reviewer | APPROVE / REQUEST CHANGES / SPAWN |

## Planner locks (P38-S03-00)

| Lock | Value |
|------|-------|
| Hypotheses | H1 (partial), H4, H8 + H6/H9 **Mempalace contrast slices** |
| Peer roots | UA, Graphify worked examples, **Mempalace** (read-only) |
| Method | Peer file:line — not README-only for mechanisms |
| Evidence | `experiments/runs/YYYY-MM-DD-p38-s03-657/evidence/` |
| Verdicts | supported \| weakened \| rejected \| inconclusive |
| Mempalace | Human-added peer 2026-08-22; mapped to existing H* — **H12+ not spawned** at planner |
| DR-NOSSEM | H4 may conclude law deferral vs product gap |
| Spawn | Unbounded peer slice → S03-01a row before S04 |
| Non-goals | No product code; no REMEDIATION-PLAN |

## Hypothesis → investigation todo map

| H | Todo(s) | Primary evidence files |
|---|---------|------------------------|
| H1 partial | T1, T2, T8 | `h1-ua-partial.md`, `h1-ua-search-mechanism.md`, `h1-mp-context-packet.md` |
| H4 | T3, T6 | `h4-gf-extracted-inferred.md`, `h4-semantic-contrast.md`, `h4-mp-hybrid-search.md` |
| H8 | T4, T5, T8 | `h8-gf-onboarding-ux.md`, `h8-ua-onboard.md`, `h8-mp-onboarding.md` |
| H6 (MP slice) | T7 | `h6-mp-mcp-surface.md` |
| H9 (MP contrast) | T9 | `h9-mp-fact-check-contrast.md` |
| Moat / Q8 | T10 | `moat-peers-lack.md` |
| All | T0, T11 | Preflight + PEER-UA-GF.md synthesis |

## Planner must-answer → PEER-UA-GF section map

| # | Question | PEER-UA-GF target |
|---|----------|-------------------|
| Q1 | UA context-builder vs Trace compiler | §1 + §4 H1 |
| Q2 | UA SearchEngine mechanism | §1 |
| Q3 | Graphify EXTRACTED/INFERRED (H4) | §2 + §4 H4 |
| Q4 | Graphify graph.html / worked examples (H8) | §2 + §4 H8 |
| Q5 | Mempalace hybrid search (H4) | §3 + §4 H4 |
| Q6 | Mempalace MCP + memory stack (H1/H6/H8) | §3 + §4 |
| Q7 | Mempalace fact_checker vs Trace intent (H9) | §3 + §4 H9 contrast |
| Q8 | Moat seed | §5 |

## Optional tools (S03)

| H | UA read | GF read | MP read | Trace |
|---|---------|---------|---------|-------|
| H1 partial | ✓ required | — | ✓ T8 | TRACE-AUDIT cite |
| H4 | optional T2 | ✓ required | ✓ required | doc.go DR-NOSSEM |
| H8 | ✓ onboard | ✓ required | ✓ T8 | App.tsx / optional GUI |
| H6 slice | — | — | ✓ required | TRACE-AUDIT H6 |
| H9 contrast | — | — | ✓ required | TRACE-AUDIT H9 |

## Peer mechanism minimum (review gate)

Reviewer expects file:line cites from peer repo for at least:

```
similar projects/Understand-Anything/understand-anything-plugin/src/context-builder.ts
similar projects/Understand-Anything/understand-anything-plugin/packages/core/src/search.ts
similar projects/graphify/graphify/validate.py
similar projects/graphify/graphify/symbol_resolution.py
similar projects/graphify/graphify/exporters/html.py  (or worked/.../graph.html)
similar projects/mempalace/mempalace/searcher.py
similar projects/mempalace/mempalace/layers.py
similar projects/mempalace/mempalace/mcp_server.py  (or service.py READ_TOOLS)
similar projects/mempalace/mempalace/fact_checker.py  (H9 contrast)
```

## Out of scope (S03)

- H2, H3, H5, H7, H10, H11 (other scopes / S01 / S02 / S04)
- Full H1/H6 matrix closure (S04)
- Codegraph peer re-audit (S02 PEER-CG)
- GAP-REGISTRY, SATURATION-NOTES, REMEDIATION-PLAN
- Any Go/TS product change in Trace

## Spawn rule

If T1–T10 reveal unbounded slice (full Graphify extraction pipeline, entire mempalace daemon), reviewer inserts **P38-S03-01a/01b** below row 658 — do not block S04 indefinitely without S05 note.
