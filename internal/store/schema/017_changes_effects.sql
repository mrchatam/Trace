-- Migration v17: changes + effects (Phase 20 S03).
-- Additive only; do not rewrite 001–016. No ALTER on files / vcs_* / S02 tables.
-- Git SHA + path refs only — no blobs, patches, diffs, or JSON path arrays.

CREATE TABLE IF NOT EXISTS changes (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    git_commit TEXT NOT NULL DEFAULT '',
    parent_change_id TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    reason TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'RECORDED', 'COMPARED', 'SUPERSEDED')),
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    last_verified_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_changes_task_id ON changes(task_id);
CREATE INDEX IF NOT EXISTS idx_changes_git_commit ON changes(git_commit);
CREATE INDEX IF NOT EXISTS idx_changes_parent ON changes(parent_change_id);

CREATE TABLE IF NOT EXISTS change_paths (
    change_id TEXT NOT NULL,
    path TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    symbol_id TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (change_id, path)
);
CREATE INDEX IF NOT EXISTS idx_change_paths_path ON change_paths(path);

CREATE TABLE IF NOT EXISTS effects (
    id TEXT PRIMARY KEY,
    change_id TEXT NOT NULL,
    dimension TEXT NOT NULL,
    expected TEXT NOT NULL DEFAULT '',
    actual TEXT NOT NULL DEFAULT '',
    comparison TEXT NOT NULL DEFAULT ''
        CHECK (comparison IN ('', 'supported', 'partially_supported', 'contradicted')),
    confidence REAL NOT NULL DEFAULT 0,
    source_type TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE(change_id, dimension)
);
CREATE INDEX IF NOT EXISTS idx_effects_change_id ON effects(change_id);
CREATE INDEX IF NOT EXISTS idx_effects_comparison ON effects(comparison);
