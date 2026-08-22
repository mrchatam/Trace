# P08 / S02 / 02 — Scope review (worktrees / project bind)

## Metadata
- id: P08-S02-02
- todo_ids: [P08-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S02 worktree / multi-root bind deliverables vs planner locks. Confirm **path-local `.trace`** + **concurrent fail-closed** without swarm / shared-DB / adapter coupling. Spawn remediations on blocker/high. Fresh subagent — do not share implementer session.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-worktrees.md](01-worktrees.md)
- Implementer board Notes on `P08-S02-01`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (review).

## Focus checklist (locked — P08-S02-00)

- [ ] Root resolve remains Abs-only (`cmd/trace` + MCP); **no** walk-up to parent `.trace` or git common-dir
- [ ] `store.Open` still binds `<absRoot>/.trace/trace.db` (per-root, not shared parent)
- [ ] Exclusive `trace.lock` (or equivalent locked path) acquired on Open, released on Close
- [ ] Second Open on same root fail-closed with clear exported error; CLI exit non-zero / clear stderr
- [ ] Isolation test: two abs roots → two DBs; no cross-visibility of project state
- [ ] Concurrent Open fail-closed test green; reopen after Close works
- [ ] **No** `011_*` for this scope; no daemon/HTTP/swarm; G19 intact
- [ ] **No** coupling to `LanguageAdapter` / analyzer registration
- [ ] `gitcli` still opens store on same abs root (worktree HEAD via `git -C`)
- [ ] Carry-forward green: Gate H (`evals/perf`), honesty A/B/C, Gates E/F/G, ablation, p0x, x0, Gate C `dry_run:false` intact
- [ ] Upcoming S03 stubs still compatible (backup targets path-local `.trace`; may use `011_*` for backup/auth only)

## Verify commands (re-run independently)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/perf/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Role work

1. Diff claims (01 Notes + locks) vs repo evidence.
2. Findings by severity: blocker | high | medium | low | nit.
3. blocker/high: small inline fix **or** spawn `P08-S02-02a` / `02b` with full prompts immediately below this row.
4. Write [REVIEW-NOTES.md](REVIEW-NOTES.md); confidence medium/high with residuals listed.
5. Light-confirm S03 Depends notes still accurate; thicken upcoming only if needed.
6. Board status + Notes; forward-only (no rewrite of `done` history).

## Exit criteria
- [ ] APPROVE (high, or medium with explicit residuals) **or** spawn with evidence
- [ ] REVIEW-NOTES under this scope folder
- [ ] Board status + Notes; next runnable is **P08-S03-00** when APPROVE

## Out of scope
- Implementing backup/auth (S03) or `evals/compat` (S04)
- Designing swarm frameworks or shared-parent `.trace` mode

## Minimal todos
- [ ] Diff claims vs repo; re-run locked tests
- [ ] Write REVIEW-NOTES; APPROVE or spawn; update board
