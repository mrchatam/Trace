# P00 / 00-PHASE-PLANNER — Foundation phase scaffold

## Metadata
- id: P00-00
- todo_ids: [P00-00]
- role: planner
- skills: [planning-and-task-breakdown, grilling]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Light replan of **Phase 00** against `docs/init/*` and current repo. Improve scope READMEs/order, ensure each scope’s prompts are coherent stubs with correct dependencies, and record phase locks. **Do not** implement Go product code. **Do not** deep-write every implement prompt — scope `00-PLANNER` rows finalize those.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [phase README](./README.md)
- [docs/init/C_FIRST_SCOPE.md](../../init/C_FIRST_SCOPE.md)
- [docs/init/D_DECISION_REGISTER.md](../../init/D_DECISION_REGISTER.md)
- [docs/TODO.md](../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner only).

## Planner work
1. Confirm S01–S09 match P0-X critical path (scaffold→store→vcs→analyzers→causal→retrieval/context→cli→fixture/harness→verify).
2. Check cross-scope assumptions (module path, `.trace/`, tree-sitter, git CLI adapter, incremental law).
3. Patch **upcoming** scope prompt stubs if gaps/conflicts found.
4. Update phase README if order/locks change.
5. Sync `docs/TODO.md` Notes for this row.

## Locked phase defaults (do not weaken)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| P0-X bar | 7/7 including incremental |
| MCP/daemon | Forbidden this phase |
| Review policy | Every scope gets `02-review` before next scope’s implement wave completes order on the board |

## Exit criteria
- [x] Phase README accurate
- [x] Each scope has 00/01/02 prompts present
- [x] TODO.md Phase 00 section matches folders
- [x] No product code
- [x] Board row `P00-00` → done + note

## Minimal todos
- [x] Inventory phase-00 tree vs init plan
- [x] Fix gaps in upcoming stubs only
- [x] Update TODO.md Notes
