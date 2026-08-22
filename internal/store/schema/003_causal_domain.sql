-- Migration v3: Task work_state (separate from provenance status) + entity_links.
-- No source-content BLOBs.

ALTER TABLE tasks ADD COLUMN work_state TEXT NOT NULL DEFAULT 'PENDING';

CREATE TABLE IF NOT EXISTS entity_links (
    id TEXT PRIMARY KEY,
    from_type TEXT NOT NULL,
    from_id TEXT NOT NULL,
    rel TEXT NOT NULL,
    to_type TEXT NOT NULL,
    to_id TEXT NOT NULL,
    source_type TEXT NOT NULL DEFAULT '',
    confidence REAL NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    UNIQUE(from_type, from_id, rel, to_type, to_id)
);

CREATE INDEX IF NOT EXISTS idx_entity_links_from ON entity_links(from_type, from_id);
CREATE INDEX IF NOT EXISTS idx_entity_links_to ON entity_links(to_type, to_id);
