# Phase 45 — MCP surface honesty

**Phase planner.** Row `P45-00`. **Done** 2026-10-04 (executed inline; planner row recorded post-hoc).

## Metadata
- id: P45-00
- todo_ids: [P45-00]
- role: planner
- skills: [planning-and-task-breakdown, api-and-interface-design, incremental-implementation]
- mcps: [trace_context, trace_why]
- verification: automated

## Mission

Turn the audit round-5 MCP findings (#149–#162) into six ordered scopes with fail-closed policies locked **before** any product code, and confirm the board row order. **No product code in this row.**

## Gate

Human promotion satisfied 2026-10-04 ("add Phase 45 board and implement"). Phase 44 closed 2026-08-23 with `no successor`; the board sat idle until this intake.

## Scope sequence

| Scope | Theme | Rows | Deliverable |
|-------|-------|------|-------------|
| S01 | FTS honesty | P45-S01-00 → 02 | `matcher` input; `graph_sync_honesty` on all walk tools |
| S02 | Write-path validation | P45-S02-00 → 02 | Fail-closed `trace_add`/`trace_link`/`trace_tasks` |
| S03 | Error quality | P45-S03-00 → 02 | No driver text; `trace.loop.gate.v1` for `loop.next` |
| S04 | Envelope + schema | P45-S04-00 → 02 | Uniform list envelopes; generated descriptions/enums |
| S05 | Store-free reads | P45-S05-00 → 02 | `trace_version` answers from the process alone |
| S06 | Guards | P45-S06-00 → 02 | Key-shape sweep + protocol conformance suite |

Ordering rationale: S04 lands **after** S02/S03 so the uniform envelope and the error wording are written once, over the final set of tools, rather than retrofitted twice. S06 is last so its guards sweep the final responses.

## Locked decisions (human, this intake)

| ID | Question | **Decision** | Rejected alternative |
|----|----------|--------------|----------------------|
| D1 | How to expose FTS operators | **Opt-in `matcher` param**, default `and` (byte-identical legacy) | Change the default to `ops` — would alter every existing caller's results silently |
| D2 | `trace_add` on an existing id | **Fail closed + explicit `replace=true`** | Silent overwrite (current, #150); or hard error always — breaks idempotent re-seeding |
| D3 | Board | **Add Phase 45 and implement** | Forward-residual queue — 14 same-class findings belong on the board |

**D1 detail.** `SearchFTS` keeps its exact current behaviour as a thin wrapper over the new `SearchFTSOpts` default. The 64-result limit cap is preserved in both modes. Operators are recognized **uppercase-only** (`AND|OR|NOT|NEAR`, `NEAR(n)`, `"phrase"`, trailing `*`); lowercase `and`/`or`/`not` stay ordinary terms, matching FTS5 and avoiding a silent change for callers who search for the word "and".

**D2 detail.** `replace=true` overwrites the entity row only. It does **not** cascade to links, changes, or evidence — that is a different, destructive operation and is out of scope. Responses gain `created: bool` so an agent can tell which happened.

## Live-repo research (verified at HEAD `35ed912`, 2026-10-04)

| File:line | Finding | Issue |
|-----------|---------|-------|
| `internal/store/fts.go:74` | `sanitizeFTSQuery` strips operators; no way to ask for them back | #149 |
| `internal/mcp/tools_search.go:28` | **Opens the store twice** for one call | (perf, fixed in passing) |
| `internal/mcp/tools_context.go:25` | The only walk tool that reports `graph_sync_honesty` | #153 |
| `internal/mcp/instructions.go:25` | Server-level instructions advertise the walk tools | #153 |
| `internal/mcp/tools_write.go:69` | `toolAdd` overwrites on id collision; returns `ok:true` | #150, #151 |
| `internal/mcp/tools_parity.go:314` | `toolVersion` calls `openStore` + the capability gate | #152 |
| `internal/mcp/tools_context.go` L34/L41/L50 | `work_state` filter unvalidated | #154 |
| `internal/domain/link.go:33,34` | `LinkBlocks`/`LinkAPIContract` accept `from == to` | #162 |
| `internal/store/entities.go:121` | Getters return `sql.ErrNoRows` verbatim | #155 |
| `internal/mcp/tools_loop.go` | `action=gate` returns an envelope; `action=next` returns prose | #156 |
| `internal/mcp/tools_changes.go` | `trace_changes show` marshals the store struct directly (PascalCase) | #161 |
| `internal/store/open.go:228` | `nowRFC3339` → second granularity ⇒ UUID-order tiebreak | #160 |
| `internal/mcp/tools_agents.go:60` | `recommend` drops the empty-catalog hint | #157 |
| `internal/loop/next.go:169` | `loop.next` preconditions fail with a domain error | #156 |
| `internal/domain/capability.go:56` | Capability kinds include `AGENT`; descriptions omit it | #159 |

### SDK facts verified in source (not from memory)

These three findings shaped the implementation and are the reason D1/S04 were designed the way they were:

1. `github.com/google/jsonschema-go/jsonschema/infer.go:~337` assigns a whole `jsonschema:"..."` struct tag to `Description`; `infer.go:399` `disallowedPrefixRegexp = ^[^ \t\n]*=` **rejects `WORD=` tags**. ⇒ **Enums cannot come from struct tags.** An explicit `*jsonschema.Schema` per tool is required.
2. `mcp/server.go:~550` honors a caller-provided `InputSchema` over the inferred one. ⇒ `schemaFor[In](constraints)` is the mechanism.
3. Schema-validation failures arrive as `CallToolResult.IsError == true` with a Go `err == nil`. ⇒ the conformance suite must assert `IsError`, not a returned error.

## Planner gate (P45-00)

- [x] Re-read live repo files in the Research anchors table
- [x] Close D1–D3 from the human intake
- [x] Thicken S01–S06 `00-PLANNER` + `01-implement` + `02-review`
- [x] Board `docs/TODO/phase-45.md` — scope rows filled, notes evidence-backed
- [x] `SCOPE-TODOS.md` per scope folder
- [x] No product code in this row

## Exit criteria

- [x] Scope stubs runnable; board points at **P45-S01-00**
- [x] D1–D3 closed, no open decisions left to a later scope
- [x] Successor after P45: **no successor** (idle)

## Next

`P45-S01-00`