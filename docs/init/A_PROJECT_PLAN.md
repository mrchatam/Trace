# A — Project Plan (Coarse)

## Goal

Build a local-first project knowledge/causal graph with progressive planning that measurably improves agent understanding, planning adaptation, and verification versus repository contents + Git alone—without replacing Git, IDEs, or coding models.

## Major architecture (early)

```text
Human / eval harness / (later) Agent
    ↓
CLI adapter  (`trace`)
    ↓
Canonical Go library
    ├── Project DB (SQLite): work/causal/events/provenance
    ├── Retrieval: exact + lexical + graph (+ temporal via Git)
    ├── Context compiler
    ├── Git adapter
    └── Code analyzer adapters (TS/JS, Python first)

Later adapters (not in P0-X): MCP → same library; optional daemon/HTTP
```

Language: **Go**. Surface: **library + CLI only** until core is validated. Semantic embeddings, environment graph, impact engine, multi-agent, UI, daemon: **later**, experiment-gated.

## Knowledge posture for this plan

| Class | Content |
|-------|---------|
| KNOWN | Design baseline; Round 1–2 answers; DR-P0X 7-point bar; module `github.com/mrchatam/Trace` |
| ASSUMED | A1–A16 |
| UNRESOLVED | Soft: Q-UNDERSTAND-N (≥5 default), Q-FIXTURE-LANG |
| EXPECTED_DISCOVERY | tree-sitter symbol granularity; schema migrations; CLI UX |
| UNKNOWN | Gate C effect size; commercial packaging; scale ceiling |

---

## Phase 0 — Foundation (docs + tiny experiment P0-X)

P0 is **not** closed by documentation alone (DR-P0).

### Phase 0a — Specification

- **Objective:** Lock laws/registers; define P0-X and X0; Round 1 decisions recorded.
- **Expected outcome:** `docs/init/*` + Round 1 answers (done); Round 2 answers.
- **Validation gate:** Round 2 blockers answered or defaults accepted → code may start at `T001`.

### Phase 0b — Tiny experiment P0-X (still Phase 0)

- **Objective:** Validate causal model, **minimal structural graph (tree-sitter)**, retrieval/context, human seeding, **several** deterministic understanding queries, and **incremental file update**—on a synthetic fixture—without LLM agent comparison.
- **Reason:** Reduce wrong-product risk (DR-RISK); prove foundation is real before declaring P0 complete.
- **Expected outcome:** All **7 DR-P0X** criteria green via `T012a`.
- **Dependencies:** Phase 0a; P0-X critical path including **T004**.
- **Main risks:** Treating CRUD-only as success; full-rebuild architecture; skipping incremental.
- **Important unknowns:** Exact “minimal symbol” grain; query set (≥5 default).
- **Validation gate:** DR-P0X 7/7 → **P0 closes**.

**Note:** Phase 1 no longer “adds structural graph later”—structure is in P0-X. Phase 1 adds honesty review, richer harness, agent X0.

## Phase 1 — Complete vertical slice & Experiment X0 readiness

- **Objective:** Finish S0 pieces not required for P0-X (honesty review path, richer agent harness) so **agent** baseline-vs-graph X0 can run via CLI.
- **Reason:** Gate C needs the measurement instrument beyond deterministic P0-X.
- **Expected outcome:** X0 executable end-to-end via CLI (MCP still after validation).
- **Dependencies:** P0 closed.
- **Main risks:** Scope creep; treating P0-X pass as Gate C pass.
- **Important unknowns:** Agent nondeterminism; seeding cost for fair compare.
- **Validation gate:** X0 dry-run emits metrics for conditions B0 and G1.

## Phase 2 — Gate C evaluation & slice hardening

- **Objective:** Run X0; measure understanding accuracy, misses, tokens, latency, task success vs baseline.
- **Reason:** Kill or continue product expansion.
- **Expected outcome:** Go/No-Go report; issue list.
- **Dependencies:** Phase 1.
- **Main risks:** Unfair baseline; contaminated ground truth; tiny sample.
- **Important unknowns:** Effect size; variance.
- **Validation gate:** Gate C decision documented (pass/fail/iterate).

## Phase 3 — Progressive planner (minimal)

- **Objective:** Coarse goal→phase→scope; deep-plan current scope; discovery→PlanChange propagation with churn controls.
- **Reason:** Test H2/H3 only after H1 shows promise (or in parallel only if Gate C strongly positive on understanding alone).
- **Expected outcome:** Replan demo on planted discovery.
- **Dependencies:** Phase 2 Go (or explicit user override).
- **Main risks:** Planning bureaucracy.
- **Important unknowns:** Replan quality metrics.
- **Validation gate:** Gate E mini-eval on fixture.

## Phase 4 — Review depth & evidence policies

- **Objective:** Scope review layer; richer verification policies; residual tracking.
- **Reason:** H5 depth.
- **Expected outcome:** Honesty suite escape-rate report.
- **Dependencies:** Phase 1 review path; preferably Phase 2.
- **Main risks:** Cost explosion.
- **Important unknowns:** ROI of independent LLM review.
- **Validation gate:** Gate G preliminary.

## Phase 5 — Decision impact & simulation

- **Objective:** Impact classes + alternatives; later `plan simulate`.
- **Reason:** H4.
- **Expected outcome:** Impact report UX; precision/recall on planted conflicts.
- **Dependencies:** Stable work graph.
- **Main risks:** False confidence in impact.
- **Important unknowns:** Graph completeness for impact.
- **Validation gate:** Gate F preliminary.

## Phase 6 — Environment/capability graph

- **Objective:** Skills/rules/MCP/tool selection for tasks.
- **Reason:** H7.
- **Expected outcome:** Task packets include only required capabilities; missing capability warnings.
- **Dependencies:** Agent adapter maturity.
- **Main risks:** Ontology bloat.
- **Important unknowns:** Skill metadata standards across harnesses.
- **Validation gate:** Capability-selection ablation.

## Phase 7 — Performance ladder & language plugins

- **Objective:** Incremental indexing quality; ignore tiers; more languages via adapters; Gate H sizes.
- **Reason:** Prove practicality.
- **Expected outcome:** Benchmark tables 10k–1M LOC.
- **Dependencies:** Solid core.
- **Main risks:** Premature optimization theater.
- **Important unknowns:** SQLite ceiling.
- **Validation gate:** Gate H thresholds (to be set after first measurements).

## Phase 8 — Ecosystem & hardening

- **Objective:** Stable plugin APIs; multi-agent worktrees; production concerns (migrations, backup, auth).
- **Reason:** Adoption beyond research.
- **Expected outcome:** Versioned APIs; contributor analyzer plugin.
- **Dependencies:** Gates C–G confidence.
- **Main risks:** Locking bad APIs early.
- **Important unknowns:** Hosted offering shape.
- **Validation gate:** Compatibility + security checklist.

---

## Milestone gates (from ROADMAP, operationalized)

| Gate | Meaning | Earliest phase |
|------|---------|----------------|
| A | Useful history reconstruction via Git index joins | 0b (P0-X) |
| B | Why-questions on seeded causal chains | 0b (P0-X) |
| C | Agents better with graph than without | 2 (X0) |
| D | Progressive planning beats static upfront | 3+ |
| E | Discoveries improve future plans | 3 |
| F | Impact predicts consequences | 5 |
| G | Review reduces false completion | 4 |
| H | Practical on large repos | 7 |

**P0 close gate:** P0-X pass (deterministic foundation), not Gate C.

## Explicit non-goals until gates say otherwise

Cloud SaaS, swarm orchestration, proprietary graph DB, universal language support, auto-summarize-everything, unrestricted autonomy, desktop dashboard-first UX.
