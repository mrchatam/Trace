# P01 / S03 / 01 — Agent X0 harness (CLI)

## Metadata
- id: P01-S03-01
- todo_ids: [P01-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Add **`evals/x0`** harness that dry-runs Experiment X0 conditions **B0** vs **G1** on synthetic `fixtures/x0`, emits **schema-valid** metrics JSON for **both** conditions, and proves G1 actually shells out to `trace` `why`/`context`. Full agent quality scoring / Gate C = **Phase 02**. Keep `evals/p0x` and `evals/honesty` untouched and green.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — Experiment X0
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 1 validation gate
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-AGENT, DR-BENCH
- Sibling [00-PLANNER.md](00-PLANNER.md) locks
- Live patterns: `evals/p0x` temp-copy + abs seed; `fixtures/x0` UUID map; CLI `why`/`context`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Conditions | **B0** = agent + ordinary repo tools (read files / git / search) — **no** `trace why` / `trace context`. **G1** = agent + `trace` CLI `context`/`why` (may still read repo) |
| MCP | **Not** required (DR-AGENT). Do not gate X0 on MCP |
| Corpus | Synthetic **`fixtures/x0`** (+ existing `seed/gt.json` v1); **no** new OSS this phase |
| Seed path | Always **absolute** path to `seed/gt.json` (cwd/abs; `-C` does not rewrite) — same as p0x |
| Stable IDs | Goal `11111111-…111`; Task `22222222-…222`; Decision `33333333-…333`; Discovery `44444444-…444`; PlanChange `55555555-…555` |
| Package | **`evals/x0`** (new) — sibling of `evals/p0x` and `evals/honesty`; **do not** merge honesty or P0-X into X0 |
| Schema | Committed **`evals/x0/schema.json`** (schema_version **1**); tests validate emitted metrics against it |
| Metrics files | Temp artifacts only: `metrics-b0.json` + `metrics-g1.json` (do **not** commit live metrics; no `.trace/` under fixtures/evals) |
| Dry-run bar | Emit schema-valid metrics for **both** B0 and G1. Stub / recorded agent OK. **Not** Gate C scoring; **not** “G1 beats B0” claim |
| Agent N | Schema may allow multiple `runs`; dry-run requires **≥1** run per condition. Real N≥3 agent repeats = Phase 02 |
| Task family (dry-run) | Primary: **`understanding`**. Do **not** fold H5 honesty into X0 metrics this phase (`evals/honesty` owns that) |
| Fixture prep (both conditions) | Temp-copy fixture → build `trace` (`CGO_ENABLED=1`) → `init` → `seed import <abs>` → `index` (reuse p0x patterns; do not import `evals/p0x` as a library) |
| G1 invocation (locked) | After prep, harness **must** invoke (exit 0, non-empty stdout): `trace -C <work> why task 22222222-2222-2222-2222-222222222222` and `trace -C <work> context 22222222-2222-2222-2222-222222222222 --format json`. Record tool use / latency in metrics |
| B0 invocation (locked) | After prep, harness **must not** call `why`/`context`. Stub agent may read fixture files under `<work>` only. Metrics still emitted with `condition: "B0"` |
| DONE policy | Prefer **not** seeding/transitioning to DONE. Existing seed stops at `IN_PROGRESS`. If DONE ever needed: Review PASS or explicit `allow_done` — never EvidenceIDs alone |
| Keep green | `evals/p0x` unchanged; `evals/honesty` unchanged; `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/honesty/... ./evals/x0/...` + `./...` |
| G19 | Harness may shell out to built `trace` binary; packages must **not** import `cmd/trace` |
| CGO | `CGO_ENABLED=1` for `./evals/x0/...` (analyzers/CLI), same as p0x |
| Out | Gate C go/no-go; real multi-model agent loops as exit criteria; MCP-required path; embeddings; rewriting `evals/p0x`; OSS corpus |

### Metrics schema (locked — implement `schema.json` to match)

JSON Schema draft that accepts objects with at least:

| Field | Rule |
|-------|------|
| `schema_version` | integer `1` |
| `experiment` | string `"X0"` |
| `condition` | `"B0"` \| `"G1"` |
| `fixture` | string (e.g. `"fixtures/x0"`) |
| `dry_run` | boolean — **true** for Phase 01 dry-run emits |
| `trace_version` | string (from `trace version`, or `"stub"` for B0-only fields if unused) |
| `runs` | array, **minItems 1**; each run: `run_id` (string), `task_family` (`"understanding"` required in dry-run), `ok` (boolean), optional `efficiency` `{ "latency_ms"?: number, "tokens"?: number }`, optional `quality` placeholders (objects/numbers OK; **unchecked for Gate C**), optional `notes` |
| `tools_used` | array of strings — **G1** must include entries that clearly indicate `why` and `context` (e.g. `"trace why"`, `"trace context"`). **B0** must **not** list those Trace CLI tools |

Implementer may add optional top-level fields (`model`, `seed`, `git_sha`) but must not remove required ones.

### Scenario (locked — implement exactly)

One primary test (suggested name: `TestX0DryRunMetricsB0AndG1`):

```text
# Shared
moduleRoot → copy fixtures/x0 → workDir (skip .trace/.git)
build ./cmd/trace with CGO_ENABLED=1 (TestMain or helper, like p0x)
for each condition independently (separate work dirs OR sequential clean copies):
  trace -C work init
  trace -C work seed import <ABS seed/gt.json>
  trace -C work index

# B0
  do NOT run why/context
  write metrics-b0.json (dry_run=true, condition=B0, ≥1 understanding run, ok=true stub)
  validate against schema.json

# G1
  run why task <TaskID> → exit 0
  run context <TaskID> --format json → exit 0
  write metrics-g1.json (dry_run=true, condition=G1, tools_used includes why+context, ≥1 understanding run)
  validate against schema.json

# Assert
  both files exist under temp; schema validation passes; B0 tools_used excludes why/context; G1 includes them
```

### Target tree

```text
evals/x0/
  doc.go           # package purpose + how to run
  schema.json      # committed metrics schema v1
  x0_test.go       # TestX0DryRunMetricsB0AndG1 (+ helpers: copy, build, validate)
```

Optional: `harness.go` helpers. Do **not** commit `.trace/` or live metrics JSON.

### How to run (implementer Notes must cite)

```bash
CGO_ENABLED=1 go test ./evals/x0/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/honesty/... ./evals/x0/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] `evals/x0` exists with `schema.json` + dry-run test emitting **B0 + G1** metrics
- [ ] G1 dry-run invokes live `trace why` + `trace context` on fixture task UUID
- [ ] B0 dry-run does **not** invoke why/context
- [ ] Metrics validate against schema; `dry_run: true`
- [ ] No committed `.trace/` / live metrics under fixtures or evals
- [ ] `evals/p0x` + `evals/honesty` still PASS; `CGO_ENABLED=1 go test ./...` green
- [ ] TODO.md status + Notes updated (test name + metrics paths)

## Minimal todos
- [ ] Scaffold `evals/x0` + `schema.json` v1
- [ ] Implement dry-run harness B0 + G1 (temp-copy, abs seed, CLI for G1)
- [ ] Self-check exit criteria / run commands above
- [ ] Board status + notes
