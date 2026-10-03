// Package impact is the Gate F preliminary planted impact harness (deterministic, no LLM).
//
// It plants conflict / clean probes via live S01 decision-impact APIs
// (AddImpactFinding, LinkDecisionTask, ImpactReport) and scores precision/recall
// against locked ground-truth assertions on ImpactReport (HasUnknown / Incomplete /
// Findings / OverallClass — never OverallClass alone when HasUnknown is required).
//
// Named test: TestPlantedImpactConflictsGateFPrelim
// Schema: schema-gate-f.json (schema_version 1)
// Metrics: temp metrics-gate-f.json under t.TempDir()
//
// Run (the test binary links the tree-sitter analyzer graph via shared test
// imports, so CGO_ENABLED=1 is required):
//
//	CGO_ENABLED=1 go test ./evals/impact/... -count=1
//
// Named subset:
//
//	CGO_ENABLED=1 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
//
// Regression with honesty / replan / p0x / x0:
//
//	CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./... -count=1
package impact
