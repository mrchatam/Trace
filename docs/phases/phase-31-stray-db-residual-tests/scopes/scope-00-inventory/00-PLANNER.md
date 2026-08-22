# P31-S00-00 — Scope planner (inventory)

## Metadata
- id: P31-S00-00
- todo_ids: [P31-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, diagnosing-bugs, test-driven-development]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock inventory scope so a fresh subagent produces `GAPS.md` that confirms / trims Phase 31 coverage candidates against live tests and P30 VERIFY residuals. **No product code. No tests yet.** Finalize sibling `01-inventory.md` if still thin; do not start S01.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [Phase 31 README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [P30 VERIFY-NOTES](../../../phase-30-stray-trace-db-hygiene/scopes/scope-03-verify/VERIFY-NOTES.md)
- [P30 S02-02 review Notes](../../../phase-30-stray-trace-db-hygiene/scopes/scope-02-implement/02-review.md) (dir-stub nit)
- Live: `internal/store/stray_trace_db_test.go`, `internal/store/open.go`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Resolve from repo evidence; do not reopen store-path design.

## Locked defaults

| Item | Value |
|------|-------|
| Output of S00-01 | `scopes/scope-00-inventory/GAPS.md` |
| Product Go / test edits | **No** on S00-00 / S00-01 |
| Canonical path | `<root>/.trace/trace.db` (locked) |
| Delete flag / silent delete | Out of scope |
| GUI | Out (Phase 32) |
| Sequence | S00 → S01 → S02 serial |
| S00 review row | **None** — S01-00 validates `GAPS.md` before implement |

## Must answer (handoff to 01)

1. Which of the six phase-planner candidates are still untested in-repo?
2. Which are must-add vs nice-to-have vs out-of-scope / defer-with-reason?
3. Preferred home for each must-add: `internal/store` unit test vs CLI harness vs script-only?
4. Any additional ignore scaffolds beyond `.gitignore` + `fixtures/x0/.gitignore` that should get `/trace.db`?
5. Confirm join + Stat-only regular-file warn still true (cite `open.go`).

## Live anchors (verify still true at planner time)

| Topic | Path / note |
|-------|-------------|
| Warn choke | `open.go` — `warnIfStrayRootTraceDB` in `openStore`; `!fi.Mode().IsRegular()` → quiet |
| Unit tests | `stray_trace_db_test.go` — Open / OpenExisting / quiet / untouched |
| Gitignore | `.gitignore`, `fixtures/x0/.gitignore` — `/trace.db` |
| P30 residuals | multi-open warn intentional; agents can still create stubs; delete future-only |

## Planner gate

- [x] `01-inventory.md` runnable (metadata, preflight, grep roots, GAPS template, exit criteria)
- [x] `SCOPE-TODOS.md` lists S00 board rows
- [x] Live anchors still accurate (adjust `01` if renamed)
- [x] Do **not** write `GAPS.md` in this planner row

## Exit criteria

- [x] Inventory implementer prompt locked enough for a fresh subagent
- [x] Board row P31-S00-00 Notes cite what was verified/thickened
- [x] Next runnable remains **P31-S00-01** (do not start S01)

## Todo updates

Status + notes on **P31-S00-00** only.

## Next

`P31-S00-01`
