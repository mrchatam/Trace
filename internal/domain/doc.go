// Package domain is the canonical work/causal library API for Trace.
//
// It persists only through *store.Store in projectRoot/.trace/trace.db
// (no second database, no in-memory-only causal facts). Domain never calls
// store.Open — callers construct the store and pass it to New.
//
// Provenance (G5): creates set source_type (default USER_ASSERTED), confidence,
// status (ACTIVE|STALE|SUPERSEDED), and timestamps. Empty titles are rejected.
// Use MarkStale rather than silently rewriting history.
//
// Task work_state is a separate column from provenance status. Vocabulary:
// PENDING | IN_PROGRESS | AWAITING_REVIEW | BLOCKED | FAILED | DONE | STALE | SKIPPED.
//
// Required entity_links rel values:
//   - decision_affects_task
//   - discovery_causes_plan_change
//   - review_judges_task (Review → Task; DONE requires linked PASS unless escape hatch)
//   - review_judges_scope (Review → plan_scope; recording only, does not gate DONE)
//
// Optional documented rels: claim_has_evidence, review_cites_evidence,
// discovery_mentions_task (DF-42; discovery→task attribution for multi-goal DPC).
// Goal→Task is persisted via tasks.goal_id only (event payload may say goal_has_task).
//
// Events (DR-EVT): entity.created, entity.linked, task.transition, review.result,
// deliberation.transition (Phase 20 S01; payload on seed task).
// Residuals are structured tracking hooks on reviews (not VerifiedFact).
//
// Decision impact (manual/planted; DR-NOIMP): findings + alternatives tables via
// AddImpactFinding / ImpactReport. No new entity_links rels — keep decision_affects_task.
//
// DONE policy: TransitionTask → DONE iff AllowDoneWithoutReview (hatch) or
// (no linked review_judges_task FAIL ∧ linked PASS ∧ AllowOperatorDone).
// Sibling FAIL blocks even when another PASS exists (DF-43); UNCERTAIN/empty
// do not. Hatch bypasses FAIL+PASS+operator. Actor string is never auth;
// --as-operator / AllowOperatorDone is a conscious flag ≠ verified identity
// (DF-44). Leaving DONE invalidates linked PASS reviews → UNCERTAIN.
// MissingCapabilities nonempty blocks any transition unless
// AllowMissingCapabilities. EvidenceIDs alone do not authorize DONE.
package domain
