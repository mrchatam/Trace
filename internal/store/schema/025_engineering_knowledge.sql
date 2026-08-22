-- Migration v25: change_patterns + engineering_knowledge (Phase 22 S06-01).
-- Additive only; do not ALTER 024 tables.

CREATE TABLE IF NOT EXISTS change_patterns (
    change_kind TEXT NOT NULL,
    outcome_kind TEXT NOT NULL,
    count_positive INTEGER NOT NULL DEFAULT 0,
    count_negative INTEGER NOT NULL DEFAULT 0,
    last_seen TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (change_kind, outcome_kind)
);

CREATE TABLE IF NOT EXISTS engineering_knowledge (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL DEFAULT '',
    body_json TEXT NOT NULL DEFAULT '{}',
    topic TEXT NOT NULL DEFAULT '',
    evidence_ids_json TEXT NOT NULL DEFAULT '[]',
    confidence REAL NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'superseded')),
    source_type TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_engineering_knowledge_topic ON engineering_knowledge(topic);
CREATE INDEX IF NOT EXISTS idx_engineering_knowledge_status ON engineering_knowledge(status);
