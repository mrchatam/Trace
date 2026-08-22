// Package store owns per-project SQLite persistence under .trace/.
//
// Open(projectRoot) resolves filepath.Abs only (no parent walk-up), creates
// <absRoot>/.trace/ if needed, acquires exclusive .trace/trace.lock, optionally
// gates on .trace/access.token vs TRACE_ACCESS_TOKEN, opens .trace/trace.db via
// modernc.org/sqlite (no CGO), and applies embedded migrations. A second Open
// on the same abs root fails with ErrLocked until Close — one writer per
// project root; parallel agents use separate -C / worktree roots.
//
// OpenExisting(projectRoot) is Abs-only like Open but never creates .trace/ or
// trace.db. It requires a regular file <abs>/.trace/trace.db and returns
// ErrNotInitialized when the file (or the whole .trace/ dir) is missing. MCP
// CallTool uses OpenExisting; CLI init and other mkdir callers keep using Open.
//
// Backup/Restore snapshot trace.db only (VACUUM INTO / file install); restore
// rebinds projects.root_path to the target Abs root. access.token is excluded
// from backups by default.
//
// The store holds causal entities (goals, tasks, …), append-only events,
// File/Symbol/Import stubs that support per-file incremental replace, and an
// FTS5 lexical index (mig 004, tokenizer unicode61) over entity text / paths /
// symbol names. SyncEntityFTS/SyncFileFTS keep the index current on Upsert*;
// Open also backfills via RebuildFTS when fts_docs is empty but content exists
// (upgrade onto a pre-004 DB). Source file contents are never stored.
package store
