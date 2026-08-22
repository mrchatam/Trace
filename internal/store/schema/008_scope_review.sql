-- Migration v8: structured review residuals for scope-level review tracking.
-- Additive only; do not rewrite 001–007.
-- Scope links use existing entity_links (rel=review_judges_scope, to_type=plan_scope).

CREATE TABLE IF NOT EXISTS review_residuals (
    id TEXT PRIMARY KEY,
    review_id TEXT NOT NULL,
    code TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    severity TEXT NOT NULL DEFAULT 'INFO',
    status TEXT NOT NULL DEFAULT 'OPEN',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_review_residuals_review ON review_residuals(review_id);
CREATE INDEX IF NOT EXISTS idx_review_residuals_status ON review_residuals(status);
