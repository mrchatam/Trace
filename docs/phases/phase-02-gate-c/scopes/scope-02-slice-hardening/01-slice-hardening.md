# P02 / S02 / 01 — Slice hardening

## Metadata
- id: P02-S02-01
- todo_ids: [P02-S02-01]
- role: implementer
- skills: [incremental-implementation, debugging-and-error-recovery, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Implement **measurement-driven** fixes for Gate C issues **GC-01** and **GC-02** only. Keep honesty Paths A/B/C, p0x 7/7, and x0 dry-run + Gate C harness tests green. Do **not** feature-factory, reopen the Gate C kill debate, or claim a new Gate C pass from this row. No daemon/HTTP/embeddings; no progressive planner product.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) locks
- Issue SoT: [`../scope-01-x0-gate-c/GATE-C-NOTES.md`](../scope-01-x0-gate-c/GATE-C-NOTES.md) § Issue list
- S01 review: [`../scope-01-x0-gate-c/REVIEW-NOTES.md`](../scope-01-x0-gate-c/REVIEW-NOTES.md)
- Pins / fairness brief: [`../../../../verification/gate-c-x0/pins.md`](../../../../verification/gate-c-x0/pins.md)
- Query bank q3: [`../../../../../evals/x0/queries.json`](../../../../../evals/x0/queries.json) (`discovery_plan_change_chain`)

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Depends-on / input
- S01-02 **APPROVE high**; Gate C verdict **Go** (mean G1 0.800 > B0 0.000); kill not fired.
- Import **only** GC-NN rows below. Do not invent Phase 03+ roadmap work.

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Verdict context | Gate C **Go** already recorded — this scope **hardens the slice**, does not re-litigate kill / scoring fairness debates |
| Issue order | **GC-01 first**, then **GC-02**. Defer **GC-03** and **GC-04** (`defer: true`) |
| Bars | `evals/honesty` Paths A/B/C; `evals/p0x` 7/7; `evals/x0` dry-run + Gate C recorded harness stay PASS |
| Out | Daemon/HTTP/embeddings; MCP-required Gate C; progressive planner; live-model pack refresh (GC-03); N-increase significance (GC-04); changing Gate C **Go** verdict text; rewriting S01 `done` prompts |
| Fixture seed | Keep `fixtures/x0/seed/gt.json` UUIDs + `discovery_causes_plan_change` link; do **not** invent a fake task↔discovery edge just to game Expand depth |
| Gate C packs | Mode-B packs under `evals/x0/testdata/gate-c/` are **historical evidence** — do **not** rewrite them to pretend q3 always passed. Product proof = new/updated unit tests on why/context |
| CGO | Product/retrieval/compiler tests as today; full bar: `CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./... -count=1` |
| G19 | Library changes in `internal/{store,retrieval,compiler}` only as needed; thin CLI already shells to them — no domain fork in MCP |

---

## Imported issue set (locked)

### Implement — GC-01 (medium) — priority 1

| Field | Value |
|-------|-------|
| id | GC-01 |
| severity | medium |
| metric | G1 `understanding_accuracy` capped at 0.8; **q3 incorrect on all 3 runs** (0 critical_miss) |
| evidence | `docs/verification/gate-c-x0/metrics-g1.json` per_query q3; GATE-C-NOTES per-run table; live `trace why`/`context` on task `2222…` omit discovery→plan_change |
| root cause (live) | `Why` expands from the **task** at `defaultExpandDepth=1`. Fixture `discovery_causes_plan_change` links discovery↔plan_change only — **not** reachable from the task hop graph (goal + `decision_affects_task`). `TaskContext` Expand/FTS likewise miss those entities for the wiring-task title |
| proposed_fix_surface | `internal/retrieval` Why (+ Expand helper) and `internal/compiler` TaskContext linkage so discovery↔plan_change appears in the **task** neighborhood |
| defer | **false** — implement |

**Locked fix approach (choose this; do not re-debate):**

1. Add store support to list links by `rel` (e.g. `ListLinksByRel("discovery_causes_plan_change")`) if missing — thin SQL on `entity_links`.
2. When `Why` seed is a **task** (and when `TaskContext` builds from a task), **attach** each `discovery_causes_plan_change` edge’s endpoints as explanation/context neighbors with `ReasonCode` = `discovery_causes_plan_change` (existing constant `ReasonDiscoveryCausesPlanChg`). Prefer a shared retrieval helper used by both Why and compiler Expand path so CLI `why`/`context` stay consistent.
3. Respect existing budgets (`MaxItems` / token trim). Prefer including the discovery↔plan_change pair over unrelated FTS noise when near budget.
4. **Do not** change Expand’s global max depth to “fix” reachability, and **do not** retarget seed GT by linking discovery/plan_change onto the task unless a review later spawns that (out of this locked approach).

**Acceptance (GC-01):**
- [ ] Regression test (retrieval and/or compiler): after seeding fixture-equivalent Goal/Task/Decision + Discovery/PlanChange + `LinkDiscoveryPlanChange`, `Why(ctx, "task", taskID)` includes **both** discovery and plan_change steps with reason `discovery_causes_plan_change` (titles or IDs present).
- [ ] `TaskContext` / ExpandContext path for that task includes those entities (JSON items and/or why_trace when `IncludeWhy`).
- [ ] Existing Expand depth-2 discovery→plan_change tests still pass; depth still rejected outside 1..2.
- [ ] No daemon/HTTP/embeddings.

### Implement — GC-02 (medium) — priority 2

| Field | Value |
|-------|-------|
| id | GC-02 |
| severity | medium |
| metric | Fairness residual — B0 could ace GT if seed/README oracle is readable (packs used oracle-exclusion brief only) |
| evidence | `fixtures/x0/seed/gt.json` + `fixtures/x0/README.md` UUID table; `pins.md` Agent brief; `evals/x0/testdata/gate-c/README.md` |
| proposed_fix_surface | fixtures/x0 layout / eval harness agent-hide of GT for live Gate C **or** documented oracle policy |
| defer | **false** — implement |

**Locked fix approach:**

1. **Agent-facing fixture README** must not publish the ground-truth UUID table / causal answer key. Move the UUID map + causal story detail to an **evaluator-only** doc under `evals/x0/` (e.g. `evals/x0/GT-MAP.md`) and/or keep SoT solely in `seed/gt.json` (harness already imports abs path).
2. Keep a short agent-safe `fixtures/x0/README.md` (layout, Apache-2.0, languages, seed-path note) **without** listing the five UUIDs / discovery→plan_change answer.
3. Strengthen oracle policy in `evals/x0/testdata/gate-c/README.md` + `docs/verification/gate-c-x0/pins.md` Agent brief: B0/G1 must not treat `seed/gt.json` or evaluator GT docs as answer oracle; harness may still read seed for import/grading.
4. Add a small automated guard (preferred): test that `fixtures/x0/README.md` does **not** contain the stable UUID literals `11111111-…` … `55555555-…` (or at least not the discovery/plan_change pair used by q3).
5. After fixture README change: recompute fixture content hash and update `pins.md` (and any Notes that cite the old hash). Do **not** alter committed Gate C metrics scores.

**Acceptance (GC-02):**
- [ ] Agent-visible `fixtures/x0/README.md` has no GT UUID oracle table.
- [ ] Evaluator map exists somewhere under `evals/x0/` (or explicit “SoT = seed/gt.json only” in gate-c README + pins).
- [ ] `pins.md` hash + Agent brief updated.
- [ ] Guard test or equivalent assert in `evals/x0` (or documented skip with reason — prefer test).

### Deferred — do not implement

| id | severity | why deferred | S02 action |
|----|----------|--------------|------------|
| GC-03 | low | Model pin is `recorded-operator-sim/v1`, not a production coding model | Leave `defer: true`; no pack refresh |
| GC-04 | low | Within-condition variance = 0 (identical N=3 grades) | Leave `defer: true`; no N-increase / significance work |

---

## Target files (expected blast radius)

```text
# GC-01 (primary)
internal/store/links.go              # ListLinksByRel (or equivalent)
internal/store/store_test.go         # optional
internal/retrieval/why.go            # attach discovery↔plan_change for task seeds
internal/retrieval/expand.go         # optional shared helper hook
internal/retrieval/*.go              # helper + ReasonCode reuse
internal/retrieval/retrieval_test.go # GC-01 regression
internal/compiler/compiler.go        # only if TaskContext needs explicit attach beyond Why/Expand
internal/compiler/compiler_test.go   # GC-01 context regression

# GC-02
fixtures/x0/README.md                # strip UUID oracle
evals/x0/GT-MAP.md                   # new evaluator-only map (name may vary)
evals/x0/testdata/gate-c/README.md   # oracle policy
evals/x0/*_test.go                   # README UUID absence guard
docs/verification/gate-c-x0/pins.md  # hash + brief
```

Do not touch MCP product surface unless a reverse-import / compile break forces a one-line fix (note it).

## How to run (Notes must cite)

```bash
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./internal/store/... -count=1
CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Optional after fixture edit:

```bash
find fixtures/x0 -type f ! -path '*/.git/*' | sort | xargs sha256sum | sha256sum
# → update pins.md Fixture content hash
```

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts. Record GC-01/GC-02 evidence (test names) in Notes. Explicitly note GC-03/GC-04 still deferred.

## Exit criteria
- [ ] GC-01 acceptance checks above met (why + context surface discovery→plan_change for task seed)
- [ ] GC-02 acceptance checks above met (agent-hide / oracle policy + pins hash)
- [ ] GC-03 / GC-04 remain deferred with reason in Notes
- [ ] Honesty + p0x + x0 + `./...` PASS (commands above)
- [ ] No daemon/HTTP/embeddings; no Gate C verdict rewrite; Mode-B packs not falsified
- [ ] Board row P02-S02-01 status + Notes only

## Minimal todos
- [ ] GC-01: store ListLinksByRel (if needed) + retrieval Why/TaskContext attach + regression tests
- [ ] GC-02: strip fixture README oracle; evaluator map + pins/hash + guard test
- [ ] Run honesty / p0x / x0 / `./...` bars
- [ ] Board status + Notes (cite test names; confirm deferrals)
