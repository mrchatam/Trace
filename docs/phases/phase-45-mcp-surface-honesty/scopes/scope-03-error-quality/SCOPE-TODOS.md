# S03 — Error quality — SCOPE-TODOS

Issues: [#155](https://github.com/mrchatam/Trace/issues/155), [#156](https://github.com/mrchatam/Trace/issues/156).

| ID | Item | Status |
|----|------|--------|
| S03-0 | P45-S03-00 thicken 01/02 | **done** 2026-10-04 |
| S03-1 | `wrapToolErr` / `wrapLinkErr`; `loop.next` → `trace.loop.gate.v1` | **done** 2026-10-04 (P45-S03-01) |
| S03-2 | Verification record (self-review) | **done** 2026-10-04 **PASS** (P45-S03-02) |

## Acceptance

- [x] No `sql: no rows in result set` or driver type names across any MCP tool — #155
- [x] Errors name the entity type actually queried — #155
- [x] `trace_link` errors name the wrong `from`/`to` **and** the rel's expected shape — #155
- [x] Hyphen and underscore rel aliases both resolve
- [x] `loop.next` returns `trace.loop.gate.v1`, not prose — #156
- [x] `reason_code` ∈ {`goal_link_missing`, `seed_task_not_found`}, each with a `remedy`
- [x] No `trace.loop.next.v1` schema introduced