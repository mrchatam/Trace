package mcp

import "fmt"

// Tool descriptions are generated from the shared vocabularies in schema.go so
// a description and its input-schema enum cannot drift apart (#159). The drift
// guard is TestToolDescriptionsCoverSchemaEnums.

var (
	descWhy = "Explain why an entity exists via causal neighborhood (retrieval.Why). Mirrors `trace why <type> <id>`."

	descContext = fmt.Sprintf(
		"Compile a task context packet (compiler.TaskContext / ExpandContext). Mirrors `trace context`. "+
			"depth %s; max_layer %s; format %s (default json; prefer json, both emits a loud warning and ~2x tokens); "+
			"optional query merges agent FTS hits; task_id required.",
		join([]string{"1", "2"}), join([]string{"1", "2", "3"}), join(contextFormats))

	descAdd = fmt.Sprintf(
		"Create a domain entity (%s). status is provenance (%s), default ACTIVE. Prefer the task/promotion path "+
			"over discovery-only edits: after a BLOCKING discovery, promote with trace_add kind=task or loop apply "+
			"spawned_tasks with discovery_id before product edits — do not discovery-only then edit. An id that "+
			"already exists is refused unless replace=true. Mirrors `trace add`.",
		join(addKinds), join(provenanceStatus))

	descLink = fmt.Sprintf(
		"Link two entities (%s). from and to must differ. Mirrors `trace link`.", join(linkRels))

	descTransition = fmt.Sprintf(
		"Transition a task work_state (to_state=%s). DONE requires linked Review PASS + as_operator with no linked "+
			"FAIL, or allow_done hatch (returns warning); as_operator is a conscious claim (flag≠identity / not "+
			"verified operator identity); evidence_ids alone do not authorize DONE; actor is not authorization; "+
			"missing caps fail-closed unless allow_missing_caps; allow_done does not bypass missing caps. "+
			"Mirrors `trace transition`.",
		join(workStates))

	descReview = fmt.Sprintf(
		"Create a review (optional task link) or set result %s. action=%s. Mirrors `trace review`.",
		join(reviewResults), join(reviewActions))

	descTasks = fmt.Sprintf(
		"List tasks (bounded). Returns {items,count,truncated,next_cursor?}. Optional goal_id, work_state (%s), "+
			"work_states, limit (default 50, cap 500), cursor. all=true is hard-capped at 500 (prefer cursor). "+
			"Ordered by created_at ASC then id ASC; rows created in the same second tie-break on id, so use cursor for "+
			"stable paging. Mirrors `trace tasks`.",
		join(workStates))

	descCapability = fmt.Sprintf(
		"Capability catalog: action=%s. kind=%s; status=%s. missing requires task or task_id (list via trace_tasks). "+
			"Mirrors `trace capability`.",
		join(capabilityActions), join(capabilityKinds), join(capabilityStatus))

	descImpact = fmt.Sprintf(
		"Decision impact: action=%s; finding/alternative op=%s. walk reports graph_sync_honesty. Mirrors `trace impact`.",
		join(impactActions), join(impactOps))

	descVersion = "Return live MCP process identity {ok,name,version,project_initialized,project_root} so agents can detect a stale " +
		"stdio server after rebuild/install. Answers without an initialized store so it still works on an uninitialized " +
		"or locked project. Does not run the capability gate."

	descSearch = fmt.Sprintf(
		"Full-text search over indexed entities. Args: query (required), limit (optional, default 32, cap 64), "+
			"matcher=%s. matcher=and (default) AND-joins the query terms; matcher=ops honors AND|OR|NOT|NEAR, "+
			"\"phrase\" and a trailing *. When a query contains operators that the active mode did not use, the "+
			"response carries query_note. Reports graph_sync_honesty. Mirrors `trace search`.",
		join(matchers))

	descChanges = fmt.Sprintf(
		"Change history: action=%s. list: optional task_id, limit; show: change_id; compare: from, to. "+
			"Returns the standard list envelope {ok,items,count,truncated,limit}. Mirrors `trace changes`.",
		join(changeActions))

	descRegressions = fmt.Sprintf(
		"Regression history: action=%s. Optional task_id, change_id, limit. Returns the standard list envelope "+
			"{ok,items,count,truncated}. Mirrors `trace regressions list`.",
		join(regressionAction))

	descLoop = fmt.Sprintf(
		"Deliberation loop: action=%s. next/status require task_id; gate requires task_id and optional for=%s "+
			"(default edit); apply requires envelope (trace.loop.apply.v1). A blocked next returns the same "+
			"trace.loop.gate.v1 envelope as gate. Mirrors `trace loop`.",
		join(loopActions), join(loopGateFors))

	descAgents = fmt.Sprintf(
		"Harness agent catalog: action=%s. describe requires slug; recommend requires task_id or phase (+ optional "+
			"keywords). Returns the standard list envelope. An empty result carries the install hint. Mirrors `trace agents`. "+
			"Recommend-only — no spawn.",
		join(agentActions))

	descPlan = fmt.Sprintf(
		"Progressive planner: action=%s. Bootstrap recovers PlanExists from plan_changes. Mirrors `trace plan`.",
		join(planActions))

	descExplore = fmt.Sprintf(
		"Unified task-aware capped explore: task context packet + search + why on top hits + neighborhood. "+
			"task_id required; limit default 32 cap 64; max_nodes default 100; depth default 2; matcher=%s. "+
			"Optional query merges via G1. Reports graph_sync_honesty. Mirrors `trace explore`.",
		join(matchers))
)
