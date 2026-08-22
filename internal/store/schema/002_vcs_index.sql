-- Migration v2: thin VCS commit/path index + refresh watermark.
-- References only: OIDs, times, short subjects, changed paths/status.
-- Never store source blobs or permanent full diff/patch bodies.

CREATE TABLE IF NOT EXISTS vcs_commits (
    oid TEXT PRIMARY KEY,
    parent_oids TEXT NOT NULL DEFAULT '',
    committed_at TEXT NOT NULL,
    subject TEXT NOT NULL DEFAULT '',
    seq INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_vcs_commits_seq ON vcs_commits(seq);

CREATE TABLE IF NOT EXISTS vcs_commit_paths (
    commit_oid TEXT NOT NULL REFERENCES vcs_commits(oid) ON DELETE CASCADE,
    path TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (commit_oid, path)
);

CREATE TABLE IF NOT EXISTS vcs_meta (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_vcs_commit_paths_path ON vcs_commit_paths(path);
CREATE INDEX IF NOT EXISTS idx_vcs_commits_committed_at ON vcs_commits(committed_at);
