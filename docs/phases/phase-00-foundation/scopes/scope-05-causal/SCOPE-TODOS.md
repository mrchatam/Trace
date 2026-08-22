# Scope S05 — Work/causal API

- [x] P00-S05-00 planner — 2026-08-15: locked `internal/domain` on `*store.Store`; mig `003` `work_state`+`entity_links`; provenance ACTIVE/STALE/SUPERSEDED; Task work machine separate; Goal→Task via `goal_id`; Decision→Task / Discovery→PlanChange via links; DONE policy stub; Claim/Evidence light only; thickened 01-causal.md; light S06/S07 Depends notes; no product Go
- [x] P00-S05-01 implement — 2026-08-15: domain on store; mig 003; Create*/Link*/TransitionTask/DONE stub/MarkStale/Claim stubs; CGO_ENABLED=0 tests ok
- [x] P00-S05-02 review — 2026-08-15: APPROVE high; no spawns; residuals non-atomic events + CreateTask WorkState escape; S06/S07 Depends thickened — [REVIEW-NOTES.md](REVIEW-NOTES.md)
