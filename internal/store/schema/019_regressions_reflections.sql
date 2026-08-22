-- Migration v19: regressions + reflections (Phase 20 S05).
-- Additive only; do not rewrite 001–018. No ALTER on outcome_results / effects /
-- hypotheses / entity_links / tasks. Attribution create-default is correlated.
-- Reflections are structured JSON arrays (no essay body). No experiment tables.

CREATE TABLE IF NOT EXISTS regressions (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    source_kind TEXT NOT NULL
        CHECK (source_kind IN ('evaluation', 'contradicted_effect')),
    source_id TEXT NOT NULL,
    dimension TEXT NOT NULL DEFAULT '',
    attribution TEXT NOT NULL DEFAULT 'correlated'
        CHECK (attribution IN ('correlated', 'hypothesized', 'caused')),
    status TEXT NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'RESOLVED', 'SUPERSEDED')),
    summary TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE(source_kind, source_id, dimension)
);
CREATE INDEX IF NOT EXISTS idx_regressions_task_id ON regressions(task_id);
CREATE INDEX IF NOT EXISTS idx_regressions_status ON regressions(task_id, status);
CREATE INDEX IF NOT EXISTS idx_regressions_attribution ON regressions(attribution);

CREATE TABLE IF NOT EXISTS reflections (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    invalidated_assumptions_json TEXT NOT NULL DEFAULT '[]',
    new_dependencies_json TEXT NOT NULL DEFAULT '[]',
    useful_tests_json TEXT NOT NULL DEFAULT '[]',
    broaden_tests_note TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_reflections_task_id ON reflections(task_id);
