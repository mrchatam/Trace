# Phase 23 — Agent enforcement choke points + harness install rules

**Status:** scaffolded (2026-08-20, P23-00 done) — human-promoted successor after Phase 22 close. Design SoT: [`ENFORCEMENT.md`](ENFORCEMENT.md). Next runnable: **`P23-S01-00`**.

Phase 22 DR-HANDOFF historical **`no successor`** at time-of-close. This phase is a **forward** queue.

## Why this phase exists

Phase 20–22 shipped deliberation policy (`SelectNext`, blocking uncertainties, verification debt) and loop stdout (`trace loop next|apply|status`). That policy only bites when agents **call** Trace. Prompts and harness defaults do not guarantee loop usage before edits, DONE transitions, or seed export.

Dogfood (D44, ab-cms-fullstack) showed **Trace won on planning quality when G1 agents used loop + seed** — but B0/G-prompt-only arms could still implement prematurely. Phase 23 adds **machine-checkable choke points** (library + CLI exit codes) and **harness install snippets** (Cursor rules/hooks, AGENTS.md block) that call those choke points — without a daemon, hosted MCP, or wrapping `git commit`.

## Architecture

```text
Harness (Cursor rules / optional pre-edit hook)
    ↓  calls before edits
trace loop gate --task <id> [--for edit|execute|…]  →  exit 0|1 + trace.loop.gate.v1 JSON
    ↓
internal/deliberation + domain.PrematureImplementation evaluator
    ↓  reuses PolicyInputs / SelectNext semantics
trace loop status  →  violations[] (warn/strict surfaces)
    ↓
Optional hard stops (--enforce):
  trace transition … DONE
  trace seed export --strict
```

| Layer | Role | Lives in |
|-------|------|----------|
| `internal/deliberation` | Reuse `PolicyInputs`, `SelectNext` for gate decisions | **extend** |
| `internal/domain` | `PrematureImplementation` (+ related gate errors) | **extend** |
| `internal/loop` | Gate evaluator, status `violations[]` | **extend** |
| CLI `trace loop gate` | Harness-agnostic exit 0/1 + JSON | **new subcommand** |
| `.trace/config.json` | `enforce`: `off` \| `warn` \| `strict` (default **off**) | **new file contract** |
| `internal/install` | Cursor/claude rules + optional cursor-hook target | **extend** |
| Git hook | Post-commit index/export (P22) | **unchanged** — separate from harness pre-edit hook |

## Phase locks (P23-00)

| Item | Lock |
|------|------|
| P19/P20 loop | **Extend** `internal/loop` + `cmd/trace/loop.go`; do not replace; P19+P20 keeper tests stay green |
| Policy reuse | Gate evaluator builds `PolicyInputs` same as `BuildPolicyInputs`; no second policy table |
| Default enforce | **`off`** — opt-in via `.trace/config.json` or explicit `--enforce` flags |
| `--enforce` | Only on `trace transition … DONE` and `trace seed export --strict` in MVP |
| `--for` values | `orient`, `edit`, `execute`, `done`, `export` — map to policy checks (S02-00 locks table) |
| JSON schema | `trace.loop.gate.v1` stdout on gate (pass or fail) |
| Status schema | Additive `violations[]` on `trace.loop.status.v1` |
| Config | `.trace/config.json` — local, gitignored; documented shape in ENFORCEMENT.md |
| Install | Extend `trace install` — cursor/claude **rules snippets** + optional **`cursor-hook`** target; AGENTS.md enforcement block template |
| Git hook | P22 `trace install git-hook` stays post-commit; **must not** wrap `git commit`; harness pre-edit hook is **separate** |
| MCP | Default **no new MCP tools** — CLI/library inherit; gate is stdout for hooks |
| Forbidden | daemon; hosted MCP; full-rebuild indexer; replacing SelectNext; LLM enforcement |

## Scope order (locked)

| Scope | Focus | Product slice |
|-------|--------|---------------|
| S00 | Architecture lock vs live repo | P23-00 |
| S01 | `domain.PrematureImplementation` + gate evaluator lib | S01 |
| S02 | CLI `trace loop gate` + `trace.loop.gate.v1` | S02 |
| S03 | `--enforce` on DONE transition + `seed export --strict` | S03 |
| S04 | `violations[]` on loop status + `.trace/config.json` | S04 |
| S05 | Harness install rules (cursor/claude/agents + optional cursor-hook) | S05 |
| S06 | VERIFY + DR-HANDOFF | S06 |

## MVP locks

- **Stdout-first**, harness-agnostic (Law 13).
- **Reuse** `internal/deliberation`, `internal/loop`, `BuildPolicyInputs`, existing transition/export seams.
- Gate exit **0** = allowed, **1** = blocked (stderr human hint + JSON stdout).
- `warn` mode: violations surfaced in status JSON; CLI choke points still exit 0 unless `--enforce`.
- `strict` mode: same as warn for status; config may tighten install defaults (document only in S05 — no silent strict default).
- No daemon, hosted MCP, or HTTP on core path.

## Out of scope unless promoted

- Hosted Trace enforcement service
- IDE-specific daemons watching filesystem
- Automatic subprocess agent spawning from gate failure
- Replacing human/orchestrator review (Law: absolute LLM enforcement impossible)

## Completion bar

A fresh agent with Trace installed via `trace install cursor --write` (S05):

1. Sees AGENTS.md / Cursor rules that say **call `trace loop gate --for edit` before product edits** when a task id is active.
2. Running gate before EXECUTE with blocking uncertainty returns exit **1** and JSON with `code=premature_implementation`.
3. `trace loop status` includes `violations[]` when policy would block edit/execute/done.
4. `trace transition … DONE --enforce` fails closed when verification debt open (unless existing escape hatches unchanged).
5. `trace seed export --strict --enforce` fails when graph would export with open violations.
6. Default install leaves `enforce: off` in `.trace/config.json` (or omits file = off).
