# P40-S02-02 — DR-HANDOFF Phase 40+ close

## Metadata
- id: P40-S02-02
- todo_ids: [P40-S02-02]
- role: closer
- skills: [documentation-and-adrs, writing-for-agents, planning-and-task-breakdown, shipping-and-launch]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed (spot-check + scaffold)

## Objective

Independent **fresh-session** review of S02-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 40+ DR-HANDOFF** with explicit **Phase 41+** successor (**never TBD**). Deliver **runnable Phase 41+ scaffold** per [REMEDIATION-PLAN.md](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) §Phase 41+ (entry themes **G8 + G9**). Update `docs/TODO.md` + `AGENTS.md`. Phase 40 complete when this row is `done`. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S02-00 locks
- [01-verify.md](01-verify.md) — locked verify floor (FINAL — S02-00)
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S02-01
- [REMEDIATION-PLAN.md](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) — G8/G9 ranks
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [INTAKE.md](../../INTAKE.md)
- Phase 39 [DR-HANDOFF.md](../../../phase-39-context-orient-harness/DR-HANDOFF.md) — predecessor close pattern
- Pattern: [P39-S03-02](../../../phase-39-context-orient-harness/scopes/scope-03-verify/02-dr-handoff.md)
- Pattern: [P38 S07-02](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-07-verify/02-dr-handoff.md)
- [docs/TODO.md](../../../../TODO.md)
- [docs/TODO/phase-40.md](../../../../TODO/phase-40.md)
- [AGENTS.md](../../../../../AGENTS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S02-01 verifier. Unattended: execute until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-02-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p40-s02-01-verify/evidence/` |
| Phase handoff | [DR-HANDOFF.md](../../DR-HANDOFF.md) |
| G5 deliverable | S00 implement + S00-02 APPROVE |
| G2 deliverable | S01 implement + S01-02 APPROVE |
| Phase board | [docs/TODO/phase-40.md](../../../../TODO/phase-40.md) |

## Locked DR-HANDOFF close policy (FINAL — S02-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S02-01** — VERIFY-NOTES + evidence dir; DR-HANDOFF stays **OPEN** |
| Who closes | **S02-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S02-01 `done`; verify blocks 0–6 green per VERIFY-NOTES + independent spot-check |
| Default successor | **Phase 41+ — Layers & intent** (human promotes **P41-00**) |
| Entry themes | **G8** Progressive layers L2–L3 + **G9** Intent pipeline |
| Secondary queue | G6, G7 per REMEDIATION-PLAN rank — document in Phase 41 INTAKE, not P40 rows |
| Idle alternative | **`no successor`** — if human explicitly defers Phase 41 (document reason in Notes) |
| Must not | Leave `Successor decision: TBD`; rewrite S00–S01 `done` history; ship product in this row |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 40 board rows `done` |
| Portable graph | Confirm `trace/graph.json` current if entity changes during P40 |

### Successor decision table (locked — pick exactly one)

| Outcome (from S02-01 + spot-check) | Decision | Next action |
|-------------------------------------|----------|-------------|
| VERIFY blocks 0–6 green; G5+G2 delivered | **Phase 41+ — G8 + G9 entry** | Close DR-HANDOFF; scaffold Phase 41+; mark Phase 40 **done** |
| Same as above but human defers implementation | **`no successor`** | Close DR-HANDOFF; Phase 40 done; Phase 41 folder optional stub-only |
| Block 0 FAIL (G5 tests/UI red) | **Do not close** — spawn repair | Keep OPEN; insert S02-02a/b |
| Block 1 FAIL (G2 / tool count) | **Do not close** — send back S01 or S02-01 | Keep OPEN |
| Block 2 FAIL (M-001) | **Do not close** — spawn repair | Keep OPEN |
| Block 3 FAIL (Laws 6–7/19) | **Do not close** — send back S00/S01 or S02-01 | Keep OPEN |
| VERIFY-NOTES missing or evidence dir absent | **Do not close** — send back S02-01 | Keep OPEN |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`**.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] **S00** G5 GUI graph orient — orient panel on `/` Explore; install hook narrative; Law 19 adapter only
- [ ] **S01** G2 unified explore — task-aware capped `trace_explore` (17th MCP tool); `compiler.Explore` library-first
- [ ] **S02** VERIFY + successor documented (VERIFY-NOTES + REVIEW-NOTES)

### P40 outcome summary (for DR-HANDOFF prose)

One paragraph covering:

1. **G5 delivered** — graph-first onboarding UX on `/` Explore route (`GraphOrientPanel`, dismiss, confidence labels, CONTRIBUTING graph-first GUI)
2. **G2 delivered** — unified `trace_explore` (17th read-only MCP tool) + CLI `trace explore`; `internal/compiler/explore.go` library-first
3. **M-001 preserved** — task loop + gates primary; explore optional convenience after compose-first; 9/17 stale hygiene
4. **Forward queue** — G6/G7 secondary; G8/G9 → Phase 41+

### Top themes for human promotion (locked default)

| Rank | Theme | GAP ids | Phase 41+ scope sketch |
|------|-------|---------|------------------------|
| 1 | **G8** Progressive layers L2–L3 | G-003 | Ship L2–L3 in compiler or revise spec with documented alternative |
| 2 | **G9** Intent pipeline | G-009 | Implement intent extraction or mark doc aspirational + supersede |

Secondary queue (document in Phase 41 INTAKE, do not scaffold implement rows at close unless human promotes): **G6** non-semantic concept retrieval (G-004b), **G7** index freshness & langs (G-005).

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| HTTP `/v1/explore` route | Not shipped — MCP+CLI sufficient; optional future adapter |
| G6/G7 | Secondary queue — human may promote before G8/G9 |
| G-004a vector | Permanent defer — DR-NOSSEM |
| Redundant double `dismissOrient()` on G5 dismiss | Low nit from S00-02 — idempotent |
| `instructions.go:30` Phase 39 S02 stub | Optional Phase 41 doc hygiene |

## Spot-check floor (minimum — do not trust Notes alone)

```bash
cd /home/ali/Desktop/Trace

# Verify artifacts exist
test -f docs/phases/phase-40-read-surface-retrieval-depth/scopes/scope-02-verify/VERIFY-NOTES.md
grep -q 'PASS' docs/phases/phase-40-read-surface-retrieval-depth/scopes/scope-02-verify/VERIFY-NOTES.md

# Evidence dir referenced in VERIFY-NOTES exists
EVID=$(grep -o 'experiments/runs/[^ ]*' \
  docs/phases/phase-40-read-surface-retrieval-depth/scopes/scope-02-verify/VERIFY-NOTES.md \
  | head -1)
test -d "$EVID/evidence" || test -d "${EVID}evidence"

# G5 artifact
test -f web/src/components/GraphOrientPanel.tsx
grep -q 'graph-orient-panel' web/src/components/GraphOrientPanel.tsx

# G2 artifact — library + MCP
test -f internal/compiler/explore.go
grep -q 'trace_explore' internal/mcp/server.go
grep -q 'compiler.Explore' internal/mcp/tools_explore.go

# Tool count 17
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered 2>&1 | tail -3

# Moat lead intact + explore optional
grep -q 'trace_tasks' internal/mcp/instructions.go
grep -q 'trace_explore' internal/mcp/instructions.go
grep -q '9/17' internal/mcp/instructions.go

# G2 acceptance subset
go test ./internal/compiler/... -count=1 -run 'TestExploreTaskRequired|TestExploreNoDump' 2>&1 | tail -3

# Stale hygiene docs
grep -q '9/17' CONTRIBUTING.md
```

## Phase 41+ scaffold (if successor = Phase 41+)

Deliver runnable minimum per agent-loop-protocol Phase handoff:

| Artifact | Minimum content |
|----------|-------------------|
| `docs/phases/phase-41-layers-intent/README.md` | Goal G8+G9, in/out, predecessor Phase 40 |
| `docs/phases/phase-41-layers-intent/00-PHASE-PLANNER.md` | Runnable planner row — lock G8+G9 scopes |
| `docs/phases/phase-41-layers-intent/INTAKE.md` | Human locks; secondary queue G6/G7 |
| `docs/phases/phase-41-layers-intent/DR-HANDOFF.md` | OPEN scaffold; successor TBD at P41-00 only |
| Per-scope stubs | S00 G8 layers + S01 G9 intent: `00-PLANNER` / `01-*` / `02-*` / `SCOPE-TODOS.md` minimal |
| `docs/TODO/phase-41.md` | Board rows; **P41-00** planner first pending |
| `docs/TODO.md` | Index link Phase 41+; orchestrator snippet |

**G6/G7:** Document in Phase 41 INTAKE secondary queue — **do not** add P40 S03/S04 rows retroactively.

### Phase 41 scope sketch (for scaffold stubs)

| Scope | Theme | Deliverable sketch |
|-------|-------|-------------------|
| S00 | G8 Progressive layers L2–L3 | Compiler layer expansion or spec revise per REMEDIATION-PLAN G8 |
| S01 | G9 Intent pipeline | Intent extraction implement or doc-revise aspirational §3 |
| S02 | VERIFY + DR-HANDOFF | Gate + successor (Phase 42+ or `no successor`) |

## DR-HANDOFF.md update template (on APPROVE — default Phase 41+)

```markdown
# DR-HANDOFF — Phase 40+

**Status:** CLOSED (P40-S02-02 YYYY-MM-DD)

| Field | Value |
|-------|-------|
| Opened | 2026-08-22 (scaffold at P39-S03-02) |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 39 CLOSED |
| Theme | Read surface & retrieval depth — G5+G2 delivered |
| Successor decision | **Phase 41+ — Layers & intent** (G8 + G9; human promotes P41-00) |

## Scope checklist (closed)

- [x] **S00** G5 GUI graph orient
- [x] **S01** G2 unified explore — 17 MCP tools
- [x] **S02** VERIFY + successor documented

## Outcome

G5 graph-first GUI orient + G2 unified trace_explore delivered. M-001 moat preserved.
Secondary queue G6/G7 forwarded to Phase 41 INTAKE.

## Successor

**Phase 41+ — Layers & intent** — G8 L2–L3 + G9 intent pipeline.
```

## Deliverables

| Artifact | Action |
|----------|--------|
| `REVIEW-NOTES.md` | Independent review verdict (APPROVE + confidence) |
| `DR-HANDOFF.md` | Status → **CLOSED**; successor explicit (**Phase 41+ G8/G9** or `no successor`) |
| `docs/TODO/phase-40.md` | All rows `done` |
| `docs/TODO.md` | Phase 40+ → done; Phase 41 index if scaffolded |
| `AGENTS.md` | Current focus → Phase 41+ scaffold pending human promotion |

### REVIEW-NOTES template

```markdown
# REVIEW-NOTES — Phase 40 / S02-02

**Date:** …
**Verdict:** APPROVE | REJECT
**Confidence:** high | medium | low

## Spot-check results
| Check | Result |
| VERIFY-NOTES PASS | |
| G5 GraphOrientPanel | |
| G2 compiler.Explore + trace_explore | |
| 17 tools | |
| Moat lead + 9/17 | |

## DR-HANDOFF
- Status: CLOSED
- Successor: Phase 41+ — G8 + G9 (never TBD)

## Scaffold delivered
- [ ] phase-41-layers-intent folder
- [ ] docs/TODO/phase-41.md
- [ ] docs/TODO.md index

## Residuals (non-blocking)
…

## Next
P41-00 (human promotion) or idle
```

## Exit criteria

- [ ] REVIEW-NOTES.md with APPROVE + confidence
- [ ] DR-HANDOFF **CLOSED** with Phase 41+ G8/G9 or `no successor`
- [ ] Phase 40 board all rows `done`
- [ ] Phase 41+ scaffold runnable (if default successor)
- [ ] Board row → `done` with evidence in Notes

## Next

**P41-00** (human promotion) or idle
