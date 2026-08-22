# Scope S06 — Retrieval + context compiler

- [x] P00-S06-00 planner — 2026-08-15: locked exact+FTS5+graph≤2 (no embeddings); packages `internal/retrieval`+`internal/compiler` (not context/contextx); mig `004_fts` (FTS absent today); Why + TaskContext/ExpandContext; Layer 0–1 JSON+MD; reason codes; untrusted labeling; budgets 4096/32; default depth 1 max 2; no dump; thickened 01-retrieval.md; light S07 Depends; no product Go
- [x] P00-S06-01 implement — 2026-08-15: mig `004_fts` + SearchFTS/Sync*/RebuildFTS; retrieval Exact+Search+Expand+Why; compiler TaskContext/ExpandContext Layer 0–1; budgets+untrusted; CGO_ENABLED=0 tests ok
- [x] P00-S06-02 review — 2026-08-15: APPROVE high; inline Open FTS backfill fix; no spawns; next P00-S07-00 — see REVIEW-NOTES.md
