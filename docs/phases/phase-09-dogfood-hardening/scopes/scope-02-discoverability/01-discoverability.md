# P09 / S02 / 01 — Discoverability (DF-02/DF-04)

## Metadata
- id: P09-S02-01
- todo_ids: [P09-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Ship **`trace tasks`** so agents can cold-start on work_state without hardcoded UUIDs (DF-02), and fix **seed import** relative-path resolve against `-C` project root (DF-04). Keep carry-forward gates green. Do **not** regress S01 review why/context.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL 2026-08-16
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-02 / DF-04
- Live: `cmd/trace/{root,help,seed}.go`; `internal/store/helpers.go` (`ListTasksByGoalID`); `fixtures/x0/README.md`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute.

## Locked defaults (FINAL — P09-S02-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Command | `trace tasks` (not `status`) |
| JSON fields | `id`, `title`, `work_state`, `goal_id` (JSON `null` if unset) |
| Shape | JSON array on stdout; empty → `[]`; exit 0 on success |
| Filter | Optional `--goal <id>` → `ListTasksByGoalID`; else `ListTasks()` |
| Store | New `ListTasks()` — no migration; reuse task columns |
| Seed resolve | Relative path → join/Abs under project root from `-C`; absolute unchanged |
| Docs | `help.go` + `fixtures/x0/README.md` (and help usage line for seed) |
| Packages | `internal/store` + `cmd/trace` (+ tests). **No** new MCP tools |
| Carry-forward | honesty A/B/C + Gate G; p0x; x0; S01 Why/context-with-review; `./...` |
| Forbidden | Daemon/HTTP/embeddings; `011_*`; MCP list tool; board spawn rights; S01 regress |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/store/helpers.go` (+ test) | `ListTasks()` all tasks ordered like ByGoalID |
| `cmd/trace/tasks.go` (new) | `cmdTasks`: parse optional `--goal`, open store, encode array |
| `cmd/trace/root.go` | Dispatch `case "tasks":` |
| `cmd/trace/help.go` | Document `tasks` + seed resolve-vs-`-C` |
| `cmd/trace/seed.go` | Resolve relative import path against project root before `ReadFile` |
| `cmd/trace/cli_test.go` | Cover `tasks` list + relative seed with `-C` when cwd ≠ project |
| `fixtures/x0/README.md` | Update harness note to new relative-resolve rule (abs still preferred for harnesses) |

## Role work

1. TDD: failing CLI test — after seed, `trace -C <root> tasks` prints JSON with expected id/title/work_state/goal_id.
2. Add `store.ListTasks`; wire thin `cmdTasks`; help text.
3. DF-04: relative seed under project root with process cwd elsewhere + `-C` must succeed; abs path regression still green.
4. Optional `--goal` filter if cheap (locked — implement).
5. Run locked verify suite; board **status + Notes only**.

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-check:

```bash
trace -C <project> tasks
# → JSON array with id, title, work_state, goal_id

# cwd outside project; seed file relative under project:
(cd /tmp && trace -C <project> seed import seed/relative.json)  # must resolve under <project>
```

## Exit criteria
- [ ] `trace tasks` JSON list: id, title, work_state, goal_id; empty `[]` OK
- [ ] Optional `--goal` filters via `ListTasksByGoalID`
- [ ] `ListTasks` in store; **no** new migration
- [ ] Relative seed resolves against `-C` project root; absolute unchanged
- [ ] Help + fixtures/x0 README match behavior
- [ ] Carry-forward green (incl. S01 review Why/context); no MCP/daemon
- [ ] Board Notes ready for **P09-S02-02**

## Out of scope
- S03 install-wire / Cursor mcp.json writer / DF-05
- New MCP `trace_tasks` tool (defer unless review promotes)
- FTS / why / context changes
- Rewriting Gate C packs or Phase 00–08 history

## Todo updates
Implementer: own row status + Notes only. Do not rewrite planner locks or spawn board rows.

## Minimal todos
- [ ] Add `store.ListTasks` + unit coverage
- [ ] Implement `trace tasks` (+ `--goal`) and help dispatch
- [ ] Resolve relative seed paths against project root; update docs
- [ ] CLI tests for list + seed footgun fix
- [ ] Run locked verify suite; mark P09-S02-01 done with Notes
