-- Migration v18: outcome_results + baselines (Phase 20 S04).
-- Additive only; do not rewrite 001–017. No ALTER on changes / tasks / S02 tables.
-- Result kinds are library records, not a test-runner product. No log blobs.

CREATE TABLE IF NOT EXISTS baselines (
    id TEXT PRIMARY KEY,
    git_commit TEXT NOT NULL,
    scores_json TEXT NOT NULL DEFAULT '{}',
    label TEXT NOT NULL DEFAULT '',
    source_type TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_baselines_git_commit ON baselines(git_commit);

CREATE TABLE IF NOT EXISTS outcome_results (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    kind TEXT NOT NULL CHECK (kind IN ('test', 'verification', 'evaluation')),
    test_name TEXT NOT NULL DEFAULT '',
    test_status TEXT NOT NULL DEFAULT ''
        CHECK (test_status IN ('', 'pass', 'fail', 'skip', 'error')),
    goal_id TEXT NOT NULL DEFAULT '',
    verification_status TEXT NOT NULL DEFAULT ''
        CHECK (verification_status IN ('', 'verified', 'failed', 'partial')),
    baseline_id TEXT NOT NULL DEFAULT '',
    scores_json TEXT NOT NULL DEFAULT '{}',
    comparison_json TEXT NOT NULL DEFAULT '{}',
    summary TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_outcome_results_task_id ON outcome_results(task_id);
CREATE INDEX IF NOT EXISTS idx_outcome_results_kind ON outcome_results(task_id, kind);
CREATE INDEX IF NOT EXISTS idx_outcome_results_goal_id ON outcome_results(goal_id);
CREATE INDEX IF NOT EXISTS idx_outcome_results_baseline_id ON outcome_results(baseline_id);
