# E02 hypotheses — P25-C validation

Maps to Phase 24 interventions and Phase 25 deliverables.

## H-P25-1 (INT-03 — default gap pass)

**Claim:** Installed `GapPassPrompt` causes agents to run a post-build gap review without custom human instruction.

**Pass signal:** Graph has ≥1 discovery **or** agent transcript shows `trace gap` / `trace loop status` after product commits and **before** final export.

**Fail signal:** E01 Session A replay — 0 discoveries, build-only graph.

## H-P25-2 (INT-04 — orchestrator Trace-first)

**Claim:** Multitask parent sets `TRACE_TASK_ID` and subagents inherit it; parent does not delegate without workspace + task UUID.

**Pass signal:** Hook fired with task ID; graph writes attributed to orchestrator turn before subagents; no “wrong path” edits.

**Fail signal:** FM-04 — workers-only Trace; parent never calls loop/gate.

## H-P25-3 (Mode collapse — FM-09)

**Claim:** Single build session + P25 install achieves **partial** Session B richness (discoveries/decisions) without Session B prompt.

**Pass signal:** discoveries ≥ 1 AND decisions ≥ 1 AND tasks still seed-only is **partial pass** (promotion still needs P25-A).

**Full pass:** discoveries ≥ 1 AND (new tasks ≥ 1 OR verify gate allows edit after gap pass).

## H-P25-4 (Product unchanged — out of scope)

P25 did **not** ship INT-01/02/05. E02 **must not** fail P25-C solely because:

- verify task still `hop_budget_exceeded`
- 0 new tasks despite discoveries

Those are **expected residuals** → Phase 26 P25-A/B if E02 otherwise passes H-P25-1/2.

## Baseline comparison

| Metric | E01 Session A | E01 Session B | E02 G1 target |
|--------|---------------|---------------|---------------|
| Human gap prompt | No | Yes | **No** |
| Discoveries | 0 | 7 | ≥ 1 |
| Decisions | 0 | 2 | ≥ 1 |
| New tasks | 0 | 0 | 0 OK (until P25-A) |
| G1 ≡ B0 code | Yes | No | No (optional) |
