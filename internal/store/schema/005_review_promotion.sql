-- Migration v5: Review result for Claim→Evidence→Review→DONE promotion.
-- Additive only; do not rewrite 001_init.sql.
-- result: PASS | FAIL | UNCERTAIN | '' (open / unset).

ALTER TABLE reviews ADD COLUMN result TEXT NOT NULL DEFAULT '';
