# P41-S02-02 — DR-HANDOFF Phase 41+ close

## Metadata
- id: P41-S02-02
- todo_ids: [P41-S02-02]
- role: closer
- skills: [documentation-and-adrs, writing-for-agents, planning-and-task-breakdown, shipping-and-launch]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed (spot-check + scaffold)

## Objective

Independent **fresh-session** review of S02-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 41+ DR-HANDOFF** with explicit **Phase 42+** successor (**never TBD**). Deliver **runnable Phase 42+ scaffold** per [REMEDIATION-PLAN.md](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) §Phase 42+ (entry themes **G6 + G7** default). Update `docs/TODO.md` + `AGENTS.md`. Phase 41 complete when this row is `done`. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S02-00 locks
- [01-verify.md](01-verify.md) — locked verify floor (FINAL — S02-00)
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S02-01
- [REMEDIATION-PLAN.md](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) — G6/G7 ranks
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [INTAKE.md](../../INTAKE.md)
- Phase 40 [DR-HANDOFF.md](../../../phase-40-read-surface-retrieval-depth/DR-HANDOFF.md) — predecessor close pattern
- Pattern: [P40 S02-02](../../../phase-40-read-surface-retrieval-depth/scopes/scope-02-verify/02-dr-handoff.md)
- Pattern: [P39-S03-02](../../../phase-39-context-orient-harness/scopes/scope-03-verify/02-dr-handoff.md)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-41.md](../../../../TODO/phase-41.md)
- [AGENTS.md](../../../../../AGENTS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S02-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-02-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p41-s02-01-verify/evidence/` |
| Phase handoff | [DR-HANDOFF.md](../../DR-HANDOFF.md) |
| G8 deliverable | S00 implement + S00-02 APPROVE (board row **697**) |
| G9 deliverable | S01 implement + S01-02 APPROVE (board row **700**) |
| Phase board | [docs/TODO/phase-41.md](../../../../TODO/phase-41.md) |

## Locked DR-HANDOFF close policy (FINAL — S02-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S02-01** — VERIFY-NOTES + evidence dir; DR-HANDOFF stays **OPEN** |
| Who closes | **S02-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S02-01 `done`; verify blocks 0–6 green per VERIFY-NOTES + independent spot-check |
| Default successor | **Phase 42+ — G6/G7 secondary queue** (human promotes **P42-00**) |
| Entry themes | **G6** Non-semantic concept retrieval + **G7** Index freshness & langs |
| Secondary queue | Per REMEDIATION-PLAN rank — document in Phase 42 INTAKE, not P41 rows |
| Idle alternative | **`no successor`** — if human explicitly defers Phase 42 (document reason in Notes) |
| Must not | Leave `Successor decision: TBD`; rewrite S00–S01 `done` history; ship product in this row |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 41 board rows `done` |
| Portable graph | Confirm `trace/graph.json` current if entity changes during P41 |

### Successor decision table (locked — pick exactly one)

| Outcome (from S02-01 + spot-check) | Decision | Next action |
|-------------------------------------|----------|-------------|
| VERIFY blocks 0–6 green; G8+G9 delivered | **Phase 42+ — G6/G7 entry** | Close DR-HANDOFF; scaffold Phase 42+; mark Phase 41 **done** |
| Same as above but human defers implementation | **`no successor`** | Close DR-HANDOFF; Phase 41 done; Phase 42 folder optional stub-only |
| Block 0 FAIL (G8 tests/layer) | **Do not close** — spawn repair | Keep OPEN; insert S02-02a/b |
| Block 1 FAIL (G9 / §3 doc) | **Do not close** — send back S01 or S02-01 | Keep OPEN |
| Block 2 FAIL (M-001) | **Do not close** — spawn repair | Keep OPEN |
| Block 3 FAIL (Laws 6–7/19) | **Do not close** — send back S00/S01 or S02-01 | Keep OPEN |
| VERIFY-NOTES missing or evidence dir absent | **Do not close** — send back S02-01 | Keep OPEN |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`**.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] **S00** G8 progressive layers L2–L3 — opt-in `max_layer` ship (`ContextOptions.MaxLayer` default 1)
- [ ] **S01** G9 intent pipeline — rule-based `ExtractIntent` + §3 revised (DR-NOSSEM semantic defer)
- [ ] **S02** VERIFY + successor documented (VERIFY-NOTES + REVIEW-NOTES)

### P41 outcome summary (for DR-HANDOFF prose)

One paragraph covering:

1. **G8 delivered** — L2/L3 opt-in via `max_layer` (CLI `--max-layer`, MCP `max_layer`); default L0–L1 preserved (`compiler.go:34–38`; G8-L1–L7 green; S00-02 APPROVE)
2. **G9 delivered** — rule-based `ExtractIntent` in `internal/retrieval/intent.go`; compiler+explore wired; §3 intent shipped + DR-NOSSEM (S01-02 APPROVE)
3. **M-001 preserved** — task loop primary; layer/intent merge into packet; compile/explore require task_id
4. **Forward queue** — G6/G7 secondary for Phase 42+

### Top themes for human promotion (locked default)

| Rank | Theme | GAP ids | Phase 42+ scope sketch |
|------|-------|---------|------------------------|
| 1 | **G6** Non-semantic concept retrieval | G-004b | Graph-label channel without vector |
| 2 | **G7** Index freshness & langs | G-005 | Analyzer/lang policy + index honesty |

Secondary queue already delivered in P41: **G8** progressive layers, **G9** intent pipeline.

**Rejects preserved:** G-004a vector, product dual-index default, query-only moat replacement, full-graph dump defaults, LLM intent extraction.

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| HTTP `max_layer` route absent | Not shipped — CLI+MCP sufficient (S00-02 low) |
| G6/G7 | Secondary queue — Phase 42+ default |
| G-004a vector | Permanent defer — DR-NOSSEM |
| `IntentSummary` JSON-only (not Markdown render) | Low nit from S01-02 |
| Search multi-OR vs `FTSQuery()` doc path | Low doc drift from S01-02 — behavior OK |
| Trim comment vs layer-only sort | Nit from S00-02 |
| `TaskContext` godoc still "L0–L1" | Nit from S00-02 |

## Spot-check floor (minimum — do not trust Notes alone)

```bash
cd /home/ali/Desktop/Trace

# Verify artifacts exist
test -f docs/phases/phase-41-layers-intent/scopes/scope-02-verify/VERIFY-NOTES.md
grep -q 'PASS' docs/phases/phase-41-layers-intent/scopes/scope-02-verify/VERIFY-NOTES.md

# Evidence dir referenced in VERIFY-NOTES exists
EVID=$(grep -o 'experiments/runs/[^ ]*' \
  docs/phases/phase-41-layers-intent/scopes/scope-02-verify/VERIFY-NOTES.md \
  | head -1)
test -d "$EVID/evidence" || test -d "${EVID}evidence"

# G8 artifact — MaxLayer default + enrich
grep -q 'MaxLayer' internal/compiler/compiler.go
test -f internal/retrieval/layer_enrich.go
grep -q 'max_layer\|MaxLayer' internal/mcp/tools_context.go

# G9 artifact — intent library + Search wiring
test -f internal/retrieval/intent.go
grep -q 'ExtractIntent' internal/retrieval/search.go
grep -q 'IntentInput' internal/compiler/compiler.go

# §3 doc revised
grep -q 'Shipped.*Phase 41' docs/RETRIEVAL_AND_CONTEXT.md
grep -q 'DR-NOSSEM' docs/RETRIEVAL_AND_CONTEXT.md

# G8 acceptance subset
go test ./internal/compiler/... -count=1 -run 'TestContextDefaultLayer1|TestContextMaxLayer2|TestNoDumpAPI' 2>&1 | tail -3

# G9 acceptance subset
go test ./internal/retrieval/... -count=1 -run 'TestExtractIntentFromTask|TestSearchUsesIntent|TestIntentNoSemantic' 2>&1 | tail -3

# MCP max_layer mirror
go test ./internal/mcp/... -count=1 -run TestMCPContextMaxLayer2 2>&1 | tail -3

# Moat lead intact
grep -q 'trace_tasks' internal/mcp/instructions.go
grep -q 'trace_context' internal/mcp/instructions.go

# Default caps unchanged
grep -q 'DefaultTokenBudget = 4096' internal/compiler/packet.go
```

## Phase 42+ scaffold (if successor = Phase 42+)

Deliver runnable minimum per agent-loop-protocol Phase handoff:

| Artifact | Minimum content |
|----------|-------------------|
| `docs/phases/phase-42-*/README.md` | Goal G6+G7, in/out, predecessor Phase 41 |
| `docs/phases/phase-42-*/00-PHASE-PLANNER.md` | Runnable planner row — lock G6+G7 scopes |
| `docs/phases/phase-42-*/INTAKE.md` | Human locks; rejects preserved |
| `docs/phases/phase-42-*/DR-HANDOFF.md` | OPEN scaffold; successor TBD at P42-00 only |
| Per-scope stubs | S00 G6 + S01 G7: `00-PLANNER` / `01-*` / `02-*` / `SCOPE-TODOS.md` minimal |
| `docs/TODO/phase-42.md` | Board rows; **P42-00** planner first pending |
| `docs/TODO.md` | Index link Phase 42+; orchestrator snippet |

**G6/G7:** Document in Phase 42 INTAKE — **do not** add P41 S03/S04 rows retroactively.

### Phase 42 scope sketch (for scaffold stubs)

| Scope | Theme | Deliverable sketch |
|-------|-------|-------------------|
| S00 | G6 Non-semantic concept retrieval | Graph-label channel without vector (G-004b) |
| S01 | G7 Index freshness & langs | Analyzer/lang policy + index honesty (G-005) |
| S02 | VERIFY + DR-HANDOFF | Gate + successor (Phase 43+ or `no successor`) |

## DR-HANDOFF.md update template (on APPROVE — default Phase 42+)

```markdown
# DR-HANDOFF — Phase 41+

**Status:** CLOSED (P41-S02-02 YYYY-MM-DD)

| Field | Value |
|-------|-------|
| Opened | 2026-08-22 (scaffold at P40-S02-02) |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 40 CLOSED |
| Theme | Layers & intent — G8+G9 delivered |
| Successor decision | **Phase 42+ — G6/G7 secondary queue** (human promotes P42-00) |

## Scope checklist (closed)

- [x] **S00** G8 progressive layers L2–L3 — opt-in max_layer
- [x] **S01** G9 intent pipeline — rule-based ExtractIntent + §3 revised
- [x] **S02** VERIFY + successor documented

## Outcome

G8 opt-in L2/L3 layers + G9 rule-based intent pipeline delivered. M-001 moat preserved.
Secondary queue G6/G7 forwarded to Phase 42 INTAKE.

## Successor

**Phase 42+ — G6/G7 secondary queue** — non-semantic concept retrieval + index freshness.
```

## Deliverables

| Artifact | Action |
|----------|--------|
| `REVIEW-NOTES.md` | Independent review verdict (APPROVE + confidence) |
| `DR-HANDOFF.md` | Status → **CLOSED**; successor explicit (**Phase 42+ G6/G7** or `no successor`) |
| `docs/TODO/phase-41.md` | All rows `done` |
| `docs/TODO.md` | Phase 41+ → done; Phase 42 index if scaffolded |
| `AGENTS.md` | Current focus → Phase 42+ scaffold pending human promotion |

### REVIEW-NOTES template

```markdown
# REVIEW-NOTES — Phase 41 / S02-02

**Date:** …
**Verdict:** APPROVE | REJECT
**Confidence:** high | medium | low

## Spot-check results
| Check | Result |
| VERIFY-NOTES PASS | |
| G8 MaxLayer + layer_enrich | |
| G9 intent.go + Search wiring | |
| §3 intent shipped + DR-NOSSEM | |
| G8-L1 default layer≤1 | |
| G9-I5 no semantic | |
| Moat lead intact | |

## DR-HANDOFF
- Status: CLOSED
- Successor: Phase 42+ — G6 + G7 (never TBD)

## Scaffold delivered
- [ ] phase-42-* folder
- [ ] docs/TODO/phase-42.md
- [ ] docs/TODO.md index

## Residuals (non-blocking)
…

## Next
P42-00 (human promotion) or idle
```

## Exit criteria

- [ ] REVIEW-NOTES.md with APPROVE + confidence
- [ ] DR-HANDOFF **CLOSED** with Phase 42+ G6/G7 or `no successor`
- [ ] Phase 41 board all rows `done`
- [ ] Phase 42+ scaffold runnable (if default successor)
- [ ] Board row → `done` with evidence in Notes

## Next

**P42-00** (human promotion) or idle
