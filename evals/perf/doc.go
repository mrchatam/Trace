// Package perf is the Gate H planted performance-ladder harness (deterministic, no LLM).
//
// It plants synthetic size ladders (smoke / ~1k / ~10k LOC), indexes them via the
// live `trace` CLI, records per-rung timings and DB size, and asserts structural
// T0-skip + incremental isolation (+ optional Go adapter) against measure-then-
// threshold regression ceilings.
//
// Named test: TestPlantedPerfLadderGateH
// Schema: schema-gate-h.json (schema_version 1)
// Metrics: temp metrics-gate-h.json under t.TempDir() (dry_run:false)
//
// Run (CGO required — analyzers + CLI):
//
//	CGO_ENABLED=1 go test ./evals/perf/... -count=1
//
// Named subset:
//
//	CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
//
// Gate H is not Phase 01 dry-run, not Gate C scores, and not commercial multi-model
// / 100k–1M CI theater (those ladders remain deferred).
package perf
