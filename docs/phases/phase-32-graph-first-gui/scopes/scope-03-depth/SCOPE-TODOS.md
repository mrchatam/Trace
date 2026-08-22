# Scope 03 — board map (depth)

Inspector / graph-home depth. Serial: **S03-00 → S03-01 → S03-02**. Visual craft deferred to S04.

## Locked for implement (from P32-S03-00)

| Lock | Value |
|------|-------|
| Home | Graph at `/`; Overview → `/overview`; `/graph` alias |
| Inspector | summary → why → context → impact → reviews → links (+ loop strip tasks) |
| Select ≠ expand | Click = select; dblclick / “Use as center” = re-center |
| Budgets | DEFAULT_MAX=50 / UI_CAP=100; truncation honesty; no full-dump CTA |
| Impact | Use S02 `getImpact(taskId)` — do not re-invent |
| Path | **No** `/v1/path` |
| Law 19 | `/v1` adapters only |
| Visual | Functional only — S04 owns craft |

## Board rows

| Board ID | Row | Prompt | Role |
|----------|-----|--------|------|
| 558 | P32-S03-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner |
| 559 | P32-S03-01 | [01-implement.md](01-implement.md) | Implementer |
| 560 | P32-S03-02 | [02-review.md](02-review.md) | Reviewer |

## Implementer checklist (mirror 01)

- [x] Routes + Nav reweight
- [x] Graph select vs expand
- [x] Inspector sections + task-scoped omit
- [x] Wire `getImpact` + other ops
- [x] Build + e2e/smoke
- [x] Board Notes

## Reviewer focus (mirror 02)

Depth + budgets + select≠expand + no path invent — **not** visual polish.
