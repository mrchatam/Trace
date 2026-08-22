// Package gitcli implements vcs.Repository using the host git CLI only.
//
// No libgit2 / go-git. Content is resolved via git show; a thin commit/path
// index lives in the project store (.trace/trace.db) and is refreshed
// incrementally via a durable watermark.
package gitcli
