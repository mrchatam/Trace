package install

// GapPassPrompt is the mandatory gap-pass instruction agents receive after a
// build session completes (INT-03). Wire this text into any agent rule or
// install target that surfaces post-build guidance.
//
// Trigger: end of build session (not async, not event-driven).
// FM-02: discoveries/decisions must be written before seed export --strict --enforce.
const GapPassPrompt = `After completing build work, run a mandatory gap pass:
1. Call ` + "`trace gap --task \"$TRACE_TASK_ID\"`" + ` (or ` + "`trace loop status`" + `) to surface open gaps.
2. For each BLOCKING gap: either fix inline or spawn a new task via ` + "`trace add task`" + `.
3. **Post-discovery nudge (FM-08):** after ` + "`trace add discovery`" + ` / MCP ` + "`trace_add`" + ` kind=discovery (especially BLOCKING), promote to a task via ` + "`trace add task --from-discovery <id>`" + ` or ` + "`trace loop apply`" + ` with ` + "`spawned_tasks[].discovery_id`" + ` **before product edits** — prefer the task/promotion path; do not discovery-only then edit.
4. After ` + "`trace seed import`" + `, read stdout ` + "`promotion_candidates`" + ` (also on ` + "`trace loop next`" + `) and promote or explicitly decline — do not invent task UUIDs.
5. **Write-before-export (required):** before any ` + "`trace seed export --strict --enforce`" + `, record ≥1 discovery OR ≥1 decision linked to ` + "`$TRACE_TASK_ID`" + ` (` + "`trace add discovery`" + ` / ` + "`trace add decision`" + ` with task links). Do not export first and backfill later — thin graphs (discoveries=0 decisions=0) fail ` + "`--strict --enforce`" + `.
6. Only after those writes: ` + "`trace seed export -o trace/graph.json --strict --enforce`" + `.
7. Do not mark the task DONE until ` + "`violations[]`" + ` is empty.`
