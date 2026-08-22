# P10-S02-00 — MCP parity + install freshness (FINAL)

## Metadata
- id: P10-S02-00
- todo_ids: [P10-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S02 implement/review prompts for **DF-21, DF-22, DF-32**. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19 adapters never fork
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1 thin `trace_tasks` + capability mirror
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-21/22/32
- S01 inherit: [../scope-01-retrieval-why-fidelity/00-PLANNER.md](../scope-01-retrieval-why-fidelity/00-PLANNER.md) + REVIEW-NOTES (do not re-litigate)
- Live: `internal/mcp/server.go` (six tools); `cmd/trace/{tasks,capability,install}.go`; `internal/domain.BuiltinMCPCapabilitySpecs`; `internal/store.Capability` (no json tags)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no user grill required (A1 + live inventory).

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `internal/mcp/server.go` | Tools: `trace_why`/`trace_context`/`trace_add`/`trace_link`/`trace_transition`/`trace_review` only — **no** tasks/capability/version |
| `cmd/trace/tasks.go` | Snake_case JSON array `[id,title,work_state,goal_id]` + optional `--goal` |
| `cmd/trace/capability.go` | Subcmds declare\|list\|require\|unrequire\|missing; **list/missing** encode `[]store.Capability` → **PascalCase** (DF-32) |
| `internal/store/capability.go` | `Capability` fields have **no** `json` tags |
| `internal/compiler/packet.go` | `CapabilityRef` already snake_case — packet OK |
| `cmd/trace/install.go` | Print/`--write` upsert `mcpServers.trace` + `.bak.<UTC>`; `--bin`/`--mcp-json`; **no** reload/version reminder |
| `internal/mcp` `serverVersion` | `"0.0.0-dev"` — not exposed as a tool |
| `domain.BuiltinMCPCapabilitySpecs` | **Six** `mcp:trace_*` slugs — must grow when tools added |
| S01 green | `plan-change`→`plan_change` alias; Exact/Why `capability`; IncludeWhy fail-closed — **inherit**, do not re-implement |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-21 | MCP cold-start parity for tasks + capability | Thin **`trace_tasks`** + **`trace_capability`** (G19); **not** plan/impact/index MCP |
| DF-22 | Stale Cursor MCP after install/rebuild | Docs + install stderr reload tip; thin **`trace_version`** so agents can verify live process |
| DF-32 | Capability JSON case vs `tasks` | CLI list/missing (+ MCP) emit **snake_case** capability rows |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Packages | **`internal/mcp`** (+ thin `cmd/trace-mcp` help text); **`cmd/trace`** for capability JSON encode + install docs/stderr; **`internal/domain`** `BuiltinMCPCapabilitySpecs` only (no domain logic fork) |
| Migration | **None** |
| G19 | MCP/CLI call `domain`/`store`/`retrieval`/`compiler` — **no** business logic in adapters |
| DF-21 `trace_tasks` | Read-only MCP tool. Params: optional `project`, optional `goal_id` (mirrors `--goal`). Output: **same shape as CLI** — JSON array of `{id,title,work_state,goal_id}` (empty `[]`). Use `store.ListTasks` / `ListTasksByGoalID` (or domain wrap if already present) — no new list API required. |
| DF-21 `trace_capability` | **One** tool with `action` ∈ `declare\|list\|require\|unrequire\|missing` (same pattern as `trace_review`). Params mirror CLI flags (`kind`/`slug`/`title`/`status`/`body`/`id`; `task`/`capability` for require/unrequire; `task` for missing; optional `project`). Mutating actions: declare/require/unrequire (`ReadOnlyHint=false`); list/missing read-only. |
| DF-21 non-goals | **No** MCP for `plan` / `impact` / `index` / residual / migrate / backup this phase |
| DF-21 Builtin specs | Update `BuiltinMCPCapabilitySpecs` to include **`mcp:trace_tasks`**, **`mcp:trace_capability`**, **`mcp:trace_version`** (still **no auto-seed**) |
| DF-21 help | Refresh `cmd/trace-mcp` `-h` tool list; MCP tool descriptions mirror CLI |
| DF-22 docs | README **Install / Cursor MCP** + `help.go` (if needed): after build/`install --write`, **rebuild `trace-mcp`**, prefer **`--bin` absolute path**, then **reload/restart Cursor MCP** (or reload window) so the stdio process is not stale |
| DF-22 install stderr | On successful `--write`, print a short stderr tip: rebuild + reload MCP / prefer abs `--bin` (DF-22). Print mode unchanged except optional one-line tip is OK if non-noisy |
| DF-22 `trace_version` | **In-scope required** thin read-only tool. JSON `{"ok":true,"name":"trace","version":"<same as mcp serverVersion / CLI version string>"}` (today `0.0.0-dev`). No process kill; no daemon; Cursor owns restart |
| DF-22 out | Auto-restart Cursor; kill PIDs; HTTP health daemon; rewriting Phase 09 install merge semantics |
| DF-32 shape | Capability **list** and **missing** arrays: each row snake_case **`id`,`kind`,`slug`,`title`,`status`** (body/timestamps **omit** from list/missing unless already on declare — keep declare/require/unrequire maps as today). Envelope may keep `ok`/`count`/`task`/`missing`/`capabilities` as today |
| DF-32 approach | Prefer explicit DTO/map encode in `cmd/trace` (mirror `taskListRow`); MCP must emit the **same keys**. Adding `json` tags on `store.Capability` is **allowed** if tests assert snake_case on CLI stdout — do not leave PascalCase public |
| S01 inherit | Do **not** re-open DF-19/23/25/27/29; `trace_why`/`trace_context` keep S01 behavior |
| Tool count | After S02: **nine** tools (prior six + `trace_tasks` + `trace_capability` + `trace_version`) |
| Carry-forward | honesty A/B/C + Gate G; Gate E/F; ablation; Gate H; compat; p0x; x0; S01 why/Exact tests; Gate C `dry_run:false` intact |
| Forbidden | New mig; daemon/HTTP/embeddings; plan/impact/index MCP; rewriting Phase 00–09 / S01 `done` history; Mode-B Gate C pack rewrite; S03/S04 product work |

## Effects on later scopes
- **S03:** serial after S02 review; no MCP coupling.
- **S04:** capability transition gates consume catalog APIs — MCP require/missing must stay G19-thin so S04 can gate without forking.
- **S05 VERIFY:** MCP checklist = **nine** tools + BuiltinMCP specs + DF-32 snake_case spot-check + DF-22 docs/stderr/`trace_version`.

## Exit
- [x] Thicken `01-mcp-parity-install.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light upcoming Depends notes (S03/S05 stubs only if needed)
- [x] Board Notes; next **P10-S02-01**
- [x] Product Go — **not** this row
