-- Migration v22: code graph relationships (Phase 22 S01).
-- Additive only; do not rewrite 001–021. This row writes rel=validates only;
-- other rel values exist so S01-03/05 need no 023.

CREATE TABLE IF NOT EXISTS code_edges (
    id TEXT PRIMARY KEY,
    from_file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    from_symbol_id TEXT REFERENCES symbols(id) ON DELETE SET NULL,
    to_file_id TEXT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    to_symbol_id TEXT REFERENCES symbols(id) ON DELETE SET NULL,
    rel TEXT NOT NULL
        CHECK (rel IN ('validates', 'contains_module', 'exports_api', 'architectural_boundary', 'depends_on')),
    provenance TEXT NOT NULL
        CHECK (provenance IN ('EXTRACTED', 'INFERRED', 'AMBIGUOUS'))
);

-- Unique identity uses IFNULL so NULL symbol ids collapse with empty (file-level edges).
CREATE UNIQUE INDEX IF NOT EXISTS idx_code_edges_unique
    ON code_edges(
        from_file_id,
        IFNULL(from_symbol_id, ''),
        to_file_id,
        IFNULL(to_symbol_id, ''),
        rel
    );

CREATE INDEX IF NOT EXISTS idx_code_edges_from_file ON code_edges(from_file_id);
CREATE INDEX IF NOT EXISTS idx_code_edges_to_file ON code_edges(to_file_id);
CREATE INDEX IF NOT EXISTS idx_code_edges_to_symbol ON code_edges(to_symbol_id);
CREATE INDEX IF NOT EXISTS idx_code_edges_rel ON code_edges(rel);
