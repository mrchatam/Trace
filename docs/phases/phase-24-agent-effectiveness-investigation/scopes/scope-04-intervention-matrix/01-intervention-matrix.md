# P24-S04-01 — Intervention matrix synthesis

## Metadata
- id: P24-S04-01
- todo_ids: [P24-S04-01]
- role: implementer (synthesizer / analyst)
- skills: [planning-and-task-breakdown, documentation-and-adrs, analyst, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- agents: [analyst, tech-lead]
- verification: manual (evidence cites + ranking rationale)
- hooks: none

## Objective

Synthesize S01–S03 into a **ranked intervention matrix** and **finalize** Phase 24 findings for handoff.

**Deliverables:**

1. **`INTERVENTION-MATRIX.md`** — ≥8 ranked rows `INT-01…` with locked schema, evidence pointers, Phase 25 theme tags
2. **`FINDINGS.md`** — all sections consolidated (status table → `done`; executive summary; recommended Phase 25 theme(s))
3. **`DR-HANDOFF.md`** — update candidate Phase 25 themes (1–3 recommended) with evidence links from matrix

**Investigation only** — no product Go commits. Interventions are **proposals** for human-promoted Phase 25; do not implement here.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — local-first, no daemon on P0-X
- SoT: [INVESTIGATION.md](../../INVESTIGATION.md) — FM taxonomy, intervention categories, two-mode model
- S04-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S01: [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md)
- S02: [CODEBASE-AUDIT.md](../scope-02-codebase-loop-audit/CODEBASE-AUDIT.md)
- S03: [EXTERNAL-RESEARCH.md](../scope-03-external-research/EXTERNAL-RESEARCH.md)
- Handoff stub: [DR-HANDOFF.md](../../DR-HANDOFF.md)

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md): Agent mode → clarify if ranking trade-offs or Phase 25 theme boundaries unclear → Plan mode → execute. **Do not re-audit code or re-run external research** — cite S01/S02/S03 deliverables.

## Locked defaults

| Item | Value |
|------|-------|
| Output file (primary) | `scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md` |
| Living findings | [FINDINGS.md](../../FINDINGS.md) — S04 consolidates all sections |
| Handoff | [DR-HANDOFF.md](../../DR-HANDOFF.md) — candidate themes table |
| Minimum interventions | **≥8** ranked rows (`INT-01` … `INT-NN`; contiguous IDs) |
| Owner taxonomy | `product` \| `harness` \| `protocol` \| `docs` |
| Impact / Effort | `high` \| `med` \| `low` each |
| Risk taxonomy | `regression` \| `agent confusion` \| `scope creep` (one primary per row; secondary in Notes column optional) |
| Evidence column | Pointer to **S01** / **S02** / **S03** section (e.g. `S01 POSTMORTEM §3 FM-10`, `S02 CODEBASE-AUDIT §2 FM-03`, `S03 EXTERNAL-RESEARCH §5`) |
| Phase 25 theme tags | `P25-A` … `P25-E` per [DR-HANDOFF.md](../../DR-HANDOFF.md) placeholder table (may add `P25-F` only if split required — note in matrix §1) |
| Product Go | **Forbidden** |
| Ranking sort | **Primary:** impact × feasibility for **live Cursor dogfood** (not harness-only wins in isolation) |
| Single-solution claim | **Forbidden** — no row may claim to fix all FM-* |

## Locked matrix schema (columns — do not alter)

Every row in `INTERVENTION-MATRIX.md` §2 **must** use these columns:

| Column | Values / rules |
|--------|----------------|
| **Rank** | 1…N after sort (1 = highest priority) |
| **ID** | `INT-01` … `INT-NN` (contiguous) |
| **Addresses** | Comma-separated `FM-*` IDs (≥1 per row) |
| **Intervention** | **One sentence** — concrete delta (what changes, not vague “improve UX”) |
| **Owner** | `product` \| `harness` \| `protocol` \| `docs` |
| **Impact** | `high` \| `med` \| `low` — on collapsing Mode A→B in dogfood |
| **Effort** | `high` \| `med` \| `low` — implementation cost in Phase 25 |
| **Risk** | Primary: `regression` \| `agent confusion` \| `scope creep` |
| **Evidence** | S01/S02/S03 pointer (section + artifact path) |
| **Phase 25 theme** | `P25-A` … `P25-E` (or new split tag with §1 note) |
| **Notes** | Optional — secondary risk, human gate, spike deferral |

## Ranking rubric (locked — apply before assigning Rank)

Score each candidate intervention on two axes, then sort **descending composite**.

### Impact (1–3) — collapses build vs directed-gap gap

| Score | Label | Criteria |
|-------|-------|----------|
| **3** | high | Directly addresses **both** sessions’ shared FMs (FM-01, FM-03, FM-08) **or** FM-09 mode collapse **or** FM-10 task promotion with Session B evidence |
| **2** | med | Addresses ≥2 FMs or one FM with strong Session A **and** B signal; partial mode collapse |
| **1** | low | Single-session FM, experiment-only, or docs-only without behavior change |

Map to column: 3→`high`, 2→`med`, 1→`low`.

### Feasibility (1–3) — live Cursor dogfood within Trace laws

| Score | Label | Criteria |
|-------|-------|----------|
| **3** | high | Library+CLI or install/hook change only; no daemon; shippable in one Phase 25 scope; peer pattern exists (S03) |
| **2** | med | Product policy change + tests; or harness prompt bundle; or protocol scorer fix |
| **1** | low | Multi-phase spike (SQLite episode model), human product call pending, or high regression surface |

Map to column: 3→`low` effort, 2→`med`, 1→`high` (effort is inverse of feasibility for sorting).

### Composite sort rule

1. Compute **priority score** = Impact_score × Feasibility_score (max 9).
2. Sort rows **descending** by priority score.
3. Tie-breakers (in order): (a) addresses FM-09 or FM-10, (b) lower effort, (c) lower regression risk, (d) more evidence sources (S01+S02 > single scope).
4. Document top-3 ranking rationale in `INTERVENTION-MATRIX.md` §1 (3–5 sentences).

**Anti-pattern:** Do not rank harness-only doc tweaks above product loop/task promotion if they do not change agent-visible behavior in dogfood.

## S03→S04 forwarded residuals (must appear as matrix rows or §4 deferred with rationale)

Each item below needs **≥1** matrix row **or** explicit deferral in §4 with human-gate note:

| Residual | Source | Matrix expectation |
|----------|--------|-------------------|
| **Hop/P19 saturation calibration** | S03 §6; S02 §3 | Trace-specific thresholds; greenfield vs post-gap pass |
| **Auto-spawn from discoveries vs explicit loop apply spawn** | S03 §6; S02 FM-10 | Human gate documented; compare AR `createTask` vs Trace `spawned_tasks[]` |
| **SQLite episode model spike** | S03 §6 (Graphiti patterns) | Defer as spike row or P25 theme split — not full Graphiti port |
| **Cursor hook API drift** | S03 §6 | Harness verification row; reference G1 `trace-loop-gate.sh` |
| **Sticky STOP reason UX** | S02 §3; S03 §5 | Unify export `p19_saturated` vs live `hop_budget_exceeded` |
| **Deliberation reset API** | S02 §3; S03 §5 | Gap-pass transition / hop_count reset / episode boundary |
| **Build-mode vs directed-gap defaults** | S01 §4; INVESTIGATION two-mode | FM-09 collapse — default gap pass or orchestrator Trace-first |

## Must-include intervention seeds (≥1 row each — may merge if one sentence covers both)

Implementer must ensure the final matrix includes rows covering (IDs assigned at synthesis time):

| Seed topic | FM focus | Evidence start |
|------------|----------|----------------|
| Discovery → **`trace add` task promotion** | FM-10 | S01 Session B; S02 `apply.go` spawned_tasks |
| **Hop budget reset** or verify-task recovery after gap pass | FM-03 | S01 Must answer §3; S02 §3 deliberation reset absence |
| **Default gap pass** in harness (collapse Mode A→B) | FM-09 | S01 §4; S03 §5 orchestrator row |
| Orchestrator Trace-first harness (Multitask) | FM-04, FM-05 | S01 Session A; S03 CUR/OH |
| Saturation/hop budget policy for single-session app builds | FM-03 | S02 `policy.go`, `select.go`; S03 §3.5 |
| Seed shape / experiment protocol + **`score.sh` graph count fix** | FM-01, FM-06, FM-07 | S01 git sparsity; protocol |
| MCP tool ordering or post-discovery nudge | FM-08 | S02 MCP tools; S03 AR/CM |
| Export/decisions/uncertainties honesty enforcement | FM-02 | S01 Session A vs B; seed export strict |
| Arm isolation in experiments | FM-06 | S01; protocol |
| **Two-session dogfood rubric** (build vs directed gap) | FM-09 | INVESTIGATION; S01 two-mode table |

## Required reads (mandatory)

Read before ranking. Every matrix **Evidence** cell must trace to one of these.

### Phase 24 SoT + handoff

| # | Path | Extract |
|---|------|---------|
| 1 | [INVESTIGATION.md](../../INVESTIGATION.md) | Intervention categories; FM-01..FM-10; two-mode model; Phase 24 outputs |
| 2 | [FINDINGS.md](../../FINDINGS.md) | Current draft sections — finalize, do not flatten Session A/B |
| 3 | [DR-HANDOFF.md](../../DR-HANDOFF.md) | Placeholder P25-A…E themes |
| 4 | [project-rules.md](../../../../rules/project-rules.md) | Trace law constraints on recommendations |
| 5 | [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) | Local-first, no full-rebuild, portable graph |

### S01 — Dogfood evidence

| # | Path | Extract |
|---|------|---------|
| 6 | [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md) | §2 two-mode; §3 FM matrix; §4 open questions; Must answer §1–§5 |
| 7 | [experiments/RESULTS.md](../../../../../experiments/RESULTS.md) | E01 verdict row |

### S02 — Product mechanisms

| # | Path | Extract |
|---|------|---------|
| 8 | [CODEBASE-AUDIT.md](../scope-02-codebase-loop-audit/CODEBASE-AUDIT.md) | §2 FM mechanism table; §3 S01 residuals; §4 cross-cutting; §6 open → S04 |
| 9 | `internal/loop/apply.go` (skim spawned_tasks) | FM-10 promotion path — cite in matrix, do not re-audit full file |
| 10 | `internal/deliberation/select.go` (skim L7–12) | Sticky STOP / reason-code order — cite only |

### S03 — External comparables

| # | Path | Extract |
|---|------|---------|
| 11 | [EXTERNAL-RESEARCH.md](../scope-03-external-research/EXTERNAL-RESEARCH.md) | §2 comparables; §3 Q-D answers; §5 S02 crosswalk; §6 open gaps |
| 12 | [ENFORCEMENT.md](../../../phase-23-enforcement-choke-points/ENFORCEMENT.md) | Baseline install enforcement — contrast interventions |

## Deliverable templates

### `INTERVENTION-MATRIX.md` — locked structure

Create in this scope folder.

#### §1 Executive summary + ranking rationale

- 3–5 sentences: top intervention + why dogfood impact
- **Ranking rationale:** why INT-01..03 beat lower rows (reference rubric scores)
- Phase 25 theme recommendation preview (1–3 themes)

#### §2 Ranked intervention table (required)

| Rank | ID | Addresses | Intervention | Owner | Impact | Effort | Risk | Evidence | Phase 25 theme | Notes |
|------|-----|-----------|--------------|-------|--------|--------|------|----------|----------------|-------|
| 1 | INT-01 | FM-… | … | product | high | med | agent confusion | S02 §2 FM-03 | P25-B | … |

**Row rules:**

- **Intervention** column: one sentence, imperative voice (“Add …”, “Reset …”, “Require …”)
- **Owner** must match who ships in Phase 25 (product = Go library/CLI; harness = install/prompts/hooks; protocol = experiment/scoring; docs = AGENTS/rules only)
- No duplicate interventions — merge or differentiate in Notes
- Every **Must-include seed** topic covered by ≥1 row

#### §3 FM coverage matrix (required)

| FM-ID | Addressed by (INT-IDs) | Residual gap |
|-------|------------------------|--------------|
| FM-01 | INT-… | … |

All FM-01..FM-10 must appear. Residual gap may be “none” or “deferred Phase 25+”.

#### §4 Deferred / human-gate items

Bullets for interventions **not** ranked (spike, product call, out of Trace law). Include:

- SQLite episode model scope boundary
- Auto-spawn approval policy
- Any peer pattern rejected per S03 §4 anti-patterns

#### §5 Phase 25 theme mapping (required)

| Theme | INT-IDs | One-line scope boundary | Out of scope |
|-------|---------|-------------------------|--------------|
| P25-A | INT-… | … | … |

Mutually scoped — **no mega-phase**. Human promotes **one theme per phase** ([DR-HANDOFF.md](../../DR-HANDOFF.md)).

### FINDINGS.md finalization checklist

Update [FINDINGS.md](../../FINDINGS.md) to **consolidated** state:

| Section | Action |
|---------|--------|
| Status table | All rows → **`done`** with link to final artifact |
| Two-mode model | Keep Session A/B separate; add 1-paragraph synthesis |
| Failure taxonomy | Retain per-session FM table; link POSTMORTEM §3 |
| Codebase audit | Brief bullets + link CODEBASE-AUDIT (no full table duplicate) |
| External comparables | Brief bullets + link EXTERNAL-RESEARCH |
| **Intervention matrix** | New subsection: top 3 INT-IDs + link INTERVENTION-MATRIX.md |
| **Executive summary** | New top section (after Status): 5–8 sentences — problem, root causes, recommended Phase 25 direction |
| **Preliminary conclusion** | Refresh to reflect full investigation (not S01-only) |

**Do not** delete Session A/B evidence tables — synthesize above them.

### DR-HANDOFF.md update template

In [DR-HANDOFF.md](../../DR-HANDOFF.md):

1. Mark S04 checklist item done when matrix + FINDINGS finalized.
2. Replace placeholder **Candidate Phase 25 themes** with **Recommended** subset (1–3 themes):
   - Each theme: one-line scope, **Evidence:** `INT-0x`, `INT-0y` links
   - Mark non-recommended placeholders `deferred` with one-line reason
3. Add **Intervention summary** bullet list (top 3 ranks)
4. Do **not** mark DR-HANDOFF `CLOSED` — S05-02 owns close

Example row format:

```markdown
| P25-A | **Discovery → task promotion** — … | **Recommended.** Evidence: INT-01, INT-04 ([INTERVENTION-MATRIX.md](scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md)) |
```

## Trace law compliance (implementer must respect)

When writing interventions and Phase 25 themes:

1. **No daemon/HTTP MCP on P0-X core path** — optional harness MCP OK to mention; product recommendations stay library+CLI first.
2. **No full-rebuild indexer** — reject Graphiti-style full re-ingest as core intervention.
3. **Portable graph** — seed export/import discipline, not git SHA as identity.
4. **Forward-only** — do not rewrite S01–S03 deliverables; link and synthesize.
5. **Human promotes Phase 25** — matrix ranks options; does not commit roadmap.

## Forbidden

- Implementing interventions (Phase 25 product Go)
- Single solution claiming to fix all FMs
- Re-running full codebase audit or external web research (cite S02/S03)
- Ranking without S01/S02/S03 evidence pointers
- Stacking all P25-A…E themes into one recommended mega-phase

## Preflight / Plan

Before writing matrix rows:

1. Confirm board row **P24-S04-01** runnable (P24-S04-00 `done`).
2. List candidate interventions from Must-include seeds + S02 §6 + S03 §5 hints.
3. Score each with Impact × Feasibility rubric; draft sort order.
4. Map rows to P25-A…E; split if any theme exceeds one MVP.
5. Plan FINDINGS executive summary narrative (two-mode → interventions → Phase 25).

## Exit criteria

- [ ] `INTERVENTION-MATRIX.md` exists with §1–§5 per locked template
- [ ] **≥8** ranked interventions; locked schema columns complete on every row
- [ ] Ranking rubric applied; §1 documents rationale for top 3
- [ ] All **Must-include seed** topics covered
- [ ] All **S03→S04 forwarded residuals** addressed (row or §4 deferral)
- [ ] §3 FM coverage includes FM-01..FM-10
- [ ] §5 Phase 25 themes mutually scoped (1–3 **Recommended** in DR-HANDOFF)
- [ ] `FINDINGS.md` consolidated per checklist (Status `done`, executive summary)
- [ ] `DR-HANDOFF.md` candidate themes updated with INT evidence links
- [ ] No product Go in diff
- [ ] Board row P24-S04-01 set `done` with artifact paths in Notes

## Minimal todos

- [ ] Read required S01/S02/S03 handoff artifacts
- [ ] Draft candidate intervention list from seeds + §6/§5 open items
- [ ] Apply ranking rubric; assign INT-01… ranks
- [ ] Write INTERVENTION-MATRIX.md §1–§5
- [ ] Finalize FINDINGS.md (status table + executive summary + intervention subsection)
- [ ] Update DR-HANDOFF.md recommended themes (1–3)
- [ ] Self-check exit criteria; set board row done

## Next

**P24-S04-02**
