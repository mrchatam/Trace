# Scope 01 todos — protocol v2 (INT-08/10)

Board: [`docs/TODO/phase-27.md`](../../../../TODO/phase-27.md)

| Order | ID | Role | Prompt | Status | Notes |
|------:|----|------|--------|--------|-------|
| 456 | P27-S01-00 | planner | [00-PLANNER.md](00-PLANNER.md) | done | Locked defaults in [01-implement.md](01-implement.md); review checklist in [02-review.md](02-review.md) |
| 457 | P27-S01-01 | implementer | [01-implement.md](01-implement.md) | pending | S01-T01..T07; experiments/ only |
| 458 | P27-S01-02 | reviewer | [02-review.md](02-review.md) | pending | Verify locks; APPROVE → P27-S02-00 |

## Task seeds (from AUDIT → implement prompt)

| Task | Files | Lock |
|------|-------|------|
| S01-T01 | prepare.sh, score.sh | Preflight export in score only; prepare guard note |
| S01-T02 | score.sh | `--strict` warn-only |
| S01-T03 | score.sh | FM-07 warn-only |
| S01-T04 | PROTOCOL.md | Export + two-session docs |
| S01-T05 | RUBRIC.md | P25-3a / P25-3b |
| S01-T06 | score.sh | `--arm build\|directed` |
| S01-T07 | PROTOCOL, RUBRIC, prompts/PROMPT-G1-DIRECTED-GAP.md | Session-B + P25-4 |
