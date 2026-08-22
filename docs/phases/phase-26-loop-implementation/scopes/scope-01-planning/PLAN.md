# Phase 26 implementation plan

## Overview

Phase 26 closes three product choke points identified in [AUDIT.md](../scope-00-loop-audit/AUDIT.md): **(1) discovery→task promotion** is confined to agent-initiated `loop apply` `spawned_tasks[]` — standalone discovery import, MCP `trace_add`, and mention links record gaps without growing the task roster (**S02**, INT-01 + INT-06); **(2) P19 saturation** fires on the first zero-write apply (including discoveries-only), persisting sticky STOP before EXECUTE, with no reset path (**S03**, INT-02 + INT-05); **(3) STOP reason divergence** — gate/status report `hop_budget_exceeded` while export retains `p19_saturated` (**S03**, INT-09). A fourth harness gap — **`ParentOrchestratorRule` defined but not wired** into install bodies (**S04**, P25-2) — closes E02 P25-2 FAIL. Human gate on promotion is non-negotiable: no background auto-spawn from discovery import alone (INTERVENTION-MATRIX §4). No daemon/HTTP; schema changes require an explicit version-bump plan.

## Architecture decisions

| Decision | Options (from AUDIT) | Recommendation for S02/S03 planner | Rationale |
|----------|----------------------|-------------------------------------|-----------|
| Promotion trigger | Manual `spawned_tasks[]` only; new `promote_discoveries[]` alias; CLI `--from-discovery <id>` | **Primary:** extend/document `spawned_tasks[]` with discovery-sourced entries via domain helper (S02-T01); **secondary:** CLI `--from-discovery` for manual promotion (S02-T03) | Reuses existing apply envelope; avoids silent spawn; helper supplies title/link from BLOCKING discovery |
| BLOCKING filter | All discoveries vs BLOCKING severity only | **BLOCKING-only** for `promotion_candidates[]` and apply promotion path | Matches INTERVENTION-MATRIX and GapPassPrompt intent |
| Idempotency key | Discovery UUID as task seed vs new UUID + link | **Discovery UUID as seed key** with link edge; re-promote is no-op | Prevents duplicate tasks on re-apply; aligns with `ImportSeedTask` upsert semantics |
| Empty apply threshold | 1 (current), 2+ consecutive, N consecutive | **≥2 consecutive** empty applies before saturation | Closes FM-03 greenfield sticky STOP at hop_count=1; S03-00 lock |
| Discoveries-only apply | Counts as empty spawn (current) vs non-saturating | **Non-saturating** when discoveries imported without spawn | Gap-pass records findings without immediate STOP; coordinate with S02-T02 |
| Hop budget | 12 fixed vs profile-specific vs per-phase | **Keep 12** unless G1 re-score proves gap-pass needs raise | AUDIT verified `HopBudget = 12`; change only with test evidence |
| First-apply greenfield exempt | Yes vs consecutive counter only | **Consecutive counter only** (no special-case exempt) | Cleaner persistence; pairs with S03-T02 counter field |
| Consecutive-empty persistence | New DB column vs derive from loop step history | **Prefer derive from loop step history** if feasible; else new column with **schema version bump** documented in S03-T02 | Phase lock: no unversioned schema |
| Canonical STOP reason | Persisted `StopReason` vs fresh `SelectNext` vs dual report | **Persisted `StopReason` wins** for gate/status `reason_code`; add `stop_reason` field to status JSON alongside `why_selected` when they differ | Closes INT-09 G1 repro; export already correct at source |
| `hop_budget_exceeded` remapping | Always when `Stopped=true` vs only when hop ≥ budget | **Only when `HopCount >= HopBudget`** without prior saturation STOP | Sticky saturation STOP keeps `p19_saturated` signal |
| Reset API surface | CLI only vs CLI + MCP vs domain-only | **CLI `trace loop reset --task <id>`** (+ domain fn); MCP deferred unless S03-00 expands | Minimal surface; human-gated recovery |
| Reset scope | Clear STOP/hop only vs also consecutive-empty counter | **Clear `Stopped`, `StopReason`, `HopCount`, consecutive-empty counter** | Prevents immediate re-STOP after reset (S03-T04) |
| ParentOrchestratorRule insert order | Before vs after GapPassPrompt | **After GapPassPrompt** in all three bodies (gate body → GapPassPrompt → ParentOrchestratorRule) | Matches AUDIT proposed pattern and S04 planner lock |
| INT-04 hook alignment | Fix `CursorLoopGateHookScript` failClosed vs TRACE_TASK_ID unset | **Deferred beyond S04 text wiring** | Document in cross-scope risks; VERIFY scores text presence only |

## Dependency graph

```
S02 (promotion) ──┐
                  ├──► S05 VERIFY (D1–D6)
S03 (saturation/reset/reason) ──┤
S04 (installer) ────────────────┘

S04 may parallel S03 after S01; default serial S02 → S03 → S04.

Critical path: S03-T01/T02 should land before or with S02-T02 (discoveries-only + empty-apply interaction).
S04 is independent of S02/S03 (install package only).
```

## S02 — P25-A: discovery→task promotion (INT-01, INT-06)

| ID | INT | Title | File(s) | Change | Acceptance criteria | Verification | Risk | Size | Depends on |
|----|-----|-------|---------|--------|---------------------|--------------|------|------|------------|
| S02-T01 | INT-01 | BLOCKING discovery promotion helper | `internal/domain/` (new helper or extend `seed_import.go`), `internal/loop/apply.go:L479–497`, `internal/mcp/tools_write.go:L121–129` | Domain helper: given BLOCKING discovery ID + goal context, produce task seed payload (title from discovery, discovery→task link). Callable from apply path when agent supplies promotion intent — **not** auto on every discovery import | Helper returns valid task UUID + link metadata; idempotent on re-promote (same discovery → same task ID, no duplicate rows); no spawn without explicit apply/add intent; BLOCKING severity filter enforced | Unit test in `internal/domain/` or `internal/loop/` | Duplicate tasks if idempotency weak; seed UUID collision with existing roster | M | None |
| S02-T02 | INT-01 | `loop apply` promotion path | `internal/loop/apply.go:L409–428`, `L479–497` | Extend apply envelope: implement `promote_discoveries[]` alias **or** document+wire discovery-sourced entries in existing `spawned_tasks[]` via S02-T01 helper; discoveries-only apply unchanged unless paired with spawn/promotion intent | Apply with promotion intent creates task row + discovery link; apply result `spawned_tasks` lists new task IDs; discoveries-only apply without spawn does not create tasks; FM-10 path closed when agent promotes | Integration test: BLOCKING discovery in store → apply with promotion → `ListTasksByGoalID` count +1 | Triggers INT-02 saturation if spawn count still 0 and discoveries-only still counts as empty (coordinate S03-T01) | M | S02-T01 |
| S02-T03 | INT-01 | CLI discoverability | `cmd/trace/add.go:L68–83`, `L104–114` | Surface promotion path in help text; optional `--from-discovery <id>` on `trace add task` **or** documented two-step flow referenced from `loop status` hints | `trace add task --help` (or equivalent) mentions discovery promotion; manual flow creates linked task from BLOCKING discovery; standalone `trace add discovery` still does not auto-spawn | CLI test or documented manual repro in scope Notes | Low — docs/UX only if no new flag | S | S02-T01 |
| S02-T04 | INT-06 | MCP `trace_add` description reorder | `internal/mcp/server.go:L77–78`, `internal/mcp/tools_write.go` | Reorder kind list: discovery first; add nudge: after BLOCKING discovery, call `trace add task` or `loop apply` with spawned_tasks before edits | MCP tool schema/description contains promotion ordering text; parity with `GapPassPrompt` intent; flat kind list no longer implies equal weight | `grep`/test on registered tool description string | Agent confusion if text too long | XS | None |
| S02-T05 | INT-01 | `loop next` promotion candidates | `internal/loop/next.go:L190–195`, packet builder | Add optional `promotion_candidates[]` (BLOCKING discoveries without linked executable task) to next packet | Next packet lists open BLOCKING discoveries eligible for promotion; empty array when none; does not auto-spawn | Unit test with fixture discoveries | Packet bloat; false positives if severity filter wrong | M | S02-T01 |
| S02-T06 | INT-01, INT-06 | End-to-end promotion test | `internal/loop/*_test.go`, optionally `internal/mcp/*_test.go` | Test: create BLOCKING discovery → promote via apply or add → task appears in `ListTasksByGoalID` | FM-10 repro closed: discovery count > 0 implies path to task count increase when agent promotes; MCP description present in registry | `go test ./internal/loop/... ./internal/mcp/...` | Flaky if UUID fixtures drift | S | S02-T02 |
| S02-T07 | INT-01 | Seed import discovery parity | `internal/domain/seed_import.go:L70–75`, `L308–334` | Ensure promotion helper integrates with `ImportSeedTask` upsert (preserve work_state on conflict); imported BLOCKING discoveries remain promotable | Clone import lands discoveries without auto-spawn; promote after import uses same idempotency key; existing seed task UUIDs unchanged | Unit test: import discovery → promote → single task row | Seed import pins roster; promotion must not clobber work_state | S | S02-T01 |

### AUDIT gap coverage (S02)

| AUDIT row | Task ID(s) |
|-----------|------------|
| `apply.go:L479–497` spawned_tasks only path | S02-T01, S02-T02 |
| `apply.go:L409–428` discoveries-only | S02-T02, S03-T01 (cross-scope) |
| `add.go:L104–114` standalone discovery | S02-T03 |
| `add.go:L68–83` task add undiscoverable | S02-T03 |
| `server.go:L77–78` flat MCP kind list | S02-T04 |
| `tools_write.go:L121–129` discovery-mentions-task | S02-T01, S02-T04 |
| `seed_import.go:L308–334`, `L70–75` | S02-T07 |
| `next.go:L190–195` static task list | S02-T05 |
| `gappass.go:L8–11` harness-only nudge | S02-T04 |
| `apply.go:L639–649` buildPromotionBlocked advisory | Cross-scope risk (misleading field; no product change required for P25-A) |

## S03 — P25-B: loop recalibration + deliberation reset (INT-02, INT-05, INT-09)

| ID | INT | Title | File(s) | Change | Acceptance criteria | Verification | Risk | Size | Depends on |
|----|-----|-------|---------|--------|---------------------|--------------|------|------|------------|
| S03-T01 | INT-02 | P19 saturation recalibration | `internal/loop/policy.go:L103–109`, `internal/loop/apply.go:L499–500`, `internal/deliberation/select.go:L10–11` | Replace immediate zero-write saturation with **consecutive empty apply counter** (named constant, default ≥2 per S03-00 lock); treat discoveries-only apply as **non-saturating** when no spawn/plan_changes (per architecture decision) | First empty apply on greenfield does **not** set `Stopped=true`; second consecutive empty apply saturates; discoveries-only apply without spawn does not increment empty counter; unit tests updated (`select_test.go:L64–67`, `L88–93`) | `go test ./internal/deliberation/... ./internal/loop/...`; G1 fixture: hop_count=1 no longer sticky STOP after single empty apply | Regression if counter logic wrong; interaction with hop budget | M | None |
| S03-T02 | INT-02 | Persist consecutive-empty counter | `internal/store/deliberation.go:L28–60`, `internal/domain/deliberation.go`, deliberation state struct | Add `ConsecutiveEmptyApplies` field **or** derive from loop step history; if new column → document schema version bump in implement Notes | Counter survives apply/status round-trip; resets on non-empty apply (spawn or plan_changes); exported in seed if persisted field | Store/domain unit tests | Schema migration if new column — requires version bump per phase lock | M | S03-T01 |
| S03-T03 | INT-05 | Deliberation reset API | `internal/domain/deliberation.go:L14–72`, `internal/store/deliberation.go`, `cmd/trace/loop.go` | New domain fn e.g. `ResetDeliberationState(taskID)` clearing `Stopped`, `StopReason`, `HopCount`, consecutive-empty counter; CLI `trace loop reset --task <id>` | After reset on G1 verify task: `stopped=false`, `hop_count=0`; phase moves toward EXECUTE (subject to `plan_critiqued` gate per AUDIT INT-05); no product path sets `Stopped=false` today — this becomes the canonical reset | `go test ./internal/domain/...`; manual CLI on G1 fixture | Bypasses human STOP intent if ungated — document when reset is allowed | M | None |
| S03-T04 | INT-05 | Reset + gap-pass integration | `internal/loop/gate.go:L188–192`, `internal/loop/apply.go:L507–512`, `apply.go:L593–594` | After reset, empty apply uses new saturation threshold (not immediate re-STOP); status packet hints reset when saturated + discoveries closed | Gap-pass flow: saturate → reset → discoveries-only apply does not immediately re-STOP (per threshold lock); status surfaces recovery guidance | Integration test sequence apply→saturate→reset→apply | Infinite loop if threshold too high | M | S03-T01, S03-T03 |
| S03-T05 | INT-09 | Unified STOP reason | `internal/deliberation/select.go:L7–8`, `L10–11`, `internal/loop/gate.go:L159`, `L191–192`, `internal/loop/deliberation_packet.go:L188–202`, `internal/loop/apply.go:L652–671` | When `Stopped=true`, gate/status/`why_selected` prefer **persisted** `StopReason` over remapped `hop_budget_exceeded`; add `stop_reason` field to status JSON alongside `why_selected`; reorder SelectNext branches so hop_budget only when hop ≥ budget | G1 repro closed: gate `reason_code` matches export `stop_reason` (`p19_saturated` case); status includes both fields when they differ during transition; next packet no longer contradicts export without hint | Unit tests for SelectNext stopped branch; gate JSON test; AUDIT jq repro passes | Breaking change for agents expecting hop_budget_exceeded on all STOP | M | None |
| S03-T06 | INT-09 | Export/gate alignment test | `internal/domain/seed_export.go:L710–719`, `internal/loop/gate_test.go` or new test | Assert export `stop_reason` equals gate `reason_code` after saturation STOP | `seed export` deliberation_states match live gate for same task | `go test ./internal/loop/... ./internal/domain/...` | Low | S | S03-T05 |
| S03-T07 | INT-02 | Hop budget policy documentation | `internal/deliberation/types.go:L8–9`, `select_test.go:L64–67` | Document `HopBudget = 12` as locked default; update tests if SelectNext branch order changes; no raise unless S03-00 re-locks with evidence | Tests document hop_budget semantics post S03-T05; HopBudget constant unchanged unless planner re-locks | `go test ./internal/deliberation/...` | Long gap-pass sessions may still hit budget before meaningful work | XS | S03-T05 |

### AUDIT gap coverage (S03)

| AUDIT row | Task ID(s) |
|-----------|------------|
| `types.go:L8–9` HopBudget = 12 | S03-T07 (architecture decision) |
| `policy.go:L103–109` zero-write saturation | S03-T01 |
| `apply.go:L499–500`, `L507–512` apply saturation inline | S03-T01, S03-T04 |
| `apply.go:L593–594` status recompute | S03-T04 |
| `select.go:L10–11` P19Saturated branch | S03-T01, S03-T05 |
| `select.go:L7–8` sticky STOP hop_budget remap | S03-T05 |
| `select_test.go:L64–67`, `L88–93` | S03-T01, S03-T07 |
| `deliberation.go:L14–72` no reset API | S03-T03 |
| `select.go:L59–62`, `L72–74` hop/STOP forward-only | S03-T03 |
| `store/deliberation.go:L28–60` upsert only | S03-T02, S03-T03 |
| `seed_import.go:L936–943` import preserves STOP | S03-T03 (reset after import) |
| `gate.go:L188–192` edit block | S03-T04 |
| `gate.go:L159`, `L191–192` reason_code from SelectNext | S03-T05 |
| `deliberation_packet.go:L188–202` dual labels | S03-T05 |
| `apply.go:L652–671` status missing stop_reason | S03-T05 |
| `cmd/trace/loop.go:L98–113` gate envelope | S03-T05 |
| `seed_export.go:L710–719` export stop_reason | S03-T06 |

## S04 — Installer fix P25-2 (ParentOrchestratorRule wiring)

| ID | INT | Title | File(s) | Change | Acceptance criteria | Verification | Risk | Size | Depends on |
|----|-----|-------|---------|--------|---------------------|--------------|------|------|------------|
| S04-T01 | P25-2 | Wire ParentOrchestratorRule into AgentsEnforcementBlock | `internal/install/enforcement.go:L65` (`AgentsEnforcementBlock`) | Append `"\n" + ParentOrchestratorRule + "\n"` after `GapPassPrompt` | `rg ParentOrchestratorRule internal/install/` hits usage in `AgentsEnforcementBlock`; AGENTS install output contains "Parent orchestrator" | `go test ./internal/install/...` | Low — text-only | XS | None |
| S04-T02 | P25-2 | Wire ParentOrchestratorRule into cursorRulesMDCContent | `internal/install/enforcement.go:L83` (`cursorRulesMDCContent`) | Append `"\n" + ParentOrchestratorRule + "\n"` after `GapPassPrompt` | Fresh `trace install cursor --write` output contains "Parent orchestrator"; MDC body order: gate body → GapPassPrompt → ParentOrchestratorRule | `go test ./internal/install/...`; manual install repro | Low — text-only | XS | None |
| S04-T03 | P25-2 | Wire ParentOrchestratorRule into claudeFallbackRulesContent | `internal/install/enforcement.go:L97` (`claudeFallbackRulesContent`) | Append `"\n" + ParentOrchestratorRule + "\n"` after `GapPassPrompt` | Claude fallback output contains "Parent orchestrator"; same concat order as MDC | `go test ./internal/install/...` | Low — text-only | XS | None |
| S04-T04 | P25-2 | Enforcement test assertions | `internal/install/enforcement_test.go` | Assert `"Parent orchestrator"` (or rule substring) present in `AgentsEnforcementBlock`, `cursorRulesMDCContent()`, and Claude fallback output | Test fails on pre-fix tree, passes after S04-T01–T03 | `go test ./internal/install/...` | Low | S | S04-T01, S04-T02, S04-T03 |
| S04-T05 | P25-2 | Install smoke verification | `internal/install/enforcement.go`, `cmd/trace/install` (read-only) | Document/automate pre/post `rg 'ParentOrchestratorRule' internal/install/` check: definition + three body usages (not tests-only) | `rg` shows ≥4 hits (definition + 3 funcs); E02 P25-2 text criterion satisfied | `rg ParentOrchestratorRule internal/install/`; optional install subtest | P25-2 PASS but hook still permissive (INT-04) — VERIFY notes hook gap | XS | S04-T04 |

### AUDIT gap coverage (S04)

| AUDIT row | Task ID(s) |
|-----------|------------|
| `ParentOrchestratorRule` defined L20–35, unused in bodies | S04-T01, S04-T02, S04-T03 |
| `AgentsEnforcementBlock` ends GapPassPrompt only | S04-T01 |
| `cursorRulesMDCContent` / `claudeFallbackRulesContent` missing rule | S04-T02, S04-T03 |
| `CursorLoopGateHookScript` L100–122 INT-04 hook gap | Cross-scope risk (deferred) |
| Phase 25 delta GapPassPrompt wired | No task — already done |

## Cross-scope risks

| Risk | Scopes | Mitigation |
|------|--------|------------|
| Discoveries-only apply still saturates under old rule | S02 + S03 | Ship S03-T01 before or with S02-T02; integration test combined flow |
| Promotion without spawn still counts as empty apply | S02 + S03 | S03 threshold must treat discoveries-only as non-saturating (architecture decision) |
| Reset re-saturates on first empty apply | S03 | Reset clears consecutive-empty counter (S03-T02, S03-T03) |
| SelectNext test breakage | S03 | Update `select_test.go` in same PR as S03-T05 / S03-T07 |
| P25-2 PASS but hook still permissive | S04 + S05 | VERIFY notes INT-04 hook gap; E03 scores text presence only |
| Misleading `promotion_blocked` in status JSON | S02 | Document advisory nature; out of scope for P25-A unless S02-00 expands |
| Schema bump omitted for new deliberation field | S03 | S03-T02 must cite version bump in implement Notes if column added |
| S02 promotion triggers saturation before S03 lands | S02 + S03 | Board serial default S02→S03; if parallel, merge S03-T01 first |

## VERIFY mapping (D1–D6)

| Deliverable | PLAN task ID(s) | Verification command |
|-------------|-----------------|----------------------|
| D1 | ParentOrchestratorRule in MDC output | S04-T02, S04-T04, S04-T05 | `go test ./internal/install/...`; `rg ParentOrchestratorRule internal/install/` |
| D2 | BLOCKING discovery → task path | S02-T01–T07 | `go test ./internal/loop/... ./internal/mcp/...` |
| D3 | Greenfield no sticky STOP at hop 1 | S03-T01, S03-T02 | G1 fixture: `trace loop status --task …0050` after single empty apply |
| D4 | Reset API clears STOP | S03-T03, S03-T04 | `trace loop reset --task …` on G1 verify task |
| D5 | Unified STOP reason | S03-T05, S03-T06 | gate vs export jq (AUDIT evidence appendix) |
| D6 | Full PASS | all | `go test ./internal/...`; `score.sh G1 --p25` |

## Open questions (human gate)

None — AUDIT and S01-00 planner lock sufficient threshold options; S02-00 / S03-00 planners lock final constants at implement time.
