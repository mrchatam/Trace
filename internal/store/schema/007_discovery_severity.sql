-- 007_discovery_severity.sql — additive severity on discoveries (P03-S02).
-- Allowed values enforced in store/domain: INFO | PLAN_AFFECTING | BLOCKING.

ALTER TABLE discoveries ADD COLUMN severity TEXT NOT NULL DEFAULT 'INFO';
