// Package vcs defines the version-control adapter interface for Trace.
//
// Consumers (analyzers, retrieval, tests) depend on Repository — not on a
// concrete Git backend. The P0 production implementation is internal/gitcli
// (host git CLI only). Content is always resolved from Git; SQLite may hold
// thin commit/path references and a refresh watermark, never source blobs.
package vcs
