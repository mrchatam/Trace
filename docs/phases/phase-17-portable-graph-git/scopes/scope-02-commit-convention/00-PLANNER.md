# P17-S02-00 — commit convention (FINAL)

## Metadata
- id: P17-S02-00
- todo_ids: [P17-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock **FINAL** defaults for **DF-82** (commit path, gitignore stays, export-before-PR) and **DF-85 docs** (SHA/author = evidence; actor ≠ identity). Point at S03 merge behavior. **DF-86** git-hook is **deferred / not this scope**. **No product Go** except help-string tweak if S01 line lacks recommended path / evidence wording.

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [DF-84-FORWARD.md](../../DF-84-FORWARD.md)
- Research § Loop 2 path convention (historical; plan-tree exclude **superseded**)
- Live: `.gitignore`; `AGENTS.md`; `CONTRIBUTING.md`; `README.md`; `cmd/trace/help.go`; MCP `actor` / `as_operator` “not authorization” / flag≠identity (DF-44 class)
- [docs/TODO.md](../../../../TODO.md)

## Depends-on S01 FINAL (2026-08-17 — do not re-lock in S02-00)

Sibling [../scope-01-seed-export/00-PLANNER.md](../scope-01-seed-export/00-PLANNER.md) **FINAL** already fixed:

- CLI **`trace seed export [-o <file>]`** (stdout default) + unchanged `seed import`
- Recommended commit path **`trace/graph.json`**
- Top-level **`exported_at_commit`** (git `rev-parse HEAD` when repo; omit when unknown; import accepts + ignores)
- Export includes **plan tree** keys (`plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state`)

S02 documents **how/when** to commit and attribution (DF-85 docs); does **not** change JSON shapes or export code.

## Session start
Follow agent-loop-protocol. Depends-on: S01 APPROVE (board). **No product Go** except locked help tweak below.

## Live inventory (confirmed 2026-08-17)

| Area | Present? | S02 gap |
|------|----------|---------|
| `.gitignore` | `.trace/` only — **`trace/` committable** | Docs must say unchanged; test read-only |
| `cmd/trace/help.go` | `seed export` line + `exported_at_commit` when git repo | Missing **recommended path** `trace/graph.json` + **evidence not identity** |
| `AGENTS.md` | P17 queued in current focus | Missing **export-before-PR** one-liner |
| `CONTRIBUTING.md` | Generic contribution model | Missing portable-graph clone recipe + merge pointer |
| `README.md` | Build + MCP install | Missing **clone/offline recipe** from `trace/graph.json` |
| S01 tests | `TestSeedExport*` green | No **`TestHelpSeedExportPath`** yet |
| DF-44 help | `flag≠identity`; `Actor string ≠ auth` on `transition` | Complementary — S02 adds git SHA/author evidence docs |
| DF-86 hook | Not implemented | Docs may mention **CONDITIONAL later**; must **not** require hook for P17 |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Path | Recommended commit path **`trace/graph.json`** (authoritative via `-o`; do not gitignore `trace/`) |
| `.gitignore` | **Unchanged** — `.trace/` ignored only; implementer **must not** add `trace/` or `graph.json` |
| Version | Git **commit SHA** is snapshot version; JSON **`exported_at_commit`** is **evidence**, not identity. Entity/plan ids stay UUIDs |
| Attribution | Git **author + SHA** (commit metadata + `exported_at_commit`) are **evidence**. `transition.actor` / review actor / `as_operator` are **not** auth (DF-44 class — do not treat actor strings as identity after clone) |
| Merge | Two-PR conflict on `trace/graph.json` → **human** resolve in git; **no** merge driver. Same-id different body → UUID **last-import-wins** (**S03** owns upsert semantics) |
| DF-86 | **Not this scope.** Do not implement `trace install git-hook`. Docs may say hook is **CONDITIONAL later**; **export-before-PR** + CI docs are the backup. Must **not** wrap `git commit` |
| Forbidden | Changing ignore of `.trace/`; encryption docs as product; rewriting done P16/P17 history; hosted MCP; pointing `trace-mcp` at internet; re-locking S01 JSON/export code |

### Exact files (implementer — docs + help only)

| File | Change |
|------|--------|
| **`AGENTS.md`** | Add one bullet under **Hard boundaries** (or adjacent) — see locked string below |
| **`CONTRIBUTING.md`** | New **`## Portable graph (git)`** section — clone recipe + export-before-PR + merge pointer |
| **`README.md`** | New **`## Portable graph (clone recipe)`** section — offline workflow from committed JSON |
| **`cmd/trace/help.go`** | Tighten **`seed export`** help block — see locked string below |
| **`cmd/trace/cli_test.go`** | Add **`TestHelpSeedExportPath`** |
| **`.gitignore`** | **Read-only verify** — no edit unless accidentally broken |

Do **not** edit `cmd/trace/seed.go`, `internal/domain/seed_export.go`, MCP tools, or import allowlist.

### Locked strings

#### AGENTS.md — export-before-PR (add under Hard boundaries)

```markdown
- **Portable graph:** Before a PR that changes Trace entities, export the semantic graph (including plan tree) to `trace/graph.json` via `trace seed export -o trace/graph.json` (`.trace/` stays local; git SHA is evidence, not identity — see CONTRIBUTING).
```

Keep existing **Current focus** paragraph intact (do not rewrite P16 complete history).

#### CONTRIBUTING.md — new section `## Portable graph (git)`

Must include (wording may tighten; substance locked):

1. **Path:** commit **`trace/graph.json`** at repo root (not `.trace/`).
2. **Before PR:** run `trace seed export -o trace/graph.json` when entity/plan graph changed.
3. **Clone recipe:** `trace init` → `trace seed import trace/graph.json` → `trace index [paths…]` → use `trace why`, `trace context`, `trace plan show` offline.
4. **Evidence:** git **author + commit SHA** and JSON **`exported_at_commit`** are snapshot **evidence**, not entity identity (UUIDs unchanged). **`transition.actor` / review actor / `as_operator` are not authentication** (same as DF-44).
5. **Merge:** parallel PRs may conflict on `trace/graph.json` — resolve in git manually; **no** custom merge driver. Re-import after merge; same UUID different body → **last-import-wins** (**Phase 17 S03**).
6. **Hook (DF-86):** optional **`trace install git-hook`** may come later; **not required** for contributions — export-before-PR remains valid without a hook.

#### README.md — new section `## Portable graph (clone recipe)`

Short fenced block:

```bash
trace init
trace seed import trace/graph.json
trace index
trace plan show
trace why goal <id>
trace context <task-id>
```

One sentence: index rebuilds derived code graph locally; causal/plan data comes from git JSON. Point to CONTRIBUTING for merge/export conventions.

#### Help — `seed export` block (replace current two-line continuation)

```text
  seed export [-o <file>]
                        Export seed JSON v1 (causal entities, links, plan tree,
                        findings/alternatives when present) to stdout or -o.
                        Recommended commit path: trace/graph.json (.trace/ stays local).
                        Sets exported_at_commit (git SHA evidence, not identity)
                        when -C root is a git repo.
```

Keep existing `seed import` and **Handoff** blocks unchanged (DF-28 complementary).

## Named tests (required)

| Test | Package | Intent |
|------|---------|--------|
| **`TestHelpSeedExportPath`** | `cmd/trace` | `trace help` output contains **`trace/graph.json`**; **`exported_at_commit`**; and evidence-not-identity phrasing (`not identity` or `evidence` near export/commit context). Must **not** remove DF-28 handoff or DF-44 actor wording |
| Keepers | `cmd/trace` | `TestHelpHandoffSoT`; `TestAsOperatorFlagIdentityDocs`; S01 `TestSeedExport*` (unchanged) |

TDD: red `TestHelpSeedExportPath` first, then help + docs.

## Owns

| DF | Intent |
|----|--------|
| DF-82 | Path convention + gitignore unchanged + export-before-PR agent docs |
| DF-85 docs | Git author+SHA + `exported_at_commit` = evidence; actor ≠ identity in contributor docs |

## Explicit deferrals

| Item | Owner |
|------|-------|
| Idempotent import / upsert | **S03** (DF-81/83) |
| Git hook auto-export | **DF-86** deferred |
| JSON shapes / export builder | **S01** (done — do not touch) |
| Two-clone VERIFY recipe | **S04** |

## Light Depends — downstream scopes

### S03 idempotent-import
- **Depends-on S02 FINAL:** contributor docs state human git merge on `trace/graph.json`, no merge driver, UUID last-import-wins pointer. S03 implements upsert + may add one paragraph cross-linking CONTRIBUTING; does not re-debate path or gitignore.

## Exit criteria
- [x] 00-PLANNER **FINAL**
- [x] 01/02 runnable thickened prompts
- [x] No product Go except locked help tweak
- [x] Next **P17-S02-01**
