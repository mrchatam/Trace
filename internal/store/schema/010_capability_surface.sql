-- Migration v10: capability catalog + task↔required attach (Phase 06 S01).
-- Additive only; do not rewrite 001–009. No ALTER on tasks.
-- Requirements live in task_capability_requirements (no new entity_links rels).

CREATE TABLE IF NOT EXISTS capabilities (
    id TEXT PRIMARY KEY,
    kind TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'UNKNOWN',
    body TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_capabilities_kind ON capabilities(kind);
CREATE INDEX IF NOT EXISTS idx_capabilities_status ON capabilities(status);

CREATE TABLE IF NOT EXISTS task_capability_requirements (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    capability_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    UNIQUE(task_id, capability_id)
);
CREATE INDEX IF NOT EXISTS idx_task_capability_requirements_task
    ON task_capability_requirements(task_id);
CREATE INDEX IF NOT EXISTS idx_task_capability_requirements_cap
    ON task_capability_requirements(capability_id);
