# scope-06-production todos

Board: [`docs/TODO/phase-29.md`](../../../../TODO/phase-29.md) orders **521–523**.

| Order | ID | Prompt | Status |
|------:|----|--------|--------|
| 521 | P29-S06-00 | [00-PLANNER.md](00-PLANNER.md) | done (planner) |
| 522 | P29-S06-01 | [01-implement.md](01-implement.md) | pending |
| 523 | P29-S06-02 | [02-review.md](02-review.md) | pending |

## Planner locks (S06-00) — carry into implement Notes

- Local-first only; loopback default; no multi-tenant SaaS this phase
- CORS deny + optional `--cors-origin` exact reflect; never `*`
- CSP on static; refuse `--static-dir` = project root
- Packaging: two-artifact primary + embed fallback when disk missing
- `mapDomainErr` UUID → `VALIDATION_ERROR`; seed HTTP strict stays 501
- Promote transition deny honesty (low); AGENTS/project-rules carve-out; cloud appendix design-only
- Budgeted graph; no unbounded dump

## Next runnable after planner

**P29-S06-01**
