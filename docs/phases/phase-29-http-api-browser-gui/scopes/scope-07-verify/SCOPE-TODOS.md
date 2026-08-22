# scope-07-verify todos

Board: [`docs/TODO/phase-29.md`](../../../../TODO/phase-29.md) orders **524–526**.

| Order | ID | Prompt | Status |
|------:|----|--------|--------|
| 524 | P29-S07-00 | [00-PLANNER.md](00-PLANNER.md) | done (planner locks VERIFY + DR-HANDOFF) |
| 525 | P29-S07-01 | [01-verify.md](01-verify.md) | pending — next runnable |
| 526 | P29-S07-02 | [02-dr-handoff.md](02-dr-handoff.md) | pending — closes Phase 29 → Phase 30 |

## Artifacts

| Artifact | Owner | Notes |
|----------|-------|-------|
| `VERIFY-NOTES.md` | S07-01 | Required; security + docs + residuals |
| `experiments/runs/YYYY-MM-DD-p29-s07-01-verify/` | S07-01 | Evidence tees |
| `DR-HANDOFF.md` | S07-02 closes | Default successor **Phase 30**; cloud ≠ Phase 30 |

## Successor lock

- Green VERIFY → **Phase 30** (stray root `trace.db`); first row **P30-00**
- Hosted SaaS → separate product/repo (CLOUD-APPENDIX); **not** Phase 30
