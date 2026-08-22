-- Migration v14: harden capability_tool_decisions.decision enum (Phase 16 S02 / DF-75 + DF-78).
-- Rebuild + CHECK; heal empty/unknown/YOLO → PENDING on copy (migrate-only).
-- Canonicalize exact builtin MCP Names to mcp:<Name>; fold dual rows fail-closed
-- (DENIED > PENDING > ALLOWED > AUTO_ALLOWED). Do not rewrite 001–013.

CREATE TABLE capability_tool_decisions_new (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    decision TEXT NOT NULL
        CHECK (decision IN ('AUTO_ALLOWED','PENDING','ALLOWED','DENIED')),
    reason TEXT NOT NULL DEFAULT '',
    actor TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE capability_tool_decisions_fold (
    id TEXT NOT NULL,
    slug TEXT NOT NULL,
    decision TEXT NOT NULL,
    reason TEXT NOT NULL,
    actor TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    prio INTEGER NOT NULL
);

INSERT INTO capability_tool_decisions_fold (id, slug, decision, reason, actor, created_at, updated_at, prio)
SELECT
    id,
    slug,
    decision,
    reason,
    actor,
    created_at,
    updated_at,
    CASE decision
        WHEN 'DENIED' THEN 4
        WHEN 'PENDING' THEN 3
        WHEN 'ALLOWED' THEN 2
        WHEN 'AUTO_ALLOWED' THEN 1
        ELSE 0
    END
FROM (
    SELECT
        id,
        CASE
            WHEN slug IN (
                'trace_why', 'trace_context', 'trace_add',
                'trace_link', 'trace_transition', 'trace_review',
                'trace_tasks', 'trace_capability', 'trace_version'
            ) THEN 'mcp:' || slug
            ELSE slug
        END AS slug,
        CASE
            WHEN decision IS NULL OR TRIM(decision) = '' THEN 'PENDING'
            WHEN decision IN ('AUTO_ALLOWED','PENDING','ALLOWED','DENIED') THEN decision
            ELSE 'PENDING'
        END AS decision,
        reason,
        actor,
        created_at,
        updated_at
    FROM capability_tool_decisions
);

INSERT INTO capability_tool_decisions_new (id, slug, decision, reason, actor, created_at, updated_at)
SELECT id, slug, decision, reason, actor, created_at, updated_at
FROM capability_tool_decisions_fold AS f
WHERE rowid = (
    SELECT f2.rowid
    FROM capability_tool_decisions_fold AS f2
    WHERE f2.slug = f.slug
    ORDER BY f2.prio DESC, f2.updated_at DESC, f2.rowid ASC
    LIMIT 1
);

DROP TABLE capability_tool_decisions_fold;
DROP TABLE capability_tool_decisions;
ALTER TABLE capability_tool_decisions_new RENAME TO capability_tool_decisions;

CREATE INDEX IF NOT EXISTS idx_capability_tool_decisions_decision
    ON capability_tool_decisions(decision);
