# I — Initial Benchmark Plan

## Hypotheses under early test

| ID | Hypothesis | When |
|----|------------|------|
| H1 | Project graph improves understanding | X0 primary |
| H5 | Evidence/review reduces false completion | X0 partial (honesty demo) + later Gate G |
| H6 | Progressive context reduces tokens without hurting success | X0 secondary metrics |
| H2–H4, H7 | Later phases | Not in X0 |

## Experiment P0-X — Foundation validation (closes P0)

**Not** an agent comparison. Proves the real foundation before Gate C.

### Pass criteria (SETTLED — DR-P0X)

1. Goal / Task / Decision / Discovery round-trip with provenance.
2. Files + minimal symbols/imports represented (tree-sitter).
3. `trace why` returns causal chain with evidence/reason codes.
4. `trace context` produces bounded task-specific context.
5. Human-seeded graph matches fixture ground truth.
6. Deterministic tests pass **several** understanding queries (no LLM; default ≥5).
7. **Incremental update** of a changed file without rebuilding the entire fixture graph.

### Architecture constraint

Full-rebuild-on-any-change is a **failed design** for P0-X, even if functionally correct. Incremental localized update must exist (need not be heavily optimized).

### Corpus

Synthetic fixture only (`fixtures/x0`).

### Kill / revise criteria (foundation)

If P0-X fails after a bounded fix budget, **do not** close P0 and **do not** add MCP/planner/impact—revise model/retrieval/seed/incremental design first.

---

## Experiment X0 — Core falsification (post-P0 / Gate C)

### Comparison

```text
Condition B0: Agent + ordinary repository tools (read files, git, search)
Condition G1: Agent + Trace CLI context/why (still can read repo)
```

MCP is **not** required for X0 if CLI invocation is scripted into the harness.

### Corpus

1. **Synthetic fixture** — required.
2. Real OSS — only after synthetic X0 dry-runs; not for P0-X.

### Seeding

- Scoring uses **human-curated** ground truth.
- Separately measure agent seeding accuracy/cost later (not required to close P0).

### Task families / metrics / protocol

Unchanged from prior plan (understanding primary; implementation + honesty secondary; quality + efficiency metrics; N≥3 runs for agent conditions).

### Kill criteria (product thesis)

- If G1 understanding ≤ B0 within error **and** seeding cost non-trivial → thesis endangered.
- Do not confuse P0-X pass with Gate C pass.

## Baselines roadmap

Track over time: raw agent → repo tools → code graph → project graph → progressive context → planner → capability selection.

P0-X requires none of the agent baselines. X0 requires (2) and (4) at minimum via CLI.

## Reproducibility

- Pin tool versions, model names, seeds, fixture SHAs.
- Store metrics JSON schema in `evals/x0/schema.json` (created during T012).
- Prefer public fixtures; never require private customer code for core claims.

## What X0 does not prove

- Multi-month progressive planning superiority (H2).
- Production monorepo scale (H Gate).
- Vendor-neutral env graph value (H7).
