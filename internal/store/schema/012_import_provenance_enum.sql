-- Migration v12: harden imports.provenance enum (Phase 13 S03 / DF-64).
-- Rebuild + CHECK; heal empty/unknown → EXTRACTED on copy (migrate-only).
-- Do not rewrite 001–011.

CREATE TABLE imports_new (
    id TEXT PRIMARY KEY,
    file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    imported_path TEXT NOT NULL,
    symbol TEXT,
    provenance TEXT NOT NULL DEFAULT 'EXTRACTED'
        CHECK (provenance IN ('EXTRACTED', 'INFERRED', 'AMBIGUOUS'))
);

INSERT INTO imports_new (id, file_id, imported_path, symbol, provenance)
SELECT
    id,
    file_id,
    imported_path,
    symbol,
    CASE
        WHEN provenance IS NULL OR provenance = '' THEN 'EXTRACTED'
        WHEN provenance IN ('EXTRACTED', 'INFERRED', 'AMBIGUOUS') THEN provenance
        ELSE 'EXTRACTED'
    END
FROM imports;

DROP TABLE imports;
ALTER TABLE imports_new RENAME TO imports;

CREATE INDEX IF NOT EXISTS idx_imports_file_id ON imports(file_id);
