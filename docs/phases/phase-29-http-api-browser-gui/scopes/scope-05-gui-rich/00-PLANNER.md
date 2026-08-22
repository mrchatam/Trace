# P29-S05-00 — Scope planner (feature-rich GUI)

## Metadata
- id: P29-S05-00
- todo_ids: [P29-S05-00]
- role: planner
- skills: [planning-and-task-breakdown, frontend-ui-engineering]
- verification: automated

## Objective

Lock P1 GUI wave prompts from UX-IA; produce/update `FEATURE-MATRIX.md` tracking.

## Session start

Follow agent-loop-protocol Session start.

## Locked defaults

| Item | Value |
|------|-------|
| Basis | [`../scope-02-ux-ia/UX-IA.md`](../scope-02-ux-ia/UX-IA.md) `gui_ship: S05` + OpenAPI |
| Matrix | [`FEATURE-MATRIX.md`](FEATURE-MATRIX.md) — P1 M01–M07 + optional/defer rows |
| Baseline | Live S04 `web/` — **extend**, do not rewrite MVP |
| Graph | `@xyflow/react`; default `max_nodes=50`; truncated banner; **no** Three.js |
| Loop | Write CTAs + confirms; independent status/gate loads (status fail ≠ blank gate) |
| Discoveries / Seed / Reviews | create+promote; export/import honesty; reviews nav list/detail |
| Deferrals | Explicit reason in FEATURE-MATRIX (O03–O06 locked deferred; O01–O02 optional) |
| Law 6–7 / 19 | Budgeted graph only; seed status/summary honesty; `/v1` → library only |
| S04 residual | Loop status UUID → `INTERNAL_ERROR` — prefer S06 `mapDomainErr`; S05 SPA must not blank gate |
| Tests | G-promote + G-export: Playwright e2e **or** browser smoke + cite httpapi seed path tests |

## P1 targets (from UX-IA ship matrix)

- Rich graph explorer (xyflow-class, expand-on-demand, truncated banner) — not Three.js / unbounded
- Loop console write: next / apply / reset + gate detail (confirm dialogs)
- Discoveries/decisions: create + promote (`createEntity` / `createLink` / `createTransition`); optional `listCapability` enrichment
- Seed export/import actions + path-confinement error honesty
- Tasks: full transitions + gate-aware DONE; reviews list/detail (API p1) — no S04 promotion
- Search chrome / filters as needed

## S04 review residuals (carry into planner notes / FEATURE-MATRIX)

- **Loop status honesty:** dogfood IDs like `rl010000-…` fail UUID validation → `GET /v1/loop/status` 500 `INTERNAL_ERROR` (CLI: `seed.task_id must be UUID`). Prefer map to `VALIDATION_ERROR` in httpapi (S06) and ensure S05 loop console still shows gate when status fails (S04 SPA now uses independent loads).
- S04 MVP already ships independent status/gate fetch + Discoveries initial search — do not regress when adding write CTAs.
- Keep SearchResponse `items` (never CLI `hits`); no agents/reviews nav in S04 leftovers — **S05 adds Reviews nav**.

## Exit criteria (scope)

- [x] P1 checklist in `FEATURE-MATRIX.md` done or deferred with reason (planner: rows authored as `planned`/`optional`/`deferred`)
- [x] E2E or integration coverage for promote + export honesty **locked** into S05-01 exit (G-promote / G-export)
- [x] `01-implement.md` + `02-review.md` + `SCOPE-TODOS.md` thickened for unattended run

## Next

P29-S05-01 → P29-S05-02
