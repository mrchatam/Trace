-- Migration v23: graph sync watermark (symbol/file index at HEAD), separate from vcs_meta.

CREATE TABLE IF NOT EXISTS graph_sync_state (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    last_indexed_commit TEXT NOT NULL DEFAULT '',
    last_indexed_at TEXT NOT NULL DEFAULT '',
    hook_installed INTEGER NOT NULL DEFAULT 0
);

INSERT OR IGNORE INTO graph_sync_state(id) VALUES (1);
