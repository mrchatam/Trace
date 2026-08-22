# P37-S03-01 — VERIFY

## Metadata
- id: P37-S03-01
- todo_ids: [P37-S03-01]
- role: implementer
- skills: [test-driven-development, debugging-and-error-recovery, browser-testing-with-devtools]
- mcps: [user-trace, cursor-ide-browser]
- verification: mixed
- hooks: []

## Objective

Run VERIFY blocks **0–5** locked by S03-00 from [PLAN.md](../scope-01-plan/PLAN.md) §6. Author `VERIFY-NOTES.md` + evidence under `experiments/runs/` (pin key artifacts to `docs/verification/phase-37-p36-residuals/` when human-gated). Prove Phase 36 regression subset still green, every S02 accept residual evidenced, feet-seller R9 refinement path documented, greenfield MCP path intact, re-defer registry updated, and R10 browser spot-check for Overview advisories + terminal surfaces. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P37-S03-02**. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Post-bootstrap critique path (R11)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor locks (FINAL — S03-00)
- [PLAN.md](../scope-01-plan/PLAN.md) — §2 accepts, §3 re-defer, §5 tests, §6 VERIFY mapping
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [RESIDUALS.md](../scope-00-triage/RESIDUALS.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S03-02
- Phase 36 [VERIFY-NOTES.md](../../../phase-36-gate-honesty-terminal-tasks/scopes/scope-03-verify/VERIFY-NOTES.md) — regression baseline + feet-seller IDs
- S02 board Notes (P37-S02-01/02) — implement + review **PASS** (high confidence)
- Code anchors: `internal/loop/apply.go` (`advisories[]`), `internal/planner/advisory.go`, `internal/mcp/tools_loop.go`, `internal/httpapi/handlers_p1.go`, `web/src/screens/Overview.tsx`, `web/src/screens/TaskDetail.tsx`

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF, decide successor, or change product bodies.

## Locked defaults (FINAL — S03-00)

| Item | Value |
|------|-------|
| Precondition | P37-00 … P37-S02-02 all `done`; S02 review **PASS** (high confidence) |
| Product / Go / TS changes | **Forbidden** (evidence + notes only). Failures → spawn remediation or leave FAIL for S03-02 |
| Trace binary | Build once: `go build -o /tmp/trace ./cmd/trace` from `/home/ali/Desktop/Trace` |
| Evidence dir (primary) | `experiments/runs/YYYY-MM-DD-p37-s03-01-verify/evidence/` |
| Pinned summary (optional) | `docs/verification/phase-37-p36-residuals/` — JSON captures, browser screenshots, R9 doc note |
| Notes artifact | `scopes/scope-03-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S03-02 closes |
| Successor | **Out of scope** — S03-02 only (Block 5 prep in VERIFY-NOTES) |
| Dogfood path | `/home/ali/Desktop/feet seller telegram app` |
| Goal ID | `353b12a4-57dd-4d68-8379-b2024e064733` |
| Step 1 UUID | `33247e2d-aa10-4b25-b194-4b7afb5a6359` |
| Loop 112 UUID | `99d8fb92-65ac-462c-82c4-21bcf198c09e` |
| GUI launch | `trace serve` + browser, or `/tmp/trace gui -C "/home/ali/Desktop/feet seller telegram app"` |

### Mutation policy (locked)

| Blocks | Feet-seller writes |
|--------|-------------------|
| **0–1, 3–5** | **Read-only** on dogfood fixture |
| **2 (R9 doc only)** | **Read-only** — document `create-coarse` / `deep` refinement path; optional quality notes; **no** task history rewrite |
| **R10 browser** | Read-only UI inspection |

**Note:** P36 VERIFY block 6 may have already bootstrapped feet-seller. Block 2 assumes post-bootstrap state; cite P36 evidence for pre-bootstrap terminal honesty if live state differs.

### S02 accept map (must tick in VERIFY-NOTES Block 1)

| ID | Evidence type | Primary proof |
|----|---------------|---------------|
| **R1** | Automated test | `TestLoopStatus_BootstrapRecommendedAdvisory` + `TestLoopStatus_BootstrapAdvisoryNeverSetsPlanExists` |
| **R2** | Automated test | `TestHTTPPlanBootstrap_CreatesPlannerRows` + OpenAPI path present |
| **R3** | Automated test | `TestMCPLoopGate_MatchesCLI` |
| **R4** | Automated test | `TestPlanHelp_MentionsRefinement` |
| **R5** | Automated test | `TestLoopStatus_IncludesGoalStructureAdvisory` |
| **R6** | Automated test | `TestWarnIfTraceDirWithoutConfig` |
| **R8** | API JSON + R10 browser | Overview `Plan advisories` banner; TaskDetail bootstrap paragraph unchanged |
| **R11** | Doc + Block 3 regression | `docs/rules/agent-loop-protocol.md` critique path; `TestGreenfield_MCPPlanBootstrap_EditGatePasses` |
| **R10** | Browser (Block 4) | Screenshot or pinned note under `docs/verification/phase-37-p36-residuals/` |

### Re-defer registry (Block 4 — must list in VERIFY-NOTES)

| ID | Item | Owner | Trigger |
|----|------|-------|---------|
| **R7** | Enforce default `warn` when `.trace/` without config | Human / product | Explicit decision to change `LoadEnforceMode` missing-config path |
| **R9** | Feet-seller deep refinement quality | Human dogfood | Post-bootstrap `create-coarse`/`deep` on fixture; quality spot-check optional |
| **R8-full** | Full plan screen / plan tree GUI | Phase 38 planner (or human) | Overview minimal surface insufficient for operators |

### Fail vs residual (locked)

**Fail VERIFY for:**

- Block 0: any P36 regression subset test fails
- Block 1: any S02 accept row lacks test/doc/API evidence
- Block 2: cannot document R9 refinement CLI path on fixture (commands fail unexpectedly)
- Block 3: `TestGreenfield_MCPPlanBootstrap_EditGatePasses` fails; MCP `trace_loop action=gate` envelope diverges from CLI
- Block 4: re-defer registry missing R7/R9/R8-full; R10 accepted but no browser evidence file when GUI reachable
- Block 5: successor prep absent from VERIFY-NOTES
- Product code shipped in this row
- R1 guard violated (`plan_exists` flipped true without planner rows)
- Active-work `plan_missing` weakened
- VERIFY-NOTES missing or evidence dir absent after claimed PASS

**Do not fail VERIFY solely for:**

| Topic | Disposition |
|-------|-------------|
| R9 planner quality subjective | Defer — document path + optional notes |
| R7 enforce default still off | Re-defer — R6 nudge test sufficient |
| R8-full plan tree GUI absent | Re-defer — Overview minimal shipped |
| Feet-seller `plan_uncritiqued` post-bootstrap | Expected — deliberation phase separate from P37 accepts |
| Vitest absent for Overview | Accept — R10 browser/manual checklist per PLAN |

## Locked verify command floor

Run from Trace repo root. Tee outputs into evidence dir. Use `date +%Y-%m-%d` for the run folder name.

```bash
cd /home/ali/Desktop/Trace
go build -o /tmp/trace ./cmd/trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p37-s03-01-verify/evidence"
mkdir -p "$EVID"
PIN="docs/verification/phase-37-p36-residuals"
mkdir -p "$PIN"
{
  echo "verify_id=P37-S03-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "trace_binary=/tmp/trace"
  echo "precondition=P37-S02-02 PASS high confidence"
  echo "fixture=/home/ali/Desktop/feet seller telegram app"
  echo "goal_id=353b12a4-57dd-4d68-8379-b2024e064733"
  echo "step1=33247e2d-aa10-4b25-b194-4b7afb5a6359"
  echo "loop112=99d8fb92-65ac-462c-82c4-21bcf198c09e"
} > "$EVID/00-run-metadata.txt"
test -d "/home/ali/Desktop/feet seller telegram app/.trace"
```

**Pass:** `$EVID` exists; metadata cites S02-02 PASS; fixture `.trace/` present.

---

### Block 0 — Phase 36 acceptance subset still green

Re-run the **7-test** P36 regression subset from PLAN §5 (must stay green after P37 S02).

```bash
cd /home/ali/Desktop/Trace
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./cmd/trace/... ./internal/config/... -count=1 \
  -run 'Greenfield_MCPPlanBootstrap_EditGatePasses|Legacy_FeetSellerExport_GateHonestyUntilBootstrap|ActiveWork_PlanMissingStillBlocksEdit|EvaluateGate_Done_TerminalPlanGapAdvisory|PlanBootstrap_Idempotent|GoalStructureWarning_OverThresholdNoPlan|RegisteredToolNames_IncludesTracePlan' \
  2>&1 | tee "$EVID/00-p36-regression-subset.txt"
```

**Pass:** exit 0. Cite each test green in VERIFY-NOTES:

- `TestGreenfield_MCPPlanBootstrap_EditGatePasses`
- `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap`
- `TestActiveWork_PlanMissingStillBlocksEdit`
- `TestEvaluateGate_Done_TerminalPlanGapAdvisory`
- `TestPlanBootstrap_Idempotent`
- `TestGoalStructureWarning_OverThresholdNoPlan`
- `TestRegisteredToolNames_IncludesTracePlan`

Optional broader S02 sweep:

```bash
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./internal/httpapi/... ./cmd/trace/... ./internal/config/... -count=1 \
  -run 'LoopStatus_.*Advisory|BootstrapAdvisoryNeverSetsPlanExists|MCPLoopGate|HTTPPlanBootstrap|PlanHelp_MentionsRefinement|WarnIfTraceDirWithoutConfig' \
  2>&1 | tee "$EVID/00b-s02-acceptance-subset.txt"
```

---

### Block 1 — Per accepted residual (test name or JSON capture)

Run S02 acceptance tests and capture representative JSON for status/gate surfaces.

```bash
# Focused acceptance re-run (feeds Block 1 table)
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./internal/httpapi/... ./cmd/trace/... ./internal/config/... -count=1 -v \
  -run 'LoopStatus_.*Advisory|BootstrapAdvisoryNeverSetsPlanExists|MCPLoopGate|HTTPPlanBootstrap|PlanHelp_MentionsRefinement|WarnIfTraceDirWithoutConfig' \
  2>&1 | tee "$EVID/01-s02-acceptance-v.txt"

# Live status advisories sample (feet-seller Step1 — read-only)
FEET="/home/ali/Desktop/feet seller telegram app"
STEP1=33247e2d-aa10-4b25-b194-4b7afb5a6359
GOAL=353b12a4-57dd-4d68-8379-b2024e064733

/tmp/trace -C "$FEET" loop status --task "$STEP1" \
  2>&1 | tee "$EVID/01-feet-loop-status.json"
jq '{schema_version, advisories, policy_inputs: .deliberation.policy_inputs.plan_exists}' \
  "$EVID/01-feet-loop-status.json" | tee "$EVID/01-feet-status-advisories-check.json"

/tmp/trace -C "$FEET" loop gate --task "$STEP1" --for done \
  2>&1 | tee "$EVID/01-feet-done-gate.json"

# R11 doc pointer (grep — no product change)
grep -n 'loop apply\|plan_changes\|TestGreenfield_MCPPlanBootstrap' \
  docs/rules/agent-loop-protocol.md \
  | tee "$EVID/01-r11-doc-cites.txt"
```

**Block 1 evidence table (fill in VERIFY-NOTES):**

| ID | Test / capture | Expected | Evidence file |
|----|----------------|----------|---------------|
| R1 | `TestLoopStatus_BootstrapRecommendedAdvisory` | `bootstrap_recommended` in `advisories[]` | `01-s02-acceptance-v.txt` |
| R1 guard | `TestLoopStatus_BootstrapAdvisoryNeverSetsPlanExists` | `plan_exists: false` | same |
| R2 | `TestHTTPPlanBootstrap_CreatesPlannerRows` | 200 + planner rows | same |
| R3 | `TestMCPLoopGate_MatchesCLI` | MCP gate ≡ CLI | same |
| R4 | `TestPlanHelp_MentionsRefinement` | create-coarse/deep mention | same |
| R5 | `TestLoopStatus_IncludesGoalStructureAdvisory` | `goal_structure_warning` | same |
| R6 | `TestWarnIfTraceDirWithoutConfig` | stderr nudge | same |
| R8 | `01-feet-loop-status.json` + R10 | `advisories[]` or gate violations consumed by Overview | Block 4 |
| R11 | `01-r11-doc-cites.txt` + Block 3 | agent-loop doc cites loop apply critique | Block 3 |

Copy `01-feet-loop-status.json` → `$PIN/feet-loop-status-advisories.json` (optional pin).

---

### Block 2 — Feet-seller spot-check (R8 Overview context + R9 defer doc path)

**Read-only** on fixture. Document post-bootstrap refinement path for R9 (defer — quality not gated).

```bash
FEET="/home/ali/Desktop/feet seller telegram app"
GOAL=353b12a4-57dd-4d68-8379-b2024e064733
mkdir -p "$EVID/02-feet-seller"

# Current planner state (post-P36 bootstrap likely)
/tmp/trace -C "$FEET" plan show --goal "$GOAL" \
  2>&1 | tee "$EVID/02-feet-seller/plan-show.json"
jq '{current_scope_id, has_deep: (.current_deep_plan != null), phase_count: (.phases | length)}' \
  "$EVID/02-feet-seller/plan-show.json" | tee "$EVID/02-feet-seller/planner-state.json"

# R9 documented refinement path (dry-run / help — do NOT rewrite task history)
{
  echo "# R9 refinement path (defer — document only)"
  echo "## Post-bootstrap state"
  jq -c '{current_scope_id, has_deep: (.current_deep_plan != null)}' "$EVID/02-feet-seller/plan-show.json"
  echo ""
  echo "## Operator commands (human dogfood)"
  echo "trace plan create-coarse --goal $GOAL --phase \"…\" --scope \"…\""
  echo "trace plan set-current --goal $GOAL --scope <scope_id>"
  echo "trace plan deep --scope <scope_id> --exit \"…\" --work \"…\""
  echo ""
  echo "## Expected post-bootstrap gate"
  echo "Edit gate may show plan_uncritiqued until deliberation critique seeded (P36 Block 1 partial — OK)."
} | tee "$EVID/02-feet-seller/r9-refinement-path.md"

# Optional: show help refinement note (R4 cross-check)
/tmp/trace plan help 2>&1 | tee "$EVID/02-feet-seller/plan-help.txt"
grep -i 'create-coarse\|deep\|minimal' "$EVID/02-feet-seller/plan-help.txt" \
  | tee "$EVID/02-feet-seller/r4-help-snippet.txt" || true

# HTTP loop status (R8 data contract — if serve running, else CLI above suffices)
# curl -sS "$BASE/v1/loop/status?task_id=$STEP1" | jq '{advisories, schema_version}'
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| `r9-refinement-path.md` exists | Documents create-coarse/deep sequence |
| Planner state captured | `current_scope_id` noted (null or UUID — cite P36 if bootstrapped) |
| No task history mutation | Task/review counts unchanged vs P36 VERIFY block 6 if applicable |
| R8 context | Status or gate JSON shows advisories/violations Overview can consume |

Copy `$EVID/02-feet-seller/r9-refinement-path.md` → `$PIN/r9-refinement-path.md` (optional).

---

### Block 3 — Greenfield MCP path (R3 gate + R11 workflow)

Prove MCP agent path still canonical; `trace_loop action=gate` matches CLI.

```bash
mkdir -p "$EVID/03-greenfield-mcp"

# Block 0 canonical test (re-run or cite Block 0 log)
go test ./internal/mcp/... -count=1 -v -run TestGreenfield_MCPPlanBootstrap_EditGatePasses \
  2>&1 | tee "$EVID/03-greenfield-mcp/greenfield-mcp-test.txt"

go test ./internal/mcp/... -count=1 -v -run TestMCPLoopGate_MatchesCLI \
  2>&1 | tee "$EVID/03-greenfield-mcp/mcp-gate-test.txt"

# R11 workflow alignment
grep -A2 'Post-bootstrap critique path\|trace loop apply' docs/rules/agent-loop-protocol.md \
  | tee "$EVID/03-greenfield-mcp/r11-agent-loop-excerpt.txt"
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| `TestGreenfield_MCPPlanBootstrap_EditGatePasses` | PASS — bootstrap → critique seed → edit gate |
| `TestMCPLoopGate_MatchesCLI` | PASS — R3 MCP gate action |
| R11 doc | agent-loop-protocol cites `trace loop apply` + plan_changes; **no** critique-seed MCP tool |
| MCP tool count | Still 16 tools (`TestRegisteredToolNames_IncludesTracePlan` from Block 0) |

---

### Block 4 — Re-defer registry + R10 browser spot-check

#### 4a — Re-defer registry (VERIFY-NOTES section)

In VERIFY-NOTES, copy/update registry from PLAN §3:

| ID | Status | Notes |
|----|--------|-------|
| R7 | re-defer | `EnforceOff` default preserved; R6 test locks stderr nudge |
| R9 | re-defer | Block 2 doc path; quality human-owned |
| R8-full | re-defer | Overview minimal only; full plan tree → Phase 38+ |

List P36 residuals **consumed** by P37 (R1–R6, R8 partial, R11 doc, R2/R3 adapters).

#### 4b — R10 browser (Overview advisories + terminal surfaces)

Launch local GUI (prefer `trace serve` on loopback + embedded SPA):

```bash
# Terminal A — from Trace repo (adjust port if needed)
/tmp/trace serve -C "$FEET" --bind 127.0.0.1
# Note URL printed (e.g. http://127.0.0.1:PORT/)
```

Use **cursor-ide-browser** MCP:

1. Navigate to Overview — confirm **Plan advisories** banner when `advisories[]` or gate violations present (`Overview.tsx` warn banner).
2. Pick active task or Step 1 DONE — TaskDetail shows bootstrap/advisory copy (`TaskDetail.tsx:205–211`).
3. GateStrip on DONE tasks: **warn** tone, not misleading red error (terminal honesty preserved from P36).
4. Capture screenshots → `$EVID/04-browser/` and optionally `$PIN/`.

Record checklist in `$EVID/04-browser/r10-spot-check-notes.txt`:

```text
# R10 browser spot-check
- Overview Plan advisories banner: (yes/no + task context)
- TaskDetail bootstrap paragraph: (yes/no)
- GateStrip DONE tone: warn | error (must be warn for terminal advisory)
- API cross-check: loop status advisories[] matches UI (Law 19)
```

**Pass:** Evidence file exists under `$EVID/04-browser/` or `$PIN/`. If GUI unreachable, mark R10 **blocked** with reason — do **not** self-claim PASS without file.

**Human-gated rule:** Pin at least one artifact to `docs/verification/phase-37-p36-residuals/` before marking R10 PASS.

---

### Block 5 — Successor table prep for DR-HANDOFF (notes only)

In VERIFY-NOTES, include successor recommendation for S03-02 (do **not** close DR-HANDOFF):

| Outcome | Recommended successor |
|---------|----------------------|
| VERIFY blocks 0–4 green; residuals only R7/R9/R8-full | **Phase 38 scaffold** — human promotes `P38-00` (investigation only) |
| Blocking product gap discovered | Spawn repair or thin follow-on — cite block FAIL |
| All accepts shipped; no investigation wanted | **`no successor`** — idle orchestrator |

Reference: Phase 38 already scaffolded at `docs/phases/phase-38-retrieval-context-peer-gaps/` + `docs/TODO/phase-38.md`. **Do not** implement Phase 38 in this row.

---

### Block 6 — WRITE VERIFY-NOTES.md

Create `docs/phases/phase-37-p36-residuals/scopes/scope-03-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — Phase 37 / S03-01

**Date:** …
**Overall:** PASS | FAIL | PARTIAL
**Git SHA:** …
**Trace binary:** /tmp/trace
**Evidence:** experiments/runs/…-p37-s03-01-verify/evidence/
**Pinned (optional):** docs/verification/phase-37-p36-residuals/

## Precondition cites
- S00 RESIDUALS.md accept set
- S01 PLAN.md §5–§6
- S02 P37-S02-01/02 PASS (high confidence)
- P36 VERIFY-NOTES regression baseline

## Block results
| Block | Result | Evidence file |
|-------|--------|---------------|
| 0 P36 regression subset | | 00-p36-regression-subset.txt |
| 1 per-residual accepts | | 01-s02-acceptance-v.txt, 01-feet-*.json |
| 2 feet-seller R9 doc | | 02-feet-seller/r9-refinement-path.md |
| 3 greenfield MCP | | 03-greenfield-mcp/ |
| 4 re-defer + R10 browser | | 04-browser/ |
| 5 successor prep | | (this section) |

## Per-residual accept map
| ID | Result | Evidence |
|----|--------|----------|
| R1 | | |
| … | | |

## Re-defer registry
| ID | Disposition |
|----|-------------|
| R7 | |
| R9 | |
| R8-full | |

## P36 residuals consumed
- …

## DR-HANDOFF
Stays OPEN — S03-02 closes. Successor recommendation: Phase 38 / no successor / repair

## Next
P37-S03-02
```

## Exit criteria

- [ ] `VERIFY-NOTES.md` with block table + per-residual map
- [ ] Evidence dir populated under `experiments/runs/…-p37-s03-01-verify/evidence/`
- [ ] Blocks 0–5 executed in order
- [ ] R10 browser evidence pinned if GUI reachable
- [ ] Board Notes on **P37-S03-01** only
- [ ] DR-HANDOFF remains OPEN
- [ ] Next: **P37-S03-02**

## Next

`P37-S03-02`
