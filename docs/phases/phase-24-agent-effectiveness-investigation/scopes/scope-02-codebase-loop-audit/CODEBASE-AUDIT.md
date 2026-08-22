# Codebase loop audit — FM→mechanism map

Phase 24 / Scope 02 / P24-S02-01 deliverable. Maps E01 failure modes (FM-01..FM-10) to live repo mechanisms with file:line evidence. Investigation only — no product changes.

---

## §1 Summary

Mode A (thin graph, G1≡B0) and Mode B (7 discoveries, no new tasks, persistent STOP) share three product choke points: **seed import pins task UUIDs** with no nudge to expand the roster; **P19 saturation fires on the first `loop apply` with zero plan_changes and zero spawned_tasks**, persisting STOP before agents reach EXECUTE; and **standalone `trace add discovery` + `trace_link discovery-mentions-task`** record gaps without creating tasks — only `loop apply` `spawned_tasks[]` creates backlog items. Harness layers add bypass: the Cursor hook **allows edits when `TRACE_TASK_ID` is unset** and installed rules mention gate/DONE only, not gap-pass replanning or task promotion. Export/git snapshots (`trace/graph.json`) can diverge from live `.trace/trace.db` when investigators clone without re-importing seed — live gate then returns `task_not_found` while export still shows deliberation STOP states.

---

## §2 FM mechanism table

| FM-ID | mechanism | file:line | agent-visible symptom | change lever |
|-------|-----------|-----------|----------------------|--------------|
| FM-01 | Seed import upserts tasks by fixed UUID; `ImportSeedTask` preserves existing `work_state` on conflict; loop next `tasks[]` lists store rows only | `internal/domain/seed_import.go:L42–55`, `L308–333`; `experiments/.../runs/G1/seed/gt.json:L10–40`; `internal/loop/next.go:L29–34`, `L275` | `trace tasks` / MCP `trace_tasks` returns same 5 seed IDs forever; next packet task list never grows unless agent calls `trace add task` or `loop apply` spawns | product |
| FM-02 | Default export omits `work_state`, reviews, transitions; uncertainties exported only if present in SQLite; `--strict` gates honesty at export time | `internal/domain/seed_export.go:L55–62`, `L347–348`, `L587–597`; `cmd/trace/seed.go:L136–193` | Session A graph looks empty (0 decisions in export); Session B has discoveries/decisions but 0 uncertainties in committed `graph.json`; clone import lands tasks PENDING regardless of live DONE | product |
| FM-03 | `SelectNext` first-match: hop budget **or already-Stopped** → `hop_budget_exceeded` before `p19_saturated`; P19 signal = last apply had `NewPlanChanges==0 && NewSpawnedTasks==0`; STOP sticky — `Stopped=true` never cleared; gate blocks edit unless phase EXECUTE | `internal/deliberation/select.go:L7–12`; `internal/loop/policy.go:L103–109`; `internal/loop/apply.go:L499–512`; `internal/loop/gate.go:L172–192`; `internal/deliberation/types.go:L8–9` | Export: `stop_reason: p19_saturated`, `hop_count: 1` (`graph.json` L344–353); live gate (DB populated): `reason_code: hop_budget_exceeded`, `recommended_phase: STOP`, edit blocked with `premature_implementation` | product |
| FM-04 | Installed enforcement rules require gate only when `TRACE_TASK_ID` set; no product rule forcing parent orchestrator to own graph; Multitask parent can delegate without Trace | `internal/install/enforcement.go:L20–48`; `experiments/.../runs/G1/.cursor/hooks/trace-loop-gate.sh:L4–7`; `PROMPT-G1-ENFORCE.md:L81–90` | Parent agent edits product code with hook `allow` (no task id); workers gate on their UUID; parent graph stays thin | harness |
| FM-05 | Install text documents `enforce` default **off**; hook allows when `trace` missing or no `TRACE_TASK_ID`; strict config in G1 does not block post-STOP coding | `internal/install/enforcement.go:L28–31`, `L44–45`; `trace-loop-gate.sh:L5–11`; G1 `.trace/config.json` | Agent sees enforcement docs but hook returns `{"permission":"allow"}` without active task; VERIFY PASS while loop status reports STOP violations | harness |
| FM-06 | E01 arms share identical starter (`f70aaea`); product has no experiment-arm isolation in loop/seed | `experiments/RESULTS.md` E01 row; **no code path** — starter git anchor only | Session A scored G1≡B0 before Session B commits; cross-arm comparison confounds harness mode with product | protocol |
| FM-07 | `loop apply` sets `plan_critiqued` only when `plan_changes[]` non-empty; deliberation transition runs after writes; export timestamps follow code commits | `internal/loop/policy.go:L60–62`; `internal/loop/apply.go:L507–512`; `seed_export.go:L722–725` | SPEC/PLANNING-MATRIX land in gap-closure commit `704e2ff` with product; graph export commits (`a37e7c0`) post-code; `plan_critiqued: false` on verify task in export | product |
| FM-08 | MCP `trace_add` lists `discovery` alongside `task` with equal weight; `trace_link` supports `discovery-mentions-task` but no “promote discovery → task”; `trace_loop` apply is only spawn path | `internal/mcp/server.go:L77–78`, `L88–89`, `L193–194`; `internal/mcp/tools_write.go:L17`, `L30`; `cmd/trace/add.go:L68–83`, `L104–114` | Agent records discoveries via MCP/CLI; tool descriptions do not prioritize task creation after discovery; Session B used mentions-task links only | product |
| FM-09 | Composite: FM-01+FM-03+FM-08 product defaults + FM-04 harness bypass + human directed-gap prompt (not in Trace binary) | See FM-01, FM-03, FM-04, FM-08 rows; `INVESTIGATION.md:L11–16` | Build mode thin graph; directed gap rich graph only after explicit human instruction; no product “gap analysis” entry point | protocol |
| FM-10 | Discoveries via `trace add discovery` / apply `discoveries[]`; task promotion only via `loop apply` `spawned_tasks[]` → `ImportSeedTask`; mentions-task link does not create tasks | `cmd/trace/add.go:L104–114`; `internal/loop/apply.go:L409–428`, `L479–497`; `internal/loop/apply_writes.go:L739–747`; `internal/mcp/tools_write.go:L190` | 7 discovery entities, 2× `discovery_mentions_task` → verify `…0050`; `tasks[]` length stays 5; `new_tasks_since_last_step: 0` in loop status | product |

---

## §3 S01 residual reconciliation

### SelectNext: `hop_budget_exceeded` vs `p19_saturated` at hop_count=1

**First STOP (initial apply):** When the first `loop apply` writes zero `plan_changes` and zero `spawned_tasks`, `p19SaturatedFromLastStep` returns true (`internal/loop/policy.go:L103–109`). With `hop_count=1` and `Stopped=false`, `SelectNext` skips the hop-budget branch and hits `p19_saturated` (`internal/deliberation/select.go:L10–11`). `ApplyTransition` persists `StopReason = "p19_saturated"` (`internal/deliberation/select.go:L73–74`).

**Subsequent gate/next (sticky STOP):** Once `Stopped=true`, the first branch fires regardless of hop count: `state.Stopped` → `ReasonHopBudgetExceeded` (`select.go:L7–8`; confirmed by `select_test.go:L88–93`). Live CLI on G1 verify task (2026-08-20, `-C runs/G1`):

```text
loop gate --for edit → reason_code: hop_budget_exceeded, recommended_phase: STOP
policy_inputs.p19_saturated: true, hop_count: 1, stopped: true
```

**Export vs live reason label:** Committed `trace/graph.json` stores persisted `stop_reason: p19_saturated` from the first transition. Live gate JSON reports `hop_budget_exceeded` because `SelectNext` re-evaluates with `Stopped=true`. `buildDeliberationSection` prefers persisted `StopReason` when `dState.Stopped` (`internal/loop/deliberation_packet.go:L198–201`) but gate uses fresh `SelectNext` reason (`gate.go:L159`, `L191–192`).

**Unit tests:** `go test ./internal/deliberation/... -run 'SelectNext|HopBudget' -count=1` → PASS (hop_count=12 → hop_budget; p19_saturated at hop_count=4; stopped → hop_budget_exceeded).

### Export vs `.trace/` SQLite divergence

| Artifact | Role | G1 evidence |
|----------|------|-------------|
| `trace/graph.json` | Git-committed export snapshot | 5 tasks, 7 discoveries, deliberation STOP states |
| `.trace/trace.db` | Live SQLite (gitignored) | 802KB; `trace tasks -C runs/G1` lists 5 seed tasks DONE |
| `seed/gt.json` | Import source for fresh clone | Fixed UUIDs L10–40 |

**Paths:** Export assembled by `BuildSeedDocument` (`seed_export.go:L347–349`) including `deliberation_states` (`L710–719`). Import via `ImportSeedDocument` (`seed_import.go:L23–55`, deliberation at `L220–221`, `L936–943`). Round-trip: `trace seed export -o …` / `trace seed import …` (`cmd/trace/seed.go`).

**S01 `task_not_found`:** Occurred when CLI ran without `-C` pointing at G1 (wrong cwd / empty `.trace/`). Hook correctly uses `trace -C "$root"` (`trace-loop-gate.sh:L13–14`). `TRACE_PROJECT_ROOT` alone is **not** read by `resolveRoot` (`cmd/trace/root.go:L130–135`) — only explicit `-C` or cwd `.trace/` matters.

**Divergence modes:** (1) clone without `seed import` → empty DB, export still rich; (2) export committed after live session → graph.json ahead/behind DB; (3) import preserves deliberation STOP from export without clearing on gap fixes.

### Deliberation reset after gap pass

**No reset API found.** Grep `HopCount\s*=` across `internal/domain/deliberation.go`, `internal/loop/`, `internal/store/` yields only `ApplyTransition` increment (`select.go:L59–70`) and test fixtures. `InitialState` sets hop_count 0 (`types.go:L94–100`) only for **missing** deliberation rows (`domain/deliberation.go:L44–45`). `ImportSeedDeliberationState` upserts exported STOP as-is (`seed_import.go:L936–943`) — no gap-pass branch.

Gap fixes in Session B did not call `loop apply` with `plan_changes` or `spawned_tasks` on verify task, so STOP persisted. **Documented gap:** no product path to clear `Stopped`, reset `hop_count`, or reopen EXECUTE after gap closure without manual deliberation row edit or new apply that adds plan/spawn output (which may still re-saturate if counts stay zero).

---

## §4 Cross-cutting observations

- **Apply vs add split:** Cognitive entities (`discovery`, `decision`, `evidence`) have standalone CLI/MCP add paths; task backlog expansion is intentionally confined to `loop apply` `spawned_tasks[]` (`apply.go:L479–497`). Agents using only `trace add` never touch deliberation or spawn.
- **MCP description gap:** `trace_add` description lists kinds flatly (`server.go:L77–78`); harness PROMPT-G1-ENFORCE L75–77 recommends task add after discovery — product install rules do not repeat this (`enforcement.go:L20–48`).
- **Install enforcement scope:** Gate-before-edit + status-before-DONE only; no nudge for uncertainties, task promotion, or parent-graph ownership.
- **FM-06 / FM-07 protocol notes:** Cross-arm leakage and post-hoc planning are scored at experiment protocol layer (shared starter, git sparsity); product-adjacent via export timestamp ordering only.
- **Execute gating:** `BuildPolicyInputs` requires `plan_critiqued` for `execute_pending` (`policy.go:L96`); Session B verify task has `plan_critiqued: false` in export — agents cannot reach EXECUTE even if STOP cleared.

---

## §5 Optional repro snippet

```bash
# Illustrates FM-03: first apply with empty writes → p19_saturated persisted;
# subsequent gate → hop_budget_exceeded (sticky Stopped).
TRACE="-C experiments/ab-incident-tracker/runs/G1"
TASK=e0100000-0000-4000-8000-000000000050
trace $TRACE loop status --task $TASK | jq '.deliberation,.violations'
trace $TRACE loop gate --task $TASK --for edit | jq '.reason_code,.violations[0].code'
# Compare export stop_reason:
jq '.deliberation_states[] | select(.task_id==env.TASK)' \
  experiments/ab-incident-tracker/runs/G1/trace/graph.json
```

---

## §6 Open questions → S04

- Should sticky STOP re-report original `stop_reason` (p19) vs remapped `hop_budget_exceeded` in gate JSON for agent clarity?
- Rank: auto-spawn task from BLOCKING discovery vs harness prompt vs loop apply UX vs deliberation reset API.
- Does `--strict --enforce` export on CI catch Session A thin-graph failure early enough?
- Portable graph workflow: should clone docs require `seed import` before gate, or should product detect export/DB drift?
- FM-09 collapse: which single intervention moves build mode toward directed-gap richness without human gap prompt?

---

## Evidence appendix (commands run 2026-08-20)

```text
go build -o bin/trace ./cmd/trace                          → OK
go test ./internal/deliberation/... -run 'SelectNext|HopBudget' -count=1 → PASS

trace -C runs/G1 loop status --task …0050
  → violations: premature_implementation; reason hop_budget_exceeded; p19_saturated true; hop_count 1

trace -C runs/G1 loop gate --task …0050 --for edit
  → allowed false; reason_code hop_budget_exceeded

trace (no -C, TRACE_PROJECT_ROOT set) loop gate …0050
  → task_not_found (wrong DB / cwd)

trace -C runs/G1 tasks → 5 seed tasks, all work_state DONE
```
