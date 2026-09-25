-- Migration v29: thin graph scope records (Phase 44 S02).
-- Additive only; do not rewrite 001–028.
-- plan_scopes (v6) remain planner hierarchy — NOT this table.

CREATE TABLE IF NOT EXISTS scopes (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    kind TEXT NOT NULL
        CHECK (kind IN ('feature', 'layer', 'business')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_scopes_kind ON scopes(kind);
