# P00 / S07 / 01 — trace CLI

## Metadata
- id: P00-S07-01
- todo_ids: [P00-S07-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Expose **P0-X operations** via a **thin `trace` CLI adapter** that calls **library APIs only** (G19 / DR-API): `init`, `index`/`reindex`, entity add + link + task transition, seed import, `why`, and `context`. Preserve existing `help` / `version`. **No** business logic forked in the CLI; **no** MCP/daemon/HTTP.

Enables scripted walkthroughs for S08 (fixture + P0-X harness) and closes the adapter gap for P0-X #3/#4 (why/context at the CLI boundary).

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md) — CLI command list; library≠CLI import; P0-X #3/#4
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G1 no blobs, G6 no dumps, G19 adapters never fork logic, incremental
- [D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-API, DR-SURFACE, DR-AGENT, DR-NAME, DR-TRACEDIR, DR-SEED, DR-INCREMENTAL, DR-P0X
- [B_INITIAL_BOARD.md](../../../../init/B_INITIAL_BOARD.md) — historical T009 (scripted CLI walkthrough)
- Live priors (S01–S06 **done**): `cmd/trace` stub help+version `0.0.0-dev`; `store.Open` → `.trace/trace.db` + mig 001–004 + FTS backfill; `gitcli.Open` / `Refresh`; `analyzers.IndexFile` / `IndexFileAtRev` / `SkipError`; `domain.New` Create*/Link*/TransitionTask; `retrieval.Why`; `compiler.TaskContext` / `ExpandContext`; `go.mod` **go 1.24.0**; no cobra/urfave

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Live substrate (do not re-guess)

| Fact | Value |
|------|-------|
| CLI today | `cmd/trace/main.go` only — `help` / `-h` / `version` → `0.0.0-dev`; unknown cmd → stderr + help + exit 1 |
| Library packages | `internal/{store,vcs,gitcli,analyzers,domain,retrieval,compiler}` — all implemented (S01–S06) |
| Store Open | Creates `.trace/`, applies migrations incl. `004_fts`, **backfills FTS** when empty — CLI must **not** call `RebuildFTS` / Sync* for normal ops |
| Analyzers | File-local incremental; `SkipError` for unsupported/binary; CGO required to **link** this package |
| Domain | Causal creates/links/transitions only via `domain.Service` — **not** raw store Upsert for goals/tasks/… |
| Why / context | `retrieval.Why` / `compiler.TaskContext`+`ExpandContext` — CLI formats output only |
| Forbidden | MCP, daemon, HTTP, TUI, embeddings, full-graph dump command |

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | `go 1.24.0` in `go.mod` (do **not** downgrade) |
| Binary | `trace` from `cmd/trace` (package `main`) |
| Framework | **stdlib only** — argv dispatch + `flag` per command. **Do not** add cobra / urfave / similar |
| Layout | Keep adapter under `cmd/trace/` (multiple `*.go` files OK). Optional tiny helpers in same package. **Do not** put product logic in `internal/*`; **do not** create `internal/cli` unless a test forces it — prefer `cmd/trace` |
| Library vs CLI | `internal/*` **must not** import `cmd/…`. CLI may import `internal/…` |
| Project root | Global `-C <dir>` (also `--project <dir>`). Default: process cwd. Resolve with `filepath.Abs` before `store.Open` |
| Version string | Keep `0.0.0-dev` unless already centralized; print single line on `version` / `--version` / `-version` |
| Store | Every command that needs DB calls `store.Open(projectRoot)` then `defer Close()`. Never reimplement migrations / FTS rebuild |
| VCS | When git repo present: `gitcli.Open(root)` → optional `Refresh` on `init` and before bulk `index`/`reindex`. Pass `retrieval.New(st).WithVCS(repo)` (and compiler `WithRetrieval`) for why/context temporal notes. If not a git repo: proceed without VCS (init/store still OK) |
| CGO | Full `./cmd/trace` binary that imports `analyzers` needs **`CGO_ENABLED=1`**. Document in help/README Build note if missing. Library packages without analyzers remain `CGO_ENABLED=0`-clean |
| Logic rule (G19) | CLI = parse args → call library → print result / exit code. **No** ranking, FTS SQL, provenance defaults beyond flag→struct mapping, parse trees, or work-state graphs in CLI |
| Out | MCP/HTTP/TUI/daemon; `dump` / list-all-entities; embeddings; claim/evidence/review **promotion** engine (stubs via `add`/`seed` OK); Phase-01 review automation |

### Command surface (locked)

```text
trace [-C <root>] <command> [args]

help | -h | --help     Show help (exit 0)
version | --version    Print 0.0.0-dev (exit 0)

init                   Open/create .trace/trace.db (store.Open); optional gitcli Refresh;
                       print DB path on success

index [paths...]       Index given paths (repo-relative or absolute under root).
                       If no paths: walk project for analyzer-supported extensions
                       (.js/.jsx/.mjs/.cjs/.ts/.tsx/.py — match DetectLanguage),
                       skipping .git/ and .trace/; honor .gitignore when walking
                       (best-effort: skip ignored paths). Per file: read bytes →
                       analyzers.IndexFile (or IndexFileAtRev when using VCS HEAD).
                       SkipError → count skip, continue. Other errors → fail command.
                       File-local only — never wipe/rebuild whole symbol graph.

reindex [paths...]     Same implementation as index for P0-X (IndexFile already
                       replaces one file’s edges). Alias allowed as shared handler.
                       Still no full-DB rebuild.

add <kind> …           Create via domain.Create* (see flags below)
link <rel> …           domain.Link* helpers only
transition …           domain.TransitionTask only (prefer this for DONE)

seed import <file>     Import seed JSON v1 (below) via domain APIs

why <type> <id>        retrieval.Why → JSON stdout (WhyResult)
context <task-id>      compiler.TaskContext or ExpandContext → stdout
```

**Kinds for `add`:** `goal` | `task` | `decision` | `assumption` | `discovery` | `plan-change` | `claim` | `evidence`  
Map CLI `plan-change` → domain entity `plan_change`.

**Rels for `link`:** `goal-task` | `decision-task` | `discovery-plan-change` | `claim-evidence`  
Map to `LinkGoalTask` / `LinkDecisionTask` / `LinkDiscoveryPlanChange` / `LinkClaimEvidence`.

### Flags (locked minimum)

| Command | Flags / args |
|---------|----------------|
| Global | `-C` / `--project` |
| `add *` | `--title` (required), `--body`, `--id` (optional UUID), `--source-type`, `--confidence`, `--status` |
| `add task` | also `--goal-id`; **do not** expose `--work-state` for non-PENDING creates in happy path — omit so default PENDING; use `transition` for DONE |
| `link *` | `--from` / `--to` (UUIDs); optional `--source-type` |
| `transition` | `--task` (id), `--to` (work_state), `--actor`, `--reason` (required non-empty), `--allow-done` (bool → `AllowDoneWithoutReview`), optional `--evidence` (repeatable ids) |
| `why` | positional `<type> <id>` — type = domain strings (`goal`, `task`, `decision`, …) |
| `context` | positional `<task-id>`; `--depth` default **1**, max **2** (depth 1 → `TaskContext`; depth 2 → `ExpandContext`); `--format` = `json` (default) \| `markdown` \| `both`; `--include-why` bool |
| `seed import` | positional path to JSON file |

Human-oriented progress (`indexed N, skipped M`) → **stderr**. Machine payloads (`why` / `context` JSON/MD, created entity ids) → **stdout**.

### Seed JSON v1 (locked for S07 + S08)

Single UTF-8 JSON object:

```json
{
  "version": 1,
  "goals": [{"id":"…","title":"…","body":"…"}],
  "tasks": [{"id":"…","title":"…","body":"…","goal_id":"…"}],
  "decisions": [{"id":"…","title":"…","body":"…"}],
  "assumptions": [],
  "discoveries": [],
  "plan_changes": [],
  "claims": [],
  "evidence": [],
  "links": [
    {"rel":"goal_has_task","from":"<goal-id>","to":"<task-id>"},
    {"rel":"decision_affects_task","from":"<decision-id>","to":"<task-id>"},
    {"rel":"discovery_causes_plan_change","from":"<discovery-id>","to":"<plan-change-id>"},
    {"rel":"claim_has_evidence","from":"<claim-id>","to":"<evidence-id>"}
  ],
  "transitions": [
    {"task_id":"…","to":"IN_PROGRESS","actor":"seed","reason":"…","allow_done":false}
  ]
}
```

Rules:

1. `version` must be `1` (reject others).
2. Create entities **before** links/transitions; honor explicit `id` when present.
3. All creates/links/transitions go through **`domain.Service`** — never raw store Upsert for causal rows.
4. Task creates: leave `work_state` empty/PENDING; apply `transitions` via `TransitionTask`.
5. Unknown fields: ignore or reject — pick one, document in code comment; prefer **reject unknown top-level keys** only if cheap, else ignore extras under entities.
6. Idempotency: re-import with same ids may Upsert via Create* (store upsert semantics) — acceptable for P0-X; do not invent a separate merge engine.

### Output + exit codes (locked)

| Code | Meaning |
|-----:|---------|
| 0 | Success |
| 1 | Usage / unknown command / bad flags |
| 2 | Operational failure (store/git/analyzer/domain/retrieval error) |

- `why`: `json.MarshalIndent` of `WhyResult` (already tagged).
- `context` `json`: `Packet.JSON()`; `markdown`: `Packet.Markdown()` (+ ensure `IncludeMarkdown` / cache as library expects); `both`: JSON block then Markdown separated by `\n---\n` (or JSON to stdout and MD note — pick JSON-then-MD separator and stick to it).
- `add` / `link` / `transition` / `seed`: print created/affected ids as simple text or one JSON object on stdout (stable enough for scripts). Prefer one JSON object `{"ok":true,...}` for seed summary.

### Wiring matrix (do not reimplement)

| CLI | Library call |
|-----|----------------|
| `init` | `store.Open` (+ optional `gitcli.Open` + `Refresh`) |
| `index`/`reindex` | `analyzers.IndexFile` / `IndexFileAtRev`; treat `*analyzers.SkipError` as skip |
| `add` / `seed` creates | `domain.CreateGoal` / `CreateTask` / … / `CreateClaim` / `CreateEvidence` |
| `link` / seed links | `LinkGoalTask` / `LinkDecisionTask` / `LinkDiscoveryPlanChange` / `LinkClaimEvidence` |
| `transition` / seed transitions | `TransitionTask` |
| `why` | `retrieval.New(st)[.WithVCS].Why` |
| `context` | `compiler.New(st)[.WithRetrieval].TaskContext` or `ExpandContext` |

### Target tree

```text
cmd/trace/
  main.go           # os.Exit(run(args)); version const
  help.go           # printHelp listing all commands
  root.go           # parse -C / global dispatch
  init.go
  index.go          # shared by reindex
  add.go
  link.go
  transition.go
  seed.go           # seed import + JSON v1 types
  why.go
  context.go
  *_test.go         # dispatch + integration on temp dirs
```

File names may vary; command coverage must not.

### Tests (required evidence)

1. **Dispatch:** `help`/`version` still exit 0; unknown → 1.
2. **Init:** temp dir → `.trace/trace.db` exists after `init`.
3. **Causal round-trip:** `add goal` + `add task --goal-id` + `link decision-task` (or seed) → `why task <id>` JSON contains steps with **reason_code** fields; `context <task>` JSON has `items` + budgets; trust labels present when excerpts included.
4. **Seed:** minimal JSON v1 import succeeds; second why/context works.
5. **Index incremental (CGO):** write two supported files A/B → index both → change only A → `index A` → B’s store symbols unchanged (query via store APIs in test or `go test` helper). Proves no full-rebuild path.
6. **`go test ./...`** passes; `CGO_ENABLED=1 go build -o bin/trace ./cmd/trace` succeeds; smoke `./bin/trace version`.

## Board rights
Implementer: **status + notes only** on this row. No spawning; no rewriting upcoming prompts.

## Out of scope
- Product MCP / daemon / HTTP
- New library ranking/FTS/domain behavior (fix only if CLI cannot call a clear bug — prefer Note for S07-02)
- Fixture corpus under `fixtures/x0` (S08) — CLI must be ready for it
- X0 agent harness

## Exit criteria
- [ ] Commands above implemented as thin adapters (G19) with locked flags/seed v1
- [ ] `why` / `context` stdout match library JSON/MD contracts; reason codes / budgets / untrusted labels not stripped
- [ ] `index`/`reindex` file-local only; SkipError skipped; incremental isolation test passes
- [ ] Causal mutations only via `domain.*` (no raw store Upsert for entities)
- [ ] No MCP/daemon/HTTP; no dump API; no cobra
- [ ] `CGO_ENABLED=1 go build -o bin/trace ./cmd/trace` + `go test ./...` evidence in Notes
- [ ] TODO.md status + Notes updated; SCOPE-TODOS checked

## Minimal todos
- [ ] Expand `cmd/trace` dispatch + help for locked commands
- [ ] Wire `init` / `index`+`reindex` / `add` / `link` / `transition` / `seed import` / `why` / `context`
- [ ] Seed JSON v1 parse + domain apply
- [ ] Tests: dispatch, init, causal why/context, seed, incremental index (CGO)
- [ ] Board status+notes
