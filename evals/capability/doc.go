// Package capability is the capability-selection ablation planted harness (deterministic, no LLM).
//
// It plants selection / missing probes via live S01 capability APIs
// (UpsertCapability, RequireCapability, MissingCapabilities) plus compiler
// TaskContext packet attach (required_capabilities / missing_capabilities) and
// scores precision/recall against locked ground-truth assertions.
//
// Named test: TestPlantedCapabilitySelectionAblation
// Schema: schema-capability.json (schema_version 1)
// Metrics: temp metrics-capability.json under t.TempDir()
//
// Run (CGO-free — domain+store+compiler):
//
//	CGO_ENABLED=0 go test ./evals/capability/... -count=1
//
// Named subset:
//
//	CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
//
// Regression with honesty / replan / impact / p0x / x0:
//
//	CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./... -count=1
package capability
