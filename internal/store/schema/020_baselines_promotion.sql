-- Migration v20: baseline promotion chain (Phase 21 S04).
-- Additive only; do not rewrite 001-019.

ALTER TABLE baselines ADD COLUMN status TEXT NOT NULL DEFAULT 'active'
    CHECK (status IN ('active', 'superseded'));
ALTER TABLE baselines ADD COLUMN supersedes_id TEXT NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_baselines_commit_label_status
    ON baselines(git_commit, label, status);
