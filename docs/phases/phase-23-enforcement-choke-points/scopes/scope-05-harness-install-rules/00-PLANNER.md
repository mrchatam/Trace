# P23-S05-00 — Harness install rules planner

## Metadata
- id: P23-S05-00
- todo_ids: [P23-S05-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock harness install extensions: Cursor/Claude rules, AGENTS.md enforcement block, optional `cursor-hook` target, and **further install-time rules** that reference loop gate. **No product Go this row.**

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- Live: `internal/install/{registry,cursor,claude,githook,agents}.go`, `cmd/trace/install.go`
- S02 gate CLI — **done**: `trace loop gate --task <id> [--for edit|…]` exit 0/1/2
- S04 status + config — **done**: `violations[]` on status; `.trace/config.json` enforce modes
- Cursor hooks contract: project `.cursor/hooks.json` + `.cursor/hooks/*` (preToolUse JSON stdin/stdout)

## Live inventory (P23-00 + S04 verified)

**Install registry today** (`internal/install/registry.go`): `cursor`, `claude`, `git-hook` only.

| Target | Today | S05 change |
|--------|-------|------------|
| `cursor` | STABLE — upserts `$HOME/.cursor/mcp.json` trace MCP entry | **Extend** — also write project rules + AGENTS.md block; pass `ProjectRoot` from CLI |
| `claude` | CONDITIONAL — `.claude/trace-mcp.json` when marker | **Extend** — also write claude rules + AGENTS.md block |
| `git-hook` | CONDITIONAL — `# begin-trace` post-commit index/export | **Unchanged** — P22 behavior preserved |
| `cursor-hook` | **Missing** | **New** CONDITIONAL target — preToolUse hook calling gate CLI |
| `agents` | Bundled harness profiles in SQLite | **Unchanged** — may cross-link in rule text only |

**Marker reuse:** `githook.go` owns `# begin-trace` / `# end-trace` for git hooks. AGENTS.md / CLAUDE.md use **different** markers (`# begin-trace-enforcement`) so git-hook and harness docs never collide.

## Live touch points (S05-01)

| File | Change |
|------|--------|
| `internal/install/types.go` | Add `TargetCursorHook = "cursor-hook"` |
| `internal/install/registry.go` | Register `cursorHookTarget{}` |
| `internal/install/markers.go` | **New** — generic `upsertMarkedFragment` / `removeMarkedFragment` with marker pair params (extract from githook or duplicate constants) |
| `internal/install/enforcement.go` | **New** — shared rules body, AGENTS.md/CLAUDE.md merge, env-var docs |
| `internal/install/cursor.go` | Extend `Install`/`Uninstall`/`Detect` for project rules file + AGENTS.md |
| `internal/install/claude.go` | Extend for claude rules surface + AGENTS.md |
| `internal/install/cursorhook.go` | **New** — hooks.json merge + shell script |
| `internal/install/githook.go` | **No behavior change** — may refactor to call shared markers helper only |
| `cmd/trace/install.go` | Pass `ProjectRoot` to cursor; add `cursor-hook` subcommand |
| `cmd/trace/help.go` | Document new target + env vars |
| `internal/install/*_test.go` | Named install tests (see 01) |
| `CONTRIBUTING.md` | Harness vs product enforcement paragraph |

## Locked defaults (S05-01 must not re-debate)

| Item | Value |
|------|-------|
| Targets | Extend **cursor**, **claude**; add **`cursor-hook`** install target |
| git-hook | **Do not change** P22 post-commit/pre-push fragment or `# begin-trace` markers |
| Print / `--write` | Same pattern as existing install targets — print-only default; `--write` upserts |
| Network | None |
| Daemon | None |
| Config on install | **Do not auto-write** `.trace/config.json` — document omission = `off`; print template in help/rules only |
| Bundled agents | `trace install agents` SQLite upsert unchanged |

### File paths (FINAL)

| Surface | Path | Tier / detect |
|---------|------|---------------|
| Cursor MCP | `$HOME/.cursor/mcp.json` (or `--mcp-json`) | unchanged — user-level |
| Cursor rules | `<projectRoot>/.cursor/rules/trace-enforcement.mdc` | created on `--write`; whole-file owned by Trace |
| Cursor hook config | `<projectRoot>/.cursor/hooks.json` | merged `preToolUse` entry only |
| Cursor hook script | `<projectRoot>/.cursor/hooks/trace-loop-gate.sh` | mode `0755` |
| Claude MCP | `<projectRoot>/.claude/trace-mcp.json` | unchanged |
| Claude rules (primary) | marker block in `<projectRoot>/CLAUDE.md` | when file exists |
| Claude rules (fallback) | `<projectRoot>/.claude/trace-enforcement-rules.md` | when no `CLAUDE.md` |
| AGENTS.md block | `<projectRoot>/AGENTS.md` | marker-delimited section |
| git-hook | `.git/hooks/post-commit`, `pre-push` | **unchanged** |

### Marker strings (FINAL)

| Use | Begin | End |
|-----|-------|-----|
| AGENTS.md / CLAUDE.md enforcement | `# begin-trace-enforcement` | `# end-trace-enforcement` |
| git-hook (P22 — do not change) | `# begin-trace` | `# end-trace` |
| Cursor rules `.mdc` | whole-file replace | N/A — uninstall deletes file |

Idempotent merge: re-run `--write` replaces content between enforcement markers (AGENTS/CLAUDE) or overwrites owned files (`.mdc`, hook script).

### Environment contract (FINAL)

| Variable | Set by | Semantics |
|----------|--------|-----------|
| `TRACE_TASK_ID` | Orchestrator / user | Active seed task UUID. When **empty**, rules say skip gate; hook **allows** (fail-open). When **set**, hook runs gate. |
| `TRACE_PROJECT_ROOT` | Optional | Project root for `trace -C`. Hook defaults to Cursor project cwd when unset. |

Document in rules body, AGENTS.md block, and help text.

### Rules snippet content (FINAL — all three surfaces share this intent)

1. **Pre-edit gate** — when `TRACE_TASK_ID` is set, before product code edits:
   ```bash
   trace loop gate --task "$TRACE_TASK_ID" --for edit
   ```
   Exit **0** → proceed. Exit **1** → read stdout JSON `recommended_phase`; do not edit. Exit **2** → surface stderr; do not parse stdout.

2. **Status violations** — before marking DONE or large scope changes:
   ```bash
   trace loop status --task "$TRACE_TASK_ID"
   ```
   Inspect `violations[]`; non-empty → follow `recommended_phase` before DONE.

3. **Opt-in DONE enforce** — team policy only:
   ```bash
   trace transition --task "$TRACE_TASK_ID" --to DONE --enforce
   ```

4. **Opt-in export strict** — CI / pre-PR only:
   ```bash
   trace seed export --strict --enforce -o trace/graph.json
   ```

5. **Config pointer** — `.trace/config.json` `{ "enforce": "off"|"warn"|"strict" }`; default off when missing; does **not** auto-enable transition/export flags.

6. **Design SoT link** — `docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md`

Keep rules **concise** — no policy fork prose in AGENTS.md; link ENFORCEMENT.md for semantics.

### Cursor rules file shape (FINAL)

Path: `.cursor/rules/trace-enforcement.mdc`

```markdown
---
description: Trace loop gate — pre-edit enforcement for active tasks
alwaysApply: true
---

## Trace enforcement (installed by trace install cursor --write)

When `TRACE_TASK_ID` is set for the active seed task:

1. Before product code edits, run `trace loop gate --task "$TRACE_TASK_ID" --for edit`.
2. Before DONE, run `trace loop status --task "$TRACE_TASK_ID"` and resolve any `violations[]`.
3. Use `--enforce` on transition DONE and `seed export --strict --enforce` only when the team opts in.

See docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md. Local config: `.trace/config.json` (`enforce`, default off).
```

### AGENTS.md block (FINAL)

```markdown
# begin-trace-enforcement
## Trace enforcement (harness)

When a Trace seed task is active, set `TRACE_TASK_ID` to its UUID.

- **Before edits:** `trace loop gate --task "$TRACE_TASK_ID" --for edit` (exit 0 = proceed).
- **Before DONE:** `trace loop status --task "$TRACE_TASK_ID"` — resolve non-empty `violations[]`.
- **Opt-in strict:** `--enforce` on `trace transition … --to DONE`; `trace seed export --strict --enforce` for CI.
- **Config:** `.trace/config.json` → `{ "enforce": "off"|"warn"|"strict" }` (default off).

Product design: docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md
# end-trace-enforcement
```

Merge: if `AGENTS.md` missing, create with block + trailing newline. If present, upsert between markers (same algorithm as git-hook fragment merge).

### Claude rules (FINAL)

- When `<projectRoot>/CLAUDE.md` exists: upsert **same marker block** as AGENTS.md (shared body function).
- When absent: write `<projectRoot>/.claude/trace-enforcement-rules.md` with equivalent markdown (no markers — whole file).

### cursor-hook target (FINAL)

| Field | Value |
|-------|-------|
| ID | `cursor-hook` |
| Tier | `CONDITIONAL` |
| Detect | `.cursor/` directory under `ProjectRoot` **or** existing `.cursor/hooks.json` with trace hook command |
| Refuse | `--write` without `.cursor/` → error cites CONDITIONAL (same pattern as claude) |

**hooks.json merge** — upsert under `hooks.preToolUse[]`:

```json
{
  "command": ".cursor/hooks/trace-loop-gate.sh",
  "matcher": "Write|StrReplace|ApplyPatch|EditNotebook",
  "failClosed": false
}
```

Preserve sibling hooks; replace entry when `command` path matches trace hook script.

**Shell script** (`.cursor/hooks/trace-loop-gate.sh`):

```bash
#!/usr/bin/env bash
set -euo pipefail
input="$(cat)"
task_id="${TRACE_TASK_ID:-}"
if [[ -z "$task_id" ]]; then
  echo '{"permission":"allow"}'
  exit 0
fi
if ! command -v trace >/dev/null 2>&1; then
  echo '{"permission":"allow"}'
  exit 0
fi
root="${TRACE_PROJECT_ROOT:-$PWD}"
if trace -C "$root" loop gate --task "$task_id" --for edit >/dev/null 2>&1; then
  echo '{"permission":"allow"}'
  exit 0
fi
msg="Trace loop gate blocked edit for task ${task_id}. Run: trace loop gate --task ${task_id} --for edit"
printf '{"permission":"deny","user_message":%q,"agent_message":"Follow recommended_phase from gate JSON before editing product code."}\n' "$msg"
exit 2
```

Hook semantics:

- `TRACE_TASK_ID` empty → **allow** (no active task).
- `trace` missing → **allow** (fail-open — rules still document manual gate).
- Gate exit **0** → allow.
- Gate exit **1** → deny (exit **2** from hook script per Cursor hook contract).
- Gate exit **2** / internal → allow (fail-open unless future `failClosed: true` — locked **false** for MVP).

**Separate from git-hook:** cursor-hook is pre-edit; git-hook remains post-commit index/export only.

### cursor target CLI change (FINAL)

`cmdInstallCursor` must pass `ProjectRoot: abs` (from `resolveRoot`) so rules + AGENTS.md merge run in the project Trace was invoked from (`trace -C … install cursor --write`).

Print-only behavior:

- **Stdout:** MCP JSON snippet (unchanged).
- **Stderr:** rules snippet header + body + existing `CursorReloadTip`.

### Uninstall / detect parity (FINAL)

| Target | Uninstall removes | Detect signal |
|--------|-------------------|---------------|
| `cursor` | MCP trace entry (existing) + `.cursor/rules/trace-enforcement.mdc` + AGENTS.md enforcement block | `.cursor` dir or mcp.json (existing) + optional reason if rules file present |
| `claude` | `trace-mcp.json` (existing) + CLAUDE.md block or fallback rules file + AGENTS.md block | `.claude/` or `CLAUDE.md` (existing) |
| `cursor-hook` | hook script + trace entry from hooks.json + AGENTS.md block | `.cursor/` or trace hook in hooks.json |
| `git-hook` | **unchanged** | **unchanged** |

Second uninstall idempotent. Sibling hooks / MCP servers preserved.

### Help text (FINAL — S05-01 add to `help.go`)

```
  install cursor [--write] [--bin path] [--mcp-json path]
        Upsert trace MCP + project enforcement rules (.cursor/rules/trace-enforcement.mdc)
        and AGENTS.md block. Print-only: MCP JSON on stdout, rules hint on stderr.

  install claude [--write] [--bin path]
        Upsert .claude/trace-mcp.json + claude enforcement rules + AGENTS.md block
        when .claude/ or CLAUDE.md marker present.

  install cursor-hook [--write]
        Install preToolUse hook calling trace loop gate (requires .cursor/ under project).
        Sets .cursor/hooks/trace-loop-gate.sh + hooks.json entry. Respects TRACE_TASK_ID.

  Environment: TRACE_TASK_ID (active task UUID), TRACE_PROJECT_ROOT (optional -C root).
```

## Planner work

1. [x] Lock marker strings and file paths (.cursor/rules, hooks, AGENTS.md, CLAUDE.md).
2. [x] Lock env contract (`TRACE_TASK_ID`, `TRACE_PROJECT_ROOT`).
3. [x] Lock cursor-hook JSON + shell contract (preToolUse, fail-open defaults).
4. [x] Lock git-hook unchanged boundary.
5. [x] Thicken `01-harness-install-rules.md` with named tests + implementation sketches.
6. [x] Thicken `02-scope-review.md` with evidence table + keeper cmds.
7. [x] Update `SCOPE-TODOS.md`.

## Exit criteria

- [x] Install targets + markers locked
- [x] Further harness rules scope explicit (gate, status, DONE enforce, export strict)
- [x] 01/02/SCOPE-TODOS runnable alone
- [x] No product Go

## Next

**P23-S05-01** after this row is `done`.
