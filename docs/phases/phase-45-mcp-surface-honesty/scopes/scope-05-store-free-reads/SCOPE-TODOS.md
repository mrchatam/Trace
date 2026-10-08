# S05 — Store-free reads — SCOPE-TODOS

Issues: [#152](https://github.com/mrchatam/Trace/issues/152), [#160](https://github.com/mrchatam/Trace/issues/160).

| ID | Item | Status |
|----|------|--------|
| S05-0 | P45-S05-00 thicken 01/02 | **done** 2026-10-04 |
| S05-1 | `store.ProjectInitialized`; store-free `toolVersion`; ordering documented | **done** 2026-10-04 (P45-S05-01) |
| S05-2 | Verification record (self-review) | **done** 2026-10-04 **PASS** (P45-S05-02) |

## Acceptance

- [x] `trace_version` answers on a virgin/locked/denied project — #152
- [x] Zero side effects: no mkdir, no lock, no gate
- [x] `project_initialized` + `project_root` on the response; description matches
- [x] `trace_tasks` ordering (and the same-second tiebreak + cursor remedy) documented — #160
- [x] No migration; no nanosecond column (locked non-goal)