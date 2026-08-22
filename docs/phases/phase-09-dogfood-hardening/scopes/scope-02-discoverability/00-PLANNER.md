# P09-S02-00 — Discoverability (DF-02/DF-04)

## Metadata
- id: P09-S02-00
- todo_ids: [P09-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S02 implement/review prompts so dogfood agents can **list tasks** (id / title / work_state / goal_id) without hardcoded UUIDs, and so **seed import** paths stop footgunning against `-C`. **No product Go in this row.**

## Depends-on (S01)

S01 APPROVE high (DF-01): `lookupEntity` `case "review"`. Discoverability must **not** assume why/context skip reviews. After reviews exist, `trace why` / `trace context` remain valid listing companions — do not regress S01.

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `cmd/trace/root.go` dispatch | No `tasks` / project-level `status` list command |
| `cmd/trace/help.go` | Seed documented as import only; no task list |
| `cmd/trace/seed.go` | `os.ReadFile(path)` — **cwd-relative** (or abs); `-C` only opens store |
| `internal/store` | `ListTasksByGoalID` exists; **no** `ListTasks` / list-all |
| MCP (`internal/mcp`) | Six tools; **no** list-tasks tool — S02 is **CLI-primary** |
| fixtures/x0 README | Documents cwd/abs seed; harness uses abs (must stay green) |
| Dogfood | DF-02 (hardcoded UUIDs); DF-04 (seed vs `-C`) |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| DF-02 command | **`trace tasks`** — not `trace status` (avoids collision with `auth status` / `migrate status`) |
| Stdout | JSON **array** of objects: `id`, `title`, `work_state`, `goal_id` (`null` when unset) |
| Empty | `[]` on success (exit 0) — not an error |
| Optional filter | `--goal <goal-id>` → filter via existing `ListTasksByGoalID`; omit flag → all tasks |
| Store API | Add `(*store.Store).ListTasks()` — `SELECT … FROM tasks ORDER BY created_at ASC, id ASC`; **no** new migration |
| Wire | Thin `cmd/trace` only (G19): open store via `-C`, list, encode JSON; no domain fork |
| MCP | **Out of scope** this scope — CLI for dogfood cold-start; S03 owns install-wire / MCP config |
| DF-04 seed resolve | Relative `<file>` → resolve against **project root** (`resolveRoot(-C)`), then `filepath.Abs`. Absolute paths unchanged |
| DF-04 docs | Update `help.go` seed line + `fixtures/x0/README.md` harness note to match new resolve rule; prefer abs in evals (already) |
| Tests | CLI: seed → `tasks` lists id/title/work_state/goal_id; relative seed with cwd≠project + `-C` succeeds; abs seed still works (p0x/x0) |
| Carry-forward | honesty A/B/C + Gate G; p0x; x0; `CGO_ENABLED=1 go test ./...`; S01 review Why/context regression |
| Forbidden | Daemon/HTTP/embeddings; new `011_*` mig; new MCP tools; rewriting Phase 00–08 history; weakening S01 review lookup |

## Exit
- [x] Thicken `01-discoverability.md` + `02-scope-review.md`
- [x] SCOPE-TODOS + board Notes; next **P09-S02-01**
- [x] Product Go — **not** this row
