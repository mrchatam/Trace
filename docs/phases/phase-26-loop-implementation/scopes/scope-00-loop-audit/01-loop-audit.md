# P26-S00-01 — Loop audit implementer

## Metadata
- id: P26-S00-01
- todo_ids: [P26-S00-01]
- role: implementer
- skills: [research, code-explorer]
- mcps: [user-codegraph]
- agents: [explore]
- verification: automated
- hooks: []

## Objective

Produce `AUDIT.md` mapping live codebase touch points for **INT-01, INT-02, INT-05, INT-06, INT-09**, and the **P25-2 installer wiring gap**. Refresh Phase 24 evidence against current repo — do not copy `CODEBASE-AUDIT.md` blindly. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [00-PLANNER.md](00-PLANNER.md)
- [Phase 25 DR-HANDOFF](../../../phase-25-orchestrator-gap-pass/DR-HANDOFF.md)
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md)
- [Phase 24 CODEBASE-AUDIT baseline](../../../phase-24-agent-effectiveness-investigation/scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md) — prior FM→mechanism map; verify line numbers still hold

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Output path | `docs/phases/phase-26-loop-implementation/scopes/scope-00-loop-audit/AUDIT.md` |
| Product code / tests | **Forbidden** |
| Threshold decisions | Document **current constants + options only** — do **not** pick final Phase 26 values (S01 owns that) |
| Tools | `rg`, Read, optional `codegraph_explore` with `projectPath` = repo root |
| Dogfood fixture | `experiments/ab-incident-tracker/runs/G1` (optional repro; use `-C` flag) |

## Preflight / Plan

Before writing `AUDIT.md`:

1. Confirm each path in **Files to audit** exists (`Glob` / `ls`); note renames in board Notes if any drift from this list.
2. Skim Phase 24 `CODEBASE-AUDIT.md` §2–§3 for expected mechanisms; re-verify file:line against live tree.
3. For each INT section, plan which symbols to grep (`SelectNext`, `spawned_tasks`, `ParentOrchestratorRule`, etc.).
4. Confirm `ParentOrchestratorRule` is defined in `enforcement.go` but **not** concatenated into `cursorRulesMDCContent()` / `claudeFallbackRulesContent()` (expected P25-2 gap).

## Files to audit (minimum — verified 2026-08-20)

> **Path note:** There is no `internal/seed/` package. Seed export/import lives under `internal/domain/` + `cmd/trace/seed.go`.

| Package / path | What to find |
|----------------|--------------|
| `internal/deliberation/` | `SelectNext`, `ApplyTransition`, `HopBudget` (N=12), reason codes |
| `internal/loop/` | `apply.go` (`spawned_tasks[]`, saturation, transition); `policy.go` (`p19SaturatedFromLastStep`, `BuildPolicyInputs`); `gate.go` (edit gate, `reason_code`); `next.go`, `deliberation_packet.go` (status/export STOP reason) |
| `internal/store/deliberation.go` | `DeliberationState` row (`hop_count`, `stopped`, `stop_reason`); upsert/get |
| `internal/domain/deliberation.go` | `GetDeliberationState`, `ApplyDeliberationTransition` — any reset surface? |
| `internal/domain/seed_export.go` | `deliberation_states` in export; discovery/task export |
| `internal/domain/seed_import.go` | `ImportSeedTask`, `ImportSeedDeliberationState` — promotion + STOP persistence |
| `cmd/trace/loop.go` | `gate`, `status`, `apply`, `next` CLI; gate JSON envelope |
| `cmd/trace/add.go` | standalone `discovery` vs `task` add paths |
| `cmd/trace/seed.go` | export/import CLI flags (`--strict`) |
| `internal/mcp/server.go` | `trace_add` tool **description ordering** |
| `internal/mcp/tools_write.go` | `trace_add` kind dispatch |
| `internal/install/enforcement.go` | `ParentOrchestratorRule`, `cursorRulesMDCContent`, `claudeFallbackRulesContent`, `AgentsEnforcementBlock`, hook script |

## INT mapping (audit scope)

| INT | Phase 26 theme | Audit focus |
|-----|----------------|-------------|
| **INT-01** | P25-A | Where BLOCKING discoveries could become tasks; `loop apply` `spawned_tasks[]` path; `ImportSeedTask`; gap vs guided promotion |
| **INT-06** | P25-A | MCP `trace_add` description text vs CLI `trace add`; harness nudge in install rules (GapPassPrompt only?) |
| **INT-02** | P25-B | P19 saturation trigger (`NewPlanChanges==0 && NewSpawnedTasks==0`); first-empty-apply sticky STOP; threshold **options** (see below) |
| **INT-05** | P25-B | `DeliberationState.Stopped`; absence/presence of reset API clearing `Stopped` + `hop_count` + reopening EXECUTE |
| **INT-09** | P25-B | Gate `reason_code` vs status/export `stop_reason` when `Stopped=true` (`hop_budget_exceeded` vs persisted `p19_saturated`) |
| **Installer P25-2** | S04 | Exact insert site for `ParentOrchestratorRule` in MDC + Claude fallback (mirror `GapPassPrompt` pattern in `AgentsEnforcementBlock`) |

### Threshold options (document only — do not decide)

Record **current** behavior and list options for S01 planning; do **not** pick final numbers:

| Knob | Current (verify live) | Options to document |
|------|----------------------|---------------------|
| Hop budget | `deliberation.HopBudget` (= 12 in `types.go`) | Keep 12; raise for gap-pass profile; separate budget per phase profile |
| P19 saturation | True when last apply had zero plan_changes **and** zero spawned_tasks (`policy.go` `p19SaturatedFromLastStep`) | Exempt first apply on greenfield; require N applies; treat discoveries-only apply as non-saturating |
| Sticky STOP | `Stopped=true` never cleared in product paths (Phase 24 finding) | New reset API (INT-05); manual DB edit only; re-apply with plan_changes/spawn |
| Reason UX | `SelectNext` first branch: `Stopped` → `hop_budget_exceeded` even if `StopReason` was `p19_saturated` | Unify on persisted `StopReason`; remap in gate only; dual report with recovery hint |

## Role work

1. Read **Files to audit** table; open each file; record **file:line** evidence.
2. For each INT (+ installer gap), fill one section in `AUDIT.md` using the template below.
3. Cross-check Phase 24 `CODEBASE-AUDIT.md` — note what changed since Phase 24 (new files, moved lines, new tests).
4. Optional repro: run read-only CLI against G1 fixture (`trace -C experiments/ab-incident-tracker/runs/G1 loop status|gate …`) — capture stdout snippets in evidence appendix.
5. Write `AUDIT.md`; update own board row (`P26-S00-01`) to `done` with evidence pointer.

## Output: AUDIT.md structure

```markdown
# Loop audit — Phase 26

## Executive summary
(3–5 sentences: top gaps blocking P25-A/B + installer)

## INT-01 / INT-06: Discovery→task promotion

| file:line | current behavior | gap vs INT target | risk |
|-----------|------------------|---------------------|------|
| … | … | … | … |

## INT-02: Saturation recalibration

(same table + threshold options subsection — no final pick)

## INT-05: Deliberation reset

(same table; explicitly state whether reset API exists)

## INT-09: Unified STOP reason

(same table; gate vs export/status divergence)

## Installer gap (P25-2)

- `ParentOrchestratorRule` definition: `internal/install/enforcement.go:L20–35`
- Used in `AgentsEnforcementBlock`? (expected: no — only GapPassPrompt)
- Missing from `cursorRulesMDCContent` / `claudeFallbackRulesContent`? (expected: yes)
- Proposed insert pattern: mirror `GapPassPrompt` concat in `AgentsEnforcementBlock` (L65)

## Delta from Phase 24 CODEBASE-AUDIT
(bullet list of confirmed / changed / new findings)

## Evidence appendix
(commands run, optional G1 repro output)
```

Each INT section **must** include at least one row with concrete `file:line`, current behavior, gap, and risk.

## Exit criteria

- [ ] `AUDIT.md` present with all five INT sections + installer gap
- [ ] Each section has files/line ranges, behavior, gap, risk (table or equivalent)
- [ ] Threshold **options** documented; no final threshold numbers chosen
- [ ] P25-2 gap confirmed (`ParentOrchestratorRule` unused in MDC/Claude content)
- [ ] No product code or test changes in repo
- [ ] Board row `P26-S00-01` Notes cite `AUDIT.md` path

## Minimal todos

- [ ] Preflight: verify all paths in **Files to audit** exist
- [ ] Map INT-01/06 paths (spawn, add, MCP descriptions)
- [ ] Map INT-02/05/09 paths (SelectNext, saturation, deliberation, gate reason)
- [ ] Confirm P25-2 wiring gap in `enforcement.go`
- [ ] Document threshold options (no final pick)
- [ ] Write `AUDIT.md` with delta from Phase 24 baseline
- [ ] Update own board row to `done` with evidence

## Do not

- Change product code or write tests
- Choose final saturation/hop thresholds (defer to S01 `PLAN.md`)
- Start S01 or any later board row
- Mark P26-S00-00 (planner row) — that row is already closed by scope planner

## Todo updates

Status + notes on **P26-S00-01** only.
