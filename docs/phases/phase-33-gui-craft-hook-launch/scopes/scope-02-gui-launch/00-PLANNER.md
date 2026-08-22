# P33-S02-00 — Scope planner (`trace gui` + PATH)

## Metadata
- id: P33-S02-00
- todo_ids: [P33-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Finalize S02 implement/review for **`trace gui`** (cwd/`-C`): start local HTTP GUI, print URL, **open default browser** (best-effort), reuse serve + P32-PORT. Lock PATH install story from S00 (`go install` / symlink / make) — distinct from `trace install` agents/MCP. **No product code in this planner row.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- `../scope-00-research/RESEARCH.md`
- Live: `cmd/trace/serve.go`, `cmd/trace/root.go`, `cmd/trace/help.go`, `cmd/trace/serve_test.go`, `internal/httpapi/`

## Session start

Follow agent-loop-protocol Session start. Prefer S00 RESEARCH leans; do not reopen Theme C.

## Locked defaults

| Item | Value |
|------|-------|
| CLI | Subcommand **`trace gui`** primary; `-gui` **secondary only** (prefer subcommand per RESEARCH / Theme C) |
| Behavior | Reuse serve listen path; loopback default `127.0.0.1:7432`; open browser after listen; print URL always |
| Open fail | Non-zero only if serve fails — browser open failure → stderr tip, still success if listening |
| Opt-out | Support UA-style **`--no-open`** (or equivalent) so CI/headless can skip browser |
| Port | Reuse P32-PORT friendly in-use messaging — **reject** UA auto-increment / silent port hop |
| PATH | **#1** document `go install github.com/mrchatam/Trace/cmd/trace@…`; **#2** contributor make/symlink from `bin/`; package later. Do **not** overload `trace install` agents/MCP |
| Docs primary flip | Minimal help/install note OK; **S05** owns full quickstart primary-story rewrite |
| Law 19 | CLI calls library/httpapi — no business-logic fork |
| Landing URL | Open browser to GUI root Explore **`/`** (Theme B / S01 UX-IA) — **not** Nav `/overview` |

## Must answer (handoff to 01) — resolved 2026-08-21

1. **Flags:** Inherit serve (`--addr`, `--allow-remote`, `--token`, `--token-file`, `--root`, `--static-dir`, `--cors-origin`) + **`--no-open`**; global `-C`/`--project`. No `-gui` this scope.
2. **Open:** `cmd/trace` helper — darwin `open` / windows `cmd /c start` / else `xdg-open`; injectable for tests; open after successful listen (optional OnListening / shared helper — no sleep race).
3. **PATH:** Help Build note + optional short quickstart tip = `go install github.com/mrchatam/Trace/cmd/trace@…`; contributor build/symlink secondary; no root Makefile required; never `trace install` for PATH.
4. **Tests:** help; remote refuse; addr-in-use; `--no-open` / mocked URL `http://{addr}/`; no auto-port.
5. **Land:** Explore **`/`** only — never `/overview`.

## Planner gate

- [x] `01-implement.md` thick enough for fresh subagent
- [x] `02-review.md` checklist includes PATH ≠ agents-install
- [x] `SCOPE-TODOS.md` accurate

## Exit criteria

- [x] Implementer locked; next **P33-S02-01**

## Todo updates

Status + notes on **P33-S02-00** only.

## Next

`P33-S02-01`
