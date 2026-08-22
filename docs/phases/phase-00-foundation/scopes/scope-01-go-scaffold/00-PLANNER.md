# P00 / S01 / 00-PLANNER — Go scaffold

## Metadata
- id: P00-S01-00
- todo_ids: [P00-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize `01-go-module-scaffold.md` against live repo. Lock Go layout defaults. No product features.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md)
- Sibling `01-*.md`, `02-*.md`; [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner).

## Planner work
- Inspect repo (docs-only vs any existing Go files).
- Thicken `01-go-module-scaffold.md` locked defaults + exit criteria.
- Note downstream effects on S02+ (paths only) in upcoming prompts if needed.

## Exit criteria
- [x] `01-*` runnable without stack guessing
- [x] SCOPE-TODOS.md updated
- [x] TODO.md Notes for P00-S01-00
- [x] No product code beyond planning docs

## Minimal todos
- [x] Repo inventory — docs-only; Go 1.24.2 on host; no `*.go`/`go.mod`; no `.gitignore` yet
- [x] Thicken 01 prompt — full taxonomy + exit criteria; `compiler` vs `contextx` locked
- [x] Sync SCOPE-TODOS + TODO.md + path locks on S02–S07 / phase README / S01-02 review focus
