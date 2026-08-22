-- Migration v27: harness agent catalog (Phase 22 S09).
-- Additive only; do not rewrite 001–026.

CREATE TABLE IF NOT EXISTS harness_agents (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,              -- agent:<name>
    title TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    subagent_type TEXT NOT NULL DEFAULT '', -- harness profile string e.g. performance-reviewer
    deliberation_phases TEXT NOT NULL DEFAULT '[]',  -- JSON array of phase strings
    task_keywords TEXT NOT NULL DEFAULT '[]',        -- JSON array of lowercase keywords
    recommend_subagent INTEGER NOT NULL DEFAULT 0,   -- 1 when profile suits fresh-subagent review
    registry_source TEXT NOT NULL DEFAULT 'bundled', -- bundled|user|host
    registry_version TEXT NOT NULL DEFAULT '',
    external_url TEXT,                               -- nullable; no fetch in P22
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_harness_agents_slug ON harness_agents(slug);
CREATE INDEX IF NOT EXISTS idx_harness_agents_subagent_type ON harness_agents(subagent_type);

CREATE TABLE IF NOT EXISTS harness_agent_requirements (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL,
    required_capability_slug TEXT NOT NULL,  -- skill:…, mcp:…, hook:…, tool:…
    created_at TEXT NOT NULL,
    UNIQUE(agent_id, required_capability_slug)
);
CREATE INDEX IF NOT EXISTS idx_harness_agent_requirements_agent
    ON harness_agent_requirements(agent_id);
CREATE INDEX IF NOT EXISTS idx_harness_agent_requirements_slug
    ON harness_agent_requirements(required_capability_slug);
