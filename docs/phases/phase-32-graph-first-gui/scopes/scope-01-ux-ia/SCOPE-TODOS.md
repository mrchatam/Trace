# Scope 01 — board map

UX IA scope. Serial: **S01-00 → S01-01 → S01-02**. Primary artifact: `UX-IA.md` (authored in **S01-01**, not planner). Graph-home shell + inspector depth map (Laws 6–7). Flag API gaps for S02 (empty library OK; client `getImpact` expected). **P32-PORT out of S01** (S02).

| Board ID | Row | Prompt | Role |
|----------|-----|--------|------|
| 552 | P32-S01-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner |
| 553 | P32-S01-01 | [01-ux-ia.md](01-ux-ia.md) | Implementer → `UX-IA.md` |
| 554 | P32-S01-02 | [02-review.md](02-review.md) | Reviewer |

## Planner locks (P32-S01-00) — do not reopen in S01-01

| Lock | Value |
|------|-------|
| Home | Graph canvas (demote Overview) |
| Inspector order | summary → why → context → impact → reviews → links (+ optional loop strip on task) |
| Select vs expand | Select opens inspector; re-center is a distinct affordance |
| Budgets | `center` + `max_nodes`; DEFAULT_MAX=50; UI_CAP=100; truncation honesty |
| S02 API | Prefer `getImpact` client wrapper only; no `/v1/path`; P32-PORT still required |
| Code | None in S01 |

## Artifact

- **S01-01 writes:** `UX-IA.md` (this folder)
- **S01-00 / S01-02 must not** write the artifact (planner/reviewer rights: thicken prompts only; reviewer may inline-fix UX-IA on findings)
