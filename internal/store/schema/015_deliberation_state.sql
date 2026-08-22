-- Migration v15: deliberation controller state (Phase 20 S01).
-- One row per seed task_id. Additive only; do not rewrite 001–014.

CREATE TABLE IF NOT EXISTS deliberation_state (
    task_id TEXT PRIMARY KEY,
    goal_id TEXT NOT NULL,
    current_phase TEXT NOT NULL DEFAULT 'ORIENT',
    hop_count INTEGER NOT NULL DEFAULT 0,
    last_phase TEXT NOT NULL DEFAULT '',
    plan_critiqued INTEGER NOT NULL DEFAULT 0,
    stopped INTEGER NOT NULL DEFAULT 0,
    stop_reason TEXT NOT NULL DEFAULT '',
    updated_at TEXT NOT NULL
);
