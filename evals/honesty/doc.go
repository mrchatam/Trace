// Package honesty is the named H5 partial honesty artifact (deterministic, no LLM).
//
// It proves fail-closed promotion against the live Claim/Review surface:
// a planted false/incomplete completion claim cannot reach Task DONE without
// an independent Review result=PASS plus AllowOperatorDone (DF-17), with any
// prior linked FAIL explicitly superseded (DF-43). Escape hatch
// AllowDoneWithoutReview is never used in TestHonestyFailClosedPlantedClaim
// (Paths A/B/C).
//
// Gate G preliminary (Phase 04 S02): TestHonestyEscapeRateGateGPrelim plants
// caught + escape cases, tallies OPEN POLICY_EXCEPTION residuals via S01
// scope-review hooks, and writes schema-valid metrics-gate-g.json (see
// schema-gate-g.json). The hatch is counted as an escape only in that report.
//
// Run both tests (CGO-free — domain+store+planner):
//
//	CGO_ENABLED=0 go test ./evals/honesty/... -count=1
//
// Named subset:
//
//	CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
//
// Regression with P0-X:
//
//	CGO_ENABLED=1 go test ./evals/p0x/... ./evals/honesty/... -count=1
package honesty
