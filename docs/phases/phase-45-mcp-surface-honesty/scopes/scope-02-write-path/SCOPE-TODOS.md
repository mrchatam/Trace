# S02 — Write-path validation — SCOPE-TODOS

Issues: [#150](https://github.com/mrchatam/Trace/issues/150), [#151](https://github.com/mrchatam/Trace/issues/151), [#154](https://github.com/mrchatam/Trace/issues/154), [#162](https://github.com/mrchatam/Trace/issues/162).

| ID | Item | Status |
|----|------|--------|
| S02-0 | P45-S02-00 thicken 01/02 | **done** 2026-10-04 |
| S02-1 | Typed `store.NotFoundError`; fail-closed `trace_add`/`trace_link`/`trace_tasks` | **done** 2026-10-04 (P45-S02-01) |
| S02-2 | Verification record (self-review) | **done** 2026-10-04 **PASS** (P45-S02-02) |

## Acceptance

- [x] `trace_add` on an existing id refuses; `replace=true` opts in and reports `created:false` — #150
- [x] Non-UUID `id` and unresolvable `goal_id` both refused before insert — #151
- [x] Invalid `work_state` is an error naming `work_state` / `work_states[i]` — #154
- [x] `from == to` refused at both the MCP and the domain boundary — #162
- [x] `errors.Is(err, sql.ErrNoRows)` still true for all converted getters
- [x] No FK migration; no `replace` cascade