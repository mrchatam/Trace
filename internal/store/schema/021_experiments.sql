-- Migration v21: experiments (Phase 21 S07).
-- Additive only; do not rewrite 001-020. Thin §16 record — no bake-off runner.

CREATE TABLE IF NOT EXISTS experiments (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    label TEXT NOT NULL DEFAULT '',
    hypothesis_summary TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'planned'
        CHECK (status IN ('planned', 'running', 'completed')),
    outcome_result_id TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_experiments_task_id ON experiments(task_id);
CREATE INDEX IF NOT EXISTS idx_experiments_status ON experiments(task_id, status);
