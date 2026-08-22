// Package compat is the Phase 08 compatibility + security checklist harness
// (deterministic planted checks — not commercial multi-model security theater).
//
// It re-proves S01 language-adapter API version, S02 path-local bind + exclusive
// trace.lock, and S03 migrate status / backup↔restore / local-auth fail-closed /
// no source BLOBs / schema through 027_harness_agents (no 028+), plus G19 (libraries do not import
// cmd/trace or cmd/trace-mcp) and no daemon/always-on HTTP as primary surface.
//
// Named test: TestCompatibilitySecurityChecklist
// Schema: schema-compat.json (schema_version 1)
// Metrics: temp metrics-compat.json under t.TempDir() (dry_run:false)
//
// Run (CGO required — analyzers + store surfaces):
//
//	CGO_ENABLED=1 go test ./evals/compat/... -count=1
//
// Named subset:
//
//	CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
//
// This checklist is not Phase 01 dry-run, not Gate C scores, and not a product
// feature surface. DR-HANDOFF after green VERIFY is "no successor" (roadmap ends
// at Phase 8) unless Notes explicitly promote a follow-on.
package compat
