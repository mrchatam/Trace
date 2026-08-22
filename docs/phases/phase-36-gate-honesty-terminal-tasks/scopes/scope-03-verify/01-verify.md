# P36-S03-01 — VERIFY

## Metadata
- id: P36-S03-01
- todo_ids: [P36-S03-01]
- role: implementer
- skills: [test-driven-development, debugging-and-error-recovery]
- mcps: [user-trace]
- verification: manual + automated
- hooks: []

## Objective

Run VERIFY blocks; author `VERIFY-NOTES.md` + evidence. Prove **greenfield agent-complete path** (MCP/bootstrap → edit gate passes), **feet-seller terminal gate honesty** (pre-bootstrap), **active-work PLAN enforcement preserved**, and **live recovery** via `trace plan bootstrap --goal`. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P36-S03-02**. **No product code.** Do **not** start S03-02 or invent a successor phase.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor locks (FINAL — S03-00)
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — §2.1–§2.8; §3 agent workflow; §6 acceptance tests
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S03-02
- [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md) — baseline repro (pre-S02 JSON)
- S02 board Notes (P36-S02-01/02) — implementation + review **PASS** (high confidence)
- Code anchors: `internal/loop/gate.go` (`goal_plan_gap_terminal_advisory`), `internal/mcp/tools_plan.go`, `internal/planner/bootstrap.go`, `web/src/components/GateStrip.tsx`, `web/src/screens/TaskDetail.tsx`

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF, decide successor, or change product bodies.

## Locked defaults (FINAL — S03-00)

| Item | Value |
|------|-------|
| Precondition | P36-00 … P36-S02-02 all `done`; S02 review **PASS** (high confidence) |
| Product / Go / TS changes | **Forbidden** (evidence + notes only). Failures → spawn remediation or leave FAIL for S03-02 |
| Trace binary | Build once: `go build -o /tmp/trace ./cmd/trace` from `/home/ali/Desktop/Trace` |
| Evidence dir (primary) | `experiments/runs/YYYY-MM-DD-p36-s03-01-verify/evidence/` |
| Pinned summary (optional) | `docs/verification/phase-36-gate-honesty/` — copy key JSON/screenshots if useful for DR-HANDOFF |
| Notes artifact | `scopes/scope-03-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S03-02 closes |
| Successor | **Out of scope** — S03-02 only (lean default **no successor**) |
| Dogfood path | `/home/ali/Desktop/feet seller telegram app` |
| Goal ID | `353b12a4-57dd-4d68-8379-b2024e064733` |
| Step 1 UUID | `33247e2d-aa10-4b25-b194-4b7afb5a6359` |
| Loop 112 UUID | `99d8fb92-65ac-462c-82c4-21bcf198c09e` |
| GUI launch | `trace gui -C "/home/ali/Desktop/feet seller telegram app"` (or `/tmp/trace gui -C "…"`) |

### Mutation policy (locked order)

| Blocks | Feet-seller writes |
|--------|-------------------|
| **0–5, 7** | **Read-only** on dogfood fixture |
| **6 only** | **Allowed** — `trace plan bootstrap --goal 353b12a4-…` mutates progressive planner state (PLAN §2.8). Run blocks **2–4 before block 6**. |

### DESIGN-LOCKS + PLAN acceptance map (must tick in VERIFY-NOTES)

| Lock / case | Acceptance |
|-------------|------------|
| Must-fix (product) | Agents can satisfy `PlanExists` via MCP/CLI bootstrap without undocumented-only CLI steps |
| Must-fix (honesty) | Terminal DONE tasks: `allowed: true` + `goal_plan_gap_terminal_advisory` — not actionable `plan_missing` block |
| Must preserve | Active non-terminal without plan → `plan_missing`, `allowed: false` |
| MCP | 16 tools including `trace_plan` (S02 test lock) |
| Recovery | Post-bootstrap on feet-seller: `plan_exists: true`, edit gate passes |
| Out of scope | PlanExists bridge; global PLAN weaken; feet-seller history delete; hosted SaaS |

### Baseline vs post-S02 (feet-seller — cite in VERIFY-NOTES)

**Pre-S02 (INVESTIGATION.md — must NOT regress to this on terminal DONE):**

```json
{"schema_version":"trace.loop.gate.v1","task_id":"33247e2d-aa10-4b25-b194-4b7afb5a6359","for":"done","allowed":false,"recommended_phase":"PLAN","reason_code":"plan_missing","violations":[{"code":"premature_implementation","for":"done","message":"done blocked: recommended phase PLAN (plan_missing)","recommended_phase":"PLAN","reason_code":"plan_missing"}]}
```

**Post-S02 terminal honesty (expected blocks 2–3):**

```json
{
  "schema_version": "trace.loop.gate.v1",
  "task_id": "33247e2d-aa10-4b25-b194-4b7afb5a6359",
  "for": "done",
  "allowed": true,
  "violations": [{
    "code": "goal_plan_gap_advisory",
    "for": "done",
    "reason_code": "goal_plan_gap_terminal_advisory",
    "message": "goal 353b12a4-57dd-4d68-8379-b2024e064733 lacks progressive plan (work already terminal); run trace plan bootstrap --goal 353b12a4-57dd-4d68-8379-b2024e064733 or MCP trace_plan action=bootstrap"
  }]
}
```

**Shape asserts (blocks 2–3):**

| Field | Expected |
|-------|----------|
| `allowed` | `true` |
| `reason_code` (top-level) | **absent** or `goal_plan_gap_terminal_advisory` (envelope sets top-level only when `!allowed`) |
| `violations[0].reason_code` | `goal_plan_gap_terminal_advisory` |
| `violations[0].code` | `goal_plan_gap_advisory` |
| CLI exit code | **0** (allowed) — not gate-blocked exit |
| Must **not** contain | `"done blocked: recommended phase PLAN (plan_missing)"` on terminal DONE |

**Active-work block (block 5 — temp dir, not feet-seller):**

```json
{
  "allowed": false,
  "reason_code": "plan_missing",
  "recommended_phase": "PLAN",
  "violations": [{
    "code": "premature_implementation",
    "reason_code": "plan_missing",
    "recommended_phase": "PLAN"
  }]
}
```

### Fail vs residual (locked)

**Fail VERIFY for:**

- Block 0 test suite non-zero exit
- Blocks 2–3: feet-seller terminal gate `allowed: false` or `plan_missing` without terminal advisory semantics
- Block 4: GUI shows red **Gate blocked** / `banner--error` on DONE Step1 or Loop112 pre-bootstrap (warn OK)
- Block 5: active task without plan passes edit gate
- Block 6: bootstrap fails or post-bootstrap `current_scope_id` still null / edit gate still `plan_missing`
- Block 1: greenfield path cannot reach `plan_exists: true` + edit gate pass
- Dogfood mutated outside block 6
- Product code shipped in this row
- Global `plan_missing` weaken on active work
- VERIFY-NOTES missing or evidence dir absent after claimed PASS

**Do not fail VERIFY solely for residuals below** (record in VERIFY-NOTES — S02-02 fold-ins):

| Residual | Disposition |
|----------|-------------|
| Bootstrap help omits explicit human-refinement note (PLAN §2.2) | Accept — document in residuals |
| `trace loop status` `advisories[]` for goal-structure warning not wired | Accept — plan show + MCP show field sufficient |
| No `WarnIfTraceDirWithoutConfig` unit test | Accept — init/install stderr nudge shipped |
| PlanExists bridge (§2.4) deferred | Accept — not Phase 36 scope |
| HTTP POST plan mutation routes deferred | Accept — MCP primary agent surface |
| MCP `trace_loop action=gate` absent | Accept — CLI/hook/GUI HTTP sufficient |

## Locked verify command floor

Run from Trace repo root unless noted. Tee outputs into evidence dir. Use `date +%Y-%m-%d` for the run folder name.

```bash
cd /home/ali/Desktop/Trace
go build -o /tmp/trace ./cmd/trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p36-s03-01-verify/evidence"
mkdir -p "$EVID"
PIN="docs/verification/phase-36-gate-honesty"
mkdir -p "$PIN"
{
  echo "verify_id=P36-S03-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "trace_binary=/tmp/trace"
  echo "precondition=P36-S02-02 PASS high confidence"
  echo "fixture=/home/ali/Desktop/feet seller telegram app"
  echo "goal_id=353b12a4-57dd-4d68-8379-b2024e064733"
  echo "step1=33247e2d-aa10-4b25-b194-4b7afb5a6359"
  echo "loop112=99d8fb92-65ac-462c-82c4-21bcf198c09e"
} > "$EVID/00-run-metadata.txt"
test -d "/home/ali/Desktop/feet seller telegram app/.trace"
```

**Pass:** `$EVID` exists; metadata cites S02-02 PASS; fixture `.trace/` present.

---

### Block 0 — S02 scoped test suite (re-check)

```bash
cd /home/ali/Desktop/Trace
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... ./cmd/trace/... \
  ./internal/config/... ./internal/install/... ./internal/domain/... -count=1 \
  2>&1 | tee "experiments/runs/$(date +%Y-%m-%d)-p36-s03-01-verify/evidence/00-s02-scoped-tests.txt"
```

**Pass:** exit 0. Cite acceptance tests green in VERIFY-NOTES:

- `TestGreenfield_MCPPlanBootstrap_EditGatePasses`
- `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap`
- `TestActiveWork_PlanMissingStillBlocksEdit`
- Recommended: `TestEvaluateGate_Done_TerminalPlanGapAdvisory`, `TestPlanBootstrap_Idempotent`, `TestGoalStructureWarning_OverThresholdNoPlan`, `TestRegisteredToolNames_IncludesTracePlan`

Optional focused re-run:

```bash
go test ./internal/planner/... ./internal/loop/... ./internal/mcp/... -count=1 \
  -run 'Greenfield|FeetSeller|ActiveWork|TerminalPlanGap|Bootstrap_Idempotent|GoalStructure|TracePlan|ToolNames' \
  2>&1 | tee "$EVID/00b-acceptance-subset.txt"
```

---

### Block 1 — Greenfield agent path (temp project → bootstrap → edit gate)

Uses CLI chain (MCP path covered by Block 0 `TestGreenfield_MCPPlanBootstrap_EditGatePasses`). Evidence: `$EVID/01-greenfield/`.

```bash
GF=$(mktemp -d)
mkdir -p "$EVID/01-greenfield"
echo "greenfield_root=$GF" > "$EVID/01-greenfield/root.txt"

/tmp/trace -C "$GF" init 2>&1 | tee "$EVID/01-greenfield/init.txt"

GOAL_JSON=$(/tmp/trace -C "$GF" add goal --title "Verify greenfield goal")
echo "$GOAL_JSON" | tee "$EVID/01-greenfield/add-goal.json"
GOAL_ID=$(echo "$GOAL_JSON" | jq -r '.id')

TASK_JSON=$(/tmp/trace -C "$GF" add task --title "Verify greenfield task" --goal-id "$GOAL_ID")
echo "$TASK_JSON" | tee "$EVID/01-greenfield/add-task.json"
TASK_ID=$(echo "$TASK_JSON" | jq -r '.id')

COARSE_JSON=$(/tmp/trace -C "$GF" plan create-coarse --goal "$GOAL_ID" --phase "Phase 1" --scope "Scope 1")
echo "$COARSE_JSON" | tee "$EVID/01-greenfield/create-coarse.json"
SCOPE_ID=$(echo "$COARSE_JSON" | jq -r '.phases[0].scopes[0].id')

/tmp/trace -C "$GF" plan set-current --goal "$GOAL_ID" --scope "$SCOPE_ID" \
  2>&1 | tee "$EVID/01-greenfield/set-current.txt"

/tmp/trace -C "$GF" plan deep --scope "$SCOPE_ID" --exit "packet ready" --work "implement" \
  2>&1 | tee "$EVID/01-greenfield/deep.txt"

/tmp/trace -C "$GF" loop status --task "$TASK_ID" \
  2>&1 | tee "$EVID/01-greenfield/loop-status.json"

/tmp/trace -C "$GF" loop gate --task "$TASK_ID" --for edit \
  2>&1 | tee "$EVID/01-greenfield/edit-gate.json"
echo "edit_gate_exit=$?" >> "$EVID/01-greenfield/edit-gate.json"
```

**Pass asserts (jq or manual):**

| Check | Expected |
|-------|----------|
| `loop status` → `policy_inputs.plan_exists` | `true` |
| `loop gate --for edit` → `allowed` | `true` |
| `reason_code` in violations | **≠** `plan_missing` |
| CLI exit (edit gate) | **0** |

Copy `$EVID/01-greenfield/edit-gate.json` → `$PIN/greenfield-edit-gate.json` (optional pin).

---

### Block 2 — Feet-seller Step 1 DONE gate (pre-bootstrap, read-only)

**Run before block 6.**

```bash
FEET="/home/ali/Desktop/feet seller telegram app"
STEP1=33247e2d-aa10-4b25-b194-4b7afb5a6359

/tmp/trace -C "$FEET" loop gate --task "$STEP1" --for done \
  | tee "$EVID/02-feet-step1-done-gate-pre-bootstrap.json"
echo "exit_code=$?" >> "$EVID/02-feet-step1-done-gate-pre-bootstrap.json"

# Sanity: progressive planner still empty pre-bootstrap
/tmp/trace -C "$FEET" plan show --goal 353b12a4-57dd-4d68-8379-b2024e064733 \
  2>"$EVID/02-feet-plan-show-stderr.txt" \
  | tee "$EVID/02-feet-plan-show-pre-bootstrap.json"
jq '{current_scope_id, current_deep_plan: (.current_deep_plan != null)}' \
  "$EVID/02-feet-plan-show-pre-bootstrap.json" \
  | tee "$EVID/02-feet-planner-empty-check.json"
```

**Pass:** `allowed: true`; violation `goal_plan_gap_terminal_advisory`; exit 0; `current_scope_id` null pre-bootstrap.

Copy → `$PIN/feet-step1-done-gate-pre-bootstrap.json` (optional).

---

### Block 3 — Feet-seller Loop 112 DONE gate (pre-bootstrap, read-only)

```bash
LOOP112=99d8fb92-65ac-462c-82c4-21bcf198c09e

/tmp/trace -C "$FEET" loop gate --task "$LOOP112" --for done \
  | tee "$EVID/03-feet-loop112-done-gate-pre-bootstrap.json"
echo "exit_code=$?" >> "$EVID/03-feet-loop112-done-gate-pre-bootstrap.json"
```

**Pass:** Same shape as block 2; only `task_id` differs. Both Step1 and Loop112 must **not** show pre-S02 `plan_missing` block JSON.

---

### Block 4 — GUI TaskDetail on DONE task (pre-bootstrap, read-only)

Launch GUI against feet-seller (stop after evidence captured):

```bash
/tmp/trace gui -C "$FEET" --no-open
# or: go run ./cmd/trace gui -C "$FEET" --no-open
```

**Checklist (record in `$EVID/04-gui-taskdetail-notes.txt` + optional PNGs under `$EVID/04-gui/`):**

1. Navigate to TaskDetail for **Step 1** (`33247e2d-…`) — e.g. `/task?task_id=33247e2d-aa10-4b25-b194-4b7afb5a6359` (exact route per SPA).
2. **DONE gate** panel: GateStrip shows **warn** tone (`Gate warnings` / `banner--warn`), **not** red `Gate blocked` / `banner--error`.
3. Advisory copy present: mentions `trace plan bootstrap --goal` or MCP `trace_plan` for terminal work (TaskDetail.tsx advisory paragraph).
4. Repeat spot-check on **Loop 112** (`99d8fb92-…`).

**Pass:** No misleading red error strip on finished tasks; warn/advisory chrome only. API gate JSON must still match blocks 2–3 (Law 19).

Optional API cross-check while GUI running:

```bash
# If serve URL known:
# curl -sS "$BASE/v1/loop/gate?task_id=$STEP1&for=done" | jq .
```

---

### Block 5 — Active non-terminal task: PLAN gate still enforced (temp dir)

```bash
ACT=$(mktemp -d)
mkdir -p "$EVID/05-active-work"
echo "active_root=$ACT" > "$EVID/05-active-work/root.txt"

/tmp/trace -C "$ACT" init
G=$(/tmp/trace -C "$ACT" add goal --title "Active verify goal" | jq -r '.id')
T=$(/tmp/trace -C "$ACT" add task --title "Active verify task" --goal-id "$G" | jq -r '.id')
/tmp/trace -C "$ACT" transition --task "$T" --to IN_PROGRESS --reason "verify start" \
  2>&1 | tee "$EVID/05-active-work/transition.txt"

/tmp/trace -C "$ACT" loop gate --task "$T" --for edit \
  | tee "$EVID/05-active-work/edit-gate.json"
echo "exit_code=$?" >> "$EVID/05-active-work/edit-gate.json"
```

**Pass:**

| Field | Expected |
|-------|----------|
| `allowed` | `false` |
| `reason_code` | `plan_missing` |
| `recommended_phase` | `PLAN` |
| CLI exit | gate-blocked (non-zero per `exitGateBlocked`) |

Copy → `$PIN/active-work-edit-gate-blocked.json` (optional).

---

### Block 6 — Feet-seller recovery: bootstrap + PlanExists (mutates fixture)

**Human-approved mutation per PLAN §2.8.** Run only after blocks 2–4 captured pre-bootstrap evidence.

```bash
GOAL=353b12a4-57dd-4d68-8379-b2024e064733
STEP1=33247e2d-aa10-4b25-b194-4b7afb5a6359

# Pre snapshot (optional)
cp "$EVID/02-feet-plan-show-pre-bootstrap.json" "$EVID/06-pre-bootstrap-plan-show.json"

/tmp/trace -C "$FEET" plan bootstrap --goal "$GOAL" \
  2>&1 | tee "$EVID/06-bootstrap-stderr.txt" | tee "$EVID/06-bootstrap-stdout.json"

/tmp/trace -C "$FEET" plan show --goal "$GOAL" \
  2>"$EVID/06-plan-show-stderr.txt" \
  | tee "$EVID/06-post-bootstrap-plan-show.json"

jq '{current_scope_id, has_deep_plan: (.current_deep_plan != null), phase_count: (.phases | length)}' \
  "$EVID/06-post-bootstrap-plan-show.json" \
  | tee "$EVID/06-plan-exists-check.json"

/tmp/trace -C "$FEET" loop status --task "$STEP1" \
  | tee "$EVID/06-post-bootstrap-loop-status.json"

/tmp/trace -C "$FEET" loop gate --task "$STEP1" --for edit \
  | tee "$EVID/06-post-bootstrap-edit-gate.json"
echo "exit_code=$?" >> "$EVID/06-post-bootstrap-edit-gate.json"

# Idempotency re-run
/tmp/trace -C "$FEET" plan bootstrap --goal "$GOAL" \
  2>&1 | tee "$EVID/06-bootstrap-idempotent.txt"
```

**Pass asserts:**

| Check | Expected |
|-------|----------|
| `plan show` → `current_scope_id` | non-null UUID |
| `current_deep_plan` | non-null |
| `loop status` → `policy_inputs.plan_exists` | `true` |
| `loop gate --for edit` → `allowed` | `true` (may need deliberation seed — if blocked on critique, note in VERIFY-NOTES; primary pass is `plan_exists: true`) |
| Bootstrap idempotent re-run | stderr note, no error |
| History preserved | 123 tasks, 11 plan-changes, 127 reviews still present (`tasks \| jq length`, `review list \| jq length`) |

Optional export pin for CONTRIBUTING:

```bash
/tmp/trace -C "$FEET" seed export -o "$PIN/feet-post-bootstrap-export-snippet.json"
# Full export optional — may omit reviews per default
```

**Note:** Post-bootstrap, terminal DONE gate may change (plan exists → no terminal advisory). Pre-bootstrap honesty is proven in blocks 2–4 only.

---

### Block 7 — Residuals + DR-HANDOFF successor prep (notes only)

In VERIFY-NOTES, list:

1. **S02-02 low findings** (bootstrap help refinement note; loop status advisories optional; enforce nudge unit test gap).
2. **Deferred by design:** PlanExists bridge (§2.4), HTTP POST plan routes, MCP `trace_loop gate`.
3. **Feet-seller post-bootstrap state:** progressive planner populated; goal structure warning may clear.
4. **Successor recommendation for S03-02:** default **`no successor`** unless VERIFY exposes blocking product gap.

Do **not** close DR-HANDOFF in this row.

---

### Block 8 — WRITE VERIFY-NOTES.md

Create `docs/phases/phase-36-gate-honesty-terminal-tasks/scopes/scope-03-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — Phase 36 / S03-01

**Date:** …
**Overall:** PASS | FAIL
**Git SHA:** …
**Trace binary:** /tmp/trace
**Evidence:** experiments/runs/…-p36-s03-01-verify/evidence/
**Pinned (optional):** docs/verification/phase-36-gate-honesty/

## Precondition cites
- S00 INVESTIGATION.md baseline repro
- S01 PLAN.md locked fix set
- S02 P36-S02-01/02 PASS (high confidence)

## Block results
| Block | Result | Evidence file |
|-------|--------|---------------|
| 0 S02 tests | | 00-s02-scoped-tests.txt |
| 1 greenfield | | 01-greenfield/ |
| 2 feet Step1 done | | 02-feet-step1-done-gate-pre-bootstrap.json |
| 3 feet Loop112 done | | 03-feet-loop112-done-gate-pre-bootstrap.json |
| 4 GUI TaskDetail | | 04-gui/ |
| 5 active PLAN block | | 05-active-work/ |
| 6 bootstrap recovery | | 06-post-bootstrap-*.json |
| 7 residuals | | (this section) |

## JSON shape spot-checks
| Case | allowed | violation reason_code |
|------|---------|----------------------|
| Step1 done pre-bootstrap | true | goal_plan_gap_terminal_advisory |
| Loop112 done pre-bootstrap | true | goal_plan_gap_terminal_advisory |
| Active edit no plan | false | plan_missing |
| Greenfield edit post-bootstrap | true | ≠ plan_missing |
| Post-bootstrap plan_exists | true | (status JSON) |

## GUI spot-check
| Task | Strip tone | Misleading red? |
|------|------------|-----------------|
| Step1 DONE | warn | no |
| Loop112 DONE | warn | no |

## Residuals (non-blocking)
- …

## DR-HANDOFF
Stays OPEN — S03-02 closes. Successor recommendation: no successor | Phase 37 …

## Next
P36-S03-02
```

## Exit criteria

- [ ] `VERIFY-NOTES.md` with command output / screenshot refs
- [ ] Evidence dir populated under `experiments/runs/…-p36-s03-01-verify/evidence/`
- [ ] Blocks 0–7 executed in order (block 6 after 2–4)
- [ ] Board Notes on **P36-S03-01** only
- [ ] DR-HANDOFF remains OPEN
- [ ] Next: **P36-S03-02**

## Next

`P36-S03-02`
