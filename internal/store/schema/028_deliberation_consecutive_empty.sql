-- Migration v28: consecutive empty-apply counter on deliberation_state (Phase 26 S03).
-- Additive only; do not rewrite 001–027.

ALTER TABLE deliberation_state ADD COLUMN consecutive_empty_applies INTEGER NOT NULL DEFAULT 0;
