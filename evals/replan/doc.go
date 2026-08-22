// Package replan is the planted discovery→PlanChange replan demo (deterministic, no LLM).
//
// It proves severity-gated auto-replan and G16/DR-CHURN budget (N=5) against the live
// S01/S02 planner surface: PLAN_AFFECTING+ supersedes via ApplyDiscoveryReplan;
// INFO does not; budget fail-closed then recovers after AckReplan.
//
// Run (CGO-free — planner+domain+store):
//
//	CGO_ENABLED=0 go test ./evals/replan/... -count=1
package replan
