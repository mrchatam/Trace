# Scope S02 — Multi-agent worktrees

**Depends-on:** `P08-S01-02` done (plugin APIs reviewed).

**S01 carry-in (2026-08-16):** `LanguageAdapter` lives in `internal/analyzers` only — orthogonal to `-C`/Open/worktree bind. Do not couple adapter registration to project-root policy.

**S02-00 locks (2026-08-16):** **Per-root `.trace`** (no shared parent / no walk-up); exclusive **`trace.lock`** on `store.Open` → fail-closed concurrent same-root Open; **no** `011_*`; not swarm.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P08-S02-00 | planner | done | 2026-08-16: inventory + FINAL locks; 01/02 thickened |
| P08-S02-01 | implement | done | 2026-08-16: path-local bind + `trace.lock`/`ErrLocked` + OpenWithStore |
| P08-S02-02 | review | done | 2026-08-16: APPROVE high; residual exit 2 vs planner 1 (low) |

## Checklist

- [x] P08-S02-00 planner
- [x] P08-S02-01 implement
- [x] P08-S02-02 review

## Phase context

- After S02: **S03** production hardening → **S04** VERIFY
- Not swarm orchestration — local worktree / multi-root bind only
- S03: backup/auth target path-local `<root>/.trace/`; respect `trace.lock`; S02 adds no mig
