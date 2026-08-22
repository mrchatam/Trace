-- Migration v26: eval_rule_sets cache (Phase 22 S07-01).
-- Additive only; do not ALTER outcome_results.

CREATE TABLE IF NOT EXISTS eval_rule_sets (
    id TEXT PRIMARY KEY,
    source_path TEXT NOT NULL DEFAULT '',
    body_json TEXT NOT NULL DEFAULT '{}',
    updated_at TEXT NOT NULL
);
