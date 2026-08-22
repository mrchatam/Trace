# P22-S05-02 — Review: search + changes

## Metadata
- id: P22-S05-02
- todo_ids: [P22-S05-02]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C29, C30, C37** via CLI and MCP (G19: no domain fork in MCP handlers).

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

1. **C29 (history questions):** `trace search` + `trace changes list|show` answer “what happened / what changed” without extending `trace why`.
2. **C30 (what changed):** `list`/`show` surface stored changes + paths; `compare` unchanged from S02.
3. **C37 (inspect historical evidence):** search returns FTS-indexed historical rows (changes, regressions, outcomes, reflections) — not task-only.
4. Grep MCP handlers — must call domain/retrieval, not duplicate SQL.
5. Limits: default **32**, cap **64** on search + changes list (match [00-PLANNER.md](00-PLANNER.md)).
6. Schema: **24** sql files, no 025+; compat PASS.
7. MCP catalog: **12** tools; `RegisteredToolNames` matches `TestToolNamesRegistered`.
8. Capability gating: `cli:search`, `cli:changes` list/show/compare respect `failCLIDenied`.
9. No source blobs in JSON (paths only on show).
10. S02 keepers: `TestChangesCompare`, `TestChangesCapture` still PASS.

## Spawn policy

If unmet: spawn **`P22-S05-02a` + `P22-S05-02b`**. Do not close with residuals.

## Re-run commands

```bash
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/mcp/... -count=1 -run 'TestCLISearch|TestCLIChanges|TestMCPSearch|TestMCPChanges|TestToolNamesRegistered|TestChangesCompare|TestChangesCapture'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/retrieval/... -count=1 -run TestSearch
ls internal/store/schema/*.sql | wc -l  # expect 24
```

## Exit criteria

- [ ] C29, C30, C37 closed or spawned
- [ ] Confidence **high** | **medium** (must spawn if medium+unmet)
- [ ] Board Notes: findings + confidence + checklist boxed when closed
