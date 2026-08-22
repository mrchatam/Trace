# P23-00 — Phase 23 scaffold: enforcement choke points + harness install

## Metadata
- id: P23-00
- todo_ids: [P23-00]
- role: planner
- skills: [planning-and-task-breakdown, documentation-and-adrs, writing-for-agents, grilling]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective

Lock Phase 23 against live repo + [`ENFORCEMENT.md`](ENFORCEMENT.md) + D44 dogfood learnings. Scaffold S01–S06. **No product Go on this row.**

Do not rewrite Phase 22 `done` history. P22 DR-HANDOFF historical `no successor` stays true at close; this phase is a forward human queue.

## References
- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md) — Laws 2, 6, 11–13, 19
- [phase README](README.md)
- [ENFORCEMENT.md](ENFORCEMENT.md)
- [experiments/ab-cms-fullstack/ENFORCEMENT.md](../../../experiments/ab-cms-fullstack/ENFORCEMENT.md)
- Live: `internal/loop/{next,apply,policy}.go`, `internal/deliberation/select.go`, `cmd/trace/{loop,transition,seed,install}.go`, `internal/install/`
- [docs/TODO.md](../../TODO.md), [docs/TODO/phase-23.md](../../TODO/phase-23.md)

## Session start

Human requested Phase 23 MVP: product choke points + harness install rules. Locks below are scaffold defaults for P23-00 — not a license to implement Go here.

## Live inventory (do not fork)

| Surface | Present | Gap vs enforcement MVP |
|---------|---------|---------------------------|
| `trace loop next/apply/status` | **Yes** P19–P20 | No `trace loop gate`; no `violations[]` on status |
| `BuildPolicyInputs` / `SelectNext` | **Yes** P20–P22 | Not exposed as harness exit-code gate |
| `trace transition … DONE` | **Yes** | No `--enforce`; escape hatches exist |
| `trace seed export` | **Yes** P17 | No `--strict` enforcement path |
| `.trace/config.json` | **No** | enforce modes unset |
| `trace install` | **Yes** cursor/claude/git-hook/agents | No loop-gate rules; no cursor-hook; no AGENTS.md enforcement block |
| Domain gate errors | **Partial** (validation errors) | No `PrematureImplementation` typed gate |
| MCP | **15** tools P22 | No gate tool (by design — stdout CLI for hooks) |
| Schema max | **027** P22 | S01–S04 likely **no SQL** or minimal config table — S01-00 decides |

## Locked defaults (phase)

| Item | Value |
|------|-------|
| Goal | Machine-checkable enforcement at edit/DONE/export boundaries + install snippets |
| Keep | P19/P20 loop schemas — **extend** additively |
| Policy | Reuse `PolicyInputs` + `SelectNext` — no second table |
| CLI gate | `trace loop gate --task <id> [--for orient\|edit\|execute\|done\|export]` exit 0/1 |
| JSON | `trace.loop.gate.v1` |
| Enforce flags | `--enforce` on DONE transition + seed export `--strict` only (MVP) |
| Config | `.trace/config.json` `enforce`: off\|warn\|strict; default **off** |
| Status | Additive `violations[]` on `trace.loop.status.v1` |
| Install | Rules + optional `cursor-hook`; git-hook **unchanged** (post-commit) |
| Transport | stdout-first; harness-agnostic |
| Forbidden | daemon; hosted MCP; wrapping git commit; replacing SelectNext |
| MCP | No new tools by default |

## Scope order (locked)

1. S01 premature-impl-gate-lib
2. S02 loop-gate-cli
3. S03 enforce-done-export
4. S04 violations-config
5. S05 harness-install-rules
6. S06 VERIFY + DR-HANDOFF

## Planner work (this row)

1. Confirm ENFORCEMENT.md covers MVP slices S01–S05.
2. Ensure S01–S06 have `00-PLANNER` / `01-*` / `02-review` / `SCOPE-TODOS.md`.
3. Board Phase 23 after Phase 22 in index + orchestrator paste.
4. Update AGENTS.md current focus.
5. Open DR-HANDOFF.md.

## Exit criteria

- [ ] README + ENFORCEMENT.md cover MVP
- [ ] S01–S06 stubs exist with locked requirements
- [ ] Board rows P23-00 … P23-S06-02
- [ ] AGENTS.md next runnable `P23-S01-00` (post P23-00 execution)
- [ ] No product Go

## Next

Orchestrator: **P23-S01-00** after this row is `done`.
