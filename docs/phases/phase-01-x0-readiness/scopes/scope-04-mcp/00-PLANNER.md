# P01 / S04 / 00-PLANNER — MCP adapter

## Metadata
- id: P01-S04-00
- todo_ids: [P01-S04-00]
- role: planner
- skills: [planning-and-task-breakdown, mcp-builder]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-mcp.md` for **thin MCP adapter** after S01–S03 validate CLI context path (DR-AGENT / DR-SURFACE). Lock package path, tool surface parity, no domain fork. No product code in this row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-AGENT, DR-SURFACE
- Live CLI surface: why/context/transition/add/link/…
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Planner work
- Lock MCP package location (prefer `cmd/trace-mcp` or `internal/mcp` + thin main — pick one in 01).
- Tool list = subset of CLI semantics; library calls only.
- Thicken `01-mcp.md`; remind S05 VERIFY docs MCP tools.
- Sync SCOPE-TODOS.md.

## Depends-on
- **S01–S03** all `done` (board order). MCP is last product scope before VERIFY.
- **S03 note (P01-S03-00):** X0 dry-run lives in **`evals/x0`** (B0/G1 via CLI `why`/`context`). MCP must **not** become required for X0; keep CLI path green (`evals/x0` + `evals/p0x` + `evals/honesty`). Tool parity may include context/why semantics, but Gate C / X0 remain CLI-capable (DR-AGENT).

## Exit criteria
- [x] `01-*` runnable without guessing
- [x] SCOPE-TODOS + TODO.md Notes updated
- [x] No product Go in this planner row

## Minimal todos
- [x] Inspect CLI public surface
- [x] Thicken 01 prompt (+ 02 review focus)
- [x] Sync todos + light S05 notes
