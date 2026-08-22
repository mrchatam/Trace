# Scope 03 — board map

**S03 VERIFY** — Serial: **P37-S03-00 → P37-S03-01 → P37-S03-02**. Artifacts: `VERIFY-NOTES.md` + DR-HANDOFF CLOSED.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 643 | P37-S03-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock VERIFY blocks 0–5 from PLAN §6 |
| 644 | P37-S03-01 | [01-verify.md](01-verify.md) | Implementer | Run VERIFY + author VERIFY-NOTES |
| 645 | P37-S03-02 | [02-dr-handoff.md](02-dr-handoff.md) | Closer | DR-HANDOFF CLOSED + index idle / P38 promote |

## VERIFY blocks (locked S03-00)

| Block | Content | Residuals / evidence |
|-------|---------|---------------------|
| **0** | Phase 36 acceptance subset still green | 7 P36 tests — `00-p36-regression-subset.txt` |
| **1** | Per accepted residual — test or JSON | R1–R6, R8, R11 |
| **2** | Feet-seller spot-check | R8 context + R9 defer doc path (`r9-refinement-path.md`) |
| **3** | Greenfield MCP path | R3 gate + R11 workflow / Block 0 test |
| **4** | Re-defer registry + R10 browser | R7, R9, R8-full in VERIFY-NOTES; Overview + TaskDetail screenshots |
| **5** | Successor table for DR-HANDOFF | Phase 38 scaffold vs no successor — S03-02 closes |

## Inputs (must exist)

- [`PLAN.md`](../scope-01-plan/PLAN.md) §5–§6
- S02 review PASS (P37-S02-02)

## Regression baseline

Phase 36 VERIFY 7-test subset — must stay green (PLAN §5 Block 0).

## Evidence dirs

- Primary: `experiments/runs/YYYY-MM-DD-p37-s03-01-verify/evidence/`
- Pinned: `docs/verification/phase-37-p36-residuals/` (R10 human-gated)

## Out of scope

- Marking human-gated browser evidence `done` without file under `docs/verification/`
- Closing DR-HANDOFF (S03-02 only)
- Product code in S03
