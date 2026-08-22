-- Migration v4: FTS5 lexical index over entity text, file paths, and symbol names.
-- Tokenizer: unicode61 (FTS5). No source-file body/BLOB columns (G1).
-- Content is maintained via store SyncEntityFTS / RebuildFTS (no triggers).

CREATE VIRTUAL TABLE IF NOT EXISTS fts_docs USING fts5(
    entity_type UNINDEXED,
    entity_id UNINDEXED,
    title,
    body,
    path,
    symbol_name,
    symbol_kind UNINDEXED,
    tokenize = 'unicode61'
);
