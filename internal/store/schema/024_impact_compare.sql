-- impact_predictions: one snapshot per change (C08 predict/compare)
CREATE TABLE impact_predictions (
    change_id TEXT PRIMARY KEY,
    predicted_json TEXT NOT NULL,
    compare_json TEXT NOT NULL DEFAULT '',
    depth INTEGER NOT NULL DEFAULT 2,
    created_at TEXT NOT NULL,
    compared_at TEXT NOT NULL DEFAULT ''
);

-- improvements: first-class C18 rows (CRUD in S04-05)
CREATE TABLE improvements (
    id TEXT PRIMARY KEY,
    change_id TEXT NOT NULL,
    task_id TEXT NOT NULL,
    dimension TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL DEFAULT '',
    evidence_ids_json TEXT NOT NULL DEFAULT '[]',
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX idx_improvements_change_id ON improvements(change_id);
CREATE INDEX idx_improvements_task_id ON improvements(task_id);
