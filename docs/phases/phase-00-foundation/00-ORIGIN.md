# P00-ORIGIN — Harness scaffold + minimal phase plan

## Metadata
- id: P00-ORIGIN
- todo_ids: [P00-ORIGIN]
- role: planner
- skills: [planning-and-task-breakdown, brainstorming, grilling]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated
- hooks: []

## Objective

Create Trace’s Perfect Planner–style harness and a **minimal** Phase 00 / Phase 01 scaffold connected to `docs/init/*`. **No Go product implementation** in this prompt (no `cmd/trace` business logic beyond empty stubs if already present).

## References
- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [docs/init/README.md](../../init/README.md)
- [docs/init/A_PROJECT_PLAN.md](../../init/A_PROJECT_PLAN.md)
- [docs/init/C_FIRST_SCOPE.md](../../init/C_FIRST_SCOPE.md)

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Deliverables
- `AGENTS.md`
- `docs/rules/{agent-loop-protocol,project-rules,skills-map}.md`
- `docs/TODO.md` master board
- `docs/phases/phase-00-foundation/**` minimal prompts for all P0-X scopes
- `docs/phases/phase-01-x0-readiness/README.md` stub
- Pointers from `docs/init/` that execution SoT is `docs/TODO.md`

## Locked Trace policies (encode in rules)
1. Implementers: board status + notes only; reviewers/planners do real backlog changes  
2. Forward-only immutability of `done` rows (best-effort)  
3. Implement → review quality loops (not deferred mega-review batches)  
4. Phase planner on every phase (light scaffold, not deep tasking)  
5. Orchestrator: one subagent per todo; re-read board after each  
6. Session: Agent mode → questions if unclear → Plan mode → execute  
7. Human verification where marked  

## Exit criteria
- [ ] Harness rules + AGENTS.md exist and match Trace policies above  
- [ ] `docs/TODO.md` lists Phase 00 rows in order; P00-ORIGIN can be marked done  
- [ ] Phase 00 scopes S01–S09 each have planner/implement/review stubs  
- [ ] `docs/init` notes that TODO.md is execution SoT  
- [ ] No P0-X product features implemented in this row  

## Minimal todos
- [ ] Write rules + AGENTS.md  
- [ ] Write Phase 00/01 folder scaffold + prompts  
- [ ] Write docs/TODO.md  
- [ ] Link init → TODO  
