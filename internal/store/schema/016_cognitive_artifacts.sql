-- Migration v16: cognitive artifacts (Phase 20 S02).
-- Additive only; do not rewrite 001–015. No ALTER on assumptions/decisions.
-- Uncertainty/Question, Hypothesis, Decision reconsideration child table.
-- Findings reuse discoveries — no Finding table.

CREATE TABLE IF NOT EXISTS uncertainties (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    severity TEXT NOT NULL DEFAULT 'INFO'
        CHECK (severity IN ('INFO', 'BLOCKING')),
    status TEXT NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'RESOLVED', 'SUPERSEDED')),
    kind TEXT NOT NULL DEFAULT ''
        CHECK (kind IN ('', 'risk', 'gap', 'unknown')),
    confidence REAL NOT NULL DEFAULT 0,
    source_type TEXT NOT NULL DEFAULT '',
    resolution TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    last_verified_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_uncertainties_status_severity
    ON uncertainties(status, severity);

CREATE TABLE IF NOT EXISTS hypotheses (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'CONFIRMED', 'REJECTED', 'SUPERSEDED')),
    confidence REAL NOT NULL DEFAULT 0,
    source_type TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    last_verified_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_hypotheses_status ON hypotheses(status);

CREATE TABLE IF NOT EXISTS decision_reconsiderations (
    id TEXT PRIMARY KEY,
    decision_id TEXT NOT NULL,
    trigger TEXT NOT NULL
        CHECK (trigger IN ('contradicted_effect', 'new_evidence', 'invalidated_assumption')),
    status TEXT NOT NULL DEFAULT 'FIRED'
        CHECK (status IN ('OPEN', 'FIRED')),
    reason TEXT NOT NULL DEFAULT '',
    related_type TEXT NOT NULL DEFAULT '',
    related_id TEXT NOT NULL DEFAULT '',
    reconsider_at TEXT NOT NULL,
    created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_decision_reconsiderations_decision
    ON decision_reconsiderations(decision_id);
