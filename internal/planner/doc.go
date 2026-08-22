// Package planner provides progressive coarse planning: Goal → phases → scopes,
// deep-planning only the current scope (+ one shallow lookahead).
//
// It is orchestration over store plan-hierarchy tables — not a domain entity dump.
// Callers supply structure; this package does not LLM-generate whole-goal backlogs.
package planner
