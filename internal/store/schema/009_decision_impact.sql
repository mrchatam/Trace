-- Migration v9: decision impact findings + alternatives (manual/planted; DR-NOIMP).
-- Additive only; do not rewrite 001–008. No ALTER on decisions.
-- Affected work stays entity_links rel=decision_affects_task (no new rels this scope).

CREATE TABLE IF NOT EXISTS decision_impact_findings (
    id TEXT PRIMARY KEY,
    decision_id TEXT NOT NULL,
    impact_class TEXT NOT NULL,
    uncertainty TEXT NOT NULL DEFAULT 'UNKNOWN',
    kind TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    related_type TEXT NOT NULL DEFAULT '',
    related_id TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_decision_impact_findings_decision
    ON decision_impact_findings(decision_id);

CREATE TABLE IF NOT EXISTS decision_alternatives (
    id TEXT PRIMARY KEY,
    decision_id TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    is_recommended INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_decision_alternatives_decision
    ON decision_alternatives(decision_id);
