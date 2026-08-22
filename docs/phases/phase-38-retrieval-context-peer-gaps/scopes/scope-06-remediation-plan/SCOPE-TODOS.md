# Scope 06 — board map

**Plan only** — after S05 saturation APPROVE.

| Order | Board ID | Prompt | Artifact |
|------:|----------|--------|----------|
| 665 | P38-S06-00 | [00-PLANNER.md](00-PLANNER.md) | (thickens 01/02) |
| 666 | P38-S06-01 | [01-plan.md](01-plan.md) | `REMEDIATION-PLAN.md` |
| 667 | P38-S06-02 | [02-review.md](02-review.md) | review verdict |

---

## Theme registry (planner lock — G1–G9)

Consolidates GAP-REGISTRY G-001…G-011 for S06-01. Severity from investigation — **build priority = rank below**.

| Theme | GAP ids | Rank | Phase sketch | H* owner |
|-------|---------|------|--------------|----------|
| **G1** Query+task orient merge | G-001, G-002 | 1 | Phase 39 | H1, H2 |
| **G3** MCP/harness orient | G-006, G-010 (+ 9/16) | 2 | Phase 39 | H6, H10 |
| **G4** Dual-stack docs | G-011 | 3 | Phase 39 docs | H11 |
| **G5** Graph onboarding UX | G-008 | 4 | Phase 39–40 | H8 |
| **G2** Read-surface strategy | G-007 | 5 | Phase 39 compose / 40+ explore | H7 |
| **G6** Concept retrieval (non-vector) | G-004b | 6 | Phase 40+ | H4 |
| **G7** Index/watch/langs | G-005 | 7 | Phase 40+ | H5 |
| **G8** Layers L2–L3 | G-003 | 8 | Phase 41+ | H3 |
| **G9** Intent pipeline | G-009 | 9 | Phase 41+ or doc-revise | H9 |

**Not themes:** G-004a vector (defer) · M-001 moat (non-gap)

---

## Ranking rubric (LOCK)

**Score = (impact × law_fit) ÷ effort** — axes 1–5 per 01-plan.md.

Tie-breakers: gap coverage → moat preservation → law-review risk → SATURATION defer owner.

---

## Locked decisions (S06-00)

| Decision | Lock |
|----------|------|
| **H7 owner** | Compose-first UX **before** unified `trace_explore` (desk-check 7/7 not equivalent) |
| **H11 stack** | **Doc-only** — CONTRIBUTING/AGENTS recipe; not product integration |
| **G-004a** | Reject/defer — not remediation theme |
| **Phase 39 entry** | G1 + G3 + G4 co-wave (human promotion) |

---

## Reject seeds (minimum 12 for REMEDIATION-PLAN §4)

| # | Reject | Source |
|---|--------|--------|
| 1 | CG detached MCP daemon as P0 | PEER-CG §4 |
| 2 | MCP-only core loop | PEER-CG §4 |
| 3 | Full-graph dump defaults | PEER-CG §4 / Law 6 |
| 4 | Graph-only product direction | PEER-CG §4 |
| 5 | Query-only replaces task packet | PEER-CG §4 / GAP §5 |
| 6 | CG benchmark % as remediation proof | PEER-CG §4 |
| 7 | Copy MP 44-tool MCP surface | GAP-REGISTRY §4 |
| 8 | Embedding/vector semantic channel | G-004a DR-NOSSEM |
| 9 | Implement any remediation in P38 | DESIGN-LOCKS |
| 10 | Product default dual-index integration | H11 investigation |
| 11 | Claim compose ≈ CG explore | h7-compose-desk-check |
| 12 | Always-on network daemon | Law / P24 |

---

## Hypothesis → S06 mapping

| H* | Primary theme | Notes |
|----|---------------|-------|
| H1, H2 | G1 | Merge query+task |
| H3 | G8 | Layers |
| H4 | G6 (+ G-004a reject) | DR-NOSSEM split |
| H5 | G7 | Index |
| H6, H10 | G3 | MCP + moat promo |
| H7 | G2 | Compose-first → explore |
| H8 | G5 | GUI orient |
| H9 | G9 | Intent |
| H11 | G4 | Doc-only stack |

---

## Tool matrix (S06)

| Tool | Use |
|------|-----|
| Read GAP-REGISTRY, SATURATION-NOTES | SoT — no re-investigate |
| Read PEER-CG §4, h7-compose-desk-check | Reject list + H7 |
| Read h11-stack-docs | H11 doc-only lock |
| `$EV/` synthesis (optional) | t1–t5 in 01-plan |
| Product code | **Forbidden** |
