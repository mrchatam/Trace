-- Migration v6: Progressive plan hierarchy (Goal → phases → scopes + deep plans).
-- Additive only; do not rewrite 001–005.

CREATE TABLE IF NOT EXISTS plan_phases (
  id TEXT PRIMARY KEY,
  goal_id TEXT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL DEFAULT '',
  ord INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_plan_phases_goal ON plan_phases(goal_id, ord);

CREATE TABLE IF NOT EXISTS plan_scopes (
  id TEXT PRIMARY KEY,
  phase_id TEXT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL DEFAULT '',
  ord INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  auto_replan_count INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_plan_scopes_phase ON plan_scopes(phase_id, ord);

CREATE TABLE IF NOT EXISTS scope_deep_plans (
  id TEXT PRIMARY KEY,
  scope_id TEXT NOT NULL,
  content_json TEXT NOT NULL DEFAULT '{}',
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_scope_deep_plans_scope ON scope_deep_plans(scope_id, status);

CREATE TABLE IF NOT EXISTS goal_plan_state (
  goal_id TEXT PRIMARY KEY,
  current_scope_id TEXT,
  updated_at TEXT NOT NULL
);
