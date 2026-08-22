# Scope S01 — Impact classes

**Depends-on:** Phase 04 complete; `P05-00` done.

**Out:** Gate F harness (S02); phase VERIFY (S03); `plan simulate`; commercial impact engine (DR-NOIMP); new entity_links rels; `internal/impact` package.

**Locked (P05-S01-00):** mig **`009_decision_impact.sql`**; enums SAFE…REVERSAL + KNOWN…UNKNOWN + finding kinds; APIs Add/List findings + alternatives + `ImpactReport`; CLI `trace impact`; package `internal/domain`+store.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P05-S01-00 | planner | done | 2026-08-16: locked mig 009 + model vs live Decision CRUD/`decision_affects_task`; thickened 01+02; light S02 Gate F hooks |
| P05-S01-01 | implement | done | 2026-08-16: mig 009 + domain/store impact APIs + thin `trace impact`; see board Notes |
| P05-S01-02 | review | done | 2026-08-16: APPROVE high; no spawns — [REVIEW-NOTES.md](REVIEW-NOTES.md) |

## Checklist

- [x] P05-S01-00 planner
- [x] P05-S01-01 implement
- [x] P05-S01-02 review
