# P39-S01-01 — Implement G3 harness orient

## Metadata
- id: P39-S01-01
- todo_ids: [P39-S01-01]
- role: implementer
- skills: [incremental-implementation, context-engineering, writing-for-agents]
- mcps: [user-trace]
- verification: automated

## Objective

Implement **G3**: MCP/harness orient playbook — ranked start-here, moat-first bootstrap messaging, Cursor 9/16 hygiene (G-006, G-010).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — **SoT** for locks
- [REMEDIATION-PLAN G3](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [INTAKE.md](../../INTAKE.md) — Q2 G2 compose-first in S01; Q5 9/16 resolution
- Live anchors (P39-S01-00 verified 2026-08-22, post-G1):
  - `internal/mcp/server.go:28–36` — `NewServer` passes `nil` ServerOptions (no Instructions yet)
  - `internal/mcp/server.go:227–235` — `RegisteredToolNames()` = 16 locked tools
  - `internal/mcp/server.go:153–161` — `trace_version` stale-process detection
  - `internal/mcp/server.go:67–75` — `trace_context` desc mentions optional `query` (G1 shipped)
  - `internal/mcp/tools_context.go:21` — `ContextInput.Query` optional (G1 shipped)
  - `internal/mcp/tools_loop.go:18–24` — `trace_loop` supports `action=gate`
  - `internal/install/cursor.go:12–13` — `CursorReloadTip` (DF-22/50)
  - `internal/install/bootstrap_hint.go:12–36` — plan bootstrap hint only (no moat line yet)
  - `internal/install/agents.go:17–57` — harness agent defaults + bootstrap hint call
  - `CONTRIBUTING.md:64–70` — agent workflow + MCP reload + `trace_version` (expand orient)
  - `go.mod` — `github.com/modelcontextprotocol/go-sdk v1.4.0`; `ServerOptions.Instructions` supported

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| GAP ids | G-006, G-010 |
| Primary deliverable | MCP server `Instructions` string (orient playbook) |
| Secondary | CONTRIBUTING + install stderr/doc hygiene for stale-server reload |
| Moat lead | `trace_tasks` → `trace_context`(+`query`) → `trace_loop` → `trace_review` → `trace_plan` |
| Gate | `trace_loop action=gate` (or CLI `trace loop gate`) before product edits |
| Compose-first (G2) | Instructions document ranked multi-tool read sequence — **not** unified explore |
| 9/16 hygiene | Instructions mention `trace_version` + rebuild/reload Cursor MCP |
| G1 dependency | **Shipped (S00)** — orient recipe references optional `trace_context.query` |
| Must not | Reduce to CG 1-tool; hide write tools; add 17th MCP tool |

## Touch-list (MCP → install → docs → tests)

| Step | File | Action |
|------|------|--------|
| 1 | `internal/mcp/instructions.go` | **Create** — `ServerInstructions() string` returning orient playbook markdown |
| 2 | `internal/mcp/server.go:31–34` | **Edit** — pass `&sdkmcp.ServerOptions{Instructions: ServerInstructions()}` as second arg to `sdkmcp.NewServer` (replace `nil`) |
| 3 | `internal/mcp/mcp_test.go` | **Edit** — add `TestServerInstructions*` (non-empty; moat tools; compose recipe; trace_version/reload) |
| 4 | `internal/install/cursor.go:12–13` | **Edit** (optional) — extend `CursorReloadTip` with 9/16 partial-tool-list note if not redundant with Instructions |
| 5 | `internal/install/bootstrap_hint.go:34` | **Edit** (optional) — moat-first one-liner before plan bootstrap hint |
| 6 | `CONTRIBUTING.md:64–70` | **Edit** — expand Agent workflow / MCP section: moat-first orient summary + reload steps + pointer to S02 dual-stack |
| 7 | `docs/rules/agent-loop-protocol.md` | **Edit** (optional cross-ref) — link MCP orient recipe if fits existing structure |

**Explicit non-touch:**

- Tool registration count/order (stay 16; no new `AddTool`)
- `internal/mcp/tools_*.go` handler logic (Instructions string only)
- `web/` GUI orient (G5 — Phase 40+)
- Bundled Codegraph MCP
- `trace_explore` implement

## Instructions content outline (minimum sections)

Implement as markdown inside `ServerInstructions()` — agent-facing, concise:

```markdown
## Trace MCP — start here (moat-first)

1. **Pick work:** trace_tasks → trace_context task_id=<uuid> [query=<optional agent query>]
2. **Deliberation loop:** trace_loop action=next|status → implement → trace_loop action=apply
3. **Gate before edits:** trace_loop action=gate (for=edit) — or CLI trace loop gate
4. **Review path:** trace_review before DONE; trace_transition with evidence
5. **Planning:** trace_plan action=bootstrap|create-coarse when goal lacks plan tree

## Read tools (compose-first — not CG single explore)

When task-scoped discovery needed, rank: trace_search → trace_why → trace_impact → trace_capability
Use progressive caps; never request full graph dump.

## Stale server hygiene (9/16)

After rebuilding trace-mcp: call trace_version; reload Cursor MCP / restart window.
Partial tool list (e.g. 9/16) = stale stdio process — not intentional surface reduction.

## Codegraph complement (optional)

For symbol-level code exploration in indexed repos, use separate codegraph MCP per CONTRIBUTING dual-stack section (Phase 39 S02).
Trace owns task loop + evidence; Codegraph owns code graph reads.
```

**String-match targets for tests:** `trace_tasks`, `trace_context`, `trace_loop`, `trace_review`, `trace_plan`, `trace_version`, `trace_search`, `9/16` or `stale`, `compose`.

## Implementation order

```text
1. internal/mcp/instructions.go — ServerInstructions()
2. internal/mcp/server.go — wire ServerOptions.Instructions in NewServer
3. internal/mcp/mcp_test.go — TestServerInstructionsNonEmpty, TestServerInstructionsMoatLead, TestServerInstructionsComposeRecipe, TestServerInstructionsStaleHygiene
4. CONTRIBUTING.md — orient + reload + dual-stack pointer (full dual-stack prose in S02)
5. Optional install stderr tweaks (cursor.go, bootstrap_hint.go)
6. go test ./internal/mcp/... ./internal/install/... -count=1
```

## Acceptance criteria (must pass)

| ID | Criterion | Evidence |
|----|-----------|----------|
| G3-A1 | MCP server exposes non-empty Instructions via go-sdk | `TestServerInstructionsNonEmpty` or initialize handler read |
| G3-A2 | Instructions lead with moat tools in order: tasks → context(+query) → loop → review → plan | `TestServerInstructionsMoatLead` string match |
| G3-A3 | Instructions include ranked compose-first read recipe (search → why → impact → capability) | `TestServerInstructionsComposeRecipe` |
| G3-A4 | Instructions mention `trace_version` + stale/reload hygiene (9/16) | `TestServerInstructionsStaleHygiene` |
| G3-A5 | CONTRIBUTING documents reload-after-rebuild + moat-first orient summary | Doc read `:64–70` region expanded |
| G3-A6 | 16 tools still registered unchanged | `TestToolNamesRegistered` green |

## Exit criteria

- [ ] G3-A1–A6 satisfied
- [ ] Board row → `done` with Notes (files + test command)

## Next

`P39-S01-02`
