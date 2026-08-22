# Scope 00 — board map

**S00 investigation index** — hypothesis register + peer map. Serial: **P38-S00-00 → P38-S00-01 → P38-S00-02**. Artifact: `INVESTIGATION-INDEX.md` (S00-01, reviewed S00-02). **No product code.**

| Order | Board ID | Prompt | Role | Artifact |
|------:|----------|--------|------|----------|
| 647 | P38-S00-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 — **done 2026-08-22** |
| 648 | P38-S00-01 | [01-investigate.md](01-investigate.md) | Implementer | Author `INVESTIGATION-INDEX.md` |
| 649 | P38-S00-02 | [02-review.md](02-review.md) | Reviewer | APPROVE / REQUEST CHANGES |

## Planner locks (P38-S00-00)

| Lock | Value |
|------|-------|
| H register | H1–H11 locked in [01-investigate.md](01-investigate.md) — method, tools, peer cites, verify/reject criteria, owner |
| Evidence | `experiments/runs/YYYY-MM-DD-p38-<scope>-<row>/evidence/` |
| Spawn | New H* or oversized slice → board row in S01–S03 (max 2 cycles/scope before S05) |
| Non-goals | No INVESTIGATION-INDEX in S00-00; no REMEDIATION-PLAN; no product code |
| Saturation prep | Index must forward-fit DESIGN-LOCKS S05 checklist (see 02-review Checklist B) |

## Hypothesis → scope routing

| Hypotheses | Primary scope | Method emphasis |
|------------|---------------|-----------------|
| H2, H3, H5, H6, H9, H10 | S01 Trace audit | Trace live + doc read |
| H1 (partial), H5, H6, H7 | S02 Codegraph | Peer read + optional CG MCP |
| H1 (partial), H4, H8 | S03 UA + Graphify | Peer read + worked examples |
| H11 | S04 matrix + docs slice | Doc read |
| All | S04 cross-matrix; S05 saturation | Aggregate + gate |

## Optional tools by hypothesis

| H | Trace CLI/MCP | Codegraph MCP | Graphify | UA read |
|---|---------------|---------------|----------|---------|
| H1 | ✓ | ✓ | — | ✓ |
| H2 | ✓ | — | — | ✓ |
| H3 | ✓ | — | — | — |
| H4 | ✓ | — | ✓ | — |
| H5 | ✓ | ✓ | — | — |
| H6 | ✓ | ✓ | — | — |
| H7 | ✓ | ✓ | — | — |
| H8 | ✓ (GUI) | — | ✓ | ✓ |
| H9 | ✓ | — | — | — |
| H10 | ✓ | — | — | — |
| H11 | ✓ (optional) | ✓ (optional) | — | — |

## Out of scope (S00)

- TRACE-AUDIT, PEER-*, GAP-REGISTRY, REMEDIATION-PLAN
- Any Go/TS product change
- Writing `INVESTIGATION-INDEX.md` in planner row (S00-01 only)
