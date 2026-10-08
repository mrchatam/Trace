# Phase 45 / S03 / Implement — Error quality

## Row `P45-S03-01` — **done** 2026-10-04. Issues [#155](https://github.com/mrchatam/Trace/issues/155), [#156](https://github.com/mrchatam/Trace/issues/156).

## 1. New file — `internal/mcp/errors.go`

### `linkRelShape` (`errors.go:14`)

Per-rel endpoint expectations. Covers both hyphen and underscore spellings because Phase 44 added `api-contract`/`api_contract` and `scope-member`/`scope_member` aliases:

```
goal-task               from=goal              to=task
decision-task           from=decision          to=task
discovery-plan-change   from=discovery         to=plan_change
discovery-mentions-task from=discovery         to=task
claim-evidence          from=claim             to=evidence
scope-member/_member    from=any linkable      to=scope
api-contract/_contract  from=task              to=task
implements              from=task or decision  to=task
blocks                  from=task              to=task
```

### `wrapToolErr(tool string, err error) error` (`errors.go:33`)

```go
if errors.As(err, &nf)  → fmt.Errorf("%s: no %s %q", tool, nf.Entity, nf.ID)
else                    → fmt.Errorf("%s: %w", tool, err)
```

The entity type comes from the **typed** error built in S02, so `GetGoal` on a missing id says `goal`, not `task` and not `sql: no rows in result set`. Non-not-found errors pass through unchanged under the tool prefix.

### `wrapLinkErr(rel, fromID, toID string, err error) error` (`errors.go:46`)

Renders the parameter the caller actually got wrong, not the order the code happened to query in:

```
trace_link: rel blocks: no task "9a1…" for to (rel blocks expects from=task to=task)
```

Logic: match `nf.ID` against the trimmed `fromID` ⇒ role `from`; against `toID` ⇒ role `to`. Unknown rel falls back to `wrapToolErr`.

## 2. Wired at five sites

| Site | Line |
|------|------|
| `toolLink` | `tools_write.go:256` — `wrapLinkErr(in.Rel, in.From, in.To, err)` |
| `trace_transition` | `tools_write.go:319` |
| `trace_review` (set) | `tools_write.go:411` |
| `trace_why` | `tools_why.go:43`, `:47` |
| `trace_changes show` | `tools_changes.go:168`, `:172` |

## 3. `loop.next` returns the gate envelope (#156)

`internal/mcp/tools_loop.go`:

```go
const mcpGateSchemaVersion = "trace.loop.gate.v1"          // :30  (unchanged — reused)
```

- `mcpGateEnvelope` gains a **`Remedy string \`json:"remedy,omitempty"\``** field (`:39`) — a blocked `next` needs to tell the caller what to do, and the existing gate did not carry that.
- `loopNextPrecondition` (`:73`) — `reasonCode` + `remedy`
- `checkLoopNextPreconditions(st, taskID)` (`:81`) — inspects the seed task **before** packet assembly
- `loopNextPreconditionResult(taskID, pre)` (`:103`) — renders the block using the gate envelope with `for: "next"`
- `loopNext` (`:140`) returns that envelope; the unknown-task case at `:162` goes through the same path

`reason_code` values: **`goal_link_missing`** (the task has no `goal_id`) and **`seed_task_not_found`**.

Imports added to `tools_loop.go`: `deliberation`, `store`.

The result is that `action=next` and `action=gate` speak one schema. An agent parsing `trace.loop.gate.v1` handles both.

## 4. Acceptance — all PASS

`internal/mcp/errors_test.go`:

| Test | Asserts |
|------|---------|
| `TestLinkErrorsNameRelVocabulary` | #155 — `trace_link` errors name the rel's expected `from`/`to`, not driver text |
| `TestToolErrorsNeverLeakDriverText` | #155 — `assertNoDriverText` sweeps every wrapped error for `sql: no rows`, driver type names, etc. |
| `TestLoopNextReturnsGateEnvelopeOnPrecondition` | #156 — missing `goal_id` ⇒ gate envelope with `reason_code: goal_link_missing` + a remedy |
| `TestLoopNextUnknownTaskReturnsGateEnvelope` | #156 — unknown task ⇒ gate envelope with `reason_code: seed_task_not_found` |

Helpers: `assertNoDriverText`, `mustDecodeID`, `seedLinkedTask`.

One redundant happy-path test was **removed** rather than kept — `TestMCPLoopNext` already covers the success path, and duplicating it would have given a false sense of coverage.

## 5. Boundaries honoured

- No new schema version; no `trace.loop.next.v1`.
- The wrapper is at the tool boundary. `store` error text is unchanged for in-process callers — only what crosses the MCP boundary is rewritten.
- Non-not-found errors are not rewritten beyond a prefix, so unrelated diagnostics keep their detail.

## Next

`P45-S03-02`