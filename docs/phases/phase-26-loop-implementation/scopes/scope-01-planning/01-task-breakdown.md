# P26-S01-01 — Task breakdown implementer

## Metadata
- id: P26-S01-01
- todo_ids: [P26-S01-01]
- role: implementer
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Write `PLAN.md` from `AUDIT.md` with one section per downstream scope (S02, S03, S04). Each task row must be runnable by a single implementer session (S/M size). **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) (planner gate — locked by P26-S01-00)
- [AUDIT.md](../scope-00-loop-audit/AUDIT.md) (source of truth for file:line evidence)
- [Phase 26 README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) (VERIFY deliverables D1–D6)
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) §4 (human gate on promotion)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Input | `scopes/scope-00-loop-audit/AUDIT.md` |
| Output | `scopes/scope-01-planning/PLAN.md` |
| Product Go | **No** on this row |
| Schema | No new tables without version-bump plan in PLAN.md |
| Daemon/HTTP | Forbidden |
| Human gate | No silent background spawn — promotion only via `loop apply` `spawned_tasks[]` or explicit `trace add task` |
| Scope order | S02 → S03 → S04 (serial default; S04 may parallel S03 only after this PLAN exists) |
| Thresholds | List options from AUDIT; **do not pick final numbers** — S02/S03 scope planners lock at implement time |

## Preflight (mandatory before writing PLAN.md)

- [ ] Confirm `AUDIT.md` contains: Executive summary, INT-01/06, INT-02, INT-05, INT-09, Installer gap (P25-2), Delta from P24, Evidence appendix
- [ ] Re-read S02/S03/S04 `00-PLANNER.md` locked defaults — PLAN tasks must satisfy their deliverable tables
- [ ] Map each AUDIT gap row to ≥1 PLAN task (no orphan findings)

## PLAN.md structure (locked — use exactly)

```markdown
# Phase 26 implementation plan

## Overview
[1 paragraph: three gaps from AUDIT executive summary + scope mapping S02/S03/S04]

## Architecture decisions
| Decision | Options (from AUDIT) | Recommendation for S02/S03 planner | Rationale |
|----------|----------------------|-------------------------------------|-----------|
| … | … | … | … |

## Dependency graph
S02 (promotion) ──┐
                  ├──► S05 VERIFY (D1–D6)
S03 (saturation/reset/reason) ──┤
S04 (installer) ────────────────┘
[S04 may parallel S03 after S01; default serial S02→S03→S04]

## S02 — P25-A: discovery→task promotion (INT-01, INT-06)
[Task table — see schema below]

## S03 — P25-B: loop recalibration + deliberation reset (INT-02, INT-05, INT-09)
[Task table]

## S04 — Installer fix P25-2 (ParentOrchestratorRule wiring)
[Task table]

## Cross-scope risks
| Risk | Scopes | Mitigation |
|------|--------|------------|

## VERIFY mapping (D1–D6)
| Deliverable | PLAN task ID(s) | Verification command |
|-------------|-----------------|----------------------|

## Open questions (human gate)
[Only if AUDIT leaves ambiguity; else "None"]
```

### Per-task row schema (required columns)

| Column | Content |
|--------|---------|
| **ID** | `S02-T01`, `S03-T02`, … (unique, stable for board Notes) |
| **INT** | INT-01, INT-06, … |
| **Title** | Short verb phrase |
| **File(s)** | Repo paths from AUDIT (with line hints) |
| **Change** | What to add/modify (function/constant/field) |
| **Acceptance criteria** | 2–4 testable bullets |
| **Verification** | `go test` package or CLI repro from AUDIT evidence appendix |
| **Risk** | From AUDIT risk column + regression note |
| **Size** | XS / S / M |
| **Depends on** | Task IDs or "None" |

---

## S02 task seeds (INT-01 + INT-06)

**Problem (AUDIT):** Only `loop apply` `spawned_tasks[]` expands backlog; standalone `trace add discovery`, MCP `trace_add`, and `discovery-mentions-task` record gaps without new tasks. MCP kind list is flat. G1 Session B: 7 discoveries, 0 spawned tasks, roster stays at 5.

**Human gate (non-negotiable):** Promotion requires agent intent via `loop apply` envelope or explicit `trace add task` — no background auto-spawn from discovery import alone.

Emit PLAN tasks covering at minimum:

### S02-T01 — BLOCKING discovery promotion helper (INT-01)

| Field | Value |
|-------|-------|
| File(s) | `internal/domain/` (new or extend seed import), `internal/loop/apply.go:L479–497`, `internal/mcp/tools_write.go:L121–129` |
| Change | Domain helper: given BLOCKING discovery ID, produce task seed payload (title from discovery, link discovery→task). Callable from apply path when agent supplies promotion intent — **not** auto on every discovery import |
| Acceptance | Helper returns valid task UUID + link; idempotent on re-promote; preserves human gate (no spawn without apply/add) |
| Verification | Unit test in `internal/domain/` or `internal/loop/` |
| Risk | Duplicate tasks if idempotency weak; seed UUID collision |
| Size | M |

### S02-T02 — `loop apply` promotion path (INT-01)

| Field | Value |
|-------|-------|
| File(s) | `internal/loop/apply.go:L409–428`, `L479–497` |
| Change | Extend apply envelope: e.g. `promote_discoveries[]` or document+implement use of existing `spawned_tasks[]` with discovery-sourced entries; discoveries-only apply must not be the **only** path for BLOCKING gaps |
| Acceptance | Apply with promotion intent creates task row + discovery link; `spawned_tasks` in apply result lists new IDs; discoveries-only apply behavior unchanged unless paired with spawn |
| Verification | Integration test: BLOCKING discovery in store → apply with promotion → task count +1 |
| Risk | Triggers INT-02 saturation if spawn count still 0 (coordinate with S03-T01) |
| Size | M |
| Depends on | S02-T01 |

### S02-T03 — CLI discoverability (INT-01)

| Field | Value |
|-------|-------|
| File(s) | `cmd/trace/add.go:L68–83`, `L104–114` |
| Change | Surface promotion path in help text; optional `--from-discovery <id>` on `trace add task` or documented two-step flow in `loop status` hints |
| Acceptance | `trace add task --help` (or equivalent) mentions discovery promotion; manual flow creates linked task |
| Verification | CLI test or documented manual repro in PLAN |
| Risk | Low — docs/UX only if no new flag |
| Size | S |

### S02-T04 — MCP `trace_add` description reorder (INT-06)

| Field | Value |
|-------|-------|
| File(s) | `internal/mcp/server.go:L77–78`, `internal/mcp/tools_write.go` |
| Change | Reorder kind list: discovery first, then explicit nudge "after BLOCKING discovery, call trace add task or loop apply with spawned_tasks before edits" |
| Acceptance | MCP tool schema/description contains promotion ordering text; parity with `GapPassPrompt` intent |
| Verification | `grep`/test on registered tool description string |
| Risk | Agent confusion if text too long |
| Size | XS |

### S02-T05 — `loop next` promotion candidates (INT-01)

| Field | Value |
|-------|-------|
| File(s) | `internal/loop/next.go:L190–195`, packet builder |
| Change | Add optional `promotion_candidates[]` (BLOCKING discoveries without linked executable task) to next packet |
| Acceptance | Next packet lists open BLOCKING discoveries eligible for promotion; empty when none |
| Verification | Unit test with fixture discoveries |
| Risk | Packet bloat; false positives if severity filter wrong |
| Size | M |
| Depends on | S02-T01 |

### S02-T06 — End-to-end promotion test (INT-01 + INT-06)

| Field | Value |
|-------|-------|
| File(s) | `internal/loop/*_test.go`, optionally `internal/mcp/*_test.go` |
| Change | Test: create BLOCKING discovery → promote via apply or add → task appears in `ListTasksByGoalID` |
| Acceptance | FM-10 repro closed: discovery count > 0 implies path to task count increase when agent promotes |
| Verification | `go test ./internal/loop/... ./internal/mcp/...` |
| Risk | Flaky if UUID fixtures drift |
| Size | S |
| Depends on | S02-T02 |

**S02 threshold decision point (for S02-00 planner — list in PLAN Architecture decisions, do not implement here):**

| Knob | AUDIT options | Planner picks at S02-00 |
|------|---------------|-------------------------|
| Promotion trigger | Manual `spawned_tasks[]` only vs `promote_discoveries[]` alias vs CLI `--from-discovery` | One primary path + documented secondary |
| BLOCKING filter | All discoveries vs BLOCKING severity only | Default BLOCKING-only per INTERVENTION-MATRIX |
| Idempotency key | Discovery UUID as task seed vs new UUID + link | Must prevent duplicate tasks on re-apply |

---

## S03 task seeds (INT-02 + INT-05 + INT-09)

**Problem (AUDIT):** First empty apply (zero plan_changes + zero spawned_tasks) sets sticky STOP; discoveries-only apply counts as empty spawn → saturates. No reset API. Gate/status report `hop_budget_exceeded` while export persists `p19_saturated`.

Emit PLAN tasks covering at minimum:

### S03-T01 — P19 saturation recalibration (INT-02)

| Field | Value |
|-------|-------|
| File(s) | `internal/loop/policy.go:L103–109`, `internal/loop/apply.go:L499–500`, `internal/deliberation/select.go:L10–11` |
| Change | Replace immediate zero-write saturation with **consecutive empty apply counter** (named constant, default ≥2 per S03-00 lock); optional: discoveries-only apply non-saturating (AUDIT option) |
| Acceptance | First empty apply on greenfield does **not** set `Stopped=true`; second consecutive empty apply saturates; unit tests updated (`select_test.go:L64–67`, `L88–93`) |
| Verification | `go test ./internal/deliberation/... ./internal/loop/...`; G1 fixture: hop_count=1 no longer sticky STOP |
| Risk | Regression if counter not persisted; interaction with hop budget |
| Size | M |

**Threshold decision point (S03-00 planner locks):**

| Knob | AUDIT options | Default in S03-00 lock |
|------|---------------|------------------------|
| Empty apply threshold | 1 (current), 2+, N consecutive | **≥2 consecutive** |
| Discoveries-only apply | Counts as empty spawn (current) vs non-saturating | PLAN must pick one; recommend non-saturating when discoveries imported without spawn |
| Hop budget | 12 fixed vs profile-specific | Keep 12 unless test proves gap-pass needs raise |
| First-apply greenfield exempt | Yes vs consecutive counter only | Prefer consecutive counter (cleaner persistence) |

### S03-T02 — Persist consecutive-empty counter (INT-02)

| Field | Value |
|-------|-------|
| File(s) | `internal/store/deliberation.go:L28–60`, `internal/domain/deliberation.go`, deliberation state struct |
| Change | Add field for consecutive empty applies (or derive from loop step history if no schema bump — if schema bump, document version bump in PLAN) |
| Acceptance | Counter survives apply/status round-trip; resets on non-empty apply |
| Verification | Store/domain unit tests |
| Risk | Schema migration if new column — requires version bump per phase lock |
| Size | M |
| Depends on | S03-T01 design choice |

### S03-T03 — Deliberation reset API (INT-05)

| Field | Value |
|-------|-------|
| File(s) | `internal/domain/deliberation.go:L14–72`, `internal/store/deliberation.go`, `cmd/trace/loop.go` |
| Change | New domain fn e.g. `ResetDeliberationState(taskID)` clearing `Stopped`, `StopReason`, `HopCount`, consecutive-empty counter; CLI `trace loop reset --task <id>` |
| Acceptance | After reset on G1 verify task: `stopped=false`, `hop_count=0`, phase returns toward EXECUTE (subject to `plan_critiqued` gate per AUDIT INT-05) |
| Verification | `go test ./internal/domain/...`; manual CLI on G1 fixture |
| Risk | Bypasses human STOP intent if ungated — document when reset is allowed |
| Size | M |

### S03-T04 — Reset + gap-pass integration (INT-05)

| Field | Value |
|-------|-------|
| File(s) | `internal/loop/gate.go:L188–192`, `internal/loop/apply.go` |
| Change | After reset, empty apply uses new saturation threshold (not immediate re-STOP); status packet hints reset when saturated + discoveries closed |
| Acceptance | Gap-pass flow: saturate → reset → discoveries-only apply does not immediately re-STOP (per threshold lock) |
| Verification | Integration test sequence apply→saturate→reset→apply |
| Risk | Infinite loop if threshold too high |
| Size | M |
| Depends on | S03-T01, S03-T03 |

### S03-T05 — Unified STOP reason (INT-09)

| Field | Value |
|-------|-------|
| File(s) | `internal/deliberation/select.go:L7–8`, `L10–11`, `internal/loop/gate.go:L159`, `L191–192`, `internal/loop/deliberation_packet.go:L188–202`, `internal/loop/apply.go:L652–671` |
| Change | When `Stopped=true`, gate/status/`why_selected` prefer **persisted** `StopReason` over remapped `hop_budget_exceeded`; add `stop_reason` field to status JSON alongside `why_selected` |
| Acceptance | G1 repro closed: gate `reason_code` matches export `stop_reason` (`p19_saturated` case); status includes both fields when they differ during transition |
| Verification | Unit tests for SelectNext stopped branch; gate JSON test; AUDIT jq repro passes |
| Risk | Breaking change for agents expecting hop_budget_exceeded |
| Size | M |

### S03-T06 — Export/gate alignment test (INT-09)

| Field | Value |
|-------|-------|
| File(s) | `internal/domain/seed_export.go:L710–719`, `internal/loop/gate_test.go` or new test |
| Change | Assert export `stop_reason` equals gate `reason_code` after saturation STOP |
| Acceptance | `seed export` deliberation_states match live gate for same task |
| Verification | `go test ./internal/loop/... ./internal/domain/...` |
| Risk | Low |
| Size | S |
| Depends on | S03-T05 |

**S03 threshold decision point (reason UX — S03-00 planner locks):**

| Knob | AUDIT options | Recommendation for PLAN |
|------|---------------|-------------------------|
| Canonical reason | Persisted `StopReason` vs fresh SelectNext vs dual report + hint | **Persisted wins** for gate/status; optional hint field for recovery |
| hop_budget_exceeded | Only when hop count ≥ budget without prior STOP | Re-order SelectNext branches |

---

## S04 task seeds (P25-2 installer)

**Problem (AUDIT):** `ParentOrchestratorRule` defined at `internal/install/enforcement.go:L20–35` but not concatenated into `AgentsEnforcementBlock`, `cursorRulesMDCContent`, or `claudeFallbackRulesContent` (only `GapPassPrompt` wired).

### S04-T01 — Wire ParentOrchestratorRule into install bodies (P25-2)

| Field | Value |
|-------|-------|
| File(s) | `internal/install/enforcement.go` (`AgentsEnforcementBlock` ~L65, `cursorRulesMDCContent` ~L83, `claudeFallbackRulesContent` ~L97) |
| Change | Append `"\n" + ParentOrchestratorRule + "\n"` after `GapPassPrompt` in all three bodies (order: gate body → GapPassPrompt → ParentOrchestratorRule per AUDIT) |
| Acceptance | `rg ParentOrchestratorRule internal/install/` hits usage in all three funcs; fresh `trace install cursor --write` output contains "Parent orchestrator" |
| Verification | `go test ./internal/install/...` |
| Risk | Low — text-only; hook alignment (INT-04) out of scope beyond text |
| Size | XS |

### S04-T02 — Enforcement test assertion (P25-2)

| Field | Value |
|-------|-------|
| File(s) | `internal/install/enforcement_test.go` |
| Change | Assert `"Parent orchestrator"` (or rule substring) present in `cursorRulesMDCContent()` and Claude fallback output |
| Acceptance | Test fails on pre-fix tree, passes after S04-T01 |
| Verification | `go test ./internal/install/...` |
| Risk | Low |
| Size | XS |
| Depends on | S04-T01 |

**S04 out of scope (document in PLAN risks):** `CursorLoopGateHookScript` failClosed vs TRACE_TASK_ID unset (`enforcement.go:L100–122`) — INT-04 hook behavior alignment deferred beyond text wiring.

---

## Cross-scope risks (include in PLAN.md)

| Risk | Scopes | Mitigation |
|------|--------|------------|
| Discoveries-only apply still saturates under old rule | S02 + S03 | Ship S03-T01 before or with S02-T02; test combined flow |
| Promotion without spawn still counts as empty apply | S02 + S03 | S03 threshold must treat "discoveries-only" per locked option |
| Reset re-saturates on first empty apply | S03 | Reset clears consecutive-empty counter (S03-T02) |
| SelectNext test breakage | S03 | Update `select_test.go` in same PR as S03-T05 |
| P25-2 PASS but hook still permissive | S04 + S05 | VERIFY notes INT-04 hook gap; E03 scores text presence only |

---

## VERIFY mapping (PLAN.md must include)

| Deliverable | Source | PLAN task IDs | Verification |
|-------------|--------|---------------|--------------|
| D1 | ParentOrchestratorRule in MDC output | S04-T01, S04-T02 | `go test ./internal/install/...` |
| D2 | BLOCKING discovery → task path | S02-T01–T06 | `go test ./internal/loop/...` |
| D3 | Greenfield no sticky STOP at hop 1 | S03-T01, S03-T02 | G1 fixture status |
| D4 | Reset API clears STOP | S03-T03, S03-T04 | `trace loop reset --task …` on G1 |
| D5 | Unified STOP reason | S03-T05, S03-T06 | gate vs export jq |
| D6 | Full PASS | all | `go test ./internal/...`; `score.sh G1 --p25` |

---

## Role work (implementer checklist)

1. Read `AUDIT.md` end-to-end; confirm preflight gates pass.
2. Create `scopes/scope-01-planning/PLAN.md` using locked structure above.
3. Expand each task seed into full PLAN rows (do not omit acceptance/verification/risk columns).
4. Fill Architecture decisions table with threshold options — mark **Recommendation** column; final lock belongs to S02-00 / S03-00 planners.
5. Confirm every AUDIT file:line gap maps to ≥1 task ID.
6. Self-check: no product Go changes; no final threshold numbers committed as implemented constants.
7. Update board row P26-S01-01 → `done` with Notes citing `PLAN.md` path.

## Exit criteria

- [ ] `PLAN.md` exists at `scopes/scope-01-planning/PLAN.md`
- [ ] Sections S02, S03, S04 each have ≥5 task rows with full schema columns
- [ ] Architecture decisions + threshold decision points documented for S02/S03 planners
- [ ] VERIFY mapping D1–D6 complete
- [ ] Cross-scope risks table present
- [ ] No product code changes
- [ ] Board Notes cite `PLAN.md`

## Minimal todos

- [ ] Preflight AUDIT.md sections
- [ ] Draft PLAN.md Overview + Architecture decisions
- [ ] Draft S02 section (INT-01/06 tasks)
- [ ] Draft S03 section (INT-02/05/09 tasks + threshold options)
- [ ] Draft S04 section (P25-2 tasks)
- [ ] Add VERIFY mapping + cross-scope risks
- [ ] Mark own row `done`

## Do not

- Implement product changes
- Pick final saturation hop numbers or consecutive threshold constants (recommend only)
- Add daemon/HTTP or unversioned schema changes
- Start S02/S03/S04 implement rows
- Modify `AUDIT.md` or S00 done history
