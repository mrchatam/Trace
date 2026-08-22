package mcp

// ServerInstructions returns the MCP orient playbook for connected agents.
func ServerInstructions() string {
	return `## Trace MCP — start here (moat-first)

1. **Pick work:** trace_tasks → trace_context task_id=<uuid> [query=<optional agent query>]
2. **Deliberation loop:** trace_loop action=next|status → implement → trace_loop action=apply
3. **Gate before edits:** trace_loop action=gate (for=edit) — or CLI trace loop gate
4. **Review path:** trace_review before DONE; trace_transition with evidence
5. **Planning:** trace_plan action=bootstrap|create-coarse when goal lacks plan tree

## Read tools (compose-first — not CG single explore)

When task-scoped discovery needed, rank: trace_search → trace_why → trace_impact → trace_capability
Use progressive caps; never request full graph dump.

## Optional convenience (after moat + compose-first)

trace_explore task_id=<uuid> [query=<optional>] — unified capped read compose (task packet + search + why + neighborhood).
Optional shortcut only; does not replace trace_tasks → trace_context moat path or manual compose for fine control.

## Stale server hygiene (9/17)

After rebuilding trace-mcp: call trace_version; reload Cursor MCP / restart window.
Partial tool list (e.g. 9/17) = stale stdio process — not intentional surface reduction.

## Codegraph complement (optional)

For symbol-level code exploration in indexed repos, use separate codegraph MCP per CONTRIBUTING dual-stack section (Phase 39 S02).
Trace owns task loop + evidence; Codegraph owns code graph reads.
`
}
