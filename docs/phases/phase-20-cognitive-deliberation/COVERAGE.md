# Phase 20 — TRACE_THOUGHTPROCESS coverage map

Source: [`docs/TRACE_THOUGHTPROCESS.md`](../../TRACE_THOUGHTPROCESS.md). Every numbered section is owned. Nothing is silently dropped.

**P20-00 audit (2026-08-18):** Walked §§1–32 against this table — all sections mapped. §16 Experiments and §18 Risk-adaptive verification remain **Future** (not S01–S07 implement). §29 subsections A–T covered via README architecture + per-scope doc maps.

**P20-S02-00 audit (2026-08-18):** Entity merge table **unchanged**. Uncertainties/hypotheses new; Assumption/Decision/Discovery/DecisionAlternative reused; Requirement/Constraint/Finding/Risk/Option stay merged (no new tables).

**P20-S03-00 audit (2026-08-18):** Change row thickened (`changes` + `change_paths` + `effects`). Paths = child table (JSON rejected). Metric/Observation still merge into effect comparison — no Metric table. Tests/baseline/scores stay S04. Contradiction may link S02 Hypothesis or spawn Discovery; does not fork Discovery as hypothesis.

**P20-S05-00 audit (2026-08-18):** Regression/Reflection/Relationship rows thickened. `regressions` derived from S04 `comparison_json` flags or S03 contradicted effect; attribution never auto-`caused`. `reflections` structured JSON arrays (no essay body). Observed vs causal stay on `entity_links` (no Relationship table). §16/§18 still Future.

Legend: **Must** this phase · **Should** this phase if it fits the locked MVP · **Future** named, not boarded as implement work · **Reuse** already in Trace

| Doc § | Topic | Disposition | Home |
|------:|-------|-------------|------|
| 1 | Externalized cognition loop | Must | S01 controller |
| 2 | Store structured outputs, not raw CoT | Must | S02 artifacts + all writes |
| 3A | Engineering graph | Must (thin links) | S02 + existing Goal/Task/File/Symbol |
| 3B | Deliberation graph / dynamic transitions | Must | S01 |
| 4 | Cognitive phases ORIENT…REPLAN | Must | S01 |
| 5 | State-driven selection, deterministic policies | Must | S01 |
| 6 | Entry/exit conditions | Must | S01 |
| 7 | Uncertainty first-class | Must | S02 |
| 8 | Assumptions create/invalidate/replan | Must | S02 (extend existing Assumption) |
| 9 | Decisions, alternatives, reconsideration | Must | S02 (extend existing Decision + DecisionAlternative) |
| 10 | Change as first-class object | Must | S03 |
| 11 | Test ≠ Verify ≠ Evaluate | Must | S04 |
| 12 | Gates; implemented ≠ verified | Must | S04 |
| 13 | Baselines | Should | S04 thin (git SHA + score snapshot) |
| 14 | Expected vs actual effects | Must | S03 |
| 15 | Regression detection + correlated_vs_caused | Must (thin) | S05 |
| 16 | Experiments as objects | Future | README Future; not S01–S07 implement |
| 17 | Historical cause/effect knowledge | Should | S05 observed_relationship vs causal |
| 18 | Risk-adaptive verification | Future | policy stub comment only |
| 19 | Reflection and learning | Must (thin) | S05 reflection artifact |
| 20 | Verification / quality debt | Must | S04 |
| 21 | Full feedback architecture | Must as state machine | S01; not a rigid workflow |
| 22 | Agent interaction model | Must | S06 extends P19 loop apply/next |
| 23 | Context selection by phase | Must | S06 |
| 24 | Model-agnostic / harness-agnostic | Must | S06; stdout-first inherited from P19 |
| 25 | Observability / audit why Trace chose X | Must | S01 events + S07 verify |
| 26 | Avoid premature complexity | Must | this matrix + MVP cuts |
| 27 | Map onto existing repo | Must | 00-PHASE-PLANNER live inventory |
| 28 | Reuse P19 loop | Must | S01/S06 wrap `internal/loop`, do not replace |
| 29A–T | Architecture through tests/migration | Must as design | README + per-scope locks |
| 30 | Challenge the concept | Must | README § Challenge |
| 31 | Target user story | Should (mini-eval) | S07 VERIFY |
| 32 | Incremental, smallest viable first | Must | MVP vs Future below |

## Entity merge (doc §29B — do not clone every noun)

| Proposed noun | Trace choice |
|---------------|--------------|
| Goal | **Reuse** `goals` |
| Requirement | **Merge** into Goal body + optional `requirement` link rel; no new table in MVP |
| Constraint | **Merge** into Assumption (`kind` or body prefix) |
| Task | **Reuse** `tasks` |
| Question / Uncertainty | **New** `uncertainties` (question text + severity + status). Do not fork Discovery |
| Assumption | **Reuse** + invalidate API + link to affected decisions/tasks |
| Finding | **Reuse** Discovery (`severity` already exists) |
| Hypothesis | **New thin** `hypotheses` (status, confidence, linked evidence). Not a second Discovery |
| Option | **Reuse** `DecisionAlternative` |
| Decision | **Reuse** + reconsideration triggers |
| Risk | **Merge** into Uncertainty (`kind=risk`) or Assumption |
| Change | **New** `changes` + `change_paths` + `effects` (git SHA/path refs, not blobs; no tests/baseline/score columns — S04) |
| Test / Verification / Evaluation | **New thin** result records, three kinds, not three runners |
| Baseline | **New thin** `baselines` (commit + scores JSON) |
| Metric / Observation | **Merge** into evaluation result + effect comparison |
| Regression | **New thin** `regressions` derived from evaluation `comparison_json` flags **or** contradicted effect; `attribution=correlated\|hypothesized\|caused` (create=`correlated`; never auto-`caused`) |
| Experiment | **Future** |
| Reflection | **New thin** `reflections` with structured JSON arrays (`invalidated_assumptions`, `new_dependencies`, `useful_tests`) — no essay-only `body` |
| Evidence / Claim | **Reuse** |
| Deliberation / Phase | **New** `deliberation_state` (one row per goal/task seed), not a graph of every thought |
| Feature | **Skip** — Goal/Task |
| Relationship | **Reuse** `entity_links` + rels `observed_relationship` (confidence) and `caused_by` (evidence required); no Relationship table |

## Future (explicitly not this phase implement)

- First-class Experiment objects and multi-candidate bake-offs (§16)
- Risk-adaptive test-selection policy engine (§18)
- Learning which deliberation sequences work per model (§5 future, §24 model bake-off)
- Performance/memory/CPU collectors as product features
- ML phase scoring
- Hosted/daemon orchestration
- Replacing Git with Trace-owned diffs
- Full-graph dumps; embeddings
- Autonomous code-edit loop that self-runs the compiler without an agent
