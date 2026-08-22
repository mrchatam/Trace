# Loop audit — Phase 26

Phase 26 / Scope 00 / P26-S00-01 deliverable. Investigation only — maps live codebase touch points for INT-01, INT-02, INT-05, INT-06, INT-09, and installer gap P25-2. Refreshed against live tree 2026-08-20; Phase 24 `CODEBASE-AUDIT.md` used as baseline (line numbers re-verified).

---

## Executive summary

Three product choke points still block P25-A/B effectiveness: **(1) discovery→task promotion** is confined to `loop apply` `spawned_tasks[]` — standalone `trace add discovery`, MCP `trace_add`, and `discovery-mentions-task` links record gaps without growing the task roster; **(2) P19 saturation** fires on the first apply with zero `plan_changes` and zero `spawned_tasks`, persisting sticky `Stopped=true` before agents reach EXECUTE; **(3) no deliberation reset API** clears `Stopped`, `hop_count`, or reopens EXECUTE after gap closure. A fourth harness gap remains: **`ParentOrchestratorRule` is defined but not wired** into `cursorRulesMDCContent()` / `claudeFallbackRulesContent()` (P25-2 E02 FAIL). Phase 25 closed the INT-03 orphan: `GapPassPrompt` is now concatenated into AGENTS/MDC/Claude install bodies. INT-09 reason divergence persists: export/next `stop_reason` retains `p19_saturated` while gate/status `reason_code` / `why_selected` report `hop_budget_exceeded` on sticky STOP.

---

## INT-01 / INT-06: Discovery→task promotion

| file:line | current behavior | gap vs INT target | risk |
|-----------|------------------|---------------------|------|
| `internal/loop/apply.go:L479–497` | `loop apply` iterates `writes.spawned_tasks[]`, calls `ImportSeedTask` per entry; only path that atomically expands backlog during deliberation | No guided promotion from BLOCKING discoveries — agent must know apply envelope schema | Session B: 7 discoveries, 0 spawned tasks, roster stays at 5 |
| `internal/loop/apply.go:L409–428` | `discoveries[]` in same apply path imports discovery entities + links only | Discoveries-only apply counts as zero spawn → triggers saturation (see INT-02) | Gap pass records findings but does not create follow-up tasks |
| `cmd/trace/add.go:L104–114` | `trace add discovery` creates discovery entity; no task side-effect | Standalone add never promotes BLOCKING discovery to task | Agents using CLI add stay on thin graph |
| `cmd/trace/add.go:L68–83` | `trace add task` creates task via `CreateTask` (separate from loop apply) | Exists but undiscoverable vs discovery in flat tool lists; no link-from-discovery helper | Manual task add bypasses loop deliberation context |
| `internal/mcp/server.go:L77–78` | `trace_add` description lists kinds flatly: `goal\|task\|decision\|assumption\|discovery\|plan-change\|claim\|evidence` | No ordering nudge (task after discovery); mirrors CLI parity, not INT-06 promotion UX | MCP-first agents treat discovery and task as equal-weight |
| `internal/mcp/tools_write.go:L121–129`, `L190–191` | `discovery` kind → `CreateDiscovery`; `discovery-mentions-task` link only records mention | No `promote discovery → task` tool or rel | Mentions-task links do not expand `tasks[]` |
| `internal/domain/seed_import.go:L308–334` | `ImportSeedTask` upserts by fixed UUID; preserves existing `work_state` on conflict | Seed import pins roster; no promotion from imported discoveries | Clone import lands same 5 seed IDs |
| `internal/domain/seed_import.go:L70–75` | Discoveries imported as entities only | No import-time task spawn from BLOCKING severity | Portable graph carries discoveries without backlog growth |
| `internal/loop/next.go:L190–195` | `loop next` `tasks[]` built from `ListTasksByGoalID` — store rows only | Packet never suggests promotion candidates from open BLOCKING discoveries | Next packet task list static unless spawn path used |
| `internal/install/gappass.go:L8–11` | `GapPassPrompt` nudges `trace add task` for BLOCKING gaps (wired in AGENTS/MDC/Claude) | Harness-only text; not in MCP `trace_add` description; no product promotion API | E02 gap-pass works when installed; MCP agents miss ordering |
| `internal/loop/apply.go:L639–649` | `buildPromotionBlocked` calls `CheckPromotionGate` (evaluation-based, advisory) | Unrelated to discovery→task; does not gate or enable spawn | Misleading `promotion_blocked` in status JSON |

---

## INT-02: Saturation recalibration

| file:line | current behavior | gap vs INT target | risk |
|-----------|------------------|---------------------|------|
| `internal/deliberation/types.go:L8–9` | `HopBudget = 12` — max phase hops per seed | Fixed N=12; no profile-specific budget (gap-pass vs build) | Long gap-pass sessions may hit budget before meaningful work |
| `internal/loop/policy.go:L103–109` | `p19SaturatedFromLastStep`: true when last loop step has `NewPlanChanges==0 && NewSpawnedTasks==0` (or `MaxIterationsReached`) | First empty apply on greenfield saturates immediately | FM-03: STOP at hop_count=1 before EXECUTE |
| `internal/loop/apply.go:L499–500` | Apply path sets `out.Saturated` from same zero-write rule inline | Same trigger at write time; feeds deliberation transition | Empty apply persists STOP in same transaction |
| `internal/loop/apply.go:L507–512` | `BuildPolicyInputs(..., out.Saturated)` then `ApplyDeliberationTransition` | Saturation flag passed into transition on every apply | First apply with discoveries-only still saturates (spawn count 0) |
| `internal/loop/apply.go:L593–594` | `loop status` recomputes `saturated` from last step payload | Status surfaces `reason: tasks_and_plan_unchanged` when saturated | Agent sees saturation but no recovery guidance in product |
| `internal/deliberation/select.go:L10–11` | When not yet stopped and `P19Saturated`, `SelectNext` → STOP / `p19_saturated` | First transition persists correct reason | — |
| `internal/deliberation/select.go:L7–8` | When `Stopped=true`, first branch → STOP / `hop_budget_exceeded` (beats p19 branch) | Sticky STOP remaps reason on subsequent reads (INT-09 overlap) | Agents lose saturation signal after first persist |
| `internal/deliberation/select_test.go:L64–67`, `L88–93` | Unit tests lock hop_budget at N=12 and stopped→hop_budget_exceeded | Tests document current semantics; no first-apply exemption | S01 must update tests if threshold changes |

### Threshold options (document only — S01 owns final pick)

| Knob | Current (verified live) | Options |
|------|------------------------|---------|
| Hop budget | `deliberation.HopBudget` = **12** (`types.go:L8–9`) | Keep 12; raise for gap-pass profile; separate budget per phase profile |
| P19 saturation | True when last apply had **zero** `plan_changes` **and** zero `spawned_tasks` (`policy.go:L103–109`; apply inline `L499–500`) | Exempt first apply on greenfield; require N consecutive empty applies; treat discoveries-only apply as non-saturating |
| Sticky STOP | `Stopped=true` never cleared in product paths (see INT-05) | New reset API; manual DB edit; re-apply with plan_changes/spawn (may re-saturate if counts stay zero) |
| Reason UX | Sticky STOP → `hop_budget_exceeded` in gate/status even when persisted `StopReason` was `p19_saturated` | Unify on persisted `StopReason`; remap in gate only; dual report with recovery hint (see INT-09) |

---

## INT-05: Deliberation reset

| file:line | current behavior | gap vs INT target | risk |
|-----------|------------------|---------------------|------|
| `internal/domain/deliberation.go:L14–72` | `ApplyDeliberationTransition` loads state, runs `ApplyTransition`, upserts — **forward-only** | **No reset API** — no function clears `Stopped` or zeroes `hop_count` | Gap closure cannot reopen EXECUTE without manual intervention |
| `internal/deliberation/select.go:L59–62` | Hop count increments on non-terminal transitions only; STOP rows do not increment past budget | Reset would need explicit hop_count write — not exposed | Sticky hop_count=1 with STOP persists forever |
| `internal/deliberation/select.go:L72–74` | `ApplyTransition` sets `Stopped=true` and `StopReason` on STOP; never clears | No branch sets `Stopped=false` | Once saturated, all gate paths block edit |
| `internal/store/deliberation.go:L28–60` | `UpsertDeliberationState` replaces row fields on conflict | Store layer supports arbitrary upsert but no domain reset caller | Manual SQL could reset; no CLI/MCP |
| `internal/domain/seed_import.go:L936–943` | `ImportSeedDeliberationState` upserts exported STOP as-is | Import preserves STOP from export; no gap-pass clear branch | Clone re-import locks STOP from committed graph |
| `internal/domain/deliberation.go:L44–45` | Missing row → `InitialState` (ORIENT, hop_count 0) | Only applies to **new** tasks, not recovery of stopped tasks | New tasks start fresh; existing seed tasks stay stopped |
| `internal/loop/gate.go:L188–192` | Edit allowed only when `phase==EXECUTE && !stopped` | Even if STOP cleared manually, `plan_critiqued=false` blocks EXECUTE via policy | Session B verify task: `plan_critiqued: false` in export |

**Reset API exists?** **No.** Repo-wide grep for `ResetDeliberation`, `ClearDeliberation`, `Stopped = false`, `HopCount = 0` in product paths (excluding tests) returns no matches. Only `ApplyTransition` increments hop count; only import/transition sets STOP.

---

## INT-09: Unified STOP reason

| file:line | current behavior | gap vs INT target | risk |
|-----------|------------------|---------------------|------|
| `internal/deliberation/select.go:L7–8` | First-match: `HopCount >= HopBudget \|\| Stopped` → `ReasonHopBudgetExceeded` | Sticky STOP always reports hop_budget, not original saturation reason | Gate/status misleading after first p19 transition |
| `internal/deliberation/select.go:L10–11` | `P19Saturated` branch only reached when `!Stopped` | First apply correctly persists `p19_saturated` | — |
| `internal/loop/gate.go:L159`, `L191–192` | Gate calls fresh `SelectNext`; violation `reason_code` = SelectNext reason | Live gate: `reason_code: hop_budget_exceeded` while export has `p19_saturated` | Agent follows gate JSON, misattributes STOP cause |
| `internal/loop/deliberation_packet.go:L188–202` | `buildDeliberationSection`: `WhySelected` from SelectNext; `StopReason` prefers **persisted** `dState.StopReason` when `dState.Stopped` | Next packet can show **both** `why_selected: hop_budget_exceeded` and `stop_reason: p19_saturated` | Dual labels without recovery hint |
| `internal/loop/apply.go:L652–671` | `buildStatusDeliberation` exposes `why_selected` from SelectNext only; **no `stop_reason` field** | Status omits persisted reason entirely | Status JSON hides original saturation cause |
| `cmd/trace/loop.go:L98–113` | Gate envelope top-level `reason_code` from first violation | Mirrors SelectNext remapping | Consistent with gate but diverges from export |
| `internal/domain/seed_export.go:L710–719` | Export writes persisted `StopReason` to `deliberation_states[]` | Git snapshot shows `p19_saturated` | Portable graph vs live gate UX mismatch |
| `experiments/.../G1/trace/graph.json:L345–353` | Verify task export: `stop_reason: p19_saturated`, `hop_count: 1` | Evidence of persisted vs reported split | E02 dogfood repro |

### Threshold options (reason UX — document only)

| Knob | Current | Options |
|------|---------|---------|
| Gate/status reason | Fresh `SelectNext` → `hop_budget_exceeded` when `Stopped=true` | Unify on persisted `StopReason`; remap in gate only; dual report + recovery hint field |
| Status packet | No `stop_reason` field | Add persisted `stop_reason` alongside `why_selected` |
| Export | Persisted `stop_reason` at transition time | Already correct at source; gate should align |

---

## Installer gap (P25-2)

- **`ParentOrchestratorRule` definition:** `internal/install/enforcement.go:L20–35` — documents parent orchestrator `TRACE_TASK_ID` ownership, preToolUse deny, failClosed semantics.
- **Used in `AgentsEnforcementBlock`?** **No.** Block ends with `GapPassPrompt` only (`enforcement.go:L65`); `ParentOrchestratorRule` not concatenated.
- **Missing from `cursorRulesMDCContent` / `claudeFallbackRulesContent`?** **Yes (confirmed).** Both append `GapPassPrompt` (`L83`, `L97`) but not `ParentOrchestratorRule`. `rg ParentOrchestratorRule` hits only definition + tests + docs — no install body reference.
- **Hook script:** `CursorLoopGateHookScript` (`L100–122`) allows edits when `TRACE_TASK_ID` unset (`L106–108`) — contradicts failClosed doc in `ParentOrchestratorRule` until rule is installed and hook behavior aligned (INT-04 scope beyond S04 text wiring).
- **Proposed insert pattern (S04):** Mirror `GapPassPrompt` concat — append `"\n" + ParentOrchestratorRule + "\n"` after `GapPassPrompt` in `AgentsEnforcementBlock`, `cursorRulesMDCContent`, and `claudeFallbackRulesContent` (order: gate body → GapPassPrompt → ParentOrchestratorRule per S04 planner lock).

**Phase 25 delta:** INT-03 fixed — `internal/install/gappass.go` defines `GapPassPrompt`; wired into AGENTS (`L65`), MDC (`L83`), Claude fallback (`L97`). P25-2 remains open.

---

## Delta from Phase 24 CODEBASE-AUDIT

**Confirmed unchanged (line numbers stable or ±1–2):**

- `SelectNext` first-match STOP branches (`select.go:L7–11`) — P24 cited L7–12 ✓
- `HopBudget = 12` (`types.go:L8–9`) — P24 L8–9 ✓
- `p19SaturatedFromLastStep` zero-write rule (`policy.go:L103–109`) — P24 L103–109 ✓
- Apply saturation + deliberation transition (`apply.go:L499–512`) — P24 L499–512 → now L499–513 (+1 line shift)
- Gate edit block (`gate.go:L172–192`) — P24 L172–192 → now L172–193
- Gate context SelectNext (`gate.go:L159`) — stable
- `buildDeliberationSection` persisted stop_reason preference (`deliberation_packet.go:L198–201`) — stable
- Seed import task UUID pin (`seed_import.go:L42–55`, `L308–334`) — stable
- MCP flat kind list (`server.go:L77–78`) — stable
- No deliberation reset API — P24 §3 finding holds

**Changed since Phase 24:**

- **`GapPassPrompt` wired (Phase 25 P25-S01-02a):** `internal/install/gappass.go` new; `AgentsEnforcementBlock`, `cursorRulesMDCContent`, `claudeFallbackRulesContent` now include gap-pass nudge. Closes INT-03 orphan noted in P24 §4.
- **`ParentOrchestratorRule` constant added (Phase 25)** but still **not** in MDC/Claude bodies — E02 P25-2 FAIL persists.
- **`buildPromotionBlocked` / `CheckPromotionGate`** in status path (`apply.go:L639–649`, `outcomes.go:L436–467`) — evaluation advisory; not discovery promotion (P24 did not emphasize this field).

**New files (not in P24 file list):**

- `internal/install/gappass.go` — GapPassPrompt source

**Re-verified P24 line drift (minor):**

| P24 reference | Live location | Notes |
|---------------|---------------|-------|
| `apply.go:L499–512` | `L499–513` | Transition error wrap added one line |
| `gate.go:L172–192` | `L172–193` | Closing brace shift only |
| `enforcement.go:L20–48` | `L20–35` (ParentOrchestratorRule), `L37–48` (EnforcementRulesMarkdown) | Split; ParentOrchestratorRule now separate const |
| `next.go:L29–34`, `L275` | `L29–34` (TaskSummary), `L190–195` (task list), `L275` (delib section) | Task list evidence at L190–195 not L29–34 |

**Preflight — all audit paths exist (2026-08-20):**

| Path | Status |
|------|--------|
| `internal/deliberation/` | ✓ (select.go, types.go, select_test.go) |
| `internal/loop/` | ✓ (16 files) |
| `internal/store/deliberation.go` | ✓ |
| `internal/domain/deliberation.go` | ✓ |
| `internal/domain/seed_export.go` | ✓ |
| `internal/domain/seed_import.go` | ✓ |
| `cmd/trace/loop.go` | ✓ |
| `cmd/trace/add.go` | ✓ |
| `cmd/trace/seed.go` | ✓ |
| `internal/mcp/server.go` | ✓ |
| `internal/mcp/tools_write.go` | ✓ |
| `internal/install/enforcement.go` | ✓ (+ `gappass.go` sibling) |

No `internal/seed/` package (seed lives under `internal/domain/` + `cmd/trace/seed.go`) — matches planner correction.

---

## Evidence appendix

### Commands run (2026-08-20)

```text
# Preflight path verification via Glob — all files present (see table above)

# G1 fixture repro (verify task …0050)
trace -C experiments/ab-incident-tracker/runs/G1 loop status --task e0100000-0000-4000-8000-000000000050
  → saturated: true, reason: tasks_and_plan_unchanged
  → deliberation.why_selected: hop_budget_exceeded, hop_count: 1, stopped: true
  → policy_inputs.p19_saturated: true, plan_critiqued: false
  → violations[0].reason_code: hop_budget_exceeded

trace -C experiments/ab-incident-tracker/runs/G1 loop gate --task e0100000-0000-4000-8000-000000000050 --for edit
  → allowed: false, reason_code: hop_budget_exceeded, recommended_phase: STOP

# Export persisted reason (committed graph.json)
jq '.deliberation_states[] | select(.task_id=="e0100000-0000-4000-8000-000000000050")' \
  experiments/ab-incident-tracker/runs/G1/trace/graph.json
  → stop_reason: p19_saturated, hop_count: 1, stopped: true, plan_critiqued: false

# P25-2 wiring check
rg 'ParentOrchestratorRule' internal/install/
  → definition enforcement.go:L26; tests enforcement_test.go; NOT in cursorRulesMDCContent/claudeFallbackRulesContent bodies

# Reset API check
rg 'ResetDeliberation|ClearDeliberation|Stopped\s*=\s*false' --glob '*.go' .
  → no product matches
```

### G1 status snippet (abbreviated)

```json
{
  "saturated": true,
  "reason": "tasks_and_plan_unchanged",
  "deliberation": {
    "why_selected": "hop_budget_exceeded",
    "hop_count": 1,
    "stopped": true,
    "policy_inputs": { "p19_saturated": true, "plan_critiqued": false }
  },
  "violations": [{ "reason_code": "hop_budget_exceeded", "recommended_phase": "STOP" }]
}
```

Export for same task: `"stop_reason": "p19_saturated"` — INT-09 divergence confirmed live.
