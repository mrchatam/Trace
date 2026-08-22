# P00 / S08 / 01 — Fixture + P0-X harness

## Metadata
- id: P00-S08-01
- todo_ids: [P00-S08-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Ship a **synthetic Apache-2.0 fixture** (`fixtures/x0`) with **human-curated** ground-truth seed JSON v1, plus a **deterministic P0-X harness** (`evals/p0x`) that proves **all 7** DR-P0X criteria without an LLM. Prefer scripted **`trace` CLI** walkthroughs against a temp copy of the fixture (library APIs allowed for store assertions only).

Closes the empirical gap for Phase 00 VERIFY (S09). **No** MCP/daemon/HTTP; **no** real OSS corpus; **no** agent Gate C (`evals/x0`).

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md) — P0-X 7/7; paths `fixtures/x0` + `evals/p0x`
- [I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — Experiment P0-X; synthetic only (DR-BENCH); human GT (DR-SEED)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G1 no blobs; G6 no dumps; G12/G21 incremental; G19 adapters
- [D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-P0X, DR-BENCH, DR-SEED, DR-INCREMENTAL, DR-PARSE, DR-AGENT
- [F_QUESTION_LEDGER.md](../../../../init/F_QUESTION_LEDGER.md) — Q-UNDERSTAND-N (≥5), Q-FIXTURE-LANG (TS+Python)
- Live priors (S01–S07 **done**): `cmd/trace` thin CLI; seed JSON **v1**; `internal/{store,vcs,gitcli,analyzers,domain,retrieval,compiler}`; `go.mod` **go 1.24.0**; CGO binary for analyzers

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Live substrate (do not re-guess)

| Fact | Value |
|------|-------|
| CLI | `trace [-C root] <cmd>` — `init` / `index`/`reindex` / `add` / `link` / `transition` / `seed import` / `why` / `context` / `help` / `version` |
| Seed path (S07 residual) | `seed import <file>` opens `<file>` via **process cwd** (or absolute). **`-C` does not rewrite the seed path.** Harness **must** pass an **absolute** seed path (or `chdir` so a relative path resolves). |
| Seed JSON v1 | Reject unknown top-level keys; `version: 1`; creates → links → `transitions` via `TransitionTask`; link `rel`: `goal_has_task` / `decision_affects_task` / `discovery_causes_plan_change` / `claim_has_evidence` |
| Index | File-local `IndexFile` / `IndexFileAtRev`; `SkipError` → skip; no full-graph wipe |
| Why / context | `why` → JSON `WhyResult` with `reason_code` on every step; `context` → Packet JSON/MD, budgets **4096** tokens / **32** items, `untrusted_data` labeling |
| Store | `.trace/trace.db` under project root; migrations 001–004 + FTS backfill on Open |
| CGO | Harness that builds/runs full `trace` needs **`CGO_ENABLED=1`** |
| Missing today | No `fixtures/` or `evals/` trees yet — create them here |

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | `go 1.24.0` in `go.mod` (do **not** downgrade) |
| Fixture root | `fixtures/x0/` (committed synthetic mini-repo; **never** commit `.trace/`) |
| License | Fixture ships **Apache-2.0** (`LICENSE` file under `fixtures/x0/`; SPDX or full text OK — match repo LICENSE spirit) |
| Corpus | **Synthetic only** (DR-BENCH). No vendored OSS projects |
| Seeding | **Human-curated** GT only (DR-SEED). No LLM-generated seed as SoT |
| Seed file | `fixtures/x0/seed/gt.json` — seed JSON **v1** with **stable fixed UUIDs** (harness asserts these ids) |
| GT narrative | `fixtures/x0/README.md` documents planted architecture + causal story (which goal/task/decision/discovery IDs mean what) |
| Languages (Q-FIXTURE-LANG **locked**) | ≥1 **TypeScript or JavaScript** file **and** ≥1 **Python** file under the fixture, both analyzer-supported extensions, with at least one import/symbol each so P0-X #2 is non-vacuous |
| Fixture shape | Small readable “app”: e.g. `src/greeter.ts` (or `.js`) + `src/math_util.py` (+ optional tiny `package.json` / empty `__init__` — keep minimal). Enough structure for import/symbol assertions; not a monorepo |
| Git | Do **not** commit a `.git` under `fixtures/x0`. Harness **may** `git init` + one commit on a **temp copy** so VCS-backed index/why paths exercise; non-git still acceptable if all 7 criteria pass |
| Workdir rule | Harness copies `fixtures/x0` → `t.TempDir()` (or equivalent) and runs `trace -C <copy> …` so the committed tree stays clean |
| Harness path | `evals/p0x/` — Go tests (`package p0x_test` or `p0x`) that are CI-able |
| Harness entry | Primary gate: `CGO_ENABLED=1 go test ./evals/p0x/... -count=1` |
| Binary | Build once in harness/`TestMain` or helper: `go build -o <tmp>/trace ./cmd/trace` (or use `exec.Command("go", "run", …)` only if flaky-free). Prefer built binary |
| Metrics artifact | On success write `metrics-p0x.json` into the test temp dir **and** assert a schema with `criteria: { "1": true, … "7": true }` (plus optional timings). Do **not** require committing a live metrics file; optional golden under `evals/p0x/testdata/` only if useful |
| Queries (Q-UNDERSTAND-N **locked**) | ≥**5** distinct deterministic understanding checks (no LLM). Minimum set below |
| Incremental (#7) | Change **one** indexed source file in the temp copy; re-`index` **only that path**; assert the other file’s symbols/imports in SQLite are **byte-identical / unchanged**; changed file’s rows update. Proves no full-fixture rebuild |
| Transitions | Seed uses `transitions` + `TransitionTask` (via CLI seed). Do **not** CreateTask with DONE shortcuts |
| Out | MCP/daemon/HTTP; `evals/x0` agent harness; real OSS corpus; embeddings; dump-all APIs; rewriting S07 CLI semantics |

### CLI walkthrough (locked order)

```text
# From Trace repo root; WORK = abs path to temp copy of fixtures/x0
# SEED = abs path to seed/gt.json  (ABS — not rewritten by -C)

trace -C "$WORK" init
trace -C "$WORK" seed import "$SEED"     # ABS or cwd-valid relative
trace -C "$WORK" index                   # walk supported sources
trace -C "$WORK" why task <task-uuid>
trace -C "$WORK" context <task-uuid> --format json
# then mutate one file + index that path only for criterion #7
```

Progress on **stderr**; machine JSON on **stdout**. Exit **0/1/2** per S07.

### Seed JSON v1 (must match S07)

Same schema as S07 `01-cli.md`. Fixture seed **must** include enough graph for criteria 1, 3, 5:

- At least one **goal**, **task** (with `goal_id`), **decision**, **discovery**, and a **plan_change** linked by `discovery_causes_plan_change`
- Links: `goal_has_task`, `decision_affects_task`, `discovery_causes_plan_change` (claim/evidence optional for P0-X)
- At least one `transitions` entry (e.g. task → `IN_PROGRESS`) with non-empty `reason`
- Stable UUIDs listed in README for harness hardcoding

### P0-X criteria → harness assertions (locked mapping)

| # | Criterion | Assert (minimum) |
|---|-----------|------------------|
| 1 | Goal/Task/Decision/Discovery round-trip + provenance | After seed: store/domain Get* (or CLI `why`) shows entities with **ACTIVE** provenance; task `work_state` reflects transition |
| 2 | Files + symbols/imports | After `index`: both TS/JS and Py paths present as files; each has ≥1 symbol **or** ≥1 import row |
| 3 | `trace why` causal chain + reason codes | `why task <id>` JSON: every step has non-empty `reason_code`; chain includes goal and/or decision neighbor (`goal_has_task` / `decision_affects_task`) |
| 4 | `trace context` bounded | `context <task>` JSON: `items` length ≤ **32**; packet respects token budget **4096**; when body/excerpts present, trust / `untrusted_data` labeling appears (JSON field and/or MD) |
| 5 | Human-seeded graph matches GT | Seeded UUIDs from `gt.json` exist; expected links recoverable via `why` steps or store `entity_links` / `goal_id`; README IDs match |
| 6 | ≥5 deterministic understanding queries | See **Query set** below — all must pass |
| 7 | Incremental file update | One-file reindex leaves sibling file’s symbol/import set unchanged; changed file updates |

### Query set for #6 (locked minimum ≥5)

Implement as five named subtests (names may vary; coverage must not):

1. **why-task** — `why task <planted-task>` returns seed + goal/decision neighbors with reason codes  
2. **why-decision** — `why decision <planted-decision>` returns a non-empty chain with reason codes  
3. **decision-constraint** — context or why shows the planted decision affecting the task (reason or entity id present)  
4. **import-or-symbol-neighbor** — after index, exact path/symbol (via store APIs or why/context excerpt) proves the TS/JS **and** Py files are in the graph  
5. **context-boundedness** — context packet respects item/token budgets (and does not dump the whole DB)

Additional queries beyond 5 are welcome; fewer than 5 fails the harness.

### Target tree

```text
fixtures/x0/
  LICENSE                 # Apache-2.0
  README.md               # architecture + GT UUID map + causal story
  seed/
    gt.json               # seed JSON v1, stable UUIDs
  src/
    greeter.ts            # (or .js) — ≥1 symbol and/or import
    math_util.py          # ≥1 symbol and/or import
  .gitignore              # optional: .trace/

evals/p0x/
  p0x_test.go             # (name flexible) — copies fixture, builds trace, runs 7/7
  harness.go              # optional helpers: abs seed path, metrics writer
  doc.go                  # package comment: deterministic P0-X; no LLM
```

File names under `src/` may vary; **both languages** and **seed path** must not.

### Tests / evidence (required)

1. `CGO_ENABLED=1 go test ./evals/p0x/... -count=1` passes and exercises all 7 criteria.  
2. `CGO_ENABLED=1 go test ./... -count=1` still passes (no regressions).  
3. Notes include: fixture paths, seed UUID examples, harness command, metrics schema confirmation, seed-path abs-path reminder.  
4. No committed `.trace/` under `fixtures/x0`; `.gitignore` at repo or fixture level covers it.

## Board rights
Implementer: **status + notes only** on this row. No spawning; no rewriting upcoming prompts.

## Out of scope
- Product MCP / daemon / HTTP
- `evals/x0` agent Gate C harness (Phase 01+)
- Real OSS benchmark corpus
- Changing CLI seed-path semantics (document/work around; do not “fix” by rewriting under `-C` unless S09/S07 spawn says so)
- New retrieval/ranking/domain features beyond what S07 already exposes

## Exit criteria
- [ ] `fixtures/x0` synthetic Apache-2.0 mini-repo with TS/JS **and** Python sources + README GT map + `seed/gt.json` v1
- [ ] `evals/p0x` deterministic harness proves **7/7** via CLI walkthrough (abs seed path) + incremental one-file assertion
- [ ] ≥5 understanding queries pass without LLM; metrics JSON schema written in temp on success
- [ ] No MCP/daemon/HTTP; no full-rebuild indexer path; committed tree has no `.trace/`
- [ ] `CGO_ENABLED=1 go test ./evals/p0x/...` and `./...` evidence in Notes
- [ ] TODO.md status + Notes updated; SCOPE-TODOS checked

## Minimal todos
- [ ] Create `fixtures/x0` (LICENSE, README, src TS/JS+Py, `seed/gt.json` with stable UUIDs)
- [ ] Create `evals/p0x` harness: temp copy → `init` → `seed import <ABS>` → `index` → why/context/GT/incremental assertions
- [ ] Map and assert all 7 P0-X criteria; ≥5 named understanding queries; write metrics-p0x.json in temp
- [ ] Ensure `.gitignore` ignores fixture `.trace/`; run CGO tests
- [ ] Board status+notes
