# P18-S02-00 — clone PENDING honesty (FINAL)

## Metadata
- id: P18-S02-00
- todo_ids: [P18-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock **FINAL** docs/help/comments/tests for **DF-88**: keep P17 export exclude; clone tasks import **PENDING** is expected. Thicken `01`/`02`/SCOPE-TODOS. **No seed-format change. No product Go on this planner row.** Stop if [../../00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) is not FINAL (it is). **Do not reverse DF-88.**

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Laws 1, 6, 9, 11
- [phase README](../../README.md)
- [DF-88-DECISION.md](../../DF-88-DECISION.md) — **SoT** (option A: keep exclude)
- [../../../phase-17-portable-graph-git/DF-84-FORWARD.md](../../../phase-17-portable-graph-git/DF-84-FORWARD.md) — historical exclude (do not rewrite)
- Live: `CONTRIBUTING.md` Portable graph; `README.md` clone recipe; `cmd/trace/help.go` seed export; `internal/domain/seed_export.go` `SeedTask` (no `work_state`); `TestSeedExportOmitsDeniedSurfaces`; `TestHelpSeedExportPath`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Depends-on: **P18-S01-02 APPROVE** (board order) even though S02 does not need S01 code. **Do not reverse export exclude.** Unattended: P18-00 locked DF-88 **wontfix + document**.

**S01 blast (light):** DF-87 is context/FTS only. Slash titles do **not** require seed-format work. S01 keepers already include `TestSeedExportOmitsDeniedSurfaces` as untouched.

## Live inventory (2026-08-17)

| Area | Present? | Gap vs DF-88 |
|------|----------|--------------|
| P17 exclude | **Yes** | `SeedTask` has `id/title/body/goal_id` only; `BuildSeedDocument` never fills `reviews` / `transitions` / task `work_state`. Fail bar: `TestSeedExportOmitsDeniedSurfaces` (no `"transitions"` key, no task `work_state`, no review/token/capability/index leak) |
| DF-84-FORWARD exclude | **Yes** | index, tokens, lock, tool decisions, capabilities, events, **reviews**, **transitions/`work_state`**. Tasks still without replaying DONE |
| Help `seed export` | Path + `exported_at_commit` evidence | **Missing** omit + clone PENDING. `PENDING` substring **absent** from `help.go` today (named test will go red first). `work_state` already appears on `transition` / `tasks` lines — test must require **omit + reviews**, not bare `work_state` |
| `TestHelpSeedExportPath` | **Yes** | DF-82/85 keeper (`trace/graph.json`, `exported_at_commit`, evidence-not-identity). **Do not overload** |
| CONTRIBUTING Portable graph | Bullets 1–6 (path/PR/clone/evidence/merge/hook) | **Missing** clone-PENDING honesty. Do **not** delete or renumber 1–6 |
| README clone recipe | init → import → index → why/context/plan | **Missing** one sentence: import lands PENDING |
| AGENTS portable-graph bullet | export-before-PR + SHA evidence | Optional one clause; do **not** rewrite Current focus from 01 |
| Include flags / new seed keys | **No** | Keep absent |

### Repro (scoring confusion, not a product bug)

D40 clone: `seed import` → `transitions: 0`; clone tasks **PENDING** while live G1/G2 were DONE/SKIPPED. P17 already excluded reviews/`work_state`/transitions. DF-88 documents that as **expected**.

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Product | **Keep omit** reviews / `work_state` / `transitions` (and P17 denied surfaces: index, tokens, lock, capabilities, tool decisions, events). No `--include-reviews`. No `--include-work-state`. No new seed keys |
| Docs | CONTRIBUTING Portable graph: after `seed import`, tasks are **PENDING** unless the clone operator transitions them; process history (reviews/DONE/SKIPPED) is **local** to the exporter’s `.trace/`, not git identity |
| Help | `seed export` block **keeps** path + `exported_at_commit` lines (DF-82/85) **and** states omit + clone tasks import PENDING |
| Fail bar / keeper | `TestSeedExportOmitsDeniedSurfaces` still PASS (must keep omitting). **Do not** weaken assertions |
| Named | **`TestHelpCloneTasksImportPending`** in `cmd/trace` — **new**. Do **not** overload `TestHelpSeedExportPath` |
| Path keeper | `TestHelpSeedExportPath` stays DF-82/85 (`trace/graph.json` + `exported_at_commit` + evidence-not-identity) |
| CGO | `cmd/trace` tests: **`CGO_ENABLED=1`**. CGO0 `cmd/trace` remains carry-forward non-fail (tree-sitter) |
| S05 | **Leave** P18-S05-00/01/02. S02 does **not** rebuild `bin/trace` / `bin/trace-mcp`. Do **not** retarget MCP CGO build-note lines (S05 owns CGO0 MCP). **Depends:** S05 still after S04 VERIFY |
| Forbidden | Include flags; rewriting P17 prompts; MCP seed tool; claiming clone DONE is a product bug; reversing export exclude; changing `SeedTask` JSON tags to add `work_state` |

### Export exclude (unchanged — fail bar)

Do **not** reverse. SoT remains P17 [DF-84-FORWARD.md](../../../phase-17-portable-graph-git/DF-84-FORWARD.md) + `TestSeedExportOmitsDeniedSurfaces`:

| Surface | Expected |
|---------|----------|
| `transitions` key | **Absent** from default export JSON |
| Task `work_state` | **Absent** from task objects |
| `reviews` / residuals | **Absent** |
| Index / token / capabilities / events | **Absent** (P17 denied) |

`why` / `plan show` work from links + plan tree **without** reviews/`work_state`.

## Locked strings (FINAL)

### CONTRIBUTING.md — append as bullet **7** (do not delete or renumber 1–6)

```markdown
7. **Clone honesty (DF-88):** Default `seed export` **omits** reviews, transitions, and task `work_state`. After `seed import`, tasks are **PENDING** until the clone operator transitions them. Live DONE/SKIPPED in the exporter’s `.trace/` is local process, not portable identity. `why` / `plan show` work from links + plan tree without reviews/`work_state`.
```

### README.md — one sentence after the clone-recipe fence (keep the bash block)

```markdown
After `seed import`, tasks are **PENDING** (default export omits reviews, transitions, and task `work_state`; live DONE/SKIPPED stays on the exporter’s `.trace/`).
```

Keep the existing index/CONTRIBUTING sentence.

### Help — `seed export` block (replace the current continuation; keep import/handoff/build-note)

```text
  seed export [-o <file>]
                        Export seed JSON v1 (causal entities, links, plan tree,
                        findings/alternatives when present) to stdout or -o.
                        Recommended commit path: trace/graph.json (.trace/ stays local).
                        Sets exported_at_commit (git SHA evidence, not identity)
                        when -C root is a git repo.
                        Default export omits reviews, transitions, and task work_state.
                        After seed import, clone tasks are PENDING until the clone
                        operator transitions them.
```

Must still contain `trace/graph.json`, `exported_at_commit`, and `not identity` so `TestHelpSeedExportPath` stays green.

### AGENTS.md — optional clause on the existing Portable graph hard-boundary bullet (append; do not replace)

```text
 Default export omits reviews/work_state; clone import tasks are PENDING — see CONTRIBUTING.
```

**01 must not rewrite Current focus.** This planner row updates Current focus to **P18-S02-01**.

### Comment — `internal/domain/seed_export.go` (behavior unchanged)

Above `type SeedTask`:

```go
// SeedTask is portable identity (id, title, body, goal_id). Default export
// omits work_state so clone import lands PENDING (DF-88; keep TestSeedExportOmitsDeniedSurfaces).
```

Above `BuildSeedDocument` (extend the existing one-liner, do not change the function body):

```go
// BuildSeedDocument assembles seed JSON v1 from the store (causal entities, links, plan tree, impact).
// Default export omits reviews, transitions, and task work_state (DF-88). Clone PENDING is expected.
```

## Named tests (FINAL)

| Test | File | Assert |
|------|------|--------|
| **`TestHelpCloneTasksImportPending`** | `cmd/trace/cli_test.go` (adjacent to `TestHelpSeedExportPath`) | `trace help` stdout, `strings.ToLower`: contains **`pending`**; contains **`import`** (clone-landing sense); contains **`omits reviews`**; contains **`transitions`**; contains **`work_state`**. Pattern: same `captureStdout` + `run([]string{"help"})` as `TestHelpSeedExportPath`. Do **not** assert bare `work_state` alone (already on `transition`/`tasks` lines) |
| Keepers | existing | `TestSeedExportOmitsDeniedSurfaces` (must still omit); `TestHelpSeedExportPath` (DF-82/85 path/evidence) |

TDD: red `TestHelpCloneTasksImportPending` first (`PENDING` absent from help today), then help + docs + comments.

## Files likely touched (implementer — P18-S02-01)

- `CONTRIBUTING.md` — bullet 7 only
- `README.md` — one sentence under Portable graph clone recipe
- `cmd/trace/help.go` — `seed export` continuation only
- `cmd/trace/cli_test.go` — named help test
- `internal/domain/seed_export.go` — comments only (no JSON/tag/builder change)
- Optionally `AGENTS.md` portable-graph bullet clause

**Do not touch:** `cmd/trace/seed.go` export path; `BuildSeedDocument` body; import allowlist; MCP tools; analyzers; FTS; P17 prompt files; S05 binary rebuild; MCP CGO build-note lines.

## Locked verify

```text
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpCloneTasksImportPending|TestHelpSeedExportPath|TestSeedExportOmitsDeniedSurfaces'
```

CGO0 `cmd/trace` is carry-forward non-fail (tree-sitter). Do not use it as this scope’s bar.

## Blast / later scopes (upcoming only)

- **S03:** No seed-format coupling. DF-89 golden only. Exclude stays.
- **S04 VERIFY:** import **`TestHelpCloneTasksImportPending`** + keepers `TestSeedExportOmitsDeniedSurfaces`, `TestHelpSeedExportPath` from S02 REVIEW-NOTES after land (names locked here; live name wins if 01 must rename — it must not).
- **S05:** Unchanged board rows after VERIFY. S02 docs/help are not a rebuild. Stale `bin/trace-mcp` remains S05.

## Non-goals
- Product Go on **this** planner row (01 owns help strings + named test + comments)
- CONDITIONAL `--include-reviews` / `--include-work-state`
- Replaying DONE onto clones; adding `CANCELLED`
- Hosted MCP; DF-86 hook; rewriting Phase 17 `done` history

## Planner work (this row)
1. [x] Live re-read CONTRIBUTING + `help.go` seed export + omit test + `TestHelpSeedExportPath` + `SeedTask`
2. [x] Lock exact help/CONTRIBUTING/README/comment substrings + named test **`TestHelpCloneTasksImportPending`** as **FINAL**
3. [x] Thicken 01/02/SCOPE-TODOS; light S03/S04 Depends; S05 rows left in place
4. [x] Mark this prompt **FINAL**; board Notes; next **P18-S02-01**
5. [ ] Product Go — **not** this row

## Exit criteria
- [x] This prompt **FINAL** with locked test name + substrings
- [x] 01/02 thickened enough to run alone
- [x] Board Notes; next **P18-S02-01**
- [x] No export-include flags; no product schema change

## Minimal todos
- [x] FINAL docs/help locks
- [x] Thicken 01/02
- [x] Board sync

## Next
Orchestrator: **P18-S02-01**. Do **not** start P18-S02-02 until S02-01 is `done`.
