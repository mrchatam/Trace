# Phase 44 / S05 / DR-HANDOFF close

## Metadata
- id: P44-S05-02
- todo_ids: [P44-S05-02]
- role: reviewer
- skills: [documentation-and-adrs, writing-for-agents, code-review-and-quality]
- mcps: [user-trace]
- verification: mixed

## Objective

Independent **fresh-session** review of S05-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 44 [`DR-HANDOFF.md`](../../DR-HANDOFF.md)** with explicit successor decision (**never TBD**). Default successor **`no successor`** unless human names next phase. Update orchestrator surfaces (`docs/TODO.md`, `AGENTS.md` current focus) per [agent-loop-protocol Phase handoff](../../../../rules/agent-loop-protocol.md). Phase 44 complete when this row is `done`. **No product code.** Do **not** implement a successor in this row.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S05-00 locks
- [01-verify.md](01-verify.md) — locked verify floor V1–V7
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S05-01
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — D1–D5 + §10
- [`DR-HANDOFF.md`](../../DR-HANDOFF.md)
- [`docs/TODO.md`](../../../../TODO.md)
- [`docs/TODO/phase-44.md`](../../../../TODO/phase-44.md)
- [`AGENTS.md`](../../../../../AGENTS.md)
- Pattern: [P34 S05-02](../../../phase-34-gui-packaging-multiproject/scopes/scope-05-verify/02-dr-handoff.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S05-01 verifier. Unattended: execute review loop until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-05-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p44-s05-01-verify/evidence/` |
| Design locks | `01-DESIGN-LOCKS.md` (LOCKED S01) |
| Phase handoff | `DR-HANDOFF.md` |
| Phase board | `docs/TODO/phase-44.md` |
| Prior scope artifacts | S00 RESEARCH; S01 locks; S02–S04 board Notes / reviews |
| Smoke fixture | `web/testdata/p44-scope-smoke-seed.json` |

## Locked DR-HANDOFF close policy (FINAL — S05-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S05-01** — VERIFY floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S05-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S05-01 `done`; V1–V7 green per VERIFY-NOTES + independent spot-check; §10 ticked; V8 backend regression green |
| Default successor | **`no successor`** — D1–D5 shipped (scope records + MVP rels + CLI infer + GUI 500 adapter); DR-HANDOFF residuals are deferred, **not** a new phase by default |
| Next phase scaffold | **Only if human promotes** — do not invent Phase 45 |
| Regression path | Spawn `P44-S05-02a` implement + `02b` review; **do not** close Phase 44 |
| Must not | Leave `Successor decision: TBD`; rewrite S00–S04 `done` history; ship product in this row; run VERIFY gate in S05-00 planner row |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 44 board rows `done` |

### Successor decision table (locked — pick exactly one)

| Outcome (from S05-01 + independent spot-check) | Decision | Next action |
|------------------------------------------------|----------|-------------|
| VERIFY V1–V7 green; residuals only as listed (plan_scopes mapping, pagination, Codegraph) | **`no successor`** | Close DR-HANDOFF; mark Phase 44 **done** in TODO/AGENTS; idle orchestrator paste |
| Auth smoke / MVP rels / infer CLI-only / Law 19 / caps FAIL | **Do not close** — spawn repair | Keep OPEN; insert 02a/02b; successor = `pending repair spawn` |
| VERIFY-NOTES missing blocks, §10 unticked, or evidence dir absent | **Do not close** — spawn repair or send back S05-01 | Keep OPEN |
| VERIFY PASS but **human** names next theme before this row | **Phase 45** — theme named by human in Notes | Scaffold per protocol **only if human promotes** |
| Seed export / INFERRED round-trip FAIL | **Do not close** — spawn repair | Keep OPEN |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`** (still not TBD).

### Successor scaffold expectations (only if human promotes)

Per [agent-loop-protocol Phase handoff](../../../../rules/agent-loop-protocol.md), the **closing** phase owns the next-phase scaffold — not a later ad-hoc session. Before marking this row `done` with a follow-on successor:

1. **Folder:** `docs/phases/phase-NN-<slug>/` with at least:
   - `README.md` (goal, scope list, in/out)
   - `00-PHASE-PLANNER.md` (runnable)
   - `DR-HANDOFF.md` (OPEN)
   - Per-scope stubs: `00-PLANNER` / `01-*` / `02-*` / `SCOPE-TODOS.md` (minimal OK)
2. **Board:** `docs/TODO/phase-NN.md` with phase planner as **first pending** row after Phase 44’s last `done` row
3. **Index:** link in `docs/TODO.md` phase boards table; orchestrator paste → Active phase NN / next runnable `PNN-00`
4. **`AGENTS.md`** Current focus → Phase NN (or idle if `no successor`)
5. Deep tasking belongs to the **next** phase planner — this row delivers a **runnable handoff**, not finished implement prompts

If default **`no successor`**: do **not** invent Phase 45 stubs. Update TODO/AGENTS to Phase 44 **complete** / idle (Phases 00–44 closed).

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-44-scoped-graph-linking/scopes/scope-05-verify/VERIFY-NOTES.md
test -d experiments/runs/*-p44-s05-01-verify/evidence 2>/dev/null || ls experiments/runs/ | grep p44-s05-01
go test ./internal/domain/ -run 'TestScopeMVPRelsAndCRUD|TestInferScopes_ExplicitWins|TestINFERREDScopeMemberSeedRoundTrip|TestSeedExportImportScopesRoundTrip' -count=1
go test ./internal/retrieval/ -run 'TestProjectGraphScopeEdgesAndFilter|TestProjectGraphTruncationHonesty' -count=1
go test ./cmd/trace/ -run 'TestGraphInferCLI_' -count=1
rg -n 'InferScopes|graph infer' internal/index internal/watch internal/mcp web/src && echo FAIL || echo PASS-cli-only
grep -n 'PROJECT_MAX_NODES = 500\|EDGE_OVERVIEW_MAX = 150' web/src/lib/overviewCompose.ts web/src/lib/graphLayout.ts
test -f web/testdata/p44-scope-smoke-seed.json
test -f trace/graph.json
```

Confirm VERIFY-NOTES: overall PASS; §10 V1–V7 ticked; V1 browser LIVE **or** WAIVED with API/board cite; DR-HANDOFF still OPEN before this row closes it.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [`DR-HANDOFF.md`](../../DR-HANDOFF.md):

- [ ] S00 research artifact committed (`RESEARCH.md`)
- [ ] S01 design locks APPROVE (`01-DESIGN-LOCKS.md` LOCKED)
- [ ] S02 backend links APPROVE (migration 029, MVP rels, export)
- [ ] S03 inference APPROVE (CLI infer, H1 seed honesty, C1 CLI-only)
- [ ] S04 GUI APPROVE (scope clustering, dashed inferred, smoke)
- [ ] S05 VERIFY + successor documented (**never TBD**; default **`no successor`**)
- [ ] Laws 6–7 / Law 19 / M-001 VERIFY blocks PASS (V4–V6)
- [ ] `trace/graph.json` export evidence if schema changed (V7)

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| Scope filter pagination/tiles at 500 cap | Deferred — DR-HANDOFF queue |
| `plan_scopes` ↔ graph scope slug auto-mapping quality | Deferred — soft-map only in S03 |
| Codegraph complement edges | Out of scope — dual-stack doc only |
| Live `trace/graph.json` empty scopes (`omitempty`) | Accept if export cmd + round-trip tests PASS |
| Synthesized goal→task SeedLinks omit provenance | Pre-existing; not INFERRED path |
| Browser smoke WAIVED with API + S04 cite | Accept if V1 substance evidenced |

### DR-HANDOFF.md update template (on APPROVE — default `no successor`)

```markdown
# Phase 44 — DR-HANDOFF

**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Opened | 2026-08-23 |
| Closed | YYYY-MM-DD |
| Predecessor | Phase 43 CLOSED |
| Theme | Scoped graph linking & scope semantics (D1–D5) |
| Outcome | Thin `scopes` + MVP rels in `entity_links`; CLI `trace graph infer`; GUI scope clustering @ 500/150; seed export honesty |
| Successor decision | **no successor** |
| Residuals (non-blocking) | plan_scopes mapping quality; scope pagination; Codegraph complement |
| Close owner | P44-S05-02 |
| Verify | Cite VERIFY-NOTES + evidence dir; §10 V1–V7 ticked |

## Handoff checklist

- [x] S00 research artifact committed
- [x] S01 design locks APPROVE
- [x] S02–S04 implement reviews APPROVE
- [x] Laws 6–7 / M-001 / Law 19 VERIFY blocks PASS
- [x] `trace/graph.json` exported (schema changed S02/S03)
- [x] `docs/TODO.md` + `AGENTS.md` orchestrator updated on close
```

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

### TODO.md / AGENTS.md close checklist (on APPROVE)

**If `no successor` (default):**

1. **`docs/TODO.md` orchestrator paste** — replace Active phase block with:

```text
Phase 00–44 complete — do not re-run closed rows.
@docs/TODO.md

- Active phase: **none** (idle — await human promotion)
- Next runnable: **none**
- Follow docs/rules/agent-loop-protocol.md
```

2. **`docs/TODO.md` phase boards table** — Phase 44 status `done`; Next `—`; parenthetical: P44 closed at S05-02; scoped graph linking shipped.

3. **`docs/TODO/phase-44.md`** — header/status note: Phase 44 **complete** at P44-S05-02; all rows `done`.

4. **`AGENTS.md` Current focus** — replace Phase 44 active block with:

```text
**Phase 44 complete** (2026-08-23) — closed at `P44-S05-02`; scoped graph linking delivered (thin scopes, MVP rels, CLI infer, GUI clustering @ 500). Successor: **no successor** (idle). Board [`docs/TODO/phase-44.md`](docs/TODO/phase-44.md).
```

5. **`AGENTS.md` Orchestrator snippet** — update paste to Phase 00–44 complete / idle (mirror TODO.md).

6. **Portable graph note** — confirm CONTRIBUTING path cited in VERIFY-NOTES; do not delete `trace/graph.json`.

**If human promotes Phase 45:** point Active phase / next runnable at scaffolded phase planner; update AGENTS Current focus to Phase 45; never TBD.

### REVIEW-NOTES.md template (required)

Write `scopes/scope-05-verify/REVIEW-NOTES.md`:

```markdown
# REVIEW-NOTES — P44-S05-02

**Date:** …
**Verdict:** APPROVE | REJECT (spawn)
**Confidence:** high | medium | low
**Successor:** no successor  (or Phase NN / PNN-00 | pending repair spawn)

## Spot-check
| Check | Result |
| VERIFY-NOTES overall | |
| Evidence dir | |
| V1 auth smoke | LIVE \| WAIVED |
| V2 MVP rels | |
| V3 infer CLI-only + explicit wins | |
| V4 caps 500/150/5000 | |
| V5 Law 19 rg | |
| V6 M-001 | |
| V7 portable graph | |
| §10 ticks | |

## Findings
…

## DR-HANDOFF
CLOSED | remains OPEN

## Orchestrator
TODO.md + AGENTS.md updated: yes | no
```

## Role work

1. Fresh-session re-verify S05-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES.md` in this folder (findings + confidence + successor pick).
3. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick scope checklist; set successor (never TBD).
4. Update `docs/TODO.md` + `AGENTS.md` per close checklist above.
5. Ensure Phase 44 board all rows `done` (or repair spawn pending — then do not close).

## Todo updates

Status + notes on **P44-S05-02** only; may spawn repair rows below this row if needed.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES.md`
- [ ] `DR-HANDOFF.md` CLOSED with successor **not** TBD
- [ ] `docs/TODO.md` + `AGENTS.md` updated (idle / `no successor` default)
- [ ] All P44 board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence

## Minimal todos

- [ ] Spot-check VERIFY-NOTES + V1–V7 floor
- [ ] Write REVIEW-NOTES.md
- [ ] Close DR-HANDOFF or spawn repair
- [ ] Update TODO.md + AGENTS.md
- [ ] Mark P44-S05-02 `done` / `failed` / `blocked`

## Next

Phase complete when this row `done` + DR-HANDOFF CLOSED (**no successor** default).
