# Scope S01 — Plugin APIs

**Depends-on:** `P08-00` done. **S01 complete** (2026-08-16) — review APPROVE high. Next runnable: **`P08-S02-00`**.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P08-S01-00 | planner | done | 2026-08-16: locked LanguageAdapterAPIVersion=1 + LanguageAdapter iface + static table; thickened 01/02; light S02/S03 Depends |
| P08-S01-01 | implement | done | 2026-08-16: LanguageAdapter surface + contribution doc + tests; VERIFY green |
| P08-S01-02 | review | done | 2026-08-16: APPROVE high; no spawns — [REVIEW-NOTES.md](REVIEW-NOTES.md) |

## Checklist

- [x] P08-S01-00 planner
- [x] P08-S01-01 implement
- [x] P08-S01-02 review

## Locked shape (S01-00)

- Package: `internal/analyzers`
- `LanguageAdapterAPIVersion = 1`
- `LanguageAdapter`: `ID` / `Extensions` / `Extract`
- Compile-time static builtin table (no `.so` / megastore)
- Doc: `docs/ANALYZER_CONTRIBUTION.md`
- No mig; IndexFile orchestration stable

## Phase context

- After S01: **S02** worktrees → **S03** production → **S04** VERIFY (`evals/compat`)
- Risk mitigated: versioned adapter + static table over premature megastore
