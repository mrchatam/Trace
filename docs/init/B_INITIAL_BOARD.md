# B — Initial Board (HISTORICAL)

> **Superseded for execution.** Run [`docs/TODO.md`](../TODO.md).  
> This file remains as the pre-harness task breakdown that Phase 00 scopes were derived from (`T001`≈S01 … `T012a`≈S08/S09).

Ordered board for Scope `S0` / P0-X (archived mapping source).


---

## P0-X critical path

```text
T001 → T002 → T003 → T004 → T005 → T007 → T008 → T009 → T011 → T012a
```

`T004` (tree-sitter symbols/imports) is **on the critical path** (DR-P0X items 2 & 7).  
`T006` / `T013` / `T012b` = Phase 1.  
`T010` (MCP) deferred until after context API validation.

---

## T001 — Repository & toolchain scaffold

| Field | Content |
|-------|---------|
| ID | T001 |
| Objective | Create Go module `github.com/mrchatam/Trace`, `cmd/trace`, internal layout, lint/test, git init if needed |
| Why | No code exists; all later tasks need a home |
| Dependencies | None remaining (Round 2 closed) |
| Scope | `go.mod` with real module path, `cmd/trace`, package layout per C_FIRST_SCOPE, smoke test, README stub, `.gitignore` (include `.trace/`) |
| Out of scope | Domain logic, schema, HTTP/MCP |
| Current project context | Docs-only; Apache-2.0 LICENSE present |
| Relevant decisions | DR-LANG, DR-NAME, DR-GOMOD, DR-SURFACE, DR-RISK, D13 |
| Relevant assumptions | A9, A16 |
| Required skills/tools/MCPs | Go toolchain; git |
| Exit criteria | `go test ./...` green; `go build -o bin/trace ./cmd/trace` succeeds; module path correct |
| Verification plan | Build + test logs; LICENSE intact; no daemon package |
| Expected evidence | Command logs; tree listing |
| Possible downstream effects | Module path fixed for all imports |
| Known unknowns | Exact internal package taxonomy |

---

## T002 — SQLite project store + thin event log

| Field | Content |
|-------|---------|
| ID | T002 |
| Objective | Implement per-project SQLite under **`.trace/`**; schema for core entities + events; migrations |
| Why | Persistence substrate for all semantics |
| Dependencies | T001 |
| Scope | Create `.trace/` (ensure gitignored); tables for Project, Goal, Decision, Assumption, Task, Discovery, Claim, Evidence, Review, PlanChange, Event, File/Symbol; provenance columns; migration runner |
| Out of scope | Retrieval, Git indexing, analyzers |
| Relevant decisions | DR-EVT, DR-CLAIM, DR-DB, DR-TRACEDIR |
| Exit criteria | Create project DB; Goal+Task+Event roundtrip tests; no source BLOB table |
| Verification plan | Unit tests; schema dump |
| Expected evidence | Test output |
| Known unknowns | Final column set |

---

## T003 — VCS adapter + Git CLI implementation + history index

| Field | Content |
|-------|---------|
| ID | T003 |
| Objective | Define **VCS adapter interface**; implement with **`git` CLI**; index commit metadata + path changes; resolve content via Git; support incremental history refresh |
| Why | Gate A; law: Git canonical; swappable backend later |
| Dependencies | T002 |
| Scope | Interface + gitcli impl; `history(file)`, `commits_between`, `changes(commit)`, `last_changed(file)`; incremental by new commits; **no full history rewrite on every call** |
| Out of scope | go-git impl; storing full diffs permanently |
| Relevant decisions | D2, DR-GIT, DR-INCREMENTAL, DR-RISK |
| Exit criteria | Tests on temp repo; content equals `git show`; incremental refresh test; interface has ≥1 fake/mock for unit tests |
| Verification plan | Deterministic git fixture |
| Expected evidence | Tests |
| Known unknowns | Large-repo performance later |

---

## T004 — Structural analyzers (TS/JS + Python via tree-sitter)

| Field | Content |
|-------|---------|
| ID | T004 |
| Objective | Emit File + **minimal Symbol/Import** edges via **tree-sitter**; **file-level incremental** updates mandatory |
| Why | P0-X items 2 & 7; structure is fundamental (DR-PARSE, DR-RISK) |
| Dependencies | T003 |
| Scope | Analyzer adapter interface; tree-sitter grammars for TS/JS + Python; ignore `.gitignore` + binary; changing one file updates only that file’s structural edges (and direct dependents if required by import graph consistency—document chosen rule) |
| Out of scope | Perfect call graphs; other languages; semantic LLM summaries; full-project rebuild path as default |
| Relevant decisions | DR-ANLANG, DR-PARSE, DR-INCREMENTAL, DR-P0X, DR-RISK, D3 |
| Relevant assumptions | A8 |
| Required skills/tools/MCPs | tree-sitter / go-tree-sitter |
| Exit criteria | Golden edges on fixture; **incremental test**: edit one file → only localized reindex (assert untouched files’ edge rows unchanged / not rewritten en masse) |
| Verification plan | Golden lists + incremental assertion (P0-X #7) |
| Expected evidence | Golden files + test logs |
| Known unknowns | Exact symbol granularity “minimal” |

---

## T005 — Work/causal API (CRUD + state transitions)

| Field | Content |
|-------|---------|
| ID | T005 |
| Objective | Canonical core API to create/link Goal, Decision, Assumption, Task, Discovery, PlanChange; task state machine with actors |
| Why | Causal layer for why-chains |
| Dependencies | T002 |
| Scope | Validation; provenance defaults; forbid implementer self-promotion to DONE without review policy hook (can stub policy) |
| Out of scope | Autoplanning; impact engine |
| Current project context | PROJECT_MODEL; PLANNING philosophy |
| Relevant decisions | D6, D8, D9, DR-CLAIM, DR-CHURN |
| Relevant assumptions | A4 |
| Required skills/tools/MCPs | none special |
| Exit criteria | API tests for links Goal→Task, Decision→Task, Discovery→PlanChange; illegal transitions rejected |
| Verification plan | Unit tests for state machine |
| Expected evidence | Test logs |
| Possible downstream effects | Planner consumes same API |
| Known unknowns | Exact state enum naming |

---

## T006 — Claim / Evidence / Review path

| Field | Content |
|-------|---------|
| ID | T006 |
| Objective | Implement Claim→Evidence→Review→VerifiedFact promotion rules; command evidence capture helper |
| Why | Anti-hallucination core |
| Dependencies | T005 |
| Scope | Evidence types: command result, file hash, test report path; Review PASS/FAIL/UNCERTAIN; promotion rules per law |
| Out of scope | Multi-model adversarial review automation |
| Current project context | REVIEW_AND_VERIFICATION; SECURITY trust |
| Relevant decisions | D6, D7 (todo only), DR-FAIL |
| Relevant assumptions | A3 |
| Required skills/tools/MCPs | process execution sandbox caution |
| Exit criteria | Cannot mark DONE without required evidence; honesty unit test with fake claim |
| Verification plan | Policy tests; secret-redaction smoke on fake AWS key in stdout |
| Expected evidence | Tests |
| Possible downstream effects | Eval honesty suite |
| Known unknowns | Redaction pattern set |

---

## T007 — Retrieval engine (exact + lexical + graph)

| Field | Content |
|-------|---------|
| ID | T007 |
| Objective | Hybrid candidate generation with reason codes; depth/budget caps |
| Why | Context without dumps |
| Dependencies | T003, T004, T005 |
| Scope | Exact; FTS; neighbor expansion; filters; **no embeddings** |
| Out of scope | Semantic embeddings; learned rankers |
| Current project context | RETRIEVAL_AND_CONTEXT; DR-NOSSEM |
| Relevant decisions | D4, D5, DR-NOSSEM |
| Relevant assumptions | A6 |
| Required skills/tools/MCPs | SQLite FTS5 |
| Exit criteria | Queries return reasons; hard cap enforced; dump API absent |
| Verification plan | Unit tests on fixture graph |
| Expected evidence | Tests + sample traces |
| Possible downstream effects | Compiler input |
| Known unknowns | Ranking weights |

---

## T008 — Context compiler

| Field | Content |
|-------|---------|
| ID | T008 |
| Objective | Compile Layer 0–1 packets (JSON+MD), token estimate, untrusted labeling |
| Why | Agent-facing value surface |
| Dependencies | T007 |
| Scope | `get_task_context`, `expand_context`, `explain_entity`/`why` |
| Out of scope | Layers 2–3 auto; fancy compression |
| Current project context | D5; SECURITY prompt injection |
| Relevant decisions | DR-PACK, D5 |
| Relevant assumptions | A1 |
| Required skills/tools/MCPs | tokenizer estimate lib optional |
| Exit criteria | Packet schema stable enough for eval; expansion explicit; injection labeling present |
| Verification plan | Snapshot tests; budget never exceeded |
| Expected evidence | Snapshots |
| Possible downstream effects | MCP tool schemas |
| Known unknowns | Q-INJECTION final wrapping |

---

## T009 — CLI adapter

| Field | Content |
|-------|---------|
| ID | T009 |
| Objective | Expose S0 operations via CLI calling core only |
| Why | Human + harness entrypoint |
| Dependencies | T003–T008 |
| Scope | Commands listed in C_FIRST_SCOPE |
| Out of scope | Rich TUI |
| Current project context | ARCHITECTURE adapters |
| Relevant decisions | DR-API, DR-NAME |
| Relevant assumptions | A9 |
| Required skills/tools/MCPs | CLI framework |
| Exit criteria | Scripted CLI walkthrough on fixture succeeds |
| Verification plan | Integration test / tape script |
| Expected evidence | Terminal transcript |
| Possible downstream effects | Dogfood |
| Known unknowns | UX flags |

---

## T010 — MCP adapter — DEFERRED (post P0-X)

| Field | Content |
|-------|---------|
| ID | T010 |
| Objective | MCP tools mirroring CLI semantics |
| Why | Agent integration — but must not shape core architecture |
| Dependencies | P0-X passed; T008/T009 validated; DR-AGENT |
| Scope | Bounded tools only; thin adapter over Go library |
| Out of scope | Raw SQL; any business logic in MCP layer |
| Status | **DEFERRED** — not on P0-X critical path |
| Relevant decisions | DR-AGENT, DR-API, DR-SURFACE |
| Exit criteria | (when scheduled) parity with CLI query/context ops |
| Known unknowns | Host quirks |

---

## T011 — Synthetic fixture + ground truth seed

| Field | Content |
|-------|---------|
| ID | T011 |
| Objective | Create `fixtures/x0` mini-repo with known architecture + JSON ground-truth causal graph |
| Why | Enables Gate C without legal ambiguity |
| Dependencies | T009 (usable CLI) or core APIs |
| Scope | Small multi-module app; README; seed file for Trace import |
| Out of scope | Huge real monorepo |
| Current project context | Q-BENCH; Q-SEED |
| Relevant decisions | DR-SLICE |
| Relevant assumptions | A11 |
| Required skills/tools/MCPs | none |
| Exit criteria | Seed import creates why-chains matching ground truth doc |
| Verification plan | Diff actual graph vs ground truth |
| Expected evidence | Export JSON |
| Possible downstream effects | X0 scoring |
| Known unknowns | Ideal fixture size |

---

## T012 — Eval harness

Split:

### T012a — Deterministic P0-X harness (P0 critical)

| Field | Content |
|-------|---------|
| ID | T012a |
| Objective | Prove **all 7 P0-X criteria** without an LLM agent |
| Why | Closes P0 (DR-P0, DR-P0X) |
| Dependencies | T004, T005, T007–T009, T011 |
| Scope | GT load; several understanding queries (≥5 default); why/context assertions; **incremental reindex test**; emit `metrics-p0x.json` |
| Out of scope | Agent baseline comparison |
| Exit criteria | Every DR-P0X item green in CI-able script |
| Verification plan | Fail the harness if any of the 7 criteria fail |
| Expected evidence | `metrics-p0x.json` + logs |

### T012b — Agent X0 harness (Phase 1, post-P0)

| Field | Content |
|-------|---------|
| ID | T012b |
| Objective | Scripts for baseline vs graph-enabled **agent** tasks |
| Why | Gate C |
| Dependencies | P0 closed; T012a; agent access |
| Scope | As in `I_BENCHMARK_PLAN.md` Experiment X0 |
| Exit criteria | Dry-run metrics for B0 and G1 |
| Known unknowns | Agent nondeterminism → repeats |

---

## T013 — Honesty review demo

| Field | Content |
|-------|---------|
| ID | T013 |
| Objective | Plant false completion claim; show review rejects without evidence / catches missing test |
| Why | Partial H5 signal inside S0 |
| Dependencies | T006, T009, T011 |
| Scope | One scripted scenario |
| Out of scope | Full Gate G |
| Current project context | Controlled honesty tests |
| Relevant decisions | D6, DR-FAIL |
| Relevant assumptions | A3 |
| Required skills/tools/MCPs | none |
| Exit criteria | Documented transcript: FAIL with evidence refs |
| Verification plan | Run scenario twice |
| Expected evidence | Transcript |
| Possible downstream effects | Review UX tweaks |
| Known unknowns | none critical |

---

## T014 — Scope S0 review & Gate C readiness report

| Field | Content |
|-------|---------|
| ID | T014 |
| Objective | Independent scope review against C_FIRST_SCOPE criteria; list gaps; declare X0 ready/not |
| Why | Close S0 properly |
| Dependencies | T001–T013 |
| Scope | Review notes; update question ledger; no new features unless blocker |
| Out of scope | Starting Phase 3 planner |
| Current project context | H_VERIFICATION_STRATEGY |
| Relevant decisions | D7, D14 |
| Relevant assumptions | all S0 |
| Required skills/tools/MCPs | fresh reviewer agent recommended |
| Exit criteria | Written Go/No-Go for running X0 experiment |
| Verification plan | Checklist in C_FIRST_SCOPE verification section |
| Expected evidence | Review document in `docs/init/` or `docs/reports/` |
| Possible downstream effects | Phase 2 begins or rework tasks filed |
| Known unknowns | none |

---

## Board order

```text
P0-X critical:
T001 → T002 → T003 → T004 → T005 → T007 → T008 → T009 → T011 → T012a
                                    (T004 on path; incremental required)

Post-P0 / Phase 1:
T006 → T013 → T012b → T014
T010 (MCP) only after context API validated
```
