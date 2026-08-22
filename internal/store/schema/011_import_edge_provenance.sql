-- Migration v11: structural import edge provenance (Phase 12 S01).
-- Additive only; do not rewrite 001–010. Does not touch causal confidence columns.

ALTER TABLE imports ADD COLUMN provenance TEXT NOT NULL DEFAULT 'EXTRACTED';
