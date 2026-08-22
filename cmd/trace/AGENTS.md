# begin-trace-enforcement
## Trace enforcement (harness)

When a Trace seed task is active, set `TRACE_TASK_ID` to its UUID.

- **Before edits:** `trace loop gate --task "$TRACE_TASK_ID" --for edit` (exit 0 = proceed).
- **Coarse plan:** New goal without progressive plan — `trace plan create-coarse` + `set-current` + `deep`, or MCP `trace_plan`; recovery: `trace plan bootstrap --goal <id>`.
- **Before DONE:** `trace loop status --task "$TRACE_TASK_ID"` — resolve non-empty `violations[]`.
- **Opt-in strict:** `--enforce` on `trace transition … --to DONE`; `trace seed export --strict --enforce` for CI.
- **Config:** `.trace/config.json` → `{ "enforce": "off"|"warn"|"strict" }` (default off).

Product design: docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md

After completing build work, run a mandatory gap pass:
1. Call `trace gap --task "$TRACE_TASK_ID"` (or `trace loop status`) to surface open gaps.
2. For each BLOCKING gap: either fix inline or spawn a new task via `trace add task`.
3. **Post-discovery nudge (FM-08):** after `trace add discovery` / MCP `trace_add` kind=discovery (especially BLOCKING), promote to a task via `trace add task --from-discovery <id>` or `trace loop apply` with `spawned_tasks[].discovery_id` **before product edits** — prefer the task/promotion path; do not discovery-only then edit.
4. After `trace seed import`, read stdout `promotion_candidates` (also on `trace loop next`) and promote or explicitly decline — do not invent task UUIDs.
5. **Write-before-export (required):** before any `trace seed export --strict --enforce`, record ≥1 discovery OR ≥1 decision linked to `$TRACE_TASK_ID` (`trace add discovery` / `trace add decision` with task links). Do not export first and backfill later — thin graphs (discoveries=0 decisions=0) fail `--strict --enforce`.
6. Only after those writes: `trace seed export -o trace/graph.json --strict --enforce`.
7. Do not mark the task DONE until `violations[]` is empty.
Parent orchestrator TRACE_TASK_ID ownership (INT-04 / FM-04):

1. The parent orchestrator MUST set TRACE_TASK_ID to the active seed task UUID
   before any product-code edit and before delegating work to any subagent.
2. Parent owns the Trace graph for the active task: gap pass, discoveries, and
   decisions. Do not complete an edit path by offloading graph-only work to
   workers while the parent edits without TRACE_TASK_ID / loop gate.
3. Worker inheritance: before each worker, set TRACE_TASK_ID and
   TRACE_PROJECT_ROOT; put workspace path + task UUID in every worker prompt;
   workers must export the same env before product edits. Do not assume Cursor
   Multitask inherits parent env automatically.
4. preToolUse deny fires when a Write/edit is attempted without an active
   TRACE_TASK_ID under enforce=strict — enforced via CursorLoopGateHookScript
   (.cursor/hooks/trace-loop-gate.sh). Option A applies per process to that
   process's TRACE_TASK_ID.
5. failClosed (Option A): When TRACE_TASK_ID is absent AND .trace/config.json
   enforce=strict, deny the edit rather than allowing untracked work.
   off/warn/missing/invalid enforce still allow (default-off preserved).
6. Multitask limit: Trace cannot product-enforce worker env inheritance or
   detect "parent orchestrator" (Option B deferred). Rules + Option A hook are
   the harness choke points; orchestrators still verify board order.

Implementation path: see CursorLoopGateHookScript in internal/install/enforcement.go.
# end-trace-enforcement
