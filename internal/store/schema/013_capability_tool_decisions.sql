-- Migration v13: capability tool-decision audit (Phase 14 S02).
-- Additive only; do not rewrite 001–012.
-- Graduated allowlist statuses: AUTO_ALLOWED | PENDING | ALLOWED | DENIED.

CREATE TABLE IF NOT EXISTS capability_tool_decisions (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    decision TEXT NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_capability_tool_decisions_decision
    ON capability_tool_decisions(decision);
