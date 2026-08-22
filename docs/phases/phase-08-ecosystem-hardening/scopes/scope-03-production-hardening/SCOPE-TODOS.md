# Scope S03 — Production hardening

**Depends-on:** `P08-S02-02` done (worktrees reviewed).

**S01 carry-in (2026-08-16):** No analyzer API version in SQLite; S01 added no mig. Optional `011_*` here is for production metadata only (prefer none).

**S02 carry-in (2026-08-16):** Path-local `.trace` + `trace.lock` fail-closed; S02 adds no mig. Backup/auth must not assume shared parent DB.

**S03-00 locks (2026-08-16):** `trace migrate status` + `trace backup`/`restore` (db snapshot + Abs rebind) + optional `.trace/access.token` + `TRACE_ACCESS_TOKEN` on Open; prefer **no** `011_*`; CLI-primary (no new MCP tools).

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P08-S03-00 | planner | done | 2026-08-16: inventory + FINAL locks; 01/02 thickened; light S04 Depends |
| P08-S03-01 | implement | done | 2026-08-16: migrate status + backup/restore + local auth; no `011_*` |
| P08-S03-02 | review | done | 2026-08-16: APPROVE high; no spawns — [REVIEW-NOTES.md](REVIEW-NOTES.md) |

## Checklist

- [x] P08-S03-00 planner
- [x] P08-S03-01 implement
- [x] P08-S03-02 review

## Phase context

- Themes: migrations + backup + **local** auth
- After S03: **S04** VERIFY owns `evals/compat` checklist (must cover migrate/backup/auth fail-closed)
