// Package p0x is the deterministic P0-X evaluation harness (no LLM).
//
// Primary gate:
//
//	CGO_ENABLED=1 go test ./evals/p0x/... -count=1
//
// It copies fixtures/x0 to a temp dir, builds cmd/trace, runs the locked CLI
// walkthrough (init → seed import with an absolute seed path → index → why /
// context → one-file reindex), and asserts all 7 DR-P0X criteria.
package p0x
