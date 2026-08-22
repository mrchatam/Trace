// Package x0 is the Experiment X0 harness (B0 vs G1).
//
// Phase 01: dry-run metrics (`TestX0DryRunMetricsB0AndG1`, dry_run:true).
// Phase 02 Gate C: recorded (or live) understanding packs graded vs GT
// (`TestX0GateCRecordedMetrics`, dry_run:false, N≥3). Phase 01 dry-run ≠ Gate C pass.
//
//	CGO_ENABLED=1 go test ./evals/x0/... -count=1
//
// Keep evals/p0x and evals/honesty separate; do not merge honesty into X0.
package x0
