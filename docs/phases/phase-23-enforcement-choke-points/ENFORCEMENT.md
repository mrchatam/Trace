# Enforcement — product design (tech-lead ruling)

**Status:** locked for Phase 23 MVP scaffold (2026-08-19). Implementers follow scope prompts; this doc is the product SoT for *what* and *why*.

## Problem

Trace deliberation policy (P20–P22) is correct but **optional** at the harness boundary. Agents can edit files, mark DONE, or export seed JSON without ever calling `trace loop next`. Prompts alone do not close that gap (see `experiments/ab-cms-fullstack/ENFORCEMENT.md`, D44).

## Two layers (do not conflate)

| Layer | What it enforces | Mechanism |
|-------|------------------|-----------|
| **Product** | Deliberation order, blocking uncertainty, verification debt, premature implementation | Library gate + CLI exit codes + optional `--enforce` on transition/export |
| **Harness** | Remind or block *before* tool edits when a task is active | `trace install` writes Cursor/Claude rules, optional pre-edit hook contract, AGENTS.md block |

Product enforcement is **honest**: it reads SQLite + git evidence, not agent claims. Harness enforcement is **best-effort**: hooks and rules call product CLI; orchestrators still verify board order.

## Choke points (MVP)

### 1. Pre-edit / pre-phase gate — `trace loop gate`

```bash
trace loop gate --task <uuid> [--for orient|edit|execute|done|export]
```

- **Stdout:** JSON `trace.loop.gate.v1` always (pass or fail).
- **Exit:** `0` allowed, `1` blocked.
- **`--for` semantics (locked intent):**

| `--for` | Blocks when (examples) |
|---------|-------------------------|
| `orient` | Task missing / no seed context |
| `edit` | Blocking uncertainty, open regression, plan missing, or phase not EXECUTE-ready |
| `execute` | Same as edit + explicit execute-pending policy from SelectNext |
| `done` | Verification debt, open regressions, failed gates for DONE promotion |
| `export` | Export would omit required honesty or carry open violations (strict subset) |

Reuse **`BuildPolicyInputs`** + **`SelectNext`** — gate is a thin adapter, not a forked policy engine.

### 2. Domain error — `PrematureImplementation`

Typed domain error (S01) with stable `code` string for JSON:

- Human message for stderr
- Machine fields: `task_id`, `for`, `recommended_phase`, `reason_code`, `violations[]`

Other gate failures may share the same envelope (`trace.loop.gate.v1`) with different `code` values — S01-00 locks the enum.

### 3. Transition DONE — `--enforce`

Existing `trace transition --task … --to DONE` gains optional **`--enforce`**:

- Without flag: current behavior (escape hatches preserved).
- With flag: call gate evaluator for `--for done`; non-zero exit, no transition on block.

Config `enforce: strict` does **not** auto-enable without explicit flag on transition (fail-safe default).

### 4. Seed export — `--strict` + `--enforce`

`trace seed export --strict` (existing or new flag — S03 locks) validates export honesty.

With **`--enforce`**: exit non-zero when violations present; no write on failure.

### 5. Loop status — `violations[]`

Additive field on `trace.loop.status.v1`:

```json
"violations": [
  {
    "code": "premature_implementation",
    "for": "edit",
    "message": "…",
    "recommended_phase": "INVESTIGATE",
    "reason_code": "blocking_uncertainty"
  }
]
```

Populated from same evaluator as gate. Empty array when clean.

### 6. Project config — `.trace/config.json`

Local file (gitignored `.trace/`):

```json
{
  "enforce": "off"
}
```

| Value | Behavior |
|-------|----------|
| `off` | Gate/status compute violations; CLI exits 0 unless `--enforce` on supported commands |
| `warn` | Status + gate JSON include violations; stderr hints; exit 0 unless `--enforce` |
| `strict` | Document as "recommended for CI/harness"; still requires `--enforce` on transition/export unless S04-00 promotes auto-block (default: **no** auto-block) |

**Default when file missing:** `off`.

## Harness install (S05)

Extend **`trace install`** (cursor, claude, new optional **`cursor-hook`**) to:

1. **Cursor / Claude rules snippet** — before product code edits on an active Trace task, run:
   `trace loop gate --task "$TRACE_TASK_ID" --for edit` (or document env var contract).
2. **AGENTS.md enforcement block** — marker-delimited section (`# begin-trace-enforcement` / `# end-trace-enforcement`) with loop gate + config pointer; idempotent merge like git-hook markers.
3. **Optional `cursor-hook` target** — pre-tool-use hook **contract** (JSON in/out documented); calls `trace loop gate`; **separate** from P22 post-commit git-hook.
   - **Option A failClosed (P28):** when `.trace/config.json` has `enforce=strict` and `TRACE_TASK_ID` is absent, the hook returns deny + non-zero exit. `off`/`warn`/missing still allow (default-off). Cursor `hooks.json` keeps `failClosed: false` (policy is script-level).
   - **Multitask / worker inheritance (FM-04):** install rules (`ParentOrchestratorRule`) require the parent to set `TRACE_TASK_ID` before edits and before delegating, put task UUID + workspace in every worker prompt, and own graph writes (gap pass / discoveries / decisions). Cursor Multitask does not guarantee env inheritance; Trace cannot product-detect parent orchestrators (Option B deferred). Option A applies per process to that process’s `TRACE_TASK_ID`.
4. **Further harness rules** (user request) — install also adds rules that reference loop status, transition enforce flag, and seed export strict — without bloating AGENTS.md body (link to ENFORCEMENT.md / phase README).

Install must remain **print / `--write`**, no daemon, no network.

## Non-goals

- No hosted MCP or HTTP enforcement API
- No wrapping `git commit`
- No filesystem watcher daemon
- No guarantee LLMs obey rules — reviewers/orchestrators still check board order

## Evidence bar (S06)

- Gate exit 1 on blocking uncertainty before `--for edit`
- Status `violations[]` matches gate for same task state
- `--enforce` on DONE fails with verification debt
- Default install: enforce off; rules present after `trace install cursor --write`
- P19/P20 loop keeper tests green
