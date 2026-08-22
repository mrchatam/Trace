# P30-S02-00 — Scope planner (implement)

## Metadata
- id: P30-S02-00
- todo_ids: [P30-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, incremental-implementation]
- verification: automated

## Objective

Lock implementer defaults from [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md). Thicken `01-implement.md` + `02-review.md`. **No product code in this row.**

## Depends

- `PLAN.md` from **P30-S01-01** (SoT for T1–T4)
- S00 verdict **agent hygiene** (INTAKE confirmed; no Trace dual-store bug) — do not invent path redesign

## Typical files (confirm against PLAN)

- `internal/store/open.go` (warn-on-stray only — **no** join change)
- Docs: `AGENTS.md`, `docs/rules/project-rules.md`, `CONTRIBUTING.md`, optionally `cmd/trace/help.go`
- `internal/install/*` only if PLAN assigns agent-rule text
- `.gitignore` → root `/trace.db` (keep `.trace/` ignored)
- Tests beside `internal/store` (and any other package PLAN touches)

## Locks (carry from S01 — re-assert in Notes)

| Item | Value |
|------|-------|
| Canonical DB | `.trace/trace.db` only |
| Path redesign | Forbidden |
| Delete | No default / silent auto-delete of root `trace.db` |
| Warn | Optional T2: stderr, once-per-open, non-fatal; never opens/deletes root file |
| Gitignore | `/trace.db` root-only |
| Tests | Required for warn/path behavior (T4) |
| Green bar | `go test ./internal/...` |
| Scope creep | No HTTP/daemon/GUI |

## Exit criteria

- [ ] Implement + review prompts runnable alone against PLAN.md
- [ ] Next: **P30-S02-01**

## Next

`P30-S02-01` → `P30-S02-02`
